## Why

目前前端採用 route-level lazy loading 與 vendor chunk 後，嵌入式前端不再是單一 entry asset；若建置與啟動流程只複製部分輸出或清空 static 目錄，fresh clone 可能無法通過 `//go:embed static` 編譯，已建置 binary 也可能在進入核心路由時收到 chunk 404。現在需要把 asset graph 完整性、stale asset 清理與 fresh-clone 可編譯性收斂成可驗證的交付契約。

## What Changes

- 建立 `embedded-frontend-delivery` capability，定義 embedded frontend 由 `frontend/dist` 到 `cmd/test_ui/static` 的完整同步與執行期資產載入行為。
- 將重型 route 改為 route-level lazy imports，提供 deterministic Suspense fallback；保留 `/studio/v2` 預設主線、`/studio` fallback、`/studio/runtime` focused monitor、`/test` 工具與 `/gateway/*` experimental route 的既有產品語意。
- 建立可重現的 Vite vendor chunk policy，使 entry、lazy route chunks、vendor chunks、CSS 與其 manifest/reference graph 一起產出且可被 embedded binary 載入。
- 保留 tracked `cmd/test_ui/static/embed-placeholder.txt`，並讓 PowerShell 與 shell 建置／啟動流程在同步前移除 stale generated assets、保留 placeholder、完整複製 `frontend/dist`。
- 加入 fresh-clone compile、frontend build、embedded asset completeness、核心路由 lazy-load unit smoke 與 Playwright browser smoke 的具體驗證目標，確認不存在由缺資產造成的 chunk 404。
- 明確要求 `start.sh` 在找不到 frontend source 時 fail closed，以 nonzero status 結束並阻止同次 backend build；shell fixture 必須保留這個 regression contract。
- 讓 `frontend/playwright.config.ts` 的 webServer 啟動設定支援 Windows 與 POSIX，確保從 `frontend` 直接執行 `npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts` 能啟動並等待測試伺服器。
- 在 final gate 檢查 `start.ps1` 與 `start.sh` 這兩個歷史上已超過 500 行的檔案，修復時行數不得高於各自的 pre-change baseline。
- 不改寫既有 route redirect、產品流程或 route 的業務語意；本 change 只處理載入分割與 embedded delivery。

## Non-Goals

- 不修改任何 route redirect 或 `/studio`、`/studio/v2`、`/studio/runtime`、`/test`、`/gateway/*` 的產品語意。
- 不處理 Workbench autosave、mapping、live-values、API、DB、Air entrypoint、鎖演算法、純測試搬遷、lint/type cleanup 或 font/browser lockfile 更新。

## Capabilities

### New Capabilities

- `embedded-frontend-delivery`: 確保 route/lazy/vendor 前端資產在 fresh clone、建置同步與 embedded runtime 間完整且可驗證地交付。

### Modified Capabilities

（無；既有 datalink route specs 的產品行為保持不變。）

## Impact

- Affected specs: `embedded-frontend-delivery`（新增 capability）；`datalink-ui`、`datalink-workbench-v2-shell`、`post-setup-runtime-dashboard` 僅作為相容性參考，不修改其 requirements。
- Affected code:
  - Modified:
    - `frontend/src/App.tsx`
    - `frontend/vite.config.ts`
    - `.gitignore`
    - `scripts/build.ps1`
    - `start.ps1`
    - `start.sh`
    - `frontend/playwright.config.ts`
  - New:
    - `cmd/test_ui/static/embed-placeholder.txt`
    - `frontend/tests/unit/app-routing-lazy-load.test.tsx`
    - `frontend/tests/e2e/embedded-frontend-delivery.spec.ts`
    - `tests/shell/embedded-frontend-delivery.sh`
  - Removed: none
