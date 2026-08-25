package api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type hydrationBlockedWorkspaceActivator struct {
	called bool
}

func (a *hydrationBlockedWorkspaceActivator) ActivateEligible(context.Context) (*workspace.ActivationResponse, error) {
	a.called = true
	return &workspace.ActivationResponse{}, nil
}

func TestConfigureModbusShareActivation_BlocksFailedHydration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	share := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
	share.SetHydrationState(modbusshare.HydrationState{
		State:             modbusshare.HydrationStateFailed,
		WorkspaceID:       "ws-1",
		WorkspaceRevision: "rev-1",
		SettingsRevision:  "set-1",
		Readiness:         false,
	})
	activator := &hydrationBlockedWorkspaceActivator{}
	handler := handlers.NewStudioV2WorkspaceActivationHandler(activator)
	configureModbusShareActivation(handler, &DatalinkServices{ModbusShare: share})

	router := gin.New()
	router.POST("/activate", handler.Activate)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/activate", bytes.NewBufferString(`{"workspace_revision":"rev-1","settings_revision":"settings-1","readiness_token":"token-1"}`))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	require.Equal(t, http.StatusUnprocessableEntity, response.Code)
	require.Contains(t, response.Body.String(), modbusshare.ErrCodeHydrationRequired)
	require.False(t, activator.called)
}
