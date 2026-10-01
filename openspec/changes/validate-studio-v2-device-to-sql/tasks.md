前置：A/B/C/D/E全部驗證；依acceptance.md矩陣，執行前確認只有simulator和disposable DB。

所有項目是尚未執行的產品實作。每個行為項先建立具名測試並確認修正前失敗（RED），再實作使其通過（GREEN），最後refactor並重跑；高複雜度狀態機／並行／交易不可跳過此順序。下列測試名是預定驗收標籤，不宣稱已存在。每項保持單次session可驗證；若發現過大，先拆出可獨立驗收的小批並更新依賴，不以未驗收的半成品勾選。

## 1. 可重跑真實路徑

- [ ] 1.1 建立ProductionDeviceSQLHarness：fresh GATEWAY_DB_PATH、loopback Modbus simulator、actual cmd/test_ui+embedded frontend與Playwright；witness含source/build/commands/platform/IDs，禁止stub final API或direct INSERT代替採集。
- [ ] 1.2 [after: 1.1] 完成DeviceToSQLiteMixedRows，由UI設定兩device同addresses、多type點位/Tags/group，獨立SQL查values/identity/UTC/quality與persisted IDs，uint64及decimal精度一致。
- [ ] 1.3 [after: 1.2] 完成DeviceToPostgresMixedRows，對可丟棄PostgreSQL重跑同一條真UI→SQL鏈及schema confirmation，環境缺失明列blocked而非借SQLite通過。

## 2. 故障與重啟

- [ ] 2.1 [after: 1.3] 跑QualityBoundaryMatrix：out-of-order/duplicate/bad/stale/missing/entirely silent/late buckets，查default skip_row與partial不造值，UI原因和SQL一致。
- [ ] 2.2 [after: 2.1] 跑OutageCrashAndUnknownCommit：target disconnect、ACK後kill、commit回覆遺失、restart/recover；actual SQL effect key無duplicates，backlog不因endpoint change改送。
- [ ] 2.3 [after: 2.2] 跑CapacityPoisonAndConcurrentEdit：縮小quota、disk failure、poison partition、stale save/activation/test-write，檢查no false ACK、健康partition、CAS、scope隔離與held draft。
- [ ] 2.4 [after: 2.3] 跑ConfirmedTestWriteCleanup：actual preview/confirm/readback/cleanup、no SELECT/no DELETE、restart及same-operation retry；SQL neighbor rows不受影響、status無假成功。

## 3. 驗收與舊案closeout

- [ ] 3.1 [after: 2.4] 收集Linux實際UI/SQL witnesses，跑Go test/vet/lint與frontend lint/Vitest/build/Playwright及line/diff/OpenSpec strict檢查，逐項記pass/fail/not-run；Windows/ARM/embedded browser/LAN/真PLC/SCADA各自列未驗。
- [ ] 3.2 [after: 3.1] 獨立review所有acceptance與migration/rollback證據，修正scope內缺陷並重跑受影響場景；完成D/E/F移交驗收後才處理舊fix-studio-v2-database-workflow的closeout，未授權不部署。
