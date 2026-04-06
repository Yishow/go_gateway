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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourceRuleHandler_ApplyDatabaseOutputs_RejectsRevisionMismatch(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	createRuleCandidatesForDecisionTest(t, fixture, "rule-db-apply-conflict", "40001", 1)
	databaseCandidates := loadDatabaseOutputCandidatesForRule(t, fixture, "rule-db-apply-conflict")
	require.Len(t, databaseCandidates, 1)

	applyBody, err := json.Marshal(sourcerule.ApplyOutputCandidatesRequest{
		RevisionID:   "stale-revision",
		CandidateIDs: []string{databaseCandidates[0].ID},
	})
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-db-apply-conflict/database-outputs/apply",
		bytes.NewBuffer(applyBody),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusConflict, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	errorPayload := payload["error"].(map[string]any)
	assert.Equal(t, "revision_mismatch", errorPayload["code"])
}

func TestSourceRuleHandler_ApplyDatabaseOutputs_ReturnsPerCandidateResultsForPartialSuccess(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	tagCandidates := createRuleCandidatesForDecisionTest(t, fixture, "rule-db-apply-partial", "40001", 2)
	require.Len(t, tagCandidates, 2)

	overrideTag, err := fixture.tagSvc.Create(context.Background(), tag.CreateTagRequest{
		Key:         "factory.db.partial",
		DisplayName: "DB Partial",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)

	fixture.ruleSvc.SetDatabaseTargetMappingReader(
		sourcerule.DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
			return []*schema.DatabaseTargetMapping{
				{
					ID:          "db-map-partial",
					TagID:       overrideTag.ID,
					ConnectorID: "connector-ready",
					TableSchema: "public",
					TableName:   "measurements",
					ColumnName:  "line_a",
					WriteMode:   schema.DatabaseWriteModeInsert,
				},
			}, nil
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

	decisionReq, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-db-apply-partial/tag-review-decisions",
		bytes.NewBuffer(decisionBody),
	)
	require.NoError(t, err)
	decisionReq.Header.Set("Content-Type", "application/json")
	decisionResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(decisionResp, decisionReq)
	require.Equal(t, http.StatusOK, decisionResp.Code)

	databaseCandidates := loadDatabaseOutputCandidatesForRule(t, fixture, "rule-db-apply-partial")
	require.Len(t, databaseCandidates, 2)

	var mappedCandidateID string
	var unmappedCandidateID string
	for _, candidate := range databaseCandidates {
		if candidate.MappingID == nil {
			unmappedCandidateID = candidate.ID
			continue
		}
		mappedCandidateID = candidate.ID
	}
	require.NotEmpty(t, mappedCandidateID)
	require.NotEmpty(t, unmappedCandidateID)

	rule, err := fixture.repo.GetByID(context.Background(), "rule-db-apply-partial")
	require.NoError(t, err)

	applyBody, err := json.Marshal(sourcerule.ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{mappedCandidateID, unmappedCandidateID},
	})
	require.NoError(t, err)

	applyReq, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-db-apply-partial/database-outputs/apply",
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

	data := payload["data"].(map[string]any)
	results := data["results"].([]any)
	require.Len(t, results, 2)

	assert.Equal(t, mappedCandidateID, results[0].(map[string]any)["candidate_id"])
	assert.Equal(t, "success", results[0].(map[string]any)["status"])
	assert.Equal(t, "db-map-partial", results[0].(map[string]any)["mapping_id"])
	assert.Equal(t, "connector-ready", results[0].(map[string]any)["connector_id"])

	assert.Equal(t, unmappedCandidateID, results[1].(map[string]any)["candidate_id"])
	assert.Equal(t, "failed", results[1].(map[string]any)["status"])
	assert.Equal(t, "schema_missing", results[1].(map[string]any)["code"])
}

func TestSourceRuleHandler_ApplyLocalModbusOutputs_DeferredCandidatesReturnSkipped(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	createRuleCandidatesForDecisionTest(t, fixture, "rule-lm-apply-deferred", "40001", 1)

	rule, err := fixture.repo.GetByID(context.Background(), "rule-lm-apply-deferred")
	require.NoError(t, err)

	applyBody, err := json.Marshal(sourcerule.ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{"lm-candidate-1"},
	})
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-lm-apply-deferred/local-modbus/apply",
		bytes.NewBuffer(applyBody),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	assert.Equal(t, true, payload["success"])
	results := payload["data"].(map[string]any)["results"].([]any)
	require.Len(t, results, 1)
	assert.Equal(t, "lm-candidate-1", results[0].(map[string]any)["candidate_id"])
	assert.Equal(t, "skipped", results[0].(map[string]any)["status"])
	assert.Equal(t, "deferred", results[0].(map[string]any)["code"])
}

func loadDatabaseOutputCandidatesForRule(
	t *testing.T,
	fixture *sourceRuleCandidatesFixture,
	ruleID string,
) []schema.SourceRuleDatabaseOutputCandidate {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/datalink/source-rules/%s/candidates", ruleID), nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	data := payload["data"].(map[string]any)
	databaseOutputs := data["database_outputs"].(map[string]any)
	candidatesRaw := databaseOutputs["candidates"]
	encoded, err := json.Marshal(candidatesRaw)
	require.NoError(t, err)

	var candidates []schema.SourceRuleDatabaseOutputCandidate
	require.NoError(t, json.Unmarshal(encoded, &candidates))
	return candidates
}
