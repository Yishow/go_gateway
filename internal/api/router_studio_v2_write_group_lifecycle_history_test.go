package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupLifecycleHistoryIncludesImmutableRevision(t *testing.T) {
	f := lifecycleAPIApplied(t)
	id := f.group["id"].(string)
	before := lifecycleAPIState(t, f.db)
	data := lifecycleAPIRead(t, f.writeGroupRouterFixture, id)
	require.Contains(t, data, "history", "lifecycle must expose immutable history, not only current actions")
	history := data["history"].(map[string]any)
	require.Equal(t, false, history["canonical"])
	require.Equal(t, id, history["current_group_id"])
	generations := history["generations"].([]any)
	require.Len(t, generations, 1)
	generation := generations[0].(map[string]any)
	require.Equal(t, id, generation["group_id"])
	require.Equal(t, true, generation["current"])
	require.Equal(t, false, generation["reserved"])
	revisions := generation["revisions"].([]any)
	require.Len(t, revisions, 1)
	revision := revisions[0].(map[string]any)
	require.Equal(t, f.group["applied_revision"], revision["group_revision"])
	require.Equal(t, []any{map[string]any{"device_id": "device-A", "point_id": "point-A", "tag_id": "tag-A"}}, revision["sources"])
	require.Equal(t, before, lifecycleAPIState(t, f.db))
}

func TestWriteGroupLifecycleHistoryKeepsFrozenDestinationAfterDraftEdit(t *testing.T) {
	f := lifecycleAPIApplied(t)
	id := f.group["id"].(string)
	applied := f.group["applied_revision"]
	changed := cloneMap(f.group)
	changed["row_policy"].(map[string]any)["interval_seconds"] = 30
	response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, groupEditRequest(f.writeGroupRouterFixture, f.data, changed))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, response))
	_, err := f.db.Exec(`UPDATE database_connectors SET identity_revision='later-private-connector',connection_config='{"dsn":"private-endpoint","password":"private-secret"}' WHERE id='connector-A'`)
	require.NoError(t, err)
	before := lifecycleAPIState(t, f.db)
	data := lifecycleAPIRead(t, f.writeGroupRouterFixture, id)
	require.Contains(t, data, "history")
	generation := data["history"].(map[string]any)["generations"].([]any)[0].(map[string]any)
	require.Equal(t, saved["revision"], generation["group_revision"])
	revisions := generation["revisions"].([]any)
	require.Len(t, revisions, 1)
	revision := revisions[0].(map[string]any)
	require.Equal(t, applied, revision["group_revision"])
	destination := revision["destination"].(map[string]any)
	require.Equal(t, "connector-A", destination["connector_id"])
	require.Equal(t, "connector-1", destination["connector_revision"])
	require.Len(t, destination["destination_scope"].(string), 64)
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	for _, private := range []string{f.file, "private-endpoint", "private-secret", "later-private-connector", "connection_config", "last_error_code", "database_name"} {
		require.NotContains(t, string(raw), private)
	}
	require.Equal(t, before, lifecycleAPIState(t, f.db))
}

func TestWriteGroupLifecycleHistoryRemovedGroupRemainsReadable(t *testing.T) {
	f := lifecycleAPIApplied(t)
	id := f.group["id"].(string)
	f.groups.WithBacklogOwnershipGuard(&lifecycleAPIGuard{})
	removed := performJSONRequest(t, f.router, http.MethodDelete, writeGroupsPath+"/"+id, groupEditRequest(f.writeGroupRouterFixture, f.data, f.group))
	require.Equal(t, http.StatusOK, removed.Code, removed.Body.String())
	before := lifecycleAPIState(t, f.db)
	data := lifecycleAPIRead(t, f.writeGroupRouterFixture, id)
	require.Equal(t, "deleted", data["status"])
	require.Contains(t, data, "history")
	generation := data["history"].(map[string]any)["generations"].([]any)[0].(map[string]any)
	require.Equal(t, "deleted", generation["status"])
	require.Len(t, generation["revisions"].([]any), 1)
	require.Equal(t, false, lifecycleAPIActions(t, data)["recreate"]["allowed"])
	require.Equal(t, before, lifecycleAPIState(t, f.db))
}

func TestWriteGroupLifecycleHistoryCorruptEvidenceFailsSafely(t *testing.T) {
	for _, mode := range []string{"invalid-payload", "foreign-workspace", "missing-version"} {
		t.Run(mode, func(t *testing.T) {
			f := lifecycleAPIApplied(t)
			id := f.group["id"].(string)
			query := `UPDATE write_group_versions SET payload='private-corrupt-payload' WHERE group_id=?`
			if mode == "foreign-workspace" {
				query = `UPDATE write_group_versions SET payload=json_set(payload,'$.workspace_id','private-foreign') WHERE group_id=?`
			} else if mode == "missing-version" {
				query = `DELETE FROM write_group_versions WHERE group_id=?`
			}
			_, err := f.db.Exec(query, id)
			require.NoError(t, err)
			before := lifecycleAPIState(t, f.db)
			response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id+"/lifecycle", nil)
			require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
			require.NotContains(t, response.Body.String(), "private-")
			require.NotContains(t, response.Body.String(), `"history"`)
			require.Equal(t, before, lifecycleAPIState(t, f.db))
		})
	}
}
