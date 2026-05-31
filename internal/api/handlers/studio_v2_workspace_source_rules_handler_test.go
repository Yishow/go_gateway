package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func seedWorkspaceRuleDevice(t *testing.T, repo *device.MemoryRepository, id string) {
	t.Helper()

	require.NoError(t, repo.Create(context.Background(), &schema.Device{
		ID:               id,
		Name:             id,
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}))
}

func TestStudioV2WorkspaceSourceRulesHandler_UpdateRejectsOwnershipMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	deviceRepo := device.NewMemoryRepository()
	seedWorkspaceRuleDevice(t, deviceRepo, "dev-A")
	seedWorkspaceRuleDevice(t, deviceRepo, "dev-B")

	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, pointSvc, nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	_, err := workspaceSvc.AttachDevice(context.Background(), "dev-A")
	require.NoError(t, err)
	_, err = workspaceSvc.AttachDevice(context.Background(), "dev-B")
	require.NoError(t, err)

	_, err = ruleSvc.Create(context.Background(), sourcerule.CreateRuleRequest{
		ID:           "rule-A",
		DeviceID:     "dev-A",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "A_",
		Enabled:      true,
	})
	require.NoError(t, err)

	handler := NewStudioV2WorkspaceSourceRulesHandler(workspaceSvc, deviceSvc, ruleSvc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/datalink/studio-v2/workspace/source-rules/rule-A", bytes.NewBufferString(`{
		"device_id":"dev-B",
		"start_address":"40005",
		"count":1,
		"data_type":"int16",
		"naming_prefix":"A_"
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "rule-A"}}

	handler.Update(c)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "ownership mismatch")
}

func TestStudioV2WorkspaceSourceRulesHandler_CreateReturnsValidationErrorMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	deviceRepo := device.NewMemoryRepository()
	seedWorkspaceRuleDevice(t, deviceRepo, "dev-A")

	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, pointSvc, nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	_, err := workspaceSvc.AttachDevice(context.Background(), "dev-A")
	require.NoError(t, err)

	handler := NewStudioV2WorkspaceSourceRulesHandler(workspaceSvc, deviceSvc, ruleSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/source-rules", strings.NewReader(`{
		"id":"rule-A",
		"device_id":"dev-A",
		"start_address":"",
		"count":1,
		"data_type":"int16",
		"naming_prefix":"A_",
		"enabled":true
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Create(c)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "start_address is required")
}
