package sourcerule

import (
	"context"
	"testing"

	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourcerule_RuleFlagCannotOverrideDisabledGlobalShare(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	svc.SetTagMappingServices(tagSvc, nil)

	dev, err := deviceSvc.Create(ctx, device.CreateDeviceRequest{
		Name:     "Test Device",
		Protocol: schema.ProtocolModbusTCP,
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     502,
			"slave_id": 1,
		},
	})
	require.NoError(t, err)

	// Create a rule with share_enabled = true
	startReg := 40001
	stride := 2
	rule, err := svc.Create(ctx, CreateRuleRequest{
		DeviceID:           dev.ID,
		StartAddress:       "40001",
		Count:              1,
		DataType:           schema.DataTypeInt16,
		NamingPrefix:       "test",
		Enabled:            true,
		ShareEnabled:       true,
		ShareStartRegister: &startReg,
		ShareStride:        &stride,
	})
	require.NoError(t, err)
	require.True(t, rule.ShareEnabled)

	// When global share is disabled
	disabledSettings := modbusshare.Settings{Enabled: false}
	desiredMappings, err := svc.BuildDesiredShareMappingsWithSettings(ctx, "ws-1", disabledSettings)
	require.NoError(t, err)

	// Result must be empty because global share is disabled!
	assert.Empty(t, desiredMappings, "rule flag must not override disabled global setting")
}
