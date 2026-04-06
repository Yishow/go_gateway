package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

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
