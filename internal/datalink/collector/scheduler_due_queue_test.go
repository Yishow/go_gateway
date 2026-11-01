package collector

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func dueSnapshot(device, group, point string) scheduledSnapshot {
	return scheduledSnapshot{device: deviceConfig{ID: device}, groupID: group, points: []pointInfo{{ID: point}}, ctx: context.Background(), capturedAt: time.Now(), cadence: time.Second}
}

func TestDueQueueCoalescesOnlySameTupleAndRejectsWithoutEviction(t *testing.T) {
	q := newDueQueue(2)
	require.True(t, q.offer(dueSnapshot("a", "g1", "old")))
	require.True(t, q.offer(dueSnapshot("a", "g2", "other")))
	require.True(t, q.offer(dueSnapshot("a", "g1", "latest")))
	require.False(t, q.offer(dueSnapshot("b", "g1", "rejected")))
	stats := q.diagnostics()
	require.Equal(t, 2, stats.Pending)
	require.EqualValues(t, 1, stats.Coalesced)
	require.EqualValues(t, 1, stats.Rejected)
	first, ok := q.next()
	require.True(t, ok)
	require.Equal(t, "g1", first.groupID)
	require.Equal(t, "latest", first.points[0].ID)
	_, ok = q.next()
	require.False(t, ok, "same device cannot create another waiting/executing scan")
	q.complete("a")
	second, ok := q.next()
	require.True(t, ok)
	require.Equal(t, "g2", second.groupID)
	require.Equal(t, "other", second.points[0].ID)
}

func TestDueQueueFairDevicesAndStableGroups(t *testing.T) {
	q := newDueQueue(10)
	for _, pair := range [][2]string{{"a", "g1"}, {"a", "g2"}, {"b", "g1"}, {"c", "g1"}} {
		require.True(t, q.offer(dueSnapshot(pair[0], pair[1], pair[1])))
	}
	for _, want := range []string{"a", "b", "c", "a"} {
		snapshot, ok := q.next()
		require.True(t, ok)
		require.Equal(t, want, snapshot.device.ID)
		q.complete(snapshot.device.ID)
	}
	require.Zero(t, q.diagnostics().Pending)
}

func TestDueQueueRunningCaptureUnchangedAndStopDiscardsPending(t *testing.T) {
	q := newDueQueue(2)
	old := dueSnapshot("a", "g", "old")
	require.True(t, q.offer(old))
	running, ok := q.next()
	require.True(t, ok)
	require.True(t, q.offer(dueSnapshot("a", "g", "new")))
	require.True(t, q.offer(dueSnapshot("a", "g", "latest")))
	require.Equal(t, old.capturedAt, running.capturedAt)
	require.Equal(t, "old", running.points[0].ID)
	q.stop()
	q.complete("a")
	_, ok = q.next()
	require.False(t, ok)
	require.False(t, q.offer(old))
	require.Zero(t, q.diagnostics().Pending)
}
