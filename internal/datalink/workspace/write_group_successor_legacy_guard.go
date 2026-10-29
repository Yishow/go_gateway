package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// checkLegacyCutoverOwner reads canonical metadata in one local snapshot before
// the legacy planners can probe targets. Commit repeats this check in its tx.
func (s *WriteGroupService) checkLegacyCutoverOwner(ctx context.Context, workspaceID, groupID string) (err error) {
	tx, err := s.repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if e := tx.Rollback(); e != nil && !errors.Is(e, sql.ErrTxDone) {
			err = errors.Join(err, e)
		}
	}()
	return s.checkLegacyCutoverOwnerInTx(ctx, tx, workspaceID, groupID)
}

func (s *WriteGroupService) checkLegacyCutoverOwnerInTx(ctx context.Context, tx *sql.Tx, workspaceID, groupID string) error {
	if err := s.repo.requireCurrentCanonicalOwnerInTx(ctx, tx, workspaceID, groupID); err != nil {
		return err
	}
	meta, found, err := s.repo.loadWriteGroupLineage(ctx, tx, workspaceID, groupID)
	if err != nil || !found || meta.Generation == 1 {
		return err
	}
	group, err := s.repo.GetInTx(ctx, tx, workspaceID, groupID)
	if err != nil || group.AppliedRevision == "" {
		return err
	}
	version, err := s.repo.VersionByRevisionInTx(ctx, tx, workspaceID, groupID, group.AppliedRevision)
	if err != nil {
		return err
	}
	if version.EffectiveAt.After(s.clockNow()) {
		return &WriteGroupLifecycleActionError{Action: "transition", ReasonCode: "successor_activation_pending", NextAction: "review_lifecycle", cause: ErrWriteGroupLifecycleBlocked}
	}
	return nil
}

// checkLegacyCutoverInTx cannot trust the earlier review's lineage or prior
// intent: both must still be eligible before any new revision/fence is written.
// Completed legacy exact replay remains outside this mutation-only guard.
func (s *WriteGroupService) checkLegacyCutoverInTx(ctx context.Context, tx *sql.Tx, workspaceID string, req CutoverRequest) error {
	if err := s.checkLegacyCutoverOwnerInTx(ctx, tx, workspaceID, req.SourceGroupID); err != nil {
		return err
	}
	if req.Kind != WriteGroupCutoverRollback {
		for _, target := range req.Targets {
			if err := s.checkLegacyCutoverOwnerInTx(ctx, tx, workspaceID, target.GroupID); err != nil {
				return err
			}
		}
		return nil
	}
	prior, err := s.repo.cutoverInTx(ctx, tx, workspaceID, req.CutoverID)
	if err != nil {
		return err
	}
	if prior.Kind != WriteGroupCutoverForward || prior.SourceGroupID != req.SourceGroupID {
		return fmt.Errorf("rollback source scope changed: %w", ErrCutoverNotApplicable)
	}
	if err := rejectSuccessorLegacyRollback(prior); err != nil {
		return err
	}
	for _, holder := range prior.Started {
		if err := s.checkLegacyCutoverOwnerInTx(ctx, tx, workspaceID, holder.GroupID); err != nil {
			return err
		}
	}
	return nil
}

func rejectSuccessorLegacyRollback(prior *WriteGroupCutover) error {
	if prior != nil && isSuccessorActivationIntent(prior.intent) {
		return fmt.Errorf("successor rollback requires a new generation: %w", ErrCutoverNotApplicable)
	}
	return nil
}
