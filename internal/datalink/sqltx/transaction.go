// Package sqltx owns SQLite connections until transaction cleanup is complete.
package sqltx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"
	"time"

	"modernc.org/sqlite"
)

// ErrCommitOutcomeUnknown forbids treating a failed acknowledgment as rollback.
var ErrCommitOutcomeUnknown = errors.New("transaction commit outcome unknown")

// ErrConnectionUnusable means cleanup could not prove the connection clean.
// Discarding it can destroy a private in-memory database; callers must fail
// closed rather than silently reopen or retry against an empty replacement.
var ErrConnectionUnusable = errors.New("transaction connection unusable; database availability may be lost")

// Transaction retains the connection beyond database/sql's terminal Tx state.
// Call Commit or Rollback, not the embedded SQL transaction's terminal methods.
type Transaction struct {
	Tx   *sql.Tx
	conn *sql.Conn
	mu   sync.Mutex
	done bool
}

// Begin retains a dedicated connection for modernc SQLite only. Other drivers
// keep database/sql's existing transaction behavior.
func Begin(ctx context.Context, db *sql.DB, opts *sql.TxOptions) (*Transaction, error) {
	if _, ok := db.Driver().(*sqlite.Driver); !ok {
		tx, err := db.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		return &Transaction{Tx: tx}, nil
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := conn.BeginTx(ctx, opts)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &Transaction{Tx: tx, conn: conn}, nil
}

// Commit preserves the original failure, even if native cleanup succeeds.
func (t *Transaction) Commit() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		return sql.ErrTxDone
	}
	t.done = true
	err := t.Tx.Commit()
	if t.conn == nil {
		return err
	}
	err = FinishSQLiteCommit(t.conn, err)
	return errors.Join(err, t.conn.Close())
}

// Rollback releases the connection only after native transaction cleanup.
func (t *Transaction) Rollback() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		return sql.ErrTxDone
	}
	t.done = true
	err := t.Tx.Rollback()
	if t.conn == nil {
		return err
	}
	if err != nil {
		cleanupErr := RollbackSQLite(t.conn)
		// Existing defer guards ignore ErrTxDone. Never attach that benign
		// sentinel to a real cleanup failure or it too would be suppressed.
		if cleanupErr != nil && errors.Is(err, sql.ErrTxDone) {
			err = cleanupErr
		} else {
			err = errors.Join(err, cleanupErr)
		}
	}
	return errors.Join(err, t.conn.Close())
}

// FinishSQLiteCommit cleans up a failed commit on a still-owned SQLite conn.
// It does not close conn: migration callers may still need to restore pragmas.
// A real SQLITE_BUSY plus a successful native rollback proves non-commit;
// all other failed acknowledgments retain an explicitly unknown outcome.
func FinishSQLiteCommit(conn *sql.Conn, commitErr error) error {
	if commitErr == nil {
		return nil
	}
	rolledBack, cleanupErr := rollbackSQLite(conn)
	var sqliteErr *sqlite.Error
	if cleanupErr == nil && rolledBack && errors.As(commitErr, &sqliteErr) && sqliteErr.Code()&255 == 5 {
		return commitErr
	}
	return errors.Join(ErrCommitOutcomeUnknown, commitErr, cleanupErr)
}

// RollbackSQLite resolves any native transaction on an exclusively owned conn.
// Call only after the sql.Tx terminal method, before returning conn to the pool.
func RollbackSQLite(conn *sql.Conn) error {
	_, err := rollbackSQLite(conn)
	return err
}

func rollbackSQLite(conn *sql.Conn) (bool, error) {
	// Caller cancellation must not prevent cleanup. Bound cleanup independently.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := conn.ExecContext(ctx, "ROLLBACK")
	if err == nil {
		return true, nil
	}
	var sqliteErr *sqlite.Error
	// Verified against modernc v1.44.2: SQLite SQLITE_ERROR (1), exact diagnostic.
	if errors.As(err, &sqliteErr) && sqliteErr.Code() == 1 && sqliteErr.Error() == "SQL logic error: cannot rollback - no transaction is active (1)" {
		return false, nil
	}
	discardErr := conn.Raw(func(any) error { return driver.ErrBadConn })
	return false, errors.Join(ErrConnectionUnusable, fmt.Errorf("cleanup SQLite transaction: %w", err), discardErr)
}
