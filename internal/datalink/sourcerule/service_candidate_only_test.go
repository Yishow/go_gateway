package sourcerule

import (
	"context"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Create_LeavesRuleTagsAndMappingsPendingUntilApply(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-candidate-only")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-candidate-only",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "MIXER",
		Enabled:      true,
	})
	require.NoError(t, err)

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 2)
	for _, link := range links {
		assert.Nil(t, link.TagID)
		assert.Nil(t, link.MappingID)
	}

	tags, err := tagSvc.List(ctx, tag.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, tags)

	mappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, mappings)

	snapshots, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, snapshots, 3)

	var payload tagSnapshotPayload
	for _, snapshot := range snapshots {
		if snapshot.CandidateType != schema.SourceRuleCandidateTypeTags {
			continue
		}
		require.NoError(t, json.Unmarshal([]byte(snapshot.Payload), &payload))
	}
	require.Len(t, payload.Candidates, 2)
	assert.Nil(t, payload.Candidates[0].TagID)
	assert.Nil(t, payload.Candidates[0].MappingID)
}

func TestService_Enable_LeavesRuleTagsAndMappingsPendingUntilApply(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-enable-candidate-only")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-enable-candidate-only",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "MIXER",
		Enabled:      false,
	})
	require.NoError(t, err)

	require.NoError(t, svc.Enable(ctx, rule.ID))

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)
	assert.Nil(t, links[0].TagID)
	assert.Nil(t, links[0].MappingID)

	pointRecord, err := pointSvc.GetByID(ctx, links[0].PointID)
	require.NoError(t, err)
	assert.True(t, pointRecord.Enabled)

	tags, err := tagSvc.List(ctx, tag.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, tags)

	mappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, mappings)

	require.NoError(t, svc.Disable(ctx, rule.ID))

	pointRecord, err = pointSvc.GetByID(ctx, links[0].PointID)
	require.NoError(t, err)
	assert.False(t, pointRecord.Enabled)
}
