package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStudioV2WorkspaceActivationHandler_SaveIncompleteBlocksActivation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	activator := &stubWorkspaceActivator{
		response: &workspace.ActivationResponse{
			WorkspaceID: "ws-1",
		},
	}
	handler := NewStudioV2WorkspaceActivationHandler(activator)

	// Case 1: Pending saves > 0
	body, _ := json.Marshal(ActivateWorkspaceRequest{
		WorkspaceRevision: "rev-1",
		SettingsRevision:  "settings-1",
		ReadinessToken:    "token-1",
		PendingSaves:      2,
	})
	req := newHandlerTestRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/activate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Activate(c)

	require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	assert.Contains(t, resp.Body.String(), modbusshare.ErrCodeSaveIncomplete)
	assert.False(t, activator.called)

	// Case 2: Save error present
	activator.called = false
	bodyErr, _ := json.Marshal(ActivateWorkspaceRequest{
		WorkspaceRevision: "rev-1",
		SettingsRevision:  "settings-1",
		ReadinessToken:    "token-1",
		SaveError:         "network failed during autosave",
	})
	reqErr := newHandlerTestRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/activate", bytes.NewReader(bodyErr))
	reqErr.Header.Set("Content-Type", "application/json")
	respErr := httptest.NewRecorder()
	cErr, _ := gin.CreateTestContext(respErr)
	cErr.Request = reqErr

	handler.Activate(cErr)

	require.Equal(t, http.StatusUnprocessableEntity, respErr.Code)
	assert.Contains(t, respErr.Body.String(), modbusshare.ErrCodeSaveIncomplete)
	assert.False(t, activator.called)
}

func TestStudioV2WorkspaceActivationHandler_RevisionConflictBlocksActivation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	activator := &stubWorkspaceActivator{
		response: &workspace.ActivationResponse{
			WorkspaceID: "ws-1",
		},
	}
	handler := NewStudioV2WorkspaceActivationHandler(activator)
	handler.WithRevisionValidator(func(ctx context.Context, expectedWsRev, expectedSetRev string) error {
		if expectedWsRev == "stale-rev" {
			return &modbusshare.Error{
				Code:      modbusshare.ErrCodeRevisionConflict,
				Message:   "workspace revision conflict",
				Retryable: true,
			}
		}
		return nil
	})

	body, _ := json.Marshal(ActivateWorkspaceRequest{
		WorkspaceRevision: "stale-rev",
		SettingsRevision:  "settings-1",
		ReadinessToken:    "token-1",
	})
	req := newHandlerTestRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/activate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Activate(c)

	require.Equal(t, http.StatusConflict, resp.Code)
	assert.Contains(t, resp.Body.String(), modbusshare.ErrCodeRevisionConflict)
	assert.False(t, activator.called)
}
