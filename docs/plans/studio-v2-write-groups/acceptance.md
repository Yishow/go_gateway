# 預定驗收矩陣（尚未執行）

以下數字是deterministic設計fixture，不是production規模、速度或穩定性測量。每案例保存command、source/build SHA、平台、run IDs、實際assertions與sanitized SQL/UI witness。

基礎fixture：兩台loopback simulator設備A/B可同時有40001；group G含temperature decimal、pressure integer、running bool、batch text及uint64 9007199254740993。測試bucket interval=10秒、lateness=0、max_age=10秒，固定clock `2026-01-01T00:00:00Z` 起；production defaults不由此推斷。

| 標籤／owner | 操作／故障 | 必須觀察到的結果 |
| --- | --- | --- |
| BasicGroupWithoutMeasurement / A | UI保存persisted Tags，無physical semantics／report／retention | basic group可保存；ID均真實，無fake measurement |
| AtomicGroupSaveAndCAS / A | 第二段save失敗；兩client同revision寫入 | 全rollback；最多一方commit，另一方409且草稿保留 |
| LegacyRowGroupMigration / A | shared-column、missing key、advanced stream、重跑migration | 只有等價layout遷移；不合者blocked且原資料不動；ID穩定不雙寫 |
| GroupLifecycleWithBacklog / A | rename、增刪member、disable、delete且有backlog | rename ID不變；semantic Apply下一bucket；tombstone保留歸屬／payload／receipts，可查及完成交付、不purge |
| OfflineDraftNavigation / E | offline或無probe，Save/Next/Back/reload | 可以保存/導航草稿，明示unverified，僅相關activation阻擋 |
| UTCWindowBoundaryAndReordering / B | observed t=8先到、t=3後到、t=10邊界；同timestamp不同sample ID | t=8勝出；t=10進下一bucket；tie-break可重現 |
| ExactMixedValueRoundTrip / B | decimal/bool/text/大整數；NaN/Inf/overflow | digits和type不變；非法值不標good、不默默變0 |
| FreshnessMissingAndSilentBucket / B | 缺一值、bad、max_age=2的stale樣本、全bucket靜默 | default零SQL row；bucket有scoped skipped/no_data原因；不補前bucket值 |
| ExplicitPartialPolicy / B | 允許NULL及quality storage後明確partial | NULL＋reason；能力不足在activation前阻擋；完整good custom row不強迫外部metadata |
| ScopedRowIdentityAndLateArrival / B | 兩group/entity同timestamp；closed後late | 兩個不同row key；late不改closed row；同sample duplicate無效應 |
| AtomicRowOutboxCheckpoint / C | ACK後kill；row/outbox/checkpoint交易中kill | ACK資料可recover；無消費checkpoint但缺outbox的裂縫 |
| TargetReceiptIdentity / C | SQLite與Postgres各自commit後丟response；local receipt保存失敗 | 相同effect key只有一個target effect；unsupported dedupe則unknown而非blind insert |
| ProductionGroupOutageRecovery / C | A目標offline跨gateway restart，B目標/Share正常 | Abacklog持久恢復；B/Share繼續；SQLcommitted只在DB證據後出現 |
| QuotaRejectsNewAck / C | fixture低quota、disk write/commit失敗 | 無假ACK，暴露scope/loss-risk，已接受資料不丟 |
| BoundedRetryAndPoisonPartition / C | 永久SQL型別錯誤、retry上限、worker重疊 | poison保留；同group順序阻擋，其他partition持續；exclusive claim有效 |
| RevisionBoundBacklog / C | endpoint/scale/group變更、disable | backlog不改送、不重算；新intake依new applied revision |
| AtomicTestWriteClaim / D | preview、wrong-kind/stale/expired/foreign、rapid duplicate | preview零target mutation；安全拒絕；same-op202或saved200，無duplicate |
| TypedReadbackEvidence / D | row exists但value mismatch；缺SELECT | written_unverified含原因，不只ID相符就verified |
| OwnedCleanupAndRestart / D | neighbor production row、缺DELETE、cleanup中kill後retry | 只處理test-owned row；cleanup status獨立；receipt保留不重寫 |
| CapabilityAndReloadReadiness / E | MQTT能力缺、connection edit、reload/back/cancel | option明示不可用；old probe失效；initial/reload gate一致 |
| DeviceScopedAddressConflict / E | A/B都40001；同device normalized range真重疊 | 跨device合法；同device衝突定位正確 |
| EmptyErrorMismatchPlanRecovery / E | load fail/empty/all mismatched | retry/新建/明確修復可達；不把error當empty |
| GroupCompletionTruth / E | managed無manual targets；Share-only；device running但DBqueued | 按backend group gate；Share保護獨立；不假DBsuccess |
| GuidedGroupKeyboardAndViewport / E | IME、Enter+blur、double click、filter+bulk、390/768/1440px | 不重送、不丟draft、明確bulk scope、焦點及修復action可達 |
| DeviceToSQLiteMixedRows / F | 真UI建立兩device/points/tags/groups，simulator採集 | 內部persisted IDs與獨立SQL查詢同鏈；不能direct INSERT替代 |
| DeviceToPostgresMixedRows / F | 同完整鏈於disposablePostgres | 真Postgres證據；不可用時blocked，不能類推pass |

## 驗證層級

1. A–E各案先有focused RED/GREEN與受影響完整基準
2. F用actual production binary＋embedded UI＋real disposable DB，回歸故障矩陣
3. Windows／ARM、外部/embedded browser、LAN、真PLC、SCADA、長時間soak與deployment-owner rollback/redeploy另列明確run或not-run；不能由Linux simulator推論
4. 不以本批文件檢查代替以上產品驗收
