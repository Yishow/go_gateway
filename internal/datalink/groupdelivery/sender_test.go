package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

const (
	plainTargetTable  = `CREATE TABLE readings (record_id TEXT, counter INTEGER NOT NULL, label TEXT, ratio REAL, running INTEGER, optional TEXT, bucket_start TEXT)`
	uniqueTargetTable = `CREATE TABLE readings (record_id TEXT PRIMARY KEY, counter INTEGER NOT NULL, label TEXT, ratio REAL, running INTEGER, optional TEXT, bucket_start TEXT)`
)

type staticResolver struct {
	target Target
	err    error
	calls  atomic.Int32
}

func (r *staticResolver) Resolve(context.Context, OutboxItem) (Target, error) {
	r.calls.Add(1)
	return r.target, r.err
}

type deliveryFixture struct {
	store    *Store
	local    *sql.DB
	target   *sql.DB
	lose     *atomic.Bool
	faults   *lostcommit.Faults
	resolver *staticResolver
	sender   *Sender
}

func newDeliveryFixture(t *testing.T, targetDDL string, receiptTable bool) *deliveryFixture {
	t.Helper()
	store, local, _ := newTestStore(t)
	faults := &lostcommit.Faults{Lose: &atomic.Bool{}}
	target := lostcommit.OpenFaults(filepath.Join(t.TempDir(), "target.db")+"?_pragma=busy_timeout(5000)", faults)
	t.Cleanup(func() { _ = target.Close() })
	_, err := target.ExecContext(t.Context(), targetDDL)
	require.NoError(t, err)
	if receiptTable {
		_, err = target.ExecContext(t.Context(), `CREATE TABLE gw_effect_receipts (effect_key TEXT PRIMARY KEY, payload_digest TEXT NOT NULL, committed_at TEXT NOT NULL)`)
		require.NoError(t, err)
	}
	resolver := &staticResolver{target: Target{DB: target, Kind: schema.DatabaseConnectorKindSQLite}}
	immediate := func(retry int) (string, time.Time) { return StateRetrying, time.Now().UTC().Add(-time.Second) }
	return &deliveryFixture{
		store: store, local: local, target: target, lose: faults.Lose, faults: faults, resolver: resolver,
		sender: NewSender(store, resolver, SenderConfig{MaxRetries: 3, Backoff: immediate}),
	}
}

func (f *deliveryFixture) enqueue(t *testing.T, capability, recordKeyColumn string, start time.Time, counter uint64) string {
	t.Helper()
	destination := testDestination
	destination.TableSchema = ""
	destination.DedupeCapability = capability
	destination.RecordKeyColumn = recordKeyColumn
	bucket := rowBucket("", start, counter)
	require.NoError(t, f.store.CommitClosure(t.Context(), testKey,
		Closure{Destination: destination, NextClose: start.Add(10 * time.Second), Buckets: []ClosedBucket{bucket}}))
	return bucket.Outcome.EffectKey
}

func (f *deliveryFixture) targetRows(t *testing.T) int { return countRows(t, f.target, "readings") }

func (f *deliveryFixture) state(t *testing.T, effectKey string) OutboxItem {
	t.Helper()
	item, err := f.store.GetOutbox(t.Context(), effectKey)
	require.NoError(t, err)
	return item
}

