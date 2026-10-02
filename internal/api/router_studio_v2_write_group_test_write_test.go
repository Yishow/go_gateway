package api

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/grouptestwrite"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// fakeTestWriter returns canned results so the route contract (status codes,
// typed errors, scope of the workspace identity) is tested independently of a
// real target; the service's own tests run against real SQLite and PostgreSQL.
type fakeTestWriter struct {
	preview    *grouptestwrite.Preview
	outcome    *grouptestwrite.Outcome
	err        error
	gotGroup   string
	gotWS      grouptestwrite.Workspace
	gotConfirm grouptestwrite.Confirmation
	previews   int
	confirms   int
}

func (f *fakeTestWriter) Preview(_ context.Context, ws grouptestwrite.Workspace, groupID string) (*grouptestwrite.Preview, error) {
	f.previews++
	f.gotWS, f.gotGroup = ws, groupID
	return f.preview, f.err
}

func (f *fakeTestWriter) Confirm(_ context.Context, ws grouptestwrite.Workspace, req grouptestwrite.Confirmation) (*grouptestwrite.Outcome, error) {
	f.confirms++
	f.gotWS, f.gotConfirm = ws, req
	return f.outcome, f.err
}

func (f writeGroupRouterFixture) routerWithTestWriter(writer *fakeTestWriter) http.Handler {
	return NewRouter(&DatalinkServices{
		Workspace: f.workspace, WriteGroups: workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db)),
		WriteGroupTestWrite: writer,
	})
}

func finishedOperation(outcome, cleanup string) *recordingplan.SchemaOperation {
	now := time.Now().UTC()
	return &recordingplan.SchemaOperation{
		OperationID: "op-1", Action: recordingplan.TestWriteAction, Status: recordingplan.SchemaOperationSucceeded,
		WriteOutcome: outcome, CleanupStatus: cleanup, CreatedAt: now, UpdatedAt: now, CompletedAt: &now,
	}
}

func TestWritePreviewRouteIsScopedToTheWorkspaceAndGroup(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	writer := &fakeTestWriter{preview: &grouptestwrite.Preview{Token: "tok-1", OperationID: "op-1", Action: recordingplan.TestWriteAction, GroupID: "g-1"}}
	router := f.routerWithTestWriter(writer)

	response := performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/g-1/test-write-preview", nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	data := decodeJSONBody(t, response)["data"].(map[string]any)
	require.Equal(t, "tok-1", data["token"])
	require.Equal(t, "test_write", data["action"])
	require.Equal(t, "g-1", writer.gotGroup)
	require.Equal(t, f.record.ID, writer.gotWS.ID, "the workspace is server-resolved, never taken from the client")
	require.Equal(t, f.record.DatabaseSetupRevision, writer.gotWS.Revision)
}

func TestWriteRoutesMapServiceErrorsToTypedStatuses(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	cases := []struct {
		name   string
		err    error
		status int
		code   string
	}{
		{"unknown group", grouptestwrite.ErrGroupNotFound, http.StatusNotFound, "WRITE_GROUP_NOT_FOUND"},
		{"unsupported ownership", &grouptestwrite.UnsupportedError{Reason: grouptestwrite.ReasonOwnershipUnsafe}, http.StatusUnprocessableEntity, "WRITE_GROUP_TEST_WRITE_UNSUPPORTED"},
		{"schema token", recordingplan.ErrPreviewTokenKind, http.StatusUnprocessableEntity, "WRITE_GROUP_TEST_WRITE_TOKEN_KIND"},
		{"unknown or foreign token", recordingplan.ErrPreviewTokenNotFound, http.StatusNotFound, "WRITE_GROUP_TEST_WRITE_PREVIEW_NOT_FOUND"},
		{"stale", recordingplan.ErrPreviewTokenStale, http.StatusConflict, "WRITE_GROUP_TEST_WRITE_PREVIEW_STALE"},
		{"preview of another group", grouptestwrite.ErrGroupMismatch, http.StatusConflict, "WRITE_GROUP_TEST_WRITE_PREVIEW_STALE"},
		{"destination changed", grouptestwrite.ErrDestinationChanged, http.StatusConflict, "WRITE_GROUP_TEST_WRITE_PREVIEW_STALE"},
		{"expired", recordingplan.ErrPreviewTokenExpired, http.StatusConflict, "WRITE_GROUP_TEST_WRITE_PREVIEW_EXPIRED"},
		{"operation mismatch", recordingplan.ErrSchemaOperationMismatch, http.StatusConflict, "WRITE_GROUP_TEST_WRITE_OPERATION_MISMATCH"},
		{"unknown failure", errors.New("driver: password=secret host=db"), http.StatusServiceUnavailable, "WRITE_GROUP_TEST_WRITE_UNAVAILABLE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			writer := &fakeTestWriter{err: tc.err}
			router := f.routerWithTestWriter(writer)
			for _, request := range []struct {
				path string
				body map[string]any
			}{
				{writeGroupsPath + "/g-1/test-write", map[string]any{"token": "tok-1", "operation_id": "op-1"}},
			} {
				response := performJSONRequest(t, router, http.MethodPost, request.path, request.body)
				require.Equal(t, tc.status, response.Code, response.Body.String())
				body := decodeJSONBody(t, response)
				require.Equal(t, false, body["success"])
				require.Equal(t, tc.code, body["error"].(map[string]any)["code"])
				require.NotContains(t, response.Body.String(), "secret", "driver details must never reach the client")
			}
		})
	}
}

