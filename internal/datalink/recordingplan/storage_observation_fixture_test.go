package recordingplan

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"modernc.org/sqlite"
)

// The ledger uses real SQLite; this wrapper injects only commit acknowledgement
// and post-commit read faults. It never fabricates a successful transaction.
type settlementFault struct {
	commitAck, readback, nextRead atomic.Bool
	queries                       atomic.Int64
	afterCommit                   func(settlementConn)
}
type settlementDriver struct{ fault *settlementFault }
type settlementConn struct {
	driver.Conn
	fault *settlementFault
}
type settlementTx struct {
	driver.Tx
	conn settlementConn
}

var settlementDriverSequence atomic.Uint64

func (d settlementDriver) Open(name string) (driver.Conn, error) {
	conn, err := (&sqlite.Driver{}).Open(name)
	if err != nil {
		return nil, err
	}
	return settlementConn{Conn: conn, fault: d.fault}, nil
}
func (c settlementConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}
func (c settlementConn) BeginTx(ctx context.Context, options driver.TxOptions) (driver.Tx, error) {
	tx, err := c.Conn.(driver.ConnBeginTx).BeginTx(ctx, options)
	if err != nil {
		return nil, err
	}
	return settlementTx{Tx: tx, conn: c}, nil
}
func (c settlementConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return c.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
}
func (c settlementConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	c.fault.queries.Add(1)
	if c.fault.nextRead.Swap(false) {
		return nil, errors.New("injected post-commit read failure")
	}
	return c.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
}
func (tx settlementTx) Commit() error {
	if err := tx.Tx.Commit(); err != nil {
		return err
	}
	if hook := tx.conn.fault.afterCommit; hook != nil {
		hook(tx.conn)
	}
	if tx.conn.fault.readback.Swap(false) {
		tx.conn.fault.nextRead.Store(true)
	}
	if tx.conn.fault.commitAck.Swap(false) {
		return errors.New("injected lost commit acknowledgement")
	}
	return nil
}
func namedValues(values ...string) []driver.NamedValue {
	out := make([]driver.NamedValue, len(values))
	for i, value := range values {
		out[i] = driver.NamedValue{Ordinal: i + 1, Value: value}
	}
	return out
}

type settlementFixture struct {
	db    *sql.DB
	repo  *SQLRepository
	fault *settlementFault
}

func newSettlementFixture(t *testing.T) settlementFixture {
	t.Helper()
	fault := &settlementFault{}
	name := fmt.Sprintf("settlement-sqlite-%d", settlementDriverSequence.Add(1))
	sql.Register(name, settlementDriver{fault: fault})
	db, err := sql.Open(name, filepath.Join(t.TempDir(), "ledger.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	applyRecordingDDL(t, db)
	return settlementFixture{db: db, repo: NewSQLRepository(db), fault: fault}
}
func (f settlementFixture) claim(t *testing.T, id string) *SchemaOperation {
	t.Helper()
	now := time.Now().UTC()
	op, outcome, err := f.repo.ClaimSchemaOperation(t.Context(), &SchemaOperation{
		OperationID: id, WorkspaceID: "workspace-1", ScopeKey: "scope-" + id,
		Token: "private-token-" + id, Owner: "private-owner-" + id, Action: StorageRepairAction,
		Status: SchemaOperationRunning, CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	return op
}
