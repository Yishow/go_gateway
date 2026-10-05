package handlers

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
)

func (h *StudioV2WorkspaceMappingsHandler) listWorkspaceRuleMappings(ctx context.Context, workspaceID, ruleID string) ([]studioV2WorkspaceMappingResponse, error) {
	ctx, release := h.ruleSvc.AcquireRuleMutation(ctx, ruleID)
	defer release()
	rule, err := h.ruleSvc.GetByID(ctx, ruleID)
	if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	links, err := h.ruleSvc.ListLinks(ctx, ruleID)
	if err != nil {
		return nil, err
	}
	payload := make([]studioV2WorkspaceMappingResponse, 0, len(links))
	linksChanged := false
	for _, link := range links {
		record, found, changed, err := h.recoverWorkspaceLinkMapping(ctx, link)
		if errors.Is(err, tag.ErrTagNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		linksChanged = linksChanged || changed
		if !found {
			continue
		}
		item, err := h.buildResponse(ctx, workspaceID, rule, link, record.ID, record.TagID)
		if err != nil {
			return nil, err
		}
		payload = append(payload, item)
	}
	if linksChanged {
		if err := h.ruleSvc.ReplaceLinks(ctx, rule.ID, links); err != nil {
			return nil, err
		}
	}
	return payload, nil
}
