package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestCutoverSourceEvidenceReportsTheActualWindow runs the production cutover over a real destination and reads the
// per source evidence of the partial record the boundary produced: the outbox keeps the nominal bucket start and the
// table keeps the real one, and the evidence must report what the table holds, with the destination the long group
// froze and a matching receipt.
func TestCutoverSourceEvidenceReportsTheActualWindow(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			r := newCutoverRig(t, kind)
			longID := r.long["id"].(string)
			r.accept(r.wideEffective.Add(2*time.Second), 1.5)
			r.advance(r.wideEffective.Add(16 * time.Second))
			r.waitFor(func() bool { return r.wideRows() == 1 }, "the wide history row is delivered")

			boundary := r.at(22*time.Second + 500*time.Millisecond)
			r.cutover(boundary)
			r.reconcile()
			r.accept(boundary, 7.0)
			r.advance(r.at(41 * time.Second)) // the long window [20s,40s) closes, cut at its start
			r.waitFor(func() bool { return len(r.longRecords()) == 1 }, "the cut window reaches the long table")

			var latest map[string]any
			var evidence map[string]any
			require.Eventually(t, func() bool {
				evidence, _ = r.delivery(longID)["source_evidence"].(map[string]any)
				if evidence == nil || len(evidence["sources"].([]any)) != 1 {
					return false
				}
				latest, _ = evidence["sources"].([]any)[0].(map[string]any)["latest"].(map[string]any)
				return latest != nil && latest["delivery"] == "sql_committed"
			}, 15*time.Second, 50*time.Millisecond, "the cut record is SQL committed with a matching receipt")

			stored := r.longRecords()[0]
			require.Equal(t, "partial", stored.Status)
			require.Equal(t, true, latest["available"], "a real cutover window is accepted: %+v", latest)
			require.Equal(t, "partial", latest["status"])
			start, end := parseEvidenceTime(t, latest["window_start"]), parseEvidenceTime(t, latest["window_end"])
			require.True(t, stored.Start.Equal(start), "evidence start %s is what the table holds %s", start, stored.Start)
			require.True(t, stored.End.Equal(end))
			require.True(t, boundary.Equal(start), "the real boundary, not the nominal bucket start")
			require.False(t, r.at(20*time.Second).Equal(start))
			require.Equal(t, destinationOf(r.long)["connector_id"], evidence["connector_id"], "the frozen applied destination is the scope")
			require.Equal(t, destinationOf(r.long)["connector_revision"], evidence["connector_revision"])
			require.Equal(t, evidence["connector_revision"], latest["connector_revision"])
			require.NotNil(t, latest["committed_at"])
		})
	}
}
