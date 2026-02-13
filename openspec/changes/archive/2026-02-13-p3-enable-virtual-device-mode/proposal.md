# p3-enable-virtual-device-mode

## 摘要

本提案建議引入**虛擬設備 (Virtual Device)** 與**服務器模式 (Server Mode)**，使 `go-gateway` 具備模擬 PLC 行為的能力，用於離線開發測試、數據代理 (Proxy) 以及協議轉換 (Protocol Conversion) 場景。

## 動機

HslCommunication 的一大亮點是提供了豐富的 `Server` 類別（如 `ModbusTcpServer`）。引入此功能可解決以下痛點：

1.  **離線開發困難**：工程師在沒有實體 PLC 的環境下，難以驗證 Datalink 映射邏輯與 UI 顯示。
2.  **協議轉換需求**：現場有多台異質設備（Fatek, MC3E），但上位 SCADA 系統僅支援標準 Modbus TCP。我們需要將 `go-gateway` 變成一個「虛擬 Modbus PLC」，彙整底層數據供上位讀取。

## 範圍

- **模組**：`internal/protocol/server`, `internal/virtual`
- **能力**：`protocol-servers` (新增), `simulation-mode` (新增)

## 預期效益

- **加速開發**：開發者可隨時啟動虛擬設備進行測試。
- **架構靈活性**：支援「南向採集 + 北向服務」的標準 Gateway 架構。
- **降低整合成本**：解決上位系統不支援特定私有協議的問題。

## 風險與考量

- **記憶體管理**：模擬 PLC 需要在記憶體中維護虛擬暫存器 (Virtual Memory)，需小心記憶體洩漏。
- **併發控制**：需處理多個 Client 同時寫入虛擬暫存器的 Race Condition。

## 測試策略 (TDD)

針對虛擬設備的開發，TDD 重點在於**行為一致性 (Behavioral Consistency)**：

1.  **記憶體單元測試**：使用 TDD 驗證 `MemoryBank` 的 Big/Little Endian 轉換邏輯與真實 PLC 完全一致。
2.  **Server 集成測試 (Integration TDD)**：啟動 `ModbusTcpServer` 後，使用真實的 Modbus Client 程式庫進行自動化讀寫測試，確保用標準工具（如 ModScan）能正常連接。
3.  **壓力測試**：模擬 100 個 Client 同時連接並讀寫同一記憶體區塊，驗證線程安全性。
