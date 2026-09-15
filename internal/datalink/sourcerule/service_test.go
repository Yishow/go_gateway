package sourcerule

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubRuntimeSync struct {
	upserted []string
	removed  []string
}

type failingLinkRepository struct {
	*MemoryRepository
	failCreateLinks bool
}

func (r *failingLinkRepository) CreateLinks(ctx context.Context, links []*schema.SourceRuleLink) error {
	if r.failCreateLinks {
		r.failCreateLinks = false
		return errors.New("create links failed")
	}
	return r.MemoryRepository.CreateLinks(ctx, links)
}

func (s *stubRuntimeSync) UpsertPoint(pointRecord *schema.Point) {
	if pointRecord == nil {
		return
	}
	s.upserted = append(s.upserted, pointRecord.ID)
}

func (s *stubRuntimeSync) RemovePoint(pointID string) {
	s.removed = append(s.removed, pointID)
}

func TestService_Create_PersistsRuleAndDerivedPoints(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	runtimeSync := &stubRuntimeSync{}
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, runtimeSync)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-1")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-1",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "MIXER",
		Enabled:      true,
	})
	require.NoError(t, err)
	assert.Equal(t, "rule-1", rule.ID)
	assert.True(t, rule.Enabled)

	rules, err := svc.List(ctx, ListFilter{DeviceID: &dev.ID})
	require.NoError(t, err)
	require.Len(t, rules, 1)

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	assert.Len(t, links, 2)

	points, err := pointSvc.List(ctx, point.ListFilter{DeviceID: &dev.ID})
	require.NoError(t, err)
	require.Len(t, points, 2)
	pointNames := []string{points[0].Name, points[1].Name}
	assert.ElementsMatch(t, []string{"MIXER_40001", "MIXER_40002"}, pointNames)
	assert.True(t, points[0].Enabled)
	assert.True(t, points[1].Enabled)
	assert.NotEmpty(t, runtimeSync.upserted)
}

func TestService_Create_AssignsDefaultPollingGroupToDerivedPoints(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	groupRepo := pollinggroup.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, groupRepo)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-1")
	require.NoError(t, err)

	_, err = svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-1",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "MIXER",
		Enabled:      true,
	})
	require.NoError(t, err)

	points, err := pointSvc.List(ctx, point.ListFilter{DeviceID: &dev.ID})
	require.NoError(t, err)
	require.Len(t, points, 2)
	require.NotNil(t, points[0].PollingGroupID)
	require.NotNil(t, points[1].PollingGroupID)
	assert.Equal(t, *points[0].PollingGroupID, *points[1].PollingGroupID)
}

func TestService_Disable_PreservesDerivedPointsButStopsCollection(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-1")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-1",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "MIXER",
		Enabled:      true,
	})
	require.NoError(t, err)

	require.NoError(t, svc.Disable(ctx, rule.ID))

	updatedRule, err := svc.GetByID(ctx, rule.ID)
	require.NoError(t, err)
	assert.False(t, updatedRule.Enabled)

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	assert.Len(t, links, 2)

	points, err := pointSvc.List(ctx, point.ListFilter{DeviceID: &dev.ID})
	require.NoError(t, err)
	require.Len(t, points, 2)
	assert.False(t, points[0].Enabled)
	assert.False(t, points[1].Enabled)
}

func TestService_RestoreDerivedPointState_ReappliesPersistedRuleEnablement(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-1")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-1",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "MIXER",
		Enabled:      false,
	})
	require.NoError(t, err)

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)

	pointID := links[0].PointID
	enabled := true
	_, err = pointSvc.Update(ctx, pointID, point.UpdatePointRequest{Enabled: &enabled})
	require.NoError(t, err)

	require.NoError(t, svc.SyncDerivedPointState(ctx))

	derivedPoint, err := pointSvc.GetByID(ctx, pointID)
	require.NoError(t, err)
	assert.False(t, derivedPoint.Enabled)
}

