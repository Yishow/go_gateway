package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
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

// SourceRuleEnvelope documents the success wrapper of a saved source rule.
type SourceRuleEnvelope struct {
	Success bool               `json:"success" example:"true"`
	Data    sourceRuleResponse `json:"data"`
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

// Create saves a workspace source rule and derives its points.
// @Summary Create workspace source rule
// @Description A homogeneous rule expands start_address, count and data_type. A mixed rule (layout_mode=mixed) expands one point per mixed_items entry; items carry their own address and data_type, may use different areas, and must not overlap. Validation errors name the field and, for item conflicts, item_ids.
// @Tags Studio V2 Source Rules
// @Accept json
// @Produce json
// @Param body body sourcerule.CreateRuleRequest true "Source rule"
// @Success 201 {object} SourceRuleEnvelope
// @Failure 400 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /v1/datalink/studio-v2/workspace/source-rules [post]
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

	rule, reconcileOutcome, err := h.ruleSvc.CreateWithRuntimeReconcile(c.Request.Context(), req)
	if err != nil {
		if renderSourceRuleShareError(c, err) {
			return
		}
		renderStudioV2WorkspaceSourceRuleError(c, err)
		return
	}

	payload := mapSourceRuleResponse(rule)
	payload.WorkspaceID = record.ID
	payload.SaveState = workspaceSaveStateSaved
	applyOutcome := mapSourceRuleRuntimeReconcileOutcome(reconcileOutcome)
	if reconcileOutcome.Status == sourcerule.RuntimeReconcileStatusDeferred {
		links, linksErr := h.ruleSvc.ListLinks(c.Request.Context(), rule.ID)
		if linksErr != nil {
			renderStudioV2WorkspaceSourceRuleError(c, linksErr)
			return
		}
		applyOutcome = resolveStudioV2ScopedRuntimeApplyOutcome(
			c.Request.Context(),
			h.workspaceSvc,
			h.deviceSvc,
			[]string{rule.DeviceID},
			workspaceRuleApplyScopes(rule.DeviceID, links),
		)
	}
	payload.RuntimeApplyStatus = applyOutcome.Status
	payload.RuntimeApplyMessage = applyOutcome.Message
	payload.RuntimeApplyIssues = applyOutcome.Issues
	c.JSON(http.StatusCreated, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}

// Update changes a workspace source rule and reconciles its derived points.
// @Summary Update workspace source rule
// @Description Omitted fields keep their stored values. layout_mode switches the layout; mixed_items replaces the complete item list of a mixed rule, and items are matched to existing points by item_id.
// @Tags Studio V2 Source Rules
// @Accept json
// @Produce json
// @Param id path string true "Source rule ID"
// @Param body body sourcerule.UpdateRuleRequest true "Changed fields"
// @Success 200 {object} SourceRuleEnvelope
// @Failure 400 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /v1/datalink/studio-v2/workspace/source-rules/{id} [put]
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
	var fieldErr *sourcerule.FieldError
	switch {
	case errors.As(err, &fieldErr):
		payload := gin.H{apiResponseCodeKey: apiValidationErrorCode, apiResponseMessageKey: err.Error(), "field": fieldErr.Field}
		if len(fieldErr.ItemIDs) > 0 {
			payload["item_ids"] = fieldErr.ItemIDs
		}
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: payload})
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
	return record != nil && slices.Contains(record.OrderedDeviceIDs, deviceID)
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
