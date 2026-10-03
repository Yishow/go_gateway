package runtime

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func openDurableDB(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	return db
}

func (f *boundaryFixture) withLedger(db *sql.DB) {
	key := groupdelivery.GroupKey{WorkspaceID: "workspace-A", GroupID: "group-G", GroupRevision: "rev-1"}
	f.config.Ledger = groupdelivery.NewStore(db).Ledger(key)
}

func countJournal(t *testing.T, db *sql.DB) int {
	t.Helper()
	var n int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_samples`).Scan(&n))
	return n
}

func TestDurableSampleAckCrashBoundaryJournalsBeforeAckAndRecovers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gateway.db")
	db := openDurableDB(t, path)
	f := newBoundaryFixture(t)
	f.withLedger(db)
	boundary := f.build(t)
	f.feedAll(t, boundary)
	require.Equal(t, 5, countJournal(t, db), "every ACKed sample is already committed")

	// Crash: the process dies after the ACKs; no Tick, no memory survives.
	require.NoError(t, db.Close())
	reopened := openDurableDB(t, path)
	again := newBoundaryFixture(t)
	again.withLedger(reopened)
	recovered := again.build(t)

	require.NoError(t, recovered.Tick(t.Context(), boundaryAt(10)))
	require.Equal(t, 1, countTable(t, reopened, "wg_delivery_outbox"), "the ACKed samples alone are enough to rebuild the row")
	cells := outboxCells(t, reopened)
	require.Equal(t, 21.5, cells["temperature"])
	require.Equal(t, int64(9007199254740993), cells["counter"])
}

func TestDurableSampleAckCrashCommitFailureIsNotAckedOrUsed(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := newBoundaryFixture(t)
	f.withLedger(db)
	boundary := f.build(t)

	_, err := db.ExecContext(t.Context(), `CREATE TRIGGER fail_samples BEFORE INSERT ON wg_delivery_samples
		BEGIN SELECT RAISE(ABORT, 'injected commit failure'); END`)
	require.NoError(t, err)

	err = boundary.AcceptSample(t.Context(), f.envelope(0, "s-temp", 4, 21.5))
	var refused *GroupSampleError
	require.ErrorAs(t, err, &refused)
	require.Equal(t, "journal-failed", refused.Reason)
	require.NotContains(t, err.Error(), "injected", "the safe error never leaks driver text")
	require.ErrorContains(t, errors.Unwrap(err), "injected commit failure", "the cause stays available to callers")
	require.Zero(t, countJournal(t, db))

	_, err = db.ExecContext(t.Context(), `DROP TRIGGER fail_samples`)
	require.NoError(t, err)
	for i, value := range []any{int64(101), true, "batch-001", uint64(1)} {
		require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(i+1, "s"+f.members[i+1].TagID, 4, value)))
	}
	require.Equal(t, 4, countJournal(t, db))

	// The sender resends the same sample. It must be journaled now: had the
	// refused attempt leaked into memory it would be treated as a duplicate and
	// skipped, leaving the ACKed sample uncommitted.
	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "s-temp", 4, 21.5)))
	require.Equal(t, 5, countJournal(t, db))

	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
	require.Equal(t, 1, countTable(t, db, "wg_delivery_outbox"))
}

func TestDurableSampleAckCrashDuplicateAfterRestartIsNoOpAndConflictRefused(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gateway.db")
	db := openDurableDB(t, path)
	f := newBoundaryFixture(t)
	f.withLedger(db)
	boundary := f.build(t)
	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "s1", 4, 21.5)))
	require.NoError(t, db.Close())

	reopened := openDurableDB(t, path)
	again := newBoundaryFixture(t)
	again.withLedger(reopened)
	recovered := again.build(t)
	require.NoError(t, recovered.AcceptSample(t.Context(), again.envelope(0, "s1", 4, 21.5)), "resend after restart is a no-op")
	require.Equal(t, 1, countJournal(t, reopened))

	var refused *GroupSampleError
	require.ErrorAs(t, recovered.AcceptSample(t.Context(), again.envelope(0, "s1", 4, 99.5)), &refused)
	require.Equal(t, snapshot.ReasonIdentityConflict, refused.Reason)
	require.Equal(t, 1, countJournal(t, reopened))
}

func TestDurableSampleAckCrashRefusedSamplesAreNeverJournaled(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := newBoundaryFixture(t)
	f.withLedger(db)
	boundary := f.build(t)

	stale := f.envelope(0, "stale", 4, 1.0)
	stale.MappingRevision = "map-A-old"
	require.Error(t, boundary.AcceptSample(t.Context(), stale))
	require.Error(t, boundary.AcceptSample(t.Context(), f.envelope(0, "far", 3600, 1.0)))
	foreign := f.envelope(0, "foreign", 4, 1.0)
	foreign.WorkspaceID = "workspace-B"
	require.Error(t, boundary.AcceptSample(t.Context(), foreign))
	other := f.envelope(0, "other", 4, 1.0)
	other.TagID = "tag-unrelated"
	require.NoError(t, boundary.AcceptSample(t.Context(), other))
	require.Zero(t, countJournal(t, db), "only accepted samples of this group reach the journal")
}

func outboxCells(t *testing.T, db *sql.DB) map[string]any {
	t.Helper()
	var payload string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_outbox`).Scan(&payload))
	row, err := groupdelivery.DecodeRowPayload([]byte(payload))
	require.NoError(t, err)
	cells := map[string]any{}
	for _, cell := range row.Cells {
		cells[cell.Column] = cell.Value
	}
	return cells
}

