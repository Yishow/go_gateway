package api

import (
	"net/http"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestNewRouter_WriteGroupSafeResourcesAndRevisionErrors(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, created))
	id := group["id"].(string)
	before, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `INSERT INTO write_groups (
		id, workspace_id, revision, name, destination_connector_id, destination_connector_revision, destination_table_name
	) VALUES ('foreign-group','foreign-workspace','foreign-revision','Foreign','connector-A','connector-1','private_table')`)
	require.NoError(t, err)

	for _, resource := range []string{"unknown-group", "foreign-group"} {
		readiness := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+resource+"/readiness", nil)
		require.Equal(t, http.StatusNotFound, readiness.Code, readiness.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, readiness), "WRITE_GROUP_NOT_FOUND")
		for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
			t.Run(method+"/"+resource, func(t *testing.T) {
				var body any
				if method != http.MethodGet {
					body = groupEditRequest(f, data, group)
				}
				response := performJSONRequest(t, f.router, method, writeGroupsPath+"/"+resource, body)
				require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
				assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_NOT_FOUND")
				require.NotContains(t, response.Body.String(), "private_table")
				require.NotContains(t, response.Body.String(), "foreign-workspace")
			})
		}
	}

	for _, field := range []string{"expected_workspace_revision", "expected_group_revision", "expected_connector_revision"} {
		t.Run("stale/"+field, func(t *testing.T) {
			request := groupEditRequest(f, data, group)
			request[field] = "stale"
			response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, request)
			require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
			assertWriteGroupError(t, decodeJSONBody(t, response), "revision_mismatch")
		})
	}
	foreign := groupEditRequest(f, data, group)
	foreign["workspace_id"] = "foreign-workspace"
	response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, foreign)
	require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_NOT_FOUND")
	after, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, before, after)
	reloaded := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	gotData, gotGroup := groupSaveData(t, decodeJSONBody(t, reloaded))
	require.Equal(t, data, gotData)
	require.Equal(t, group, gotGroup)
}

func TestNewRouter_WriteGroupRequestValidationAndTombstone(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	for _, field := range []string{"workspace_id", "expected_workspace_revision", "expected_connector_revision", "group"} {
		t.Run("missing/"+field, func(t *testing.T) {
			request := f.createRequest()
			delete(request, field)
			response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
			require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
			assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_INVALID_REQUEST")
		})
	}
	invalid := f.createRequest()
	invalid["group"].(map[string]any)["name"] = ""
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, invalid)
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_INVALID")
	foreign := f.createRequest()
	foreign["group"].(map[string]any)["workspace_id"] = "foreign-workspace"
	response = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, foreign)
	require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())

	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, created))
	id := group["id"].(string)
	listed := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
	require.Equal(t, http.StatusOK, listed.Code, listed.Body.String())
	listData := decodeJSONBody(t, listed)["data"].(map[string]any)
	require.Equal(t, f.record.ID, listData["workspace_id"])
	require.Equal(t, data["workspace_revision"], listData["workspace_revision"])
	require.Len(t, listData["groups"], 1)

	missingRevision := groupEditRequest(f, data, group)
	delete(missingRevision, "expected_group_revision")
	response = performJSONRequest(t, f.router, http.MethodDelete, writeGroupsPath+"/"+id, missingRevision)
	require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
	deleted := performJSONRequest(t, f.router, http.MethodDelete, writeGroupsPath+"/"+id, groupEditRequest(f, data, group))
	require.Equal(t, http.StatusOK, deleted.Code, deleted.Body.String())
	deletedData, tombstone := groupSaveData(t, decodeJSONBody(t, deleted))
	require.Equal(t, string(workspace.WriteGroupStatusDeleted), tombstone["status"])
	require.Equal(t, id, tombstone["id"])
	read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	gotData, gotGroup := groupSaveData(t, decodeJSONBody(t, read))
	require.Equal(t, deletedData, gotData)
	require.Equal(t, tombstone, gotGroup)
	var members int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_group_members WHERE group_id = ?`, id).Scan(&members))
	require.Equal(t, 1, members)
}

func TestNewRouter_WriteGroupMissingServiceFailsClosed(t *testing.T) {
	router := NewRouter(&DatalinkServices{Workspace: workspace.NewService(workspace.NewMemoryRepository())})
	response := performJSONRequest(t, router, http.MethodGet, writeGroupsPath, nil)
	require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_UNAVAILABLE")
}
