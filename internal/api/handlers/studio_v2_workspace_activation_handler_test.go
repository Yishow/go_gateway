package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type stubWorkspaceActivator struct {
	response *workspace.ActivationResponse
	err      error
}

func (s *stubWorkspaceActivator) ActivateEligible(_ context.Context) (*workspace.ActivationResponse, error) {
	return s.response, s.err
}

func TestStudioV2WorkspaceActivationHandler_ReturnsPerDeviceResultsOnPartialSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewStudioV2WorkspaceActivationHandler(&stubWorkspaceActivator{
		response: &workspace.ActivationResponse{
			WorkspaceID: "workspace-1",
			Results: []workspace.ActivationResult{
				{DeviceID: "dev-A", Status: workspace.ActivationResultStatusSuccess, Message: "activated"},
				{DeviceID: "dev-B", Status: workspace.ActivationResultStatusFailed, Message: "activation timeout"},
			},
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/activate", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Activate(c)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}
	if body := resp.Body.String(); body == "" || !containsAll(body, "workspace-1", "dev-A", "success", "dev-B", "failed", "activation timeout") {
		t.Fatalf("expected per-device partial success payload, got %s", body)
	}
}

func containsAll(body string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(body, part) {
			return false
		}
	}
	return true
}
