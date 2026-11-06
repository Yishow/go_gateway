package sourcerule

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

// A removed mixed item may occupy the destination of a retained or new item.
// Keep receipts for the exact resources removed before moving the other points.
type retiredRuleLink struct {
	point          *schema.Point
	tag            *schema.Tag
	mapping        *schema.Mapping
	tagID          string
	mappingID      string
	pointDeleted   bool
	tagDeleted     bool
	mappingDeleted bool
}

func (s *Service) retiredRuleLink(ctx context.Context, linkID string) *retiredRuleLink {
	journal := s.pointJournal(ctx)
	if journal == nil {
		return nil
	}
	return journal.retiredLinks[linkID]
}

func (s *Service) retireCollidingRuleLink(ctx context.Context, link *schema.SourceRuleLink) error {
	journal := s.pointJournal(ctx)
	if journal == nil {
		return fmt.Errorf("retired point ownership unavailable")
	}
	retired := &retiredRuleLink{}
	if link.TagID != nil {
		retired.tagID = *link.TagID
	}
	if link.MappingID != nil {
		retired.mappingID = *link.MappingID
	}
	if baseline := journal.retirementBaseline; baseline != nil {
		retired.point = clonePointSnapshot(baseline.points[link.PointID])
		retired.tag = cloneTagForRollback(baseline.tags[retired.tagID])
		retired.mapping = cloneMappingForRollback(baseline.mappings[retired.mappingID])
	} else {
		var err error
		retired.point, err = s.pointSvc.GetByID(ctx, link.PointID)
		if err != nil {
			return err
		}
		if s.tagSvc != nil && retired.tagID != "" {
			retired.tag, err = s.tagSvc.GetByID(ctx, retired.tagID)
			if err != nil && !errors.Is(err, tag.ErrTagNotFound) {
				return err
			}
		}
		if s.mappingSvc != nil && retired.mappingID != "" {
			retired.mapping, err = s.mappingSvc.GetByID(ctx, retired.mappingID)
			if err != nil && !errors.Is(err, mapping.ErrMappingNotFound) {
				return err
			}
		}
	}
	if retired.point == nil || (retired.mappingID != "" && s.mappingSvc == nil) {
		return point.ErrFormatConflict
	}
	if retired.mapping != nil && (retired.mapping.PointID != link.PointID || (retired.tagID != "" && retired.mapping.TagID != retired.tagID)) {
		return point.ErrFormatConflict
	}
	if s.mappingSvc != nil {
		maps, listErr := s.mappingSvc.List(ctx, mapping.ListFilter{PointID: &link.PointID})
		if listErr != nil {
			return listErr
		}
		if len(maps) > 1 || (len(maps) == 1 && (retired.mapping == nil || maps[0].ID != retired.mapping.ID)) {
			return point.ErrFormatConflict
		}
	}
	if journal.retiredLinks == nil {
		journal.retiredLinks = map[string]*retiredRuleLink{}
	}
	journal.retiredLinks[link.ID] = retired
	if retired.mapping != nil {
		if err := s.mappingSvc.DeleteRetiredMapping(ctx, retired.mapping); err != nil {
			if errors.Is(err, mapping.ErrWorkspaceConfirmationConflict) {
				return errors.Join(point.ErrFormatConflict, err)
			}
			return s.dirtyUnknown(ctx)
		}
		retired.mappingDeleted = true
	}
	if err := s.pointSvc.DeleteUnmappedMaterializedObserved(ctx, retired.point, s.observePointWrite(ctx)); err != nil {
		return err
	}
	retired.pointDeleted = true
	if s.runtimeSync != nil && !s.removeOwnedPointRuntime(ctx, retired.point.ID) {
		return s.dirtyUnknown(ctx)
	}
	return nil
}

func (s *Service) cleanupRetiredRuleTag(ctx context.Context, rule *schema.SourceRule, link *schema.SourceRuleLink) error {
	retired := s.retiredRuleLink(ctx, link.ID)
	if retired == nil || retired.tag == nil || s.tagSvc == nil || s.mappingSvc == nil || !isRuleManagedTag(retired.tag, rule.ID, link.Address) {
		return nil
	}
	maps, err := s.mappingSvc.List(ctx, mapping.ListFilter{TagID: &retired.tag.ID})
	if err != nil || len(maps) > 0 {
		return err
	}
	retired.tagDeleted, err = s.tagSvc.DeleteRetiredTag(ctx, retired.tag)
	if err != nil {
		return s.dirtyUnknown(ctx)
	}
	return nil
}

// Restore points first, then tags and mappings, so restored links satisfy FKs.
// Every insert is conditional on the identity remaining absent and unbound.
func (s *Service) restoreRetiredRuleResources(ctx context.Context) error {
	journal := s.pointJournal(ctx)
	if journal == nil || len(journal.retiredLinks) == 0 {
		return nil
	}
	for _, retired := range journal.retiredLinks {
		if !retired.pointDeleted {
			continue
		}
		entry := journal.points[retired.point.ID]
		if entry == nil {
			return fmt.Errorf("retired point write ownership unavailable")
		}
		if entry.current == nil {
			if err := s.pointSvc.RestoreAbsent(ctx, retired.point); err != nil {
				return err
			}
			s.notePointRestored(ctx, retired.point)
		}
	}
	if journal.restoredIdentityTags == nil {
		journal.restoredIdentityTags = map[string]bool{}
		journal.restoredIdentityMaps = map[string]bool{}
	}
	for _, retired := range journal.retiredLinks {
		if retired.tag != nil {
			if retired.tagDeleted {
				if err := s.tagSvc.RestoreRetiredTag(ctx, retired.tag); err != nil {
					return err
				}
				retired.tagDeleted = false
			}
		}
		if retired.tagID != "" {
			journal.restoredIdentityTags[retired.tagID] = true
		}
	}
	for _, retired := range journal.retiredLinks {
		if retired.mapping != nil {
			if retired.mappingDeleted {
				if err := s.mappingSvc.RestoreRetiredMapping(ctx, retired.mapping); err != nil {
					return err
				}
				retired.mappingDeleted = false
			}
		}
		if retired.mappingID != "" {
			journal.restoredIdentityMaps[retired.mappingID] = true
		}
	}
	return nil
}

func (s *Service) restoreRetiredRuleRuntime(ctx context.Context) error {
	journal := s.pointJournal(ctx)
	if journal == nil {
		return nil
	}
	var points []*schema.Point
	for _, retired := range journal.retiredLinks {
		if retired.pointDeleted {
			points = append(points, retired.point)
		}
	}
	if len(points) == 0 {
		return nil
	}
	return s.restoreOwnedRuntime(ctx, points)
}
