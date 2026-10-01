前置：無新案前置；先核對既有verified-schema/result-truthfulness與workspace安全契約。

所有項目是尚未執行的產品實作。每個行為項先建立具名測試並確認修正前失敗（RED），再實作使其通過（GREEN），最後refactor並重跑；高複雜度狀態機／並行／交易不可跳過此順序。下列測試名是預定驗收標籤，不宣稱已存在。每項保持單次session可驗證；若發現過大，先拆出可獨立驗收的小批並更新依賴，不以未驗收的半成品勾選。

## 1. 擴充可保存群組

- [ ] 1.1 新增BasicGroupWithoutMeasurement測試，沿internal/datalink/workspace/及schema/migrations擴充WriteGroup欄位與repository；保存真實device/point/tag IDs、optional measurement及destination/group revision，重載完全一致，foreign/missing member不落資料。
- [ ] 1.2 [after: 1.1] 新增AtomicGroupSaveAndCAS，實作workspace write-groups API與同交易projection/CAS；注入第二段保存失敗、兩client同revision與unknown/foreign資源，確認全回滾/409/安全404，更新API/types。
- [ ] 1.3 [after: 1.2] 新增BasicGroupReadiness，讓group的managed/custom模式共用persisted authority；未選advanced semantics不要求MeasurementDefinition，保存draft不改applied revision或建表。
- [ ] 1.4 [after: 1.3] 以GroupLifecycleWithBacklog測試rename穩定ID、增刪member/改policy先draft且Apply於下一bucket切換、disable停新intake、delete tombstone仍可查並交付accepted payload/receipt；無法保護歸屬時拒delete，不cascade/purge。

## 2. 分批相容移轉

- [ ] 2.1 [after: 1.4] 以LegacySingleMappingMigration fixture完成第一批single-point mapping dry-run/review/apply與stable ID map；重跑不重複，舊read API反映同canonical資料。
- [ ] 2.2 [after: 2.1] 以LegacyRowGroupMigration完成第二批row-group migration；保留member/key metadata及共享column的row identity，會合併覆蓋或缺識別的情況blocked，原資料不動。
- [ ] 2.3 [after: 2.2] 以BasicPlanMigration完成第三批可等價basic recording plans；advanced stream/多target/不明measurement不猜測轉換，保存完整provenance與修復狀態。

## 3. 收斂authority與交接

- [ ] 3.1 [after: 2.3] 新增LegacyWriteAdapterConflict，將舊API寫入轉成canonical transaction或actionable409；證明不會產生第二份可寫設定，舊clientreload仍一致。
- [ ] 3.2 [after: 3.1] 新增WriterOwnershipActivationBarrier，準備revision-checked原子owner切換與rollback契約；C接線前不啟用新runtime，Share safety/token與已有backlog保留。
- [ ] 3.3 [after: 3.2] 執行focused API/repository/migration測試及AGENTS後端基準，驗證每批before/after payload；同步schema/API文件與migration備份/rollback說明，標明C/E/F尚未完成，不archive產品。
