package workspace

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"go-gateway/internal/datalink/measurement"
)

type successorActivationSnapshot struct {
	tx             *sql.Tx
	predecessor    *WriteGroup
	source         *WriteGroupAppliedSnapshot
	reserved       *WriteGroupSuccessorPreparation
	ready          *writeGroupReadinessSnapshot
	activationKind string
}

// SuccessorActivationReview reads a coherent local snapshot without target probes.
func (s *WriteGroupService) SuccessorActivationReview(ctx context.Context, id string, req WriteGroupSuccessorActivationRequest) (*WriteGroupSuccessorActivationReview, error) {
	view, _, err := s.readSuccessorActivation(ctx, id, req)
	return view, err
}

func (s *WriteGroupService) readSuccessorActivation(ctx context.Context, id string, req WriteGroupSuccessorActivationRequest) (view *WriteGroupSuccessorActivationReview, snapshot *successorActivationSnapshot, err error) {
	if err = s.validate(); err != nil {
		return nil, nil, err
	}
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, nil, ErrWriteGroupServiceUnavailable
	}
	setup, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if e := setup.Rollback(); e != nil && !errors.Is(e, sql.ErrTxDone) {
			view = nil
			snapshot = nil
			err = errors.Join(err, e)
		}
	}()
	record, err := setup.Get(ctx)
	if err != nil {
		return nil, nil, err
	}
	return s.activationReviewInTx(ctx, setup.SQLTx(), record, strings.TrimSpace(id), req)
}

