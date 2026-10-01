前置：A、B驗證完成；必須先讀現有delivery SQL存儲和worker，重用其可用部分。

所有項目是尚未執行的產品實作。每個行為項先建立具名測試並確認修正前失敗（RED），再實作使其通過（GREEN），最後refactor並重跑；高複雜度狀態機／並行／交易不可跳過此順序。下列測試名是預定驗收標籤，不宣稱已存在。每項保持單次session可驗證；若發現過大，先拆出可獨立驗收的小批並更新依賴，不以未驗收的半成品勾選。

## 1. 耐久接受與transaction

- [ ] 1.1 以DurableSampleAckCrash測試補足SQLJournal交易介面並在sample ACK前commit；注入commit失敗及ACK後crash，確認失敗不ACK、成功可recover。
- [ ] 1.2 [after: 1.1] 以AtomicRowOutboxCheckpoint在row closure同交易保存payload、frozen revisions、outbox及checkpoint；crash每個邊界後可replay且不重用sample，silent bucket outcome也持久化。
- [ ] 1.3 [after: 1.2] 以SQLiteTargetReceiptIdentity實作SQLite sender的row+effect receipt同target transaction與digest核對；commit response lost/restart只一份effect，local receipt失敗不盲目重送。
- [ ] 1.4 [after: 1.3] 以PostgresTargetReceiptIdentity完成PostgreSQL同等策略及actual DB integration fixture；custom無dedupe的unknown回覆進blocked而非重插。

## 2. 有界恢復與隔離

- [ ] 2.1 [after: 1.4] 以BoundedRetryAndPoisonPartition驗證retry/backoff/batch limits、retry耗盡保留資料、poison quarantine及同group/entity順序阻擋；健康partition不中斷。
- [ ] 2.2 [after: 2.1] 以QuotaRejectsNewAck驗證warning/hard limit與disk-full故障，保存exact affected scope及loss-risk，不能刪accepted pending資料騰空間。
- [ ] 2.3 [after: 2.2] 以WorkerFencingAndShutdown實作durable exclusive claim/lease、restart stale worker防護及bounded shutdown；deadline後未交付intent仍可恢復。

## 3. Production接線與truth

- [ ] 3.1 [after: 2.3] 在cmd/test_ui/service_wiring.go/main.go/target_writer.go接group pipeline及worker lifecycle；以ProductionGroupOutageRecovery測試確認新group只一writer，Share及另一DB在故障期間持續。
- [ ] 3.2 [after: 3.1] 以RevisionBoundBacklogAndDeliveryStages同步runtime/API/SSE/types；buffered不等於SQL committed，endpoint edit/disable不改送old backlog，credential失效scoped blocked。
- [ ] 3.3 [after: 3.2] 執行實際SQLite/Postgres故障及重啟測試、Go完整基準與受影響前端檢查；記錄sender dedupe能力限制與roll-back資料保護，再移交D，不宣稱universal exactly-once。
