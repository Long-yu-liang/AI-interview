<template>
  <div class="ai-avatar-fallback" :class="type">
    <span>AI</span>
  </div>
</template>

<script setup>
defineProps({
  type: {
    type: String,
    default: 'default'
  }
})
</script>

<style scoped>
.ai-avatar-fallback {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: #ffffff;
  border-radius: 50%;
}

.interview {
  background: linear-gradient(135deg, #0e7490, #0369a1);
}

.default, .super {
  background: linear-gradient(135deg, #0369a1, #075985);
}
</style> 
