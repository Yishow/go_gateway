package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"go-gateway/internal/datalink/modbusshare"

	"github.com/gin-gonic/gin"
)

// @Summary Upsert a workspace-scoped Modbus Share mapping
// @Tags datalink
// @Accept json
// @Produce json
// @Param tagId path string true "Tag identifier"
// @Param request body UpsertMirrorMappingRequest false "Register allocation"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 422 {object} APIErrorResponse
// @Failure 503 {object} APIErrorResponse
// @Router /datalink/modbus-share/mappings/{tagId} [put]
func (h *ModbusShareHandler) UpsertMapping(c *gin.Context) {
	hs, ok := h.checkGates(c)
	if !ok {
		return
	}

	tagID := c.Param("tagId")
	if h.gate != nil && hs.WorkspaceID != "" {
		if !h.gate.ValidateOwnership(c.Request.Context(), hs.WorkspaceID, tagID) {
			renderModbusShareAPIError(c, http.StatusForbidden, modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "tag is not owned by the current workspace", false), modbusshare.ErrCodeWorkspaceScope, false)
			return
		}
	}
	if h.rejectDirectMutation(c, hs, tagID) {
		return
	}

	var reg uint16
	var req UpsertMirrorMappingRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.Register != nil {
		reg = *req.Register
	} else if regStr := c.Query("register"); regStr != "" {
		parsed, parseErr := strconv.ParseUint(regStr, 10, 16)
		if parseErr != nil {
			renderSafeError(c, http.StatusBadRequest, ErrCodeSettingsInvalid, false)
			return
		}
		reg = uint16(parsed)
	} else {
		renderSafeError(c, http.StatusBadRequest, ErrCodeSettingsInvalid, false)
		return
	}

	m, err := h.svc.UpsertMapping(c.Request.Context(), tagID, reg)
	if err != nil {
		var shareErr *modbusshare.Error
		if errors.As(err, &shareErr) {
			renderModbusShareAPIError(c, http.StatusUnprocessableEntity, shareErr, modbusshare.ErrCodeReconcileFailed, true)
			return
		}
		renderModbusShareAPIError(c, http.StatusUnprocessableEntity, err, modbusshare.ErrCodeReconcileFailed, true)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: m})
}

// @Summary Delete a workspace-scoped Modbus Share mapping
// @Tags datalink
// @Produce json
// @Param tagId path string true "Tag identifier"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} APIErrorResponse
// @Failure 422 {object} APIErrorResponse
// @Failure 503 {object} APIErrorResponse
// @Router /datalink/modbus-share/mappings/{tagId} [delete]
func (h *ModbusShareHandler) DeleteMapping(c *gin.Context) {
	hs, ok := h.checkGates(c)
	if !ok {
		return
	}

	tagID := c.Param("tagId")
	if h.gate != nil && hs.WorkspaceID != "" {
		if !h.gate.ValidateOwnership(c.Request.Context(), hs.WorkspaceID, tagID) {
			renderModbusShareAPIError(c, http.StatusForbidden, modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "tag is not owned by the current workspace", false), modbusshare.ErrCodeWorkspaceScope, false)
			return
		}
	}
	if h.rejectDirectMutation(c, hs, tagID) {
		return
	}

	h.svc.RemoveMapping(tagID)
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: gin.H{apiResponseDeletedKey: true}})
}
