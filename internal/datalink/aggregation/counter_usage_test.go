package aggregation

import (
	"math"
	"testing"
	"time"
)

func TestCounterUsage_C01_StandardMonotonicCounter(t *testing.T) {
	// C01: 12000.0→12000.4→12000.7 kWh -> 初筆只建基準，0.4+0.3=0.7
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	windowStart := baseTime
	windowEnd := baseTime.Add(60 * time.Minute)

	samples := []TelemetrySample{
		{
			MeasurementID: "meas-kwh",
			SeriesEpoch:   1,
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(12000.0),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-kwh",
			SeriesEpoch:   1,
			ObservedAt:    baseTime.Add(10 * time.Minute),
			ValueNumeric:  floatPtr(12000.4),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-kwh",
			SeriesEpoch:   1,
			ObservedAt:    baseTime.Add(20 * time.Minute),
			ValueNumeric:  floatPtr(12000.7),
			Quality:       QualityGood,
		},
	}

	config := CounterUsageConfig{
		IsSigned: false,
	}

	result := CalculateCounterUsage("meas-kwh", windowStart, windowEnd, samples, config)

	if !result.IsComplete {
		t.Errorf("expected complete result")
	}
	if result.UsageDelta == nil {
		t.Fatalf("expected non-nil usage delta")
	}
	if math.Abs(*result.UsageDelta-0.7) > 1e-6 {
		t.Errorf("expected usage delta 0.7, got %.6f", *result.UsageDelta)
	}
	if result.StartReading == nil || *result.StartReading != 12000.0 {
		t.Errorf("expected start reading 12000.0, got %v", result.StartReading)
	}
	if result.EndReading == nil || *result.EndReading != 12000.7 {
		t.Errorf("expected end reading 12000.7, got %v", result.EndReading)
	}
	if result.HasDiscontinuity {
		t.Errorf("unexpected discontinuity")
	}
}

func TestCounterUsage_C02_ZeroUsage(t *testing.T) {
	// C02: 123→123 -> 有效零用量，不是 missing
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	samples := []TelemetrySample{
		{
			MeasurementID: "meas-kwh",
			SeriesEpoch:   1,
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(123.0),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-kwh",
			SeriesEpoch:   1,
			ObservedAt:    baseTime.Add(10 * time.Minute),
			ValueNumeric:  floatPtr(123.0),
			Quality:       QualityGood,
		},
	}

	result := CalculateCounterUsage("meas-kwh", baseTime, baseTime.Add(15*time.Minute), samples, CounterUsageConfig{})

	if result.UsageDelta == nil {
		t.Fatalf("expected non-nil usage delta")
	}
	if *result.UsageDelta != 0.0 {
		t.Errorf("expected 0.0 usage, got %v", *result.UsageDelta)
	}
	if result.IsUncertain {
		t.Errorf("unexpected uncertain flag")
	}
}

func TestCounterUsage_C03_UnexplainedNegativeDrop(t *testing.T) {
	// C03: counter 999→4，無 reset 證據 -> uncertain，不用 abs，不自動把 4 算用量
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	samples := []TelemetrySample{
		{
			MeasurementID: "meas-kwh",
			SeriesEpoch:   1,
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(999.0),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-kwh",
			SeriesEpoch:   1,
			ObservedAt:    baseTime.Add(10 * time.Minute),
			ValueNumeric:  floatPtr(4.0),
			Quality:       QualityGood,
		},
	}

	result := CalculateCounterUsage("meas-kwh", baseTime, baseTime.Add(15*time.Minute), samples, CounterUsageConfig{})

	if !result.IsUncertain {
		t.Errorf("expected uncertain result on unexplained drop")
	}
	if !result.HasDiscontinuity {
		t.Errorf("expected discontinuity flag")
	}
	// 不應產生 995 或 4 作為合法 delta
	if result.UsageDelta != nil {
		t.Errorf("expected nil usage delta for uncertain calculation, got %v", *result.UsageDelta)
	}
}

func TestCounterUsage_C04_ConfirmedNewEpoch(t *testing.T) {
	// C04: 已確認新 epoch 4→6 -> 新段 2，舊新中間仍不完整
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	samples := []TelemetrySample{
		{
			MeasurementID: "meas-kwh",
			SeriesEpoch:   1,
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(999.0),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-kwh",
			SeriesEpoch:   2, // 新 epoch
			ObservedAt:    baseTime.Add(10 * time.Minute),
			ValueNumeric:  floatPtr(4.0),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-kwh",
			SeriesEpoch:   2,
			ObservedAt:    baseTime.Add(20 * time.Minute),
			ValueNumeric:  floatPtr(6.0),
			Quality:       QualityGood,
		},
	}

	result := CalculateCounterUsage("meas-kwh", baseTime, baseTime.Add(30*time.Minute), samples, CounterUsageConfig{
		TargetEpoch: 2,
	})

	if result.UsageDelta == nil || *result.UsageDelta != 2.0 {
		t.Errorf("expected usage delta 2.0 in epoch 2, got %v", result.UsageDelta)
	}
	if !result.HasDiscontinuity {
		t.Errorf("expected discontinuity flag due to epoch transition")
	}
}

