package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type studioV2WorkspaceBootstrapResponse struct {
	ID               string                        `json:"id"`
	Kind             workspace.Kind                `json:"kind"`
	Status           workspace.Status              `json:"status"`
	OrderedDeviceIDs []string                      `json:"ordered_device_ids"`
	CreatedAt        any                           `json:"created_at"`
	UpdatedAt        any                           `json:"updated_at"`
	ReadinessSummary *workspace.ReadinessSummary   `json:"readiness_summary"`
	ModbusShare      *modbusShareBootstrapResponse `json:"modbus_share,omitempty"`
}

type modbusShareBootstrapResponse struct {
	HydrationState           string                          `json:"hydration_state"`
	WorkspaceID              string                          `json:"workspace_id"`
	WorkspaceRevision        string                          `json:"workspace_revision"`
	SettingsRevision         string                          `json:"settings_revision"`
	Readiness                bool                            `json:"readiness"`
	ReadinessToken           string                          `json:"readiness_token,omitempty"`
	Status                   modbusshare.Status              `json:"status"`
	CanonicalPlan            *modbusshare.CanonicalSharePlan `json:"canonical_plan,omitempty"`
	CanonicalDesiredMappings []modbusshare.DesiredMapping    `json:"canonical_desired_mappings,omitempty"`
	Error                    *TypedAPIErrorEnvelope          `json:"error,omitempty"`
}

// StudioV2WorkspaceHandler serves the Studio v2 workspace bootstrap payload.
type StudioV2WorkspaceHandler struct {
	svc            *workspace.Service
	modbusShare    *modbusshare.Service
	reconciler     *modbusshare.Reconciler
	desiredBuilder DesiredMappingBuilder
}

// NewStudioV2WorkspaceHandler creates a Studio v2 workspace handler.
func NewStudioV2WorkspaceHandler(svc *workspace.Service) *StudioV2WorkspaceHandler {
	return &StudioV2WorkspaceHandler{svc: svc}
}

// WithModbusShare adds the backend-authoritative Share hydration snapshot to
// the same workspace bootstrap response consumed by /studio/v2.
func (h *StudioV2WorkspaceHandler) WithModbusShare(svc *modbusshare.Service) *StudioV2WorkspaceHandler {
	h.modbusShare = svc
	return h
}

// WithModbusShareReconciler wires the durable revision/dirty boundary into
// bootstrap so restart and status responses expose the same scope.
func (h *StudioV2WorkspaceHandler) WithModbusShareReconciler(reconciler *modbusshare.Reconciler) *StudioV2WorkspaceHandler {
	h.reconciler = reconciler
	return h
}

// WithCanonicalPlanBuilder wires the persisted source-rule authority into
// workspace bootstrap so clients can reconcile by plan identity.
func (h *StudioV2WorkspaceHandler) WithCanonicalPlanBuilder(builder DesiredMappingBuilder) *StudioV2WorkspaceHandler {
	h.desiredBuilder = builder
	return h
}

