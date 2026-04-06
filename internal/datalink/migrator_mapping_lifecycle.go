package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func ensureSQLiteMappingLifecycleColumns(db *sql.DB) error {
	const migrationName = "011_mapping_lifecycle_sqlite.up.sql"

	tableExists, err := sqliteTableExists(db, "mappings")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite mappings table for migration %s: %w", migrationName, err)
	}
	if !tableExists {
		return nil
	}

	columns := []struct {
		name string
		ddl  string
	}{
		{name: "status", ddl: "ALTER TABLE mappings ADD COLUMN status TEXT NOT NULL DEFAULT 'active'"},
		{name: "rule_candidate_id", ddl: "ALTER TABLE mappings ADD COLUMN rule_candidate_id TEXT"},
		{name: "proposed_signature", ddl: "ALTER TABLE mappings ADD COLUMN proposed_signature TEXT"},
		{name: "last_applied_signature", ddl: "ALTER TABLE mappings ADD COLUMN last_applied_signature TEXT"},
		{name: "blocking_reason", ddl: "ALTER TABLE mappings ADD COLUMN blocking_reason TEXT"},
	}

	for _, column := range columns {
		exists, err := sqliteColumnExists(db, "mappings", column.name)
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

	if _, err := db.ExecContext(context.Background(), `
		UPDATE mappings
		SET status = CASE WHEN enabled = 1 THEN 'active' ELSE 'draft' END
		WHERE status IS NULL OR status = '' OR (enabled = 0 AND status = 'active')
	`); err != nil {
		return fmt.Errorf("failed to normalize sqlite mapping lifecycle states for migration %s: %w", migrationName, err)
	}

	if _, err := db.ExecContext(context.Background(), `
		CREATE INDEX IF NOT EXISTS idx_mappings_rule_candidate_id ON mappings(rule_candidate_id)
	`); err != nil {
		return fmt.Errorf("failed to create sqlite mapping lifecycle index for migration %s: %w", migrationName, err)
	}

	return nil
}
