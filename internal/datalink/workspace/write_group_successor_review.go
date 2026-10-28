package workspace

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

// SuccessorReview reads only local state. It never reserves a group, persists
// a schema token, probes a connector, or authorizes a future mutation.
func (s *WriteGroupService) SuccessorReview(ctx context.Context, id string, req WriteGroupSuccessorReviewRequest) (result *WriteGroupSuccessorReview, err error) {
	if err = s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrWriteGroupValidation
	}
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupServiceUnavailable
	}
	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := setupTx.Rollback(); e != nil && !errors.Is(e, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, e)
		}
	}()
	record, err := setupTx.Get(ctx)
	if err != nil {
		return nil, err
	}
	return s.reviewSuccessorInTx(ctx, setupTx.SQLTx(), record, id, req)
}

func (s *WriteGroupService) reviewSuccessorInTx(ctx context.Context, tx *sql.Tx, record *Record, id string, req WriteGroupSuccessorReviewRequest) (*WriteGroupSuccessorReview, error) {
	req, err := normalizeSuccessorReviewRequest(id, req)
	if err != nil {
		return nil, err
	}
	if record == nil || tx == nil || record.ID != req.WorkspaceID {
		return nil, recreateReservationConflict("scope_changed")
	}
	if record.DatabaseSetupRevision != req.ExpectedWorkspaceRevision {
		return nil, recreateReservationConflict("workspace_changed")
	}
	metadata, found, err := s.repo.loadWriteGroupLineage(ctx, tx, record.ID, id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrWriteGroupLifecycleBlocked
	}
	if metadata.LineageID != req.LineageID || metadata.CurrentGroupID != id || metadata.Revision != req.ExpectedLineageRevision {
		return nil, recreateReservationConflict("lineage_changed")
	}
	predecessor, err := s.repo.GetInTx(ctx, tx, record.ID, id)
	if err != nil {
		return nil, err
	}
	if predecessor.Revision != req.ExpectedPredecessorRevision {
		return nil, recreateReservationConflict("predecessor_changed")
	}
	if predecessor.Destination.ConnectorRevision != req.ExpectedPredecessorConnectorRevision {
		return nil, recreateReservationConflict("source_connector_changed")
	}
	if err := validateLiveWriteGroupConnector(ctx, tx, s.repo, req.Group.Destination.ConnectorID, req.ExpectedTargetConnectorRevision); err != nil {
		return nil, err
	}
	source, err := s.repo.recreateStructuralSourceInTx(ctx, tx, predecessor)
	if err != nil {
		return nil, err
	}
	definitionDigest, err := successorDomainDigest("go-gateway/successor-source-definition/v1", source)
	if err != nil {
		return nil, err
	}
	history, err := s.repo.lifecycleHistoryInTx(ctx, tx, record.ID, id)
	if err != nil {
		return nil, err
	}
	result := &WriteGroupSuccessorReview{History: history, Blockers: []WriteGroupSuccessorBlocker{}, RequiredChecks: []WriteGroupSuccessorBlocker{}}
	block := func(reason string) {
		result.Blockers = append(result.Blockers, WriteGroupSuccessorBlocker{ReasonCode: reason, NextAction: "review_lifecycle"})
	}
	if metadata.ReservedGroupID.Valid {
		block("successor_reserved")
	}
	switch predecessor.Status {
	case WriteGroupStatusDeleted:
	case WriteGroupStatusDraft, WriteGroupStatusReady, WriteGroupStatusRunning, WriteGroupStatusDisabled:
		if predecessor.RecordingProfile == nil {
			block("active_basic_unsupported")
		}
		if predecessor.AppliedRevision == "" {
			block("source_definition_unavailable")
		}
	default:
		return nil, ErrWriteGroupLifecycleBlocked
	}
	if predecessor.AppliedRevision != "" {
		if s.backlogGuard == nil {
			block("ownership_unavailable")
		} else if err := s.backlogGuard.CheckWriteGroupBacklog(ctx, tx, predecessor); err != nil {
			block("ownership_unavailable")
		}
	}
	if len(result.Blockers) == 0 {
		candidate := cloneWriteGroup(req.Group)
		required, validationErr := s.repo.validateSuccessorStructureInTx(ctx, tx, predecessor, candidate)
		if err := validationErr; err != nil {
			var rejected *WriteGroupTransitionError
			switch {
			case errors.As(err, &rejected):
				block(rejected.ReasonCode)
			case errors.Is(err, ErrWriteGroupValidation):
				block("target_definition_invalid")
			default:
				return nil, err
			}
		} else {
			if predecessor.Status != WriteGroupStatusDeleted && source.isLongLayout() && candidate.isLongLayout() && longStorageScopeOf(source) == longStorageScopeOf(candidate) {
				block("destination-unchanged")
			}
			if required && len(result.Blockers) == 0 {
				if predecessor.AppliedRevision == "" {
					result.RequiredChecks = append(result.RequiredChecks, WriteGroupSuccessorBlocker{ReasonCode: "source-family-check-required", NextAction: "confirm_reservation"})
				}
				result.RequiredChecks = append(result.RequiredChecks, WriteGroupSuccessorBlocker{ReasonCode: "target-family-check-required", NextAction: "confirm_reservation"})
			}
			req.Group = candidate
			req, err = normalizeSuccessorReviewRequest(id, req)
			if err != nil {
				return nil, err
			}
		}
	}
	result.Envelope = WriteGroupSuccessorEnvelope{Version: successorReviewVersion, PredecessorGroupID: id, Generation: metadata.Generation, SourceAppliedRevision: predecessor.AppliedRevision, SourceDefinitionDigest: definitionDigest, Request: req}
	result.Envelope, result.ReviewDigest, _, err = normalizeSuccessorEnvelope(result.Envelope)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
