package groupdelivery

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

// journalBytes is the payload size the quota measures for one sample.
func journalBytes(t *testing.T, db *sql.DB) int64 {
	t.Helper()
	var n sql.NullInt64
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT SUM(LENGTH(payload)) FROM wg_delivery_samples WHERE consumed = 0`).Scan(&n))
	return n.Int64
}

// uniformSample keeps every payload the same size so tests can reason in units.
func uniformSample(i int) snapshot.Sample {
	return sampleAt("sample-"+string(rune('a'+i/26))+string(rune('a'+i%26)), "temp", 3, int64(1000+i%9000))
}

func appendN(t *testing.T, store *Store, key GroupKey, from, to int) error {
	t.Helper()
	for i := from; i < to; i++ {
		if err := store.AppendSample(t.Context(), key, t0, uniformSample(i)); err != nil {
			return err
		}
	}
	return nil
}

// appendUntilRefused appends until the store refuses and returns how many were accepted.
func appendUntilRefused(t *testing.T, store *Store, key GroupKey) (accepted int, refusal error) {
	t.Helper()
	for i := 200; i < 400; i++ {
		if err := store.AppendSample(t.Context(), key, t0, uniformSample(i)); err != nil {
			return accepted, err
		}
		accepted++
	}
	return accepted, nil
}

func TestQuotaRejectsNewAckStatusFollowsMeasuredUsage(t *testing.T) {
	store, db, _ := newTestStore(t)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, uniformSample(0)))
	one := journalBytes(t, db)
	require.Positive(t, one)

	status, err := store.QuotaStatus(t.Context(), testKey)
	require.NoError(t, err)
	require.False(t, status.Configured, "without an explicit policy usage is not judged")

	quota, err := store.WithQuota(QuotaConfig{GlobalMaxBytes: one * 10, WarningRatio: 0.5, HardRatio: 0.9})
	require.NoError(t, err)
	status, err = quota.QuotaStatus(t.Context(), testKey)
	require.NoError(t, err)
	require.True(t, status.Configured)
	require.Equal(t, QuotaOK, status.State)
	require.Equal(t, one, status.UsedBytes, "usage is the measured payload bytes of accepted pending samples")
	require.EqualValues(t, one*10, status.MaxBytes)

	require.NoError(t, appendN(t, quota, testKey, 1, 6)) // 6 samples ≈ 60% of the limit
	status, err = quota.QuotaStatus(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, QuotaWarning, status.State)
	require.Equal(t, QuotaScopeGlobal, status.Scope)
	require.False(t, status.IntakeRefused, "a warning still accepts data")
}

func TestQuotaRejectsNewAckHardLimitRefusesWithoutTouchingAcceptedData(t *testing.T) {
	store, db, _ := newTestStore(t)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, uniformSample(0)))
	one := journalBytes(t, db)
	quota, err := store.WithQuota(QuotaConfig{GlobalMaxBytes: one * 10, WarningRatio: 0.5, HardRatio: 0.9})
	require.NoError(t, err)

	accepted, refusal := appendUntilRefused(t, quota, testKey)
	var quotaErr *QuotaError
	require.ErrorAs(t, refusal, &quotaErr)
	require.ErrorIs(t, refusal, ErrQuotaExceeded)
	require.Equal(t, ReasonQuotaHardLimit, quotaErr.Reason)
	require.Equal(t, QuotaScopeGlobal, quotaErr.Scope)
	require.Equal(t, one*10, quotaErr.MaxBytes)
	require.GreaterOrEqual(t, float64(quotaErr.UsedBytes), 0.9*float64(one*10), "intake stops once the hard threshold is reached")
	require.LessOrEqual(t, quotaErr.UsedBytes, one*10, "and accepted data stays within the maximum here")
	require.Equal(t, 8, accepted, "nine samples reach 90% of ten units")

	before := countRows(t, db, "wg_delivery_samples")
	require.Error(t, appendN(t, quota, testKey, 300, 302))
	require.Equal(t, before, countRows(t, db, "wg_delivery_samples"), "a refused sample is never partially stored")
	status, err := quota.QuotaStatus(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, QuotaHard, status.State, "the status agrees with the refusal")
	require.True(t, status.IntakeRefused)
	require.Contains(t, status.LossRiskNotice, "not recorded", "the operator is told new readings are at risk")

	restored, err := quota.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Len(t, restored.Samples, before, "every previously accepted sample is still there")
}

func TestQuotaRejectsNewAckDuplicateOfAcceptedSampleIsNotRefused(t *testing.T) {
	store, db, _ := newTestStore(t)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, uniformSample(0)))
	one := journalBytes(t, db)
	quota, err := store.WithQuota(QuotaConfig{GlobalMaxBytes: one * 2, WarningRatio: 0.3, HardRatio: 0.5})
	require.NoError(t, err)

	require.NoError(t, quota.AppendSample(t.Context(), testKey, t0, uniformSample(0)),
		"a resend of data that is already durable needs no new space")
	require.ErrorIs(t, quota.AppendSample(t.Context(), testKey, t0, uniformSample(1)), ErrQuotaExceeded)
}

func TestQuotaRejectsNewAckGroupScopeLimitsOnlyThatGroup(t *testing.T) {
	store, db, _ := newTestStore(t)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, uniformSample(0)))
	one := journalBytes(t, db)
	quota, err := store.WithQuota(QuotaConfig{GlobalMaxBytes: one * 100, GroupMaxBytes: one * 3, HardRatio: 1})
	require.NoError(t, err)
	require.NoError(t, appendN(t, quota, testKey, 1, 3))

	err = appendN(t, quota, testKey, 3, 5)
	var quotaErr *QuotaError
	require.ErrorAs(t, err, &quotaErr)
	require.Equal(t, QuotaScopeGroup, quotaErr.Scope)
	require.Equal(t, testKey.GroupID, quotaErr.GroupID)

	require.NoError(t, quota.AppendSample(t.Context(), otherKey, t0, uniformSample(9)), "another group is not affected")
	status, err := quota.QuotaStatus(t.Context(), testKey)
	require.NoError(t, err)
	require.True(t, status.IntakeRefused)
	require.Equal(t, QuotaScopeGroup, status.Scope)
	other, err := quota.QuotaStatus(t.Context(), otherKey)
	require.NoError(t, err)
	require.False(t, other.IntakeRefused)
}

func TestQuotaRejectsNewAckUsageCountsOnlyUndeliveredData(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedSamples(t, store) // three unconsumed samples
	samples := journalBytes(t, db)
	usage, err := store.Usage(t.Context(), "")
	require.NoError(t, err)
	require.Equal(t, samples, usage)

	closure := Closure{Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{rowBucket("", t0, 7)}}
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closure))
	var payload int64
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT LENGTH(payload) FROM wg_delivery_outbox`).Scan(&payload))
	remainingSamples := journalBytes(t, db)
	usage, err = store.Usage(t.Context(), "")
	require.NoError(t, err)
	require.Equal(t, remainingSamples+payload, usage, "consumed samples stop counting; the pending outbox row starts")

	effect := closure.Buckets[0].Outcome.EffectKey
	claim, err := store.BeginDelivery(t.Context(), effect, "worker-1", time.Minute)
	require.NoError(t, err)
	require.NoError(t, store.CompleteDelivery(t.Context(), claim, "d"))
	usage, err = store.Usage(t.Context(), "")
	require.NoError(t, err)
	require.Equal(t, remainingSamples, usage, "a committed row no longer counts")

	groupUsage, err := store.Usage(t.Context(), testKey.GroupID)
	require.NoError(t, err)
	require.Equal(t, remainingSamples, groupUsage)
	otherUsage, err := store.Usage(t.Context(), otherKey.GroupID)
	require.NoError(t, err)
	require.Zero(t, otherUsage)
}