func TestService_Update_ReconcilesDerivedPoints(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-1")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-1",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	startAddress := "40003"
	count := 3
	dataType := schema.DataTypeInt32
	namingPrefix := "MIXER"
	updatedRule, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{
		StartAddress: &startAddress,
		Count:        &count,
		DataType:     &dataType,
		NamingPrefix: &namingPrefix,
	})
	require.NoError(t, err)
	assert.Equal(t, "40003", updatedRule.StartAddress)
	assert.Equal(t, 3, updatedRule.Count)
	assert.Equal(t, schema.DataTypeInt32, updatedRule.DataType)

	deviceID := dev.ID
	points, err := pointSvc.List(ctx, point.ListFilter{DeviceID: &deviceID})
	require.NoError(t, err)
	require.Len(t, points, 3)

	addresses := []string{points[0].Address, points[1].Address, points[2].Address}
	assert.ElementsMatch(t, []string{"40003", "40005", "40007"}, addresses)
}

func TestService_Update_RecreatesWhenDerivedPointManuallyDeleted(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-1")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-1",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)
	deadPointID := links[0].PointID

	require.NoError(t, pointSvc.Delete(ctx, deadPointID))

	namingPrefix := "NEW"
	updated, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: &namingPrefix})
	require.NoError(t, err)
	assert.Equal(t, "NEW", updated.NamingPrefix)

	links2, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links2, 1)
	assert.NotEqual(t, deadPointID, links2[0].PointID)

	p, err := pointSvc.GetByID(ctx, links2[0].PointID)
	require.NoError(t, err)
	assert.Equal(t, "NEW_40001", p.Name)
}

func TestService_Enable_BlocksWhenDeviceNotActive(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedDeviceWithStatus(ctx, deviceRepo, "device-1", schema.DeviceStatusDraft)
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-1",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)
	assert.False(t, rule.Enabled)

	err = svc.Enable(ctx, rule.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "activation readiness")
}

func TestService_SyncDerivedPointState_DoesNotBackfillTagMappingsForPendingLinks(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	legacySvc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-1")
	require.NoError(t, err)

	rule, err := legacySvc.Create(ctx, CreateRuleRequest{
		ID:           "rule-legacy",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "MIXER",
		Enabled:      false,
	})
	require.NoError(t, err)

	legacyLinks, err := legacySvc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, legacyLinks, 1)
	assert.Nil(t, legacyLinks[0].TagID)
	assert.Nil(t, legacyLinks[0].MappingID)

	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	syncSvc := NewService(repo, deviceSvc, pointSvc, nil)
	syncSvc.SetTagMappingServices(tagSvc, mappingSvc)

	require.NoError(t, syncSvc.SyncDerivedPointState(ctx))

	syncedLinks, err := syncSvc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, syncedLinks, 1)
	assert.Nil(t, syncedLinks[0].TagID)
	assert.Nil(t, syncedLinks[0].MappingID)

	mappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, mappings)

	tags, err := tagSvc.List(ctx, tag.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestService_Create_AllowsCandidateTagKeyCollisionUntilApply(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-1")
	require.NoError(t, err)

	_, err = tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:      "MIXER_40001",
		DataType: schema.DataTypeFloat32,
	})
	require.NoError(t, err)

	_, err = svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-collision",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "MIXER",
		Enabled:      true,
	})
	require.NoError(t, err)

	rules, err := svc.List(ctx, ListFilter{DeviceID: &dev.ID})
	require.NoError(t, err)
	require.Len(t, rules, 1)

	points, err := pointSvc.List(ctx, point.ListFilter{DeviceID: &dev.ID})
	require.NoError(t, err)
	require.Len(t, points, 1)

	mappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, mappings)

	tags, err := tagSvc.List(ctx, tag.ListFilter{})
	require.NoError(t, err)
	require.Len(t, tags, 1)
	assert.Equal(t, schema.DataTypeFloat32, tags[0].DataType)
}

func TestService_Update_ShrunkRangeKeepsCandidatesPendingUntilApply(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-1")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-shrink",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "MIXER",
		Enabled:      true,
	})
	require.NoError(t, err)

	count := 1
	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{Count: &count})
	require.NoError(t, err)

	tags, err := tagSvc.List(ctx, tag.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, tags)

	mappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, mappings)

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)
	assert.Nil(t, links[0].TagID)
	assert.Nil(t, links[0].MappingID)
}

