package runtime

import (
	"errors"
	"testing"

	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionErrorFinalizationIsIndependentOfObserverAndPreservesOwnedRestore(t *testing.T) {
	for _, restore := range []bool{false, true} {
		t.Run(map[bool]string{false: "failed-forward", true: "restored-old"}[restore], func(t *testing.T) {
			live, _, projection := newAcquisitionInputSQLFixture(t)
			require.NoError(t, live.ApplyWorkspaceProjection(t.Context(), projection))
			live.sourceRuleObserver = nil
			live.running.Store(true)
			snapshot := live.SnapshotSourceRuleProjection("d", nil)
			scope := sourcerule.RuntimeReconcileScope{DeviceID: "d", RuleID: "zero-point-scope"}
			ctx, finish := live.BeginSourceRuleAcquisitionMutation(t.Context(), scope)
			out := live.ReconcileSourceRule(ctx, sourcerule.RuntimeReconcileRequest{Scope: scope, Operation: sourcerule.RuntimeReconcileOperationCreate, ProjectionSnapshot: snapshot})
			require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, out.Status)
			operation := ctx.Value(sourceAcquisitionMutationKey{}).(*sourceAcquisitionMutation)
			require.True(t, operation.forward.Published())
			if restore {
				require.NoError(t, live.RestoreSourceRuleProjectionContext(ctx, snapshot, nil))
			}
			live.FinalizeSourceRuleObservation(ctx, out, errors.New("source authority refused"))
			finish()
			next := live.scheduler.BeginAcquisitionGeometryMutation("d")
			require.Equal(t, restore, live.scheduler.CanRestoreAcquisitionGeometry(next), "only a separately accepted old geometry survives final error cleanup")
			live.scheduler.FinishAcquisitionGeometryMutation(next)
		})
	}
}
