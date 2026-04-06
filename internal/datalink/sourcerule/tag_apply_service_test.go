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

func TestService_ApplyTagCandidates_UsesApprovedRenameDecision(t *testing.T) {
	t.Parallel()

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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-apply-rename")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-apply-rename",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	candidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID: candidate.ID,
		Action:      schema.SourceRuleTagReviewDecisionActionRename,
		TagKey:      "factory.temperature",
	})
	require.NoError(t, err)

	response, err := svc.ApplyTagCandidates(ctx, rule.ID, ApplyTagCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{candidate.ID},
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Len(t, response.Results, 1)
	assert.Equal(t, "applied", response.Results[0].Status)
	assert.Equal(t, candidate.ID, response.Results[0].CandidateID)

	tagRecord, err := tagSvc.GetByID(ctx, response.Results[0].TagID)
	require.NoError(t, err)
	assert.Equal(t, "factory.temperature", tagRecord.Key)

	mappingRecord, err := mappingSvc.GetByID(ctx, response.Results[0].MappingID)
	require.NoError(t, err)
	assert.Equal(t, tagRecord.ID, mappingRecord.TagID)
	assert.Equal(t, schema.MappingStatusActive, mappingRecord.Status)
	assert.Equal(t, mappingRecord.ProposedSignature, mappingRecord.LastAppliedSignature)
	assert.Empty(t, mappingRecord.BlockingReason)

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)
	require.NotNil(t, links[0].TagID)
	require.NotNil(t, links[0].MappingID)
	assert.Equal(t, tagRecord.ID, *links[0].TagID)
	assert.Equal(t, mappingRecord.ID, *links[0].MappingID)
}

func TestService_ApplyTagCandidates_RejectsRevisionMismatch(t *testing.T) {
	t.Parallel()

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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-apply-conflict")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-apply-conflict",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	candidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)

	_, err = svc.ApplyTagCandidates(ctx, rule.ID, ApplyTagCandidatesRequest{
		RevisionID:   "stale-revision",
		CandidateIDs: []string{candidate.ID},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrTagApplyRevisionConflict)

	tags, err := tagSvc.List(ctx, tag.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, tags)

	mappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, mappings)
}
