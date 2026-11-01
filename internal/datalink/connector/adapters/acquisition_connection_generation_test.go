package adapters

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/mcprotocol"

	"github.com/stretchr/testify/require"
)

type acquisitionGenerationTransport struct {
	connect func(context.Context) error
	live    bool
}

func (tr *acquisitionGenerationTransport) Connect() error {
	return tr.ConnectContext(context.Background())
}
func (tr *acquisitionGenerationTransport) ConnectContext(ctx context.Context) error {
	if tr.connect != nil {
		return tr.connect(ctx)
	}
	return ctx.Err()
}
func (*acquisitionGenerationTransport) Close() error { return nil }
func (*acquisitionGenerationTransport) SendReceive([]byte) ([]byte, error) {
	return nil, errors.New("unexpected send")
}
func (tr *acquisitionGenerationTransport) IsConnectedContext(context.Context) (bool, error) {
	return tr.live, nil
}

func TestAcquisitionTransportGenerationPublishedByReconnect(t *testing.T) {
	for _, protocol := range []string{"mc", "fatek"} {
		t.Run(protocol, func(t *testing.T) {
			transport := &acquisitionGenerationTransport{}
			var generation func() uint64
			var reconnect func(context.Context) error
			if protocol == "mc" {
				c := &MC3EConnector{client: mcprotocol.NewClientWithTransport(transport), connected: true, persistentMode: true}
				generation, reconnect = c.AcquisitionTransportGeneration, c.Reconnect
			} else {
				c := &FatekConnector{client: fatek.NewClient(transport, 1), connected: true, persistentMode: true}
				generation, reconnect = c.AcquisitionTransportGeneration, c.Reconnect
			}
			require.Zero(t, generation(), "manually assembled fixture is not proof of successful connect")
			require.NoError(t, reconnect(t.Context()))
			first := generation()
			require.NotZero(t, first)
			failure := errors.New("dial failed")
			transport.connect = func(context.Context) error { return failure }
			require.ErrorIs(t, reconnect(t.Context()), failure)
			require.Equal(t, first, generation(), "failed reconnect cannot invalidate a cached downgrade")
			ctx, cancel := context.WithCancel(t.Context())
			transport.connect = func(context.Context) error { cancel(); return nil }
			require.ErrorIs(t, reconnect(ctx), context.Canceled)
			require.Equal(t, first, generation(), "canceled late success is not a new usable connection")
			transport.connect = nil
			require.NoError(t, reconnect(t.Context()))
			require.Greater(t, generation(), first)
		})
	}
}

func TestAcquisitionTransportGenerationIncludesEnsureConnection(t *testing.T) {
	transport := &acquisitionGenerationTransport{}
	mc := &MC3EConnector{client: mcprotocol.NewClientWithTransport(transport), connected: true, persistentMode: true}
	require.NoError(t, mc.ensureConnectionContext(t.Context()))
	require.NotZero(t, mc.AcquisitionTransportGeneration())
	fc := &FatekConnector{client: fatek.NewClient(transport, 1), connected: true, persistentMode: true}
	require.NoError(t, fc.ensureConnectionContext(t.Context()))
	require.NotZero(t, fc.AcquisitionTransportGeneration())
	transport.live = true
	before := fc.AcquisitionTransportGeneration()
	require.NoError(t, fc.ensureConnectionContext(t.Context()))
	require.Equal(t, before, fc.AcquisitionTransportGeneration(), "a liveness check is not a reconnect")
}

func TestAcquisitionTransportGenerationRejectsLateSuccess(t *testing.T) {
	var generation acquisitionGeneration
	entered, release := make(chan struct{}), make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- generation.connect(t.Context(), func(context.Context) error {
			close(entered)
			<-release
			return nil
		})
	}()
	<-entered
	require.NoError(t, generation.connect(t.Context(), func(context.Context) error { return nil }))
	newer := generation.published.Load()
	require.NotZero(t, newer)
	close(release)
	require.ErrorIs(t, <-done, errAcquisitionConnectSuperseded)
	require.Equal(t, newer, generation.published.Load(), "late success must not clear a newer session's downgrade")
	failed := errors.New("failed reconnect")
	require.ErrorIs(t, generation.connect(t.Context(), func(context.Context) error { return failed }), failed)
	require.Equal(t, newer, generation.published.Load())
}

func TestAcquisitionTransportGenerationExhaustionDoesNotConnect(t *testing.T) {
	var generation acquisitionGeneration
	generation.attempt = ^uint64(0)
	called := false
	require.Error(t, generation.connect(t.Context(), func(context.Context) error { called = true; return nil }))
	require.False(t, called)
	require.Zero(t, generation.published.Load())
}

func TestAcquisitionTransportGenerationPublicationRejectsNewerPendingOrFailedAttempt(t *testing.T) {
	for _, newerFailed := range []bool{false, true} {
		t.Run(map[bool]string{false: "pending", true: "failed"}[newerFailed], func(t *testing.T) {
			var generation acquisitionGeneration
			require.NoError(t, generation.connect(t.Context(), func(context.Context) error { return nil }))
			previous := generation.published.Load()
			older, err := generation.begin(t.Context())
			require.NoError(t, err)
			// Freeze the older successful operation just before its publication.
			// A newer begin must invalidate that completion even if it never succeeds.
			entered, release := make(chan struct{}), make(chan struct{})
			unblock := sync.OnceFunc(func() { close(release) })
			t.Cleanup(unblock)
			failure := errors.New("newer connect failed")
			done := make(chan error, 1)
			go func() {
				done <- generation.connect(t.Context(), func(context.Context) error {
					close(entered)
					<-release
					return failure
				})
			}()
			<-entered
			if newerFailed {
				unblock()
				require.ErrorIs(t, <-done, failure)
			}
			published := make(chan error, 1)
			go func() { published <- generation.publish(t.Context(), older) }()
			select {
			case err := <-published:
				require.ErrorIs(t, err, errAcquisitionConnectSuperseded)
			case <-time.After(time.Second):
				t.Fatal("publication lock was held across the newer connection operation")
			}
			require.Equal(t, previous, generation.published.Load(), "late success must preserve the existing downgrade token")
			if !newerFailed {
				unblock()
				require.ErrorIs(t, <-done, failure)
			}
			require.NoError(t, generation.connect(t.Context(), func(context.Context) error { return nil }))
			require.Greater(t, generation.published.Load(), older)
		})
	}
}
