package groupdelivery

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func (f *deliveryFixture) enqueueEntity(t *testing.T, entity string, start time.Time, counter uint64) string {
	t.Helper()
	destination := testDestination
	destination.TableSchema = ""
	destination.DedupeCapability = DedupeNone
	bucket := rowBucket(entity, start, counter)
	require.NoError(t, f.store.CommitClosure(t.Context(), testKey,
		Closure{Destination: destination, NextClose: start.Add(10 * time.Second), Buckets: []ClosedBucket{bucket}}))
	return bucket.Outcome.EffectKey
}

func (f *deliveryFixture) delivered(t *testing.T) []int64 {
	t.Helper()
	rows, err := f.target.QueryContext(t.Context(), `SELECT counter FROM readings ORDER BY rowid`)
	require.NoError(t, err)
	defer rows.Close()
	var counters []int64
	for rows.Next() {
		var counter int64
		require.NoError(t, rows.Scan(&counter))
		counters = append(counters, counter)
	}
	require.NoError(t, rows.Err())
	return counters
}

// poisonValue is the counter the destination rejects in these tests.
const poisonValue = 13

func (f *deliveryFixture) poisonDestination(t *testing.T) {
	t.Helper()
	_, err := f.target.ExecContext(t.Context(), `CREATE TRIGGER poison BEFORE INSERT ON readings WHEN NEW.counter = 13 BEGIN SELECT RAISE(ABORT, 'refused'); END`)
	require.NoError(t, err)
}

func newDispatcherFixture(t *testing.T, config DispatcherConfig) (*deliveryFixture, *Dispatcher) {
	t.Helper()
	f := newDeliveryFixture(t, plainTargetTable, false)
	return f, NewDispatcher(f.store, f.sender, config)
}

func TestBoundedRetryAndPoisonPartitionPoisonRowIsQuarantinedImmediately(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	f.poisonDestination(t)
	effect := f.enqueue(t, DedupeNone, "", t0, poisonValue)

	result, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateQuarantined, result.State)
	item := f.state(t, effect)
	require.Equal(t, "destination-rejected-row", item.LastErrorCode)
	require.NotEmpty(t, item.Payload, "the payload is kept for the operator")
	require.Zero(t, item.RetryCount, "a row that can never succeed is not retried")
	require.NotContains(t, item.LastErrorCode, "refused")
	_, err = f.sender.Deliver(t.Context(), effect)
	require.ErrorIs(t, err, ErrNotDeliverable)
	require.Zero(t, f.targetRows(t))
}

func TestBoundedRetryAndPoisonPartitionUnusableTargetBlocksInsteadOfQuarantining(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueue(t, DedupeNone, "", t0, 5)
	_, err := f.target.ExecContext(t.Context(), `DROP TABLE readings`)
	require.NoError(t, err)

	result, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateBlocked, result.State)
	require.Equal(t, "target-unusable", f.state(t, effect).LastErrorCode)
}

func TestBoundedRetryAndPoisonPartitionPoisonBlocksOnlyItsOwnPartition(t *testing.T) {
	f, dispatcher := newDispatcherFixture(t, DispatcherConfig{})
	f.poisonDestination(t)
	a1 := f.enqueueEntity(t, "plant-A", at(0), 1)
	a2 := f.enqueueEntity(t, "plant-A", at(10), poisonValue) // poison
	a3 := f.enqueueEntity(t, "plant-A", at(20), 3)
	b1 := f.enqueueEntity(t, "plant-B", at(0), 21)
	b2 := f.enqueueEntity(t, "plant-B", at(10), 22)

	report, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Equal(t, 2, report.Partitions)
	require.Equal(t, 3, report.Delivered, "plant-A's first row and both plant-B rows")
	require.Equal(t, 1, report.Quarantined)

	require.Equal(t, StateCommitted, f.state(t, a1).State)
	require.Equal(t, StateQuarantined, f.state(t, a2).State)
	require.Equal(t, StatePending, f.state(t, a3).State, "later rows of the same partition wait; they are not skipped")
	require.Equal(t, StateCommitted, f.state(t, b1).State)
	require.Equal(t, StateCommitted, f.state(t, b2).State, "a healthy partition is unaffected")
	require.ElementsMatch(t, []int64{1, 21, 22}, f.delivered(t))

	again, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Zero(t, again.Delivered, "nothing passes the quarantined row")
	require.Equal(t, StatePending, f.state(t, a3).State)
}

