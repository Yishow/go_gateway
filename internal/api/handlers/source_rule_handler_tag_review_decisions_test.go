package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

type tagReviewDecisionResponse struct {
	Success bool                               `json:"success"`
	Data    schema.SourceRuleTagReviewDecision `json:"data"`
}

type tagReviewDecisionListResponse struct {
	Success bool                                 `json:"success"`
	Data    []schema.SourceRuleTagReviewDecision `json:"data"`
}

func TestSourceRuleHandler_TagReviewDecisions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		ruleID         string
		buildRequest   func(t *testing.T, fixture *sourceRuleCandidatesFixture, candidate schema.SourceRuleTagCandidate) sourcerule.UpsertTagReviewDecisionRequest
		assertDecision func(t *testing.T, decision schema.SourceRuleTagReviewDecision, candidate schema.SourceRuleTagCandidate)
	}{
		{
			name:   "rename",
			ruleID: "rule-review-rename",
			buildRequest: func(_ *testing.T, _ *sourceRuleCandidatesFixture, candidate schema.SourceRuleTagCandidate) sourcerule.UpsertTagReviewDecisionRequest {
				return sourcerule.UpsertTagReviewDecisionRequest{
					CandidateID: candidate.ID,
					Action:      schema.SourceRuleTagReviewDecisionActionRename,
					TagKey:      "RENAMED_40001",
				}
			},
			assertDecision: func(t *testing.T, decision schema.SourceRuleTagReviewDecision, candidate schema.SourceRuleTagCandidate) {
				t.Helper()
				require.Equal(t, candidate.ID, decision.CandidateID)
				require.Equal(t, schema.SourceRuleTagReviewDecisionActionRename, decision.Action)
				require.Equal(t, "RENAMED_40001", decision.TagKey)
				require.Nil(t, decision.OverrideTagID)
				require.False(t, decision.Stale)
			},
		},
		{
			name:   "skip",
			ruleID: "rule-review-skip",
			buildRequest: func(_ *testing.T, _ *sourceRuleCandidatesFixture, candidate schema.SourceRuleTagCandidate) sourcerule.UpsertTagReviewDecisionRequest {
				return sourcerule.UpsertTagReviewDecisionRequest{
					CandidateID: candidate.ID,
					Action:      schema.SourceRuleTagReviewDecisionActionSkip,
				}
			},
			assertDecision: func(t *testing.T, decision schema.SourceRuleTagReviewDecision, candidate schema.SourceRuleTagCandidate) {
				t.Helper()
				require.Equal(t, candidate.ID, decision.CandidateID)
				require.Equal(t, schema.SourceRuleTagReviewDecisionActionSkip, decision.Action)
				require.Empty(t, decision.TagKey)
				require.Nil(t, decision.OverrideTagID)
				require.False(t, decision.Stale)
			},
		},
		{
			name:   "override",
			ruleID: "rule-review-override",
			buildRequest: func(t *testing.T, fixture *sourceRuleCandidatesFixture, candidate schema.SourceRuleTagCandidate) sourcerule.UpsertTagReviewDecisionRequest {
				t.Helper()
				overrideTag, err := fixture.tagSvc.Create(context.Background(), tag.CreateTagRequest{
					Key:         "OVERRIDE_40001",
					DisplayName: "Override tag",
					DataType:    candidate.DataType,
				})
				require.NoError(t, err)
				return sourcerule.UpsertTagReviewDecisionRequest{
					CandidateID:   candidate.ID,
					Action:        schema.SourceRuleTagReviewDecisionActionOverride,
					OverrideTagID: overrideTag.ID,
				}
			},
			assertDecision: func(t *testing.T, decision schema.SourceRuleTagReviewDecision, candidate schema.SourceRuleTagCandidate) {
				t.Helper()
				require.Equal(t, candidate.ID, decision.CandidateID)
				require.Equal(t, schema.SourceRuleTagReviewDecisionActionOverride, decision.Action)
				require.Equal(t, "OVERRIDE_40001", decision.TagKey)
				require.NotNil(t, decision.OverrideTagID)
				require.False(t, decision.Stale)
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			fixture := setupSourceRuleCandidatesFixture(t)
			candidate := createRuleCandidateForDecisionTest(t, fixture, testCase.ruleID)
			reqBody, err := json.Marshal(testCase.buildRequest(t, fixture, candidate))
			require.NoError(t, err)

			upsertReq, err := http.NewRequestWithContext(
				context.Background(), http.MethodPost,
				fmt.Sprintf("/datalink/source-rules/%s/tag-review-decisions", testCase.ruleID),
				bytes.NewBuffer(reqBody),
			)
			require.NoError(t, err)
			upsertReq.Header.Set("Content-Type", "application/json")

			upsertResp := httptest.NewRecorder()
			fixture.router.ServeHTTP(upsertResp, upsertReq)
			require.Equal(t, http.StatusOK, upsertResp.Code)

			var upsertPayload tagReviewDecisionResponse
			require.NoError(t, json.Unmarshal(upsertResp.Body.Bytes(), &upsertPayload))
			require.True(t, upsertPayload.Success)
			testCase.assertDecision(t, upsertPayload.Data, candidate)

			listReq, err := http.NewRequestWithContext(
				t.Context(), http.MethodGet,
				fmt.Sprintf("/datalink/source-rules/%s/tag-review-decisions", testCase.ruleID),
				http.NoBody,
			)
			require.NoError(t, err)

			listResp := httptest.NewRecorder()
			fixture.router.ServeHTTP(listResp, listReq)
			require.Equal(t, http.StatusOK, listResp.Code)

			var listPayload tagReviewDecisionListResponse
			require.NoError(t, json.Unmarshal(listResp.Body.Bytes(), &listPayload))
			require.True(t, listPayload.Success)
			require.Len(t, listPayload.Data, 1)
			testCase.assertDecision(t, listPayload.Data[0], candidate)
		})
	}
}

func createRuleCandidateForDecisionTest(t *testing.T, fixture *sourceRuleCandidatesFixture, ruleID string) schema.SourceRuleTagCandidate {
	t.Helper()

	candidates := createRuleCandidatesForDecisionTest(t, fixture, ruleID, "40001", 1)
	require.Len(t, candidates, 1)
	return candidates[0]
}

func createRuleCandidatesForDecisionTest(
	t *testing.T,
	fixture *sourceRuleCandidatesFixture,
	ruleID string,
	startAddress string,
	count int,
) []schema.SourceRuleTagCandidate {
	t.Helper()

	createBody, err := json.Marshal(sourcerule.CreateRuleRequest{
		ID:           ruleID,
		DeviceID:     "device-1",
		StartAddress: startAddress,
		Count:        count,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	createReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/datalink/source-rules", bytes.NewBuffer(createBody))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "application/json")

	createResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(createResp, createReq)
	require.Equal(t, http.StatusCreated, createResp.Code)

	rule, err := fixture.repo.GetByID(context.Background(), ruleID)
	require.NoError(t, err)

	snapshots, err := fixture.repo.ListCandidateSnapshots(context.Background(), rule.ID, rule.RevisionID)
	require.NoError(t, err)

	for _, snapshot := range snapshots {
		if snapshot == nil || snapshot.CandidateType != schema.SourceRuleCandidateTypeTags {
			continue
		}

		var payload handlerTagSnapshotPayload
		require.NoError(t, json.Unmarshal([]byte(snapshot.Payload), &payload))
		require.Len(t, payload.Candidates, count)
		return payload.Candidates
	}

	t.Fatalf("expected tag candidate snapshot for %s", ruleID)
	return nil
}
