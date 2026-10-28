package groupdelivery

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/testutil/lostcommit"
)

func TestClosureOwnershipCutoverWinsWithoutReceipts(t *testing.T) {
	for _, mode := range []string{"empty-cursor", "empty-bucket", "pre-boundary-only"} {
		t.Run(mode, func(t *testing.T) {
			store, db, path := newTestStore(t)
			if mode == "pre-boundary-only" {
				require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sampleAt("original", ownershipMember, 8, 20)))
			}
			before, err := store.Restore(t.Context(), testKey)
			require.NoError(t, err)
			require.Equal(t, HandoffView{}, before.HandoffView)
			closure := Closure{HandoffView: before.HandoffView, Destination: testDestination, NextClose: at(20)}
			if mode != "empty-cursor" {
				closure.Buckets = []ClosedBucket{rowBucket("", at(10), 20)}
				closure.Baselines = []CounterBaseline{baselineAt(t, "temp", at(20), 20)}
			}
			entered, release := make(chan struct{}), make(chan struct{})
			faults := &lostcommit.Faults{}
			hook := func(query string) error {
				if strings.Contains(query, "UPDATE wg_delivery_checkpoints SET updated_at") {
					close(entered)
					<-release
				}
				return nil
			}
			faults.Reject.Store(&hook)
			concurrent := lostcommit.OpenFaults(path+"?_pragma=busy_timeout(5000)", faults)
			defer concurrent.Close()
			result := make(chan error, 1)
			go func() { result <- NewStore(concurrent).CommitClosure(t.Context(), testKey, closure) }()
			<-entered
			seedOwnershipCutover(t, db, "f-cutover-wins", testKey, otherKey, t0.Add(12500*time.Millisecond))
			close(release)
			require.ErrorIs(t, <-result, ErrOwnershipConflict, "new durable owner epoch must fence even an empty stale closure")
			for _, table := range []string{"wg_delivery_checkpoints", "wg_delivery_buckets", "wg_delivery_outbox", "wg_delivery_counter_baselines", "wg_delivery_sample_handoffs"} {
				require.Zero(t, countRows(t, db, table), table)
			}
			require.Zero(t, openCount(t, db, "wg_delivery_samples", "consumed=1"))
			after, err := store.Restore(t.Context(), testKey)
			require.NoError(t, err)
			require.Equal(t, before.Samples, after.Samples)
		})
	}
}

func ownershipPartialClosure(t *testing.T, restored Restored, end time.Time) Closure {
	t.Helper()
	bucket := rowBucket("", at(10), 20)
	bucket.Outcome.WindowStart, bucket.Outcome.WindowEnd, bucket.Outcome.Truncated = at(10), end, "end"
	bucket.Outcome.Members[0].Sample.ObservedAt = at(11)
	return Closure{OwnershipView: restored.OwnershipView, HandoffView: restored.HandoffView,
		Destination: testDestination, NextClose: at(20), Buckets: []ClosedBucket{bucket},
		Baselines: []CounterBaseline{baselineAt(t, "temp", end, 20)}}
}

func TestClosureOwnershipRefreshedPartialAdvancesAtomically(t *testing.T) {
	store, db, path := newTestStore(t)
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(10), sampleAt("original", ownershipMember, 11, 20)))
	before, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	end := t0.Add(12500 * time.Millisecond)
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, end)
	restored, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.NotEqual(t, before.OwnershipView, restored.OwnershipView)
	require.NoError(t, restored.OwnershipView.ValidateBounds(testKey, t0, time.Time{}, at(20), end, 10*time.Second))
	require.NoError(t, restored.OwnershipView.ValidateBounds(testKey, t0, time.Time{}, at(10), time.Time{}, 10*time.Second), "an earlier disabled aligned stop must not extend to the later cutover")
	require.ErrorIs(t, restored.OwnershipView.ValidateBounds(testKey, t0, time.Time{}, at(20), time.Time{}, 10*time.Second), ErrOwnershipConflict, "the cutover-containing window still needs its exact actual end")
	require.ErrorIs(t, restored.OwnershipView.ValidateBounds(testKey, t0, time.Time{}, at(15), end, 10*time.Second), ErrOwnershipConflict, "an earlier stop must remain aligned")
	require.NoError(t, restored.OwnershipView.ValidateBounds(testKey, at(20), time.Time{}, at(20), end, 10*time.Second), "completed checkpoint still reconstructs a stopped owner")
	require.ErrorIs(t, restored.OwnershipView.ValidateBounds(testKey, t0, time.Time{}, time.Time{}, time.Time{}, 10*time.Second), ErrOwnershipConflict)
	require.ErrorIs(t, restored.OwnershipView.ValidateBounds(testKey, t0, time.Time{}, at(30), end, 10*time.Second), ErrOwnershipConflict)
	require.ErrorIs(t, restored.OwnershipView.ValidateBounds(otherKey, t0, time.Time{}, at(20), end, 10*time.Second), ErrOwnershipConflict)
	require.ErrorIs(t, restored.OwnershipView.ValidateBounds(testKey, t0, time.Time{}, at(15), end, 5*time.Second), ErrOwnershipConflict)
	closure := ownershipPartialClosure(t, restored, end)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closure), "nominal NextClose 20 may exceed actual end 12.5")
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closure), "unchanged authoritative retry is idempotent")
	for _, table := range []string{"wg_delivery_checkpoints", "wg_delivery_buckets", "wg_delivery_outbox", "wg_delivery_counter_baselines"} {
		require.Equal(t, 1, countRows(t, db, table), table)
	}
	require.Equal(t, 1, openCount(t, db, "wg_delivery_samples", "consumed=1"))
	require.Zero(t, countRows(t, db, "wg_delivery_sample_handoffs"))
	reopened := NewStore(openStoreDB(t, path))
	after, err := reopened.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, at(20), after.NextClose)
	require.Equal(t, restored.OwnershipView, after.OwnershipView)
	require.Equal(t, closure.Baselines, after.Baselines)
	require.Empty(t, after.Samples)
	target, err := reopened.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.NoError(t, target.OwnershipView.ValidateBounds(otherKey, at(10), end, time.Time{}, time.Time{}, 10*time.Second))
	require.ErrorIs(t, target.OwnershipView.ValidateBounds(otherKey, at(10), time.Time{}, time.Time{}, time.Time{}, 10*time.Second), ErrOwnershipConflict)
	require.ErrorIs(t, target.OwnershipView.ValidateBounds(otherKey, t0, end, time.Time{}, time.Time{}, 10*time.Second), ErrOwnershipConflict)
}

