package connector

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

type generationProtocol struct {
	stubProtocol
	dial   func(context.Context, string) error
	closed atomic.Int32
	live   atomic.Bool
}

func (p *generationProtocol) Connect(ctx context.Context, c string) error {
	if p.dial != nil {
		if err := p.dial(ctx, c); err != nil {
			return err
		}
	}
	p.live.Store(true)
	return nil
}
func (p *generationProtocol) Close() error      { p.closed.Add(1); p.live.Store(false); return nil }
func (p *generationProtocol) IsConnected() bool { return p.live.Load() }

type generationResult struct {
	conn *ManagedConnection
	err  error
}

func acquire(cm *ConnectionManager, ctx context.Context, id string, typ schema.ProtocolType, c string) <-chan generationResult {
	ch := make(chan generationResult, 1)
	go func() { v, e := cm.GetOrCreate(ctx, id, typ, c); ch <- generationResult{v, e} }()
	return ch
}
func result(t *testing.T, ch <-chan generationResult) generationResult {
	t.Helper()
	select {
	case r := <-ch:
		return r
	case <-time.After(time.Second):
		t.Fatal("connection operation blocked behind another device/dial")
		return generationResult{}
	}
}
func registerGeneration(t *testing.T, typ schema.ProtocolType, f ConnectorFactory) {
	t.Helper()
	Register(typ, f)
	t.Cleanup(func() { Unregister(typ) })
}

func TestGenerationSlowDialDoesNotBlockOtherDevice(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	started := make(chan struct{})
	unblock := make(chan struct{})
	registerGeneration(t, typ, func() Protocol {
		return &generationProtocol{dial: func(_ context.Context, c string) error {
			if strings.Contains(c, "slow") {
				close(started)
				<-unblock
			}
			return nil
		}}
	})
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	a := acquire(cm, t.Context(), "a", typ, `{"host":"slow"}`)
	defer func() { close(unblock); result(t, a) }()
	<-started
	require.NoError(t, result(t, acquire(cm, t.Context(), "b", typ, `{"host":"fast"}`)).err)
}
func TestGenerationConfigSwapRejectsLateDial(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	started := make(chan struct{})
	unblock := make(chan struct{})
	var first *generationProtocol
	registerGeneration(t, typ, func() Protocol {
		p := &generationProtocol{dial: func(_ context.Context, c string) error {
			if strings.Contains(c, "old") {
				close(started)
				<-unblock
			}
			return nil
		}}
		if first == nil {
			first = p
		}
		return p
	})
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	a := acquire(cm, t.Context(), "a", typ, `{"host":"old"}`)
	<-started
	bch := acquire(cm, t.Context(), "a", typ, `{"host":"new"}`)
	select {
	case b := <-bch:
		require.NoError(t, b.err)
		close(unblock)
		old := result(t, a)
		require.Error(t, old.err)
		require.Nil(t, old.conn)
		require.EqualValues(t, 1, first.closed.Load())
		current, e := cm.GetOrCreate(t.Context(), "a", typ, `{"host":"new"}`)
		require.NoError(t, e)
		require.Same(t, b.conn, current)
		require.True(t, current.Protocol.IsConnected())
	case <-time.After(time.Second):
		close(unblock)
		result(t, a)
		result(t, bch)
		t.Fatal("changed configuration blocked by old dial")
	}
}
func TestGenerationConfigAndProtocolIdentity(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	other := schema.ProtocolType(t.Name() + "2")
	var calls atomic.Int32
	f := func() Protocol { calls.Add(1); return &generationProtocol{} }
	registerGeneration(t, typ, f)
	registerGeneration(t, other, f)
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	a, e := cm.GetOrCreate(t.Context(), "a", typ, `{"host":"a","password":"secret"}`)
	require.NoError(t, e)
	b, e := cm.GetOrCreate(t.Context(), "a", typ, `{ "password":"secret", "host":"a" }`)
	require.NoError(t, e)
	require.Same(t, a, b)
	c, e := cm.GetOrCreate(t.Context(), "a", typ, `{"host":"a","password":"changed"}`)
	require.NoError(t, e)
	require.NotSame(t, a, c)
	d, e := cm.GetOrCreate(t.Context(), "a", other, c.Config)
	require.NoError(t, e)
	require.NotSame(t, c, d)
	require.EqualValues(t, 3, calls.Load())
}
func TestGenerationCloseAndCancelInvalidateDial(t *testing.T) {
	for _, action := range []string{"close", "all", "cancel"} {
		t.Run(action, func(t *testing.T) {
			typ := schema.ProtocolType(t.Name())
			started := make(chan struct{})
			unblock := make(chan struct{})
			p := &generationProtocol{dial: func(context.Context, string) error { close(started); <-unblock; return nil }}
			registerGeneration(t, typ, func() Protocol { return p })
			cm := NewConnectionManager(DefaultConnectionManagerConfig())
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			ch := acquire(cm, ctx, "a", typ, `{}`)
			<-started
			done := make(chan struct{})
			go func() {
				switch action {
				case "close":
					_ = cm.Close("a")
				case "all":
					_ = cm.CloseAll()
				case "cancel":
					cancel()
				}
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(time.Second):
				close(unblock)
				result(t, ch)
				t.Fatal("close blocked by dial")
			}
			close(unblock)
			r := result(t, ch)
			require.Error(t, r.err)
			require.Nil(t, r.conn)
			require.EqualValues(t, 1, p.closed.Load())
			_, ok := cm.GetStatus("a")
			require.False(t, ok)
		})
	}
}
func TestGenerationFailedDialCleansAndRedacts(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	p := &generationProtocol{dial: func(context.Context, string) error { return errors.New("private-host secret-password") }}
	registerGeneration(t, typ, func() Protocol { return p })
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	_, e := cm.GetOrCreate(t.Context(), "a", typ, `{}`)
	require.Error(t, e)
	require.NotContains(t, e.Error(), "private-host")
	require.NotContains(t, e.Error(), "secret-password")
	require.EqualValues(t, 1, p.closed.Load())
}
func TestGenerationSerialAliasRejected(t *testing.T) {
	for _, typ := range []schema.ProtocolType{schema.ProtocolModbusRTU, schema.ProtocolFatekFBs} {
		t.Run(string(typ), func(t *testing.T) {
			var calls atomic.Int32
			registerGeneration(t, typ, func() Protocol { calls.Add(1); return &generationProtocol{} })
			cm := NewConnectionManager(DefaultConnectionManagerConfig())
			_, e := cm.GetOrCreate(t.Context(), "a", typ, `{"mode":"serial","serial_port":"COM9","slave_id":1}`)
			require.NoError(t, e)
			_, e = cm.GetOrCreate(t.Context(), "b", typ, `{"mode":"serial","serial_port":"com9","slave_id":2}`)
			require.Error(t, e)
			if e != nil {
				require.NotContains(t, e.Error(), "COM9")
			}
			require.EqualValues(t, 1, calls.Load())
			_, e = cm.GetOrCreate(t.Context(), "c", typ, `{"mode":"serial","serial_port":"COM10"}`)
			require.NoError(t, e)
		})
	}
}

