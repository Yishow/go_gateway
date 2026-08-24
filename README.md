# go-gateway

工業資料採集閘道（Go + React）：從 PLC 協議讀取資料，經 Datalink 映射後輸出到資料儲存或訊息系統，並提供嵌入式 Web UI。

## 專案總覽
- 目標：OT 到 IT 的資料橋接與可視化配置，支援設備連線、規則規劃、Tag/Mapping 管理與輸出綁定。
- 技術棧：Go `1.25.5` + Gin（後端），React `19` + TypeScript `5` + Vite `7`（前端）。
- 部署型態：單一可執行檔整合 API 與前端（主要入口 `cmd/test_ui`）。
- 前端嵌入：前端 build 產物嵌入 `cmd/test_ui/static`，由 `internal/web/embed.go` 提供 SPA 靜態資源。
- 產品主線路由：
  - `/studio/v2`：datalink 唯一面向使用者的 setup 入口。
  - `/studio/runtime`：setup 後的 focused runtime 觀察面。
  - `/test`：工程測試工具入口；`/gateway/*`：experimental surfaces。
  - `/studio`：不再提供 dedicated route；刪除後與任意 unknown route 共用既有 generic policy。
- `/datalink/*` 舊路由：generic landing 收斂到 `/studio/v2`；其他相容路徑須以 route inventory 與測試證據為準，不得重新引入 `/studio`。
- 目前 output 主線：Local Modbus 與 Database，Database 正式支援 `SQLite` / `PostgreSQL`。

## 文件閱讀要求
- 本 repo 不能只讀單一規範文件。
- 共通規範：`AGENTS.md`。
- Agent 專屬補充：`CLAUDE.md`（或其他 agent 文件）。
- 實際執行時需同時參考兩者，不能以單一文件取代另一份。

## 模組結構
- `cmd/`：可執行入口（`cmd/test_ui`、`cmd/test_all`、`cmd/fatek_test` 等）。
- `internal/api/`：HTTP / SSE 路由與 handlers，Swagger 掛載。
- `internal/datalink/`：核心 domain（device、point、tag、mapping、pollinggroup、runtime、dbtarget、sourcerule 等）。
- `internal/protocol/`：協議客戶端與傳輸層（Modbus、FATEK、MC Protocol）。
- `internal/datalink/schema/migrations/`：資料庫 migration。
- `lib/`：共用協議與演算法函式庫。
- `frontend/`：Vite + React + TypeScript 前端與測試。
- `docs/`：架構、流程、維運與 swagger 文件。
- `openspec/`：需求規格、變更提案、任務追蹤。

## 建置與測試指令
### 常用建置
- `make build`：安裝前端依賴 + 前端建置 + 後端建置（輸出 `bin/test-ui.exe`）。
- `powershell -File scripts/build.ps1`：Windows 本機完整建置（含複製 embed static）。
- `make clean`：清理 `frontend/dist`、`frontend/node_modules`、`bin/`。
- `make gen-docs`：重產 Swagger 到 `docs/swagger/`。
- `make test-ui`：建置完成後執行 `bin/test-ui.exe`。

### 後端驗證
- `go test ./...`
- `go vet ./...`
- `golangci-lint run ./...`

### 前端驗證
- `cd frontend && npm run dev`
- `cd frontend && npm run lint`
- `cd frontend && npm run test`
- `cd frontend && npm run build`
- `cd frontend && npm run test:e2e`
- `cd frontend && npm run test:gateway:unit`
- `cd frontend && npm run test:gateway:e2e`
- `cd frontend && npm run test:gateway:gate`

### Gate / Migration 工具
- `make gate-smoke`
- `make gate-final`
- `make gate-soak`
- `make gatev11`
- `make gatev12`
- `make points-precheck-up`
- `make points-precheck-down`
- `make points-migrate-up`
- `make points-migrate-down`（需明確確認 down 策略）
- `make longtask-smoke`
- `make check-lines`（檔案行數規範檢查）

## 程式碼樣式與命名規則
### Go
- 使用 `gofmt` / `goimports`，package 名稱維持小寫、單一語意。
- exported symbol 必須有註解；維持 happy path 左對齊，優先早退。
- package 宣告每檔案只能有一個；避免重複宣告。
- 命名：
  - package：小寫單字，避免 `util/common/base` 類泛稱。
  - 變數/函式：`camelCase`，匯出名稱 `PascalCase`。
  - 介面：優先 `-er` 語意（如 `Reader`、`Writer`）。

### Frontend / TypeScript
- 以 TypeScript 為主，2 空白縮排，ES modules。
- 元件檔 `PascalCase.tsx`，hooks `useX.ts`。
- 避免 `any`；必要時用 `unknown` + narrowing。
- 前端資料流：
  - 型別放 `frontend/src/types/`
  - API 呼叫集中 `frontend/src/services/`
  - 伺服器狀態走 React Query hooks
  - 使用者文字走 `frontend/src/i18n/locales/`

## Error Handling Pattern
### 後端（Go）
- 函式呼叫後立即檢查 `err`，不要忽略錯誤。
- 錯誤往上傳遞時用 `fmt.Errorf("...: %w", err)` 補上下文。
- 錯誤訊息維持小寫、無句點。
- 需要語意判斷時使用 sentinel/custom error，並透過 `errors.Is` / `errors.As` 檢查。
- 避免同時「log + return」同一錯誤，統一在合適層級處理。