func TestClosureOwnershipFreshViewDoesNotAuthorizeWrongPlan(t *testing.T) {
	for _, mode := range []string{"empty-past-end", "untruncated", "wrong-actual-end", "outside-sample", "outside-baseline"} {
		t.Run(mode, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			end := t0.Add(12500 * time.Millisecond)
			seedOwnershipCutover(t, db, "f-2", testKey, otherKey, end)
			restored, err := store.Restore(t.Context(), testKey)
			require.NoError(t, err)
			closure := ownershipPartialClosure(t, restored, end)
			switch mode {
			case "empty-past-end":
				closure.Buckets, closure.Baselines, closure.NextClose = nil, nil, at(30)
			case "untruncated":
				closure.Buckets[0].Outcome.Truncated = ""
			case "wrong-actual-end":
				closure.Buckets[0].Outcome.WindowEnd = at(13)
			case "outside-sample":
				closure.Buckets[0].Outcome.Members[0].Sample.ObservedAt = end
			case "outside-baseline":
				closure.Baselines = []CounterBaseline{baselineAt(t, "temp", at(20), 20)}
			}
			require.ErrorIs(t, store.CommitClosure(t.Context(), testKey, closure), ErrOwnershipConflict)
			for _, table := range []string{"wg_delivery_checkpoints", "wg_delivery_buckets", "wg_delivery_outbox", "wg_delivery_counter_baselines"} {
				require.Zero(t, countRows(t, db, table), table)
			}
		})
	}
}

func TestClosureOwnershipMissingAuthorityIsScoped(t *testing.T) {
	for _, mutation := range []string{
		`DELETE FROM write_group_ownership_fences`,
		`DELETE FROM write_group_revision_ends`,
		`UPDATE write_group_versions SET payload=json_remove(payload,'$.row_policy.interval_seconds') WHERE group_id='group-G'`,
	} {
		t.Run(mutation, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			end := t0.Add(12500 * time.Millisecond)
			seedOwnershipCutover(t, db, "f-2", testKey, otherKey, end)
			restored, err := store.Restore(t.Context(), testKey)
			require.NoError(t, err)
			_, err = db.Exec(mutation)
			require.NoError(t, err)
			err = store.CommitClosure(t.Context(), testKey, ownershipPartialClosure(t, restored, end))
			require.ErrorIs(t, err, ErrOwnershipUnavailable)
			require.Equal(t, "recording-ownership-unavailable", err.Error())
			_, err = store.Restore(t.Context(), testKey)
			require.ErrorIs(t, err, ErrOwnershipUnavailable)
			for _, table := range []string{"wg_delivery_checkpoints", "wg_delivery_buckets", "wg_delivery_outbox", "wg_delivery_counter_baselines"} {
				require.Zero(t, countRows(t, db, table), table)
			}
			healthy := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "healthy", GroupRevision: "healthy-1"}
			require.NoError(t, store.CommitClosure(t.Context(), healthy, Closure{Destination: testDestination, NextClose: at(30)}))
			require.Equal(t, 1, countRows(t, db, "wg_delivery_checkpoints"))
		})
	}
}

