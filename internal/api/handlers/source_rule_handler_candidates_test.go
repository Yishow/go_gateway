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
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSourceRuleCandidatesFixture(t *testing.T) (*gin.Engine, *sourcerule.MemoryRepository) {
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
	repo := sourcerule.NewMemoryRepository()
	ruleSvc := sourcerule.NewService(repo, deviceSvc, pointSvc, nil)
	ruleSvc.SetTagMappingServices(tagSvc, mappingSvc)
	handler := NewSourceRuleHandler(ruleSvc)

	router := gin.Default()
	router.POST("/datalink/source-rules", handler.Create)
	router.PUT("/datalink/source-rules/:id", handler.Update)
	router.GET("/datalink/source-rules/:id/candidates", handler.Candidates)
	router.POST("/datalink/source-rules/:id/candidates/recompute", handler.RecomputeCandidates)
	return router, repo
}

func setupSourceRuleCandidatesRouter(t *testing.T) *gin.Engine {
	t.Helper()
	router, _ := setupSourceRuleCandidatesFixture(t)
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

func TestSourceRuleHandler_RecomputeCandidates_UpdatesPersistedCurrentSnapshot(t *testing.T) {
	router, repo := setupSourceRuleCandidatesFixture(t)

	createBody, err := json.Marshal(sourcerule.CreateRuleRequest{
		ID:           "rule-recompute-persist",
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

	rule, err := repo.GetByID(context.Background(), "rule-recompute-persist")
	require.NoError(t, err)
	require.NoError(t, repo.ReplaceCandidateSnapshots(context.Background(), []*schema.SourceRuleCandidateSnapshot{
		{
			SourceRuleID:  rule.ID,
			RevisionID:    rule.RevisionID,
			CandidateType: schema.SourceRuleCandidateTypeTags,
			Payload:       `{"candidates":[]}`,
			Status:        schema.SourceRuleCandidateStatusReady,
			GeneratedAt:   time.Now().UTC(),
		},
		{
			SourceRuleID:  rule.ID,
			RevisionID:    rule.RevisionID,
			CandidateType: schema.SourceRuleCandidateTypeDatabaseOutputs,
			Payload:       `{"candidates":[]}`,
			Status:        schema.SourceRuleCandidateStatusDeferred,
			Reason:        "stale database snapshot",
			GeneratedAt:   time.Now().UTC(),
		},
		{
			SourceRuleID:  rule.ID,
			RevisionID:    rule.RevisionID,
			CandidateType: schema.SourceRuleCandidateTypeLocalModbusOutputs,
			Payload:       `{"candidates":[]}`,
			Status:        schema.SourceRuleCandidateStatusDeferred,
			Reason:        "stale local modbus snapshot",
			GeneratedAt:   time.Now().UTC(),
		},
	}))

	beforeReq, err := http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-recompute-persist/candidates", nil)
	require.NoError(t, err)
	beforeResp := httptest.NewRecorder()
	router.ServeHTTP(beforeResp, beforeReq)
	require.Equal(t, http.StatusOK, beforeResp.Code)

	var beforePayload map[string]any
	require.NoError(t, json.Unmarshal(beforeResp.Body.Bytes(), &beforePayload))
	beforeTags := beforePayload["data"].(map[string]any)["tags"].(map[string]any)
	assert.Empty(t, beforeTags["candidates"].([]any))

	recomputeReq, err := http.NewRequest(http.MethodPost, "/datalink/source-rules/rule-recompute-persist/candidates/recompute", nil)
	require.NoError(t, err)
	recomputeResp := httptest.NewRecorder()
	router.ServeHTTP(recomputeResp, recomputeReq)
	require.Equal(t, http.StatusOK, recomputeResp.Code)

	afterReq, err := http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-recompute-persist/candidates", nil)
	require.NoError(t, err)
	afterResp := httptest.NewRecorder()
	router.ServeHTTP(afterResp, afterReq)
	require.Equal(t, http.StatusOK, afterResp.Code)

	var afterPayload map[string]any
	require.NoError(t, json.Unmarshal(afterResp.Body.Bytes(), &afterPayload))
	afterData := afterPayload["data"].(map[string]any)
	assert.Equal(t, rule.RevisionID, afterData["revision_id"])
	afterTags := afterData["tags"].(map[string]any)
	require.Len(t, afterTags["candidates"].([]any), 1)
}
