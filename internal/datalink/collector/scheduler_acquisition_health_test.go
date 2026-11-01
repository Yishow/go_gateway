package collector

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"

	"github.com/stretchr/testify/require"
)

func TestScheduledSafetyCanceledSnapshotCannotUpdateDeviceHealth(t *testing.T) {
	s, snapshot := scheduledPlanFixture(t, `{}`, "D0")
	ctx, cancel := context.WithCancel(t.Context())
	snapshot.ctx = ctx
	cancel()
	s.reportScheduledHealth(snapshot, snapshot.points, s.deviceBreakers["device"], context.Canceled)
	require.Zero(t, s.deviceBreakers["device"].GetStats().TrackerStats.TotalRequests, "caller cancellation is not a device communication failure")
}

func TestScheduledPlanFailurePreservesFinalPointHealthUnit(t *testing.T) {
	for _, addresses := range [][]string{{"D0"}, {"D0", "D10"}} {
		s, snapshot := scheduledPlanFixture(t, `{"acquisition":{"version":1,"profile":"safe-v1","coverage":{"revision":"v1","forbidden":[{"area":"D","start":0,"end":1}]}}}`, addresses...)
		scan, err := buildScheduledPlan(snapshot)
		require.NoError(t, err)
		require.Len(t, scan.failures, 1)
		if len(addresses) == 1 {
			s.pollScheduledSnapshot(snapshot)
		} else {
			p := &scheduledExecutionProtocol{timeout: time.Second}
			conn := &connector.ManagedConnection{DeviceID: snapshot.device.ID, ProtocolType: snapshot.device.Protocol, Config: snapshot.device.Config, Protocol: p}
			s.executeScheduledPlan(t.Context(), conn, scan, "local-plan-failure", s.deviceBreakers["device"])
		}
		require.Len(t, collectedExecutionValues(s), len(addresses))
		stats := s.deviceBreakers["device"].GetStats().TrackerStats
		require.Equal(t, len(addresses), stats.TotalRequests, "local invalid points remain final Point results, not invisible wire units")
		require.Equal(t, 1, stats.FailureCount)
		require.Equal(t, len(addresses)-1, stats.SuccessCount)
	}
}
