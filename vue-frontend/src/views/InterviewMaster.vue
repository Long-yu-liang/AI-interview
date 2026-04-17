<template>
  <div class="studio-page">
    <div class="ambient ambient-a"></div>
    <div class="ambient ambient-b"></div>

    <header class="studio-header">
      <button class="back-btn" @click="goBack">返回首页</button>
      <div class="title-wrap">
        <h1>AI 面试官</h1>
        <p>会话 ID：{{ chatId || "生成中" }}</p>
      </div>
      <div class="status-chip" :class="connectionStatus">
        {{ statusLabel }}
      </div>
    </header>

    <main class="studio-main">
      <ChatRoom
        :messages="messages"
        :connection-status="connectionStatus"
        ai-type="interview"
        @send-message="sendMessage"
      />
    </main>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, onBeforeUnmount } from "vue";
import { useRouter } from "vue-router";
import { useHead } from "@vueuse/head";
import ChatRoom from "../components/ChatRoom.vue";
import { chatWithLoveApp } from "../api";

useHead({
  title: "AI面试助手",
  meta: [
    {
      name: "description",
      content:
        "AI面试助手是超级智能体应用平台的专业面试助手，帮你解答各种面试问题，提供面试问题回答",
    },
  ],
});

const router = useRouter();
const messages = ref([]);
const chatId = ref("");
const connectionStatus = ref("disconnected");
let eventSource = null;

const statusLabel = computed(() => {
  if (connectionStatus.value === "connecting") return "思考中";
  if (connectionStatus.value === "error") return "连接异常";
  return "就绪";
});

const generateChatId = () => {
  return "interview_" + Math.random().toString(36).substring(2, 10);
};

const addMessage = (content, isUser) => {
  messages.value.push({
    content,
    isUser,
    time: new Date().getTime(),
  });
};

const sendMessage = (formData) => {
  const message = formData.get("message") || "";
  addMessage(message, true);

  if (eventSource) {
    eventSource.close();
  }

  addMessage("", false);
  const aiMessageIndex = messages.value.length - 1;

  formData.append("chatId", chatId.value);
  connectionStatus.value = "connecting";

  const handleMessage = (data) => {
    if (data && data !== "[DONE]") {
      if (aiMessageIndex < messages.value.length) {
        messages.value[aiMessageIndex].content += data;
      }
    }

    if (data === "[DONE]") {
      connectionStatus.value = "disconnected";
      eventSource.close();
    }
  };

  const handleError = () => {
    connectionStatus.value = "error";
    eventSource.close();
  };

  eventSource = chatWithLoveApp(formData, chatId.value);

  if (eventSource.onmessage !== undefined) {
    eventSource.onmessage = handleMessage;
    eventSource.onerror = handleError;
  } else {
    eventSource.onmessage = (event) => {
      handleMessage(event.data);
    };
    eventSource.onerror = handleError;
  }
};

const goBack = () => {
  router.push("/");
};

onMounted(() => {
  chatId.value = generateChatId();
  addMessage(
    "我是你的 AI 面试官，请上传你的简历，我来帮你解析。",
    false
  );
});

onBeforeUnmount(() => {
  if (eventSource) {
    eventSource.close();
  }
});
</script>

<style scoped>
.studio-page {
  position: relative;
  min-height: 100vh;
  background: radial-gradient(circle at 10% 18%, rgba(249, 115, 22, 0.09), transparent 34%),
    radial-gradient(circle at 90% 6%, rgba(14, 116, 144, 0.15), transparent 40%),
    linear-gradient(160deg, #f7f8fb 0%, #eef2f7 56%, #f8fafc 100%);
}

.ambient {
  position: absolute;
  z-index: 0;
  border-radius: 999px;
  filter: blur(48px);
  pointer-events: none;
}

.ambient-a {
  width: 18rem;
  height: 18rem;
  left: -4rem;
  top: 7rem;
  background: rgba(249, 115, 22, 0.14);
}

.ambient-b {
  width: 22rem;
  height: 22rem;
  right: -5rem;
  bottom: 18rem;
  background: rgba(14, 116, 144, 0.16);
}

.studio-header,
.studio-main {
  position: relative;
  z-index: 1;
  max-width: 1120px;
  margin: 0 auto;
  padding-inline: 1.1rem;
}

.studio-header {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 0.8rem;
  padding-top: 1rem;
}

.back-btn {
  border: 1px solid rgba(15, 23, 42, 0.12);
  background: rgba(255, 255, 255, 0.8);
  color: #0f172a;
  border-radius: 0.7rem;
  padding: 0.48rem 0.86rem;
  font-size: 0.86rem;
  cursor: pointer;
}

.title-wrap {
  text-align: center;
}

.title-wrap h1 {
  margin: 0;
  font-family: "Outfit", "Noto Sans SC", sans-serif;
  font-size: clamp(1.2rem, 3.1vw, 1.7rem);
}

.title-wrap p {
  margin: 0.2rem 0 0;
  color: var(--ink-600);
  font-size: 0.86rem;
}

.status-chip {
  justify-self: end;
  border-radius: 999px;
  padding: 0.28rem 0.66rem;
  font-size: 0.78rem;
  border: 1px solid rgba(15, 23, 42, 0.12);
  background: rgba(255, 255, 255, 0.78);
  color: #334155;
}

.status-chip.connecting {
  color: #0369a1;
  border-color: rgba(3, 105, 161, 0.3);
}

.status-chip.error {
  color: #b91c1c;
  border-color: rgba(185, 28, 28, 0.3);
}

.studio-main {
  margin-top: 0.85rem;
  height: calc(100vh - 92px);
  min-height: 600px;
}

@media (max-width: 768px) {
  .studio-header {
    grid-template-columns: auto 1fr;
  }

  .status-chip {
    display: none;
  }

  .title-wrap {
    text-align: left;
  }

  .studio-main {
    height: calc(100vh - 86px);
    min-height: 520px;
  }
}
</style>
