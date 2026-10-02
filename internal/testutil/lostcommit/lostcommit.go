// Package lostcommit provides a SQLite database whose transactions can commit
// for real and then report a failure, which is what a lost network response
// looks like to the caller. It exists for delivery fault tests only.
package lostcommit

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"sync/atomic"

	"modernc.org/sqlite"
)

// ErrResponseLost is returned by Commit while the switch is on.
var ErrResponseLost = errors.New("connection reset by peer while reading the commit response")

// Open opens a SQLite file. While lose is true every Commit succeeds on the
// database and then returns ErrResponseLost.
func Open(path string, lose *atomic.Bool) *sql.DB {
	return Wrap(&sqlite.Driver{}, path, lose)
}

// Wrap applies the same lost-response behavior to any database/sql driver, for
// example the pgx stdlib driver against a disposable PostgreSQL.
func Wrap(inner driver.Driver, dsn string, lose *atomic.Bool) *sql.DB {
	return WrapFaults(inner, dsn, &Faults{Lose: lose})
}

// Faults are the switches a test can flip while the database is in use.
type Faults struct {
	// Lose makes Commit succeed on the database and then return ErrResponseLost.
	Lose *atomic.Bool
	// ExecError, when it holds a non-nil error, makes every Exec fail with it
	// before reaching the database (a transient failure such as a reset).
	ExecError atomic.Pointer[error]
	// AfterCommit, when set, runs right after a transaction committed for real
	// and before Commit returns, e.g. to cancel the caller's context.
	AfterCommit atomic.Pointer[func()]
	// Reject, when set, is asked about every statement before it reaches the
	// database; a non-nil error is returned instead (e.g. a permission denial
	// for SELECT or DELETE on one table).
	Reject atomic.Pointer[func(query string) error]
}

// WrapFaults is Wrap with the full set of fault switches.
func WrapFaults(inner driver.Driver, dsn string, faults *Faults) *sql.DB {
	if faults.Lose == nil {
		faults.Lose = &atomic.Bool{}
	}
	return sql.OpenDB(&connector{inner: inner, path: dsn, faults: faults})
}

// OpenFaults opens a SQLite file with the full set of fault switches.
func OpenFaults(path string, faults *Faults) *sql.DB {
	return WrapFaults(&sqlite.Driver{}, path, faults)
}

type connector struct {
	inner  driver.Driver
	path   string
	faults *Faults
}

func (c *connector) Connect(context.Context) (driver.Conn, error) {
	conn, err := c.inner.Open(c.path)
	if err != nil {
		return nil, err
	}
	return &lostConn{Conn: conn, faults: c.faults}, nil
}

func (c *connector) Driver() driver.Driver { return c.inner }

type lostConn struct {
	driver.Conn
	faults *Faults
}

func (c *lostConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	begin, ok := c.Conn.(driver.ConnBeginTx)
	if !ok {
		return nil, errors.New("inner connection cannot begin transactions")
	}
	tx, err := begin.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &lostTx{Tx: tx, lose: c.faults.Lose, faults: c.faults}, nil
}

func (c *lostConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if injected := c.faults.ExecError.Load(); injected != nil {
		return nil, *injected
	}
	if err := c.faults.rejected(query); err != nil {
		return nil, err
	}
	exec, ok := c.Conn.(driver.ExecerContext)
	if !ok {
		return nil, errors.New("inner connection cannot exec")
	}
	return exec.ExecContext(ctx, query, args)
}

func (c *lostConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if err := c.faults.rejected(query); err != nil {
		return nil, err
	}
	queryer, ok := c.Conn.(driver.QueryerContext)
	if !ok {
		return nil, errors.New("inner connection cannot query")
	}
	return queryer.QueryContext(ctx, query, args)
}

func (f *Faults) rejected(query string) error {
	if reject := f.Reject.Load(); reject != nil {
		return (*reject)(query)
	}
	return nil
}

type lostTx struct {
	driver.Tx
	lose   *atomic.Bool
	faults *Faults
}

func (t *lostTx) Commit() error {
	if err := t.Tx.Commit(); err != nil {
		return err
	}
	if hook := t.faults.AfterCommit.Load(); hook != nil {
		(*hook)()
	}
	if t.lose.Load() {
		return ErrResponseLost
	}
	return nil
}
