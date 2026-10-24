package workspace

import (
	"context"
	"time"
)

// retryWhileSQLiteBusy repeats a create-once operation while SQLite reports a
// transient lock. The operations are idempotent, so a repeat is safe.
func retryWhileSQLiteBusy(ctx context.Context, operation func() (*WriteGroupSaveResult, error)) (*WriteGroupSaveResult, error) {
	for attempt := 0; ; attempt++ {
		result, err := operation()
		if err == nil || !isSQLiteBusyError(err) || attempt == basicManagedRetryAttempts-1 {
			return result, err
		}
		timer := time.NewTimer(10 * time.Millisecond)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

// currentGroupResult pairs an already saved group with the live workspace
// revision, as a replayed create-once call must return.
func (s *WriteGroupService) currentGroupResult(ctx context.Context, group *WriteGroup, operation string) (*WriteGroupSaveResult, error) {
	workspace, err := s.workspaceSvc.GetOrCreate(ctx)
	if err != nil {
		return nil, normalizeWriteGroupServiceError(operation, err)
	}
	if workspace.ID != group.WorkspaceID {
		return nil, writeGroupNotFound(operation)
	}
	return &WriteGroupSaveResult{WorkspaceRevision: workspace.DatabaseSetupRevision, Group: cloneWriteGroup(group)}, nil
}
