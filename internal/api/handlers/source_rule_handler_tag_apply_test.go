package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tagApplyResponse struct {
	Success bool `json:"success"`
	Data    struct {
		SourceRuleID string `json:"source_rule_id"`
		RevisionID   string `json:"revision_id"`
		Results      []struct {
			CandidateID string `json:"candidate_id"`
			Status      string `json:"status"`
			TagID       string `json:"tag_id"`
			MappingID   string `json:"mapping_id"`
			Error       string `json:"error"`
		} `json:"results"`
	} `json:"data"`
}

func TestSourceRuleHandler_ApplyTags_AppliesApprovedRenameDecision(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	candidate := createRuleCandidateForDecisionTest(t, fixture, "rule-apply-handler")

	decisionBody, err := json.Marshal(sourcerule.UpsertTagReviewDecisionRequest{
		CandidateID: candidate.ID,
		Action:      schema.SourceRuleTagReviewDecisionActionRename,
		TagKey:      "factory.temperature.handler",
	})
	require.NoError(t, err)

	decisionReq, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-apply-handler/tag-review-decisions",
		bytes.NewBuffer(decisionBody),
	)
	require.NoError(t, err)
	decisionReq.Header.Set("Content-Type", "application/json")
	decisionResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(decisionResp, decisionReq)
	require.Equal(t, http.StatusOK, decisionResp.Code)

	rule, err := fixture.repo.GetByID(decisionReq.Context(), "rule-apply-handler")
	require.NoError(t, err)

	applyBody, err := json.Marshal(sourcerule.ApplyTagCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{candidate.ID},
	})
	require.NoError(t, err)

	applyReq, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-apply-handler/tags/apply",
		bytes.NewBuffer(applyBody),
	)
	require.NoError(t, err)
	applyReq.Header.Set("Content-Type", "application/json")
	applyResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(applyResp, applyReq)
	require.Equal(t, http.StatusOK, applyResp.Code)

	var payload tagApplyResponse
	require.NoError(t, json.Unmarshal(applyResp.Body.Bytes(), &payload))
	require.True(t, payload.Success)
	require.Len(t, payload.Data.Results, 1)
	assert.Equal(t, candidate.ID, payload.Data.Results[0].CandidateID)
	assert.Equal(t, "applied", payload.Data.Results[0].Status)

	tagRecord, err := fixture.tagSvc.GetByID(applyReq.Context(), payload.Data.Results[0].TagID)
	require.NoError(t, err)
	assert.Equal(t, "factory.temperature.handler", tagRecord.Key)

	mappingRecord, err := fixture.mappingSvc.GetByID(applyReq.Context(), payload.Data.Results[0].MappingID)
	require.NoError(t, err)
	assert.Equal(t, tagRecord.ID, mappingRecord.TagID)

	links, err := fixture.repo.ListLinks(applyReq.Context(), "rule-apply-handler")
	require.NoError(t, err)
	require.Len(t, links, 1)
	require.NotNil(t, links[0].TagID)
	require.NotNil(t, links[0].MappingID)
	assert.Equal(t, tagRecord.ID, *links[0].TagID)
	assert.Equal(t, mappingRecord.ID, *links[0].MappingID)
}

func TestSourceRuleHandler_ApplyTags_RejectsRevisionMismatch(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	candidate := createRuleCandidateForDecisionTest(t, fixture, "rule-apply-conflict")

	applyBody, err := json.Marshal(sourcerule.ApplyTagCandidatesRequest{
		RevisionID:   "stale-revision",
		CandidateIDs: []string{candidate.ID},
	})
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/datalink/source-rules/%s/tags/apply", "rule-apply-conflict"),
		bytes.NewBuffer(applyBody),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusConflict, resp.Code)
}

