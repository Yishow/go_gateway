package recordingplan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// GetPlanByWorkspaceInTx reads one plan through the caller's local
// transaction without applying service validation or default mutation.
func (r *SQLRepository) GetPlanByWorkspaceInTx(
	ctx context.Context,
	tx *sql.Tx,
	id string,
	workspaceID string,
) (*RecordingPlan, error) {
	if r == nil || tx == nil {
		return nil, fmt.Errorf("read recording plan in transaction: %w", ErrPlanNotFound)
	}
	query := adaptPlaceholders(`
		SELECT id, workspace_id, revision, applied_revision, name, status, timezone,
			members, streams, destinations, retention, limits, created_at, updated_at
		FROM recording_plans
		WHERE id = $1 AND workspace_id = $2
	`)
	plan, err := scanRecordingPlan(tx.QueryRowContext(ctx, query, id, workspaceID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", ErrPlanNotFound, id)
	}
	if err != nil {
		return nil, fmt.Errorf("read recording plan in transaction: %w", err)
	}
	return plan, nil
}
