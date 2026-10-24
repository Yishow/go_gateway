# 前端最新版套件相容性修復

`pnpm update --latest` 後，Tailwind 4 的 PostCSS 介面變更造成 CSS HTTP 500；
TypeScript 7 移除 `baseUrl`，而 typescript-eslint 尚需 TypeScript 6 的 API。
本次保留升級版本，修復設定、測試型別及 ESLint 10 新規則攔下的程式。

## 相容設定

- `@tailwindcss/postcss` 接手 CSS 編譯；以 `@config` 載入既有色彩、字型與深色設定。
- 保留既有小陰影、模糊、圓角、預設邊框及強制色彩模式的 outline 行為。
  試寫結果的 grid 欄位改用 v4 支援的空白語法。
- `@typescript/native` 是 TypeScript 7.0.2 的 npm alias，提供真正的 `tsc`。
  `typescript` alias 指向官方 `@typescript/typescript6`，只提供 lint 所需 API 與 `tsc6`。
  建置仍使用 TS 7；單檔建置 helper 同時辨識新舊 compiler 路徑。
- `eslint-plugin-react` 7.37.5 的 peer range 尚未列出 ESLint 10。
  npm override 僅將此外掛的 ESLint 解析指向根目錄版本；既有設定只有停用的
  React 規則，實際啟用的 Hooks 規則來自 `eslint-plugin-react-hooks`。
  未放寬 lint 規則；npm 與 pnpm lockfile 均同步。
- Vitest 5 使用 `jest-dom/vitest`，測試 actions 沿用元件實際 props 型別。

## 本機驗證

- PASS：`pnpm run lint`、`pnpm run build`（含 TS 7 型別檢查）。
- PASS：`pnpm exec vitest run`，218 files / 1,441 tests。
- PASS：`node --test scripts/tests/frontend-styles.test.mjs scripts/tests/build-single.test.mjs`，22 tests。
- PASS：`pnpm install --frozen-lockfile`；獨立暫存目錄的實際 `npm ci`。
- PASS：`scripts/tests/launcher-live-smoke.py`，獨立副本與可丟棄 SQLite；
  fresh bootstrap、API proxy、Air reload、編譯錯誤後恢復、Vite HMR、TERM 清理與埠／lock 釋放。
  使用暫存 GOCACHE；macOS 原生檔案事件需主機權限，沙盒內無法驗證 HMR。
- PASS：專用 headless Chromium 載入 `/studio/v2`，無 JavaScript error／Vite overlay；
  實測背景、邊框、圓角、陰影與模糊值。截圖保留於本次交接的暫存產物。
- PASS：`git diff --check`、`make check-lines`；歷史超長檔未增行。
- NOT RUN：Windows、真設備、正式資料庫及現場驗收。

非阻斷訊息：Vitest 的 jsdom navigation 提示，以及開發模式中 PostCSS plugin
未傳遞 `from` 的上游警告；CSS 已編譯、HTTP 200，瀏覽器無載入錯誤。

參考：[Tailwind 4 升級指南](https://tailwindcss.com/docs/upgrade-guide)、
[TypeScript 7 官方並存方案](https://devblogs.microsoft.com/typescript/announcing-typescript-7-0/#running-side-by-side-with-typescript-60)。

提交前的現況覆核與限制見 [分批驗證紀錄](remaining-wip-review-20261025.md)；上列數字保留原驗證日期的結果。
