package workspace

import (
	"context"
	"database/sql"
	"encoding/json"

	"go-gateway/internal/datalink/recordingplan"
)

// requireSuccessorNoHistoryInTx proves absence across the entire immutable
// lineage, not just its current head or a user-facing pending row count.
func (r *SQLWriteGroupRepository) requireSuccessorNoHistoryInTx(ctx context.Context, tx *sql.Tx, workspaceID, lineageID, predecessorID string) error {
	if tx == nil {
		return ErrWriteGroupValidation
	}
	meta, known, err := r.loadWriteGroupLineage(ctx, tx, workspaceID, predecessorID)
	if err != nil {
		return err
	}
	if !known || meta.LineageID != lineageID || meta.CurrentGroupID != predecessorID {
		return ErrWriteGroupLineageCorrupt
	}
	ids, err := r.lifecycleGenerationIDs(ctx, tx, lineageID)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return ErrWriteGroupLineageCorrupt
	}
	members := make(map[string]bool, len(ids))
	for _, id := range ids {
		members[id] = true
		group, err := r.GetInTx(ctx, tx, workspaceID, id)
		if err != nil {
			return err
		}
		member, known, err := r.loadWriteGroupLineage(ctx, tx, workspaceID, id)
		if err != nil {
			return err
		}
		if !known || member.LineageID != lineageID || group.AppliedRevision != "" || (id == meta.ReservedGroupID.String && group.Status != WriteGroupStatusDraft) || (id != meta.ReservedGroupID.String && group.Status != WriteGroupStatusDeleted) {
			return transitionRejection("successor-history-present")
		}
		// Deliberately no workspace predicate: a foreign row for this globally unique
		// group identity is contradictory evidence, not proof that history is absent.
		for _, table := range []string{"write_group_versions", "wg_runtime_versions", "wg_delivery_samples", "wg_delivery_outbox", "wg_delivery_checkpoints", "wg_delivery_buckets", "wg_delivery_counter_baselines", "wg_counter_proof_associations", "write_group_revision_ends", "write_group_ownership_fences"} {
			var count int
			if err := tx.QueryRowContext(ctx, r.query("SELECT COUNT(*) FROM "+table+" WHERE group_id=$1"), id).Scan(&count); err != nil {
				return err
			}
			if count != 0 {
				return transitionRejection("successor-history-present")
			}
		}
		var count int
		if err := tx.QueryRowContext(ctx, r.query(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs WHERE source_group_id=$1 OR target_group_id=$2`), id, id).Scan(&count); err != nil {
			return err
		}
		if count != 0 {
			return transitionRejection("successor-history-present")
		}
		if err := tx.QueryRowContext(ctx, r.query(`SELECT COUNT(*) FROM wg_delivery_receipts r JOIN wg_delivery_outbox o ON o.effect_key=r.effect_key WHERE o.group_id=$1`), id).Scan(&count); err != nil {
			return err
		}
		if count != 0 {
			return transitionRejection("successor-history-present")
		}
	}
	if err := r.requireNoLineageCutoverInTx(ctx, tx, workspaceID, members); err != nil {
		return err
	}
	return requireNoLineageStartInTx(ctx, tx, r, workspaceID, lineageID, members)
}

func (r *SQLWriteGroupRepository) requireNoLineageCutoverInTx(ctx context.Context, tx *sql.Tx, workspaceID string, ids map[string]bool) error {
	rows, err := tx.QueryContext(ctx, r.query(`SELECT workspace_id,source_group_id,detail FROM write_group_cutovers`))
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var storedWorkspace, source, raw string
		if err := rows.Scan(&storedWorkspace, &source, &raw); err != nil {
			return err
		}
		var detail cutoverDetail
		if json.Unmarshal([]byte(raw), &detail) != nil || len(detail.Ended) == 0 || len(detail.Started) == 0 {
			if storedWorkspace == workspaceID || ids[source] {
				return transitionRejection("successor-history-unverifiable")
			}
			continue
		}
		if ids[source] {
			return transitionRejection("successor-history-present")
		}
		for _, end := range detail.Ended {
			if ids[end.GroupID] {
				return transitionRejection("successor-history-present")
			}
		}
		for _, start := range detail.Started {
			if ids[start.GroupID] {
				return transitionRejection("successor-history-present")
			}
		}
	}
	return rows.Err()
}

func requireNoLineageStartInTx(ctx context.Context, tx *sql.Tx, r *SQLWriteGroupRepository, workspaceID, lineageID string, ids map[string]bool) error {
	rows, err := tx.QueryContext(ctx, r.query(`SELECT workspace_id,action,detail,payload_digest FROM managed_schema_operations WHERE action=$1`), recordingplan.RecordingStartAction)
	if err != nil {
		return err
	}
	for rows.Next() {
		var storedWorkspace, action, detail, digest string
		if err = rows.Scan(&storedWorkspace, &action, &detail, &digest); err != nil {
			rows.Close()
			return err
		}
		p, e := decodeRecordingStart(&recordingplan.SchemaOperation{Action: action, WorkspaceID: storedWorkspace, Detail: detail, PayloadDigest: digest})
		if e != nil {
			var partial recordingStartProgress
			attributed := storedWorkspace == workspaceID
			if json.Unmarshal([]byte(detail), &partial) == nil {
				for _, g := range partial.Request.Groups {
					attributed = attributed || ids[g.GroupID]
				}
				for _, g := range partial.View.Groups {
					attributed = attributed || ids[g.GroupID]
				}
			}
			if attributed {
				rows.Close()
				return transitionRejection("successor-history-unverifiable")
			}
			continue
		}
		for _, g := range p.Request.Groups {
			if ids[g.GroupID] {
				rows.Close()
				return transitionRejection("successor-history-present")
			}
		}
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err = rows.Close(); err != nil {
		return err
	}
	groupIDs := make([]string, 0, len(ids))
	for id := range ids {
		groupIDs = append(groupIDs, id)
	}
	starts, err := recordingplan.ReadSuccessorFirstStartsForLineageInTx(ctx, tx, workspaceID, lineageID, groupIDs)
	if err != nil {
		return err
	}
	for _, start := range starts {
		if start.Scope.LineageID == lineageID || ids[start.Scope.PredecessorGroupID] || ids[start.Scope.SuccessorGroupID] {
			return transitionRejection("successor-history-present")
		}
	}
	return nil
}
