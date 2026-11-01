package collector

import (
	"context"
	"runtime"
	"testing"
	"testing/synctest"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type queuedIdentityProtocol struct {
	connector.Protocol
	sent chan string
}

func (*queuedIdentityProtocol) Connect(context.Context, string) error { return nil }
func (*queuedIdentityProtocol) Close() error                          { return nil }
func (*queuedIdentityProtocol) IsConnected() bool                     { return true }
func (p *queuedIdentityProtocol) Read(_ context.Context, request connector.ReadRequest) (connector.ReadResult, error) {
	p.sent <- request.Address
	return connector.ReadResult{Quality: schema.QualityGood, Value: request.Address}, nil
}

func queuedScheduler(t *testing.T) (scheduler *Scheduler, sent chan string) {
	t.Helper()
	typ := schema.ProtocolType("queued-identity")
	sent = make(chan string, 20)
	connector.Register(typ, func() connector.Protocol { return &queuedIdentityProtocol{sent: sent} })
	t.Cleanup(func() { connector.Unregister(typ) })
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	config := DefaultSchedulerConfig()
	capacity := 2
	config.PendingScanCapacity = &capacity
	s := NewScheduler(config, manager)
	s.AddDevice(&schema.Device{ID: "a", Protocol: typ, ConnectionConfig: "{}"})
	return s, sent
}

func TestBoundedSchedulerPendingRevisionCannotSendOldPoint(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, sent := queuedScheduler(t)
		group := "g"
		point := &schema.Point{ID: "p", DeviceID: "a", RevisionID: "v1", Address: "1", DataType: schema.DataTypeUint16, PollingGroupID: &group}
		s.AddPoint(point)
		s.deviceLocks["a"].Lock()
		require.NoError(t, s.Start([]*schema.PollingGroup{{ID: group, Enabled: true, IntervalMs: 10}}))
		defer func() { require.NoError(t, s.Stop()) }()
		time.Sleep(15 * time.Millisecond)
		synctest.Wait()
		point.RevisionID, point.Address = "v2", "2"
		s.AddPoint(point)
		time.Sleep(10 * time.Millisecond)
		synctest.Wait()
		require.Equal(t, 1, s.DueDiagnostics().Pending)
		s.deviceLocks["a"].Unlock()
		synctest.Wait()
		require.Len(t, sent, 1)
		require.Equal(t, "2", <-sent)
		values := collectedExecutionValues(s)
		require.Len(t, values, 1)
		require.Equal(t, "v2", values[0].PointRevision)
	})
}

func TestBoundedSchedulerGroupsStayBoundedWhileDeviceBusy(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, sent := queuedScheduler(t)
		for _, id := range []string{"g1", "g2"} {
			group := id
			s.AddPoint(&schema.Point{ID: id, DeviceID: "a", RevisionID: "v1", Address: id, DataType: schema.DataTypeUint16, PollingGroupID: &group})
		}
		before := runtime.NumGoroutine()
		s.deviceLocks["a"].Lock()
		require.NoError(t, s.Start([]*schema.PollingGroup{{ID: "g1", Enabled: true, IntervalMs: 10}, {ID: "g2", Enabled: true, IntervalMs: 15}}))
		defer func() { require.NoError(t, s.Stop()) }()
		time.Sleep(95 * time.Millisecond)
		synctest.Wait()
		stats := s.DueDiagnostics()
		require.Equal(t, 2, stats.Pending)
		require.Equal(t, 1, stats.Running)
		require.Positive(t, stats.Coalesced)
		require.LessOrEqual(t, runtime.NumGoroutine()-before, 6, "ticks cannot create scan goroutines per pending tuple")
		s.deviceLocks["a"].Unlock()
		synctest.Wait()
		got := map[string]bool{}
		for len(sent) > 0 {
			got[<-sent] = true
		}
		require.Equal(t, map[string]bool{"g1": true, "g2": true}, got)
		require.Zero(t, s.DueDiagnostics().Pending)
		require.Zero(t, s.DueDiagnostics().Running)
		require.Positive(t, s.DueDiagnostics().Overruns, "device-gate wait must remain visible as a measured scan overrun")
	})
}
