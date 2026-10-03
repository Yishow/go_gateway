package dbtarget

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"modernc.org/sqlite"
)

// InsertErrorClass says what a failed insert means for delivery.
type InsertErrorClass string

const (
	// InsertErrorTransient may succeed later without any change (connection,
	// timeout, lock, resource pressure). Unknown errors are transient so a
	// bounded retry, not an assumption, decides their fate.
	InsertErrorTransient InsertErrorClass = "transient"
	// InsertErrorRow means the destination rejects this row's values (data or
	// integrity violation); retrying the same row can never succeed.
	InsertErrorRow InsertErrorClass = "row_rejected"
	// InsertErrorTarget means the destination cannot accept any row until it is
	// repaired (credential, permission, missing table or column).
	InsertErrorTarget InsertErrorClass = "target_unusable"
)

// SQLite primary result codes that matter for classification.
const (
	sqliteError    = 1
	sqliteBusy     = 5
	sqliteLocked   = 6
	sqliteReadOnly = 8
	sqliteIOErr    = 10
	sqliteFull     = 13
	sqliteCantOpen = 14
	sqliteTooBig   = 18
	sqliteConstr   = 19
	sqliteMismatch = 20
	sqliteAuth     = 23
)

// ClassifyInsertError classifies a failed InsertGroupRow error. It relies on
// driver error codes (PostgreSQL SQLSTATE, SQLite result codes), never on
// message text, so classification cannot leak or depend on row contents.
func ClassifyInsertError(err error) InsertErrorClass {
	if err == nil {
		return InsertErrorTransient
	}
	if errors.Is(err, ErrReceiptDigestMismatch) || errors.Is(err, ErrExistingRowDiffers) {
		return InsertErrorRow
	}
	if errors.Is(err, ErrGroupInsertUnsupported) {
		return InsertErrorTarget
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return classifySQLSTATE(pgErr.Code)
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return classifySQLSTATE(string(pqErr.Code))
	}
	var liteErr *sqlite.Error
	if errors.As(err, &liteErr) {
		return classifySQLiteCode(liteErr)
	}
	return InsertErrorTransient
}

func classifySQLSTATE(code string) InsertErrorClass {
	if len(code) < 2 {
		return InsertErrorTransient
	}
	switch code[:2] {
	case "22", "23": // data exception, integrity constraint violation
		return InsertErrorRow
	case "28", "3D", "3F", "42": // authorization, catalog, schema, syntax/access rule
		return InsertErrorTarget
	}
	return InsertErrorTransient
}

func classifySQLiteCode(err *sqlite.Error) InsertErrorClass {
	switch err.Code() & 0xff {
	case sqliteConstr, sqliteMismatch, sqliteTooBig:
		return InsertErrorRow
	case sqliteReadOnly, sqliteAuth:
		return InsertErrorTarget
	case sqliteBusy, sqliteLocked, sqliteIOErr, sqliteFull, sqliteCantOpen:
		return InsertErrorTransient
	case sqliteError:
		// The generic code covers a missing table or column as well as other
		// statement errors; all of them need a schema repair, not a retry.
		return InsertErrorTarget
	}
	return InsertErrorTransient
}
