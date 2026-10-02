package groupdelivery

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// controlledClock lets a test move the store's notion of time.
type controlledClock struct{ nanos atomic.Int64 }

func newControlledClock() *controlledClock {
	c := &controlledClock{}
	c.nanos.Store(time.Now().UTC().UnixNano())
	return c
}

func (c *controlledClock) now() time.Time          { return time.Unix(0, c.nanos.Load()).UTC() }
func (c *controlledClock) advance(d time.Duration) { c.nanos.Add(int64(d)) }

func (f *deliveryFixture) enqueueRaw(t *testing.T, capability string, start time.Time, counter uint64) string {
	t.Helper()
	return f.enqueue(t, capability, "", start, counter)
}

func TestWorkerFencingAndShutdownOnlyOneOwnerHoldsAClaimAndEpochsIncrease(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, true)
	effect := f.enqueueRaw(t, DedupeReceipt, t0, 1)

	first, err := f.store.BeginDelivery(t.Context(), effect, "node-1/a", time.Minute)
	require.NoError(t, err)
	require.Equal(t, "node-1/a", first.Owner)
	require.Equal(t, int64(1), first.Epoch)
	_, err = f.store.BeginDelivery(t.Context(), effect, "node-1/b", time.Minute)
	require.ErrorIs(t, err, ErrNotDeliverable, "a held claim cannot be taken")

	require.NoError(t, f.store.MarkRetry(t.Context(), first, StateRetrying, "insert-failed", time.Now().UTC().Add(-time.Second)))
	second, err := f.store.BeginDelivery(t.Context(), effect, "node-1/b", time.Minute)
	require.NoError(t, err)
	require.Equal(t, int64(2), second.Epoch, "every claim has a strictly newer fencing epoch")
}

func TestWorkerFencingAndShutdownStaleWorkerCannotWriteAfterTakeover(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, true)
	clock := newControlledClock()
	f.store.now = clock.now
	effect := f.enqueueRaw(t, DedupeReceipt, t0, 5)

	stale, err := f.store.BeginDelivery(t.Context(), effect, "node-1/old", 10*time.Second)
	require.NoError(t, err)
	clock.advance(30 * time.Second) // the old worker stalled past its lease

	recovered, err := f.store.RecoverStaleClaims(t.Context(), "node-9", "node-9/self", clock.now())
	require.NoError(t, err)
	require.Equal(t, 1, recovered)
	require.Equal(t, StateRetrying, f.state(t, effect).State, "a destination receipt makes recovery a safe retry")

	immediate := func(int) (string, time.Time) { return StateRetrying, clock.now().Add(-time.Second) }
	fresh := NewSender(f.store, f.resolver, SenderConfig{Owner: "node-2/new", Backoff: immediate, Now: clock.now})
	result, err := fresh.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, result.State)

	// The stalled worker wakes up and tries to record its outcome.
	require.ErrorIs(t, f.store.CompleteDelivery(t.Context(), stale, "digest"), ErrFenced)
	require.ErrorIs(t, f.store.MarkRetry(t.Context(), stale, StateRetrying, "x", clock.now()), ErrFenced)
	require.ErrorIs(t, f.store.MarkUnknown(t.Context(), stale, "x"), ErrFenced)
	require.ErrorIs(t, f.store.MarkQuarantined(t.Context(), stale, "x"), ErrFenced)
	require.ErrorIs(t, f.store.BlockClaimed(t.Context(), stale, "x"), ErrFenced)
	require.Equal(t, StateCommitted, f.state(t, effect).State, "the stale worker changed nothing")
	require.Equal(t, 1, f.targetRows(t))
	require.Equal(t, 1, countRows(t, f.local, "wg_delivery_receipts"))
}

