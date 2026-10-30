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
