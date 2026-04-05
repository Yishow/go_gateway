package sourcerule

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_SQLRevisionPersistenceRestoresLatestRevisionAfterRestart(t *testing.T) {
	ctx := context.Background()
	db := setupSQLRepoDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	groupRepo := pollinggroup.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, deviceRepo, "device-1"))

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, groupRepo)
	svc := NewService(NewSQLRepository(db), deviceSvc, pointSvc, nil)

	created, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-1",
		DeviceID:     "device-1",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)
	require.NotEmpty(t, created.RevisionID)

	initialRevision := created.RevisionID
	updatedPrefix := "UPDATED"
	enabled := false
	updated, err := svc.Update(ctx, created.ID, UpdateRuleRequest{
		NamingPrefix: &updatedPrefix,
		Enabled:      &enabled,
	})
	require.NoError(t, err)
	require.NotEmpty(t, updated.RevisionID)
	assert.NotEqual(t, initialRevision, updated.RevisionID)
	assert.False(t, updated.Enabled)

	restartedDeviceSvc := device.NewService(device.NewSQLRepository(db), nil)
	restartedPointSvc := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	restartedSvc := NewService(NewSQLRepository(db), restartedDeviceSvc, restartedPointSvc, nil)

	restored, err := restartedSvc.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, updated.RevisionID, restored.RevisionID)
	assert.Equal(t, updatedPrefix, restored.NamingPrefix)
	assert.False(t, restored.Enabled)
}

func TestService_SetEnabled_PersistsStateWithoutChangingRevisionID(t *testing.T) {
	ctx := context.Background()
	db := setupSQLRepoDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	groupRepo := pollinggroup.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, deviceRepo, "device-1"))

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, groupRepo)
	svc := NewService(NewSQLRepository(db), deviceSvc, pointSvc, nil)

	created, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-activation",
		DeviceID:     "device-1",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)
	initialRevision := created.RevisionID

	require.NoError(t, svc.Disable(ctx, created.ID))

	restartedDeviceSvc := device.NewService(device.NewSQLRepository(db), nil)
	restartedPointSvc := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	restartedSvc := NewService(NewSQLRepository(db), restartedDeviceSvc, restartedPointSvc, nil)

	restored, err := restartedSvc.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, initialRevision, restored.RevisionID)
	assert.False(t, restored.Enabled)
}

func TestService_ListRestoresLatestRevisionMetadataAfterRestart(t *testing.T) {
	ctx := context.Background()
	db := setupSQLRepoDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	groupRepo := pollinggroup.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, deviceRepo, "device-list"))

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, groupRepo)
	svc := NewService(NewSQLRepository(db), deviceSvc, pointSvc, nil)

	created, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-list-restore",
		DeviceID:     "device-list",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	updatedPrefix := "RESTARTED"
	enabled := false
	updated, err := svc.Update(ctx, created.ID, UpdateRuleRequest{
		NamingPrefix: &updatedPrefix,
		Enabled:      &enabled,
	})
	require.NoError(t, err)

	restartedSvc := NewService(NewSQLRepository(db), device.NewService(device.NewSQLRepository(db), nil), point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db)), nil)
	deviceID := "device-list"
	rules, err := restartedSvc.List(ctx, ListFilter{DeviceID: &deviceID})
	require.NoError(t, err)
	require.Len(t, rules, 1)
	assert.Equal(t, updated.ID, rules[0].ID)
	assert.Equal(t, updated.RevisionID, rules[0].RevisionID)
	assert.Equal(t, updatedPrefix, rules[0].NamingPrefix)
	assert.False(t, rules[0].Enabled)
}

func seedSQLSourceRuleDevice(ctx context.Context, repo *device.SQLRepository, id string) error {
	now := time.Now().UTC()
	return repo.Create(ctx, &schema.Device{
		ID:               id,
		Name:             "Mixer PLC",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	})
}
