## Why

目前 repo 已有 `GET /api/v1/datalink/runtime/status`、`GET /api/v1/datalink/runtime/stream` 與舊的 `dashboard/*` 統計端點，但它們對 post-setup runtime dashboard 來說仍缺少一致且可依賴的契約。若直接讓前端依現況開做，Gemini 需要一邊猜事件形狀、一邊補 fallback 語意，之後高機率返工。

## What Changes

- 定義一組專供 post-setup runtime dashboard 使用的後端監看契約，涵蓋 device-scoped runtime snapshot 與 SSE stream。
- 明確化 `GET /api/v1/datalink/runtime/status` 的回應欄位、單 device 過濾語意、stale/error 判定與 dashboard fallback 用途。
- 明確化 `GET /api/v1/datalink/runtime/stream` 的事件集合與資料形狀，至少包含 `value`、`status`、`heartbeat`，並定義何時推送設備狀態變化。
- 對齊 scheduler freshness、circuit breaker、last error、point live value 的推導規則，讓 snapshot 與 stream 不再各說各話。
- 補齊 router / handler / runtime service 層的契約測試，讓後續 dashboard route 與 Step 4 handoff 可以直接依賴。

## Capabilities

### New Capabilities

- `runtime-dashboard-backend`: 定義 post-setup runtime dashboard 所依賴的後端 snapshot 與 SSE 監看契約。

### Modified Capabilities

- `datalink-api`: 補強 runtime monitoring API requirements，讓 `/api/v1/datalink/runtime/status` 與 `/api/v1/datalink/runtime/stream` 成為正式可依賴的 dashboard contract。

## Impact

- Affected specs: `runtime-dashboard-backend`, `datalink-api`
- Affected code:
  - Modified:
    - `internal/api/router.go`
    - `internal/api/handlers/runtime_handler.go`
    - `internal/api/handlers/runtime_stream_handler.go`
    - `internal/api/router_runtime_test.go`
    - `internal/datalink/runtime/service.go`
    - `internal/datalink/runtime/stream.go`
    - `internal/datalink/runtime/ingestor.go`
    - `frontend/src/services/datalink.ts`
  - New:
    - `internal/api/handlers/runtime_handler_contract_test.go`
    - `internal/api/handlers/runtime_stream_handler_contract_test.go`
