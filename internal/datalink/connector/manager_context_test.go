package connector

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestCanceledCleanupBeforeLockRetainsExactSerialGeneration(t *testing.T) {
	typ := schema.ProtocolModbusRTU
	protocol := &generationProtocol{}
	registerGeneration(t, typ, func() Protocol { return protocol })
	manager := NewConnectionManager(DefaultConnectionManagerConfig())
	_, err := manager.GetOrCreate(t.Context(), "a", typ, `{"serial_port":"COM9"}`)
	require.NoError(t, err)
	entry := manager.entries["a"]
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	require.ErrorIs(t, manager.cleanupEntryContext(ctx, entry), context.Canceled)
	require.Same(t, entry.conn, manager.cleanup[entry])
	require.False(t, entry.cleaned)
	require.Zero(t, protocol.closed.Load())
	_, err = manager.GetOrCreate(t.Context(), "b", typ, `{"serial_port":"COM9"}`)
	require.Error(t, err)
	require.NoError(t, manager.CloseAll())
	require.EqualValues(t, 1, protocol.closed.Load())
	require.Empty(t, manager.cleanup)
}
