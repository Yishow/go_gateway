package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

const (
	recordingSchemaNotImplementedCode    = "RECORDING_SCHEMA_NOT_IMPLEMENTED"
	recordingSchemaNotImplementedMessage = "managed schema apply is not implemented"
	recordingSchemaPreviewAgainAction    = "preview the schema again"
)

var errRecordingTargetExecutionEmpty = errors.New("target execution returned no result")

type schemaApplyRequest struct {
	Token       string `json:"token" binding:"required"`
	OperationID string `json:"operation_id" binding:"required"`
	expectedPlanRevisions
}

func (r *schemaApplyRequest) normalize() {
	r.Token = strings.TrimSpace(r.Token)
	r.OperationID = strings.TrimSpace(r.OperationID)
	r.expectedPlanRevisions.normalize()
}

func (r *schemaApplyRequest) complete() bool {
	return r.Token != "" && r.OperationID != "" && r.expectedPlanRevisions.complete()
}

// SchemaApply confirms a schema preview. A repeated confirmation gets the
// operation already recorded for its token (202 while running, 200 once
// finished); another operation on the same target scope gives 409. A new
// confirmation is re-validated against the saved setup and the target before it
// executes once on a verified adapter, and its result is checked on the target.
// @Summary Apply the managed recording schema
// @Description Executes a current schema preview once on a verified adapter and returns the verified operation; repeated confirmations return the recorded operation.
// @Tags Studio V2 Recording Plans
// @Accept json
// @Produce json
// @Param request body schemaApplyRequest true "Schema apply request"
// @Success 200 {object} map[string]interface{} "success=true; data is the finished schema operation (succeeded, partial, failed or unknown)"
// @Success 202 {object} map[string]interface{} "success=true; data is the schema operation that is still running"
// @Failure 400 {object} APIErrorResponse "Missing token, operation_id or expected revisions"
// @Failure 404 {object} APIErrorResponse "success=false; code=RECORDING_SCHEMA_PREVIEW_NOT_FOUND or RECORDING_PLAN_NOT_FOUND"
// @Failure 409 {object} APIErrorResponse "success=false; code=RECORDING_SCHEMA_PREVIEW_STALE, RECORDING_SCHEMA_PREVIEW_EXPIRED, RECORDING_CONNECTOR_REVISION_CONFLICT, RECORDING_SCHEMA_OPERATION_MISMATCH or RECORDING_SCHEMA_OPERATION_BUSY"
// @Failure 422 {object} APIErrorResponse "success=false; incompatible target, missing permission or disabled connector"
// @Failure 501 {object} APIErrorResponse "success=false; code=RECORDING_SCHEMA_NOT_IMPLEMENTED; the adapter has no verified managed schema execution"
// @Failure 503 {object} APIErrorResponse "success=false; code=RECORDING_SCHEMA_OPERATION_UNKNOWN (operation_id names the unresolved operation) or RECORDING_SCHEMA_TARGET_UNAVAILABLE"
// @Router /v1/datalink/studio-v2/workspace/recording-plans/schema-apply [post]
func (h *StudioV2WorkspaceRecordingPlansHandler) SchemaApply(c *gin.Context) {
	var req schemaApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}
	req.normalize()
	if !req.complete() {
		renderStudioV2WorkspaceValidationError(c, errors.New("token, operation_id and the expected workspace, plan and connector revisions are required"))
		return
	}

	ctx := c.Request.Context()
	record, err := h.workspaceSvc.GetOrCreate(ctx)
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}
	op, outcome, err := h.planSvc.ResolveSchemaApplyReplay(ctx, record.ID, req.Token, req.OperationID)
	// A running operation may have lost its owning execution; the apply path
	// resolves a stale claim from the target instead of answering 202 forever.
	if err == nil && (outcome == recordingplan.ClaimNone || outcome == recordingplan.ClaimInProgress) {
		op, outcome, err = h.applyConfirmedPreview(ctx, record, &req)
	}
	if err != nil {
		renderSchemaApplyError(c, err, op)
		return
	}
	renderSchemaApplyOutcome(c, op, outcome)
}

// applyConfirmedPreview re-validates a confirmation against the current saved
// setup and target, then executes it once.
func (h *StudioV2WorkspaceRecordingPlansHandler) applyConfirmedPreview(ctx context.Context, record *workspace.Record, req *schemaApplyRequest) (*recordingplan.SchemaOperation, recordingplan.ClaimOutcome, error) {
	token, err := h.planSvc.PreviewTokenForWorkspace(ctx, record.ID, req.Token)
	if err != nil {
		return nil, "", err
	}
	plan, err := h.getPlanForWorkspace(ctx, record.ID, token.PlanID)
	if err != nil {
		return nil, "", err
	}
	if record.DatabaseSetupRevision != req.ExpectedWorkspaceRevision || plan.Revision != req.ExpectedPlanRevision ||
		strings.TrimSpace(record.DatabaseConnectorID) != token.ConnectorID {
		return nil, "", errRecordingSchemaPreviewStale
	}
	connector, err := h.resolvePreviewConnector(ctx, record, plan, token.ConnectorID, req.ExpectedConnectorRevision)
	if err != nil {
		return nil, "", err
	}
	tablePrefix, err := savedPreviewTablePrefix(plan, token.ConnectorID, token.TablePrefix)
	if err != nil {
		return nil, "", err
	}
	scope := recordingPreviewScope(record, plan, connector, tablePrefix)
	validated, err := h.planSvc.ValidatePreviewTokenForApply(ctx, req.Token, scope)
	if err != nil {
		return nil, "", err
	}
	return h.planSvc.ApplySchemaPreview(ctx, validated, recordingplan.SchemaApplyTarget{
		Inspect: h.recordingTargetInspector(connector.ID, scope.Schema),
		Execute: h.recordingTargetExecutor(connector.ID, scope.Schema),
	})
}

