package collector

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestManualCapacitySharedAcrossSchedulerLifecycle(t *testing.T) {
	for _, phase := range []string{"before-start", "started", "restarted"} {
		t.Run(phase, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				typ := schema.ProtocolType("manual-capacity-lifecycle")
				entered := make(chan context.Context, 8)
				connector.Register(typ, func() connector.Protocol { return &shutdownProtocol{entered: entered, cooperate: true} })
				defer connector.Unregister(typ)
				manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
				defer manager.CloseAll()
				config := DefaultSchedulerConfig()
				one := 1
				config.MaxActiveScans = &one
				config.DefaultRetryCount = 0
				s := NewScheduler(config, manager)
				defer s.Stop()
				group := "scheduled-b"
				for _, id := range []string{"a", "b"} {
					p := &schema.Point{ID: id, DeviceID: id, Address: "0", DataType: schema.DataTypeUint16}
					if id == "b" {
						p.PollingGroupID = &group
					}
					s.AddDevice(&schema.Device{ID: id, Protocol: typ, ConnectionConfig: "{}"})
					s.AddPoint(p)
				}
				manualCtx, cancel := context.WithCancel(t.Context())
				defer cancel()
				first := make(chan []CollectedValue, 1)
				go func() { first <- s.PollNowContext(manualCtx, []string{"a"}) }()
				synctest.Wait()
				require.Len(t, entered, 1)
				manualRead := <-entered
				if phase == "before-start" {
					second := make(chan []CollectedValue, 1)
					go func() { second <- s.PollNowContext(manualCtx, []string{"b"}) }()
					synctest.Wait()
					require.Empty(t, entered, "manual-before-Start still shares global capacity")
					cancel()
					require.Len(t, <-second, 1)
				} else {
					groups := []*schema.PollingGroup{{ID: group, Enabled: true, IntervalMs: 10}}
					require.NoError(t, s.Start(groups))
					if phase == "restarted" {
						require.NoError(t, s.Stop())
						require.NoError(t, manualRead.Err(), "Stop must not cancel caller-owned manual read")
						require.NoError(t, s.Start(groups))
					}
					time.Sleep(15 * time.Millisecond)
					synctest.Wait()
					require.Empty(t, entered, "Start/restart must not split the capacity pool from an in-flight manual read")
					require.NoError(t, s.Stop(), "Stop joins scheduled permit waiters, not unrelated manual calls")
					require.NoError(t, manualRead.Err())
					cancel()
				}
				values := <-first
				require.Len(t, values, 1)
				require.Equal(t, schema.QualityBad, values[0].Quality)
				require.Contains(t, values[0].Error, context.Canceled.Error())
				require.Empty(t, collectedExecutionValues(s), "manual result has caller delivery only")
			})
		})
	}
}

func TestManualLegacyProtocolRetainsSingleSchedulerInvocation(t *testing.T) {
	typ := schema.ProtocolType("manual-legacy-single-invocation")
	protocol := &retryFactsProtocol{}
	connector.Register(typ, func() connector.Protocol { return protocol })
	defer connector.Unregister(typ)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	defer manager.CloseAll()
	config := DefaultSchedulerConfig()
	config.DefaultRetryCount, config.DefaultRetryDelay = 3, 0
	s := NewScheduler(config, manager)
	s.AddDevice(&schema.Device{ID: "device", Protocol: typ, ConnectionConfig: "{}"})
	s.AddPoint(&schema.Point{ID: "point", DeviceID: "device", Address: "0", DataType: schema.DataTypeUint64})
	values := s.PollNowContext(t.Context(), []string{"point"})
	require.Len(t, values, 1)
	require.Equal(t, schema.QualityBad, values[0].Quality)
	require.EqualValues(t, 1, protocol.calls.Load(), "migration must not add scheduler retries to a custom legacy manual read")
}
