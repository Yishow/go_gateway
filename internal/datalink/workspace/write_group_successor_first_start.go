package workspace

import (
	"context"
	"database/sql"
	"math"
	"strings"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
)

// ConfirmSuccessorFirstStart applies a prepared successor with no predecessor
// history. No predecessor version, revision end, fence or cutover is fabricated.
func (s *WriteGroupService) ConfirmSuccessorFirstStart(ctx context.Context, id string, req WriteGroupSuccessorActivationConfirmation) (*WriteGroupSuccessorFirstStartResult, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	req, intent, err := successorFirstStartIntent(id, req)
	if err != nil {
		return nil, err
	}
	if prior, err := s.recordedSuccessorFirstStart(ctx, req, intent); err != nil || prior != nil {
		return prior, err
	}
	checked, err := s.SuccessorActivationReadiness(ctx, id, req)
	if err != nil {
		if prior, e := s.recordedSuccessorFirstStart(ctx, req, intent); e == nil && prior != nil {
			return prior, nil
		}
		return nil, err
	}
	if checked == nil || checked.snapshot == nil || checked.Review == nil || checked.Readiness == nil || len(checked.Review.Blockers) > 0 || !checked.Readiness.Ready {
		return nil, ErrWriteGroupApplyNotReady
	}
	if checked.Review.ReviewDigest != req.ReviewDigest || checked.ReadinessDigest != req.ReadinessDigest {
		return nil, recreateReservationConflict("activation_review_changed")
	}
	var result *WriteGroupSuccessorFirstStartResult
	_, err = retryWhileSQLiteBusy(ctx, func() (*WriteGroupSaveResult, error) {
		if prior, e := s.recordedSuccessorFirstStart(ctx, req, intent); e != nil || prior != nil {
			result = prior
			return nil, e
		}
		_, e := s.workspaceSvc.updateDatabaseSetupWithCheckpoint(ctx, req.Envelope.Request.Stage.Current.WorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
			if _, e := tx.ExecContext(ctx, `UPDATE wg_delivery_checkpoints SET updated_at=updated_at WHERE 0`); e != nil {
				return e
			}
			review, snap, e := s.activationReviewInTx(ctx, tx, record, id, req.Envelope.Request)
			if e != nil {
				return e
			}
			if review.Envelope.ActivationKind != successorActivationFirstStart || len(review.Blockers) > 0 || review.ReviewDigest != req.ReviewDigest || !sameWriteGroupReadinessSnapshot(checked.snapshot.ready, snap.ready) {
				return recreateReservationConflict("activation_review_changed")
			}
			result, e = s.firstStartSuccessorInTx(ctx, tx, record, req, snap, checked.Readiness)
			return e
		}, func(ctx context.Context, tx *sql.Tx, record *Record) error {
			return s.checkpointSuccessorFirstStart(ctx, tx, record, req, intent, result)
		})
		if e != nil {
			if prior, re := s.recordedSuccessorFirstStart(ctx, req, intent); re == nil && prior != nil {
				result = prior
				return nil, nil
			}
		}
		return nil, e
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("first start successor", err)
	}
	return result, nil
}

func (s *WriteGroupService) firstStartSuccessorInTx(ctx context.Context, tx *sql.Tx, record *Record, req WriteGroupSuccessorActivationConfirmation, snap *successorActivationSnapshot, readiness *WriteGroupReadiness) (*WriteGroupSuccessorFirstStartResult, error) {
	if tx == nil || snap == nil || snap.tx != tx || snap.source != nil || snap.reserved == nil || snap.activationKind != successorActivationFirstStart || readiness == nil {
		return nil, ErrWriteGroupValidation
	}
	boundary := req.Envelope.Request.Boundary
	if !boundary.After(s.clockNow()) {
		return nil, ErrWriteGroupLifecycleBlocked
	}
	if snap.reserved.LineageRevision == math.MaxInt64 {
		return nil, ErrWriteGroupLineageCorrupt
	}
	if err := s.repo.requireSuccessorNoHistoryInTx(ctx, tx, record.ID, snap.reserved.LineageID, snap.predecessor.ID); err != nil {
		return nil, err
	}
	candidate := cloneWriteGroup(snap.reserved.Group)
	if err := s.requireNoOtherSuccessorActiveOwnerInTx(ctx, tx, snap.predecessor, candidate); err != nil {
		return nil, err
	}
	definitions := make(map[string]*measurement.MeasurementDefinition)
	if err := s.repo.validateDefinitionWithMeasurements(ctx, tx, candidate, definitions); err != nil {
		return nil, err
	}
	if err := s.validateSuccessorActivationDefinitionInTx(ctx, tx, snap, candidate); err != nil {
		return nil, err
	}
	for i := range candidate.Members {
		if !sameExactRecordingSource(candidate.Members[i].ExactRecordingSource, snap.reserved.Group.Members[i].ExactRecordingSource) {
			return nil, recreateReservationConflict("source_definition_changed")
		}
	}
	mutation := WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: req.Envelope.Request.Stage.Current.WorkspaceRevision, ExpectedGroupRevision: candidate.Revision, ExpectedConnectorRevision: candidate.Destination.ConnectorRevision}
	applied, err := s.applyValidatedGroupInTx(ctx, tx, record, snap.reserved.Group, candidate, mutation, readiness, definitions, func(*WriteGroup, *WriteGroupAppliedSnapshot) (time.Time, error) { return boundary, nil }, true)
	if err != nil {
		return nil, err
	}
	updated, err := tx.ExecContext(ctx, s.repo.query(`UPDATE write_group_lineages SET current_group_id=$1,reserved_group_id=NULL,reservation_operation_id=NULL,reservation_intent_digest=NULL,revision=revision+1 WHERE workspace_id=$2 AND lineage_id=$3 AND current_group_id=$4 AND reserved_group_id=$5 AND reservation_operation_id=$6 AND reservation_intent_digest=$7 AND revision=$8`), applied.ID, record.ID, snap.reserved.LineageID, snap.predecessor.ID, applied.ID, snap.reserved.OperationID, snap.reserved.IntentDigest, snap.reserved.LineageRevision)
	if err != nil {
		return nil, err
	}
	n, err := updated.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n != 1 {
		return nil, recreateReservationConflict("lineage_changed")
	}
	rootID := recordingplan.SuccessorReserveOperationID(record.ID, snap.reserved.OperationID)
	return &WriteGroupSuccessorFirstStartResult{OperationID: recordingplan.SuccessorFirstStartOperationID(record.ID, rootID), RootOperationID: snap.reserved.OperationID, LineageID: snap.reserved.LineageID, LineageRevision: snap.reserved.LineageRevision + 1, PredecessorGroupID: snap.predecessor.ID, GroupID: applied.ID, AppliedRevision: applied.AppliedRevision, Boundary: boundary, Generation: snap.reserved.Generation, Status: "succeeded"}, nil
}
