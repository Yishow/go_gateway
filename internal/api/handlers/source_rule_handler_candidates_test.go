package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSourceRuleCandidatesRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	require.NoError(t, seedSourceRuleDevice(context.Background(), deviceRepo, "device-1"))
	record, err := deviceRepo.GetByID(context.Background(), "device-1")
	require.NoError(t, err)
	success := true
	record.LastTestSuccess = &success
	record.LastTestError = ""
	require.NoError(t, deviceRepo.Update(context.Background(), record))

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, pointSvc, nil)
	ruleSvc.SetTagMappingServices(tagSvc, mappingSvc)
	handler := NewSourceRuleHandler(ruleSvc)

	router := gin.Default()
	router.POST("/datalink/source-rules", handler.Create)
	router.PUT("/datalink/source-rules/:id", handler.Update)
	router.GET("/datalink/source-rules/:id/candidates", handler.Candidates)
	router.POST("/datalink/source-rules/:id/candidates/recompute", handler.RecomputeCandidates)
	return router
}

func TestSourceRuleHandler_Candidates_ReturnsCurrentRevisionSnapshot(t *testing.T) {
	router := setupSourceRuleCandidatesRouter(t)

	createBody, err := json.Marshal(sourcerule.CreateRuleRequest{
		ID:           "rule-candidates",
		DeviceID:     "device-1",
		StartAddress: "40001",
		Count:        1,
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

	updateReq, err := http.NewRequest(http.MethodPut, "/datalink/source-rules/rule-candidates", bytes.NewBufferString(`{"scale_multiplier":2}`))
	require.NoError(t, err)
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	router.ServeHTTP(updateResp, updateReq)
	require.Equal(t, http.StatusOK, updateResp.Code)

	var updatePayload map[string]any
	require.NoError(t, json.Unmarshal(updateResp.Body.Bytes(), &updatePayload))
	updatedRule := updatePayload["data"].(map[string]any)

	req, err := http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-candidates/candidates", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	data := payload["data"].(map[string]any)
	assert.Equal(t, "rule-candidates", data["source_rule_id"])
	assert.Equal(t, updatedRule["revision_id"], data["revision_id"])

	tags := data["tags"].(map[string]any)
	assert.Equal(t, "ready", tags["status"])
	tagCandidates := tags["candidates"].([]any)
	require.Len(t, tagCandidates, 1)

	candidate := tagCandidates[0].(map[string]any)
	assert.NotEmpty(t, candidate["id"])
	assert.Equal(t, "out_of_sync", candidate["status"])
	assert.NotEmpty(t, candidate["proposed_signature"])
	assert.NotEmpty(t, candidate["last_applied_signature"])
	assert.NotEmpty(t, candidate["blocking_reason"])

	databaseOutputs := data["database_outputs"].(map[string]any)
	assert.Equal(t, "deferred", databaseOutputs["status"])
	assert.NotEmpty(t, databaseOutputs["reason"])

	localModbusOutputs := data["local_modbus_outputs"].(map[string]any)
	assert.Equal(t, "deferred", localModbusOutputs["status"])
	assert.NotEmpty(t, localModbusOutputs["reason"])
}

func TestSourceRuleHandler_RecomputeCandidates_ReturnsCurrentRevisionSnapshot(t *testing.T) {
	router := setupSourceRuleCandidatesRouter(t)

	createBody, err := json.Marshal(sourcerule.CreateRuleRequest{
		ID:           "rule-recompute",
		DeviceID:     "device-1",
		StartAddress: "40001",
		Count:        1,
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

	var createPayload map[string]any
	require.NoError(t, json.Unmarshal(createResp.Body.Bytes(), &createPayload))
	createdRule := createPayload["data"].(map[string]any)

	req, err := http.NewRequest(http.MethodPost, "/datalink/source-rules/rule-recompute/candidates/recompute", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	data := payload["data"].(map[string]any)
	assert.Equal(t, "rule-recompute", data["source_rule_id"])
	assert.Equal(t, createdRule["revision_id"], data["revision_id"])

	tags := data["tags"].(map[string]any)
	assert.Equal(t, "ready", tags["status"])
	tagCandidates := tags["candidates"].([]any)
	require.Len(t, tagCandidates, 1)

	candidate := tagCandidates[0].(map[string]any)
	assert.NotEmpty(t, candidate["id"])
	assert.Equal(t, "active", candidate["status"])
	assert.NotEmpty(t, candidate["proposed_signature"])
	assert.NotEmpty(t, candidate["last_applied_signature"])
}
