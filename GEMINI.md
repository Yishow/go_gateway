# GEMINI.md

本檔案提供 Gemini 系列 Agent 在本倉庫工作的補充說明。

## 文件分工（先讀）
1. 先讀 `AGENTS.md`：專案共通規範、測試要求、UI/UX 主線與文件化工作流。
2. 依修改檔案類型讀取 `.github/instructions/` 下對應規範：
   - Go：`go.instructions.md`
   - React：`reactjs.instructions.md`
   - TypeScript：`typescript-5-es2022.instructions.md`
3. 再讀 `GEMINI.md`：Gemini Agent 的專案脈絡與執行補充。

原則：
- `AGENTS.md` 管專案怎麼做。
- `.github/instructions/` 管程式怎麼寫。
- `GEMINI.md` 補充 Agent 如何在本專案內工作。

## 專案脈絡
- 專案目標：工業資料採集閘道，從 PLC 協議讀取資料，經 Datalink 映射後輸出到資料儲存與訊息系統。
- 技術棧：後端為 Go 1.25.x + Gin；前端為 React 19 + TypeScript + Vite。
- 部署型態：以單一可執行檔整合 API 與嵌入式前端，主要入口為 `cmd/test_ui`。
- 關鍵協議：Modbus、FATEK、MC Protocol；對應實作位於 `internal/protocol/*` 與 `internal/datalink/connector/adapters/*`。

## 前端 UI/UX 目前主線
- `SmartDashboard` 是 datalink 主產品介面，應優先對齊單人操作情境：
  `建立/選擇資料來源 -> 在格子上看到資料 -> 設定 Tag -> 對應本地 Modbus -> 寫入資料庫供其他 UI 專案使用`
- `LocalModbusWorkbenchPage` 是後段工作台，應與 `SmartDashboard` 形成清楚接力。
- `TestPage` 是專用工程測試工具，只做風格一致化，不承載產品主流程。
- UI/UX 改造前需先盤查 `SmartDashboard` 的 legacy 結構、舊 redirect、未引用元件與過時設計。

## Gemini 工作流程
1. 先確認本次任務涉及的規範來源：`AGENTS.md`、`.github/instructions/`、必要時 `openspec/specs/*`。
2. 多步驟任務預設採 `planning-with-files`，在專案根目錄維持 `task_plan.md`、`findings.md`、`progress.md`。
3. 依 phase 執行，不跳步、不只交付最小可動版本。
4. 實作後執行受影響範圍的測試、lint 與建置。
5. 回報修改檔案、驗證結果、風險與後續建議。

## 驗證基準
- 後端最低：`go test ./...`、`golangci-lint run ./...`
- 前端最低：`cd frontend && npm run lint && npm run test && npm run build`
- 任何未執行項都必須在回報中說明原因、風險與建議補驗證步驟。

## 語言要求
- 回答與回報使用繁體中文。
