package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// inTx runs fn in one local transaction. A nil result means the transaction
// committed; any error rolls everything back, so partial state never persists.
func (s *Store) inTx(ctx context.Context, name string, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin %s transaction: %w", name, err)
	}
	// Take the write lock first. A deferred transaction that reads and then
	// writes can fail immediately with SQLITE_BUSY when another writer commits
	// in between (busy_timeout does not apply to that upgrade); a statement that
	// starts as a write waits on the lock instead. The UPDATE matches no row.
	if _, err := tx.ExecContext(ctx, `UPDATE wg_delivery_checkpoints SET updated_at = updated_at WHERE 0`); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("rollback %s transaction: %w", name, rollbackErr))
		}
		return fmt.Errorf("lock %s transaction: %w", name, err)
	}
	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			return errors.Join(err, fmt.Errorf("rollback %s transaction: %w", name, rollbackErr))
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit %s transaction: %w", name, err)
	}
	return nil
}
