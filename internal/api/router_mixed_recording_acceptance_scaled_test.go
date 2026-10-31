package api

import (
	"encoding/json"
	"go-gateway/internal/datalink/dbtarget"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMixedAcceptanceScaledDeltaDurableSQLite(t *testing.T) {
	f, request := mixedScaledSaveFixture(t, "delta")
	base := time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)
	r := applyMixedAcceptance(t, f, request, base.Add(7*time.Second))
	r.offer(t, 0, "baseline", base.Add(10*time.Second), uint64(90071992547409930))
	r.tick(t, base.Add(time.Minute))
	var outcomes, outbox int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM wg_delivery_mixed_outcomes WHERE group_id=? AND kind='skipped'`, r.frozen.GroupID).Scan(&outcomes))
	require.Equal(t, 1, outcomes)
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM wg_delivery_outbox WHERE group_id=?`, r.frozen.GroupID).Scan(&outbox))
	require.Zero(t, outbox)
	r.reconcile(t) // A fresh production worker restores the reliable skipped baseline.
	r.offer(t, 0, "next", base.Add(70*time.Second), uint64(90071992547410010))
	r.tick(t, base.Add(2*time.Minute))
	keys := r.deliver(t)
	require.Len(t, keys, 1)
	table := r.frozen.Group.Destination.TableName
	var value, kind, profile string
	require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT CAST(value_decimal AS TEXT),value_type,profile_kind FROM "`+table+`"`).Scan(&value, &kind, &profile))
	require.Equal(t, "8.0", value)
	require.Equal(t, "decimal", kind)
	require.Equal(t, "mixed", profile)
	var affinity, provenance string
	require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT typeof(value_decimal),provenance FROM "`+table+`"`).Scan(&affinity, &provenance))
	require.Equal(t, "text", affinity)
	var evidence map[string]any
	require.NoError(t, json.Unmarshal([]byte(provenance), &evidence))
	require.Equal(t, "9007199254741001.0", evidence["counter"].(map[string]any)["raw"].(map[string]any)["value"])
	require.Equal(t, "8.0", evidence["counter"].(map[string]any)["delta"])
	require.Equal(t, "baseline", evidence["counter"].(map[string]any)["baseline_sample_id"])
	require.Equal(t, "next", evidence["counter"].(map[string]any)["end_sample_id"])
	var targetReceipts int
	require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT count(*) FROM "`+dbtarget.EffectReceiptTable+`" WHERE effect_key=?`, keys[0]).Scan(&targetReceipts))
	require.Equal(t, 1, targetReceipts)

	var receipts int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM wg_delivery_receipts WHERE effect_key=?`, keys[0]).Scan(&receipts))
	require.Equal(t, 1, receipts)
	assertMixedScaledReadback(t, r)
	assertMixedReadbackRejectsForeignProof(t, r)
	r.reconcile(t)
	r.tick(t, base.Add(2*time.Minute))
	require.Empty(t, r.deliver(t))
	var count int
	require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT count(*) FROM "`+table+`"`).Scan(&count))
	require.Equal(t, 1, count)
}
