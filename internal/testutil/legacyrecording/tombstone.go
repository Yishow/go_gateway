// Package legacyrecording builds persisted legacy recording fixtures in tests.
package legacyrecording

import (
	"context"
	"database/sql"
	"fmt"
)

// Tombstone models a group saved by the former removal behavior. It does not
// call the current permanent-delete API or alter the group's accepted data.
func Tombstone(ctx context.Context, db *sql.DB, groupID string) error {
	result, err := db.ExecContext(ctx, `UPDATE write_groups SET status='deleted' WHERE id=?`, groupID)
	if err != nil {
		return fmt.Errorf("seed legacy tombstone: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count legacy tombstone rows: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("legacy tombstone group not found")
	}
	return nil
}
