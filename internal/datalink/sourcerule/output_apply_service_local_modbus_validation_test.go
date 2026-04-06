package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_ApplyLocalModbusOutputCandidates_RejectsRevisionMismatch(t *testing.T) {
	t.Parallel()

	ctx, svc, rule := createLocalModbusApplyValidationRule(t, "device-lm-apply-conflict", "rule-lm-apply-conflict", "40001")

	_, err := svc.ApplyLocalModbusOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   "stale-revision",
		CandidateIDs: []string{"lm-candidate-1"},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrOutputApplyRevisionConflict)
}

func TestService_ApplyLocalModbusOutputCandidates_RejectsMissingRevisionID(t *testing.T) {
	t.Parallel()

	ctx, svc, rule := createLocalModbusApplyValidationRule(t, "device-lm-apply-validation", "rule-lm-apply-validation", "40021")

	_, err := svc.ApplyLocalModbusOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		CandidateIDs: []string{"lm-candidate-1"},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidOutputApplyRequest)
}

func createLocalModbusApplyValidationRule(
	t *testing.T,
	deviceID string,
	ruleID string,
	startAddress string,
) (context.Context, *Service, *schema.SourceRule) {
	t.Helper()

	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, deviceID)
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           ruleID,
		DeviceID:     dev.ID,
		StartAddress: startAddress,
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	return ctx, svc, rule
}
