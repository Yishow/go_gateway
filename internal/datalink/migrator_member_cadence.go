package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema/migrations"
)

func ensureSQLiteMemberCadence(db *sql.DB) error {
	return runSQLiteConfigurationExpansion(db, func(ctx context.Context, tx *sql.Tx) error {
		var count int
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM pragma_table_info('write_group_members') WHERE name='cadence'`).Scan(&count); err != nil {
			return err
		}
		raw, err := migrations.FS.ReadFile("050_member_cadence_sqlite.up.sql")
		if err != nil {
			return err
		}
		if count == 0 {
			if _, err := tx.ExecContext(ctx, string(raw)); err != nil {
				return fmt.Errorf("expand member cadence: %w", err)
			}
		}
		var ddl string
		if err := tx.QueryRowContext(ctx, `SELECT sql FROM sqlite_master WHERE type='table' AND name='write_group_members'`).Scan(&ddl); err != nil {
			return err
		}
		declaration := strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(string(raw), "ALTER TABLE write_group_members ADD COLUMN ")), ";")
		if strings.Count(ddl, declaration) != 1 {
			return fmt.Errorf("unknown member cadence column; migration refused")
		}
		if err := migrateMixedRecordProfileKeysTx(ctx, tx); err != nil {
			return err
		}
		return checkSQLiteMixedForeignKeys(ctx, tx)
	})
}
