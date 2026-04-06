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

func TestSourceRuleHandler_Candidates_RevalidatesDatabaseScopeWithoutResettingLocalModbusSnapshot(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	overrideTag, err := fixture.tagSvc.Create(context.Background(), tag.CreateTagRequest{
		Key:         "factory.db.recovered.isolated",
		DisplayName: "DB Recovered Isolated",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	fixture.ruleSvc.SetDatabaseTargetMappingReader(
		sourcerule.DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
			return []*schema.DatabaseTargetMapping{
				{
					ID:          "db-map-recovered-isolated",
					TagID:       overrideTag.ID,
					ConnectorID: "connector-recovered-isolated",
					TableSchema: "public",
					TableName:   "measurements",
					ColumnName:  "line_a",
					WriteMode:   schema.DatabaseWriteModeInsert,
				},
			}, nil
		}),
	)

	validation := &sourcerule.DatabaseTargetConnectorValidation{
		Ready: false,
		Issues: []sourcerule.DatabaseTargetValidationIssue{
			{
				Severity:  "error",
				MappingID: "db-map-recovered-isolated",
				Code:      "table_missing",
				Message:   "找不到資料表: public.measurements",
			},
		},
	}
	fixture.ruleSvc.SetDatabaseTargetConnectorValidator(
		sourcerule.DatabaseTargetConnectorValidatorFunc(func(context.Context, string) (*sourcerule.DatabaseTargetConnectorValidation, error) {
			return validation, nil
		}),
	)

	tagCandidates := createRuleCandidatesForDecisionTest(t, fixture, "rule-db-revalidate-isolated", "40001", 1)
	body, err := json.Marshal(sourcerule.UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[0].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-db-revalidate-isolated/tag-review-decisions",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	ctx := context.Background()
	rule, err := fixture.repo.GetByID(ctx, "rule-db-revalidate-isolated")
	require.NoError(t, err)
	snapshots, err := fixture.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
	require.NoError(t, err)

	localModbusPayload := `{"candidates":[{"tag_id":"tag-local","tag_key":"LM_LOCAL","register":7}]}`
	nextSnapshots := make([]*schema.SourceRuleCandidateSnapshot, 0, len(snapshots))
	for _, snapshot := range snapshots {
		candidateSnapshot := *snapshot
		if candidateSnapshot.CandidateType == schema.SourceRuleCandidateTypeLocalModbusOutputs {
			candidateSnapshot.Payload = localModbusPayload
			candidateSnapshot.Status = schema.SourceRuleCandidateStatusBlocked
			candidateSnapshot.Reason = "preserved local modbus state"
		}
		nextSnapshots = append(nextSnapshots, &candidateSnapshot)
	}
	require.NoError(t, fixture.repo.ReplaceCandidateSnapshots(ctx, nextSnapshots))

	validation = &sourcerule.DatabaseTargetConnectorValidation{
		Ready:  true,
		Issues: []sourcerule.DatabaseTargetValidationIssue{},
	}

	req, err = http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-db-revalidate-isolated/candidates", nil)
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	data := payload["data"].(map[string]any)
	databaseOutputs := data["database_outputs"].(map[string]any)
	assert.Equal(t, "ready", databaseOutputs["status"])
	candidate := databaseOutputs["candidates"].([]any)[0].(map[string]any)
	assert.Equal(t, "ready", candidate["status"])

	localModbusOutputs := data["local_modbus_outputs"].(map[string]any)
	assert.Equal(t, "blocked", localModbusOutputs["status"])
	assert.Equal(t, "preserved local modbus state", localModbusOutputs["reason"])
	localModbusCandidate := localModbusOutputs["candidates"].([]any)[0].(map[string]any)
	assert.Equal(t, "LM_LOCAL", localModbusCandidate["tag_key"])
	assert.Equal(t, float64(7), localModbusCandidate["register"])

	persistedSnapshots, err := fixture.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
	require.NoError(t, err)
	for _, snapshot := range persistedSnapshots {
		if snapshot.CandidateType != schema.SourceRuleCandidateTypeLocalModbusOutputs {
			continue
		}
		assert.Equal(t, schema.SourceRuleCandidateStatusBlocked, snapshot.Status)
		assert.Equal(t, "preserved local modbus state", snapshot.Reason)
		assert.JSONEq(t, localModbusPayload, snapshot.Payload)
		return
	}

	t.Fatal("expected persisted local modbus snapshot")
}
