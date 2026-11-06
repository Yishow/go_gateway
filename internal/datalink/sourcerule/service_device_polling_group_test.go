package sourcerule

import (
	"testing"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestService_CreateUpdateEnable_PrefersDeviceOwnedPollingGroupAndFallsBackToDefault(t *testing.T) {
	ctx := t.Context()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	groupRepo := pollinggroup.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, groupRepo)
	deviceSvc := device.NewService(deviceRepo, nil)
	svc := NewService(NewMemoryRepository(), deviceSvc, pointSvc, nil)

	devA, err := seedActiveDevice(ctx, deviceRepo, "dev-a")
	require.NoError(t, err)
	devB, err := seedActiveDevice(ctx, deviceRepo, "dev-b")
	require.NoError(t, err)

	groupA, err := pointSvc.SetDeviceAcquisitionInterval(ctx, devA.ID, 500)
	require.NoError(t, err)

	ruleA, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-a",
		DeviceID:     devA.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "A",
		Enabled:      true,
	})
	require.NoError(t, err)

	pointsA, err := pointSvc.List(ctx, point.ListFilter{DeviceID: &devA.ID})
	require.NoError(t, err)
	require.Len(t, pointsA, 2)
	for _, pt := range pointsA {
		require.Equal(t, common.Ptr(groupA.ID), pt.PollingGroupID)
	}

	_, err = svc.Update(ctx, ruleA.ID, UpdateRuleRequest{
		Count: common.Ptr(3),
	})
	require.NoError(t, err)

	pointsA, err = pointSvc.List(ctx, point.ListFilter{DeviceID: &devA.ID})
	require.NoError(t, err)
	require.Len(t, pointsA, 3)
	for _, pt := range pointsA {
		require.Equal(t, common.Ptr(groupA.ID), pt.PollingGroupID)
	}

	require.NoError(t, pointSvc.RemoveFromGroup(ctx, pointsA[0].ID))
	require.NoError(t, svc.Disable(ctx, ruleA.ID))
	require.NoError(t, svc.Enable(ctx, ruleA.ID))

	reloadedPt0, err := pointSvc.GetByID(ctx, pointsA[0].ID)
	require.NoError(t, err)
	require.Equal(t, common.Ptr(groupA.ID), reloadedPt0.PollingGroupID)

	_, err = svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-b",
		DeviceID:     devB.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "B",
		Enabled:      true,
	})
	require.NoError(t, err)

	pointsB, err := pointSvc.List(ctx, point.ListFilter{DeviceID: &devB.ID})
	require.NoError(t, err)
	require.Len(t, pointsB, 1)
	require.NotNil(t, pointsB[0].PollingGroupID)
	require.NotEqual(t, groupA.ID, *pointsB[0].PollingGroupID)

	groupB, err := groupRepo.GetByID(ctx, *pointsB[0].PollingGroupID)
	require.NoError(t, err)
	require.Nil(t, groupB.OwnerDeviceID)
	require.Equal(t, 1000, groupB.IntervalMs)
}
