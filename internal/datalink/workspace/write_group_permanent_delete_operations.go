package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

// A successful schema operation has no test-write fields; a successful test
// write must have confirmed cleanup. A failed operation is safe only when it
// had no executed statements and its test-write cleanup has been resolved.
const settledDeleteOperation = `((o.status='succeeded' AND
	((o.write_outcome='' AND o.cleanup_status='') OR (o.write_outcome='written_verified' AND o.cleanup_status='cleaned')))
	OR (o.status='failed' AND o.executed_statements=0 AND
	((o.write_outcome='' AND o.cleanup_status='') OR (o.write_outcome='failed' AND o.cleanup_status='not_attempted')
	OR (o.write_outcome IN ('written_verified','written_unverified','failed') AND o.cleanup_status='cleaned'))))`

func (s *WriteGroupService) checkPermanentDeleteOperations(ctx context.Context, tx *sql.Tx, group *WriteGroup) error {
	lineage, _, err := s.repo.loadWriteGroupLineage(ctx, tx, group.WorkspaceID, group.ID)
	if err != nil {
		return err
	}
	ids, _, err := s.permanentDeleteIDsForLineage(ctx, tx, group.ID, lineage)
	if err != nil {
		return err
	}
	return s.checkPermanentDeleteOperationsForScope(ctx, tx, group, ids, lineage)
}

func (s *WriteGroupService) checkPermanentDeleteOperationsForScope(ctx context.Context, tx *sql.Tx, group *WriteGroup, ids map[string]bool, lineage *writeGroupLineageMetadata) error {
	if err := s.checkPermanentDeleteRecordingStarts(ctx, tx, group.WorkspaceID, ids); err != nil {
		return err
	}
	for id := range ids {
		var pending int
		err := tx.QueryRowContext(ctx, s.repo.query(`SELECT COUNT(*) FROM managed_schema_operations o
			JOIN managed_schema_preview_tokens p ON p.token=o.token
			WHERE p.workspace_id=$1 AND p.plan_id=$2 AND o.workspace_id=$3 AND NOT `+settledDeleteOperation), group.WorkspaceID, id, group.WorkspaceID).Scan(&pending)
		if err != nil {
			return fmt.Errorf("read recording schema operations: %w", err)
		}
		if pending > 0 {
			return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
		}
	}
	rows, err := tx.QueryContext(ctx, s.repo.query(`SELECT o.detail FROM managed_schema_operations o
		WHERE o.workspace_id=$1 AND o.action IN ('successor_reserve','successor_first_start') AND NOT `+settledDeleteOperation), group.WorkspaceID)
	if err != nil {
		return fmt.Errorf("read pending successor operations: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return fmt.Errorf("scan successor operation: %w", err)
		}
		// Successor operation scopes use snake_case JSON field names.
		var envelope struct {
			Scope struct {
				WorkspaceID        string `json:"workspace_id"`
				LineageID          string `json:"lineage_id"`
				PredecessorGroupID string `json:"predecessor_group_id"`
				SuccessorGroupID   string `json:"successor_group_id"`
			} `json:"scope"`
		}
		if json.Unmarshal([]byte(raw), &envelope) != nil || envelope.Scope.WorkspaceID != group.WorkspaceID || envelope.Scope.LineageID == "" || envelope.Scope.PredecessorGroupID == "" || envelope.Scope.SuccessorGroupID == "" {
			return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
		}
		if ids[envelope.Scope.PredecessorGroupID] || ids[envelope.Scope.SuccessorGroupID] || (lineage != nil && lineage.LineageID == envelope.Scope.LineageID) {
			return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate successor operations: %w", err)
	}
	return nil
}

func (s *WriteGroupService) purgeCompletedSuccessorOperationsInTx(ctx context.Context, tx *sql.Tx, groups []*WriteGroup, lineageID string) error {
	if lineageID == "" || len(groups) == 0 {
		return nil
	}
	ids := make(map[string]bool, len(groups))
	for _, group := range groups {
		ids[group.ID] = true
	}
	rows, err := tx.QueryContext(ctx, s.repo.query(`SELECT o.operation_id,o.detail FROM managed_schema_operations o
		WHERE o.workspace_id=$1 AND o.action IN ('successor_reserve','successor_first_start') AND `+settledDeleteOperation), groups[0].WorkspaceID)
	if err != nil {
		return fmt.Errorf("read completed successor operations: %w", err)
	}
	var operations []string
	for rows.Next() {
		var id, raw string
		if err := rows.Scan(&id, &raw); err != nil {
			rows.Close()
			return fmt.Errorf("scan completed successor operation: %w", err)
		}
		var detail struct {
			Scope struct {
				WorkspaceID        string `json:"workspace_id"`
				LineageID          string `json:"lineage_id"`
				PredecessorGroupID string `json:"predecessor_group_id"`
				SuccessorGroupID   string `json:"successor_group_id"`
			} `json:"scope"`
		}
		if json.Unmarshal([]byte(raw), &detail) != nil || detail.Scope.WorkspaceID != groups[0].WorkspaceID || detail.Scope.LineageID == "" {
			rows.Close()
			return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
		}
		if detail.Scope.LineageID != lineageID {
			continue
		}
		if !ids[detail.Scope.PredecessorGroupID] || !ids[detail.Scope.SuccessorGroupID] {
			rows.Close()
			return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
		}
		operations = append(operations, id)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate completed successor operations: %w", err)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, id := range operations {
		if _, err := tx.ExecContext(ctx, s.repo.query(`DELETE FROM managed_schema_operations WHERE workspace_id=$1 AND operation_id=$2`), groups[0].WorkspaceID, id); err != nil {
			return fmt.Errorf("delete completed successor operation: %w", err)
		}
	}
	return nil
}
