package workspace

import (
	"context"
	"database/sql"
	"errors"
	"math"
)

func (s *WriteGroupService) reserveRecreatedDraft(ctx context.Context, req writeGroupRecreateReservationRequest) (*writeGroupReservedDraft, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	req, digest, err := normalizeRecreateReservationRequest(req)
	if err != nil {
		return nil, err
	}
	var reserved *writeGroupReservedDraft
	_, err = retryWhileSQLiteBusy(ctx, func() (*WriteGroupSaveResult, error) {
		replay, found, readErr := s.readPendingRecreateReservation(ctx, req, digest)
		if readErr != nil {
			return nil, readErr
		}
		if found {
			reserved = replay
			return nil, nil
		}
		created, mutationErr := s.createRecreatedDraft(ctx, req, digest)
		if mutationErr == nil {
			reserved = created
			return nil, nil
		}
		// A competing identical winner or a lost commit acknowledgement is success
		// only when the exact original durable operation/intent can be read back.
		replay, found, readErr = s.readPendingRecreateReservation(ctx, req, digest)
		if readErr == nil && found {
			reserved = replay
			return nil, nil
		}
		if errors.Is(readErr, ErrWriteGroupLineageCorrupt) {
			return nil, readErr
		}
		return nil, mutationErr
	})
	if err != nil {
		var conflict *writeGroupReservationConflictError
		if errors.Is(err, ErrWriteGroupRevisionConflict) && !errors.As(err, &conflict) {
			err = errors.Join(recreateReservationConflict("revision_changed"), err)
		}
		return nil, &writeGroupReservationFailure{cause: err}
	}
	return reserved, nil
}

func (s *WriteGroupService) createRecreatedDraft(ctx context.Context, req writeGroupRecreateReservationRequest, digest string) (*writeGroupReservedDraft, error) {
	var saved *writeGroupReservedDraft
	updated, err := s.workspaceSvc.UpdateDatabaseSetup(ctx, req.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		if record.ID != req.WorkspaceID {
			return recreateReservationConflict("scope_changed")
		}
		metadata, found, err := s.repo.loadWriteGroupLineage(ctx, tx, req.WorkspaceID, req.PredecessorGroupID)
		if err != nil {
			return err
		}
		if !found {
			return lineageCorrupt("recreation predecessor has no canonical lineage")
		}
		if metadata.LineageID != req.LineageID || metadata.CurrentGroupID != req.PredecessorGroupID || metadata.Revision != req.ExpectedLineageRevision {
			return recreateReservationConflict("lineage_changed")
		}
		if metadata.ReservedGroupID.Valid {
			return recreateReservationConflict("successor_reserved")
		}
		predecessor, err := s.repo.GetInTx(ctx, tx, req.WorkspaceID, req.PredecessorGroupID)
		if err != nil {
			return err
		}
		if predecessor.Revision != req.ExpectedPredecessorRevision {
			return recreateReservationConflict("predecessor_changed")
		}
		if predecessor.Status != WriteGroupStatusDeleted {
			return ErrWriteGroupLifecycleBlocked
		}
		if predecessor.Destination.ConnectorRevision != req.ExpectedPredecessorConnectorRevision {
			return recreateReservationConflict("source_connector_changed")
		}
		if predecessor.AppliedRevision != "" {
			if s.backlogGuard == nil {
				return writeGroupRemoveActionError(predecessor, "ownership_unavailable", ErrWriteGroupLifecycleBlocked)
			}
			if err := s.backlogGuard.CheckWriteGroupBacklog(ctx, tx, predecessor); err != nil {
				return writeGroupRemoveActionError(predecessor, "ownership_unavailable", errors.Join(ErrWriteGroupLifecycleBlocked, err))
			}
		}
		saved, err = s.repo.reserveRecreatedDraftInTx(ctx, tx, req, digest)
		if err != nil {
			return err
		}
		return projectWriteGroup(record, saved.Group)
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("reserve successor draft", err)
	}
	saved.WorkspaceRevision = updated.DatabaseSetupRevision
	return saved, nil
}

// Error text is closed: callers can classify the original sentinel/cause but
// never accidentally display raw SQL, target configuration or source payloads.
type writeGroupReservationFailure struct{ cause error }

func (e *writeGroupReservationFailure) Error() string {
	return "write group successor reservation unavailable"
}
func (e *writeGroupReservationFailure) Unwrap() error { return e.cause }

func (s *WriteGroupService) readPendingRecreateReservation(ctx context.Context, req writeGroupRecreateReservationRequest, digest string) (result *writeGroupReservedDraft, found bool, err error) {
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, false, ErrWriteGroupServiceUnavailable
	}
	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, false, err
	}
	defer func() {
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result, found, err = nil, false, errors.Join(err, rollbackErr)
		}
	}()
	record, err := setupTx.Get(ctx)
	if err != nil {
		return nil, false, err
	}
	if record.ID != req.WorkspaceID {
		return nil, false, recreateReservationConflict("scope_changed")
	}
	metadata, known, err := s.repo.loadWriteGroupLineage(ctx, setupTx.SQLTx(), req.WorkspaceID, req.PredecessorGroupID)
	if err != nil {
		return nil, false, err
	}
	if !known {
		return nil, false, lineageCorrupt("recreation predecessor has no canonical lineage")
	}
	if metadata.LineageID != req.LineageID || metadata.CurrentGroupID != req.PredecessorGroupID {
		return nil, false, recreateReservationConflict("lineage_changed")
	}
	if !metadata.ReservedGroupID.Valid {
		return nil, false, nil
	}
	if metadata.ReservationOperationID.String != req.OperationID || metadata.ReservationIntentDigest.String != digest {
		return nil, false, recreateReservationConflict("reservation_intent_changed")
	}
	member, known, err := s.repo.loadWriteGroupLineage(ctx, setupTx.SQLTx(), req.WorkspaceID, metadata.ReservedGroupID.String)
	if err != nil {
		return nil, false, err
	}
	if !known || member.LineageID != metadata.LineageID || metadata.Generation == math.MaxInt64 || member.Generation != metadata.Generation+1 ||
		!member.PredecessorGroupID.Valid || member.PredecessorGroupID.String != req.PredecessorGroupID ||
		req.ExpectedLineageRevision == math.MaxInt64 || metadata.Revision != req.ExpectedLineageRevision+1 {
		return nil, false, lineageCorrupt("pending successor membership is inconsistent")
	}
	predecessor, err := s.repo.GetInTx(ctx, setupTx.SQLTx(), req.WorkspaceID, req.PredecessorGroupID)
	if err != nil {
		return nil, false, err
	}
	if predecessor.Status != WriteGroupStatusDeleted || predecessor.Revision != req.ExpectedPredecessorRevision || predecessor.Destination.ConnectorRevision != req.ExpectedPredecessorConnectorRevision {
		return nil, false, lineageCorrupt("pending successor predecessor is inconsistent")
	}
	group, err := s.repo.GetInTx(ctx, setupTx.SQLTx(), req.WorkspaceID, metadata.ReservedGroupID.String)
	if err != nil {
		return nil, false, err
	}
	if group.Status != WriteGroupStatusDraft || group.AppliedRevision != "" {
		return nil, false, lineageCorrupt("pending successor is not an unapplied draft")
	}
	return &writeGroupReservedDraft{WorkspaceRevision: record.DatabaseSetupRevision, LineageID: metadata.LineageID,
		LineageRevision: metadata.Revision, Generation: member.Generation, PredecessorGroupID: req.PredecessorGroupID,
		OperationID: req.OperationID, IntentDigest: digest, Group: cloneWriteGroup(group)}, true, nil
}
