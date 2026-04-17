package svc

import (
	"ai-gozero-agent/api/internal/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	openai "github.com/sashabaranov/go-openai"
	"log"
)

type ServiceContext struct {
	Config       config.Config
	OpenAIClient *openai.Client
	VectorStore  *VectorStore // 替换SessionStore
	PdfClient    *PdfClient
	Redis        *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 创建OpenAI客户端
	openaiConfig := openai.DefaultConfig(c.OpenAI.ApiKey)
	openaiConfig.BaseURL = c.OpenAI.BaseURL
	openAIClient := openai.NewClientWithConfig(openaiConfig)

	// 初始化Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	// 测试Redis连接
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	} else {
		log.Println("Redis连接成功")
	}

	// 初始化向量存储
	vectorStore, err := NewVectorStore(c.VectorDB, openAIClient)
	if err != nil {
		log.Fatalf("初始化向量数据库失败: %v", err)
	}

	// 测试数据库连接
	if err := vectorStore.TestConnection(); err != nil {
		log.Fatalf("向量数据库连接失败: %v", err)
	} else {
		log.Println("向量数据库连接成功")
	}

	return &ServiceContext{
		Config:       c,
		OpenAIClient: openAIClient,
		VectorStore:  vectorStore,
		PdfClient:    NewPdfClient(c.MCP.Endpoint),
		Redis:        rdb,
	}
}
