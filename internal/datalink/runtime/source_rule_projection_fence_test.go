package runtime

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
	"testing"
)

func TestSourceRuleProjectionUsesOwnedReceiptNotCurrentLease(t *testing.T) {
	for _, mode := range []string{"foreign-before-reconcile", "absent-aba-in-provider", "deleted-foreign-before-reconcile"} {
		t.Run(mode, func(t *testing.T) {
			scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
			d := &schema.Device{ID: "D", Status: schema.DeviceStatusActive}
			p := &schema.Point{ID: "P", DeviceID: d.ID, RevisionID: "r", Address: "40001", DataType: schema.DataTypeInt16, Enabled: true}
			old := scheduler.AddPointWithObservationLease(p)
			proof := &sourcerule.RuntimeProjectionProof{Epoch: old, Points: map[string]sourcerule.RuntimePointProof{p.ID: {Revision: p.RevisionID, Lease: old}}}
			projection := &workspace.RuntimeProjection{DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}, Points: []*schema.Point{p}, RuleLinks: []*schema.SourceRuleLink{{RuleID: "R", PointID: p.ID}}}
			reader := &compensationProjectionReader{projection: projection}
			var foreign uint64
			if mode == "absent-aba-in-provider" {
				require.True(t, scheduler.RemovePointAtLease(p.ID, old))
				proof.Points[p.ID] = sourcerule.RuntimePointProof{Revision: p.RevisionID}
				reader.hook = func() {
					foreign = scheduler.AddPointWithObservationLease(p)
					require.True(t, scheduler.RemovePointAtLease(p.ID, foreign))
				}
			} else {
				foreign = scheduler.AddPointWithObservationLease(p)
				if mode == "deleted-foreign-before-reconcile" {
					projection.Points = nil
					projection.RuleLinks = nil
					proof.Points[p.ID] = sourcerule.RuntimePointProof{Revision: p.RevisionID, Lease: old, Deleted: true}
				}
			}
			svc := &Service{scheduler: scheduler, workspace: reader, mappingIndex: map[string][]mappingBinding{}, pointMetaIndex: map[string]pointMeta{p.ID: {DeviceID: d.ID, Address: p.Address}}}
			svc.running.Store(true)
			outcome := svc.ReconcileSourceRule(t.Context(), sourcerule.RuntimeReconcileRequest{Operation: sourcerule.RuntimeReconcileOperationUpdate, Scope: sourcerule.RuntimeReconcileScope{RuleID: "R", DeviceID: d.ID}, ProjectionProof: proof})
			require.Equal(t, sourcerule.RuntimeReconcileStatusStale, outcome.Status)
			if mode == "absent-aba-in-provider" {
				require.Zero(t, scheduler.CapturePointLease(p.ID, p.RevisionID))
			} else {
				require.Equal(t, foreign, scheduler.CapturePointLease(p.ID, p.RevisionID))
			}
		})
	}
}

func TestSourceRuleEmptyScopeDoesNotProjectWholeDevice(t *testing.T) {
	f := newSourceRuleReconcileRuntime(t, false, false)
	defer func() { require.NoError(t, f.service.Stop(t.Context())) }()
	f.service.RemovePoint("point-B")
	for _, op := range []sourcerule.RuntimeReconcileOperation{sourcerule.RuntimeReconcileOperationCreate, sourcerule.RuntimeReconcileOperationUpdate, sourcerule.RuntimeReconcileOperationDelete, sourcerule.RuntimeReconcileOperationDisable} {
		out := f.service.ReconcileSourceRule(t.Context(), sourcerule.RuntimeReconcileRequest{Operation: op, Scope: sourcerule.RuntimeReconcileScope{RuleID: "empty", DeviceID: f.device.ID}})
		require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, out.Status)
		require.Len(t, f.service.scheduler.CollectionStatus(), 1)
	}
}
