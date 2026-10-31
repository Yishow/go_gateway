package modbusshare

import "context"

type recoveryWorkspaceLockKey struct{}
type recoveryWorkspaceLock struct {
	reconciler  *Reconciler
	workspaceID string
}

// WithFormatRecovery serializes the whole explicit attempt with the existing
// workspace reconcile authority. The private context prevents recursive locking;
// it is neither a second lock registry nor browser-controlled authorization.
func (r *Reconciler) WithFormatRecovery(ctx context.Context, workspaceID string, run func(context.Context) error) error {
	if r == nil || workspaceID == "" || run == nil {
		return NewError(ErrCodeWorkspaceScope, "format recovery scope is unavailable", false)
	}
	lock := r.getWorkspaceLock(workspaceID)
	lock.Lock()
	defer lock.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	return run(context.WithValue(ctx, recoveryWorkspaceLockKey{}, recoveryWorkspaceLock{r, workspaceID}))
}

func (r *Reconciler) lockReconcileWorkspace(ctx context.Context, workspaceID string) func() {
	if held, ok := ctx.Value(recoveryWorkspaceLockKey{}).(recoveryWorkspaceLock); ok && held.reconciler == r && held.workspaceID == workspaceID {
		return func() {}
	}
	lock := r.getWorkspaceLock(workspaceID)
	lock.Lock()
	return lock.Unlock
}
