package mapping

import (
	"context"
	"database/sql"
	"go-gateway/internal/datalink/schema"
)

// ReadByPointIDTx includes disabled mappings and preserves the empty set as
// evidence. It reuses the canonical row decoder without service/cache reads.
func ReadByPointIDTx(ctx context.Context, tx *sql.Tx, pointID string) ([]*schema.Mapping, error) {
	if tx == nil {
		return nil, ErrMappingNotFound
	}
	rows, err := tx.QueryContext(ctx, `SELECT id,point_id,tag_id,transform_pipeline,status,rule_candidate_id,proposed_signature,last_applied_signature,blocking_reason,enabled,created_at,updated_at FROM mappings WHERE point_id=? ORDER BY id`, pointID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result, err := (&SQLRepository{}).scanMappings(rows)
	if err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
