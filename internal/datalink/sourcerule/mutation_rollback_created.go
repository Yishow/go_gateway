package sourcerule

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/tag"
)

func (s *Service) rollbackCreatedPoints(ctx context.Context, pointIDs []string) error {
	if err := s.storageCompensationError(ctx); err != nil {
		return err
	}

	var rollbackErrs []error
	for _, pointID := range pointIDs {
		journal := s.pointJournal(ctx)
		if journal == nil || journal.points[pointID] == nil {
			rollbackErrs = append(rollbackErrs, fmt.Errorf("created point ownership unavailable"))
			continue
		}
		entry := journal.points[pointID]
		if entry.uncertain {
			rollbackErrs = append(rollbackErrs, fmt.Errorf("created point outcome uncertain"))
			continue
		}
		if entry.current == nil {
			continue
		}
		if err := s.pointSvc.DeleteRevision(ctx, entry.current); err != nil {
			rollbackErrs = append(rollbackErrs, err)
			continue
		}
		s.suspendOwnedPoint(ctx, pointID)
		entry.current = nil
	}

	return errors.Join(rollbackErrs...)
}

func (s *Service) rollbackTagMappingSync(ctx context.Context, result tagMappingSyncResult) error {
	if err := s.storageCompensationError(ctx); err != nil {
		return err
	}

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
	if err := s.storageCompensationError(ctx); err != nil {
		return err
	}

	current, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return err
	}
	if err = s.deleteOwnedRule(ctx, current); err != nil {
		return err
	}
	var rollbackErrs []error
	if err := s.rollbackTagMappingSync(ctx, tagMappings); err != nil {
		rollbackErrs = append(rollbackErrs, err)
	}
	if err := s.rollbackCreatedPoints(ctx, pointIDs); err != nil {
		rollbackErrs = append(rollbackErrs, err)
	}

	return errors.Join(rollbackErrs...)
}