// recordingTargetExecutor runs confirmed statements through the saved target
// service, in the same schema the preview was inspected in.
func (h *StudioV2WorkspaceRecordingPlansHandler) recordingTargetExecutor(connectorID, schemaName string) recordingplan.SchemaExecutor {
	return func(ctx context.Context, statements []string) (recordingplan.SchemaExecution, error) {
		execution, err := h.dbTargetSvc.ExecuteSchemaStatements(ctx, connectorID, schemaName, statements)
		if errors.Is(err, dbtarget.ErrSchemaExecutionPermissionDenied) {
			err = fmt.Errorf("%w: %w", recordingplan.ErrTargetPermissionDenied, err)
		}
		if execution == nil {
			if err == nil {
				err = errRecordingTargetExecutionEmpty
			}
			return recordingplan.SchemaExecution{RolledBack: true}, err
		}
		return recordingplan.SchemaExecution{Committed: execution.Committed, RolledBack: execution.RolledBack}, err
	}
}

func renderSchemaApplyOutcome(c *gin.Context, op *recordingplan.SchemaOperation, outcome recordingplan.ClaimOutcome) {
	switch outcome {
	case recordingplan.ClaimAcquired, recordingplan.ClaimCompleted:
		c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: op})
	case recordingplan.ClaimInProgress:
		c.JSON(http.StatusAccepted, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: op})
	case recordingplan.ClaimScopeBusy:
		renderSchemaOperationError(c, http.StatusConflict, "RECORDING_SCHEMA_OPERATION_BUSY", "another schema operation is using this target",
			"check the running operation, then retry", op.OperationID)
	default:
		renderRecordingOperationNotImplemented(c, recordingSchemaNotImplementedCode, recordingSchemaNotImplementedMessage)
	}
}

func renderSchemaApplyError(c *gin.Context, err error, op *recordingplan.SchemaOperation) {
	switch {
	case errors.Is(err, recordingplan.ErrSchemaOperationUnacknowledged) && op != nil:
		renderSchemaOperationError(c, http.StatusServiceUnavailable, "RECORDING_SCHEMA_OPERATION_UNKNOWN", "the schema operation result could not be recorded",
			"check the operation status before retrying", op.OperationID)
	case errors.Is(err, recordingplan.ErrPreviewTokenNotFound):
		renderRecordingPlanError(c, http.StatusNotFound, "RECORDING_SCHEMA_PREVIEW_NOT_FOUND", "schema preview not found", false, recordingSchemaPreviewAgainAction)
	case errors.Is(err, recordingplan.ErrSchemaOperationMismatch):
		renderRecordingPlanError(c, http.StatusConflict, "RECORDING_SCHEMA_OPERATION_MISMATCH", "operation does not belong to this schema preview", false, "use the operation id returned with the preview")
	case errors.Is(err, recordingplan.ErrPreviewTokenExpired):
		renderRecordingPlanError(c, http.StatusConflict, "RECORDING_SCHEMA_PREVIEW_EXPIRED", "schema preview expired", false, recordingSchemaPreviewAgainAction)
	case errors.Is(err, recordingplan.ErrPreviewTokenLegacy), errors.Is(err, recordingplan.ErrPreviewTokenStale), errors.Is(err, errRecordingSchemaPreviewStale):
		renderRecordingPlanError(c, http.StatusConflict, recordingSchemaPreviewStaleCode, "recording schema preview is stale", false, recordingSchemaPreviewAgainAction)
	case errors.Is(err, recordingplan.ErrSchemaAdapterUnverified):
		renderRecordingOperationNotImplemented(c, recordingSchemaNotImplementedCode, recordingSchemaNotImplementedMessage)
	case errors.Is(err, errRecordingPlanNotFound):
		renderRecordingPlanReadError(c, err)
	case isRecordingSchemaPreviewError(err):
		renderRecordingSchemaPreviewError(c, err)
	default:
		renderRecordingPlanTargetError(c, err)
	}
}

func isRecordingSchemaPreviewError(err error) bool {
	for _, target := range []error{
		errRecordingSchemaPrefixAmbiguous, recordingplan.ErrInvalidTablePrefix, recordingplan.ErrIncompatibleExistingTable,
		recordingplan.ErrTargetPermissionDenied, recordingplan.ErrTargetInspectionUnconfirmed, recordingplan.ErrUnsupportedSchemaDialect,
	} {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func renderSchemaOperationError(c *gin.Context, status int, code, message, action, operationID string) {
	c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: TypedAPIErrorEnvelope{
		Code: code, Message: message, Retryable: true, RequestID: getOrGenerateRequestID(c), Action: action, OperationID: operationID,
	}})
}
