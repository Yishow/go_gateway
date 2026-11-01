package connector

import (
	"context"
	"sync"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type closeAllContextKey struct{}

type closeAllContextProtocol struct {
	generationProtocol
	entered chan context.Context
	release chan struct{}
}

func (p *closeAllContextProtocol) CloseContext(ctx context.Context) error {
	p.entered <- ctx
	<-p.release
	return p.Close()
}

func TestCloseAllContextCanceledCleanupRetainsExactGeneration(t *testing.T) {
	protocol := &generationProtocol{}
	registerGeneration(t, schema.ProtocolModbusRTU, func() Protocol { return protocol })
	manager := NewConnectionManager(DefaultConnectionManagerConfig())
	_, err := manager.GetOrCreate(t.Context(), "a", schema.ProtocolModbusRTU, `{"serial_port":"COM9"}`)
	require.NoError(t, err)
	entry := manager.entries["a"]
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	require.ErrorIs(t, manager.CloseAllContext(ctx), context.Canceled)
	require.Empty(t, manager.entries)
	require.Empty(t, manager.connections)
	require.Same(t, entry.conn, manager.cleanup[entry])
	require.False(t, entry.cleaned)
	require.Zero(t, protocol.closed.Load())
	_, err = manager.GetOrCreate(t.Context(), "b", schema.ProtocolModbusRTU, `{"serial_port":"COM9"}`)
	require.Error(t, err, "cancellation cannot release an unresolved serial handle")
	require.NoError(t, manager.CloseAll())
	require.EqualValues(t, 1, protocol.closed.Load())
	require.Empty(t, manager.cleanup)
}

func TestCloseAllContextWithoutCancelWaitsForOriginalCleanup(t *testing.T) {
	protocol := &closeAllContextProtocol{entered: make(chan context.Context, 1), release: make(chan struct{})}
	release := sync.OnceFunc(func() { close(protocol.release) })
	defer release()
	typ := schema.ProtocolType(t.Name())
	registerGeneration(t, typ, func() Protocol { return protocol })
	manager := NewConnectionManager(DefaultConnectionManagerConfig())
	_, err := manager.GetOrCreate(t.Context(), "a", typ, `{}`)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.WithValue(t.Context(), closeAllContextKey{}, "shutdown-owner"))
	cancel()
	completionCtx := context.WithoutCancel(ctx)
	done := make(chan error, 1)
	go func() { done <- manager.CloseAllContext(completionCtx) }()
	observed := <-protocol.entered
	require.Equal(t, completionCtx, observed, "cleanup must receive the completion context unchanged")
	require.NoError(t, observed.Err())
	require.Nil(t, observed.Done())
	require.Equal(t, "shutdown-owner", observed.Value(closeAllContextKey{}))
	select {
	case err := <-done:
		t.Fatalf("cleanup returned before its original handle closed: %v", err)
	default:
	}
	release()
	require.NoError(t, <-done)
	require.EqualValues(t, 1, protocol.closed.Load())
	require.Empty(t, manager.cleanup)
}

func TestCloseAllContextFailureRemainsObservableAndRetryable(t *testing.T) {
	protocol := &recoverableCleanupProtocol{}
	protocol.failClose.Store(true)
	typ := schema.ProtocolType(t.Name())
	registerGeneration(t, typ, func() Protocol { return protocol })
	manager := NewConnectionManager(DefaultConnectionManagerConfig())
	_, err := manager.GetOrCreate(t.Context(), "a", typ, `{}`)
	require.NoError(t, err)
	entry := manager.entries["a"]
	require.ErrorIs(t, manager.CloseAllContext(t.Context()), errConnectionCleanup)
	require.Same(t, entry.conn, manager.cleanup[entry])
	require.False(t, entry.cleaned)
	protocol.failClose.Store(false)
	require.NoError(t, manager.CloseAllContext(t.Context()))
	require.EqualValues(t, 2, protocol.closed.Load())
	require.Empty(t, manager.cleanup)
}
