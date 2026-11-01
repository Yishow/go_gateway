package workspace

import (
	"context"
	"database/sql"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
)

// validateSuccessorStructureInTx is local eligibility, NEVER target-family proof.
// The source family is used only to ask whether a same-family transition could
// be supported. Public reserve must independently discharge the required check.
func (r *SQLWriteGroupRepository) validateSuccessorStructureInTx(ctx context.Context, tx *sql.Tx, predecessor, candidate *WriteGroup) (bool, error) {
	if err := r.validateSuccessorLocalDefinitionInTx(ctx, tx, predecessor, candidate); err != nil {
		return false, err
	}
	var kind string
	if err := tx.QueryRowContext(ctx, r.query(`SELECT kind FROM database_connectors WHERE id=$1`), candidate.Destination.ConnectorID).Scan(&kind); err != nil {
		return false, err
	}
	if strings.ToLower(strings.TrimSpace(kind)) != writeGroupConnectorKindMySQL {
		return false, r.validateRecreatedDraftIdentityInTx(ctx, tx, predecessor, candidate)
	}
	if predecessor.RecordingProfile == nil {
		return false, transitionRejection("source-engine-evidence-unavailable")
	}
	source, err := r.recreatedProfileIdentityInTx(ctx, tx, predecessor, candidate)
	if err != nil {
		return false, err
	}
	family := writeGroupConnectorKindMySQL // Provisional structural check only; both families require fresh probes.
	if predecessor.AppliedRevision == "" {
		if predecessor.Status != WriteGroupStatusDeleted {
			return false, transitionRejection("source-engine-evidence-unavailable")
		}
		meta, known, e := r.loadWriteGroupLineage(ctx, tx, predecessor.WorkspaceID, predecessor.ID)
		if e != nil {
			return false, e
		}
		if !known {
			return false, ErrWriteGroupLineageCorrupt
		}
		if e := r.requireSuccessorNoHistoryInTx(ctx, tx, predecessor.WorkspaceID, meta.LineageID, predecessor.ID); e != nil {
			return false, e
		}
		var sourceKind string
		if e := tx.QueryRowContext(ctx, r.query(`SELECT kind FROM database_connectors WHERE id=$1`), source.Destination.ConnectorID).Scan(&sourceKind); e != nil {
			return false, e
		}
		if strings.ToLower(strings.TrimSpace(sourceKind)) != writeGroupConnectorKindMySQL {
			return false, transitionRejection("cross-engine-unsupported")
		}
		if e := validateLiveWriteGroupConnector(ctx, tx, r, source.Destination.ConnectorID, source.Destination.ConnectorRevision); e != nil {
			return false, e
		}
	} else {
		family, err = r.recreateSourceEngineInTx(ctx, tx, source)
		if err != nil {
			return false, err
		}
	}
	if family != writeGroupConnectorKindMySQL && family != writeGroupConnectorFamilyMariaDB {
		return false, transitionRejection("cross-engine-unsupported")
	}
	tags, err := loadWriteGroupTagTypes(ctx, tx, r.query, candidate)
	if err != nil {
		return false, err
	}
	if err := ValidateWriteGroupTransitionStructure(source, candidate, family, family, tags); err != nil {
		return false, err
	}
	if longStorageScopeOf(source) == longStorageScopeOf(candidate) {
		return false, transitionRejection("destination-unchanged")
	}
	return true, nil
}

// successorMySQLReserveEvidence is minted after a fresh probe and full in-tx
// review. It is private, ephemeral, and bound to exactly this transaction/intent.
type successorMySQLReserveEvidence struct {
	tx               *sql.Tx
	digest           string
	workspaceID      string
	capability       dbtarget.MySQLRecordingCapability
	sourceCapability *dbtarget.MySQLRecordingCapability
}

type successorMySQLCapabilities struct {
	Source, Target *dbtarget.MySQLRecordingCapability
}

