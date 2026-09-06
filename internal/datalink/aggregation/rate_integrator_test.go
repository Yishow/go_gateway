package aggregation

import (
	"math"
	"testing"
	"time"
)

func TestRateIntegrator_R01_PowerToEnergy(t *testing.T) {
	// R01: 2 kW 有效持續 1800 秒 -> 積分估算 1 kWh
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	windowStart := baseTime
	windowEnd := baseTime.Add(1800 * time.Second)

	samples := []TelemetrySample{
		{
			MeasurementID: "meas-kw",
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(2.0),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-kw",
			ObservedAt:    baseTime.Add(1800 * time.Second),
			ValueNumeric:  floatPtr(2.0),
			Quality:       QualityGood,
		},
	}

	config := RateIntegrationConfig{
		InputUnit:      "kW",
		OutputUnit:     "kWh",
		TimeDivisor:    3600.0, // seconds to hours
		MaxGapSeconds:  1900.0,
	}

	result := IntegrateRate("meas-kw", windowStart, windowEnd, samples, config)

	if !result.IsEstimated {
		t.Errorf("expected is_estimated to be true")
	}
	if result.Unit != "kWh" {
		t.Errorf("expected unit kWh, got %s", result.Unit)
	}
	if math.Abs(result.IntegratedValue-1.0) > 1e-6 {
		t.Errorf("expected 1.0 kWh, got %.6f", result.IntegratedValue)
	}
}

func TestRateIntegrator_R02_FlowRateToVolume(t *testing.T) {
	// R02: 120 L/min 有效持續 60 秒 -> 積分估算 120 L
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	windowStart := baseTime
	windowEnd := baseTime.Add(60 * time.Second)

	samples := []TelemetrySample{
		{
			MeasurementID: "meas-flow-rate",
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(120.0),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-flow-rate",
			ObservedAt:    baseTime.Add(60 * time.Second),
			ValueNumeric:  floatPtr(120.0),
			Quality:       QualityGood,
		},
	}

	config := RateIntegrationConfig{
		InputUnit:     "L/min",
		OutputUnit:    "L",
		TimeDivisor:   60.0, // seconds to minutes
		MaxGapSeconds: 100.0,
	}

	result := IntegrateRate("meas-flow-rate", windowStart, windowEnd, samples, config)

	if math.Abs(result.IntegratedValue-120.0) > 1e-6 {
		t.Errorf("expected 120.0 L, got %.6f", result.IntegratedValue)
	}
}
