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

func TestService_ApplyTagCandidates_ReturnsPerItemResultsForPartialSuccess(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-apply-partial")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-apply-partial",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	candidates := tagCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, candidates, 2)

	otherRule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-apply-partial-other",
		DeviceID:     dev.ID,
		StartAddress: "40101",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	otherCandidate := firstTagCandidateFromMemoryRepo(t, repo, otherRule.ID, otherRule.RevisionID)
	otherLinks, err := repo.ListLinks(ctx, otherRule.ID)
	require.NoError(t, err)
	require.Len(t, otherLinks, 1)

	existingTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "factory.shared.partial",
		DisplayName: "Shared Partial",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)

	enabled := true
	status := schema.MappingStatusActive
	_, err = mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID:           otherLinks[0].PointID,
		TagID:             existingTag.ID,
		Enabled:           &enabled,
		TransformPipeline: otherCandidate.TransformPipeline,
		Status:            &status,
	})
	require.NoError(t, err)

	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID:   candidates[1].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: existingTag.ID,
	})
	require.NoError(t, err)

	response, err := svc.ApplyTagCandidates(ctx, rule.ID, ApplyTagCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{candidates[0].ID, candidates[1].ID},
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Len(t, response.Results, 2)

	assert.Equal(t, candidates[0].ID, response.Results[0].CandidateID)
	assert.Equal(t, "applied", response.Results[0].Status)
	assert.NotEmpty(t, response.Results[0].TagID)
	assert.NotEmpty(t, response.Results[0].MappingID)
	assert.Empty(t, response.Results[0].Error)

	assert.Equal(t, candidates[1].ID, response.Results[1].CandidateID)
	assert.Equal(t, "failed", response.Results[1].Status)
	assert.Contains(t, response.Results[1].Error, "已綁定其他 point")
	assert.Empty(t, response.Results[1].TagID)
	assert.Empty(t, response.Results[1].MappingID)

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 2)
	linkByAddress := make(map[string]*schema.SourceRuleLink, len(links))
	for _, link := range links {
		linkByAddress[normalizeAddressKey(link.Address)] = link
	}

	firstLink := linkByAddress[normalizeAddressKey(candidates[0].Address)]
	require.NotNil(t, firstLink)
	require.NotNil(t, firstLink.TagID)
	require.NotNil(t, firstLink.MappingID)
	assert.Equal(t, response.Results[0].TagID, *firstLink.TagID)
	assert.Equal(t, response.Results[0].MappingID, *firstLink.MappingID)

	secondLink := linkByAddress[normalizeAddressKey(candidates[1].Address)]
	require.NotNil(t, secondLink)
	assert.Nil(t, secondLink.TagID)
	assert.Nil(t, secondLink.MappingID)

	refreshedCandidates := tagCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, refreshedCandidates, 2)
	refreshedByID := make(map[string]schema.SourceRuleTagCandidate, len(refreshedCandidates))
	for _, refreshedCandidate := range refreshedCandidates {
		refreshedByID[refreshedCandidate.ID] = refreshedCandidate
	}

	require.Contains(t, refreshedByID, candidates[0].ID)
	require.NotNil(t, refreshedByID[candidates[0].ID].TagID)
	require.NotNil(t, refreshedByID[candidates[0].ID].MappingID)
	assert.Equal(t, response.Results[0].TagID, *refreshedByID[candidates[0].ID].TagID)
	assert.Equal(t, response.Results[0].MappingID, *refreshedByID[candidates[0].ID].MappingID)

	require.Contains(t, refreshedByID, candidates[1].ID)
	assert.Nil(t, refreshedByID[candidates[1].ID].TagID)
	assert.Nil(t, refreshedByID[candidates[1].ID].MappingID)
}

func TestService_ApplyTagCandidates_RecomputesCandidateSnapshotsAfterApply(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-apply-refresh")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-apply-refresh",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	candidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	response, err := svc.ApplyTagCandidates(ctx, rule.ID, ApplyTagCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{candidate.ID},
	})
	require.NoError(t, err)
	require.Len(t, response.Results, 1)

	snapshots, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, snapshots, 3)

	refreshedCandidates := tagCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, refreshedCandidates, 1)
	require.NotNil(t, refreshedCandidates[0].TagID)
	require.NotNil(t, refreshedCandidates[0].MappingID)
	assert.Equal(t, response.Results[0].TagID, *refreshedCandidates[0].TagID)
	assert.Equal(t, response.Results[0].MappingID, *refreshedCandidates[0].MappingID)

	snapshotByType := make(map[schema.SourceRuleCandidateType]*schema.SourceRuleCandidateSnapshot, len(snapshots))
	for _, snapshot := range snapshots {
		snapshotByType[snapshot.CandidateType] = snapshot
	}

	require.Contains(t, snapshotByType, schema.SourceRuleCandidateTypeDatabaseOutputs)
	assert.Equal(t, schema.SourceRuleCandidateStatusDeferred, snapshotByType[schema.SourceRuleCandidateTypeDatabaseOutputs].Status)
	assert.Equal(t, databaseOutputsDeferredReason, snapshotByType[schema.SourceRuleCandidateTypeDatabaseOutputs].Reason)

	require.Contains(t, snapshotByType, schema.SourceRuleCandidateTypeLocalModbusOutputs)
	assert.Equal(t, schema.SourceRuleCandidateStatusDeferred, snapshotByType[schema.SourceRuleCandidateTypeLocalModbusOutputs].Status)
	assert.Equal(t, localModbusOutputsDeferredReason, snapshotByType[schema.SourceRuleCandidateTypeLocalModbusOutputs].Reason)
}
