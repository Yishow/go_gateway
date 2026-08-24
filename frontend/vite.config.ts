import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

function resolveProxyTarget(env: Record<string, string>): string {
  if (env.VITE_API_PROXY_TARGET) {
    return env.VITE_API_PROXY_TARGET;
  }

  const backendPort = env.PORT || "8080";
  return `http://127.0.0.1:${backendPort}`;
}

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // 載入環境變數（從 .env 文件）
  const env = loadEnv(mode, process.cwd(), "");

  // 從環境變數讀取配置，提供預設值
  const devPort = parseInt(env.VITE_DEV_PORT || "5173", 10);
  const proxyTarget = resolveProxyTarget(env);

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
      rollupOptions: {
        output: {
          /**
           * 將大型第三方依賴拆分為獨立 chunk：
           * - 避免單一 chunk 超過 500 kB 觸發 Rollup 警告
           * - 各 vendor hash 獨立，嵌入式部署時可善用瀏覽器快取
           */
          manualChunks(id: string) {
            if (!id.includes("node_modules")) {
              return undefined;
            }
            if (/[\\/]node_modules[\\/](react|react-dom|scheduler)[\\/]/.test(id)) {
              return "react-vendor";
            }
            if (/[\\/]node_modules[\\/](@mui|@emotion)[\\/]/.test(id)) {
              return "mui-vendor";
            }
            if (
              /[\\/]node_modules[\\/](recharts|victory-vendor|react-smooth|recharts-scale|d3-[^\\/]+|internmap)[\\/]/.test(
                id,
              )
            ) {
              return "charts-vendor";
            }
            if (
              /[\\/]node_modules[\\/](?:@monaco-editor|monaco-editor|monaco-editor-core|mermaid|@mermaid-js|reactflow|@xyflow|elkjs|dagre)[\\/]/.test(
                id,
              )
            ) {
              return "editor-vendor";
            }
            if (/[\\/]node_modules[\\/](@tanstack|react-router|react-i18next|i18next)[\\/]/.test(id)) {
              return "app-vendor";
            }
            return "misc-vendor";
          },
        },
      },
    },
    server: {
      fs: {
        allow: [".."],
      },
      // 從環境變數讀取端口，預設 5173
      port: devPort,
      proxy: {
        "/api": {
          // 開發模式時讓 proxy 跟隨 backend PORT，避免 start.sh/start.ps1 啟動到不同 port。
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
      include: [
        "tests/unit/**/*.{test,spec}.{ts,tsx}",
        "tests/integration/**/*.{test,spec}.{ts,tsx}",
      ],
      exclude: ["tests/e2e/**"],
      maxWorkers: 6,
    },
  };
});
