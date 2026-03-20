# Repository Guidelines

## 專案結構與模組分工
- `cmd/`：可執行入口（如 `cmd/gateway`、`cmd/test_ui`，以及 `cmd/test_all`、`cmd/fatek_test` 測試工具）。
- `internal/`：核心應用程式碼（API handlers、datalink services、protocol adapters、virtual simulation）。
- `lib/`：共用協議與演算法函式庫。
- `frontend/`：Vite + React + TypeScript 前端（`src/`、`public/`、ESLint/Vitest 設定）。
- `docs/`：架構、API 與運維文件（含 `docs/swagger/`）。
- `openspec/`：需求規格、變更提案與任務追蹤；重大變更請遵循此流程。

## 建置、測試與開發指令
- `make build`：安裝前端依賴、建置前端，再建置 `bin/test-ui.exe`。
- `make test-ui`：執行已建置的 UI 測試工具。
- `powershell -File scripts/build.ps1`：完整本機建置（前端 + 靜態檔嵌入 + 後端）。
- `go test ./...`：執行所有後端單元/整合測試。
- `golangci-lint run ./...`：依 `.golangci.yml` 執行 Go 靜態分析。
- `cd frontend; npm run dev`：啟動前端開發伺服器。
- `cd frontend; npm run build | npm run test | npm run lint`：前端建置、測試、Lint。

## 程式風格與命名慣例
- Go：使用 `gofmt`/`goimports`；package 名稱維持小寫；錯誤處理優先明確回傳並附上下文。
- Frontend：以 TypeScript 為主、2 空白縮排；元件檔用 `PascalCase.tsx`，hooks 用 `useX.ts`。
- 測試命名：Go 使用 `*_test.go`；前端使用 `*.test.ts(x)`，必要時放在 `__tests__/`。
- 維持 import 與死碼乾淨；Lint 警告視為待處理工作，不帶入主分支。

## 規範來源分層
- `AGENTS.md` 是專案共通規範入口，負責專案目標、工作流、測試與文件要求。
- `CLAUDE.md`、`GEMINI.md` 是不同 Agent 的執行補充，負責工作方式、脈絡與執行邊界。
- `.github/instructions/` 是正式的語言與框架實作規範來源；Agent 在修改對應檔案時必須一併遵守。
- 檔案類型對應如下：
  - `*.go`、`go.mod`、`go.sum`：遵守 `.github/instructions/go.instructions.md`
  - `*.tsx`、`*.jsx`、`*.js`、`*.css`、`*.scss`：遵守 `.github/instructions/reactjs.instructions.md`
  - `*.ts`：同時遵守 `.github/instructions/reactjs.instructions.md` 與 `.github/instructions/typescript-5-es2022.instructions.md`
- 規範衝突時，優先順序為：`AGENTS.md` 專案目標與工作流 -> Agent 專屬文件 -> `.github/instructions/` 細部實作規範。

## 測試規範
- 後端使用 Go `testing`，必要時搭配 `testify`。
- 前端使用 Vitest + Testing Library（`frontend/src/setupTests.ts`）。
- 任何行為變更都要同步新增或調整測試，尤其是協議解析、映射流程、排程邏輯。
- 先補精準單元測試，再補資料庫/協議邊界的整合測試。
- 前端 UI/UX 調整採 TDD 先行：先補測試，再改介面；至少覆蓋 `Studio` 主流程（`device -> source -> tag -> output`）與 `TestPage` 的主要互動流程。
- 前端測試統一由 `frontend/tests/` 作為正式入口並分類管理：
  - `frontend/tests/unit/`：Vitest 單元與頁面互動測試
  - `frontend/tests/integration/`：Vitest 跨模組整合測試
  - `frontend/tests/e2e/`：Playwright 端對端測試
- 新增前端測試時，一律放在 `frontend/tests/` 對應分類；utils、hooks、features、components 已完成實體遷移，頁面測試以 `frontend/tests/` 為正式入口；目前 Gateway 仍有部分 wrapper 匯入 `src/pages/.../__tests__/`。
- Go `*_test.go` 維持與實作檔相鄰，不搬到 root `tests/`，以符合 Go 工具鏈與 package 慣例。

## Commit 與 Pull Request 規範
- 近期提交慣例為簡短、祈使語氣的繁中主旨（例：`修正...`、`完成...`、`補齊...`）。
- 每個 commit 聚焦單一主題，並將對應測試一併提交。
- PR 請附：變更摘要、影響模組、測試證據（`go test ./...`、`npm run test`）、相關 OpenSpec 連結；前端變更需附截圖。

## 前端 UI/UX 專案目標
- `/studio` 是目前 datalink 主產品介面，設計與重構都要優先對齊單人操作情境：`建立/選擇資料來源 -> 在格子上看到資料 -> 簡單設定 Tag -> 對應本地 Modbus -> 寫入資料庫供其他 UI 專案使用`。
- `Studio` 內的 `device -> source -> tag -> output` 四步驟是唯一主流程；不再維持 `SmartDashboard` / `LocalModbusWorkbenchPage` 這類舊頁面的平行心智模型。
- `TestPage` 是專用工程測試工具；其核心用途與操作行為應保持穩定，僅做 UI 風格一致化，不作產品主流程承載。
- 前端改版與 cleanup 時，應優先盤查並清除舊 datalink legacy 結構、舊 redirect、未引用元件與過時設計，避免新主線旁再殘留第二套產品入口。

## 文件化工作流
- 多步驟 UI/UX、架構整理或大型重構任務，預設採 SKILL `planning-with-files` 工作法。
- 專案根目錄需維持 `task_plan.md`、`findings.md`、`progress.md` 三份文件，分別記錄階段計畫、關鍵發現與執行/驗證過程。
- 開始執行前先做 session catchup；每完成一個 phase，要同步更新計畫狀態、測試結果、修改檔案與錯誤紀錄。
- 重要發現、legacy 清單、風險與失敗嘗試不得只留在對話上下文，必須寫入文件以利續作。

## OpenSpec 與 Agent 工作注意事項
- 依規劃 phase 與步驟執行，不可跳步或僅交付最小可動版本。
- 規格導向工作需同步更新 `openspec/changes/.../tasks.md` 的實作進度。
