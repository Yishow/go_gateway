package collector

import (
	"encoding/json"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// Minimal real-adapter all-timeout pressure case. The existing independent MC
// frame peer consumes complete requests but never replies. No fake-clock elapsed
// time, fake successes or physical-PLC timing is reported as performance.
func TestAcquisitionAcceptanceAllTimeoutQueuePressure(t *testing.T) {
	previous, existed := connector.GetFactory(schema.ProtocolMC3E)
	connector.Register(schema.ProtocolMC3E, adapters.NewMC3EConnector)
	t.Cleanup(func() {
		if existed {
			connector.Register(schema.ProtocolMC3E, previous)
		} else {
			connector.Unregister(schema.ProtocolMC3E)
		}
	})
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	config := DefaultSchedulerConfig()
	one := 1
	config.MaxActiveScans, config.PendingScanCapacity = &one, &one
	config.DefaultRetryCount, config.DefaultRetryDelay = 0, 0
	s := NewScheduler(config, manager)
	counters := map[string]*atomic.Int32{}
	group := "acceptance-timeout"
	for _, id := range []string{"a", "b"} {
		address, sends := mcBudgetServer(t, "stall")
		counters[id] = sends
		host, portText, err := net.SplitHostPort(address)
		require.NoError(t, err)
		port, err := net.DefaultResolver.LookupPort(t.Context(), "tcp", portText)
		require.NoError(t, err)
		raw, err := json.Marshal(schema.ConnectionConfigMC3E{Host: host, Port: port, Timeout: 1})
		require.NoError(t, err)
		s.AddDevice(&schema.Device{ID: id, Protocol: schema.ProtocolMC3E, ConnectionConfig: string(raw)})
		s.AddPoint(&schema.Point{ID: id, DeviceID: id, Address: "D0", DataType: schema.DataTypeUint16, PollingGroupID: &group})
	}
	t.Cleanup(func() {
		require.NoError(t, s.Stop())
		require.NoError(t, manager.CloseAll())
	})
	started := time.Now()
	require.NoError(t, s.Start([]*schema.PollingGroup{{ID: group, Enabled: true, IntervalMs: 10}}))
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	deadline := time.NewTimer(8 * time.Second)
	defer deadline.Stop()
	failed := map[string]bool{}
	peakPending := 0
	for len(failed) < 2 {
		select {
		case value := <-s.ValueChannel():
			require.Equal(t, schema.QualityBad, value.Quality)
			failed[value.DeviceID] = true
		case <-ticker.C:
			stats := s.DueDiagnostics()
			peakPending = max(peakPending, stats.Pending)
			require.LessOrEqual(t, stats.Pending, 1)
		case <-deadline.C:
			t.Fatal("both stalled devices must receive bounded timeout service")
		}
	}
	stats := s.DueDiagnostics()
	require.Positive(t, stats.Coalesced)
	require.Positive(t, stats.Rejected)
	require.Equal(t, 1, peakPending)
	require.Zero(t, s.OutputLossDiagnostics().Total, "due rejection is not output loss")
	for id, sends := range counters {
		require.Positive(t, sends.Load(), id)
	}
	require.NoError(t, s.Stop())
	require.NoError(t, manager.CloseAll())
	beforeA, beforeB := counters["a"].Load(), counters["b"].Load()
	time.Sleep(30 * time.Millisecond)
	require.Equal(t, beforeA, counters["a"].Load())
	require.Equal(t, beforeB, counters["b"].Load())
	for _, value := range collectedExecutionValues(s) {
		require.Equal(t, schema.QualityBad, value.Quality, "no late Good after Stop")
	}
	t.Logf("all-timeout real-adapter evidence: elapsed=%s sends_a=%d sends_b=%d due=%+v output_loss=%+v", time.Since(started), beforeA, beforeB, stats, s.OutputLossDiagnostics())
}
