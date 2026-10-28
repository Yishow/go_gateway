package grouppipeline

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestBacklogGuardIterationFailureIsNotEmptyBacklog(t *testing.T) {
	for _, table := range []string{"wg_delivery_outbox", "wg_delivery_samples"} {
		t.Run(table, func(t *testing.T) {
			failure := errors.New("injected iteration failure")
			db := sql.OpenDB(backlogIterationConnector{table: table, failure: failure})
			t.Cleanup(func() { require.NoError(t, db.Close()) })
			err := checkBacklogOwnership(t, db, &workspace.WriteGroup{ID: "G1", WorkspaceID: "workspace-A"})
			require.ErrorIs(t, err, failure, "unknown rows cannot be reported as a proven empty backlog")
		})
	}
}

type backlogIterationConnector struct {
	table   string
	failure error
}

func (c backlogIterationConnector) Connect(context.Context) (driver.Conn, error) {
	return backlogIterationConn{c}, nil
}
func (backlogIterationConnector) Driver() driver.Driver { return backlogIterationDriver{} }

type backlogIterationDriver struct{}

func (backlogIterationDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("use test connector")
}

type backlogIterationConn struct{ backlogIterationConnector }

func (backlogIterationConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepare")
}
func (backlogIterationConn) Close() error              { return nil }
func (backlogIterationConn) Begin() (driver.Tx, error) { return backlogIterationTx{}, nil }
func (c backlogIterationConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	if !strings.Contains(query, "FROM wg_delivery_outbox") && !strings.Contains(query, "FROM wg_delivery_samples") {
		return nil, errors.New("unexpected query")
	}
	var failure error
	if strings.Contains(query, "FROM "+c.table) {
		failure = c.failure
	}
	return &backlogIterationRows{failure: failure}, nil
}

type backlogIterationTx struct{}

func (backlogIterationTx) Commit() error   { return nil }
func (backlogIterationTx) Rollback() error { return nil }

type backlogIterationRows struct{ failure error }

func (*backlogIterationRows) Columns() []string { return []string{"unused"} }
func (*backlogIterationRows) Close() error      { return nil }
func (r *backlogIterationRows) Next([]driver.Value) error {
	if r.failure != nil {
		return r.failure
	}
	return io.EOF
}
