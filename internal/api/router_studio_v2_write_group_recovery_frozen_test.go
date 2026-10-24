package api

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// outboxDigests returns effect key → payload digest over the local outbox.
func (r *recoveryRig) outboxDigests() map[string]string {
	r.t.Helper()
	rows, err := r.f.db.QueryContext(r.t.Context(), `SELECT effect_key, payload_digest FROM wg_delivery_outbox`)
	require.NoError(r.t, err)
	defer rows.Close()
	digests := map[string]string{}
	for rows.Next() {
		var key, digest string
		require.NoError(r.t, rows.Scan(&key, &digest))
		digests[key] = digest
	}
	require.NoError(r.t, rows.Err())
	return digests
}

// stored returns the destination's long rows as record_id → "status|value_float",
// read back through the destination's own catalog and SQL.
func (r *recoveryRig) stored() map[string]string {
	r.t.Helper()
	rows, err := r.dest.QueryContext(r.t.Context(), `SELECT record_id, status, COALESCE(CAST(value_float AS TEXT), 'NULL') FROM `+r.quoted(r.table))
	require.NoError(r.t, err)
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var id, status, value string
		require.NoError(r.t, rows.Scan(&id, &status, &value))
		out[id] = status + "|" + value
	}
	require.NoError(r.t, rows.Err())
	return out
}

// TestLongStorageMissingRecoveryKeepsFrozenRows drops the confirmed table while
// closed windows of four groups sharing one scope are still waiting, restarts the
// gateway, and requires the waiting rows to arrive exactly as they were frozen:
// the same payload digests, the same values, no rewritten receipt, and one
// effective rebuild for the whole scope.
func TestLongStorageMissingRecoveryKeepsFrozenRows(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres", "timescale"} {
		t.Run(kind, func(t *testing.T) {
			r := newRecoveryRig(t, kind, "aux", "aux2", "aux3")
			groups := len(r.groups)
			require.Equal(t, 4, groups)
			r.start()
			r.record(10.5)
			r.committed(groups, "the first window is delivered")
			before := r.receipts()
			firstRows := r.stored()
			require.Len(t, firstRows, groups)

			// A long pause between delivery cycles keeps the next windows pending.
			r.restart(time.Hour)
			r.drop(r.table)
			r.record(11.5)
			r.record(12.5)
			r.waitFor(func() bool { return r.outbox()["pending"] == 2*groups }, "the closed rows of every group wait")
			require.False(t, r.tableExists(r.table), "nothing restores the table before a delivery asks for it")
			frozen := r.outboxDigests()

			r.restart(recoveryWorkerEvery)
			r.committed(3*groups, "the waiting rows are delivered after the restart restored the table")

			require.Equal(t, frozen, r.outboxDigests(), "a waiting row is never rebuilt or recomputed")
			after := r.receipts()
			require.Len(t, after, 3*groups)
			for key, value := range before {
				require.Equal(t, value, after[key], "an existing receipt is neither deleted nor rewritten")
			}
			for key, digest := range frozen {
				if _, existed := before[key]; !existed {
					require.True(t, strings.HasPrefix(after[key], digest+"|"), "the receipt of %s proves the frozen payload", key)
				}
			}

			rows := r.stored()
			require.Len(t, rows, 2*groups, "only what was still waiting: the first window went down with the table")
			values := map[string]int{}
			for id, row := range rows {
				require.NotContains(t, firstRows, id, "no earlier record is invented")
				values[row]++
			}
			require.Equal(t, map[string]int{"ok|11.5": groups, "ok|12.5": groups}, values)

			lineage := r.lineage()
			require.Equal(t, 2, lineage.Generation, "the generation moves once")
			require.Equal(t, 1, lineage.Reconstructions)
			require.True(t, lineage.HistoryMayBeLost)
			require.Equal(t, 1, r.repairOperations(), "the groups delivered at the same time shared one effective rebuild")
		})
	}
}
