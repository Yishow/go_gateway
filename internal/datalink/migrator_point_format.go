package datalink

import (
	"context"
	"database/sql"
	"fmt"

	"go-gateway/internal/datalink/schema/migrations"
)

// ensureSQLitePointFormatIntent adds nullable legacy intent atomically. Existing
// data_format values, revisions, and accepted measurements are not rewritten.
func ensureSQLitePointFormatIntent(db *sql.DB) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin point format migration: %w", err)
	}
	defer tx.Rollback()
	// Acquire a write lock before inspecting the schema on concurrent startup.
	if _, err = tx.ExecContext(ctx, "UPDATE points SET id=id WHERE 0"); err != nil {
		return fmt.Errorf("lock point format migration: %w", err)
	}
	var columns int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM pragma_table_info('points') WHERE name IN ('format_mode','format_override','format_version','revision_id')`).Scan(&columns); err != nil {
		return fmt.Errorf("inspect point format migration: %w", err)
	}
	switch columns {
	case 0:
		raw, err := migrations.FS.ReadFile("048_point_format_intent_sqlite.up.sql")
		if err != nil {
			return fmt.Errorf("read point format migration: %w", err)
		}
		if _, err = tx.ExecContext(ctx, string(raw)); err != nil {
			return fmt.Errorf("expand point format schema: %w", err)
		}
	case 4:
		// A repeat must not backfill intent or reset explicit choices.
	default:
		return fmt.Errorf("point format migration is incomplete: %d of 4 columns", columns)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit point format migration: %w", err)
	}
	return nil
}
