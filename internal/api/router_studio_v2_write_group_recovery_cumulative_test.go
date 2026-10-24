package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// storedCumulative returns the destination's long rows in window order as
// "status|raw|delta", read back as exact decimal text.
func (r *recoveryRig) storedCumulative() []string {
	r.t.Helper()
	rows, err := r.dest.QueryContext(r.t.Context(),
		`SELECT status, CAST(value_decimal AS TEXT), COALESCE(CAST(delta_value AS TEXT), 'NULL') FROM `+r.quoted(r.table)+` ORDER BY bucket_start`)
	require.NoError(r.t, err)
	defer rows.Close()
	var out []string
	for rows.Next() {
		var status, raw, delta string
		require.NoError(r.t, rows.Scan(&status, &raw, &delta))
		out = append(out, status+"|"+raw+"|"+delta)
	}
	require.NoError(r.t, rows.Err())
	return out
}

// TestLongStorageCumulativeRecovery loses the confirmed table of a cumulative
// group between two windows. The structure comes back and the waiting windows
// are stored with exact raw values above 2^53 and exact differences: the counter
// baseline lives in the local database, so it neither moves nor restarts because
// the destination's table did.
func TestLongStorageCumulativeRecovery(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres", "timescale"} {
		t.Run(kind, func(t *testing.T) {
			r := newCumulativeRecoveryRig(t, kind)
			r.start()

			r.record(uint64(9007199254740992))
			r.committed(1, "the baseline window is delivered")
			require.Equal(t, []string{"baseline|9007199254740992|NULL"}, r.storedCumulative())
			before := r.receipts()

			r.restart(recoveryWorkerEvery * 50)
			r.drop(r.table)
			r.record(uint64(9007199254740993))
			r.record(uint64(9007199254740995))
			r.committed(3, "the windows waiting when the table vanished are delivered after it is restored")

			require.Equal(t, []string{
				"ok|9007199254740993|1",
				"ok|9007199254740995|2",
			}, r.storedCumulative(), "the first stored window after the loss is a difference, not a new baseline; the lost window is not invented")
			after := r.receipts()
			require.Len(t, after, 3)
			for key, value := range before {
				require.Equal(t, value, after[key], "an existing receipt is neither deleted nor rewritten")
			}
			lineage := r.lineage()
			require.Equal(t, 2, lineage.Generation)
			require.True(t, lineage.HistoryMayBeLost)
			require.Equal(t, 1, lineage.Reconstructions)
		})
	}
}
