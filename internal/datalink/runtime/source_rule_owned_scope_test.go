package runtime

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
	"testing"
)

func TestSourceRuleOwnedMutationKeepsRecordingPointScope(t *testing.T) {
	ctx := t.Context()
	dr := device.NewMemoryRepository()
	ready := true
	d := &schema.Device{ID: "D", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"host":"127.0.0.1","data_format":"ABCD"}`, LastTestSuccess: &ready}
	require.NoError(t, dr.Create(ctx, d))
	ps := point.NewService(point.NewMemoryRepository(), nil)
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	live := &Service{scheduler: scheduler, pointSvc: ps, mappingIndex: map[string][]mappingBinding{}, pointMetaIndex: map[string]pointMeta{}}
	rules := sourcerule.NewService(sourcerule.NewMemoryRepository(), device.NewService(dr, nil), ps, live)
	live.workspace = compensationWorkspaceFunc(func(ctx context.Context) (*workspace.RuntimeProjection, error) {
		return persistedRuleProjection(ctx, rules, ps, d)
	})
	require.NoError(t, live.UpsertDevice(ctx, d))
	a, err := rules.Create(ctx, sourcerule.CreateRuleRequest{ID: "A", DeviceID: d.ID, StartAddress: "40001", Count: 8, DataType: schema.DataTypeInt16, NamingPrefix: "A", Enabled: true})
	require.NoError(t, err)
	links, err := rules.ListLinks(ctx, a.ID)
	require.NoError(t, err)
	for _, link := range links[1:] {
		live.RemovePoint(link.PointID)
	}
	projection, err := live.workspace.RuntimeProjection(ctx)
	require.NoError(t, err)
	require.NoError(t, live.ApplyWorkspaceProjectionForDevicesAndPoints(ctx, projection, []string{d.ID}, []string{links[0].PointID}))
	live.running.Store(true)
	check := func(count int) {
		t.Helper()
		require.Len(t, scheduler.CollectionStatus(), count)
		require.Len(t, live.deviceProjectionState(d.ID).AppliedPointIDs, count)
		require.Zero(t, live.CapturePointLease(links[1].PointID, ""))
		require.NotContains(t, scheduler.CollectionStatus(), links[1].PointID)
	}
	check(1)
	b, outcome, err := rules.CreateWithRuntimeReconcile(ctx, sourcerule.CreateRuleRequest{ID: "B", DeviceID: d.ID, StartAddress: "40009", Count: 8, DataType: schema.DataTypeInt16, NamingPrefix: "B", Enabled: true})
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	check(9)
	outcome, err = rules.DisableWithRuntimeReconcile(ctx, b.ID)
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	check(1)
	outcome, err = rules.EnableWithRuntimeReconcile(ctx, b.ID)
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	check(9)
	count := 2
	_, outcome, err = rules.UpdateWithRuntimeReconcile(ctx, b.ID, sourcerule.UpdateRuleRequest{Count: &count})
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	check(3)
	outcome, err = rules.DeleteWithRuntimeReconcile(ctx, b.ID)
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	check(1)
}
