package collector

import (
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestScheduledSafetyBoundedProtocolCannotBypassPlanWithoutPlannedReader(t *testing.T) {
	protocol := &deviceGateProtocol{}
	previous, exists := connector.GetFactory(schema.ProtocolMC3E)
	connector.Register(schema.ProtocolMC3E, func() connector.Protocol { return protocol })
	t.Cleanup(func() {
		if exists {
			connector.Register(schema.ProtocolMC3E, previous)
		} else {
			connector.Unregister(schema.ProtocolMC3E)
		}
	})
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	s, snapshot := scheduledPlanFixture(t, `{"acquisition":{"version":1,"profile":"safe-v1","coverage":{"revision":"v1","forbidden":[{"area":"D","start":0,"end":1}]}}}`, "D0", "D10")
	s.connMgr = manager
	s.pollScheduledSnapshot(snapshot)
	require.Zero(t, protocol.reads.Load(), "an optional interface mismatch cannot bypass planned forbidden ranges, original deadlines or actual-send admission through the scalar executor")
	values := collectedExecutionValues(s)
	require.Len(t, values, 2)
	for _, value := range values {
		require.Equal(t, schema.QualityBad, value.Quality)
		require.Nil(t, value.Value)
	}
}