func (s *WriteGroupService) activationReviewInTx(ctx context.Context, tx *sql.Tx, record *Record, id string, req WriteGroupSuccessorActivationRequest) (*WriteGroupSuccessorActivationReview, *successorActivationSnapshot, error) {
	if tx == nil || record == nil || id == "" || req.Boundary.IsZero() {
		return nil, nil, ErrWriteGroupValidation
	}
	req.Boundary = req.Boundary.UTC()
	reserved, err := s.successorReservationInTx(ctx, tx, record, id, req.Stage)
	if err != nil {
		return nil, nil, err
	}
	pred, err := s.repo.GetInTx(ctx, tx, record.ID, id)
	if err != nil {
		return nil, nil, err
	}
	result := &WriteGroupSuccessorActivationReview{Blockers: []WriteGroupSuccessorBlocker{}}
	block := func(code string) {
		result.Blockers = append(result.Blockers, WriteGroupSuccessorBlocker{ReasonCode: code, NextAction: "review_lifecycle"})
	}
	snap := &successorActivationSnapshot{tx: tx, predecessor: pred, reserved: reserved, activationKind: successorActivationCutover}
	if pred.AppliedRevision == "" {
		snap.activationKind = successorActivationFirstStart
		if err := s.repo.requireSuccessorNoHistoryInTx(ctx, tx, record.ID, reserved.LineageID, id); err != nil {
			return nil, nil, err
		}
	} else {
		snap.source, err = s.repo.VersionByRevisionInTx(ctx, tx, record.ID, id, pred.AppliedRevision)
		if err != nil {
			return nil, nil, err
		}
		if !req.Boundary.After(snap.source.EffectiveAt) {
			block(CutoverIssueBoundaryBeforeActive)
		}
		end, e := s.repo.revisionEndInTx(ctx, tx, id, pred.AppliedRevision)
		if e != nil {
			return nil, nil, e
		}
		if !end.IsZero() {
			block("source_revision_ended")
		}
	}
	if !req.Boundary.After(s.clockNow()) {
		block(CutoverIssueBoundaryNotFuture)
	}
	group := reserved.Group
	if err = s.requireNoOtherSuccessorActiveOwnerInTx(ctx, tx, pred, group); err != nil {
		return nil, nil, err
	}
	kind, err := s.successorDriverInTx(ctx, tx, group.Destination.ConnectorID)
	if err != nil {
		return nil, nil, err
	}
	if !successorActivationEngineSupported(group, kind) {
		block("successor_engine_unavailable")
	}
	if snap.source != nil {
		engine, e := s.repo.recreateSourceEngineInTx(ctx, tx, snap.source.Group)
		if e != nil {
			return nil, nil, e
		}
		if !successorActivationEngineSupported(snap.source.Group, engine) {
			block("successor_engine_unavailable")
		}
	}
	validated := cloneWriteGroup(group)
	definitions := make(map[string]*measurement.MeasurementDefinition)
	if err = s.repo.validateDefinitionWithMeasurements(ctx, tx, validated, definitions); err != nil {
		return nil, nil, err
	}
	types, err := loadWriteGroupTagTypes(ctx, tx, s.repo.query, validated)
	if err != nil {
		return nil, nil, err
	}
	var config string
	if err = tx.QueryRowContext(ctx, s.repo.query("SELECT connection_config FROM database_connectors WHERE id=$1"), group.Destination.ConnectorID).Scan(&config); err != nil {
		return nil, nil, err
	}
	proofs, issues := s.counterBindingEvidenceInTx(ctx, tx, validated, types, definitions)
	snap.ready = &writeGroupReadinessSnapshot{record: record, group: group, validated: validated, tagTypes: types, connectorKind: kind, connectorConfig: config, counterBindings: proofs, counterIssues: issues}
	for _, issue := range issues {
		block(issue.Code)
	}
	if group.Destination.SchemaRevision == "" || group.Destination.SchemaDigest == "" {
		block("schema-confirmation-required")
	} else {
		if err = s.activationSchemaProofInTx(ctx, tx, snap); err != nil {
			block("schema-confirmation-unverified")
		}
	}
	result.History, err = s.repo.lifecycleHistoryInTx(ctx, tx, record.ID, id)
	if err != nil {
		return nil, nil, err
	}
	result.Envelope = WriteGroupSuccessorActivationEnvelope{Version: successorActivationVersion, ActivationKind: snap.activationKind, Request: req, SourceAppliedRevision: pred.AppliedRevision, SchemaRevision: group.Destination.SchemaRevision, SchemaDigest: group.Destination.SchemaDigest}
	if snap.source == nil {
		result.Envelope.SourceDefinitionDigest = req.Stage.Envelope.SourceDefinitionDigest
	}
	if snap.source != nil {
		result.Envelope.SourceDefinitionDigest, err = successorDomainDigest("successor-activation-source/v1", snap.source)
		if err != nil {
			return nil, nil, err
		}
	}
	result.Envelope.TargetDefinitionDigest, err = successorDomainDigest("successor-activation-target/v1", validated)
	if err != nil {
		return nil, nil, err
	}
	result.Envelope.CounterProofDigest, err = s.activationCounterDigestInTx(ctx, tx, snap)
	if err != nil {
		block("counter-evidence-unavailable")
		result.Envelope.CounterProofDigest = ""
	}
	result.ReviewDigest, err = successorActivationDigest(result.Envelope)
	if err != nil {
		return nil, nil, err
	}
	return result, snap, ctx.Err()
}

// The reservation reloader already verifies the immutable source, exact meaning
// and same-engine transition matrix in this transaction. This gate only limits
// activation to adapters with implemented runtime support, on BOTH sides.
// Timescale capability/hypertable proof remains mandatory in readiness; allowing
// its engine here neither substitutes plain storage nor grants intake authority.
func successorActivationEngineSupported(group *WriteGroup, engine string) bool {
	if group == nil {
		return false
	}
	switch engine {
	case writeGroupConnectorKindSQLite:
		return group.StorageMode == "" || group.StorageMode == WriteGroupStorageModePlain
	case writeGroupConnectorKindMySQL, "mariadb", writeGroupConnectorKindSQLServer:
		return managedLongGroup(group) && (group.StorageMode == WriteGroupStorageModePlain || group.StorageMode == WriteGroupStorageModeNativeMonthly)
	case writeGroupConnectorKindPostgres:
		_, supported := transitionStorageMode(group, engine)
		return supported
	default:
		return false
	}
}
