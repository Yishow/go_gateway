package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/measurement"
	"slices"
)

// validateSuccessorDefinitionInTx validates a candidate without creating any
// authority outside this transaction. Only the exact canonical predecessor may
// overlap. Ordinary validation never calls this path and still rejects it.
func (r *SQLWriteGroupRepository) validateSuccessorDefinitionInTx(ctx context.Context, tx *sql.Tx, predecessor, candidate *WriteGroup) error {
	if err := r.validateSuccessorLocalDefinitionInTx(ctx, tx, predecessor, candidate); err != nil {
		return err
	}
	return r.validateRecreatedDraftIdentityInTx(ctx, tx, predecessor, candidate)
}

// validateSuccessorLocalDefinitionInTx validates ownership and the local definition only.
func (r *SQLWriteGroupRepository) validateSuccessorLocalDefinitionInTx(ctx context.Context, tx *sql.Tx, predecessor, candidate *WriteGroup) error {
	if tx == nil || predecessor == nil || candidate == nil || predecessor.WorkspaceID != candidate.WorkspaceID {
		return ErrWriteGroupValidation
	}
	defaultRecreateMissing(candidate)
	if candidate.Destination.TableName == "" && !candidate.isLongLayout() && candidate.Destination.StorageStrategy == WriteGroupStorageStrategyManaged {
		// Validation-only identifier: never an existing owner or a reserved table.
		candidate.Destination.TableName = managedGroupTableName("successor-review")
		defer func() { candidate.Destination.TableName = "" }()
	}
	metadata, found, err := r.loadWriteGroupLineage(ctx, tx, predecessor.WorkspaceID, predecessor.ID)
	if err != nil {
		return err
	}
	if !found || metadata.CurrentGroupID != predecessor.ID {
		return ErrWriteGroupLineageCorrupt
	}
	if predecessor.Status == WriteGroupStatusDeleted {
		if predecessor.RecordingProfile == nil {
			if err := r.validateRecreatedBasicInTx(ctx, tx, predecessor, candidate); err != nil {
				return err
			}
		}
		return r.validate(ctx, tx, candidate)
	}
	if predecessor.RecordingProfile == nil || predecessor.AppliedRevision == "" {
		return ErrWriteGroupLifecycleBlocked
	}
	if predecessor.Status != WriteGroupStatusReady && predecessor.Status != WriteGroupStatusRunning && predecessor.Status != WriteGroupStatusDisabled {
		return ErrWriteGroupLifecycleBlocked
	}
	if err := r.validateDefinitionWithMeasurements(ctx, tx, candidate, make(map[string]*measurement.MeasurementDefinition)); err != nil {
		return err
	}
	return r.checkSuccessorProfileOwnersInTx(ctx, tx, predecessor, candidate)
}

