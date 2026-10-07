package collector

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestSchedulerShutdownRepeatWaitsForOriginalWorkers(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
	require.NoError(t, s.Start(nil))
	release := make(chan struct{})
	unblock := sync.OnceFunc(func() { close(release) })
	defer unblock()
	s.wg.Go(func() { <-release })
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	require.ErrorIs(t, s.StopContext(ctx), context.Canceled)
	require.ErrorIs(t, s.StopContext(ctx), context.Canceled, "repeat stop must not report a live worker as finished")
	require.Error(t, s.Start(nil), "unfinished stop must reject restart")
	waiter, ok := any(s).(interface{ WaitStopped(context.Context) error })
	require.True(t, ok)
	unblock()
	require.NoError(t, waiter.WaitStopped(t.Context()))
	require.NoError(t, s.Start(nil))
	require.NoError(t, s.Stop())
}

type shutdownProtocol struct {
	connector.Protocol
	entered   chan context.Context
	release   chan struct{}
	cooperate bool
	reads     atomic.Int32
}

func (*shutdownProtocol) Connect(context.Context, string) error { return nil }
func (*shutdownProtocol) Close() error                          { return nil }
func (*shutdownProtocol) IsConnected() bool                     { return true }
func (p *shutdownProtocol) Read(ctx context.Context, _ connector.ReadRequest) (connector.ReadResult, error) {
	p.reads.Add(1)
	p.entered <- ctx
	if p.cooperate {
		<-ctx.Done()
	} else {
		<-p.release
	}
	return connector.ReadResult{}, ctx.Err()
}

func TestSchedulerShutdownCancelsPollingAndStopsRetries(t *testing.T) {
	for _, cooperate := range []bool{true, false} {
		t.Run(map[bool]string{true: "cooperative", false: "late-completion"}[cooperate], func(t *testing.T) {
			protocolType := schema.ProtocolType("shutdown-lifecycle")
			protocol := &shutdownProtocol{entered: make(chan context.Context, 1), release: make(chan struct{}), cooperate: cooperate}
			release := sync.OnceFunc(func() { close(protocol.release) })
			defer release()
			connector.Register(protocolType, func() connector.Protocol { return protocol })
			defer connector.Unregister(protocolType)
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			defer manager.CloseAll()
			scheduler := NewScheduler(DefaultSchedulerConfig(), manager)
			scheduler.AddDevice(&schema.Device{ID: "d", Protocol: protocolType, ConnectionConfig: "{}"})
			groupID := "g"
			scheduler.AddPoint(&schema.Point{ID: "p", DeviceID: "d", Address: "0", DataType: schema.DataTypeUint16, PollingGroupID: &groupID})
			require.NoError(t, scheduler.Start([]*schema.PollingGroup{{ID: groupID, Enabled: true, IntervalMs: 1}}))
			var pollCtx context.Context
			select {
			case pollCtx = <-protocol.entered:
			case <-time.After(time.Second):
				t.Fatal("poll did not reach protocol")
			}
			ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
			defer cancel()
			err := scheduler.StopContext(ctx)
			if cooperate {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, context.DeadlineExceeded)
				require.ErrorIs(t, scheduler.WaitStopped(ctx), context.DeadlineExceeded)
				require.Error(t, scheduler.Start(nil))
			}
			require.ErrorIs(t, pollCtx.Err(), context.Canceled)
			release()
			require.NoError(t, scheduler.WaitStopped(t.Context()))
			require.EqualValues(t, 1, protocol.reads.Load(), "no retry or new poll starts after stop")
		})
	}
}
