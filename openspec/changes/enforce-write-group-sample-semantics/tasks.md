前置：A unify-studio-v2-write-group-contract 的tasks與驗證通過後開始。

所有項目是尚未執行的產品實作。每個行為項先建立具名測試並確認修正前失敗（RED），再實作使其通過（GREEN），最後refactor並重跑；高複雜度狀態機／並行／交易不可跳過此順序。下列測試名是預定驗收標籤，不宣稱已存在。每項保持單次session可驗證；若發現過大，先拆出可獨立驗收的小批並更新依賴，不以未驗收的半成品勾選。

## 1. 傳遞真實sample

- [ ] 1.1 建立TypedAcquisitionEnvelope與GatewayTimeOrigin測試；沿collector/runtime/measurement types保留source IDs/revisions、真observed_at/received_at、quality及reason，不用flush時間補造觀察時間。
- [ ] 1.2 [after: 1.1] 以ExactMixedValueRoundTrip驗證bool/text/int64/uint64 9007199254740993/decimal跨JSON及SQL codec；NaN/Inf/overflow標bad或拒絕，不截斷、不以0補值。

## 2. 一次且可重現的資料列

- [ ] 2.1 [after: 1.2] 以UTCWindowBoundaryAndReordering實作pure snapshot selector；驗證半開10秒fixture、max observed_at、same-time sample_id tie-break、duplicate same/different payload與configured future skew。
- [ ] 2.2 [after: 2.1] 以FreshnessMissingAndSilentBucket實作time-driven closure和skip_row預設；缺值/bad/stale/整bucket無sample留下scoped no_data/skipped，零SQL row、無carry-forward。
- [ ] 2.3 [after: 2.2] 以ExplicitPartialPolicy驗證partial只有nullable/quality storage能力足夠才接受；missing為NULL+reason；全good custom row的provenance可留local durable envelope，不強迫外部metadata表。
- [ ] 2.4 [after: 2.3] 以ScopedRowIdentityAndLateArrival實作workspace/group/revision/entity/bucket identity與destination namespacing；同timestamp不同entity兩row，closed bucket拒late改寫，scale/revision改變不重解釋舊row。

## 3. 接入並驗證

- [ ] 3.1 [after: 2.4] 將typed sample與row result接到新group runtime boundary（internal/datalink/runtime/、dbtarget/）；對unsupported SQL type/identity回安全阻擋，legacy未移轉路徑保持原樣。
- [ ] 3.2 [after: 3.1] 跑focused sample/row/SQLcodec測試與AGENTS後端基準，保存fixture expected值與status；不把B的記憶體assembler稱為durable，移交journal/checkpoint給C。
