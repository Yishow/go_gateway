# p1-implement-circuit-breaker

## 摘要

本提案旨在引入**設備健康畫像 (Health Profiling)** 與**自適應熔斷機制 (Adaptive Circuit Breaking)**，以提升 Gateway 在部分設備故障時的系統整體韌性 (Resiliency)。

## 動機

在工業物聯網現場，設備斷線、電源故障或網路波動是常態。目前 `go-gateway` 若遇到設備超時，往往會持續嘗試連線，導致執行緒資源（Goroutines）被卡在 I/O 等待上，形成 Head-of-line blocking，進而拖慢其他正常設備的採集頻率。

我們需要一種自我保護機制，能夠：

1.  **快速失敗**：當設備被判定為「不健康」時，立即跳過後續請求。
2.  **降級服務**：自動將問題設備移至「慢速佇列」，降低對系統資源的消耗。
3.  **自動恢復**：在背景探測成功後，自動恢復正常採集頻率。

這符合 HslCommunication 文件中對於網路異常處理的最佳實踐建議。

## 範圍

- **模組**：`internal/collection`, `internal/device-registry`
- **能力**：`collection-scheduler`, `device-monitoring` (新增)

## 預期效益

- **系統穩定性**：單一設備故障不會影響整體 CPU/Memory 資源。
- **採集延遲保證**：確保健康設備的採集週期不被拖累。
- **可觀測性提升**：提供設備健康度指標（Availability, Latency, Error Rate）。

## 風險與考量

- 需調整調度器邏輯，增加狀態機複雜度。
- 需小心設定熔斷閾值，避免網路抖動造成頻繁的狀態切換（Flapping）。

## 測試策略 (TDD)

本變更採用 **TDD (Test-Driven Development)** 流程進行開發，確保狀態機邏輯的絕對正確性。

1.  **狀態機測試**：模擬 Mock Clock 與 Mock Network，驗證 Healthy -> Unstable -> Dead 的所有狀態轉換路徑。
2.  **併發測試**：使用 Go 的 `-race` detector 驗證多執行緒下的滑動視窗統計是否安全。
3.  **整合測試**：在 `Scheduler` 層級注入 Mock Connector，驗證高層調度是否正確響應底層的熔斷訊號。
