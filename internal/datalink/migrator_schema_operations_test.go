package datalink

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestMigrator_CreatesSchemaOperationLedger(t *testing.T) {
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

	exists, err := sqliteTableExists(db, "managed_schema_operations")
	if err != nil || !exists {
		t.Fatalf("managed_schema_operations must exist after migration: exists=%v err=%v", exists, err)
	}
	var indexes int
	if err := db.QueryRowContext(context.Background(), `
		SELECT COUNT(*) FROM sqlite_master
		WHERE type = 'index' AND name IN ('idx_managed_schema_operations_token', 'idx_managed_schema_operations_active_scope')
	`).Scan(&indexes); err != nil {
		t.Fatalf("inspect ledger indexes: %v", err)
	}
	if indexes != 2 {
		t.Fatalf("the ledger must keep unique token and active-scope indexes, found %d", indexes)
	}
}
