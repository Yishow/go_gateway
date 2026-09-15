package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_CandidateSnapshots_BlocksOnlyConflictingLocalModbusCandidates(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-lm-conflicts")
	require.NoError(t, err)

	ruleA, linksA := createLocalModbusConflictRule(ctx, t, svc, repo, dev.ID, "rule-lm-conflict-a", "40001", schema.DataTypeInt32)
	ruleB, linksB := createLocalModbusConflictRule(ctx, t, svc, repo, dev.ID, "rule-lm-conflict-b", "40101", schema.DataTypeInt16)
	ruleC, linksC := createLocalModbusConflictRule(ctx, t, svc, repo, dev.ID, "rule-lm-conflict-c", "40201", schema.DataTypeInt16)

	require.NotNil(t, linksA[0].TagID)
	require.NotNil(t, linksB[0].TagID)
	require.NotNil(t, linksC[0].TagID)

	svc.SetLocalModbusMappingReader(LocalModbusMappingListFunc(func(context.Context) ([]LocalModbusMappingRecord, error) {
		return []LocalModbusMappingRecord{
			{TagID: *linksA[0].TagID, Register: 10, DataType: schema.DataTypeInt32},
			{TagID: *linksB[0].TagID, Register: 11, DataType: schema.DataTypeInt16},
			{TagID: *linksC[0].TagID, Register: 20, DataType: schema.DataTypeInt16},
		}, nil
	}))

	require.NoError(t, svc.persistCandidateSnapshots(ctx, ruleA, linksA))
	require.NoError(t, svc.persistCandidateSnapshots(ctx, ruleB, linksB))
	require.NoError(t, svc.persistCandidateSnapshots(ctx, ruleC, linksC))

	candidateA, snapshotA := localModbusSingleCandidateFromRule(ctx, t, svc, ruleA.ID)
	candidateB, snapshotB := localModbusSingleCandidateFromRule(ctx, t, svc, ruleB.ID)
	candidateC, snapshotC := localModbusSingleCandidateFromRule(ctx, t, svc, ruleC.ID)

	assert.Equal(t, schema.SourceRuleCandidateStatusDeferred, snapshotA.Status)
	assert.Equal(t, schema.SourceRuleCandidateStatusDeferred, snapshotB.Status)
	assert.Equal(t, schema.SourceRuleCandidateStatusDeferred, snapshotC.Status)

	assert.Equal(t, schema.SourceRuleLocalModbusOutputStatusBlockedConflict, candidateA.Status)
	assert.NotEmpty(t, candidateA.BlockingReason)
	assert.Equal(t, schema.SourceRuleLocalModbusOutputStatusBlockedConflict, candidateB.Status)
	assert.NotEmpty(t, candidateB.BlockingReason)

	assert.Equal(t, schema.SourceRuleLocalModbusOutputStatusDeferred, candidateC.Status)
	assert.Empty(t, candidateC.BlockingReason)
}

func createLocalModbusConflictRule(
	ctx context.Context,
	t *testing.T,
	svc *Service,
	repo *MemoryRepository,
	deviceID, ruleID, startAddress string,
	dataType schema.DataType,
) (*schema.SourceRule, []*schema.SourceRuleLink) {
	t.Helper()

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           ruleID,
		DeviceID:     deviceID,
		StartAddress: startAddress,
		Count:        1,
		DataType:     dataType,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	links := applyRuleManagedLinks(ctx, t, repo, svc, rule.ID)
	require.Len(t, links, 1)
	return rule, links
}

func localModbusSingleCandidateFromRule(
	ctx context.Context,
	t *testing.T,
	svc *Service,
	ruleID string,
) (schema.SourceRuleLocalModbusOutputCandidate, *schema.SourceRuleCandidateSnapshot) {
	t.Helper()

	snapshots, err := svc.ListCandidateSnapshots(ctx, ruleID)
	require.NoError(t, err)

	localModbusSnapshot := localModbusSnapshotFromSnapshots(t, snapshots)
	candidates, err := decodeCandidatePayload[schema.SourceRuleLocalModbusOutputCandidate](localModbusSnapshot.Payload)
	require.NoError(t, err)
	require.Len(t, candidates, 1)
	return candidates[0], localModbusSnapshot
}
