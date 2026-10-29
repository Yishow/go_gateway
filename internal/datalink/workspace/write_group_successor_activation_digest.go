package workspace

import (
	"context"
	"database/sql"
	"reflect"
	"strings"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/recordingplan"
)

const successorActivationVersion = "successor-activation-review-v2"
const successorActivationLegacyVersion = "successor-activation-review-v1"
const successorActivationCutover = "cutover"
const successorActivationFirstStart = "first_start"

func successorActivationKind(e WriteGroupSuccessorActivationEnvelope) string {
	if e.Version == successorActivationLegacyVersion {
		return successorActivationCutover
	}
	return e.ActivationKind
}

// normalizeSuccessorActivationEnvelope is state-free for exact replay after the
// reservation has been consumed. Mutation authority still needs transaction checks.
func normalizeSuccessorActivationEnvelope(e WriteGroupSuccessorActivationEnvelope) (WriteGroupSuccessorActivationEnvelope, string, error) {
	if e.Request.Boundary.IsZero() || (e.Version != successorActivationVersion && e.Version != successorActivationLegacyVersion) || (e.Version == successorActivationLegacyVersion && e.ActivationKind != "") || (e.Version == successorActivationVersion && e.ActivationKind != successorActivationCutover && e.ActivationKind != successorActivationFirstStart) {
		return e, "", ErrWriteGroupValidation
	}
	original, review, _, err := normalizeSuccessorEnvelope(e.Request.Stage.Envelope)
	if err != nil {
		return e, "", err
	}
	if review != e.Request.Stage.ReviewDigest {
		return e, "", recreateReservationConflict("review_changed")
	}
	e.Request.Stage.Envelope = original
	e.Request.Boundary = e.Request.Boundary.UTC()
	c := &e.Request.Stage.Current
	for _, v := range []*string{&c.WorkspaceRevision, &c.PredecessorRevision, &c.PredecessorConnectorRevision, &c.GroupID, &c.GroupRevision, &c.TargetConnectorRevision} {
		if strings.TrimSpace(*v) == "" || strings.TrimSpace(*v) != *v {
			return e, "", ErrWriteGroupValidation
		}
	}
	if c.LineageRevision < 1 {
		return e, "", ErrWriteGroupValidation
	}
	domain := "go-gateway/successor-activation-review/v2"
	if e.Version == successorActivationLegacyVersion {
		domain = "go-gateway/successor-activation-review/v1"
	}
	digest, err := successorDomainDigest(domain, e)
	return e, digest, err
}
func successorActivationDigest(e WriteGroupSuccessorActivationEnvelope) (string, error) {
	_, d, err := normalizeSuccessorActivationEnvelope(e)
	return d, err
}

func (s *WriteGroupService) activationSchemaProofInTx(ctx context.Context, tx *sql.Tx, snap *successorActivationSnapshot) error {
	group := snap.reserved.Group
	var layout *recordingplan.GroupSchemaLayout
	var err error
	if snap.ready.connectorKind == writeGroupConnectorKindMySQL {
		token, e := recordingplan.SucceededSchemaTokenInTx(ctx, tx, group.WorkspaceID, group.Destination.SchemaRevision, group.Destination.SchemaDigest)
		if e != nil {
			return e
		}
		family, e := s.successorSourceEngineInTx(ctx, tx, snap.predecessor, snap.reserved)
		if e != nil {
			return e
		}
		layout, err = successorMySQLProofLayout(group, family, token.GroupLayout)
	} else if snap.ready.connectorKind == writeGroupConnectorKindSQLServer {
		token, e := recordingplan.SucceededSchemaTokenInTx(ctx, tx, group.WorkspaceID, group.Destination.SchemaRevision, group.Destination.SchemaDigest)
		if e != nil {
			return e
		}
		layout, err = successorSQLServerProofLayoutInTx(ctx, tx, group, token.GroupLayout)
	} else {
		layout, err = managedSchemaLayout(group, snap.ready.tagTypes, snap.ready.connectorKind)
	}
	if err != nil {
		return err
	}
	digest, err := successorSchemaSourceDigest(&successorSchemaSnapshot{reserved: snap.reserved, types: snap.ready.tagTypes, dialect: snap.ready.connectorKind})
	if err != nil {
		return err
	}
	scope := recordingplan.SchemaPreviewScope{WorkspaceID: group.WorkspaceID, WorkspaceRevision: snap.reserved.WorkspaceRevision, PlanID: group.ID, PlanRevision: group.Revision, ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision, Dialect: snap.ready.connectorKind, Database: group.Destination.Database, Schema: group.Destination.TableSchema, GroupLayout: layout, SourceDigest: digest, SchemaRevision: group.Destination.SchemaRevision, SchemaDigest: group.Destination.SchemaDigest, RecreateMissing: group.WritePolicy.RecreateMissing}
	_, err = recordingplan.VerifiedSchemaProofInTx(ctx, tx, scope, group.Destination.SchemaRevision, group.Destination.SchemaDigest)
	if err == nil && snap.ready.connectorKind == writeGroupConnectorKindMySQL {
		snap.ready.confirmedMySQLLayout = layout
	}
	if err == nil && snap.ready.connectorKind == writeGroupConnectorKindSQLServer {
		snap.ready.confirmedSQLServerLayout = layout
	}
	return err
}

