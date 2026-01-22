## Context

Datalink UI 前端已完整實現（types、services、components、pages），但後端 API handlers 使用記憶體儲存且缺少多個端點，導致前後端無法對齊。此變更將完整實現後端 API 並增強前端功能。

### 現狀分析

| 層級          | 現況                                      | 問題           |
| ------------- | ----------------------------------------- | -------------- |
| 後端 Handlers | 使用 MemoryRepository                     | 重啟即遺失資料 |
| 後端 API      | 缺少 polling-groups、protocols、health 等 | 前端呼叫失敗   |
| 後端 Settings | 硬編碼 Mock                               | 無法實際設定   |
| 前端          | 完整但無法運作                            | 依賴缺失的後端 |

## Goals / Non-Goals

### Goals

- 完整實現 datalink-ui spec 所有需求
- 同時支援 SQLite 與 Postgres 持久化
- 使用 SSE 實現即時預覽功能
- 提供導引式工作流程（6 步驟 Wizard）

### Non-Goals

- Role-based access control 或多租戶
- CSV/Excel 匯入匯出
- 內建歷史查詢 API（由外部 CMS 負責）
- WebSocket 雙向通訊（改用 SSE）

## Decisions

### Decision: 使用 sqlx 統一 SQL Repository

相容 SQLite 與 Postgres，避免維護兩套 Repository 實現。

- **Alternatives considered**: 分別實現 sqlite_repo.go 和 postgres_repo.go
- **Rationale**: 兩者 SQL 語法差異極小（僅 Partition），用條件判斷處理即可

### Decision: 使用 SSE 而非 WebSocket 實現即時預覽

- **Rationale**: SSE 為單向推送，符合預覽需求；後端已有 sse.go 可重用
- **Alternatives considered**: WebSocket 雙向通訊
- **Rationale for rejection**: 預覽不需要客戶端上行，SSE 更簡單

### Decision: 使用 Stepperize 庫實現 Wizard

- **Rationale**: TypeScript 原生、Headless 設計、支援 Zod 驗證、Lifecycle hooks
- **Alternatives considered**: 自行實現、react-use-wizard
- **Rationale for selection**: Stepperize 功能最完整且維護良好

### Decision: Repository 工廠模式

根據配置（環境變數或 YAML）自動選擇 Memory/SQLite/Postgres Repository。

- **Rationale**: 開發時用 Memory，測試用 SQLite，生產用 Postgres

## Risks / Trade-offs

| 風險                    | 影響            | 緩解措施                       |
| ----------------------- | --------------- | ------------------------------ |
| SQLite 與 Postgres 差異 | 部分 SQL 不相容 | 使用標準 SQL，分區邏輯條件處理 |
| SSE 連線中斷            | 預覽停止        | 前端自動重連機制               |
| Wizard 複雜度           | 開發時間增加    | 使用成熟庫 Stepperize          |
| 遷移現有資料            | 記憶體資料遺失  | 記憶體資料為測試用，可重建     |

## Migration Plan

1. 新增 SQL Repositories 但保留 Memory Repositories
2. 新增 Repository 工廠，預設使用 Memory
3. 完成所有 API 擴展
4. 切換至 SQLite 測試
5. 切換至 Postgres 生產
6. 移除 Memory Repository（可選）

## Open Questions

1. ~~儲存層選擇~~ → 已決定同時支援 SQLite + Postgres
2. ~~即時預覽方式~~ → 已決定使用 SSE
3. ~~Wizard UI 模式~~ → 已決定使用 Stepper + Side Preview
