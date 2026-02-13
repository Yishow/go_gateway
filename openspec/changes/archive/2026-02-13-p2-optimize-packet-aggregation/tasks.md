# p2-optimize-packet-aggregation 任務清單 (TDD Edition)

## 階段 1：聚合演算法核心 (Optimizer Core)

- [x] **1.1** [RED] 定義 `BlockMerger` 基礎行為測試。 ✅
  - 建立 `internal/datalink/optimizer/merger_test.go`。
  - **Case A (連續)**: 輸入 `[D100, D101, D102]` -> 期望 `[Block(Start:D100, Len:3)]`。
  - **Case B (分離)**: 輸入 `[D100, D200]` -> 期望 `[Block(D100, 1), Block(D200, 1)]`。
  - 確認測試失敗。

- [x] **1.2** [GREEN] 實作基礎合併邏輯。 ✅
  - 實作 `MergePoints([]Point) []Block`。
  - 採用簡單的排序與貪婪合併算法。
  - 通過 1.1 測試。

- [x] **1.3** [RED] 定義 PDU 限制與空洞測試。 ✅
  - **Case C (空洞)**: 輸入 `[D100, D105]`, MaxHole=10 -> 期望 `[Block(D100, 6)]` (讀取 D100-D105)。
  - **Case D (PDU限制)**: 輸入 `[D100...D300]`, MaxLen=100 -> 期望拆分為多個 Blocks。
- [x] **1.4** [GREEN] 增強合併邏輯。 ✅
  - 加入 `MaxHole` 與 `MaxPDU` 參數處理。
  - 通過 1.3 測試。

## 階段 2：資料分發 (Response Dispatcher)

- [x] **2.1** [RED] 測試 Block 資料解析回填。 ✅
  - 建立 `dispatcher_test.go`。
  - 輸入：一個 `Block(D100, Len=3)` 的 Raw Bytes `[00 01 00 02 00 03]`。
  - 期望：
    - Point(D100) = 1
    - Point(D101) = 2
    - Point(D102) = 3

- [x] **2.2** [GREEN] 實作 `Dispatcher`。 ✅
  - 根據 Point 的 Offset 從 Block Raw Bytes 中提取數據。
  - 處理 ByteOrder (BigEndian/LittleEndian)。

## 階段 3：調度器整合 (Integration)

- [x] **3.1** 修改 `collection-scheduler` 啟動流程。 ✅ (待整合)
  - 注入 `Optimizer`。
  - 在 `Start()` 前執行靜態優化產生 `ExecutionPlan`。

- [x] **3.2** [Verify] 整合測試。 ✅
  - 模擬一個擁有 1000 個 Tags 的設備。
  - 紀錄 `Connector.Read()` 被呼叫的次數。
  - 斷言呼叫次數大幅減少 (例如從 1000 次降為 20 次)。

## 階段 4：效能基準測試 (Benchmark)

- [x] **4.1** 建立 Benchmark。 ✅
  - `BenchmarkOptimizer_1000Tags`：測量優化演算法本身的耗時（應在 ms 等級）。
  - 結果：1000 點位約 409µs，10000 點位約 11ms

## 實作摘要

| 元件 | 檔案 | 說明 |
|------|------|------|
| BlockMerger | `optimizer/merger.go` | 貪婪合併演算法，支援 MaxHoleSize 和 MaxBlockLength |
| Dispatcher | `optimizer/dispatcher.go` | 原始數據分發，支援 BigEndian/LittleEndian 和多種資料類型 |
| Benchmark | `optimizer/benchmark_test.go` | 效能基準測試 |

## 效能數據

```
BenchmarkMerger_100Points-22      33646    40402 ns/op    68976 B/op   275 allocs/op
BenchmarkMerger_1000Points-22      2767   409174 ns/op   756924 B/op  2419 allocs/op
BenchmarkMerger_10000Points-22      100 11332112 ns/op  8347129 B/op 23834 allocs/op
```
