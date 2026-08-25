package main

import (
	"os"
	"strings"

	"go-gateway/internal/datalink"
)

// embeddedSQLiteDSN resolves the EXE harness database path without changing
// the normal embedded default. GATEWAY_DB_PATH is preferred; DB_PATH and
// SQLITE_PATH are retained as compatibility aliases used by local scripts.
func embeddedSQLiteDSN() string {
	path := strings.TrimSpace(os.Getenv("GATEWAY_DB_PATH"))
	if path == "" {
		path = strings.TrimSpace(os.Getenv("DB_PATH"))
	}
	if path == "" {
		path = strings.TrimSpace(os.Getenv("SQLITE_PATH"))
	}
	if path == "" {
		return datalink.DefaultEmbeddedSQLiteDSN
	}
	return strings.Replace(datalink.DefaultEmbeddedSQLiteDSN, "datalink.db", path, 1)
}
