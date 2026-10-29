package handlers

import (
	"context"
	"crypto/subtle"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) scope(ctx context.Context, id string, req workspace.WriteGroupSuccessorStageRequest) (recordingplan.SchemaPreviewScope, error) {
	scope, err := h.groups.SuccessorSchemaScope(ctx, id, req)
	if err != nil {
		return scope, err
	}
	if scope.GroupLayout == nil {
		return scope, recordingplan.ErrInvalidGroupSchemaLayout
	}
	if scope.GroupLayout.MySQL == nil && scope.GroupLayout.SQLServer == nil {
		_, err = h.inspector.ValidateDestination(ctx, scope.ConnectorID, scope.ConnectorRevision)
	}
	return scope, err
}
func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) target(id string, req workspace.WriteGroupSuccessorStageRequest, scope recordingplan.SchemaPreviewScope) *checkedGroupSchemaTarget {
	return &checkedGroupSchemaTarget{groups: h.groups, inspector: h.inspector, revalidate: func(ctx context.Context) error {
		current, err := h.scope(ctx, id, req)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(current, scope) {
			return recordingplan.ErrPreviewTokenStale
		}
		return nil
	}}
}
func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) preview(c *gin.Context, req workspace.WriteGroupSuccessorStageRequest) {
	ctx, id := c.Request.Context(), c.Param("id")
	scope, err := h.scope(ctx, id, req)
	if err != nil {
		renderGroupSchemaError(c, err, nil)
		return
	}
	target := h.target(id, req, scope)
	token, err := h.ledger.PrepareSchemaPreview(ctx, scope, target.targetInspector(scope))
	if err == nil {
		err = target.check(ctx)
	}
	if err != nil {
		renderGroupSchemaError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: token})
}
func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) confirm(c *gin.Context, req successorConfirmRequest) {
	ctx, id := c.Request.Context(), c.Param("id")
	stage := req.WriteGroupSuccessorStageRequest
	token, err := h.ledger.PreviewTokenForWorkspace(ctx, stage.Envelope.Request.WorkspaceID, req.Token)
	if err != nil {
		renderGroupSchemaError(c, err, nil)
		return
	}
	if token.PlanID != stage.Current.GroupID || token.GroupLayout == nil || token.OperationID != req.SchemaOperationID || subtle.ConstantTimeCompare([]byte(token.Digest), []byte(req.Digest)) != 1 {
		renderGroupSchemaError(c, recordingplan.ErrSchemaOperationMismatch, nil)
		return
	}
	op, outcome, err := h.ledger.ResolveSchemaApplyReplay(ctx, stage.Envelope.Request.WorkspaceID, req.Token, req.SchemaOperationID)
	if err == nil && mysqlSchemaResumeRequired(token, op) {
		outcome = recordingplan.ClaimInProgress
	}
	proofSaved := false
	if err == nil && outcome == recordingplan.ClaimCompleted {
		proofSaved, err = h.validateCompleted(ctx, id, stage, token, op)
	}
	if err == nil && (outcome == recordingplan.ClaimNone || outcome == recordingplan.ClaimInProgress) {
		var scope recordingplan.SchemaPreviewScope
		scope, err = h.scope(ctx, id, stage)
		if err == nil {
			token, err = h.ledger.ValidatePreviewTokenForApply(ctx, req.Token, scope)
		}
		if err == nil {
			target := h.target(id, stage, scope)
			if scope.GroupLayout.SQLServer != nil {
				op, outcome, err = h.ledger.ApplySQLServerSchemaPreview(ctx, token, recordingplan.SQLServerSchemaApplyTarget{Inspect: target.targetInspector(scope), Execute: func(ctx context.Context, token *recordingplan.SchemaPreviewToken, statements []string) (recordingplan.SchemaExecution, error) {
					return target.executeSQLServerGroupSchemaStatements(ctx, scope, token, statements)
				}})
			} else {
				op, outcome, err = h.ledger.ApplySchemaPreview(ctx, token, recordingplan.SchemaApplyTarget{Inspect: target.targetInspector(scope), Execute: func(ctx context.Context, statements []string) (recordingplan.SchemaExecution, error) {
					return target.executeGroupSchemaStatements(ctx, scope, statements)
				}})
			}
		}
	}
	if err != nil {
		renderGroupSchemaError(c, err, op)
		return
	}
	if op != nil && op.Status == recordingplan.SchemaOperationSucceeded && !proofSaved {
		evidence := &workspace.StorageEvidence{Layout: *token.GroupLayout, CreatedTable: createsTable(token)}
		_, err = h.groups.ConfirmSuccessorSchema(ctx, id, stage, op.OperationID, op.VerifiedDigest, evidence)
		if err != nil {
			renderGroupSchemaError(c, err, op)
			return
		}
	}
	renderSchemaApplyOutcome(c, op, outcome)
}

func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) validateCompleted(ctx context.Context, id string, req workspace.WriteGroupSuccessorStageRequest, token *recordingplan.SchemaPreviewToken, op *recordingplan.SchemaOperation) (bool, error) {
	saved, err := h.groups.Get(ctx, req.Current.GroupID)
	if err != nil {
		return false, err
	}
	proofSaved := op != nil && op.Status == recordingplan.SchemaOperationSucceeded && saved.Group.Destination.SchemaRevision == op.OperationID && saved.Group.Destination.SchemaDigest == op.VerifiedDigest
	currentReq := req
	if proofSaved {
		currentReq.Current.WorkspaceRevision = saved.WorkspaceRevision
	}
	scope, err := h.scope(ctx, id, currentReq)
	if err != nil {
		return false, err
	}
	if proofSaved {
		scope.SchemaRevision, scope.SchemaDigest = token.SchemaRevision, token.SchemaDigest
		scope.WorkspaceRevision = token.WorkspaceRevision
	}
	_, err = h.ledger.ValidatePreviewTokenForApply(ctx, token.Token, scope)
	return proofSaved, err
}
