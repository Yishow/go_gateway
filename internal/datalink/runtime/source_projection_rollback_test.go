package runtime

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
)

func TestSourceRuleShareFailureRestoresAffectedProjection(t *testing.T) {
	for _, operation := range []string{"disable", "delete", "grow", "create"} {
		t.Run(operation, func(t *testing.T) {
			ctx := t.Context()
			dr := device.NewMemoryRepository()
			ready := true
			d := &schema.Device{ID: "d", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"host":"127.0.0.1","port":502}`, LastTestSuccess: &ready}
			require.NoError(t, dr.Create(ctx, d))
			ps := point.NewService(point.NewMemoryRepository(), nil)
			scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
			live := &Service{scheduler: scheduler, pointSvc: ps, mappingIndex: map[string][]mappingBinding{}, pointMetaIndex: map[string]pointMeta{}}
			rules := sourcerule.NewService(sourcerule.NewMemoryRepository(), device.NewService(dr, nil), ps, live)
			scheduler.AddDevice(d)
			live.workspace = compensationWorkspaceFunc(func(ctx context.Context) (*workspace.RuntimeProjection, error) {
				return persistedRuleProjection(ctx, rules, ps, d)
			})
			main, err := rules.Create(ctx, sourcerule.CreateRuleRequest{ID: "main", DeviceID: d.ID, StartAddress: "40001", Count: 2, DataType: schema.DataTypeInt16, Enabled: true})
			require.NoError(t, err)
			links, err := rules.ListLinks(ctx, main.ID)
			require.NoError(t, err)
			projection, err := live.workspace.RuntimeProjection(ctx)
			require.NoError(t, err)
			require.NoError(t, live.ApplyWorkspaceProjection(ctx, projection))
			live.mappingIndex[links[0].PointID] = []mappingBinding{{TagID: "old-tag"}}
			live.mappingIndex["unrelated"] = []mappingBinding{{TagID: "unrelated-tag"}}
			live.projections[d.ID].AppliedPointIDs["unrelated"] = struct{}{}
			live.running.Store(true)
			rules.SetShareRuntimeReconciler(&compensationShareFailure{})
			switch operation {
			case "create":
				_, _, err = rules.CreateWithRuntimeReconcile(ctx, sourcerule.CreateRuleRequest{ID: "new", DeviceID: d.ID, StartAddress: "40009", Count: 2, DataType: schema.DataTypeInt16, Enabled: true})
			case "disable":
				_, err = rules.DisableWithRuntimeReconcile(ctx, main.ID)
			case "delete":
				_, err = rules.DeleteWithRuntimeReconcile(ctx, main.ID)
			case "grow":
				count := 3
				_, _, err = rules.UpdateWithRuntimeReconcile(ctx, main.ID, sourcerule.UpdateRuleRequest{Count: &count})
			}
			require.Error(t, err)
			want := map[string]struct{}{links[0].PointID: {}, links[1].PointID: {}, "unrelated": {}}
			require.Equal(t, want, live.deviceProjectionState(d.ID).AppliedPointIDs, "compensation restores exact applied keys and removes newly added keys")
			require.Equal(t, []mappingBinding{{TagID: "old-tag"}}, live.mappingBindingsForPoint(links[0].PointID), "compensation restores runtime mappings")
			require.Equal(t, []mappingBinding{{TagID: "unrelated-tag"}}, live.mappingBindingsForPoint("unrelated"))
		})
	}
}

func TestSourceProjectionRollbackRejectsForeignLeaseAndAbsentABA(t *testing.T) {
	for _, mode := range []string{"lease", "absent-aba", "mapping-only", "mapping-same-values"} {
		t.Run(mode, func(t *testing.T) {
			scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
			live := &Service{scheduler: scheduler, mappingIndex: map[string][]mappingBinding{"p": {{TagID: "old"}}}, pointMetaIndex: map[string]pointMeta{}, projections: map[string]runtimeProjectionDeviceState{"d": {AppliedPointIDs: map[string]struct{}{"p": {}}}}}
			point := &schema.Point{ID: "p", DeviceID: "d", Address: "40001", DataType: schema.DataTypeInt16, Enabled: true}
			if mode == "lease" {
				scheduler.AddPoint(point)
			}
			snapshot := live.SnapshotSourceRuleProjection("d", []string{"p"})
			current, captured := scheduler.CaptureOrdinaryPoint("p")
			require.True(t, captured)
			proof := &sourcerule.RuntimeProjectionProof{Epoch: current.Epoch, Points: map[string]sourcerule.RuntimePointProof{"p": {Lease: current.Lease}}}
			if mode != "mapping-only" && mode != "mapping-same-values" {
				scheduler.AddPoint(point)
			}
			if mode == "absent-aba" {
				scheduler.RemovePoint("p")
			}
			foreignTag := "foreign"
			if mode == "mapping-same-values" {
				foreignTag = "old"
			}
			live.mappingIndex["p"] = []mappingBinding{{TagID: foreignTag}}
			delete(live.projections["d"].AppliedPointIDs, "p")
			require.Error(t, live.RestoreSourceRuleProjection(snapshot, proof))
			require.Equal(t, []mappingBinding{{TagID: foreignTag}}, live.mappingBindingsForPoint("p"))
			require.NotContains(t, live.deviceProjectionState("d").AppliedPointIDs, "p")
		})
	}
}