func TestSQLiteTargetReceiptIdentitySenderCommitsAndRecordsLocalReceipt(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, true)
	effect := f.enqueue(t, DedupeReceipt, "", t0, 9007199254740993)

	result, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, result.State)
	require.Equal(t, 1, f.targetRows(t))
	require.Equal(t, 1, countRows(t, f.target, "gw_effect_receipts"), "the destination holds its own receipt")

	item := f.state(t, effect)
	require.Equal(t, StateCommitted, item.State)
	require.False(t, item.CommittedAt.IsZero())
	require.Equal(t, 1, countRows(t, f.local, "wg_delivery_receipts"))
	var counter int64
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT counter FROM readings`).Scan(&counter))
	require.Equal(t, int64(9007199254740993), counter, "values survive outbox, sender and destination exactly")

	// Delivering a committed item again never touches the destination.
	again, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, again.State)
	require.Equal(t, 1, f.targetRows(t))
	require.EqualValues(t, 1, f.resolver.calls.Load())
}

func TestSQLiteTargetReceiptIdentityResponseLostRetriesAsOneEffect(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, true)
	effect := f.enqueue(t, DedupeReceipt, "", t0, 5)

	f.lose.Store(true)
	first, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateRetrying, first.State, "a receipt makes the ambiguous commit safe to retry")
	require.Equal(t, 1, f.targetRows(t), "the destination committed even though the response was lost")
	require.Equal(t, 0, countRows(t, f.local, "wg_delivery_receipts"), "no local proof until the destination confirms")

	f.lose.Store(false)
	second, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, second.State)
	require.Equal(t, 1, f.targetRows(t), "exactly one effect")
	require.Equal(t, StateCommitted, f.state(t, effect).State)
}

func TestSQLiteTargetReceiptIdentityResponseLostWithoutDedupeIsUnknownAndNeverResent(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueue(t, DedupeNone, "", t0, 5)

	f.lose.Store(true)
	first, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateUnknown, first.State)
	require.Equal(t, 1, f.targetRows(t))

	f.lose.Store(false)
	second, err := f.sender.Deliver(t.Context(), effect)
	require.ErrorIs(t, err, ErrNotDeliverable)
	require.Equal(t, StateUnknown, second.State, "an unknown outcome waits for reconciliation instead of blind insert")
	require.Equal(t, 1, f.targetRows(t))
	require.Equal(t, StateUnknown, f.state(t, effect).State)
	require.Equal(t, 0, countRows(t, f.local, "wg_delivery_receipts"), "unknown is never recorded as committed")
}

func TestSQLiteTargetReceiptIdentityLocalReceiptFailureDoesNotResendWithoutDedupe(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueue(t, DedupeNone, "", t0, 5)
	_, err := f.local.ExecContext(t.Context(), `CREATE TRIGGER no_local_receipt BEFORE INSERT ON wg_delivery_receipts
		BEGIN SELECT RAISE(ABORT, 'injected local failure'); END`)
	require.NoError(t, err)

	result, err := f.sender.Deliver(t.Context(), effect)
	require.ErrorIs(t, err, ErrLocalReceiptFailed)
	require.Equal(t, StateSending, result.State, "the local state could not be updated")
	require.Equal(t, 1, f.targetRows(t), "the destination did commit")

	// Interrupted-delivery resolution (also what a restart runs) must not turn
	// this into a retry: without dedupe the earlier attempt may have committed.
	resolved, err := f.store.ResolveInterrupted(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateUnknown, resolved)
	_, err = f.sender.Deliver(t.Context(), effect)
	require.ErrorIs(t, err, ErrNotDeliverable)
	require.Equal(t, 1, f.targetRows(t), "no second row")
}

func TestSQLiteTargetReceiptIdentityLocalReceiptFailureWithReceiptConvergesWithoutDuplicate(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, true)
	effect := f.enqueue(t, DedupeReceipt, "", t0, 5)
	_, err := f.local.ExecContext(t.Context(), `CREATE TRIGGER no_local_receipt BEFORE INSERT ON wg_delivery_receipts
		BEGIN SELECT RAISE(ABORT, 'injected local failure'); END`)
	require.NoError(t, err)
	_, err = f.sender.Deliver(t.Context(), effect)
	require.ErrorIs(t, err, ErrLocalReceiptFailed)

	resolved, err := f.store.ResolveInterrupted(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateRetrying, resolved, "a destination receipt makes the retry safe")
	_, err = f.local.ExecContext(t.Context(), `DROP TRIGGER no_local_receipt`)
	require.NoError(t, err)

	result, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, result.State)
	require.Equal(t, 1, f.targetRows(t))
	require.Equal(t, 1, countRows(t, f.local, "wg_delivery_receipts"))
}

func TestSQLiteTargetReceiptIdentityUniqueKeyDestination(t *testing.T) {
	f := newDeliveryFixture(t, uniqueTargetTable, false)
	effect := f.enqueue(t, DedupeUniqueKey, "record_id", t0, 5)

	f.lose.Store(true)
	first, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateRetrying, first.State)
	f.lose.Store(false)
	second, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, second.State)
	require.Equal(t, 1, f.targetRows(t))
}

func TestSQLiteTargetReceiptIdentityPreCommitFailureRetriesAndThenCommitsOnce(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueue(t, DedupeNone, "", t0, 5)
	transient := errors.New("connection reset by peer")
	f.faults.ExecError.Store(&transient)

	result, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateRetrying, result.State)
	item := f.state(t, effect)
	require.Equal(t, 1, item.RetryCount)
	require.Equal(t, "insert-failed", item.LastErrorCode, "only a safe code is stored, never driver text")
	require.Zero(t, f.targetRows(t))

	f.faults.ExecError.Store(nil)
	result, err = f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, result.State)
	require.Equal(t, 1, f.targetRows(t))
}

func TestSQLiteTargetReceiptIdentityRetriesAreBoundedAndKeepTheData(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueue(t, DedupeNone, "", t0, 5)
	transient := errors.New("connection reset by peer")
	f.faults.ExecError.Store(&transient)
	f.sender = NewSender(f.store, f.resolver, SenderConfig{MaxRetries: 2, Backoff: func(retry int) (string, time.Time) {
		if retry >= 2 {
			return StateBlocked, time.Now().UTC()
		}
		return StateRetrying, time.Now().UTC().Add(-time.Second)
	}})

	for i := 0; i < 2; i++ {
		_, deliverErr := f.sender.Deliver(t.Context(), effect)
		require.NoError(t, deliverErr)
	}
	require.Equal(t, StateBlocked, f.state(t, effect).State, "retry exhaustion blocks instead of deleting")
	require.Equal(t, 1, countRows(t, f.local, "wg_delivery_outbox"), "the accepted row is still there")
	_, err := f.sender.Deliver(t.Context(), effect)
	require.ErrorIs(t, err, ErrNotDeliverable)
}

func TestSQLiteTargetReceiptIdentityDigestMismatchBlocksWithoutOverwriting(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, true)
	effect := f.enqueue(t, DedupeReceipt, "", t0, 5)
	_, err := f.target.ExecContext(t.Context(), `INSERT INTO gw_effect_receipts VALUES (?, 'someone-elses-digest', 'x')`, effect)
	require.NoError(t, err)

	result, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateBlocked, result.State)
	require.Equal(t, "destination-identity-conflict", f.state(t, effect).LastErrorCode)
	require.Zero(t, f.targetRows(t))
	var digest string
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT payload_digest FROM gw_effect_receipts`).Scan(&digest))
	require.Equal(t, "someone-elses-digest", digest, "the destination receipt is never overwritten")
}

