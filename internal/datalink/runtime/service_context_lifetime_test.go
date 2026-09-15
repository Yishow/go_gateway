package runtime

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"

	"github.com/stretchr/testify/require"
)

type runtimeContextKey struct{}
type runtimeContextProtocol struct{ connector.Protocol }

func (*runtimeContextProtocol) Connect(context.Context, string) error { return nil }
func (*runtimeContextProtocol) Close() error                          { return nil }
func (*runtimeContextProtocol) IsConnected() bool                     { return true }
func (*runtimeContextProtocol) Read(context.Context, connector.ReadRequest) (connector.ReadResult, error) {
	return connector.ReadResult{Value: 42, Quality: schema.QualityGood, Timestamp: time.Now()}, nil
}

type runtimeWriteContext struct {
	value any
	err   error
}

type runtimeContextWriter struct {
	storage.Writer
	contexts chan runtimeWriteContext
}

func (w *runtimeContextWriter) Write(ctx context.Context, _ storage.TimeSeriesRecord) error {
	select {
	case w.contexts <- runtimeWriteContext{value: ctx.Value(runtimeContextKey{}), err: ctx.Err()}:
	default:
	}
	return nil
}
func (*runtimeContextWriter) Flush(context.Context) error { return nil }
func (*runtimeContextWriter) Close() error                { return nil }

func TestRuntimeConsumerRetainsValuesAfterStartRequestEnds(t *testing.T) {
	protocolType := schema.ProtocolType("runtime-context-test")
	connector.Register(protocolType, func() connector.Protocol { return &runtimeContextProtocol{} })
	defer connector.Unregister(protocolType)
	groupID := "context-group"
	writer := &runtimeContextWriter{contexts: make(chan runtimeWriteContext, 1)}
	svc, err := NewService(Config{Writer: writer, Snapshot: Snapshot{
		Devices:       []*schema.Device{{ID: "context-device", Protocol: protocolType, Status: schema.DeviceStatusActive, ConnectionConfig: "{}"}},
		Points:        []*schema.Point{{ID: "context-point", DeviceID: "context-device", Address: "0", Enabled: true, PollingGroupID: &groupID, DataType: schema.DataTypeUint16}},
		PollingGroups: []*schema.PollingGroup{{ID: groupID, Enabled: true, IntervalMs: 10}},
		Mappings:      []*schema.Mapping{{ID: "context-mapping", PointID: "context-point", TagID: "context-tag", Enabled: true}},
		Tags:          []*schema.Tag{{ID: "context-tag", DataType: schema.DataTypeUint16}},
	}})
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.WithValue(t.Context(), runtimeContextKey{}, "start-request"))
	require.NoError(t, svc.Start(ctx))
	cancel()
	defer func() { require.NoError(t, svc.Stop(t.Context())) }()
	select {
	case observed := <-writer.contexts:
		require.Equal(t, "start-request", observed.value)
		require.NoError(t, observed.err)
	case <-time.After(2 * time.Second):
		t.Fatal("runtime stopped consuming after the start request was canceled")
	}
}
