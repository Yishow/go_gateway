package runtime

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
)

func TestSourceProjectionForeignMappingBeforeForwardMustSurvive(t *testing.T) {
	for _, operation := range []string{"disable", "delete", "grow", "disable-same-values", "delete-same-values", "grow-same-values"} {
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
			armed := false
			var targetID string
			foreignTag := "foreign-accepted"
			if strings.HasSuffix(operation, "-same-values") {
				foreignTag = "old-tag"
			}
			foreign := []mappingBinding{{TagID: foreignTag}}
			live.workspace = compensationWorkspaceFunc(func(ctx context.Context) (*workspace.RuntimeProjection, error) {
				if armed {
					armed = false
					live.mappingMu.Lock()
					live.mappingIndex[targetID] = foreign
					live.mappingMu.Unlock()
				}
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
			targetID = links[0].PointID
			armed = true
			var outcome sourcerule.RuntimeReconcileOutcome
			switch strings.TrimSuffix(operation, "-same-values") {
			case "disable":
				outcome, err = rules.DisableWithRuntimeReconcile(ctx, main.ID)
			case "delete":
				outcome, err = rules.DeleteWithRuntimeReconcile(ctx, main.ID)
			case "grow":
				count := 3
				_, outcome, err = rules.UpdateWithRuntimeReconcile(ctx, main.ID, sourcerule.UpdateRuleRequest{Count: &count})
			}
			require.Error(t, err)
			require.Equal(t, sourcerule.RuntimeReconcileStatusStale, outcome.Status)
			want := map[string]struct{}{links[0].PointID: {}, links[1].PointID: {}, "unrelated": {}}
			require.Equal(t, want, live.deviceProjectionState(d.ID).AppliedPointIDs, "compensation restores exact applied keys and removes newly added keys")
			require.Same(t, &foreign[0], &live.mappingIndex[targetID][0], "neither forward publication nor compensation may replace the foreign installation")
			require.Equal(t, []mappingBinding{{TagID: "unrelated-tag"}}, live.mappingBindingsForPoint("unrelated"))
		})
	}
}

func TestSourceProjectionForwardMappingCASIsAtomicForSelectedKeys(t *testing.T) {
	for _, sameValues := range []bool{false, true} {
		t.Run(fmt.Sprint(sameValues), func(t *testing.T) {
			live := &Service{scheduler: collector.NewScheduler(collector.DefaultSchedulerConfig(), nil), mappingIndex: map[string][]mappingBinding{"p": {{TagID: "old-p"}}, "q": {{TagID: "old-q"}}}}
			snapshot := live.SnapshotSourceRuleProjection("d", []string{"p", "q"})
			ctx := context.WithValue(t.Context(), sourceProjectionRollbackKey{}, snapshot)
			foreignTag := "foreign"
			if sameValues {
				foreignTag = "old-q"
			}
			foreign := []mappingBinding{{TagID: foreignTag}}
			live.mappingIndex["q"] = foreign
			err := live.applyWorkspaceProjectionMappingsForPoints(ctx, &workspace.RuntimeProjection{}, map[string]struct{}{"p": {}, "q": {}})
			require.Error(t, err)
			require.Equal(t, []mappingBinding{{TagID: "old-p"}}, live.mappingBindingsForPoint("p"), "a conflict on q must not publish p first")
			require.Same(t, &foreign[0], &live.mappingIndex["q"][0], "foreign install identity survives")
			proof := &sourcerule.RuntimeProjectionProof{Points: map[string]sourcerule.RuntimePointProof{"p": {}, "q": {}}}
			require.Error(t, live.RestoreSourceRuleProjection(snapshot, proof))
			require.Same(t, &foreign[0], &live.mappingIndex["q"][0])
		})
	}
}

func TestSourceProjectionForwardRejectsForeignAfterOwnPublication(t *testing.T) {
	live := &Service{scheduler: collector.NewScheduler(collector.DefaultSchedulerConfig(), nil), mappingIndex: map[string][]mappingBinding{"p": {{TagID: "before"}}}}
	snapshot := live.SnapshotSourceRuleProjection("d", []string{"p"})
	ctx := context.WithValue(t.Context(), sourceProjectionRollbackKey{}, snapshot)
	require.NoError(t, live.applyWorkspaceProjectionMappingsForPoints(ctx, &workspace.RuntimeProjection{}, map[string]struct{}{"p": {}}))
	foreign := []mappingBinding{{TagID: "later"}}
	live.mappingIndex["p"] = foreign
	require.Error(t, live.applyWorkspaceProjectionMappingsForPoints(ctx, &workspace.RuntimeProjection{}, map[string]struct{}{"p": {}}))
	require.Same(t, &foreign[0], &live.mappingIndex["p"][0])
	require.Error(t, live.RestoreSourceRuleProjection(snapshot, &sourcerule.RuntimeProjectionProof{Points: map[string]sourcerule.RuntimePointProof{"p": {}}}))
	require.Same(t, &foreign[0], &live.mappingIndex["p"][0])
}