func TestWorkerFencingAndShutdownRecoveredRowWithoutDedupeIsUnknownNotResent(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	clock := newControlledClock()
	f.store.now = clock.now
	effect := f.enqueueRaw(t, DedupeNone, t0, 5)
	stale, err := f.store.BeginDelivery(t.Context(), effect, "node-1/old", 10*time.Second)
	require.NoError(t, err)
	clock.advance(time.Minute)

	recovered, err := f.store.RecoverStaleClaims(t.Context(), "node-9", "node-9/self", clock.now())
	require.NoError(t, err)
	require.Equal(t, 1, recovered)
	require.Equal(t, StateUnknown, f.state(t, effect).State, "the earlier attempt may have committed")

	fresh := NewSender(f.store, f.resolver, SenderConfig{Owner: "node-2/new"})
	_, err = fresh.Deliver(t.Context(), effect)
	require.ErrorIs(t, err, ErrNotDeliverable)
	require.ErrorIs(t, f.store.CompleteDelivery(t.Context(), stale, "digest"), ErrFenced, "and the stale worker cannot mark it committed")
	require.Equal(t, StateUnknown, f.state(t, effect).State)
	require.Zero(t, f.targetRows(t))
}

func TestWorkerFencingAndShutdownRestartRecoversPriorIncarnationBeforeLeaseExpiry(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, true)
	mine := f.enqueueRaw(t, DedupeReceipt, t0, 1)
	earlier := f.enqueueRaw(t, DedupeReceipt, at(10), 2)
	otherNode := f.enqueueRaw(t, DedupeReceipt, at(20), 3)

	_, err := f.store.BeginDelivery(t.Context(), earlier, "node-1/previous", time.Hour)
	require.NoError(t, err)
	_, err = f.store.BeginDelivery(t.Context(), mine, "node-1/current", time.Hour)
	require.NoError(t, err)
	_, err = f.store.BeginDelivery(t.Context(), otherNode, "node-2/running", time.Hour)
	require.NoError(t, err)

	recovered, err := f.store.RecoverStaleClaims(t.Context(), "node-1", "node-1/current", time.Now().UTC())
	require.NoError(t, err)
	require.Equal(t, 1, recovered, "only this node's earlier incarnation is provably gone")
	require.Equal(t, StateRetrying, f.state(t, earlier).State)
	require.Equal(t, StateSending, f.state(t, mine).State, "the current incarnation keeps its own claims")
	require.Equal(t, StateSending, f.state(t, otherNode).State, "another node's unexpired claim is never taken")

	again, err := f.store.RecoverStaleClaims(t.Context(), "node-1", "node-1/current", time.Now().UTC())
	require.NoError(t, err)
	require.Zero(t, again, "recovery is idempotent")
}

func TestWorkerFencingAndShutdownTwoWorkersNeverSendTheSameRowTwice(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	const items = 24
	for i := 0; i < items; i++ {
		f.enqueueEntity(t, "plant-"+string(rune('A'+i%4)), at(10*(i/4)), uint64(i+1))
	}
	var senders [2]*Dispatcher
	for i := range senders {
		sender := NewSender(f.store, f.resolver, SenderConfig{Owner: "node-1/w" + string(rune('1'+i)), MaxRetries: 5})
		senders[i] = NewDispatcher(f.store, sender, DispatcherConfig{BatchSize: 3, MaxPartitions: 2})
	}
	var wg sync.WaitGroup
	for _, dispatcher := range senders {
		wg.Add(1)
		go func(d *Dispatcher) {
			defer wg.Done()
			for round := 0; round < 30; round++ {
				_, _ = d.RunOnce(t.Context())
			}
		}(dispatcher)
	}
	wg.Wait()
	for i := 0; i < 3; i++ { // drain what the racing workers left behind
		_, err := senders[0].RunOnce(t.Context())
		require.NoError(t, err)
	}
	require.Equal(t, items, f.targetRows(t), "even without dedupe, exclusive claims keep every row to one send")
	require.Equal(t, items, openCount(t, f.local, "wg_delivery_outbox", "state = 'sql_committed'"))
}

type gatedResolver struct {
	inner   *staticResolver
	hangFor string
	entered chan struct{}
	once    sync.Once
}

func (g *gatedResolver) Resolve(ctx context.Context, item OutboxItem) (Target, error) {
	if item.EffectKey == g.hangFor {
		g.once.Do(func() { close(g.entered) })
		<-ctx.Done() // a hung destination: only cancellation ends it
		return Target{}, ctx.Err()
	}
	return g.inner.Resolve(ctx, item)
}

