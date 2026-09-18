package datalink

import (
	"database/sql"
	"fmt"
)

func ensureSQLiteTelemetryRecordingMigrations(db *sql.DB) error {
	if err := ensureSQLiteMeasurementSemanticsTable(db); err != nil {
		return err
	}
	if err := ensureSQLiteSourceRuleMixedLayoutColumns(db); err != nil {
		return err
	}
	if err := ensureSQLiteRecordingPlansTable(db); err != nil {
		return err
	}
	if err := ensureSQLiteSchemaPreviewScopeColumns(db); err != nil {
		return err
	}
	if err := ensureSQLiteSchemaOperationsTable(db); err != nil {
		return err
	}
	return ensureSQLiteDurableDeliveryTables(db)
}

// ensureSQLiteSchemaOperationsTable creates the durable schema operation ledger.
func ensureSQLiteSchemaOperationsTable(db *sql.DB) error {
	const migrationName = "023_schema_operations_sqlite.up.sql"

	exists, err := sqliteTableExists(db, "managed_schema_operations")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table managed_schema_operations: %w", err)
	}
	if exists {
		return nil
	}

	return applySQLiteMigrationOnce(db, migrationName)
}

// ensureSQLiteSchemaPreviewScopeColumns adds preview protection columns with
// empty defaults; existing tokens stay legacy and require a new preview.
func ensureSQLiteSchemaPreviewScopeColumns(db *sql.DB) error {
	const migrationName = "022_schema_preview_scope_sqlite.up.sql"

	exists, err := sqliteColumnExists(db, "managed_schema_preview_tokens", "digest")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite preview token digest column: %w", err)
	}
	if exists {
		return nil
	}

	return applySQLiteMigrationOnce(db, migrationName)
}

func ensureSQLiteDurableDeliveryTables(db *sql.DB) error {
	const migrationName = "020_durable_delivery_sqlite.up.sql"

	exists, err := sqliteTableExists(db, "gw_delivery_outbox")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table gw_delivery_outbox: %w", err)
	}
	if exists {
		return nil
	}

	return applySQLiteMigrationOnce(db, migrationName)
}

func ensureSQLiteRecordingPlansTable(db *sql.DB) error {
	const migrationName = "019_recording_plans_sqlite.up.sql"

	exists, err := sqliteTableExists(db, "recording_plans")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table recording_plans for migration %s: %w", migrationName, err)
	}
	if exists {
		return nil
	}

	return applySQLiteMigrationOnce(db, migrationName)
}
