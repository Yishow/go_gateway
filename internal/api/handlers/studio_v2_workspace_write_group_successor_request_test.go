package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func TestSuccessorHTTPRejectsAmbiguousStageInputs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewStudioV2WorkspaceWriteGroupSuccessorHandler(nil, nil, nil)
	base := map[string]any{"action": "reserve", "envelope": map[string]any{"version": "v1"}, "review_digest": "digest"}
	stage := workspace.WriteGroupSuccessorStageRequest{Envelope: workspace.WriteGroupSuccessorEnvelope{Version: "v1"}, ReviewDigest: "digest", Current: workspace.WriteGroupSuccessorCurrentCAS{WorkspaceRevision: "w", LineageRevision: 2, PredecessorRevision: "p", PredecessorConnectorRevision: "c", GroupID: "g2", GroupRevision: "r", TargetConnectorRevision: "t"}}
	stageBytes, err := json.Marshal(stage)
	require.NoError(t, err)
	var preview map[string]any
	require.NoError(t, json.Unmarshal(stageBytes, &preview))
	preview["action"] = "schema-preview"
	for _, tc := range []struct{ name, body string }{
		{"empty", `{}`}, {"null", `null`}, {"unknown-action", `{"action":"activate"}`}, {"trailing", `{"action":"reserve"} {}`},
		{"missing-scope", `{"action":"schema-preview"}`}, {"missing-confirmation", `{"action":"schema-confirm"}`},
		{"oversized", strings.Repeat(" ", 1<<20+1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, http.StatusBadRequest, successorHandlerStatus(t, h.Successor, tc.body))
		})
	}
	for _, tc := range []struct {
		name    string
		request map[string]any
		field   string
		value   any
	}{
		{"reserve-cannot-carry-schema-token", base, "token", "fake"},
		{"preview-cannot-carry-confirmation", preview, "schema_operation_id", "fake"},
		{"preview-cannot-carry-sql", preview, "sql", "DROP TABLE forbidden"},
		{"preview-cannot-select-bypass", preview, "allow_reserved", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			copy := map[string]any{}
			for k, v := range tc.request {
				copy[k] = v
			}
			copy[tc.field] = tc.value
			raw, err := json.Marshal(copy)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, successorHandlerStatus(t, h.Successor, string(raw)))
		})
	}
	// A syntactically valid envelope still cannot bypass missing service authority.
	raw, err := json.Marshal(base)
	require.NoError(t, err)
	require.Equal(t, http.StatusServiceUnavailable, successorHandlerStatus(t, h.Successor, string(raw)))
	require.Equal(t, http.StatusBadRequest, successorHandlerStatus(t, h.Review, `{"workspace_id":"w","allow_reserved":true}`))
}

func successorHandlerStatus(t *testing.T, handler gin.HandlerFunc, body string) int {
	t.Helper()
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request = httptest.NewRequest(http.MethodPost, "/successor", strings.NewReader(body))
	c.Params = gin.Params{{Key: "id", Value: "g1"}}
	handler(c)
	require.NotContains(t, response.Body.String(), "DROP TABLE")
	return response.Code
}
