package collector

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type manualContextKey struct{}
type manualContextProtocol struct {
	connector.Protocol
	connectContext context.Context
	readContext    context.Context
	cancelOnRead   context.CancelFunc
}

func (p *manualContextProtocol) Connect(ctx context.Context, _ string) error {
	p.connectContext = ctx
	return nil
}
func (p *manualContextProtocol) Read(ctx context.Context, _ connector.ReadRequest) (connector.ReadResult, error) {
	p.readContext = ctx
	if p.cancelOnRead != nil {
		p.cancelOnRead()
	}
	return connector.ReadResult{Value: 42, Quality: schema.QualityGood}, ctx.Err()
}
func (p *manualContextProtocol) Close() error      { return nil }
func (p *manualContextProtocol) IsConnected() bool { return true }

func TestManualPollPropagatesRequestContext(t *testing.T) {
	for _, mode := range []string{"active", "pre-canceled", "cancel-during-read"} {
		t.Run(mode, func(t *testing.T) {
			protocolType := schema.ProtocolType("manual-context-test")
			protocol := &manualContextProtocol{}
			connector.Register(protocolType, func() connector.Protocol { return protocol })
			defer connector.Unregister(protocolType)
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			defer manager.CloseAll()
			scheduler := NewScheduler(DefaultSchedulerConfig(), manager)
			scheduler.AddDevice(&schema.Device{ID: "context-device", Protocol: protocolType, ConnectionConfig: "{}"})
			scheduler.AddPoint(&schema.Point{ID: "context-point", DeviceID: "context-device", Address: "0", DataType: schema.DataTypeUint16})
			ctx, cancel := context.WithCancel(context.WithValue(t.Context(), manualContextKey{}, "request"))
			defer cancel()
			if mode == "pre-canceled" {
				cancel()
			}
			if mode == "cancel-during-read" {
				protocol.cancelOnRead = cancel
			}
			results := scheduler.PollNowContext(ctx, []string{"context-point"})
			require.Len(t, results, 1)
			if mode == "pre-canceled" {
				require.Nil(t, protocol.connectContext)
				require.Nil(t, protocol.readContext)
				require.Equal(t, schema.QualityBad, results[0].Quality)
				require.Contains(t, results[0].Error, context.Canceled.Error())
				return
			}
			require.NotNil(t, protocol.connectContext)
			require.NotNil(t, protocol.readContext)
			require.Equal(t, "request", protocol.connectContext.Value(manualContextKey{}))
			require.Equal(t, "request", protocol.readContext.Value(manualContextKey{}))
			if mode == "cancel-during-read" {
				require.ErrorIs(t, protocol.readContext.Err(), context.Canceled)
				require.Equal(t, schema.QualityBad, results[0].Quality)
				require.Contains(t, results[0].Error, context.Canceled.Error())
			} else {
				require.NoError(t, protocol.readContext.Err())
				require.Empty(t, results[0].Error)
			}
		})
	}
}

func TestManualCallerContextPreservesNilFallbackAndCallerIdentity(t *testing.T) {
	// A zero context is a supported legacy input, not a new caller convention.
	for _, tc := range []struct {
		name  string
		input context.Context
	}{{name: "legacy zero context"}} {
		t.Run(tc.name, func(t *testing.T) {
			fallback := acquisitionCallerContext(tc.input)
			require.NotNil(t, fallback)
			require.Nil(t, fallback.Done())
			require.NoError(t, fallback.Err())
			_, hasDeadline := fallback.Deadline()
			require.False(t, hasDeadline)
		})
	}
	ctx, cancel := context.WithCancel(context.WithValue(t.Context(), manualContextKey{}, "caller"))
	defer cancel()
	require.Same(t, ctx, acquisitionCallerContext(ctx))
	require.Equal(t, "caller", acquisitionCallerContext(ctx).Value(manualContextKey{}))
	cancel()
	require.Same(t, ctx, acquisitionCallerContext(ctx))
	require.ErrorIs(t, acquisitionCallerContext(ctx).Err(), context.Canceled)
}

func TestAcquisitionDispatchRetainsCapturedCallerContext(t *testing.T) {
	for _, mode := range []string{"legacy-zero", "active", "canceled"} {
		t.Run(mode, func(t *testing.T) {
			typ := schema.ProtocolType(t.Name())
			protocol := &manualContextProtocol{}
			connector.Register(typ, func() connector.Protocol { return protocol })
			defer connector.Unregister(typ)
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			defer manager.CloseAll()
			scheduler := NewScheduler(DefaultSchedulerConfig(), manager)
			scheduler.AddDevice(&schema.Device{ID: "d", Protocol: typ, ConnectionConfig: "{}"})
			scheduler.AddPoint(&schema.Point{ID: "p", DeviceID: "d", Address: "0", DataType: schema.DataTypeUint16})
			ctx, cancel := context.WithCancel(context.WithValue(t.Context(), manualContextKey{}, "captured"))
			defer cancel()
			snapshots := scheduler.captureManualPoints(ctx, []string{"p"})
			require.Len(t, snapshots, 1)
			snapshot := snapshots[0]
			require.Same(t, ctx, snapshot.ctx)
			switch mode {
			case "legacy-zero":
				snapshot.ctx = nil // Exercise the existing zero-snapshot compatibility path.
			case "canceled":
				cancel()
			}
			var values []CollectedValue
			scheduler.pollAcquisitionSnapshot(snapshot.ctx, snapshot, "captured-call", true, func(value CollectedValue) { values = append(values, value) })
			if mode == "canceled" {
				require.Empty(t, values)
				require.Nil(t, protocol.connectContext)
				require.Nil(t, protocol.readContext)
				return
			}
			require.Len(t, values, 1)
			require.Equal(t, schema.QualityGood, values[0].Quality)
			require.NotNil(t, protocol.connectContext)
			require.NotNil(t, protocol.readContext)
			if mode == "active" {
				require.Equal(t, "captured", protocol.connectContext.Value(manualContextKey{}))
				require.Equal(t, "captured", protocol.readContext.Value(manualContextKey{}))
			} else {
				require.Nil(t, protocol.connectContext.Value(manualContextKey{}))
				require.Nil(t, protocol.readContext.Value(manualContextKey{}))
			}
		})
	}
}
