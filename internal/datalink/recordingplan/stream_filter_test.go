package recordingplan

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

func TestStreamFilter_EverySample(t *testing.T) {
	filter := NewStreamFilter()
	stream := PlanStream{
		StreamID:      "s-raw",
		MeasurementID: "meas-1",
		Mode:          StreamModeRawHistory,
		RawPolicy:     RawPolicyEverySample,
	}

	now := time.Now()
	for i := 0; i < 5; i++ {
		sample := measurement.SampleEnvelope{
			SampleID:      "samp",
			MeasurementID: "meas-1",
			ObservedAt:    now.Add(time.Duration(i) * time.Second),
			Value:         10.0,
			Quality:       schema.QualityGood,
		}
		emit, reason := filter.ShouldEmitSample(stream, sample)
		if !emit {
			t.Errorf("iteration %d: expected emit for every_sample, got false (%s)", i, reason)
		}
	}
}

func TestStreamFilter_OnChange_DeadbandAndHeartbeat(t *testing.T) {
	filter := NewStreamFilter()
	deadband := 0.5
	heartbeat := 10
	stream := PlanStream{
		StreamID:            "s-change",
		MeasurementID:       "meas-1",
		Mode:                StreamModeRawHistory,
		RawPolicy:           RawPolicyOnChange,
		OnChangeDeadband:    &deadband,
		MaxHeartbeatSeconds: &heartbeat,
	}

	t0 := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)

	// 1. First sample -> must emit
	s1 := measurement.SampleEnvelope{
		MeasurementID: "meas-1",
		ObservedAt:    t0,
		Value:         100.0,
		Quality:       schema.QualityGood,
	}
	emit, reason := filter.ShouldEmitSample(stream, s1)
	if !emit || reason != "initial_sample" {
		t.Errorf("expected initial_sample emit, got %v (%s)", emit, reason)
	}

	// 2. Small change (100.2, delta 0.2 <= 0.5) 2 seconds later -> suppressed
	s2 := measurement.SampleEnvelope{
		MeasurementID: "meas-1",
		ObservedAt:    t0.Add(2 * time.Second),
		Value:         100.2,
		Quality:       schema.QualityGood,
	}
	emit, reason = filter.ShouldEmitSample(stream, s2)
	if emit {
		t.Errorf("expected suppression for small delta <= deadband, got emit (%s)", reason)
	}

	// 3. Significant change (100.8, delta 0.8 > 0.5) 4 seconds later -> emit
	s3 := measurement.SampleEnvelope{
		MeasurementID: "meas-1",
		ObservedAt:    t0.Add(4 * time.Second),
		Value:         100.8,
		Quality:       schema.QualityGood,
	}
	emit, reason = filter.ShouldEmitSample(stream, s3)
	if !emit || reason != "deadband_exceeded" {
		t.Errorf("expected deadband_exceeded emit, got %v (%s)", emit, reason)
	}

	// 4. Quality change (Good -> Bad) with same value -> must emit immediately!
	s4 := measurement.SampleEnvelope{
		MeasurementID: "meas-1",
		ObservedAt:    t0.Add(6 * time.Second),
		Value:         100.8,
		Quality:       schema.QualityBad,
	}
	emit, reason = filter.ShouldEmitSample(stream, s4)
	if !emit || reason != "quality_change" {
		t.Errorf("expected quality_change emit, got %v (%s)", emit, reason)
	}

	// 5. Heartbeat exceeded (15s after last emit with no value change) -> emit
	s5 := measurement.SampleEnvelope{
		MeasurementID: "meas-1",
		ObservedAt:    t0.Add(25 * time.Second),
		Value:         100.8,
		Quality:       schema.QualityBad,
	}
	emit, reason = filter.ShouldEmitSample(stream, s5)
	if !emit || reason != "heartbeat" {
		t.Errorf("expected heartbeat emit, got %v (%s)", emit, reason)
	}
}

func TestStreamFilter_LatestOnly_TimestampGuard(t *testing.T) {
	filter := NewStreamFilter()
	stream := PlanStream{
		StreamID:      "s-latest",
		MeasurementID: "meas-1",
		Mode:          StreamModeLatestOnly,
	}

	t0 := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)

	// Sample at t0
	s1 := measurement.SampleEnvelope{
		MeasurementID: "meas-1",
		ObservedAt:    t0,
		Value:         10,
		Quality:       schema.QualityGood,
	}
	emit, _ := filter.ShouldEmitSample(stream, s1)
	if !emit {
		t.Errorf("expected emit for initial latest_only")
	}

	// Late-arriving sample at t0 - 5s -> should NOT overwrite newer data!
	sLate := measurement.SampleEnvelope{
		MeasurementID: "meas-1",
		ObservedAt:    t0.Add(-5 * time.Second),
		Value:         9,
		Quality:       schema.QualityGood,
	}
	emit, reason := filter.ShouldEmitSample(stream, sLate)
	if emit {
		t.Errorf("expected rejection of late-arriving sample for latest_only, got %s", reason)
	}

	// Newer sample at t0 + 5s -> emit
	sNew := measurement.SampleEnvelope{
		MeasurementID: "meas-1",
		ObservedAt:    t0.Add(5 * time.Second),
		Value:         12,
		Quality:       schema.QualityGood,
	}
	emit, _ = filter.ShouldEmitSample(stream, sNew)
	if !emit {
		t.Errorf("expected emit for newer sample in latest_only")
	}
}

func TestStreamFilter_BatchSnapshot_HeldTriggerDeduplication(t *testing.T) {
	filter := NewStreamFilter()
	triggerMember := "meas-trigger"
	stream := PlanStream{
		StreamID:             "s-batch",
		MeasurementID:        "meas-val",
		Mode:                 StreamModeBatchSnapshot,
		BatchTriggerMemberID: &triggerMember,
	}

	t0 := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)

	// Trigger goes high (0 -> 1): Rising edge -> trigger batch
	triggerSample1 := measurement.SampleEnvelope{
		MeasurementID: "meas-trigger",
		ObservedAt:    t0,
		Value:         1,
		Quality:       schema.QualityGood,
	}
	isTriggered := filter.EvaluateBatchTrigger(stream, triggerSample1)
	if !isTriggered {
		t.Errorf("expected rising edge trigger to activate batch")
	}

	// Trigger held high (1 -> 1) on next poll: Should NOT trigger duplicate batch
	triggerSample2 := measurement.SampleEnvelope{
		MeasurementID: "meas-trigger",
		ObservedAt:    t0.Add(1 * time.Second),
		Value:         1,
		Quality:       schema.QualityGood,
	}
	isTriggeredAgain := filter.EvaluateBatchTrigger(stream, triggerSample2)
	if isTriggeredAgain {
		t.Errorf("expected held-high trigger to NOT re-trigger batch")
	}

	// Trigger drops low (1 -> 0)
	triggerSample3 := measurement.SampleEnvelope{
		MeasurementID: "meas-trigger",
		ObservedAt:    t0.Add(2 * time.Second),
		Value:         0,
		Quality:       schema.QualityGood,
	}
	filter.EvaluateBatchTrigger(stream, triggerSample3)

	// Next rising edge (0 -> 1) -> triggers new batch!
	triggerSample4 := measurement.SampleEnvelope{
		MeasurementID: "meas-trigger",
		ObservedAt:    t0.Add(3 * time.Second),
		Value:         1,
		Quality:       schema.QualityGood,
	}
	if !filter.EvaluateBatchTrigger(stream, triggerSample4) {
		t.Errorf("expected next rising edge to trigger new batch")
	}
}
