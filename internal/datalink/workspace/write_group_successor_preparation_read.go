package workspace

import (
	"context"
	"database/sql"
	"errors"
	"math"
)

func (s *WriteGroupService) readSuccessorPreparation(ctx context.Context, id string, e WriteGroupSuccessorEnvelope, digest string) (result *WriteGroupSuccessorPreparation, found bool, err error) {
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, false, ErrWriteGroupServiceUnavailable
	}
	setup, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, false, err
	}
	defer func() {
		if x := setup.Rollback(); x != nil && !errors.Is(x, sql.ErrTxDone) {
			result = nil
			found = false
			err = errors.Join(err, x)
		}
	}()
	record, err := setup.Get(ctx)
	if err != nil {
		return nil, false, err
	}
	if record.ID != e.Request.WorkspaceID || id != e.PredecessorGroupID {
		return nil, false, recreateReservationConflict("scope_changed")
	}
	if result, found, err := s.recordedSuccessorReserveInTx(ctx, setup.SQLTx(), e, digest); err != nil || found {
		return result, found, err
	}
	return s.pendingSuccessorInTx(ctx, setup.SQLTx(), record, id, e, digest)
}

func (s *WriteGroupService) pendingSuccessorInTx(ctx context.Context, tx *sql.Tx, record *Record, id string, e WriteGroupSuccessorEnvelope, digest string) (*WriteGroupSuccessorPreparation, bool, error) {
	if record == nil || record.ID != e.Request.WorkspaceID || id != e.PredecessorGroupID {
		return nil, false, recreateReservationConflict("scope_changed")
	}
	meta, known, err := s.repo.loadWriteGroupLineage(ctx, tx, record.ID, id)
	if err != nil {
		return nil, false, err
	}
	if !known {
		return nil, false, ErrWriteGroupLineageCorrupt
	}
	if meta.LineageID != e.Request.LineageID || meta.CurrentGroupID != id {
		return nil, false, recreateReservationConflict("lineage_changed")
	}
	if !meta.ReservedGroupID.Valid {
		return nil, false, nil
	}
	if meta.ReservationOperationID.String != e.Request.OperationID || meta.ReservationIntentDigest.String != digest {
		return nil, false, recreateReservationConflict("reservation_intent_changed")
	}
	successor, known, err := s.repo.loadWriteGroupLineage(ctx, tx, record.ID, meta.ReservedGroupID.String)
	if err != nil {
		return nil, false, err
	}
	if !known || e.Generation == math.MaxInt64 || successor.Generation != e.Generation+1 || successor.LineageID != meta.LineageID || successor.PredecessorGroupID.String != id || e.Request.ExpectedLineageRevision == math.MaxInt64 || meta.Revision != e.Request.ExpectedLineageRevision+1 {
		return nil, false, ErrWriteGroupLineageCorrupt
	}
	group, err := s.repo.GetInTx(ctx, tx, record.ID, meta.ReservedGroupID.String)
	if err != nil {
		return nil, false, err
	}
	if group.Status != WriteGroupStatusDraft || group.AppliedRevision != "" {
		return nil, false, ErrWriteGroupLineageCorrupt
	}
	// Schema proof and server-owned revisions may advance, but the reserved
	// recording meaning cannot be replaced by a different draft.
	actual := e
	actual.Request.Group = cloneWriteGroup(group)
	if group.isWideProfile() || group.RecordingProfile == nil {
		actual.Request.Group.Destination.TableName = ""
	}
	_, _, actualDigest, err := normalizeSuccessorEnvelope(actual)
	if err != nil || actualDigest != digest {
		return nil, false, ErrWriteGroupLineageCorrupt
	}
	return &WriteGroupSuccessorPreparation{WorkspaceRevision: record.DatabaseSetupRevision, LineageID: meta.LineageID, LineageRevision: meta.Revision, Generation: successor.Generation, PredecessorGroupID: id, OperationID: e.Request.OperationID, IntentDigest: digest, Group: group}, true, nil
}

// successorReservationInTx authorizes one new preparation stage with CURRENT
// CAS, independently of the frozen original reservation replay identity.
func (s *WriteGroupService) successorReservationInTx(ctx context.Context, tx *sql.Tx, record *Record, id string, stage WriteGroupSuccessorStageRequest) (*WriteGroupSuccessorPreparation, error) {
	e, review, digest, err := normalizeSuccessorEnvelope(stage.Envelope)
	if err != nil {
		return nil, err
	}
	if review != stage.ReviewDigest || id != e.PredecessorGroupID {
		return nil, recreateReservationConflict("review_changed")
	}
	// Current-stage CAS is checked before inspecting immutable membership. A
	// stale client revision is a conflict; a caller that names the current but
	// inconsistent registry still reaches the corruption checks below.
	metadata, known, err := s.repo.loadWriteGroupLineage(ctx, tx, record.ID, id)
	if err != nil {
		return nil, err
	}
	if known && stage.Current.LineageRevision != metadata.Revision {
		return nil, recreateReservationConflict("preparation_revision_changed")
	}
	saved, found, err := s.pendingSuccessorInTx(ctx, tx, record, id, e, digest)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, recreateReservationConflict("reservation_missing")
	}
	cas := stage.Current
	if cas.WorkspaceRevision != record.DatabaseSetupRevision || cas.LineageRevision != saved.LineageRevision || cas.GroupID != saved.Group.ID || cas.GroupRevision != saved.Group.Revision || cas.TargetConnectorRevision != saved.Group.Destination.ConnectorRevision {
		return nil, recreateReservationConflict("preparation_revision_changed")
	}
	pred, err := s.repo.GetInTx(ctx, tx, record.ID, id)
	if err != nil {
		return nil, err
	}
	if pred.Revision != cas.PredecessorRevision || pred.Destination.ConnectorRevision != cas.PredecessorConnectorRevision || pred.AppliedRevision != e.SourceAppliedRevision {
		return nil, recreateReservationConflict("predecessor_changed")
	}
	source, err := s.repo.recreateStructuralSourceInTx(ctx, tx, pred)
	if err != nil {
		return nil, err
	}
	sourceDigest, err := successorDomainDigest("go-gateway/successor-source-definition/v1", source)
	if err != nil {
		return nil, err
	}
	if sourceDigest != e.SourceDefinitionDigest {
		return nil, recreateReservationConflict("source_definition_changed")
	}
	if pred.AppliedRevision != "" {
		if s.backlogGuard == nil {
			return nil, ErrWriteGroupLifecycleBlocked
		}
		if err := s.backlogGuard.CheckWriteGroupBacklog(ctx, tx, pred); err != nil {
			return nil, err
		}
	}
	candidate := cloneWriteGroup(saved.Group)
	if _, err := s.repo.validateSuccessorStructureInTx(ctx, tx, pred, candidate); err != nil {
		return nil, err
	}
	return saved, nil
}
