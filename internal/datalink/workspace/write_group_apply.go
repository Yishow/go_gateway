package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/measurement"
)

// applySchedule picks the instant a newly applied revision takes effect.
type applySchedule func(candidate *WriteGroup, active *WriteGroupAppliedSnapshot) (time.Time, error)

// nextAlignedCutover is the ordinary schedule: the next UTC boundary that both
// the old and the new interval share.
func nextAlignedCutover(candidate *WriteGroup, active *WriteGroupAppliedSnapshot) (time.Time, error) {
	if isMixedWriteGroup(candidate) || (active != nil && isMixedWriteGroup(active.Group)) {
		return mixedApplyCutover(candidate, active)
	}
	oldIntervalSeconds := 0
	if active != nil {
		oldIntervalSeconds = active.Group.RowPolicy.IntervalSeconds
	}
	return nextWriteGroupCutover(candidate.UpdatedAt, oldIntervalSeconds, candidate.RowPolicy.IntervalSeconds)
}

// applyGroupInTx applies one saved draft inside the caller's setup transaction:
// the revision checks, the immutable version, the runtime descriptor and the
// writer-ownership barrier. Apply and a reviewed cutover both run it, so a
// cutover is never a second, competing way to activate a group.
//
// laterThanActive is set by a cutover, which has already checked that its
// instant is after the active revision takes effect; the guard against a second
// pending revision then has nothing left to protect.
func (s *WriteGroupService) applyGroupInTx(ctx context.Context, tx *sql.Tx, record *Record, id string, mutation WriteGroupMutation,
	readiness *WriteGroupReadiness, schedule applySchedule, laterThanActive bool,
) (*WriteGroup, error) {
	existing, err := s.repo.GetInTx(ctx, tx, record.ID, id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.requireUnreservedWriteGroupInTx(ctx, tx, record.ID, id); err != nil {
		return nil, err
	}
	if err := s.repo.requireCurrentCanonicalOwnerInTx(ctx, tx, record.ID, id); err != nil {
		return nil, err
	}
	if existing.Revision != mutation.ExpectedGroupRevision {
		return nil, writeGroupRevisionConflict("write group revision is stale")
	}
	if existing.Status == WriteGroupStatusDeleted {
		return nil, ErrWriteGroupLifecycleBlocked
	}
	if existing.Destination.ConnectorRevision != mutation.ExpectedConnectorRevision {
		return nil, writeGroupRevisionConflict("write group connector revision is stale")
	}
	if err := validateLiveWriteGroupConnector(ctx, tx, s.repo, existing.Destination.ConnectorID, mutation.ExpectedConnectorRevision); err != nil {
		return nil, err
	}
	candidate := cloneWriteGroup(existing)
	definitions := make(map[string]*measurement.MeasurementDefinition)
	if err := s.repo.validateWithMeasurements(ctx, tx, candidate, definitions); err != nil {
		return nil, err
	}
	for i := range candidate.Members {
		if !sameExactRecordingSource(candidate.Members[i].ExactRecordingSource, existing.Members[i].ExactRecordingSource) {
			return nil, writeGroupRevisionConflict("write group exact recording source changed during verification")
		}
	}
	return s.applyValidatedGroupInTx(ctx, tx, record, existing, candidate, mutation, readiness, definitions, schedule, laterThanActive)
}

// applyValidatedGroupInTx is shared only by ordinary validation and the private
// successor transaction after exact reservation/source/owner validation.
func (s *WriteGroupService) applyValidatedGroupInTx(ctx context.Context, tx *sql.Tx, record *Record, existing, candidate *WriteGroup, mutation WriteGroupMutation, readiness *WriteGroupReadiness, definitions map[string]*measurement.MeasurementDefinition, schedule applySchedule, laterThanActive bool) (*WriteGroup, error) {
	proofs, err := s.verifyCounterBindingsInTx(ctx, tx, candidate, readiness, definitions)
	if err != nil {
		return nil, err
	}
	candidate.Status = WriteGroupStatusReady
	candidate.UpdatedAt = s.clockNow()
	activeSnapshot, snapshotErr := s.repo.versionByRevisionInTx(ctx, tx, record.ID, existing.ID, existing.AppliedRevision)
	if snapshotErr != nil && !errors.Is(snapshotErr, ErrWriteGroupAppliedRevisionUnavailable) {
		return nil, snapshotErr
	}
	// A retained applied identity without its immutable snapshot is not a new
	// recording. In particular, disabling/resuming must never reset its start.
	if isMixedWriteGroup(candidate) && activeSnapshot == nil {
		if existing.AppliedRevision != "" {
			return nil, ErrWriteGroupAppliedRevisionUnavailable
		}
		var historical bool
		err := tx.QueryRowContext(ctx, s.repo.query(`SELECT EXISTS(SELECT 1 FROM write_group_versions WHERE group_id=$1) OR EXISTS(SELECT 1 FROM wg_runtime_versions WHERE workspace_id=$2 AND group_id=$3)`), existing.ID, record.ID, existing.ID).Scan(&historical)
		if err != nil {
			return nil, err
		}
		if historical {
			return nil, ErrWriteGroupAppliedRevisionUnavailable
		}
	}
	reuseApplied := existing.Status != WriteGroupStatusDisabled && activeSnapshot != nil && sameWriteGroupAppliedSemantics(candidate, activeSnapshot.Group)
	if reuseApplied {
		reuseApplied, err = s.sameAppliedCounterMeaningInTx(ctx, tx, activeSnapshot, proofs)
		if err != nil {
			var proof *CounterBindingProofError
			if !errors.As(err, &proof) {
				return nil, err
			}
			if candidate.Revision == activeSnapshot.AppliedRevision {
				return nil, errors.Join(ErrWriteGroupApplyNotReady, err)
			}
			// A distinct reviewed owner does not repair or reinterpret unproven history.
			reuseApplied = false
		}
	}
	if reuseApplied {
		// A display-only rename must not restart the applied writer or create
		// a second bucket snapshot. Keep the prior active revision while
		// explicitly moving the saved draft back to ready. A disabled group
		// follows the normal new-version path so re-enabling cannot reuse its
		// retired intake boundary.
		candidate.AppliedRevision = existing.AppliedRevision
		activeSnapshot.RuntimeLayout = readiness.runtimeLayout
		if err := s.saveHistoricalRuntimeVersionInTx(ctx, tx, activeSnapshot); err != nil {
			return nil, err
		}
		if err := s.repo.setAppliedInTx(ctx, tx, candidate); err != nil {
			return nil, err
		}
		if err := projectWriteGroup(record, candidate); err != nil {
			return nil, err
		}
		return cloneWriteGroup(candidate), nil
	}
	if activeSnapshot != nil && candidate.Revision == activeSnapshot.AppliedRevision {
		return nil, writeGroupRevisionConflict("changed counter meaning requires a fresh saved write-group revision")
	}
	if !laterThanActive && activeSnapshot != nil && activeSnapshot.EffectiveAt.After(candidate.UpdatedAt) {
		// Keep one pending semantic cutover so an older scheduled snapshot
		// cannot become effective again after a newer policy was applied.
		return nil, fmt.Errorf("write-group cutover is still pending: %w", ErrWriteGroupLifecycleBlocked)
	}
	candidate.AppliedRevision = candidate.Revision
	effectiveAt, err := schedule(candidate, activeSnapshot)
	if err != nil {
		return nil, fmt.Errorf("schedule write-group cutover: %w", err)
	}
	if err := s.repo.insertVersionInTx(ctx, tx, candidate, effectiveAt); err != nil {
		return nil, err
	}
	if err := s.saveRuntimeVersionInTx(ctx, tx, &WriteGroupAppliedSnapshot{
		WorkspaceID: record.ID, GroupID: candidate.ID, GroupRevision: candidate.Revision,
		AppliedRevision: candidate.AppliedRevision, EffectiveAt: effectiveAt.UTC(),
		Group: candidate, RuntimeLayout: readiness.runtimeLayout, CounterBindings: proofs,
	}); err != nil {
		return nil, err
	}
	if err := s.repo.setAppliedInTx(ctx, tx, candidate); err != nil {
		return nil, err
	}
	if err := projectWriteGroup(record, candidate); err != nil {
		return nil, err
	}
	if err := s.activateWriterOwnershipInTx(ctx, tx, WriterOwnershipActivationRequest{
		WorkspaceID:                   record.ID,
		GroupID:                       candidate.ID,
		ExpectedDatabaseSetupRevision: mutation.ExpectedWorkspaceRevision,
		ExpectedGroupRevision:         mutation.ExpectedGroupRevision,
		ExpectedConnectorRevision:     mutation.ExpectedConnectorRevision,
		PreviousAppliedRevision:       existing.AppliedRevision,
		AppliedRevision:               candidate.AppliedRevision,
		EffectiveAt:                   effectiveAt.UTC(),
	}); err != nil {
		return nil, err
	}
	return cloneWriteGroup(candidate), nil
}
