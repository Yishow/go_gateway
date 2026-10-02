package dbtarget

import (
	"database/sql"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

var insertBucket = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func insertRow(counter int64, recordID string) EncodedRow {
	return EncodedRow{
		RecordID: recordID, EffectKey: "effect-" + recordID, BucketStart: insertBucket,
		Cells: []EncodedCell{
			{Column: "record_id", Value: recordID},
			{Column: "counter", Value: counter},
			{Column: "label", Value: "batch-001"},
			{Column: "optional", Value: nil},
		},
	}
}

func insertRequest(strategy, recordID string, counter int64) GroupInsertRequest {
	return GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindSQLite, TableName: "readings", Row: insertRow(counter, recordID),
		PayloadDigest: "digest-" + recordID + "-" + string(rune('a'+counter)), Strategy: strategy,
		RecordKeyColumn: "record_id", CommittedAt: "2026-01-01T00:00:10Z",
	}
}

func openTargetSQLite(t *testing.T, createTable string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "target.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	_, err = db.ExecContext(t.Context(), createTable)
	require.NoError(t, err)
	return db
}

const targetTable = `CREATE TABLE readings (record_id TEXT, counter INTEGER NOT NULL, label TEXT, optional TEXT)`
const uniqueTargetTable = `CREATE TABLE readings (record_id TEXT PRIMARY KEY, counter INTEGER NOT NULL, label TEXT, optional TEXT)`

func targetCount(t *testing.T, db *sql.DB, table string) int {
	t.Helper()
	var n int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+table).Scan(&n))
	return n
}

func TestSQLiteTargetReceiptIdentityRowAndReceiptShareOneTransaction(t *testing.T) {
	db := openTargetSQLite(t, targetTable)
	require.NoError(t, CreateEffectReceiptTable(t.Context(), db, schema.DatabaseConnectorKindSQLite, ""))
	request := insertRequest(GroupEffectReceipt, "r1", 9007199254740993)

	result, err := InsertGroupRow(t.Context(), db, request)
	require.NoError(t, err)
	require.False(t, result.AlreadyCommitted)
	require.Equal(t, 1, targetCount(t, db, "readings"))
	require.Equal(t, 1, targetCount(t, db, EffectReceiptTable))
	var counter int64
	var optional sql.NullString
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT counter, optional FROM readings`).Scan(&counter, &optional))
	require.Equal(t, int64(9007199254740993), counter)
	require.False(t, optional.Valid, "NULL cells stay NULL")
	var digest, committedAt string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT payload_digest, committed_at FROM `+EffectReceiptTable+` WHERE effect_key = ?`, request.Row.EffectKey).Scan(&digest, &committedAt))
	require.Equal(t, request.PayloadDigest, digest)
	require.Equal(t, "2026-01-01T00:00:10Z", committedAt)
}

func TestSQLiteTargetReceiptIdentityRepeatedSendIsOneEffect(t *testing.T) {
	db := openTargetSQLite(t, targetTable)
	require.NoError(t, CreateEffectReceiptTable(t.Context(), db, schema.DatabaseConnectorKindSQLite, ""))
	request := insertRequest(GroupEffectReceipt, "r1", 5)
	_, err := InsertGroupRow(t.Context(), db, request)
	require.NoError(t, err)

	again, err := InsertGroupRow(t.Context(), db, request)
	require.NoError(t, err)
	require.True(t, again.AlreadyCommitted, "the receipt proves the earlier commit; nothing is inserted twice")
	require.Equal(t, 1, targetCount(t, db, "readings"))

	changed := request
	changed.PayloadDigest = "a-different-digest"
	_, err = InsertGroupRow(t.Context(), db, changed)
	require.ErrorIs(t, err, ErrReceiptDigestMismatch, "an existing receipt is never overwritten")
	require.Equal(t, 1, targetCount(t, db, "readings"))
}

func TestSQLiteTargetReceiptIdentityReceiptFailureRollsBackTheRow(t *testing.T) {
	db := openTargetSQLite(t, targetTable)
	require.NoError(t, CreateEffectReceiptTable(t.Context(), db, schema.DatabaseConnectorKindSQLite, ""))
	_, err := db.ExecContext(t.Context(), `CREATE TRIGGER no_receipt BEFORE INSERT ON `+EffectReceiptTable+
		` BEGIN SELECT RAISE(ABORT, 'injected receipt failure'); END`)
	require.NoError(t, err)

	_, err = InsertGroupRow(t.Context(), db, insertRequest(GroupEffectReceipt, "r1", 5))
	var insertErr *GroupInsertError
	require.ErrorAs(t, err, &insertErr)
	require.Equal(t, InsertPhasePreCommit, insertErr.Phase)
	require.False(t, insertErr.Ambiguous(), "nothing was committed, so a retry is safe")
	require.Zero(t, targetCount(t, db, "readings"), "a row without its receipt never survives")
	require.Zero(t, targetCount(t, db, EffectReceiptTable))
}

func TestSQLiteTargetReceiptIdentityMissingReceiptTableIsPreCommit(t *testing.T) {
	db := openTargetSQLite(t, targetTable)
	_, err := InsertGroupRow(t.Context(), db, insertRequest(GroupEffectReceipt, "r1", 5))
	var insertErr *GroupInsertError
	require.ErrorAs(t, err, &insertErr)
	require.Equal(t, InsertPhasePreCommit, insertErr.Phase)
	require.Zero(t, targetCount(t, db, "readings"))
}

