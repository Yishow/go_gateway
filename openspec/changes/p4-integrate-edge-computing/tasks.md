# p4-integrate-edge-computing 任務清單 (TDD Edition)

## 階段 1：演算法庫 (Math TDD)

- [x] **1.1** [RED] 建立 FFT 演算法測試。 ✅
  - 建立 `lib/algorithms/fft_test.go`。
  - **Case A (DC Offset)**: 輸入常數 `[1, 1, 1, 1]` -> 輸出頻譜在 0Hz 為 4，其餘為 0。
  - **Case B (Sine Wave)**: 輸入 10Hz 正弦波 -> 輸出在對應 Bin 有峰值。
- [x] **1.2** [GREEN] 實作 FFT 演算法。 ✅
  - 實作 Cooley-Tukey 算法。
  - 支援 IFFT 逆變換。
  - 通過 1.1 測試。

- [x] **1.3** [RED/GREEN] 實作 PID 邏輯。 ✅
  - 測試：設定 Target=100, Input=0 -> Output 應為正值 (P作用)。
  - 測試：設定積分項，長時間誤差應累積 Output。
  - 實作 Anti-windup 機制。

## 階段 2：Pipeline 架構 (Data Flow TDD)

- [x] **2.1** [RED] 定義 Pipeline 處理測試。 ✅
  - 建立 `internal/edge/pipeline_test.go`。
  - 定義 Mock Processor (將輸入 +1)。
  - 測試：Pipeline([1, 2, 3]) -> Result [2, 3, 4]。

- [x] **2.2** [GREEN] 實作 `PipelineRunner`。 ✅
  - 實作介面 `type Processor interface`。
  - 支援多處理器串聯。
  - 通過 2.1 測試。

## 階段 3：整合與配置

- [x] **3.1** [RED] 測試配置解析。 ✅
  - 輸入 YAML: `processors: [{type: fft}]`。
  - 斷言：解析出正確的 Processor Chain。

- [x] **3.2** [GREEN] 實作配置解析邏輯。 ✅
  - 使用 gopkg.in/yaml.v3。
  - 實作 ProcessorRegistry 工廠模式。

## 階段 4：系統保護 (System Protection)

- [x] **4.1** [RED] 測試超時丟棄機制。 ✅
  - Mock 一個耗時 1s 的 Processor，但 Pipeline 設定 Timeout=100ms。
  - 測試：輸入數據後，應回傳 Error 或被丟棄，且不阻塞下一個輸入。
- [x] **4.2** [GREEN] 實作 Context Timeout 與 Worker Pool 管理。 ✅
  - 使用 context.WithTimeout。
  - 快速返回不阻塞。

## 實作摘要

| 模組 | 檔案 | 說明 |
|------|------|------|
| FFT | `lib/algorithms/fft.go` | Cooley-Tukey FFT/IFFT 演算法 |
| PID | `lib/algorithms/pid.go` | PID 控制器 (含 Anti-windup) |
| Pipeline | `internal/edge/pipeline.go` | 資料處理管道 |
| Processors | `internal/edge/processors.go` | 內建處理器 (Scale, Filter, Threshold, Deadband, FFT, PID) |

## 測試覆蓋

- FFT: 4 個測試
- PID: 5 個測試
- Pipeline: 7 個測試

總計: 16 個測試全部通過
