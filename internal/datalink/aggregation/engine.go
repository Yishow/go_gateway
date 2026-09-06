package aggregation

import (
	"fmt"
	"sync"
	"time"
)

// EngineConfig 聚合引擎全域配置。
type EngineConfig struct {
	CorrectionHorizonHours int `json:"correction_horizon_hours"`
	GracePeriodSeconds     int `json:"grace_period_seconds"`
}

// Checkpoint 記錄單一量測項目的計算水位與最新基準。
type Checkpoint struct {
	MeasurementID       string    `json:"measurement_id"`
	Watermark           time.Time `json:"watermark"`
	LastEpoch           int64     `json:"last_epoch"`
	LastReadingNumeric  *float64  `json:"last_reading_numeric,omitempty"`
	LastObservedAt      time.Time `json:"last_observed_at"`
	CalculationRevision int64     `json:"calculation_revision"`
}

// TelemetryAggregationEngine 遙測聚合與計量核心引擎。
type TelemetryAggregationEngine struct {
	mu          sync.RWMutex
	config      EngineConfig
	samples     map[string][]TelemetrySample
	revisions   map[string]int64
	watermark   time.Time
	checkpoints map[string]*Checkpoint
}

// NewTelemetryAggregationEngine 建立新的聚合引擎實例。
func NewTelemetryAggregationEngine(config EngineConfig) *TelemetryAggregationEngine {
	if config.CorrectionHorizonHours <= 0 {
		config.CorrectionHorizonHours = 24
	}
	if config.GracePeriodSeconds <= 0 {
		config.GracePeriodSeconds = 10
	}
	return &TelemetryAggregationEngine{
		config:      config,
		samples:     make(map[string][]TelemetrySample),
		revisions:   make(map[string]int64),
		checkpoints: make(map[string]*Checkpoint),
	}
}

// AdvanceWatermark 前進全域處理水位時間。
func (e *TelemetryAggregationEngine) AdvanceWatermark(t time.Time) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if t.After(e.watermark) {
		e.watermark = t
	}
}

// IngestSample 接收標準遙測樣本，並驗證未超過可修正期限 (Correction Horizon)。
func (e *TelemetryAggregationEngine) IngestSample(s TelemetrySample) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	horizon := time.Duration(e.config.CorrectionHorizonHours) * time.Hour
	if !e.watermark.IsZero() && s.ObservedAt.Before(e.watermark.Add(-horizon)) {
		return fmt.Errorf("sample observed at %v is beyond correction horizon (%d hours)", s.ObservedAt, e.config.CorrectionHorizonHours)
	}

	e.samples[s.MeasurementID] = append(e.samples[s.MeasurementID], s)
	if s.ObservedAt.After(e.watermark) {
		e.watermark = s.ObservedAt
	}

	// 遞增該 measurement 的修訂版次
	e.revisions[s.MeasurementID]++
	return nil
}

// ComputeWindowSummary 計算指定量測項目與時間視窗的時間加權統計。
func (e *TelemetryAggregationEngine) ComputeWindowSummary(
	measurementID string,
	windowStart time.Time,
	windowEnd time.Time,
) (SummaryResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	samples := e.samples[measurementID]
	rev := e.revisions[measurementID]
	if rev == 0 {
		rev = 1
	}

	config := WindowSummaryConfig{
		MaxHoldSeconds:      60,
		CalculationRevision: rev,
	}

	result := CalculateWindowSummary(measurementID, windowStart, windowEnd, samples, config)
	return result, nil
}

// ComputeCounterUsage 計算指定量測項目的累積計量差值。
func (e *TelemetryAggregationEngine) ComputeCounterUsage(
	measurementID string,
	windowStart time.Time,
	windowEnd time.Time,
	config CounterUsageConfig,
) (UsageResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	samples := e.samples[measurementID]
	rev := e.revisions[measurementID]
	if rev == 0 {
		rev = 1
	}
	config.CalculationRevision = rev

	result := CalculateCounterUsage(measurementID, windowStart, windowEnd, samples, config)
	return result, nil
}
