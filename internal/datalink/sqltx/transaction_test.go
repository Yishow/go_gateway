package sqltx

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestCanceledMemoryTransactionPreservesDatabase(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	_, err = db.Exec("CREATE TABLE probe(value INTEGER); INSERT INTO probe VALUES(1)")
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(t.Context())
	tx, err := Begin(ctx, db, nil)
	require.NoError(t, err)
	_, err = tx.Tx.ExecContext(ctx, "UPDATE probe SET value=2")
	require.NoError(t, err)
	cancel()
	require.Error(t, tx.Commit())
	next, err := Begin(t.Context(), db, nil)
	require.NoError(t, err)
	defer next.Rollback()
	var value int
	require.NoError(t, next.Tx.QueryRow("SELECT value FROM probe").Scan(&value))
	require.Equal(t, 1, value)
}

func TestCleanupAlreadyCommittedIsNotCommitSuccess(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	conn, err := db.Conn(t.Context())
	require.NoError(t, err)
	defer conn.Close()
	_, err = conn.ExecContext(t.Context(), "CREATE TABLE probe(value INTEGER)")
	require.NoError(t, err)
	tx, err := conn.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	_, err = tx.Exec("INSERT INTO probe VALUES(7)")
	require.NoError(t, err)
	require.NoError(t, tx.Commit())
	uncertain := errors.New("database is locked: acknowledgment lost")
	err = FinishSQLiteCommit(conn, uncertain)
	require.ErrorIs(t, err, ErrCommitOutcomeUnknown)
	require.ErrorIs(t, err, uncertain)
	next, err := conn.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer next.Rollback()
	var value int
	require.NoError(t, next.QueryRow("SELECT value FROM probe").Scan(&value))
	require.Equal(t, 7, value)
}

func TestCleanupUnavailableConnectionFailsClosed(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	conn, err := db.Conn(t.Context())
	require.NoError(t, err)
	require.NoError(t, conn.Close())
	uncertain := errors.New("commit failed")
	err = FinishSQLiteCommit(conn, uncertain)
	require.ErrorIs(t, err, ErrConnectionUnusable)
	require.ErrorIs(t, err, ErrCommitOutcomeUnknown)
	require.ErrorIs(t, err, uncertain)
}
