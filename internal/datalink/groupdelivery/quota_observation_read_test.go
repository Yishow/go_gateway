package groupdelivery

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type quotaScopeReader interface {
	QuotaStatusForScope(context.Context, GroupKey, QuotaScope) (QuotaStatus, error)
}

// The optional interface makes the initial missing-read-path test behavioral,
// rather than failing to compile before the additive method exists.
func exactQuotaReader(t *testing.T, store *Store) quotaScopeReader {
	t.Helper()
	reader, ok := any(store).(quotaScopeReader)
	require.True(t, ok, "exact-scope quota evidence is unavailable; aggregate most-constrained status cannot substitute")
	return reader
}

func TestQuotaStatusForScopePreservesHistoricalGroupAndGlobalUsageReadOnly(t *testing.T) {
	base, db, _ := newTestStore(t)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	current := testKey
	current.GroupRevision = "rev-2"
	for i, key := range []GroupKey{testKey, current, otherKey} {
		require.NoError(t, base.AppendSample(t.Context(), key, at(20), sampleAt("pending-"+itoa(i), "temp", 23, int64(i))))
	}
	enqueueFor(t, base, testKey, testDestination, "historical", at(0), 1)
	enqueueFor(t, base, otherKey, testDestination, "other-group", at(0), 2)
	committed := enqueueFor(t, base, testKey, testDestination, "committed", at(10), 3)
	claim, err := base.BeginDelivery(t.Context(), committed, "quota-read-fixture", time.Minute)
	require.NoError(t, err)
	item, err := base.GetOutbox(t.Context(), committed)
	require.NoError(t, err)
	require.NoError(t, base.CompleteDelivery(t.Context(), claim, item.PayloadDigest))
	groupBytes := quotaFixtureUsage(t, db, testKey.GroupID)
	globalBytes := quotaFixtureUsage(t, db, "")
	require.Greater(t, globalBytes, groupBytes)
	store, err := base.WithQuota(QuotaConfig{GroupMaxBytes: 1, GlobalMaxBytes: globalBytes * 2, WarningRatio: 0.5, HardRatio: 0.9})
	require.NoError(t, err)
	reader := exactQuotaReader(t, store)
	var before, after int64
	require.NoError(t, db.QueryRow(`SELECT total_changes()`).Scan(&before))
	group, err := reader.QuotaStatusForScope(t.Context(), current, QuotaScopeGroup)
	require.NoError(t, err)
	require.Equal(t, QuotaScopeGroup, group.Scope)
	require.Equal(t, current.GroupID, group.GroupID)
	require.Equal(t, groupBytes, group.UsedBytes, "the requested revision does not exclude older accepted inputs/outbox")
	require.EqualValues(t, 1, group.MaxBytes)
	require.Equal(t, QuotaHard, group.State)
	require.True(t, group.IntakeRefused)
	global, err := reader.QuotaStatusForScope(t.Context(), current, QuotaScopeGlobal)
	require.NoError(t, err)
	require.Equal(t, QuotaScopeGlobal, global.Scope)
	require.Empty(t, global.GroupID)
	require.Equal(t, globalBytes, global.UsedBytes)
	require.Equal(t, globalBytes*2, global.MaxBytes)
	require.Equal(t, QuotaWarning, global.State)
	require.False(t, global.IntakeRefused)
	aggregate, err := store.QuotaStatus(t.Context(), current)
	require.NoError(t, err)
	require.Equal(t, group, aggregate, "aggregate selects hard group; exact global must not borrow that scope")
	require.NoError(t, db.QueryRow(`SELECT total_changes()`).Scan(&after))
	require.Equal(t, before, after, "the reader must not mutate journals, outbox, receipts or quota state")
	require.Zero(t, store.usage.measurements.Load(), "pure reads must not populate the intake usage cache")
}

func quotaFixtureUsage(t *testing.T, db *sql.DB, groupID string) int64 {
	t.Helper()
	var used int64
	require.NoError(t, db.QueryRow(`SELECT COALESCE(SUM(bytes),0) FROM (
		SELECT LENGTH(CAST(payload AS BLOB)) AS bytes FROM wg_delivery_samples WHERE consumed=0 AND (?='' OR group_id=?)
		UNION ALL SELECT LENGTH(CAST(payload AS BLOB)) FROM wg_delivery_outbox WHERE state!='sql_committed' AND (?='' OR group_id=?)
	)`, groupID, groupID, groupID, groupID).Scan(&used))
	return used
}

