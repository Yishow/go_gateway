package api

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"

	"github.com/stretchr/testify/require"
)

// Cleanup accepts only records whose immutable outbox cells and receipt digest
// are available. Validate the entire set under locks before deleting any row.
func cleanupSQLServerCircuitEffects(t *testing.T, db *sql.DB, schemaName, table string, recordIDs []string, effectDigests map[string]string, frozenRows map[string]dbtarget.EncodedRow) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	require.NoError(t, err)
	defer tx.Rollback() // An assertion failure must roll back the entire set.
	type deletion struct {
		query string
		args  []any
	}
	deletions := make([]deletion, 0, len(recordIDs)*2)
	seenEffects := map[string]bool{}
	for _, id := range recordIDs {
		row, ok := frozenRows[id]
		require.True(t, ok, "cleanup requires the immutable outbox record")
		require.Equal(t, id, row.RecordID)
		digest, err := groupdelivery.RowPayloadDigest(row)
		require.NoError(t, err)
		require.Equal(t, effectDigests[row.EffectKey], digest)
		cells := map[string]any{}
		var owner string
		for _, cell := range row.Cells {
			cells[cell.Column] = cell.Value
			if strings.HasPrefix(cell.Column, "_gw_owner_") {
				require.Empty(t, owner)
				owner = cell.Column
			}
		}
		require.Regexp(t, `^_gw_owner_[0-9a-f]{24}$`, owner)
		require.Equal(t, "gateway", cells[owner])
		args := []any{id, cells["workspace_id"], cells["group_id"], cells["group_revision"], cells["provenance"], "gateway"}
		predicate := "record_id=@p1 AND workspace_id=@p2 AND group_id=@p3 AND group_revision=@p4 AND provenance=@p5 AND [" + owner + "]=@p6"
		var matches int
		require.NoError(t, tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+sqlServerQualified(schemaName, table)+
			" WITH (UPDLOCK,HOLDLOCK) WHERE "+predicate, args...).Scan(&matches))
		require.Equal(t, 1, matches, "cleanup refuses changed or foreign record identity/provenance")
		deletions = append(deletions, deletion{"DELETE FROM " + sqlServerQualified(schemaName, table) + " WHERE " + predicate, args})
		if !seenEffects[row.EffectKey] {
			receiptArgs := []any{row.EffectKey, digest, "gateway"}
			receiptPredicate := "effect_key=@p1 AND payload_digest=@p2 AND [_gw_effect_receipts_v1]=@p3"
			require.NoError(t, tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+sqlServerQualified(schemaName, dbtarget.EffectReceiptTable)+
				" WITH (UPDLOCK,HOLDLOCK) WHERE "+receiptPredicate, receiptArgs...).Scan(&matches))
			require.Equal(t, 1, matches, "cleanup refuses a missing or conflicting receipt")
			deletions = append(deletions, deletion{"DELETE FROM " + sqlServerQualified(schemaName, dbtarget.EffectReceiptTable) + " WHERE " + receiptPredicate, receiptArgs})
			seenEffects[row.EffectKey] = true
		}
	}
	for _, statement := range deletions {
		result, err := tx.ExecContext(ctx, statement.query, statement.args...)
		require.NoError(t, err)
		affected, err := result.RowsAffected()
		require.NoError(t, err)
		require.EqualValues(t, 1, affected)
	}
	require.NoError(t, tx.Commit())
}
