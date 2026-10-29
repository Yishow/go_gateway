package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

// StudioV2WorkspaceWriteGroupSuccessorHandler separates preparation from
// explicit reviewed activation; schema confirmation never activates intake.
type StudioV2WorkspaceWriteGroupSuccessorHandler struct {
	groups    *workspace.WriteGroupService
	ledger    *recordingplan.Service
	inspector *dbtarget.ManagedTableInspector
}

// NewStudioV2WorkspaceWriteGroupSuccessorHandler reuses the existing schema ledger.
func NewStudioV2WorkspaceWriteGroupSuccessorHandler(groups *workspace.WriteGroupService, ledger *recordingplan.Service, targets *dbtarget.ConnectorService) *StudioV2WorkspaceWriteGroupSuccessorHandler {
	h := &StudioV2WorkspaceWriteGroupSuccessorHandler{groups: groups, ledger: ledger}
	if groups != nil && targets != nil {
		h.inspector = groups.ManagedTableInspector(targets)
	}
	return h
}

type successorReserveRequest struct {
	Action string `json:"action"`
	workspace.WriteGroupSuccessorReservationRequest
}
type successorPreviewRequest struct {
	Action string `json:"action"`
	workspace.WriteGroupSuccessorStageRequest
}
type successorConfirmRequest struct {
	Action string `json:"action"`
	workspace.WriteGroupSuccessorStageRequest
	Token             string `json:"token"`
	Digest            string `json:"digest"`
	SchemaOperationID string `json:"schema_operation_id"`
}

