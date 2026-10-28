package workspace

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

// ReserveSuccessor confirms one reviewed intent. It only reserves a draft;
// schema preparation and effective ownership are separate stages.
func (s *WriteGroupService) ReserveSuccessor(ctx context.Context, id string, req WriteGroupSuccessorReservationRequest) (*WriteGroupSuccessorPreparation, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	envelope, review, digest, err := normalizeSuccessorEnvelope(req.Envelope)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(id) != envelope.PredecessorGroupID || req.ReviewDigest != review {
		return nil, recreateReservationConflict("review_changed")
	}
	var result *WriteGroupSuccessorPreparation
	_, err = retryWhileSQLiteBusy(ctx, func() (*WriteGroupSaveResult, error) {
		replay, found, readErr := s.readSuccessorPreparation(ctx, id, envelope, digest)
		if readErr != nil {
			return nil, readErr
		}
		if found {
			result = replay
			return nil, nil
		}
		capability, probeErr := s.probeSuccessorMySQL(ctx, id, envelope, review)
		if probeErr != nil {
			return nil, probeErr
		}
		saved, mutationErr := s.createSuccessorPreparation(ctx, id, envelope, review, digest, capability)
		if mutationErr == nil {
			result = saved
			return nil, nil
		}
		replay, found, readErr = s.readSuccessorPreparation(ctx, id, envelope, digest)
		if readErr == nil && found {
			result = replay
			return nil, nil
		}
		if errors.Is(readErr, ErrWriteGroupLineageCorrupt) {
			return nil, readErr
		}
		return nil, mutationErr
	})
	if err != nil {
		return nil, &writeGroupReservationFailure{cause: err}
	}
	return result, nil
}

func (s *WriteGroupService) createSuccessorPreparation(ctx context.Context, id string, envelope WriteGroupSuccessorEnvelope, review, digest string, capability *successorMySQLCapabilities) (*WriteGroupSuccessorPreparation, error) {
	var result *WriteGroupSuccessorPreparation
	var sourceEvidence *successorReserveSourceEvidence
	updated, err := s.workspaceSvc.updateDatabaseSetupWithCheckpoint(ctx, envelope.Request.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		checked, err := s.reviewSuccessorInTx(ctx, tx, record, id, envelope.Request)
		if err != nil {
			return err
		}
		if checked.ReviewDigest != review || len(checked.Blockers) != 0 {
			return recreateReservationConflict("review_changed")
		}
		pred, err := s.repo.GetInTx(ctx, tx, record.ID, id)
		if err != nil {
			return err
		}
		source, err := s.repo.recreateStructuralSourceInTx(ctx, tx, pred)
		if err != nil {
			return err
		}
		sourceEvidence = &successorReserveSourceEvidence{Version: "successor-reserve-source/v1", Source: source}
		if capability != nil {
			sourceEvidence.SourceMySQL = capability.Source
			sourceEvidence.TargetMySQL = capability.Target
		}
		req := successorRecreateRequest(envelope)
		var saved *writeGroupReservedDraft
		if len(checked.RequiredChecks) > 0 {
			if capability == nil || capability.Target == nil || !successorMySQLRequiredChecks(checked.RequiredChecks, pred.AppliedRevision == "") {
				return ErrWriteGroupValidation
			}
			if err := s.repo.verifySuccessorMySQLInTx(ctx, tx, pred, cloneWriteGroup(envelope.Request.Group), capability.Target, capability.Source); err != nil {
				return err
			}
			evidence := successorMySQLReserveEvidence{tx: tx, digest: digest, workspaceID: record.ID, capability: *capability.Target, sourceCapability: capability.Source}
			authority := successorTransitionAuthority{tx: tx, workspaceID: record.ID, lineageID: envelope.Request.LineageID, predecessorID: id, predecessorRevision: pred.Revision, digest: digest}
			saved, err = s.repo.reserveVerifiedMySQLSuccessorInTx(ctx, tx, req, digest, authority, evidence)
		} else if pred.Status == WriteGroupStatusDeleted {
			saved, err = s.repo.reserveRecreatedDraftInTx(ctx, tx, req, digest)
		} else {
			authority := successorTransitionAuthority{tx: tx, workspaceID: record.ID, lineageID: envelope.Request.LineageID, predecessorID: id, predecessorRevision: pred.Revision, digest: digest}
			saved, err = s.repo.reserveTransitionDraftInTx(ctx, tx, req, digest, authority)
		}
		if err != nil {
			return err
		}
		result = publicSuccessorPreparation(saved)
		return projectWriteGroup(record, saved.Group)
	}, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		return s.checkpointSuccessorReserve(ctx, tx, record, envelope, digest, result, sourceEvidence)
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("reserve reviewed successor", err)
	}
	result.WorkspaceRevision = updated.DatabaseSetupRevision
	return result, nil
}

func successorRecreateRequest(e WriteGroupSuccessorEnvelope) writeGroupRecreateReservationRequest {
	r := e.Request
	return writeGroupRecreateReservationRequest{WorkspaceID: r.WorkspaceID, ExpectedWorkspaceRevision: r.ExpectedWorkspaceRevision, LineageID: r.LineageID, ExpectedLineageRevision: r.ExpectedLineageRevision, PredecessorGroupID: e.PredecessorGroupID, ExpectedPredecessorRevision: r.ExpectedPredecessorRevision, ExpectedPredecessorConnectorRevision: r.ExpectedPredecessorConnectorRevision, ExpectedTargetConnectorRevision: r.ExpectedTargetConnectorRevision, OperationID: r.OperationID, Group: cloneWriteGroup(r.Group)}
}

func publicSuccessorPreparation(r *writeGroupReservedDraft) *WriteGroupSuccessorPreparation {
	return &WriteGroupSuccessorPreparation{WorkspaceRevision: r.WorkspaceRevision, LineageID: r.LineageID, LineageRevision: r.LineageRevision, Generation: r.Generation, PredecessorGroupID: r.PredecessorGroupID, OperationID: r.OperationID, IntentDigest: r.IntentDigest, Group: cloneWriteGroup(r.Group)}
}
