package aggregation

import (
	"time"
)

// QualityIndicator 表示樣本品質狀態。
type QualityIndicator string

const (
	QualityGood      QualityIndicator = "good"
	QualityBad       QualityIndicator = "bad"
	QualityUncertain QualityIndicator = "uncertain"
)

// TelemetrySample 供計算引擎消費的統一標準樣本。
type TelemetrySample struct {
	MeasurementID string           `json:"measurement_id"`
	EquipmentID   string           `json:"equipment_id,omitempty"`
	SeriesEpoch   int64            `json:"series_epoch"`
	ObservedAt    time.Time        `json:"observed_at"`
	ReceivedAt    time.Time        `json:"received_at"`
	Value         interface{}      `json:"value"`
	ValueNumeric  *float64         `json:"value_numeric,omitempty"`
	ValueString   *string          `json:"value_string,omitempty"`
	ValueBool     *bool            `json:"value_bool,omitempty"`
	ValueInt64    *int64           `json:"value_int64,omitempty"`
	ValueUint64   *uint64          `json:"value_uint64,omitempty"`
	Quality       QualityIndicator `json:"quality"`
	IntervalID    string           `json:"interval_id,omitempty"`
	Sequence      int64            `json:"sequence,omitempty"`
}

// Coverage 記錄視窗內有效採樣時長與覆蓋比例。
type Coverage struct {
	ValidSeconds float64 `json:"valid_seconds"`
	TotalSeconds float64 `json:"total_seconds"`
	Ratio        float64 `json:"ratio"`
}

// SummaryResult 時間加權統計結果。
type SummaryResult struct {
	MeasurementID       string    `json:"measurement_id"`
	WindowStart         time.Time `json:"window_start"`
	WindowEnd           time.Time `json:"window_end"`
	TimeWeightedMean    *float64  `json:"time_weighted_mean"`
	SampledMin          *float64  `json:"sampled_min"`
	SampledMax          *float64  `json:"sampled_max"`
	FirstValue          *float64  `json:"first_value"`
	LastValue           *float64  `json:"last_value"`
	SampleCount         int       `json:"sample_count"`
	Coverage            Coverage  `json:"coverage"`
	IsProvisional       bool      `json:"is_provisional"`
	CalculationRevision int64     `json:"calculation_revision"`
}

// UsageResult 累積量與差分計算結果。
type UsageResult struct {
	MeasurementID       string    `json:"measurement_id"`
	SeriesEpoch         int64     `json:"series_epoch"`
	WindowStart         time.Time `json:"window_start"`
	WindowEnd           time.Time `json:"window_end"`
	StartReading        *float64  `json:"start_reading,omitempty"`
	EndReading          *float64  `json:"end_reading,omitempty"`
	UsageDelta          *float64  `json:"usage_delta"`
	KnownSubtotal       *float64  `json:"known_subtotal,omitempty"`
	IsComplete          bool      `json:"is_complete"`
	IsEstimated         bool      `json:"is_estimated"`
	HasDiscontinuity    bool      `json:"has_discontinuity"`
	HasRollover         bool      `json:"has_rollover"`
	IsUncertain         bool      `json:"is_uncertain"`
	UncertainReason     string    `json:"uncertain_reason,omitempty"`
	CalculationRevision int64     `json:"calculation_revision"`
}

// StateDuration 記錄單一狀態持續時間。
type StateDuration struct {
	State           string  `json:"state"`
	DurationSeconds float64 `json:"duration_seconds"`
	Percentage      float64 `json:"percentage"`
}

// StateTrackingResult 狀態變更與時間統計結果。
type StateTrackingResult struct {
	MeasurementID       string          `json:"measurement_id"`
	WindowStart         time.Time       `json:"window_start"`
	WindowEnd           time.Time       `json:"window_end"`
	InitialState        string          `json:"initial_state"`
	FinalState          string          `json:"final_state"`
	TransitionsCount    int             `json:"transitions_count"`
	StateDurations      []StateDuration `json:"state_durations"`
	UnknownDurationSec  float64         `json:"unknown_duration_seconds"`
	CalculationRevision int64           `json:"calculation_revision"`
}

// RateIntegrationResult 速率積分估算結果。
type RateIntegrationResult struct {
	MeasurementID       string    `json:"measurement_id"`
	WindowStart         time.Time `json:"window_start"`
	WindowEnd           time.Time `json:"window_end"`
	IntegratedValue     float64   `json:"integrated_value"`
	Unit                string    `json:"unit"`
	Method              string    `json:"method"` // e.g. "trapezoidal"
	IsEstimated         bool      `json:"is_estimated"`
	CalculationRevision int64     `json:"calculation_revision"`
}

// DerivedCalculationResult 衍生量測計算結果。
type DerivedCalculationResult struct {
	FormulaID           string    `json:"formula_id"`
	OutputMeasurementID string    `json:"output_measurement_id"`
	ObservedAt          time.Time `json:"observed_at"`
	ValueNumeric        *float64  `json:"value_numeric"`
	IsValid             bool      `json:"is_valid"`
	ErrorMessage        string    `json:"error_message,omitempty"`
	CalculationRevision int64     `json:"calculation_revision"`
}
