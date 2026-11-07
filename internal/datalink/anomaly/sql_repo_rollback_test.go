package anomaly

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type listFailureConnector struct{ queryErr, rollbackErr error }

func (c listFailureConnector) Connect(context.Context) (driver.Conn, error) {
	return listFailureConn{c}, nil
}
func (c listFailureConnector) Driver() driver.Driver { return listFailureDriver{c} }

type listFailureDriver struct{ connector listFailureConnector }

func (d listFailureDriver) Open(string) (driver.Conn, error) {
	return d.connector.Connect(context.Background())
}

type listFailureConn struct{ connector listFailureConnector }

func (c listFailureConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare unsupported")
}
func (c listFailureConn) Close() error { return nil }
func (c listFailureConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return c.Begin()
}
func (c listFailureConn) Begin() (driver.Tx, error) {
	return listFailureTx{c.connector.rollbackErr}, nil
}
func (c listFailureConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return nil, c.connector.queryErr
}

type listFailureTx struct{ rollbackErr error }

func (listFailureTx) Commit() error      { return nil }
func (tx listFailureTx) Rollback() error { return tx.rollbackErr }

func TestSQLRepositoryListPreservesRollbackFailure(t *testing.T) {
	queryErr, rollbackErr := errors.New("query failed"), errors.New("rollback failed")
	db := sql.OpenDB(listFailureConnector{queryErr: queryErr, rollbackErr: rollbackErr})
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	events, total, err := NewSQLRepository(db).List(t.Context(), ListFilter{DeviceID: "dev-a"})
	require.Nil(t, events)
	require.Zero(t, total)
	require.ErrorIs(t, err, queryErr)
	require.ErrorIs(t, err, rollbackErr)
}
