package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema/migrations"
)

// ensureSQLiteMemberCalculationIntent adds nullable intent atomically. Legacy
// intervals, reductions, roles, revisions and SQL values are not rewritten.
func ensureSQLiteMemberCalculationIntent(db *sql.DB) error {
	return expandSQLiteMixedProfile(db, nil)
}

func addSQLiteMemberCalculationColumns(ctx context.Context, tx *sql.Tx) error {
	// Each nullable column is independently additive. Existing repair migrations
	// may recreate the member table; restore a missing column without backfill.
	raw, err := migrations.FS.ReadFile("049_member_calculation_intent_sqlite.up.sql")
	if err != nil {
		return fmt.Errorf("read member calculation migration: %w", err)
	}
	statements := strings.Split(string(raw), ";")
	for i, item := range []struct{ table, column string }{
		{"write_groups", "result_shape"},
		{"write_group_members", "calculation"},
	} {
		var count int
		if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM pragma_table_info(?) WHERE name=?", item.table, item.column).Scan(&count); err != nil {
			return fmt.Errorf("inspect member calculation column: %w", err)
		}
		if count == 0 {
			if _, err := tx.ExecContext(ctx, statements[i]); err != nil {
				return fmt.Errorf("expand member calculation column: %w", err)
			}
		}
	}
	return nil
}
