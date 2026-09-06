package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"go-gateway/internal/datalink/schema/migrations"
)

func ensureSQLiteMeasurementSemanticsTable(db *sql.DB) error {
	const migrationName = "018_measurement_semantics_sqlite.up.sql"

	exists, err := sqliteTableExists(db, "measurement_definitions")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table measurement_definitions for migration %s: %w", migrationName, err)
	}
	if exists {
		return nil
	}

	content, err := migrations.FS.ReadFile(migrationName)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", migrationName, err)
	}

	log.Printf("Executing SQLite migration: %s", migrationName)
	if _, err := db.ExecContext(context.Background(), string(content)); err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", migrationName, err)
	}

	return nil
}

func ensureSQLiteSourceRuleMixedLayoutColumns(db *sql.DB) error {
	const migrationName = "018_source_rule_mixed_layout_sqlite"

	columns := []struct {
		name string
		ddl  string
	}{
		{
			name: "layout_mode",
			ddl:  "ALTER TABLE source_rules ADD COLUMN layout_mode TEXT NOT NULL DEFAULT 'homogeneous'",
		},
		{
			name: "mixed_items",
			ddl:  "ALTER TABLE source_rules ADD COLUMN mixed_items TEXT",
		},
	}

	for _, column := range columns {
		exists, err := sqliteColumnExists(db, "source_rules", column.name)
		if err != nil {
			return fmt.Errorf("failed to inspect sqlite column %s for migration %s: %w", column.name, migrationName, err)
		}
		if exists {
			continue
		}
		log.Printf("Executing SQLite migration: %s (%s)", migrationName, column.name)
		if _, err := db.ExecContext(context.Background(), column.ddl); err != nil {
			return fmt.Errorf("failed to execute migration %s for column %s: %w", migrationName, column.name, err)
		}
	}

	return nil
}
