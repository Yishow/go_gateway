package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestOffsetAddress_MC3EHexContactUsesHexadecimalRadix(t *testing.T) {
	got, err := offsetAddress("X0", 16, schema.ProtocolMC3E)
	require.NoError(t, err)
	require.Equal(t, "X10", got)
}

func TestOffsetAddress_ProtocolRadixBoundarySpecTable(t *testing.T) {
	tests := []struct {
		protocol schema.ProtocolType
		start    string
		offset   int
		expected string
	}{
		{protocol: schema.ProtocolMC3E, start: "X0", offset: 16, expected: "X10"},
		{protocol: schema.ProtocolMC3E, start: "D0", offset: 16, expected: "D16"},
		{protocol: schema.ProtocolFatekFBs, start: "R0", offset: 16, expected: "R16"},
		{protocol: schema.ProtocolModbusTCP, start: "40001", offset: 2, expected: "40003"},
	}

	for _, tt := range tests {
		got, err := offsetAddress(tt.start, tt.offset, tt.protocol)
		require.NoError(t, err)
		require.Equal(t, tt.expected, got)
	}
}

func TestOffsetAddress_InvalidMC3EAddressDoesNotFallBackToOriginal(t *testing.T) {
	got, err := offsetAddress("Z999", 1, schema.ProtocolMC3E)
	require.Error(t, err)
	require.Empty(t, got)
}

func TestBuildPlannedPointAddresses_InvalidMC3EAddressFailsClosed(t *testing.T) {
	got, err := buildPlannedPointAddresses("Z999", 1, schema.DataTypeInt16, schema.ProtocolMC3E)
	require.Error(t, err)
	require.Empty(t, got)
}

func TestServiceCreate_RejectsInvalidMC3EAddressBeforePersistence(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-mc-invalid")
	require.NoError(t, err)
	dev.Protocol = schema.ProtocolMC3E
	require.NoError(t, deviceRepo.Update(ctx, dev))

	_, err = svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-invalid-address",
		DeviceID:     dev.ID,
		StartAddress: "Z999",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "PLC_",
		Enabled:      true,
	})
	require.ErrorIs(t, err, ErrValidation)

	rules, listErr := repo.List(ctx, ListFilter{DeviceID: &dev.ID})
	require.NoError(t, listErr)
	require.Empty(t, rules)
}