func countTable(t *testing.T, db *sql.DB, table string) int {
	t.Helper()
	var n int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+table).Scan(&n))
	return n
}

func TestAtomicRowOutboxCheckpointBoundaryClosesIntoOutboxAndSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gateway.db")
	db := openDurableDB(t, path)
	f := newBoundaryFixture(t)
	f.withLedger(db)
	boundary := f.build(t)
	f.feedAll(t, boundary)

	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
	require.Empty(t, f.sink.rows, "in durable mode the outbox, not the sink, is the handoff")
	require.Equal(t, 1, countTable(t, db, "wg_delivery_outbox"))
	require.Equal(t, 1, countTable(t, db, "wg_delivery_buckets"))
	var consumed int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed = 1`).Scan(&consumed))
	require.Equal(t, 5, consumed)

	var scope, connector string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT destination_scope, connector_revision FROM wg_delivery_outbox`).Scan(&scope, &connector))
	wantScope, err := snapshot.DestinationScope("connector-1", "crev-1", "db", "", "readings")
	require.NoError(t, err)
	require.Equal(t, wantScope, scope)
	require.Equal(t, "crev-1", connector)

	// Restart: nothing is re-emitted and a late sample of the closed bucket is refused.
	require.NoError(t, db.Close())
	reopened := openDurableDB(t, path)
	again := newBoundaryFixture(t)
	again.withLedger(reopened)
	recovered := again.build(t)
	require.NoError(t, recovered.Tick(t.Context(), boundaryAt(10)))
	require.Equal(t, 1, countTable(t, reopened, "wg_delivery_outbox"), "the closed bucket is not closed twice")
	again.clock = boundaryAt(12)
	var refused *GroupSampleError
	require.ErrorAs(t, recovered.AcceptSample(t.Context(), again.envelope(0, "late", 9, 99.0)), &refused)
	require.Equal(t, snapshot.ReasonLate, refused.Reason, "a consumed bucket cannot take samples after restart")
}

