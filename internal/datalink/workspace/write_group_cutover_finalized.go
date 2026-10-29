package workspace

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
)

// ErrCutoverFinalizedOverlap requires a fresh review at a boundary that does
// not reassign any part of an already finalized owner's range.
var ErrCutoverFinalizedOverlap = errors.New("recording-cutover-finalized-overlap")

// requireNoFinalizedOverlap runs only in the configuration commit transaction.
// Reclaim removes journal/outbox/bucket history, but never the closure cursor:
// next_close is the first UNclosed bucket start, not the last observation time.
// Runtime never persists the initial in-memory cursor; CommitClosure advances
// it atomically with effects, consumed samples and counter baselines. Equality
// is legal: a closed half-open range ending at boundary takes no new-owner time.
func requireNoFinalizedOverlap(ctx context.Context, tx *sql.Tx, plan *cutoverPlan) error {
	// Share Store.inTx's SQLite write authority before reading closure evidence.
	// If a deferred transaction already has an obsolete read snapshot, its lock
	// upgrade fails safely rather than authorizing a stale cutover.
	if _, err := tx.ExecContext(ctx, `UPDATE wg_delivery_checkpoints SET updated_at=updated_at WHERE 0`); err != nil {
		return groupdelivery.ErrOwnershipUnavailable
	}
	for _, ended := range plan.ended {
		var next string
		err := tx.QueryRowContext(ctx, `SELECT next_close FROM wg_delivery_checkpoints WHERE group_id=? AND group_revision=?`, ended.GroupID, ended.Revision).Scan(&next)
		if errors.Is(err, sql.ErrNoRows) {
			// No cursor is normal only before the first closure. Retained
			// effects without their atomically committed cursor mean authority
			// is incomplete, not permission to rewrite historical ownership.
			if err := requireNoOrphanClosure(ctx, tx, ended); err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return groupdelivery.ErrOwnershipUnavailable
		}
		cursor, err := time.Parse(time.RFC3339Nano, next)
		if err != nil || cursor.IsZero() {
			return groupdelivery.ErrOwnershipUnavailable
		}
		if cursor.After(plan.req.Boundary) {
			return ErrCutoverFinalizedOverlap
		}
	}
	return nil
}

func requireNoOrphanClosure(ctx context.Context, tx *sql.Tx, ended CutoverEnd) error {
	var exists bool
	err := tx.QueryRowContext(ctx, `SELECT
 EXISTS(SELECT 1 FROM wg_delivery_samples WHERE group_id=? AND group_revision=? AND consumed=1)
 OR EXISTS(SELECT 1 FROM wg_delivery_buckets WHERE group_id=? AND group_revision=?)
 OR EXISTS(SELECT 1 FROM wg_delivery_outbox WHERE group_id=? AND group_revision=?)
 OR EXISTS(SELECT 1 FROM wg_delivery_counter_baselines WHERE group_id=? AND group_revision=?)`,
		ended.GroupID, ended.Revision, ended.GroupID, ended.Revision,
		ended.GroupID, ended.Revision, ended.GroupID, ended.Revision).Scan(&exists)
	if err != nil || exists {
		return groupdelivery.ErrOwnershipUnavailable
	}
	return nil
}
