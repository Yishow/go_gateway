## Why

目前 dbtarget 的 grouped buffer 只按到達順序覆寫欄位，UTC bucket 並不代表各欄位同時、最新或品質良好。群組真正寫 SQL 前需要可測的觀察時間、來源、型別、完整性與資料列身分契約，避免把舊值當新值或不同設備互相覆蓋。

## What Changes

- 依唯一群組 revision 產生具來源與品質的 typed sample；連到現有 SampleEnvelope 可重用處，補齊 point/tag 身分，不強迫 raw Tag 擁有物理語意。
- 明確採 UTC 半開區間 snapshot，按 observation time 選值，預設不跨區間補舊值；缺少／bad／stale 成員依已存政策阻擋資料列或明確記 partial。
- 資料列 identity 包含 workspace、write-group、group revision、entity key 與 bucket start；同時間不同群組／設備仍是不同資料列。
- 保留 bool、text、精確整數與 decimal；SQL型別不相容及非有限值不得轉成正常值。late、duplicate、out-of-order 均有具名案例。

## Capabilities

### New Capabilities

- `write-group-row-semantics`: 實際採集到群組資料列的時間、品質、完整性、精確值與identity

### Modified Capabilities

無。沿用既有相關契約，新增production範圍要求，不重寫其他能力。

## Impact

internal/datalink/runtime/、collector/、measurement/types.go、dbtarget/writer.go 與 writer_statements.go，以及受影響 API／frontend types。前置：unify-studio-v2-write-group-contract 完成並驗證。只定義基本 snapshot，衍生統計、區間電量與歷史報表不在範圍。

本次已授權依 A→F 實作與驗證；完成狀態以 tasks 與 validation 為準。來源、依賴與移交見 [總覽](../../../docs/plans/studio-v2-write-groups/README.md) 及 [現況證據](../../../docs/plans/studio-v2-write-groups/evidence.md)。
