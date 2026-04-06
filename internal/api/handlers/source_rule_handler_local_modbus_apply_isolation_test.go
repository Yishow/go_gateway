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

func TestSourceRuleHandler_ApplyLocalModbusOutputs_PreservesDatabaseSnapshot(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	overrideTag, err := fixture.tagSvc.Create(context.Background(), tag.CreateTagRequest{
		Key:         "factory.lm.apply.handler.isolated",
		DisplayName: "LM Apply Handler Isolated",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	fixture.ruleSvc.SetDatabaseTargetMappingReader(sourcerule.DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return []*schema.DatabaseTargetMapping{
			{
				ID:          "db-map-lm-handler-isolated",
				TagID:       overrideTag.ID,
				ConnectorID: "connector-lm-handler-isolated",
				TableSchema: "public",
				TableName:   "measurements",
				ColumnName:  "line_a",
				WriteMode:   schema.DatabaseWriteModeInsert,
			},
		}, nil
	}))

	validation := &sourcerule.DatabaseTargetConnectorValidation{
		Ready: false,
		Issues: []sourcerule.DatabaseTargetValidationIssue{
			{
				Severity:  "error",
				MappingID: "db-map-lm-handler-isolated",
				Code:      "table_missing",
				Message:   "找不到資料表: public.measurements",
			},
		},
	}
	fixture.ruleSvc.SetDatabaseTargetConnectorValidator(sourcerule.DatabaseTargetConnectorValidatorFunc(func(context.Context, string) (*sourcerule.DatabaseTargetConnectorValidation, error) {
		return validation, nil
	}))

	tagCandidates := createRuleCandidatesForDecisionTest(t, fixture, "rule-lm-apply-handler-isolated", "40001", 1)
	decisionBody, err := json.Marshal(sourcerule.UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[0].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	decisionReq, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/datalink/source-rules/rule-lm-apply-handler-isolated/tag-review-decisions",
		bytes.NewBuffer(decisionBody),
	)
	require.NoError(t, err)
	decisionReq.Header.Set("Content-Type", "application/json")
	decisionResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(decisionResp, decisionReq)
	require.Equal(t, http.StatusOK, decisionResp.Code)

	ctx := context.Background()
	rule, err := fixture.repo.GetByID(ctx, "rule-lm-apply-handler-isolated")
	require.NoError(t, err)

	databasePayload := `{"candidates":[{"id":"db-handler-preserved","tag_key":"DB_HANDLER","connector_id":"connector-handler","table_name":"measurements","column_name":"line_a","status":"blocked","blocking_reason":"database context pending"}]}`
	setHandlerDatabaseSnapshot(t, fixture, rule.ID, rule.RevisionID, databasePayload, schema.SourceRuleCandidateStatusBlocked, "preserved database handler state")

	localModbusPayload := `{"candidates":[{"id":"lm-handler-preserved","tag_key":"LM_HANDLER","register":11}]}`
	setHandlerLocalModbusSnapshot(t, fixture, rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusReady, "")

	applyBody, err := json.Marshal(sourcerule.ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{"lm-handler-preserved"},
	})
	require.NoError(t, err)

	applyReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/datalink/source-rules/rule-lm-apply-handler-isolated/local-modbus/apply",
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
	assert.Equal(t, "lm-handler-preserved", results[0].(map[string]any)["candidate_id"])

	assertHandlerDatabaseSnapshot(t, fixture, rule.ID, rule.RevisionID, databasePayload, schema.SourceRuleCandidateStatusBlocked, "preserved database handler state")

	validation = &sourcerule.DatabaseTargetConnectorValidation{
		Ready:  true,
		Issues: []sourcerule.DatabaseTargetValidationIssue{},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/datalink/source-rules/rule-lm-apply-handler-isolated/candidates", http.NoBody)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	payload = map[string]any{}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	data := payload["data"].(map[string]any)
	databaseOutputs := data["database_outputs"].(map[string]any)
	assert.Equal(t, "ready", databaseOutputs["status"])
	assert.Empty(t, databaseOutputs["reason"])
	databaseCandidate := databaseOutputs["candidates"].([]any)[0].(map[string]any)
	assert.Equal(t, "factory.lm.apply.handler.isolated", databaseCandidate["tag_key"])
	assert.Equal(t, "connector-lm-handler-isolated", databaseCandidate["connector_id"])
	assert.Equal(t, "public", databaseCandidate["table_schema"])
	assert.Equal(t, "measurements", databaseCandidate["table_name"])
	assert.Equal(t, "line_a", databaseCandidate["column_name"])

	localModbusOutputs := data["local_modbus_outputs"].(map[string]any)
	assert.Equal(t, "ready", localModbusOutputs["status"])
	localModbusCandidate := localModbusOutputs["candidates"].([]any)[0].(map[string]any)
	assert.Equal(t, "LM_HANDLER", localModbusCandidate["tag_key"])
	assert.Equal(t, float64(11), localModbusCandidate["register"])

	assertHandlerLocalModbusSnapshot(t, fixture, rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusReady, "")
}

func setHandlerDatabaseSnapshot(
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
		if candidateSnapshot.CandidateType == schema.SourceRuleCandidateTypeDatabaseOutputs {
			candidateSnapshot.Payload = payload
			candidateSnapshot.Status = status
			candidateSnapshot.Reason = reason
		}
		nextSnapshots = append(nextSnapshots, &candidateSnapshot)
	}
	require.NoError(t, fixture.repo.ReplaceCandidateSnapshots(context.Background(), nextSnapshots))
}

func assertHandlerDatabaseSnapshot(
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
		if snapshot.CandidateType != schema.SourceRuleCandidateTypeDatabaseOutputs {
			continue
		}
		assert.Equal(t, status, snapshot.Status)
		assert.Equal(t, reason, snapshot.Reason)
		assert.JSONEq(t, payload, snapshot.Payload)
		return
	}

	t.Fatalf("database output snapshot not found for rule %s revision %s", ruleID, revisionID)
}
