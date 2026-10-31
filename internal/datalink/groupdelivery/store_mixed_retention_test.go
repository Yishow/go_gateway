package groupdelivery

import (
	"go-gateway/internal/datalink/snapshot"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mixedRowClosure(t *testing.T, store *Store) Closure {
	_, c := mixedStorePlan(t, store)
	b := &c.Buckets[0]
	b.Outcome.Kind = snapshot.OutcomeRow
	b.Outcome.EffectKey, _ = snapshot.EffectKey(c.Destination.Scope, b.Outcome.RecordID)
	row := rowBucket("", t0, 9007199254740993).Row
	row.RecordID = b.Outcome.RecordID
	row.EffectKey = b.Outcome.EffectKey
	row.BucketStart = time.Time{}
	row.DueAt = b.Outcome.DueAt
	b.Row = row
	return c
}
func TestMixedRetentionPinsUndeliveredOutcomeAndCurrentProgress(t *testing.T) {
	store, db, _ := newTestStore(t)
	store = store.WithClock(func() time.Time { return at(100) })
	c := mixedRowClosure(t, store)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, c))
	store = store.WithClock(func() time.Time { return at(200) })
	_, err := store.Reclaim(t.Context(), time.Second)
	require.NoError(t, err)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_mixed_outcomes"))
	require.Equal(t, 1, countRows(t, db, "wg_delivery_mixed_progress"))
	require.Equal(t, 1, countRows(t, db, "wg_runtime_versions"))
	_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET state='sql_committed',committed_at=?`, timestamp(at(100)))
	require.NoError(t, err)
	_, err = store.Reclaim(t.Context(), time.Second)
	require.NoError(t, err)
	require.Zero(t, countRows(t, db, "wg_delivery_mixed_outcomes"))
	require.Zero(t, countRows(t, db, "wg_delivery_outbox"))
	require.Equal(t, 1, countRows(t, db, "wg_delivery_mixed_progress"))
	r, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.NotNil(t, r.Mixed)
	require.Len(t, r.Samples, 1)
}
func TestMixedRowAtomicFailureAtEveryInsert(t *testing.T) {
	for _, table := range []string{"wg_delivery_mixed_outcomes", "wg_delivery_outbox", "wg_delivery_mixed_progress"} {
		t.Run(table, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			c := mixedRowClosure(t, store)
			failOn(t, db, table)
			require.Error(t, store.CommitClosure(t.Context(), testKey, c))
			for _, name := range []string{"wg_delivery_mixed_outcomes", "wg_delivery_outbox", "wg_delivery_mixed_progress"} {
				require.Zero(t, countRows(t, db, name))
			}
			r, err := store.Restore(t.Context(), testKey)
			require.NoError(t, err)
			require.Len(t, r.Samples, 2)
		})
	}
}
func TestMixedOutboxRejectsMissingOrForeignDispatch(t *testing.T) {
	for _, change := range []string{`UPDATE wg_delivery_outbox SET due_at=NULL`, `UPDATE wg_delivery_outbox SET due_at='2026-01-01T00:00:30Z'`, `UPDATE wg_delivery_outbox SET workspace_id='foreign'`, `UPDATE wg_delivery_outbox SET bucket_start='2026-01-01T00:00:00Z'`} {
		t.Run(change, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			c := mixedRowClosure(t, store)
			require.NoError(t, store.CommitClosure(t.Context(), testKey, c))
			_, err := db.ExecContext(t.Context(), change)
			require.NoError(t, err)
			_, err = store.GetOutbox(t.Context(), c.Buckets[0].Outcome.EffectKey)
			require.Error(t, err)
		})
	}
}
