package api

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
)

// Register after rig construction: LIFO runs proof-scoped effect cleanup before
// the fixture's schema cleanup. No wildcard DELETE or DROP bypasses ownership.
func sqlServerSuccessorRegisterEffectCleanup(t *testing.T, r *sqlServerSuccessorRig) {
	t.Helper()
	t.Cleanup(func() {
		require.NoError(t, r.pipe.Stop(time.Second))
		sqlServerSuccessorCleanupGroupEffects(t, r, r.f.target, r.snapshot.Group)
		sqlServerSuccessorCleanupGroupEffects(t, r, r.targetDB, r.prepared.Group)
	})
}

func sqlServerSuccessorCleanupGroupEffects(t *testing.T, r *sqlServerSuccessorRig, target *sql.DB, group *workspace.WriteGroup) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err := r.f.db.QueryContext(ctx, `SELECT effect_key,payload_digest,payload FROM wg_delivery_outbox WHERE workspace_id=? AND group_id=? ORDER BY effect_key`, group.WorkspaceID, group.ID)
	require.NoError(t, err)
	var recordIDs []string
	frozenRows := map[string]dbtarget.EncodedRow{}
	digests := map[string]string{}
	for rows.Next() {
		var key, digest, payload string
		require.NoError(t, rows.Scan(&key, &digest, &payload))
		row, err := groupdelivery.DecodeRowPayload([]byte(payload))
		require.NoError(t, err)
		require.Equal(t, key, row.EffectKey)
		actual, err := groupdelivery.RowPayloadDigest(row)
		require.NoError(t, err)
		require.Equal(t, digest, actual)
		frozenRows[row.RecordID], digests[key] = row, digest
		recordIDs = append(recordIDs, row.RecordID)
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	committed := make([]string, 0, len(recordIDs))
	for _, recordID := range recordIDs {
		row := frozenRows[recordID]
		var records, receipts int
		require.NoError(t, target.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+sqlServerQualified(group.Destination.TableSchema, group.Destination.TableName)+` WHERE record_id=@p1`, recordID).Scan(&records))
		require.NoError(t, target.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+sqlServerQualified(group.Destination.TableSchema, dbtarget.EffectReceiptTable)+` WHERE effect_key=@p1`, row.EffectKey).Scan(&receipts))
		if records == 0 && receipts == 0 {
			continue // A frozen but undelivered outbox is not an external effect.
		}
		require.Equal(t, 1, records, "cleanup refuses ambiguous SQL row identity")
		require.Equal(t, 1, receipts, "cleanup refuses partial/missing receipt evidence")
		committed = append(committed, recordID)
	}
	if len(committed) != 0 {
		cleanupSQLServerCircuitEffects(t, target, group.Destination.TableSchema, group.Destination.TableName, committed, digests, frozenRows)
	}
	var remaining int
	require.NoError(t, target.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+sqlServerQualified(group.Destination.TableSchema, group.Destination.TableName)).Scan(&remaining))
	require.Zero(t, remaining, "cleanup must not hide unowned or unmatched rows")
}
