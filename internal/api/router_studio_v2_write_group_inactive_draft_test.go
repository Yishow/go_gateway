package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRouter_WriteGroupKnownInactiveSourceDraftRename(t *testing.T) {
	f := newInactiveDraftFixture(t)
	f.stopSource(t)
	f.group["name"] = "Known inactive draft renamed"
	response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+f.group["id"].(string), f.editRequest(t))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, response))
	require.Equal(t, "Known inactive draft renamed", group["name"])
	require.Equal(t, "draft", group["status"])
	require.Empty(t, group["applied_revision"])
	require.NotEqual(t, f.group["revision"], group["revision"])
	require.NotEqual(t, f.data["workspace_revision"], data["workspace_revision"])
	id := group["id"].(string)
	reloaded := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	_, persisted := groupSaveData(t, decodeJSONBody(t, reloaded))
	require.Equal(t, group, persisted)
	readiness := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id+"/readiness", nil)
	require.Equal(t, http.StatusOK, readiness.Code, readiness.Body.String())
	ready := decodeJSONBody(t, readiness)["data"].(map[string]any)
	require.Equal(t, false, ready["ready"])
	require.Equal(t, false, ready["config_ready"])
	before := inactiveDraftState(t, f.db)
	applyRequest := groupEditRequest(f.writeGroupRouterFixture, data, group)
	delete(applyRequest, "group")
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", applyRequest)
	require.Equal(t, http.StatusConflict, applied.Code, applied.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, applied), "WRITE_GROUP_NOT_READY")
	require.Equal(t, before, inactiveDraftState(t, f.db))
	// A second edit binds to the saved inactive hash, not the old active hash.
	f.data, f.group = data, group
	f.group["name"] = "Known inactive draft renamed again"
	repeated := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, f.editRequest(t))
	require.Equal(t, http.StatusOK, repeated.Code, repeated.Body.String())
}

func TestNewRouter_WriteGroupInactiveDraftRejectsChangedMapping(t *testing.T) {
	for _, mutation := range []string{
		`DELETE FROM mappings WHERE id=?`,
		`UPDATE mappings SET id='replacement-mapping' WHERE id=?`,
		`UPDATE mappings SET tag_id='tag-A' WHERE id=?`,
		`UPDATE mappings SET point_id='point-A' WHERE id=?`,
		`UPDATE mappings SET transform_pipeline='[{"type":"scale","params":{"scale":2}}]' WHERE id=?`,
		`UPDATE mappings SET proposed_signature='changed' WHERE id=?`,
		`UPDATE mappings SET last_applied_signature='changed' WHERE id=?`,
		`UPDATE mappings SET rule_candidate_id='changed' WHERE id=?`,
		`UPDATE mappings SET status='out_of_sync',blocking_reason='needs confirmation' WHERE id=?`,
	} {
		t.Run(mutation, func(t *testing.T) {
			f := newInactiveDraftFixture(t)
			f.stopSource(t)
			_, err := f.db.ExecContext(t.Context(), mutation, *f.link.MappingID)
			require.NoError(t, err)
			before := inactiveDraftState(t, f.db)
			f.group["name"] = "Must not save"
			response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+f.group["id"].(string), f.editRequest(t))
			require.Contains(t, []int{http.StatusNotFound, http.StatusConflict}, response.Code, response.Body.String())
			require.Equal(t, before, inactiveDraftState(t, f.db))
		})
	}
}
