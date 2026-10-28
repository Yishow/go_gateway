package workspace

import (
	"context"
	"database/sql"
	"errors"
)

// requireNoOtherSuccessorActiveOwnerInTx preserves ordinary active-wide/custom
// ownership checks while excluding exactly the predecessor validated in this tx.
// It is not a caller-selectable exception for ordinary Apply.
func (s *WriteGroupService) requireNoOtherSuccessorActiveOwnerInTx(ctx context.Context, tx *sql.Tx, pred, candidate *WriteGroup) error {
	if tx == nil || pred == nil || candidate == nil || pred.WorkspaceID != candidate.WorkspaceID {
		return ErrWriteGroupValidation
	}
	if !candidate.isLongLayout() && !candidate.isWideProfile() {
		return nil
	}
	meta, found, err := s.repo.loadWriteGroupLineage(ctx, tx, pred.WorkspaceID, pred.ID)
	if err != nil {
		return err
	}
	if !found || meta.CurrentGroupID != pred.ID || meta.ReservedGroupID.String != candidate.ID {
		return ErrWriteGroupLineageCorrupt
	}
	rows, err := tx.QueryContext(ctx, s.repo.query(`SELECT id,applied_revision FROM write_groups WHERE workspace_id=$1 AND id<>$2 AND id<>$3 AND applied_revision<>'' AND status IN ('draft','ready','running')`), candidate.WorkspaceID, candidate.ID, pred.ID)
	if err != nil {
		return err
	}
	type owner struct{ id, revision string }
	var owners []owner
	for rows.Next() {
		var o owner
		if err := rows.Scan(&o.id, &o.revision); err != nil {
			rows.Close()
			return err
		}
		owners = append(owners, o)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	wanted := make(map[CutoverMemberRef]bool, len(candidate.Members))
	for _, m := range candidate.Members {
		wanted[memberRefOf(m)] = true
	}
	for _, o := range owners {
		end, err := s.repo.revisionEndInTx(ctx, tx, o.id, o.revision)
		if err != nil {
			return err
		}
		if !end.IsZero() && !end.After(s.clockNow()) {
			continue
		}
		version, err := s.repo.versionByRevisionInTx(ctx, tx, candidate.WorkspaceID, o.id, o.revision)
		if errors.Is(err, ErrWriteGroupAppliedRevisionUnavailable) {
			return ErrWriteGroupLifecycleBlocked
		}
		if err != nil {
			return err
		}
		if version.Group.isLongLayout() {
			continue
		}
		if candidate.isWideProfile() && version.Group.isWideProfile() && profileKindOf(candidate) != profileKindOf(version.Group) {
			continue
		}
		for _, m := range version.Group.Members {
			if wanted[memberRefOf(m)] {
				return ErrWriteGroupCutoverRequired
			}
		}
	}
	return nil
}
