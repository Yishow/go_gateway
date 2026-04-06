package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"go-gateway/internal/datalink/schema/migrations"
)

func ensureSQLiteSourceRuleTagReviewDecisionsTable(db *sql.DB) error {
	const migrationName = "012_source_rule_tag_review_decision_sqlite.up.sql"

	exists, err := sqliteTableExists(db, "source_rule_tag_review_decisions")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table source_rule_tag_review_decisions for migration %s: %w", migrationName, err)
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
