package api

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func postgresSuccessorRuntimeValues(t *testing.T, db *sql.DB, g *workspace.WriteGroup) []string {
	t.Helper()
	column := `COALESCE(CAST(value_decimal AS TEXT),CAST(value_float AS TEXT),CAST(value_int AS TEXT))`
	if g.EffectiveLayoutVersion() == workspace.WriteGroupLayoutManagedWideProfileV1 {
		column = `CAST("` + g.Members[0].TargetColumn + `" AS TEXT)`
	}
	rows, err := db.Query(`SELECT ` + column + ` FROM ` + qualified("postgres", g.Destination.TableSchema, g.Destination.TableName) + ` WHERE ` + column + ` IS NOT NULL ORDER BY bucket_start`)
	require.NoError(t, err)
	defer rows.Close()
	values := []string{}
	for rows.Next() {
		var value string
		require.NoError(t, rows.Scan(&value))
		values = append(values, plainDecimal(value))
	}
	require.NoError(t, rows.Err())
	return values
}

func postgresSuccessorRuntimeRows(t *testing.T, db *sql.DB, g *workspace.WriteGroup) map[string]string {
	t.Helper()
	return postgresSuccessorRuntimeEvidenceRows(t, db, g.Destination.TableSchema, g.Destination.TableName, "bucket_start")
}

func postgresSuccessorRuntimeEvidenceRows(t *testing.T, db *sql.DB, namespace, table, identity string) map[string]string {
	t.Helper()
	rows, err := db.Query(`SELECT * FROM ` + qualified("postgres", namespace, table) + ` ORDER BY "` + identity + `"`)
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
	result := map[string]string{}
	for rows.Next() {
		values := make([]any, len(columns))
		dest := make([]any, len(values))
		for i := range values {
			dest[i] = &values[i]
		}
		require.NoError(t, rows.Scan(dest...))
		raw, err := json.Marshal(values)
		require.NoError(t, err)
		// The stored bucket_start is immutable identity, not the mutable row count.
		for i, name := range columns {
			if name == identity {
				key, err := json.Marshal(values[i])
				require.NoError(t, err)
				result[string(key)] = string(raw)
			}
		}
	}
	require.NoError(t, rows.Err())
	return result
}

var postgresSuccessorRuntimeTransitions = []struct {
	name                              string
	sourceLong, targetLong, timescale bool
}{
	{"plain-wide-to-long", false, true, false},
	{"plain-long-to-wide", true, false, false},
	{"plain-wide-A-to-B", false, false, false},
	{"plain-long-A-to-B", true, true, false},
	{"timescale-long-A-to-B", true, true, true},
}

func postgresSuccessorRuntimeCounterRows(t *testing.T, r *postgresSuccessorRig, g *workspace.WriteGroup) []counterWindow {
	t.Helper()
	from := qualified("postgres", g.Destination.TableSchema, g.Destination.TableName)
	if g.EffectiveLayoutVersion() == workspace.WriteGroupLayoutManagedLongV1 {
		return readLongCounterWindows(t, r.target, from)
	}
	return readWideCounterWindows(t, r.target, from, g.Members[0].TargetColumn)
}