func TestWriteConfirmRouteStatusesForRunningRetainedAndBusy(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	body := map[string]any{"token": "tok-1", "operation_id": "op-1"}
	path := writeGroupsPath + "/g-1/test-write"

	running := &fakeTestWriter{outcome: &grouptestwrite.Outcome{Claim: recordingplan.ClaimInProgress, Operation: &recordingplan.SchemaOperation{OperationID: "op-1", Status: recordingplan.SchemaOperationRunning}}}
	response := performJSONRequest(t, f.routerWithTestWriter(running), http.MethodPost, path, body)
	require.Equal(t, http.StatusAccepted, response.Code, "a running duplicate is 202: %s", response.Body.String())

	retained := &fakeTestWriter{outcome: &grouptestwrite.Outcome{Claim: recordingplan.ClaimCompleted, Operation: finishedOperation("written_unverified", "failed")}}
	response = performJSONRequest(t, f.routerWithTestWriter(retained), http.MethodPost, path, body)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	data := decodeJSONBody(t, response)["data"].(map[string]any)
	require.Equal(t, "written_unverified", data["write_outcome"])
	require.Equal(t, "failed", data["cleanup_status"], "write verification and cleanup are separate fields")
	require.NotContains(t, data, "owner")
	require.NotContains(t, data, "detail")

	busy := &fakeTestWriter{outcome: &grouptestwrite.Outcome{Claim: recordingplan.ClaimScopeBusy, Operation: &recordingplan.SchemaOperation{OperationID: "op-holder"}}}
	response = performJSONRequest(t, f.routerWithTestWriter(busy), http.MethodPost, path, body)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	envelope := decodeJSONBody(t, response)["error"].(map[string]any)
	require.Equal(t, "WRITE_GROUP_TEST_WRITE_BUSY", envelope["code"])
	require.Equal(t, "op-holder", envelope["operation_id"], "the client can follow the holder through the shared status endpoint")

	unknown := &fakeTestWriter{
		outcome: &grouptestwrite.Outcome{Operation: &recordingplan.SchemaOperation{OperationID: "op-1"}},
		err:     grouptestwrite.ErrResultUnacknowledged,
	}
	response = performJSONRequest(t, f.routerWithTestWriter(unknown), http.MethodPost, path, body)
	require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
	envelope = decodeJSONBody(t, response)["error"].(map[string]any)
	require.Equal(t, "WRITE_GROUP_TEST_WRITE_RESULT_UNKNOWN", envelope["code"])
	require.Equal(t, "op-1", envelope["operation_id"])
}

func TestWriteConfirmRouteRejectsIncompleteRequestsBeforeTheService(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	writer := &fakeTestWriter{}
	router := f.routerWithTestWriter(writer)
	for name, body := range map[string]map[string]any{
		"empty":               {},
		"missing operation":   {"token": "tok-1"},
		"missing token":       {"operation_id": "op-1"},
		"blank values":        {"token": " ", "operation_id": " "},
		"bare legacy plan id": {"plan_id": "plan-1"},
	} {
		t.Run(name, func(t *testing.T) {
			for _, path := range []string{writeGroupsPath + "/g-1/test-write", "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write"} {
				response := performJSONRequest(t, router, http.MethodPost, path, body)
				require.Equal(t, http.StatusBadRequest, response.Code, "%s: %s", path, response.Body.String())
			}
		})
	}
	require.Zero(t, writer.confirms, "nothing reaches the service for an incomplete confirmation")
}

func TestWriteLegacyPlanRoutesResolveToTheSameGroupOrRefuse(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	writer := &fakeTestWriter{preview: &grouptestwrite.Preview{Token: "tok-1", OperationID: "op-1"}}
	router := f.routerWithTestWriter(writer)

	// No group carries provenance for this plan, so nothing may be guessed.
	response := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write-preview", map[string]any{"plan_id": "plan-1"})
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	require.Equal(t, "RECORDING_TEST_WRITE_PLAN_UNRESOLVED", decodeJSONBody(t, response)["error"].(map[string]any)["code"])
	require.Zero(t, writer.previews)

	response = performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write",
		map[string]any{"plan_id": "plan-1", "token": "tok-1", "operation_id": "op-1"})
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, "an unresolvable plan is refused before the confirmation runs: %s", response.Body.String())
	require.Zero(t, writer.confirms)

	response = performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write-preview", map[string]any{})
	require.Equal(t, http.StatusBadRequest, response.Code)
}

func TestWriteRoutesAnswerUnavailableWithoutAService(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	for _, path := range []string{writeGroupsPath + "/g-1/test-write-preview", writeGroupsPath + "/g-1/test-write"} {
		response := performJSONRequest(t, f.router, http.MethodPost, path, map[string]any{"token": "tok-1", "operation_id": "op-1"})
		require.Equal(t, http.StatusServiceUnavailable, response.Code, "%s: %s", path, response.Body.String())
		require.Equal(t, "WRITE_GROUP_TEST_WRITE_UNAVAILABLE", decodeJSONBody(t, response)["error"].(map[string]any)["code"])
	}
}

func TestWriteConfirmRoutesPassTheAddressedGroupToTheService(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	writer := &fakeTestWriter{outcome: &grouptestwrite.Outcome{Claim: recordingplan.ClaimCompleted, Operation: finishedOperation("written_verified", "cleaned")}}
	router := f.routerWithTestWriter(writer)
	response := performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/g-9/test-write", map[string]any{"token": "tok-1", "operation_id": "op-1"})
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.Equal(t, "g-9", writer.gotConfirm.GroupID, "the URL group is checked against the preview's group by the service")
}
