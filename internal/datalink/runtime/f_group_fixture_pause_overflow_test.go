//go:build f_write_group_fixture

package runtime

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

type fixtureBlockingSink struct {
	started  chan struct{}
	release  chan struct{}
	once     sync.Once
	accepted atomic.Int32
}

func (s *fixtureBlockingSink) AcceptSample(ctx context.Context, _ measurement.SampleEnvelope) error {
	s.once.Do(func() { close(s.started) })
	select {
	case <-s.release:
		s.accepted.Add(1)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type fixtureOverflowSink struct{}

func (fixtureOverflowSink) AcceptSample(ctx context.Context, _ measurement.SampleEnvelope) error {
	ReportFixtureSampleError(ctx)
	return errors.New("fixture capture queue full")
}

func TestPauseFixtureWaitsForInFlightConsumer(t *testing.T) {
	sink := &fixtureBlockingSink{started: make(chan struct{}), release: make(chan struct{})}
	svc, binding, _ := typedRuntimeFixture(t, sink)
	svc.scheduler = collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	svc.stopCh = make(chan struct{})
	var releaseOnce sync.Once
	var stopOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(sink.release) }) }
	stop := func() { stopOnce.Do(func() { close(svc.stopCh) }) }
	t.Cleanup(func() {
		release()
		stop()
		svc.wg.Wait()
	})
	// consumeLoop owns its matching Done call, so it must be paired with Add
	// rather than WaitGroup.Go (which also calls Done).
	svc.wg.Add(1)
	go svc.consumeLoop(context.Background())

	value := typedRuntimeValue(binding)
	if !svc.scheduler.EmitFixtureValue(value) {
		t.Fatal("timed out enqueueing fixture value")
	}
	select {
	case <-sink.started:
	case <-time.After(time.Second):
		t.Fatal("runtime consumer did not reach the blocking sink")
	}

	paused := make(chan error, 1)
	go func() { paused <- svc.PauseFixture(t.Context()) }()
	select {
	case err := <-paused:
		t.Fatalf("pause returned while consumer was in flight: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	release()
	require.NoError(t, <-paused)
	require.Equal(t, int32(1), sink.accepted.Load())
	require.Equal(t, uint64(1), svc.Snapshot().CollectedTotal)
	require.True(t, svc.scheduler.EmitFixtureValue(value))
	require.False(t, fixtureConsumerBegin(svc))
	require.Equal(t, int32(1), sink.accepted.Load())
	require.Equal(t, uint64(1), svc.Snapshot().CollectedTotal)
	stop()
	svc.wg.Wait()
}

func TestAcceptFixtureCollectedValueReportsCaptureOverflow(t *testing.T) {
	svc, binding, _ := typedRuntimeFixture(t, fixtureOverflowSink{})
	err := svc.AcceptFixtureCollectedValue(t.Context(), typedRuntimeValue(binding))
	require.Error(t, err)
	require.NotContains(t, err.Error(), "queue full")
}
