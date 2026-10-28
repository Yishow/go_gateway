package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
)

func successorStatusBody(f successorAPIFixture) map[string]any {
	return map[string]any{"action": "operation-status", "workspace_id": f.request.WorkspaceID, "lineage_id": f.request.LineageID, "operation_id": f.request.OperationID}
}
func successorStatus(t *testing.T, f successorAPIFixture) map[string]any {
	t.Helper()
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", successorStatusBody(f))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	for _, field := range []string{`"token":`, `"detail":`, `"dsn":`, `"connection_config":`, `"owner":`} {
		require.NotContains(t, response.Body.String(), field)
	}
	return decodeJSONBody(t, response)["data"].(map[string]any)
}
func successorSchemaStatuses(t *testing.T, status map[string]any) map[string]string {
	t.Helper()
	result := map[string]string{}
	for _, raw := range status["schema_operations"].([]any) {
		operation := raw.(map[string]any)
		id := operation["operation_id"].(string)
		require.NotContains(t, result, id, "related operation cannot appear more than once")
		result[id] = operation["status"].(string)
	}
	return result
}

func TestSuccessorOperationStatusKeepsUnknownWhenNewPreviewAndHeadMove(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	review := f.review(t)
	prepared := f.reserve(t, review)
	first := successorAPIPreview(t, f, review, prepared)
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(f.f.db))
	operation, outcome, err := ledger.ClaimSchemaApply(t.Context(), &first)
	require.NoError(t, err)
	require.Equal(t, recordingplan.ClaimAcquired, outcome)
	_, err = ledger.FinishSchemaOperation(t.Context(), operation.OperationID, operation.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationUnknown, Reason: "outcome_unknown", NextAction: "inspect_original_operation"})
	require.NoError(t, err)
	second := successorAPIPreview(t, f, review, prepared)
	require.NotEqual(t, first.OperationID, second.OperationID)
	before := successorAPIAllState(t, f.f.db)
	status := successorStatus(t, f)
	require.Equal(t, map[string]string{first.OperationID: "unknown", second.OperationID: "preview_not_submitted"}, successorSchemaStatuses(t, status))
	require.Equal(t, before, successorAPIAllState(t, f.f.db), "status may not repair or mutate original outcomes")
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", successorAPIConfirmRequest(t, f, review, prepared, second))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	activation := successorActivationAPIFixture{successorAPIFixture: f, reviewResult: review, prepared: prepared, boundary: f.request.Boundary}
	event := activation.postActivation(t, activation.activationConfirmation(t))
	before = successorAPIAllState(t, f.f.db)
	status = successorStatus(t, f)
	require.Equal(t, map[string]string{first.OperationID: "unknown", second.OperationID: "succeeded"}, successorSchemaStatuses(t, status))
	require.Equal(t, event, status["cutover"].(map[string]any)["result"])
	require.Equal(t, "not_submitted", status["first_start"].(map[string]any)["status"])
	require.Equal(t, before, successorAPIAllState(t, f.f.db))
	require.NotContains(t, status["cutover"], "intent")
}

func TestSuccessorOperationStatusRejectsForeignRootScopeWithoutMutation(t *testing.T) {
	f := newSuccessorActivationAPIFixture(t)
	before := successorAPIAllState(t, f.f.db)
	for _, field := range []string{"workspace_id", "lineage_id", "operation_id", "path", "root_action"} {
		t.Run(field, func(t *testing.T) {
			request := successorStatusBody(f.successorAPIFixture)
			id := f.groupID
			switch field {
			case "path":
				id = f.prepared.Group.ID
			case "root_action":
				var schema string
				require.NoError(t, f.f.db.QueryRow(`SELECT operation_id FROM managed_schema_operations WHERE action='schema_apply' LIMIT 1`).Scan(&schema))
				request["operation_id"] = schema
			default:
				request[field] = "foreign"
			}
			response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+id+"/successor", request)
			require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
			require.Less(t, response.Code, 500, response.Body.String())
			require.Equal(t, before, successorAPIAllState(t, f.f.db))
			require.NotContains(t, response.Body.String(), f.prepared.Group.ID)
		})
	}
}
