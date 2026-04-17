<template>
  <div class="chat-container">
    <!-- 聊天记录区域 -->
    <div class="chat-messages" ref="messagesContainer">
      <div
        v-for="(msg, index) in messages"
        :key="index"
        class="message-wrapper"
      >
        <!-- AI消息 -->
        <div v-if="!msg.isUser" class="message ai-message" :class="[msg.type]">
          <div class="avatar ai-avatar">
            <AiAvatarFallback :type="aiType" />
          </div>
          <div class="message-bubble">
            <!-- 使用 XMarkdown 渲染 AI 消息 -->
            <XMarkdown
              :markdown="processMessageContent(msg.content || '')"
              :enableLatex="true"
              :enableBreaks="true"
              :allowHtml="false"
              class="message-content"
            />
            <div class="message-time">{{ formatTime(msg.time) }}</div>
          </div>
        </div>

        <!-- 用户消息 -->
        <div v-else class="message user-message" :class="[msg.type]">
          <div class="message-bubble">
            <div class="message-content">{{ msg.content }}</div>
            <div class="message-time">{{ formatTime(msg.time) }}</div>
          </div>
          <div class="avatar user-avatar">
            <div class="avatar-placeholder">我</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 使用 Sender 组件作为输入区域 -->
    <div class="chat-input-container">
      <Sender
        v-model="inputMessage"
        variant="updown"
        :disabled="connectionStatus === 'connecting'"
        :loading="connectionStatus === 'connecting'"
        submit-type="enter"
        :auto-size="{ minRows: 3, maxRows: 4 }"
        clearable
        placeholder="请输入消息..."
        class="chat-sender"
        @submit="handleSendMessage"
      >
        <template #prefix>
          <div class="sender-prefix">
            <el-upload
              v-model:file-list="file"
              :limit="1"
              class="upload-chip"
              :before-remove="beforeRemove"
              :auto-upload="false"
              :on-exceed="handleExceed"
              :show-file-list="false"
            >
              <el-icon class="upload-icon"><Paperclip /></el-icon>
              <el-text v-if="file[0]?.name" class="upload-name" size="small">{{ file[0]?.name }}</el-text>
            </el-upload>
            <el-upload
              v-model:file-list="knowledgeFile"
              :limit="1"
              accept=".pdf"
              class="upload-chip knowledge-chip"
              :before-remove="beforeRemoveKnowledge"
              :auto-upload="false"
              :on-exceed="handleExceedKnowledge"
              :show-file-list="false"
              :on-change="handleKnowledgeFileChange"
            >
              <el-icon class="upload-icon knowledge-icon"><Document /></el-icon>
              <el-text v-if="knowledgeFile[0]?.name" class="upload-name knowledge-name" size="small">{{ knowledgeFile[0]?.name }}</el-text>
              <el-text v-else class="upload-name knowledge-name" size="small">知识库PDF</el-text>
            </el-upload>
            <div :class="['thinking-chip', { isSelect }]" @click="isSelect = !isSelect">
              <el-icon><ElementPlus /></el-icon>
              <span>深度思考</span>
            </div>
          </div>
        </template>

        <!-- <template #action-list>
          <div style="display: flex; align-items: center; gap: 8px">
            <el-button round color="#626aef">
              <el-icon><Promotion /></el-icon>
            </el-button>
          </div>
        </template> -->
      </Sender>
    </div>
  </div>
</template>

<script setup>
import { ElementPlus, Paperclip, Document } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { ref } from "vue";
import AiAvatarFallback from "./AiAvatarFallback.vue";
import { uploadKnowledge } from "../api/index.js";

const props = defineProps({
  messages: {
    type: Array,
    default: () => [],
  },
  connectionStatus: {
    type: String,
    default: "disconnected",
  },
  aiType: {
    type: String,
    default: "default", // 'interview' 或 'super'
  },
});

const file = ref([]);
const knowledgeFile = ref([]);
const isSelect = ref(false);
const emit = defineEmits(["send-message", "upload-knowledge"]);

const inputMessage = ref("");
const messagesContainer = ref(null);

