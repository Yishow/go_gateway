package health

import (
	"sync"
	"time"
)

// =============================================================================
// 狀態定義
// =============================================================================

// State 熔斷器狀態
type State string

const (
	// StateHealthy 健康狀態 (Closed Circuit) - 正常運作
	StateHealthy State = "healthy"

	// StateUnstable 不穩定狀態 (Half-Open) - 偶發失敗，降頻採集
	StateUnstable State = "unstable"

	// StateDead 死亡狀態 (Open Circuit) - 熔斷，僅背景探測
	StateDead State = "dead"
)

// =============================================================================
// 熔斷器配置
// =============================================================================

// BreakerConfig 熔斷器配置
type BreakerConfig struct {
	// ErrorThreshold 錯誤率閾值，超過此值從 Healthy 轉為 Unstable
	ErrorThreshold float64

	// UnstableThreshold 不穩定閾值，超過此值從 Unstable 轉為 Dead
	UnstableThreshold float64

	// WindowSize 滑動視窗大小
	WindowSize int

	// CooldownPeriod 冷卻時間 (Dead 狀態多久後允許探測)
	CooldownPeriod time.Duration
}

// DefaultBreakerConfig 預設熔斷器配置
func DefaultBreakerConfig() BreakerConfig {
	return BreakerConfig{
		ErrorThreshold:    0.2,           // 20% 錯誤率
		UnstableThreshold: 0.5,           // 50% 錯誤率
		WindowSize:        10,            // 10 次請求
		CooldownPeriod:    30 * time.Second,
	}
}

// =============================================================================
// 熔斷器
// =============================================================================

// CircuitBreaker 熔斷器
type CircuitBreaker struct {
	mu sync.RWMutex

	config  BreakerConfig
	tracker *Tracker
	state   State

	// Dead 狀態進入時間
	deadSince time.Time

	// 是否正在探測中
	probing bool
}

// NewCircuitBreaker 建立新的熔斷器
func NewCircuitBreaker(config BreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config:  config,
		tracker: NewTracker(config.WindowSize),
		state:   StateHealthy,
	}
}

// State 取得當前狀態
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// AllowRequest 檢查是否允許請求
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateHealthy:
		return true

	case StateUnstable:
		return true

	case StateDead:
		// 檢查是否已過冷卻時間
		if time.Since(cb.deadSince) >= cb.config.CooldownPeriod {
			// 允許一次探測，轉為 Half-Open (Unstable)
			cb.state = StateUnstable
			cb.probing = true
			cb.tracker.Reset()
			return true
		}
		return false
	}

	return false
}

// ReportResult 報告請求結果
func (cb *CircuitBreaker) ReportResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.tracker.RecordFailure(err)
	} else {
		cb.tracker.RecordSuccess(0) // 延遲由外部追蹤
	}

	// 更新狀態
	cb.updateState(err)
}

// updateState 更新狀態機 (必須持有鎖)
func (cb *CircuitBreaker) updateState(lastErr error) {
	errorRate := cb.tracker.ErrorRate()

	switch cb.state {
	case StateHealthy:
		if errorRate >= cb.config.ErrorThreshold {
			cb.state = StateUnstable
		}

	case StateUnstable:
		if cb.probing {
			// 探測模式
			if lastErr != nil {
				// 探測失敗，回到 Dead
				cb.state = StateDead
				cb.deadSince = time.Now()
				cb.probing = false
			} else {
				// 探測成功，恢復 Healthy
				cb.state = StateHealthy
				cb.probing = false
				cb.tracker.Reset()
			}
		} else {
			// 正常 Unstable 模式
			if errorRate >= cb.config.UnstableThreshold {
				cb.state = StateDead
				cb.deadSince = time.Now()
			} else if errorRate < cb.config.ErrorThreshold {
				cb.state = StateHealthy
			}
		}

	case StateDead:
		// Dead 狀態不會自動轉換，需等待冷卻後探測
	}
}

// Reset 重置熔斷器
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateHealthy
	cb.tracker.Reset()
	cb.probing = false
}

// GetStats 取得統計資訊
func (cb *CircuitBreaker) GetStats() BreakerStats {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return BreakerStats{
		State:      cb.state,
		ErrorRate:  cb.tracker.ErrorRate(),
		TrackerStats: cb.tracker.GetStats(),
		DeadSince:  cb.deadSince,
		IsProbing:  cb.probing,
	}
}

// BreakerStats 熔斷器統計資訊
type BreakerStats struct {
	State        State
	ErrorRate    float64
	TrackerStats Stats
	DeadSince    time.Time
	IsProbing    bool
}
