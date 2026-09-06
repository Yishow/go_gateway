package aggregation

import (
	"fmt"
	"sort"
	"time"
)

// CounterUsageConfig 累積量差分配置。
type CounterUsageConfig struct {
	IsSigned            bool     `json:"is_signed"`
	RolloverModulus     *float64 `json:"rollover_modulus,omitempty"`
	MaxRatePerSec       *float64 `json:"max_rate_per_sec,omitempty"`
	TargetEpoch         int64    `json:"target_epoch,omitempty"`
	CalculationRevision int64    `json:"calculation_revision"`
}

// CalculateCounterUsage 計算視窗內的累積計量差值、回捲與不連續性。
func CalculateCounterUsage(
	measurementID string,
	windowStart time.Time,
	windowEnd time.Time,
	samples []TelemetrySample,
	config CounterUsageConfig,
) UsageResult {
	// 過濾出視窗內、有效品質的數值樣本
	var validSamples []TelemetrySample
	for _, s := range samples {
		if s.Quality == QualityGood && s.ValueNumeric != nil {
			if !s.ObservedAt.Before(windowStart) && s.ObservedAt.Before(windowEnd) {
				validSamples = append(validSamples, s)
			}
		}
	}

	sort.SliceStable(validSamples, func(i, j int) bool {
		return validSamples[i].ObservedAt.Before(validSamples[j].ObservedAt)
	})

	result := UsageResult{
		MeasurementID:       measurementID,
		WindowStart:         windowStart,
		WindowEnd:           windowEnd,
		CalculationRevision: config.CalculationRevision,
		IsComplete:          true,
	}

	if len(validSamples) == 0 {
		result.IsComplete = false
		return result
	}

	result.StartReading = validSamples[0].ValueNumeric
	result.EndReading = validSamples[len(validSamples)-1].ValueNumeric
	result.SeriesEpoch = validSamples[0].SeriesEpoch

	if len(validSamples) == 1 {
		zero := 0.0
		result.UsageDelta = &zero
		return result
	}

	var totalDelta float64
	hasDiscontinuity := false
	hasRollover := false
	isUncertain := false
	var uncertainReason string

	for i := 1; i < len(validSamples); i++ {
		prev := validSamples[i-1]
		curr := validSamples[i]

		// 檢查 epoch 邊界
		if prev.SeriesEpoch != curr.SeriesEpoch {
			hasDiscontinuity = true
			if config.TargetEpoch > 0 && curr.SeriesEpoch != config.TargetEpoch {
				continue
			}
			if config.TargetEpoch > 0 && prev.SeriesEpoch != config.TargetEpoch {
				// 新 epoch 第一筆作為基準，不相減
				continue
			}
			// 無指定 targetEpoch 時，跨 epoch 不相減
			continue
		}

		if config.TargetEpoch > 0 && curr.SeriesEpoch != config.TargetEpoch {
			continue
		}

		prevVal := *prev.ValueNumeric
		currVal := *curr.ValueNumeric
		diff := currVal - prevVal

		if diff >= 0 {
			totalDelta += diff
		} else {
			// 負差值處理
			if config.IsSigned {
				totalDelta += diff
			} else if config.RolloverModulus != nil && *config.RolloverModulus > 0 {
				mod := *config.RolloverModulus
				wrapDiff := (mod - prevVal) + currVal
				if config.MaxRatePerSec != nil && *config.MaxRatePerSec > 0 {
					dt := curr.ObservedAt.Sub(prev.ObservedAt).Seconds()
					if dt <= 0 {
						dt = 1.0
					}
					maxAllowed := *config.MaxRatePerSec * dt
					if wrapDiff <= maxAllowed {
						totalDelta += wrapDiff
						hasRollover = true
					} else {
						isUncertain = true
						hasDiscontinuity = true
						uncertainReason = fmt.Sprintf("wrap difference %.2f exceeds max allowable rate", wrapDiff)
					}
				} else {
					totalDelta += wrapDiff
					hasRollover = true
				}
			} else {
				isUncertain = true
				hasDiscontinuity = true
				uncertainReason = fmt.Sprintf("unexplained drop from %.2f to %.2f", prevVal, currVal)
			}
		}
	}

	result.HasDiscontinuity = hasDiscontinuity
	result.HasRollover = hasRollover
	result.IsUncertain = isUncertain
	result.UncertainReason = uncertainReason

	if !isUncertain {
		result.UsageDelta = &totalDelta
	}

	return result
}

// CalculateDeltaUsage 針對具備唯一 interval/event 標識的增量序列進行去重累加。
func CalculateDeltaUsage(
	measurementID string,
	windowStart time.Time,
	windowEnd time.Time,
	samples []TelemetrySample,
) UsageResult {
	seenIntervals := make(map[string]bool)
	var totalDelta float64

	for _, s := range samples {
		if s.Quality != QualityGood || s.ValueNumeric == nil {
			continue
		}
		if s.ObservedAt.Before(windowStart) || !s.ObservedAt.Before(windowEnd) {
			continue
		}

		intervalKey := s.IntervalID
		if intervalKey == "" {
			intervalKey = fmt.Sprintf("%d_%d", s.Sequence, s.ObservedAt.UnixNano())
		}

		if !seenIntervals[intervalKey] {
			seenIntervals[intervalKey] = true
			totalDelta += *s.ValueNumeric
		}
	}

	return UsageResult{
		MeasurementID: measurementID,
		WindowStart:   windowStart,
		WindowEnd:     windowEnd,
		UsageDelta:    &totalDelta,
		IsComplete:    true,
	}
}
