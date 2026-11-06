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
export default defineConfig(({ command, mode }) => {
  // 載入環境變數（從 .env 文件）
  const env = loadEnv(mode, process.cwd(), "");

  // 從環境變數讀取配置，提供預設值
  const devPort = parseInt(env.VITE_DEV_PORT || "5173", 10);
  const proxyTarget = resolveProxyTarget(env);
  // dev 預設先打包再送瀏覽器：/studio/v2 由約 577 個模組請求降為約 38 個，
  // 遠端（Tailscale/LAN）開 dev server 時 F5 不再被逐模組往返拖慢。
  // 實驗功能；遇到 HMR 異常可設 VITE_BUNDLED_DEV=false 回到逐模組模式。
  // 只用於 dev server：vitest 也會載入本設定，bundledDev 下其 SSR transform 不處理 TS，
  // setupTests.ts 會解析失敗，因此測試（VITEST / mode=test）一律關閉。
  const isVitest = Boolean(process.env.VITEST) || mode === "test";
  const bundledDev = command === "serve" && !isVitest && env.VITE_BUNDLED_DEV !== "false";

  return {
    plugins: [react()],
    base: "/",
    experimental: {
      bundledDev,
    },
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
        "@": path.resolve(import.meta.dirname, "./src"),
      },
    },
    build: {
      outDir: "dist",
      emptyOutDir: true,
      manifest: true,
      rolldownOptions: {
        output: {
          codeSplitting: {
            // Group only matched dependencies, never their importing route modules.
            includeDependenciesRecursively: false,
            groups: [{ name(id: string) {
              const dependency = id.replaceAll('\\', '/').split('/node_modules/').at(-1);
              if (!id.includes('node_modules') || !dependency) return null;
              if (/^(react|react-dom|scheduler|react-is|use-sync-external-store)\//.test(dependency)) return 'react-vendor';
              if (/^(recharts|recharts-scale|d3-[^/]+|victory-vendor|internmap|decimal.js-light)\//.test(dependency)) return 'charts-vendor';
              if (/^(@mui|@emotion)\//.test(dependency)) return 'mui-vendor';
              if (/^(@tanstack|react-router|react-router-dom|i18next|react-i18next|i18next-browser-languagedetector)\//.test(dependency)) return 'app-vendor';
              return 'misc-vendor';
            } }],
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
