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

func TestSourceRuleHandler_Candidates_UsesEffectiveTagReviewStateForDatabaseOutputs(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	tagCandidates := createRuleCandidatesForDecisionTest(t, fixture, "rule-db-effective", "40001", 3)

	overrideTag, err := fixture.tagSvc.Create(context.Background(), tag.CreateTagRequest{
		Key:         "factory.db.override",
		DisplayName: "DB Override",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)

	decisionRequests := []sourcerule.UpsertTagReviewDecisionRequest{
		{
			CandidateID: tagCandidates[0].ID,
			Action:      schema.SourceRuleTagReviewDecisionActionRename,
			TagKey:      "factory.db.renamed",
		},
		{
			CandidateID: tagCandidates[1].ID,
			Action:      schema.SourceRuleTagReviewDecisionActionSkip,
		},
		{
			CandidateID:   tagCandidates[2].ID,
			Action:        schema.SourceRuleTagReviewDecisionActionOverride,
			OverrideTagID: overrideTag.ID,
		},
	}
	for _, request := range decisionRequests {
		body, marshalErr := json.Marshal(request)
		require.NoError(t, marshalErr)

		req, newReqErr := http.NewRequest(
			http.MethodPost,
			"/datalink/source-rules/rule-db-effective/tag-review-decisions",
			bytes.NewBuffer(body),
		)
		require.NoError(t, newReqErr)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		fixture.router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
	}

	req, err := http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-db-effective/candidates", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	data := payload["data"].(map[string]any)
	databaseOutputs := data["database_outputs"].(map[string]any)
	assert.Equal(t, "ready", databaseOutputs["status"])
	assert.Empty(t, databaseOutputs["reason"])

	candidates := databaseOutputs["candidates"].([]any)
	require.Len(t, candidates, 2)
	candidateByAddress := make(map[string]map[string]any, len(candidates))
	for _, item := range candidates {
		record := item.(map[string]any)
		candidateByAddress[record["address"].(string)] = record
	}

	require.Contains(t, candidateByAddress, "40001")
	assert.Equal(t, "factory.db.renamed", candidateByAddress["40001"]["tag_key"])
	assert.Nil(t, candidateByAddress["40001"]["tag_id"])

	require.Contains(t, candidateByAddress, "40003")
	assert.Equal(t, "factory.db.override", candidateByAddress["40003"]["tag_key"])
	assert.Equal(t, overrideTag.ID, candidateByAddress["40003"]["tag_id"])
}

func TestSourceRuleHandler_Candidates_InfersGroupedDatabaseScopeForSlashTagKeys(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	tagCandidates := createRuleCandidatesForDecisionTest(t, fixture, "rule-db-grouped", "40001", 1)
	require.Len(t, tagCandidates, 1)

	body, err := json.Marshal(sourcerule.UpsertTagReviewDecisionRequest{
		CandidateID: tagCandidates[0].ID,
		Action:      schema.SourceRuleTagReviewDecisionActionRename,
		TagKey:      "meter/A1",
	})
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-db-grouped/tag-review-decisions",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	req, err = http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-db-grouped/candidates", nil)
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	data := payload["data"].(map[string]any)
	databaseOutputs := data["database_outputs"].(map[string]any)
	require.Len(t, databaseOutputs["candidates"].([]any), 1)

	candidate := databaseOutputs["candidates"].([]any)[0].(map[string]any)
	assert.Equal(t, "meter", candidate["group_key"])
	assert.Equal(t, "a1", candidate["column_name"])
	assert.Equal(t, float64(15), candidate["write_interval_seconds"])
}

func TestSourceRuleHandler_Candidates_UsesPersistedDatabaseScopeForDatabaseOutputs(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	timestampColumn := "ts"
	writeIntervalSeconds := 15
	overrideTag, err := fixture.tagSvc.Create(context.Background(), tag.CreateTagRequest{
		Key:         "factory.db.bound",
		DisplayName: "DB Bound",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	fixture.ruleSvc.SetDatabaseTargetMappingReader(
		sourcerule.DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
			return []*schema.DatabaseTargetMapping{
				{
					ID:                   "db-map-1",
					TagID:                overrideTag.ID,
					ConnectorID:          "connector-1",
					TableSchema:          "public",
					TableName:            "measurements",
					ColumnName:           "line_a",
					WriteMode:            schema.DatabaseWriteModeUpsert,
					WriteIntervalSeconds: &writeIntervalSeconds,
					TimestampColumn:      &timestampColumn,
				},
			}, nil
		}),
	)

	tagCandidates := createRuleCandidatesForDecisionTest(t, fixture, "rule-db-scope", "40001", 1)
	body, err := json.Marshal(sourcerule.UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[0].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-db-scope/tag-review-decisions",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	req, err = http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-db-scope/candidates", nil)
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	data := payload["data"].(map[string]any)
	databaseOutputs := data["database_outputs"].(map[string]any)
	require.Len(t, databaseOutputs["candidates"].([]any), 1)

	candidate := databaseOutputs["candidates"].([]any)[0].(map[string]any)
	assert.Equal(t, "db-map-1", candidate["mapping_id"])
	assert.Equal(t, "connector-1", candidate["connector_id"])
	assert.Equal(t, "public", candidate["table_schema"])
	assert.Equal(t, "measurements", candidate["table_name"])
	assert.Equal(t, "line_a", candidate["column_name"])
	_, hasGroupKey := candidate["group_key"]
	assert.True(t, hasGroupKey)
	assert.Nil(t, candidate["group_key"])
	assert.Equal(t, float64(15), candidate["write_interval_seconds"])
	assert.Equal(t, "upsert", candidate["write_mode"])
	assert.Equal(t, "ts", candidate["timestamp_column"])
}

