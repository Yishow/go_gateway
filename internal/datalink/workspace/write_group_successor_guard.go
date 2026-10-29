package workspace

import (
	"context"
	"database/sql"
	"errors"
)

func successorReservedError() error {
	return &WriteGroupLifecycleActionError{Action: "successor", ReasonCode: "successor_reserved", NextAction: "review_lifecycle", cause: ErrWriteGroupLifecycleBlocked}
}

func isSuccessorReservedError(err error) bool {
	var action *WriteGroupLifecycleActionError
	return errors.As(err, &action) && action.ReasonCode == "successor_reserved"
}

// Reservation authority is read from the lineage, never inferred from the
// draft status, caller fields, schema proof or a process-local mutex.
func (r *SQLWriteGroupRepository) requireUnreservedWriteGroupInTx(ctx context.Context, runner writeGroupSQLRunner, workspaceID, groupID string) error {
	metadata, found, err := r.loadWriteGroupLineage(ctx, runner, workspaceID, groupID)
	if err != nil {
		return err
	}
	if found && metadata.ReservedGroupID.Valid && metadata.ReservedGroupID.String == groupID {
		return successorReservedError()
	}
	return nil
}

// This read-only preflight prevents probing a reserved destination. Mutations
// still recheck reservation authority inside their own transaction.
func (s *WriteGroupService) requireUnreservedWriteGroup(ctx context.Context, workspaceID, groupID string) (err error) {
	if err = s.validate(); err != nil {
		return err
	}
	tx, err := s.repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, rollbackErr)
		}
	}()
	if _, err = s.repo.GetInTx(ctx, tx, workspaceID, groupID); err != nil {
		return err
	}
	return s.repo.requireUnreservedWriteGroupInTx(ctx, tx, workspaceID, groupID)
}
