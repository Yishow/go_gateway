package handlers

import (
	"context"
	"errors"
	"net/http"

	"go-gateway/internal/datalink/modbusshare"

	"github.com/gin-gonic/gin"
)

// WorkspaceShareGate defines the gate contract for hydration, settings, and ownership checks.
type WorkspaceShareGate interface {
	CheckHydration(ctx context.Context) (modbusshare.HydrationState, error)
	GetSettings(ctx context.Context) (modbusshare.Settings, error)
	ValidateOwnership(ctx context.Context, workspaceID, tagID string) bool
}

type directMutationGate interface {
	ValidateDirectMutation(context.Context, string, string) error
}

// Register 使用指標：若為值型別 uint16，validator 的 required 會把合法位址 0（對應顯示 40001）當成零值而拒絕。
type UpsertMirrorMappingRequest struct {
	Register *uint16 `json:"register" binding:"required"`
}

// WriteTagValueRequest is the direct tag write payload.
type WriteTagValueRequest struct {
	TagID string      `json:"tag_id" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}

// StartModbusShareRequest is the Modbus Share start payload.
type StartModbusShareRequest struct {
	Port                     int    `json:"port"`
	ExpectedSettingsRevision string `json:"expected_settings_revision"`
}

// StopModbusShareRequest is the Modbus Share stop payload.
type StopModbusShareRequest struct {
	ExpectedSettingsRevision string `json:"expected_settings_revision"`
}

func (h *ModbusShareHandler) checkGates(c *gin.Context) (modbusshare.HydrationState, bool) {
	if h.gate == nil {
		return modbusshare.HydrationState{State: modbusshare.HydrationStateReady, Readiness: true}, true
	}

	hs, err := h.gate.CheckHydration(c.Request.Context())
	if err != nil {
		renderModbusShareAPIError(c, http.StatusServiceUnavailable, err, modbusshare.ErrCodeHydrationRequired, true)
		return hs, false
	}
	if hs.State != "ready" || !hs.Readiness {
		renderModbusShareAPIError(c, http.StatusUnprocessableEntity, &modbusshare.Error{Code: modbusshare.ErrCodeHydrationRequired, Message: "workspace hydration is required before accessing Modbus Share", Retryable: true, Action: "complete workspace bootstrap and hydration first"}, modbusshare.ErrCodeHydrationRequired, true)
		return hs, false
	}

	st, err := h.gate.GetSettings(c.Request.Context())
	if err != nil {
		renderModbusShareAPIError(c, http.StatusServiceUnavailable, err, ErrCodeSettingsUnavailable, true)
		return hs, false
	}
	if !st.Enabled {
		renderModbusShareAPIError(c, http.StatusUnprocessableEntity, modbusshare.NewError(modbusshare.ErrCodeDisabled, "Modbus Share is disabled in global settings", false), modbusshare.ErrCodeDisabled, false)
		return hs, false
	}

	return hs, true
}

func (h *ModbusShareHandler) rejectDirectMutation(c *gin.Context, hs modbusshare.HydrationState, tagID string) bool {
	gate, ok := h.gate.(directMutationGate)
	if !ok {
		return false
	}
	if err := gate.ValidateDirectMutation(c.Request.Context(), hs.WorkspaceID, tagID); err != nil {
		var shareErr *modbusshare.Error
		if !errors.As(err, &shareErr) {
			shareErr = modbusshare.NewError(modbusshare.ErrCodeProjectionRequired, "local Modbus projection is required", false)
		}
		renderModbusShareAPIError(c, http.StatusUnprocessableEntity, shareErr, modbusshare.ErrCodeProjectionRequired, false)
		return true
	}
	return false
}
