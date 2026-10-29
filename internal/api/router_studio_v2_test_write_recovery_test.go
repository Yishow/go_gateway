package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
	"net/http"
	"testing"
)

func TestStudioV2TestWriteRecoveryListsLostReference(t *testing.T) {
	f := newRecordingRouteFixture(t)
	token, err := f.planSvc.PrepareTestWritePreview(t.Context(), recordingplan.TestWriteScope{
		WorkspaceID: f.workspaceID, WorkspaceRevision: "historical-workspace", GroupID: "old-group", GroupRevision: "old-revision",
		ConnectorID: "old-connector", ConnectorRevision: "old-connector-revision", Dialect: "sqlite", Database: "private-endpoint", Schema: "main", Table: "historical_table",
	}, "content")
	require.NoError(t, err)
	op, _, err := f.planSvc.ClaimTestWrite(t.Context(), token)
	require.NoError(t, err)
	response := serveRecordingScopedRequest(f.router, http.MethodGet, databaseOperationsPath+"test-writes?workspace_id=foreign", "")
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var result struct {
		Success bool `json:"success"`
		Data    struct {
			WorkspaceID string `json:"workspace_id"`
			Complete    bool   `json:"complete"`
			Operations  []struct {
				OperationID   string `json:"operation_id"`
				GroupID       string `json:"group_id"`
				GroupRevision string `json:"group_revision"`
			} `json:"operations"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
	require.True(t, result.Success)
	require.True(t, result.Data.Complete)
	require.Equal(t, f.workspaceID, result.Data.WorkspaceID)
	require.Len(t, result.Data.Operations, 1)
	require.Equal(t, op.OperationID, result.Data.Operations[0].OperationID)
	require.Equal(t, "old-group", result.Data.Operations[0].GroupID)
	require.Equal(t, "old-revision", result.Data.Operations[0].GroupRevision)
	for _, secret := range []string{token.Token, op.Owner, "private-endpoint", "historical_table", "scope_key", "detail"} {
		require.NotContains(t, response.Body.String(), secret)
	}
}

func TestStudioV2TestWriteRecoveryForeignIsolationEmptyAndMissingBinding(t *testing.T) {
	f := newRecordingRouteFixture(t)
	foreignToken, err := f.planSvc.PrepareTestWritePreview(t.Context(), recordingplan.TestWriteScope{
		WorkspaceID: "foreign", WorkspaceRevision: "foreign-rev", GroupID: "foreign-group", GroupRevision: "foreign-group-rev",
		ConnectorID: "foreign-connector", ConnectorRevision: "foreign-connector-rev", Dialect: "sqlite", Database: "secret-foreign", Table: "foreign-table",
	}, "content")
	require.NoError(t, err)
	foreignOp, _, err := f.planSvc.ClaimTestWrite(t.Context(), foreignToken)
	require.NoError(t, err)
	empty := serveRecordingScopedRequest(f.router, http.MethodGet, databaseOperationsPath+"test-writes?workspace_id=foreign", "")
	require.Equal(t, http.StatusOK, empty.Code)
	require.Contains(t, empty.Body.String(), `"operations":[]`)
	require.Contains(t, empty.Body.String(), `"complete":true`)
	require.NotContains(t, empty.Body.String(), foreignOp.OperationID)
	for _, id := range []string{foreignOp.OperationID, "missing"} {
		response := serveRecordingScopedRequest(f.router, http.MethodGet, databaseOperationsPath+id, "")
		assertRecordingPlanError(t, response, http.StatusNotFound, "RECORDING_OPERATION_NOT_FOUND", false)
	}
	token, err := f.planSvc.PrepareTestWritePreview(t.Context(), recordingplan.TestWriteScope{
		WorkspaceID: f.workspaceID, WorkspaceRevision: "original", GroupID: "original-group", GroupRevision: "original-rev", ConnectorID: "connector", ConnectorRevision: "revision", Dialect: "sqlite", Table: "table",
	}, "content")
	require.NoError(t, err)
	original, _, err := f.planSvc.ClaimTestWrite(t.Context(), token)
	require.NoError(t, err)
	require.NoError(t, f.planRepo.DeletePreviewToken(t.Context(), token.Token))
	response := serveRecordingScopedRequest(f.router, http.MethodGet, databaseOperationsPath+"test-writes", "")
	assertRecordingPlanError(t, response, http.StatusServiceUnavailable, "RECORDING_PLAN_UNAVAILABLE", true)
	require.NotContains(t, response.Body.String(), `"operations"`)
	after, err := f.planRepo.GetSchemaOperation(t.Context(), f.workspaceID, original.OperationID)
	require.NoError(t, err)
	require.Equal(t, original, after)
}
