package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStudioV2WorkspaceSourceRulesHandler_ShareConfigRoundTripsAndClears(t *testing.T) {
	fixture := newWorkspaceSourceRuleReconcileFixture(t)

	updateReq := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/api/v1/datalink/studio-v2/workspace/source-rules/rule-A", strings.NewReader(`{
		"share_enabled":true,
		"share_start_register":40001,
		"share_stride":2
	}`))
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	updateContext, _ := gin.CreateTestContext(updateResp)
	updateContext.Request = updateReq
	updateContext.Params = gin.Params{{Key: "id", Value: "rule-A"}}
	fixture.handler.Update(updateContext)

	require.Equal(t, http.StatusOK, updateResp.Code, updateResp.Body.String())
	assert.Contains(t, updateResp.Body.String(), `"share_enabled":true`)
	assert.Contains(t, updateResp.Body.String(), `"share_start_register":40001`)
	assert.Contains(t, updateResp.Body.String(), `"share_stride":2`)

	clearReq := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/api/v1/datalink/studio-v2/workspace/source-rules/rule-A", strings.NewReader(`{
		"share_enabled":false,
		"share_start_register":null,
		"share_stride":null
	}`))
	clearReq.Header.Set("Content-Type", "application/json")
	clearResp := httptest.NewRecorder()
	clearContext, _ := gin.CreateTestContext(clearResp)
	clearContext.Request = clearReq
	clearContext.Params = gin.Params{{Key: "id", Value: "rule-A"}}
	fixture.handler.Update(clearContext)

	require.Equal(t, http.StatusOK, clearResp.Code, clearResp.Body.String())
	assert.Contains(t, clearResp.Body.String(), `"share_enabled":false`)
	assert.Contains(t, clearResp.Body.String(), `"share_start_register":null`)
	assert.Contains(t, clearResp.Body.String(), `"share_stride":null`)
}