// Review is local and read-only: no reservation, preview token or target probe.
// @Summary Review a saved recording successor
// @Description Local read only. Requires original workspace/lineage/G1/connector CAS, candidate group, original operation_id and explicit boundary. Returns a frozen envelope and review_digest, advisory history/blockers/required_checks; no reservation, preview token, external probe, DDL or activation. Later actions must preserve the original envelope and revalidate current scope.
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Original predecessor group ID (G1)"
// @Param request body WriteGroupSuccessorReviewRequest true "Complete local review intent; not mutation authority"
// @Success 200 {object} WriteGroupSuccessorResponse{data=workspace.WriteGroupSuccessorReview}
// @Failure 400 {object} APIErrorResponse "WRITE_GROUP_INVALID_REQUEST: malformed, unknown or missing request fields"
// @Failure 404,409,422,500,503 {object} APIErrorResponse "Typed scope/CAS, lifecycle, definition or availability error; inspect code, action, retryable and request_id"
// @Router /v1/datalink/studio-v2/workspace/write-groups/{id}/successor-review [post]
func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) Review(c *gin.Context) {
	raw, ok := successorRequestBytes(c)
	if !ok {
		return
	}
	var req workspace.WriteGroupSuccessorReviewRequest
	if decodeSuccessorRequest(raw, &req) != nil || strings.TrimSpace(req.WorkspaceID) == "" || strings.TrimSpace(req.ExpectedWorkspaceRevision) == "" || req.Group == nil || req.OperationID == "" || req.Boundary.IsZero() {
		invalidSuccessorRequest(c)
		return
	}
	if h.groups == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	result, err := h.groups.SuccessorReview(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}

// Successor uses explicit stages; unknown actions cannot imply activation.
// @Summary Execute one explicit successor action or read its original operation status
// @Description Swagger 2 cannot express this action-dependent request/response union. The body schema is ONLY a discriminator, not a complete payload: select exactly one complete request and response_data model from x-action-contracts below. Unknown fields and fields from other actions are rejected; there are no per-stage URLs. Every successful response wraps the selected response_data as {success:true,data:...}.
// @Description Scope: path id is original G1. Preserve the immutable original envelope and envelope.request.operation_id (root) for reserve/schema/activation; refreshed current CAS belongs in current for schema actions or stage.current for activation-review (then inside the activation envelope), never rebased into the original envelope. The immutable reserve result is not current CAS. A preview's operation_id is used as schema_operation_id for schema-confirm, never as the root; first-start result operation_id is also distinct from root_operation_id.
// @Description reserve requires action/envelope/review_digest and records only a reserved draft; required capability checks may probe. schema-preview requires action/envelope/review_digest/current, may perform external read-only probes and save a local token, but never DDL. schema-confirm additionally requires token/digest/schema_operation_id and explicit permission for the stored server plan; inspect returned SchemaOperation.status even on HTTP 200, or 202 while in progress.
// @Description activation-review requires action/stage/boundary and is local read only with explicit future B. activation-readiness requires action/activation envelope/review_digest and may inspect the target without DDL or activation. activation-confirm and first-start-confirm both additionally require readiness_digest and recheck scope/CAS/schema/readiness/future B. activation-confirm returns the existing WriteGroupCutover for cutover (valid old v1 may omit activation_kind); first-start-confirm requires first_start and returns WriteGroupSuccessorFirstStartResult, never a fake cutover. The two confirmations are not interchangeable.
// @Description operation-status requires only action/workspace_id/lineage_id/original operation_id and path G1, no current CAS or target access. It is a coherent local read of immutable reserve, all related schema operations, cutover and first-start outcomes. A succeeded stage is not runtime running and not SQL committed. On unknown transport/outcome, retain the original request and query status before any permitted exact retry; do not switch intent/root or assume rollback. No general rebase/cancel endpoint exists.
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Original predecessor group ID (G1), not reserved G2 or current head"
// @Param request body WriteGroupSuccessorAction true "Discriminator only; use exactly one complete request from x-action-contracts"
// @Success 200 {object} WriteGroupSuccessorResponse "success=true; data is exactly the selected x-action-contracts response_data model; schema status may be succeeded, failed, partial or unknown"
// @Success 202 {object} WriteGroupSuccessorSchemaResponse "schema-confirm only: schema operation is still in progress, not confirmed success"
// @Failure 400 {object} APIErrorResponse "WRITE_GROUP_INVALID_REQUEST: unknown/mixed action fields, missing scope/confirmation, malformed JSON or wrong activation kind"
// @Failure 404 {object} APIErrorResponse "WRITE_GROUP_NOT_FOUND or RECORDING_SCHEMA_PREVIEW_NOT_FOUND: original scoped resource missing"
// @Failure 409 {object} APIErrorResponse "revision_mismatch, lifecycle/ownership conflict, stale/expired preview, operation mismatch or busy; preserve original root and inspect typed code/action"
// @Failure 422 {object} APIErrorResponse "Invalid/unsupported transition, boundary, capability, schema, permission or connector; no implicit fallback"
// @Failure 500 {object} APIErrorResponse "Typed local service failure; do not infer that a timed-out mutation rolled back"
// @Failure 501 {object} APIErrorResponse "RECORDING_SCHEMA_NOT_IMPLEMENTED: schema-confirm adapter has no verified execution"
// @Failure 503 {object} APIErrorResponse "Service/target/ownership unavailable or RECORDING_SCHEMA_OPERATION_UNKNOWN; error.operation_id, when present, identifies the schema operation, not the root"
// @x-action-contracts {"reserve":{"request":"#/definitions/handlers.WriteGroupSuccessorReserveRequest","response_data":"#/definitions/workspace.WriteGroupSuccessorPreparation","success_statuses":[200],"error":"#/definitions/handlers.APIErrorResponse"},"schema-preview":{"request":"#/definitions/handlers.WriteGroupSuccessorPreviewRequest","response_data":"#/definitions/recordingplan.SchemaPreviewToken","success_statuses":[200],"error":"#/definitions/handlers.APIErrorResponse"},"schema-confirm":{"request":"#/definitions/handlers.WriteGroupSuccessorConfirmRequest","response_data":"#/definitions/recordingplan.SchemaOperation","success_statuses":[200,202],"error":"#/definitions/handlers.APIErrorResponse"},"activation-review":{"request":"#/definitions/handlers.WriteGroupSuccessorActivationReviewRequest","response_data":"#/definitions/workspace.WriteGroupSuccessorActivationReview","success_statuses":[200],"error":"#/definitions/handlers.APIErrorResponse"},"activation-readiness":{"request":"#/definitions/handlers.WriteGroupSuccessorReadinessRequest","response_data":"#/definitions/workspace.WriteGroupSuccessorActivationReadiness","success_statuses":[200],"error":"#/definitions/handlers.APIErrorResponse"},"activation-confirm":{"request":"#/definitions/handlers.WriteGroupSuccessorCutoverConfirmRequest","response_data":"#/definitions/workspace.WriteGroupCutover","success_statuses":[200],"error":"#/definitions/handlers.APIErrorResponse"},"first-start-confirm":{"request":"#/definitions/handlers.WriteGroupSuccessorFirstStartConfirmRequest","response_data":"#/definitions/workspace.WriteGroupSuccessorFirstStartResult","success_statuses":[200],"error":"#/definitions/handlers.APIErrorResponse"},"operation-status":{"request":"#/definitions/handlers.WriteGroupSuccessorStatusRequest","response_data":"#/definitions/workspace.WriteGroupSuccessorOperationStatus","success_statuses":[200],"error":"#/definitions/handlers.APIErrorResponse"}}
// @Router /v1/datalink/studio-v2/workspace/write-groups/{id}/successor [post]
func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) Successor(c *gin.Context) {
	raw, ok := successorRequestBytes(c)
	if !ok {
		return
	}
	var action struct {
		Action string `json:"action"`
	}
	if json.Unmarshal(raw, &action) != nil {
		invalidSuccessorRequest(c)
		return
	}
	switch action.Action {
	case "first-start-confirm":
		h.firstStartConfirm(c, raw)
	case "operation-status":
		h.operationStatus(c, raw)
	case "activation-confirm":
		h.activationConfirm(c, raw)
	case "activation-review":
		h.activationReview(c, raw)
	case "activation-readiness":
		h.activationReadiness(c, raw)
	case "reserve":
		var req successorReserveRequest
		if decodeSuccessorRequest(raw, &req) != nil || req.ReviewDigest == "" || req.Envelope.Version == "" {
			invalidSuccessorRequest(c)
			return
		}
		if h.groups == nil {
			renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
			return
		}
		result, err := h.groups.ReserveSuccessor(c.Request.Context(), c.Param("id"), req.WriteGroupSuccessorReservationRequest)
		if err != nil {
			renderWriteGroupError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
	case "schema-preview":
		var req successorPreviewRequest
		if decodeSuccessorRequest(raw, &req) != nil || !validSuccessorStage(req.WriteGroupSuccessorStageRequest) {
			invalidSuccessorRequest(c)
			return
		}
		if !h.schemaAvailable(c) {
			return
		}
		h.preview(c, req.WriteGroupSuccessorStageRequest)
	case "schema-confirm":
		var req successorConfirmRequest
		if decodeSuccessorRequest(raw, &req) != nil || !validSuccessorStage(req.WriteGroupSuccessorStageRequest) || req.Token == "" || req.SchemaOperationID == "" || len(req.Digest) != 64 {
			invalidSuccessorRequest(c)
			return
		}
		if !h.schemaAvailable(c) {
			return
		}
		h.confirm(c, req)
	default:
		invalidSuccessorRequest(c)
	}
}

func validSuccessorStage(req workspace.WriteGroupSuccessorStageRequest) bool {
	c := req.Current
	return req.Envelope.Version != "" && req.ReviewDigest != "" && c.WorkspaceRevision != "" && c.LineageRevision > 0 && c.PredecessorRevision != "" && c.PredecessorConnectorRevision != "" && c.GroupID != "" && c.GroupRevision != "" && c.TargetConnectorRevision != ""
}
func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) schemaAvailable(c *gin.Context) bool {
	if h.groups == nil || h.ledger == nil || h.inspector == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return false
	}
	return true
}
func successorRequestBytes(c *gin.Context) ([]byte, bool) {
	const maximum = 1 << 20
	raw, err := io.ReadAll(io.LimitReader(c.Request.Body, maximum+1))
	if err != nil || len(raw) > maximum {
		invalidSuccessorRequest(c)
		return nil, false
	}
	return raw, true
}
func decodeSuccessorRequest(raw []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var trailing any
	if decoder.Decode(&trailing) != io.EOF {
		return errors.New("unexpected request content")
	}
	return nil
}
func invalidSuccessorRequest(c *gin.Context) {
	renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "successor action requires an explicit stage and its reviewed scope", false, "review_request")
}
