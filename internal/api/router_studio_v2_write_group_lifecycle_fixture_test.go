package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type lifecycleAPIInspector struct{ calls int }

func (i *lifecycleAPIInspector) InspectTable(context.Context, string, string, string) (*dbtarget.TableInspection, error) {
	i.calls++
	return nil, nil
}

type lifecycleAPIGuard struct {
	calls int
	err   error
}

func (g *lifecycleAPIGuard) CheckWriteGroupBacklog(_ context.Context, tx *sql.Tx, group *workspace.WriteGroup) error {
	g.calls++
	var id string
	if err := tx.QueryRow(`SELECT id FROM write_groups WHERE id=? AND revision=?`, group.ID, group.Revision).Scan(&id); err != nil {
		return err
	}
	return g.err
}

func lifecycleAPIRead(t *testing.T, f writeGroupRouterFixture, id string) map[string]any {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id+"/lifecycle", nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	body := decodeJSONBody(t, response)
	require.Equal(t, true, body["success"])
	return body["data"].(map[string]any)
}

func lifecycleAPIActions(t *testing.T, data map[string]any) map[string]map[string]any {
	t.Helper()
	entries := data["actions"].([]any)
	require.Len(t, entries, 5)
	result := map[string]map[string]any{}
	for i, name := range []string{"edit", "stop", "remove", "recreate", "transition"} {
		action := entries[i].(map[string]any)
		require.Equal(t, name, action["action"])
		for _, field := range []string{"action", "allowed", "reason_code", "next_action", "requires_review"} {
			require.Contains(t, action, field)
		}
		result[name] = action
	}
	return result
}

func lifecycleAPIState(t *testing.T, db *sql.DB) map[string][]string {
	t.Helper()
	state := map[string][]string{}
	for _, table := range []string{"system_settings", "database_connectors", "write_groups", "write_group_versions", "wg_runtime_versions", "wg_delivery_samples", "wg_delivery_outbox", "wg_delivery_receipts"} {
		rows, err := db.Query(`SELECT * FROM ` + table + ` ORDER BY rowid`)
		require.NoError(t, err)
		cols, err := rows.Columns()
		require.NoError(t, err)
		for rows.Next() {
			values := make([]any, len(cols))
			dest := make([]any, len(cols))
			for i := range values {
				dest[i] = &values[i]
			}
			require.NoError(t, rows.Scan(dest...))
			raw, err := json.Marshal(values)
			require.NoError(t, err)
			state[table] = append(state[table], string(raw))
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
	}
	return state
}

func lifecycleAPIApplied(t *testing.T) managedGroupFixture {
	t.Helper()
	f := newManagedGroupFixture(t, "sqlite")
	token := f.preview(t)
	request := f.schemaRequest()
	request["token"], request["operation_id"] = token.Token, token.OperationID
	response := performJSONRequest(t, f.router, http.MethodPost, f.confirmPath(), request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	id := f.group["id"].(string)
	response = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	f.data = getGroupData(t, f, id)
	f.group = f.data["group"].(map[string]any)
	return f
}
