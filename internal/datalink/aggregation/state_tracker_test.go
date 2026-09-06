package aggregation

import (
	"testing"
	"time"
)

func strPtr(s string) *string {
	return &s
}

func TestStateTracker_S01_InitialSnapshotAndTransition(t *testing.T) {
	// S01: 初筆run，30秒後stop，中間有效 -> initial snapshot及轉態，觀察run 30秒；無虛構先前事件
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	windowStart := baseTime
	windowEnd := baseTime.Add(60 * time.Second)

	samples := []TelemetrySample{
		{
			MeasurementID: "meas-state",
			ObservedAt:    baseTime,
			ValueString:   strPtr("RUN"),
			Quality:       QualityGood,
		},
		{
			MeasurementID: "meas-state",
			ObservedAt:    baseTime.Add(30 * time.Second),
			ValueString:   strPtr("STOP"),
			Quality:       QualityGood,
		},
	}

	config := StateTrackingConfig{
		MaxHoldSeconds: 60,
	}

	result := TrackStateDurations("meas-state", windowStart, windowEnd, samples, config)

	if result.InitialState != "RUN" {
		t.Errorf("expected initial state RUN, got %s", result.InitialState)
	}
	if result.FinalState != "STOP" {
		t.Errorf("expected final state STOP, got %s", result.FinalState)
	}
	if result.TransitionsCount != 1 {
		t.Errorf("expected 1 transition, got %d", result.TransitionsCount)
	}

	// RUN 應持續 30 秒，STOP 應持續 30 秒
	var runDur, stopDur float64
	for _, d := range result.StateDurations {
		if d.State == "RUN" {
			runDur = d.DurationSeconds
		}
		if d.State == "STOP" {
			stopDur = d.DurationSeconds
		}
	}

	if runDur != 30.0 {
		t.Errorf("expected RUN 30s, got %.2f", runDur)
	}
	if stopDur != 30.0 {
		t.Errorf("expected STOP 30s, got %.2f", stopDur)
	}
}

func TestStateTracker_S02_LongDisconnectionBeyondMaxHold(t *testing.T) {
	// S02: run後長斷線 -> 超過max_hold的時間unknown，不能全部算run
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	windowStart := baseTime
	windowEnd := baseTime.Add(60 * time.Second)

	samples := []TelemetrySample{
		{
			MeasurementID: "meas-state",
			ObservedAt:    baseTime,
			ValueString:   strPtr("RUN"),
			Quality:       QualityGood,
		},
	}

	config := StateTrackingConfig{
		MaxHoldSeconds: 15, // RUN 只能 hold 15s，剩下的 45s 為 unknown
	}

	result := TrackStateDurations("meas-state", windowStart, windowEnd, samples, config)

	var runDur float64
	for _, d := range result.StateDurations {
		if d.State == "RUN" {
			runDur = d.DurationSeconds
		}
	}

	if runDur != 15.0 {
		t.Errorf("expected RUN 15s, got %.2f", runDur)
	}
	if result.UnknownDurationSec != 45.0 {
		t.Errorf("expected unknown duration 45s, got %.2f", result.UnknownDurationSec)
	}
}
