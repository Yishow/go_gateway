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
}

func (p *manualContextProtocol) Connect(ctx context.Context, _ string) error {
	p.connectContext = ctx
	return nil
}
func (p *manualContextProtocol) Read(ctx context.Context, _ connector.ReadRequest) (connector.ReadResult, error) {
	p.readContext = ctx
	return connector.ReadResult{Value: 42, Quality: schema.QualityGood}, ctx.Err()
}
func (p *manualContextProtocol) Close() error      { return nil }
func (p *manualContextProtocol) IsConnected() bool { return true }

func TestManualPollPropagatesRequestContext(t *testing.T) {
	for _, canceled := range []bool{false, true} {
		t.Run(map[bool]string{false: "active", true: "canceled"}[canceled], func(t *testing.T) {
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
			if canceled {
				cancel()
			}
			results := scheduler.PollNowContext(ctx, []string{"context-point"})
			require.Len(t, results, 1)
			require.NotNil(t, protocol.connectContext)
			require.NotNil(t, protocol.readContext)
			require.Equal(t, "request", protocol.connectContext.Value(manualContextKey{}))
			require.Equal(t, "request", protocol.readContext.Value(manualContextKey{}))
			if canceled {
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
