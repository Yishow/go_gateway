package collector

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/mcprotocol"

	"github.com/stretchr/testify/require"
)

func executionManagedConnection(t *testing.T, s *Scheduler, scan scheduledPlan, protocol *scheduledExecutionProtocol) *connector.ManagedConnection {
	t.Helper()
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
	s.connMgr = manager
	conn, err := manager.GetOrCreate(t.Context(), scan.snapshot.device.ID, scan.snapshot.device.Protocol, scan.snapshot.device.Config)
	require.NoError(t, err)
	return conn
}

func TestScheduledExecutionDowngradeRequiresCurrentFinalConnectionToken(t *testing.T) {
	s, scan, p, _ := scheduledExecutionFixture(t, `{"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`, "D0", "D10", "D20")
	conn := executionManagedConnection(t, s, scan, p)
	initial, err := conn.CurrentAcquisitionToken(t.Context())
	require.NoError(t, err)
	require.True(t, initial.Valid())
	p.hook = func(_ context.Context, req connector.PlannedReadRequest, send int) ([]connector.PlannedPointResult, error) {
		r, _ := req.Plan.Request(req.RequestIndex)
		if r.Kind == "random" {
			err := &mcprotocol.MCError{Code: 0xC059}
			return executionResults(req, time.Unix(int64(send), 0), err), err
		}
		return executionResults(req, time.Unix(int64(send), 0), nil), nil
	}
	s.executeScheduledPlan(t.Context(), conn, scan, "first", nil)
	require.Len(t, collectedExecutionValues(s), 3)
	require.EqualValues(t, 4, p.sends.Load())
	require.Equal(t, initial, s.acquisitionDowngrades["device"].token)
	s.executeScheduledPlan(t.Context(), conn, scan, "cached", nil)
	values := collectedExecutionValues(s)
	require.Len(t, values, 3)
	for _, value := range values {
		require.Contains(t, value.RequestID, "/safe/")
	}
	require.EqualValues(t, 7, p.sends.Load())
	// Simulate the provider's next successful publication. Adapter publication
	// itself has separate real-connect tests; this collector must only consume it.
	p.generation.Add(1)
	s.executeScheduledPlan(t.Context(), conn, scan, "reconnected", nil)
	require.Len(t, collectedExecutionValues(s), 3)
	require.EqualValues(t, 11, p.sends.Load())
	newer := s.acquisitionDowngrades["device"]
	require.Greater(t, newer.token.TransportGeneration, initial.TransportGeneration)
	s.rememberScheduledDowngrade(t.Context(), conn, scan, initial)
	require.Equal(t, newer, s.acquisitionDowngrades["device"], "late old token cannot overwrite the new downgrade")
}

func TestScheduledExecutionPostSendTokenChangeCannotClaimNewDowngrade(t *testing.T) {
	s, scan, p, _ := scheduledExecutionFixture(t, `{"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`, "D0", "D10")
	conn := executionManagedConnection(t, s, scan, p)
	p.hook = func(_ context.Context, req connector.PlannedReadRequest, send int) ([]connector.PlannedPointResult, error) {
		r, _ := req.Plan.Request(req.RequestIndex)
		if r.Kind == "random" {
			p.generation.Add(1)
			err := &mcprotocol.MCError{Code: 0xC059}
			return executionResults(req, time.Unix(int64(send), 0), err), err
		}
		return executionResults(req, time.Unix(int64(send), 0), nil), nil
	}
	s.executeScheduledPlan(t.Context(), conn, scan, "changed-after-send", nil)
	require.Len(t, collectedExecutionValues(s), 2)
	require.EqualValues(t, 3, p.sends.Load(), "one scan may still fall back safely")
	require.Empty(t, s.acquisitionDowngrades, "pre-read/current token cannot relabel the rejected final wire attempt")
}
