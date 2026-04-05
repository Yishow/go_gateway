package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestSourceRuleHandler_RecomputeCandidates_PreservesOutOfSyncReviewState(t *testing.T) {
	router := setupSourceRuleCandidatesRouter(t)

	createBody, err := json.Marshal(sourcerule.CreateRuleRequest{
		ID:           "rule-out-of-sync-review",
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

	updateReq, err := http.NewRequest(http.MethodPut, "/datalink/source-rules/rule-out-of-sync-review", bytes.NewBufferString(`{"scale_multiplier":2}`))
	require.NoError(t, err)
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	router.ServeHTTP(updateResp, updateReq)
	require.Equal(t, http.StatusOK, updateResp.Code)

	beforeReq, err := http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-out-of-sync-review/candidates", nil)
	require.NoError(t, err)
	beforeResp := httptest.NewRecorder()
	router.ServeHTTP(beforeResp, beforeReq)
	require.Equal(t, http.StatusOK, beforeResp.Code)

	var beforePayload map[string]any
	require.NoError(t, json.Unmarshal(beforeResp.Body.Bytes(), &beforePayload))
	beforeCandidate := beforePayload["data"].(map[string]any)["tags"].(map[string]any)["candidates"].([]any)[0].(map[string]any)
	require.Equal(t, "out_of_sync", beforeCandidate["status"])
	require.NotEmpty(t, beforeCandidate["blocking_reason"])
	proposedSignature := beforeCandidate["proposed_signature"]

	recomputeReq, err := http.NewRequest(http.MethodPost, "/datalink/source-rules/rule-out-of-sync-review/candidates/recompute", nil)
	require.NoError(t, err)
	recomputeResp := httptest.NewRecorder()
	router.ServeHTTP(recomputeResp, recomputeReq)
	require.Equal(t, http.StatusOK, recomputeResp.Code)

	afterReq, err := http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-out-of-sync-review/candidates", nil)
	require.NoError(t, err)
	afterResp := httptest.NewRecorder()
	router.ServeHTTP(afterResp, afterReq)
	require.Equal(t, http.StatusOK, afterResp.Code)

	var afterPayload map[string]any
	require.NoError(t, json.Unmarshal(afterResp.Body.Bytes(), &afterPayload))
	afterCandidate := afterPayload["data"].(map[string]any)["tags"].(map[string]any)["candidates"].([]any)[0].(map[string]any)
	assert.Equal(t, "out_of_sync", afterCandidate["status"])
	assert.Equal(t, proposedSignature, afterCandidate["proposed_signature"])
	assert.NotEmpty(t, afterCandidate["last_applied_signature"])
	assert.NotEmpty(t, afterCandidate["blocking_reason"])
}