func TestBoundedRetryAndPoisonPartitionOperatorDispositionAdvancesInOrder(t *testing.T) {
	f, dispatcher := newDispatcherFixture(t, DispatcherConfig{})
	f.poisonDestination(t)
	f.enqueueEntity(t, "plant-A", at(0), 1)
	poison := f.enqueueEntity(t, "plant-A", at(10), 13)
	f.enqueueEntity(t, "plant-A", at(20), 3)
	_, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Equal(t, StateQuarantined, f.state(t, poison).State)

	require.ErrorIs(t, f.store.ResolveQuarantine(t.Context(), "missing", ResolutionSkip), ErrOutboxItemNotFound)
	require.NoError(t, f.store.ResolveQuarantine(t.Context(), poison, ResolutionSkip))
	skipped := f.state(t, poison)
	require.Equal(t, StateSkipped, skipped.State)
	require.NotEmpty(t, skipped.Payload, "a skipped row stays on record")

	report, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Equal(t, 1, report.Delivered)
	require.Equal(t, []int64{1, 3}, f.delivered(t), "the partition moved past the skipped row and kept its order")
	require.ErrorIs(t, f.store.ResolveQuarantine(t.Context(), poison, ResolutionRetry), ErrNotDeliverable, "a resolved row cannot be resolved again")
}

func TestBoundedRetryAndPoisonPartitionOperatorRetryAfterRepairKeepsOrder(t *testing.T) {
	f, dispatcher := newDispatcherFixture(t, DispatcherConfig{})
	f.poisonDestination(t)
	f.enqueueEntity(t, "plant-A", at(0), poisonValue)
	f.enqueueEntity(t, "plant-A", at(10), 3)
	first, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Equal(t, 1, first.Quarantined)

	_, err = f.target.ExecContext(t.Context(), `DROP TRIGGER poison`)
	require.NoError(t, err)
	poison := rowBucket("plant-A", at(0), poisonValue).Outcome.EffectKey
	require.NoError(t, f.store.ResolveQuarantine(t.Context(), poison, ResolutionRetry))
	require.Equal(t, StatePending, f.state(t, poison).State)
	require.Zero(t, f.state(t, poison).RetryCount, "the operator's retry starts a fresh budget")

	report, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Equal(t, 2, report.Delivered)
	require.Equal(t, []int64{13, 3}, f.delivered(t))
}

func TestBoundedRetryAndPoisonPartitionTransientHeadBlocksSuccessorsOnlyUntilDue(t *testing.T) {
	var clock atomic.Int64
	now := func() time.Time { return time.Unix(0, clock.Load()).UTC() }
	f := newDeliveryFixture(t, plainTargetTable, false)
	f.sender = NewSender(f.store, f.resolver, SenderConfig{MaxRetries: 5, Now: now, Backoff: func(retry int) (string, time.Time) {
		return StateRetrying, now().Add(time.Minute)
	}})
	dispatcher := NewDispatcher(f.store, f.sender, DispatcherConfig{Now: now})
	a1 := f.enqueueEntity(t, "plant-A", at(0), 1)
	f.enqueueEntity(t, "plant-A", at(10), 2)
	f.enqueueEntity(t, "plant-B", at(0), 21)
	clock.Store(time.Now().UTC().Add(time.Second).UnixNano()) // after the rows became due

	transient := errors.New("connection reset by peer")
	f.faults.ExecError.Store(&transient)
	report, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Equal(t, 2, report.Retrying, "both partitions' heads failed transiently")
	require.Zero(t, report.Delivered)

	f.faults.ExecError.Store(nil)
	again, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Zero(t, again.Delivered, "backoff has not elapsed; nothing is retried early and nothing overtakes")
	require.Equal(t, StateRetrying, f.state(t, a1).State)

	clock.Add(int64(2 * time.Minute))
	final, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Equal(t, 3, final.Delivered)
	require.ElementsMatch(t, []int64{1, 2, 21}, f.delivered(t))
	counters := f.delivered(t)
	var indexOf1, indexOf2 int
	for i, c := range counters {
		if c == 1 {
			indexOf1 = i
		}
		if c == 2 {
			indexOf2 = i
		}
	}
	require.Less(t, indexOf1, indexOf2, "order within a partition is preserved across retries")
}

func TestBoundedRetryAndPoisonPartitionRetryExhaustionBlocksPartitionAndKeepsData(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	f.sender = NewSender(f.store, f.resolver, SenderConfig{MaxRetries: 2, Backoff: func(retry int) (string, time.Time) {
		if retry >= 2 {
			return StateBlocked, time.Now().UTC()
		}
		return StateRetrying, time.Now().UTC().Add(-time.Second)
	}})
	dispatcher := NewDispatcher(f.store, f.sender, DispatcherConfig{})
	head := f.enqueueEntity(t, "plant-A", at(0), 1)
	tail := f.enqueueEntity(t, "plant-A", at(10), 2)
	transient := errors.New("connection reset by peer")
	f.faults.ExecError.Store(&transient)

	for i := 0; i < 3; i++ {
		_, err := dispatcher.RunOnce(t.Context())
		require.NoError(t, err)
	}
	require.Equal(t, StateBlocked, f.state(t, head).State)
	require.Equal(t, StatePending, f.state(t, tail).State)
	require.Equal(t, 2, countRows(t, f.local, "wg_delivery_outbox"), "no accepted row is ever deleted")

	f.faults.ExecError.Store(nil)
	report, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Zero(t, report.Delivered, "a blocked head keeps blocking until an operator acts")
	require.NoError(t, f.store.ResolveQuarantine(t.Context(), head, ResolutionRetry))
	report, err = dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Equal(t, 2, report.Delivered)
	require.Equal(t, []int64{1, 2}, f.delivered(t))
}

