package sourcerule

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_ApplyLocalModbusOutputCandidates_ReturnsPerItemResultsForPartialSuccess(t *testing.T) {
	t.Parallel()

	ctx, svc, rule := createLocalModbusApplyValidationRule(t, "device-lm-apply-partial", "rule-lm-apply-partial", "40041")

	localModbusPayload := `{"candidates":[{"id":"lm-candidate-ready","tag_key":"LM_READY","register":11}]}`
	setLocalModbusSnapshotForRule(t, svc.repo.(*MemoryRepository), rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusReady, "")

	response, err := svc.ApplyLocalModbusOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{"lm-candidate-ready", "lm-candidate-missing"},
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Len(t, response.Results, 2)

	assert.Equal(t, "lm-candidate-ready", response.Results[0].CandidateID)
	assert.Equal(t, "success", response.Results[0].Status)
	assert.Empty(t, response.Results[0].Code)
	assert.Empty(t, response.Results[0].Reason)

	assert.Equal(t, "lm-candidate-missing", response.Results[1].CandidateID)
	assert.Equal(t, "failed", response.Results[1].Status)
	assert.Equal(t, "validation", response.Results[1].Code)
	assert.Equal(t, "candidate lm-candidate-missing not found in revision "+rule.RevisionID, response.Results[1].Reason)
}
