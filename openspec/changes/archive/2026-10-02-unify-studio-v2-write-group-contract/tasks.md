前置：無新案前置；先核對既有verified-schema/result-truthfulness與workspace安全契約。

完成項目以勾選與驗證紀錄為準，其餘產品實作仍待執行。每個行為項先建立具名測試並確認修正前失敗（RED），再實作使其通過（GREEN），最後refactor並重跑；高複雜度狀態機／並行／交易不可跳過此順序。下列測試名是預定驗收標籤，不宣稱已存在。每項保持單次session可驗證；若發現過大，先拆出可獨立驗收的小批並更新依賴，不以未驗收的半成品勾選。

## 1. 擴充可保存群組

- [x] 1.1 新增BasicGroupWithoutMeasurement測試，沿internal/datalink/workspace/及schema/migrations擴充WriteGroup欄位與repository；保存真實device/point/tag IDs、optional measurement及destination/group revision，重載完全一致，foreign/missing member不落資料。
- [x] 1.2 [after: 1.1] 新增AtomicGroupSaveAndCAS，實作workspace write-groups API與同交易projection/CAS；注入第二段保存失敗、兩client同revision與unknown/foreign資源，確認全回滾/409/安全404，更新API/types。
- [x] 1.3 [after: 1.2] 新增BasicGroupReadiness，讓group的managed/custom模式共用persisted authority；未選advanced semantics不要求MeasurementDefinition，保存draft不改applied revision或建表。
- [x] 1.4 [after: 1.3] 以GroupLifecycleWithBacklog測試rename穩定ID、增刪member/改policy先draft且Apply於下一bucket切換、disable停新intake、delete tombstone仍可查並交付accepted payload/receipt；無法保護歸屬時拒delete，不cascade/purge。

## 2. 分批相容移轉

- [x] 2.1 [after: 1.4] 以LegacySingleMappingMigration fixture完成Reviewable repeatable compatibility migration與Reviewed single-mapping snapshot conversion：第一批single-point dry-run明列逐筆→snapshot差異，explicit confirmation及同revision/digest review保存draft與stable ID map；stale/blocked/foreign及交易失敗不落資料，重跑不重複、不覆寫canonical edit，舊read反映同canonical資料，保留原設定及舊writer到有效Apply。
- [x] 2.2 [after: 2.1] 以LegacyRowGroupMigration驗證Workspace database row groups capture shared-column intent並完成第二批row-group preview／confirmed review；保存entity_key、member/key metadata與stable provenance，舊read反映canonical member，會合併覆蓋或缺識別的情況blocked；驗證同revision/digest、重跑不覆寫及交易rollback，原target與writer不動。
- [x] 2.3 [after: 2.2] 以BasicPlanMigration完成Blocked recording-plan migration：完整原 intent/source revisions 預覽及修復提示，single every_sample 與 unsupported/多target/不明measurement 不猜測轉換；current review422、stale409、foreign/unknown404，原plan與所有group/map/projection/runtime不變，新基本寫入沿canonical Create；驗證完整重載、source/destination digest變動、無副作用與SDK bounded blocked契約。

## 3. 收斂authority與交接

- [x] 3.1 [after: 2.3] 新增LegacyWriteAdapterConflict，驗證One production writer owner per output：global/Studio target及owned row-group修改回actionable409，同transaction ownership guard涵蓋Create原/新Tag-connector scope、Update/Delete、explicit-empty replacement、rollback與並行移轉；未移轉CRUD及canonical Save/Review/reload保持一致，不切runtime。
- [x] 3.2 [after: 3.1] 新增WriterOwnershipActivationBarrier，以Apply同transaction注入owner transition及精確revision/effective_at；驗證stale／未ready不呼叫、barrier及後段workspace失敗全rollback、rename不重置與accepted payload/receipt保留；實跑Share settings/hydration/token/revision回歸，C前不新增Apply route或啟用新runtime。
- [x] 3.3 [after: 3.2] 回查One persisted write-group authority、Atomic revision-safe group saves、Safe basic group lifecycle與Open Questions已有決策／證據；執行focused API/repository/migration測試及AGENTS後端基準，驗證每批before/after payload；同步schema/API文件與migration備份/rollback說明，標明C/E/F尚未完成，不archive產品。

變更說明：A2.1原「single-point等價映射」前提已由design的Reviewed single-mapping snapshot conversion取代；已完成1.1–1.4描述與provenance原樣保留。

變更說明：使用者確認以 Blocked recording-plan migration 取代原 pending 2.3 的自動等價移轉；原 pending 描述不算完成。2.3 編號及 3.1 前置保留，已完成 1.1–2.2 描述與 provenance 不變。
