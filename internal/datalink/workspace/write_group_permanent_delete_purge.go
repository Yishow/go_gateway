package workspace

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *WriteGroupService) purgeWriteGroupLineageInTx(ctx context.Context, tx *sql.Tx, groups []*WriteGroup, lineageID string) error {
	if err := s.purgeCompletedSuccessorOperationsInTx(ctx, tx, groups, lineageID); err != nil {
		return err
	}
	if err := s.purgeUnsharedStorageScopesInTx(ctx, tx, groups); err != nil {
		return err
	}
	for _, group := range groups {
		for _, table := range []string{"gw_delivery_receipts", "gw_delivery_outbox"} {
			if _, err := tx.ExecContext(ctx, s.repo.query(`DELETE FROM `+table+` WHERE destination_id=$1`), group.ID); err != nil {
				return fmt.Errorf("delete legacy recording rows from %s: %w", table, err)
			}
		}
		// A receipt has no group column: resolve it through the still-present outbox.
		if _, err := tx.ExecContext(ctx, s.repo.query(`DELETE FROM wg_delivery_receipts WHERE effect_key IN (SELECT effect_key FROM wg_delivery_outbox WHERE group_id=$1)`), group.ID); err != nil {
			return fmt.Errorf("delete recording receipts: %w", err)
		}
		if _, err := tx.ExecContext(ctx, s.repo.query(`DELETE FROM wg_delivery_sample_handoffs WHERE source_group_id=$1 OR target_group_id=$2`), group.ID, group.ID); err != nil {
			return fmt.Errorf("delete recording handoffs: %w", err)
		}
		for _, table := range []string{
			"wg_delivery_mixed_outcomes", "wg_delivery_mixed_progress", "wg_counter_proof_associations",
			"wg_delivery_counter_baselines", "wg_delivery_samples", "wg_delivery_checkpoints",
			"wg_delivery_buckets", "wg_delivery_outbox", "wg_runtime_versions",
			"write_group_ownership_fences", "write_group_revision_ends", "write_group_storage_events",
			"write_group_migration_maps", "write_group_members", "write_group_versions",
			"write_group_basic_keys", "write_group_profile_keys",
		} {
			if _, err := tx.ExecContext(ctx, s.repo.query(`DELETE FROM `+table+` WHERE group_id=$1`), group.ID); err != nil {
				return fmt.Errorf("delete recording rows from %s: %w", table, err)
			}
		}
		if _, err := tx.ExecContext(ctx, s.repo.query(`DELETE FROM managed_schema_operations WHERE operation_id IN
			(SELECT o.operation_id FROM managed_schema_operations o WHERE o.workspace_id=$1 AND `+settledDeleteOperation+` AND o.token IN
			(SELECT token FROM managed_schema_preview_tokens WHERE workspace_id=$2 AND plan_id=$3))`), group.WorkspaceID, group.WorkspaceID, group.ID); err != nil {
			return fmt.Errorf("delete recording schema operations: %w", err)
		}
		if _, err := tx.ExecContext(ctx, s.repo.query(`DELETE FROM managed_schema_preview_tokens WHERE workspace_id=$1 AND plan_id=$2`), group.WorkspaceID, group.ID); err != nil {
			return fmt.Errorf("delete recording schema previews: %w", err)
		}
	}
	for _, group := range groups {
		if _, err := tx.ExecContext(ctx, s.repo.query(`DELETE FROM write_group_cutovers WHERE workspace_id=$1 AND source_group_id=$2`), group.WorkspaceID, group.ID); err != nil {
			return fmt.Errorf("delete recording cutovers: %w", err)
		}
	}
	if lineageID != "" {
		for _, group := range groups {
			if _, err := tx.ExecContext(ctx, s.repo.query(`DELETE FROM write_group_generations WHERE group_id=$1`), group.ID); err != nil {
				return fmt.Errorf("delete recording generation: %w", err)
			}
		}
		if _, err := tx.ExecContext(ctx, s.repo.query(`DELETE FROM write_group_lineages WHERE lineage_id=$1`), lineageID); err != nil {
			return fmt.Errorf("delete recording lineage: %w", err)
		}
	}
	for _, group := range groups {
		if _, err := tx.ExecContext(ctx, s.repo.query(`DELETE FROM write_groups WHERE workspace_id=$1 AND id=$2`), group.WorkspaceID, group.ID); err != nil {
			return fmt.Errorf("delete recording group: %w", err)
		}
	}
	return nil
}
