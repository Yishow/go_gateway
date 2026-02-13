# CLAUDE.md

本檔案提供 AI Agent（Claude/Codex）在本倉庫工作的專屬操作說明。

## 文件分工（先讀）
1. 先讀 `AGENTS.md`：通用貢獻規範（結構、命令、測試、PR）。
2. 再讀 `CLAUDE.md`：AI 工作流程、架構脈絡與實作注意事項。

原則：`AGENTS.md` 是共用規範；本檔補充「AI 需要的上下文與執行邊界」。

## 專案快速脈絡
- 專案目標：工業資料採集閘道，從 PLC 協議讀取資料，經 Datalink 映射後輸出到資料儲存/訊息系統。
- 後端：Go 1.25.x + Gin；前端：React 19 + TypeScript + Vite。
- 部署型態：單一可執行檔整合 API 與嵌入式前端（主要入口 `cmd/test_ui`）。
- 協議重點：Modbus、FATEK、MC Protocol；對應實作在 `internal/protocol/*` 與 `internal/datalink/connector/adapters/*`。

## 關鍵目錄（AI 常用）
- `cmd/test_ui/`：主程式入口與靜態資源。
- `internal/api/`：HTTP/SSE 路由與 handlers。
- `internal/datalink/`：核心業務（device/tag/mapping/point/pollinggroup/collector/connector）。
- `internal/protocol/`：協議客戶端與傳輸層。
- `internal/datalink/schema/migrations/`：資料表與遷移。
- `frontend/src/`：頁面、組件、hooks、services、types、i18n。
- `openspec/`：規格、提案、tasks 驗收來源。

## 系統資料流（高階）
1. 協議連線：`connector/manager` 依設備配置建立/管理連線。
2. 資料採集：`collector/scheduler` 依 polling 設定輪詢 Tag。
3. 映射處理：`mapping` 與 transform 將來源值轉為 Point/Tag 目標格式。
4. 儲存與對外：寫入 storage，並透過 API/SSE 提供前端。

## 常用命令（以目前 repo 為準）
```bash
# 全量建置（前端 + 後端）
make build

# PowerShell 建置腳本（嵌入前端後建置）
powershell -File scripts/build.ps1

# 後端測試與 lint
go test ./...
golangci-lint run ./...

# 前端
cd frontend
npm run dev
npm run lint
npm run test
npm run build
```

## AI 實作流程（必遵守）
1. **Context Check**：先讀 `AGENTS.md` 與對應 `openspec/specs/*/spec.md`。
2. **Plan by Phase**：明確本次 phase、輸入、輸出、驗收條件。
3. **Implement Fully**：不可跳步、不可僅最小可動；需完成錯誤處理與邊界條件。
4. **Validate**：執行受影響範圍測試與 lint。
5. **Report**：回報修改檔案、驗證結果、風險與後續建議。

## 開發路徑指引
### 新增協議適配器
1. 在 `internal/protocol/<protocol>/` 實作 client/transport/frame（視需求）。
2. 在 `internal/datalink/connector/adapters/` 增加對應 adapter。
3. 在 `internal/datalink/connector/registry.go` 註冊協議。
4. 補上單元與整合測試（連線、讀寫、錯誤路徑）。

### 新增 API 能力
1. 在 `internal/api/handlers/` 新增 handler。
2. 於 `internal/api/router.go` 掛載路由。
3. 在 `internal/datalink/<module>/` 補 service/repo。
4. 若影響契約，更新 OpenSpec 與 Swagger。

### 前端擴充
1. 型別優先更新 `frontend/src/types/`。
2. API 呼叫集中於 `frontend/src/services/`。
3. 伺服器狀態使用 React Query hooks（`frontend/src/hooks/datalink/`）。
4. 使用者文字走 i18n 字典（`frontend/src/i18n/locales/`）。

## OpenSpec 規格流程（必對齊）
- 先看現行 spec：`openspec/specs/`。
- 規格導向實作時，必須同步更新 `openspec/changes/.../tasks.md`。
- 新功能/架構變更需提案；完成後歸檔至 `openspec/changes/archive/`。
- 實作與驗收衝突時，以 spec 與 tasks 為準並先回報差異。

## 驗證基準
- 後端最低：`go test ./...`、`golangci-lint run ./...`
- 前端最低：`cd frontend && npm run lint && npm run test && npm run build`
- 任何未執行項必須在回報中說明原因、風險與建議補驗證步驟。

## 文件維護規則
- 通用規範調整：改 `AGENTS.md`。
- AI 流程、架構脈絡、執行邊界調整：改 `CLAUDE.md`。
