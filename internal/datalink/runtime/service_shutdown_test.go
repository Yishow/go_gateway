package runtime

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type shutdownWriter struct {
	storage.Writer
	flushes atomic.Int32
	closes  atomic.Int32
}

func (w *shutdownWriter) Flush(context.Context) error { w.flushes.Add(1); return nil }
func (w *shutdownWriter) Close() error                { w.closes.Add(1); return nil }

func TestServiceShutdownClosesNeverStartedResourcesOnce(t *testing.T) {
	writer := &shutdownWriter{}
	svc, err := NewService(Config{Writer: writer})
	require.NoError(t, err)
	require.NoError(t, svc.Stop(t.Context()))
	require.NoError(t, svc.Stop(t.Context()))
	require.EqualValues(t, 1, writer.closes.Load())
}

func TestServiceShutdownDeadlineDoesNotCloseUnderLiveConsumer(t *testing.T) {
	writer := &shutdownWriter{}
	svc, err := NewService(Config{Writer: writer})
	require.NoError(t, err)
	require.NoError(t, svc.Start(t.Context()))
	release := make(chan struct{})
	unblock := sync.OnceFunc(func() { close(release) })
	defer unblock()
	svc.wg.Go(func() { <-release })
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	stopped := make(chan error, 1)
	go func() { stopped <- svc.Stop(ctx) }()
	select {
	case err := <-stopped:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(250 * time.Millisecond):
		t.Fatal("Stop ignored the caller deadline while a consumer remained live")
	}
	require.Zero(t, writer.closes.Load())
	require.Error(t, svc.Start(t.Context()), "must not restart with unfinished workers")
	unblock()
	waiter, ok := any(svc).(interface{ WaitStopped(context.Context) error })
	require.True(t, ok, "runtime must expose original stop completion")
	require.NoError(t, waiter.WaitStopped(t.Context()))
	require.EqualValues(t, 1, writer.closes.Load())
}

func TestServiceShutdownFailedStartClosesBatchTimer(t *testing.T) {
	underlying := &shutdownWriter{}
	batch := storage.NewBatchWriter(underlying, storage.DefaultBatchWriterConfig())
	svc, err := NewService(Config{Writer: batch})
	require.NoError(t, err)
	require.NoError(t, svc.scheduler.Start(nil))
	require.Error(t, svc.Start(t.Context()), "existing scheduler makes bootstrap fail")
	require.NoError(t, svc.Stop(t.Context()))
	require.NoError(t, svc.WaitStopped(t.Context()))
	require.NoError(t, batch.WaitClosed(t.Context()))
	require.EqualValues(t, 1, underlying.closes.Load())
	require.NoError(t, svc.Stop(t.Context()))
	require.EqualValues(t, 1, underlying.closes.Load())
}

func TestServiceShutdownConcurrentCallersCloseOnce(t *testing.T) {
	writer := &shutdownWriter{}
	svc, err := NewService(Config{Writer: writer})
	require.NoError(t, err)
	require.NoError(t, svc.Start(t.Context()))
	results := make(chan error, 20)
	var wg sync.WaitGroup
	for range cap(results) {
		wg.Go(func() { results <- svc.Stop(t.Context()) })
	}
	wg.Wait()
	close(results)
	for err := range results {
		require.NoError(t, err)
	}
	require.EqualValues(t, 1, writer.closes.Load())
	require.EqualValues(t, 1, writer.flushes.Load())
}

type stoppingContextWriter struct {
	shutdownWriter
	entered chan context.Context
}

func (w *stoppingContextWriter) Write(ctx context.Context, _ storage.TimeSeriesRecord) error {
	w.entered <- ctx
	<-ctx.Done()
	return ctx.Err()
}

