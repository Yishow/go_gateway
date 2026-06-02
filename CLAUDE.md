<!-- SPECTRA:START v1.0.2 -->

# Spectra Instructions

This project uses Spectra for Spec-Driven Development(SDD). Specs live in `openspec/specs/`, change proposals in `openspec/changes/`.

## Use `/spectra-*` skills when:

- A discussion needs structure before coding → `/spectra-discuss`
- User wants to plan, propose, or design a change → `/spectra-propose`
- Tasks are ready to implement → `/spectra-apply`
- There's an in-progress change to continue → `/spectra-ingest`
- User asks about specs or how something works → `/spectra-ask`
- Implementation is done → `/spectra-archive`
- Commit only files related to a specific change → `/spectra-commit`

## Workflow

discuss? → propose → apply ⇄ ingest → archive

- `discuss` is optional — skip if requirements are clear
- Requirements change mid-work? Plan mode → `ingest` → resume `apply`

## Parked Changes

Changes can be parked（暫存）— temporarily moved out of `openspec/changes/`. Parked changes won't appear in `spectra list` but can be found with `spectra list --parked`. To restore: `spectra unpark <name>`. The `/spectra-apply` and `/spectra-ingest` skills handle parked changes automatically.

<!-- SPECTRA:END -->

# CLAUDE.md

本檔案提供 AI Agent（Claude / Codex）在本倉庫工作的專屬操作說明，內容已與 `AGENTS.md` 對齊；若兩者重疊，`AGENTS.md` 是共通規範來源，`CLAUDE.md` 則補充 agent 執行脈絡與邊界。

## 文件分工（先讀）
1. 先讀 `AGENTS.md`：通用貢獻規範（專案總覽、命令、樣式、測試、安全、文件工作流）。
2. 依修改檔案類型讀取 `.github/instructions/` 下對應規範：
   - Go：`go.instructions.md`
   - React：`reactjs.instructions.md`
   - TypeScript：`typescript-5-es2022.instructions.md`
3. 再讀 `CLAUDE.md`：AI 工作流程、架構脈絡與實作注意事項。

原則：
- `AGENTS.md` 是共用規範與專案目標入口。
- `CLAUDE.md` 補充 AI 需要的上下文、工作流程與執行邊界。
- `.github/instructions/` 提供語言與框架層的細部實作規範。
- **不可只讀單一文件**：即使在 Claude Code，也不能只讀 `CLAUDE.md`；必須先讀 `AGENTS.md` 再讀 `CLAUDE.md`。

## 專案總覽
- 專案目標：工業資料採集閘道，從 PLC 協議讀取資料，經 Datalink 映射後輸出到資料儲存 / 訊息系統。
- 後端：Go 1.25.x（目前 `go.mod` 為 1.25.5）+ Gin；前端：React 19 + TypeScript 5 + Vite 7。
- 部署型態：單一可執行檔整合 API 與嵌入式前端，主要入口為 `cmd/test_ui`。
- 前端主線入口：`/studio/v2`；`/studio` 保留為既有完整工作台 fallback，舊 datalink generic landing routes 應收斂 redirect 到 `/studio/v2`，工程測試工具集中於 `/test`。
- Database output 正式支援 `SQLite` 與 `PostgreSQL`。

## 關鍵目錄（AI 常用）
- `cmd/test_ui/`：主程式入口與嵌入式靜態資源。
- `internal/api/`：HTTP / SSE 路由、handlers、Swagger 掛載。
- `internal/datalink/`：核心業務（device、point、tag、mapping、pollinggroup、collector、connector、runtime、dbtarget、sourcerule 等）。
- `internal/protocol/`：協議客戶端與傳輸層。
- `internal/datalink/schema/migrations/`：資料表與遷移。
- `frontend/src/`：頁面、組件、hooks、services、types、i18n。
- `frontend/tests/`：前端 unit / integration / e2e 正式測試入口。
- `openspec/`：規格、提案、tasks 驗收來源。

