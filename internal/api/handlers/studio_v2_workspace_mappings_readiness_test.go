package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStudioV2WorkspaceMappingsHandler_DeleteDefersLiveApplyWhenReadinessBlocks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	deviceRepo := device.NewMemoryRepository()
	lastTestSuccess := true
	require.NoError(t, deviceRepo.Create(context.Background(), &schema.Device{
		ID:               "dev-A",
		Name:             "dev-A",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}`,
		LastTestSuccess:  &lastTestSuccess,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}))

	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	deviceSvc := device.NewService(deviceRepo, nil)
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, pointSvc, nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository()).WithReadinessServices(deviceSvc, ruleSvc, nil, nil)
	_, err := workspaceSvc.AttachDevice(context.Background(), "dev-A")
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

	handler := NewStudioV2WorkspaceMappingsHandler(workspaceSvc, deviceSvc, ruleSvc, pointSvc, tagSvc, mappingSvc)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/mappings", strings.NewReader(`{
		"rule_id":"rule-A",
		"address":"40001",
		"tag_key":"line.a.temp",
		"display_name":"Line A Temp",
		"unit":"C",
		"target_type":"float64",
		"scale":1,
		"offset":0,
		"enabled":true
	}`))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	createCtx, _ := gin.CreateTestContext(createResp)
	createCtx.Request = createReq
	handler.Create(createCtx)
	require.Equal(t, http.StatusCreated, createResp.Code, createResp.Body.String())

	createBody := decodeWorkspaceMappingBody(t, createResp)
	mappingID := createBody["data"].(map[string]any)["id"].(string)
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/datalink/studio-v2/workspace/mappings/"+mappingID, nil)
	deleteResp := httptest.NewRecorder()
	deleteCtx, _ := gin.CreateTestContext(deleteResp)
	deleteCtx.Request = deleteReq
	deleteCtx.Params = gin.Params{{Key: "id", Value: mappingID}}

	handler.Delete(deleteCtx)

	require.Equal(t, http.StatusOK, deleteResp.Code, deleteResp.Body.String())
	require.Contains(t, deleteResp.Body.String(), `"runtime_apply_status":"deferred"`)
	require.Contains(t, deleteResp.Body.String(), "tag-missing")
}