func TestQuotaRejectsNewAckReclaimNeverDeletesPendingData(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedSamples(t, store)
	pending := rowBucket("", t0, 1)
	done := rowBucket("", at(10), 2)
	require.NoError(t, store.CommitClosure(t.Context(), testKey,
		Closure{Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{pending}}))
	require.NoError(t, store.CommitClosure(t.Context(), testKey,
		Closure{Destination: testDestination, NextClose: at(20), Buckets: []ClosedBucket{done}}))
	claim, err := store.BeginDelivery(t.Context(), done.Outcome.EffectKey, "worker-1", time.Minute)
	require.NoError(t, err)
	require.NoError(t, store.CompleteDelivery(t.Context(), claim, "d"))
	_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET committed_at = ? WHERE effect_key = ?`,
		timestamp(time.Now().UTC().Add(-48*time.Hour)), done.Outcome.EffectKey)
	require.NoError(t, err)
	unconsumed := openCount(t, db, "wg_delivery_samples", "consumed = 0")

	removed, err := store.Reclaim(t.Context(), 24*time.Hour)
	require.NoError(t, err)
	require.Positive(t, removed)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_outbox"), "only the old committed row went")
	require.Equal(t, 1, openCount(t, db, "wg_delivery_outbox", "state = 'pending'"), "the pending row is untouched")
	require.Equal(t, unconsumed, openCount(t, db, "wg_delivery_samples", "consumed = 0"), "unconsumed samples are untouched")
	require.Zero(t, openCount(t, db, "wg_delivery_samples", "consumed = 1"), "finished samples are reclaimed")

	// Even with a hopeless quota, Reclaim stays within finished data.
	quota, err := store.WithQuota(QuotaConfig{GlobalMaxBytes: 1, HardRatio: 1})
	require.NoError(t, err)
	removed, err = quota.Reclaim(t.Context(), 0)
	require.NoError(t, err)
	require.Equal(t, 1, openCount(t, db, "wg_delivery_outbox", "state = 'pending'"))
	require.GreaterOrEqual(t, removed, int64(0))
}

func TestQuotaRejectsNewAckRejectsInvalidPolicy(t *testing.T) {
	store, _, _ := newTestStore(t)
	for name, config := range map[string]QuotaConfig{
		"negative global":   {GlobalMaxBytes: -1},
		"negative group":    {GroupMaxBytes: -1},
		"warning above 1":   {GlobalMaxBytes: 10, WarningRatio: 1.5},
		"hard above 1":      {GlobalMaxBytes: 10, HardRatio: 1.5},
		"warning over hard": {GlobalMaxBytes: 10, WarningRatio: 0.9, HardRatio: 0.5},
		"nothing limited":   {},
	} {
		_, err := store.WithQuota(config)
		require.Error(t, err, name)
	}
}

func openSingleConnStore(t *testing.T) (*Store, *sql.DB) {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "gateway.db")+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	return NewStore(db), db
}

func TestQuotaRejectsNewAckDiskFullIsRefusedAndKeepsAcceptedData(t *testing.T) {
	store, db := openSingleConnStore(t)
	require.NoError(t, appendN(t, store, testKey, 0, 3))
	closure := Closure{Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{rowBucket("", t0, 1)}}
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closure))
	before := countRows(t, db, "wg_delivery_samples")

	var pages int
	require.NoError(t, db.QueryRowContext(t.Context(), `PRAGMA page_count`).Scan(&pages))
	_, err := db.ExecContext(t.Context(), `PRAGMA max_page_count = `+itoa(pages))
	require.NoError(t, err)

	var refused error
	for i := 100; i < 400 && refused == nil; i++ { // keep going until the file genuinely cannot grow
		refused = store.AppendSample(t.Context(), testKey, at(10), sampleAt("full-"+itoa(i), "temp", 13, int64(i)))
	}
	var quotaErr *QuotaError
	require.ErrorAs(t, refused, &quotaErr, "a full disk must surface as a capacity refusal, not a generic failure")
	require.Equal(t, ReasonDiskFull, quotaErr.Reason)
	require.ErrorIs(t, refused, ErrQuotaExceeded)
	require.GreaterOrEqual(t, countRows(t, db, "wg_delivery_samples"), before, "nothing accepted earlier was dropped to make room")
	require.Equal(t, 1, countRows(t, db, "wg_delivery_outbox"))

	// Finished data can be reclaimed; once space exists intake resumes.
	removed, err := store.Reclaim(t.Context(), 0)
	require.NoError(t, err)
	_ = removed
	_, err = db.ExecContext(t.Context(), `PRAGMA max_page_count = `+itoa(pages+64))
	require.NoError(t, err)
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(10), sampleAt("after-space", "temp", 14, 1)))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}

func TestQuotaRejectsNewAckBoundaryRefusalCarriesScopeAndLossRisk(t *testing.T) {
	// The boundary-level mapping lives in the runtime package tests; here we
	// pin the error contract it depends on.
	err := error(&QuotaError{Reason: ReasonQuotaHardLimit, Scope: QuotaScopeGroup, GroupID: "group-G", UsedBytes: 9, MaxBytes: 10})
	var quotaErr *QuotaError
	require.True(t, errors.As(err, &quotaErr))
	require.ErrorIs(t, err, ErrQuotaExceeded)
	require.NotContains(t, err.Error(), "group-G", "the message stays value-free; scope is in the fields")
	_ = snapshot.OutcomeRow
}

func TestQuotaRejectsNewAckReclaimPrunesOldBucketStatusButNotRecentOnes(t *testing.T) {
	store, db, _ := newTestStore(t)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, Closure{
		Destination: testDestination, NextClose: at(20), Buckets: []ClosedBucket{silentBucket("", t0), silentBucket("", at(10))},
	}))
	_, err := db.ExecContext(t.Context(), `UPDATE wg_delivery_buckets SET created_at = ? WHERE bucket_start = ?`,
		timestamp(time.Now().UTC().Add(-48*time.Hour)), timestamp(t0))
	require.NoError(t, err)

	removed, err := store.Reclaim(t.Context(), 24*time.Hour)
	require.NoError(t, err)
	require.Equal(t, int64(1), removed)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_buckets"), "status history older than the retention is pruned, recent history stays")
}

func TestQuotaRejectsNewAckIntakeChecksDoNotRescanTheBacklogPerSample(t *testing.T) {
	store, db, _ := newTestStore(t)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, uniformSample(0)))
	one := journalBytes(t, db)
	quota, err := store.WithQuota(QuotaConfig{GlobalMaxBytes: one * 1000, WarningRatio: 0.5, HardRatio: 0.9})
	require.NoError(t, err)

	// A counter makes every SUM over the backlog visible: a trigger-less way is
	// to count statements that read the usage. Here the usage check goes through
	// cachedUsage, so repeated accepts within the TTL measure at most once.
	before := quota.usage.measurements.Load()
	require.NoError(t, appendN(t, quota, testKey, 1, 40))
	measured := quota.usage.measurements.Load() - before
	require.LessOrEqual(t, measured, int64(2), "40 accepted samples must not trigger 40 full scans")

	// The cache is conservative: freed space is noticed once the TTL passes.
	quota.usage.ttl = time.Millisecond
	time.Sleep(5 * time.Millisecond)
	status, err := quota.QuotaStatus(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, QuotaOK, status.State)
}

func TestQuotaRejectsNewAckCachedUsageStillRefusesAtTheHardLimit(t *testing.T) {
	store, db, _ := newTestStore(t)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, uniformSample(0)))
	one := journalBytes(t, db)
	quota, err := store.WithQuota(QuotaConfig{GlobalMaxBytes: one * 10, WarningRatio: 0.5, HardRatio: 0.9})
	require.NoError(t, err)
	_, refusal := appendUntilRefused(t, quota, testKey)
	require.ErrorIs(t, refusal, ErrQuotaExceeded, "incremental accounting between measurements still reaches the hard limit")

	// A refreshed measurement agrees (the cache never under-counts accepted data).
	quota.usage.invalidate()
	_, again := appendUntilRefused(t, quota, testKey)
	require.ErrorIs(t, again, ErrQuotaExceeded)
}
