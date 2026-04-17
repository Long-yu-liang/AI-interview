<template>
  <div class="home-page">
    <div class="ambient ambient-one"></div>
    <div class="ambient ambient-two"></div>

    <header class="topbar">
      <div class="brand">
        <div class="brand-mark">AI</div>
        <div class="brand-text">
          <strong>面试官工作台</strong>
          <span>Interview Studio</span>
        </div>
      </div>
      <router-link to="/interview-master" class="cta ghost">立即体验</router-link>
    </header>

    <main class="hero">
      <section class="hero-copy">
        <p class="eyebrow">AI 模拟面试平台</p>
        <h1>把每一次面试，变成可复盘的训练。</h1>
        <p class="summary">
          覆盖算法、系统设计与表达能力，实时对话反馈，帮助你在真正的技术面试中更稳、更清晰。
        </p>

        <div class="hero-actions">
          <router-link to="/interview-master" class="cta primary">开始模拟面试</router-link>
          <span class="cta link" role="text">查看学习资料</span>
        </div>

        <ul class="metrics">
          <li>
            <strong>10k+</strong>
            <span>训练用户</span>
          </li>
          <li>
            <strong>95%</strong>
            <span>反馈满意度</span>
          </li>
          <li>
            <strong>24/7</strong>
            <span>随时练习</span>
          </li>
        </ul>
      </section>

      <section class="hero-panel">
        <article class="interview-card">
          <div class="card-head">
            <span>实时模拟</span>
            <span class="online">online</span>
          </div>
          <div class="question">
            设计一个高并发下的秒杀库存扣减方案，并说明如何保证最终一致性。
          </div>
          <div class="feedback">
            <h3>AI 反馈</h3>
            <p>思路完整，建议补充幂等设计、重试机制与监控告警闭环。</p>
          </div>
          <div class="tags">
            <span>系统设计</span>
            <span>分布式</span>
            <span>高可用</span>
          </div>
        </article>
      </section>
    </main>

    <section class="features">
      <article class="feature">
        <h3>智能追问</h3>
        <p>根据你的回答动态调整难度，不是固定题库，而是更接近真实面试官的连续追问。</p>
      </article>
      <article class="feature">
        <h3>结构化复盘</h3>
        <p>从表达清晰度、技术深度、方案权衡三个维度打分，知道自己卡在什么位置。</p>
      </article>
      <article class="feature">
        <h3>知识库注入</h3>
        <p>支持上传 PDF 简历或项目材料，让面试问题围绕你真实经历展开，更有针对性。</p>
      </article>
    </section>

    <AppFooter />
  </div>
</template>

<script setup>
import AppFooter from "../components/AppFooter.vue";
</script>

