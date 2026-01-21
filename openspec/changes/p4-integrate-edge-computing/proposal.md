# p4-integrate-edge-computing

## 摘要

本提案建議整合 HslCommunication 中的 `Algorithms` 模組概念，為 Gateway 引入**邊緣運算 (Edge Computing)** 能力，包括 FFT (快速傅立葉變換)、PID 模擬與數位過濾器。

## 動機

隨著工業 4.0 的發展，單純的「數據透傳」(Pass-through) 已不足以滿足需求。雲端處理高頻 raw data (如振動感測器的 10kHz 訊號) 成本過高且延遲過長。
HslCommunication 文件中提供了豐富的演算法支援 (`Fourier`, `PID`, `Filter`)。將這些能力下放到 Gateway 端，可以：

1.  **降低頻寬成本**：僅上傳特徵值（如 FFT 後的峰值頻率），而非原始波形。
2.  **即時反應**：在本地進行 PID 運算或異常檢測，實現毫秒級控制響應（若硬體允許）。
3.  **數據清洗**：透過濾波器去除雜訊，提高數據品質。

## 範圍

- **模組**：`internal/edge` (新增), `lib/algorithms` (新增)
- **能力**：`edge-analytics` (新增)

## 預期效益

- 拓展 Gateway 應用場景至「狀態監測 (Condition Monitoring)」。
- 大幅減少上傳至 MQTT/Database 的資料量。
- 提供加值服務功能。

## 風險與考量

- **CPU 資源**：FFT 等運算對 Raspberry Pi 等邊緣設備負擔較大，需提供開啟/關閉開關。
- **複雜度**：需設計靈活的 Pipeline 機制來串接不同的運算單元。

## 測試策略 (TDD)

邊緣運算涉及複雜數學邏輯，**Math TDD** 是確保演算法正確性的唯一途徑：

1.  **演算法正確性測試**：
    - 針對 FFT，使用已知訊號（如純 DC、單頻正弦波）作為 Test Case，驗證頻譜輸出是否符合理論值。
    - 針對 PID，使用標準階躍響應 (Step Response) 驗證控制輸出的收斂性。
2.  **Pipeline 數據流測試**：
    - Mock 各個 Processor 節點，驗證數據在 Pipeline 中的傳遞、型別轉換與錯誤處理機制。
3.  **資源限制測試**：
    - 撰寫測試案例模擬「運算過載」場景，驗證系統是否能依照預期丟棄舊數據，而非造成記憶體溢出 (OOM)。
