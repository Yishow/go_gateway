package mapping

import (
	"context"
	"fmt"
	"go-gateway/internal/datalink/schema"
)

type workspaceRollbackRepository interface {
	RestoreWorkspaceSave(context.Context, *schema.Mapping, *schema.Mapping) error
}

// RestoreWorkspaceSave restores a failed save only while the exact last draft
// is still current. An intervening direct edit must never be overwritten.
func (s *Service) RestoreWorkspaceSave(ctx context.Context, expected, previous *schema.Mapping) error {
	repo, ok := s.repo.(workspaceRollbackRepository)
	if !ok || expected == nil || previous == nil || expected.ID != previous.ID {
		return fmt.Errorf("mapping save rollback unavailable")
	}
	return repo.RestoreWorkspaceSave(ctx, expected, previous)
}

// RestoreWorkspaceSave conditionally restores the last authored persisted fields.
func (r *SQLRepository) RestoreWorkspaceSave(ctx context.Context, expected, previous *schema.Mapping) error {
	result, err := r.db.ExecContext(ctx, `UPDATE mappings SET tag_id=?,transform_pipeline=?,status=?,rule_candidate_id=?,proposed_signature=?,last_applied_signature=?,blocking_reason=?,enabled=?,updated_at=? WHERE id=? AND point_id=? AND tag_id=? AND transform_pipeline=? AND status=? AND rule_candidate_id=? AND proposed_signature=? AND last_applied_signature=? AND blocking_reason=? AND enabled=?`, previous.TagID, previous.TransformPipeline, previous.Status, previous.RuleCandidateID, previous.ProposedSignature, previous.LastAppliedSignature, previous.BlockingReason, previous.Enabled, previous.UpdatedAt, expected.ID, expected.PointID, expected.TagID, expected.TransformPipeline, expected.Status, expected.RuleCandidateID, expected.ProposedSignature, expected.LastAppliedSignature, expected.BlockingReason, expected.Enabled)
	if err != nil {
		return fmt.Errorf("restore failed mapping save: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrWorkspaceConfirmationConflict
	}
	return nil
}

// RestoreWorkspaceSave compares and restores the owned fields under one mutex.
func (r *MemoryRepository) RestoreWorkspaceSave(_ context.Context, expected, previous *schema.Mapping) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	current := r.mappings[expected.ID]
	if current == nil || current.PointID != expected.PointID || current.TagID != expected.TagID || current.TransformPipeline != expected.TransformPipeline || current.Status != expected.Status || current.RuleCandidateID != expected.RuleCandidateID || current.ProposedSignature != expected.ProposedSignature || current.LastAppliedSignature != expected.LastAppliedSignature || current.BlockingReason != expected.BlockingReason || current.Enabled != expected.Enabled {
		return ErrWorkspaceConfirmationConflict
	}
	restored := *previous
	r.mappings[expected.ID] = &restored
	return nil
}
