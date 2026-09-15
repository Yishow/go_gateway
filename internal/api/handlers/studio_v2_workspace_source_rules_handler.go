package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type StudioV2WorkspaceSourceRulesHandler struct {
	workspaceSvc *workspace.Service
	deviceSvc    *device.Service
	ruleSvc      *sourcerule.Service
}

type studioV2WorkspaceUpdateRuleRequest struct {
	DeviceID *string `json:"device_id,omitempty"`
	sourcerule.UpdateRuleRequest
}

func NewStudioV2WorkspaceSourceRulesHandler(workspaceSvc *workspace.Service, deviceSvc *device.Service, ruleSvc *sourcerule.Service) *StudioV2WorkspaceSourceRulesHandler {
	return &StudioV2WorkspaceSourceRulesHandler{
		workspaceSvc: workspaceSvc,
		deviceSvc:    deviceSvc,
		ruleSvc:      ruleSvc,
	}
}

func (h *StudioV2WorkspaceSourceRulesHandler) List(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	rules, err := h.ruleSvc.ListByDeviceIDs(c.Request.Context(), record.OrderedDeviceIDs)
	if err != nil {
		renderStudioV2WorkspaceSourceRuleError(c, err)
		return
	}

	payload := make([]sourceRuleResponse, 0, len(rules))
	for _, rule := range rules {
		item := mapSourceRuleResponse(rule)
		item.WorkspaceID = record.ID
		item.SaveState = workspaceSaveStateSaved
		payload = append(payload, item)
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    payload,
	})
}

func (h *StudioV2WorkspaceSourceRulesHandler) Create(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	var req sourcerule.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}
	if !workspaceOwnsDevice(record, req.DeviceID) {
		renderStudioV2WorkspaceValidationError(c, fmt.Errorf("device_id does not belong to workspace: device_id=%s", strings.TrimSpace(req.DeviceID)))
		return
	}

	rule, err := h.ruleSvc.Create(c.Request.Context(), req)
	if err != nil {
		renderStudioV2WorkspaceSourceRuleError(c, err)
		return
	}

	payload := mapSourceRuleResponse(rule)
	payload.WorkspaceID = record.ID
	payload.SaveState = workspaceSaveStateSaved
	links, err := h.ruleSvc.ListLinks(c.Request.Context(), rule.ID)
	if err != nil {
		renderStudioV2WorkspaceSourceRuleError(c, err)
		return
	}
	applyOutcome := resolveStudioV2ScopedRuntimeApplyOutcome(c.Request.Context(), h.workspaceSvc, h.deviceSvc, []string{rule.DeviceID}, workspaceRuleApplyScopes(rule.DeviceID, links))
	payload.RuntimeApplyStatus = applyOutcome.Status
	payload.RuntimeApplyMessage = applyOutcome.Message
	payload.RuntimeApplyIssues = applyOutcome.Issues
	c.JSON(http.StatusCreated, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}

func (h *StudioV2WorkspaceSourceRulesHandler) Update(c *gin.Context) {
	record, rule, ok := h.requireWorkspaceRule(c)
	if !ok {
		return
	}

	req, err := parseStudioV2WorkspaceUpdateRuleRequest(c)
	if err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}
	if req.DeviceID != nil && *req.DeviceID != rule.DeviceID {
		renderStudioV2WorkspaceValidationError(c, fmt.Errorf("workspace rule ownership mismatch: device_id=%s", strings.TrimSpace(*req.DeviceID)))
		return
	}

	updatedRule, reconcileOutcome, err := h.ruleSvc.UpdateWithRuntimeReconcile(c.Request.Context(), c.Param("id"), req.UpdateRuleRequest)
	if err != nil {
		renderStudioV2WorkspaceSourceRuleError(c, err)
		return
	}

	payload := mapSourceRuleResponse(updatedRule)
	payload.WorkspaceID = record.ID
	payload.SaveState = workspaceSaveStateSaved
	applyOutcome := mapSourceRuleRuntimeReconcileOutcome(reconcileOutcome)
	payload.RuntimeApplyStatus = applyOutcome.Status
	payload.RuntimeApplyMessage = applyOutcome.Message
	payload.RuntimeApplyIssues = applyOutcome.Issues
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}

func (h *StudioV2WorkspaceSourceRulesHandler) Delete(c *gin.Context) {
	_, rule, ok := h.requireWorkspaceRule(c)
	if !ok {
		return
	}

	reconcileOutcome, err := h.ruleSvc.DeleteWithRuntimeReconcile(c.Request.Context(), rule.ID)
	if err != nil {
		renderStudioV2WorkspaceSourceRuleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    mapStudioV2RuntimeApplyResponse(mapSourceRuleRuntimeReconcileOutcome(reconcileOutcome)),
	})
}

func (h *StudioV2WorkspaceSourceRulesHandler) requireWorkspaceRule(c *gin.Context) (*workspace.Record, *schema.SourceRule, bool) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return nil, nil, false
	}

	rule, err := h.ruleSvc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: studioV2SourceRuleNotFoundMessage},
		})
		return nil, nil, false
	}
	if !workspaceOwnsDevice(record, rule.DeviceID) {
		c.JSON(http.StatusNotFound, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: studioV2SourceRuleNotFoundMessage},
		})
		return nil, nil, false
	}

	return record, rule, true
}

func parseStudioV2WorkspaceUpdateRuleRequest(c *gin.Context) (studioV2WorkspaceUpdateRuleRequest, error) {
	body, err := c.GetRawData()
	if err != nil {
		return studioV2WorkspaceUpdateRuleRequest{}, err
	}

	var req studioV2WorkspaceUpdateRuleRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return studioV2WorkspaceUpdateRuleRequest{}, err
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return studioV2WorkspaceUpdateRuleRequest{}, err
	}
	_, req.TargetDataTypeSet = raw["target_data_type"]
	_, req.ScaleMultiplierSet = raw["scale_multiplier"]
	_, req.ScaleOffsetSet = raw["scale_offset"]
	_, req.DataFormatSet = raw["data_format"]
	_, req.ShareStartRegisterSet = raw["share_start_register"]
	_, req.ShareStrideSet = raw["share_stride"]

	return req, nil
}

func renderStudioV2WorkspaceSourceRuleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, sourcerule.ErrValidation):
		renderStudioV2WorkspaceValidationError(c, err)
	case errors.Is(err, sourcerule.ErrSourceRuleNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: studioV2SourceRuleNotFoundMessage},
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: err.Error()},
		})
	}
}

func workspaceOwnsDevice(record *workspace.Record, deviceID string) bool {
	for _, ownedDeviceID := range record.OrderedDeviceIDs {
		if ownedDeviceID == deviceID {
			return true
		}
	}
	return false
}

func workspaceRuleApplyScopes(deviceID string, links []*schema.SourceRuleLink) []string {
	scopes := make([]string, 0, len(links)+1)
	if deviceID != "" {
		scopes = append(scopes, deviceID)
	}
	for _, link := range links {
		if link == nil {
			continue
		}
		if link.PointID != "" {
			scopes = append(scopes, link.PointID)
		}
	}
	return scopes
}
