package sqltx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type foreignConnector struct {
	err                       error
	commits, rollbacks, execs int
}

func (c *foreignConnector) Connect(context.Context) (driver.Conn, error) { return foreignConn{c}, nil }
func (c *foreignConnector) Driver() driver.Driver                        { return foreignDriver{c} }

type foreignDriver struct{ c *foreignConnector }

func (d foreignDriver) Open(string) (driver.Conn, error) { return foreignConn{d.c}, nil }

type foreignConn struct{ c *foreignConnector }

func (c foreignConn) Close() error { return nil }
func (c foreignConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepare")
}
func (c foreignConn) Begin() (driver.Tx, error) { return foreignTx{c.c}, nil }
func (c foreignConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	c.c.execs++
	return nil, errors.New("unexpected SQLite cleanup on foreign driver")
}

type foreignTx struct{ c *foreignConnector }

func (tx foreignTx) Commit() error   { tx.c.commits++; return tx.c.err }
func (tx foreignTx) Rollback() error { tx.c.rollbacks++; return tx.c.err }

func TestNonSQLiteRetainsDriverTransactionBehavior(t *testing.T) {
	for _, commit := range []bool{true, false} {
		t.Run(map[bool]string{true: "commit", false: "rollback"}[commit], func(t *testing.T) {
			expected := errors.New("foreign driver failure")
			connector := &foreignConnector{err: expected}
			db := sql.OpenDB(connector)
			defer db.Close()
			tx, err := Begin(context.Background(), db, nil)
			require.NoError(t, err)
			if commit {
				err = tx.Commit()
				require.Equal(t, 1, connector.commits)
			} else {
				err = tx.Rollback()
				require.Equal(t, 1, connector.rollbacks)
			}
			require.Same(t, expected, err)
			require.ErrorIs(t, tx.Rollback(), sql.ErrTxDone)
			require.Zero(t, connector.execs)
		})
	}
}
