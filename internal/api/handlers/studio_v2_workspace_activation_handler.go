package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type workspaceActivator interface {
	ActivateEligible(ctx context.Context) (*workspace.ActivationResponse, error)
}

// workspaceSchemaEnsurer 在啟動採集前確保 workspace database target 對應的
// 資料表已存在（CREATE TABLE IF NOT EXISTS 為冪等操作）。
type workspaceSchemaEnsurer interface {
	EnsureWorkspaceSchema(ctx context.Context) error
}

// ActivateWorkspaceRequest contains the autosave barrier tokens and revisions
// required before activation.
type ActivateWorkspaceRequest struct {
	WorkspaceRevision string `json:"workspace_revision" binding:"required"`
	SettingsRevision  string `json:"settings_revision" binding:"required"`
	ReadinessToken    string `json:"readiness_token" binding:"required"`
	PendingSaves      int    `json:"pending_saves,omitempty"`
	SaveError         string `json:"save_error,omitempty"`
}

// ActivationBarrierValidator re-reads persisted readiness and revisions on the
// server immediately before activation.
type ActivationBarrierValidator func(context.Context, ActivateWorkspaceRequest) error

// ShareRestoreBarrier restores the durable Share projection before activation.
type ShareRestoreBarrier func(context.Context, ActivateWorkspaceRequest) error

// StudioV2WorkspaceActivationHandler activates a workspace after durable
// autosave, revision, and Share restore barriers pass.
type StudioV2WorkspaceActivationHandler struct {
	activator         workspaceActivator
	schemaEnsurer     workspaceSchemaEnsurer
	auditSvc          workspaceAuditRecorder
	revisionValidator func(ctx context.Context, expectedWsRev, expectedSetRev string) error
	barrierValidator  ActivationBarrierValidator
	shareRestore      ShareRestoreBarrier
}

// WithShareRestore installs the canonical persisted Share restore seam that
// must complete immediately before device activation.
func (h *StudioV2WorkspaceActivationHandler) WithShareRestore(restore ShareRestoreBarrier) *StudioV2WorkspaceActivationHandler {
	h.shareRestore = restore
	return h
}

// WithActivationBarrierValidator installs the backend persisted-state barrier.
func (h *StudioV2WorkspaceActivationHandler) WithActivationBarrierValidator(validator ActivationBarrierValidator) *StudioV2WorkspaceActivationHandler {
	h.barrierValidator = validator
	return h
}

// NewStudioV2WorkspaceActivationHandler creates a workspace activation handler.
func NewStudioV2WorkspaceActivationHandler(activator workspaceActivator, schemaEnsurer ...workspaceSchemaEnsurer) *StudioV2WorkspaceActivationHandler {
	handler := &StudioV2WorkspaceActivationHandler{activator: activator}
	if len(schemaEnsurer) > 0 {
		handler.schemaEnsurer = schemaEnsurer[0]
	}
	return handler
}

// WithRevisionValidator installs the persisted workspace and settings revision
// check that runs immediately before activation.
func (h *StudioV2WorkspaceActivationHandler) WithRevisionValidator(validator func(ctx context.Context, expectedWsRev, expectedSetRev string) error) *StudioV2WorkspaceActivationHandler {
	h.revisionValidator = validator
	return h
}

// Activate starts eligible devices after validating the durable autosave barrier.
// @Summary Activate the Studio v2 workspace
// @Description Activates the persisted workspace after hydration, save, and revision checks.
// @Tags studio-v2
// @Accept json
// @Produce json
// @Param request body ActivateWorkspaceRequest true "Activation barrier revisions and save state"
// @Success 200 {object} map[string]interface{}
// @Failure 409 {object} APIErrorResponse "Workspace or settings revision conflict"
// @Failure 422 {object} APIErrorResponse "Activation barrier or readiness is incomplete"
// @Failure 500 {object} APIErrorResponse "Activation failed"
// @Router /datalink/studio-v2/workspace/activate [post]
func (h *StudioV2WorkspaceActivationHandler) Activate(c *gin.Context) {
	activationReq := ActivateWorkspaceRequest{}
	// 1. 檢查 autosave barrier
	if c.Request.Body != nil {
		var req ActivateWorkspaceRequest
		if err := c.ShouldBindJSON(&req); err == nil {
			activationReq = req
			if req.SaveError != "" || req.PendingSaves > 0 {
				renderSafeError(c, http.StatusUnprocessableEntity, modbusshare.ErrCodeSaveIncomplete, true)
				return
			}
			if h.revisionValidator != nil {
				if err := h.revisionValidator(c.Request.Context(), req.WorkspaceRevision, req.SettingsRevision); err != nil {
					var shareErr *modbusshare.Error
					if errors.As(err, &shareErr) {
						renderModbusShareAPIError(c, http.StatusConflict, shareErr, modbusshare.ErrCodeRevisionConflict, true)
						return
					}
					renderSafeError(c, http.StatusConflict, modbusshare.ErrCodeRevisionConflict, true)
					return
				}
			}
			if h.barrierValidator != nil {
				if err := h.barrierValidator(c.Request.Context(), req); err != nil {
					renderActivationBarrierError(c, err)
					return
				}
			}
		} else if errors.Is(err, io.EOF) {
			renderSafeError(c, http.StatusUnprocessableEntity, modbusshare.ErrCodeSaveIncomplete, true)
			return
		} else if !errors.Is(err, io.EOF) {
			renderSafeError(c, http.StatusBadRequest, ErrCodeActivationRequestInvalid, false)
			return
		}
	}

	// 2. 先確保目標資料表存在，避免採集到值卻因缺表而寫入失敗。
	if h.schemaEnsurer != nil {
		if err := h.schemaEnsurer.EnsureWorkspaceSchema(c.Request.Context()); err != nil {
			renderStudioV2WorkspaceSchemaEnsureError(c, err)
			return
		}
	}
	if h.shareRestore != nil {
		if err := h.shareRestore(c.Request.Context(), activationReq); err != nil {
			renderModbusShareAPIError(c, http.StatusUnprocessableEntity, err, modbusshare.ErrCodeReconcileFailed, true)
			return
		}
	}

	response, err := h.activator.ActivateEligible(c.Request.Context())
	if err != nil {
		var blockedErr *workspace.ReadinessBlockedError
		if errors.As(err, &blockedErr) {
			renderTypedAPIErrorWithData(c, http.StatusUnprocessableEntity, ErrCodeReadinessBlocked, true, gin.H{"issues": blockedErr.BlockingIssues(), "summary": blockedErr.Summary})
			return
		}
		renderSafeError(c, http.StatusInternalServerError, ErrCodeActivationFailed, true)
		return
	}
	h.recordActivationAudit(c.Request.Context(), response)

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    response,
	})
}

func renderActivationBarrierError(c *gin.Context, err error) {
	var shareErr *modbusshare.Error
	if errors.As(err, &shareErr) {
		status := http.StatusUnprocessableEntity
		if shareErr.Code == modbusshare.ErrCodeRevisionConflict {
			status = http.StatusConflict
		}
		renderModbusShareAPIError(c, status, shareErr, modbusshare.ErrCodeSaveIncomplete, true)
		return
	}
	renderSafeError(c, http.StatusUnprocessableEntity, modbusshare.ErrCodeSaveIncomplete, true)
}
