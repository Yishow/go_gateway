package point

import (
	"testing"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestEnsureDefaultPollingGroupIgnoresDeviceOwnedGroups(t *testing.T) {
	groups := NewMemoryPollingGroupRepository()
	owner := "dev-a"
	// The owned group sorts first (lower priority and interval) and must still
	// never become another device's fallback cadence.
	require.NoError(t, groups.Create(t.Context(), &schema.PollingGroup{ID: "owned", Name: "device:dev-a", IntervalMs: 100, Priority: 1, Enabled: true, OwnerDeviceID: &owner}))
	require.NoError(t, groups.Create(t.Context(), &schema.PollingGroup{ID: "auto", Name: "Auto Polling", IntervalMs: 1000, Priority: 50, Enabled: true}))
	service := NewService(NewMemoryRepository(), groups)
	got, err := service.EnsureDefaultPollingGroup(t.Context())
	require.NoError(t, err)
	require.Equal(t, "auto", got.ID)
}

func TestEnsureDefaultPollingGroupCreatesSharedGroupWhenOnlyOwnedExist(t *testing.T) {
	groups := NewMemoryPollingGroupRepository()
	owner := "dev-a"
	require.NoError(t, groups.Create(t.Context(), &schema.PollingGroup{ID: "owned", Name: "device:dev-a", IntervalMs: 500, Priority: 1, Enabled: true, OwnerDeviceID: &owner}))
	service := NewService(NewMemoryRepository(), groups)
	got, err := service.EnsureDefaultPollingGroup(t.Context())
	require.NoError(t, err)
	require.NotEqual(t, "owned", got.ID)
	require.Nil(t, got.OwnerDeviceID)
	require.Equal(t, defaultPollingGroupName, got.Name)
	require.Equal(t, defaultPollingGroupIntervalMs, got.IntervalMs)
	owned, err := groups.GetByID(t.Context(), "owned")
	require.NoError(t, err)
	require.Equal(t, 500, owned.IntervalMs, "owned group untouched")
}

func TestEnsureDefaultPollingGroupDoesNotEnableDisabledOwnedGroup(t *testing.T) {
	groups := NewMemoryPollingGroupRepository()
	owner := "dev-a"
	require.NoError(t, groups.Create(t.Context(), &schema.PollingGroup{ID: "owned", Name: "device:dev-a", IntervalMs: 500, Enabled: false, OwnerDeviceID: &owner}))
	require.NoError(t, groups.Create(t.Context(), &schema.PollingGroup{ID: "auto", Name: "Auto Polling", IntervalMs: 1000, Priority: 50, Enabled: false}))
	service := NewService(NewMemoryRepository(), groups)
	got, err := service.EnsureDefaultPollingGroup(t.Context())
	require.NoError(t, err)
	require.Equal(t, "auto", got.ID)
	require.True(t, got.Enabled)
	owned, err := groups.GetByID(t.Context(), "owned")
	require.NoError(t, err)
	require.False(t, owned.Enabled)
}

func TestSetDeviceAcquisitionInterval_CreatesAndUpdatesOwnedGroupAndReassignsOnlyDevicePoints(t *testing.T) {
	pointRepo := NewMemoryRepository()
	groupRepo := NewMemoryPollingGroupRepository()
	svc := NewService(pointRepo, groupRepo)
	syncer := &stubPollingGroupSyncer{}
	svc.SetPollingGroupSyncer(syncer)
	ctx := t.Context()

	fast, err := svc.CreatePollingGroup(ctx, CreatePollingGroupRequest{Name: "fast", IntervalMs: 250})
	require.NoError(t, err)
	slow, err := svc.CreatePollingGroup(ctx, CreatePollingGroupRequest{Name: "slow", IntervalMs: 2000})
	require.NoError(t, err)
	groupB, err := svc.CreatePollingGroup(ctx, CreatePollingGroupRequest{Name: "device:dev-b", IntervalMs: 1000, OwnerDeviceID: common.Ptr("dev-b")})
	require.NoError(t, err)

	a1, err := svc.Create(ctx, CreatePointRequest{DeviceID: "dev-a", Name: "A1", Address: "40001", DataType: schema.DataTypeUint16, PollingGroupID: &fast.ID})
	require.NoError(t, err)
	a2, err := svc.Create(ctx, CreatePointRequest{DeviceID: "dev-a", Name: "A2", Address: "40002", DataType: schema.DataTypeUint16, PollingGroupID: &slow.ID})
	require.NoError(t, err)
	a3, err := svc.Create(ctx, CreatePointRequest{DeviceID: "dev-a", Name: "A3", Address: "40003", DataType: schema.DataTypeUint16})
	require.NoError(t, err)
	b1, err := svc.Create(ctx, CreatePointRequest{DeviceID: "dev-b", Name: "B1", Address: "40001", DataType: schema.DataTypeUint16, PollingGroupID: &groupB.ID})
	require.NoError(t, err)

	groupA, err := svc.SetDeviceAcquisitionInterval(ctx, "dev-a", 500)
	require.NoError(t, err)
	require.NotNil(t, groupA)
	require.Equal(t, "device:dev-a", groupA.Name)
	require.Equal(t, common.Ptr("dev-a"), groupA.OwnerDeviceID)
	require.Equal(t, 500, groupA.IntervalMs)
	require.True(t, groupA.Enabled)
	require.NotEmpty(t, syncer.groups)
	require.Equal(t, groupA.ID, syncer.groups[len(syncer.groups)-1].ID)
	require.Equal(t, 500, syncer.groups[len(syncer.groups)-1].IntervalMs)

	for _, id := range []string{a1.ID, a2.ID, a3.ID} {
		pt, getErr := svc.GetByID(ctx, id)
		require.NoError(t, getErr)
		require.Equal(t, common.Ptr(groupA.ID), pt.PollingGroupID)
	}
	ptB, err := svc.GetByID(ctx, b1.ID)
	require.NoError(t, err)
	require.Equal(t, common.Ptr(groupB.ID), ptB.PollingGroupID)

	updatedGroupA, err := svc.SetDeviceAcquisitionInterval(ctx, "dev-a", 750)
	require.NoError(t, err)
	require.Equal(t, groupA.ID, updatedGroupA.ID)
	require.Equal(t, 750, updatedGroupA.IntervalMs)
	require.Equal(t, 750, syncer.groups[len(syncer.groups)-1].IntervalMs)
}

func TestSetDeviceAcquisitionInterval_ResolvesNameCollisionWithSuffix(t *testing.T) {
	svc := NewService(NewMemoryRepository(), NewMemoryPollingGroupRepository())
	ctx := t.Context()

	_, err := svc.CreatePollingGroup(ctx, CreatePollingGroupRequest{Name: "device:dev-a", IntervalMs: 1000})
	require.NoError(t, err)
	_, err = svc.CreatePollingGroup(ctx, CreatePollingGroupRequest{Name: "device:dev-a#2", IntervalMs: 1000})
	require.NoError(t, err)

	group, err := svc.SetDeviceAcquisitionInterval(ctx, "dev-a", 500)
	require.NoError(t, err)
	require.Equal(t, "device:dev-a#3", group.Name)
	require.Equal(t, common.Ptr("dev-a"), group.OwnerDeviceID)
}

func TestSetDeviceAcquisitionInterval_ValidatesIntervalBounds(t *testing.T) {
	svc := NewService(NewMemoryRepository(), NewMemoryPollingGroupRepository())
	ctx := t.Context()

	for _, invalid := range []int{99, 3600001} {
		_, err := svc.SetDeviceAcquisitionInterval(ctx, "dev-a", invalid)
		require.Error(t, err)
	}
	group, err := svc.GetPollingGroupByOwnerDeviceID(ctx, "dev-a")
	require.NoError(t, err)
	require.Nil(t, group)
}

func TestEnsureDevicePollingGroup_PrefersOwnedGroupAndFallsBackToSharedAutoPolling(t *testing.T) {
	svc := NewService(NewMemoryRepository(), NewMemoryPollingGroupRepository())
	ctx := t.Context()

	ownedA, err := svc.SetDeviceAcquisitionInterval(ctx, "dev-a", 500)
	require.NoError(t, err)

	gotA, err := svc.EnsureDevicePollingGroup(ctx, "dev-a")
	require.NoError(t, err)
	require.Equal(t, ownedA.ID, gotA.ID)
	require.Equal(t, 500, gotA.IntervalMs)

	gotB, err := svc.EnsureDevicePollingGroup(ctx, "dev-b")
	require.NoError(t, err)
	require.NotNil(t, gotB)
	require.NotEqual(t, ownedA.ID, gotB.ID)
	require.Nil(t, gotB.OwnerDeviceID)
	require.Equal(t, defaultPollingGroupIntervalMs, gotB.IntervalMs)
}
