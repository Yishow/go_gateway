package workspace

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
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
		if through, found, err := mixedFinalizedThroughInTx(ctx, tx, ended); err != nil {
			return err
		} else if found {
			if through.After(plan.req.Boundary) {
				return ErrCutoverFinalizedOverlap
			}
			continue
		}
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
 OR EXISTS(SELECT 1 FROM wg_delivery_counter_baselines WHERE group_id=? AND group_revision=?)
 OR EXISTS(SELECT 1 FROM wg_delivery_mixed_outcomes WHERE group_id=? AND group_revision=?)`,
		ended.GroupID, ended.Revision, ended.GroupID, ended.Revision,
		ended.GroupID, ended.Revision, ended.GroupID, ended.Revision, ended.GroupID, ended.Revision).Scan(&exists)
	if err != nil || exists {
		return groupdelivery.ErrOwnershipUnavailable
	}
	return nil
}

// Mixed closure keeps its original per-member progress, never a legacy cursor.
// Authenticate that existing scheduling projection in the cutover transaction;
// no sample, baseline or ledger record is rewritten by this adapter.
func mixedFinalizedThroughInTx(ctx context.Context, tx *sql.Tx, ended CutoverEnd) (time.Time, bool, error) {
	unavailable := func() (time.Time, bool, error) { return time.Time{}, false, groupdelivery.ErrOwnershipUnavailable }
	var count, legacy int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_mixed_progress WHERE group_id=? AND group_revision=?`, ended.GroupID, ended.Revision).Scan(&count); err != nil {
		return unavailable()
	}
	if count == 0 {
		// A mixed record never substitutes a legacy cursor for missing progress.
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_checkpoints c JOIN write_groups g ON g.id=c.group_id WHERE c.group_id=? AND c.group_revision=? AND g.profile_kind='mixed'`, ended.GroupID, ended.Revision).Scan(&legacy); err != nil || legacy != 0 {
			return unavailable()
		}
		return time.Time{}, false, nil
	}
	if count != 1 {
		return unavailable()
	}
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_checkpoints WHERE group_id=? AND group_revision=?`, ended.GroupID, ended.Revision).Scan(&legacy); err != nil || legacy != 0 {
		return unavailable()
	}
	var workspace, ownerWorkspace, raw, digest, descriptorDigest, through, payload string
	var version int
	if err := tx.QueryRowContext(ctx, `SELECT p.workspace_id,g.workspace_id,p.version,p.checkpoint,p.checkpoint_digest,p.descriptor_digest,p.closed_through,v.payload
 FROM wg_delivery_mixed_progress p JOIN wg_runtime_versions v ON v.workspace_id=p.workspace_id AND v.group_id=p.group_id AND v.group_revision=p.group_revision
 JOIN write_groups g ON g.id=p.group_id WHERE p.group_id=? AND p.group_revision=?`, ended.GroupID, ended.Revision).Scan(&workspace, &ownerWorkspace, &version, &raw, &digest, &descriptorDigest, &through, &payload); err != nil {
		return unavailable()
	}
	if workspace != ownerWorkspace || version != 1 || fmt.Sprintf("%x", sha256.Sum256([]byte(raw))) != digest || fmt.Sprintf("%x", sha256.Sum256([]byte(payload))) != descriptorDigest || !counterproof.UniqueJSONObject([]byte(raw), false) {
		return unavailable()
	}
	evidence, err := counterproof.ResolveSnapshotEvidence([]byte(payload), counterproof.SnapshotScope{WorkspaceID: workspace, GroupID: ended.GroupID, AppliedRevision: ended.Revision})
	if err != nil {
		return unavailable()
	}
	var frozen WriteGroupAppliedSnapshot
	if json.Unmarshal([]byte(payload), &frozen) != nil || frozen.Group == nil || frozen.GroupRevision != ended.Revision || frozen.WorkspaceID != workspace || frozen.GroupID != ended.GroupID || frozen.AppliedRevision != ended.Revision {
		return unavailable()
	}
	types := map[string]schema.DataType{}
	intervals := map[string]time.Duration{}
	for _, m := range frozen.Group.Members {
		if m.Cadence == nil {
			return unavailable()
		}
		duration, err := m.Cadence.Duration()
		if err != nil {
			return unavailable()
		}
		key := (CounterBindingMemberIdentity{DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID}).SampleMemberKey()
		if _, duplicate := intervals[key]; duplicate {
			return unavailable()
		}
		intervals[key], types[m.TagID] = duration, m.SourceValueType
	}
	if ValidateMixedRecordingSnapshot(frozen.Group, types, evidence.Bindings) != nil {
		return unavailable()
	}
	var cp snapshot.MixedCheckpoint
	if decodeSuccessorJSON([]byte(raw), &cp) != nil || cp.Version != 1 || cp.Progress.WorkspaceID != workspace || cp.Progress.GroupID != ended.GroupID || cp.Progress.GroupRevision != ended.Revision || !cp.StartsAt.Equal(frozen.EffectiveAt) || snapshot.ValidateMixedCheckpointCadences(cp, intervals) != nil {
		return unavailable()
	}
	cursor, err := time.Parse(time.RFC3339Nano, through)
	if err != nil || cursor.IsZero() || !cursor.Equal(cp.Progress.ClosedThrough) {
		return unavailable()
	}
	return cursor, true, nil
}
