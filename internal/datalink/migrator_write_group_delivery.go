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
			if err := applySQLiteMigrationOnce(db, migrationName); err != nil {
				return err
			}
			break
		}
	}
	if err := ensureSQLiteCounterBaselineTable(db); err != nil {
		return err
	}
	return ensureSQLiteCounterProofAssociations(db)
}

// ensureSQLiteCounterBaselineTable adds the durable counter baseline table
// beside the delivery tables of an installation that already has them.
func ensureSQLiteCounterBaselineTable(db *sql.DB) error {
	const migrationName = "033_write_group_counter_baselines_sqlite.up.sql"
	exists, err := sqliteTableExists(db, "wg_delivery_counter_baselines")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table wg_delivery_counter_baselines for migration %s: %w", migrationName, err)
	}
	if exists {
		return nil
	}
	return applySQLiteMigrationOnce(db, migrationName)
}

// ensureSQLiteCounterProofAssociations adds audit metadata without rewriting history.
func ensureSQLiteCounterProofAssociations(db *sql.DB) error {
	exists, err := sqliteTableExists(db, "wg_counter_proof_associations")
	if err != nil {
		return fmt.Errorf("inspect counter proof associations: %w", err)
	}
	if exists {
		return nil
	}
	return applySQLiteMigrationOnce(db, "043_counter_proof_associations_sqlite.up.sql")
}
