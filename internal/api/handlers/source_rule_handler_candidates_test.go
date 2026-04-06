package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

type sourceRuleCandidatesFixture struct {
	router     *gin.Engine
	repo       *sourcerule.MemoryRepository
	tagSvc     *tag.Service
	mappingSvc *mapping.Service
}

type handlerTagSnapshotPayload struct {
	Candidates []schema.SourceRuleTagCandidate `json:"candidates"`
}

func setupSourceRuleCandidatesFixture(t *testing.T) *sourceRuleCandidatesFixture {
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
	router.GET("/datalink/source-rules/:id/tag-review-decisions", handler.ListTagReviewDecisions)
	router.POST("/datalink/source-rules/:id/tag-review-decisions", handler.UpsertTagReviewDecision)
	return &sourceRuleCandidatesFixture{
		router:     router,
		repo:       repo,
		tagSvc:     tagSvc,
		mappingSvc: mappingSvc,
	}
}

func setupSourceRuleCandidatesRouter(t *testing.T) *gin.Engine {
	t.Helper()
	return setupSourceRuleCandidatesFixture(t).router
}

func seedAppliedRuleManagedLink(t *testing.T, fixture *sourceRuleCandidatesFixture, ruleID string) {
	t.Helper()

	ctx := context.Background()
	rule, err := fixture.repo.GetByID(ctx, ruleID)
	require.NoError(t, err)

	links, err := fixture.repo.ListLinks(ctx, ruleID)
	require.NoError(t, err)
	require.Len(t, links, 1)

	snapshots, err := fixture.repo.ListCandidateSnapshots(ctx, ruleID, rule.RevisionID)
	require.NoError(t, err)

	var candidate schema.SourceRuleTagCandidate
	for _, snapshot := range snapshots {
		if snapshot.CandidateType != schema.SourceRuleCandidateTypeTags {
			continue
		}
		var payload handlerTagSnapshotPayload
		require.NoError(t, json.Unmarshal([]byte(snapshot.Payload), &payload))
		require.Len(t, payload.Candidates, 1)
		candidate = payload.Candidates[0]
		break
	}
	require.NotEmpty(t, candidate.ID)

	tagRecord, err := fixture.tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         candidate.TagKey,
		DisplayName: candidate.DisplayName,
		DataType:    candidate.DataType,
		Labels: map[string]string{
			"source":              "source-rule",
			"source_rule_id":      rule.ID,
			"source_rule_address": strings.ToUpper(strings.TrimSpace(candidate.Address)),
		},
	})
	require.NoError(t, err)

	enabled := rule.Enabled
	status := schema.MappingStatusActive
	if !enabled {
		status = schema.MappingStatusDraft
	}
	mappingRecord, err := fixture.mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID:           links[0].PointID,
		TagID:             tagRecord.ID,
		Enabled:           &enabled,
		TransformPipeline: candidate.TransformPipeline,
		Status:            &status,
	})
	require.NoError(t, err)

	links[0].TagID = stringPtr(tagRecord.ID)
	links[0].MappingID = stringPtr(mappingRecord.ID)
	require.NoError(t, fixture.repo.DeleteLinks(ctx, ruleID))
	require.NoError(t, fixture.repo.CreateLinks(ctx, links))
}

func stringPtr(value string) *string {
	return &value
}

func TestSourceRuleHandler_Candidates_ReturnsCurrentRevisionSnapshot(t *testing.T) {
	fixture := setupSourceRuleCandidatesFixture(t)
	router := fixture.router

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

	seedAppliedRuleManagedLink(t, fixture, "rule-candidates")

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
	assert.Equal(t, "draft", candidate["status"])
	assert.NotEmpty(t, candidate["proposed_signature"])
	assert.Nil(t, candidate["last_applied_signature"])
}
