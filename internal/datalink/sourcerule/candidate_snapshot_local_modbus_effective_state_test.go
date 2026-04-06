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

func TestService_CandidateSnapshots_BuildLocalModbusOutputsFromEffectiveTagReviewState(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-lm-effective")
	require.NoError(t, err)

	overrideTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "factory.lm.override",
		DisplayName: "LM Override",
		DataType:    schema.DataTypeInt32,
	})
	require.NoError(t, err)
	svc.SetLocalModbusMappingReader(LocalModbusMappingListFunc(func(context.Context) ([]LocalModbusMappingRecord, error) {
		return []LocalModbusMappingRecord{
			{
				TagID:    overrideTag.ID,
				Register: 11,
				DataType: schema.DataTypeInt32,
			},
		}, nil
	}))

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-lm-effective",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        3,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	_ = applyRuleManagedLinks(t, ctx, repo, svc, rule.ID)
	tagCandidates := tagCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, tagCandidates, 3)

	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID: tagCandidates[0].ID,
		Action:      schema.SourceRuleTagReviewDecisionActionRename,
		TagKey:      "factory.lm.renamed",
	})
	require.NoError(t, err)
	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID: tagCandidates[1].ID,
		Action:      schema.SourceRuleTagReviewDecisionActionSkip,
	})
	require.NoError(t, err)
	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[2].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	snapshots, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)

	localModbusSnapshot := localModbusSnapshotFromSnapshots(t, snapshots)
	assert.Equal(t, schema.SourceRuleCandidateStatusDeferred, localModbusSnapshot.Status)
	assert.Equal(t, rule.RevisionID, localModbusSnapshot.RevisionID)

	candidates, err := decodeCandidatePayload[schema.SourceRuleLocalModbusOutputCandidate](localModbusSnapshot.Payload)
	require.NoError(t, err)
	require.Len(t, candidates, 2)

	candidateByAddress := make(map[string]schema.SourceRuleLocalModbusOutputCandidate, len(candidates))
	for _, candidate := range candidates {
		candidateByAddress[candidate.Address] = candidate
	}

	require.Contains(t, candidateByAddress, "40001")
	assert.Nil(t, candidateByAddress["40001"].TagID)
	assert.Equal(t, "factory.lm.renamed", candidateByAddress["40001"].TagKey)

	require.Contains(t, candidateByAddress, "40003")
	require.NotNil(t, candidateByAddress["40003"].TagID)
	require.NotNil(t, candidateByAddress["40003"].Register)
	assert.Equal(t, overrideTag.ID, *candidateByAddress["40003"].TagID)
	assert.Equal(t, "factory.lm.override", candidateByAddress["40003"].TagKey)
	assert.Equal(t, "LM Override", candidateByAddress["40003"].DisplayName)
	assert.Equal(t, schema.DataTypeInt32, candidateByAddress["40003"].DataType)
	assert.Equal(t, 2, candidateByAddress["40003"].RegisterCount)
	assert.Equal(t, uint16(11), *candidateByAddress["40003"].Register)
}

func TestService_ListCandidateSnapshots_UsesActiveRevisionLocalModbusEffectiveState(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-lm-active-revision")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-lm-active-revision",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	_ = applyRuleManagedLinks(t, ctx, repo, svc, rule.ID)
	tagCandidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID: tagCandidate.ID,
		Action:      schema.SourceRuleTagReviewDecisionActionRename,
		TagKey:      "factory.lm.active-revision",
	})
	require.NoError(t, err)

	oldRevisionID := rule.RevisionID
	updatedRule, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{
		ScaleMultiplier:    float64Ptr(2),
		ScaleMultiplierSet: true,
	})
	require.NoError(t, err)
	require.NotEqual(t, oldRevisionID, updatedRule.RevisionID)

	snapshots, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)

	localModbusSnapshot := localModbusSnapshotFromSnapshots(t, snapshots)
	assert.Equal(t, updatedRule.RevisionID, localModbusSnapshot.RevisionID)

	candidates, err := decodeCandidatePayload[schema.SourceRuleLocalModbusOutputCandidate](localModbusSnapshot.Payload)
	require.NoError(t, err)
	require.Len(t, candidates, 1)
	assert.Nil(t, candidates[0].TagID)
	assert.Equal(t, "factory.lm.active-revision", candidates[0].TagKey)
}