func TestSQLiteTargetReceiptIdentityBlockedTargetAndCorruptPayloadNeverSend(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueue(t, DedupeNone, "", t0, 5)
	f.resolver.err = errors.Join(ErrTargetBlocked, errors.New("credential rejected"))
	result, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateBlocked, result.State)
	require.Equal(t, "target-blocked", f.state(t, effect).LastErrorCode)
	require.Zero(t, f.targetRows(t))

	other := newDeliveryFixture(t, plainTargetTable, false)
	corrupt := other.enqueue(t, DedupeNone, "", t0, 5)
	_, err = other.local.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET payload = replace(payload, 'batch-001', 'tampered')`)
	require.NoError(t, err)
	result, err = other.sender.Deliver(t.Context(), corrupt)
	require.NoError(t, err)
	require.Equal(t, StateBlocked, result.State)
	require.Equal(t, "payload-digest-mismatch", other.state(t, corrupt).LastErrorCode)
	require.Zero(t, other.targetRows(t))
	require.EqualValues(t, 0, other.resolver.calls.Load(), "a corrupt payload is detected before any connection is opened")
}

func TestSQLiteTargetReceiptIdentityOnlyPendingAndRetryingItemsMaySend(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueue(t, DedupeNone, "", t0, 5)
	_, err := f.store.BeginDelivery(t.Context(), effect, "worker-1", time.Minute)
	require.NoError(t, err)
	_, err = f.store.BeginDelivery(t.Context(), effect, "worker-2", time.Minute)
	require.ErrorIs(t, err, ErrNotDeliverable, "a claimed item cannot be claimed again")
	_, err = f.store.BeginDelivery(t.Context(), "missing", "worker-1", time.Minute)
	require.ErrorIs(t, err, ErrOutboxItemNotFound)

	_, err = f.sender.Deliver(t.Context(), "missing")
	require.ErrorIs(t, err, ErrOutboxItemNotFound)
}
