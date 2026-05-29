## Context

目前 repo 已經同時存在三套彼此沒有完全對齊的監看資料來源：`dashboard/stats` 與 `dashboard/device-statuses` 提供偏靜態的列表統計、`runtime/status` 提供 collector 摘要、`runtime/stream` 提供 point value SSE。對 post-setup runtime dashboard 來說，這種狀態會讓前端同時依賴多個端點並自行猜測 freshness、breaker、last error 的語意。

這個 change 的目標不是先決定 dashboard 長相，而是先建立 backend-first contract，讓後續 route 與 Step 4 handoff 都能依賴同一套 runtime truth。它也刻意不沿用現有 `runtime-dashboard` change 的 UI-first 範圍，避免一次把 route、TopBar、五個監控區塊與後端事件設計綁死。

## Goals / Non-Goals

**Goals:**

- 定義一套專供 post-setup runtime dashboard 使用的後端契約，支援「剛 commit 的那台 device」作為預設監看對象。
- 讓 `GET /api/v1/datalink/runtime/status` 成為可輪詢的 snapshot fallback，並回傳 dashboard 需要的 top-level runtime metrics 與 device collector 摘要。
- 讓 `GET /api/v1/datalink/runtime/stream` 成為可長連線的 live feed，並將事件集合收斂為 `value`、`status`、`heartbeat` 三種正式事件。
- 明確規定 freshness、breaker state、last error 與 point live value 的推導規則，避免 snapshot 與 stream 解讀不同。
- 用 router / handler / runtime service 契約測試把行為鎖定，讓後續前端實作不需要猜 API。

**Non-Goals:**

- 不包含新的 dashboard route、TopBar 導覽入口、Step 4 導頁 wiring。
- 不包含 Prometheus `/metrics` exporter、歷史趨勢查詢、告警規則引擎或外部 observability pipeline。
- 不包含 queue backlog、diagnostic logs、fleet-wide overview page 等尚未有穩定資料來源的面板需求。
- 不重寫或刪除既有 `dashboard/*` endpoints；本 change 只定義 post-setup dashboard 的正式依賴來源。

## Decisions

### Decision: Keep runtime monitoring device-scoped

- **Rationale**: 你已經決定 post-setup dashboard 預設只看剛完成設定的 device。這與 `runtime/stream` 目前要求 `device_id` 的現況一致，也能避免把整個 fleet 的高頻 point 更新推給單一 EventSource 連線。
- **Alternatives considered**:
  - **Fleet-wide stream by default**: 可以少一次 device 切換，但會讓事件量與前端 state 爆炸，且與 post-setup 使用情境不符。
  - **Point-only stream without device scope**: 查詢彈性較大，但呼叫端要先知道 point 集合，無法支援「剛 setup 完直接開看」的入口。

### Decision: Reuse runtime service and scheduler health as the source of truth

- **Rationale**: repo 已有 `runtime.Service`、`collector.Scheduler`、point `last_*` 狀態與 breaker stats，可直接組出 snapshot 與 stream。新增平行 telemetry service 只會複製狀態並引入漂移。
- **Alternatives considered**:
  - **Introduce a dedicated telemetry aggregator**: 只有在需要跨程序或歷史聚合時才值得，目前只會增加同步成本。
  - **Build dashboard contract from `dashboard/*` handlers**: 這些端點偏靜態統計，無法成為 post-setup live monitoring 的 single source of truth。

### Decision: Align runtime snapshot and stream under one derived status model

- **Rationale**: dashboard 的困難不在於少一個 endpoint，而在於同一台設備在 snapshot 顯示 `warning`、stream 又只看到 `value` 事件時，前端無法知道該信誰。這次要把 `status`、`points_stale`、`points_error`、`last_read_at`、`last_error`、`breaker_state` 的推導規則集中。
- **Alternatives considered**:
  - **Let frontend derive status from raw events**: 會讓每個 client 重新實作 stale/error 規則，驗證成本高且容易漂移。
  - **Stream values only and poll status separately without alignment rule**: 看似簡單，但實際上會留下 race condition 與不同步問題。

### Decision: Keep the stream loss-tolerant instead of blocking runtime collection

- **Rationale**: 低資源閘道器上的首要目標是採集不中斷，不是保證每個瀏覽器都拿到每一筆 live event。現有 subscriber channel 已採 non-blocking broadcast，這個 change 只把這個取捨正式化，讓前端知道 stream 是近即時監看，不是精確事件回放。
- **Alternatives considered**:
  - **Block producer until every subscriber consumes events**: 會把瀏覽器速度反向影響 scheduler 與寫入鏈路，風險不可接受。
  - **Persist every live event for replay**: 這是另一個產品級需求，超出本 change 範圍。

## Implementation Contract

