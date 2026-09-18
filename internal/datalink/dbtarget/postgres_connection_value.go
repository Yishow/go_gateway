package dbtarget

import "strings"

// quotePostgresConnectionValue preserves literal bytes in libpq keyword/value
// connection strings, including whitespace, quotes and backslashes.
func quotePostgresConnectionValue(value string) string {
	return "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(value) + "'"
}
