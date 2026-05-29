## 1. 契約測試先行

- [ ] [P] 1.1 以 Red-Green 方式鎖定 `Runtime status snapshot API` 行為：在 `internal/api/router_runtime_test.go` 或等價契約測試中先寫出 top-level `metrics` 與 `device_id` 過濾斷言，完成後可觀察到 `GET /api/v1/datalink/runtime/status?device_id=...` 回傳 runtime metrics 且 `collectors` 只含指定 device；驗證：執行 `go test ./internal/api/... -run 'TestNewRouter_RuntimeStatus'`。
- [ ] [P] 1.2 以 Red-Green 方式鎖定 `Runtime value stream API`、`Device-scoped runtime monitoring contract` 與 `Decision: Keep runtime monitoring device-scoped`：新增或擴充 `internal/api/handlers` 契約測試，先證明缺少 `device_id` 會回 400、帶 `point_ids` 時只接收指定 points 的 `value` 事件、且同一連線仍可收到 `status` / `heartbeat`；驗證：執行 `go test ./internal/api/handlers/... -run 'TestRuntimeStreamHandler'`。
- [ ] [P] 1.3 以 Red-Green 方式鎖定 `Live monitoring stream remains loss-tolerant` 與 `Decision: Keep the stream loss-tolerant instead of blocking runtime collection`：在 `internal/datalink/runtime` 測試中先模擬 subscriber buffer 飽和，完成後可觀察到採集流程持續前進且其他 subscriber 不被拖慢；驗證：執行 `go test ./internal/datalink/runtime/... -run 'Test.*LossTolerant|Test.*Subscriber'`。

## 2. Snapshot 契約實作

- [ ] 2.1 落實 `Decision: Keep runtime monitoring device-scoped` 與 `Runtime status snapshot API`：更新 runtime status handler，讓 `GET /api/v1/datalink/runtime/status` 在 additive 相容前提下回傳 `metrics`、保留既有 collector 欄位、並正確套用 `device_id` filter；驗證：重新執行 `go test ./internal/api/... -run 'TestNewRouter_RuntimeStatus'`，並以 `curl /api/v1/datalink/runtime/status?device_id=<id>` 人工檢查 payload。
- [ ] 2.2 落實 `Decision: Reuse runtime service and scheduler health as the source of truth`：將 collector status、stale/error 計數、`last_read_at`、`last_error`、`breaker_state` 的推導集中為單一來源，完成後 snapshot 不再依賴前端自行推算；驗證：執行 `go test ./internal/api/... ./internal/datalink/runtime/... -run 'Test.*RuntimeStatus|Test.*Breaker'`。

## 3. Stream 對齊實作

- [ ] 3.1 落實 `Runtime value stream API` 與 `Keep runtime monitoring device-scoped`：更新 runtime stream handler 與訂閱介面，讓 `value`、`status`、`heartbeat` 三種事件都遵守 `device_id` 必填與 `point_ids` 過濾契約；驗證：執行 `go test ./internal/api/handlers/... -run 'TestRuntimeStreamHandler'`，並用 `curl -N /api/v1/datalink/runtime/stream?device_id=<id>` 檢查 SSE framing。
- [ ] 3.2 落實 Decision: Align runtime snapshot and stream under one derived status model，並交付 `Snapshot and stream use one derived device status model`：在 runtime service 中於 stale、point error、breaker 狀態變化時發出與 snapshot 同欄位的 `status` 事件，完成後前端不需為 snapshot 與 stream 維護兩套狀態推導；驗證：執行 `go test ./internal/datalink/runtime/... -run 'Test.*StatusEvent|Test.*AlignedStatus'`。

## 4. 收斂與驗證

- [ ] 4.1 完成 `datalink-api` 與 `runtime-dashboard-backend` 兩份 spec 的實作收斂：執行 `spectra analyze add-runtime-dashboard-backend-contract --json`、`spectra validate add-runtime-dashboard-backend-contract`、`go test ./internal/api/... ./internal/datalink/runtime/...`，確認無 Critical/Warning、spec 條文與測試證據一致，再交給後續 dashboard route change 使用。
