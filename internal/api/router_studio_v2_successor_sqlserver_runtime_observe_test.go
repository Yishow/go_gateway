package api

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

var sqlServerSuccessorRuntimeTransitions = []struct {
	name string
	mode workspace.WriteGroupStorageMode
}{
	{"plain-A-to-B", workspace.WriteGroupStorageModePlain},
	{"native-monthly-A-to-B", workspace.WriteGroupStorageModeNativeMonthly},
}

func sqlServerSuccessorRuntimeValues(t *testing.T, db *sql.DB, g *workspace.WriteGroup) []string {
	t.Helper()
	column := `COALESCE(CONVERT(varchar(100),value_decimal),CONVERT(varchar(100),value_float),CONVERT(varchar(100),value_int))`
	rows, err := db.QueryContext(t.Context(), `SELECT `+column+` FROM `+sqlServerQualified(g.Destination.TableSchema, g.Destination.TableName)+` WHERE `+column+` IS NOT NULL ORDER BY bucket_start`)
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

func sqlServerSuccessorRuntimeRows(t *testing.T, db *sql.DB, g *workspace.WriteGroup) map[string]string {
	t.Helper()
	return sqlServerSuccessorRuntimeEvidenceRows(t, db, g.Destination.TableSchema, g.Destination.TableName, "record_id")
}

func sqlServerSuccessorRuntimeEvidenceRows(t *testing.T, db *sql.DB, namespace, table, identity string) map[string]string {
	t.Helper()
	require.Contains(t, []string{"record_id", "effect_key"}, identity)
	rows, err := db.QueryContext(t.Context(), `SELECT * FROM `+sqlServerQualified(namespace, table)+` ORDER BY [`+identity+`]`)
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
	require.Contains(t, columns, identity)
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
		for i, name := range columns {
			if name == identity {
				key, err := json.Marshal(values[i])
				require.NoError(t, err)
				require.NotContains(t, result, string(key), "readback must not hide duplicate immutable identities")
				result[string(key)] = string(raw)
			}
		}
	}
	require.NoError(t, rows.Err())
	return result
}