func TestSQLiteTargetReceiptIdentityUniqueKeyStrategy(t *testing.T) {
	db := openTargetSQLite(t, uniqueTargetTable)
	request := insertRequest(GroupEffectUniqueKey, "r1", 5)
	result, err := InsertGroupRow(t.Context(), db, request)
	require.NoError(t, err)
	require.False(t, result.AlreadyCommitted)

	again, err := InsertGroupRow(t.Context(), db, request)
	require.NoError(t, err)
	require.True(t, again.AlreadyCommitted, "the unique record key recognizes the same row")
	require.Equal(t, 1, targetCount(t, db, "readings"))

	different := insertRequest(GroupEffectUniqueKey, "r1", 6)
	_, err = InsertGroupRow(t.Context(), db, different)
	require.ErrorIs(t, err, ErrExistingRowDiffers)
	var counter int64
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT counter FROM readings`).Scan(&counter))
	require.Equal(t, int64(5), counter, "the existing row is never changed")
}

func TestSQLiteTargetReceiptIdentityNoDedupeInsertsOnce(t *testing.T) {
	db := openTargetSQLite(t, targetTable)
	_, err := InsertGroupRow(t.Context(), db, insertRequest(GroupEffectNone, "r1", 5))
	require.NoError(t, err)
	require.Equal(t, 1, targetCount(t, db, "readings"))

	_, err = db.ExecContext(t.Context(), `CREATE TRIGGER no_insert BEFORE INSERT ON readings BEGIN SELECT RAISE(ABORT, 'injected insert failure'); END`)
	require.NoError(t, err)
	_, err = InsertGroupRow(t.Context(), db, insertRequest(GroupEffectNone, "r2", 6))
	var insertErr *GroupInsertError
	require.ErrorAs(t, err, &insertErr)
	require.Equal(t, InsertPhasePreCommit, insertErr.Phase)
	require.Equal(t, 1, targetCount(t, db, "readings"))
}

func TestSQLiteTargetReceiptIdentityRejectsUnsupportedRequests(t *testing.T) {
	db := openTargetSQLite(t, targetTable)
	for name, mutate := range map[string]func(*GroupInsertRequest){
		"unknown strategy":          func(r *GroupInsertRequest) { r.Strategy = "magic" },
		"mysql":                     func(r *GroupInsertRequest) { r.Kind = schema.DatabaseConnectorKindMySQL },
		"no table":                  func(r *GroupInsertRequest) { r.TableName = "" },
		"no cells":                  func(r *GroupInsertRequest) { r.Row.Cells = nil },
		"no effect key":             func(r *GroupInsertRequest) { r.Row.EffectKey = "" },
		"unique key without column": func(r *GroupInsertRequest) { r.Strategy, r.RecordKeyColumn = GroupEffectUniqueKey, "" },
		"receipt without digest":    func(r *GroupInsertRequest) { r.PayloadDigest = "" },
		"unsafe-looking identifier": func(r *GroupInsertRequest) { r.Row.Cells[1].Column = `counter"; DROP TABLE readings; --` },
	} {
		request := insertRequest(GroupEffectReceipt, "r1", 5)
		request.Row.Cells = append([]EncodedCell(nil), request.Row.Cells...)
		mutate(&request)
		_, err := InsertGroupRow(t.Context(), db, request)
		require.Error(t, err, name)
	}
	require.Equal(t, 0, targetCount(t, db, "readings"))
}

func openLostCommitDB(t *testing.T, createTable string) (*sql.DB, *atomic.Bool) {
	t.Helper()
	enabled := &atomic.Bool{}
	db := lostcommit.Open(filepath.Join(t.TempDir(), "target.db"), enabled)
	t.Cleanup(func() { _ = db.Close() })
	_, err := db.ExecContext(t.Context(), createTable)
	require.NoError(t, err)
	return db, enabled
}

func TestSQLiteTargetReceiptIdentityCommitResponseLostWithReceiptIsOneEffect(t *testing.T) {
	db, lose := openLostCommitDB(t, targetTable)
	require.NoError(t, CreateEffectReceiptTable(t.Context(), db, schema.DatabaseConnectorKindSQLite, ""))
	request := insertRequest(GroupEffectReceipt, "r1", 5)

	lose.Store(true)
	_, err := InsertGroupRow(t.Context(), db, request)
	var insertErr *GroupInsertError
	require.ErrorAs(t, err, &insertErr)
	require.Equal(t, InsertPhaseCommit, insertErr.Phase)
	require.True(t, insertErr.Ambiguous(), "the caller cannot know whether the destination committed")
	require.Equal(t, 1, targetCount(t, db, "readings"), "the destination did commit")

	lose.Store(false)
	again, err := InsertGroupRow(t.Context(), db, request)
	require.NoError(t, err)
	require.True(t, again.AlreadyCommitted, "the receipt resolves the ambiguity without a second insert")
	require.Equal(t, 1, targetCount(t, db, "readings"))
}

func TestSQLiteTargetReceiptIdentityCommitResponseLostWithoutDedupeIsAmbiguousNotRetried(t *testing.T) {
	db, lose := openLostCommitDB(t, targetTable)
	request := insertRequest(GroupEffectNone, "r1", 5)

	lose.Store(true)
	_, err := InsertGroupRow(t.Context(), db, request)
	var insertErr *GroupInsertError
	require.ErrorAs(t, err, &insertErr)
	require.True(t, insertErr.Ambiguous())
	require.Equal(t, 1, targetCount(t, db, "readings"))
	// The caller must not call again: the destination cannot recognize the
	// repeat, so a retry would create a second row. The test documents the
	// hazard the delivery state machine guards against.
	lose.Store(false)
	_, err = InsertGroupRow(t.Context(), db, request)
	require.NoError(t, err)
	require.Equal(t, 2, targetCount(t, db, "readings"), "blind retry without dedupe duplicates the effect")
}
