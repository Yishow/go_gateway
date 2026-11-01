package recordingplan

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type recoveryRollbackError struct{ phase string }

func (e *recoveryRollbackError) Error() string { return "private recovery storage fault: " + e.phase }

type recoveryRollbackConnector struct {
	queryErr, rollbackErr error
	commits, rollbacks    int
}

func (c *recoveryRollbackConnector) Connect(context.Context) (driver.Conn, error) {
	return &recoveryRollbackConn{owner: c}, nil
}
func (c *recoveryRollbackConnector) Driver() driver.Driver { return recoveryRollbackDriver{owner: c} }

type recoveryRollbackDriver struct{ owner *recoveryRollbackConnector }

func (d recoveryRollbackDriver) Open(string) (driver.Conn, error) {
	return &recoveryRollbackConn{owner: d.owner}, nil
}

type recoveryRollbackConn struct {
	driver.Conn
	owner *recoveryRollbackConnector
}

func (*recoveryRollbackConn) Close() error { return nil }
func (c *recoveryRollbackConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return recoveryRollbackTx{owner: c.owner}, nil
}
func (c *recoveryRollbackConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	if c.owner.queryErr != nil {
		return nil, c.owner.queryErr
	}
	return recoveryEmptyRows{}, nil
}

type recoveryRollbackTx struct{ owner *recoveryRollbackConnector }

func (tx recoveryRollbackTx) Commit() error   { tx.owner.commits++; return nil }
func (tx recoveryRollbackTx) Rollback() error { tx.owner.rollbacks++; return tx.owner.rollbackErr }

type recoveryEmptyRows struct{}

func (recoveryEmptyRows) Columns() []string         { return nil }
func (recoveryEmptyRows) Close() error              { return nil }
func (recoveryEmptyRows) Next([]driver.Value) error { return io.EOF }

func TestTestWriteRecoveryRollbackFailureRetainsPrimaryAndSanitizes(t *testing.T) {
	primary := &recoveryRollbackError{phase: "query"}
	cleanup := &recoveryRollbackError{phase: "rollback"}
	owner := &recoveryRollbackConnector{queryErr: primary, rollbackErr: cleanup}
	db := sql.OpenDB(owner)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	repo := NewSQLRepository(db)
	rows, err := repo.readTestWriteRecovery(t.Context(), "workspace")
	require.Nil(t, rows)
	require.ErrorIs(t, err, primary)
	require.ErrorIs(t, err, cleanup)
	var typed *recoveryRollbackError
	require.ErrorAs(t, err, &typed)
	require.Same(t, primary, typed)
	require.Zero(t, owner.commits)
	require.Equal(t, 1, owner.rollbacks)
	result, err := NewService(repo).ListUnresolvedTestWrites(t.Context(), "workspace")
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrTestWriteRecoveryUnavailable)
	require.NotContains(t, err.Error(), "private recovery storage fault")
	require.Equal(t, 2, owner.rollbacks)
}
func TestTestWriteRecoveryCommittedReadDoesNotBecomeCleanupFailure(t *testing.T) {
	owner := &recoveryRollbackConnector{rollbackErr: &recoveryRollbackError{phase: "unused"}}
	db := sql.OpenDB(owner)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	records, err := NewSQLRepository(db).readTestWriteRecovery(t.Context(), "workspace")
	require.NoError(t, err)
	require.Empty(t, records)
	require.Equal(t, 1, owner.commits)
	require.Zero(t, owner.rollbacks, "database/sql returns ErrTxDone without invoking rollback after Commit")
}
