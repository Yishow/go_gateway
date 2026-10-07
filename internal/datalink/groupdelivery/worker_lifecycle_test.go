package groupdelivery

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type stubbornResolver struct {
	entered chan context.Context
	release chan struct{}
}

func (r *stubbornResolver) Resolve(ctx context.Context, _ OutboxItem) (Target, error) {
	r.entered <- ctx
	<-r.release // Models a driver that ignores cancellation.
	return Target{}, errors.New("offline test destination")
}

func TestWorkerShutdownDeadlineLeavesOriginalWorkerObservable(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueueEntity(t, "plant-A", at(0), 1)
	behind := f.enqueueEntity(t, "plant-A", at(10), 2)
	gate := &stubbornResolver{entered: make(chan context.Context, 1), release: make(chan struct{})}
	release := sync.OnceFunc(func() { close(gate.release) })
	defer release()
	sender := NewSender(f.store, gate, SenderConfig{Owner: "node/current", DeliveryTimeout: time.Minute})
	worker := NewWorker(f.store, NewDispatcher(f.store, sender, DispatcherConfig{}), WorkerConfig{NodeID: "node", Owner: "node/current"})
	require.NoError(t, worker.Start(t.Context()))
	var deliveryCtx context.Context
	select {
	case deliveryCtx = <-gate.entered:
	case <-time.After(3 * time.Second):
		t.Fatal("worker did not enter resolver")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	began := time.Now()
	require.ErrorIs(t, worker.StopContext(ctx), context.DeadlineExceeded)
	require.Less(t, time.Since(began), time.Second, "Stop does not add a second grace budget")
	require.ErrorIs(t, worker.StopContext(ctx), ErrShutdownForced)
	require.ErrorIs(t, worker.WaitStopped(ctx), context.DeadlineExceeded)
	require.ErrorIs(t, worker.Start(t.Context()), ErrWorkerRunning)
	select {
	case <-deliveryCtx.Done():
	case <-time.After(time.Second):
		t.Fatal("original shutdown did not cancel the attempt")
	}
	require.Equal(t, StateSending, f.state(t, effect).State)
	require.Equal(t, StatePending, f.state(t, behind).State)
	release()
	require.NoError(t, worker.WaitStopped(t.Context()))
	require.NoError(t, worker.Stop(time.Second))
	require.Equal(t, StateRetrying, f.state(t, effect).State)
	require.Equal(t, StatePending, f.state(t, behind).State)
	require.Zero(t, f.targetRows(t), "shutdown never claims backlog was sent")
}

func TestWorkerShutdownPreservesAcceptedAndUnknownAfterReopen(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	unknown := f.enqueueEntity(t, "ambiguous", at(0), 9)
	f.lose.Store(true)
	result, err := f.sender.Deliver(t.Context(), unknown)
	require.NoError(t, err)
	require.Equal(t, StateUnknown, result.State)
	f.lose.Store(false)
	pending := f.enqueueEntity(t, "ambiguous", at(10), 10)
	accepted := sampleAt("accepted-sample", "temperature", 23, 77)
	require.NoError(t, f.store.AppendSample(t.Context(), otherKey, at(20), accepted))
	// Use the fixture's actual durable path, then reopen only after stop completes.
	var seq int
	var name, path string
	require.NoError(t, f.local.QueryRowContext(t.Context(), "PRAGMA database_list").Scan(&seq, &name, &path))
	worker := NewWorker(f.store, NewDispatcher(f.store, f.sender, DispatcherConfig{}), WorkerConfig{NodeID: "node", Owner: "node/current"})
	require.NoError(t, worker.Start(t.Context()))
	require.NoError(t, worker.StopContext(t.Context()))
	require.NoError(t, worker.WaitStopped(t.Context()))
	beforeUnknown, beforePending := f.state(t, unknown), f.state(t, pending)
	require.NoError(t, f.local.Close())
	reopened := NewStore(openStoreDB(t, path))
	_, err = reopened.RecoverStaleClaims(t.Context(), "node", "node/next", time.Now())
	require.NoError(t, err)
	for _, before := range []OutboxItem{beforeUnknown, beforePending} {
		after, readErr := reopened.GetOutbox(t.Context(), before.EffectKey)
		require.NoError(t, readErr)
		require.Equal(t, before.State, after.State)
		require.Equal(t, before.Payload, after.Payload)
		require.Equal(t, before.PayloadDigest, after.PayloadDigest)
	}
	require.Equal(t, StateUnknown, beforeUnknown.State)
	require.Equal(t, StatePending, beforePending.State)
	restored, err := reopened.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
	require.Equal(t, accepted.SampleID, restored.Samples[0].SampleID)
	require.True(t, accepted.Value.Equal(restored.Samples[0].Value))
}

func TestWorkerShutdownConcurrentCallersObserveOneStop(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	worker := NewWorker(f.store, NewDispatcher(f.store, f.sender, DispatcherConfig{}), WorkerConfig{})
	require.NoError(t, worker.Start(t.Context()))
	results := make(chan error, 16)
	var wg sync.WaitGroup
	for range cap(results) {
		wg.Go(func() { results <- worker.StopContext(t.Context()) })
	}
	wg.Wait()
	close(results)
	for err := range results {
		require.NoError(t, err)
	}
	require.NoError(t, worker.WaitStopped(t.Context()))
	require.NoError(t, worker.Start(t.Context()))
	require.NoError(t, worker.Stop(time.Second))
}

func TestWorkerStartupCancellationInterruptsRecoveryWithoutStartingDelivery(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueueEntity(t, "pending", at(0), 1)
	f.local.SetMaxOpenConns(1)
	conn, err := f.local.Conn(t.Context())
	require.NoError(t, err)
	defer conn.Close()
	baseline := f.local.Stats().WaitCount
	worker := NewWorker(f.store, NewDispatcher(f.store, f.sender, DispatcherConfig{}), WorkerConfig{NodeID: "node", Owner: "node/new"})
	startup, cancel := context.WithCancel(t.Context())
	defer cancel()
	lifetime, stopLifetime := context.WithCancel(context.WithoutCancel(t.Context()))
	defer stopLifetime()
	started := make(chan error, 1)
	go func() { started <- worker.StartWithLifetime(startup, lifetime) }()
	require.Eventually(t, func() bool { return f.local.Stats().WaitCount > baseline }, time.Second, time.Millisecond, "recovery is blocked waiting for the only SQLite connection")
	cancel()
	select {
	case err := <-started:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("startup cancellation did not interrupt recovery")
	}
	require.NoError(t, lifetime.Err(), "startup and loop lifetimes are separate")
	require.NoError(t, worker.WaitStopped(t.Context()))
	require.NoError(t, conn.Close())
	require.Equal(t, StatePending, f.state(t, effect).State)
	require.Zero(t, f.targetRows(t))
	require.Zero(t, f.resolver.calls.Load(), "canceled recovery must never start delivery")
}

func TestWorkerStartupRequestCanEndWhileOwnedLifetimeContinues(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	worker := NewWorker(f.store, NewDispatcher(f.store, f.sender, DispatcherConfig{}), WorkerConfig{NodeID: "node", Owner: "node/current", Interval: time.Millisecond})
	startup, cancel := context.WithCancel(t.Context())
	require.NoError(t, worker.StartWithLifetime(startup, t.Context()))
	cancel()
	defer worker.Stop(time.Second)
	effect := f.enqueueEntity(t, "live", at(0), 1)
	require.Eventually(t, func() bool { return f.state(t, effect).State == StateCommitted }, 3*time.Second, 10*time.Millisecond)
	require.NoError(t, worker.StopContext(t.Context()))
	require.NoError(t, worker.WaitStopped(t.Context()))
	require.Equal(t, 1, f.targetRows(t))
}
