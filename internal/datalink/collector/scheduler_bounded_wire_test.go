package collector

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// This fixture uses the production MC adapter and loopback wire peers. It is not
// physical-PLC, SQL-commit, or throughput evidence.
func TestBoundedSchedulerActualAdapterBackoffYieldsToOtherDevice(t *testing.T) {
	addressA, sendsA := mcBudgetServer(t, "eof")
	addressB, sendsB := mcBudgetServer(t, "unsupported")
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
	config.MaxActiveScans = &one
	config.DefaultRetryCount, config.DefaultRetryDelay = 1, 2*time.Second
	s := NewScheduler(config, manager)
	t.Cleanup(func() { require.NoError(t, s.Stop()) })
	for id, address := range map[string]string{"a": addressA, "b": addressB} {
		host, portText, err := net.SplitHostPort(address)
		require.NoError(t, err)
		port, err := net.DefaultResolver.LookupPort(t.Context(), "tcp", portText)
		require.NoError(t, err)
		raw, err := json.Marshal(schema.ConnectionConfigMC3E{Host: host, Port: port, Timeout: 5})
		require.NoError(t, err)
		s.AddDevice(&schema.Device{ID: id, Protocol: schema.ProtocolMC3E, ConnectionConfig: string(raw)})
		group := id
		s.AddPoint(&schema.Point{ID: id, DeviceID: id, Address: "D0", DataType: schema.DataTypeUint16, PollingGroupID: &group})
	}
	require.NoError(t, s.Start([]*schema.PollingGroup{{ID: "a", Enabled: true, IntervalMs: 10}}))
	require.Eventually(t, func() bool { return sendsA.Load() == 1 }, time.Second, time.Millisecond)
	// A's EOF is handled inside the real adapter's wirebudget.Wait. B becomes due
	// only after A has actually sent, so it cannot win merely by starting first.
	s.AddPollingGroup(&schema.PollingGroup{ID: "b", Enabled: true, IntervalMs: 10})
	require.Eventually(t, func() bool { return sendsB.Load() > 0 }, time.Second, time.Millisecond, "active capacity one must be available during actual adapter backoff")
	manualCtx, cancel := context.WithTimeout(t.Context(), 30*time.Millisecond)
	defer cancel()
	values := s.PollNowContext(manualCtx, []string{"a"})
	require.Len(t, values, 1)
	require.Equal(t, schema.QualityBad, values[0].Quality)
	require.EqualValues(t, 1, sendsA.Load(), "manual operation cannot enter A while scheduled backoff owns its device gate")
	require.NoError(t, s.Stop())
	beforeA, beforeB := sendsA.Load(), sendsB.Load()
	time.Sleep(30 * time.Millisecond)
	require.Equal(t, beforeA, sendsA.Load())
	require.Equal(t, beforeB, sendsB.Load())
}
