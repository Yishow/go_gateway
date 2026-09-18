package datalink

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestMigrator_AddsSchemaPreviewScopeColumns(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "main.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	for run := 0; run < 2; run++ {
		if err := NewMigrator().Migrate(db); err != nil {
			t.Fatalf("migrate run %d: %v", run+1, err)
		}
	}

	for _, column := range []string{
		"operation_id", "action", "workspace_revision", "connector_revision", "dialect",
		"database_name", "schema_name", "no_change_reason", "tables", "digest",
	} {
		exists, err := sqliteColumnExists(db, "managed_schema_preview_tokens", column)
		if err != nil {
			t.Fatalf("inspect column %s: %v", column, err)
		}
		if !exists {
			t.Fatalf("managed_schema_preview_tokens.%s must exist after migration", column)
		}
	}
}
