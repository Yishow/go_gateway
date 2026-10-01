前置：A/B/C全部驗證；保留既有schema ledger及501，直到本案end-to-end service測試通過。

所有項目是尚未執行的產品實作。每個行為項先建立具名測試並確認修正前失敗（RED），再實作使其通過（GREEN），最後refactor並重跑；高複雜度狀態機／並行／交易不可跳過此順序。下列測試名是預定驗收標籤，不宣稱已存在。每項保持單次session可驗證；若發現過大，先拆出可獨立驗收的小批並更新依賴，不以未驗收的半成品勾選。

## 1. 沿用確認與operation

- [ ] 1.1 建立TestWritePreviewIsReadOnlyAndScoped測試，擴充既有operation kind/payload；group preview及legacy plan adapter只解析scope/revisions，不作target mutation，錯誤kind/foreign/unknown/legacy裸plan_id安全拒絕。
- [ ] 1.2 [after: 1.1] 以AtomicTestWriteClaim實作confirmation digest/expiry/CAS與same-scope互斥；同operationrunning202、saved200、stale首次claim409，expired已完成token不重寫。

## 2. 實際寫入與清理

- [ ] 2.1 [after: 1.2] 以IdempotentProductionTestWrite讓SQLite及Postgres經C的codec/sender產生operation-owned test row；queued不能顯示written，lost response/restart查同receipt無duplicate。
- [ ] 2.2 [after: 2.1] 以TypedReadbackEvidence比較真實type/value/scope/identity/provenance；缺SELECT和payload mismatch均written_unverified，write不確定unknown，不自行填table/time。
- [ ] 2.3 [after: 2.2] 以OwnedCleanupAndRestart驗證只刪本operation rows/metadata，production neighbor不變；缺DELETE/cleanup timeout獨立status，cleanup後retry/restart不再寫入。

## 3. 完整契約開放

- [ ] 3.1 [after: 2.3] 更新Go/TypeScript/API/Swagger、service/hooks及capability；unsupported ownership target維持封鎖，schema token不能跨kind；完整測試前不移除501。
- [ ] 3.2 [after: 3.1] 執行API/operation fault matrix與實際DB測試，重跑AGENTS相關完整基準；保存write/readback/cleanup各別證據，接手舊3.3/3.4需求後交E整合UI。
