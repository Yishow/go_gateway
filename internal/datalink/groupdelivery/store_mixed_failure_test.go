package groupdelivery

import (
	"encoding/json"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/sqltx"
	"go-gateway/internal/testutil/lostcommit"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedClosureLostAckRestoresAcceptedProgress(t *testing.T) {
	store, db, path := newTestStore(t)
	_, c := mixedStorePlan(t, store)
	lose := &atomic.Bool{}
	wrapped := lostcommit.Open(path, lose)
	t.Cleanup(func() { wrapped.Close() })
	lose.Store(true)
	err := NewStore(wrapped).CommitClosure(t.Context(), testKey, c)
	require.ErrorIs(t, err, lostcommit.ErrResponseLost)
	require.ErrorIs(t, err, sqltx.ErrCommitOutcomeUnknown)
	lose.Store(false)
	r, err := NewStore(openStoreDB(t, path)).Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, c.Mixed.After, *r.Mixed)
	require.Len(t, r.Samples, 1)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, c))
	require.Equal(t, 1, countRows(t, db, "wg_delivery_mixed_outcomes"))
}
func TestMixedClosureCorruptionFailsClosed(t *testing.T) {
	for _, change := range []string{
		`UPDATE wg_delivery_mixed_progress SET checkpoint_digest='bad'`,
		`UPDATE wg_delivery_mixed_progress SET descriptor_digest='bad'`,
		`UPDATE wg_delivery_mixed_progress SET closed_through='2026-01-01T00:01:00.000000000Z'`,
		`INSERT INTO wg_delivery_checkpoints(group_id,group_revision,next_close,updated_at) VALUES('group-G','rev-1','2026-01-01T00:00:15Z','now')`,
	} {
		t.Run(change, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			_, c := mixedStorePlan(t, store)
			require.NoError(t, store.CommitClosure(t.Context(), testKey, c))
			_, err := db.ExecContext(t.Context(), change)
			require.NoError(t, err)
			_, err = store.Restore(t.Context(), testKey)
			require.Error(t, err)
			require.Error(t, store.CommitClosure(t.Context(), testKey, c))
		})
	}
}
func TestMixedClosureCannotSkipUnionWindowsOrChangeMethod(t *testing.T) {
	for _, kind := range []string{"missing union", "changed method"} {
		t.Run(kind, func(t *testing.T) {
			store, _, _ := newTestStore(t)
			a, c := mixedStorePlan(t, store)
			if kind == "missing union" {
				cp, err := a.MixedCheckpointForPlan(a.Plan(at(60)))
				require.NoError(t, err)
				c.Mixed.After = cp
			} else {
				c.Buckets[0].Outcome.Members[0].Window.Calculation = snapshot.CalculationMean
			}
			require.Error(t, store.CommitClosure(t.Context(), testKey, c))
		})
	}
}
func TestMixedClosureJSONBoundsTamperingWithRehashedBytesRejected(t *testing.T) {
	store, db, _ := newTestStore(t)
	_, c := mixedStorePlan(t, store)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, c))
	cp := c.Mixed.After
	for key, m := range cp.Members {
		m.Start = at(1)
		cp.Members[key] = m
		break
	}
	raw, err := json.Marshal(cp)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_mixed_progress SET checkpoint=?,checkpoint_digest=?`, string(raw), handoffHash(string(raw)))
	require.NoError(t, err)
	_, err = store.Restore(t.Context(), testKey)
	require.Error(t, err)
}

func TestMixedOutcomeCodecRetainsLatestAndSealedCounterAbsence(t *testing.T) {
	store, db, _ := newTestStore(t)
	_, c := mixedStorePlan(t, store)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, c))
	var raw, digest string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT outcome,outcome_digest FROM wg_delivery_mixed_outcomes`).Scan(&raw, &digest))
	out, err := DecodeMixedOutcome([]byte(raw), digest)
	require.NoError(t, err)
	require.Equal(t, c.Buckets[0].Outcome, out)
}

func TestMixedJournalRefusesForeignWindowAndLegacyCursor(t *testing.T) {
	store, db, _ := newTestStore(t)
	_, c := mixedStorePlan(t, store)
	sample := sampleAt("wrong-window", `["d","p","fast"]`, 20, 99)
	require.Error(t, store.AppendSample(t.Context(), testKey, t0, sample), "member20 belongs to[15,30), not common zero")
	legacy := Closure{Destination: testDestination, NextClose: at(15)}
	require.Error(t, store.CommitClosure(t.Context(), testKey, legacy), "mixed descriptor cannot write legacy cursor")
	require.NoError(t, store.CommitClosure(t.Context(), testKey, c))
	require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"))
}
