package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

func (r *SQLWriteGroupRepository) insertVersionInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup, effectiveAt time.Time) error {
	if tx == nil || group == nil || group.ID == "" || group.Revision == "" {
		return fmt.Errorf("insert write-group version: %w", ErrWriteGroupValidation)
	}
	payload, err := json.Marshal(cloneWriteGroup(group))
	if err != nil {
		return fmt.Errorf("encode write-group version: %w", err)
	}
	_, err = tx.ExecContext(ctx, r.query(`
		INSERT INTO write_group_versions (group_id, group_revision, effective_at, payload, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`), group.ID, group.Revision, effectiveAt.UTC(), string(payload), r.now())
	if err != nil {
		return fmt.Errorf("insert write-group version: %w", err)
	}
	return nil
}

func (r *SQLWriteGroupRepository) setAppliedInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup) error {
	if tx == nil || group == nil || group.ID == "" || group.WorkspaceID == "" || group.Revision == "" {
		return fmt.Errorf("set write-group applied revision: %w", ErrWriteGroupValidation)
	}
	result, err := tx.ExecContext(ctx, r.query(`
		UPDATE write_groups
		SET applied_revision = $1, status = $2, updated_at = $3
		WHERE workspace_id = $4 AND id = $5 AND revision = $6
	`), group.AppliedRevision, group.Status, group.UpdatedAt, group.WorkspaceID, group.ID, group.Revision)
	if err != nil {
		return fmt.Errorf("set write-group applied revision: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("confirm write-group applied revision: %w", err)
	}
	if affected != 1 {
		return fmt.Errorf("set write-group applied revision: %w", ErrWriteGroupRevisionConflict)
	}
	return nil
}

func (r *SQLWriteGroupRepository) resolveVersionInTx(ctx context.Context, runner writeGroupSQLRunner, workspaceID, groupID string, at time.Time) (*WriteGroupAppliedSnapshot, error) {
	row := runner.QueryRowContext(ctx, r.query(`
		SELECT v.group_revision, v.effective_at, v.payload
		FROM write_group_versions v
		JOIN write_groups g ON g.id = v.group_id
		WHERE g.workspace_id = $1 AND v.group_id = $2 AND v.effective_at <= $3
		ORDER BY v.effective_at DESC, v.created_at DESC, v.group_revision DESC
		LIMIT 1
	`), workspaceID, groupID, at.UTC())
	var revision string
	var effectiveAt time.Time
	var payload string
	if err := row.Scan(&revision, &effectiveAt, &payload); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("resolve write-group version: %w", ErrWriteGroupAppliedRevisionUnavailable)
		}
		return nil, fmt.Errorf("read write-group version: %w", err)
	}
	var group WriteGroup
	if err := json.Unmarshal([]byte(payload), &group); err != nil {
		return nil, fmt.Errorf("decode write-group version: %w", err)
	}
	if group.ID != groupID || group.WorkspaceID != workspaceID || group.Revision != revision {
		return nil, fmt.Errorf("resolve write-group version: %w", ErrWriteGroupValidation)
	}
	return &WriteGroupAppliedSnapshot{
		WorkspaceID:     workspaceID,
		GroupID:         groupID,
		GroupRevision:   revision,
		AppliedRevision: group.AppliedRevision,
		EffectiveAt:     effectiveAt.UTC(),
		Group:           cloneWriteGroup(&group),
	}, nil
}

func (r *SQLWriteGroupRepository) versionByRevisionInTx(ctx context.Context, runner writeGroupSQLRunner, workspaceID, groupID, revision string) (*WriteGroupAppliedSnapshot, error) {
	if revision == "" {
		return nil, fmt.Errorf("read write-group version: %w", ErrWriteGroupAppliedRevisionUnavailable)
	}
	row := runner.QueryRowContext(ctx, r.query(`
		SELECT v.group_revision, v.effective_at, v.payload
		FROM write_group_versions v
		JOIN write_groups g ON g.id = v.group_id
		WHERE g.workspace_id = $1 AND v.group_id = $2 AND v.group_revision = $3
	`), workspaceID, groupID, revision)
	var version WriteGroupAppliedSnapshot
	var payload string
	if err := row.Scan(&version.GroupRevision, &version.EffectiveAt, &payload); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("read write-group version: %w", ErrWriteGroupAppliedRevisionUnavailable)
		}
		return nil, fmt.Errorf("read write-group version: %w", err)
	}
	var group WriteGroup
	if err := json.Unmarshal([]byte(payload), &group); err != nil {
		return nil, fmt.Errorf("decode write-group version: %w", err)
	}
	if group.ID != groupID || group.WorkspaceID != workspaceID || group.Revision != version.GroupRevision {
		return nil, fmt.Errorf("read write-group version: %w", ErrWriteGroupValidation)
	}
	version.WorkspaceID = workspaceID
	version.GroupID = groupID
	version.AppliedRevision = group.AppliedRevision
	version.EffectiveAt = version.EffectiveAt.UTC()
	version.Group = cloneWriteGroup(&group)
	return &version, nil
}
