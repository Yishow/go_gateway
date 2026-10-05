package sourcerule

import (
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorkspaceConfirmationOneMappingRestartAndSourceEdit(t *testing.T) {
	ctx := t.Context()
	devices := device.NewMemoryRepository()
	tags := tag.NewService(tag.NewMemoryRepository())
	maps := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tags.GetByID)
	points := point.NewService(point.NewMemoryRepository(), nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, device.NewService(devices, nil), points, nil)
	svc.SetTagMappingServices(tags, maps)
	dev, err := seedActiveDevice(ctx, devices, "workspace-accept")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{ID: "workspace-accept-rule", DeviceID: dev.ID, StartAddress: "40001", Count: 2, DataType: schema.DataTypeInt16, NamingPrefix: "ACCEPT", Enabled: true})
	require.NoError(t, err)
	links := applyRuleManagedLinks(ctx, t, repo, svc, rule.ID)
	require.Len(t, links, 2)
	first, err := maps.GetByID(ctx, *links[0].MappingID)
	require.NoError(t, err)
	second, err := maps.GetByID(ctx, *links[1].MappingID)
	require.NoError(t, err)
	secondAccepted := second.LastAppliedSignature
	require.NoError(t, svc.ValidateWorkspaceMappingSave(ctx, rule, first))
	floatType := schema.DataTypeFloat64
	_, err = tags.Update(ctx, first.TagID, tag.UpdateTagRequest{DataType: &floatType})
	require.NoError(t, err)
	scaled := []schema.TransformStep{{Type: schema.TransformScale, Order: 0, Params: map[string]any{"scale": 0.5, "offset": 10.0}}, {Type: schema.TransformCast, Order: 1, Params: map[string]any{"target_type": "float64"}}}
	saved, err := maps.Update(ctx, first.ID, mapping.UpdateMappingRequest{TransformPipeline: scaled})
	require.NoError(t, err)
	// Another same-rule mapping has an unconfirmed manual edit.
	_, err = maps.Update(ctx, second.ID, mapping.UpdateMappingRequest{TransformPipeline: []schema.TransformStep{{Type: schema.TransformScale, Params: map[string]any{"scale": 3.0, "offset": 0.0}}}})
	require.NoError(t, err)
	require.NoError(t, svc.ConfirmWorkspaceMapping(ctx, rule.ID, rule.RevisionID, saved, floatType))
	confirmed, err := maps.GetByID(ctx, first.ID)
	require.NoError(t, err)
	require.NotEqual(t, confirmed.ProposedSignature, confirmed.LastAppliedSignature)
	for range 2 {
		require.NoError(t, svc.SyncDerivedPointState(ctx))
		current, e := maps.GetByID(ctx, first.ID)
		require.NoError(t, e)
		require.Equal(t, schema.MappingStatusActive, current.Status)
		require.Equal(t, confirmed.LastAppliedSignature, current.LastAppliedSignature)
		require.Equal(t, confirmed.ProposedSignature, current.ProposedSignature)
		require.Equal(t, confirmed.TransformPipeline, current.TransformPipeline)
		tagRecord, e := tags.GetByID(ctx, first.TagID)
		require.NoError(t, e)
		require.Equal(t, floatType, tagRecord.DataType)
	}
	changed, err := maps.GetByID(ctx, second.ID)
	require.NoError(t, err)
	require.Equal(t, schema.MappingStatusOutOfSync, changed.Status)
	require.Equal(t, secondAccepted, changed.LastAppliedSignature, "sync never accepts another manual edit")
	require.ErrorIs(t, svc.ValidateWorkspaceMappingSave(ctx, rule, changed), mapping.ErrWorkspaceConfirmationConflict)
	multiplier := 2.0
	updated, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{ScaleMultiplier: &multiplier, ScaleMultiplierSet: true, TargetDataType: &floatType})
	require.NoError(t, err)
	stale, err := maps.GetByID(ctx, first.ID)
	require.NoError(t, err)
	require.Equal(t, schema.MappingStatusOutOfSync, stale.Status)
	require.NotEmpty(t, stale.BlockingReason)
	require.Equal(t, confirmed.TransformPipeline, stale.TransformPipeline)
	require.ErrorIs(t, svc.ValidateWorkspaceMappingSave(ctx, updated, stale), mapping.ErrWorkspaceConfirmationConflict)
	require.ErrorIs(t, svc.ConfirmWorkspaceMapping(ctx, rule.ID, rule.RevisionID, confirmed, floatType), mapping.ErrWorkspaceConfirmationConflict)
	// Only the existing explicit candidate apply changes the published pipeline.
	currentPoint, err := points.GetByID(ctx, stale.PointID)
	require.NoError(t, err)
	result := newTagMappingSyncResult()
	_, reapplied, err := svc.applyRuleTagMapping(ctx, updated, currentPoint, links[0], "", stale.TagID, &result)
	require.NoError(t, err)
	require.Equal(t, schema.MappingStatusActive, reapplied.Status)
	require.Equal(t, reapplied.ProposedSignature, reapplied.LastAppliedSignature)
	require.NoError(t, svc.SyncDerivedPointState(ctx))
	final, err := maps.GetByID(ctx, first.ID)
	require.NoError(t, err)
	require.Equal(t, schema.MappingStatusActive, final.Status)
	value, err := mapping.ExecutePipeline(int16(243), final.TransformPipeline)
	require.NoError(t, err)
	require.Equal(t, float64(486), value.CurrentValue)
}
