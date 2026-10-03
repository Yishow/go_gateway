package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// WriteGroupSchemaRequest confirms only server-resolved saved group state.
// SQL, table names and connection settings are never request inputs.
type WriteGroupSchemaRequest struct {
	WorkspaceID               string `json:"workspace_id"`
	ExpectedWorkspaceRevision string `json:"expected_workspace_revision"`
	ExpectedGroupRevision     string `json:"expected_group_revision"`
	ExpectedConnectorRevision string `json:"expected_connector_revision"`
	Token                     string `json:"token,omitempty"`
	OperationID               string `json:"operation_id,omitempty"`
}

func (r WriteGroupSchemaRequest) mutation() workspace.WriteGroupMutation {
	return workspace.WriteGroupMutation{WorkspaceID: r.WorkspaceID, ExpectedWorkspaceRevision: r.ExpectedWorkspaceRevision, ExpectedGroupRevision: r.ExpectedGroupRevision, ExpectedConnectorRevision: r.ExpectedConnectorRevision}
}

// StudioV2WorkspaceWriteGroupSchemaHandler adapts canonical groups to the
// existing recording schema token and operation ledger.
type StudioV2WorkspaceWriteGroupSchemaHandler struct {
	groups    *workspace.WriteGroupService
	ledger    *recordingplan.Service
	inspector *dbtarget.ManagedTableInspector
	targets   *dbtarget.ConnectorService
}

// NewStudioV2WorkspaceWriteGroupSchemaHandler uses the same schema service as
// advanced recording plans; it owns no token store or ledger.
func NewStudioV2WorkspaceWriteGroupSchemaHandler(_ *workspace.Service, groups *workspace.WriteGroupService, ledger *recordingplan.Service, targets *dbtarget.ConnectorService) *StudioV2WorkspaceWriteGroupSchemaHandler {
	h := &StudioV2WorkspaceWriteGroupSchemaHandler{groups: groups, ledger: ledger, targets: targets}
	if groups != nil {
		h.inspector = groups.ManagedTableInspector(targets)
	}
	return h
}

// Preview reads the saved managed target without creating a SQLite file.
// @Summary Preview a canonical write group's managed schema
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body WriteGroupSchemaRequest true "Saved group revisions"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404,409,422,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id}/schema-preview [post]
func (h *StudioV2WorkspaceWriteGroupSchemaHandler) Preview(c *gin.Context) {
	req, ok := h.request(c, false)
	if !ok {
		return
	}
	scope, err := h.scope(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		renderGroupSchemaError(c, err, nil)
		return
	}
	token, err := h.ledger.PrepareSchemaPreview(c.Request.Context(), scope, h.targetInspector(c.Param("id"), req, scope))
	if err != nil {
		renderGroupSchemaError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: token})
}

// Confirm executes only the stored, explicitly confirmed statement batch.
// @Summary Confirm a canonical write group's managed schema
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body WriteGroupSchemaRequest true "Preview token, operation and revisions"
// @Success 200,202 {object} map[string]interface{}
// @Failure 400,404,409,422,501,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id}/schema-apply [post]
func (h *StudioV2WorkspaceWriteGroupSchemaHandler) Confirm(c *gin.Context) {
	req, ok := h.request(c, true)
	if !ok {
		return
	}
	ctx, id := c.Request.Context(), c.Param("id")
	saved, err := h.groups.Get(ctx, id)
	if err != nil {
		renderGroupSchemaError(c, err, nil)
		return
	}
	if saved.Group.WorkspaceID != req.WorkspaceID || saved.Group.Status == workspace.WriteGroupStatusDeleted {
		renderGroupSchemaError(c, workspace.ErrWriteGroupNotFound, nil)
		return
	}
	token, err := h.ledger.PreviewTokenForWorkspace(ctx, req.WorkspaceID, req.Token)
	if err != nil {
		renderGroupSchemaError(c, err, nil)
		return
	}
	if token.PlanID != id || token.GroupLayout == nil {
		renderGroupSchemaError(c, recordingplan.ErrPreviewTokenNotFound, nil)
		return
	}
	op, outcome, err := h.ledger.ResolveSchemaApplyReplay(ctx, req.WorkspaceID, req.Token, req.OperationID)
	if err == nil && outcome == recordingplan.ClaimCompleted {
		err = h.validateCompletedReplay(ctx, id, req, saved.WorkspaceRevision, token, op)
	}
	if err == nil && (outcome == recordingplan.ClaimNone || outcome == recordingplan.ClaimInProgress) {
		var scope recordingplan.SchemaPreviewScope
		scope, err = h.scope(ctx, id, req)
		if err == nil {
			token, err = h.ledger.ValidatePreviewTokenForApply(ctx, req.Token, scope)
		}
		if err == nil {
			op, outcome, err = h.ledger.ApplySchemaPreview(ctx, token, recordingplan.SchemaApplyTarget{
				Inspect: h.targetInspector(id, req, scope),
				Execute: func(ctx context.Context, statements []string) (recordingplan.SchemaExecution, error) {
					if err := h.checkScope(ctx, id, req, scope); err != nil {
						return recordingplan.SchemaExecution{RolledBack: true}, err
					}
					execution, err := h.inspector.ExecuteSchemaStatements(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, statements)
					return recordingSchemaExecution(execution, err)
				},
			})
		}
	}
	if err != nil {
		renderGroupSchemaError(c, err, op)
		return
	}
	if op != nil && op.Status == recordingplan.SchemaOperationSucceeded {
		if _, err = h.groups.ConfirmManagedSchema(ctx, id, req.mutation(), op.OperationID, op.VerifiedDigest); err != nil {
			renderGroupSchemaError(c, err, op)
			return
		}
	}
	renderSchemaApplyOutcome(c, op, outcome)
}

