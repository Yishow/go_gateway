package sourcerule

import (
	"context"
	"fmt"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
)

// ValidateWorkspaceMappingSave refuses stale source proposals or unconfirmed
// edits before the workspace route mutates the tag/pipeline.
func (s *Service) ValidateWorkspaceMappingSave(ctx context.Context, rule *schema.SourceRule, record *schema.Mapping) error {
	ctx, release := s.AcquireRuleMutation(ctx, rule.ID)
	defer release()
	current, err := s.repo.GetByID(ctx, rule.ID)
	if err != nil {
		return err
	}
	if current.RevisionID != rule.RevisionID || record.Status == schema.MappingStatusOutOfSync {
		return mapping.ErrWorkspaceConfirmationConflict
	}
	if record.RuleCandidateID == "" {
		return nil
	}
	point, err := s.pointSvc.GetByID(ctx, record.PointID)
	if err != nil {
		return err
	}
	proposed, err := mappingCandidateSignature(s.buildRuleTransformPipeline(current, point))
	if err != nil {
		return err
	}
	steps, err := decodeTransformPipeline(record.TransformPipeline)
	if err != nil {
		return err
	}
	actual, err := mappingCandidateSignature(steps)
	if err != nil {
		return err
	}
	if record.ProposedSignature != proposed || record.LastAppliedSignature != actual {
		return mapping.ErrWorkspaceConfirmationConflict
	}
	return nil
}

// ConfirmWorkspaceMapping is used only after an explicit owned workspace save.
// It confirms one mapping; ordinary derived sync cannot confirm any other edit.
func (s *Service) ConfirmWorkspaceMapping(ctx context.Context, ruleID, expectedRevision string, saved *schema.Mapping, targetType schema.DataType) error {
	ctx, release := s.AcquireRuleMutation(ctx, ruleID)
	defer release()
	if s.mappingSvc == nil || s.tagSvc == nil || saved == nil {
		return fmt.Errorf("workspace mapping confirmation services unavailable")
	}
	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return err
	}
	if rule.RevisionID != expectedRevision || saved.Status == schema.MappingStatusOutOfSync {
		return mapping.ErrWorkspaceConfirmationConflict
	}
	links, err := s.repo.ListLinks(ctx, ruleID)
	if err != nil {
		return err
	}
	var owned *schema.SourceRuleLink
	for _, link := range links {
		if link.MappingID != nil && *link.MappingID == saved.ID && link.PointID == saved.PointID && link.TagID != nil && *link.TagID == saved.TagID {
			owned = link
			break
		}
	}
	if owned == nil {
		return mapping.ErrWorkspaceConfirmationConflict
	}
	point, err := s.pointSvc.GetByID(ctx, saved.PointID)
	if err != nil {
		return err
	}
	candidateID, proposed, err := ruleManagedMappingMetadata(rule, point, owned, s.buildRuleTransformPipeline(rule, point))
	if err != nil {
		return err
	}
	appliedSteps, err := decodeTransformPipeline(saved.TransformPipeline)
	if err != nil {
		return err
	}
	applied, err := mappingCandidateSignature(appliedSteps)
	if err != nil {
		return err
	}
	guard := func() error {
		current, err := s.repo.GetByID(ctx, ruleID)
		if err != nil {
			return err
		}
		if current.RevisionID != expectedRevision {
			return mapping.ErrWorkspaceConfirmationConflict
		}
		return nil
	}
	// Settle other links first; any failure leaves the selected accepted hash unchanged.
	if err := s.syncDerivedRuleState(ctx, rule, saved.ID); err != nil {
		return err
	}
	enabled, err := s.runtimeRuleEnabled(ctx, rule)
	if err != nil {
		return err
	}
	return s.mappingSvc.ConfirmWorkspace(ctx, saved, ruleID, expectedRevision, candidateID, proposed, applied, targetType, enabled, guard)
}

func confirmedWorkspacePipeline(record *schema.Mapping, proposed, actual string) bool {
	return record.RuleCandidateID != "" && record.ProposedSignature == proposed && record.LastAppliedSignature == actual &&
		(record.Status == schema.MappingStatusActive || record.Status == schema.MappingStatusDraft)
}