## 系統資料流（高階）
1. 協議連線：`connector/manager` 依設備配置建立 / 管理連線。
2. 資料採集：`collector/scheduler` 依 polling 設定輪詢 Tag。
3. 映射處理：`mapping` 與 transform 將來源值轉為 Point / Tag / Output 目標格式。
4. 儲存與對外：寫入 storage，並透過 API / SSE 提供前端。

## 常用命令（以目前 repo 為準）
```bash
# 全量建置（frontend install/build + backend build）
make build

# Windows 本機完整建置（含 embed static）
powershell -File scripts/build.ps1

# 清理建置產物
make clean

# 重新產生 Swagger
make gen-docs

# 後端驗證
go test ./...
go vet ./...
golangci-lint run ./...

# 前端
cd frontend
npm run dev
npm run lint
npm run test
npm run build
npm run test:e2e
npm run test:gateway:unit
npm run test:gateway:e2e
npm run test:gateway:gate

# Gate / migration / smoke workflows
make gate-smoke
make gate-final
make gate-soak
make gatev11
make gatev12
make points-precheck-up
make points-precheck-down
make points-migrate-up
make points-migrate-down
make longtask-smoke
make check-lines
```

## 程式碼樣式與設計準則
- Go：遵守 `gofmt` / `goimports`、小寫 package 名稱、明確錯誤回傳與上下文包裝；避免重複 `package` 宣告。
- Go：優先早退、保持 happy path 左對齊；先重用標準函式庫與既有 service / helper，再考慮新增 abstraction。
- Frontend：TypeScript 為主、2 空白縮排、元件 `PascalCase.tsx`、hooks `useX.ts`。
- Frontend：遵守 ES modules、TypeScript `strict`、既有 ESLint / React Hooks 規則；避免 `any`，必要時以更窄的型別或 guard 取代。
- 前端資料流：型別放 `frontend/src/types/`，API 呼叫集中在 `frontend/src/services/`，伺服器狀態優先經由 React Query hooks 管理，使用者文字走 i18n 字典。

## 命名規則
- Go package：小寫單字、單數語意，避免 `util/common/base`。
- Go 匯出符號：`PascalCase`；非匯出符號：`camelCase`。
- Go 介面：優先以行為命名（例如 `Reader`、`Writer`）。
- Frontend 元件：`PascalCase.tsx`；hooks：`useX.ts`。
- 測試檔：Go `*_test.go`；前端 `*.test.ts` / `*.test.tsx`。

## Error Handling Pattern
- Go：每次呼叫後立即檢查 `err`；往上傳遞錯誤時用 `fmt.Errorf("...: %w", err)` 包裝上下文。
- Go：錯誤訊息維持小寫、無句點；需語意判斷時使用 sentinel/custom error + `errors.Is` / `errors.As`。
- Go：避免同一錯誤在多層重複記錄（`log + return`），應在適當層級單點處理。
- Frontend：非同步流程必須顯式處理失敗（`try/catch` + error state / fallback），禁止 silent failure。
- Frontend：user-facing error 文案需可行動，不可暴露敏感資訊（憑證、內網位址、SQL 細節等）。

## 測試與驗證指引
- 後端最低驗證基準：`go test ./...`、`go vet ./...`、`golangci-lint run ./...`。
- 前端最低驗證基準：`cd frontend && npm run lint && npm run test && npm run build`；CI / 非 watch 模式可用 `npm test -- --run`。
- 任何行為變更都要同步新增或調整測試，尤其是協議解析、映射流程、排程邏輯、runtime lifecycle、output binding。
- 前端 UI / UX 調整採 TDD 先行：先補測試，再改介面；至少覆蓋 `Studio` 主流程與 `TestPage` 的主要互動。
- 前端正式測試入口以 `frontend/tests/` 為主；Go 測試檔維持與實作檔相鄰。

## 測試要求（提交前）
- 行為變更必須附對應測試或明確說明無法補測的技術原因。
- 文件任務至少執行 `git diff --check`，確認無格式/空白異常。
- 若未能在當前環境執行完整測試，回報必須明確列出「已執行」與「未執行」項目及風險。

