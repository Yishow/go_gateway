package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRouter_WriteGroupMissingSourcesAndConnectorDoNotPersist(t *testing.T) {
	for _, resource := range []string{"device", "point", "tag", "mapping", "connector"} {
		t.Run(resource, func(t *testing.T) {
			f := newWriteGroupRouterFixture(t)
			created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
			require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
			data, group := groupSaveData(t, decodeJSONBody(t, created))
			id := group["id"].(string)
			before, err := f.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)
			_, err = f.db.ExecContext(t.Context(), map[string]string{
				"device":    "DELETE FROM devices WHERE id = 'device-A'",
				"point":     "DELETE FROM points WHERE id = 'point-A'",
				"tag":       "DELETE FROM tags WHERE id = 'tag-A'",
				"mapping":   "DELETE FROM mappings WHERE id = 'mapping-A'",
				"connector": "DELETE FROM database_connectors WHERE id = 'connector-A'",
			}[resource])
			require.NoError(t, err)
			for _, method := range []string{http.MethodPost, http.MethodPut} {
				var path string
				request := groupEditRequest(f, data, group)
				if method == http.MethodPost {
					path = writeGroupsPath
					delete(request, "expected_group_revision")
				} else {
					path = writeGroupsPath + "/" + id
				}
				response := performJSONRequest(t, f.router, method, path, request)
				require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
				assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_NOT_FOUND")
			}
			after, err := f.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)
			require.Equal(t, before, after)
			var count int
			require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_groups WHERE workspace_id = ?`, f.record.ID).Scan(&count))
			require.Equal(t, 1, count)
			read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
			actualData, actualGroup := groupSaveData(t, decodeJSONBody(t, read))
			require.Equal(t, data, actualData)
			require.Equal(t, group, actualGroup)
		})
	}
}

func TestNewRouter_WriteGroupCreateAndDeleteRejectStaleConnector(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, created))
	id := group["id"].(string)
	before, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE database_connectors SET identity_revision = 'connector-2' WHERE id = 'connector-A'`)
	require.NoError(t, err)
	for _, method := range []string{http.MethodPost, http.MethodDelete} {
		request := groupEditRequest(f, data, group)
		path := writeGroupsPath
		if method == http.MethodDelete {
			path += "/" + id
		}
		response := performJSONRequest(t, f.router, method, path, request)
		require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, response), "revision_mismatch")
	}
	after, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, before, after)
	read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	actualData, actualGroup := groupSaveData(t, decodeJSONBody(t, read))
	require.Equal(t, data, actualData)
	require.Equal(t, group, actualGroup)
}

func TestNewRouter_WriteGroupUpdatePreservesServerFields(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, created))
	id := group["id"].(string)
	_, err := f.db.ExecContext(t.Context(), `UPDATE write_groups SET applied_revision = revision,
		migration = '{"source_kind":"fixture","source_ids":["legacy-1"],"adapter_version":"v1"}',
		destination_schema_revision = 'verified-schema', destination_schema_digest = 'verified-digest' WHERE id = ?`, id)
	require.NoError(t, err)
	before := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	data, stored := groupSaveData(t, decodeJSONBody(t, before))
	payload := decodeJSONBody(t, before)["data"].(map[string]any)["group"].(map[string]any)
	payload["id"] = "client-replacement"
	payload["revision"] = "client-revision"
	payload["applied_revision"] = "client-applied"
	payload["created_at"] = "2030-01-01T00:00:00Z"
	payload["status"] = "running"
	payload["name"] = "Renamed"
	payload["migration"] = map[string]any{"source_kind": "client", "source_ids": []string{"client"}}
	destination := payload["destination"].(map[string]any)
	destination["schema_revision"] = "client-proof"
	destination["schema_digest"] = "client-proof"
	request := groupEditRequest(f, data, payload)
	request["expected_group_revision"] = stored["revision"]
	response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	savedData, saved := groupSaveData(t, decodeJSONBody(t, response))
	require.Equal(t, id, saved["id"])
	require.Equal(t, stored["created_at"], saved["created_at"])
	require.Equal(t, stored["applied_revision"], saved["applied_revision"])
	require.Equal(t, stored["migration"], saved["migration"])
	require.Equal(t, "draft", saved["status"])
	require.Equal(t, "Renamed", saved["name"])
	require.NotEqual(t, stored["revision"], saved["revision"])
	savedDestination := saved["destination"].(map[string]any)
	require.Equal(t, "verified-schema", savedDestination["schema_revision"])
	require.Equal(t, "verified-digest", savedDestination["schema_digest"])
	read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	actualData, actualGroup := groupSaveData(t, decodeJSONBody(t, read))
	require.Equal(t, savedData, actualData)
	require.Equal(t, saved, actualGroup)
}
