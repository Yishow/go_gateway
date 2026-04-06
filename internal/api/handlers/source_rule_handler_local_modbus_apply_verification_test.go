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

func TestSourceRuleHandler_ApplyLocalModbusOutputs_MissingRegisterReturnsFailedResult(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	createRuleCandidatesForDecisionTest(t, fixture, "rule-lm-apply-verify-register", "40062", 1)

	ctx := context.Background()
	rule, err := fixture.repo.GetByID(ctx, "rule-lm-apply-verify-register")
	require.NoError(t, err)

	localModbusPayload := `{"candidates":[{"id":"lm-candidate-pending-register","tag_id":"tag-register","tag_key":"LM_PENDING_REGISTER","status":"deferred"}]}`
	setHandlerLocalModbusSnapshot(t, fixture, rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusReady, "")

	applyBody, err := json.Marshal(sourcerule.ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{"lm-candidate-pending-register"},
	})
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"/datalink/source-rules/rule-lm-apply-verify-register/local-modbus/apply",
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
	assert.Equal(t, "lm-candidate-pending-register", results[0].(map[string]any)["candidate_id"])
	assert.Equal(t, "failed", results[0].(map[string]any)["status"])
	assert.Equal(t, "register_missing", results[0].(map[string]any)["code"])
	assert.Equal(t, "local modbus register is not configured", results[0].(map[string]any)["reason"])
}
