package workspace

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// A rollback-journal reader deterministically makes SQLite COMMIT return BUSY.
// database/sql marks Tx done even though the native transaction remains open.
func TestDatabaseSetupBusyCommitDoesNotPoisonNextTransaction(t *testing.T) {
	path := filepath.Join(t.TempDir(), "commit-busy.db")
	writerDB, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer writerDB.Close()
	writerDB.SetMaxOpenConns(1)
	readerDB, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer readerDB.Close()
	_, err = writerDB.ExecContext(t.Context(), `CREATE TABLE probe(value INTEGER); INSERT INTO probe VALUES(1)`)
	require.NoError(t, err)
	repo := NewSQLRepository(writerDB)
	writer, err := repo.BeginDatabaseSetup(t.Context())
	require.NoError(t, err)
	_, err = writer.SQLTx().ExecContext(t.Context(), `UPDATE probe SET value=2`)
	require.NoError(t, err)
	reader, err := readerDB.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	var value int
	require.NoError(t, reader.QueryRowContext(t.Context(), `SELECT value FROM probe`).Scan(&value))
	require.Equal(t, 1, value)
	err = writer.Commit()
	require.Error(t, err)
	require.True(t, isSQLiteBusyError(err), "known busy remains retryable only after proven rollback")
	_ = writer.Rollback()
	require.NoError(t, reader.Rollback())
	next, err := repo.BeginDatabaseSetup(t.Context())
	require.NoError(t, err, "failed commit must not leave a native transaction in the pool")
	require.NoError(t, next.SQLTx().QueryRowContext(t.Context(), `SELECT value FROM probe`).Scan(&value))
	require.Equal(t, 1, value, "failed write was rolled back")
	require.NoError(t, next.Rollback())
}
