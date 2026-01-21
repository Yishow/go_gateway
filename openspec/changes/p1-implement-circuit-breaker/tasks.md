# p1-implement-circuit-breaker 任務清單 (TDD Edition)

## 階段 1：健康度追蹤器 (Health Profiling)

- [ ] **1.1** [RED] 建立 `internal/collection/health` 套件與測試。
  - 建立 `health/tracker_test.go`。
  - 撰寫 `TestTracker_SlidingWindow`：模擬輸入 10 次結果 (8成2敗)，斷言錯誤率為 20%。
  - 執行測試確認失敗 (或編譯錯誤)。

- [ ] **1.2** [GREEN] 實作 `HealthTracker` 核心邏輯。
  - 實作滑動視窗 (Sliding Window) 或 Ring Buffer。
  - 實作 `RecordSuccess`, `RecordFailure`, `ErrorRate` 方法。
  - 確保通過 1.1 的測試。

- [ ] **1.3** [REFACTOR] 優化併發安全性。
  - 增加 `TestTracker_Concurrency`：模擬多個 Goroutine 同時寫入狀態。
  - 加入 `sync.Mutex` 或 atomic 操作確保線程安全。

## 階段 2：熔斷狀態機 (Circuit Breaker State Machine)

- [ ] **2.1** [RED] 定義狀態機行為測試。
  - 建立 `health/breaker_test.go`。
  - 撰寫 `TestBreaker_Trip_To_Dead`：設定閾值 20%，輸入錯誤直到觸發熔斷，斷言狀態轉為 `Dead`。
  - 撰寫 `TestBreaker_CoolDown`：斷言在 `Dead` 狀態下，未過冷卻時間前 `AllowRequest()` 回傳 false。

- [ ] **2.2** [GREEN] 實作 `CircuitBreaker` 結構與狀態流轉。
  - 定義 `State` Enum (Healthy, Unstable, Dead)。
  - 實作 `AllowRequest()`, `ReportResult()`。
  - 通過 2.1 所有測試。

- [ ] **2.3** [RED/GREEN] 實作半開 (Half-Open) 探測邏輯。
  - 測試：過了冷卻時間後，`AllowRequest()` 應回傳 true (一次)，且狀態暫轉為 `Unstable`。
  - 實作邏輯並通過測試。

## 階段 3：調度器整合 (Integration)

- [ ] **3.1** [RED] 建立集成測試 `collection/scheduler_test.go`。
  - 使用 `MockConnector` 模擬一個總是超時的設備。
  - 測試：調度器應在 N 次失敗後，停止呼叫 `MockConnector.Read()`。

- [ ] **3.2** [GREEN] 修改 `Scheduler` 整合 `CircuitBreaker`。
  - 在 `worker` 迴圈中加入 `breaker.AllowRequest()` 檢查。
  - 在 `read` 返回後呼叫 `breaker.ReportResult(err)`。

## 階段 4：背景探測與 UI (Probe & UI)

- [ ] **4.1** [TDD] 實作 `BackgroundProber`。
  - 測試：模擬 `Dead` 設備，Prober 應定期呼叫 `Connect()`。
  - 實作 Prober Worker。

- [ ] **4.2** API 與 UI 更新。
  - 暴露 `BreakerState` 至 `Device` 模型。
  - 更新 Dashboard 顯示邏輯。