## 安全考量
- 不要把真實帳號、密碼、Token、DSN、設備憑證或其他 secrets 寫入程式碼、測試、文件、OpenSpec 或 commit message。
- 對所有外部輸入做驗證與正規化；優先重用既有型別、validator、schema、service 層檢查，不要在 UI 或 handler 端偷偷吞錯。
- 前端顯示外部資料時維持 React 預設 escaping，不拼接不受信任 HTML，不引入動態程式碼執行。
- 資料庫與儲存層修改應使用既有 repository / service abstraction 或參數化查詢；避免手刻可注入 SQL。
- Go 靜態分析以 `golangci-lint run ./...` 為基準，`gosec` 必須維持乾淨；涉及 shell、檔案路徑、網路連線與序列埠時要特別保守。
- 連線測試 / probe 是由 backend 主機發起；在文件、診斷與 UX 文案中要避免誤導成瀏覽器端直連。

## 禁止事項
- 禁止把 secrets 寫入任何版本化檔案或 commit message。
- 禁止在未對齊 OpenSpec 的情況下自行改寫需求語意。
- 禁止把 lint 警告、未使用 import 或 dead code 直接合入主分支。
- 禁止以手刻字串 SQL 或繞過 service/repository abstraction 快速上線功能。

## 其他 Repo 特定規則
- `cmd/test_ui` 為單一可執行檔入口，前端資產嵌入 `cmd/test_ui/static`。
- 若任務涉及 `/studio`、`/studio/v2`、`/studio/runtime`、`/test`、`/gateway/*` 或 `docs/technical/studio-surface-inventory/`，先讀 `docs/technical/studio-surface-inventory/START_HERE.md`、`docs/technical/studio-surface-inventory/context.json`、`docs/technical/studio-surface-inventory/CURRENT_STATE.md`，只在不足以回答問題時再展開完整 md/html 文件。
- 只要修改 `docs/technical/studio-surface-inventory/` 內任何文件，必須同步寫入 `docs/technical/studio-surface-inventory/changelog.sqlite`，避免 inventory 在沒有明確紀錄下被悄悄改動。
- `studio-surface-inventory` changelog 一律使用 `go run ./cmd/studio_inventory_changelog ...` 管理；至少要留下 `summary`、`surface`、`files`、`reason`。

## 檔案行數規範（強制）
- 目標：單檔不超過 300 行；硬上限 500 行，MD/HTML不在此限(超過500行仍需注意，要分檔案)。
- 規則：
  - `> 300` 行：警告，必須在 PR 提供原因與拆分計畫。
  - `> 500` 行：CI 阻擋。
  - 歷史超長檔僅允許不增加行數的變更，應逐步縮減。
- 工具與落地：
  - `scripts/check_file_lines.sh`
  - `.line-limit-ignore`
  - `.github/workflows/file-line-limit.yml`
  - `.githooks/pre-commit`（可用 `git config core.hooksPath .githooks` 啟用）

## AI 實作流程（必遵守）
1. **Context Check**：先讀 `AGENTS.md`、對應 `.github/instructions/`。
2. **Implement Fully**：不可跳步、不可僅交付最小可動；需完成錯誤處理與邊界條件。
3. **Validate**：執行受影響範圍測試與 lint；未執行項必須明確說明風險。
4. **Report**：回報修改檔案、驗證結果、風險與後續建議。

## 文件與規範優先順序
1. `AGENTS.md`
   專案目標、測試要求、安全原則、文件工作流與 UI / UX 主線。
2. `CLAUDE.md`
   Agent 執行邊界、架構脈絡、工作方式與 repo 特定注意事項。
3. `.github/instructions/*.md`
   依檔案類型套用的語言與框架實作規範。

## 文件維護規則
- 通用規範調整：改 `AGENTS.md`。
- Agent 執行脈絡、工作方式、邊界或 repo 補充上下文調整：改 `CLAUDE.md`。
