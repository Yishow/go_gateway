package workspace

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"go-gateway/internal/datalink/recordingplan"
)

// SuccessorOperationStatus reads only persisted local evidence in one coherent
// transaction. It never probes a target, reconciles DDL, or modifies a ledger.
func (s *WriteGroupService) SuccessorOperationStatus(ctx context.Context, id string, req WriteGroupSuccessorOperationStatusRequest) (result *WriteGroupSuccessorOperationStatus, err error) {
	if err = s.validate(); err != nil {
		return nil, err
	}
	for _, value := range []string{id, req.WorkspaceID, req.LineageID, req.OperationID} {
		if value == "" || strings.TrimSpace(value) != value {
			return nil, ErrWriteGroupValidation
		}
	}
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupServiceUnavailable
	}
	setup, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := setup.Rollback(); e != nil && !errors.Is(e, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, e)
		}
	}()
	record, err := setup.Get(ctx)
	if err != nil {
		return nil, err
	}
	if record == nil || record.ID != req.WorkspaceID {
		return nil, ErrWriteGroupNotFound
	}
	tx := setup.SQLTx()
	root, err := recordingplan.ReadSuccessorReserveInTx(ctx, tx, req.WorkspaceID, req.OperationID)
	if errors.Is(err, recordingplan.ErrSchemaOperationMismatch) || errors.Is(err, recordingplan.ErrSchemaOperationNotFound) {
		return nil, ErrWriteGroupNotFound
	}
	if err != nil {
		return nil, err
	}
	if root == nil || root.Scope.WorkspaceID != req.WorkspaceID || root.Scope.RequestID != req.OperationID ||
		root.Scope.LineageID != req.LineageID || root.Scope.PredecessorGroupID != id {
		return nil, ErrWriteGroupNotFound
	}
	original, reserved, err := decodeSuccessorReserveRecord(root)
	if err != nil {
		return nil, err
	}
	binding, err := decodeSuccessorReserveSchemaBinding(root.SchemaSourceBinding, reserved)
	if err != nil {
		return nil, err
	}
	result = &WriteGroupSuccessorOperationStatus{WorkspaceID: req.WorkspaceID, LineageID: req.LineageID, PredecessorGroupID: id, OperationID: req.OperationID,
		Reserve: WriteGroupSuccessorReserveStatus{Status: "succeeded", Result: reserved},
		Cutover: WriteGroupSuccessorCutoverStatus{Status: "not_submitted"}, FirstStart: WriteGroupSuccessorFirstStartStatus{Status: "not_submitted"}}
	result.SchemaOperations, err = s.successorSchemaOperationsInTx(ctx, tx, reserved, binding)
	if err != nil {
		return nil, err
	}
	result.Cutover.Result, err = s.successorStatusCutoverInTx(ctx, tx, original, reserved, binding)
	if err != nil {
		return nil, err
	}
	if result.Cutover.Result != nil {
		result.Cutover.Status = "succeeded"
	}
	first, err := recordingplan.ReadSuccessorFirstStartInTx(ctx, tx, req.WorkspaceID, recordingplan.SuccessorReserveOperationID(req.WorkspaceID, req.OperationID))
	if err != nil {
		return nil, err
	}
	if first != nil {
		confirmation, started, e := decodeSuccessorFirstStartRecord(first)
		if e != nil {
			return nil, e
		}
		if first.Scope != root.Scope || !sameSuccessorJSON(confirmation.Envelope.Request.Stage.Envelope, *original) ||
			started.GroupID != reserved.Group.ID || result.Cutover.Result != nil {
			return nil, ErrWriteGroupLineageCorrupt
		}
		if e := s.validateSuccessorStatusFirstStart(ctx, tx, confirmation, started, reserved, binding); e != nil {
			return nil, e
		}
		result.FirstStart = WriteGroupSuccessorFirstStartStatus{Status: "succeeded", Result: started}
	}
	return result, nil
}
