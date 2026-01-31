// Package health 提供設備健康度追蹤與熔斷機制
//
// 本套件實作滑動視窗統計、熔斷狀態機與背景探測等功能，
// 用於保護系統免受故障設備的影響。
package health

import (
	"sync"
	"time"
)

// =============================================================================
// 類型定義
// =============================================================================

// Record 單次請求記錄
type Record struct {
	Success   bool
	Duration  time.Duration
	Timestamp time.Time
}

// Stats 健康統計資訊
type Stats struct {
	TotalRequests int
	SuccessCount  int
	FailureCount  int
	ErrorRate     float64
	AvgLatency    time.Duration
	P99Latency    time.Duration
}

// =============================================================================
// Tracker 健康度追蹤器
// =============================================================================

// Tracker 滑動視窗健康度追蹤器
type Tracker struct {
	mu sync.RWMutex

	// 滑動視窗大小
	windowSize int

	// 環形緩衝區
	records []Record
	index   int
	count   int
}

// NewTracker 建立新的健康度追蹤器
func NewTracker(windowSize int) *Tracker {
	if windowSize <= 0 {
		windowSize = 10
	}
	return &Tracker{
		windowSize: windowSize,
		records:    make([]Record, windowSize),
		index:      0,
		count:      0,
	}
}

// RecordSuccess 記錄成功請求
func (t *Tracker) RecordSuccess(duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.records[t.index] = Record{
		Success:   true,
		Duration:  duration,
		Timestamp: time.Now(),
	}
	t.advance()
}

// RecordFailure 記錄失敗請求
func (t *Tracker) RecordFailure(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.records[t.index] = Record{
		Success:   false,
		Duration:  0,
		Timestamp: time.Now(),
	}
	t.advance()
}

// advance 推進環形緩衝區索引
func (t *Tracker) advance() {
	t.index = (t.index + 1) % t.windowSize
	if t.count < t.windowSize {
		t.count++
	}
}

// ErrorRate 計算錯誤率
func (t *Tracker) ErrorRate() float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.count == 0 {
		return 0
	}

	failures := 0
	for i := 0; i < t.count; i++ {
		if !t.records[i].Success {
			failures++
		}
	}

	return float64(failures) / float64(t.count)
}

// GetStats 取得完整統計資訊
func (t *Tracker) GetStats() Stats {
	t.mu.RLock()
	defer t.mu.RUnlock()

	stats := Stats{
		TotalRequests: t.count,
	}

	if t.count == 0 {
		return stats
	}

	var totalDuration time.Duration
	durations := make([]time.Duration, 0, t.count)

	for i := 0; i < t.count; i++ {
		record := t.records[i]
		if record.Success {
			stats.SuccessCount++
			totalDuration += record.Duration
			durations = append(durations, record.Duration)
		} else {
			stats.FailureCount++
		}
	}

	stats.ErrorRate = float64(stats.FailureCount) / float64(t.count)

	if stats.SuccessCount > 0 {
		stats.AvgLatency = totalDuration / time.Duration(stats.SuccessCount)
	}

	// 計算 P99 延遲 (簡化版：取最大值)
	if len(durations) > 0 {
		maxDuration := durations[0]
		for _, d := range durations {
			if d > maxDuration {
				maxDuration = d
			}
		}
		stats.P99Latency = maxDuration
	}

	return stats
}

// Reset 重置追蹤器
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.records = make([]Record, t.windowSize)
	t.index = 0
	t.count = 0
}