<style scoped>
.home-page {
  position: relative;
  min-height: 100vh;
  background: radial-gradient(circle at 15% 20%, rgba(249, 115, 22, 0.08), transparent 35%),
    radial-gradient(circle at 85% 0%, rgba(14, 116, 144, 0.14), transparent 42%),
    linear-gradient(160deg, #f7f8fb 0%, #eef2f7 55%, #f7fafc 100%);
  color: var(--ink-900);
  overflow: hidden;
}

.ambient {
  position: absolute;
  z-index: 0;
  border-radius: 999px;
  filter: blur(42px);
  pointer-events: none;
}

.ambient-one {
  width: 22rem;
  height: 22rem;
  right: -4rem;
  top: 4rem;
  background: rgba(14, 116, 144, 0.16);
}

.ambient-two {
  width: 18rem;
  height: 18rem;
  left: -5rem;
  bottom: 24rem;
  background: rgba(249, 115, 22, 0.14);
}

.topbar,
.hero,
.features {
  position: relative;
  z-index: 1;
  max-width: 1120px;
  margin: 0 auto;
  padding-inline: 1.25rem;
}

.topbar {
  padding-top: 1.2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.brand {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.brand-mark {
  width: 2.3rem;
  height: 2.3rem;
  display: grid;
  place-items: center;
  border-radius: 0.7rem;
  font-weight: 700;
  color: #fff;
  background: linear-gradient(135deg, #0e7490, #0369a1);
  box-shadow: 0 8px 24px rgba(14, 116, 144, 0.25);
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.brand-text strong {
  font-size: 0.96rem;
  line-height: 1.1;
}

.brand-text span {
  color: var(--ink-600);
  font-size: 0.78rem;
}

.hero {
  padding-top: 3.5rem;
  display: grid;
  grid-template-columns: 1.05fr 0.95fr;
  gap: 2.4rem;
  align-items: center;
}

.eyebrow {
  display: inline-flex;
  margin: 0 0 1rem;
  border-radius: 999px;
  background: rgba(14, 116, 144, 0.12);
  color: #0e7490;
  padding: 0.3rem 0.75rem;
  font-size: 0.85rem;
  font-weight: 600;
}

h1 {
  margin: 0;
  font-family: "Outfit", "Noto Sans SC", sans-serif;
  font-size: clamp(1.9rem, 4.6vw, 3.35rem);
  line-height: 1.12;
  letter-spacing: -0.03em;
}

.summary {
  margin-top: 1.1rem;
  max-width: 38rem;
  color: var(--ink-700);
  line-height: 1.72;
}

.hero-actions {
  margin-top: 1.5rem;
  display: flex;
  gap: 0.8rem;
  flex-wrap: wrap;
}

.cta {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  text-decoration: none;
  border-radius: 0.8rem;
  font-weight: 600;
  transition: transform 180ms ease, box-shadow 180ms ease, background-color 180ms ease;
}

.cta.primary {
  background: linear-gradient(135deg, #0e7490, #0369a1);
  color: #fff;
  padding: 0.75rem 1.25rem;
  box-shadow: 0 12px 26px rgba(14, 116, 144, 0.24);
}

.cta.primary:hover,
.cta.ghost:hover {
  transform: translateY(-2px);
}

.cta.ghost {
  padding: 0.6rem 1rem;
  color: #0f172a;
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(15, 23, 42, 0.08);
}

.cta.link {
  color: #0e7490;
  padding: 0.75rem 1rem;
  background: rgba(255, 255, 255, 0.66);
  border: 1px solid rgba(14, 116, 144, 0.18);
}

.metrics {
  list-style: none;
  margin: 1.5rem 0 0;
  padding: 0;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.75rem;
  max-width: 34rem;
}

.metrics li {
  background: rgba(255, 255, 255, 0.74);
  border: 1px solid rgba(15, 23, 42, 0.08);
  border-radius: 0.8rem;
  padding: 0.75rem 0.8rem;
}

.metrics strong {
  display: block;
  font-size: 1.2rem;
  color: #0f172a;
}

.metrics span {
  color: var(--ink-600);
  font-size: 0.83rem;
}

.interview-card {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(250, 252, 255, 0.96));
  border: 1px solid rgba(15, 23, 42, 0.08);
  border-radius: 1.1rem;
  padding: 1.2rem;
  box-shadow: 0 24px 52px rgba(15, 23, 42, 0.12);
  animation: riseIn 560ms ease;
}

.card-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.85rem;
  color: var(--ink-600);
}

.online {
  color: #15803d;
  font-weight: 600;
}

.question {
  margin-top: 0.9rem;
  background: rgba(14, 116, 144, 0.08);
  border: 1px solid rgba(14, 116, 144, 0.2);
  border-radius: 0.75rem;
  padding: 0.9rem;
  line-height: 1.6;
}

.feedback {
  margin-top: 0.9rem;
  padding: 0.9rem;
  border-radius: 0.75rem;
  background: rgba(249, 115, 22, 0.08);
  border: 1px solid rgba(249, 115, 22, 0.2);
}

.feedback h3 {
  margin: 0 0 0.4rem;
  font-size: 0.95rem;
}

.feedback p {
  margin: 0;
  line-height: 1.6;
  color: var(--ink-700);
}

.tags {
  margin-top: 0.9rem;
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.tags span {
  font-size: 0.77rem;
  color: #334155;
  padding: 0.22rem 0.58rem;
  border-radius: 999px;
  border: 1px solid rgba(15, 23, 42, 0.12);
  background: rgba(255, 255, 255, 0.9);
}

.features {
  margin-top: 3rem;
  margin-bottom: 3.5rem;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.9rem;
}

.feature {
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(15, 23, 42, 0.08);
  border-radius: 0.95rem;
  padding: 1rem;
  transition: transform 180ms ease;
}

.feature:hover {
  transform: translateY(-3px);
}

.feature h3 {
  margin: 0 0 0.45rem;
  font-size: 1rem;
}

.feature p {
  margin: 0;
  color: var(--ink-700);
  line-height: 1.62;
  font-size: 0.92rem;
}

@keyframes riseIn {
  from {
    opacity: 0;
    transform: translateY(14px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 980px) {
  .hero {
    grid-template-columns: 1fr;
    padding-top: 2.2rem;
  }

  .features {
    grid-template-columns: 1fr;
    margin-top: 2.1rem;
  }
}

@media (max-width: 640px) {
  .metrics {
    grid-template-columns: 1fr;
  }

  .topbar {
    padding-top: 1rem;
  }
}
</style>
