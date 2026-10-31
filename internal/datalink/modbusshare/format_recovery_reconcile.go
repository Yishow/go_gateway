package modbusshare

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/sqltx"
)

type formatRecoveryWitnessStore interface {
	CheckFormatRecoveryWitness(context.Context, string, string) error
}

// ReconcileFormatRecovery consumes only a server-built desired projection. It
// wraps every ordinary reconciler success exit (including aligned/replay) with
// the captured hydration and durable-clean evidence; it never publishes Ready.
func (r *Reconciler) ReconcileFormatRecovery(ctx context.Context, req ReconcileRequest, fence FormatRecoveryFence, witness string) (ReconcileOutcome, error) {
	store, ok := r.revStore.(formatRecoveryWitnessStore)
	if !ok || witness == "" || req.WorkspaceID == "" {
		return ReconcileOutcome{}, NewError(ErrCodeProjectionRequired, "format recovery evidence is unavailable", true)
	}
	check := func() error {
		if err := r.svc.CheckFormatRecoveryFence(ctx, fence); err != nil {
			return err
		}
		return store.CheckFormatRecoveryWitness(ctx, req.WorkspaceID, witness)
	}
	if err := check(); err != nil {
		return ReconcileOutcome{}, err
	}
	req.Restore = true
	out, err := r.Reconcile(WithFormatRecoveryWitness(ctx, req.WorkspaceID, witness), req)
	if err != nil {
		if errors.Is(err, sqltx.ErrCommitOutcomeUnknown) || errors.Is(err, sqltx.ErrConnectionUnusable) {
			unknown, stopErr := r.rollbackUnknown(context.WithoutCancel(ctx), req.WorkspaceID, fence.hydration)
			return unknown, errors.Join(err, stopErr)
		}
		return out, err
	}
	if err = check(); err != nil {
		return ReconcileOutcome{}, err
	}
	return out, nil
}
