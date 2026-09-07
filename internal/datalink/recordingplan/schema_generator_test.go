package recordingplan

import (
	"strings"
	"testing"
)

func TestGenerateManagedSchemaDDL_SQLite(t *testing.T) {
	stmts, err := GenerateManagedSchemaDDL("sqlite", "gw_record_")
	if err != nil {
		t.Fatalf("GenerateManagedSchemaDDL sqlite failed: %v", err)
	}

	if len(stmts) == 0 {
		t.Fatalf("expected non-empty statements for sqlite")
	}

	joined := strings.Join(stmts, "\n")
	if !strings.Contains(joined, "gw_record_samples") {
		t.Errorf("expected gw_record_samples table in ddl")
	}
	if !strings.Contains(joined, "gw_record_intervals") {
		t.Errorf("expected gw_record_intervals table in ddl")
	}
	if !strings.Contains(joined, "gw_record_receipts") {
		t.Errorf("expected gw_record_receipts table in ddl")
	}
}

func TestGenerateManagedSchemaDDL_PostgreSQL(t *testing.T) {
	stmts, err := GenerateManagedSchemaDDL("postgres", "gw_record_")
	if err != nil {
		t.Fatalf("GenerateManagedSchemaDDL postgres failed: %v", err)
	}

	joined := strings.Join(stmts, "\n")
	if !strings.Contains(joined, "TIMESTAMP WITH TIME ZONE") {
		t.Errorf("expected TIMESTAMP WITH TIME ZONE in postgres ddl")
	}
	if !strings.Contains(joined, "CREATE INDEX IF NOT EXISTS") {
		t.Errorf("expected index creation in postgres ddl")
	}
}

func TestGenerateManagedSchemaDDL_MySQL(t *testing.T) {
	_, err := GenerateManagedSchemaDDL("mysql", "gw_record_")
	if err == nil {
		t.Fatalf("expected unsupported dialect error for mysql until capability suite is verified")
	}
}
