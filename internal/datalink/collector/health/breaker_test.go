package health

import (
	"testing"
	"time"
)

// =============================================================================
// 任務 2.1: [RED] 定義狀態機行為測試
// =============================================================================

func TestBreaker_InitialState(t *testing.T) {
	breaker := NewCircuitBreaker(BreakerConfig{
		ErrorThreshold:    0.2, // 20% 錯誤率觸發
		WindowSize:        10,
		UnstableThreshold: 0.5, // 50% 錯誤率轉為 Dead
		CooldownPeriod:    30 * time.Second,
	})

	if breaker.State() != StateHealthy {
		t.Errorf("初始狀態應為 Healthy，實際為 %s", breaker.State())
	}

	if !breaker.AllowRequest() {
		t.Error("Healthy 狀態應允許請求")
	}
}

func TestBreaker_Trip_To_Unstable(t *testing.T) {
	breaker := NewCircuitBreaker(BreakerConfig{
		ErrorThreshold:    0.2, // 20% 錯誤率觸發 Unstable
		WindowSize:        10,
		UnstableThreshold: 0.5,
		CooldownPeriod:    30 * time.Second,
	})

	// 輸入 10 次結果，3 次失敗 (30% > 20%)
	for i := 0; i < 7; i++ {
		breaker.ReportResult(nil)
	}
	for i := 0; i < 3; i++ {
		breaker.ReportResult(errTest)
	}

	if breaker.State() != StateUnstable {
		t.Errorf("應轉為 Unstable 狀態，實際為 %s", breaker.State())
	}
}

func TestBreaker_Trip_To_Dead(t *testing.T) {
	breaker := NewCircuitBreaker(BreakerConfig{
		ErrorThreshold:    0.2, // 20% 錯誤率觸發 Unstable
		WindowSize:        10,
		UnstableThreshold: 0.5, // 50% 錯誤率觸發 Dead
		CooldownPeriod:    30 * time.Second,
	})

	// 輸入 10 次失敗 (100% > 50%)
	for i := 0; i < 10; i++ {
		breaker.ReportResult(errTest)
	}

	if breaker.State() != StateDead {
		t.Errorf("應轉為 Dead 狀態，實際為 %s", breaker.State())
	}
}

func TestBreaker_Dead_Blocks_Requests(t *testing.T) {
	breaker := NewCircuitBreaker(BreakerConfig{
		ErrorThreshold:    0.2,
		WindowSize:        10,
		UnstableThreshold: 0.5,
		CooldownPeriod:    1 * time.Hour, // 長冷卻時間
	})

	// 觸發 Dead 狀態
	for i := 0; i < 10; i++ {
		breaker.ReportResult(errTest)
	}

	// Dead 狀態下應拒絕請求
	if breaker.AllowRequest() {
		t.Error("Dead 狀態應拒絕請求")
	}
}

func TestBreaker_CoolDown(t *testing.T) {
	breaker := NewCircuitBreaker(BreakerConfig{
		ErrorThreshold:    0.2,
		WindowSize:        10,
		UnstableThreshold: 0.5,
		CooldownPeriod:    100 * time.Millisecond, // 短冷卻時間
	})

	// 觸發 Dead 狀態
	for i := 0; i < 10; i++ {
		breaker.ReportResult(errTest)
	}

	// 等待冷卻時間
	time.Sleep(150 * time.Millisecond)

	// 冷卻後應允許一次探測請求
	if !breaker.AllowRequest() {
		t.Error("冷卻後應允許探測請求")
	}

	// 允許後狀態應轉為 Unstable (Half-Open)
	if breaker.State() != StateUnstable {
		t.Errorf("探測後應轉為 Unstable 狀態，實際為 %s", breaker.State())
	}
}

// =============================================================================
// 任務 2.3: [RED/GREEN] 半開探測邏輯測試
// =============================================================================

func TestBreaker_HalfOpen_Recovery(t *testing.T) {
	breaker := NewCircuitBreaker(BreakerConfig{
		ErrorThreshold:    0.2,
		WindowSize:        10,
		UnstableThreshold: 0.5,
		CooldownPeriod:    100 * time.Millisecond,
	})

	// 觸發 Dead 狀態
	for i := 0; i < 10; i++ {
		breaker.ReportResult(errTest)
	}

	// 等待冷卻
	time.Sleep(150 * time.Millisecond)

	// 允許探測
	breaker.AllowRequest()

	// 探測成功
	breaker.ReportResult(nil)

	// 應恢復為 Healthy
	if breaker.State() != StateHealthy {
		t.Errorf("探測成功後應恢復為 Healthy，實際為 %s", breaker.State())
	}
}

func TestBreaker_HalfOpen_Failure(t *testing.T) {
	breaker := NewCircuitBreaker(BreakerConfig{
		ErrorThreshold:    0.2,
		WindowSize:        10,
		UnstableThreshold: 0.5,
		CooldownPeriod:    100 * time.Millisecond,
	})

	// 觸發 Dead 狀態
	for i := 0; i < 10; i++ {
		breaker.ReportResult(errTest)
	}

	// 等待冷卻
	time.Sleep(150 * time.Millisecond)

	// 允許探測
	breaker.AllowRequest()

	// 探測失敗
	breaker.ReportResult(errTest)

	// 應回到 Dead
	if breaker.State() != StateDead {
		t.Errorf("探測失敗後應回到 Dead，實際為 %s", breaker.State())
	}
}

func TestBreaker_Unstable_Recovery(t *testing.T) {
	breaker := NewCircuitBreaker(BreakerConfig{
		ErrorThreshold:    0.2, // 20%
		WindowSize:        10,
		UnstableThreshold: 0.5,
		CooldownPeriod:    30 * time.Second,
	})

	// 觸發 Unstable 狀態 (30%)
	for i := 0; i < 7; i++ {
		breaker.ReportResult(nil)
	}
	for i := 0; i < 3; i++ {
		breaker.ReportResult(errTest)
	}

	// 連續成功請求應恢復
	for i := 0; i < 10; i++ {
		breaker.ReportResult(nil)
	}

	if breaker.State() != StateHealthy {
		t.Errorf("連續成功後應恢復為 Healthy，實際為 %s", breaker.State())
	}
}

// 測試用錯誤
var errTest = &testError{}

type testError struct{}

func (e *testError) Error() string { return "test error" }
