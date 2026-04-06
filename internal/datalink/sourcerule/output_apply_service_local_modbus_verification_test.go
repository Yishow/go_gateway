package sourcerule

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_ApplyLocalModbusOutputCandidates_BlockedConflictReturnsFailedResult(t *testing.T) {
	t.Parallel()

	ctx, svc, rule := createLocalModbusApplyValidationRule(t, "device-lm-apply-verify-conflict", "rule-lm-apply-verify-conflict", "40061")

	localModbusPayload := `{"candidates":[{"id":"lm-candidate-conflict","tag_id":"tag-conflict","tag_key":"LM_CONFLICT","register":11,"status":"blocked_conflict","blocking_reason":"local 5020 register conflict: overlaps rule other"}]}`
	setLocalModbusSnapshotForRule(t, svc.repo.(*MemoryRepository), rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusReady, "")

	response, err := svc.ApplyLocalModbusOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{"lm-candidate-conflict"},
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Len(t, response.Results, 1)

	assert.Equal(t, "lm-candidate-conflict", response.Results[0].CandidateID)
	assert.Equal(t, "failed", response.Results[0].Status)
	assert.Equal(t, "conflict", response.Results[0].Code)
	assert.Equal(t, "local 5020 register conflict: overlaps rule other", response.Results[0].Reason)
}

func TestService_ApplyLocalModbusOutputCandidates_MissingRegisterReturnsFailedResult(t *testing.T) {
	t.Parallel()

	ctx, svc, rule := createLocalModbusApplyValidationRule(t, "device-lm-apply-verify-register", "rule-lm-apply-verify-register", "40062")

	localModbusPayload := `{"candidates":[{"id":"lm-candidate-pending-register","tag_id":"tag-register","tag_key":"LM_PENDING_REGISTER","status":"deferred"}]}`
	setLocalModbusSnapshotForRule(t, svc.repo.(*MemoryRepository), rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusReady, "")

	response, err := svc.ApplyLocalModbusOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{"lm-candidate-pending-register"},
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Len(t, response.Results, 1)

	assert.Equal(t, "lm-candidate-pending-register", response.Results[0].CandidateID)
	assert.Equal(t, "failed", response.Results[0].Status)
	assert.Equal(t, "register_missing", response.Results[0].Code)
	assert.Equal(t, "local modbus register is not configured", response.Results[0].Reason)
}

func TestService_ApplyLocalModbusOutputCandidates_MissingTagReturnsFailedResult(t *testing.T) {
	t.Parallel()

	ctx, svc, rule := createLocalModbusApplyValidationRule(t, "device-lm-apply-verify-tag", "rule-lm-apply-verify-tag", "40063")

	localModbusPayload := `{"candidates":[{"id":"lm-candidate-pending-tag","tag_key":"LM_PENDING_TAG","register":11,"status":"deferred"}]}`
	setLocalModbusSnapshotForRule(t, svc.repo.(*MemoryRepository), rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusReady, "")

	response, err := svc.ApplyLocalModbusOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{"lm-candidate-pending-tag"},
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Len(t, response.Results, 1)

	assert.Equal(t, "lm-candidate-pending-tag", response.Results[0].CandidateID)
	assert.Equal(t, "failed", response.Results[0].Status)
	assert.Equal(t, "tag_missing", response.Results[0].Code)
	assert.Equal(t, "local modbus tag is not configured", response.Results[0].Reason)
}
