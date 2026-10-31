package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
)

// FormatRecoveryResult describes configuration recovery, never recording start
// or proof of an acquired Good value. Disabled Share is explicitly not applied.
type FormatRecoveryResult struct {
	Outcome      string                     `json:"outcome"`
	RuntimeState string                     `json:"runtime_state"`
	ShareApplied bool                       `json:"share_applied"`
	Hydration    modbusshare.HydrationState `json:"hydration"`
}
type PointFormatRecovery func(context.Context, point.FormatRecoveryRequest) (FormatRecoveryResult, error)

func (h *ModbusShareHandler) WithFormatRecovery(recover PointFormatRecovery) *ModbusShareHandler {
	h.formatRecovery = recover
	return h
}

// handleFormatRecovery recognizes only the explicit nested intent on the existing
// route. Ordinary reconcile keeps its existing readiness gate and payload policy.
func (h *ModbusShareHandler) handleFormatRecovery(c *gin.Context) bool {
	var body map[string]json.RawMessage
	if c.ShouldBindBodyWith(&body, binding.JSON) != nil {
		return false
	}
	raw, present := body["format_recovery"]
	if !present {
		return false
	}
	if len(body) != 1 || string(raw) == "null" {
		renderPointFormatError(c, point.ErrFormatInvalid)
		return true
	}
	var req point.FormatRecoveryRequest
	if json.Unmarshal(raw, &req) != nil {
		renderPointFormatError(c, point.ErrFormatInvalid)
		return true
	}
	// Reject injected restore/desired fields instead of silently ignoring them.
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil {
		renderPointFormatError(c, point.ErrFormatInvalid)
		return true
	}
	for key := range fields {
		switch key {
		case "workspace_id", "device_id", "point_id", "rule_id", "expected_workspace_revision", "expected_point_revision", "expected_rule_revision", "request_id", "expected_connection_revision", "expected_settings_revision":
		default:
			renderPointFormatError(c, point.ErrFormatInvalid)
			return true
		}
	}
	if h.formatRecovery == nil {
		renderPointFormatError(c, point.ErrFormatUnavailable)
		return true
	}
	out, err := h.formatRecovery(c.Request.Context(), req)
	if err != nil {
		renderPointFormatError(c, err)
		return true
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: out})
	return true
}
