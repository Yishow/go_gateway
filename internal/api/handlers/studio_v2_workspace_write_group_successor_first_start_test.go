package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func successorFirstStartEnvelope(t *testing.T) map[string]any {
	t.Helper()
	stage := workspace.WriteGroupSuccessorStageRequest{Envelope: workspace.WriteGroupSuccessorEnvelope{Version: "original-v1"}, ReviewDigest: "original-digest", Current: workspace.WriteGroupSuccessorCurrentCAS{
		WorkspaceRevision: "workspace", LineageRevision: 2, PredecessorRevision: "predecessor", PredecessorConnectorRevision: "source", GroupID: "reserved", GroupRevision: "target", TargetConnectorRevision: "connector",
	}}
	raw, err := json.Marshal(workspace.WriteGroupSuccessorActivationEnvelope{Version: "successor-activation-review-v2", Request: workspace.WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC)}})
	require.NoError(t, err)
	var result map[string]any
	require.NoError(t, json.Unmarshal(raw, &result))
	result["activation_kind"] = "first_start"
	return result
}

func TestSuccessorFirstStartHTTPRequiresExplicitCompleteStage(t *testing.T) {
	h := NewStudioV2WorkspaceWriteGroupSuccessorHandler(nil, nil, nil)
	request := map[string]any{"action": "first-start-confirm", "envelope": successorFirstStartEnvelope(t), "review_digest": strings.Repeat("a", 64), "readiness_digest": strings.Repeat("b", 64)}
	raw, err := json.Marshal(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, successorHandlerStatus(t, h.Successor, string(raw)), "well-formed first-start must dispatch to server authority, not unknown-action")
	for _, field := range []string{"token", "sql", "allow_reserved", "detail"} {
		request[field] = "forbidden"
		raw, err = json.Marshal(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, successorHandlerStatus(t, h.Successor, string(raw)))
		delete(request, field)
	}
	for _, body := range []string{`{"action":"first-start-confirm"}`, `{"action":"first-start-confirm","envelope":null}`, `{"action":"first-start-confirm"} {}`} {
		require.Equal(t, http.StatusBadRequest, successorHandlerStatus(t, h.Successor, body))
	}
}

func TestSuccessorOperationStatusHTTPRequiresExactLocalScope(t *testing.T) {
	h := NewStudioV2WorkspaceWriteGroupSuccessorHandler(nil, nil, nil)
	request := map[string]any{"action": "operation-status", "workspace_id": "workspace", "lineage_id": "lineage", "operation_id": "original-operation"}
	raw, err := json.Marshal(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, successorHandlerStatus(t, h.Successor, string(raw)), "well-formed recovery action must dispatch without live target authority")
	for _, field := range []string{"workspace_id", "lineage_id", "operation_id"} {
		original := request[field]
		request[field] = " "
		raw, err = json.Marshal(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, successorHandlerStatus(t, h.Successor, string(raw)))
		request[field] = original
	}
	for _, field := range []string{"token", "latest", "detail", "envelope", "activation_kind", "schema_operation_id"} {
		request[field] = "forbidden"
		raw, err = json.Marshal(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, successorHandlerStatus(t, h.Successor, string(raw)))
		delete(request, field)
	}
}
