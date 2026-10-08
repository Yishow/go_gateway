package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

const workspaceMappingCleanupFailedStatus = "failed"

func (h *StudioV2WorkspaceMappingsHandler) saveExistingWorkspaceMapping(ctx context.Context, rule *schema.SourceRule, links []*schema.SourceRuleLink, link *schema.SourceRuleLink, record *schema.Mapping, req studioV2WorkspaceMappingRequest) (saved *schema.Mapping, cleanupStatus string, saveErr error) {
	if err := h.ruleSvc.ValidateWorkspaceMappingSave(ctx, rule, record); err != nil {
		return nil, "", err
	}
	pointRecord, err := h.pointSvc.GetByID(ctx, link.PointID)
	if err != nil {
		return nil, "", err
	}
	previousMapping := *record
	previousLinks := cloneWorkspaceLinks(links)
	previousTag, err := h.tagSvc.GetByKey(ctx, tag.NormalizeTagKey(req.TagKey))
	if err != nil && !errors.Is(err, tag.ErrTagNotFound) {
		return nil, "", err
	}
	var lastMapping *schema.Mapping
	var lastTag *schema.Tag
	linksReplaced := false
	defer func() {
		if saveErr == nil || lastMapping == nil {
			return
		}
		// Request cancellation must not strand an unconfirmed pipeline. Compensation
		// remains bounded and protected by the caller's synchronous rule guard.
		rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		restoreErr := h.mappingSvc.RestoreWorkspaceSave(rollbackCtx, lastMapping, &previousMapping)
		if restoreErr != nil {
			saveErr = errors.Join(saveErr, fmt.Errorf("restore workspace mapping: %w", restoreErr))
			return
		}
		var rollbackErrs []error
		if linksReplaced {
			rollbackErrs = append(rollbackErrs, h.ruleSvc.ReplaceLinks(rollbackCtx, rule.ID, previousLinks))
		}
		if lastTag != nil {
			if previousTag != nil {
				rollbackErrs = append(rollbackErrs, h.tagSvc.RestoreWorkspaceSave(rollbackCtx, lastTag, previousTag))
			} else {
				rollbackErrs = append(rollbackErrs, h.deleteOrphanRuleManagedTag(rollbackCtx, lastTag.ID))
			}
		}
		saveErr = errors.Join(saveErr, errors.Join(rollbackErrs...))
	}()
	draftStatus := schema.MappingStatusDraft
	lastMapping, err = h.mappingSvc.Update(ctx, record.ID, mapping.UpdateMappingRequest{Enabled: common.Ptr(false), Status: &draftStatus})
	if err != nil {
		return nil, "", err
	}
	lastTag, err = h.upsertTag(ctx, rule.ID, req, link.Address)
	if err != nil {
		return nil, "", err
	}
	updated, err := h.mappingSvc.Update(ctx, record.ID, mapping.UpdateMappingRequest{TagID: common.Ptr(lastTag.ID), Enabled: common.Ptr(false), TransformPipeline: buildWorkspaceMappingPipeline(pointRecord.DataType, req.TargetType, req.Scale, req.Offset)})
	if err != nil {
		return nil, "", err
	}
	lastMapping = updated
	syncWorkspaceLinkToMapping(link, lastMapping)
	if err := h.ruleSvc.ReplaceLinks(ctx, rule.ID, links); err != nil {
		return nil, "", err
	}
	linksReplaced = true
	if err := h.ruleSvc.ConfirmWorkspaceMapping(ctx, rule.ID, rule.RevisionID, lastMapping, req.TargetType); err != nil {
		return nil, "", err
	}
	lastMapping = nil
	// Confirmation has succeeded. Orphan cleanup cannot roll back an accepted save.
	if previousMapping.TagID != lastTag.ID {
		if err := h.deleteOrphanRuleManagedTag(ctx, previousMapping.TagID); err != nil {
			return updated, workspaceMappingCleanupFailedStatus, nil
		}
	}
	return updated, "", nil
}

func cloneWorkspaceLinks(links []*schema.SourceRuleLink) []*schema.SourceRuleLink {
	result := make([]*schema.SourceRuleLink, 0, len(links))
	for _, link := range links {
		clonedLink := *link
		if link.TagID != nil {
			clonedLink.TagID = common.Ptr(*link.TagID)
		}
		if link.MappingID != nil {
			clonedLink.MappingID = common.Ptr(*link.MappingID)
		}
		result = append(result, &clonedLink)
	}
	return result
}