func TestService_Update_RollsBackBackfilledPollingGroupOnFailure(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	legacyPointSvc := point.NewService(pointRepo, nil)
	baseRepo := NewMemoryRepository()
	legacySvc := NewService(baseRepo, deviceSvc, legacyPointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-rollback-update")
	require.NoError(t, err)

	rule, err := legacySvc.Create(ctx, CreateRuleRequest{
		ID:           "rule-rollback-update",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	links, err := legacySvc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)

	legacyPoint, err := legacyPointSvc.GetByID(ctx, links[0].PointID)
	require.NoError(t, err)
	require.Nil(t, legacyPoint.PollingGroupID)

	groupRepo := pollinggroup.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, groupRepo)
	repo := &failingLinkRepository{MemoryRepository: baseRepo, failCreateLinks: true}
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	namingPrefix := "MIXER"
	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: &namingPrefix})
	require.Error(t, err)

	rolledBackPoint, err := pointSvc.GetByID(ctx, links[0].PointID)
	require.NoError(t, err)
	assert.Equal(t, "SRC_40001", rolledBackPoint.Name)
	assert.Nil(t, rolledBackPoint.PollingGroupID)
}

func TestService_Enable_RollsBackBackfilledPollingGroupOnFailure(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	legacyPointSvc := point.NewService(pointRepo, nil)
	baseRepo := NewMemoryRepository()
	legacySvc := NewService(baseRepo, deviceSvc, legacyPointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-rollback-enable")
	require.NoError(t, err)

	rule, err := legacySvc.Create(ctx, CreateRuleRequest{
		ID:           "rule-rollback-enable",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      false,
	})
	require.NoError(t, err)

	links, err := legacySvc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)

	legacyPoint, err := legacyPointSvc.GetByID(ctx, links[0].PointID)
	require.NoError(t, err)
	require.False(t, legacyPoint.Enabled)
	require.Nil(t, legacyPoint.PollingGroupID)

	groupRepo := pollinggroup.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, groupRepo)
	repo := &failingLinkRepository{MemoryRepository: baseRepo, failCreateLinks: true}
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	err = svc.Enable(ctx, rule.ID)
	require.Error(t, err)

	rolledBackPoint, err := pointSvc.GetByID(ctx, links[0].PointID)
	require.NoError(t, err)
	assert.False(t, rolledBackPoint.Enabled)
	assert.Nil(t, rolledBackPoint.PollingGroupID)
}

func TestService_Delete_RemovesOrphanedAutoTags(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-1")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-delete",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "MIXER",
		Enabled:      true,
	})
	require.NoError(t, err)

	require.NoError(t, svc.Delete(ctx, rule.ID))

	tags, err := tagSvc.List(ctx, tag.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, tags)

	mappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	assert.Empty(t, mappings)
}

func seedActiveDevice(ctx context.Context, repo *device.MemoryRepository, id string) (*schema.Device, error) {
	return seedActivationReadyDevice(ctx, repo, id)
}

