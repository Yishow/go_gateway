package modbusshare

import (
	"context"
	"strings"
	"sync"
)

// WorkspaceRevisionStore handles workspace revision CAS and dirty-state tracking.
type WorkspaceRevisionStore interface {
	GetRevision(ctx context.Context, wsID string) (revision string, isDirty bool, err error)
	UpdateRevision(ctx context.Context, wsID, expectedRev, newRev string) error
	MarkDirty(ctx context.Context, wsID string) error
}

// RevisionStatus is the durable CAS boundary used by status/bootstrap
// consumers. Dirty is intentionally read from the existing revision store;
// no parallel recovery record is maintained.
type RevisionStatus struct {
	Revision string
	Dirty    bool
}

// Reconciler coordinates atomic, serialized, idempotent, CAS-checked Share reconciliations.
type Reconciler struct {
	svc            *Service
	revStore       WorkspaceRevisionStore
	workspaceLocks sync.Map
	// projectionMu serializes the global memory-bank transaction. Workspace
	// locks still preserve per-workspace ownership/CAS semantics, while this
	// lock prevents a rollback snapshot from crossing another workspace commit.
	projectionMu  sync.Mutex
	idempotencyMu sync.Mutex
	idempotency   map[string]idempotencyRecord
	ownership     func(context.Context, DesiredMapping) error
}

type idempotencyRecord struct {
	outcome               ReconcileOutcome
	postWorkspaceRevision string
	postSettingsRevision  string
}

// NewReconciler creates a new Modbus Share Reconciler.
func NewReconciler(svc *Service, revStore WorkspaceRevisionStore) *Reconciler {
	return &Reconciler{
		svc:         svc,
		revStore:    revStore,
		idempotency: make(map[string]idempotencyRecord),
		ownership:   rejectUnprovenOwnership,
	}
}

// WorkspaceRevisionStore exposes the durable CAS seam to adjacent lifecycle
// boundaries such as workspace membership invalidation.
func (r *Reconciler) WorkspaceRevisionStore() WorkspaceRevisionStore {
	if r == nil {
		return nil
	}
	return r.revStore
}

// RevisionStatus reads the durable workspace revision and dirty marker for
// status and bootstrap projections after process restart.
func (r *Reconciler) RevisionStatus(ctx context.Context, workspaceID string) (RevisionStatus, error) {
	if r == nil || r.revStore == nil {
		return RevisionStatus{}, NewError(ErrCodeProjectionRequired, "durable Share revision store is required", true)
	}
	if workspaceID == "" {
		return RevisionStatus{}, NewError(ErrCodeWorkspaceScope, "workspace id is required", false)
	}
	revision, dirty, err := r.revStore.GetRevision(ctx, workspaceID)
	if err != nil {
		return RevisionStatus{}, err
	}
	return RevisionStatus{Revision: revision, Dirty: dirty}, nil
}

// WithOwnershipValidator requires every desired mapping to be proven by the
// persisted workspace/source-rule/tag relationship before projection.
func (r *Reconciler) WithOwnershipValidator(validator func(context.Context, DesiredMapping) error) *Reconciler {
	if validator == nil {
		r.ownership = rejectUnprovenOwnership
		return r
	}
	r.ownership = validator
	return r
}

func rejectUnprovenOwnership(_ context.Context, desired DesiredMapping) error {
	return NewError(ErrCodeWorkspaceScope, "mapping ownership could not be proven", true)
}

func (r *Reconciler) getWorkspaceLock(wsID string) *sync.Mutex {
	lock, _ := r.workspaceLocks.LoadOrStore(wsID, &sync.Mutex{})
	mutex, ok := lock.(*sync.Mutex)
	if !ok {
		panic("modbus share workspace lock has unexpected type")
	}
	return mutex
}

// markDirty records an uncertain projection when durable revision storage is
// available. A nil store cannot persist the marker, but rollback callers still
// return the typed dirty_unknown error rather than panicking.
func (r *Reconciler) markDirty(ctx context.Context, wsID string) error {
	if r.revStore == nil {
		return nil
	}
	return r.revStore.MarkDirty(ctx, wsID)
}

func (r *Reconciler) rollbackUnknown(ctx context.Context, workspaceID string, hydration HydrationState) (ReconcileOutcome, error) {
	// Fail closed first. The durable dirty marker is useful for recovery, but a
	// failed marker must never leave a possibly active listener/readiness gate.
	r.svc.FailClosed(hydration)
	message := "runtime or durable rollback is uncertain"
	if err := r.markDirty(ctx, workspaceID); err != nil {
		message = failedMarkWorkspaceDirtyMessage
	}
	return ReconcileOutcome{Outcome: dirtyUnknownState, NewWorkspaceRevision: hydration.WorkspaceRevision, SettingsRevision: hydration.SettingsRevision}, &Error{
		Code:              ErrCodeDirtyUnknown,
		Message:           message,
		Retryable:         true,
		DirtyState:        dirtyUnknownState,
		WorkspaceRevision: hydration.WorkspaceRevision,
		SettingsRevision:  hydration.SettingsRevision,
		Action:            runRecoveryReconcileAction,
	}
}

func (r *Reconciler) validateDurableInputs(req ReconcileRequest) (DurableDesiredMappingStore, error) {
	if r == nil || r.revStore == nil {
		return nil, NewError(ErrCodeProjectionRequired, "durable Share revision store is required", false)
	}
	durable, ok := r.revStore.(DurableDesiredMappingStore)
	if !ok || durable == nil {
		return nil, NewError(ErrCodeProjectionRequired, "durable Share revision store is required", false)
	}
	if strings.TrimSpace(req.ExpectedWorkspaceRevision) == "" || strings.TrimSpace(req.ExpectedSettingsRevision) == "" {
		return nil, NewError(ErrCodeRevisionRequired, "workspace and settings revisions are required", false)
	}
	return durable, nil
}