func TestWorkerFencingAndShutdownStopIsBoundedAndLeavesIntentRecoverable(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, true)
	hung := f.enqueueEntity(t, "plant-A", at(0), 1)
	behindHung := f.enqueueEntity(t, "plant-A", at(10), 2)
	healthy := f.enqueueEntity(t, "plant-B", at(0), 3)
	gate := &gatedResolver{inner: f.resolver, hangFor: hung, entered: make(chan struct{})}
	sender := NewSender(f.store, gate, SenderConfig{Owner: "node-1/first", MaxRetries: 5})
	worker := NewWorker(f.store, NewDispatcher(f.store, sender, DispatcherConfig{}),
		WorkerConfig{Interval: 10 * time.Millisecond, NodeID: "node-1", Owner: "node-1/first"})
	require.NoError(t, worker.Start(t.Context()))
	require.ErrorIs(t, worker.Start(t.Context()), ErrWorkerRunning)
	select {
	case <-gate.entered:
	case <-time.After(5 * time.Second):
		t.Fatal("the delivery never reached the hung destination")
	}

	began := time.Now()
	err := worker.Stop(200 * time.Millisecond)
	require.ErrorIs(t, err, ErrShutdownForced)
	require.Less(t, time.Since(began), 3*time.Second, "shutdown is bounded by its deadline, not by the hung destination")
	require.NoError(t, worker.Stop(time.Second), "stopping again is harmless")

	require.Contains(t, []string{StateSending, StateRetrying}, f.state(t, hung).State, "the interrupted row is not lost or marked done")
	require.Equal(t, StatePending, f.state(t, behindHung).State, "rows behind it never started")
	require.Equal(t, 1, countRows(t, f.local, "wg_delivery_receipts")+0*countRows(t, f.local, "wg_delivery_outbox"), "only the healthy partition's row was delivered")
	require.Equal(t, StateCommitted, f.state(t, healthy).State)

	// A new incarnation recovers and finishes everything exactly once.
	restarted := NewSender(f.store, f.resolver, SenderConfig{Owner: "node-1/second", MaxRetries: 5})
	next := NewWorker(f.store, NewDispatcher(f.store, restarted, DispatcherConfig{}),
		WorkerConfig{Interval: 10 * time.Millisecond, NodeID: "node-1", Owner: "node-1/second"})
	require.NoError(t, next.Start(t.Context()))
	require.Eventually(t, func() bool { return f.targetRows(t) == 3 }, 5*time.Second, 20*time.Millisecond)
	require.NoError(t, next.Stop(2*time.Second))
	require.Equal(t, 3, f.targetRows(t), "every accepted row reached the destination exactly once")
	require.Equal(t, StateCommitted, f.state(t, hung).State)
	require.Equal(t, StateCommitted, f.state(t, behindHung).State)
}

func TestWorkerFencingAndShutdownIdleStopIsImmediateAndGraceful(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	sender := NewSender(f.store, f.resolver, SenderConfig{Owner: "node-1/a"})
	worker := NewWorker(f.store, NewDispatcher(f.store, sender, DispatcherConfig{}), WorkerConfig{Interval: 10 * time.Millisecond, NodeID: "node-1", Owner: "node-1/a"})
	require.NoError(t, worker.Stop(time.Second), "stopping a worker that never started is a no-op")
	require.NoError(t, worker.Start(t.Context()))
	began := time.Now()
	require.NoError(t, worker.Stop(5*time.Second))
	require.Less(t, time.Since(began), time.Second)
}

func TestWorkerFencingAndShutdownDeliveryTimeoutBoundsAHungDestination(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueueRaw(t, DedupeNone, t0, 5)
	gate := &gatedResolver{inner: f.resolver, hangFor: effect, entered: make(chan struct{})}
	immediate := func(int) (string, time.Time) { return StateRetrying, time.Now().UTC().Add(-time.Second) }
	sender := NewSender(f.store, gate, SenderConfig{Owner: "node-1/a", DeliveryTimeout: 150 * time.Millisecond, Backoff: immediate})

	began := time.Now()
	result, err := sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Less(t, time.Since(began), 3*time.Second)
	require.Equal(t, StateRetrying, result.State, "a timeout before anything was sent is a safe retry")
	require.Equal(t, "target-unavailable", f.state(t, effect).LastErrorCode)
	require.Zero(t, f.targetRows(t))
}

