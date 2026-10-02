package dbtarget

import (
	"database/sql"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// Live PostgreSQL fixtures for the group row insert. They need POSTGRES_DSN
// (skipped otherwise) and use one isolated schema per test.

type postgresGroupTarget struct {
	db     *sql.DB
	admin  *sql.DB
	schema string
	lose   *atomic.Bool
}

func newPostgresGroupTarget(t *testing.T, tableDDL string, withReceipts bool) *postgresGroupTarget {
	t.Helper()
	dsn, _ := postgresKeywordDSN(t)
	admin := postgresAdminDB(t, dsn)
	schemaName := postgresScopedSchema(t, admin)
	lose := &atomic.Bool{}
	db := lostcommit.Wrap(stdlib.GetDefaultDriver(), dsn, lose)
	t.Cleanup(func() { _ = db.Close() })
	_, err := admin.ExecContext(t.Context(), strings.ReplaceAll(tableDDL, "{schema}", pgIdentifier(t, schemaName)))
	require.NoError(t, err)
	if withReceipts {
		require.NoError(t, CreateEffectReceiptTable(t.Context(), admin, schema.DatabaseConnectorKindPostgres, schemaName))
	}
	return &postgresGroupTarget{db: db, admin: admin, schema: schemaName, lose: lose}
}

const pgPlainTable = `CREATE TABLE {schema}.readings (
	record_id TEXT, counter BIGINT NOT NULL, label TEXT, ratio DOUBLE PRECISION, running BOOLEAN,
	optional TEXT, bucket_start TIMESTAMPTZ)`

const pgUniqueTable = `CREATE TABLE {schema}.readings (
	record_id TEXT PRIMARY KEY, counter BIGINT NOT NULL, label TEXT, ratio DOUBLE PRECISION, running BOOLEAN,
	optional TEXT, bucket_start TIMESTAMPTZ)`

func (p *postgresGroupTarget) request(strategy string, counter int64) GroupInsertRequest {
	const recordID = "r1"
	return GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindPostgres, SchemaName: p.schema, TableName: "readings",
		Row: EncodedRow{
			RecordID: recordID, EffectKey: "effect-" + recordID, BucketStart: insertBucket,
			Cells: []EncodedCell{
				{Column: "record_id", Value: recordID},
				{Column: "counter", Value: counter},
				{Column: "label", Value: "batch-001"},
				{Column: "ratio", Value: 0.1},
				{Column: "running", Value: true},
				{Column: "optional", Value: nil},
				{Column: "bucket_start", Value: insertBucket},
			},
		},
		PayloadDigest: "digest-" + recordID, Strategy: strategy, RecordKeyColumn: "record_id", CommittedAt: "2026-01-01T00:00:10Z",
	}
}

func (p *postgresGroupTarget) count(t *testing.T, table string) int {
	t.Helper()
	var n int
	require.NoError(t, p.admin.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+pgIdentifier(t, p.schema)+`.`+table).Scan(&n))
	return n
}

func TestPostgresTargetReceiptIdentityRowAndReceiptShareOneTransaction(t *testing.T) {
	p := newPostgresGroupTarget(t, pgPlainTable, true)
	request := p.request(GroupEffectReceipt, 9007199254740993)

	result, err := InsertGroupRow(t.Context(), p.db, request)
	require.NoError(t, err)
	require.False(t, result.AlreadyCommitted)
	require.Equal(t, 1, p.count(t, "readings"))
	require.Equal(t, 1, p.count(t, EffectReceiptTable))

	var counter int64
	var running bool
	var ratio float64
	var optional sql.NullString
	var bucket time.Time
	require.NoError(t, p.admin.QueryRowContext(t.Context(),
		`SELECT counter, running, ratio, optional, bucket_start FROM `+pgIdentifier(t, p.schema)+`.readings`).
		Scan(&counter, &running, &ratio, &optional, &bucket))
	require.Equal(t, int64(9007199254740993), counter, "the exact bigint value reaches real PostgreSQL")
	require.True(t, running)
	require.Equal(t, 0.1, ratio)
	require.False(t, optional.Valid)
	require.True(t, insertBucket.Equal(bucket))
}

func TestPostgresTargetReceiptIdentityRepeatedAndConflictingSends(t *testing.T) {
	p := newPostgresGroupTarget(t, pgPlainTable, true)
	request := p.request(GroupEffectReceipt, 5)
	_, err := InsertGroupRow(t.Context(), p.db, request)
	require.NoError(t, err)

	again, err := InsertGroupRow(t.Context(), p.db, request)
	require.NoError(t, err)
	require.True(t, again.AlreadyCommitted)
	require.Equal(t, 1, p.count(t, "readings"))

	changed := request
	changed.PayloadDigest = "a-different-digest"
	_, err = InsertGroupRow(t.Context(), p.db, changed)
	require.ErrorIs(t, err, ErrReceiptDigestMismatch)
	require.Equal(t, 1, p.count(t, "readings"))
}