func TestBoundedRetryAndPoisonPartitionBatchLimitBoundsEachCycle(t *testing.T) {
	f, dispatcher := newDispatcherFixture(t, DispatcherConfig{BatchSize: 2})
	for i := 0; i < 5; i++ {
		f.enqueueEntity(t, "plant-A", at(10*i), uint64(i+1))
	}
	for _, want := range []int{2, 2, 1, 0} {
		report, err := dispatcher.RunOnce(t.Context())
		require.NoError(t, err)
		require.Equal(t, want, report.Delivered)
	}
	require.Equal(t, []int64{1, 2, 3, 4, 5}, f.delivered(t))
}

type gaugeResolver struct {
	inner   *staticResolver
	current atomic.Int32
	peak    atomic.Int32
}

func (g *gaugeResolver) Resolve(ctx context.Context, item OutboxItem) (Target, error) {
	now := g.current.Add(1)
	defer g.current.Add(-1)
	for {
		peak := g.peak.Load()
		if now <= peak || g.peak.CompareAndSwap(peak, now) {
			break
		}
	}
	time.Sleep(40 * time.Millisecond)
	return g.inner.Resolve(ctx, item)
}

func TestBoundedRetryAndPoisonPartitionConcurrencyIsBounded(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	gauge := &gaugeResolver{inner: f.resolver}
	f.sender = NewSender(f.store, gauge, SenderConfig{MaxRetries: 3})
	dispatcher := NewDispatcher(f.store, f.sender, DispatcherConfig{MaxPartitions: 2})
	for i := 0; i < 6; i++ {
		f.enqueueEntity(t, "plant-"+string(rune('A'+i)), at(0), uint64(i+1))
	}

	report, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Equal(t, 6, report.Partitions)
	require.Equal(t, 6, report.Delivered)
	require.EqualValues(t, 2, gauge.peak.Load(), "at most MaxPartitions partitions run at once, and they do run in parallel")
}

func TestBoundedRetryAndPoisonPartitionDefaultsAndConcurrentRunsNeverDoubleSend(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	dispatcher := NewDispatcher(f.store, f.sender, DispatcherConfig{BatchSize: -1, MaxPartitions: 0})
	require.Equal(t, DefaultBatchSize, dispatcher.config.BatchSize)
	require.Equal(t, DefaultMaxPartitions, dispatcher.config.MaxPartitions)

	for i := 0; i < 4; i++ {
		f.enqueueEntity(t, "plant-A", at(10*i), uint64(i+1))
	}
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = dispatcher.RunOnce(t.Context())
		}()
	}
	wg.Wait()
	_, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Equal(t, []int64{1, 2, 3, 4}, f.delivered(t), "overlapping cycles never send a row twice nor out of order")
	_ = schema.QualityGood
}

func TestProductionGroupOutageRecoveryDefaultsNeverTurnAnOutageIntoABlockedBacklog(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueueEntity(t, "plant-A", at(0), 1)
	sender := NewSender(f.store, f.resolver, SenderConfig{Owner: "node-1/a"}) // production defaults: no custom backoff, no retry cap
	transient := errors.New("connection refused")
	f.faults.ExecError.Store(&transient)

	// Far more attempts than the old five-retry budget, i.e. an outage lasting
	// minutes to hours. Transient failures must never exhaust into `blocked`.
	for i := 0; i < 15; i++ {
		result, err := sender.Deliver(t.Context(), effect)
		require.NoError(t, err)
		require.Equal(t, StateRetrying, result.State, "attempt %d", i+1)
	}
	item := f.state(t, effect)
	require.Equal(t, 15, item.RetryCount)
	require.Equal(t, "insert-failed", item.LastErrorCode)
	require.WithinDuration(t, time.Now().UTC(), item.NextRetryAt, 8*time.Minute, "backoff is capped, so recovery is noticed within minutes")

	// The destination returns: the very same row is delivered, no operator needed.
	f.faults.ExecError.Store(nil)
	result, err := sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, result.State)
	require.Equal(t, 1, f.targetRows(t))
}

func TestProductionGroupOutageRecoveryAnExplicitRetryCapStillBlocksAndKeepsTheData(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueueEntity(t, "plant-A", at(0), 1)
	sender := NewSender(f.store, f.resolver, SenderConfig{Owner: "node-1/a", MaxRetries: 3})
	transient := errors.New("connection refused")
	f.faults.ExecError.Store(&transient)
	for i := 0; i < 3; i++ {
		_, err := sender.Deliver(t.Context(), effect)
		if err != nil {
			require.ErrorIs(t, err, ErrNotDeliverable)
		}
	}
	require.Equal(t, StateBlocked, f.state(t, effect).State, "an operator-chosen cap still blocks, with the payload kept")
	require.NotEmpty(t, f.state(t, effect).Payload)
}
