package sourcerule

import (
	"context"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

type finalSourceRecorder struct {
	*reconcileRuntimeSync
	final            func(RuntimeReconcileOutcome, error)
	calls            int
	failOwnedRestore bool
}

func (r *finalSourceRecorder) RestoreRemovedPointOwned(p *schema.Point, removedLease, epoch uint64) (uint64, bool) {
	if r.failOwnedRestore {
		return 0, false
	}
	return r.reconcileRuntimeSync.RestoreRemovedPointOwned(p, removedLease, epoch)
}

func (r *finalSourceRecorder) FinalizeSourceRuleObservation(_ context.Context, o RuntimeReconcileOutcome, err error) {
	r.calls++
	if r.final != nil {
		r.final(o, err)
	}
}

func TestSourceFinalizationWaitsForShareAndCompensation(t *testing.T) {
	f := newRuntimeReconcileRuleFixture(t, true)
	r := &finalSourceRecorder{reconcileRuntimeSync: f.runtimeSync}
	f.svc.runtimeSync = r
	share := false
	f.svc.SetShareRuntimeReconciler(ShareRuntimeReconcilerFunc(func(context.Context, RuntimeReconcileRequest) (RuntimeReconcileOutcome, error) {
		share = true
		require.Zero(t, r.calls)
		return RuntimeReconcileOutcome{}, errShareReconcile
	}))
	r.final = func(_ RuntimeReconcileOutcome, err error) {
		require.True(t, share)
		requireShareReconcileFailure(t, err)
		restored, getErr := f.svc.GetByID(f.ctx, f.rule.ID)
		require.NoError(t, getErr)
		require.Equal(t, f.rule.NamingPrefix, restored.NamingPrefix)
	}
	prefix := "failed-change"
	_, _, err := f.svc.UpdateWithRuntimeReconcile(f.ctx, f.rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})
	requireShareReconcileFailure(t, err)
	require.Equal(t, 1, r.calls)
}

func TestSourceFinalizationPanicDoesNotChangeSuccessfulOutcome(t *testing.T) {
	f := newRuntimeReconcileRuleFixture(t, true)
	r := &finalSourceRecorder{reconcileRuntimeSync: f.runtimeSync, final: func(RuntimeReconcileOutcome, error) { panic("password_fixture") }}
	f.svc.runtimeSync = r
	outcome, err := f.svc.DisableWithRuntimeReconcile(f.ctx, f.rule.ID)
	require.NoError(t, err)
	require.Equal(t, RuntimeReconcileStatusAligned, outcome.Status)
	require.Equal(t, 1, r.calls)
	rule, err := f.svc.GetByID(f.ctx, f.rule.ID)
	require.NoError(t, err)
	require.False(t, rule.Enabled)
}

func TestSourceFinalizationDirtyUnknownIsPreserved(t *testing.T) {
	f := newRuntimeReconcileRuleFixture(t, true)
	f.runtimeSync.reconcileOutcomes = []RuntimeReconcileStatus{RuntimeReconcileStatusStale, RuntimeReconcileStatusFailed}
	r := &finalSourceRecorder{reconcileRuntimeSync: f.runtimeSync, failOwnedRestore: true}
	f.svc.runtimeSync = r
	r.final = func(_ RuntimeReconcileOutcome, err error) {
		var typed *modbusshare.Error
		require.ErrorAs(t, err, &typed)
		require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
	}
	_, err := f.svc.DisableWithRuntimeReconcile(f.ctx, f.rule.ID)
	require.Error(t, err)
	require.Equal(t, 1, r.calls)
}

func TestSourceFinalizationCarriesGenericProofThroughShareOutcome(t *testing.T) {
	f := newRuntimeReconcileRuleFixture(t, true)
	proof := &RuntimeSourceObservation{RuleID: f.rule.ID, RuleRevision: f.rule.RevisionID, ReferenceID: "actual-attempt", Verified: true, Attempted: true, Applied: true}
	f.runtimeSync.reconcileOutcome = RuntimeReconcileOutcome{Status: RuntimeReconcileStatusAligned, Observation: proof}
	r := &finalSourceRecorder{reconcileRuntimeSync: f.runtimeSync}
	f.svc.runtimeSync = r
	f.svc.SetShareRuntimeReconciler(&shareReconcileRecorder{})
	r.final = func(out RuntimeReconcileOutcome, err error) {
		require.NoError(t, err)
		require.Same(t, proof, out.Observation)
	}
	_, err := f.svc.DeleteWithRuntimeReconcile(f.ctx, f.rule.ID)
	require.NoError(t, err)
	require.Equal(t, 1, r.calls)
	require.Equal(t, f.rule.RevisionID, r.reconcileRequests[0].RuleRevision)
}

func TestSourceFinalizationDoesNotRunBeforeRuntimeAttempt(t *testing.T) {
	f := newRuntimeReconcileRuleFixture(t, true)
	r := &finalSourceRecorder{reconcileRuntimeSync: f.runtimeSync}
	f.svc.runtimeSync = r
	_, err := f.svc.DeleteWithRuntimeReconcile(f.ctx, "missing")
	require.Error(t, err)
	require.Zero(t, r.calls)
}
