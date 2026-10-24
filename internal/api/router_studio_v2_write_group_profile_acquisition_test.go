package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// dumpTable returns every row of a table as strings so that two snapshots can
// be compared without knowing the schema.
func dumpTable(t *testing.T, db *sql.DB, table string) []string {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), "SELECT * FROM "+table+" ORDER BY 1")
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
	var dump []string
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		require.NoError(t, rows.Scan(pointers...))
		dump = append(dump, fmt.Sprint(values...))
	}
	require.NoError(t, rows.Err())
	return dump
}

func TestStudioV2WriteGroupsProfileLeavesAcquisitionUntouched(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	seedProfileSource(t, f, "energy")
	tables := []string{"devices", "points", "tags", "mappings", "database_connectors"}
	before := map[string][]string{}
	for _, table := range tables {
		before[table] = dumpTable(t, f.db, table)
	}

	for _, request := range []map[string]any{
		profileRequest(t, f, "instantaneous", 1, "second", "fast", ""),
		profileRequest(t, f, "cumulative", 600, "minute", "fast", "energy"),
		profileRequest(t, f, "instantaneous", 301, "second", "slow", ""),
	} {
		performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	}

	for _, table := range tables {
		require.Equal(t, before[table], dumpTable(t, f.db, table),
			"saving a recording profile must not change %s: it is not an acquisition setting", table)
	}
}
