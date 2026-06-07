package runtime

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestService_StartUsesWorkspaceProjectionInsteadOfAllActiveDevices(t *testing.T) {
	ctx := context.Background()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, connMgr)
	pointRepo := point.NewMemoryRepository()
	groupRepo := pollinggroup.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, groupRepo)
	mappingRepo := mapping.NewMemoryRepository()
	mappingSvc := mapping.NewService(mappingRepo)
	tagRepo := tag.NewMemoryRepository()
	tagSvc := tag.NewService(tagRepo)
	groupSvc := pollinggroup.NewService(groupRepo)
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), connMgr)

	projectedDevice := &schema.Device{
		ID:               "dev-A",
		Name:             "Projected device",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":1}`,
	}
	outsideDevice := &schema.Device{
		ID:               "dev-Z",
		Name:             "Outside active device",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":1}`,
	}
	require.NoError(t, deviceRepo.Create(ctx, projectedDevice))
	require.NoError(t, deviceRepo.Create(ctx, outsideDevice))

	groupRecord := &schema.PollingGroup{ID: "group-fast", Name: "fast", IntervalMs: 60000, Enabled: true}
	require.NoError(t, groupRepo.Create(ctx, groupRecord))
	pointRecord := &schema.Point{
		ID:             "point-A",
		DeviceID:       projectedDevice.ID,
		Name:           "Projected point",
		Address:        "40001",
		DataType:       schema.DataTypeInt16,
		Mode:           schema.PointModeReadOnly,
		PollingGroupID: &groupRecord.ID,
		Enabled:        true,
	}
	require.NoError(t, pointRepo.Create(ctx, pointRecord))
	tagRecord := &schema.Tag{
		ID:          "tag-A",
		Key:         "projected.a",
		KeyLower:    "projected.a",
		DisplayName: "Projected A",
		DataType:    schema.DataTypeInt16,
		Status:      schema.TagStatusActive,
	}
	require.NoError(t, tagRepo.Create(ctx, tagRecord))
	mappingRecord := &schema.Mapping{
		ID:                "map-A",
		PointID:           pointRecord.ID,
		TagID:             tagRecord.ID,
		TransformPipeline: "[]",
		Enabled:           true,
	}
	require.NoError(t, mappingRepo.Create(ctx, mappingRecord))

	projection := &workspace.RuntimeProjection{
		WorkspaceID:   "ws-1",
		Version:       "projection-v1",
		Alignment:     workspace.RuntimeProjectionAlignmentAligned,
		DeviceIDs:     []string{projectedDevice.ID},
		Devices:       []*schema.Device{projectedDevice},
		Points:        []*schema.Point{pointRecord},
		PollingGroups: []*schema.PollingGroup{groupRecord},
		Mappings:      []*schema.Mapping{mappingRecord},
		Tags:          []*schema.Tag{tagRecord},
	}
	runtimeSvc, err := NewService(DefaultConfig(), Dependencies{
		Scheduler:           scheduler,
		Writer:              storage.NewMemoryStorage(16),
		DeviceService:       deviceSvc,
		PointService:        pointSvc,
		MappingService:      mappingSvc,
		TagService:          tagSvc,
		PollingGroupService: groupSvc,
		WorkspaceProjection: workspaceProjectionStub{projection: projection},
	})
	require.NoError(t, err)

	require.NoError(t, runtimeSvc.Start(ctx))
	defer func() {
		require.NoError(t, runtimeSvc.Stop(context.Background()))
	}()

	_, projectedExists := scheduler.GetDeviceBreakerState(projectedDevice.ID)
	_, outsideExists := scheduler.GetDeviceBreakerState(outsideDevice.ID)
	require.True(t, projectedExists)
	require.False(t, outsideExists)
}

type workspaceProjectionStub struct {
	projection *workspace.RuntimeProjection
}

func (s workspaceProjectionStub) RuntimeProjection(context.Context) (*workspace.RuntimeProjection, error) {
	return s.projection, nil
}
