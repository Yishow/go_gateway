import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // 載入環境變數（從 .env 文件）
  const env = loadEnv(mode, process.cwd(), "");

  // 從環境變數讀取配置，提供預設值
  const devPort = parseInt(env.VITE_DEV_PORT || "5173", 10);
  const proxyTarget = env.VITE_API_PROXY_TARGET || "http://localhost:8080";

  return {
    plugins: [react()],
    base: "/",
    /**
     * 強制所有依賴使用同一個 React 實例
     * 解決 react-i18next "Cannot read properties of null (reading 'useMemo')" 錯誤
     * 此錯誤發生於依賴包使用不同的 React 實例時
     */
    resolve: {
      dedupe: ["react", "react-dom"],
      alias: {
        react: "react",
        "react-dom": "react-dom",
        "@": path.resolve(__dirname, "./src"),
      },
    },
    build: {
      outDir: "dist",
      emptyOutDir: true,
    },
    server: {
      fs: {
        allow: [".."],
      },
      // 從環境變數讀取端口，預設 5173
      port: devPort,
      proxy: {
        "/api": {
          // 從環境變數讀取代理目標，預設 http://localhost:8080
          target: proxyTarget,
          changeOrigin: true,
          secure: false,
        },
      },
    },
    // Force CSS rebuild
    css: {
      devSourcemap: true,
    },
    test: {
      globals: true,
      environment: "jsdom",
      setupFiles: "./src/setupTests.ts",
    },
  };
});
