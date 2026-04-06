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

func TestSourceRuleHandler_ApplyDatabaseOutputs_PreservesLocalModbusSnapshot(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	overrideTag, err := fixture.tagSvc.Create(context.Background(), tag.CreateTagRequest{
		Key:         "factory.db.apply.handler.isolated",
		DisplayName: "DB Apply Handler Isolated",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	fixture.ruleSvc.SetDatabaseTargetMappingReader(sourcerule.DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return []*schema.DatabaseTargetMapping{
			{
				ID:          "db-map-handler-isolated",
				TagID:       overrideTag.ID,
				ConnectorID: "connector-handler-isolated",
				TableSchema: "public",
				TableName:   "measurements",
				ColumnName:  "line_a",
				WriteMode:   schema.DatabaseWriteModeInsert,
			},
		}, nil
	}))
	fixture.ruleSvc.SetDatabaseTargetConnectorValidator(sourcerule.DatabaseTargetConnectorValidatorFunc(func(context.Context, string) (*sourcerule.DatabaseTargetConnectorValidation, error) {
		return &sourcerule.DatabaseTargetConnectorValidation{Ready: true}, nil
	}))

	tagCandidates := createRuleCandidatesForDecisionTest(t, fixture, "rule-db-apply-handler-isolated", "40001", 1)
	decisionBody, err := json.Marshal(sourcerule.UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[0].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	decisionReq, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-db-apply-handler-isolated/tag-review-decisions",
		bytes.NewBuffer(decisionBody),
	)
	require.NoError(t, err)
	decisionReq.Header.Set("Content-Type", "application/json")
	decisionResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(decisionResp, decisionReq)
	require.Equal(t, http.StatusOK, decisionResp.Code)

	databaseCandidates := loadDatabaseOutputCandidatesForRule(t, fixture, "rule-db-apply-handler-isolated")
	require.Len(t, databaseCandidates, 1)
	require.NotNil(t, databaseCandidates[0].MappingID)

	ctx := context.Background()
	rule, err := fixture.repo.GetByID(ctx, "rule-db-apply-handler-isolated")
	require.NoError(t, err)

	localModbusPayload := `{"candidates":[{"tag_id":"tag-local","tag_key":"LM_HANDLER","register":11}]}`
	setHandlerLocalModbusSnapshot(t, fixture, rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusBlocked, "preserved local modbus handler state")

	applyBody, err := json.Marshal(sourcerule.ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{databaseCandidates[0].ID},
	})
	require.NoError(t, err)

	applyReq, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-db-apply-handler-isolated/database-outputs/apply",
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
	assert.Equal(t, "success", results[0].(map[string]any)["status"])
	assert.Equal(t, "db-map-handler-isolated", results[0].(map[string]any)["mapping_id"])
	assert.Equal(t, "connector-handler-isolated", results[0].(map[string]any)["connector_id"])

	req, err := http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-db-apply-handler-isolated/candidates", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	payload = map[string]any{}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	data := payload["data"].(map[string]any)
	localModbusOutputs := data["local_modbus_outputs"].(map[string]any)
	assert.Equal(t, "blocked", localModbusOutputs["status"])
	assert.Equal(t, "preserved local modbus handler state", localModbusOutputs["reason"])
	localModbusCandidate := localModbusOutputs["candidates"].([]any)[0].(map[string]any)
	assert.Equal(t, "LM_HANDLER", localModbusCandidate["tag_key"])
	assert.Equal(t, float64(11), localModbusCandidate["register"])

	assertHandlerLocalModbusSnapshot(t, fixture, rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusBlocked, "preserved local modbus handler state")
}

func setHandlerLocalModbusSnapshot(
	t *testing.T,
	fixture *sourceRuleCandidatesFixture,
	ruleID string,
	revisionID string,
	payload string,
	status schema.SourceRuleCandidateStatus,
	reason string,
) {
	t.Helper()

	snapshots, err := fixture.repo.ListCandidateSnapshots(context.Background(), ruleID, revisionID)
	require.NoError(t, err)

	nextSnapshots := make([]*schema.SourceRuleCandidateSnapshot, 0, len(snapshots))
	for _, snapshot := range snapshots {
		candidateSnapshot := *snapshot
		if candidateSnapshot.CandidateType == schema.SourceRuleCandidateTypeLocalModbusOutputs {
			candidateSnapshot.Payload = payload
			candidateSnapshot.Status = status
			candidateSnapshot.Reason = reason
		}
		nextSnapshots = append(nextSnapshots, &candidateSnapshot)
	}
	require.NoError(t, fixture.repo.ReplaceCandidateSnapshots(context.Background(), nextSnapshots))
}

func assertHandlerLocalModbusSnapshot(
	t *testing.T,
	fixture *sourceRuleCandidatesFixture,
	ruleID string,
	revisionID string,
	payload string,
	status schema.SourceRuleCandidateStatus,
	reason string,
) {
	t.Helper()

	snapshots, err := fixture.repo.ListCandidateSnapshots(context.Background(), ruleID, revisionID)
	require.NoError(t, err)

	for _, snapshot := range snapshots {
		if snapshot.CandidateType != schema.SourceRuleCandidateTypeLocalModbusOutputs {
			continue
		}
		assert.Equal(t, status, snapshot.Status)
		assert.Equal(t, reason, snapshot.Reason)
		assert.JSONEq(t, payload, snapshot.Payload)
		return
	}

	t.Fatalf("local modbus snapshot not found for rule %s revision %s", ruleID, revisionID)
}
