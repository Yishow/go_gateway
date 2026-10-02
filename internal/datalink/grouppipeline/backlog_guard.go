package grouppipeline

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	if tx == nil || group == nil || group.ID == "" {
		return ErrBacklogOwnershipUnproven
	}
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT group_revision, connector_id, connector_revision, table_name
		FROM wg_delivery_outbox
		WHERE group_id = ? AND state NOT IN ('sql_committed', 'operator_skipped')`, group.ID)
	if err != nil {
		return fmt.Errorf("read undelivered rows: %w", err)
	}
	type owner struct{ revision, connector, connectorRevision, table string }
	var owners []owner
	for rows.Next() {
		var o owner
		if err := rows.Scan(&o.revision, &o.connector, &o.connectorRevision, &o.table); err != nil {
			_ = rows.Close()
			return fmt.Errorf("scan undelivered row: %w", err)
		}
		owners = append(owners, o)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, o := range owners {
		if o.revision == "" || o.connector == "" || o.connectorRevision == "" || o.table == "" {
			return ErrBacklogOwnershipUnproven
		}
		if err := requireVersion(ctx, tx, group.ID, o.revision); err != nil {
			return err
		}
	}
	revisions, err := tx.QueryContext(ctx,
		`SELECT DISTINCT group_revision FROM wg_delivery_samples WHERE group_id = ? AND consumed = 0`, group.ID)
	if err != nil {
		return fmt.Errorf("read uncollected samples: %w", err)
	}
	var pending []string
	for revisions.Next() {
		var revision string
		if err := revisions.Scan(&revision); err != nil {
			_ = revisions.Close()
			return fmt.Errorf("scan uncollected sample: %w", err)
		}
		pending = append(pending, revision)
	}
	if err := revisions.Close(); err != nil {
		return err
	}
	for _, revision := range pending {
		if err := requireVersion(ctx, tx, group.ID, revision); err != nil {
			return err
		}
	}
	return nil
}

func requireVersion(ctx context.Context, tx *sql.Tx, groupID, revision string) error {
	var found int
	err := tx.QueryRowContext(ctx,
		`SELECT 1 FROM write_group_versions WHERE group_id = ? AND group_revision = ?`, groupID, revision).Scan(&found)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrBacklogOwnershipUnproven
	}
	if err != nil {
		return fmt.Errorf("read immutable group version: %w", err)
	}
	return nil
}
