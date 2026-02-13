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

## 測試規範
- 後端使用 Go `testing`，必要時搭配 `testify`。
- 前端使用 Vitest + Testing Library（`frontend/src/setupTests.ts`）。
- 任何行為變更都要同步新增或調整測試，尤其是協議解析、映射流程、排程邏輯。
- 先補精準單元測試，再補資料庫/協議邊界的整合測試。

## Commit 與 Pull Request 規範
- 近期提交慣例為簡短、祈使語氣的繁中主旨（例：`修正...`、`完成...`、`補齊...`）。
- 每個 commit 聚焦單一主題，並將對應測試一併提交。
- PR 請附：變更摘要、影響模組、測試證據（`go test ./...`、`npm run test`）、相關 OpenSpec 連結；前端變更需附截圖。

## OpenSpec 與 Agent 工作注意事項
- 依規劃 phase 與步驟執行，不可跳步或僅交付最小可動版本。
- 規格導向工作需同步更新 `openspec/changes/.../tasks.md` 的實作進度。
