package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Production SQL repositories + real HTTP handlers. A revisioned proposal is
// not the applied workspace mapping projection; clients must join exact IDs
// through the current source-rule link projection, not invent proposal IDs.
func TestRecordingEditorRangeReadbackUsesLiveWorkspaceMappingIdentity(t *testing.T) {
	f := newWorkspaceConcurrentFixture(t, 1, 2, nil)
	router := gin.New()
	router.GET("/rules/:id/candidates", NewSourceRuleHandler(f.rules).Candidates)
	router.GET("/workspace/mappings", f.h.List)
	router.GET("/points/:id", NewPointHandler(f.h.pointSvc).Get)
	router.GET("/tags/:id", NewTagHandler(f.tags).Get)
	get := func(path string) map[string]any {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequestWithContext(f.ctx, http.MethodGet, path, nil))
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		var body map[string]any
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
		return body
	}
	candidates := get("/rules/diag-rule/candidates")["data"].(map[string]any)["tags"].(map[string]any)["candidates"].([]any)
	mappings := get("/workspace/mappings")["data"].([]any)
	require.Len(t, candidates, 2)
	require.Len(t, mappings, 2)
	for _, raw := range candidates {
		proposal := raw.(map[string]any)
		require.Empty(t, proposal["mapping_id"])
		require.Empty(t, proposal["tag_id"])
		require.Equal(t, "draft", proposal["status"])
		var exact map[string]any
		for _, row := range mappings {
			applied := row.(map[string]any)
			if applied["point_id"] == proposal["point_id"] && applied["address"] == proposal["address"] {
				exact = applied
			}
		}
		require.NotNil(t, exact)
		require.Equal(t, "diag-rule", exact["rule_id"])
		require.Equal(t, "diag-dev", exact["device_id"])
		require.Contains(t, f.ids, exact["id"])
		point := get("/points/" + exact["point_id"].(string))["data"].(map[string]any)
		tag := get("/tags/" + exact["tag_id"].(string))["data"].(map[string]any)
		require.Equal(t, exact["point_id"], point["id"])
		require.Equal(t, exact["device_id"], point["device_id"])
		require.Equal(t, exact["tag_id"], tag["id"])
		require.Equal(t, exact["tag_key"], tag["key"])
		require.Equal(t, "", exact["unit"], "workspace mapping serializes an empty unit")
		require.NotContains(t, tag, "unit", "Tag GET omits empty unit")
		wire, err := json.Marshal(tag)
		require.NoError(t, err)
		t.Logf("actual-empty-unit-tag-wire=%s", wire)
	}
}
