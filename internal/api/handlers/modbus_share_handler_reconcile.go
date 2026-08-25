package handlers

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"sort"

	"go-gateway/internal/datalink/modbusshare"

	"github.com/gin-gonic/gin"
)

// WithReconciler wires the workspace-scoped projection authority.
func (h *ModbusShareHandler) WithReconciler(reconciler *modbusshare.Reconciler) *ModbusShareHandler {
	h.reconciler = reconciler
	return h
}

// ReconcileRequest is the complete workspace desired projection submitted by
// the candidate/apply seam; direct mapping mutation is never inferred here.
type ReconcileRequest struct {
	WorkspaceID               string                       `json:"workspace_id"`
	ExpectedWorkspaceRevision string                       `json:"expected_workspace_revision"`
	ExpectedSettingsRevision  string                       `json:"expected_settings_revision"`
	ReadinessToken            string                       `json:"readiness_token"`
	DesiredMappings           []modbusshare.DesiredMapping `json:"desired_mappings"`
	CanonicalPlanSignature    string                       `json:"canonical_plan_signature,omitempty"`
}

// DesiredMappingBuilder builds the server-authoritative desired mapping set.
type DesiredMappingBuilder func(context.Context, string, modbusshare.Settings) ([]modbusshare.DesiredMapping, error)

// WithDesiredMappingBuilder configures the server-authoritative mapping builder.
func (h *ModbusShareHandler) WithDesiredMappingBuilder(builder DesiredMappingBuilder) *ModbusShareHandler {
	h.desiredBuilder = builder
	return h
}

// @Summary Reconcile Modbus Share projection
// @Tags datalink
// @Accept json
// @Produce json
// @Param request body ReconcileRequest true "Complete workspace desired mapping set"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 409 {object} APIErrorResponse
// @Failure 422 {object} APIErrorResponse
// @Failure 503 {object} APIErrorResponse
// @Router /datalink/modbus-share/reconcile [post]
func (h *ModbusShareHandler) Reconcile(c *gin.Context) {
	hs, ok := h.checkGates(c)
	if !ok {
		return
	}
	if h.reconciler == nil {
		renderSafeError(c, http.StatusServiceUnavailable, modbusshare.ErrCodeProjectionRequired, true)
		return
	}
	var req ReconcileRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.WorkspaceID == "" || req.ExpectedWorkspaceRevision == "" || req.ExpectedSettingsRevision == "" || req.ReadinessToken == "" {
		renderSafeError(c, http.StatusBadRequest, ErrCodeSettingsInvalid, false)
		return
	}
	if req.WorkspaceID != hs.WorkspaceID {
		renderSafeError(c, http.StatusForbidden, modbusshare.ErrCodeWorkspaceScope, false)
		return
	}
	if h.gate == nil || h.desiredBuilder == nil {
		renderSafeError(c, http.StatusUnprocessableEntity, modbusshare.ErrCodeProjectionRequired, false)
		return
	}
	settings, settingsErr := h.gate.GetSettings(c.Request.Context())
	if settingsErr != nil {
		renderModbusShareAPIError(c, http.StatusServiceUnavailable, settingsErr, ErrCodeSettingsUnavailable, true)
		return
	}
	serverDesired, buildErr := h.desiredBuilder(c.Request.Context(), req.WorkspaceID, settings)
	if buildErr != nil {
		renderModbusShareAPIError(c, http.StatusUnprocessableEntity, buildErr, modbusshare.ErrCodeProjectionRequired, true)
		return
	}
	plan := modbusshare.NewCanonicalSharePlan(req.WorkspaceID, req.ExpectedWorkspaceRevision, req.ExpectedSettingsRevision, serverDesired)
	desired := plan.DesiredMappings
	if req.CanonicalPlanSignature != "" {
		if !plan.Matches(req.CanonicalPlanSignature) {
			renderSafeError(c, http.StatusUnprocessableEntity, modbusshare.ErrCodeProjectionRequired, false)
			return
		}
	} else {
		left, right := append([]modbusshare.DesiredMapping(nil), req.DesiredMappings...), append([]modbusshare.DesiredMapping(nil), serverDesired...)
		sort.Slice(left, func(i, j int) bool { return left[i].TagID < left[j].TagID })
		sort.Slice(right, func(i, j int) bool { return right[i].TagID < right[j].TagID })
		if !reflect.DeepEqual(left, right) {
			renderSafeError(c, http.StatusUnprocessableEntity, modbusshare.ErrCodeProjectionRequired, false)
			return
		}
	}
	outcome, err := h.reconciler.Reconcile(c.Request.Context(), modbusshare.ReconcileRequest{WorkspaceID: req.WorkspaceID, ExpectedWorkspaceRevision: req.ExpectedWorkspaceRevision, ExpectedSettingsRevision: req.ExpectedSettingsRevision, ReadinessToken: req.ReadinessToken, DesiredMappings: desired, CanonicalPlanSignature: plan.Signature})
	if err != nil {
		statusCode := http.StatusUnprocessableEntity
		var shareErr *modbusshare.Error
		if errors.As(err, &shareErr) && shareErr.Code == modbusshare.ErrCodeRevisionConflict {
			statusCode = http.StatusConflict
		}
		renderModbusShareAPIError(c, statusCode, err, modbusshare.ErrCodeReconcileFailed, true)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: outcome})
}
