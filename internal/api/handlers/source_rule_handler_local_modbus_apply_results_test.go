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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourceRuleHandler_ApplyLocalModbusOutputs_ReturnsPerCandidateResultsForPartialSuccess(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	createRuleCandidatesForDecisionTest(t, fixture, "rule-lm-apply-partial", "40041", 1)

	ctx := context.Background()
	rule, err := fixture.repo.GetByID(ctx, "rule-lm-apply-partial")
	require.NoError(t, err)

	localModbusPayload := `{"candidates":[{"id":"lm-candidate-ready","tag_key":"LM_READY","register":11}]}`
	setHandlerLocalModbusSnapshot(t, fixture, rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusReady, "")

	applyBody, err := json.Marshal(sourcerule.ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{"lm-candidate-ready", "lm-candidate-missing"},
	})
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/datalink/source-rules/rule-lm-apply-partial/local-modbus/apply",
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
	require.Len(t, results, 2)

	assert.Equal(t, "lm-candidate-ready", results[0].(map[string]any)["candidate_id"])
	assert.Equal(t, "success", results[0].(map[string]any)["status"])

	assert.Equal(t, "lm-candidate-missing", results[1].(map[string]any)["candidate_id"])
	assert.Equal(t, "failed", results[1].(map[string]any)["status"])
	assert.Equal(t, "validation", results[1].(map[string]any)["code"])
	assert.Equal(t, "candidate lm-candidate-missing not found in revision "+rule.RevisionID, results[1].(map[string]any)["reason"])
}