func seedDeviceWithStatus(ctx context.Context, repo *device.MemoryRepository, id string, status schema.DeviceStatus) (*schema.Device, error) {
	record := &schema.Device{
		ID:               id,
		Name:             "Mixer PLC",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           status,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1}`,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := repo.Create(ctx, record); err != nil {
		return nil, err
	}
	return record, nil
}

// =============================================================================
// Target Data Type 與 Scale 相關測試
// =============================================================================

func TestService_Create_WithTargetDataType_PublishesCandidateWithTargetType(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	runtimeSync := &stubRuntimeSync{}
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, runtimeSync)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-target-dtype")
	require.NoError(t, err)

	// 建立 uint16 規則，目標型別為 float64
	targetType := schema.DataTypeFloat64
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:             "rule-target-dtype",
		DeviceID:       dev.ID,
		StartAddress:   "40001",
		Count:          1,
		DataType:       schema.DataTypeUint16,
		NamingPrefix:   "VAL",
		Enabled:        true,
		TargetDataType: &targetType,
	})
	require.NoError(t, err)
	assert.Equal(t, schema.DataTypeUint16, rule.DataType)
	require.NotNil(t, rule.TargetDataType)
	assert.Equal(t, schema.DataTypeFloat64, *rule.TargetDataType)

	// 驗證 Point 使用協議型別 uint16
	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)

	pointRecord, err := pointSvc.GetByID(ctx, links[0].PointID)
	require.NoError(t, err)
	assert.Equal(t, schema.DataTypeUint16, pointRecord.DataType)

	assert.Nil(t, links[0].TagID)
	assert.Nil(t, links[0].MappingID)

	candidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	assert.Equal(t, schema.DataTypeFloat64, candidate.DataType)
	require.Len(t, candidate.TransformPipeline, 1)
	assert.Equal(t, schema.TransformCast, candidate.TransformPipeline[0].Type)
}

func TestService_Create_WithScale_PublishesScalePipelineCandidate(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	runtimeSync := &stubRuntimeSync{}
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, runtimeSync)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-scale")
	require.NoError(t, err)

	// 建立帶 scale 的規則
	multiplier := 0.1
	offset := 0.0
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:              "rule-scale",
		DeviceID:        dev.ID,
		StartAddress:    "40001",
		Count:           1,
		DataType:        schema.DataTypeUint16,
		NamingPrefix:    "TEMP",
		Enabled:         true,
		ScaleMultiplier: &multiplier,
		ScaleOffset:     &offset,
	})
	require.NoError(t, err)
	require.NotNil(t, rule.ScaleMultiplier)
	assert.Equal(t, 0.1, *rule.ScaleMultiplier)

	// 驗證 Mapping 包含 scale 步驟
	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)
	assert.Nil(t, links[0].TagID)
	assert.Nil(t, links[0].MappingID)

	candidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, candidate.TransformPipeline, 1)
	assert.Equal(t, schema.TransformScale, candidate.TransformPipeline[0].Type)
}

func TestService_Create_WithTargetTypeAndScale_PublishesCastThenScaleCandidate(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	runtimeSync := &stubRuntimeSync{}
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, runtimeSync)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-cast-scale")
	require.NoError(t, err)

	// 建立 uint16→float64 且帶 scale 的規則
	targetType := schema.DataTypeFloat64
	multiplier := 0.01
	offset := -273.15
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:              "rule-cast-scale",
		DeviceID:        dev.ID,
		StartAddress:    "40001",
		Count:           1,
		DataType:        schema.DataTypeUint16,
		NamingPrefix:    "KELVIN",
		Enabled:         true,
		TargetDataType:  &targetType,
		ScaleMultiplier: &multiplier,
		ScaleOffset:     &offset,
	})
	require.NoError(t, err)

	// 驗證 Mapping 包含 cast + scale 兩個步驟
	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)
	assert.Nil(t, links[0].TagID)
	assert.Nil(t, links[0].MappingID)

	candidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	steps := candidate.TransformPipeline
	require.Len(t, steps, 2)

	// 驗證順序：先 cast，再 scale
	assert.Equal(t, schema.TransformCast, steps[0].Type)
	assert.Equal(t, schema.TransformScale, steps[1].Type)

	// 驗證 cast 參數
	assert.Equal(t, "float64", steps[0].Params["target_type"])

	// 驗證 scale 參數
	assert.InDelta(t, 0.01, steps[1].Params["scale"], 0.0001)
	assert.InDelta(t, -273.15, steps[1].Params["offset"], 0.0001)
}

func TestService_Update_WithTargetDataType_UpdatesTagType(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	runtimeSync := &stubRuntimeSync{}
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, runtimeSync)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-update-target")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-update-target",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeUint16,
		NamingPrefix: "VAL",
		Enabled:      true,
	})
	require.NoError(t, err)
	links := applyRuleManagedLinks(ctx, t, repo, svc, rule.ID)
	tagRecord, err := tagSvc.GetByID(ctx, *links[0].TagID)
	require.NoError(t, err)
	assert.Equal(t, schema.DataTypeUint16, tagRecord.DataType)
	initialMapping, err := mappingSvc.GetByID(ctx, *links[0].MappingID)
	require.NoError(t, err)
	initialSignature := initialMapping.ProposedSignature
	targetType := schema.DataTypeFloat64
	updatedRule, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{
		TargetDataType: &targetType,
	})
	require.NoError(t, err)
	require.NotNil(t, updatedRule.TargetDataType)
	assert.Equal(t, schema.DataTypeFloat64, *updatedRule.TargetDataType)
	savedRule, err := svc.GetByID(ctx, rule.ID)
	require.NoError(t, err)
	require.NotNil(t, savedRule.TargetDataType)
	assert.Equal(t, schema.DataTypeFloat64, *savedRule.TargetDataType)
	tagRecord, err = tagSvc.GetByID(ctx, *links[0].TagID)
	require.NoError(t, err)
	assert.Equal(t, schema.DataTypeFloat64, tagRecord.DataType)
	mappingRecord, err := mappingSvc.GetByID(ctx, *links[0].MappingID)
	require.NoError(t, err)
	var steps []schema.TransformStep
	require.NoError(t, json.Unmarshal([]byte(mappingRecord.TransformPipeline), &steps))
	assert.Empty(t, steps)
	assert.Equal(t, schema.MappingStatusOutOfSync, mappingRecord.Status)
	assert.Equal(t, initialSignature, mappingRecord.LastAppliedSignature)
	assert.NotEqual(t, initialSignature, mappingRecord.ProposedSignature)
	assert.NotEmpty(t, mappingRecord.BlockingReason)
}

func TestService_Update_ClearConversionSettings_RemovesTagCastAndScale(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	runtimeSync := &stubRuntimeSync{}
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, runtimeSync)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-clear-conversion")
	require.NoError(t, err)

	targetType := schema.DataTypeFloat64
	multiplier := 0.1
	offset := 2.5
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:              "rule-clear-conversion",
		DeviceID:        dev.ID,
		StartAddress:    "40001",
		Count:           1,
		DataType:        schema.DataTypeUint16,
		NamingPrefix:    "TEMP",
		Enabled:         true,
		TargetDataType:  &targetType,
		ScaleMultiplier: &multiplier,
		ScaleOffset:     &offset,
	})
	require.NoError(t, err)
	links := applyRuleManagedLinks(ctx, t, repo, svc, rule.ID)
	initialMapping, err := mappingSvc.GetByID(ctx, *links[0].MappingID)
	require.NoError(t, err)
	initialPipeline := initialMapping.TransformPipeline
	initialSignature := initialMapping.ProposedSignature
	updatedRule, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{
		TargetDataTypeSet:  true,
		ScaleMultiplierSet: true,
		ScaleOffsetSet:     true,
	})
	require.NoError(t, err)
	assert.Nil(t, updatedRule.TargetDataType)
	assert.Nil(t, updatedRule.ScaleMultiplier)
	assert.Nil(t, updatedRule.ScaleOffset)
	savedRule, err := svc.GetByID(ctx, rule.ID)
	require.NoError(t, err)
	assert.Nil(t, savedRule.TargetDataType)
	assert.Nil(t, savedRule.ScaleMultiplier)
	assert.Nil(t, savedRule.ScaleOffset)
	tagRecord, err := tagSvc.GetByID(ctx, *links[0].TagID)
	require.NoError(t, err)
	assert.Equal(t, schema.DataTypeUint16, tagRecord.DataType)
	mappingRecord, err := mappingSvc.GetByID(ctx, *links[0].MappingID)
	require.NoError(t, err)
	assert.Equal(t, initialPipeline, mappingRecord.TransformPipeline)
	assert.Equal(t, schema.MappingStatusOutOfSync, mappingRecord.Status)
	assert.Equal(t, initialSignature, mappingRecord.LastAppliedSignature)
	assert.NotEqual(t, initialSignature, mappingRecord.ProposedSignature)
	assert.NotEmpty(t, mappingRecord.BlockingReason)
}