// 处理消息内容中的转义符
const processMessageContent = (content) => {
  if (!content) return "";

  // 处理各种转义符，按照正确的顺序处理避免冲突
  let processedContent = content;

  // 首先处理双重转义（避免 \\n 被错误处理）
  processedContent = processedContent.replace(/\\\\\\/g, "\\TEMP_BACKSLASH\\");

  // 处理常见的转义符
  processedContent = processedContent
    .replace(/\\n/g, "\n") // 将 \n 转换为真正的换行符
    .replace(/\\t/g, "\t") // 将 \t 转换为制表符
    .replace(/\\r/g, "\r") // 将 \r 转换为回车符
    .replace(/\\"/g, '"') // 将 \" 转换为双引号
    .replace(/\\'/g, "'") // 将 \' 转换为单引号
    .replace(/\\\//g, "/") // 将 \/ 转换为斜杠
    .replace(/\\b/g, "\b") // 将 \b 转换为退格符
    .replace(/\\f/g, "\f") // 将 \f 转换为换页符
    .replace(/\\v/g, "\v") // 将 \v 转换为垂直制表符
    .replace(/\\\\/g, "\\") // 将 \\ 转换为反斜杠
    .replace(/\\TEMP_BACKSLASH\\/g, "\\"); // 恢复临时标记的反斜杠

  // 处理可能的 Unicode 转义序列
  processedContent = processedContent.replace(
    /\\u([0-9a-fA-F]{4})/g,
    (match, hex) => {
      return String.fromCharCode(parseInt(hex, 16));
    }
  );

  // 处理特殊的 markdown 字符转义
  processedContent = processedContent
    .replace(/\\`/g, "`") // 将 \` 转换为反引号
    .replace(/\\~/g, "~") // 将 \~ 转换为波浪号
    .replace(/\\#/g, "#") // 将 \# 转换为井号
    .replace(/\\>/g, ">") // 将 \> 转换为大于号
    .replace(/\\</g, "<") // 将 \< 转换为小于号
    .replace(/\\&/g, "&") // 将 \& 转换为和号
    .replace(/\\\*/g, "*") // 将 \* 转换为星号
    .replace(/\\-/g, "-") // 将 \- 转换为横线
    .replace(/\\\+/g, "+") // 将 \+ 转换为加号
    .replace(/\\=/g, "=") // 将 \= 转换为等号
    .replace(/\\\|/g, "|") // 将 \| 转换为竖线
    .replace(/\\\[/g, "[") // 将 \[ 转换为左方括号
    .replace(/\\\]/g, "]") // 将 \] 转换为右方括号
    .replace(/\\\{/g, "{") // 将 \{ 转换为左大括号
    .replace(/\\\}/g, "}") // 将 \} 转换为右大括号
    .replace(/\\\(/g, "(") // 将 \( 转换为左小括号
    .replace(/\\\)/g, ")"); // 将 \) 转换为右小括号

  // 处理列表格式，确保列表项前后有适当的换行
  processedContent = processedContent.replace(/^(\s*)([\-\*\+])\s+/gm, "$1$2 ");
  processedContent = processedContent.replace(/^(\s*)(\d+\.)\s+/gm, "$1$2 ");

  // 确保列表项之间有适当的空行处理
  processedContent = processedContent.replace(
    /(\n\s*[\-\*\+\d\.]\s+)/g,
    "\n$1"
  );

  // 去除前后空格并返回
  return processedContent.trim();
};

// 处理发送消息
const handleSendMessage = (message) => {
  if (!message.trim()) return;

  // 创建 FormData 格式
  const formData = new FormData();
  formData.append("message", message);

  // 只有当文件存在时才添加文件
  if (file.value && file.value.length > 0 && file.value[0]) {
    formData.append("file", file.value[0].raw);
  }

  emit("send-message", formData);
  file.value = [];
  inputMessage.value = "";
};

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp);
  return date.toLocaleTimeString("zh-CN", {
    hour: "2-digit",
    minute: "2-digit",
  });
};

// 处理文件上传相关函数
const beforeRemove = (uploadFile, uploadFiles) => {
  return true; // 允许删除
};

const handleExceed = (files, uploadFiles) => {
  file.value = [files[0]];
};

// 处理知识库PDF上传相关函数
const beforeRemoveKnowledge = (uploadFile, uploadFiles) => {
  return true; // 允许删除
};

const handleExceedKnowledge = (files, uploadFiles) => {
  knowledgeFile.value = [files[0]];
};

const handleKnowledgeFileChange = async (uploadFile, uploadFiles) => {
  if (uploadFile.status === "ready" && uploadFile.raw) {
    // 验证文件类型
    if (uploadFile.raw.type !== "application/pdf") {
      ElMessage.error("只支持PDF格式的文件");
      knowledgeFile.value = [];
      return;
    }

    // 验证文件大小（限制为500KB）
    if (uploadFile.raw.size > 500 * 1024) {
      ElMessage.error("文件大小不能超过500KB");
      knowledgeFile.value = [];
      return;
    }

    try {
      // 显示上传中提示
      const loadingMessage = ElMessage({
        message: "正在上传知识库文件...",
        type: "info",
        duration: 0, // 不自动关闭
      });

      // 创建FormData并上传
      const formData = new FormData();
      formData.append("file", uploadFile.raw);

      // 调用上传接口
      await uploadKnowledge(formData);

      // 关闭加载提示
      loadingMessage.close();

      // 显示上传成功提示
      ElMessage.success("知识库文件上传成功");
    } catch (error) {
      // 关闭加载提示（如果存在）
      ElMessage.closeAll();

      // 显示错误提示
      ElMessage.error(
        error.response?.data?.message || "知识库文件上传失败，请重试"
      );

      // 清空文件列表
      knowledgeFile.value = [];
    }
  }
};
</script>

<style scoped>
.chat-container {
  position: relative;
  height: 100%;
  border-radius: 1rem;
  overflow: hidden;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: rgba(255, 255, 255, 0.76);
  box-shadow: 0 20px 48px rgba(15, 23, 42, 0.12);
  backdrop-filter: blur(10px);
}

.chat-messages {
  position: absolute;
  inset: 0 0 112px;
  overflow-y: auto;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.55rem;
}

.message-wrapper {
  display: flex;
  width: 100%;
}

.message {
  display: flex;
  align-items: flex-start;
  max-width: min(86%, 760px);
}

.user-message {
  margin-left: auto;
}

.ai-message {
  margin-right: auto;
}

.avatar {
  width: 2rem;
  height: 2rem;
  border-radius: 999px;
  overflow: hidden;
  flex-shrink: 0;
  display: grid;
  place-items: center;
}

.user-avatar {
  margin-left: 0.5rem;
}

.ai-avatar {
  margin-right: 0.5rem;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: grid;
  place-items: center;
  color: #fff;
  font-weight: 700;
  font-size: 0.75rem;
  background: linear-gradient(135deg, #0e7490, #0369a1);
}

.message-bubble {
  border-radius: 1rem;
  padding: 0.7rem 0.82rem;
  line-height: 1.6;
}

.user-message .message-bubble {
  background: linear-gradient(135deg, #0e7490, #0369a1);
  color: #fff;
  border-bottom-right-radius: 0.36rem;
}

.ai-message .message-bubble {
  background: rgba(255, 255, 255, 0.92);
  color: #0f172a;
  border: 1px solid rgba(15, 23, 42, 0.08);
  border-bottom-left-radius: 0.36rem;
}

.message-content {
  font-size: 0.95rem;
}

.message-time {
  margin-top: 0.35rem;
  font-size: 0.72rem;
  text-align: right;
  opacity: 0.74;
}

.chat-input-container {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 0.78rem;
  border-top: 1px solid rgba(15, 23, 42, 0.08);
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(9px);
}

.chat-sender {
  width: 100%;
}

.sender-prefix {
  display: flex;
  align-items: center;
  gap: 0.48rem;
  flex-wrap: wrap;
}

.upload-chip {
  height: 1.7rem;
  width: auto;
  padding: 0 0.52rem;
  border: 1px solid rgba(15, 23, 42, 0.14);
  border-radius: 999px;
  background: rgba(248, 250, 252, 0.9);
}

.knowledge-chip {
  border-color: rgba(249, 115, 22, 0.34);
  background: rgba(255, 247, 237, 0.95);
}

.upload-icon {
  height: 1.35rem;
  color: #334155;
}

.knowledge-icon,
.knowledge-name {
  color: #c2410c;
}

.upload-name {
  margin-left: 0.25rem;
}

.thinking-chip {
  display: flex;
  align-items: center;
  gap: 0.28rem;
  padding: 0.1rem 0.72rem;
  border: 1px solid rgba(15, 23, 42, 0.22);
  border-radius: 999px;
  cursor: pointer;
  font-size: 0.76rem;
  background: rgba(255, 255, 255, 0.9);
  transition: all 160ms ease;
}

.thinking-chip.isSelect {
  color: #0369a1;
  border-color: rgba(3, 105, 161, 0.45);
  background: rgba(224, 242, 254, 0.92);
  font-weight: 700;
}

.chat-sender :deep(.el-textarea__inner) {
  border-radius: 0.9rem;
  border: 1px solid rgba(15, 23, 42, 0.14);
  padding: 0.72rem 0.88rem;
  font-size: 0.96rem;
  resize: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.chat-sender :deep(.el-textarea__inner:focus) {
  border-color: #0369a1;
  box-shadow: 0 0 0 3px rgba(3, 105, 161, 0.12);
}

.chat-sender :deep(.el-button--primary) {
  border-radius: 999px;
  border: none;
  background: linear-gradient(135deg, #0e7490, #0369a1);
}

.ai-answer,
.ai-final {
  animation: riseIn 0.28s ease;
}

@keyframes riseIn {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.ai-message .message-content :deep(div),
.ai-message .message-content :deep(p),
.ai-message .message-content :deep(li) {
  color: #1e293b;
}

.ai-message .message-content :deep(pre) {
  margin: 0.55rem 0;
  padding: 0.65rem;
  border-radius: 0.65rem;
  background: rgba(15, 23, 42, 0.06);
  overflow-x: auto;
}

.ai-message .message-content :deep(code) {
  font-family: "JetBrains Mono", monospace;
  font-size: 0.84em;
}

@media (max-width: 768px) {
  .chat-messages {
    inset: 0 0 106px;
    padding: 0.78rem;
  }

  .message {
    max-width: 95%;
  }

  .chat-input-container {
    padding: 0.62rem;
  }
}
</style>
