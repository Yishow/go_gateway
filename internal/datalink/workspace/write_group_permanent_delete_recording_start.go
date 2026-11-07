package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

// checkPermanentDeleteRecordingStarts runs in the deletion transaction. An
// unresolved start can still save/apply/activate a selected recording.
func (s *WriteGroupService) checkPermanentDeleteRecordingStarts(ctx context.Context, tx *sql.Tx, workspaceID string, ids map[string]bool) error {
	rows, err := tx.QueryContext(ctx, s.repo.query(`SELECT detail FROM managed_schema_operations
		WHERE workspace_id=$1 AND action='recording_start' AND status NOT IN ('succeeded','failed')`), workspaceID)
	if err != nil {
		return fmt.Errorf("read recording start operations: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return fmt.Errorf("scan recording start operation: %w", err)
		}
		var detail struct {
			Version int `json:"version"`
			Request struct {
				WorkspaceID string          `json:"workspace_id"`
				Groups      json.RawMessage `json:"groups"`
			} `json:"request"`
		}
		if json.Unmarshal([]byte(raw), &detail) != nil || detail.Version != 1 || detail.Request.WorkspaceID != workspaceID || len(detail.Request.Groups) == 0 {
			return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
		}
		var groups []RecordingStartGroupIntent
		if json.Unmarshal(detail.Request.Groups, &groups) != nil {
			return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
		}
		for _, group := range groups {
			if strings.TrimSpace(group.GroupID) == "" || ids[group.GroupID] {
				return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate recording start operations: %w", err)
	}
	return nil
}
