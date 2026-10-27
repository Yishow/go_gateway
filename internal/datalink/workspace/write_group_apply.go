package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"
)

// applySchedule picks the instant a newly applied revision takes effect.
type applySchedule func(candidate *WriteGroup, active *WriteGroupAppliedSnapshot) (time.Time, error)

// nextAlignedCutover is the ordinary schedule: the next UTC boundary that both
// the old and the new interval share.
func nextAlignedCutover(candidate *WriteGroup, active *WriteGroupAppliedSnapshot) (time.Time, error) {
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
	if err := s.repo.validate(ctx, tx, candidate); err != nil {
		return nil, err
	}
	for i := range candidate.Members {
		if !sameExactRecordingSource(candidate.Members[i].ExactRecordingSource, existing.Members[i].ExactRecordingSource) {
			return nil, writeGroupRevisionConflict("write group exact recording source changed during verification")
		}
	}
	if readiness.runtimeLayout != nil {
		currentTypes, err := loadWriteGroupTagTypes(ctx, tx, s.repo.query, candidate)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(currentTypes, readiness.runtimeLayout.TagTypes) {
			return nil, writeGroupRevisionConflict("write group tag types changed during schema verification")
		}
	}
	candidate.Status = WriteGroupStatusReady
	candidate.UpdatedAt = s.clockNow()
	activeSnapshot, snapshotErr := s.repo.versionByRevisionInTx(ctx, tx, record.ID, existing.ID, existing.AppliedRevision)
	if snapshotErr != nil && !errors.Is(snapshotErr, ErrWriteGroupAppliedRevisionUnavailable) {
		return nil, snapshotErr
	}
	if existing.Status != WriteGroupStatusDisabled && activeSnapshot != nil && sameWriteGroupAppliedSemantics(candidate, activeSnapshot.Group) {
		// A display-only rename must not restart the applied writer or create
		// a second bucket snapshot. Keep the prior active revision while
		// explicitly moving the saved draft back to ready. A disabled group
		// follows the normal new-version path so re-enabling cannot reuse its
		// retired intake boundary.
		candidate.AppliedRevision = existing.AppliedRevision
		activeSnapshot.RuntimeLayout = readiness.runtimeLayout
		if err := s.saveRuntimeVersionInTx(ctx, tx, activeSnapshot); err != nil {
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
		Group: candidate, RuntimeLayout: readiness.runtimeLayout,
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
