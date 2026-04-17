<template>
  <div class="studio-page">
    <div class="ambient ambient-a"></div>
    <div class="ambient ambient-b"></div>

    <header class="studio-header">
      <button class="back-btn" @click="goBack">返回首页</button>
      <div class="title-wrap">
        <h1>AI 超级智能体</h1>
        <p>开放问答 · 深度思考</p>
      </div>
      <div class="status-chip" :class="connectionStatus">
        {{ statusLabel }}
      </div>
    </header>

    <main class="studio-main">
      <ChatRoom
        :messages="messages"
        :connection-status="connectionStatus"
        ai-type="super"
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
import { chatWithManus } from "../api";

useHead({
  title: "AI超级智能体 - 大鱼AI超级智能体应用平台",
  meta: [
    {
      name: "description",
      content:
        "AI超级智能体是大鱼AI超级智能体应用平台的全能助手，能解答各类专业问题，提供精准建议和解决方案",
    },
  ],
});

const router = useRouter();
const messages = ref([]);
const connectionStatus = ref("disconnected");
let eventSource = null;

const statusLabel = computed(() => {
  if (connectionStatus.value === "connecting") return "思考中";
  if (connectionStatus.value === "error") return "连接异常";
  return "就绪";
});

const addMessage = (content, isUser, type = "") => {
  messages.value.push({
    content,
    isUser,
    type,
    time: new Date().getTime(),
  });
};

const sendMessage = (message) => {
  addMessage(message, true, "user-question");

  if (eventSource) {
    eventSource.close();
  }

  connectionStatus.value = "connecting";

  let messageBuffer = [];
  let lastBubbleTime = Date.now();
  let isFirstResponse = true;

  const chineseEndPunctuation = ["。", "！", "？", "…"];
  const minBubbleInterval = 800;

  const createBubble = (content, type = "ai-answer") => {
    if (!content.trim()) return;

    const now = Date.now();
    const timeSinceLastBubble = now - lastBubbleTime;

    if (isFirstResponse) {
      addMessage(content, false, type);
      isFirstResponse = false;
    } else if (timeSinceLastBubble < minBubbleInterval) {
      setTimeout(() => {
        addMessage(content, false, type);
      }, minBubbleInterval - timeSinceLastBubble);
    } else {
      addMessage(content, false, type);
    }

    lastBubbleTime = now;
    messageBuffer = [];
  };

  eventSource = chatWithManus(message);

  eventSource.onmessage = (event) => {
    const data = event.data;

    if (data && data !== "[DONE]") {
      messageBuffer.push(data);
      const combinedText = messageBuffer.join("");

      const lastChar = data.charAt(data.length - 1);
      const hasCompleteSentence =
        chineseEndPunctuation.includes(lastChar) || data.includes("\n\n");
      const isLongEnough = combinedText.length > 40;

      if (hasCompleteSentence || isLongEnough) {
        createBubble(combinedText);
      }
    }

    if (data === "[DONE]") {
      if (messageBuffer.length > 0) {
        createBubble(messageBuffer.join(""), "ai-final");
      }
      connectionStatus.value = "disconnected";
      eventSource.close();
    }
  };

  eventSource.onerror = () => {
    connectionStatus.value = "error";
    eventSource.close();

    if (messageBuffer.length > 0) {
      createBubble(messageBuffer.join(""), "ai-error");
    }
  };
};

const goBack = () => {
  router.push("/");
};

onMounted(() => {
  addMessage(
    "你好，我是 AI 超级智能体。你可以直接问我技术问题、项目设计或面试准备建议。",
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
