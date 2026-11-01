package runtime

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestStatusRefreshPreservesFreshReadsAndDerivedTransitions(t *testing.T) {
	svc, counter, points := newStatusCadenceFixture(t)
	deviceID := points[0].DeviceID
	svc.refreshDerivedStatuses()
	require.Zero(t, counter.reads.Load(), "no subscribers need no periodic summary read")
	counter.projection.Version = "v2"
	points[0].LastError = "persistent fault"
	fresh, found, err := svc.DeviceStatus(t.Context(), deviceID)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, int64(1), counter.reads.Load(), "direct API read remains fresh, not last-status cache")
	require.Equal(t, "error", fresh.Status)
	require.Equal(t, string(workspace.RuntimeProjectionAlignmentStale), fresh.ProjectionAlignment)
	stream, cancel := svc.SubscribeStatusEvents(deviceID)
	defer cancel()
	require.Equal(t, fresh, receiveCadenceStatus(t, stream), "single-device subscription reads an immediate fresh initial summary")
	require.Equal(t, int64(2), counter.reads.Load())

	points[0].LastError = ""
	old := time.Now().Add(-2 * time.Hour)
	points[1].LastReadAt = &old
	counter.reads.Store(0)
	svc.refreshDerivedStatuses()
	require.Equal(t, int64(1), counter.reads.Load(), "one periodic refresh shares one workspace projection across device summary")
	event := receiveCadenceStatus(t, stream)
	require.Equal(t, "warning", event.Status)
	require.Equal(t, 1, event.PointsStale)
	snapshot, err := svc.RuntimeStatusSnapshot(t.Context(), deviceID)
	require.NoError(t, err)
	require.Equal(t, snapshot.Collectors[0], event)

	now := time.Now()
	points[1].LastReadAt = &now
	// The existing breaker path fails only against this test's loopback endpoint.
	// Raw health/quality handling is unchanged by summary scheduling.
	for range 3 {
		svc.scheduler.PollNow([]string{points[0].ID})
	}
	svc.refreshDerivedStatuses()
	event = receiveCadenceStatus(t, stream)
	require.Equal(t, "open", event.BreakerState)
	require.Equal(t, "error", event.Status)
	snapshot, err = svc.RuntimeStatusSnapshot(t.Context(), deviceID)
	require.NoError(t, err)
	require.Equal(t, snapshot.Collectors[0], event)
}

func TestStatusRefreshWorkerStopsWithService(t *testing.T) {
	svc, counter, points := newStatusCadenceFixture(t)
	stream, cancel := svc.SubscribeStatusEvents(points[0].DeviceID)
	defer cancel()
	require.Equal(t, "running", receiveCadenceStatus(t, stream).Status)
	points[0].LastError = "persistent fault"
	svc.wg.Go(svc.statusLoop)
	select {
	case event := <-stream:
		require.Equal(t, "error", event.Status)
		require.Equal(t, 1, event.PointsError)
	case <-time.After(3 * time.Second):
		t.Fatal("existing status worker did not publish persistent fault")
	}
	require.NoError(t, svc.Stop(t.Context()))
	require.NoError(t, svc.WaitStopped(t.Context()))
	reads := counter.reads.Load()
	// Observe one whole former tick interval after Stop, not a guessed goroutine count.
	time.Sleep(statusRefreshInterval + 20*time.Millisecond)
	require.Equal(t, reads, counter.reads.Load(), "Stop must finish the original status worker")
	stopped, found, err := svc.DeviceStatus(context.Background(), points[0].DeviceID)
	require.NoError(t, err)
	require.True(t, found)
	require.False(t, stopped.Running)
}

func receiveCadenceStatus(t *testing.T, stream <-chan DeviceStatusEvent) DeviceStatusEvent {
	t.Helper()
	select {
	case event := <-stream:
		return event
	case <-time.After(3 * time.Second):
		t.Fatal("missing periodic status event")
		return DeviceStatusEvent{}
	}
}
