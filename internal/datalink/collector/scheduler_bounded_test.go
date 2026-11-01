package collector

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// Controlled protocol isolates scheduling; it is not physical-wire evidence.
type boundedSchedulerProtocol struct {
	connector.Protocol
	slow  bool
	reads *atomic.Int32
}

func (p *boundedSchedulerProtocol) Connect(_ context.Context, config string) error {
	p.slow = strings.Contains(config, "slow")
	return nil
}
func (*boundedSchedulerProtocol) Close() error      { return nil }
func (*boundedSchedulerProtocol) IsConnected() bool { return true }
func (p *boundedSchedulerProtocol) Read(ctx context.Context, _ connector.ReadRequest) (connector.ReadResult, error) {
	if p.slow {
		<-ctx.Done()
		return connector.ReadResult{}, ctx.Err()
	}
	p.reads.Add(1)
	return connector.ReadResult{Quality: schema.QualityGood, Value: 1}, nil
}
func TestBoundedSchedulerSlowDeviceDoesNotBlockGroup(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var fast atomic.Int32
		typ := schema.ProtocolType("bounded-scheduler-progress")
		connector.Register(typ, func() connector.Protocol { return &boundedSchedulerProtocol{reads: &fast} })
		defer connector.Unregister(typ)
		manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
		defer manager.CloseAll()
		config := DefaultSchedulerConfig()
		config.DefaultRetryCount = 0
		s := NewScheduler(config, manager)
		group := "g"
		for _, id := range []string{"slow", "fast"} {
			s.AddDevice(&schema.Device{ID: id, Protocol: typ, ConnectionConfig: `{"mode":"` + id + `"}`})
			s.AddPoint(&schema.Point{ID: id, DeviceID: id, Address: "0", DataType: schema.DataTypeUint16, PollingGroupID: &group})
		}
		require.NoError(t, s.Start([]*schema.PollingGroup{{ID: group, Enabled: true, IntervalMs: 10}}))
		defer func() { require.NoError(t, s.Stop()) }()
		time.Sleep(35 * time.Millisecond)
		synctest.Wait()
		require.GreaterOrEqual(t, fast.Load(), int32(3), "slow A must not prevent B's next group due scans")
	})
}

func TestBoundedSchedulerActiveLimitAndStopPermitWait(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		typ := schema.ProtocolType("bounded-scheduler-active")
		entered := make(chan context.Context, 8)
		connector.Register(typ, func() connector.Protocol { return &shutdownProtocol{entered: entered, cooperate: true} })
		defer connector.Unregister(typ)
		manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
		defer manager.CloseAll()
		config := DefaultSchedulerConfig()
		one := 1
		config.MaxActiveScans = &one
		s := NewScheduler(config, manager)
		group := "g"
		for _, id := range []string{"a", "b"} {
			s.AddDevice(&schema.Device{ID: id, Protocol: typ, ConnectionConfig: "{}"})
			s.AddPoint(&schema.Point{ID: id, DeviceID: id, Address: "0", DataType: schema.DataTypeUint16, PollingGroupID: &group})
		}
		require.NoError(t, s.Start([]*schema.PollingGroup{{ID: group, Enabled: true, IntervalMs: 10}}))
		defer func() { require.NoError(t, s.Stop()) }()
		time.Sleep(15 * time.Millisecond)
		synctest.Wait()
		require.Len(t, entered, 1, "active phase is globally bounded after each device gate")
		require.NoError(t, s.Stop())
		require.Len(t, entered, 1, "canceling permit waiter must not start another read")
	})
}

func TestBoundedSchedulerGateWaitDoesNotOccupyCapacity(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var reads atomic.Int32
		typ := schema.ProtocolType("bounded-scheduler-gate")
		connector.Register(typ, func() connector.Protocol { return &boundedSchedulerProtocol{reads: &reads} })
		defer connector.Unregister(typ)
		manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
		defer manager.CloseAll()
		config := DefaultSchedulerConfig()
		one := 1
		config.MaxActiveScans = &one
		s := NewScheduler(config, manager)
		group := "g"
		for _, id := range []string{"a", "b"} {
			s.AddDevice(&schema.Device{ID: id, Protocol: typ, ConnectionConfig: `{"mode":"fast"}`})
			s.AddPoint(&schema.Point{ID: id, DeviceID: id, Address: "0", DataType: schema.DataTypeUint16, PollingGroupID: &group})
		}
		s.deviceLocks["a"].Lock()
		defer s.deviceLocks["a"].Unlock()
		require.NoError(t, s.Start([]*schema.PollingGroup{{ID: group, Enabled: true, IntervalMs: 10}}))
		defer func() { require.NoError(t, s.Stop()) }()
		time.Sleep(35 * time.Millisecond)
		synctest.Wait()
		require.GreaterOrEqual(t, reads.Load(), int32(3), "A gate wait must leave one active permit for B")
		require.NoError(t, s.Stop())
		require.NoError(t, s.Start(nil))
		time.Sleep(50 * time.Millisecond)
		require.EqualValues(t, 3, reads.Load(), "restart must not resurrect pending due work")
	})
}
