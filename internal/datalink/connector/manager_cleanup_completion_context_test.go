package connector

import (
	"context"
	"go-gateway/internal/datalink/schema"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIdleCleanupInFlightOutlivesRoutineCancellation(t *testing.T) {
	protocol := &closeAllContextProtocol{entered: make(chan context.Context, 1), release: make(chan struct{})}
	release := sync.OnceFunc(func() { close(protocol.release) })
	defer release()
	typ := schema.ProtocolType(t.Name())
	registerGeneration(t, typ, func() Protocol { return protocol })
	config := DefaultConnectionManagerConfig()
	config.IdleTimeout = 10 * time.Millisecond
	manager := NewConnectionManager(config)
	conn, err := manager.GetOrCreate(t.Context(), "a", typ, `{}`)
	require.NoError(t, err)
	conn.LastUsed = time.Now().Add(-time.Hour)
	ctx, cancel := context.WithCancel(context.WithValue(t.Context(), closeAllContextKey{}, "idle-owner"))
	defer cancel()
	manager.StartCleanupRoutine(ctx)
	var cleanupCtx context.Context
	select {
	case cleanupCtx = <-protocol.entered:
	case <-time.After(time.Second):
		t.Fatal("idle cleanup did not enter")
	}
	cancel()
	require.Nil(t, cleanupCtx.Done())
	require.NoError(t, cleanupCtx.Err())
	require.Equal(t, "idle-owner", cleanupCtx.Value(closeAllContextKey{}))
	require.Zero(t, protocol.closed.Load())
	release()
	require.Eventually(t, func() bool { return protocol.closed.Load() == 1 }, time.Second, time.Millisecond)
	manager.mu.RLock()
	require.Empty(t, manager.entries)
	manager.mu.RUnlock()
}

func TestIdleCleanupLegacyUnownedConnectionPreservesCloseFallback(t *testing.T) {
	protocol := &closeAllContextProtocol{entered: make(chan context.Context, 1), release: make(chan struct{})}
	close(protocol.release)
	manager := NewConnectionManager(DefaultConnectionManagerConfig())
	manager.connections["legacy"] = &ManagedConnection{Protocol: protocol, LastUsed: time.Now().Add(-24 * time.Hour)}
	require.Equal(t, 1, manager.CleanupIdleContext(t.Context()))
	require.EqualValues(t, 1, protocol.closed.Load())
	require.Empty(t, protocol.entered, "legacy unowned fallback must not switch to CloseContext")
	require.Empty(t, manager.connections)
}
