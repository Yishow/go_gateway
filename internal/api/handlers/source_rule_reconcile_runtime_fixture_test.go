package handlers

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	runtimepkg "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	"testing"
)

func wireSourceRuleReconcileRuntimeFixture(t *testing.T, recorder *sourceRuleRuntimeReconcileRecorder, devices *device.Service, points *point.Service, rules *sourcerule.Service) {
	t.Helper()
	live, err := runtimepkg.NewService(runtimepkg.DefaultConfig(), runtimepkg.Dependencies{
		Scheduler: collector.NewScheduler(collector.DefaultSchedulerConfig(), nil), Writer: storage.NewMemoryStorage(10),
		DeviceService: devices, PointService: points, MappingService: mapping.NewService(mapping.NewMemoryRepository()),
		TagService: tag.NewService(tag.NewMemoryRepository()), PollingGroupService: pollinggroup.NewService(pollinggroup.NewMemoryRepository()),
	})
	require.NoError(t, err)
	recorder.Service = live
	d, err := devices.GetByID(t.Context(), "dev-A")
	require.NoError(t, err)
	require.NoError(t, live.ApplyWorkspaceProjection(t.Context(), &workspace.RuntimeProjection{WorkspaceID: "fixture", DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}}))
	live.SetWorkspaceProjectionReader(pointProfileProjectionFunc(func(ctx context.Context) (*workspace.RuntimeProjection, error) {
		d, err := devices.GetByID(ctx, "dev-A")
		if err != nil {
			return nil, err
		}
		p := &workspace.RuntimeProjection{WorkspaceID: "fixture", DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}}
		p.Points, err = points.ListByDevice(ctx, d.ID)
		if err != nil {
			return nil, err
		}
		p.Rules, err = rules.ListByDeviceIDs(ctx, []string{d.ID})
		if err != nil {
			return nil, err
		}
		for _, rule := range p.Rules {
			links, e := rules.ListLinks(ctx, rule.ID)
			if e != nil {
				return nil, e
			}
			p.RuleLinks = append(p.RuleLinks, links...)
		}
		return p, nil
	}))
}