func TestAtomicRowOutboxCheckpointBoundaryPreservesEntityScopedRows(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := newBoundaryFixture(t)
	group := copyGroup(f.config.Group)
	group.Members = group.Members[:2]
	group.Members[0].EntityKey = "line-a"
	group.Members[0].TargetColumn = "reading"
	group.Members[1].EntityKey = "line-b"
	group.Members[1].TargetColumn = "reading"
	group.RowPolicy.EntityKeyColumn = "entity"
	f.config.Group = group
	f.members = group.Members
	f.config.TagTypes["tag-B"] = f.config.TagTypes["tag-A"]
	f.config.Columns = []dbtarget.ColumnInfo{
		{Name: "reading", DataType: "REAL"},
		{Name: "entity", DataType: "TEXT"},
	}
	f.withLedger(db)
	boundary := f.build(t)

	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "line-a-sample", 4, 21.5)))
	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(1, "line-b-sample", 4, 42.5)))
	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
	require.Empty(t, f.sink.rows, "durable rows hand off through the outbox")
	require.Equal(t, 2, countTable(t, db, "wg_delivery_outbox"))
	var skipped int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_buckets WHERE kind = 'skipped'`).Scan(&skipped))
	require.Zero(t, skipped, "complete entity rows are not encoded as member-missing")

	rows, err := db.QueryContext(t.Context(), `SELECT entity_key, payload FROM wg_delivery_outbox ORDER BY entity_key`)
	require.NoError(t, err)
	defer rows.Close()
	var entities []string
	for rows.Next() {
		var entity, payload string
		require.NoError(t, rows.Scan(&entity, &payload))
		encoded, decodeErr := groupdelivery.DecodeRowPayload([]byte(payload))
		require.NoError(t, decodeErr)
		require.Equal(t, entity, encoded.EntityKey)
		cells := map[string]any{}
		for _, cell := range encoded.Cells {
			cells[cell.Column] = cell.Value
		}
		require.Equal(t, entity, cells["entity"])
		require.Contains(t, cells, "reading")
		entities = append(entities, entity)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, []string{"line-a", "line-b"}, entities)
}

func TestAtomicRowOutboxCheckpointBoundaryCommitFailureKeepsBucketOpenForReplay(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := newBoundaryFixture(t)
	f.withLedger(db)
	boundary := f.build(t)
	f.feedAll(t, boundary)

	_, err := db.ExecContext(t.Context(), `CREATE TRIGGER fail_outbox BEFORE INSERT ON wg_delivery_outbox
		BEGIN SELECT RAISE(ABORT, 'injected commit failure'); END`)
	require.NoError(t, err)
	require.ErrorContains(t, boundary.Tick(t.Context(), boundaryAt(10)), "injected commit failure")
	require.Zero(t, countTable(t, db, "wg_delivery_outbox"))
	require.Zero(t, countTable(t, db, "wg_delivery_checkpoints"), "the checkpoint moves only with the outbox")

	_, err = db.ExecContext(t.Context(), `DROP TRIGGER fail_outbox`)
	require.NoError(t, err)
	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)), "the same closure is retried from unchanged memory")
	require.Equal(t, 1, countTable(t, db, "wg_delivery_outbox"))
	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
	require.Equal(t, 1, countTable(t, db, "wg_delivery_outbox"), "and never closed twice")
}

func TestAtomicRowOutboxCheckpointBoundaryPersistsSilentAndSkippedBuckets(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := newBoundaryFixture(t)
	f.withLedger(db)
	boundary := f.build(t)
	f.clock = boundaryAt(15)
	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "only-temp", 14, 21.5)))

	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(30)))
	require.Zero(t, countTable(t, db, "wg_delivery_outbox"))
	var noData, skipped int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_buckets WHERE kind = 'no_data'`).Scan(&noData))
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_buckets WHERE kind = 'skipped'`).Scan(&skipped))
	require.Equal(t, 2, noData, "buckets [0,10) and [20,30) were silent")
	require.Equal(t, 1, skipped, "bucket [10,20) lacked required members")
	require.Len(t, f.sink.outcomes, 3, "operators are also told after the commit")
}

func TestAtomicRowOutboxCheckpointBoundaryBlockedEncodingIsPersistedAsSkipped(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := newBoundaryFixture(t)
	f.withLedger(db)
	boundary := f.build(t)
	f.feedAll(t, boundary)
	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(4, "huge", 6, uint64(1<<63))))

	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
	require.Zero(t, countTable(t, db, "wg_delivery_outbox"))
	var reason string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT reason FROM wg_delivery_buckets WHERE kind = 'skipped'`).Scan(&reason))
	require.Equal(t, "encode-blocked:sql-value-blocked", reason)
}

