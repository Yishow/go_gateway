package workspace

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPermanentDeleteAfterReclaimLeavesNoReceiptForRecording(t *testing.T) {
	db, service, group := lifecyclePreservationFixture(t)
	store := seedLifecycleAcceptedData(t, db, group)
	ctx := t.Context()
	_, err := db.ExecContext(ctx, `DELETE FROM wg_delivery_outbox WHERE effect_key NOT IN ('Qcommitted','Qretry','Qblocked')`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE wg_delivery_outbox SET committed_at='2026-10-02T12:02:00.000000000Z' WHERE effect_key='Qcommitted'`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE wg_delivery_outbox SET group_id='other', state='sql_committed', committed_at='2026-10-05T11:59:00.000000000Z' WHERE effect_key='Qretry'`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE wg_delivery_outbox SET group_id='other', state='pending' WHERE effect_key='Qblocked'`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO wg_delivery_receipts(effect_key,payload_digest,committed_at) SELECT effect_key,payload_digest,'2026-10-05T11:59:00.000000000Z' FROM wg_delivery_outbox WHERE effect_key IN ('Qretry','Qblocked')`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `DELETE FROM wg_delivery_samples`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `DELETE FROM managed_schema_operations`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE write_groups SET status='deleted' WHERE id=?`, group.ID)
	require.NoError(t, err)
	service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))

	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	_, err = store.WithClock(func() time.Time { return now }).Reclaim(ctx, time.Hour)
	require.NoError(t, err)
	var count int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_outbox WHERE effect_key='Qcommitted'`).Scan(&count))
	require.Zero(t, count, "the old committed outbox row was reclaimed")

	result, err := service.Delete(ctx, group.ID, lifecycleReadMutation(t, service, group))
	require.NoError(t, err)
	require.True(t, result.Deleted)
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key='Qcommitted'`).Scan(&count))
	require.Zero(t, count, "no receipt remains for the reclaimed recording")
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key IN ('Qretry','Qblocked')`).Scan(&count))
	require.Equal(t, 2, count, "recent and pending outbox receipts remain")
}
