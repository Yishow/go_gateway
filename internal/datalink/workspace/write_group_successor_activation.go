package workspace

import (
	"context"
	"database/sql"
	"math"
	"strings"
	"time"

	"go-gateway/internal/datalink/measurement"
)

// ConfirmSuccessorActivation commits one reviewed observation-time transition.
// Schema confirmation alone never calls this method or enables collection.
func (s *WriteGroupService) ConfirmSuccessorActivation(ctx context.Context, id string, req WriteGroupSuccessorActivationConfirmation) (*WriteGroupCutover, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	req, intent, err := successorActivationIntent(id, req)
	if err != nil {
		return nil, err
	}
	if successorActivationKind(req.Envelope) != successorActivationCutover {
		return nil, ErrWriteGroupValidation
	}
	if prior, err := s.recordedSuccessorActivation(ctx, req, intent); err != nil || prior != nil {
		return prior, err
	}
	checked, err := s.SuccessorActivationReadiness(ctx, id, req)
	if err != nil {
		if prior, replayErr := s.recordedSuccessorActivation(ctx, req, intent); replayErr == nil && prior != nil {
			return prior, nil
		}
		return nil, err
	}
	if checked == nil || checked.snapshot == nil || checked.Review == nil || checked.Readiness == nil || len(checked.Review.Blockers) > 0 || !checked.Readiness.Ready {
		return nil, ErrWriteGroupApplyNotReady
	}
	if checked.Review.ReviewDigest != req.ReviewDigest || checked.ReadinessDigest != req.ReadinessDigest {
		return nil, writeGroupRevisionConflict("successor activation readiness is stale")
	}
	var result *WriteGroupCutover
	_, err = retryWhileSQLiteBusy(ctx, func() (*WriteGroupSaveResult, error) {
		if prior, e := s.recordedSuccessorActivation(ctx, req, intent); e != nil || prior != nil {
			result = prior
			return nil, e
		}
		_, e := s.workspaceSvc.UpdateDatabaseSetup(ctx, req.Envelope.Request.Stage.Current.WorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
			// Acquire write authority before review, preserving raw SQLite busy errors
			// for retry instead of masking an obsolete snapshot as unavailable.
			if _, e := tx.ExecContext(ctx, `UPDATE wg_delivery_checkpoints SET updated_at=updated_at WHERE 0`); e != nil {
				return e
			}
			review, snap, e := s.activationReviewInTx(ctx, tx, record, id, req.Envelope.Request)
			if e != nil {
				return e
			}
			if len(review.Blockers) > 0 || review.ReviewDigest != req.ReviewDigest || !sameWriteGroupReadinessSnapshot(checked.snapshot.ready, snap.ready) {
				return writeGroupRevisionConflict("successor activation evidence changed")
			}
			result, e = s.activateSuccessorInTx(ctx, tx, record, req, intent, snap, checked.Readiness)
			return e
		})
		if e != nil {
			if prior, re := s.recordedSuccessorActivation(ctx, req, intent); re == nil && prior != nil {
				result = prior
				return nil, nil
			}
		}
		return nil, e
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("activate successor", err)
	}
	return result, nil
}

func (s *WriteGroupService) activateSuccessorInTx(ctx context.Context, tx *sql.Tx, record *Record, req WriteGroupSuccessorActivationConfirmation, intent string, snap *successorActivationSnapshot, readiness *WriteGroupReadiness) (*WriteGroupCutover, error) {
	if tx == nil || snap == nil || snap.tx != tx || snap.source == nil || snap.reserved == nil || readiness == nil {
		return nil, ErrWriteGroupValidation
	}
	boundary := req.Envelope.Request.Boundary
	if !boundary.After(s.clockNow()) {
		return nil, ErrWriteGroupLifecycleBlocked
	}
	if snap.reserved.LineageRevision == math.MaxInt64 {
		return nil, ErrWriteGroupLineageCorrupt
	}
	ended := []CutoverEnd{{GroupID: snap.predecessor.ID, Revision: snap.source.AppliedRevision}}
	if err := requireNoFinalizedOverlap(ctx, tx, &cutoverPlan{req: CutoverRequest{Boundary: boundary}, ended: ended}); err != nil {
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
			return nil, writeGroupRevisionConflict("successor source changed during activation")
		}
	}
	mutation := WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: req.Envelope.Request.Stage.Current.WorkspaceRevision, ExpectedGroupRevision: candidate.Revision, ExpectedConnectorRevision: candidate.Destination.ConnectorRevision}
	applied, err := s.applyValidatedGroupInTx(ctx, tx, record, snap.reserved.Group, candidate, mutation, readiness, definitions, func(*WriteGroup, *WriteGroupAppliedSnapshot) (time.Time, error) { return boundary, nil }, true)
	if err != nil {
		return nil, err
	}
	eventID := s.repo.newID()
	moved := make([]CutoverMemberRef, 0, len(snap.source.Group.Members))
	for _, m := range snap.source.Group.Members {
		moved = append(moved, memberRefOf(m))
	}
	sortMemberRefs(moved)
	event := &WriteGroupCutover{ID: eventID, OwnershipFenceRevision: eventID, WorkspaceID: record.ID, Kind: WriteGroupCutoverForward, SourceGroupID: snap.predecessor.ID, Boundary: boundary, ReviewDigest: req.ReviewDigest, Ended: ended, Started: []CutoverStart{{GroupID: applied.ID, Revision: applied.AppliedRevision}}, Moved: moved, CreatedAt: s.clockNow(), intent: intent}
	if err := s.repo.insertRevisionEndInTx(ctx, tx, record.ID, snap.predecessor.ID, snap.source.AppliedRevision, eventID, boundary); err != nil {
		return nil, err
	}
	if err := s.repo.insertCutoverInTx(ctx, tx, event); err != nil {
		return nil, err
	}
	if err := s.repo.insertOwnershipFenceInTx(ctx, tx, event); err != nil {
		return nil, err
	}
	result, err := tx.ExecContext(ctx, s.repo.query(`UPDATE write_group_lineages SET current_group_id=$1,reserved_group_id=NULL,reservation_operation_id=NULL,reservation_intent_digest=NULL,revision=revision+1
 WHERE workspace_id=$2 AND lineage_id=$3 AND current_group_id=$4 AND reserved_group_id=$5 AND reservation_operation_id=$6 AND reservation_intent_digest=$7 AND revision=$8`), applied.ID, record.ID, snap.reserved.LineageID, snap.predecessor.ID, applied.ID, snap.reserved.OperationID, snap.reserved.IntentDigest, snap.reserved.LineageRevision)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n != 1 {
		return nil, writeGroupRevisionConflict("successor lineage changed")
	}
	return event, nil
}
