package sqltx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"modernc.org/sqlite"
)

// Inject only cleanup I/O failure; transactions and storage are real SQLite.
type cleanupFaultConnector struct {
	fault       error
	closes      int
	sawDeadline bool
}

func (c *cleanupFaultConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.Driver().Open(":memory:")
	if err != nil {
		return nil, err
	}
	return &cleanupFaultConn{Conn: conn, owner: c}, nil
}
func (c *cleanupFaultConnector) Driver() driver.Driver { return &sqlite.Driver{} }

type cleanupFaultConn struct {
	driver.Conn
	owner *cleanupFaultConnector
}

func (c *cleanupFaultConn) Close() error { c.owner.closes++; return c.Conn.Close() }
func (c *cleanupFaultConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if query == "ROLLBACK" {
		_, c.owner.sawDeadline = ctx.Deadline()
		if c.owner.fault != nil {
			return nil, c.owner.fault
		}
	}
	return c.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
}

func TestCleanupFailureDiscardsAndReportsMemoryDatabaseLoss(t *testing.T) {
	fault := errors.New("injected cleanup I/O failure")
	connector := &cleanupFaultConnector{fault: fault}
	db := sql.OpenDB(connector)
	defer db.Close()
	db.SetMaxOpenConns(1)
	conn, err := db.Conn(t.Context())
	require.NoError(t, err)
	_, err = conn.ExecContext(t.Context(), "CREATE TABLE probe(value INTEGER)")
	require.NoError(t, err)
	tx, err := conn.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	_, err = tx.Exec("INSERT INTO probe VALUES(7)")
	require.NoError(t, err)
	require.NoError(t, tx.Commit())
	lostAck := errors.New("commit acknowledgment lost")
	err = FinishSQLiteCommit(conn, lostAck)
	require.ErrorIs(t, err, ErrConnectionUnusable)
	require.ErrorIs(t, err, ErrCommitOutcomeUnknown)
	require.ErrorIs(t, err, fault)
	require.ErrorIs(t, err, lostAck)
	require.True(t, connector.sawDeadline, "cleanup uses a bounded context")
	require.Equal(t, 1, connector.closes, "unproven connection must not return to pool")
	require.ErrorIs(t, conn.PingContext(t.Context()), sql.ErrConnDone, "migration must not restore pragmas on discarded conn")
	// This is deliberately NOT recovery: the last private-memory handle was lost.
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM probe").Scan(&count)
	require.ErrorContains(t, err, "no such table: probe")
}

func TestRollbackCleanupFailureCannotBeHiddenAsTxDone(t *testing.T) {
	connector := &cleanupFaultConnector{fault: errors.New("cleanup failed")}
	db := sql.OpenDB(connector)
	defer db.Close()
	conn, err := db.Conn(t.Context())
	require.NoError(t, err)
	tx, err := conn.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, tx.Rollback())
	owned := &Transaction{Tx: tx, conn: conn}
	err = owned.Rollback()
	require.ErrorIs(t, err, ErrConnectionUnusable)
	require.False(t, errors.Is(err, sql.ErrTxDone), "usual defer guards must not suppress cleanup failure")
}
