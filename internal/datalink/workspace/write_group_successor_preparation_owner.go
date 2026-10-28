package workspace

import (
	"context"
	"database/sql"
	"fmt"
	"go-gateway/internal/datalink/measurement"
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
		if err := r.validate(ctx, tx, candidate); err != nil {
			return err
		}
		return nil
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
	for rows.Next() {
		var owner string
		var source [3]string
		if err := rows.Scan(&owner, &source[0], &source[1], &source[2]); err != nil {
			return err
		}
		if wanted[source] && owner != predecessor.ID {
			return fmt.Errorf("%w: another recording owns a successor source", ErrWriteGroupCircuitProfileConflict)
		}
	}
	return rows.Err()
}
