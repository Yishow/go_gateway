package device

import (
	"context"
	"errors"
	"testing"
	"time"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type stubAcquisitionSyncer struct {
	groups []*schema.PollingGroup
}

func (s *stubAcquisitionSyncer) AddPollingGroup(group *schema.PollingGroup) {
	if group != nil {
		s.groups = append(s.groups, group)
	}
}

func TestService_Update_AcquisitionIntervalReassignsOnlyTargetDevicePointsAndUpdatesScheduler(t *testing.T) {
	ctx := context.Background()
	devRepo := NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	groupRepo := pollinggroup.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, groupRepo)
	syncer := &stubAcquisitionSyncer{}
	pointSvc.SetPollingGroupSyncer(syncer)

	svc := NewService(devRepo, nil)
	svc.SetAcquisitionIntervalCoordinator(pointSvc)

	now := time.Now()
	for _, id := range []string{"dev-a", "dev-b"} {
		require.NoError(t, devRepo.Create(ctx, &schema.Device{
			ID:               id,
			Name:             id,
			Protocol:         schema.ProtocolModbusTCP,
			Status:           schema.DeviceStatusActive,
			ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}`,
			CreatedAt:        now,
			UpdatedAt:        now,
		}))
	}

	groupAInit, err := pointSvc.SetDeviceAcquisitionInterval(ctx, "dev-a", 1000)
	require.NoError(t, err)
	groupBInit, err := pointSvc.SetDeviceAcquisitionInterval(ctx, "dev-b", 1000)
	require.NoError(t, err)

	ptA1, err := pointSvc.Create(ctx, point.CreatePointRequest{DeviceID: "dev-a", Name: "A1", Address: "40001", DataType: schema.DataTypeUint16})
	require.NoError(t, err)
	ptB1, err := pointSvc.Create(ctx, point.CreatePointRequest{DeviceID: "dev-b", Name: "B1", Address: "40001", DataType: schema.DataTypeUint16, PollingGroupID: &groupBInit.ID})
	require.NoError(t, err)

	updatedA, err := svc.Update(ctx, "dev-a", UpdateDeviceRequest{
		AcquisitionIntervalMs:   common.Ptr(500),
		AcquisitionAutoInterval: common.Ptr(true),
		StaleThresholdMs:        common.Ptr(2500),
		FrozenThresholdMs:       common.Ptr(10000),
	})
	require.NoError(t, err)
	require.Equal(t, common.Ptr(500), updatedA.AcquisitionIntervalMs)
	require.True(t, updatedA.AcquisitionAutoInterval)
	require.Equal(t, common.Ptr(2500), updatedA.StaleThresholdMs)
	require.Equal(t, 10000, updatedA.FrozenThresholdMs)

	gotA, err := svc.GetByID(ctx, "dev-a")
	require.NoError(t, err)
	require.Equal(t, common.Ptr(500), gotA.AcquisitionIntervalMs)
	require.True(t, gotA.AcquisitionAutoInterval)
	require.Equal(t, common.Ptr(2500), gotA.StaleThresholdMs)
	require.Equal(t, 10000, gotA.FrozenThresholdMs)

	reloadedPtA1, err := pointSvc.GetByID(ctx, ptA1.ID)
	require.NoError(t, err)
	require.Equal(t, common.Ptr(groupAInit.ID), reloadedPtA1.PollingGroupID)

	groupAUpdated, err := pointSvc.GetPollingGroupByOwnerDeviceID(ctx, "dev-a")
	require.NoError(t, err)
	require.Equal(t, 500, groupAUpdated.IntervalMs)
	require.Equal(t, 500, syncer.groups[len(syncer.groups)-1].IntervalMs)

	reloadedPtB1, err := pointSvc.GetByID(ctx, ptB1.ID)
	require.NoError(t, err)
	require.Equal(t, common.Ptr(groupBInit.ID), reloadedPtB1.PollingGroupID)
	groupB, err := pointSvc.GetPollingGroupByOwnerDeviceID(ctx, "dev-b")
	require.NoError(t, err)
	require.Equal(t, 1000, groupB.IntervalMs)
}

func TestService_Update_RejectsInvalidAcquisitionIntervalAndThresholds(t *testing.T) {
	ctx := context.Background()
	devRepo := NewMemoryRepository()
	pointSvc := point.NewService(point.NewMemoryRepository(), pollinggroup.NewMemoryRepository())
	svc := NewService(devRepo, nil)
	svc.SetAcquisitionIntervalCoordinator(pointSvc)

	now := time.Now()
	require.NoError(t, devRepo.Create(ctx, &schema.Device{
		ID:               "dev-a",
		Name:             "dev-a",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}))

	cases := []struct {
		name string
		req  UpdateDeviceRequest
		want string
	}{
		{name: "interval 99", req: UpdateDeviceRequest{AcquisitionIntervalMs: common.Ptr(99)}, want: "acquisition_interval_ms"},
		{name: "interval 3600001", req: UpdateDeviceRequest{AcquisitionIntervalMs: common.Ptr(3600001)}, want: "acquisition_interval_ms"},
		{name: "stale 0", req: UpdateDeviceRequest{StaleThresholdMs: common.Ptr(0)}, want: "stale_threshold_ms"},
		{name: "stale -1", req: UpdateDeviceRequest{StaleThresholdMs: common.Ptr(-1)}, want: "stale_threshold_ms"},
		{name: "frozen -1", req: UpdateDeviceRequest{FrozenThresholdMs: common.Ptr(-1)}, want: "frozen_threshold_ms"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Update(ctx, "dev-a", tc.req)
			require.Error(t, err)
			require.True(t, errors.Is(err, ErrValidation), "expected ErrValidation, got %v", err)
			require.ErrorContains(t, err, tc.want)
		})
	}

	group, err := pointSvc.GetPollingGroupByOwnerDeviceID(ctx, "dev-a")
	require.NoError(t, err)
	require.Nil(t, group)
}

type countingAcquisitionCoordinator struct {
	groups      []*schema.PollingGroup
	listCalls   int
	lookupCalls int
}

func (c *countingAcquisitionCoordinator) SetDeviceAcquisitionInterval(context.Context, string, int) (*schema.PollingGroup, error) {
	return nil, nil
}

func (c *countingAcquisitionCoordinator) GetPollingGroupByOwnerDeviceID(context.Context, string) (*schema.PollingGroup, error) {
	c.lookupCalls++
	return nil, nil
}

func (c *countingAcquisitionCoordinator) ListPollingGroups(context.Context) ([]*schema.PollingGroup, error) {
	c.listCalls++
	return c.groups, nil
}

func TestServiceList_HydratesIntervalsWithOneGroupRead(t *testing.T) {
	ctx := t.Context()
	repo := NewMemoryRepository()
	for _, id := range []string{"a", "b", "c"} {
		require.NoError(t, repo.Create(ctx, &schema.Device{ID: id, Name: id, Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{}`}))
	}
	owner := func(id string) *string { return &id }
	coord := &countingAcquisitionCoordinator{groups: []*schema.PollingGroup{
		{ID: "g-a", IntervalMs: 500, OwnerDeviceID: owner("a")},
		{ID: "g-b", IntervalMs: 2000, OwnerDeviceID: owner("b")},
		{ID: "shared", IntervalMs: 1000},
	}}
	svc := NewService(repo, nil)
	svc.SetAcquisitionIntervalCoordinator(coord)

	devices, err := svc.List(ctx, ListFilter{})
	require.NoError(t, err)
	intervals := map[string]*int{}
	for _, dev := range devices {
		intervals[dev.ID] = dev.AcquisitionIntervalMs
	}
	require.Equal(t, 500, *intervals["a"])
	require.Equal(t, 2000, *intervals["b"])
	require.Nil(t, intervals["c"])
	require.Equal(t, 1, coord.listCalls)
	require.Zero(t, coord.lookupCalls)
}
