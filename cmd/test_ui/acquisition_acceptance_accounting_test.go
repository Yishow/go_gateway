package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/mapping"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type acceptanceDeliveryCut struct {
	native, sink, dropped uint64
	queued                int
	ingestor              datalinkruntime.Stats
	rejected              uint64
	deliveries            map[string]map[string]int
}

// Capture the first coherent boundary after the producer has joined. Retry only
// inconsistent reads; do NOT wait for delivery/queue exhaustion. The single
// original runtime consumer remains running throughout this snapshot.
func acceptanceDeliverySnapshot(t *testing.T, f scheduledProductionFixture, native *acceptanceNativeCapture, drain *acceptanceSampleDrain) acceptanceDeliveryCut {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for {
		interestBefore := f.capture.acceptanceInterest.snapshot()
		before := f.services.runtime.Snapshot()
		sinkBefore, deliveries, sampleErr := drain.deliverySnapshot()
		require.NoError(t, sampleErr)
		rejectedBefore := f.services.runtime.RejectedObservationCount()
		queuedBefore := len(f.services.scheduler.ValueChannel())
		dropped := f.services.scheduler.OutputLossDiagnostics().Total
		queuedAfter := len(f.services.scheduler.ValueChannel())
		rejectedAfter := f.services.runtime.RejectedObservationCount()
		sinkAfter := drain.count()
		after := f.services.runtime.Snapshot()
		interestAfter := f.capture.acceptanceInterest.snapshot()
		produced := native.count()
		if interestBefore == interestAfter && before == after && sinkBefore == sinkAfter && rejectedBefore == rejectedAfter && queuedBefore == queuedAfter && produced == dropped+after.CollectedTotal+uint64(queuedAfter) {
			require.GreaterOrEqual(t, after.CollectedTotal, sinkAfter)
			// This fixture has exactly one simple binding per Point. All entered values
			// must be delivered, rejected, fail point-state persistence, or be the one
			// serial consumer's current pre-SampleSink work. Typed/write errors are
			// separately visible, never reclassified as ValueBuffer drops. WriteError
			// also counts post-SampleSink errors: this is an upper bound, NOT an
			// exact causal partition of unoffered values.
			require.Empty(t, interestAfter.bad)
			require.GreaterOrEqual(t, after.CollectedTotal, sinkAfter+interestAfter.declined)
			unoffered := after.CollectedTotal - sinkAfter - interestAfter.declined
			t.Logf("real_pipeline_interest accepted=%d declined=%d; post_ingestor_ValueEvents_received_checked=%d (loss-tolerant32 buffer; no AcquisitionID; not complete scheduled-value evidence)", interestAfter.accepted, interestAfter.declined, interestAfter.events)
			require.LessOrEqual(t, unoffered, after.PointStateError+after.WriteError+rejectedAfter+1,
				"unexplained ingestion-to-SampleSink gap; do not relabel as queued/drop")
			return acceptanceDeliveryCut{produced, sinkAfter, dropped, queuedAfter, after, rejectedAfter, deliveries}
		}
		if time.Now().After(deadline) {
			t.Fatal("could not obtain a coherent native/queue/ingestor boundary; no fabricated pending count")
		}
		time.Sleep(time.Millisecond)
	}
}

type acceptanceInterestState struct {
	accepted, declined, events uint64
	bad                        string
}
type acceptanceInterestCapture struct {
	mu    sync.Mutex
	state acceptanceInterestState
	want  map[string]map[string]int16
}

func (c *acceptanceInterestCapture) wants(pipe *grouppipeline.Pipeline, device, point, tag string) bool {
	interested := pipe.WantsSample(device, point, tag)
	c.mu.Lock()
	if interested {
		c.state.accepted++
	} else {
		c.state.declined++
	}
	c.mu.Unlock()
	return interested
}
func (c *acceptanceInterestCapture) snapshot() acceptanceInterestState {
	if c == nil {
		return acceptanceInterestState{}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}
func startAcceptanceInterestCapture(t *testing.T, f scheduledProductionFixture, load []acceptanceScheduledDevice) *acceptanceInterestCapture {
	t.Helper()
	c := &acceptanceInterestCapture{want: map[string]map[string]int16{}}
	for _, d := range load {
		c.want[d.id] = d.want
	}
	events, unsubscribe := f.services.runtime.SubscribeValueEvents("", nil)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for event := range events {
			c.mu.Lock()
			want, known := c.want[event.DeviceID][event.PointID]
			if !known || event.RawValue != want || event.Quality != schema.QualityGood || event.Timestamp.IsZero() {
				c.state.bad = fmt.Sprintf("invalid received post-ingestor value device=%s point=%s", event.DeviceID, event.PointID)
			}
			c.state.events++
			c.mu.Unlock()
		}
	}()
	t.Cleanup(func() { unsubscribe(); <-done })
	return c
}

func TestAcquisitionAcceptanceInterestDelegation(t *testing.T) {
	f, _, _ := newAcceptanceScheduledLoad(t, 1, 3, time.Millisecond)
	require.NotNil(t, f.capture.acceptanceInterest)
	// No output groups/cutovers are configured. Compare the actual delegated
	// result, not a hand-authored decision or fake pipeline implementation.
	mappings, err := f.services.mapping.List(t.Context(), mapping.ListFilter{})
	require.NoError(t, err)
	require.Len(t, mappings, 3)
	for _, binding := range mappings {
		expected := f.capture.pipe.WantsSample("batch-device", binding.PointID, binding.TagID)
		require.Equal(t, expected, f.capture.WantsSample("batch-device", binding.PointID, binding.TagID))
		require.False(t, expected)
	}
	f.capture.acceptanceInterest = nil
	require.True(t, f.capture.WantsSample("batch-device", "unmapped-point", "unmapped-tag"), "legacy wiring fixtures retain capture-all semantics")
}
