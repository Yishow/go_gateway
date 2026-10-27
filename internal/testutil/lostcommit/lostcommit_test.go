package lostcommit

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"testing"
)

// purePrepareConn simulates a driver connection like github.com/microsoft/go-mssqldb
// which implements ConnPrepareContext and ConnBeginTx, but NOT ExecerContext or QueryerContext.
type purePrepareConn struct {
	lastQuery string
}

func (c *purePrepareConn) Prepare(query string) (driver.Stmt, error) {
	c.lastQuery = query
	return &purePrepareStmt{query: query}, nil
}

func (c *purePrepareConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	c.lastQuery = query
	return &purePrepareStmt{query: query}, nil
}

func (c *purePrepareConn) Close() error { return nil }

func (c *purePrepareConn) Begin() (driver.Tx, error) {
	return &purePrepareTx{}, nil
}

func (c *purePrepareConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return &purePrepareTx{}, nil
}

type purePrepareTx struct{}

func (t *purePrepareTx) Commit() error   { return nil }
func (t *purePrepareTx) Rollback() error { return nil }

type purePrepareStmt struct {
	query string
}

func (s *purePrepareStmt) Close() error { return nil }
func (s *purePrepareStmt) NumInput() int { return -1 }

func (s *purePrepareStmt) Exec(args []driver.Value) (driver.Result, error) {
	return driver.RowsAffected(1), nil
}

func (s *purePrepareStmt) Query(args []driver.Value) (driver.Rows, error) {
	return &purePrepareRows{read: false}, nil
}

type purePrepareRows struct {
	read bool
}

func (r *purePrepareRows) Columns() []string { return []string{"val"} }
func (r *purePrepareRows) Close() error      { return nil }
func (r *purePrepareRows) Next(dest []driver.Value) error {
	if r.read {
		return io.EOF
	}
	r.read = true
	dest[0] = int64(42)
	return nil
}

type purePrepareDriver struct{}

func (d *purePrepareDriver) Open(name string) (driver.Conn, error) {
	return &purePrepareConn{}, nil
}

func TestLostConnPurePrepareDriverFallbackAndReject(t *testing.T) {
	faults := &Faults{}
	var db *sql.DB = WrapFaults(&purePrepareDriver{}, "test-dsn", faults)
	t.Cleanup(func() { _ = db.Close() })

	// 1. Without fast-path QueryerContext, database/sql should fallback cleanly through ErrSkip.
	var val int64
	err := db.QueryRowContext(t.Context(), "SELECT 42").Scan(&val)
	if err != nil {
		t.Fatalf("expected fallback query to succeed, got %v", err)
	}
	if val != 42 {
		t.Fatalf("expected val=42, got %d", val)
	}

	// 2. Fallback ExecContext should also succeed.
	res, err := db.ExecContext(t.Context(), "UPDATE tbl SET x=1")
	if err != nil {
		t.Fatalf("expected fallback exec to succeed, got %v", err)
	}
	affected, err := res.RowsAffected()
	if err != nil || affected != 1 {
		t.Fatalf("expected 1 affected row, got %d (err: %v)", affected, err)
	}

	// 3. Reject hook must fire and reject statements even under the fallback path.
	errDenied := errors.New("denied by hook")
	rejectHook := func(query string) error {
		if query == "SELECT FORBIDDEN" {
			return errDenied
		}
		return nil
	}
	faults.Reject.Store(&rejectHook)

	err = db.QueryRowContext(t.Context(), "SELECT FORBIDDEN").Scan(&val)
	if !errors.Is(err, errDenied) {
		t.Fatalf("expected errDenied, got %v", err)
	}

	// Unrejected query still succeeds.
	err = db.QueryRowContext(t.Context(), "SELECT ALLOWED").Scan(&val)
	if err != nil {
		t.Fatalf("expected allowed query to succeed, got %v", err)
	}
}
