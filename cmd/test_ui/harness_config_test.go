package main

import (
	"strings"
	"testing"

	"go-gateway/internal/datalink"
)

func TestEmbeddedSQLiteDSNUsesDefaultWhenHarnessPathIsUnset(t *testing.T) {
	t.Setenv("GATEWAY_DB_PATH", "")
	t.Setenv("DB_PATH", "")
	t.Setenv("SQLITE_PATH", "")
	if got := embeddedSQLiteDSN(); got != datalink.DefaultEmbeddedSQLiteDSN {
		t.Fatalf("expected default embedded DSN, got %q", got)
	}
}

func TestEmbeddedSQLiteDSNUsesConfiguredHarnessPath(t *testing.T) {
	t.Setenv("GATEWAY_DB_PATH", "/tmp/gateway-harness.db")
	t.Setenv("SQLITE_PATH", "/tmp/ignored.db")
	got := embeddedSQLiteDSN()
	if !strings.Contains(got, "/tmp/gateway-harness.db") || strings.Contains(got, "ignored.db") {
		t.Fatalf("expected preferred configured path, got %q", got)
	}
}

func TestEmbeddedSQLiteDSNUsesDBPathCompatibilityAlias(t *testing.T) {
	t.Setenv("GATEWAY_DB_PATH", "")
	t.Setenv("DB_PATH", "/tmp/gateway-db-path.db")
	t.Setenv("SQLITE_PATH", "/tmp/ignored.db")
	got := embeddedSQLiteDSN()
	if !strings.Contains(got, "/tmp/gateway-db-path.db") || strings.Contains(got, "ignored.db") {
		t.Fatalf("expected DB_PATH alias, got %q", got)
	}
}
