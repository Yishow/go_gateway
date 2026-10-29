package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func TestSuccessorActivationHTTPRequiresExplicitCompleteStage(t *testing.T) {
	h := NewStudioV2WorkspaceWriteGroupSuccessorHandler(nil, nil, nil)
	stage := workspace.WriteGroupSuccessorStageRequest{Envelope: workspace.WriteGroupSuccessorEnvelope{Version: "original-v1"}, ReviewDigest: "original-digest", Current: workspace.WriteGroupSuccessorCurrentCAS{
		WorkspaceRevision: "workspace", LineageRevision: 2, PredecessorRevision: "predecessor", PredecessorConnectorRevision: "source", GroupID: "reserved", GroupRevision: "target", TargetConnectorRevision: "connector",
	}}
	envelope := workspace.WriteGroupSuccessorActivationEnvelope{Version: "activation-v1", Request: workspace.WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC)}}
	digest := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	for _, request := range []map[string]any{
		{"action": "activation-review", "stage": stage, "boundary": envelope.Request.Boundary},
		{"action": "activation-readiness", "envelope": envelope, "review_digest": digest},
		{"action": "activation-confirm", "envelope": envelope, "review_digest": digest, "readiness_digest": digest},
	} {
		raw, err := json.Marshal(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusServiceUnavailable, successorHandlerStatus(t, h.Successor, string(raw)), "valid stage still needs server authority")
		request["allow_reserved"] = true
		raw, err = json.Marshal(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, successorHandlerStatus(t, h.Successor, string(raw)))
	}
	for _, body := range []string{
		`{"action":"activation-review"}`, `{"action":"activation-readiness"}`, `{"action":"activation-confirm"}`,
		`{"action":"activation-confirm","review_digest":"legacy","readiness_digest":"legacy"}`, `{"action":"activation-confirm"} {}`,
	} {
		require.Equal(t, http.StatusBadRequest, successorHandlerStatus(t, h.Successor, body))
	}
}
