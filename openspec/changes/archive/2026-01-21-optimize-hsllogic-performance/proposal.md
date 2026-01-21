# optimize-hsllogic-performance

## 摘要

此變更旨在優化 `lib/hsllogic` 套件及協議連接器的效能，主要基於 HslCommunication 文件中識別出的最佳實踐。

## 動機

經過對 HslCommunication API 文件的分析，發現以下核心效能瓶頸與改進機會：

1.  **缺少長連接模式**：目前的協議連接器缺乏類似 HSL `NetworkDoubleBase` 的雙模式（短連接/長連接）支援，導致高頻輪詢時握手開銷過大。
2.  **PacketLogger 物件分配**：`PacketLogger` 為每個封包分配新的 `PacketLog` 結構，在高吞吐量場景下會造成 GC 壓力。
3.  **位址解析重複**：每次輪詢都重新解析位址字串（如 `D100`），浪費 CPU 資源。
4.  **缺少連接池**：多個 Goroutine 存取同一設備時，沒有連接池機制來重用 Socket。

## 範圍

- **模組**：`lib/hsllogic`, `internal/protocol/*`
- **能力**：`protocol-connectors`（修改）

## 預期效益

- 降低高頻採集時的 CPU 使用率與記憶體分配。
- 減少網路連線握手延遲。
- 提升系統在低資源硬體（Raspberry Pi）上的穩定性。

## 風險與考量

- 長連接模式需要增加 Heartbeat/KeepAlive 機制，增加複雜度。
- `sync.Pool` 的使用需謹慎處理物件狀態重置。

## 相關連結

- 優化分析報告：[optimization_analysis.md](../../docs/optimization_analysis.md)
- HSL 文件：`docs/HslDocs_TW/Introduction.md`
