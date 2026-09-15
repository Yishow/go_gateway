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

func TestService_Update_MarksRuleDerivedMappingOutOfSyncWhenSignatureChanges(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-out-of-sync")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-out-of-sync",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "MIXER",
		Enabled:      true,
	})
	require.NoError(t, err)

	applyRuleManagedLinks(ctx, t, repo, svc, rule.ID)

	mappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	require.Len(t, mappings, 1)
	original := mappings[0]
	originalPipeline := original.TransformPipeline
	originalSignature := original.ProposedSignature

	updatedRule, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{
		ScaleMultiplier:    float64Ptr(2),
		ScaleMultiplierSet: true,
	})
	require.NoError(t, err)
	assert.Equal(t, float64Ptr(2), updatedRule.ScaleMultiplier)

	mappings, err = mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	require.Len(t, mappings, 1)

	current := mappings[0]
	assert.Equal(t, originalPipeline, current.TransformPipeline)
	assert.Equal(t, schema.MappingStatusOutOfSync, current.Status)
	assert.Equal(t, originalSignature, current.LastAppliedSignature)
	assert.NotEqual(t, originalSignature, current.ProposedSignature)
	assert.NotEmpty(t, current.BlockingReason)
}

func TestService_Update_ReactivatesOutOfSyncMappingWhenSignatureMatchesApplied(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-out-of-sync-reset")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-out-of-sync-reset",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "MIXER",
		Enabled:      true,
	})
	require.NoError(t, err)

	applyRuleManagedLinks(ctx, t, repo, svc, rule.ID)

	mappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	require.Len(t, mappings, 1)
	appliedSignature := mappings[0].ProposedSignature

	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{
		ScaleMultiplier:    float64Ptr(2),
		ScaleMultiplierSet: true,
	})
	require.NoError(t, err)

	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{
		ScaleMultiplier:    nil,
		ScaleMultiplierSet: true,
	})
	require.NoError(t, err)

	mappings, err = mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	require.Len(t, mappings, 1)

	current := mappings[0]
	assert.Equal(t, schema.MappingStatusActive, current.Status)
	assert.Equal(t, appliedSignature, current.ProposedSignature)
	assert.Equal(t, appliedSignature, current.LastAppliedSignature)
	assert.Empty(t, current.BlockingReason)
}
