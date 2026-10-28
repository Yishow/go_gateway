package workspace

import (
	"context"
	"database/sql"
	"math"
	"strings"
)

// successorTransitionAuthority is minted only after service review/backlog checks
// and cannot authorize another transaction, predecessor, lineage or intent.
type successorTransitionAuthority struct {
	tx                                                                 *sql.Tx
	workspaceID, lineageID, predecessorID, predecessorRevision, digest string
}

// reserveTransitionDraftInTx is distinct from deleted-only recreation. Its
// caller has just revalidated the reviewed envelope and backlog in this tx.
func (r *SQLWriteGroupRepository) reserveTransitionDraftInTx(ctx context.Context, tx *sql.Tx, req writeGroupRecreateReservationRequest, digest string, authority successorTransitionAuthority) (*writeGroupReservedDraft, error) {
	return r.reserveTransitionCandidateInTx(ctx, tx, req, digest, authority, nil)
}

// reserveVerifiedMySQLSuccessorInTx is a distinct public adapter; the original
// deleted-only recreate primitive remains local and refuses MySQL-driver targets.
func (r *SQLWriteGroupRepository) reserveVerifiedMySQLSuccessorInTx(ctx context.Context, tx *sql.Tx, req writeGroupRecreateReservationRequest, digest string, authority successorTransitionAuthority, evidence successorMySQLReserveEvidence) (*writeGroupReservedDraft, error) {
	if evidence.tx != tx || evidence.digest != digest || evidence.workspaceID != req.WorkspaceID {
		return nil, ErrWriteGroupValidation
	}
	return r.reserveTransitionCandidateInTx(ctx, tx, req, digest, authority, &evidence)
}

func (r *SQLWriteGroupRepository) reserveTransitionCandidateInTx(ctx context.Context, tx *sql.Tx, req writeGroupRecreateReservationRequest, digest string, authority successorTransitionAuthority, evidence *successorMySQLReserveEvidence) (*writeGroupReservedDraft, error) {
	if tx == nil || req.Group == nil || digest == "" || authority.tx != tx || authority.workspaceID != req.WorkspaceID || authority.lineageID != req.LineageID || authority.predecessorID != req.PredecessorGroupID || authority.predecessorRevision != req.ExpectedPredecessorRevision || authority.digest != digest {
		return nil, ErrWriteGroupValidation
	}
	meta, found, err := r.loadWriteGroupLineage(ctx, tx, req.WorkspaceID, req.PredecessorGroupID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrWriteGroupLineageCorrupt
	}
	if meta.LineageID != req.LineageID || meta.CurrentGroupID != req.PredecessorGroupID || meta.Revision != req.ExpectedLineageRevision || meta.ReservedGroupID.Valid {
		return nil, recreateReservationConflict("lineage_changed")
	}
	pred, err := r.GetInTx(ctx, tx, req.WorkspaceID, req.PredecessorGroupID)
	if err != nil {
		return nil, err
	}
	if (pred.Status == WriteGroupStatusDeleted && evidence == nil) || pred.RecordingProfile == nil || (pred.AppliedRevision == "" && (evidence == nil || evidence.sourceCapability == nil || pred.Status != WriteGroupStatusDeleted)) {
		return nil, ErrWriteGroupLifecycleBlocked
	}
	if pred.Revision != req.ExpectedPredecessorRevision || pred.Destination.ConnectorRevision != req.ExpectedPredecessorConnectorRevision {
		return nil, recreateReservationConflict("predecessor_changed")
	}
	if meta.Generation == math.MaxInt64 || meta.Revision == math.MaxInt64 {
		return nil, ErrWriteGroupLineageCorrupt
	}
	var future int
	if err := tx.QueryRowContext(ctx, r.query(`SELECT COUNT(*) FROM write_group_generations WHERE lineage_id=$1 AND generation>$2`), meta.LineageID, meta.Generation).Scan(&future); err != nil {
		return nil, err
	}
	if future != 0 {
		return nil, ErrWriteGroupLineageCorrupt
	}
	candidate := cloneWriteGroup(req.Group)
	candidate.ID = r.newID()
	candidate.BasicManagedDeviceID = ""
	candidate.Migration = WriteGroupMigration{}
	candidate.Destination.SchemaRevision, candidate.Destination.SchemaDigest = "", ""
	if candidate.isWideProfile() {
		table := strings.TrimSpace(candidate.Destination.TableName)
		if table != "" && table != managedGroupTableName(pred.ID) {
			return nil, transitionRejection("wide-table-identity-change")
		}
		candidate.Destination.TableName = managedGroupTableName(candidate.ID)
	}
	defaultRecreateMissing(candidate)
	if evidence != nil {
		if err := r.verifySuccessorMySQLInTx(ctx, tx, pred, candidate, &evidence.capability, evidence.sourceCapability); err != nil {
			return nil, err
		}
	} else if err := r.validateSuccessorDefinitionInTx(ctx, tx, pred, candidate); err != nil {
		return nil, err
	}
	source, err := r.recreateStructuralSourceInTx(ctx, tx, pred)
	if err != nil {
		return nil, err
	}
	if source.isLongLayout() && candidate.isLongLayout() && longStorageScopeOf(source) == longStorageScopeOf(candidate) {
		return nil, transitionRejection("destination-unchanged")
	}
	now := r.now()
	candidate.Revision = r.newID()
	candidate.AppliedRevision = ""
	candidate.Status = WriteGroupStatusDraft
	candidate.CreatedAt = now
	candidate.UpdatedAt = now
	if err := r.insertGroup(ctx, tx, candidate); err != nil {
		return nil, err
	}
	if _, err := tx.ExecContext(ctx, r.query(`INSERT INTO write_group_generations(group_id,lineage_id,generation,predecessor_group_id) VALUES($1,$2,$3,$4)`), candidate.ID, meta.LineageID, meta.Generation+1, pred.ID); err != nil {
		return nil, err
	}
	result, err := tx.ExecContext(ctx, r.query(`UPDATE write_group_lineages SET reserved_group_id=$1,reservation_operation_id=$2,reservation_intent_digest=$3,revision=revision+1 WHERE lineage_id=$4 AND workspace_id=$5 AND current_group_id=$6 AND revision=$7 AND reserved_group_id IS NULL AND reservation_operation_id IS NULL AND reservation_intent_digest IS NULL`), candidate.ID, req.OperationID, digest, meta.LineageID, req.WorkspaceID, pred.ID, meta.Revision)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n != 1 {
		return nil, recreateReservationConflict("lineage_changed")
	}
	if err := r.ensureStorageScopeInTx(ctx, tx, candidate); err != nil {
		return nil, err
	}
	candidate.BasicManagedDeviceID = pred.BasicManagedDeviceID
	return &writeGroupReservedDraft{LineageID: meta.LineageID, LineageRevision: meta.Revision + 1, Generation: meta.Generation + 1, PredecessorGroupID: pred.ID, OperationID: req.OperationID, IntentDigest: digest, Group: cloneWriteGroup(candidate)}, nil
}
