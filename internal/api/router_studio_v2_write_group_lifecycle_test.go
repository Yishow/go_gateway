package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRouter_GroupLifecycleDisableUsesCASAndPreservesMembership(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, created))
	id := group["id"].(string)
	request := groupEditRequest(f, data, group)
	delete(request, "group")
	disabled := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/disable", request)
	require.Equal(t, http.StatusOK, disabled.Code, disabled.Body.String())
	disabledData, disabledGroup := groupSaveData(t, decodeJSONBody(t, disabled))
	require.Equal(t, "disabled", disabledGroup["status"])
	require.Equal(t, id, disabledGroup["id"])
	require.Equal(t, group["members"], disabledGroup["members"])
	require.Equal(t, group["applied_revision"], disabledGroup["applied_revision"])
	require.NotEqual(t, group["revision"], disabledGroup["revision"])
	require.NotEqual(t, data["workspace_revision"], disabledData["workspace_revision"])
	reloaded := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	require.Equal(t, decodeJSONBody(t, disabled), decodeJSONBody(t, reloaded))
	stale := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/disable", request)
	require.Equal(t, http.StatusConflict, stale.Code, stale.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, stale), "revision_mismatch")
	unknown := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/unknown/disable", groupEditRequest(f, disabledData, disabledGroup))
	require.Equal(t, http.StatusNotFound, unknown.Code, unknown.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, unknown), "WRITE_GROUP_NOT_FOUND")
}

func TestNewRouter_GroupLifecycleDisableValidatesRequestAndCannotReviveTombstone(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code)
	data, group := groupSaveData(t, decodeJSONBody(t, created))
	id := group["id"].(string)
	for _, field := range []string{"workspace_id", "expected_workspace_revision", "expected_group_revision", "expected_connector_revision"} {
		request := groupEditRequest(f, data, group)
		delete(request, "group")
		delete(request, field)
		response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/disable", request)
		require.Equal(t, http.StatusBadRequest, response.Code, field)
		assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_INVALID_REQUEST")
	}
	deleted := performJSONRequest(t, f.router, http.MethodDelete, writeGroupsPath+"/"+id, groupEditRequest(f, data, group))
	require.Equal(t, http.StatusOK, deleted.Code, deleted.Body.String())
	deletedData, tombstone := groupSaveData(t, decodeJSONBody(t, deleted))
	disabled := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/disable", groupEditRequest(f, deletedData, tombstone))
	require.Equal(t, http.StatusConflict, disabled.Code, disabled.Body.String())
	body := decodeJSONBody(t, disabled)
	assertWriteGroupError(t, body, "WRITE_GROUP_LIFECYCLE_BLOCKED")
	require.Equal(t, "reload", body["error"].(map[string]any)["action"])
	reloaded := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	require.Equal(t, decodeJSONBody(t, deleted), decodeJSONBody(t, reloaded))
}

func TestNewRouter_GroupApplyRequiresReadinessAndCurrentRevisions(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, created))
	id := group["id"].(string)

	for _, field := range []string{"workspace_id", "expected_workspace_revision", "expected_group_revision", "expected_connector_revision"} {
		request := groupEditRequest(f, data, group)
		delete(request, "group")
		delete(request, field)
		response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", request)
		require.Equal(t, http.StatusBadRequest, response.Code, field)
		assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_INVALID_REQUEST")
	}

	// A group payload is never accepted: Apply schedules what is already saved.
	withPayload := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", groupEditRequest(f, data, group))
	require.Equal(t, http.StatusUnprocessableEntity, withPayload.Code, withPayload.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, withPayload), "WRITE_GROUP_INVALID")

	// The fixture has no readable destination table, so the draft is not ready:
	// nothing is applied and the group keeps its revision and applied state.
	request := groupEditRequest(f, data, group)
	delete(request, "group")
	notReady := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", request)
	require.Equal(t, http.StatusConflict, notReady.Code, notReady.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, notReady), "WRITE_GROUP_NOT_READY")
	reloaded := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	_, after := groupSaveData(t, decodeJSONBody(t, reloaded))
	require.Equal(t, group["revision"], after["revision"])
	require.Equal(t, group["applied_revision"], after["applied_revision"])

	unknown := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/unknown/apply", request)
	require.Equal(t, http.StatusNotFound, unknown.Code, unknown.Body.String())
}
