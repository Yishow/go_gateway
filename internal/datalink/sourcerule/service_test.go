package sourcerule

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubRuntimeSync struct {
	upserted []string
	removed  []string
}

func (s *stubRuntimeSync) UpsertPoint(point *schema.Point) {
	if point == nil {
		return
	}
	s.upserted = append(s.upserted, point.ID)
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
	assert.Contains(t, err.Error(), "probe readiness")
}

func seedActiveDevice(ctx context.Context, repo *device.MemoryRepository, id string) (*schema.Device, error) {
	return seedDeviceWithStatus(ctx, repo, id, schema.DeviceStatusActive)
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
