package workspace

import (
	"context"
	"database/sql"

	"go-gateway/internal/datalink/recordingplan"
)

func decodeSuccessorReserveSourceEvidence(root *recordingplan.SuccessorReserveRecord, e *WriteGroupSuccessorEnvelope, result *WriteGroupSuccessorPreparation) (*successorReserveSourceEvidence, error) {
	var proof successorReserveSourceEvidence
	if root == nil || e == nil || result == nil || decodeSuccessorJSON(root.SourceEvidence, &proof) != nil || proof.Version != "successor-reserve-source/v1" || proof.Source == nil {
		return nil, ErrWriteGroupLineageCorrupt
	}
	source := proof.Source
	digest, err := successorDomainDigest("go-gateway/successor-source-definition/v1", source)
	if err != nil || digest != e.SourceDefinitionDigest || source.ID != root.Scope.PredecessorGroupID || source.WorkspaceID != root.Scope.WorkspaceID || source.AppliedRevision != e.SourceAppliedRevision || (e.SourceAppliedRevision == "" && source.Destination.ConnectorRevision != e.Request.ExpectedPredecessorConnectorRevision) {
		return nil, ErrWriteGroupLineageCorrupt
	}
	if proof.TargetMySQL != nil && !successorMySQLCapabilityMatches(proof.TargetMySQL, result.Group) {
		return nil, ErrWriteGroupLineageCorrupt
	}
	if proof.SourceMySQL != nil {
		if e.SourceAppliedRevision != "" || source.Status != WriteGroupStatusDeleted || !successorMySQLCapabilityMatches(proof.SourceMySQL, source) || proof.TargetMySQL == nil || proof.SourceMySQL.ServerFamily != proof.TargetMySQL.ServerFamily {
			return nil, ErrWriteGroupLineageCorrupt
		}
	}
	binding, err := decodeSuccessorReserveSchemaBinding(root.SchemaSourceBinding, result)
	if err != nil {
		return nil, err
	}
	if binding.Dialect != writeGroupConnectorKindMySQL && (proof.SourceMySQL != nil || proof.TargetMySQL != nil) {
		return nil, ErrWriteGroupLineageCorrupt
	}
	if binding.Dialect == writeGroupConnectorKindMySQL && (proof.TargetMySQL == nil || proof.TargetMySQL.ServerFamily != binding.SourceFamily || e.SourceAppliedRevision == "" && proof.SourceMySQL == nil) {
		return nil, ErrWriteGroupLineageCorrupt
	}
	return &proof, nil
}

// successorSourceEngineInTx only permits private reserve proof for an entirely
// unapplied source. Applied sources keep the existing immutable runtime path.
func (s *WriteGroupService) successorSourceEngineInTx(ctx context.Context, tx *sql.Tx, pred *WriteGroup, reserved *WriteGroupSuccessorPreparation) (string, error) {
	source, err := s.repo.recreateStructuralSourceInTx(ctx, tx, pred)
	if err != nil {
		return "", err
	}
	if pred.AppliedRevision != "" {
		return s.repo.recreateSourceEngineInTx(ctx, tx, source)
	}
	kind, err := s.successorDriverInTx(ctx, tx, source.Destination.ConnectorID)
	if err != nil {
		return "", err
	}
	if kind != writeGroupConnectorKindMySQL {
		return s.repo.recreateSourceEngineInTx(ctx, tx, source)
	}
	root, err := recordingplan.ReadSuccessorReserveInTx(ctx, tx, pred.WorkspaceID, reserved.OperationID)
	if err != nil {
		return "", err
	}
	if root == nil {
		return "", transitionRejection("source-engine-evidence-unavailable")
	}
	envelope, original, err := decodeSuccessorReserveRecord(root)
	if err != nil {
		return "", err
	}
	if original.Group.ID != reserved.Group.ID || original.LineageID != reserved.LineageID || root.Scope.PredecessorGroupID != pred.ID {
		return "", ErrWriteGroupLineageCorrupt
	}
	proof, err := decodeSuccessorReserveSourceEvidence(root, envelope, original)
	if err != nil {
		return "", err
	}
	if !sameSuccessorJSON(proof.Source, source) || proof.SourceMySQL == nil {
		return "", transitionRejection("source-engine-evidence-unavailable")
	}
	if err := s.repo.requireSuccessorNoHistoryInTx(ctx, tx, pred.WorkspaceID, reserved.LineageID, pred.ID); err != nil {
		return "", err
	}
	return proof.SourceMySQL.ServerFamily, nil
}
