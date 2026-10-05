package mapping

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/schema"
	"time"
)

var ErrWorkspaceConfirmationConflict = errors.New("workspace mapping needs source candidate review or changed during save")

type workspaceConfirmationRepository interface {
	ConfirmWorkspace(context.Context, *schema.Mapping, *schema.Mapping, string, string, schema.DataType, func() error) error
}

// ConfirmWorkspace saves confirmation metadata only for the exact successfully
// saved pipeline. The SQL guard also checks the current source revision and link.
func (s *Service) ConfirmWorkspace(ctx context.Context, saved *schema.Mapping, ruleID, revision, candidateID, proposed, applied string, targetType schema.DataType, enabled bool, guard func() error) error {
	repo, ok := s.repo.(workspaceConfirmationRepository)
	if !ok {
		return ErrWorkspaceConfirmationConflict
	}
	confirmed := *saved
	confirmed.Enabled = enabled
	confirmed.Status = schema.MappingStatusDraft
	if enabled {
		confirmed.Status = schema.MappingStatusActive
	}
	confirmed.RuleCandidateID = candidateID
	confirmed.ProposedSignature = proposed
	confirmed.LastAppliedSignature = applied
	confirmed.BlockingReason = ""
	confirmed.UpdatedAt = time.Now()
	return repo.ConfirmWorkspace(ctx, saved, &confirmed, ruleID, revision, targetType, guard)
}

func (r *SQLRepository) ConfirmWorkspace(ctx context.Context, saved, confirmed *schema.Mapping, ruleID, revision string, targetType schema.DataType, _ func() error) error {
	result, err := r.db.ExecContext(ctx, `UPDATE mappings SET rule_candidate_id=?,proposed_signature=?,last_applied_signature=?,blocking_reason='',enabled=?,status=?,updated_at=? WHERE id=? AND point_id=? AND tag_id=? AND transform_pipeline=? AND status=? AND enabled=? AND EXISTS(SELECT 1 FROM source_rules r JOIN source_rule_links l ON l.rule_id=r.id WHERE r.id=? AND r.revision_id=? AND l.point_id=mappings.point_id AND l.mapping_id=mappings.id AND l.tag_id=mappings.tag_id) AND EXISTS(SELECT 1 FROM tags WHERE id=mappings.tag_id AND data_type=?)`, confirmed.RuleCandidateID, confirmed.ProposedSignature, confirmed.LastAppliedSignature, confirmed.Enabled, confirmed.Status, confirmed.UpdatedAt, saved.ID, saved.PointID, saved.TagID, saved.TransformPipeline, saved.Status, saved.Enabled, ruleID, revision, string(targetType))
	if err != nil {
		return err
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

func (r *MemoryRepository) ConfirmWorkspace(ctx context.Context, saved, confirmed *schema.Mapping, _, _ string, _ schema.DataType, guard func() error) error {
	if err := guard(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	current := r.mappings[saved.ID]
	if current == nil || current.PointID != saved.PointID || current.TagID != saved.TagID || current.TransformPipeline != saved.TransformPipeline || current.Status != saved.Status || current.Enabled != saved.Enabled {
		return ErrWorkspaceConfirmationConflict
	}
	record := *confirmed
	r.mappings[saved.ID] = &record
	return nil
}