func (h *StudioV2WorkspaceWriteGroupSchemaHandler) validateCompletedReplay(ctx context.Context, id string, req WriteGroupSchemaRequest, workspaceRevision string, token *recordingplan.SchemaPreviewToken, op *recordingplan.SchemaOperation) error {
	// Recording this operation's verified proof advances the workspace revision.
	// Completed replay is read-only once this proof is saved; later unrelated
	// workspace revisions do not change the recorded operation. Source,
	// group, connector, table and layout must still match the original preview.
	replay := req
	replay.ExpectedWorkspaceRevision = workspaceRevision
	current, err := h.scope(ctx, id, replay)
	if err != nil {
		return err
	}
	if op != nil && op.Status == recordingplan.SchemaOperationSucceeded && current.SchemaRevision == op.OperationID && current.SchemaDigest == op.VerifiedDigest {
		current.SchemaRevision, current.SchemaDigest = token.SchemaRevision, token.SchemaDigest
	}
	current.WorkspaceRevision = token.WorkspaceRevision
	_, err = h.ledger.ValidatePreviewTokenForApply(ctx, req.Token, current)
	return err
}

func (h *StudioV2WorkspaceWriteGroupSchemaHandler) request(c *gin.Context, confirm bool) (WriteGroupSchemaRequest, bool) {
	var req WriteGroupSchemaRequest
	if h.groups == nil || h.ledger == nil || h.targets == nil || h.inspector == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return req, false
	}
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	var extra any
	if err == nil && decoder.Decode(&extra) != io.EOF {
		err = errors.New("unexpected request content")
	}
	if err != nil || strings.TrimSpace(req.WorkspaceID) == "" || strings.TrimSpace(req.ExpectedWorkspaceRevision) == "" || strings.TrimSpace(req.ExpectedGroupRevision) == "" || strings.TrimSpace(req.ExpectedConnectorRevision) == "" || (confirm && (strings.TrimSpace(req.Token) == "" || strings.TrimSpace(req.OperationID) == "")) || (!confirm && (req.Token != "" || req.OperationID != "")) {
		renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "schema confirmation requires the saved group revisions", false, "review_request")
		return req, false
	}
	return req, true
}

func (h *StudioV2WorkspaceWriteGroupSchemaHandler) scope(ctx context.Context, id string, req WriteGroupSchemaRequest) (recordingplan.SchemaPreviewScope, error) {
	scope, err := h.groups.ManagedSchemaScope(ctx, id, req.mutation())
	if err != nil {
		return scope, err
	}
	_, err = h.inspector.ValidateDestination(ctx, scope.ConnectorID, scope.ConnectorRevision)
	return scope, err
}

func (h *StudioV2WorkspaceWriteGroupSchemaHandler) checkScope(ctx context.Context, id string, req WriteGroupSchemaRequest, scope recordingplan.SchemaPreviewScope) error {
	current, err := h.scope(ctx, id, req)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(current, scope) {
		return recordingplan.ErrPreviewTokenStale
	}
	return nil
}

func (h *StudioV2WorkspaceWriteGroupSchemaHandler) targetInspector(id string, req WriteGroupSchemaRequest, scope recordingplan.SchemaPreviewScope) recordingplan.TargetInspector {
	return func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
		if err := h.checkScope(ctx, id, req, scope); err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		result, err := h.inspector.InspectTableAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, table)
		if err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		if result == nil {
			return recordingplan.TargetTableInspection{}, errRecordingTargetInspectionEmpty
		}
		state := recordingplan.TargetTableInspection{Status: string(result.Status), ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{}}
		for _, column := range result.Columns {
			state.Columns = append(state.Columns, column.Name)
			state.ColumnTypes[column.Name] = column.DataType
			state.ColumnNullable[column.Name] = column.Nullable
			state.ColumnPrimaryKey[column.Name] = column.PrimaryKey
		}
		return state, h.checkScope(ctx, id, req, scope)
	}
}

func renderGroupSchemaError(c *gin.Context, err error, op *recordingplan.SchemaOperation) {
	switch {
	case errors.Is(err, dbtarget.ErrManagedDestinationInternal):
		renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "WRITE_GROUP_INTERNAL_DATABASE", "recording requires a distinct application-owned destination", false, "select_recording_destination")
	case errors.Is(err, dbtarget.ErrManagedDestinationInvalid):
		renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "WRITE_GROUP_DESTINATION_INVALID", "the saved recording destination cannot be verified", false, "select_recording_destination")
	case errors.Is(err, recordingplan.ErrIncompatibleExistingTable):
		renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "RECORDING_SCHEMA_INCOMPATIBLE", "an existing managed table is incompatible or has unverified ownership", false, "select_new_managed_table_or_advanced_table")
	case errors.Is(err, workspace.ErrWriteGroupValidation), errors.Is(err, workspace.ErrWriteGroupNotFound), errors.Is(err, workspace.ErrWriteGroupRevisionConflict), errors.Is(err, workspace.ErrSetupRevisionConflict), errors.Is(err, workspace.ErrWriteGroupSourceRevisionConflict):
		renderWriteGroupError(c, err)
	default:
		renderSchemaApplyError(c, err, op)
	}
}
