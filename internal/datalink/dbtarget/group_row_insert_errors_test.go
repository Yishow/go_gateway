package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestBoundedRetryAndPoisonPartitionClassifiesPostgresSQLSTATE(t *testing.T) {
	cases := map[string]InsertErrorClass{
		"22P02": InsertErrorRow,       // invalid text representation
		"22003": InsertErrorRow,       // numeric value out of range
		"22001": InsertErrorRow,       // string data right truncation
		"23502": InsertErrorRow,       // not null violation
		"23514": InsertErrorRow,       // check violation
		"23505": InsertErrorRow,       // unique violation
		"23503": InsertErrorRow,       // foreign key violation
		"42P01": InsertErrorTarget,    // undefined table
		"42703": InsertErrorTarget,    // undefined column
		"42501": InsertErrorTarget,    // insufficient privilege
		"28P01": InsertErrorTarget,    // invalid password
		"28000": InsertErrorTarget,    // invalid authorization
		"3D000": InsertErrorTarget,    // invalid catalog name
		"3F000": InsertErrorTarget,    // invalid schema name
		"40001": InsertErrorTransient, // serialization failure
		"40P01": InsertErrorTransient, // deadlock detected
		"53300": InsertErrorTransient, // too many connections
		"57P01": InsertErrorTransient, // admin shutdown
		"08006": InsertErrorTransient, // connection failure
		"XX000": InsertErrorTransient, // unknown stays transient and bounded
	}
	for code, want := range cases {
		wrapped := &GroupInsertError{Phase: InsertPhasePreCommit, Cause: fmt.Errorf("insert destination row: %w", &pgconn.PgError{Code: code, Message: "secret detail"})}
		require.Equal(t, want, ClassifyInsertError(wrapped), code)
	}
}

func TestBoundedRetryAndPoisonPartitionTransientAndUnknownErrors(t *testing.T) {
	require.Equal(t, InsertErrorTransient, ClassifyInsertError(nil))
	require.Equal(t, InsertErrorTransient, ClassifyInsertError(errors.New("connection reset by peer")))
	require.Equal(t, InsertErrorTransient, ClassifyInsertError(context.DeadlineExceeded))
	require.Equal(t, InsertErrorTransient, ClassifyInsertError(&GroupInsertError{Phase: InsertPhaseCommit, Cause: errors.New("lost")}))
	require.Equal(t, InsertErrorTarget, ClassifyInsertError(fmt.Errorf("%w: kind", ErrGroupInsertUnsupported)))
	require.Equal(t, InsertErrorRow, ClassifyInsertError(ErrExistingRowDiffers))
	require.Equal(t, InsertErrorRow, ClassifyInsertError(ErrReceiptDigestMismatch))
}

func TestBoundedRetryAndPoisonPartitionClassifiesRealSQLiteErrors(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "target.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	_, err = db.ExecContext(t.Context(), `CREATE TABLE readings (record_id TEXT, counter INTEGER NOT NULL CHECK (counter >= 0), label TEXT)`)
	require.NoError(t, err)

	request := func(cells ...EncodedCell) GroupInsertRequest {
		return GroupInsertRequest{
			Kind: schema.DatabaseConnectorKindSQLite, TableName: "readings", Strategy: GroupEffectNone,
			Row: EncodedRow{EffectKey: "effect-1", Cells: cells},
		}
	}
	_, err = InsertGroupRow(t.Context(), db, request(EncodedCell{Column: "counter", Value: nil}))
	require.Equal(t, InsertErrorRow, ClassifyInsertError(err), "NOT NULL violation is the row's fault")
	_, err = InsertGroupRow(t.Context(), db, request(EncodedCell{Column: "counter", Value: int64(-1)}))
	require.Equal(t, InsertErrorRow, ClassifyInsertError(err), "CHECK violation")
	_, err = InsertGroupRow(t.Context(), db, request(EncodedCell{Column: "no_such_column", Value: int64(1)}))
	require.Equal(t, InsertErrorTarget, ClassifyInsertError(err), "a missing column needs a schema repair")
	missingTable := request(EncodedCell{Column: "counter", Value: int64(1)})
	missingTable.TableName = "absent"
	_, err = InsertGroupRow(t.Context(), db, missingTable)
	require.Equal(t, InsertErrorTarget, ClassifyInsertError(err))

	_, err = db.ExecContext(t.Context(), `CREATE TRIGGER refuse BEFORE INSERT ON readings BEGIN SELECT RAISE(ABORT, 'refused'); END`)
	require.NoError(t, err)
	_, err = InsertGroupRow(t.Context(), db, request(EncodedCell{Column: "counter", Value: int64(1)}))
	require.Equal(t, InsertErrorRow, ClassifyInsertError(err), "a destination rule that rejects the row")
	require.NotContains(t, err.Error(), "refused", "the classification does not depend on and never exposes message text")
}

func TestBoundedRetryAndPoisonPartitionClassifiesRealPostgresErrors(t *testing.T) {
	p := newPostgresGroupTarget(t, pgPlainTable, false)
	notNull := p.request(GroupEffectNone, 5)
	notNull.Row.Cells[1].Value = nil // counter BIGINT NOT NULL
	_, err := InsertGroupRow(t.Context(), p.db, notNull)
	require.Equal(t, InsertErrorRow, ClassifyInsertError(err))

	badType := p.request(GroupEffectNone, 5)
	badType.Row.Cells[1].Value = "not a number"
	_, err = InsertGroupRow(t.Context(), p.db, badType)
	require.Equal(t, InsertErrorRow, ClassifyInsertError(err))

	missingColumn := p.request(GroupEffectNone, 5)
	missingColumn.Row.Cells[2].Column = "no_such_column"
	_, err = InsertGroupRow(t.Context(), p.db, missingColumn)
	require.Equal(t, InsertErrorTarget, ClassifyInsertError(err))

	missingTable := p.request(GroupEffectNone, 5)
	missingTable.TableName = "absent"
	_, err = InsertGroupRow(t.Context(), p.db, missingTable)
	require.Equal(t, InsertErrorTarget, ClassifyInsertError(err))
	require.Zero(t, p.count(t, "readings"))
}
