package groupdelivery

import (
	"context"
	"database/sql"
)

// reclaimMixedOutcomes follows existing status retention, but retains original
// outcome evidence while any corresponding outbox effect is still undelivered.
// The current progress row and its immutable runtime descriptor are not history.
func reclaimMixedOutcomes(ctx context.Context, tx *sql.Tx, cutoff string) (sql.Result, error) {
	return tx.ExecContext(ctx, `DELETE FROM wg_delivery_mixed_outcomes AS m WHERE created_at<=? AND NOT EXISTS(SELECT 1 FROM wg_delivery_outbox o WHERE o.workspace_id=m.workspace_id AND o.group_id=m.group_id AND o.group_revision=m.group_revision AND o.record_id=m.record_id)`, cutoff)
}
