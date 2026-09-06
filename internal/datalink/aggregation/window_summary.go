package aggregation

import (
	"sort"
	"time"
)

// WindowSummaryConfig 統計視窗配置。
type WindowSummaryConfig struct {
	MaxHoldSeconds      float64 `json:"max_hold_seconds"`
	CalculationRevision int64   `json:"calculation_revision"`
}

// CalculateWindowSummary 計算給定視窗的時間加權平均、採樣極值與覆蓋時長。
func CalculateWindowSummary(
	measurementID string,
	windowStart time.Time,
	windowEnd time.Time,
	samples []TelemetrySample,
	config WindowSummaryConfig,
) SummaryResult {
	totalSeconds := windowEnd.Sub(windowStart).Seconds()
	if totalSeconds <= 0 {
		totalSeconds = 1.0
	}

	maxHold := config.MaxHoldSeconds
	if maxHold <= 0 {
		maxHold = 60.0 // 預設 60 秒
	}

	// 依時間排序
	sortedSamples := make([]TelemetrySample, len(samples))
	copy(sortedSamples, samples)
	sort.SliceStable(sortedSamples, func(i, j int) bool {
		return sortedSamples[i].ObservedAt.Before(sortedSamples[j].ObservedAt)
	})

	var (
		weightedSum      float64
		validDuration    float64
		minVal           *float64
		maxVal           *float64
		firstVal         *float64
		lastVal          *float64
		validSampleCount int
	)

	for i, s := range sortedSamples {
		if s.Quality != QualityGood || s.ValueNumeric == nil {
			continue
		}

		val := *s.ValueNumeric
		obsTime := s.ObservedAt

		if !obsTime.Before(windowEnd) {
			continue
		}

		segStart := obsTime
		if segStart.Before(windowStart) {
			segStart = windowStart
		}

		// 計算 hold 結束時間
		holdEnd := obsTime.Add(time.Duration(maxHold * float64(time.Second)))

		// 尋找下一個樣本的時間或視窗終點
		segEnd := holdEnd
		if i+1 < len(sortedSamples) {
			nextTime := sortedSamples[i+1].ObservedAt
			if nextTime.Before(segEnd) {
				segEnd = nextTime
			}
		}

		if segEnd.After(windowEnd) {
			segEnd = windowEnd
		}

		if segEnd.After(segStart) {
			duration := segEnd.Sub(segStart).Seconds()
			weightedSum += val * duration
			validDuration += duration

			if minVal == nil || val < *minVal {
				minVal = &val
			}
			if maxVal == nil || val > *maxVal {
				maxVal = &val
			}
			if firstVal == nil {
				firstVal = &val
			}
			lastVal = &val
			validSampleCount++
		}
	}

	result := SummaryResult{
		MeasurementID: measurementID,
		WindowStart:   windowStart,
		WindowEnd:     windowEnd,
		SampleCount:   validSampleCount,
		Coverage: Coverage{
			ValidSeconds: validDuration,
			TotalSeconds: totalSeconds,
			Ratio:        validDuration / totalSeconds,
		},
		CalculationRevision: config.CalculationRevision,
	}

	if validDuration > 0 {
		mean := weightedSum / validDuration
		result.TimeWeightedMean = &mean
		result.SampledMin = minVal
		result.SampledMax = maxVal
		result.FirstValue = firstVal
		result.LastValue = lastVal
	}

	return result
}
