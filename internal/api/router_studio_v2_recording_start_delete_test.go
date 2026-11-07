package api

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestRecordingStartReadersAfterPermanentDelete(t *testing.T) {
	f := newRecordingStartFixture(t, recordingStartSuccess())
	request := f.request("completed-before-permanent-delete")
	completed, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, completed.Status)
	groupID := request.Groups[0].GroupID
	f.groups.WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{})
	current, err := f.groups.Get(t.Context(), groupID)
	require.NoError(t, err)
	stopped, err := f.groups.Disable(t.Context(), groupID, workspace.WriteGroupMutation{
		WorkspaceID: current.Group.WorkspaceID, ExpectedWorkspaceRevision: current.WorkspaceRevision,
		ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: current.Group.Destination.ConnectorRevision,
	})
	require.NoError(t, err)
	deleted, err := f.groups.Delete(t.Context(), groupID, workspace.WriteGroupMutation{
		WorkspaceID: stopped.Group.WorkspaceID, ExpectedWorkspaceRevision: stopped.WorkspaceRevision,
		ExpectedGroupRevision: stopped.Group.Revision, ExpectedConnectorRevision: stopped.Group.Destination.ConnectorRevision,
	})
	require.NoError(t, err)
	require.True(t, deleted.Deleted)
	var count int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM managed_schema_operations WHERE operation_id=?`, completed.OperationID).Scan(&count))
	require.Equal(t, 1, count, "completed start remains workspace-level historical evidence")

	response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+groupID, nil)
	require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
	response = performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.NotContains(t, response.Body.String(), groupID, "the home recording list must not recreate the deleted group")
	startPath := "/api/v1/datalink/studio-v2/workspace/recording-start"
	for _, read := range []struct {
		method, path string
		body         any
	}{
		{http.MethodGet, startPath + "/operations/" + completed.OperationID, nil},
		{http.MethodPost, startPath + "/lookup", request},
		{http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-operations/" + completed.OperationID, nil},
	} {
		response = performJSONRequest(t, f.router, read.method, read.path, read.body)
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		data := decodeJSONBody(t, response)["data"].(map[string]any)
		require.Equal(t, completed.OperationID, data["operation_id"])
		require.Equal(t, "succeeded", data["status"], "the retained result describes the historical start")
	}
	replayed, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, completed.OperationID, replayed.OperationID)
	require.Equal(t, 1, f.activation.Calls(), "replay must not reactivate the deleted group")

	_, err = f.db.ExecContext(t.Context(), `UPDATE devices SET description='' WHERE id='device-A'`)
	require.NoError(t, err)
	pipe := grouppipeline.New(grouppipeline.Dependencies{Groups: f.groups, Store: groupdelivery.NewStore(f.db)}, grouppipeline.Config{})
	rig := &runtimeRuleEvidenceRig{f: f.managedGroupFixture, pipe: pipe, clock: time.Now()}
	q := url.Values{"workspace_id": {f.record.ID}, "workspace_revision": {deleted.WorkspaceRevision}, "device_id": {"device-A"}}
	response = performJSONRequest(t, rig.deterministicRouter(f.groups), http.MethodGet, runtimeRuleContextPath+"?"+q.Encode(), nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	evidence := decodeJSONBody(t, response)["data"].(map[string]any)["rule_evidence"].(map[string]any)
	require.Empty(t, evidence["groups"], "runtime evidence must not report the deleted group as effective")
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_groups WHERE id=?`, groupID).Scan(&count))
	require.Zero(t, count)
}

func TestRecordingStartFailedReplayCannotRecreateDeletedGroup(t *testing.T) {
	f := newRecordingStartServiceFixture(t, newManagedGroupFixture(t, "sqlite"))
	request := f.request("failed-before-permanent-delete")
	failed, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationFailed, failed.Status)
	current, err := f.groups.Get(t.Context(), request.Groups[0].GroupID)
	require.NoError(t, err)
	_, err = f.groups.Delete(t.Context(), current.Group.ID, workspace.WriteGroupMutation{
		WorkspaceID: current.Group.WorkspaceID, ExpectedWorkspaceRevision: current.WorkspaceRevision,
		ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: current.Group.Destination.ConnectorRevision,
	})
	require.NoError(t, err)
	lookup := performJSONRequest(t, f.router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-start/lookup", request)
	require.Equal(t, http.StatusOK, lookup.Code, lookup.Body.String())
	require.Equal(t, "failed", decodeJSONBody(t, lookup)["data"].(map[string]any)["status"])
	_, err = f.start.Start(t.Context(), request)
	require.ErrorIs(t, err, workspace.ErrWriteGroupNotFound)
	require.Zero(t, f.activation.Calls())
	var operations, groups int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM managed_schema_operations WHERE operation_id=?`, failed.OperationID).Scan(&operations))
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_groups WHERE id=?`, current.Group.ID).Scan(&groups))
	require.Equal(t, 1, operations)
	require.Zero(t, groups)
}
