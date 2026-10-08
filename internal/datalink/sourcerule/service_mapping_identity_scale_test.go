package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A rule saved with the neutral scale (x1 + 0) describes the same mapping as
// one with no scale step. A restart re-syncs every rule, and that sync must
// not turn an unchanged mapping into out_of_sync: its revision feeds every
// write group that uses the tag, so the drift stops the group's data.
func TestService_SyncDerivedPointState_KeepsMappingActiveForNeutralScaleRule(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-neutral-scale")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID: "rule-neutral-scale", DeviceID: dev.ID, StartAddress: "40001", Count: 1,
		DataType: schema.DataTypeInt16, NamingPrefix: "MIXER", Enabled: true,
		ScaleMultiplier: common.Ptr(float64(1)), ScaleOffset: common.Ptr(float64(0)),
	})
	require.NoError(t, err)
	applyRuleManagedLinks(ctx, t, repo, svc, rule.ID)

	// The mapping the operator confirmed carries no transform at all.
	mappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	require.Len(t, mappings, 1)
	_, err = mappingSvc.Update(ctx, mappings[0].ID, mapping.UpdateMappingRequest{TransformPipeline: []schema.TransformStep{}})
	require.NoError(t, err)
	before, err := mappingSvc.GetByID(ctx, mappings[0].ID)
	require.NoError(t, err)
	require.Equal(t, schema.MappingStatusActive, before.Status)

	require.NoError(t, svc.SyncDerivedPointState(ctx))

	after, err := mappingSvc.GetByID(ctx, mappings[0].ID)
	require.NoError(t, err)
	assert.Equal(t, schema.MappingStatusActive, after.Status, "a restart sync must not flip a neutral-scale mapping to out_of_sync")
	assert.Empty(t, after.BlockingReason)
	assert.Equal(t, before.TransformPipeline, after.TransformPipeline, "the confirmed pipeline must not be rewritten")
}

func TestMappingCandidateSignature_TreatsNeutralScaleAsNoStep(t *testing.T) {
	neutral := []schema.TransformStep{{Type: schema.TransformScale, Order: 0, Params: map[string]any{"scale": 1.0, "offset": 0.0}}}
	scaled := []schema.TransformStep{{Type: schema.TransformScale, Order: 0, Params: map[string]any{"scale": 2.0, "offset": 0.0}}}

	none, err := mappingCandidateSignature([]schema.TransformStep{})
	require.NoError(t, err)
	nilSteps, err := mappingCandidateSignature(nil)
	require.NoError(t, err)
	withNeutral, err := mappingCandidateSignature(neutral)
	require.NoError(t, err)
	withScale, err := mappingCandidateSignature(scaled)
	require.NoError(t, err)

	assert.Equal(t, none, nilSteps)
	assert.Equal(t, none, withNeutral)
	assert.NotEqual(t, none, withScale, "a real scale must still change the signature")
}

// Mappings saved from the workspace carry no rule-candidate metadata. The first
// sync adopts them; doing that when they are saved means a restart changes
// nothing a write group is bound to.
func TestService_SyncRuleDerivedState_AdoptsWorkspaceMappingSoRestartChangesNothing(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-adopt")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID: "rule-adopt", DeviceID: dev.ID, StartAddress: "40001", Count: 1,
		DataType: schema.DataTypeInt16, NamingPrefix: "MIXER", Enabled: true,
		ScaleMultiplier: common.Ptr(float64(1)), ScaleOffset: common.Ptr(float64(0)),
	})
	require.NoError(t, err)
	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)

	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key: "line.adopt", DisplayName: "Adopt", DataType: schema.DataTypeInt16,
		Labels: map[string]string{"source": "source-rule", "source_rule_id": rule.ID, "source_rule_address": "40001"},
	})
	require.NoError(t, err)
	enabled := true
	created, err := mappingSvc.Create(ctx, mapping.CreateMappingRequest{PointID: links[0].PointID, TagID: tagRecord.ID, Enabled: &enabled, TransformPipeline: []schema.TransformStep{}})
	require.NoError(t, err)
	links[0].TagID, links[0].MappingID = &tagRecord.ID, &created.ID
	require.NoError(t, svc.ReplaceLinks(ctx, rule.ID, links))
	require.Empty(t, created.RuleCandidateID)

	require.NoError(t, svc.SyncRuleDerivedState(ctx, rule.ID))
	settled, err := mappingSvc.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, settled.RuleCandidateID, "the saved mapping is adopted straight away")
	assert.Equal(t, schema.MappingStatusActive, settled.Status)

	revisionBefore, err := measurementMappingRevision(settled)
	require.NoError(t, err)
	require.NoError(t, svc.SyncDerivedPointState(ctx)) // what a restart does
	afterRestart, err := mappingSvc.GetByID(ctx, created.ID)
	require.NoError(t, err)
	revisionAfter, err := measurementMappingRevision(afterRestart)
	require.NoError(t, err)
	assert.Equal(t, revisionBefore, revisionAfter, "a restart must not change the mapping revision write groups bind to")
}

func measurementMappingRevision(record *schema.Mapping) (string, error) {
	return measurement.MappingRevision(*record)
}
