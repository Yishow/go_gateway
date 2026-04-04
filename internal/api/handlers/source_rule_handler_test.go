package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSourceRuleRouter(t *testing.T) (*gin.Engine, *point.Service) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	require.NoError(t, seedSourceRuleDevice(context.Background(), deviceRepo, "device-1"))

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, pointSvc, nil)
	handler := NewSourceRuleHandler(ruleSvc)

	router := gin.Default()
	router.GET("/datalink/source-rules", handler.List)
	router.GET("/datalink/source-rules/:id", handler.Get)
	router.POST("/datalink/source-rules", handler.Create)
	router.PUT("/datalink/source-rules/:id", handler.Update)
	router.DELETE("/datalink/source-rules/:id", handler.Delete)
	router.POST("/datalink/source-rules/:id/enable", handler.Enable)
	router.POST("/datalink/source-rules/:id/disable", handler.Disable)

	return router, pointSvc
}

func TestSourceRuleHandler_CreateAndList(t *testing.T) {
	router, pointSvc := setupSourceRuleRouter(t)

	body, err := json.Marshal(sourcerule.CreateRuleRequest{
		ID:           "rule-1",
		DeviceID:     "device-1",
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/datalink/source-rules", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)

	var createPayload map[string]interface{}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &createPayload))
	assert.True(t, createPayload["success"].(bool))

	listReq, err := http.NewRequest(http.MethodGet, "/datalink/source-rules?device_id=device-1", nil)
	require.NoError(t, err)
	listResp := httptest.NewRecorder()
	router.ServeHTTP(listResp, listReq)

	require.Equal(t, http.StatusOK, listResp.Code)

	var listPayload map[string]interface{}
	require.NoError(t, json.Unmarshal(listResp.Body.Bytes(), &listPayload))
	items := listPayload["data"].([]interface{})
	require.Len(t, items, 1)

	deviceID := "device-1"
	points, err := pointSvc.List(context.Background(), point.ListFilter{DeviceID: &deviceID})
	require.NoError(t, err)
	assert.Len(t, points, 2)
}

func TestSourceRuleHandler_Disable(t *testing.T) {
	router, pointSvc := setupSourceRuleRouter(t)

	createBody, err := json.Marshal(sourcerule.CreateRuleRequest{
		ID:           "rule-1",
		DeviceID:     "device-1",
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	createReq, err := http.NewRequest(http.MethodPost, "/datalink/source-rules", bytes.NewBuffer(createBody))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	router.ServeHTTP(createResp, createReq)
	require.Equal(t, http.StatusCreated, createResp.Code)

	disableReq, err := http.NewRequest(http.MethodPost, "/datalink/source-rules/rule-1/disable", nil)
	require.NoError(t, err)
	disableResp := httptest.NewRecorder()
	router.ServeHTTP(disableResp, disableReq)
	require.Equal(t, http.StatusOK, disableResp.Code)

	var disablePayload map[string]interface{}
	require.NoError(t, json.Unmarshal(disableResp.Body.Bytes(), &disablePayload))
	ruleData := disablePayload["data"].(map[string]interface{})
	assert.False(t, ruleData["enabled"].(bool))

	deviceID := "device-1"
	points, err := pointSvc.List(context.Background(), point.ListFilter{DeviceID: &deviceID})
	require.NoError(t, err)
	require.Len(t, points, 2)
	assert.False(t, points[0].Enabled)
	assert.False(t, points[1].Enabled)
}

func TestSourceRuleHandler_Get_NotFound(t *testing.T) {
	router, _ := setupSourceRuleRouter(t)

	req, err := http.NewRequest(http.MethodGet, "/datalink/source-rules/missing-rule", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestSourceRuleHandler_Update_ClearsOptionalConversionSettings(t *testing.T) {
	router, _ := setupSourceRuleRouter(t)

	targetType := schema.DataTypeFloat64
	multiplier := 0.25
	offset := -5.0
	createBody, err := json.Marshal(sourcerule.CreateRuleRequest{
		ID:              "rule-clear-conversion",
		DeviceID:        "device-1",
		StartAddress:    "40001",
		Count:           1,
		DataType:        schema.DataTypeUint16,
		NamingPrefix:    "SRC",
		Enabled:         true,
		TargetDataType:  &targetType,
		ScaleMultiplier: &multiplier,
		ScaleOffset:     &offset,
	})
	require.NoError(t, err)

	createReq, err := http.NewRequest(http.MethodPost, "/datalink/source-rules", bytes.NewBuffer(createBody))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	router.ServeHTTP(createResp, createReq)
	require.Equal(t, http.StatusCreated, createResp.Code)

	updateReq, err := http.NewRequest(http.MethodPut, "/datalink/source-rules/rule-clear-conversion", bytes.NewBufferString(`{"target_data_type":null,"scale_multiplier":null,"scale_offset":null}`))
	require.NoError(t, err)
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	router.ServeHTTP(updateResp, updateReq)
	require.Equal(t, http.StatusOK, updateResp.Code)

	var updatePayload map[string]any
	require.NoError(t, json.Unmarshal(updateResp.Body.Bytes(), &updatePayload))
	data := updatePayload["data"].(map[string]any)
	_, hasTargetDataType := data["target_data_type"]
	_, hasScaleMultiplier := data["scale_multiplier"]
	_, hasScaleOffset := data["scale_offset"]
	assert.False(t, hasTargetDataType)
	assert.False(t, hasScaleMultiplier)
	assert.False(t, hasScaleOffset)

	getReq, err := http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-clear-conversion", nil)
	require.NoError(t, err)
	getResp := httptest.NewRecorder()
	router.ServeHTTP(getResp, getReq)
	require.Equal(t, http.StatusOK, getResp.Code)

	var getPayload map[string]any
	require.NoError(t, json.Unmarshal(getResp.Body.Bytes(), &getPayload))
	getData := getPayload["data"].(map[string]any)
	_, hasPersistedTargetDataType := getData["target_data_type"]
	_, hasPersistedScaleMultiplier := getData["scale_multiplier"]
	_, hasPersistedScaleOffset := getData["scale_offset"]
	assert.False(t, hasPersistedTargetDataType)
	assert.False(t, hasPersistedScaleMultiplier)
	assert.False(t, hasPersistedScaleOffset)
}

func seedSourceRuleDevice(ctx context.Context, repo *device.MemoryRepository, id string) error {
	return repo.Create(ctx, &schema.Device{
		ID:               id,
		Name:             "Mixer PLC",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1}`,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	})
}
