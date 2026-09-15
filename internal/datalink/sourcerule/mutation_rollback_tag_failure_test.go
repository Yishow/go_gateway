package sourcerule

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_ApplyTagCandidates_TagRollbackFailureReturnsDirtyUnknown(t *testing.T) {
	ctx := context.Background()
	pointRepo := &rollbackPointRepository{
		Repository:     point.NewMemoryRepository(),
		getFailures:    map[int]error{},
		updateFailures: map[int]error{},
	}
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	baseRepo := NewMemoryRepository()
	svc := NewService(baseRepo, deviceSvc, point.NewService(pointRepo, nil), nil)
	tagRepo := &rollbackTagRepository{Repository: tag.NewMemoryRepository(), deleteFailures: 1}
	tagSvc := tag.NewService(tagRepo)
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	share := &rollbackFailCloser{}

	dev, err := seedActiveDevice(ctx, deviceRepo, "tag-rollback-device")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "tag-rollback-rule",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)
	svc.SetTagMappingServices(tagSvc, mappingSvc)
	svc.SetShareRuntimeReconciler(share)
	pointRepo.getFailures[pointRepo.getCalls+2] = errors.New("injected point read failure")
	candidates := tagCandidatesFromMemoryRepo(t, baseRepo, rule.ID, rule.RevisionID)
	require.Len(t, candidates, 2)

	_, err = svc.ApplyTagCandidates(ctx, rule.ID, ApplyTagCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{candidates[0].ID, candidates[1].ID},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dirty_unknown")
	assert.Equal(t, 1, share.failClosedCalls)
}

func TestService_ApplyTagCandidates_CandidateRollbackCleansPriorResults(t *testing.T) {
	ctx := context.Background()
	pointRepo := &rollbackPointRepository{
		Repository:     point.NewMemoryRepository(),
		getFailures:    map[int]error{},
		updateFailures: map[int]error{},
	}
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	baseRepo := NewMemoryRepository()
	svc := NewService(baseRepo, deviceSvc, point.NewService(pointRepo, nil), nil)
	tagRepo := &rollbackTagRepository{Repository: tag.NewMemoryRepository(), deleteFailures: 1}
	tagSvc := tag.NewService(tagRepo)
	mappingRepo := &rollbackMappingRepository{Repository: mapping.NewMemoryRepository(), failCreateOn: 2}
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
	share := &rollbackFailCloser{}

	dev, err := seedActiveDevice(ctx, deviceRepo, "tag-rollback-prior-device")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "tag-rollback-prior-rule",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)
	svc.SetTagMappingServices(tagSvc, mappingSvc)
	svc.SetShareRuntimeReconciler(share)
	candidates := tagCandidatesFromMemoryRepo(t, baseRepo, rule.ID, rule.RevisionID)
	require.Len(t, candidates, 2)

	_, err = svc.ApplyTagCandidates(ctx, rule.ID, ApplyTagCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{candidates[0].ID, candidates[1].ID},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dirty_unknown")
	assert.Equal(t, 1, share.failClosedCalls)
	tags, tagErr := tagSvc.List(ctx, tag.ListFilter{})
	require.NoError(t, tagErr)
	require.Len(t, tags, 1)
	assert.Equal(t, buildPointName(rule.NamingPrefix, candidates[1].Address), tags[0].Key)
	mappings, mappingErr := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, mappingErr)
	assert.Empty(t, mappings)
}

func TestService_ApplyTagCandidates_PreservesPriorResultsWhenCandidateRollbackSucceeds(t *testing.T) {
	ctx := context.Background()
	pointRepo := &rollbackPointRepository{
		Repository:     point.NewMemoryRepository(),
		getFailures:    map[int]error{},
		updateFailures: map[int]error{},
	}
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	baseRepo := NewMemoryRepository()
	svc := NewService(baseRepo, deviceSvc, point.NewService(pointRepo, nil), nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingRepo := &rollbackMappingRepository{Repository: mapping.NewMemoryRepository(), failCreateOn: 2}
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)

	dev, err := seedActiveDevice(ctx, deviceRepo, "tag-rollback-preserve-device")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "tag-rollback-preserve-rule",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)
	svc.SetTagMappingServices(tagSvc, mappingSvc)
	candidates := tagCandidatesFromMemoryRepo(t, baseRepo, rule.ID, rule.RevisionID)
	require.Len(t, candidates, 2)

	response, err := svc.ApplyTagCandidates(ctx, rule.ID, ApplyTagCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{candidates[0].ID, candidates[1].ID},
	})

	require.NoError(t, err)
	require.Len(t, response.Results, 2)
	assert.Equal(t, "applied", response.Results[0].Status)
	assert.Equal(t, "failed", response.Results[1].Status)

	tags, tagErr := tagSvc.List(ctx, tag.ListFilter{})
	require.NoError(t, tagErr)
	require.Len(t, tags, 1)
	assert.Equal(t, response.Results[0].TagID, tags[0].ID)
	mappings, mappingErr := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, mappingErr)
	require.Len(t, mappings, 1)
	assert.Equal(t, response.Results[0].MappingID, mappings[0].ID)
}
