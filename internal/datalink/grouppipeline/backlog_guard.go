package grouppipeline

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

// ErrBacklogOwnershipUnproven means accepted data of a group cannot be tied to
// an immutable group revision and a frozen destination, so deleting the group
// could orphan it.
var ErrBacklogOwnershipUnproven = errors.New("accepted data ownership cannot be proven")

// BacklogGuard lets an applied group be deleted only when every accepted but
// undelivered row and every uncollected sample still names the immutable
// revision and the frozen destination it was accepted for. The tombstoned
// group's backlog is then delivered by the same workers from those frozen
// revisions; deleting never reroutes, replays or drops it.
type BacklogGuard struct{}

var _ workspace.WriteGroupBacklogOwnershipGuard = BacklogGuard{}

// CheckWriteGroupBacklog reads only through tx.
func (BacklogGuard) CheckWriteGroupBacklog(ctx context.Context, tx *sql.Tx, group *workspace.WriteGroup) error {
	if tx == nil || group == nil || group.ID == "" || group.WorkspaceID == "" {
		return ErrBacklogOwnershipUnproven
	}
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT workspace_id, group_revision, connector_id, connector_revision,
		 database_name, table_schema, table_name, destination_scope, dedupe_capability
		FROM wg_delivery_outbox
		WHERE group_id = ? AND state NOT IN ('sql_committed', 'operator_skipped')`, group.ID)
	if err != nil {
		return fmt.Errorf("read undelivered rows: %w", err)
	}
	type owner struct{ workspace, revision, connector, connectorRevision, database, schema, table, scope, dedupe string }
	var owners []owner
	for rows.Next() {
		var o owner
		if err := rows.Scan(&o.workspace, &o.revision, &o.connector, &o.connectorRevision, &o.database, &o.schema, &o.table, &o.scope, &o.dedupe); err != nil {
			_ = rows.Close()
			return fmt.Errorf("scan undelivered row: %w", err)
		}
		owners = append(owners, o)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return fmt.Errorf("iterate undelivered rows: %w", err)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, o := range owners {
		version, err := requireVersion(ctx, tx, group, o.workspace, o.revision)
		if err != nil {
			return err
		}
		d := version.Destination
		scope, err := snapshot.DestinationScope(d.ConnectorID, d.ConnectorRevision, d.Database, d.TableSchema, d.TableName)
		dedupe := strings.ToLower(strings.TrimSpace(version.WritePolicy.DedupeCapability))
		if dedupe == "" {
			dedupe = groupdelivery.DedupeNone
		}
		if err != nil || o.connector != d.ConnectorID || o.connectorRevision != d.ConnectorRevision ||
			o.database != d.Database || o.schema != d.TableSchema || o.table != d.TableName || o.scope != scope || o.dedupe != dedupe {
			return ErrBacklogOwnershipUnproven
		}
	}
	revisions, err := tx.QueryContext(ctx,
		`SELECT DISTINCT workspace_id, group_revision FROM wg_delivery_samples WHERE group_id = ? AND consumed = 0`, group.ID)
	if err != nil {
		return fmt.Errorf("read uncollected samples: %w", err)
	}
	type pendingOwner struct{ workspace, revision string }
	var pending []pendingOwner
	for revisions.Next() {
		var owner pendingOwner
		if err := revisions.Scan(&owner.workspace, &owner.revision); err != nil {
			_ = revisions.Close()
			return fmt.Errorf("scan uncollected sample: %w", err)
		}
		pending = append(pending, owner)
	}
	if err := revisions.Err(); err != nil {
		_ = revisions.Close()
		return fmt.Errorf("iterate uncollected samples: %w", err)
	}
	if err := revisions.Close(); err != nil {
		return err
	}
	for _, owner := range pending {
		if _, err := requireVersion(ctx, tx, group, owner.workspace, owner.revision); err != nil {
			return err
		}
	}
	return nil
}

func requireVersion(ctx context.Context, tx *sql.Tx, group *workspace.WriteGroup, workspaceID, revision string) (*workspace.WriteGroup, error) {
	if workspaceID != group.WorkspaceID {
		return nil, ErrBacklogOwnershipUnproven
	}
	version, err := workspace.NewSQLWriteGroupRepository(nil).VersionByRevisionInTx(ctx, tx, workspaceID, group.ID, revision)
	if err != nil {
		return nil, errors.Join(ErrBacklogOwnershipUnproven, err)
	}
	d := version.Group.Destination
	if _, err := snapshot.DestinationScope(d.ConnectorID, d.ConnectorRevision, d.Database, d.TableSchema, d.TableName); err != nil {
		return nil, ErrBacklogOwnershipUnproven
	}
	return version.Group, nil
}
