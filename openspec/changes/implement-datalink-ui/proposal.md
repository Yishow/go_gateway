# Change: Implement Complete Datalink UI

## Why

目前 Datalink UI 前後端嚴重不對齊：前端已完整定義 API 客戶端與類型，但後端使用記憶體儲存（重啟遺失）、多個 API 端點缺失、Settings 為硬編碼 Mock。這導致 UI 完全無法正常運作。

## What Changes

### 後端 API 完善

- **BREAKING** 所有 handlers 從 MemoryRepository 切換至 SQLRepository
- 新增 `GET /datalink/polling-groups` CRUD API
- 新增 `GET /datalink/protocols` API
- 新增 `GET /datalink/health` API
- 擴展 Tag Handler: `POST /tags/:id/activate`, `POST /tags/:id/retire`, `POST /tags/batch`, `POST /tags/validate-key`
- 擴展 Device Handler: `POST /devices/test-batch`
- 擴展 Point Handler: `POST /points/:id/poll`, `POST /points/poll`
- 擴展 Mapping Handler: `POST /mappings/validate-pipeline`
- 重構 Settings Handler 為真實實現

### 儲存層持久化

- 實現 SQL Repository（相容 SQLite + Postgres）
- 整合現有 migrations
- Repository 工廠模式

### SSE 即時預覽

- 新增 `GET /datalink/preview/stream` SSE 端點
- 推送 raw → step results → final value

### 導引式工作流程 UI

- 使用 Stepperize 庫實現 6 步驟 Wizard
- 整合 SSE 即時預覽面板

## Impact

- **Affected specs**: `datalink-ui`, `datalink-api`, `timeseries-storage`
- **Affected code**:
  - `internal/api/handlers/*.go` - 所有 Datalink handlers
  - `internal/api/router.go` - 路由註冊
  - `internal/datalink/*/sql_repo.go` - 新增 SQL repositories
  - `internal/datalink/pollinggroup/` - 新增模組
  - `internal/datalink/settings/` - 新增模組
  - `web/test-ui/src/components/datalink/wizard/` - 新增 Wizard 組件
