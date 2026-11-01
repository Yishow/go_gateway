package collector

import (
	"go-gateway/internal/datalink/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDueAdmissionLateOldCaptureCannotReplaceNewPending(t *testing.T) {
	for _, change := range []string{"group", "point"} {
		t.Run(change, func(t *testing.T) {
			s, _ := queuedScheduler(t)
			group := "g"
			point := &schema.Point{ID: "p", DeviceID: "a", RevisionID: "v1", Address: "1", DataType: schema.DataTypeUint16, PollingGroupID: &group}
			s.AddPoint(point)
			// Deterministically order the two capture/enqueue halves around a mutation.
			// No dispatcher is needed to prove which immutable pending entry survives.
			s.running, s.pollCtx, s.due = true, t.Context(), newDueQueue(2)
			oldTicker := &groupTicker{group: &schema.PollingGroup{ID: group, IntervalMs: 10}}
			s.groupTickers[group] = oldTicker
			old := s.captureScheduledGroup(group, oldTicker)
			require.Len(t, old, 1)
			if change == "group" {
				s.groupTickers[group] = &groupTicker{group: &schema.PollingGroup{ID: group, IntervalMs: 20}}
			} else {
				point.RevisionID = "v2"
				s.AddPoint(point)
			}
			latest := s.captureScheduledGroup(group, s.groupTickers[group])
			require.Len(t, latest, 1)
			s.enqueueDueSnapshots(latest)
			s.enqueueDueSnapshots(old)
			got, ok := s.due.next()
			require.True(t, ok)
			require.Equal(t, latest[0].ticker, got.ticker, "late old group capture cannot evict newer admitted work")
			require.Equal(t, latest[0].points[0].observationLease, got.points[0].observationLease, "late old installation cannot evict newer admitted work")
			require.Zero(t, s.DueDiagnostics().Coalesced, "rejected stale capture is not legitimate due coalescing")
		})
	}
}

func TestDueAdmissionCapacityOneRotatesAcrossGroupDevices(t *testing.T) {
	s, _ := queuedScheduler(t)
	typ := s.deviceConfigs["a"].Protocol
	s.AddDevice(&schema.Device{ID: "b", Protocol: typ, ConnectionConfig: "{}"})
	group := "g"
	for _, id := range []string{"a", "b"} {
		s.AddPoint(&schema.Point{ID: id, DeviceID: id, Address: "0", DataType: schema.DataTypeUint16, PollingGroupID: &group})
	}
	s.running, s.pollCtx, s.due = true, t.Context(), newDueQueue(1)
	ticker := &groupTicker{group: &schema.PollingGroup{ID: group, IntervalMs: 10}}
	s.groupTickers[group] = ticker
	served := make([]string, 0, 4)
	for range 4 {
		s.enqueueDueGroup(group, ticker)
		next, ok := s.due.next()
		require.True(t, ok)
		served = append(served, next.device.ID)
		s.due.complete(next.device.ID)
	}
	require.Equal(t, []string{"a", "b", "a", "b"}, served, "capacity rejection must not permanently privilege the first sorted device")
	require.EqualValues(t, 4, s.DueDiagnostics().Rejected)
}
