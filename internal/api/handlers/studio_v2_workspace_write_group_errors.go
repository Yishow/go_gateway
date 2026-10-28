package handlers

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

func renderWriteGroupError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, workspace.ErrWriteGroupNotFound), errors.Is(err, workspace.ErrWriteGroupConnectorNotFound), errors.Is(err, workspace.ErrNotFound):
		renderWriteGroupTyped(c, http.StatusNotFound, "WRITE_GROUP_NOT_FOUND", "write group or resource is not available", false, "reload")
	case errors.Is(err, workspace.ErrSetupRevisionConflict), errors.Is(err, workspace.ErrWriteGroupRevisionConflict):
		renderWriteGroupTyped(c, http.StatusConflict, "revision_mismatch", "saved configuration changed; reload before saving", false, "reload")
	case errors.Is(err, workspace.ErrWriteGroupBasicManagedConflict):
		renderWriteGroupTyped(c, http.StatusConflict, "WRITE_GROUP_BASIC_INTENT_CHANGED", "basic recording configuration changed; reload and explicitly update the saved group", false, "review_group")
	case errors.Is(err, workspace.ErrWriteGroupCircuitProfileConflict):
		renderWriteGroupTyped(c, http.StatusConflict, "WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT", "circuit recording profile already exists with a different configuration or source ownership; reload and update the saved group", false, "review_group")
	case errors.Is(err, workspace.ErrWriteGroupLifecycleBlocked):
		renderWriteGroupTyped(c, http.StatusConflict, "WRITE_GROUP_LIFECYCLE_BLOCKED", "write group cannot be deleted while accepted data ownership is unresolved", false, "disable_group")
	case errors.Is(err, workspace.ErrCutoverFinalizedOverlap):
		renderWriteGroupTyped(c, http.StatusConflict, "recording-cutover-finalized-overlap", "finalized recording history overlaps this boundary; review a new cutover boundary", false, "review_cutover")
	case errors.Is(err, groupdelivery.ErrOwnershipUnavailable):
		renderWriteGroupTyped(c, http.StatusServiceUnavailable, "recording-ownership-unavailable", "recording ownership cannot be verified; review the cutover when its history is available", false, "review_cutover")
	case errors.Is(err, workspace.ErrCutoverNotApplicable):
		renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "cutover_not_applicable", "the reviewed cutover cannot be applied; review its issues", false, "review_cutover")
	case errors.Is(err, workspace.ErrWriteGroupCutoverRequired):
		renderWriteGroupTyped(c, http.StatusConflict, "cutover_required", "these sources are still recorded by an active group; move them with a reviewed cutover", false, "review_cutover")
	case errors.Is(err, workspace.ErrWriteGroupApplyNotReady):
		renderWriteGroupTyped(c, http.StatusConflict, "WRITE_GROUP_NOT_READY", "write group is not ready to apply; review its readiness issues", false, "review_group")
	case errors.Is(err, workspace.ErrWriteGroupProfileIntervalInvalid):
		renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "profile_interval_invalid", "recording interval must be a whole number within its profile range and unit", false, "review_group")
	case errors.Is(err, workspace.ErrWriteGroupMemberMetadataInvalid):
		renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "member_metadata_invalid", "member recording metadata must use the known roles and types", false, "review_group")
	case errors.Is(err, recordingplan.ErrTimescaleUnavailable):
		renderTimescaleUnavailable(c, err)
	case errors.Is(err, workspace.ErrWriteGroupValidation):
		renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "WRITE_GROUP_INVALID", "write-group configuration is invalid", false, "review_group")
	case errors.Is(err, workspace.ErrDatabaseSetupUnavailable), errors.Is(err, workspace.ErrWriteGroupServiceUnavailable), errors.Is(err, workspace.ErrWriteGroupReadinessUnavailable):
		renderWriteGroupTyped(c, http.StatusServiceUnavailable, "WRITE_GROUP_UNAVAILABLE", "write-group service is unavailable", true, "retry")
	default:
		renderWriteGroupTyped(c, http.StatusInternalServerError, "WRITE_GROUP_UNAVAILABLE", "write group could not be saved or read", true, "reload")
	}
}

// renderTimescaleUnavailable answers a selected timescale mode that cannot be
// prepared: the stable timescale_unavailable code, the stable reasons, and one
// repair action. The selection is never rewritten to plain storage; SQLite can
// only be fixed by choosing plain, any other cause by repairing the destination.
func renderTimescaleUnavailable(c *gin.Context, err error) {
	var unavailable *recordingplan.TimescaleUnavailableError
	reasons := "capability-unconfirmed"
	action := "repair_timescale_capability"
	if errors.As(err, &unavailable) && len(unavailable.Reasons) > 0 {
		reasons = strings.Join(unavailable.Reasons, ", ")
		if slices.Contains(unavailable.Reasons, dbtarget.TimescaleReasonDialect) {
			action = "select_plain_storage"
		}
	}
	renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "timescale_unavailable",
		"TimescaleDB storage cannot be prepared on the selected destination ("+reasons+")", false, action)
}

func renderWriteGroupTyped(c *gin.Context, status int, code, message string, retryable bool, action string) {
	renderStudioV2WorkspaceDatabaseTyped(c, status, TypedAPIErrorEnvelope{
		Code: code, Message: message, Retryable: retryable,
		RequestID: getOrGenerateRequestID(c), Action: action,
	})
}
