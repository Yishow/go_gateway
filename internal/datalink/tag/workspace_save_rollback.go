package tag

import (
	"context"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/schema"
)

// ErrWorkspaceRollbackConflict indicates an intervening tag edit.
var ErrWorkspaceRollbackConflict = errors.New("workspace tag changed during save rollback")

type workspaceRollbackRepository interface {
	RestoreWorkspaceSave(context.Context, *schema.Tag, *schema.Tag) error
}

// RestoreWorkspaceSave restores owned fields only if the last authored tag
// values still match; it never overwrites an intervening direct tag edit.
func (s *Service) RestoreWorkspaceSave(ctx context.Context, expected, previous *schema.Tag) error {
	repo, ok := s.repo.(workspaceRollbackRepository)
	if !ok || expected == nil || previous == nil || expected.ID != previous.ID {
		return fmt.Errorf("tag save rollback unavailable")
	}
	return repo.RestoreWorkspaceSave(ctx, expected, previous)
}

// RestoreWorkspaceSave conditionally restores the last authored persisted fields.
func (r *SQLRepository) RestoreWorkspaceSave(ctx context.Context, expected, previous *schema.Tag) error {
	result, err := r.db.ExecContext(ctx, `UPDATE tags SET display_name=?,data_type=?,unit=?,description=?,status=?,labels=?,updated_at=? WHERE id=? AND key=? AND display_name=? AND data_type=? AND unit=? AND description=? AND status=? AND labels=?`, previous.DisplayName, previous.DataType, previous.Unit, previous.Description, previous.Status, previous.Labels, previous.UpdatedAt, expected.ID, expected.Key, expected.DisplayName, expected.DataType, expected.Unit, expected.Description, expected.Status, expected.Labels)
	if err != nil {
		return fmt.Errorf("restore failed workspace tag save: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrWorkspaceRollbackConflict
	}
	return nil
}

// RestoreWorkspaceSave compares and restores the owned fields under one mutex.
func (r *MemoryRepository) RestoreWorkspaceSave(_ context.Context, expected, previous *schema.Tag) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	current := r.tags[expected.ID]
	if current == nil || current.Key != expected.Key || current.DisplayName != expected.DisplayName || current.DataType != expected.DataType || current.Unit != expected.Unit || current.Description != expected.Description || current.Status != expected.Status || current.Labels != expected.Labels {
		return ErrWorkspaceRollbackConflict
	}
	restored := *previous
	r.tags[expected.ID] = &restored
	return nil
}
