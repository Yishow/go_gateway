# go-gateway

工業資料採集閘道（Go + React），提供 PLC 協議連接、資料映射與嵌入式 Web UI。

## 30 秒上手
- 先讀通用貢獻規範：[`AGENTS.md`](./AGENTS.md)
- AI Agent 再讀執行補充：[`CLAUDE.md`](./CLAUDE.md)

## 專案文件
- 架構與分析：[`docs/`](./docs)
- OpenSpec 需求與變更流程：[`openspec/`](./openspec)

## Datalink Studio 目前主線
- 新主線入口：`/studio`
- 建議操作順序：
  1. Step 1 設定資料來源，並分開查看 `connect` / `probe` 診斷結果。
  2. Step 2 建立或還原 `Source Rule`，在格狀畫布查看 `planned / used / unmanaged / conflict` 狀態與即時值。
  3. Step 3 以 **review-first** 方式檢查系統自動建立的 Tag / Mapping；只有例外狀況才進行修正、補綁或解綁。
  4. Step 4 直接在 Local Modbus register 或 Database schema/column 表面上綁定輸出。
- 目前 Database output flow 以 `SQLite` 與 `PostgreSQL` 為正式支援範圍。
- 舊 `datalink` legacy routes 會收斂 redirect 到 `/studio`；測試工具則統一收斂到 `/test`。
- `SmartDashboard` 與其他舊 datalink 頁面不再是對外主入口。
