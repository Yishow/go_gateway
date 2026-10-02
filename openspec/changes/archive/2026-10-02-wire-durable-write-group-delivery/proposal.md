## Why

delivery journal/outbox/worker 元件已存在，但 cmd/test_ui 現有 writer 仍直接使用記憶體 grouped buffer，flush 失敗後不保留同一 bucket 重試。需要把已定義的 durable 契約真正接到 production entrypoint，才能區分已接收與實際 SQL 提交，並在斷線、重啟與不確定回覆後安全恢復。

## What Changes

- 把具版本的 sample intake、bucket recovery、finalized row 與 outbox checkpoint 接入本地交易，成功 ACK 只在 durable commit 後發出。
- 重用 delivery 元件並補足 production sender、destination-side identity/receipt、worker ownership、bounded shutdown 與 startup recovery。
- 採每個目的地獨立 bounded retry；disk-full／quota blocking 與 poison row 不靜默丟棄，不因資料庫故障擋住不相關 Share。
- 接受後凍結目標與group revisions；endpoint變更、停用與schema改變不得把backlog暗中改送。UI/API分開local durable、queued/retrying、SQL committed、unknown和verified。

## Capabilities

### New Capabilities

- `runtime-write-group-delivery`: production entrypoint上的持久接收、交付、恢復與可觀察狀態

### Modified Capabilities

無。沿用既有相關契約，新增production範圍要求，不重寫其他能力。

## Impact

cmd/test_ui/service_wiring.go、main.go、target_writer.go、internal/datalink/delivery/、dbtarget/、runtime/、schema/migrations/。前置：unify-studio-v2-write-group-contract、enforce-write-group-sample-semantics 全部驗證。延伸既有 durable-recording-delivery，不另建平行queue；不承諾跨所有目的地exactly-once。

本次僅起草文件，沒有產品實作。來源、依賴與移交見 [總覽](../../../docs/plans/studio-v2-write-groups/README.md) 及 [現況證據](../../../docs/plans/studio-v2-write-groups/evidence.md)。
