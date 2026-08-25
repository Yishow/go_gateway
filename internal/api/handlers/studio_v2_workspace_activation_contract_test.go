package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go-gateway/internal/datalink/modbusshare"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newValidActivationRequest(t *testing.T) *http.Request {
	t.Helper()
	body, err := json.Marshal(ActivateWorkspaceRequest{
		WorkspaceRevision: "workspace-revision-1",
		SettingsRevision:  "settings-revision-1",
		ReadinessToken:    "readiness-token-1",
	})
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/activate", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestActivateWorkspaceRequestContractRequiresBarrierFields(t *testing.T) {
	typeOfRequest := reflect.TypeOf(ActivateWorkspaceRequest{})
	for _, fieldName := range []string{"WorkspaceRevision", "SettingsRevision", "ReadinessToken"} {
		field, ok := typeOfRequest.FieldByName(fieldName)
		require.Truef(t, ok, "missing request field %s", fieldName)
		require.Equal(t, "required", field.Tag.Get("binding"), "%s must be required by Gin validation", fieldName)
		jsonName := field.Tag.Get("json")
		require.NotContains(t, jsonName, "omitempty", "%s must be present in the public JSON contract", fieldName)
	}
}

func TestStudioV2WorkspaceActivationHandler_MissingBodyReturnsTyped422(t *testing.T) {
	gin.SetMode(gin.TestMode)
	activator := &stubWorkspaceActivator{}
	handler := NewStudioV2WorkspaceActivationHandler(activator)
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/activate", http.NoBody)
	require.NoError(t, err)
	ctx.Request = req

	handler.Activate(ctx)

	require.Equal(t, http.StatusUnprocessableEntity, response.Code)
	require.Contains(t, response.Body.String(), modbusshare.ErrCodeSaveIncomplete)
	require.False(t, activator.called)
}

func TestStudioV2WorkspaceActivationHandler_MissingBarrierFieldReturnsTyped400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	activator := &stubWorkspaceActivator{}
	handler := NewStudioV2WorkspaceActivationHandler(activator)
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/activate", bytes.NewBufferString(`{"workspace_revision":"workspace-revision-1","settings_revision":"settings-revision-1"}`))
	require.NoError(t, err)
	ctx.Request = req
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.Activate(ctx)

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Contains(t, response.Body.String(), ErrCodeActivationRequestInvalid)
	require.False(t, activator.called)
}
