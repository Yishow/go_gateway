package sourcerule

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/tag"
)

func (s *Service) rollbackCreatedPoints(ctx context.Context, pointIDs []string) error {
	var rollbackErrs []error
	for _, pointID := range pointIDs {
		if err := s.pointSvc.Delete(ctx, pointID); err != nil {
			rollbackErrs = append(rollbackErrs, fmt.Errorf("delete created point %s: %w", pointID, err))
		}
		if s.runtimeSync != nil {
			s.runtimeSync.RemovePoint(pointID)
		}
	}
	return errors.Join(rollbackErrs...)
}

func (s *Service) rollbackTagMappingSync(ctx context.Context, result tagMappingSyncResult) error {
	var rollbackErrs []error
	if s.mappingSvc != nil {
		for mappingID, state := range result.updatedMappings {
			enabled := state.enabled
			if _, err := s.mappingSvc.Update(ctx, mappingID, mapping.UpdateMappingRequest{Enabled: &enabled, TransformPipeline: state.transformPipeline, Status: &state.status, RuleCandidateID: &state.ruleCandidateID, ProposedSignature: &state.proposedSignature, LastAppliedSignature: &state.lastAppliedSignature, BlockingReason: &state.blockingReason}); err != nil {
				rollbackErrs = append(rollbackErrs, fmt.Errorf("restore mapping %s: %w", mappingID, err))
			}
		}
		for _, mappingID := range result.createdMappingIDs {
			if err := s.mappingSvc.Delete(ctx, mappingID); err != nil {
				rollbackErrs = append(rollbackErrs, fmt.Errorf("delete created mapping %s: %w", mappingID, err))
			}
		}
	}
	if s.tagSvc != nil {
		for tagID, state := range result.updatedTags {
			dataType := state.dataType
			if _, err := s.tagSvc.Update(ctx, tagID, tag.UpdateTagRequest{DataType: &dataType}); err != nil {
				rollbackErrs = append(rollbackErrs, fmt.Errorf("restore tag %s: %w", tagID, err))
			}
		}
		for _, tagID := range result.createdTagIDs {
			if err := s.tagSvc.Delete(ctx, tagID); err != nil {
				rollbackErrs = append(rollbackErrs, fmt.Errorf("delete created tag %s: %w", tagID, err))
			}
		}
	}
	return errors.Join(rollbackErrs...)
}

func (s *Service) rollbackCreatedRule(ctx context.Context, ruleID, revisionID string, pointIDs []string, tagMappings tagMappingSyncResult) error {
	var rollbackErrs []error
	if err := s.rollbackTagMappingSync(ctx, tagMappings); err != nil {
		rollbackErrs = append(rollbackErrs, err)
	}
	if err := s.rollbackCreatedPoints(ctx, pointIDs); err != nil {
		rollbackErrs = append(rollbackErrs, err)
	}
	if err := s.repo.DeleteLinks(ctx, ruleID); err != nil {
		rollbackErrs = append(rollbackErrs, fmt.Errorf("delete created source-rule links %s: %w", ruleID, err))
	}
	if err := s.repo.DeleteCandidateSnapshots(ctx, ruleID, revisionID); err != nil {
		rollbackErrs = append(rollbackErrs, fmt.Errorf("delete created source-rule candidate snapshots %s: %w", ruleID, err))
	}
	if err := s.repo.Delete(ctx, ruleID); err != nil {
		rollbackErrs = append(rollbackErrs, fmt.Errorf("delete created source rule %s: %w", ruleID, err))
	}
	return errors.Join(rollbackErrs...)
}
