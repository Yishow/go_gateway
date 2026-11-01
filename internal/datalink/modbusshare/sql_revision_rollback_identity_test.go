package modbusshare

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type revisionRollbackConnector struct {
	queryErr, rollbackErr error
	queries, rollbacks    int
}

func (c *revisionRollbackConnector) Connect(context.Context) (driver.Conn, error) {
	return revisionRollbackConn{c}, nil
}
func (c *revisionRollbackConnector) Driver() driver.Driver { return revisionRollbackDriver{c} }

type revisionRollbackDriver struct{ owner *revisionRollbackConnector }

func (d revisionRollbackDriver) Open(string) (driver.Conn, error) {
	return revisionRollbackConn(d), nil
}

type revisionRollbackConn struct{ owner *revisionRollbackConnector }

func (c revisionRollbackConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepare")
}
func (c revisionRollbackConn) Close() error              { return nil }
func (c revisionRollbackConn) Begin() (driver.Tx, error) { return revisionRollbackTx(c), nil }
func (c revisionRollbackConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	c.owner.queries++
	return nil, c.owner.queryErr
}

type revisionRollbackTx struct{ owner *revisionRollbackConnector }

func (tx revisionRollbackTx) Commit() error   { return errors.New("unexpected commit") }
func (tx revisionRollbackTx) Rollback() error { tx.owner.rollbacks++; return tx.owner.rollbackErr }

func TestWorkspaceRevisionRollbackIgnoresOnlyBareTxDone(t *testing.T) {
	primary := errors.New("read failed")
	cleanup := errors.New("rollback driver cause")
	wrapped := fmt.Errorf("driver wrapped terminal: %w", sql.ErrTxDone)
	joined := errors.Join(sql.ErrTxDone, cleanup)
	for _, operation := range []string{"desired", "revision", "dirty"} {
		t.Run(operation, func(t *testing.T) {
			for _, tc := range []struct {
				name     string
				rollback error
				retained bool
			}{
				{"nil", nil, false}, {"bare", sql.ErrTxDone, false}, {"wrapped", wrapped, true}, {"joined", joined, true},
			} {
				t.Run(tc.name, func(t *testing.T) {
					connector := &revisionRollbackConnector{queryErr: primary, rollbackErr: tc.rollback}
					db := sql.OpenDB(connector)
					t.Cleanup(func() { require.NoError(t, db.Close()) })
					store := NewSQLWorkspaceRevisionStore(db)
					var err error
					switch operation {
					case "desired":
						err = store.CommitDesiredMappings(t.Context(), "ws", "rev-1", "rev-2", nil)
					case "revision":
						err = store.UpdateRevision(t.Context(), "ws", "rev-1", "rev-2")
					default:
						err = store.MarkDirty(t.Context(), "ws")
					}
					require.ErrorIs(t, err, primary)
					require.Equal(t, 1, connector.queries)
					require.Equal(t, 1, connector.rollbacks)
					require.Equal(t, tc.retained, errors.Is(err, sql.ErrTxDone))
					if tc.retained {
						require.ErrorIs(t, err, tc.rollback)
					}
					require.Equal(t, tc.name == "joined", errors.Is(err, cleanup))
				})
			}
		})
	}
}
