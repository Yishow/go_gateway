package workspace

import (
	"context"
	"database/sql"
	"fmt"
)

// insertOwnershipFenceInTx binds the existing intent to immutable membership
// on BOTH sides, including an already reviewed retained-source successor. The
// cutover ID is an opaque fence revision, not a clock or a global counter.
//
// No current group/destination reconstruction is allowed: the same applied
// versions used by cutover/revision-end authority supply the frozen scopes.
// These internal rows are not operator-facing data. Acceptance and closure
// consumers must later validate the fence in their own local transaction.
func (r *SQLWriteGroupRepository) insertOwnershipFenceInTx(ctx context.Context, tx *sql.Tx, c *WriteGroupCutover) error {
	for _, end := range c.Ended {
		if err := r.insertOwnershipFenceScopeInTx(ctx, tx, c, "ended", end.GroupID, end.Revision); err != nil {
			return err
		}
	}
	for _, start := range c.Started {
		if err := r.insertOwnershipFenceScopeInTx(ctx, tx, c, "started", start.GroupID, start.Revision); err != nil {
			return err
		}
	}
	return nil
}

func (r *SQLWriteGroupRepository) insertOwnershipFenceScopeInTx(ctx context.Context, tx *sql.Tx, c *WriteGroupCutover, direction, groupID, revision string) error {
	version, err := r.versionByRevisionInTx(ctx, tx, c.WorkspaceID, groupID, revision)
	if err != nil {
		return fmt.Errorf("read ownership fence scope: %w", err)
	}
	group := version.Group
	if len(group.Members) == 0 {
		return fmt.Errorf("empty ownership fence scope: %w", ErrWriteGroupValidation)
	}
	destination := group.Destination
	for _, member := range group.Members {
		_, err := tx.ExecContext(ctx, r.query(`INSERT INTO write_group_ownership_fences
   (fence_revision, workspace_id, direction, group_id, group_revision, device_id, point_id, tag_id,
    profile_kind, connector_id, connector_revision, database_name, table_schema, table_name)
   VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`),
			c.OwnershipFenceRevision, c.WorkspaceID, direction, groupID, revision, member.DeviceID, member.PointID, member.TagID,
			profileKindOf(group), destination.ConnectorID, destination.ConnectorRevision, destination.Database, destination.TableSchema, destination.TableName)
		if err != nil {
			return fmt.Errorf("record ownership fence scope: %w", err)
		}
	}
	return nil
}