func (f *boundaryFixture) withQuotaLedger(db *sql.DB, config groupdelivery.QuotaConfig) {
	key := groupdelivery.GroupKey{WorkspaceID: "workspace-A", GroupID: "group-G", GroupRevision: "rev-1"}
	quota, err := groupdelivery.NewStore(db).WithQuota(config)
	if err != nil {
		panic(err)
	}
	f.config.Ledger = quota.Ledger(key)
}

func TestQuotaRejectsNewAckBoundaryRefusesWithScopeAndKeepsAcceptedData(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := newBoundaryFixture(t)
	// Probe the size of one journaled sample, then size the quota to a few of them.
	f.withLedger(db)
	probe := f.build(t)
	require.NoError(t, probe.AcceptSample(t.Context(), f.envelope(0, "probe", 4, 21.5)))
	var one int64
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT LENGTH(payload) FROM wg_delivery_samples`).Scan(&one))

	f.withQuotaLedger(db, groupdelivery.QuotaConfig{GlobalMaxBytes: one * 3, WarningRatio: 0.5, HardRatio: 0.9})
	boundary := f.build(t)
	var refusal error
	accepted := 0
	for i, value := range []any{int64(101), true, "batch-001", uint64(9)} {
		if err := boundary.AcceptSample(t.Context(), f.envelope(i+1, "q"+f.members[i+1].TagID, 4, value)); err != nil {
			refusal = err
			break
		}
		accepted++
	}
	require.Error(t, refusal, "the store runs out of budget before every member is accepted")
	var refused *GroupSampleError
	require.ErrorAs(t, refusal, &refused)
	require.Equal(t, groupdelivery.ReasonQuotaHardLimit, refused.Reason)
	require.Equal(t, snapshot.OfferRejected, refused.Outcome)
	var quotaErr *groupdelivery.QuotaError
	require.ErrorAs(t, refusal, &quotaErr)
	require.Equal(t, groupdelivery.QuotaScopeGlobal, quotaErr.Scope)
	require.Equal(t, one*3, quotaErr.MaxBytes)
	require.Equal(t, 1+accepted, countJournal(t, db), "a refused sample is not journaled and nothing accepted earlier was dropped")
}

func TestQuotaRejectsNewAckClosureStillCommitsAtTheHardLimit(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := newBoundaryFixture(t)
	f.withLedger(db)
	boundary := f.build(t)
	f.feedAll(t, boundary)
	var total int64
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT SUM(LENGTH(payload)) FROM wg_delivery_samples`).Scan(&total))

	// Re-open the same journal under a quota it already exceeds.
	f.withQuotaLedger(db, groupdelivery.QuotaConfig{GlobalMaxBytes: total, WarningRatio: 0.5, HardRatio: 0.5})
	capped := f.build(t)
	var refused *GroupSampleError
	require.ErrorAs(t, capped.AcceptSample(t.Context(), f.envelope(0, "extra", 5, 22.5)), &refused)
	require.Equal(t, groupdelivery.ReasonQuotaHardLimit, refused.Reason)

	require.NoError(t, capped.Tick(t.Context(), boundaryAt(10)), "turning accepted samples into a deliverable row must not be blocked by the intake limit")
	require.Equal(t, 1, countTable(t, db, "wg_delivery_outbox"))
	require.NoError(t, groupdeliveryReclaim(t, db))
}

func groupdeliveryReclaim(t *testing.T, db *sql.DB) error {
	t.Helper()
	removed, err := groupdelivery.NewStore(db).Reclaim(t.Context(), 0)
	require.Positive(t, removed, "the consumed samples can now be reclaimed")
	require.Equal(t, 1, countTable(t, db, "wg_delivery_outbox"), "the pending row is kept")
	return err
}
