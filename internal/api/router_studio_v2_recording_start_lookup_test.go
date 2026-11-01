package api

import (
	"net/http"
	"testing"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

const recordingStartLookupPath = "/api/v1/datalink/studio-v2/workspace/recording-start/lookup"

// A query-only facade must not silently become a second start attempt. These
// triggers also catch readiness/probe persistence and workspace initialization.
func forbidRecordingLookupWrites(t *testing.T, f recordingStartFixture) {
	t.Helper()
	for _, table := range []string{"system_settings", "managed_schema_operations", "write_groups", "devices", "points", "tags", "mappings", "measurement_definitions"} {
		for _, verb := range []string{"INSERT", "UPDATE", "DELETE"} {
			_, err := f.db.ExecContext(t.Context(), `CREATE TRIGGER lookup_no_`+verb+`_`+table+` BEFORE `+verb+` ON `+table+` BEGIN SELECT RAISE(ABORT,'lookup must be read only'); END`)
			require.NoError(t, err)
		}
	}
}

func TestRecordingStartLookupRecoversLostAcknowledgementReadOnly(t *testing.T) {
	f := newRecordingStartFixture(t)
	request := f.request("lost-ack-no-operation-id")
	original, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, original.Status)
	// Historical results survive current source and group head edits.
	_, err = f.db.ExecContext(t.Context(), `UPDATE devices SET connection_config='{"changed":true}' WHERE id='device-A'`)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE write_groups SET revision='later-revision' WHERE id=?`, request.Groups[0].GroupID)
	require.NoError(t, err)
	forbidRecordingLookupWrites(t, f)
	response := performJSONRequest(t, f.router, http.MethodPost, recordingStartLookupPath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	result := decodeJSONBody(t, response)["data"].(map[string]any)
	require.Equal(t, original.OperationID, result["operation_id"])
	require.Equal(t, string(original.Status), result["status"])
	require.Equal(t, original.IntentDigest, result["intent_digest"])
	require.Equal(t, 1, f.activation.Calls())
	get := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/recording-start/operations/"+original.OperationID, nil)
	require.Equal(t, http.StatusOK, get.Code, get.Body.String())
	require.Equal(t, decodeJSONBody(t, get)["data"], result)
}

func TestRecordingStartLookupRejectsUnknownOrMismatchedScopeWithoutWrites(t *testing.T) {
	for _, mode := range []string{"missing", "foreign workspace", "changed payload", "wrong action", "corrupt detail", "missing workspace"} {
		t.Run(mode, func(t *testing.T) {
			f := newRecordingStartFixture(t)
			request := f.request("lookup-guards")
			original, err := f.start.Start(t.Context(), request)
			require.NoError(t, err)
			expected := http.StatusNotFound
			switch mode {
			case "missing":
				request.RequestID = "never-submitted"
			case "foreign workspace":
				request.WorkspaceID = "foreign-workspace"
				expected = http.StatusBadRequest
			case "changed payload":
				request.Groups[0].ExpectedGroupRevision = "different"
				expected = http.StatusConflict
			case "wrong action":
				_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET action='schema_apply' WHERE operation_id=?`, original.OperationID)
				require.NoError(t, err)
				expected = http.StatusConflict
			case "corrupt detail":
				_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET detail='{}' WHERE operation_id=?`, original.OperationID)
				require.NoError(t, err)
				expected = http.StatusInternalServerError
			case "missing workspace":
				_, err = f.db.ExecContext(t.Context(), `DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
				require.NoError(t, err)
			}
			forbidRecordingLookupWrites(t, f)
			response := performJSONRequest(t, f.router, http.MethodPost, recordingStartLookupPath, request)
			require.Equal(t, expected, response.Code, response.Body.String())
			body := decodeJSONBody(t, response)
			require.NotContains(t, body, "data")
			if expected == http.StatusNotFound {
				require.Equal(t, "RECORDING_START_NOT_FOUND", body["error"].(map[string]any)["code"])
			}
			require.Equal(t, 1, f.activation.Calls())
		})
	}
}

func TestRecordingStartLookupPreservesUnresolvedStatusWithoutResuming(t *testing.T) {
	for _, state := range []string{"running", "unknown", "partial"} {
		t.Run(state, func(t *testing.T) {
			f := newRecordingStartFixture(t)
			request := f.request("unresolved-lookup")
			original, err := f.start.Start(t.Context(), request)
			require.NoError(t, err)
			_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET status=? WHERE operation_id=?`, state, original.OperationID)
			require.NoError(t, err)
			forbidRecordingLookupWrites(t, f)
			response := performJSONRequest(t, f.router, http.MethodPost, recordingStartLookupPath, request)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			result := decodeJSONBody(t, response)["data"].(map[string]any)
			require.Equal(t, state, result["status"])
			require.Equal(t, original.OperationID, result["operation_id"])
			require.Equal(t, 1, f.activation.Calls())
		})
	}
}
