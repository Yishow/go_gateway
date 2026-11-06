package mapping

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/schema"
)

type sourceRuleRetirementRepository interface {
	deleteRetiredMapping(context.Context, *schema.Mapping) error
	restoreRetiredMapping(context.Context, *schema.Mapping) error
}

// DeleteRetiredMapping removes only the mapping revision observed by a rule edit.
func (s *Service) DeleteRetiredMapping(ctx context.Context, expected *schema.Mapping) error {
	repo, ok := s.repo.(sourceRuleRetirementRepository)
	if !ok || expected == nil {
		return fmt.Errorf("mapping retirement unavailable")
	}
	return repo.deleteRetiredMapping(ctx, expected)
}

// RestoreRetiredMapping inserts the removed mapping only while its endpoints remain unbound.
func (s *Service) RestoreRetiredMapping(ctx context.Context, record *schema.Mapping) error {
	repo, ok := s.repo.(sourceRuleRetirementRepository)
	if !ok || record == nil {
		return fmt.Errorf("mapping retirement recovery unavailable")
	}
	return repo.restoreRetiredMapping(ctx, record)
}

func (r *SQLRepository) deleteRetiredMapping(ctx context.Context, m *schema.Mapping) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM mappings WHERE id=? AND point_id=? AND tag_id=? AND transform_pipeline=? AND status=? AND COALESCE(rule_candidate_id,'')=? AND COALESCE(proposed_signature,'')=? AND COALESCE(last_applied_signature,'')=? AND COALESCE(blocking_reason,'')=? AND enabled=?`, m.ID, m.PointID, m.TagID, m.TransformPipeline, m.Status, m.RuleCandidateID, m.ProposedSignature, m.LastAppliedSignature, m.BlockingReason, m.Enabled)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrWorkspaceConfirmationConflict
	}
	return nil
}

func (r *SQLRepository) restoreRetiredMapping(ctx context.Context, m *schema.Mapping) error {
	result, err := r.db.ExecContext(ctx, `INSERT INTO mappings(id,point_id,tag_id,transform_pipeline,status,rule_candidate_id,proposed_signature,last_applied_signature,blocking_reason,enabled,created_at,updated_at)
		SELECT ?,?,?,?,?,?,?,?,?,?,?,? WHERE NOT EXISTS (SELECT 1 FROM mappings WHERE point_id=? OR tag_id=?)`, m.ID, m.PointID, m.TagID, m.TransformPipeline, m.Status, m.RuleCandidateID, m.ProposedSignature, m.LastAppliedSignature, m.BlockingReason, m.Enabled, m.CreatedAt, m.UpdatedAt, m.PointID, m.TagID)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrWorkspaceConfirmationConflict
	}
	return nil
}

func (r *MemoryRepository) deleteRetiredMapping(_ context.Context, expected *schema.Mapping) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	current := r.mappings[expected.ID]
	if current == nil || current.PointID != expected.PointID || current.TagID != expected.TagID || current.TransformPipeline != expected.TransformPipeline || current.Status != expected.Status || current.RuleCandidateID != expected.RuleCandidateID || current.ProposedSignature != expected.ProposedSignature || current.LastAppliedSignature != expected.LastAppliedSignature || current.BlockingReason != expected.BlockingReason || current.Enabled != expected.Enabled {
		return ErrWorkspaceConfirmationConflict
	}
	delete(r.mappings, expected.ID)
	return nil
}

func (r *MemoryRepository) restoreRetiredMapping(_ context.Context, record *schema.Mapping) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, current := range r.mappings {
		if id == record.ID || current.PointID == record.PointID || current.TagID == record.TagID {
			return ErrWorkspaceConfirmationConflict
		}
	}
	r.mappings[record.ID] = cloneMapping(record)
	return nil
}
