package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

// UpdateFormat applies only explicitly scoped configuration intent. Runtime
// activation is separately reported, never implied by a successful save.
func (h *PointHandler) UpdateFormat(c *gin.Context) {
	var req point.FormatMutationRequest
	if c.ShouldBindJSON(&req) != nil {
		renderPointFormatError(c, point.ErrFormatInvalid)
		return
	}
	if req.PointID != c.Param("id") {
		renderPointFormatError(c, point.ErrFormatScope)
		return
	}
	h.mutatePointFormat(c, req, false)
}

func (h *PointHandler) updateLegacyPointFormat(c *gin.Context, id string, patch point.UpdatePointRequest) {
	req := point.FormatMutationRequest{FormatScope: patch.FormatScope}
	if req.PointID == "" {
		renderPointFormatError(c, point.ErrFormatInvalid)
		return
	}
	if req.PointID != id {
		renderPointFormatError(c, point.ErrFormatScope)
		return
	}
	if patch.DataFormat != nil {
		if patch.FormatMode != nil || patch.ConfiguredFormat != nil {
			renderPointFormatError(c, point.ErrFormatInvalid)
			return
		}
		req.FormatMode = schema.PointFormatOverride
		req.ConfiguredFormat = strings.ToUpper(strings.TrimSpace(*patch.DataFormat))
		if req.ConfiguredFormat == "" {
			req.FormatMode = schema.PointFormatInherit
		}
	} else {
		if patch.FormatMode == nil {
			renderPointFormatError(c, point.ErrFormatInvalid)
			return
		}
		req.FormatMode = *patch.FormatMode
		if patch.ConfiguredFormat != nil {
			req.ConfiguredFormat = *patch.ConfiguredFormat
		}
	}
	patch.FormatScope = point.FormatScope{}
	patch.DataFormat = nil
	patch.FormatMode = nil
	patch.ConfiguredFormat = nil
	// Keep mixed legacy updates atomic; no preliminary unguarded ordinary update.
	req.Patch = &patch
	h.mutatePointFormat(c, req, true)
}

func (h *PointHandler) mutatePointFormat(c *gin.Context, req point.FormatMutationRequest, legacy bool) {
	result, err := h.svc.MutateFormat(c.Request.Context(), req)
	if err != nil {
		renderPointFormatError(c, err)
		return
	}
	// Do not UpsertPoint: formal adapter application is a separate lifecycle step.
	if legacy {
		c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: safeRuntimePointRead(&result.Point), "format_resolution": result.Resolution, "workspace_revision": result.WorkspaceRevision, "runtime_apply_status": result.RuntimeApplyStatus, "invalidation_status": result.InvalidationStatus})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}

func renderPointFormatError(c *gin.Context, err error) {
	status, code, message := http.StatusServiceUnavailable, "point_format_unavailable", point.ErrFormatUnavailable.Error()
	switch {
	case errors.Is(err, point.ErrFormatScope):
		status, code, message = http.StatusForbidden, "point_format_scope", point.ErrFormatScope.Error()
	case errors.Is(err, point.ErrFormatConflict):
		status, code, message = http.StatusConflict, "point_format_conflict", point.ErrFormatConflict.Error()
	case errors.Is(err, point.ErrFormatInvalid):
		status, code, message = http.StatusUnprocessableEntity, "point_format_invalid", point.ErrFormatInvalid.Error()
	}
	c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{"code": code, apiResponseMessageKey: message}})
}

// InspectFormatOperation reports a durable pending/invalidated receipt without replay.
func (h *PointHandler) InspectFormatOperation(c *gin.Context) {
	scope := point.FormatScope{WorkspaceID: c.Query("workspace_id"), DeviceID: c.Query("device_id"), RuleID: c.Query("rule_id"), PointID: c.Param("id"), RequestID: c.Param("request_id")}
	result, err := h.svc.InspectFormatMutation(c.Request.Context(), scope)
	if err != nil {
		if errors.Is(err, point.ErrPointNotFound) {
			c.JSON(http.StatusNotFound, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{"code": "point_format_operation_not_found"}})
			return
		}
		renderPointFormatError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}
