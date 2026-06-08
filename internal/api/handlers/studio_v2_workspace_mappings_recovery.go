package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

func (h *StudioV2WorkspaceMappingsHandler) recoverWorkspaceLinkMapping(ctx context.Context, link *schema.SourceRuleLink) (*schema.Mapping, bool, bool, error) {
	return h.recoverWorkspaceLinkMappingWithPolicy(ctx, link, false)
}

func (h *StudioV2WorkspaceMappingsHandler) recoverWorkspaceLinkMappingForMutation(ctx context.Context, link *schema.SourceRuleLink) (*schema.Mapping, bool, bool, error) {
	return h.recoverWorkspaceLinkMappingWithPolicy(ctx, link, true)
}

func (h *StudioV2WorkspaceMappingsHandler) recoverWorkspaceLinkMappingWithPolicy(ctx context.Context, link *schema.SourceRuleLink, allowMissingTag bool) (*schema.Mapping, bool, bool, error) {
	if link == nil {
		return nil, false, false, nil
	}

	if link.MappingID != nil && strings.TrimSpace(*link.MappingID) != "" {
		mappingRecord, err := h.mappingSvc.GetByID(ctx, strings.TrimSpace(*link.MappingID))
		if err == nil {
			mappingRecord, matched, conflict, err := h.matchRecoveredWorkspaceLinkMapping(ctx, link, mappingRecord, allowMissingTag)
			if err != nil {
				return nil, false, false, err
			}
			if matched {
				return mappingRecord, true, syncWorkspaceLinkToMapping(link, mappingRecord), nil
			}
			if conflict {
				return nil, false, false, fmt.Errorf("persisted mapping %s does not belong to source rule link %s@%s", mappingRecord.ID, strings.TrimSpace(link.RuleID), normalizeWorkspaceAddress(link.Address))
			}
		}
		if !errors.Is(err, mapping.ErrMappingNotFound) {
			return nil, false, false, err
		}
	}

	pointID := strings.TrimSpace(link.PointID)
	if pointID == "" {
		return nil, false, false, nil
	}
	records, err := h.mappingSvc.List(ctx, mapping.ListFilter{PointID: &pointID})
	if err != nil {
		return nil, false, false, err
	}
	matches := make([]*schema.Mapping, 0, len(records))
	conflictFound := false
	for _, record := range records {
		candidate, matched, conflict, matchErr := h.matchRecoveredWorkspaceLinkMapping(ctx, link, record, allowMissingTag)
		if matchErr != nil {
			return nil, false, false, matchErr
		}
		if matched {
			matches = append(matches, candidate)
			continue
		}
		if conflict {
			conflictFound = true
		}
	}
	if len(matches) == 0 {
		if conflictFound {
			return nil, false, false, fmt.Errorf("persisted mappings for point %s belong to another source rule link", pointID)
		}
		return nil, false, false, nil
	}
	if len(matches) > 1 {
		return nil, false, false, fmt.Errorf("recover workspace mapping for point %s: multiple persisted mappings", pointID)
	}
	mappingRecord := matches[0]
	return mappingRecord, true, syncWorkspaceLinkToMapping(link, mappingRecord), nil
}

func (h *StudioV2WorkspaceMappingsHandler) matchRecoveredWorkspaceLinkMapping(ctx context.Context, link *schema.SourceRuleLink, mappingRecord *schema.Mapping, allowMissingTag bool) (*schema.Mapping, bool, bool, error) {
	if mappingRecord == nil {
		return nil, false, false, nil
	}
	if strings.TrimSpace(mappingRecord.PointID) != strings.TrimSpace(link.PointID) {
		return nil, false, true, nil
	}
	tagRecord, err := h.tagSvc.GetByID(ctx, mappingRecord.TagID)
	if err != nil {
		if allowMissingTag && errors.Is(err, tag.ErrTagNotFound) {
			return mappingRecord, true, false, nil
		}
		if errors.Is(err, tag.ErrTagNotFound) {
			return nil, false, false, nil
		}
		return nil, false, false, fmt.Errorf("load recovered workspace tag %s: %w", mappingRecord.TagID, err)
	}
	if !isRuleManagedWorkspaceTagOwnedBy(tagRecord, link.RuleID, link.Address) && !isLegacyDirectWorkspaceBinding(link, mappingRecord, tagRecord) {
		return nil, false, true, nil
	}
	return mappingRecord, true, false, nil
}

func syncWorkspaceLinkToMapping(link *schema.SourceRuleLink, mappingRecord *schema.Mapping) bool {
	changed := false
	if link.MappingID == nil || *link.MappingID != mappingRecord.ID {
		link.MappingID = cloneStringPtr(mappingRecord.ID)
		changed = true
	}
	if link.TagID == nil || *link.TagID != mappingRecord.TagID {
		link.TagID = cloneStringPtr(mappingRecord.TagID)
		changed = true
	}
	if changed {
		link.UpdatedAt = time.Now()
	}
	return changed
}

func (h *StudioV2WorkspaceMappingsHandler) saveExistingWorkspaceMapping(ctx context.Context, rule *schema.SourceRule, links []*schema.SourceRuleLink, link *schema.SourceRuleLink, mappingRecord *schema.Mapping, req studioV2WorkspaceMappingRequest) (*schema.Mapping, error) {
	pointRecord, err := h.pointSvc.GetByID(ctx, link.PointID)
	if err != nil {
		return nil, err
	}
	oldTagID := mappingRecord.TagID
	tagRecord, err := h.upsertTag(ctx, rule.ID, req, link.Address)
	if err != nil {
		return nil, err
	}
	mappingRecord, err = h.mappingSvc.Update(ctx, mappingRecord.ID, mapping.UpdateMappingRequest{
		TagID:             cloneStringPtr(tagRecord.ID),
		Enabled:           boolPtr(req.Enabled),
		TransformPipeline: buildWorkspaceMappingPipeline(pointRecord.DataType, req.TargetType, req.Scale, req.Offset),
	})
	if err != nil {
		return nil, err
	}

	syncWorkspaceLinkToMapping(link, mappingRecord)
	if err := h.ruleSvc.ReplaceLinks(ctx, rule.ID, links); err != nil {
		return nil, err
	}
	if oldTagID != tagRecord.ID {
		if err := h.deleteOrphanRuleManagedTag(ctx, oldTagID); err != nil {
			return nil, err
		}
	}
	return mappingRecord, nil
}