func TestPostgresTargetReceiptIdentityReceiptFailureRollsBackTheRow(t *testing.T) {
	p := newPostgresGroupTarget(t, pgPlainTable, true)
	schemaQuoted := pgIdentifier(t, p.schema)
	_, err := p.admin.ExecContext(t.Context(), `
		CREATE FUNCTION `+schemaQuoted+`.refuse_receipt() RETURNS trigger LANGUAGE plpgsql AS $$
		BEGIN RAISE EXCEPTION 'injected receipt failure'; END $$;
		CREATE TRIGGER refuse_receipt BEFORE INSERT ON `+schemaQuoted+`.`+EffectReceiptTable+`
		FOR EACH ROW EXECUTE FUNCTION `+schemaQuoted+`.refuse_receipt()`)
	require.NoError(t, err)

	_, err = InsertGroupRow(t.Context(), p.db, p.request(GroupEffectReceipt, 5))
	var insertErr *GroupInsertError
	require.ErrorAs(t, err, &insertErr)
	require.Equal(t, InsertPhasePreCommit, insertErr.Phase)
	require.Zero(t, p.count(t, "readings"), "a row without its receipt never survives")
	require.Zero(t, p.count(t, EffectReceiptTable))
}

func TestPostgresTargetReceiptIdentityCommitResponseLostIsOneEffect(t *testing.T) {
	p := newPostgresGroupTarget(t, pgPlainTable, true)
	request := p.request(GroupEffectReceipt, 5)

	p.lose.Store(true)
	_, err := InsertGroupRow(t.Context(), p.db, request)
	var insertErr *GroupInsertError
	require.ErrorAs(t, err, &insertErr)
	require.True(t, insertErr.Ambiguous())
	require.Equal(t, 1, p.count(t, "readings"), "PostgreSQL committed although the caller saw an error")

	p.lose.Store(false)
	again, err := InsertGroupRow(t.Context(), p.db, request)
	require.NoError(t, err)
	require.True(t, again.AlreadyCommitted)
	require.Equal(t, 1, p.count(t, "readings"))
}

func TestPostgresTargetReceiptIdentityConcurrentSendersCommitOneEffect(t *testing.T) {
	p := newPostgresGroupTarget(t, pgPlainTable, true)
	request := p.request(GroupEffectReceipt, 5)

	const senders = 8
	var wg sync.WaitGroup
	var committed, already, failed atomic.Int32
	for i := 0; i < senders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := InsertGroupRow(t.Context(), p.db, request)
			switch {
			case err != nil:
				failed.Add(1)
			case result.AlreadyCommitted:
				already.Add(1)
			default:
				committed.Add(1)
			}
		}()
	}
	wg.Wait()
	require.Equal(t, 1, p.count(t, "readings"), "racing senders never create a second effect")
	require.EqualValues(t, 1, committed.Load())
	require.EqualValues(t, senders, committed.Load()+already.Load()+failed.Load())

	// Every loser converges on the same effect when it retries.
	for i := int32(0); i < failed.Load(); i++ {
		result, err := InsertGroupRow(t.Context(), p.db, request)
		require.NoError(t, err)
		require.True(t, result.AlreadyCommitted)
	}
	require.Equal(t, 1, p.count(t, "readings"))
}

func TestPostgresTargetReceiptIdentityUniqueKeyStrategy(t *testing.T) {
	p := newPostgresGroupTarget(t, pgUniqueTable, false)
	request := p.request(GroupEffectUniqueKey, 5)

	p.lose.Store(true)
	_, err := InsertGroupRow(t.Context(), p.db, request)
	var insertErr *GroupInsertError
	require.ErrorAs(t, err, &insertErr)
	require.True(t, insertErr.Ambiguous())
	p.lose.Store(false)

	again, err := InsertGroupRow(t.Context(), p.db, request)
	require.NoError(t, err)
	require.True(t, again.AlreadyCommitted, "the unique key recognizes the committed row (all typed cells compare equal)")
	require.Equal(t, 1, p.count(t, "readings"))

	_, err = InsertGroupRow(t.Context(), p.db, p.request(GroupEffectUniqueKey, 6))
	require.ErrorIs(t, err, ErrExistingRowDiffers)
}

func TestPostgresTargetReceiptIdentityNoDedupeAmbiguityIsNotRetriedByTheDatabaseLayer(t *testing.T) {
	p := newPostgresGroupTarget(t, pgPlainTable, false)
	request := p.request(GroupEffectNone, 5)

	p.lose.Store(true)
	_, err := InsertGroupRow(t.Context(), p.db, request)
	var insertErr *GroupInsertError
	require.ErrorAs(t, err, &insertErr)
	require.True(t, insertErr.Ambiguous())
	require.Equal(t, 1, p.count(t, "readings"))
	// The delivery state machine turns this into `unknown`; calling again here
	// proves why: PostgreSQL happily stores the repeat as a second row.
	p.lose.Store(false)
	_, err = InsertGroupRow(t.Context(), p.db, request)
	require.NoError(t, err)
	require.Equal(t, 2, p.count(t, "readings"))
}

func TestPostgresTargetReceiptIdentityTypeAndRangeErrorsArePreCommit(t *testing.T) {
	p := newPostgresGroupTarget(t, pgPlainTable, true)
	request := p.request(GroupEffectReceipt, math.MaxInt64)
	request.Row.Cells[1].Value = "not a number" // BIGINT column
	_, err := InsertGroupRow(t.Context(), p.db, request)
	var insertErr *GroupInsertError
	require.ErrorAs(t, err, &insertErr)
	require.Equal(t, InsertPhasePreCommit, insertErr.Phase)
	require.Zero(t, p.count(t, "readings"))
	require.Zero(t, p.count(t, EffectReceiptTable))
}
