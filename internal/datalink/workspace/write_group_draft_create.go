package workspace

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// prepareNewWriteGroupDraftInTx is shared by ordinary create and private
// reservation. Source ownership validation is never bypassed. Callers choose
// which canonical references to insert only after this preparation succeeds.
func (r *SQLWriteGroupRepository) prepareNewWriteGroupDraftInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup) (*WriteGroup, error) {
	if tx == nil {
		return nil, fmt.Errorf("create write group in transaction: %w", ErrWriteGroupValidation)
	}
	if group == nil {
		return nil, fmt.Errorf("create write group: %w: group is nil", ErrWriteGroupValidation)
	}
	candidate := cloneWriteGroup(group)
	// This field is derived from the persistent basic-managed key. Never let a
	// caller turn an ordinary group into a basic owner through the payload.
	candidate.BasicManagedDeviceID = ""
	candidate.ID = r.newID()
	defaultRecreateMissing(candidate)
	if candidate.Destination.StorageStrategy == WriteGroupStorageStrategyManaged && candidate.RowPolicy.RecordKeyColumn != "" && strings.TrimSpace(candidate.Destination.TableName) == "" && !candidate.isLongLayout() {
		candidate.Destination.TableName = managedGroupTableName(candidate.ID)
	}
	if err := r.validate(ctx, tx, candidate); err != nil {
		return nil, err
	}
	now := r.now()
	candidate.Revision = r.newID()
	candidate.AppliedRevision = ""
	candidate.Status = WriteGroupStatusDraft
	candidate.CreatedAt = now
	candidate.UpdatedAt = now
	return candidate, nil
}
