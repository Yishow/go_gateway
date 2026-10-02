package datalink

import (
	"database/sql"
	"fmt"
)

// ensureSQLiteOperationTestWriteColumns adds the test-write result columns to
// the existing schema operation ledger without touching its rows.
func ensureSQLiteOperationTestWriteColumns(db *sql.DB) error {
	const migrationName = "026_operation_test_write_sqlite.up.sql"

	exists, err := sqliteColumnExists(db, "managed_schema_operations", "write_outcome")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite operation ledger for migration %s: %w", migrationName, err)
	}
	if exists {
		return nil
	}
	return applySQLiteMigrationOnce(db, migrationName)
}
