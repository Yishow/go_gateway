package measurement

import (
	"context"
	"database/sql"
	"fmt"
)

// GetByIDInTx reads one measurement through the caller's local transaction.
// It deliberately bypasses Validate so persisted payloads remain unchanged.
func (r *SQLRepository) GetByIDInTx(ctx context.Context, tx *sql.Tx, id string) (*MeasurementDefinition, error) {
	if r == nil || tx == nil {
		return nil, fmt.Errorf("read measurement in transaction: %w", sql.ErrNoRows)
	}
	query := adaptPlaceholders(`
		SELECT id, workspace_id, device_id, point_id, tag_id, equipment_id,
			definition_revision, source_binding_revision, series_epoch,
			name, quantity, unit, semantic_kind, numeric_encoding,
			counter_policy, state_map, bitmask_labels, created_at, updated_at
		FROM measurement_definitions
		WHERE id = $1
	`)
	return scanMeasurementDefinition(tx.QueryRowContext(ctx, query, id))
}
