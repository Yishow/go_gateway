package groupdelivery

import (
	"context"
	"database/sql"
	"time"
)

// WithCanonicalIntake returns a clone whose new samples require canonical
// write-group lifecycle authority in the journal transaction. NewStore remains
// available for standalone ledgers whose providers do not own write_groups.
// Configure this before sharing the store; clock/quota clones retain the mode.
func (s *Store) WithCanonicalIntake() *Store {
	clone := *s
	clone.canonicalIntake = true
	return &clone
}

// Only new native acceptance consults lifecycle. Already accepted retries and
// closure/delivery retain their original immutable authority. inTx holds the
// SQLite write lock, so a stop either precedes this check or follows the append
// commit; a cached boundary can never authorize a post-stop journal insertion.
func (s *Store) validateLifecycleIntake(ctx context.Context, tx *sql.Tx, key GroupKey, observedAt time.Time) error {
	if !s.canonicalIntake {
		return nil
	}
	var workspaceID, status, applied string
	err := tx.QueryRowContext(ctx, `SELECT workspace_id,status,applied_revision FROM write_groups WHERE id=?`, key.GroupID).Scan(&workspaceID, &status, &applied)
	if err != nil || workspaceID != key.WorkspaceID {
		return ErrOwnershipUnavailable
	}
	switch status {
	case "disabled", "deleted":
		return ErrIntakeStopped
	case "draft", "ready", "running":
		if applied == "" {
			return ErrOwnershipUnavailable
		}
	default:
		return ErrOwnershipUnavailable
	}
	// Current draft/applied revisions do not replace historical ownership. The
	// exact immutable version must exist. First start has no cutover fence, so
	// its own start bound must also be enforced in this journal transaction.
	_, effectiveAt, err := readOwnershipVersion(ctx, tx, key.WorkspaceID, ownershipRef{GroupID: key.GroupID, Revision: key.GroupRevision})
	if err != nil || effectiveAt.IsZero() {
		return ErrOwnershipUnavailable
	}
	if observedAt.Before(effectiveAt) {
		return ErrOwnershipConflict
	}
	return nil
}