func (s *WriteGroupService) probeSuccessorMySQL(ctx context.Context, id string, envelope WriteGroupSuccessorEnvelope, review string) (*successorMySQLCapabilities, error) {
	checked, err := s.SuccessorReview(ctx, id, envelope.Request)
	if err != nil {
		return nil, err
	}
	if checked.ReviewDigest != review || len(checked.Blockers) != 0 {
		return nil, recreateReservationConflict("review_changed")
	}
	if len(checked.RequiredChecks) == 0 {
		return nil, nil
	}
	if !successorMySQLRequiredChecks(checked.RequiredChecks, envelope.SourceAppliedRevision == "") {
		return nil, ErrWriteGroupValidation
	}
	probe, ok := s.managedTableInspector.(ManagedWriteGroupMySQLProbeInspector)
	if !ok {
		return nil, transitionRejection("target-family-evidence-unavailable")
	}
	result := &successorMySQLCapabilities{}
	if envelope.SourceAppliedRevision == "" {
		pred, e := s.repo.Get(ctx, envelope.Request.WorkspaceID, id)
		if e != nil {
			return nil, e
		}
		if pred.Revision != envelope.Request.ExpectedPredecessorRevision || pred.AppliedRevision != "" || pred.Destination.ConnectorRevision != envelope.Request.ExpectedPredecessorConnectorRevision {
			return nil, recreateReservationConflict("predecessor_changed")
		}
		source, e := probe.ProbeMySQLRecordingAtRevision(ctx, pred.Destination.ConnectorID, envelope.Request.ExpectedPredecessorConnectorRevision)
		if e != nil || source == nil {
			return nil, transitionRejection("source-family-evidence-unavailable")
		}
		capabilityCopy := *source
		capabilityCopy.Reasons = append([]string(nil), source.Reasons...)
		result.Source = &capabilityCopy
	}
	capability, err := probe.ProbeMySQLRecordingAtRevision(ctx, envelope.Request.Group.Destination.ConnectorID, envelope.Request.ExpectedTargetConnectorRevision)
	if err != nil {
		return nil, transitionRejection("target-family-evidence-unavailable")
	}
	if capability == nil {
		return nil, transitionRejection("target-family-evidence-unavailable")
	}
	// Copy the facts; never retain an inspector-owned mutable pointer as authority.
	capabilityCopy := *capability
	capabilityCopy.Reasons = append([]string(nil), capability.Reasons...)
	result.Target = &capabilityCopy
	return result, nil
}

func (r *SQLWriteGroupRepository) verifySuccessorMySQLInTx(ctx context.Context, tx *sql.Tx, pred, candidate *WriteGroup, targetCapability, sourceCap *dbtarget.MySQLRecordingCapability) error {
	required, err := r.validateSuccessorStructureInTx(ctx, tx, pred, candidate)
	if err != nil {
		return err
	}
	if !required || targetCapability == nil {
		return transitionRejection("target-family-evidence-unavailable")
	}
	if !targetCapability.Connected || targetCapability.MetadataStatus != writeGroupMySQLMetadataVerified || !targetCapability.TransactionalEngine || !targetCapability.SupportsManagedLong || targetCapability.CapabilityRevision == "" || !dbtarget.MySQLRecordingVersionVerified(targetCapability.ServerFamily, targetCapability.ServerVersion) ||
		targetCapability.ConnectorID != candidate.Destination.ConnectorID || targetCapability.ConnectorRevision != candidate.Destination.ConnectorRevision || targetCapability.Database != candidate.Destination.Database ||
		(candidate.StorageMode == WriteGroupStorageModeNativeMonthly && !targetCapability.SupportsNativePartition) {
		return transitionRejection("target-family-evidence-unavailable")
	}
	if err := validateLiveWriteGroupConnector(ctx, tx, r, targetCapability.ConnectorID, targetCapability.ConnectorRevision); err != nil {
		return err
	}
	source, err := r.recreateStructuralSourceInTx(ctx, tx, pred)
	if err != nil {
		return err
	}
	var family string
	if pred.AppliedRevision == "" {
		if !successorMySQLCapabilityMatches(sourceCap, source) {
			return transitionRejection("source-family-evidence-unavailable")
		}
		if err := validateLiveWriteGroupConnector(ctx, tx, r, sourceCap.ConnectorID, sourceCap.ConnectorRevision); err != nil {
			return err
		}
		family = sourceCap.ServerFamily
	} else {
		if sourceCap != nil {
			return ErrWriteGroupValidation
		}
		family, err = r.recreateSourceEngineInTx(ctx, tx, source)
		if err != nil {
			return err
		}
	}
	tags, err := loadWriteGroupTagTypes(ctx, tx, r.query, candidate)
	if err != nil {
		return err
	}
	return ValidateWriteGroupTransitionStructure(source, candidate, family, targetCapability.ServerFamily, tags)
}

func successorMySQLRequiredChecks(checks []WriteGroupSuccessorBlocker, first bool) bool {
	expected := 1
	if first {
		expected = 2
	}
	if len(checks) != expected {
		return false
	}
	source, target := false, false
	for _, check := range checks {
		switch check.ReasonCode {
		case "source-family-check-required":
			source = true
		case "target-family-check-required":
			target = true
		default:
			return false
		}
	}
	return target && source == first
}

func successorMySQLCapabilityMatches(targetCapability *dbtarget.MySQLRecordingCapability, g *WriteGroup) bool {
	return targetCapability != nil && g != nil && targetCapability.Connected && targetCapability.MetadataStatus == writeGroupMySQLMetadataVerified && targetCapability.TransactionalEngine && targetCapability.SupportsManagedLong && targetCapability.CapabilityRevision != "" && dbtarget.MySQLRecordingVersionVerified(targetCapability.ServerFamily, targetCapability.ServerVersion) && targetCapability.ConnectorID == g.Destination.ConnectorID && targetCapability.ConnectorRevision == g.Destination.ConnectorRevision && targetCapability.Database == g.Destination.Database && (g.StorageMode != WriteGroupStorageModeNativeMonthly || targetCapability.SupportsNativePartition)
}
