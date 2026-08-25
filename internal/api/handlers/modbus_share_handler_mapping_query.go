package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/modbusshare"

	"github.com/gin-gonic/gin"
)

// @Summary List workspace-scoped Modbus Share mappings
// @Tags datalink
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} APIErrorResponse
// @Failure 422 {object} APIErrorResponse
// @Failure 503 {object} APIErrorResponse
// @Router /datalink/modbus-share/mappings [get]
func (h *ModbusShareHandler) ListMappings(c *gin.Context) {
	hs, ok := h.checkGates(c)
	if !ok {
		return
	}
	mappings, err := h.svc.ListMappingsForWorkspace(c.Request.Context(), hs.WorkspaceID)
	if err != nil {
		renderModbusShareAPIError(c, http.StatusForbidden, err, modbusshare.ErrCodeWorkspaceScope, false)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: mappings})
}
