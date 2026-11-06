package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type blockedRecoveryConnector struct{ entered, release chan struct{} }

func (c blockedRecoveryConnector) Connect(context.Context) (driver.Conn, error) { return c, nil }
func (c blockedRecoveryConnector) Driver() driver.Driver                        { return blockedRecoveryDriver{} }
func (c blockedRecoveryConnector) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("not supported")
}
func (c blockedRecoveryConnector) Close() error              { return nil }
func (c blockedRecoveryConnector) Begin() (driver.Tx, error) { return nil, errors.New("not supported") }
func (c blockedRecoveryConnector) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	close(c.entered)
	<-c.release
	return &recoveryExistsRows{}, nil
}

type blockedRecoveryDriver struct{}

func (blockedRecoveryDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("not supported")
}

type recoveryExistsRows struct{ returned bool }

func (*recoveryExistsRows) Columns() []string { return []string{"exists"} }
func (*recoveryExistsRows) Close() error      { return nil }
func (r *recoveryExistsRows) Next(values []driver.Value) error {
	if r.returned {
		return io.EOF
	}
	r.returned = true
	values[0] = true
	return nil
}

func TestWriter_RecoveryOfOneConnectorDoesNotBlockAnother(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	slow := sql.OpenDB(blockedRecoveryConnector{entered: entered, release: release})
	t.Cleanup(func() { _ = slow.Close() })
	fast, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "fast.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = fast.Close() })
	_, err = fast.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
	require.NoError(t, err)
	writer := NewWriter(nil, nil)
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	slowResult, fastResult := make(chan error, 1), make(chan error, 1)
	go func() {
		slowResult <- writer.recoverWriterTable(t.Context(), slow, &schema.DatabaseConnector{ID: "slow", Kind: schema.DatabaseConnectorKindPostgres}, "public", "snapshot_values")
	}()
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		close(release)
		t.Fatal("slow recovery did not start")
	}
	go func() {
		fastResult <- writer.recoverWriterTable(t.Context(), fast, &schema.DatabaseConnector{ID: "fast", Kind: schema.DatabaseConnectorKindSQLite}, "main", "snapshot_values")
	}()
	select {
	case err := <-fastResult:
		require.NoError(t, err)
	case <-time.After(50 * time.Millisecond):
		close(release)
		<-slowResult
		<-fastResult
		t.Fatal("one connector's table recovery blocked another connector")
	}
	close(release)
	require.NoError(t, <-slowResult)
}
