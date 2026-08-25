package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/modbusshare"

	"github.com/gin-gonic/gin"
)

// @Summary Get Modbus Share status
// @Tags datalink
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} APIErrorResponse
// @Router /datalink/modbus-share/status [get]
func (h *ModbusShareHandler) Status(c *gin.Context) {
	status := h.svc.Status()
	if h.gate != nil {
		hs, err := h.gate.CheckHydration(c.Request.Context())
		if err != nil {
			renderModbusShareAPIError(c, http.StatusServiceUnavailable, err, modbusshare.ErrCodeHydrationRequired, true)
			return
		}
		status = h.svc.StatusForWorkspace(hs.WorkspaceID)
		status.WorkspaceID = hs.WorkspaceID
		status.HydrationState = hs.State
		status.WorkspaceRevision = hs.WorkspaceRevision
		status.SettingsRevision = hs.SettingsRevision
		status.Readiness = hs.State == modbusshare.HydrationStateReady && hs.Readiness
		status.ReadinessToken = hs.ReadinessToken
		ready := status.Readiness
		if !ready {
			status.Enabled = false
			status.BindState = modbusShareBindFailureState
		}
		dirty := false
		if h.reconciler != nil && hs.WorkspaceID != "" {
			revision, revisionErr := h.reconciler.RevisionStatus(c.Request.Context(), hs.WorkspaceID)
			if revisionErr != nil {
				status.Readiness = false
				status.Enabled = false
				status.BindState = modbusShareBindFailureState
				status.Recovery = &modbusshare.RecoveryStatus{
					Code: modbusshare.ErrCodeStorageFailure, Retryable: true,
					Action:    "retry Modbus Share status after durable storage is available",
					RequestID: getOrGenerateRequestID(c),
				}
			} else {
				if revision.Revision != "" {
					status.WorkspaceRevision = revision.Revision
				}
				dirty = revision.Dirty
				if dirty {
					status.DirtyState = modbusShareDirtyUnknownState
					status.Readiness = false
					status.Enabled = false
					status.BindState = modbusShareBindFailureState
					status.Recovery = &modbusshare.RecoveryStatus{
						Code: modbusshare.ErrCodeDirtyUnknown, Retryable: true,
						Action:    "run recovery reconcile",
						RequestID: getOrGenerateRequestID(c),
					}
				}
			}
		}
		st, err := h.gate.GetSettings(c.Request.Context())
		if err != nil {
			renderModbusShareAPIError(c, http.StatusServiceUnavailable, err, ErrCodeSettingsUnavailable, true)
			return
		}
		status.SettingsRevision = st.SettingsRevision
		if ready && !dirty && status.Recovery == nil {
			if !st.Enabled {
				status.Enabled = false
				status.BindState = "disabled"
			}
		}
		if ready && !dirty && status.Recovery == nil {
			if h.desiredBuilder != nil && st.Enabled {
				desired, buildErr := h.desiredBuilder(c.Request.Context(), hs.WorkspaceID, st)
				if buildErr != nil {
					renderModbusShareAPIError(c, http.StatusServiceUnavailable, buildErr, modbusshare.ErrCodeProjectionRequired, true)
					return
				}
				plan := modbusshare.NewCanonicalSharePlan(hs.WorkspaceID, status.WorkspaceRevision, st.SettingsRevision, desired)
				status.CanonicalPlan = &plan
				status.CanonicalDesiredMappings = plan.DesiredMappings
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: status})
}