func (r *SQLWriteGroupRepository) checkSuccessorProfileOwnersInTx(ctx context.Context, tx *sql.Tx, predecessor, candidate *WriteGroup) error {
	rows, err := tx.QueryContext(ctx, r.query(`SELECT g.id,m.device_id,m.point_id,m.tag_id FROM write_groups g JOIN write_group_members m ON m.group_id=g.id
 WHERE g.workspace_id=$1 AND g.id<>$2 AND g.status<>$3 AND g.profile_kind=$4
 AND g.destination_connector_id=$5 AND g.destination_database=$6 AND g.destination_table_schema=$7`), candidate.WorkspaceID, candidate.ID, WriteGroupStatusDeleted, string(candidate.RecordingProfile.Kind), candidate.Destination.ConnectorID, candidate.Destination.Database, candidate.Destination.TableSchema)
	if err != nil {
		return err
	}
	defer rows.Close()
	wanted := make(map[[3]string]bool, len(candidate.Members))
	for _, m := range candidate.Members {
		wanted[[3]string{m.DeviceID, m.PointID, m.TagID}] = true
	}
	owners := map[string]bool{}
	for rows.Next() {
		var owner string
		var source [3]string
		if err := rows.Scan(&owner, &source[0], &source[1], &source[2]); err != nil {
			return err
		}
		if wanted[source] && owner != predecessor.ID {
			owners[owner] = true
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	var record *mixedActiveRecordIdentity
	if isMixedWriteGroup(predecessor) && isMixedWriteGroup(candidate) && predecessor.RecordID == candidate.RecordID {
		record, err = r.mixedRecordIdentityInTx(ctx, tx, predecessor)
		if err != nil {
			return err
		}
	}
	for owner := range owners {
		if record != nil {
			group, err := r.GetInTx(ctx, tx, candidate.WorkspaceID, owner)
			if err != nil {
				return err
			}
			retired, err := r.retiredMixedSuccessorOwnerInTx(ctx, tx, predecessor, group, record)
			if err != nil {
				return err
			}
			if retired {
				continue
			}
			other, err := r.mixedRecordIdentityInTx(ctx, tx, group)
			if err != nil {
				return err
			}
			if independentMixedRecordIdentities(record, other) {
				continue
			}
		}
		return fmt.Errorf("%w: another recording owns a successor source", ErrWriteGroupCircuitProfileConflict)
	}
	return nil
}

// A completed earlier generation is history, not a new source owner. This is
// restricted to the same canonical record and authenticated forward lineage;
// unknown or still-future ends do not qualify for the exception.
func (r *SQLWriteGroupRepository) retiredMixedSuccessorOwnerInTx(ctx context.Context, tx *sql.Tx, predecessor, owner *WriteGroup, record *mixedActiveRecordIdentity) (bool, error) {
	if !isMixedWriteGroup(owner) || owner.RecordID != record.recordID || owner.AppliedRevision == "" {
		return false, nil
	}
	prior, found, err := r.loadWriteGroupLineage(ctx, tx, owner.WorkspaceID, owner.ID)
	if err != nil {
		return false, err
	}
	if !found || prior.LineageID != record.lineageID || prior.CurrentGroupID != predecessor.ID {
		return false, nil
	}
	current, found, err := r.loadWriteGroupLineage(ctx, tx, predecessor.WorkspaceID, predecessor.ID)
	if err != nil {
		return false, err
	}
	if !found || prior.Generation >= current.Generation {
		return false, nil
	}
	end, err := r.revisionEndInTx(ctx, tx, owner.ID, owner.AppliedRevision)
	if err != nil {
		return false, err
	}
	if end.IsZero() || end.After(r.now()) {
		return false, nil
	}
	var eventID string
	err = tx.QueryRowContext(ctx, r.query(`SELECT cutover_id FROM write_group_revision_ends WHERE workspace_id=$1 AND group_id=$2 AND group_revision=$3`), owner.WorkspaceID, owner.ID, owner.AppliedRevision).Scan(&eventID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrWriteGroupLineageCorrupt
	}
	if err != nil {
		return false, err
	}
	event, err := r.cutoverInTx(ctx, tx, owner.WorkspaceID, eventID)
	if err != nil {
		return false, err
	}
	if event == nil || event.Kind != WriteGroupCutoverForward || event.SourceGroupID != owner.ID || event.OwnershipFenceRevision != event.ID || !isSuccessorActivationIntent(event.intent) || !event.Boundary.Equal(end) || len(event.Ended) != 1 || !slices.Contains(event.Ended, CutoverEnd{GroupID: owner.ID, Revision: owner.AppliedRevision}) || len(event.Started) != 1 {
		return false, ErrWriteGroupLineageCorrupt
	}
	next, found, err := r.loadWriteGroupLineage(ctx, tx, owner.WorkspaceID, event.Started[0].GroupID)
	if err != nil {
		return false, err
	}
	if !found || next.LineageID != record.lineageID || next.CurrentGroupID != predecessor.ID || next.Generation != prior.Generation+1 || next.PredecessorGroupID.String != owner.ID {
		return false, ErrWriteGroupLineageCorrupt
	}
	version, err := r.versionByRevisionInTx(ctx, tx, owner.WorkspaceID, event.Started[0].GroupID, event.Started[0].Revision)
	if err != nil {
		return false, err
	}
	if !isMixedWriteGroup(version.Group) || version.Group.RecordID != record.recordID {
		return false, ErrWriteGroupLineageCorrupt
	}
	return true, nil
}
