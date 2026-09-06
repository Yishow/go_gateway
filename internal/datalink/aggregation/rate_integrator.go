package aggregation

import (
	"sort"
	"time"
)

// RateIntegrationConfig 速率積分配置。
type RateIntegrationConfig struct {
	InputUnit           string  `json:"input_unit"`
	OutputUnit          string  `json:"output_unit"`
	TimeDivisor         float64 `json:"time_divisor"` // 例如秒轉小時為 3600，秒轉分維 60
	MaxGapSeconds       float64 `json:"max_gap_seconds"`
	CalculationRevision int64   `json:"calculation_revision"`
}

// IntegrateRate 採用梯形積分法對相鄰有效採樣進行速率積分。
func IntegrateRate(
	measurementID string,
	windowStart time.Time,
	windowEnd time.Time,
	samples []TelemetrySample,
	config RateIntegrationConfig,
) RateIntegrationResult {
	result := RateIntegrationResult{
		MeasurementID:       measurementID,
		WindowStart:         windowStart,
		WindowEnd:           windowEnd,
		Unit:                config.OutputUnit,
		Method:              "trapezoidal",
		IsEstimated:         true,
		CalculationRevision: config.CalculationRevision,
	}

	divisor := config.TimeDivisor
	if divisor <= 0 {
		divisor = 3600.0 // 預設秒轉小時
	}

	maxGap := config.MaxGapSeconds
	if maxGap <= 0 {
		maxGap = 300.0 // 預設 5 分鐘
	}

	var validSamples []TelemetrySample
	for _, s := range samples {
		if s.Quality == QualityGood && s.ValueNumeric != nil {
			if !s.ObservedAt.Before(windowStart) && !s.ObservedAt.After(windowEnd) {
				validSamples = append(validSamples, s)
			}
		}
	}

	if len(validSamples) < 2 {
		return result
	}

	sort.SliceStable(validSamples, func(i, j int) bool {
		return validSamples[i].ObservedAt.Before(validSamples[j].ObservedAt)
	})

	var totalArea float64
	for i := 1; i < len(validSamples); i++ {
		prev := validSamples[i-1]
		curr := validSamples[i]

		dt := curr.ObservedAt.Sub(prev.ObservedAt).Seconds()
		if dt <= 0 || dt > maxGap {
			continue // 超過 maxGap 視為中斷，不作積分
		}

		v1 := *prev.ValueNumeric
		v2 := *curr.ValueNumeric

		// 梯形面積 = ((v1 + v2) / 2) * (dt / divisor)
		area := ((v1 + v2) / 2.0) * (dt / divisor)
		totalArea += area
	}

	result.IntegratedValue = totalArea
	return result
}
