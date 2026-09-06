package aggregation

import (
	"math"
	"testing"
	"time"
)

func floatPtr(v float64) *float64 {
	return &v
}

func TestWindowSummary_W01_UnequalSamplingIntervals(t *testing.T) {
	// W01: 10持續10秒、20持續50秒 -> mean 1100/60 = 18.3333..., min 10, max 20, coverage 60秒
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	windowStart := baseTime
	windowEnd := baseTime.Add(60 * time.Second)

	samples := []TelemetrySample{
		{
			MeasurementID: "meas-gauge-1",
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(10.0),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-gauge-1",
			ObservedAt:    baseTime.Add(10 * time.Second),
			ValueNumeric:  floatPtr(20.0),
			Quality:       QualityGood,
		},
	}

	config := WindowSummaryConfig{
		MaxHoldSeconds: 60, // 足以覆蓋剩餘時間
	}

	result := CalculateWindowSummary("meas-gauge-1", windowStart, windowEnd, samples, config)

	if result.SampleCount != 2 {
		t.Fatalf("expected sample count 2, got %d", result.SampleCount)
	}
	if result.TimeWeightedMean == nil {
		t.Fatalf("expected non-nil mean")
	}

	expectedMean := 1100.0 / 60.0
	if math.Abs(*result.TimeWeightedMean-expectedMean) > 1e-6 {
		t.Errorf("expected mean %.6f, got %.6f", expectedMean, *result.TimeWeightedMean)
	}
	if result.SampledMin == nil || *result.SampledMin != 10.0 {
		t.Errorf("expected min 10, got %v", result.SampledMin)
	}
	if result.SampledMax == nil || *result.SampledMax != 20.0 {
		t.Errorf("expected max 20, got %v", result.SampledMax)
	}
	if result.Coverage.ValidSeconds != 60.0 {
		t.Errorf("expected 60 valid seconds, got %.2f", result.Coverage.ValidSeconds)
	}
	if result.Coverage.Ratio != 1.0 {
		t.Errorf("expected coverage ratio 1.0, got %.2f", result.Coverage.Ratio)
	}
}

func TestWindowSummary_W02_NoValidData(t *testing.T) {
	// W02: 60秒全缺資料 -> count 0, mean/min/max null, 非 0 值
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	windowStart := baseTime
	windowEnd := baseTime.Add(60 * time.Second)

	var samples []TelemetrySample

	config := WindowSummaryConfig{
		MaxHoldSeconds: 10,
	}

	result := CalculateWindowSummary("meas-gauge-1", windowStart, windowEnd, samples, config)

	if result.SampleCount != 0 {
		t.Errorf("expected sample count 0, got %d", result.SampleCount)
	}
	if result.TimeWeightedMean != nil {
		t.Errorf("expected nil mean, got %v", *result.TimeWeightedMean)
	}
	if result.SampledMin != nil {
		t.Errorf("expected nil min, got %v", *result.SampledMin)
	}
	if result.SampledMax != nil {
		t.Errorf("expected nil max, got %v", *result.SampledMax)
	}
	if result.Coverage.ValidSeconds != 0 {
		t.Errorf("expected 0 valid seconds, got %v", result.Coverage.ValidSeconds)
	}
	if result.Coverage.Ratio != 0 {
		t.Errorf("expected 0 coverage ratio, got %v", result.Coverage.Ratio)
	}
}

func TestWindowSummary_QualityInterruptionAndMaxHold(t *testing.T) {
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	windowStart := baseTime
	windowEnd := baseTime.Add(60 * time.Second)

	// 10s valid with value 100, then bad quality at 10s, then good quality value 200 at 30s with max_hold=10s
	samples := []TelemetrySample{
		{
			MeasurementID: "meas-gauge-1",
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(100.0),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-gauge-1",
			ObservedAt:    baseTime.Add(10 * time.Second),
			Quality:       QualityBad, // 立即中斷 100 的 hold
		},
		{
			MeasurementID: "meas-gauge-1",
			ObservedAt:    baseTime.Add(30 * time.Second),
			ValueNumeric:  floatPtr(200.0),
			Quality:       QualityGood,
		},
	}

	config := WindowSummaryConfig{
		MaxHoldSeconds: 10, // 30s sample 只 hold 到 40s
	}

	result := CalculateWindowSummary("meas-gauge-1", windowStart, windowEnd, samples, config)

	// 有效時長: 100 持續 10s (0~10s), 200 持續 10s (30~40s) = 20s
	// 加權平均: (100*10 + 200*10) / 20 = 3000 / 20 = 150
	if result.Coverage.ValidSeconds != 20.0 {
		t.Errorf("expected 20 valid seconds, got %.2f", result.Coverage.ValidSeconds)
	}
	if result.TimeWeightedMean == nil || *result.TimeWeightedMean != 150.0 {
		t.Errorf("expected mean 150, got %v", result.TimeWeightedMean)
	}
	if result.SampledMin == nil || *result.SampledMin != 100.0 {
		t.Errorf("expected min 100, got %v", result.SampledMin)
	}
	if result.SampledMax == nil || *result.SampledMax != 200.0 {
		t.Errorf("expected max 200, got %v", result.SampledMax)
	}
}
