package handlers

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type workspaceMappingRow struct {
	record  *workspace.Record
	rule    *schema.SourceRule
	links   []*schema.SourceRuleLink
	link    *schema.SourceRuleLink
	mapping *schema.Mapping
	request studioV2WorkspaceMappingRequest
	release func()
}

func (h *StudioV2WorkspaceMappingsHandler) resolveRequestRow(c *gin.Context) (workspaceMappingRow, bool) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return workspaceMappingRow{}, false
	}
	var req studioV2WorkspaceMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return workspaceMappingRow{}, false
	}
	if err := validateWorkspaceMappingRequest(req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return workspaceMappingRow{}, false
	}
	ctx, release := h.ruleSvc.AcquireRuleMutation(c.Request.Context(), req.RuleID)
	c.Request = c.Request.WithContext(ctx)
	resolved := false
	defer func() {
		if !resolved {
			release()
		}
	}()
	rule, err := h.ruleSvc.GetByID(c.Request.Context(), req.RuleID)
	if err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return workspaceMappingRow{}, false
	}
	if !workspaceOwnsDevice(record, rule.DeviceID) {
		renderStudioV2WorkspaceMappingError(c, mapping.ErrMappingNotFound)
		return workspaceMappingRow{}, false
	}
	links, err := h.ruleSvc.ListLinks(c.Request.Context(), rule.ID)
	if err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return workspaceMappingRow{}, false
	}
	for _, link := range links {
		if normalizeWorkspaceAddress(link.Address) == normalizeWorkspaceAddress(req.Address) {
			resolved = true
			return workspaceMappingRow{record: record, rule: rule, links: links, link: link, request: req, release: release}, true
		}
	}
	renderStudioV2WorkspaceMappingError(c, mapping.ErrMappingNotFound)
	return workspaceMappingRow{}, false
}

func (h *StudioV2WorkspaceMappingsHandler) requireWorkspaceMapping(c *gin.Context) (workspaceMappingRow, bool) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return workspaceMappingRow{}, false
	}
	var req studioV2WorkspaceMappingRequest
	if c.Request.ContentLength != 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			renderStudioV2WorkspaceValidationError(c, err)
			return workspaceMappingRow{}, false
		}
		if err := validateWorkspaceMappingRequest(req); err != nil {
			renderStudioV2WorkspaceValidationError(c, err)
			return workspaceMappingRow{}, false
		}
	}

	ruleID := ""
	{
		// Locate ownership using atomic link reads before waiting for any guard.
		// Do not serialize a request behind unrelated workspace rules.
		// Derive actual ownership even when a mismatched rule_id was supplied.
		recordMapping, err := h.mappingSvc.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			renderStudioV2WorkspaceMappingError(c, err)
			return workspaceMappingRow{}, false
		}
		rules, err := h.ruleSvc.ListByDeviceIDs(c.Request.Context(), record.OrderedDeviceIDs)
		if err != nil {
			renderStudioV2WorkspaceMappingError(c, err)
			return workspaceMappingRow{}, false
		}
		for _, rule := range rules {
			links, err := h.ruleSvc.ListLinks(c.Request.Context(), rule.ID)
			if err != nil {
				renderStudioV2WorkspaceMappingError(c, err)
				return workspaceMappingRow{}, false
			}
			for _, link := range links {
				if link.PointID == recordMapping.PointID {
					ruleID = rule.ID
					break
				}
			}
			if ruleID != "" {
				break
			}
		}
	}
	if ruleID != "" {
		ctx, release := h.ruleSvc.AcquireRuleMutation(c.Request.Context(), ruleID)
		row, found, err := h.findWorkspaceMappingRow(ctx, record, ruleID, c.Param("id"), req)
		if err != nil {
			release()
			renderStudioV2WorkspaceMappingError(c, err)
			return workspaceMappingRow{}, false
		}
		if found {
			row.release = release
			c.Request = c.Request.WithContext(ctx)
			return row, true
		}
		release()
	}

	renderStudioV2WorkspaceMappingError(c, mapping.ErrMappingNotFound)
	return workspaceMappingRow{}, false
}

func (h *StudioV2WorkspaceMappingsHandler) findWorkspaceMappingRow(ctx context.Context, record *workspace.Record, ruleID, mappingID string, req studioV2WorkspaceMappingRequest) (workspaceMappingRow, bool, error) {
	rule, err := h.ruleSvc.GetByID(ctx, ruleID)
	if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
		return workspaceMappingRow{}, false, nil
	}
	if err != nil {
		return workspaceMappingRow{}, false, err
	}
	if !workspaceOwnsDevice(record, rule.DeviceID) {
		return workspaceMappingRow{}, false, nil
	}
	links, err := h.ruleSvc.ListLinks(ctx, rule.ID)
	if err != nil {
		return workspaceMappingRow{}, false, err
	}
	for _, link := range links {
		mappingRecord, found, changed, err := h.recoverWorkspaceLinkMappingForMutation(ctx, link)
		if err != nil {
			return workspaceMappingRow{}, false, err
		}
		if found && mappingRecord.ID == mappingID {
			if changed {
				if err := h.ruleSvc.ReplaceLinks(ctx, rule.ID, links); err != nil {
					return workspaceMappingRow{}, false, err
				}
			}
			return workspaceMappingRow{record: record, rule: rule, links: links, link: link, mapping: mappingRecord, request: req}, true, nil
		}
	}
	return workspaceMappingRow{}, false, nil
}
