package datalink

import (
	"database/sql"
	"fmt"
)

// ensureSQLiteWriteGroupDeliveryTables creates the durable write-group
// delivery tables once; every table is CREATE IF NOT EXISTS.
func ensureSQLiteWriteGroupDeliveryTables(db *sql.DB) error {
	const migrationName = "025_write_group_delivery_sqlite.up.sql"

	for _, table := range []string{
		"wg_delivery_samples", "wg_delivery_checkpoints", "wg_delivery_buckets",
		"wg_delivery_outbox", "wg_delivery_receipts",
	} {
		exists, err := sqliteTableExists(db, table)
		if err != nil {
			return fmt.Errorf("failed to inspect sqlite table %s for migration %s: %w", table, migrationName, err)
		}
		if !exists {
			return applySQLiteMigrationOnce(db, migrationName)
		}
	}
	return nil
}