func TestSourceRuleHandler_ApplyTags_ReturnsPerCandidateResultsForPartialSuccess(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	candidates := createRuleCandidatesForDecisionTest(t, fixture, "rule-apply-partial", "40001", 2)
	otherCandidates := createRuleCandidatesForDecisionTest(t, fixture, "rule-apply-partial-existing", "40101", 1)

	otherLinks, err := fixture.repo.ListLinks(context.Background(), "rule-apply-partial-existing")
	require.NoError(t, err)
	require.Len(t, otherLinks, 1)

	existingTag, err := fixture.tagSvc.Create(context.Background(), tag.CreateTagRequest{
		Key:         "factory.partial.shared",
		DisplayName: "Partial Shared",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)

	enabled := true
	status := schema.MappingStatusActive
	_, err = fixture.mappingSvc.Create(context.Background(), mapping.CreateMappingRequest{
		PointID:           otherLinks[0].PointID,
		TagID:             existingTag.ID,
		Enabled:           &enabled,
		TransformPipeline: otherCandidates[0].TransformPipeline,
		Status:            &status,
	})
	require.NoError(t, err)

	decisionBody, err := json.Marshal(sourcerule.UpsertTagReviewDecisionRequest{
		CandidateID:   candidates[1].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: existingTag.ID,
	})
	require.NoError(t, err)

	decisionReq, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-apply-partial/tag-review-decisions",
		bytes.NewBuffer(decisionBody),
	)
	require.NoError(t, err)
	decisionReq.Header.Set("Content-Type", "application/json")
	decisionResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(decisionResp, decisionReq)
	require.Equal(t, http.StatusOK, decisionResp.Code)

	rule, err := fixture.repo.GetByID(context.Background(), "rule-apply-partial")
	require.NoError(t, err)

	applyBody, err := json.Marshal(sourcerule.ApplyTagCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{candidates[0].ID, candidates[1].ID},
	})
	require.NoError(t, err)

	applyReq, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-apply-partial/tags/apply",
		bytes.NewBuffer(applyBody),
	)
	require.NoError(t, err)
	applyReq.Header.Set("Content-Type", "application/json")
	applyResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(applyResp, applyReq)
	require.Equal(t, http.StatusOK, applyResp.Code)

	var payload tagApplyResponse
	require.NoError(t, json.Unmarshal(applyResp.Body.Bytes(), &payload))
	require.True(t, payload.Success)
	require.Len(t, payload.Data.Results, 2)

	assert.Equal(t, candidates[0].ID, payload.Data.Results[0].CandidateID)
	assert.Equal(t, "applied", payload.Data.Results[0].Status)
	assert.NotEmpty(t, payload.Data.Results[0].TagID)
	assert.NotEmpty(t, payload.Data.Results[0].MappingID)
	assert.Empty(t, payload.Data.Results[0].Error)

	assert.Equal(t, candidates[1].ID, payload.Data.Results[1].CandidateID)
	assert.Equal(t, "failed", payload.Data.Results[1].Status)
	assert.Contains(t, payload.Data.Results[1].Error, "已綁定其他 point")
	assert.Empty(t, payload.Data.Results[1].TagID)
	assert.Empty(t, payload.Data.Results[1].MappingID)

	candidatesReq, err := http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-apply-partial/candidates", nil)
	require.NoError(t, err)
	candidatesResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(candidatesResp, candidatesReq)
	require.Equal(t, http.StatusOK, candidatesResp.Code)

	var candidatesPayload map[string]any
	require.NoError(t, json.Unmarshal(candidatesResp.Body.Bytes(), &candidatesPayload))
	data := candidatesPayload["data"].(map[string]any)
	tagCandidates := data["tags"].(map[string]any)["candidates"].([]any)
	require.Len(t, tagCandidates, 2)

	refreshedByID := make(map[string]map[string]any, len(tagCandidates))
	for _, candidate := range tagCandidates {
		record := candidate.(map[string]any)
		refreshedByID[record["id"].(string)] = record
	}

	require.Contains(t, refreshedByID, candidates[0].ID)
	assert.Equal(t, payload.Data.Results[0].TagID, refreshedByID[candidates[0].ID]["tag_id"])
	assert.Equal(t, payload.Data.Results[0].MappingID, refreshedByID[candidates[0].ID]["mapping_id"])

	require.Contains(t, refreshedByID, candidates[1].ID)
	assert.Nil(t, refreshedByID[candidates[1].ID]["tag_id"])
	assert.Nil(t, refreshedByID[candidates[1].ID]["mapping_id"])
}

func TestSourceRuleHandler_ApplyTags_RefreshesCandidateViewAfterApply(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	candidate := createRuleCandidateForDecisionTest(t, fixture, "rule-apply-refresh")

	rule, err := fixture.repo.GetByID(context.Background(), "rule-apply-refresh")
	require.NoError(t, err)

	applyBody, err := json.Marshal(sourcerule.ApplyTagCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{candidate.ID},
	})
	require.NoError(t, err)

	applyReq, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-apply-refresh/tags/apply",
		bytes.NewBuffer(applyBody),
	)
	require.NoError(t, err)
	applyReq.Header.Set("Content-Type", "application/json")
	applyResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(applyResp, applyReq)
	require.Equal(t, http.StatusOK, applyResp.Code)

	candidatesReq, err := http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-apply-refresh/candidates", nil)
	require.NoError(t, err)
	candidatesResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(candidatesResp, candidatesReq)
	require.Equal(t, http.StatusOK, candidatesResp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(candidatesResp.Body.Bytes(), &payload))
	data := payload["data"].(map[string]any)

	tags := data["tags"].(map[string]any)
	tagCandidates := tags["candidates"].([]any)
	require.Len(t, tagCandidates, 1)

	refreshed := tagCandidates[0].(map[string]any)
	assert.Equal(t, candidate.ID, refreshed["id"])
	assert.NotEmpty(t, refreshed["tag_id"])
	assert.NotEmpty(t, refreshed["mapping_id"])

	databaseOutputs := data["database_outputs"].(map[string]any)
	assert.Equal(t, "ready", databaseOutputs["status"])
	assert.Empty(t, databaseOutputs["reason"])
	require.Len(t, databaseOutputs["candidates"].([]any), 1)

	localModbusOutputs := data["local_modbus_outputs"].(map[string]any)
	assert.Equal(t, "deferred", localModbusOutputs["status"])
	assert.NotEmpty(t, localModbusOutputs["reason"])
}
