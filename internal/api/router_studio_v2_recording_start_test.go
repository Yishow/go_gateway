package api

import (
	"go-gateway/internal/datalink/recordingplan"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecordingStartMissingPreparationHasNoDDL(t *testing.T) {
	f := newManagedGroupFixture(t, "sqlite")
	response := performJSONRequest(t, f.router, http.MethodPost,
		strings.TrimSuffix(writeGroupsPath, "/write-groups")+"/recording-start", map[string]any{
			"request_id": "start-request-1", "workspace_id": f.record.ID,
			"expected_workspace_revision": f.data["workspace_revision"],
			"device_ids":                  []string{"device-A"},
			"groups": []map[string]any{{
				"group_id": f.group["id"], "expected_group_revision": f.group["revision"],
				"expected_connector_revision": "connector-1",
			}},
		})
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	body := decodeJSONBody(t, response)
	op := body["data"].(map[string]any)
	require.Equal(t, "recording_start", op["action"])
	require.Equal(t, "failed", op["status"])
	require.Equal(t, "prepare_schema", op["next_action"])
	require.NoFileExists(t, f.file)
	var plans int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM recording_plans`).Scan(&plans))
	require.Zero(t, plans)
}

func TestRecordingStartWithoutDatabaseRequiresEnabledShareOutput(t *testing.T) {
	f := newManagedGroupFixture(t, "sqlite")
	response := performJSONRequest(t, f.router, http.MethodPost,
		strings.TrimSuffix(writeGroupsPath, "/write-groups")+"/recording-start", map[string]any{
			"request_id": "start-no-output", "workspace_id": f.record.ID,
			"expected_workspace_revision": f.data["workspace_revision"],
			"device_ids":                  []string{"device-A"}, "groups": []any{},
		})
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	op := decodeJSONBody(t, response)["data"].(map[string]any)
	require.Equal(t, "failed", op["status"])
	require.Equal(t, "invalid_scope", op["reason"])
	require.Equal(t, "review_selection", op["next_action"])
	require.NoFileExists(t, f.file)
}

func TestRecordingStartBusyProvidesRetryInsteadOfDisablingHealthyGroups(t *testing.T) {
	f := newManagedGroupFixture(t, "sqlite")
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(f.db))
	_, _, err := ledger.ClaimRecordingStart(t.Context(), f.record.ID, "other-active-request", strings.Repeat("a", 64), "{}")
	require.NoError(t, err)
	response := performJSONRequest(t, f.router, http.MethodPost,
		strings.TrimSuffix(writeGroupsPath, "/write-groups")+"/recording-start", map[string]any{
			"request_id": "new-request", "workspace_id": f.record.ID,
			"expected_workspace_revision": f.data["workspace_revision"],
			"device_ids":                  []string{"device-A"},
			"groups":                      []map[string]any{{"group_id": f.group["id"], "expected_group_revision": f.group["revision"], "expected_connector_revision": "connector-1"}},
		})
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, response), "RECORDING_START_BUSY")
	require.Contains(t, response.Body.String(), "retry_same_request")
	require.NotContains(t, response.Body.String(), "disable_group")
}

func TestRecordingStartCompletedReplaySurvivesLaterGroupDeletion(t *testing.T) {
	f := newRecordingStartFixture(t)
	request := f.request("completed-before-delete")
	completed, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, completed.Status)
	_, err = f.db.ExecContext(t.Context(), `UPDATE write_groups SET status='deleted' WHERE id=?`, request.Groups[0].GroupID)
	require.NoError(t, err)
	replayed, err := f.start.Start(t.Context(), request)
	require.NoError(t, err, "a completed result is historical evidence, independent of later resource edits")
	require.Equal(t, completed.OperationID, replayed.OperationID)
	require.Equal(t, completed.Groups, replayed.Groups)
	require.Equal(t, 1, f.activation.Calls(), "replay never reactivates the deleted group")
}

func TestRecordingStartRevalidatesAfterExplicitPreparation(t *testing.T) {
	managed := newManagedGroupFixture(t, "sqlite")
	f := newRecordingStartServiceFixture(t, managed)
	request := f.request("prepare-then-resume")
	first, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, "preparation_required", first.Reason)
	require.NoFileExists(t, f.file)
	token := managed.preview(t)
	confirm := managed.schemaRequest()
	confirm["token"], confirm["operation_id"] = token.Token, token.OperationID
	response := performJSONRequest(t, managed.router, http.MethodPost, managed.confirmPath(), confirm)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	current, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	saved, err := f.groups.Get(t.Context(), request.Groups[0].GroupID)
	require.NoError(t, err)
	resumed, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, "stale_intent", resumed.Reason)
	require.Equal(t, "revalidate", resumed.NextAction)
	require.Equal(t, first.OperationID, resumed.OperationID)
	require.Zero(t, f.activation.Calls())
	request.RequestID = "start-confirmed-preparation"
	request.ExpectedWorkspaceRevision = current.DatabaseSetupRevision
	request.Groups[0].ExpectedGroupRevision = saved.Group.Revision
	started, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, started.Status, "%+v", started)
	require.NotEqual(t, first.OperationID, started.OperationID)
	require.Equal(t, 1, f.activation.Calls())
	var starts int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM managed_schema_operations WHERE action='recording_start'`).Scan(&starts))
	require.Equal(t, 2, starts, "the unprepared attempt remains historical; one confirmed intent applies")
}
