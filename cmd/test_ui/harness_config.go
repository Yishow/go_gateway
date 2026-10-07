package main

import (
	"net/url"
	"os"
	"path/filepath"
	"runtime"
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
	return embeddedSQLiteDSNForPath(path)
}

// embeddedSQLiteDSNForPath opens exactly the literal file whose owner was locked.
func embeddedSQLiteDSNForPath(path string) string {
	_, query, _ := strings.Cut(datalink.DefaultEmbeddedSQLiteDSN, "?")
	literal := filepath.ToSlash(path)
	if runtime.GOOS == "windows" && filepath.IsAbs(path) && !strings.HasPrefix(literal, "/") {
		literal = "/" + literal
	}
	uri := url.URL{Scheme: "file", Path: literal, RawQuery: query, OmitHost: !filepath.IsAbs(path)}
	return uri.String()
}