- **Behavior**:
  - 後端 SHALL 支援 `GET /api/v1/datalink/runtime/status` 作為 post-setup dashboard 的 snapshot fallback。當帶 `device_id=<id>` 時，回應 SHALL 僅包含該 device 的 collector 摘要；未帶 filter 時，回應 SHALL 可列出全部 collector 摘要供未來全局頁面重用。
  - 後端 SHALL 支援 `GET /api/v1/datalink/runtime/stream?device_id=<id>` 作為該 device 的 live feed。若帶 `point_ids=a,b,c`，stream SHALL 只推送這些 points 的 `value` 事件。
  - stream SHALL 持續推送 `heartbeat`，並在設備摘要狀態變動時推送 `status` 事件；`value` 事件 SHALL 在 runtime 收到 point collected value 後發出。
- **Interface / data shape**:
  - `GET /api/v1/datalink/runtime/status`
    - Query: `device_id` optional
    - Response top-level fields SHALL include `running`, `uptime_seconds`, `metrics`, `collectors`
    - `metrics` SHALL include at least `collected_total`, `write_success_total`, `write_error_total`, `mapping_error_total`, `point_state_error_total`
    - each collector SHALL include at least `device_id`, `device_name`, `protocol`, `status`, `points_total`, `points_healthy`, `points_stale`, `points_error`, `last_read_at`, `last_error`, `breaker_state`
  - `GET /api/v1/datalink/runtime/stream`
    - Query: `device_id` required, `point_ids` optional comma-separated list
    - `event: value` payload SHALL include `device_id`, `point_id`, `address`, `raw_value`, `transformed_value`, `quality`, `stale`, `timestamp`
    - `event: status` payload SHALL include `device_id`, `status`, `points_total`, `points_healthy`, `points_stale`, `points_error`, `last_read_at`, `last_error`, `breaker_state`
    - `event: heartbeat` payload SHALL include `ts`
- **Failure modes**:
  - Missing `device_id` on `runtime/stream` SHALL return HTTP 400.
  - Runtime stream source unavailable SHALL return HTTP 503 instead of silently hanging.
  - Non-stream-capable response writer SHALL return HTTP 500.
  - Slow subscribers SHALL NOT block runtime collection; intermediate live events are allowed to be dropped when the subscriber buffer is saturated.
- **Acceptance criteria**:
  - Router tests prove `runtime/status` returns metrics plus filtered collectors for a requested device.
  - Stream handler tests prove missing `device_id`, unavailable runtime source, and successful `value` / `status` / `heartbeat` framing.
  - Runtime service tests prove a collected value updates stream payloads and a saturated subscriber does not stall collection.
  - `spectra analyze add-runtime-dashboard-backend-contract --json` reports no Critical/Warning findings and `spectra validate add-runtime-dashboard-backend-contract` passes.
- **Scope boundaries**:
  - In scope: backend runtime contract, event shapes, derived status semantics, contract tests.
  - Out of scope: dashboard page layout, Step 4 navigation, queue/log/history panels, new storage schema.

## Risks / Trade-offs

- `[Risk]` 舊 client 可能把 `runtime/status` 當成只有 collector 清單的端點，新增 `metrics` 後若有脆弱 decoder 可能出現相容性問題 -> `Mitigation`: 只做 additive fields，不移除既有欄位，也不改 collector 既有 key 名稱。
- `[Risk]` 以 point `last_read_at` 加 polling interval 推導 stale，若某些 point 沒有正確同步 polling group，dashboard 可能誤判 -> `Mitigation`: 契約測試要覆蓋有 group / 無 group 兩種情況，並沿用 handler 既有 fallback interval。
- `[Risk]` loss-tolerant stream 會讓前端看不到每一筆值變化 -> `Mitigation`: 明確把 stream 定位成 live monitoring feed，前端需要精確狀態時應重新拉 snapshot。
- `[Risk]` `status` 事件若只在 value path 觸發，部分 breaker-only 狀態變化可能延遲被看見 -> `Mitigation`: 在設計與 tasks 中要求 handler/runtime service 將 breaker-derived status 納入狀態推送條件與測試覆蓋。

## Migration Plan

1. 先補 contract tests，鎖住 `runtime/status` 與 `runtime/stream` 期望行為。
2. 在 additive 相容前提下擴充 runtime handler 與 stream handler payload。
3. 對齊 runtime service 的 status event 廣播條件後，再讓後續 dashboard route change 接上。
4. 若部署後新 stream 行為出現問題，可先讓前端僅使用 `runtime/status` 輪詢，後端不需 schema rollback。

## Open Questions

- `status` 事件是否需要在沒有新 point value 的 breaker 狀態切換時立即推送，或接受下一次 value update 再同步；此 change 會先以「需要立即推送」為預設方向，實作時若受限於 scheduler hook 再明確記錄。
