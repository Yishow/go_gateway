package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/sqltx"
)

// inTx runs fn in one local transaction. A nil result means the transaction
// committed. Statement failures roll back atomically; a failed commit response
// may be unknown and must never be treated as proof of rollback.
func (s *Store) inTx(ctx context.Context, name string, fn func(context.Context, *sql.Tx) error) (resultErr error) {
	ctx, release, err := acquireStoreWriter(ctx, s.db)
	if err != nil {
		return fmt.Errorf("acquire %s transaction: %w", name, err)
	}
	defer release()
	owned, err := sqltx.Begin(ctx, s.db, nil)
	if err != nil {
		return fmt.Errorf("begin %s transaction: %w", name, err)
	}
	// Also run on callback panic/Goexit. The owned transaction's terminal
	// operation performs native cleanup before this defer returns the permit.
	defer func() { //nolint:contextcheck // sqltx terminal cleanup uses its own bounded context so caller cancellation cannot leak a transaction.
		panicValue := recover()
		cleanupErr := owned.Rollback()
		if errors.Is(cleanupErr, sql.ErrTxDone) {
			cleanupErr = nil
		}
		if panicValue != nil {
			if cleanupErr != nil {
				panic(fmt.Errorf("transaction callback panic (%v); cleanup failed: %w", panicValue, cleanupErr))
			}
			panic(panicValue)
		}
		if cleanupErr != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("rollback %s transaction: %w", name, cleanupErr))
		}
	}()
	tx := owned.Tx
	// Take the write lock first. A deferred transaction that reads and then
	// writes can fail immediately with SQLITE_BUSY when another writer commits
	// in between (busy_timeout does not apply to that upgrade); a statement that
	// starts as a write waits on the lock instead. The UPDATE matches no row.
	if _, err := tx.ExecContext(ctx, `UPDATE wg_delivery_checkpoints SET updated_at = updated_at WHERE 0`); err != nil {
		if rollbackErr := owned.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) { //nolint:contextcheck // sqltx terminal cleanup uses its own bounded context so caller cancellation cannot leak a transaction.
			err = errors.Join(err, fmt.Errorf("rollback %s transaction: %w", name, rollbackErr))
		}
		return fmt.Errorf("lock %s transaction: %w", name, err)
	}
	if err := fn(ctx, tx); err != nil {
		if rollbackErr := owned.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) { //nolint:contextcheck // sqltx terminal cleanup uses its own bounded context so caller cancellation cannot leak a transaction.
			return errors.Join(err, fmt.Errorf("rollback %s transaction: %w", name, rollbackErr))
		}
		return err
	}
	if err := owned.Commit(); err != nil { //nolint:contextcheck // sqltx terminal cleanup uses its own bounded context so caller cancellation cannot leak a transaction.
		return fmt.Errorf("commit %s transaction: %w", name, err)
	}
	return nil
}
