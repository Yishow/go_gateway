package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourceRuleHandler_ApplyDatabaseOutputs_DoesNotSilentlyRebindToNewScope(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	tagCandidates := createRuleCandidatesForDecisionTest(t, fixture, "rule-db-apply-rebind", "40001", 1)
	require.Len(t, tagCandidates, 1)

	overrideTag, err := fixture.tagSvc.Create(context.Background(), tag.CreateTagRequest{
		Key:         "factory.db.rebind",
		DisplayName: "DB Rebind",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)

	currentMappings := []*schema.DatabaseTargetMapping{
		{
			ID:          "db-map-rebind",
			TagID:       overrideTag.ID,
			ConnectorID: "connector-a",
			TableSchema: "public",
			TableName:   "measurements",
			ColumnName:  "line_a",
			WriteMode:   schema.DatabaseWriteModeInsert,
		},
	}
	fixture.ruleSvc.SetDatabaseTargetMappingReader(
		sourcerule.DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
			return currentMappings, nil
		}),
	)
	fixture.ruleSvc.SetDatabaseTargetConnectorValidator(
		sourcerule.DatabaseTargetConnectorValidatorFunc(func(context.Context, string) (*sourcerule.DatabaseTargetConnectorValidation, error) {
			return &sourcerule.DatabaseTargetConnectorValidation{Ready: true}, nil
		}),
	)

	decisionBody, err := json.Marshal(sourcerule.UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[0].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	decisionReq, err := http.NewRequestWithContext(
		context.Background(), http.MethodPost,
		"/datalink/source-rules/rule-db-apply-rebind/tag-review-decisions",
		bytes.NewBuffer(decisionBody),
	)
	require.NoError(t, err)
	decisionReq.Header.Set("Content-Type", "application/json")
	decisionResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(decisionResp, decisionReq)
	require.Equal(t, http.StatusOK, decisionResp.Code)

	rule, err := fixture.repo.GetByID(context.Background(), "rule-db-apply-rebind")
	require.NoError(t, err)

	originalCandidates := loadDatabaseOutputCandidatesForRule(t, fixture, "rule-db-apply-rebind")
	require.Len(t, originalCandidates, 1)
	originalCandidateID := originalCandidates[0].ID
	require.Equal(t, "connector-a", originalCandidates[0].ConnectorID)
	require.Equal(t, "measurements", originalCandidates[0].TableName)

	currentMappings = []*schema.DatabaseTargetMapping{
		{
			ID:          "db-map-rebind",
			TagID:       overrideTag.ID,
			ConnectorID: "connector-b",
			TableSchema: "analytics",
			TableName:   "measurements_v2",
			ColumnName:  "line_b",
			WriteMode:   schema.DatabaseWriteModeInsert,
		},
	}
	_, err = fixture.ruleSvc.RecomputeCandidateView(context.Background(), rule.ID)
	require.NoError(t, err)

	reboundCandidates := loadDatabaseOutputCandidatesForRule(t, fixture, "rule-db-apply-rebind")
	require.Len(t, reboundCandidates, 1)
	require.NotEqual(t, originalCandidateID, reboundCandidates[0].ID)
	require.Equal(t, "connector-b", reboundCandidates[0].ConnectorID)
	require.Equal(t, "measurements_v2", reboundCandidates[0].TableName)

	applyBody, err := json.Marshal(sourcerule.ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{originalCandidateID},
	})
	require.NoError(t, err)

	applyReq, err := http.NewRequestWithContext(
		context.Background(), http.MethodPost,
		"/datalink/source-rules/rule-db-apply-rebind/database-outputs/apply",
		bytes.NewBuffer(applyBody),
	)
	require.NoError(t, err)
	applyReq.Header.Set("Content-Type", "application/json")
	applyResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(applyResp, applyReq)
	require.Equal(t, http.StatusOK, applyResp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(applyResp.Body.Bytes(), &payload))
	assert.Equal(t, true, payload["success"])

	results := payload["data"].(map[string]any)["results"].([]any)
	require.Len(t, results, 1)
	result := results[0].(map[string]any)
	assert.Equal(t, originalCandidateID, result["candidate_id"])
	assert.Equal(t, "failed", result["status"])
	assert.Equal(t, "validation", result["code"])
	assert.Contains(t, result["reason"], "not found in revision")
}