func TestGenerationSameKeySingleFlight(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	started := make(chan struct{})
	unblock := make(chan struct{})
	var calls atomic.Int32
	var once sync.Once
	registerGeneration(t, typ, func() Protocol {
		calls.Add(1)
		return &generationProtocol{dial: func(context.Context, string) error { once.Do(func() { close(started) }); <-unblock; return nil }}
	})
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	first := acquire(cm, t.Context(), "a", typ, `{}`)
	<-started
	var waiting []<-chan generationResult
	for i := 0; i < 20; i++ {
		waiting = append(waiting, acquire(cm, t.Context(), "a", typ, `{}`))
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	require.Error(t, result(t, acquire(cm, ctx, "a", typ, `{}`)).err)
	close(unblock)
	a := result(t, first)
	require.NoError(t, a.err)
	for _, ch := range waiting {
		b := result(t, ch)
		require.NoError(t, b.err)
		require.Same(t, a.conn, b.conn)
	}
	require.EqualValues(t, 1, calls.Load())
}

func TestGenerationCleanupRetiresIdentity(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	var calls atomic.Int32
	registerGeneration(t, typ, func() Protocol { calls.Add(1); return &generationProtocol{} })
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	a, e := cm.GetOrCreate(t.Context(), "a", typ, `{}`)
	require.NoError(t, e)
	a.mu.Lock()
	a.LastUsed = time.Now().Add(-time.Hour)
	a.mu.Unlock()
	require.Equal(t, 1, cm.CleanupIdle())
	b, e := cm.GetOrCreate(t.Context(), "a", typ, `{}`)
	require.NoError(t, e)
	require.NotSame(t, a, b)
	require.EqualValues(t, 1, a.Protocol.(*generationProtocol).closed.Load())
}

func TestGenerationSerialReservationDuringDialAndAfterClose(t *testing.T) {
	typ := schema.ProtocolModbusRTU
	started := make(chan struct{})
	unblock := make(chan struct{})
	var calls atomic.Int32
	registerGeneration(t, typ, func() Protocol {
		n := calls.Add(1)
		return &generationProtocol{dial: func(context.Context, string) error {
			if n == 1 {
				close(started)
				<-unblock
			}
			return nil
		}}
	})
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	ch := acquire(cm, t.Context(), "a", typ, `{"serial_port":"COM9"}`)
	<-started
	_, e := cm.GetOrCreate(t.Context(), "b", typ, `{"serial_port":"COM9"}`)
	require.Error(t, e)
	require.NoError(t, cm.Close("a"))
	_, e = cm.GetOrCreate(t.Context(), "b", typ, `{"serial_port":"COM9"}`)
	require.Error(t, e)
	close(unblock)
	require.Error(t, result(t, ch).err)
	_, e = cm.GetOrCreate(t.Context(), "b", typ, `{"serial_port":"COM9"}`)
	require.NoError(t, e)
}
