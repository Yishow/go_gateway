package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func ensureSQLiteDatabaseDeliveryOutcomeColumns(db *sql.DB) error {
	const migrationName = "015_database_delivery_outcomes_sqlite.up.sql"

	columns := []struct {
		name string
		ddl  string
	}{
		{name: "last_schema_ensure_at", ddl: "ALTER TABLE database_connectors ADD COLUMN last_schema_ensure_at DATETIME"},
		{name: "last_schema_ensure_status", ddl: "ALTER TABLE database_connectors ADD COLUMN last_schema_ensure_status TEXT NOT NULL DEFAULT ''"},
		{name: "last_schema_ensure_error", ddl: "ALTER TABLE database_connectors ADD COLUMN last_schema_ensure_error TEXT NOT NULL DEFAULT ''"},
		{name: "last_write_at", ddl: "ALTER TABLE database_connectors ADD COLUMN last_write_at DATETIME"},
		{name: "last_write_status", ddl: "ALTER TABLE database_connectors ADD COLUMN last_write_status TEXT NOT NULL DEFAULT ''"},
		{name: "last_write_error", ddl: "ALTER TABLE database_connectors ADD COLUMN last_write_error TEXT NOT NULL DEFAULT ''"},
		{name: "last_flush_at", ddl: "ALTER TABLE database_connectors ADD COLUMN last_flush_at DATETIME"},
		{name: "last_flush_status", ddl: "ALTER TABLE database_connectors ADD COLUMN last_flush_status TEXT NOT NULL DEFAULT ''"},
		{name: "last_flush_error", ddl: "ALTER TABLE database_connectors ADD COLUMN last_flush_error TEXT NOT NULL DEFAULT ''"},
	}

	for _, column := range columns {
		exists, err := sqliteColumnExists(db, "database_connectors", column.name)
		if err != nil {
			return fmt.Errorf("failed to inspect sqlite column %s for migration %s: %w", column.name, migrationName, err)
		}
		if exists {
			continue
		}
		log.Printf("Executing SQLite migration: %s (database_connectors.%s)", migrationName, column.name)
		if _, err := db.ExecContext(context.Background(), column.ddl); err != nil {
			return fmt.Errorf("failed to execute migration %s for column %s: %w", migrationName, column.name, err)
		}
	}

	return nil
}
