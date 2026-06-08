import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  server: {
    port: 3000,
    cors: true,
    host: '0.0.0.0',
    // 添加新域名到 allowedHosts 数组
    allowedHosts: ['794jx56302cy.vicp.fun'],
    proxy: {
      '/api': {
        target: 'http://localhost:8123/api',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
})
