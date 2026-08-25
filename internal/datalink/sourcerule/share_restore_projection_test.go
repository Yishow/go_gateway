package sourcerule

import (
	"context"
	"testing"
	"time"

	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourcerule_RestoreReconstructsExactShareProjectionFromPersistedState(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := deviceSvc.Create(ctx, device.CreateDeviceRequest{
		Name:     "Dev1",
		Protocol: schema.ProtocolModbusTCP,
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     502,
			"slave_id": 1,
		},
	})
	require.NoError(t, err)

	startReg := 40001
	stride := 2
	rule, err := svc.Create(ctx, CreateRuleRequest{
		DeviceID:           dev.ID,
		StartAddress:       "40001",
		Count:              1,
		DataType:           schema.DataTypeFloat32,
		NamingPrefix:       "temp",
		Enabled:            true,
		ShareEnabled:       true,
		ShareStartRegister: &startReg,
		ShareStride:        &stride,
	})
	require.NoError(t, err)

	links := applyRuleManagedLinksWithRepository(t, ctx, repo, svc, rule.ID)
	require.Len(t, links, 1)
	require.NotNil(t, links[0].TagID)

	svc.SetLocalModbusMappingReader(LocalModbusMappingListFunc(func(context.Context) ([]LocalModbusMappingRecord, error) {
		return []LocalModbusMappingRecord{
			{
				TagID:     *links[0].TagID,
				Register:  0,
				DataType:  schema.DataTypeFloat32,
				UpdatedAt: time.Now().UTC(),
			},
		}, nil
	}))
	require.NoError(t, svc.persistCandidateSnapshots(ctx, rule, links))

	// Recreate service (simulating restart)
	restartedSvc := NewService(repo, deviceSvc, pointSvc, nil)
	restartedSvc.SetTagMappingServices(tagSvc, mappingSvc)

	enabledSettings := modbusshare.Settings{Enabled: true, Port: 5020}
	desired, err := restartedSvc.BuildDesiredShareMappingsWithSettings(ctx, "ws-1", enabledSettings)
	require.NoError(t, err)

	require.Len(t, desired, 1)
	assert.Equal(t, *links[0].TagID, desired[0].TagID)
	assert.Equal(t, uint16(0), desired[0].ZeroBasedRegister)
	assert.Equal(t, 2, desired[0].SpanRegisters)
	assert.Equal(t, schema.DataTypeFloat32, desired[0].DataType)
}
