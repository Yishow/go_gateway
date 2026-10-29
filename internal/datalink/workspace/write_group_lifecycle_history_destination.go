package workspace

import (
	"context"
	"database/sql"
	"encoding/json"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
)

// lifecycleRecordKeyInTx mirrors the real boundary's frozen destination rule:
// only MySQL/SQL Server persist RecordKeyColumn in outbox metadata. The current
// connector cannot establish the dialect of an older applied revision.
func (r *SQLWriteGroupRepository) lifecycleRecordKeyInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup) (string, error) {
	if group.RowPolicy.RecordKeyColumn == "" {
		return "", nil
	}
	// No outbox means this comparison is irrelevant. Reading historical identity
	// or collecting samples does not assert verified delivery/storage capability.
	var hasOutbox bool
	if err := tx.QueryRowContext(ctx, r.query(`SELECT EXISTS(SELECT 1 FROM wg_delivery_outbox WHERE workspace_id=$1 AND group_id=$2 AND group_revision=$3)`), group.WorkspaceID, group.ID, group.Revision).Scan(&hasOutbox); err != nil {
		return "", err
	}
	if !hasOutbox {
		return "", nil
	}
	var count int
	if err := tx.QueryRowContext(ctx, r.query(`SELECT COUNT(*) FROM wg_runtime_versions WHERE group_id=$1 AND group_revision=$2`), group.ID, group.Revision).Scan(&count); err != nil {
		return "", err
	}
	if count != 1 {
		return "", groupdelivery.ErrOwnershipUnavailable
	}
	var payload string
	if err := tx.QueryRowContext(ctx, r.query(`SELECT payload FROM wg_runtime_versions WHERE workspace_id=$1 AND group_id=$2 AND group_revision=$3`), group.WorkspaceID, group.ID, group.Revision).Scan(&payload); err != nil {
		return "", groupdelivery.ErrOwnershipUnavailable
	}
	var frozen WriteGroupAppliedSnapshot
	if json.Unmarshal([]byte(payload), &frozen) != nil || frozen.Group == nil || frozen.RuntimeLayout == nil ||
		frozen.WorkspaceID != group.WorkspaceID || frozen.GroupID != group.ID || frozen.GroupRevision != group.Revision || frozen.AppliedRevision != group.Revision ||
		frozen.Group.WorkspaceID != group.WorkspaceID || frozen.Group.ID != group.ID || frozen.Group.Revision != group.Revision || frozen.Group.AppliedRevision != group.Revision ||
		frozen.Group.Destination != group.Destination || frozen.Group.RowPolicy.RecordKeyColumn != group.RowPolicy.RecordKeyColumn {
		return "", groupdelivery.ErrOwnershipUnavailable
	}
	switch frozen.RuntimeLayout.Dialect {
	case dbtarget.SQLDialectSQLite, dbtarget.SQLDialectPostgres:
		return "", nil
	case dbtarget.SQLDialectMySQL, dbtarget.SQLDialectSQLServer:
		return group.RowPolicy.RecordKeyColumn, nil
	default:
		return "", groupdelivery.ErrOwnershipUnavailable
	}
}
