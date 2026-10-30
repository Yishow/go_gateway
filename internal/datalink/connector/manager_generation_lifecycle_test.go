package connector

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

func TestGenerationCleanupDoesNotHoldGlobalLockWhileDeviceBusy(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	registerGeneration(t, typ, func() Protocol { return &generationProtocol{} })
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	a, e := cm.GetOrCreate(t.Context(), "a", typ, `{}`)
	require.NoError(t, e)
	a.Lock()
	done := make(chan struct{})
	defer func() { a.Unlock(); <-done }()
	go func() { cm.CleanupIdle(); close(done) }()
	// Cleanup may wait for A's request lock, but must not obstruct a different device.
	time.Sleep(20 * time.Millisecond)
	require.NoError(t, result(t, acquire(cm, t.Context(), "b", typ, `{}`)).err)
}

func TestGenerationReconnectDoesNotBlockOtherDevice(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	started := make(chan struct{})
	unblock := make(chan struct{})
	done := make(chan error, 1)
	defer func() { close(unblock); <-done }()
	calls := 0
	registerGeneration(t, typ, func() Protocol {
		calls++
		n := calls
		return &generationProtocol{dial: func(context.Context, string) error {
			if n == 2 {
				close(started)
				<-unblock
			}
			return nil
		}}
	})
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	_, e := cm.GetOrCreate(t.Context(), "a", typ, `{}`)
	require.NoError(t, e)
	go func() { done <- cm.Reconnect(t.Context(), "a") }()
	<-started
	require.NoError(t, result(t, acquire(cm, t.Context(), "b", typ, `{}`)).err)
}

func TestGenerationSerialOwnConfigChangeClosesBeforeReplacement(t *testing.T) {
	typ := schema.ProtocolModbusRTU
	var first *generationProtocol
	registerGeneration(t, typ, func() Protocol {
		p := &generationProtocol{dial: func(context.Context, string) error {
			if first != nil {
				require.EqualValues(t, 1, first.closed.Load())
			}
			return nil
		}}
		return p
	})
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	a, e := cm.GetOrCreate(t.Context(), "a", typ, `{"serial_port":"COM9","baud_rate":9600}`)
	require.NoError(t, e)
	first = a.Protocol.(*generationProtocol)
	b, e := cm.GetOrCreate(t.Context(), "a", typ, `{"serial_port":"COM9","baud_rate":19200}`)
	require.NoError(t, e)
	require.NotSame(t, a, b)
	require.NoError(t, cm.Close("a"))
	require.EqualValues(t, 1, first.closed.Load())
	require.EqualValues(t, 1, b.Protocol.(*generationProtocol).closed.Load())
}

func TestGenerationReconnectCannotResurrectAfterNewerClose(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	failed := make(chan struct{})
	var calls atomic.Int32
	registerGeneration(t, typ, func() Protocol {
		n := calls.Add(1)
		return &generationProtocol{dial: func(context.Context, string) error {
			if n == 2 {
				close(failed)
				return errors.New("failed")
			}
			return nil
		}}
	})
	cfg := DefaultConnectionManagerConfig()
	cfg.RetryDelay = 100 * time.Millisecond
	cm := NewConnectionManager(cfg)
	_, e := cm.GetOrCreate(t.Context(), "a", typ, `{"host":"old"}`)
	require.NoError(t, e)
	done := make(chan error, 1)
	go func() { done <- cm.Reconnect(t.Context(), "a") }()
	<-failed
	// Wait until the failed attempt has been removed before installing a newer generation.
	require.Eventually(t, func() bool { cm.mu.RLock(); defer cm.mu.RUnlock(); return cm.entries["a"] == nil }, time.Second, time.Millisecond)
	_, e = cm.GetOrCreate(t.Context(), "a", typ, `{"host":"new"}`)
	require.NoError(t, e)
	require.NoError(t, cm.Close("a"))
	require.Error(t, <-done)
	_, ok := cm.GetStatus("a")
	require.False(t, ok)
}

func TestGenerationCanceledWaiterDoesNotCancelOwner(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	started := make(chan struct{})
	unblock := make(chan struct{})
	var calls atomic.Int32
	registerGeneration(t, typ, func() Protocol {
		calls.Add(1)
		return &generationProtocol{dial: func(ctx context.Context, _ string) error { close(started); <-unblock; return ctx.Err() }}
	})
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	owner := acquire(cm, t.Context(), "a", typ, `{}`)
	<-started
	ctx, cancel := context.WithCancel(t.Context())
	waiter := acquire(cm, ctx, "a", typ, `{}`)
	cancel()
	require.ErrorIs(t, result(t, waiter).err, context.Canceled)
	close(unblock)
	require.NoError(t, result(t, owner).err)
	require.EqualValues(t, 1, calls.Load())
}

func TestGenerationFactoryAndDeadlineFailures(t *testing.T) {
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	_, e := cm.GetOrCreate(t.Context(), "a", schema.ProtocolType("not-registered-generation"), `{}`)
	require.Error(t, e)
	typ := schema.ProtocolType(t.Name())
	p := &generationProtocol{dial: func(ctx context.Context, _ string) error { <-ctx.Done(); return ctx.Err() }}
	registerGeneration(t, typ, func() Protocol { return p })
	cfg := DefaultConnectionManagerConfig()
	cfg.ConnectionTimeout = time.Millisecond
	cm = NewConnectionManager(cfg)
	_, e = cm.GetOrCreate(t.Context(), "a", typ, `{}`)
	require.ErrorIs(t, e, context.DeadlineExceeded)
	require.EqualValues(t, 1, p.closed.Load())
}
