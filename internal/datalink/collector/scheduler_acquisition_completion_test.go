package collector

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionCompletionPublishedDetachedAndUnlocked(t *testing.T) {
	s, _ := previewExecutionFixture(t, `{}`, "D0", "D1")
	var events []AcquisitionCompletionObservation
	s.config.AcquisitionCompletionObserver = func(o AcquisitionCompletionObservation) {
		// A scheduler read must be safe: publication mutex is already released.
		require.True(t, s.CaptureAcquisitionGeometry("device").authorized)
		events = append(events, o)
	}
	s.pollAcquisitionSnapshot(s.captureScheduledGroup("source-cadence")[0], "scan", false, func(CollectedValue) {})
	require.Len(t, events, 1)
	event := events[0]
	require.Equal(t, "w", event.WorkspaceID)
	require.Equal(t, "source", event.SourceRevision)
	require.Equal(t, "device", event.DeviceID)
	require.Equal(t, "source-cadence", event.PollingGroupID)
	require.Equal(t, s.InstallationScope(), event.InstallationScope)
	require.Equal(t, 2, event.PointCount)
	require.Equal(t, "scan", event.Metrics.AcquisitionID)
	original := s.acquisitionEvidence[acquisitionEvidenceTuple{"device", "source-cadence"}]
	require.Equal(t, original.head, event.IdentityHead)
	require.Equal(t, original.metrics.CycleLatencyMS, event.Metrics.CycleLatencyMS, "copy the original due-to-finish duration, never recalculate")
	require.Equal(t, original.metrics.ObservedAt, event.Metrics.ObservedAt)
	require.NotNil(t, event.Metrics.LastGoodAgeMS)
	require.GreaterOrEqual(t, *event.Metrics.LastGoodAgeMS, float64(0))
	event.Metrics.Phases[0].Phase = "mutated"
	event.Metrics.LogicalRequests.Succeeded = 999
	*event.Metrics.LastGoodAgeMS = -999
	view, _, err := s.PreviewAcquisition(t.Context(), "w", "device", "source")
	require.NoError(t, err)
	require.Equal(t, "initial", view.Groups[0].Actual.Phases[0].Phase)
	require.EqualValues(t, 1, view.Groups[0].Actual.LogicalRequests.Succeeded)
	require.GreaterOrEqual(t, *view.Groups[0].Actual.LastGoodAgeMS, float64(0))
}

func TestAcquisitionCompletionNilPanicAndRejected(t *testing.T) {
	t.Run("nil-no-clone", func(t *testing.T) {
		s, _ := previewExecutionFixture(t, `{}`, "D0")
		require.Nil(t, DefaultSchedulerConfig().AcquisitionCompletionObserver)
		require.Zero(t, testing.AllocsPerRun(10, func() { s.observeAcquisitionCompletion(nil) }))
		s.pollAcquisitionSnapshot(s.captureScheduledGroup("source-cadence")[0], "nil", false, func(CollectedValue) {})
		original := s.acquisitionEvidence[acquisitionEvidenceTuple{"device", "source-cadence"}]
		require.Equal(t, "nil", original.metrics.AcquisitionID)
		require.Zero(t, testing.AllocsPerRun(10, func() { s.observeAcquisitionCompletion(original) }))
	})
	t.Run("panic-keeps-publication", func(t *testing.T) {
		s, _ := previewExecutionFixture(t, `{}`, "D0")
		var calls int
		s.config.AcquisitionCompletionObserver = func(AcquisitionCompletionObservation) { calls++; panic("diagnostic failure") }
		require.NotPanics(t, func() {
			s.pollAcquisitionSnapshot(s.captureScheduledGroup("source-cadence")[0], "panic", false, func(CollectedValue) {})
		})
		require.Equal(t, 1, calls)
		require.Equal(t, "panic", s.acquisitionEvidence[acquisitionEvidenceTuple{"device", "source-cadence"}].metrics.AcquisitionID)
	})
	t.Run("manual-late-stale", func(t *testing.T) {
		s, _ := previewExecutionFixture(t, `{}`, "D0")
		var ids []string
		s.config.AcquisitionCompletionObserver = func(o AcquisitionCompletionObservation) { ids = append(ids, o.Metrics.AcquisitionID) }
		snapshot := s.captureScheduledGroup("source-cadence")[0]
		s.pollAcquisitionSnapshot(snapshot, "manual", true, func(CollectedValue) {})
		require.Empty(t, ids)
		scan, err := buildScheduledPlan(snapshot)
		require.NoError(t, err)
		_, old := s.beginAcquisitionEvidence(t.Context(), scan, "old")
		_, newer := s.beginAcquisitionEvidence(t.Context(), scan, "new")
		newer.snapshot.capturedAt = old.snapshot.capturedAt.Add(time.Nanosecond)
		s.finishAcquisitionEvidence(newer)
		s.finishAcquisitionEvidence(old)
		s.finishAcquisitionEvidence(newer)
		require.Equal(t, []string{"new"}, ids)
		s.RemovePoint("point-0")
		_, stale := s.beginAcquisitionEvidence(t.Context(), scan, "stale")
		stale.snapshot.capturedAt = newer.snapshot.capturedAt.Add(time.Nanosecond)
		s.finishAcquisitionEvidence(stale)
		require.Equal(t, []string{"new"}, ids)
	})
}

func TestAcquisitionCompletionConcurrentStopFact(t *testing.T) {
	s, _ := previewExecutionFixture(t, `{}`, "D0")
	var calls atomic.Int64
	events := make(chan AcquisitionCompletionObservation, 64)
	s.config.AcquisitionCompletionObserver = func(o AcquisitionCompletionObservation) {
		_ = s.IsRunning() // legal read after publication unlock, not callback-time authority
		calls.Add(1)
		select {
		case events <- o:
		default:
		}
	}
	require.NoError(t, s.Start(nil))
	t.Cleanup(func() { require.NoError(t, s.Stop()) })
	snapshot := s.captureScheduledGroup("source-cadence")[0]
	scan, err := buildScheduledPlan(snapshot)
	require.NoError(t, err)
	_, first := s.beginAcquisitionEvidence(t.Context(), scan, "first")
	s.finishAcquisitionEvidence(first)
	require.EqualValues(t, 1, calls.Load())
	var workers sync.WaitGroup
	workers.Add(1)
	go func() {
		defer workers.Done()
		for i := range 20 {
			e := *first
			e.metrics = first.metrics
			e.metrics.AcquisitionID = fmt.Sprintf("next-%d", i)
			e.snapshot.capturedAt = first.snapshot.capturedAt.Add(time.Duration(i+1) * time.Nanosecond)
			s.finishAcquisitionEvidence(&e)
		}
	}()
	require.NoError(t, s.Stop())
	workers.Wait()
	close(events)
	require.Empty(t, s.acquisitionEvidence)
	for event := range events {
		require.Equal(t, "device", event.DeviceID)
		require.NotEmpty(t, event.IdentityHead)
		require.NotZero(t, event.Metrics.ObservedAt)
	}
	// These are immutable publication facts even if Stop completed before a
	// consumer reads them. They cannot authorize the new/stopped source scope.
	before := calls.Load()
	s.finishAcquisitionEvidence(first)
	require.Equal(t, before, calls.Load(), "late stopped completion must not publish or notify")
}
