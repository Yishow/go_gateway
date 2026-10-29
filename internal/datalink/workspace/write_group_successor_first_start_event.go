package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"

	"go-gateway/internal/datalink/recordingplan"
)

const successorFirstStartIntentVersion = "go-gateway/successor-first-start-intent/v1"

type successorFirstStartIntentEnvelope struct {
	Version            string                                    `json:"version"`
	PredecessorGroupID string                                    `json:"predecessor_group_id"`
	Confirmation       WriteGroupSuccessorActivationConfirmation `json:"confirmation"`
}

func successorFirstStartIntent(id string, req WriteGroupSuccessorActivationConfirmation) (WriteGroupSuccessorActivationConfirmation, []byte, error) {
	frozen, _, err := successorActivationIntent(id, req)
	if err != nil {
		return req, nil, err
	}
	if successorActivationKind(frozen.Envelope) != successorActivationFirstStart || frozen.Envelope.SourceAppliedRevision != "" || frozen.Envelope.Request.Stage.Envelope.SourceAppliedRevision != "" || frozen.Envelope.SourceDefinitionDigest != frozen.Envelope.Request.Stage.Envelope.SourceDefinitionDigest {
		return req, nil, ErrWriteGroupValidation
	}
	intent, err := json.Marshal(successorFirstStartIntentEnvelope{successorFirstStartIntentVersion, id, frozen})
	return frozen, intent, err
}

func decodeSuccessorFirstStartRecord(record *recordingplan.SuccessorFirstStartRecord) (*WriteGroupSuccessorActivationConfirmation, *WriteGroupSuccessorFirstStartResult, error) {
	if record == nil {
		return nil, nil, ErrWriteGroupLineageCorrupt
	}
	var envelope successorFirstStartIntentEnvelope
	var result WriteGroupSuccessorFirstStartResult
	if decodeSuccessorJSON(record.Intent, &envelope) != nil || decodeSuccessorJSON(record.Result, &result) != nil || envelope.Version != successorFirstStartIntentVersion {
		return nil, nil, ErrWriteGroupLineageCorrupt
	}
	req, _, err := successorFirstStartIntent(envelope.PredecessorGroupID, envelope.Confirmation)
	if err != nil || !sameSuccessorJSON(envelope.Confirmation, req) || record.IntentDigest != req.ReviewDigest {
		return nil, nil, ErrWriteGroupLineageCorrupt
	}
	e := req.Envelope.Request.Stage.Envelope
	if e.Generation == math.MaxInt64 || req.Envelope.Request.Stage.Current.LineageRevision == math.MaxInt64 || result.Generation < 2 || result.LineageRevision < 2 {
		return nil, nil, ErrWriteGroupLineageCorrupt
	}
	scope := record.Scope
	root := recordingplan.SuccessorReserveOperationID(scope.WorkspaceID, scope.RequestID)
	if scope.WorkspaceID != e.Request.WorkspaceID || scope.RequestID != e.Request.OperationID || scope.LineageID != e.Request.LineageID || scope.PredecessorGroupID != e.PredecessorGroupID || scope.SuccessorGroupID != req.Envelope.Request.Stage.Current.GroupID || record.ReserveOperationID != root ||
		result.OperationID != recordingplan.SuccessorFirstStartOperationID(scope.WorkspaceID, root) || result.RootOperationID != scope.RequestID || result.LineageID != scope.LineageID || result.PredecessorGroupID != scope.PredecessorGroupID || result.GroupID != scope.SuccessorGroupID || result.AppliedRevision != req.Envelope.Request.Stage.Current.GroupRevision || result.AppliedRevision == "" || result.WorkspaceRevision == "" || result.LineageRevision != req.Envelope.Request.Stage.Current.LineageRevision+1 || result.Generation != e.Generation+1 || !result.Boundary.Equal(req.Envelope.Request.Boundary) || result.Status != "succeeded" {
		return nil, nil, ErrWriteGroupLineageCorrupt
	}
	return &req, &result, nil
}

func (s *WriteGroupService) recordedSuccessorFirstStart(ctx context.Context, req WriteGroupSuccessorActivationConfirmation, intent []byte) (result *WriteGroupSuccessorFirstStartResult, err error) {
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
	if len(intent) == 0 {
		return nil, ErrWriteGroupValidation
	}
	e := req.Envelope.Request.Stage.Envelope
	if record.ID != e.Request.WorkspaceID {
		return nil, recreateReservationConflict("scope_changed")
	}
	root, err := recordingplan.ReadSuccessorReserveInTx(ctx, setup.SQLTx(), record.ID, e.Request.OperationID)
	if err != nil {
		return nil, err
	}
	if root == nil {
		return nil, transitionRejection("successor-operation-unrecoverable")
	}
	original, reserved, err := decodeSuccessorReserveRecord(root)
	if err != nil {
		return nil, err
	}
	if !sameSuccessorJSON(*original, e) || root.Scope.SuccessorGroupID != req.Envelope.Request.Stage.Current.GroupID || reserved.Group.Revision != req.Envelope.Request.Stage.Current.GroupRevision {
		return nil, recreateReservationConflict("reservation_intent_changed")
	}
	first, err := recordingplan.ReadSuccessorFirstStartInTx(ctx, setup.SQLTx(), record.ID, recordingplan.SuccessorReserveOperationID(record.ID, e.Request.OperationID))
	if err != nil || first == nil {
		return nil, err
	}
	originalConfirmation, decodedResult, decodeErr := decodeSuccessorFirstStartRecord(first)
	result, err = decodedResult, decodeErr
	if err != nil {
		return nil, err
	}
	if !sameSuccessorJSON(*originalConfirmation, req) {
		return nil, recreateReservationConflict("first_start_intent_changed")
	}
	return result, nil
}

func (s *WriteGroupService) checkpointSuccessorFirstStart(ctx context.Context, tx *sql.Tx, record *Record, req WriteGroupSuccessorActivationConfirmation, intent []byte, result *WriteGroupSuccessorFirstStartResult) error {
	e := req.Envelope.Request.Stage.Envelope
	rootID := recordingplan.SuccessorReserveOperationID(record.ID, e.Request.OperationID)
	result.WorkspaceRevision = record.DatabaseSetupRevision
	payload, err := json.Marshal(result)
	if err != nil {
		return err
	}
	_, err = recordingplan.RecordSuccessorFirstStartInTx(ctx, tx, &recordingplan.SuccessorFirstStartRecord{Scope: recordingplan.SuccessorOperationScope{WorkspaceID: record.ID, RequestID: e.Request.OperationID, LineageID: e.Request.LineageID, PredecessorGroupID: e.PredecessorGroupID, SuccessorGroupID: result.GroupID}, ReserveOperationID: rootID, IntentDigest: req.ReviewDigest, Intent: intent, Result: payload})
	return err
}
