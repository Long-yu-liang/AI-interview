package logic

import (
	"context"
	"errors"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"strings"

	"ai-gozero-agent/api/internal/svc"
	"ai-gozero-agent/api/internal/types"
	"ai-gozero-agent/api/internal/utils"
)

type ChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatLogic) Chat(req *types.InterviewAPPChatReq) (<-chan *types.ChatResponse, error) {
	ch := make(chan *types.ChatResponse)

	go func() {
		defer close(ch)
		stateManager := NewStateManager(l.svcCtx)

		// 1. 保存用户消息
		if err := l.svcCtx.VectorStore.SaveMessage(
			req.ChatId,
			openai.ChatMessageRoleUser,
			req.Message,
		); err != nil {
			l.Logger.Errorf("保存用户消息失败: %v", err)
		}

		// 2. 获取当前状态（确保初始化）
		currentState, err := stateManager.GetOrInitState(req.ChatId)
		if err != nil {
			l.Logger.Errorf("获取状态失败: %v", err)
			currentState = types.StateStart
		}

		// 3. 知识检索
		knowledge, err := l.svcCtx.VectorStore.RetrieveKnowledge(req.Message, 3)
		if err != nil {
			l.Logger.Errorf("知识检索失败: %v", err)
			knowledge = []types.KnowledgeChunk{}
		}

		// 4. 构建系统消息（带状态）
		messages, err := l.buildMessagesWithState(req.ChatId, currentState, knowledge)
		if err != nil {
			l.Logger.Errorf("构建消息失败: %v", err)
			ch <- &types.ChatResponse{Content: "系统错误：无法构建对话", IsLast: true}
			return
		}

		// 5. 创建OpenAI请求
		request := openai.ChatCompletionRequest{
			Model:       l.svcCtx.Config.OpenAI.Model,
			Messages:    messages,
			Stream:      true,
			MaxTokens:   l.svcCtx.Config.OpenAI.MaxTokens,
			Temperature: l.svcCtx.Config.OpenAI.Temperature,
		}

		// 6. 处理流式响应
		var fullResponse strings.Builder
		stream, err := l.svcCtx.OpenAIClient.CreateChatCompletionStream(l.ctx, request)
		if err != nil {
			l.Logger.Error("创建聊天完成流失败: ", err)
			ch <- &types.ChatResponse{Content: "系统错误：无法连接AI服务", IsLast: true}
			return
		}
		defer stream.Close()

		for {
			select {
			case <-l.ctx.Done():
				return
			default:
				response, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					// 流结束后处理状态更新
					finalResponse := fullResponse.String()
					if finalResponse != "" {
						// 保存AI回复
						if err := l.svcCtx.VectorStore.SaveMessage(
							req.ChatId,
							openai.ChatMessageRoleAssistant,
							finalResponse,
						); err != nil {
							l.Logger.Errorf("保存助手消息失败: %v", err)
						}

						// 更新状态
						newState, err := stateManager.EvaluateAndUpdateState(req.ChatId, finalResponse)
						if err != nil {
							l.Logger.Errorf("更新状态失败: %v", err)
						} else {
							l.Logger.Infof("状态更新: %s -> %s", currentState, newState)
						}
					}

					ch <- &types.ChatResponse{IsLast: true}
					return
				}
				if err != nil {
					l.Logger.Error("接收流数据失败: ", err)
					return
				}

				if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
					content := response.Choices[0].Delta.Content
					fullResponse.WriteString(content)
					ch <- &types.ChatResponse{
						Content: content,
						IsLast:  false,
					}
				}
			}
		}
	}()

	return ch, nil
}

// 构建带状态的消息
func (l *ChatLogic) buildMessagesWithState(chatId, currentState string, knowledge []types.KnowledgeChunk) ([]openai.ChatCompletionMessage, error) {
	// 构建状态特定的系统消息
	systemMessage := "你是一个专业的Go语言面试官，负责评估候选人的Go语言能力。"
	systemMessage += "\n\n当前状态: " + currentState

	switch currentState {
	case types.StateStart:
		systemMessage += "\n目标: 欢迎候选人并开始面试流程"
	case types.StateQuestion:
		systemMessage += "\n目标: 提出有深度的问题考察Go语言核心概念"
	case types.StateFollowUp:
		systemMessage += "\n目标: 基于候选人的回答进行追问，深入考察理解深度"
	case types.StateEvaluate:
		systemMessage += "\n目标: 全面评估候选人的技术能力"
	case types.StateEnd:
		systemMessage += "\n目标: 结束面试并提供反馈"
	}

	// 注入知识库
	if len(knowledge) > 0 {
		systemMessage += "\n\n相关背景知识："
		for i, k := range knowledge {
			truncatedContent := utils.TruncateText(k.Content, 500)
			systemMessage += fmt.Sprintf("\n[知识片段%d] %s: %s", i+1, k.Title, truncatedContent)
		}
	}

	// 获取历史消息
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemMessage,
		},
	}

	history, err := l.svcCtx.VectorStore.GetMessages(chatId, 10)
	if err != nil {
		return nil, err
	}

	for _, msg := range history {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return messages, nil
}
