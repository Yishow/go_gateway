package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestService_ListByDeviceIDsPreservesWorkspaceDeviceOwnershipOrder(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	devA, err := seedActiveDevice(ctx, deviceRepo, "device-A")
	require.NoError(t, err)
	devB, err := seedActiveDevice(ctx, deviceRepo, "device-B")
	require.NoError(t, err)

	_, err = svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-A",
		DeviceID:     devA.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "A_",
		Enabled:      true,
	})
	require.NoError(t, err)

	_, err = svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-B1",
		DeviceID:     devB.ID,
		StartAddress: "40011",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "B1_",
		Enabled:      true,
	})
	require.NoError(t, err)

	_, err = svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-B2",
		DeviceID:     devB.ID,
		StartAddress: "40021",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "B2_",
		Enabled:      true,
	})
	require.NoError(t, err)

	rules, err := svc.ListByDeviceIDs(ctx, []string{devB.ID, devA.ID})
	require.NoError(t, err)
	require.Len(t, rules, 3)
	require.Equal(t, []string{"rule-B1", "rule-B2", "rule-A"}, []string{rules[0].ID, rules[1].ID, rules[2].ID})
	require.Equal(t, []string{devB.ID, devB.ID, devA.ID}, []string{rules[0].DeviceID, rules[1].DeviceID, rules[2].DeviceID})
}