func TestCounterUsage_C05_ConfirmedSingleRollover(t *testing.T) {
	// C05: 65530→4，M=65536，可證單次回捲 -> 10；不能排除多次則 unknown
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	samples := []TelemetrySample{
		{
			MeasurementID: "meas-16bit",
			SeriesEpoch:   1,
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(65530.0),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-16bit",
			SeriesEpoch:   1,
			ObservedAt:    baseTime.Add(10 * time.Second),
			ValueNumeric:  floatPtr(4.0),
			Quality:       QualityGood,
		},
	}

	config := CounterUsageConfig{
		RolloverModulus: floatPtr(65536.0),
		MaxRatePerSec:   floatPtr(100.0), // 10s 內最多 1000，(65536-65530)+4 = 10 <= 1000，合理解析單次回捲
	}

	result := CalculateCounterUsage("meas-16bit", baseTime, baseTime.Add(1*time.Minute), samples, config)

	if !result.HasRollover {
		t.Errorf("expected rollover flag")
	}
	if result.UsageDelta == nil || *result.UsageDelta != 10.0 {
		t.Errorf("expected usage delta 10.0, got %v", result.UsageDelta)
	}
}

func TestCounterUsage_C06_SignedCounter(t *testing.T) {
	// C06: signed net 100→98 -> -2，不當 reset
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	samples := []TelemetrySample{
		{
			MeasurementID: "meas-net-flow",
			SeriesEpoch:   1,
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(100.0),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-net-flow",
			SeriesEpoch:   1,
			ObservedAt:    baseTime.Add(10 * time.Second),
			ValueNumeric:  floatPtr(98.0),
			Quality:       QualityGood,
		},
	}

	config := CounterUsageConfig{
		IsSigned: true,
	}

	result := CalculateCounterUsage("meas-net-flow", baseTime, baseTime.Add(1*time.Minute), samples, config)

	if result.IsUncertain {
		t.Errorf("unexpected uncertain flag for signed counter")
	}
	if result.UsageDelta == nil || *result.UsageDelta != -2.0 {
		t.Errorf("expected -2.0 delta, got %v", result.UsageDelta)
	}
}

func TestCounterUsage_C07_C08_DeltaDeduplicationAndAccumulation(t *testing.T) {
	// C07: delta 1 重送同 identity -> 只算 1
	// C08: 不同非重疊 interval，各 delta 1 -> 合計 2
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	samples := []TelemetrySample{
		{
			MeasurementID: "meas-batch-delta",
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(1.0),
			IntervalID:    "interval-1",
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-batch-delta",
			ObservedAt:    baseTime.Add(5 * time.Second),
			ValueNumeric:  floatPtr(1.0),
			IntervalID:    "interval-1", // 重複傳送同一個 interval
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-batch-delta",
			ObservedAt:    baseTime.Add(10 * time.Second),
			ValueNumeric:  floatPtr(1.0),
			IntervalID:    "interval-2", // 獨立的另一個 interval
			Quality:       QualityGood,
		},
	}

	result := CalculateDeltaUsage("meas-batch-delta", baseTime, baseTime.Add(1*time.Minute), samples)

	if result.UsageDelta == nil || *result.UsageDelta != 2.0 {
		t.Errorf("expected total delta 2.0 (interval-1 deduplicated to 1 + interval-2 as 1), got %v", result.UsageDelta)
	}
}

func TestCalculateDeltaUsage_BlockedWhenMissingIntervalID(t *testing.T) {
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	samples := []TelemetrySample{
		{
			MeasurementID: "meas-delta-untrusted",
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(1.5),
			IntervalID:    "", // 缺少可信 interval ID
			Quality:       QualityGood,
		},
	}

	result := CalculateDeltaUsage("meas-delta-untrusted", baseTime, baseTime.Add(1*time.Minute), samples)
	if !result.IsUncertain || result.UsageDelta != nil {
		t.Errorf("expected accumulation blocked with uncertain result, got delta: %v, uncertain: %v", result.UsageDelta, result.IsUncertain)
	}
}