func TestQuotaStatusForScopePreservesObservedZeroAndNativeOvershoot(t *testing.T) {
	base, _, _ := newTestStore(t)
	store, err := base.WithQuota(QuotaConfig{GroupMaxBytes: 1, GlobalMaxBytes: 1, WarningRatio: 0.5, HardRatio: 0.9})
	require.NoError(t, err)
	reader := exactQuotaReader(t, store)
	for _, scope := range []QuotaScope{QuotaScopeGroup, QuotaScopeGlobal} {
		q, err := reader.QuotaStatusForScope(t.Context(), testKey, scope)
		require.NoError(t, err)
		require.True(t, q.Configured)
		require.Equal(t, scope, q.Scope)
		require.Zero(t, q.UsedBytes)
		require.Equal(t, QuotaOK, q.State)
		require.False(t, q.IntakeRefused)
	}
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(0), sampleAt("overshoot", "temp", 3, 1)), "native policy permits the sample that crosses the limit")
	require.ErrorIs(t, store.AppendSample(t.Context(), testKey, at(0), sampleAt("refused", "temp", 4, 2)), ErrQuotaExceeded)
	for _, scope := range []QuotaScope{QuotaScopeGroup, QuotaScopeGlobal} {
		q, err := reader.QuotaStatusForScope(t.Context(), testKey, scope)
		require.NoError(t, err)
		want, err := store.Usage(t.Context(), q.GroupID)
		require.NoError(t, err)
		require.Equal(t, want, q.UsedBytes)
		require.Greater(t, q.UsedBytes, q.MaxBytes, "do not clamp an actual overshoot")
		require.Equal(t, QuotaHard, q.State)
		require.True(t, q.IntakeRefused)
	}
}

func TestQuotaStatusForScopeUnconfiguredNeverBorrowsOtherScope(t *testing.T) {
	base, db, _ := newTestStore(t)
	groupOnly, err := base.WithQuota(QuotaConfig{GroupMaxBytes: 1})
	require.NoError(t, err)
	globalOnly, err := base.WithQuota(QuotaConfig{GlobalMaxBytes: 1})
	require.NoError(t, err)
	// A missing policy needs no measurement, even when a DB read cannot succeed.
	require.NoError(t, db.Close())
	for _, tc := range []struct {
		store *Store
		scope QuotaScope
	}{{base, QuotaScopeGroup}, {base, QuotaScopeGlobal}, {groupOnly, QuotaScopeGlobal}, {globalOnly, QuotaScopeGroup}} {
		q, err := exactQuotaReader(t, tc.store).QuotaStatusForScope(t.Context(), testKey, tc.scope)
		require.NoError(t, err)
		require.Equal(t, QuotaStatus{}, q)
	}
}

func TestQuotaStatusForScopeRejectsInvalidInputWithoutSuccessfulZero(t *testing.T) {
	base, _, _ := newTestStore(t)
	store, err := base.WithQuota(QuotaConfig{GlobalMaxBytes: 100, GroupMaxBytes: 100})
	require.NoError(t, err)
	reader := exactQuotaReader(t, store)
	for _, scope := range []QuotaScope{"", "GLOBAL", "other"} {
		q, err := reader.QuotaStatusForScope(t.Context(), testKey, scope)
		require.ErrorIs(t, err, ErrInvalidQuota)
		require.Equal(t, QuotaStatus{}, q)
	}
	for _, key := range []GroupKey{{}, {GroupID: testKey.GroupID, GroupRevision: testKey.GroupRevision}, {WorkspaceID: testKey.WorkspaceID, GroupRevision: testKey.GroupRevision}, {WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID}} {
		q, err := reader.QuotaStatusForScope(t.Context(), key, QuotaScopeGlobal)
		require.ErrorIs(t, err, ErrInvalidGroupKey)
		require.Equal(t, QuotaStatus{}, q)
	}
	for _, invalid := range []*Store{nil, NewStore(nil)} {
		q, err := exactQuotaReader(t, invalid).QuotaStatusForScope(t.Context(), testKey, QuotaScopeGlobal)
		require.Error(t, err)
		require.Equal(t, QuotaStatus{}, q)
	}
}

func TestQuotaStatusForScopeReadFailureAndCancellationRemainErrors(t *testing.T) {
	for _, failure := range []string{"closed-db", "missing-table", "cancelled"} {
		t.Run(failure, func(t *testing.T) {
			base, db, _ := newTestStore(t)
			store, err := base.WithQuota(QuotaConfig{GlobalMaxBytes: 100})
			require.NoError(t, err)
			reader := exactQuotaReader(t, store)
			ctx := t.Context()
			switch failure {
			case "closed-db":
				require.NoError(t, db.Close())
			case "missing-table":
				_, err := db.Exec(`DROP TABLE wg_delivery_samples`)
				require.NoError(t, err)
			case "cancelled":
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}
			q, err := reader.QuotaStatusForScope(ctx, testKey, QuotaScopeGlobal)
			require.Error(t, err)
			require.Equal(t, QuotaStatus{}, q)
			if failure == "cancelled" {
				require.ErrorIs(t, err, context.Canceled)
			}
		})
	}
}