func TestWorkerFencingAndShutdownCommittedEffectIsRecordedEvenWhenTheCallerIsCancelled(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, true)
	effect := f.enqueueRaw(t, DedupeReceipt, t0, 5)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	hook := func() { cancel() } // shutdown begins the instant the destination commits
	f.faults.AfterCommit.Store(&hook)
	sender := NewSender(f.store, f.resolver, SenderConfig{Owner: "node-1/a"})

	result, err := sender.Deliver(ctx, effect)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, result.State, "the bounded settle context records a committed effect despite cancellation")
	require.Equal(t, 1, f.targetRows(t))
	require.Equal(t, 1, countRows(t, f.local, "wg_delivery_receipts"))
	require.NotErrorIs(t, err, context.Canceled)
	_ = errors.New
}

func TestWorkerFencingAndShutdownRunningWorkerReapsAnExpiredClaimWithoutARestart(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, true)
	sender := NewSender(f.store, f.resolver, SenderConfig{Owner: "node-9/live"})
	worker := NewWorker(f.store, NewDispatcher(f.store, sender, DispatcherConfig{}),
		WorkerConfig{Interval: 20 * time.Millisecond, NodeID: "node-9", Owner: "node-9/live"})
	require.NoError(t, worker.Start(t.Context()))
	defer func() { _ = worker.Stop(time.Second) }()

	// A claim of some other worker expires while this worker keeps running: the
	// row would otherwise hold its whole partition until the next gateway restart.
	effect := f.enqueueRaw(t, DedupeReceipt, t0, 5)
	_, err := f.store.BeginDelivery(t.Context(), effect, "node-2/gone", time.Millisecond)
	require.NoError(t, err)

	require.Eventually(t, func() bool { return f.state(t, effect).State == StateCommitted }, 10*time.Second, 25*time.Millisecond)
	require.Equal(t, 1, f.targetRows(t), "the recovered row was delivered exactly once")
}

func TestWorkerFencingAndShutdownWorkerCanBeStartedAgainAfterItsContextEnds(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	sender := NewSender(f.store, f.resolver, SenderConfig{Owner: "node-1/a"})
	worker := NewWorker(f.store, NewDispatcher(f.store, sender, DispatcherConfig{}), WorkerConfig{Interval: 10 * time.Millisecond})
	ctx, cancel := context.WithCancel(t.Context())
	require.NoError(t, worker.Start(ctx))
	cancel() // the owner of the context goes away without calling Stop

	require.Eventually(t, func() bool {
		err := worker.Start(t.Context())
		if err == nil {
			_ = worker.Stop(time.Second)
		}
		return err == nil
	}, 5*time.Second, 20*time.Millisecond, "a worker whose loop ended must not stay marked as running")
}

func TestWorkerFencingAndShutdownWorkerReclaimsFinishedDataPeriodically(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	done := f.enqueueRaw(t, DedupeNone, t0, 1)
	claim, err := f.store.BeginDelivery(t.Context(), done, "w", time.Minute)
	require.NoError(t, err)
	require.NoError(t, f.store.CompleteDelivery(t.Context(), claim, "d"))
	_, err = f.local.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET committed_at = ? WHERE effect_key = ?`,
		timestamp(time.Now().UTC().Add(-48*time.Hour)), done)
	require.NoError(t, err)
	pending := f.enqueueRaw(t, DedupeNone, at(10), 2)
	require.NoError(t, f.store.MarkBlocked(t.Context(), pending, "target-blocked"))

	sender := NewSender(f.store, f.resolver, SenderConfig{Owner: "node-1/a"})
	worker := NewWorker(f.store, NewDispatcher(f.store, sender, DispatcherConfig{}), WorkerConfig{
		Interval: 10 * time.Millisecond, ReclaimEvery: 20 * time.Millisecond, ReclaimRetention: 24 * time.Hour,
	})
	require.NoError(t, worker.Start(t.Context()))
	defer func() { _ = worker.Stop(time.Second) }()

	require.Eventually(t, func() bool { return countRows(t, f.local, "wg_delivery_outbox") == 1 }, 5*time.Second, 20*time.Millisecond,
		"the old committed row is reclaimed by the running worker")
	require.Equal(t, StateBlocked, f.state(t, pending).State, "undelivered data is never reclaimed")
}