func TestSourceRuleHandler_Candidates_MarksInvalidDatabaseScopeOutOfSync(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	overrideTag, err := fixture.tagSvc.Create(context.Background(), tag.CreateTagRequest{
		Key:         "factory.db.broken",
		DisplayName: "DB Broken",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	fixture.ruleSvc.SetDatabaseTargetMappingReader(
		sourcerule.DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
			return []*schema.DatabaseTargetMapping{
				{
					ID:          "db-map-broken",
					TagID:       overrideTag.ID,
					ConnectorID: "connector-broken",
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
			return &sourcerule.DatabaseTargetConnectorValidation{
				Ready: false,
				Issues: []sourcerule.DatabaseTargetValidationIssue{
					{
						Severity:  "error",
						MappingID: "db-map-broken",
						Code:      "table_missing",
						Message:   "找不到資料表: public.measurements",
					},
				},
			}, nil
		}),
	)

	tagCandidates := createRuleCandidatesForDecisionTest(t, fixture, "rule-db-invalid", "40001", 1)
	body, err := json.Marshal(sourcerule.UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[0].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-db-invalid/tag-review-decisions",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	req, err = http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-db-invalid/candidates", nil)
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	data := payload["data"].(map[string]any)
	databaseOutputs := data["database_outputs"].(map[string]any)
	assert.Equal(t, "blocked", databaseOutputs["status"])
	assert.Contains(t, databaseOutputs["reason"], "找不到資料表")

	candidate := databaseOutputs["candidates"].([]any)[0].(map[string]any)
	assert.Equal(t, "out_of_sync", candidate["status"])
	assert.Contains(t, candidate["blocking_reason"], "找不到資料表")
}

func TestSourceRuleHandler_Candidates_RevalidatesDatabaseScopeWhenConnectorBecomesReady(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	overrideTag, err := fixture.tagSvc.Create(context.Background(), tag.CreateTagRequest{
		Key:         "factory.db.recovered",
		DisplayName: "DB Recovered",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	fixture.ruleSvc.SetDatabaseTargetMappingReader(
		sourcerule.DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
			return []*schema.DatabaseTargetMapping{
				{
					ID:          "db-map-recovered",
					TagID:       overrideTag.ID,
					ConnectorID: "connector-recovered",
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
				MappingID: "db-map-recovered",
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

	tagCandidates := createRuleCandidatesForDecisionTest(t, fixture, "rule-db-revalidate", "40001", 1)
	body, err := json.Marshal(sourcerule.UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[0].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	req, err := http.NewRequest(
		http.MethodPost,
		"/datalink/source-rules/rule-db-revalidate/tag-review-decisions",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	req, err = http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-db-revalidate/candidates", nil)
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	data := payload["data"].(map[string]any)
	databaseOutputs := data["database_outputs"].(map[string]any)
	assert.Equal(t, "blocked", databaseOutputs["status"])
	candidate := databaseOutputs["candidates"].([]any)[0].(map[string]any)
	assert.Equal(t, "out_of_sync", candidate["status"])

	validation = &sourcerule.DatabaseTargetConnectorValidation{
		Ready:  true,
		Issues: []sourcerule.DatabaseTargetValidationIssue{},
	}

	req, err = http.NewRequest(http.MethodGet, "/datalink/source-rules/rule-db-revalidate/candidates", nil)
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	payload = map[string]any{}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	data = payload["data"].(map[string]any)
	databaseOutputs = data["database_outputs"].(map[string]any)
	assert.Equal(t, "ready", databaseOutputs["status"])
	assert.Empty(t, databaseOutputs["reason"])
	candidate = databaseOutputs["candidates"].([]any)[0].(map[string]any)
	assert.Equal(t, "ready", candidate["status"])
	assert.Empty(t, candidate["blocking_reason"])
}
