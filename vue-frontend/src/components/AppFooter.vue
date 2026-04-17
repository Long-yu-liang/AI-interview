<template>
  <footer class="app-footer">
    <div class="footer-shell">
      <section class="brand-col">
        <div class="brand-head">
          <div class="brand-mark">AI</div>
          <div>
            <h3>面试官工作台</h3>
            <p>Interview Studio</p>
          </div>
        </div>
        <p class="brand-desc">
          专注程序员面试模拟训练，帮助你在算法、系统设计与表达能力上持续进阶。
        </p>
      </section>

      <section class="link-col">
        <h4>训练内容</h4>
        <a href="#">算法面试</a>
        <a href="#">系统设计</a>
        <a href="#">项目深挖</a>
      </section>

      <section class="link-col">
        <h4>学习资源</h4>
        <span class="footer-text">LeetCode</span>
        <span class="footer-text">GitHub</span>
        <span class="footer-text">Stack Overflow</span>
      </section>

      <section class="link-col">
        <h4>联系与支持</h4>
        <span class="footer-text">商务合作</span>
        <span class="footer-text">社区交流</span>
        <span class="footer-text">用户协议</span>
      </section>
    </div>

    <div class="footer-bottom">
      <span>© {{ currentYear }} AI 面试官 · 让每次面试都成为成长机会</span>
      <span class="badge">Built for Developers</span>
    </div>
  </footer>
</template>

<script setup>
import { computed } from "vue";

const currentYear = computed(() => new Date().getFullYear());
</script>

<style scoped>
.app-footer {
  margin-top: auto;
  border-top: 1px solid rgba(15, 23, 42, 0.08);
  background: linear-gradient(160deg, rgba(255, 255, 255, 0.9), rgba(248, 250, 252, 0.9));
  backdrop-filter: blur(8px);
}

.footer-shell {
  max-width: 1120px;
  margin: 0 auto;
  padding: 2.1rem 1.2rem 1.3rem;
  display: grid;
  grid-template-columns: 1.3fr 1fr 1fr 1fr;
  gap: 1rem;
}

.brand-head {
  display: flex;
  align-items: center;
  gap: 0.65rem;
}

.brand-mark {
  width: 2.1rem;
  height: 2.1rem;
  border-radius: 0.62rem;
  display: grid;
  place-items: center;
  color: #fff;
  font-weight: 700;
  font-size: 0.8rem;
  letter-spacing: 0.08em;
  background: linear-gradient(135deg, #0e7490, #0369a1);
}

.brand-head h3 {
  margin: 0;
  font-size: 1rem;
  font-family: "Outfit", "Noto Sans SC", sans-serif;
}

.brand-head p {
  margin: 0.1rem 0 0;
  font-size: 0.76rem;
  color: #64748b;
}

.brand-desc {
  margin: 0.8rem 0 0;
  color: #475569;
  line-height: 1.7;
  font-size: 0.9rem;
}

.link-col h4 {
  margin: 0;
  font-size: 0.9rem;
  color: #0f172a;
}

.link-col {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.link-col a {
  color: #475569;
  font-size: 0.88rem;
  transition: color 160ms ease, transform 160ms ease;
}

.footer-text {
  color: #475569;
  font-size: 0.88rem;
}

.link-col a:hover {
  color: #0369a1;
  transform: translateX(2px);
}

.footer-bottom {
  max-width: 1120px;
  margin: 0 auto;
  padding: 0.85rem 1.2rem 1rem;
  border-top: 1px solid rgba(15, 23, 42, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #64748b;
  font-size: 0.8rem;
}

.badge {
  padding: 0.22rem 0.55rem;
  border-radius: 999px;
  border: 1px solid rgba(3, 105, 161, 0.25);
  color: #0369a1;
  background: rgba(224, 242, 254, 0.75);
  font-weight: 600;
}

@media (max-width: 980px) {
  .footer-shell {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 640px) {
  .footer-shell {
    grid-template-columns: 1fr;
    gap: 1.1rem;
    padding-top: 1.5rem;
  }

  .footer-bottom {
    flex-direction: column;
    gap: 0.45rem;
    text-align: center;
  }
}
</style>
