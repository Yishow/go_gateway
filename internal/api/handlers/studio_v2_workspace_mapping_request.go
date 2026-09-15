package handlers

import (
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"net/http"

	"github.com/gin-gonic/gin"
)

type workspaceMappingRow struct {
	record  *workspace.Record
	rule    *schema.SourceRule
	links   []*schema.SourceRuleLink
	link    *schema.SourceRuleLink
	mapping *schema.Mapping
	request studioV2WorkspaceMappingRequest
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
	rule, err := h.ruleSvc.GetByID(c.Request.Context(), req.RuleID)
	if err != nil || !workspaceOwnsDevice(record, rule.DeviceID) {
		c.JSON(http.StatusNotFound, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: studioV2MappingRowNotFoundMessage}})
		return workspaceMappingRow{}, false
	}
	links, err := h.ruleSvc.ListLinks(c.Request.Context(), rule.ID)
	if err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return workspaceMappingRow{}, false
	}
	for _, link := range links {
		if normalizeWorkspaceAddress(link.Address) == normalizeWorkspaceAddress(req.Address) {
			return workspaceMappingRow{record: record, rule: rule, links: links, link: link, request: req}, true
		}
	}
	c.JSON(http.StatusNotFound, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: studioV2MappingRowNotFoundMessage}})
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
			mappingRecord, found, changed, err := h.recoverWorkspaceLinkMappingForMutation(c.Request.Context(), link)
			if err != nil {
				renderStudioV2WorkspaceMappingError(c, err)
				return workspaceMappingRow{}, false
			}
			if found && mappingRecord.ID == c.Param("id") {
				if changed {
					if err := h.ruleSvc.ReplaceLinks(c.Request.Context(), rule.ID, links); err != nil {
						renderStudioV2WorkspaceMappingError(c, err)
						return workspaceMappingRow{}, false
					}
				}
				return workspaceMappingRow{record: record, rule: rule, links: links, link: link, mapping: mappingRecord, request: req}, true
			}
		}
	}
	c.JSON(http.StatusNotFound, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: studioV2MappingRowNotFoundMessage}})
	return workspaceMappingRow{}, false
}
