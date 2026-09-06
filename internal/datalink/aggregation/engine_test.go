package aggregation

import (
	"math"
	"testing"
	"time"
)

func TestEngine_LateSampleRevisionIncrement(t *testing.T) {
	// 測試晚到樣本觸發計算修訂版次 (CalculationRevision) 遞增
	engine := NewTelemetryAggregationEngine(EngineConfig{
		CorrectionHorizonHours: 24,
		GracePeriodSeconds:     10,
	})

	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	windowStart := baseTime
	windowEnd := baseTime.Add(60 * time.Second)

	// 第一次輸入 1 筆樣本並結算
	s1 := TelemetrySample{
		MeasurementID: "m1",
		ObservedAt:    baseTime,
		ValueNumeric:  floatPtr(10.0),
		Quality:       QualityGood,
	}
	engine.IngestSample(s1)

	res1, err := engine.ComputeWindowSummary("m1", windowStart, windowEnd)
	if err != nil {
		t.Fatalf("first compute failed: %v", err)
	}
	if res1.CalculationRevision != 1 {
		t.Errorf("expected revision 1, got %d", res1.CalculationRevision)
	}

	// 晚到樣本 (在 24 小時 horizon 內)
	lateSample := TelemetrySample{
		MeasurementID: "m1",
		ObservedAt:    baseTime.Add(30 * time.Second),
		ValueNumeric:  floatPtr(20.0),
		Quality:       QualityGood,
	}
	err = engine.IngestSample(lateSample)
	if err != nil {
		t.Fatalf("failed to ingest late sample: %v", err)
	}

	// 重算
	res2, err := engine.ComputeWindowSummary("m1", windowStart, windowEnd)
	if err != nil {
		t.Fatalf("recompute failed: %v", err)
	}
	if res2.CalculationRevision != 2 {
		t.Errorf("expected revision 2 after late sample, got %d", res2.CalculationRevision)
	}
	// 加權平均: 10*30s + 20*30s = 900 / 60 = 15.0
	if res2.TimeWeightedMean == nil || math.Abs(*res2.TimeWeightedMean-15.0) > 1e-6 {
		t.Errorf("expected updated mean 15.0, got %v", res2.TimeWeightedMean)
	}
}

func TestEngine_P08_ExpiredEvidenceRejection(t *testing.T) {
	// P08: 證據已過期要求重算 -> 拒絕並說明限制
	engine := NewTelemetryAggregationEngine(EngineConfig{
		CorrectionHorizonHours: 1, // 1 小時期限
	})

	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	// 當前時間已前進到 10:00 + 2 小時
	engine.AdvanceWatermark(baseTime.Add(2 * time.Hour))

	// 嘗試注入 2 小時前的過期樣本
	expiredSample := TelemetrySample{
		MeasurementID: "m1",
		ObservedAt:    baseTime,
		ValueNumeric:  floatPtr(10.0),
		Quality:       QualityGood,
	}

	err := engine.IngestSample(expiredSample)
	if err == nil {
		t.Errorf("expected error for expired sample beyond correction horizon")
	}
}
