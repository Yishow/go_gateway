package groupdelivery

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	t0       = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	testKey  = GroupKey{WorkspaceID: "workspace-A", GroupID: "group-G", GroupRevision: "rev-1"}
	otherKey = GroupKey{WorkspaceID: "workspace-A", GroupID: "group-H", GroupRevision: "rev-1"}
)

func at(seconds int) time.Time { return t0.Add(time.Duration(seconds) * time.Second) }

// openStoreDB opens a migrated on-disk SQLite file so tests can reopen it to
// model a process restart.
func openStoreDB(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	return db
}

func newTestStore(t *testing.T) (*Store, *sql.DB, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "gateway.db")
	db := openStoreDB(t, path)
	return NewStore(db), db, path
}

func sampleAt(id, member string, seconds int, value int64) snapshot.Sample {
	return snapshot.Sample{
		SampleID: id, MemberKey: member, ObservedAt: at(seconds), Quality: schema.QualityGood,
		Value: measurement.NewInt64(value), SourceRevision: "src-1", MappingRevision: "map-1",
	}
}

func failOn(t *testing.T, db *sql.DB, table string) {
	t.Helper()
	_, err := db.ExecContext(t.Context(), `CREATE TRIGGER fail_`+table+` BEFORE INSERT ON `+table+
		` BEGIN SELECT RAISE(ABORT, 'injected commit failure'); END`)
	require.NoError(t, err)
}

func countRows(t *testing.T, db *sql.DB, table string) int {
	t.Helper()
	var n int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+table).Scan(&n))
	return n
}

func TestDurableSampleAckCrashAppendThenRestore(t *testing.T) {
	store, _, _ := newTestStore(t)
	first := sampleAt("s1", "temp", 3, 30)
	second := sampleAt("s2", "pressure", 8, 80)
	require.NoError(t, store.AppendSample(t.Context(), testKey, snapshot.BucketStart(first.ObservedAt, 10*time.Second), first))
	require.NoError(t, store.AppendSample(t.Context(), testKey, snapshot.BucketStart(second.ObservedAt, 10*time.Second), second))
	require.NoError(t, store.AppendSample(t.Context(), otherKey, t0, sampleAt("other", "temp", 4, 1)))

	restored, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.True(t, restored.NextClose.IsZero(), "no bucket has closed yet")
	require.Len(t, restored.Samples, 2, "another group's samples are never restored")
	require.Equal(t, "s1", restored.Samples[0].SampleID, "acceptance order is preserved")
	require.Equal(t, "s2", restored.Samples[1].SampleID)
	require.True(t, measurement.NewInt64(80).Equal(restored.Samples[1].Value), "exact value survives the journal")
	require.Equal(t, at(8), restored.Samples[1].ObservedAt)
	require.Equal(t, "map-1", restored.Samples[1].MappingRevision)
	require.Equal(t, "pressure", restored.Samples[1].MemberKey)
}

func TestDurableSampleAckCrashDuplicateAndConflict(t *testing.T) {
	store, db, _ := newTestStore(t)
	sample := sampleAt("s1", "temp", 3, 30)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample), "an identical resend is a no-op")
	require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))

	changed := sample
	changed.Value = measurement.NewInt64(31)
	require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, changed), ErrSampleConflict)
	restored, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.True(t, measurement.NewInt64(30).Equal(restored.Samples[0].Value), "a conflicting resend never replaces the journaled sample")

	// Same sample ID in a different group revision is a different identity.
	next := GroupKey{WorkspaceID: "workspace-A", GroupID: "group-G", GroupRevision: "rev-2"}
	require.NoError(t, store.AppendSample(t.Context(), next, t0, changed))
}

func TestDurableSampleAckCrashCommitFailureIsNotAcknowledged(t *testing.T) {
	store, db, _ := newTestStore(t)
	failOn(t, db, "wg_delivery_samples")

	err := store.AppendSample(t.Context(), testKey, t0, sampleAt("s1", "temp", 3, 30))
	require.Error(t, err, "a failed commit must not look like an ACK")
	require.ErrorContains(t, err, "injected commit failure")
	require.Zero(t, countRows(t, db, "wg_delivery_samples"))
}

func TestDurableSampleAckCrashSurvivesReopeningTheDatabase(t *testing.T) {
	store, db, path := newTestStore(t)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sampleAt("s1", "temp", 3, 30)))
	require.NoError(t, db.Close())

	reopened := NewStore(openStoreDB(t, path))
	restored, err := reopened.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
	require.Equal(t, "s1", restored.Samples[0].SampleID)
}

func TestDurableSampleAckCrashRejectsIncompleteKeysAndSamples(t *testing.T) {
	store, _, _ := newTestStore(t)
	require.ErrorIs(t, store.AppendSample(t.Context(), GroupKey{}, t0, sampleAt("s1", "temp", 3, 1)), ErrInvalidGroupKey)
	require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, snapshot.Sample{}), ErrInvalidSample)
	_, err := store.Restore(t.Context(), GroupKey{})
	require.ErrorIs(t, err, ErrInvalidGroupKey)
}

func TestDurableSampleAckCrashConcurrentWritersNeverSeeBusy(t *testing.T) {
	store, db, _ := newTestStore(t)
	db.SetMaxOpenConns(8)
	const writers, perWriter = 8, 25
	errs := make(chan error, writers*perWriter+writers)
	var wg sync.WaitGroup
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			key := GroupKey{WorkspaceID: "workspace-A", GroupID: fmt.Sprintf("group-%d", w%4), GroupRevision: "rev-1"}
			for i := 0; i < perWriter; i++ {
				id := fmt.Sprintf("sample-%d-%d", w, i)
				errs <- store.AppendSample(t.Context(), key, t0, sampleAt(id, "temp", 3, int64(i)))
			}
			// Closure commits interleave with journal appends from other groups.
			bucket := rowBucket("", at(10*w), uint64(w))
			errs <- store.CommitClosure(t.Context(), key, Closure{Destination: testDestination, NextClose: at(10*w + 10), Buckets: []ClosedBucket{bucket}})
		}(w)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err, "concurrent local writers must queue on the write lock, not fail with SQLITE_BUSY")
	}
	require.Equal(t, writers*perWriter, countRows(t, db, "wg_delivery_samples"))
	require.Equal(t, writers, countRows(t, db, "wg_delivery_outbox"))
}

func TestRevisionBoundBacklogMigrationRerunKeepsAcceptedData(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedSamples(t, store)
	closure := Closure{Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{rowBucket("", t0, 1)}}
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closure))
	samples, outbox, buckets := countRows(t, db, "wg_delivery_samples"), countRows(t, db, "wg_delivery_outbox"), countRows(t, db, "wg_delivery_buckets")

	// A binary that does not know these tables, or a repeated upgrade, runs the
	// migrator again: it must only add what is missing and never alter accepted data.
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	require.Equal(t, samples, countRows(t, db, "wg_delivery_samples"))
	require.Equal(t, outbox, countRows(t, db, "wg_delivery_outbox"))
	require.Equal(t, buckets, countRows(t, db, "wg_delivery_buckets"))
	restored, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, at(10), restored.NextClose)
	require.Len(t, restored.Samples, 1)
}