func TestClosureOwnershipLaterEpochInvalidatesPriorView(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedOwnershipCutover(t, db, "A-B", testKey, otherKey, t0.Add(12500*time.Millisecond))
	before, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	hint, err := store.OwnershipView(t.Context(), otherKey)
	require.NoError(t, err)
	require.Equal(t, before.OwnershipView, hint)
	third := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "owner-C", GroupRevision: "rev-C"}
	seedOwnershipCutover(t, db, "B-C", otherKey, third, t0.Add(22500*time.Millisecond))
	stale := Closure{OwnershipView: before.OwnershipView, HandoffView: before.HandoffView, Destination: testDestination, NextClose: at(30)}
	require.ErrorIs(t, store.CommitClosure(t.Context(), otherKey, stale), ErrOwnershipConflict)
	require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"))
	after, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.NotEqual(t, before.OwnershipView, after.OwnershipView)
	require.Equal(t, before.HandoffView, after.HandoffView, "no receipts changed in this race")
	hint, err = store.OwnershipView(t.Context(), otherKey)
	require.NoError(t, err)
	require.Equal(t, after.OwnershipView, hint)
	require.NoError(t, after.OwnershipView.ValidateBounds(otherKey, at(10), t0.Add(12500*time.Millisecond), at(30), t0.Add(22500*time.Millisecond), 10*time.Second))
	stale.OwnershipView = after.OwnershipView // no assembler: explicitly refreshed empty Store fixture
	require.NoError(t, store.CommitClosure(t.Context(), otherKey, stale))
	require.Equal(t, 1, countRows(t, db, "wg_delivery_checkpoints"))
	_, err = store.OwnershipView(t.Context(), GroupKey{})
	require.ErrorIs(t, err, ErrInvalidGroupKey)
}

func TestClosureOwnershipRefusalPreservesExistingProgress(t *testing.T) {
	store, db, _ := newTestStore(t)
	baseline := baselineAt(t, "temp", at(10), 10)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closureWith(t, at(10), 10, baseline)))
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(10), sampleAt("still-open", ownershipMember, 11, 20)))
	before, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	var originalPayload, originalDigest string
	require.NoError(t, db.QueryRow(`SELECT payload,payload_digest FROM wg_delivery_outbox`).Scan(&originalPayload, &originalDigest))
	seedOwnershipCutover(t, db, "f-after-progress", testKey, otherKey, t0.Add(12500*time.Millisecond))
	stale := closureWith(t, at(20), 20, baselineAt(t, "temp", at(20), 20))
	stale.OwnershipView = before.OwnershipView
	require.ErrorIs(t, store.CommitClosure(t.Context(), testKey, stale), ErrOwnershipConflict)
	after, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, before.NextClose, after.NextClose)
	require.Equal(t, before.Baselines, after.Baselines)
	require.Equal(t, before.Samples, after.Samples)
	var payload, digest string
	require.NoError(t, db.QueryRow(`SELECT payload,payload_digest FROM wg_delivery_outbox`).Scan(&payload, &digest))
	require.Equal(t, originalPayload, payload)
	require.Equal(t, originalDigest, digest)
	for _, table := range []string{"wg_delivery_checkpoints", "wg_delivery_buckets", "wg_delivery_outbox", "wg_delivery_counter_baselines"} {
		require.Equal(t, 1, countRows(t, db, table), table)
	}
	require.Zero(t, openCount(t, db, "wg_delivery_samples", "consumed=1"))
}

func TestClosureOwnershipZeroWidthOwnerCannotProduceBucket(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedOwnershipCutover(t, db, "A-B", testKey, otherKey, at(10))
	third := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "owner-C", GroupRevision: "rev-C"}
	seedOwnershipCutover(t, db, "B-C", otherKey, third, at(10))
	restored, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.NoError(t, restored.OwnershipView.ValidateBounds(otherKey, at(10), time.Time{}, at(10), at(10), 10*time.Second), "an empty owner can restore as already stopped")
	bucket := silentBucket("", t0)
	bucket.Outcome.WindowStart, bucket.Outcome.WindowEnd, bucket.Outcome.Truncated = at(10), at(10), "start"
	closure := Closure{OwnershipView: restored.OwnershipView, HandoffView: restored.HandoffView, Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{bucket}}
	require.ErrorIs(t, store.CommitClosure(t.Context(), otherKey, closure), ErrOwnershipConflict, "a zero-width owner cannot invent a closed bucket")
	require.Zero(t, countRows(t, db, "wg_delivery_buckets"))
	require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"))
	require.Zero(t, countRows(t, db, "wg_delivery_outbox"))
	closure.Buckets = nil
	require.NoError(t, store.CommitClosure(t.Context(), otherKey, closure), "an empty cursor may record the stop without producing a bucket")
}
