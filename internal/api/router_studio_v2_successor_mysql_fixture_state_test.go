package api

import (
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

// Include table DDL and every row in both owned databases, not merely counts.
func (r *mysqlSuccessorRig) targetState(t *testing.T) map[string][]string {
	t.Helper()
	state := map[string][]string{}
	for i, db := range []*sql.DB{r.f.target, r.targetDB} {
		scope := []string{r.f.namespace, r.targetNamespace}[i]
		rows, err := db.Query(`SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE() ORDER BY TABLE_NAME`)
		require.NoError(t, err)
		var names []string
		for rows.Next() {
			var name string
			require.NoError(t, rows.Scan(&name))
			names = append(names, name)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		for _, name := range names {
			require.Regexp(t, `^(gw_long_[0-9a-f]{16}|gw_effect_receipts)$`, name)
			var table, ddl string
			require.NoError(t, db.QueryRow("SHOW CREATE TABLE `"+name+"`").Scan(&table, &ddl))
			key := scope + "." + name
			state[key] = []string{ddl}
			rows, err := db.Query("SELECT * FROM `" + name + "`")
			require.NoError(t, err)
			columns, err := rows.Columns()
			require.NoError(t, err)
			var bodies []string
			for rows.Next() {
				values := make([]any, len(columns))
				refs := make([]any, len(columns))
				for j := range values {
					refs[j] = &values[j]
				}
				require.NoError(t, rows.Scan(refs...))
				raw, err := json.Marshal(values)
				require.NoError(t, err)
				bodies = append(bodies, string(raw))
			}
			require.NoError(t, rows.Err())
			require.NoError(t, rows.Close())
			sort.Strings(bodies)
			state[key] = append(state[key], bodies...)
		}
	}
	return state
}
