package groupdelivery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReclaimRollsBackReceiptWhenOutboxDeleteFails(t *testing.T) {
	store, db, _ := newTestStore(t)
	done := rowBucket("", t0, 1)
	require.NoError(t, store.CommitClosure(t.Context(), testKey,
		Closure{Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{done}}))
	claim, err := store.BeginDelivery(t.Context(), done.Outcome.EffectKey, "worker-1", time.Minute)
	require.NoError(t, err)
	require.NoError(t, store.CompleteDelivery(t.Context(), claim, "d"))
	_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET committed_at='2026-01-01T00:00:00.000000000Z'`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `CREATE TRIGGER fail_reclaim BEFORE DELETE ON wg_delivery_outbox BEGIN SELECT RAISE(ABORT,'reclaim failed'); END`)
	require.NoError(t, err)

	_, err = store.Reclaim(t.Context(), 0)
	require.Error(t, err)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_outbox"))
	require.Equal(t, 1, countRows(t, db, "wg_delivery_receipts"), "receipt deletion must roll back with the outbox")
}
