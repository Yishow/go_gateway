package sourcerule

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestService_Create_InvalidAddressErrorIncludesRuleAddressAndProtocol(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	svc := NewService(NewMemoryRepository(), deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-create-address-context")
	require.NoError(t, err)

	_, err = svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-create-address-context",
		DeviceID:     dev.ID,
		StartAddress: "40O01",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC_",
		Enabled:      true,
	})

	require.ErrorIs(t, err, ErrValidation)
	require.Error(t, err)
	for _, fragment := range []string{
		"rule_id=rule-create-address-context",
		"start_address=40O01",
		"protocol=modbus_tcp",
	} {
		require.Contains(t, err.Error(), fragment)
	}
}

func TestService_Update_InvalidAddressErrorIncludesRuleAddressAndProtocol(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-update-address-context")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-update-address-context",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC_",
		Enabled:      true,
	})
	require.NoError(t, err)

	invalidAddress := "40O01"
	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{StartAddress: &invalidAddress})

	require.ErrorIs(t, err, ErrValidation)
	require.Error(t, err)
	for _, fragment := range []string{
		"rule_id=rule-update-address-context",
		"start_address=40O01",
		"protocol=modbus_tcp",
	} {
		require.Contains(t, err.Error(), fragment)
	}
}

func TestService_Create_MissingDeviceErrorIncludesDeviceIDAndReason(t *testing.T) {
	ctx := context.Background()
	deviceSvc := device.NewService(device.NewMemoryRepository(), nil)
	svc := NewService(NewMemoryRepository(), deviceSvc, point.NewService(point.NewMemoryRepository(), nil), nil)

	_, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-missing-device-context",
		DeviceID:     "device-does-not-exist",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC_",
		Enabled:      true,
	})

	require.Error(t, err)
	require.False(t, errors.Is(err, ErrValidation))
	for _, fragment := range []string{
		"device_id=device-does-not-exist",
		"device not found",
	} {
		require.Contains(t, err.Error(), fragment)
	}
}
