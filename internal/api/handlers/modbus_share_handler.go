package handlers

import (
	"errors"
	"net/http"

	"go-gateway/internal/datalink/modbusshare"

	"github.com/gin-gonic/gin"
)

type ModbusShareHandler struct {
	svc            *modbusshare.Service
	gate           WorkspaceShareGate
	reconciler     *modbusshare.Reconciler
	desiredBuilder DesiredMappingBuilder
}

func NewModbusShareHandler(svc *modbusshare.Service) *ModbusShareHandler {
	return &ModbusShareHandler{svc: svc}
}

func (h *ModbusShareHandler) WithGate(gate WorkspaceShareGate) *ModbusShareHandler {
	h.gate = gate
	return h
}

// @Summary Start Modbus Share listener
// @Tags datalink
// @Accept json
// @Produce json
// @Param request body StartModbusShareRequest true "Listener settings"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} APIErrorResponse
// @Failure 409 {object} APIErrorResponse
// @Failure 422 {object} APIErrorResponse
// @Failure 503 {object} APIErrorResponse
// @Router /datalink/modbus-share/start [post]
func (h *ModbusShareHandler) Start(c *gin.Context) {
	if _, ok := h.checkGates(c); !ok {
		return
	}

	var req StartModbusShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderSafeError(c, http.StatusBadRequest, ErrCodeSettingsInvalid, false)
		return
	}

	port := req.Port
	if port == 0 {
		if h.gate != nil {
			settings, err := h.gate.GetSettings(c.Request.Context())
			if err != nil {
				renderModbusShareAPIError(c, http.StatusServiceUnavailable, err, ErrCodeSettingsUnavailable, true)
				return
			}
			if settings.Port == 0 {
				renderSafeError(c, http.StatusUnprocessableEntity, modbusshare.ErrCodeListenerBindFailed, false)
				return
			}
			port = settings.Port
		} else {
			renderSafeError(c, http.StatusUnprocessableEntity, modbusshare.ErrCodeListenerBindFailed, false)
			return
		}
	}
	if h.gate != nil && req.ExpectedSettingsRevision == "" && !h.svc.SettingsAuthoritative() {
		status := h.svc.Status()
		if status.Enabled && status.Port == port {
			c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: status})
			return
		}
	}
	if err := h.svc.StartCAS(c.Request.Context(), port, req.ExpectedSettingsRevision); err != nil {
		statusCode := http.StatusUnprocessableEntity
		var shareErr *modbusshare.Error
		if errors.As(err, &shareErr) && (shareErr.Message == "modbus share listener is already running" || shareErr.Code == modbusshare.ErrCodeRevisionConflict) {
			statusCode = http.StatusConflict
		}
		renderModbusShareAPIError(c, statusCode, err, modbusshare.ErrCodeListenerBindFailed, true)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": h.svc.Status()})
}

// @Summary Stop Modbus Share listener and persist disabled state
// @Tags datalink
// @Accept json
// @Produce json
// @Param request body StopModbusShareRequest false "Optional settings revision guard"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} APIErrorResponse
// @Failure 409 {object} APIErrorResponse
// @Failure 422 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Failure 503 {object} APIErrorResponse
// @Router /datalink/modbus-share/stop [post]
func (h *ModbusShareHandler) Stop(c *gin.Context) {
	if _, ok := h.checkGates(c); !ok {
		return
	}
	var req StopModbusShareRequest
	if c.Request.ContentLength != 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			renderSafeError(c, http.StatusBadRequest, ErrCodeSettingsInvalid, false)
			return
		}
	}
	if err := h.svc.StopCAS(c.Request.Context(), req.ExpectedSettingsRevision); err != nil {
		statusCode := http.StatusInternalServerError
		var shareErr *modbusshare.Error
		if errors.As(err, &shareErr) && shareErr.Code == modbusshare.ErrCodeRevisionConflict {
			statusCode = http.StatusConflict
		}
		renderModbusShareAPIError(c, statusCode, err, modbusshare.ErrCodeListenerBindFailed, true)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": h.svc.Status()})
}

// @Summary Write a mapped tag value to Modbus Share memory
// @Tags datalink
// @Accept json
// @Produce json
// @Param request body WriteTagValueRequest true "Mapped tag value"
// @Success 200 {object} map[string]interface{}
// @Failure 422 {object} APIErrorResponse
// @Router /datalink/modbus-share/write-tag-value [post]
func (h *ModbusShareHandler) WriteTagValue(c *gin.Context) {
	if _, ok := h.checkGates(c); !ok {
		return
	}
	renderSafeError(c, http.StatusUnprocessableEntity, modbusshare.ErrCodeProjectionRequired, false)
}

// SyncFromMappings syncs latest point values to modbus mirror by enabled mappings.
// @Summary Synchronize mapped point values to Modbus Share memory
// @Tags datalink
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} APIErrorResponse
// @Failure 422 {object} APIErrorResponse
// @Router /datalink/modbus-share/sync [post]
func (h *ModbusShareHandler) SyncFromMappings(c *gin.Context) {
	if _, ok := h.checkGates(c); !ok {
		return
	}
	renderSafeError(c, http.StatusUnprocessableEntity, modbusshare.ErrCodeProjectionRequired, false)
}