func TestServiceShutdownCancelsConsumerOwnedContext(t *testing.T) {
	protocolType := schema.ProtocolType("runtime-owned-shutdown")
	connector.Register(protocolType, func() connector.Protocol { return &runtimeContextProtocol{} })
	defer connector.Unregister(protocolType)
	groupID := "shutdown-group"
	writer := &stoppingContextWriter{entered: make(chan context.Context, 1)}
	svc, err := NewService(Config{Writer: writer, Snapshot: Snapshot{
		Devices:       []*schema.Device{{ID: "d", Protocol: protocolType, Status: schema.DeviceStatusActive, ConnectionConfig: "{}"}},
		Points:        []*schema.Point{{ID: "p", DeviceID: "d", Address: "0", Enabled: true, PollingGroupID: &groupID, DataType: schema.DataTypeUint16}},
		PollingGroups: []*schema.PollingGroup{{ID: groupID, Enabled: true, IntervalMs: 1}},
		Mappings:      []*schema.Mapping{{ID: "m", PointID: "p", TagID: "t", Enabled: true}},
		Tags:          []*schema.Tag{{ID: "t", DataType: schema.DataTypeUint16}},
	}})
	require.NoError(t, err)
	startCtx, cancel := context.WithCancel(context.WithValue(t.Context(), runtimeContextKey{}, "preserved"))
	require.NoError(t, svc.Start(startCtx))
	cancel()
	var consumer context.Context
	select {
	case consumer = <-writer.entered:
	case <-time.After(time.Second):
		t.Fatal("consumer did not reach writer")
	}
	require.NoError(t, consumer.Err())
	require.Equal(t, "preserved", consumer.Value(runtimeContextKey{}))
	ctx, stopCancel := context.WithTimeout(t.Context(), time.Second)
	defer stopCancel()
	require.NoError(t, svc.Stop(ctx))
	require.ErrorIs(t, consumer.Err(), context.Canceled)
	require.NoError(t, svc.WaitStopped(t.Context()))
	require.EqualValues(t, 1, writer.closes.Load())
}

type blockedProjection struct {
	entered chan struct{}
	release chan struct{}
}

func (p *blockedProjection) RuntimeProjection(context.Context) (*workspace.RuntimeProjection, error) {
	close(p.entered)
	<-p.release
	return nil, errors.New("startup unavailable")
}

func TestServiceShutdownRacesBlockedStartup(t *testing.T) {
	writer := &shutdownWriter{}
	svc, err := NewService(Config{Writer: writer})
	require.NoError(t, err)
	projection := &blockedProjection{entered: make(chan struct{}), release: make(chan struct{})}
	release := sync.OnceFunc(func() { close(projection.release) })
	defer release()
	svc.deviceSvc = &device.Service{}
	svc.workspace = projection
	started := make(chan error, 1)
	go func() { started <- svc.Start(t.Context()) }()
	select {
	case <-projection.entered:
	case <-time.After(time.Second):
		t.Fatal("startup did not reach projection")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	require.ErrorIs(t, svc.Stop(ctx), context.Canceled)
	require.Zero(t, writer.closes.Load())
	require.Error(t, svc.Start(t.Context()))
	release()
	require.Error(t, <-started)
	require.NoError(t, svc.WaitStopped(t.Context()))
	require.EqualValues(t, 1, writer.closes.Load())
}

type failedCloseTarget struct{ err error }

func (*failedCloseTarget) WriteTagValue(context.Context, string, any, time.Time) error { return nil }
func (w *failedCloseTarget) Close(context.Context) error                               { return w.err }

func TestServiceShutdownTargetErrorDoesNotSkipWriterClose(t *testing.T) {
	writer := &shutdownWriter{}
	svc, err := NewService(Config{Writer: writer})
	require.NoError(t, err)
	failure := errors.New("target close failure")
	svc.target = &failedCloseTarget{err: failure}
	require.ErrorIs(t, svc.Stop(t.Context()), failure)
	require.ErrorIs(t, svc.WaitStopped(t.Context()), failure)
	require.EqualValues(t, 1, writer.closes.Load())
}

func TestServiceShutdownFailedStartupCanRetryBeforeStop(t *testing.T) {
	writer := &shutdownWriter{}
	svc, err := NewService(Config{Writer: writer})
	require.NoError(t, err)
	require.NoError(t, svc.scheduler.Start(nil))
	require.Error(t, svc.Start(t.Context()))
	require.NoError(t, svc.scheduler.StopContext(t.Context()))
	require.NoError(t, svc.Start(t.Context()))
	require.NoError(t, svc.Stop(t.Context()))
	require.EqualValues(t, 1, writer.closes.Load())
}

func TestServiceStartupCancellationDoesNotLaunchConsumers(t *testing.T) {
	writer := &shutdownWriter{}
	svc, err := NewService(Config{Writer: writer})
	require.NoError(t, err)
	startup, cancel := context.WithCancel(t.Context())
	cancel()
	require.ErrorIs(t, svc.Start(startup), context.Canceled)
	require.False(t, svc.IsRunning())
	require.NoError(t, svc.Stop(t.Context()))
	require.NoError(t, svc.WaitStopped(t.Context()))
	require.EqualValues(t, 1, writer.closes.Load())
}
