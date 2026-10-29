package api

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func mysqlSuccessorRuntimeValues(t *testing.T, db *sql.DB, g *workspace.WriteGroup) []string {
	t.Helper()
	column := `COALESCE(CAST(value_decimal AS CHAR),CAST(value_float AS CHAR),CAST(value_int AS CHAR))`

	rows, err := db.Query(`SELECT ` + column + ` FROM ` + qualified("mysql", g.Destination.TableSchema, g.Destination.TableName) + ` WHERE ` + column + ` IS NOT NULL ORDER BY bucket_start`)
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

func mysqlSuccessorRuntimeRows(t *testing.T, db *sql.DB, g *workspace.WriteGroup) map[string]string {
	t.Helper()
	return mysqlSuccessorRuntimeEvidenceRows(t, db, g.Destination.TableSchema, g.Destination.TableName, "bucket_start")
}

func mysqlSuccessorRuntimeEvidenceRows(t *testing.T, db *sql.DB, namespace, table, identity string) map[string]string {
	t.Helper()
	rows, err := db.Query(`SELECT * FROM ` + qualified("mysql", namespace, table) + " ORDER BY `" + identity + "`")
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
				require.NotContains(t, result, string(key), "SQL readback must not hide duplicate logical identities")
				result[string(key)] = string(raw)
			}
		}
	}
	require.NoError(t, rows.Err())
	return result
}

var mysqlSuccessorRuntimeTransitions = []struct {
	name, family string
	mode         workspace.WriteGroupStorageMode
}{
	{"mysql/plain-A-to-B", "mysql", workspace.WriteGroupStorageModePlain},
	{"mysql/native-monthly-A-to-B", "mysql", workspace.WriteGroupStorageModeNativeMonthly},
	{"mariadb/plain-A-to-B", "mariadb", workspace.WriteGroupStorageModePlain},
	{"mariadb/native-monthly-A-to-B", "mariadb", workspace.WriteGroupStorageModeNativeMonthly},
}
