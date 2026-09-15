package collector

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

func TestEmitValueDoesNotBlock(t *testing.T) {
	config := DefaultSchedulerConfig()
	config.ValueBufferSize = 1
	scheduler := NewScheduler(config, nil)

	scheduler.valueChan <- CollectedValue{}

	done := make(chan struct{})
	go func() {
		scheduler.emitValue(CollectedValue{})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("emitValue blocked with full buffer")
	}
}

// =============================================================================
// 任務 3.1: [RED] 建立集成測試
// =============================================================================

// MockConnector 模擬總是超時的連接器
type MockConnector struct {
	readCount  atomic.Int32
	shouldFail bool
}

func (m *MockConnector) Connect(ctx context.Context) error {
	return nil
}

func (m *MockConnector) Disconnect() error {
	return nil
}

func (m *MockConnector) Read(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	m.readCount.Add(1)
	if m.shouldFail {
		return connector.ReadResult{
			Quality: schema.QualityBad,
			Error:   "connection timeout",
		}, errors.New("connection timeout")
	}
	return connector.ReadResult{
		Value:     42,
		Quality:   schema.QualityGood,
		Timestamp: time.Now(),
	}, nil
}

func (m *MockConnector) Write(ctx context.Context, req connector.WriteRequest) error {
	return nil
}

func (m *MockConnector) IsConnected() bool {
	return true
}

func (m *MockConnector) Protocol() schema.ProtocolType {
	return schema.ProtocolModbusTCP
}

func TestScheduler_CircuitBreaker_StopsAfterFailures(t *testing.T) {
	// 建立熔斷器配置
	breakerConfig := health.BreakerConfig{
		ErrorThreshold:    0.2,           // 20% 錯誤率觸發
		WindowSize:        5,             // 5 次請求視窗
		UnstableThreshold: 0.5,           // 50% 觸發 Dead
		CooldownPeriod:    1 * time.Hour, // 長冷卻避免自動恢復
	}

	// 建立熔斷器
	breaker := health.NewCircuitBreaker(breakerConfig)

	// 模擬 5 次連續失敗
	for i := 0; i < 5; i++ {
		if breaker.AllowRequest() {
			breaker.ReportResult(errors.New("timeout"))
		}
	}

	// 斷言：熔斷器應處於 Dead 狀態
	if breaker.State() != health.StateDead {
		t.Errorf("預期狀態為 Dead，實際為 %s", breaker.State())
	}

	// 斷言：應拒絕後續請求
	if breaker.AllowRequest() {
		t.Error("Dead 狀態應拒絕請求")
	}
}

func TestScheduler_CircuitBreaker_RecoveryAfterSuccess(t *testing.T) {
	breakerConfig := health.BreakerConfig{
		ErrorThreshold:    0.2,
		WindowSize:        5,
		UnstableThreshold: 0.5,
		CooldownPeriod:    100 * time.Millisecond,
	}

	breaker := health.NewCircuitBreaker(breakerConfig)

	// 觸發 Dead 狀態
	for i := 0; i < 5; i++ {
		breaker.ReportResult(errors.New("timeout"))
	}

	// 等待冷卻
	time.Sleep(150 * time.Millisecond)

	// 嘗試探測
	if !breaker.AllowRequest() {
		t.Fatal("冷卻後應允許探測")
	}

	// 探測成功
	breaker.ReportResult(nil)

	// 應恢復 Healthy
	if breaker.State() != health.StateHealthy {
		t.Errorf("探測成功後應恢復 Healthy，實際為 %s", breaker.State())
	}
}
