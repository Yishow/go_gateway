package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupLifecycleReadRoute(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, created))
	response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+group["id"].(string)+"/lifecycle", nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	body := decodeJSONBody(t, response)
	require.Equal(t, true, body["success"])
	data := body["data"].(map[string]any)
	require.Equal(t, group["revision"], data["group_revision"])
	require.Equal(t, "", data["applied_revision"])
	require.Equal(t, "connector-1", data["connector_revision"])
	require.Len(t, data["actions"], 5)
}

func TestWriteGroupLifecycleReadIsLocalScopedAndNonMutating(t *testing.T) {
	f := newManagedGroupFixture(t, "sqlite")
	id := f.group["id"].(string)
	inspector := &lifecycleAPIInspector{}
	f.groups.WithTableInspector(inspector)
	before := lifecycleAPIState(t, f.db)
	for range 2 {
		data := lifecycleAPIRead(t, f.writeGroupRouterFixture, id)
		require.Equal(t, f.record.ID, data["workspace_id"])
		require.Equal(t, f.data["workspace_revision"], data["workspace_revision"])
		actions := lifecycleAPIActions(t, data)
		for _, name := range []string{"edit", "stop", "remove"} {
			require.Equal(t, true, actions[name]["allowed"])
			require.Equal(t, "", actions[name]["reason_code"])
		}
		for _, name := range []string{"recreate", "transition"} {
			require.Equal(t, false, actions[name]["allowed"])
			require.Equal(t, "canonical_required", actions[name]["reason_code"])
			require.Equal(t, true, actions[name]["requires_review"])
		}
	}
	require.Equal(t, before, lifecycleAPIState(t, f.db))
	require.Zero(t, inspector.calls)
	require.NoFileExists(t, f.file)
	_, err := f.db.Exec(`INSERT INTO write_groups(id,workspace_id,revision,name,destination_connector_id,destination_connector_revision,destination_table_name) VALUES('foreign-lifecycle','private-workspace','private-revision','Private','connector-A','connector-1','private_table')`)
	require.NoError(t, err)
	for _, id := range []string{"missing", "foreign-lifecycle"} {
		response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id+"/lifecycle", nil)
		require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_NOT_FOUND")
		require.NotContains(t, response.Body.String(), "private-")
	}
}

func TestWriteGroupLifecycleStaleAndRemovedConnectorKeepStop(t *testing.T) {
	for _, mode := range []string{"changed", "removed"} {
		t.Run(mode, func(t *testing.T) {
			f := newManagedGroupFixture(t, "sqlite")
			id := f.group["id"].(string)
			var err error
			reason := "connector_identity_stale"
			if mode == "changed" {
				_, err = f.db.Exec(`UPDATE database_connectors SET identity_revision='connector-2',connection_config='{"password":"private-secret","dsn":"private-endpoint"}' WHERE id='connector-A'`)
			} else {
				_, err = f.db.Exec(`DELETE FROM database_connectors WHERE id='connector-A'`)
				reason = "connector_unavailable"
			}
			require.NoError(t, err)
			before := lifecycleAPIState(t, f.db)
			data := lifecycleAPIRead(t, f.writeGroupRouterFixture, id)
			actions := lifecycleAPIActions(t, data)
			require.Equal(t, "connector-1", data["connector_revision"])
			require.Equal(t, false, actions["remove"]["allowed"])
			require.Equal(t, reason, actions["remove"]["reason_code"])
			require.Equal(t, "stop", actions["remove"]["next_action"])
			require.Equal(t, true, actions["stop"]["allowed"])
			encoded, err := json.Marshal(data)
			require.NoError(t, err)
			require.NotContains(t, string(encoded), "private-")
			require.Equal(t, before, lifecycleAPIState(t, f.db))
			stopped := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/disable", groupEditRequest(f.writeGroupRouterFixture, f.data, f.group))
			require.Equal(t, http.StatusOK, stopped.Code, stopped.Body.String())
			after := lifecycleAPIActions(t, lifecycleAPIRead(t, f.writeGroupRouterFixture, id))
			require.Equal(t, false, after["stop"]["allowed"])
			require.Equal(t, "already_stopped", after["stop"]["reason_code"])
		})
	}
}

func TestWriteGroupLifecycleRemovedAndUnavailableService(t *testing.T) {
	f := newManagedGroupFixture(t, "sqlite")
	id := f.group["id"].(string)
	removed := performJSONRequest(t, f.router, http.MethodDelete, writeGroupsPath+"/"+id, groupEditRequest(f.writeGroupRouterFixture, f.data, f.group))
	require.Equal(t, http.StatusOK, removed.Code, removed.Body.String())
	data := lifecycleAPIRead(t, f.writeGroupRouterFixture, id)
	require.Equal(t, "deleted", data["status"])
	actions := lifecycleAPIActions(t, data)
	for _, name := range []string{"edit", "stop", "remove"} {
		require.Equal(t, false, actions[name]["allowed"])
		require.Equal(t, "recording_removed", actions[name]["reason_code"])
	}
	require.Equal(t, false, actions["recreate"]["allowed"])
	unavailable := performJSONRequest(t, NewRouter(&DatalinkServices{Workspace: f.workspace}), http.MethodGet, writeGroupsPath+"/"+id+"/lifecycle", nil)
	require.Equal(t, http.StatusServiceUnavailable, unavailable.Code, unavailable.Body.String())
}

func TestWriteGroupLifecycleUnknownStatusFailsSafely(t *testing.T) {
	f := newManagedGroupFixture(t, "sqlite")
	id := f.group["id"].(string)
	_, err := f.db.Exec(`UPDATE write_groups SET status='private-unknown-status' WHERE id=?`, id)
	require.NoError(t, err)
	before := lifecycleAPIState(t, f.db)
	response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id+"/lifecycle", nil)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_LIFECYCLE_BLOCKED")
	require.NotContains(t, response.Body.String(), "private-unknown-status")
	require.NotContains(t, response.Body.String(), `"actions"`)
	require.Equal(t, before, lifecycleAPIState(t, f.db))
}