func TestSourceProjectionEmptyScopePreservesForeignMappingAndLease(t *testing.T) {
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	live := &Service{scheduler: scheduler, mappingIndex: map[string][]mappingBinding{"p": {{TagID: "before"}}}}
	snapshot := live.SnapshotSourceRuleProjection("d", []string{"p"})
	foreign := []mappingBinding{{TagID: "foreign"}}
	live.mappingIndex["p"] = foreign
	p := &schema.Point{ID: "p", DeviceID: "d", RevisionID: "foreign-r", DataType: schema.DataTypeInt16, Enabled: true}
	lease := scheduler.AddPointWithObservationLease(p)
	ctx := context.WithValue(t.Context(), sourceProjectionRollbackKey{}, snapshot)
	require.NoError(t, live.applyAffectedSourceRuleProjection(ctx, &workspace.RuntimeProjection{}, sourcerule.RuntimeReconcileRequest{Scope: sourcerule.RuntimeReconcileScope{RuleID: "empty", DeviceID: "d"}}))
	require.Same(t, &foreign[0], &live.mappingIndex["p"][0])
	require.Equal(t, lease, scheduler.CapturePointLease(p.ID, p.RevisionID))
}

func TestSourceProjectionPartialRemovalPreservesLaterForeignTarget(t *testing.T) {
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	live := &Service{scheduler: scheduler, mappingIndex: map[string][]mappingBinding{"p": {{TagID: "p-before"}}, "q": {{TagID: "q-before"}}}, pointMetaIndex: map[string]pointMeta{}, projections: map[string]runtimeProjectionDeviceState{"d": {AppliedPointIDs: map[string]struct{}{"p": {}, "q": {}}}}}
	leases := map[string]uint64{}
	for _, id := range []string{"p", "q"} {
		leases[id] = scheduler.AddPointWithObservationLease(&schema.Point{ID: id, DeviceID: "d", RevisionID: "r", Enabled: true, DataType: schema.DataTypeInt16})
	}
	snapshot := live.SnapshotSourceRuleProjection("d", []string{"p", "q"})
	proof := &sourcerule.RuntimeProjectionProof{Epoch: leases["q"], Points: map[string]sourcerule.RuntimePointProof{"p": {Revision: "r", Lease: leases["p"], Deleted: true}, "q": {Revision: "r", Lease: leases["q"], Deleted: true}}}
	projection := &workspace.RuntimeProjection{Version: "after", DeviceIDs: []string{"d"}, Devices: []*schema.Device{{ID: "d", Status: schema.DeviceStatusActive}}}
	live.workspace = workspaceProjectionStub{projection: projection}
	live.running.Store(true)
	var foreignID, removedID string
	var receipt uint64
	foreign := []mappingBinding{{TagID: "foreign-after-first-target"}}
	out := live.ReconcileSourceRule(t.Context(), sourcerule.RuntimeReconcileRequest{Operation: sourcerule.RuntimeReconcileOperationDelete, Scope: sourcerule.RuntimeReconcileScope{RuleID: "R", DeviceID: "d"}, AffectedPointIDs: []string{"p", "q"}, ProjectionProof: proof, ProjectionSnapshot: snapshot, PointRemoved: func(id string, epoch uint64) {
		// This synchronous receipt hook runs under the publisher's mapping lock.
		// Replace the other target between the first and second scoped removals.
		require.Empty(t, removedID)
		removedID = id
		receipt = epoch
		foreignID = "p"
		if id == "p" {
			foreignID = "q"
		}
		live.mappingIndex[foreignID] = foreign
	}})
	require.Equal(t, sourcerule.RuntimeReconcileStatusStale, out.Status)
	require.NotEmpty(t, removedID)
	require.Same(t, &foreign[0], &live.mappingIndex[foreignID][0])
	require.Equal(t, leases[foreignID], scheduler.CapturePointLease(foreignID, "r"), "second target lease must survive the mapping CAS refusal")
	require.Zero(t, scheduler.CapturePointLease(removedID, "r"))
	proof.Epoch = receipt
	proof.Points[removedID] = sourcerule.RuntimePointProof{Revision: "r", Deleted: true}
	require.Error(t, live.RestoreSourceRuleProjection(snapshot, proof))
	require.Same(t, &foreign[0], &live.mappingIndex[foreignID][0])
	require.Equal(t, leases[foreignID], scheduler.CapturePointLease(foreignID, "r"))
}

func TestSourceProjectionForwardRejectsForeignSnapshotOwner(t *testing.T) {
	live := &Service{mappingIndex: map[string][]mappingBinding{"p": {{TagID: "untouched"}}}}
	other := &Service{}
	for _, value := range []any{other.SnapshotSourceRuleProjection("d", []string{"p"}), struct{}{}} {
		ctx := context.WithValue(t.Context(), sourceProjectionRollbackKey{}, value)
		require.Error(t, live.applyWorkspaceProjectionMappingsForPoints(ctx, &workspace.RuntimeProjection{}, map[string]struct{}{"p": {}}))
		require.Equal(t, []mappingBinding{{TagID: "untouched"}}, live.mappingBindingsForPoint("p"))
	}
}
