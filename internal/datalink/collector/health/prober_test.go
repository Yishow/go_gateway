package health

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// =============================================================================
// 任務 4.1: [TDD] BackgroundProber 測試
// =============================================================================

func TestProber_ProbesDeadDevices(t *testing.T) {
	var probeCount atomic.Int32
	probeFunc := func(ctx context.Context, deviceID string) error {
		probeCount.Add(1)
		return nil
	}

	prober := NewBackgroundProber(ProberConfig{
		ProbeInterval: 50 * time.Millisecond,
		ProbeTimeout:  1 * time.Second,
	}, probeFunc)

	// 建立 Dead 狀態的熔斷器
	breaker := NewCircuitBreaker(BreakerConfig{
		ErrorThreshold:    0.2,
		WindowSize:        5,
		UnstableThreshold: 0.5,
		CooldownPeriod:    10 * time.Millisecond, // 短冷卻
	})

	// 觸發 Dead 狀態
	for i := 0; i < 5; i++ {
		breaker.ReportResult(errors.New("error"))
	}

	if breaker.State() != StateDead {
		t.Fatalf("預期 Dead 狀態，實際為 %s", breaker.State())
	}

	prober.RegisterDevice("dev-1", breaker)
	prober.Start()

	// 等待至少一次探測
	time.Sleep(100 * time.Millisecond)

	prober.Stop()

	if probeCount.Load() == 0 {
		t.Error("應該至少執行一次探測")
	}
}

func TestProber_SuccessfulProbeRecovery(t *testing.T) {
	probeFunc := func(ctx context.Context, deviceID string) error {
		return nil // 探測成功
	}

	prober := NewBackgroundProber(ProberConfig{
		ProbeInterval: 50 * time.Millisecond,
		ProbeTimeout:  1 * time.Second,
	}, probeFunc)

	// 建立 Dead 狀態的熔斷器
	breaker := NewCircuitBreaker(BreakerConfig{
		ErrorThreshold:    0.2,
		WindowSize:        5,
		UnstableThreshold: 0.5,
		CooldownPeriod:    10 * time.Millisecond,
	})

	// 觸發 Dead 狀態
	for i := 0; i < 5; i++ {
		breaker.ReportResult(errors.New("error"))
	}

	prober.RegisterDevice("dev-1", breaker)
	prober.Start()

	// 等待探測與恢復
	time.Sleep(150 * time.Millisecond)

	prober.Stop()

	// 探測成功後應恢復 Healthy
	if breaker.State() != StateHealthy {
		t.Errorf("探測成功後應恢復 Healthy，實際為 %s", breaker.State())
	}
}

func TestProber_FailedProbeStaysDead(t *testing.T) {
	probeFunc := func(ctx context.Context, deviceID string) error {
		return errors.New("still failing") // 探測失敗
	}

	prober := NewBackgroundProber(ProberConfig{
		ProbeInterval: 50 * time.Millisecond,
		ProbeTimeout:  1 * time.Second,
	}, probeFunc)

	// 建立 Dead 狀態的熔斷器
	breaker := NewCircuitBreaker(BreakerConfig{
		ErrorThreshold:    0.2,
		WindowSize:        5,
		UnstableThreshold: 0.5,
		CooldownPeriod:    10 * time.Millisecond,
	})

	// 觸發 Dead 狀態
	for i := 0; i < 5; i++ {
		breaker.ReportResult(errors.New("error"))
	}

	prober.RegisterDevice("dev-1", breaker)
	prober.Start()

	// 等待探測
	time.Sleep(150 * time.Millisecond)

	prober.Stop()

	// 探測失敗後應保持 Dead
	if breaker.State() != StateDead {
		t.Errorf("探測失敗後應保持 Dead，實際為 %s", breaker.State())
	}
}

func TestProber_GetDeadDevices(t *testing.T) {
	prober := NewBackgroundProber(DefaultProberConfig(), nil)

	// 建立兩個熔斷器
	healthyBreaker := NewCircuitBreaker(DefaultBreakerConfig())
	deadBreaker := NewCircuitBreaker(BreakerConfig{
		ErrorThreshold:    0.2,
		WindowSize:        5,
		UnstableThreshold: 0.5,
		CooldownPeriod:    1 * time.Hour,
	})

	// 觸發 Dead 狀態
	for i := 0; i < 5; i++ {
		deadBreaker.ReportResult(errors.New("error"))
	}

	prober.RegisterDevice("healthy-dev", healthyBreaker)
	prober.RegisterDevice("dead-dev", deadBreaker)

	deadDevices := prober.GetDeadDevices()

	if len(deadDevices) != 1 {
		t.Errorf("預期 1 個 Dead 設備，實際有 %d 個", len(deadDevices))
	}

	if len(deadDevices) > 0 && deadDevices[0] != "dead-dev" {
		t.Errorf("預期 dead-dev，實際為 %s", deadDevices[0])
	}
}
