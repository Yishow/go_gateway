import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // 載入環境變數（從 .env 文件）
  const env = loadEnv(mode, process.cwd(), '')
  
  // 從環境變數讀取配置，提供預設值
  const devPort = parseInt(env.VITE_DEV_PORT || '5173', 10)
  const proxyTarget = env.VITE_API_PROXY_TARGET || 'http://localhost:8080'
  
  return {
    plugins: [react()],
    base: '/',
    build: {
      outDir: 'dist',
      emptyOutDir: true,
    },
    server: {
      // 從環境變數讀取端口，預設 5173
      port: devPort,
      proxy: {
        '/api': {
          // 從環境變數讀取代理目標，預設 http://localhost:8080
          target: proxyTarget,
          changeOrigin: true,
          secure: false,
        },
      },
    },
  }
})
