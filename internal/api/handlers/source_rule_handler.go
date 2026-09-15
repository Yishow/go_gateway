package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type SourceRuleHandler struct {
	svc *sourcerule.Service
}

type sourceRuleResponse struct {
	ID               string   `json:"id"`
	DeviceID         string   `json:"device_id"`
	WorkspaceID      string   `json:"workspace_id,omitempty"`
	StartAddress     string   `json:"start_address"`
	Count            int      `json:"count"`
	DataType         string   `json:"data_type"`
	NamingPrefix     string   `json:"naming_prefix"`
	Enabled          bool     `json:"enabled"`
	Locked           bool     `json:"locked"`
	Origin           string   `json:"origin"`
	TemplateName     string   `json:"template_name,omitempty"`
	SkippedAddresses []string `json:"skipped_addresses"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
	RevisionID       string   `json:"revision_id"`
	SaveState        string   `json:"save_state,omitempty"`
	// TargetDataType 目標資料型態（可空）
	TargetDataType *string `json:"target_data_type,omitempty"`
	// ScaleMultiplier 縮放倍率（可空）
	ScaleMultiplier *float64 `json:"scale_multiplier,omitempty"`
	// ScaleOffset 偏移量（可空）
	ScaleOffset *float64 `json:"scale_offset,omitempty"`
	// DataFormat 字節序格式（可空）
	DataFormat         string `json:"data_format,omitempty"`
	ShareEnabled       bool   `json:"share_enabled"`
	ShareStartRegister *int   `json:"share_start_register"`
	ShareStride        *int   `json:"share_stride"`
	// RuntimeApplyStatus V2 autosave 成功後的 runtime 套用狀態。
	RuntimeApplyStatus string `json:"runtime_apply_status,omitempty"`
	// RuntimeApplyMessage 只在 apply_failed 時帶出錯誤說明。
	RuntimeApplyMessage string `json:"runtime_apply_message,omitempty"`
	// RuntimeApplyIssues carries normalized readiness issues when live apply is deferred or warned.
	RuntimeApplyIssues []workspace.ReadinessIssue `json:"runtime_apply_issues,omitempty"`
}

func NewSourceRuleHandler(svc *sourcerule.Service) *SourceRuleHandler {
	return &SourceRuleHandler{svc: svc}
}

func (h *SourceRuleHandler) List(c *gin.Context) {
	filter := sourcerule.ListFilter{}
	if deviceID := c.Query("device_id"); deviceID != "" {
		filter.DeviceID = &deviceID
	}

	rules, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	payload := make([]sourceRuleResponse, 0, len(rules))
	for _, rule := range rules {
		payload = append(payload, mapSourceRuleResponse(rule))
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}

func (h *SourceRuleHandler) Get(c *gin.Context) {
	rule, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: mapSourceRuleResponse(rule)})
}

func (h *SourceRuleHandler) Create(c *gin.Context) {
	var req sourcerule.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	rule, reconcileOutcome, err := h.svc.CreateWithRuntimeReconcile(c.Request.Context(), req)
	if err != nil {
		if renderSourceRuleShareError(c, err) {
			return
		}
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sourcerule.ErrValidation) {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	payload := mapSourceRuleResponse(rule)
	runtimeApply := mapSourceRuleRuntimeReconcileOutcome(reconcileOutcome)
	payload.RuntimeApplyStatus = runtimeApply.Status
	payload.RuntimeApplyMessage = runtimeApply.Message
	payload.RuntimeApplyIssues = runtimeApply.Issues
	c.JSON(http.StatusCreated, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}

func (h *SourceRuleHandler) Update(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	var req sourcerule.UpdateRuleRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	_, req.TargetDataTypeSet = raw["target_data_type"]
	_, req.ScaleMultiplierSet = raw["scale_multiplier"]
	_, req.ScaleOffsetSet = raw["scale_offset"]
	_, req.DataFormatSet = raw["data_format"]
	_, req.ShareStartRegisterSet = raw["share_start_register"]
	_, req.ShareStrideSet = raw["share_stride"]

	rule, reconcileOutcome, err := h.svc.UpdateWithRuntimeReconcile(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if renderSourceRuleShareError(c, err) {
			return
		}
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sourcerule.ErrValidation) {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	payload := mapSourceRuleResponse(rule)
	runtimeApply := mapSourceRuleRuntimeReconcileOutcome(reconcileOutcome)
	payload.RuntimeApplyStatus = runtimeApply.Status
	payload.RuntimeApplyMessage = runtimeApply.Message
	payload.RuntimeApplyIssues = runtimeApply.Issues
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}

func (h *SourceRuleHandler) Delete(c *gin.Context) {
	reconcileOutcome, err := h.svc.DeleteWithRuntimeReconcile(c.Request.Context(), c.Param("id"))
	if err != nil {
		if renderSourceRuleShareError(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: mapStudioV2RuntimeApplyResponse(mapSourceRuleRuntimeReconcileOutcome(reconcileOutcome))})
}

func (h *SourceRuleHandler) Enable(c *gin.Context) {
	id := c.Param("id")
	reconcileOutcome, err := h.svc.EnableWithRuntimeReconcile(c.Request.Context(), id)
	if err != nil {
		if renderSourceRuleShareError(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	rule, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	payload := mapSourceRuleResponse(rule)
	runtimeApply := mapSourceRuleRuntimeReconcileOutcome(reconcileOutcome)
	payload.RuntimeApplyStatus = runtimeApply.Status
	payload.RuntimeApplyMessage = runtimeApply.Message
	payload.RuntimeApplyIssues = runtimeApply.Issues
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}

func (h *SourceRuleHandler) Disable(c *gin.Context) {
	id := c.Param("id")
	reconcileOutcome, err := h.svc.DisableWithRuntimeReconcile(c.Request.Context(), id)
	if err != nil {
		if renderSourceRuleShareError(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	rule, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	payload := mapSourceRuleResponse(rule)
	runtimeApply := mapSourceRuleRuntimeReconcileOutcome(reconcileOutcome)
	payload.RuntimeApplyStatus = runtimeApply.Status
	payload.RuntimeApplyMessage = runtimeApply.Message
	payload.RuntimeApplyIssues = runtimeApply.Issues
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}

func renderSourceRuleShareError(c *gin.Context, err error) bool {
	var shareErr *modbusshare.Error
	if !errors.As(err, &shareErr) {
		return false
	}
	status := http.StatusUnprocessableEntity
	if shareErr.Retryable {
		status = http.StatusServiceUnavailable
	}
	renderModbusShareAPIError(c, status, shareErr, modbusshare.ErrCodeReconcileFailed, shareErr.Retryable)
	return true
}

func mapSourceRuleResponse(rule *schema.SourceRule) sourceRuleResponse {
	var skipped []string
	if rule != nil && rule.SkippedAddresses != "" {
		if err := json.Unmarshal([]byte(rule.SkippedAddresses), &skipped); err != nil {
			skipped = nil
		}
	}

	var targetDataType *string
	if rule.TargetDataType != nil {
		value := string(*rule.TargetDataType)
		targetDataType = &value
	}

	return sourceRuleResponse{
		ID:                 rule.ID,
		DeviceID:           rule.DeviceID,
		StartAddress:       rule.StartAddress,
		Count:              rule.Count,
		DataType:           string(rule.DataType),
		NamingPrefix:       rule.NamingPrefix,
		Enabled:            rule.Enabled,
		Locked:             rule.Locked,
		Origin:             rule.Origin,
		TemplateName:       rule.TemplateName,
		SkippedAddresses:   skipped,
		CreatedAt:          rule.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:          rule.UpdatedAt.Format(time.RFC3339Nano),
		RevisionID:         rule.RevisionID,
		TargetDataType:     targetDataType,
		ScaleMultiplier:    rule.ScaleMultiplier,
		ScaleOffset:        rule.ScaleOffset,
		DataFormat:         rule.DataFormat,
		ShareEnabled:       rule.ShareEnabled,
		ShareStartRegister: rule.ShareStartRegister,
		ShareStride:        rule.ShareStride,
	}
}
