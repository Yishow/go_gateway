package datalink

import (
	"context"
	"database/sql"
	"fmt"
)

func ensureSQLiteSourceRuleTagReviewDecisionStaleColumns(db *sql.DB) error {
	const migrationName = "013_source_rule_tag_review_decision_stale_sqlite"

	columns := []struct {
		name string
		ddl  string
	}{
		{
			name: "stale",
			ddl:  "ALTER TABLE source_rule_tag_review_decisions ADD COLUMN stale INTEGER NOT NULL DEFAULT 0",
		},
		{
			name: "stale_revision_id",
			ddl:  "ALTER TABLE source_rule_tag_review_decisions ADD COLUMN stale_revision_id TEXT NOT NULL DEFAULT ''",
		},
		{
			name: "stale_at",
			ddl:  "ALTER TABLE source_rule_tag_review_decisions ADD COLUMN stale_at TEXT",
		},
	}

	for _, column := range columns {
		exists, err := sqliteColumnExists(db, "source_rule_tag_review_decisions", column.name)
		if err != nil {
			return fmt.Errorf("failed to inspect sqlite column %s for migration %s: %w", column.name, migrationName, err)
		}
		if exists {
			continue
		}
		if _, err := db.ExecContext(context.Background(), column.ddl); err != nil {
			return fmt.Errorf("failed to execute migration %s for column %s: %w", migrationName, column.name, err)
		}
	}

	return nil
}