// Get returns the workspace record, readiness summary, and Share hydration
// state used to bootstrap Studio v2.
func (h *StudioV2WorkspaceHandler) Get(c *gin.Context) {
	record, err := h.svc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}
	readinessSummary, err := h.svc.Readiness(c.Request.Context())
	if err != nil {
		readinessSummary = studioV2WorkspaceReadinessUnavailableSummary()
	}

	response := studioV2WorkspaceBootstrapResponse{
		ID:               record.ID,
		Kind:             record.Kind,
		Status:           record.Status,
		OrderedDeviceIDs: append([]string{}, record.OrderedDeviceIDs...),
		CreatedAt:        record.CreatedAt,
		UpdatedAt:        record.UpdatedAt,
		ReadinessSummary: readinessSummary,
	}
	if h.modbusShare != nil {
		hs, hydrationErr := h.modbusShare.CheckHydration(c.Request.Context())
		shareWorkspaceID := hs.WorkspaceID
		if shareWorkspaceID == "" {
			shareWorkspaceID = record.ID
		}
		status := h.modbusShare.StatusForWorkspace(shareWorkspaceID)
		share := &modbusShareBootstrapResponse{
			HydrationState: hs.State, WorkspaceID: shareWorkspaceID, WorkspaceRevision: hs.WorkspaceRevision,
			SettingsRevision: hs.SettingsRevision, Readiness: hs.Readiness, ReadinessToken: hs.ReadinessToken,
			Status: status,
		}
		share.Status.WorkspaceID = shareWorkspaceID
		share.Status.WorkspaceRevision = share.WorkspaceRevision
		share.Status.SettingsRevision = share.SettingsRevision
		share.Status.Readiness = share.Readiness
		durableDirty := false
		if h.reconciler != nil && shareWorkspaceID != "" {
			if revision, revisionErr := h.reconciler.RevisionStatus(c.Request.Context(), shareWorkspaceID); revisionErr != nil {
				share.Readiness = false
				share.Status.Readiness = false
				share.Status.Enabled = false
				share.Status.BindState = modbusShareBindFailureState
				share.Status.Recovery = &modbusshare.RecoveryStatus{Code: modbusshare.ErrCodeStorageFailure, Retryable: true, Action: "retry workspace bootstrap after durable storage is available", RequestID: getOrGenerateRequestID(c)}
				share.Error = &TypedAPIErrorEnvelope{Code: modbusshare.ErrCodeStorageFailure, Message: "workspace Share status is unavailable", Retryable: true, RequestID: getOrGenerateRequestID(c), Action: "Retry workspace bootstrap"}
			} else {
				if revision.Revision != "" {
					share.WorkspaceRevision = revision.Revision
					share.Status.WorkspaceRevision = revision.Revision
				}
				durableDirty = revision.Dirty
				if durableDirty {
					share.Readiness = false
					share.Status.Readiness = false
					share.Status.Enabled = false
					share.Status.DirtyState = modbusShareDirtyUnknownState
					share.Status.BindState = modbusShareBindFailureState
					share.Status.Recovery = &modbusshare.RecoveryStatus{Code: modbusshare.ErrCodeDirtyUnknown, Retryable: true, Action: "run recovery reconcile", RequestID: getOrGenerateRequestID(c)}
					share.Error = &TypedAPIErrorEnvelope{Code: modbusshare.ErrCodeDirtyUnknown, Message: "workspace Share projection requires recovery", Retryable: true, DirtyState: modbusShareDirtyUnknownState, WorkspaceRevision: share.WorkspaceRevision, RequestID: getOrGenerateRequestID(c), Action: "Run recovery reconcile"}
				}
			}
		}
		if hydrationErr == nil && !durableDirty && hs.State == modbusshare.HydrationStateReady && hs.Readiness {
			settings, settingsErr := h.modbusShare.GetSettings(c.Request.Context())
			if settingsErr != nil {
				share.Error = &TypedAPIErrorEnvelope{Code: modbusshare.ErrCodeProjectionRequired, Message: "workspace Share canonical plan is unavailable", Retryable: true, RequestID: getOrGenerateRequestID(c), Action: retryWorkspaceBootstrapAction}
			} else {
				share.SettingsRevision = settings.SettingsRevision
				share.Status.SettingsRevision = settings.SettingsRevision
				if !settings.Enabled {
					share.Error = &TypedAPIErrorEnvelope{Code: modbusshare.ErrCodeDisabled, Message: "workspace Share canonical plan is unavailable while Modbus Share is disabled", Retryable: false, RequestID: getOrGenerateRequestID(c), Action: "Enable Modbus Share before requesting its canonical plan"}
				} else if h.desiredBuilder != nil {
					if desired, buildErr := h.desiredBuilder(c.Request.Context(), record.ID, settings); buildErr != nil {
						share.Error = &TypedAPIErrorEnvelope{Code: modbusshare.ErrCodeProjectionRequired, Message: "workspace Share canonical plan is unavailable", Retryable: true, RequestID: getOrGenerateRequestID(c), Action: retryWorkspaceBootstrapAction}
					} else {
						plan := modbusshare.NewCanonicalSharePlan(record.ID, share.WorkspaceRevision, settings.SettingsRevision, desired)
						share.CanonicalPlan = &plan
						share.CanonicalDesiredMappings = plan.DesiredMappings
						share.Status.CanonicalPlan = &plan
						share.Status.CanonicalDesiredMappings = plan.DesiredMappings
					}
				}
			}
		}
		if !durableDirty && (hydrationErr != nil || (hs.State != "" && hs.State != modbusshare.HydrationStateReady)) {
			if hydrationErr != nil {
				share.HydrationState = modbusshare.HydrationStateFailed
			}
			share.Readiness = false
			share.Error = &TypedAPIErrorEnvelope{Code: modbusshare.ErrCodeHydrationRequired, Message: "workspace Share hydration is unavailable", Retryable: true, RequestID: getOrGenerateRequestID(c), Action: "Retry workspace bootstrap"}
		}
		response.ModbusShare = share
	}
	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    response,
	})
}

func studioV2WorkspaceReadinessUnavailableSummary() *workspace.ReadinessSummary {
	return &workspace.ReadinessSummary{
		Ready:         false,
		BlockingCount: 1,
		Issues: []workspace.ReadinessIssue{
			{
				Code:     "readiness-unavailable",
				Severity: workspace.ReadinessSeverityBlocking,
				Step:     workspace.ReadinessStep1,
				Scope:    "workspace",
				Message:  "workspace readiness is unavailable",
			},
		},
	}
}