### 前端（React/TS）
- 非同步流程使用 `try/catch`，不可吞錯（UI/handler/service 都一樣）。
- 區分 user-facing error 與 internal error，顯示訊息需可理解且不洩漏敏感資訊。
- 保留 fallback/error state（含 Error Boundary）與 loading state，避免 silent failure。

## 測試指引與測試要求
- 後端測試：Go `testing`（必要時 `testify`），`*_test.go` 與實作檔相鄰。
- 前端測試：Vitest + Testing Library；E2E 使用 Playwright。
- 前端正式測試入口：`frontend/tests/unit`、`frontend/tests/integration`、`frontend/tests/e2e`。
- 任何行為變更都必須同步新增/調整測試，特別是：
  - 協議解析
  - 映射流程
  - 排程邏輯
  - runtime lifecycle
  - output binding
- 前端 UI/UX 調整採 TDD 先行：先補測試再改 UI，至少覆蓋 `Studio` 主流程與 `TestPage` 主要互動。
- 建議本地最小驗證：
  - `go test ./... && go vet ./... && golangci-lint run ./...`
  - `cd frontend && npm run lint && npm run test && npm run build`

## 安全考量
- 禁止把真實帳密、Token、DSN、設備憑證等 secrets 放入程式碼、測試、文件、OpenSpec、commit message。
- 所有外部輸入（HTTP payload、query/route params、protocol config、DB connector config、檔案內容）都要驗證與正規化。
- 前端不可拼接不受信任 HTML，不可引入動態程式碼執行。
- DB 層必須走既有 abstraction 或參數化查詢，避免可注入 SQL。
- `golangci-lint run ./...`（含 `gosec`）需維持乾淨。
- 對外連線 probe 由 backend 主機發起；文件與 UX 文案必須避免誤導成瀏覽器直連。

## 禁止事項
- 禁止新增 `/studio` 專用 route、handler、tombstone 或 special redirect；保留 `/studio/v2`、`/studio/runtime`、`/test`、`/gateway/*` 的 route identity。
- 禁止把 lint 警告、未使用程式碼、未清理 import 帶入主分支。
- 禁止跳過 OpenSpec 既有規格就直接改需求語意。
- 禁止在未對齊資料模型時，先行放入不完整契約（例如 parser 支援先行但 runtime/model 尚未打通）。
- 禁止以手刻字串 SQL 或跳過輸入驗證來快速修補功能。

## Repo 特定規則
### 規範優先順序
1. `AGENTS.md`（專案目標、流程、測試、安全）
2. Agent 專屬文件（`CLAUDE.md`、`GEMINI.md`）
3. `.github/instructions/*.md`（語言/框架實作規範）

### 檔案類型規範對應
- `*.go`、`go.mod`、`go.sum`：`.github/instructions/go.instructions.md`
- `*.tsx`、`*.jsx`、`*.js`、`*.css`、`*.scss`：`.github/instructions/reactjs.instructions.md`
- `*.ts`：同時遵守 react + TypeScript instructions

### 文件化工作流
- 多步驟任務預設採 `planning-with-files`。
- 維持 `task_plan.md`、`findings.md`、`progress.md` 三份追蹤文件並持續更新。

### 檔案行數規範（強制）
- 目標：單檔 `<= 300` 行；硬上限 `<= 500` 行。
- `> 300` 行：警告，PR 需提供原因與拆分計畫。
- `> 500` 行：CI 阻擋。
- 歷史超長檔僅允許不增加行數的修改，需逐步縮減。
- 強制工具：
  - `scripts/check_file_lines.sh`
  - `.line-limit-ignore`
  - `.github/workflows/file-line-limit.yml`
  - `.githooks/pre-commit`（啟用：`git config core.hooksPath .githooks`）

### OpenSpec 流程
- 先看 `openspec/specs/` 再處理 `openspec/changes/`。
- 規格導向工作需同步更新 `openspec/changes/.../tasks.md` 實作進度。
- 實作與 spec 衝突時先回報，不直接改寫需求。

### Commit / PR
- commit 主旨採簡短祈使語氣繁中（例如：`修正...`、`補齊...`）。
- 每個 commit 聚焦單一主題並附測試。
- PR 需附：
  - 變更摘要
  - 影響模組
  - 測試證據
  - 相關 OpenSpec 連結
  - 前端變更截圖

## 參考文件
- 共通規範：[`AGENTS.md`](./AGENTS.md)
- Agent 補充規範：[`CLAUDE.md`](./CLAUDE.md)
- Go 規範：[`./.github/instructions/go.instructions.md`](./.github/instructions/go.instructions.md)
- React 規範：[`./.github/instructions/reactjs.instructions.md`](./.github/instructions/reactjs.instructions.md)
- TypeScript 規範：[`./.github/instructions/typescript-5-es2022.instructions.md`](./.github/instructions/typescript-5-es2022.instructions.md)
- OpenSpec：[`openspec/`](./openspec)
- 專案文件：[`docs/`](./docs)
