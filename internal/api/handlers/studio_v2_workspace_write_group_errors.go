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
	render := func(status int, code, message string, retryable bool, action string) {
		envelope := TypedAPIErrorEnvelope{
			Code: code, Message: message, Retryable: retryable,
			RequestID: getOrGenerateRequestID(c), Action: action,
		}
		if detail, ok := writeGroupLifecycleErrorDetails(err); ok {
			envelope.Message = detail.message
			envelope.Action = "review_group"
			if detail.next == "stop" {
				envelope.Action = "disable_group"
			}
			c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: struct {
				TypedAPIErrorEnvelope
				ReasonCode string `json:"reason_code"`
				NextAction string `json:"next_action"`
			}{envelope, detail.reason, detail.next}})
			return
		}
		renderStudioV2WorkspaceDatabaseTyped(c, status, envelope)
	}
	switch {
	case errors.Is(err, workspace.ErrWriteGroupNotFound), errors.Is(err, workspace.ErrWriteGroupConnectorNotFound), errors.Is(err, workspace.ErrNotFound):
		render(http.StatusNotFound, "WRITE_GROUP_NOT_FOUND", "write group or resource is not available", false, "reload")
	case errors.Is(err, workspace.ErrSetupRevisionConflict), errors.Is(err, workspace.ErrWriteGroupRevisionConflict):
		render(http.StatusConflict, "revision_mismatch", "saved configuration changed; reload before saving", false, "reload")
	case errors.Is(err, workspace.ErrWriteGroupBasicManagedConflict):
		render(http.StatusConflict, "WRITE_GROUP_BASIC_INTENT_CHANGED", "basic recording configuration changed; reload and explicitly update the saved group", false, "review_group")
	case errors.Is(err, workspace.ErrWriteGroupCircuitProfileConflict):
		render(http.StatusConflict, "WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT", "circuit recording profile already exists with a different configuration or source ownership; reload and update the saved group", false, "review_group")
	case errors.Is(err, workspace.ErrWriteGroupLifecycleBlocked):
		render(http.StatusConflict, "WRITE_GROUP_LIFECYCLE_BLOCKED", "write group cannot be deleted while accepted data ownership is unresolved", false, "disable_group")
	case errors.Is(err, workspace.ErrCutoverFinalizedOverlap):
		render(http.StatusConflict, "recording-cutover-finalized-overlap", "finalized recording history overlaps this boundary; review a new cutover boundary", false, "review_cutover")
	case errors.Is(err, groupdelivery.ErrOwnershipUnavailable):
		render(http.StatusServiceUnavailable, "recording-ownership-unavailable", "recording ownership cannot be verified; review the cutover when its history is available", false, "review_cutover")
	case errors.Is(err, workspace.ErrCutoverNotApplicable):
		render(http.StatusUnprocessableEntity, "cutover_not_applicable", "the reviewed cutover cannot be applied; review its issues", false, "review_cutover")
	case errors.Is(err, workspace.ErrWriteGroupCutoverRequired):
		render(http.StatusConflict, "cutover_required", "these sources are still recorded by an active group; move them with a reviewed cutover", false, "review_cutover")
	case errors.Is(err, workspace.ErrWriteGroupApplyNotReady):
		render(http.StatusConflict, "WRITE_GROUP_NOT_READY", "write group is not ready to apply; review its readiness issues", false, "review_group")
	case errors.Is(err, workspace.ErrWriteGroupProfileIntervalInvalid):
		render(http.StatusUnprocessableEntity, "profile_interval_invalid", "recording interval must be a whole number within its profile range and unit", false, "review_group")
	case errors.Is(err, workspace.ErrWriteGroupMemberMetadataInvalid):
		render(http.StatusUnprocessableEntity, "member_metadata_invalid", "member recording metadata must use the known roles and types", false, "review_group")
	case errors.Is(err, recordingplan.ErrTimescaleUnavailable):
		renderTimescaleUnavailable(c, err)
	case errors.Is(err, workspace.ErrWriteGroupValidation):
		render(http.StatusUnprocessableEntity, "WRITE_GROUP_INVALID", "write-group configuration is invalid", false, "review_group")
	case errors.Is(err, workspace.ErrDatabaseSetupUnavailable), errors.Is(err, workspace.ErrWriteGroupServiceUnavailable), errors.Is(err, workspace.ErrWriteGroupReadinessUnavailable):
		render(http.StatusServiceUnavailable, "WRITE_GROUP_UNAVAILABLE", "write-group service is unavailable", true, "retry")
	default:
		render(http.StatusInternalServerError, "WRITE_GROUP_UNAVAILABLE", "write group could not be saved or read", true, "reload")
	}
}

// Lifecycle metadata extends the existing envelope without changing its status,
// code or retry classification. Only known remove/repair combinations are safe.
func writeGroupLifecycleErrorDetails(err error) (detail struct{ reason, next, message string }, ok bool) {
	var unavailable *workspace.WriteGroupLifecycleActionError
	if !errors.As(err, &unavailable) || unavailable == nil || unavailable.Action != "remove" {
		return detail, false
	}
	switch unavailable.ReasonCode {
	case "connector_identity_stale":
		detail.message = "recording removal is blocked because the saved connector identity changed"
	case "connector_unavailable":
		detail.message = "recording removal is blocked because the saved connector is unavailable"
	case "ownership_unavailable":
		detail.message = "recording removal is blocked because accepted data ownership cannot be verified"
	default:
		return detail, false
	}
	if unavailable.NextAction != "stop" {
		expected := "review_connector"
		if unavailable.ReasonCode == "ownership_unavailable" {
			expected = "inspect_ownership"
		}
		if unavailable.NextAction != expected {
			return detail, false
		}
	}
	detail.reason, detail.next = unavailable.ReasonCode, unavailable.NextAction
	return detail, true
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
