package workspace

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
)

// reserveRecreatedDraftInTx never commits or opens another connection. Its
// caller owns workspace CAS, backlog validation and compatibility projection.
func (r *SQLWriteGroupRepository) reserveRecreatedDraftInTx(ctx context.Context, tx *sql.Tx, req writeGroupRecreateReservationRequest, digest string) (*writeGroupReservedDraft, error) {
	if tx == nil || req.Group == nil || strings.TrimSpace(digest) == "" || strings.TrimSpace(req.OperationID) == "" ||
		strings.TrimSpace(req.LineageID) == "" || req.ExpectedLineageRevision < 1 ||
		req.WorkspaceID == "" || req.Group.WorkspaceID != req.WorkspaceID {
		return nil, fmt.Errorf("reserve successor draft: %w", ErrWriteGroupValidation)
	}
	metadata, found, err := r.loadWriteGroupLineage(ctx, tx, req.WorkspaceID, req.PredecessorGroupID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, lineageCorrupt("recreation predecessor has no canonical lineage")
	}
	if metadata.LineageID != req.LineageID || metadata.CurrentGroupID != req.PredecessorGroupID ||
		metadata.Revision != req.ExpectedLineageRevision {
		return nil, recreateReservationConflict("lineage_changed")
	}
	if metadata.ReservedGroupID.Valid {
		return nil, recreateReservationConflict("successor_reserved")
	}
	predecessor, err := r.get(ctx, tx, req.WorkspaceID, req.PredecessorGroupID)
	if err != nil {
		return nil, lineageCorrupt("read recreation predecessor: %v", err)
	}
	if predecessor.Revision != req.ExpectedPredecessorRevision {
		return nil, recreateReservationConflict("predecessor_changed")
	}
	if predecessor.Status != WriteGroupStatusDeleted {
		return nil, ErrWriteGroupLifecycleBlocked
	}
	if predecessor.Destination.ConnectorRevision != req.ExpectedPredecessorConnectorRevision ||
		req.Group.Destination.ConnectorRevision != req.ExpectedTargetConnectorRevision {
		return nil, fmt.Errorf("recreate connector revision: %w", ErrWriteGroupConnectorRevisionConflict)
	}
	if metadata.Generation == math.MaxInt64 || metadata.Revision == math.MaxInt64 {
		return nil, lineageCorrupt("recreation generation or revision is exhausted")
	}
	var future int
	if err := tx.QueryRowContext(ctx, r.query(`SELECT COUNT(*) FROM write_group_generations WHERE lineage_id=$1 AND generation>$2`),
		metadata.LineageID, metadata.Generation).Scan(&future); err != nil {
		return nil, fmt.Errorf("read successor generations: %w", err)
	}
	if future != 0 {
		return nil, lineageCorrupt("unreserved future generation exists")
	}
	candidate := cloneWriteGroup(req.Group)
	candidate.Migration = WriteGroupMigration{}
	candidate.Destination.SchemaRevision, candidate.Destination.SchemaDigest = "", ""
	// A wide generation never inherits the physical identity of its predecessor.
	if candidate.isWideProfile() || predecessor.RecordingProfile == nil {
		table := strings.TrimSpace(candidate.Destination.TableName)
		if table != "" && table != managedGroupTableName(predecessor.ID) {
			return nil, transitionRejection("wide-table-identity-change")
		}
		candidate.Destination.TableName = ""
	}
	// Reject changed Basic semantics before generic metadata normalization.
	if predecessor.RecordingProfile == nil {
		if err := r.validateRecreatedBasicInTx(ctx, tx, predecessor, candidate); err != nil {
			return nil, err
		}
	}
	candidate, err = r.prepareNewWriteGroupDraftInTx(ctx, tx, candidate)
	if err != nil {
		return nil, err
	}
	if err := r.validateRecreatedDraftIdentityInTx(ctx, tx, predecessor, candidate); err != nil {
		return nil, err
	}
	if err := r.insertGroup(ctx, tx, candidate); err != nil {
		return nil, err
	}
	if _, err := tx.ExecContext(ctx, r.query(`INSERT INTO write_group_generations (group_id,lineage_id,generation,predecessor_group_id)
  VALUES ($1,$2,$3,$4)`), candidate.ID, metadata.LineageID, metadata.Generation+1, predecessor.ID); err != nil {
		return nil, fmt.Errorf("insert successor generation: %w", err)
	}
	result, err := tx.ExecContext(ctx, r.query(`UPDATE write_group_lineages
  SET reserved_group_id=$1,reservation_operation_id=$2,reservation_intent_digest=$3,revision=revision+1
  WHERE lineage_id=$4 AND workspace_id=$5 AND current_group_id=$6 AND revision=$7
  AND reserved_group_id IS NULL AND reservation_operation_id IS NULL AND reservation_intent_digest IS NULL`),
		candidate.ID, req.OperationID, digest, metadata.LineageID, req.WorkspaceID, predecessor.ID, req.ExpectedLineageRevision)
	if err != nil {
		return nil, fmt.Errorf("reserve successor generation: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("confirm successor reservation: %w", err)
	}
	if affected != 1 {
		return nil, recreateReservationConflict("lineage_changed")
	}
	if err := r.ensureStorageScopeInTx(ctx, tx, candidate); err != nil {
		return nil, err
	}
	// The derived compatibility field follows the unchanged genesis alias.
	candidate.BasicManagedDeviceID = predecessor.BasicManagedDeviceID
	return &writeGroupReservedDraft{LineageID: metadata.LineageID, LineageRevision: metadata.Revision + 1,
		Generation: metadata.Generation + 1, PredecessorGroupID: predecessor.ID, OperationID: req.OperationID, IntentDigest: digest,
		Group: cloneWriteGroup(candidate)}, nil
}
