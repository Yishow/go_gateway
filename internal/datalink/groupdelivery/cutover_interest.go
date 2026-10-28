package groupdelivery

import (
	"context"
	"database/sql"
	"time"
)

// HasCutoverSource is only an interest hint for a missing intake worker. It
// never authorizes acceptance; AppendSample still checks the durable fence in
// its write transaction. Disabled/deleted groups are not reactivated by a hint.
// All validated affected revisions are candidates, including ended sources that
// can still own delayed pre-boundary observations. A canonical active historical
// revision can cause extra routing attempts, never acceptance outside its fence.
// No scheduler clock is needed: observed_at is checked only by decisive intake.
func (s *Store) HasCutoverSource(ctx context.Context, source OwnershipSource) (bool, error) {
	if source.WorkspaceID == "" || source.DeviceID == "" || source.PointID == "" || source.TagID == "" {
		return false, ErrOwnershipUnavailable
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return false, ErrOwnershipUnavailable
	}
	defer tx.Rollback()
	ownership, err := resolveCutoverOwnership(ctx, tx, source, time.Time{}, GroupKey{})
	if err != nil {
		return false, err
	}
	interested := false
	for owner := range ownership.Affected {
		var status, applied string
		if tx.QueryRowContext(ctx, `SELECT status,applied_revision FROM write_groups WHERE workspace_id=? AND id=?`, owner.WorkspaceID, owner.GroupID).Scan(&status, &applied) != nil {
			return false, ErrOwnershipUnavailable
		}
		switch status {
		case "disabled", "deleted":
		case "draft", "ready", "running":
			interested = interested || applied == owner.GroupRevision
		default:
			return false, ErrOwnershipUnavailable
		}
	}
	if tx.Commit() != nil {
		return false, ErrOwnershipUnavailable
	}
	return interested, nil
}