func (s *WriteGroupService) activationCounterDigestInTx(ctx context.Context, tx *sql.Tx, snap *successorActivationSnapshot) (string, error) {
	target := snap.ready.counterBindings
	var original CounterBindingProofs
	if snap.source != nil && len(target) > 0 {
		var raw string
		if err := tx.QueryRowContext(ctx, s.repo.query("SELECT payload FROM wg_runtime_versions WHERE workspace_id=$1 AND group_id=$2 AND group_revision=$3"), snap.source.WorkspaceID, snap.source.GroupID, snap.source.AppliedRevision).Scan(&raw); err != nil {
			return "", err
		}
		var err error
		original, err = counterproof.ResolveSnapshot([]byte(raw), counterproof.SnapshotScope{WorkspaceID: snap.source.WorkspaceID, GroupID: snap.source.GroupID, AppliedRevision: snap.source.AppliedRevision})
		if err != nil {
			return "", err
		}
		comparable := make(CounterBindingProofs, len(original))
		for _, proof := range original {
			facts := proof.CounterBindingFacts
			facts.Member.GroupID = snap.reserved.Group.ID
			facts.Member.AppliedRevision = snap.reserved.Group.Revision
			p, err := NewCounterBindingProof(facts)
			if err != nil {
				return "", err
			}
			comparable[p.Member.Key()] = p
		}
		if !reflect.DeepEqual(comparable, target) {
			return "", ErrWriteGroupApplyNotReady
		}
	}
	return successorDomainDigest("go-gateway/successor-activation-counter/v1", struct{ Source, Target CounterBindingProofs }{original, target})
}

// Token metadata selects capability facts only. Group ownership, source family,
// supported mode and canonical shape remain independently derived local facts.
func successorMySQLProofLayout(group *WriteGroup, family string, recorded *recordingplan.GroupSchemaLayout) (*recordingplan.GroupSchemaLayout, error) {
	if recorded == nil || recorded.MySQL == nil || recorded.MySQL.ServerFamily != family ||
		(family != "mysql" && family != "mariadb") || !successorActivationEngineSupported(group, family) ||
		(group.StorageMode == WriteGroupStorageModeNativeMonthly) != (recorded.MySQL.NativeMonthly != nil) {
		return nil, ErrWriteGroupValidation
	}
	table, owner, ok := LongStorageIdentity(group)
	if !ok || table != group.Destination.TableName {
		return nil, ErrWriteGroupValidation
	}
	layout, err := recordingplan.ManagedMySQLLongGroupLayout(table, owner, *recorded.MySQL)
	if err != nil {
		return nil, err
	}
	return &layout, nil
}

// MySQL's final transaction validates the same immutable matrix using the fully
// authenticated, reservation-bound schema token, never the driver label.
func (s *WriteGroupService) validateSuccessorActivationDefinitionInTx(ctx context.Context, tx *sql.Tx, snap *successorActivationSnapshot, candidate *WriteGroup) error {
	if snap.ready.connectorKind != writeGroupConnectorKindMySQL {
		return s.repo.validateSuccessorDefinitionInTx(ctx, tx, snap.predecessor, candidate)
	}
	if _, err := s.repo.validateSuccessorStructureInTx(ctx, tx, snap.predecessor, candidate); err != nil {
		return err
	}
	if err := s.activationSchemaProofInTx(ctx, tx, snap); err != nil {
		return err
	}
	family, err := s.successorSourceEngineInTx(ctx, tx, snap.predecessor, snap.reserved)
	if err != nil {
		return err
	}
	source, err := s.repo.recreateStructuralSourceInTx(ctx, tx, snap.predecessor)
	if err != nil {
		return err
	}
	return ValidateWriteGroupTransitionStructure(source, candidate, family, snap.ready.confirmedMySQLLayout.MySQL.ServerFamily, snap.ready.tagTypes)
}
