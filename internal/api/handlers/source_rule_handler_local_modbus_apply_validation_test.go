package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourceRuleHandler_ApplyLocalModbusOutputs_RejectsRevisionMismatch(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	createRuleCandidatesForDecisionTest(t, fixture, "rule-lm-apply-conflict", "40001", 1)

	applyBody, err := json.Marshal(sourcerule.ApplyOutputCandidatesRequest{
		RevisionID:   "stale-revision",
		CandidateIDs: []string{"lm-candidate-1"},
	})
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/datalink/source-rules/rule-lm-apply-conflict/local-modbus/apply",
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

func TestSourceRuleHandler_ApplyLocalModbusOutputs_RejectsMissingRevisionID(t *testing.T) {
	t.Parallel()

	fixture := setupSourceRuleCandidatesFixture(t)
	createRuleCandidatesForDecisionTest(t, fixture, "rule-lm-apply-validation", "40021", 1)

	applyBody, err := json.Marshal(sourcerule.ApplyOutputCandidatesRequest{
		CandidateIDs: []string{"lm-candidate-1"},
	})
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/datalink/source-rules/rule-lm-apply-validation/local-modbus/apply",
		bytes.NewBuffer(applyBody),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusBadRequest, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	errorPayload := payload["error"].(map[string]any)
	assert.Equal(t, "validation", errorPayload["code"])
}
