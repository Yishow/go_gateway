package modbusshare

import (
	"context"
	"fmt"
)

// Reconcile performs atomic, serialized, and revision-checked reconciliation of desired mappings.
func (r *Reconciler) Reconcile(ctx context.Context, req ReconcileRequest) (ReconcileOutcome, error) {
	if req.WorkspaceID == "" {
		return ReconcileOutcome{Outcome: HydrationStateFailed}, NewError(ErrCodeWorkspaceScope, "workspace id is required", false)
	}
	durable, inputErr := r.validateDurableInputs(req)
	if inputErr != nil {
		return ReconcileOutcome{Outcome: HydrationStateFailed}, inputErr
	}
	lock := r.getWorkspaceLock(req.WorkspaceID)
	lock.Lock()
	defer lock.Unlock()

	state, outcome, done, err := r.prepareReconcile(ctx, req, durable)
	if done {
		return outcome, err
	}
	r.projectionMu.Lock()
	defer r.projectionMu.Unlock()

	transaction, outcome, err := r.stageTransaction(ctx, req, state)
	if err != nil {
		return outcome, err
	}
	if outcome, replay := r.replayOutcome(state.key, state.currentRev, state.settings.SettingsRevision); replay {
		return outcome, nil
	}
	if req.ExpectedWorkspaceRevision != "" && state.currentRev != "" && req.ExpectedWorkspaceRevision != state.currentRev {
		outcome := ReconcileOutcome{Outcome: HydrationStateFailed, NewWorkspaceRevision: state.currentRev}
		err := &Error{Code: ErrCodeRevisionConflict, Message: fmt.Sprintf("workspace revision conflict: expected %s, observed %s", req.ExpectedWorkspaceRevision, state.currentRev), Retryable: true, WorkspaceRevision: state.currentRev}
		return outcome, err
	}
	if state.hydration.ReadinessToken != "" && req.ReadinessToken != state.hydration.ReadinessToken && !req.Restore {
		outcome := ReconcileOutcome{Outcome: HydrationStateFailed, NewWorkspaceRevision: state.hydration.WorkspaceRevision, SettingsRevision: state.settings.SettingsRevision}
		err := &Error{Code: ErrCodeHydrationRequired, Message: "share readiness token is stale", Retryable: true, Action: "hydrate the workspace again before reconcile"}
		return outcome, err
	}
	if isMappingsEquivalent(transaction.currentMappings, transaction.staged) {
		outcome := ReconcileOutcome{Outcome: "aligned", NewWorkspaceRevision: state.currentRev, SettingsRevision: state.settings.SettingsRevision, AppliedCount: len(transaction.staged), Mappings: transaction.currentMappings}
		r.rememberOutcome(state.key, outcome, state.currentRev, state.settings.SettingsRevision)
		return outcome, nil
	}
	if outcome, err := r.applyTransaction(ctx, req, state, transaction); err != nil {
		return outcome, err
	}
	outcome = r.buildAppliedOutcome(state, transaction, req.WorkspaceID)
	r.rememberOutcome(state.key, outcome, transaction.newRev, state.settings.SettingsRevision)
	return outcome, nil
}
