package groupdelivery

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedActualPartitionHeadQueryUsesIndexAndKeepsFIFO(t *testing.T) {
	store, db, _ := newTestStore(t)
	for _, row := range []struct {
		id     string
		second int
		mixed  bool
	}{{"legacy-early", 10, false}, {"mixed-a", 15, true}, {"mixed-b", 15, true}, {"legacy-last", 20, false}} {
		bucket := timestamp(at(row.second))
		var due any
		if row.mixed {
			due = bucket
			bucket = ""
		}
		_, err := db.ExecContext(t.Context(), `INSERT INTO wg_delivery_outbox(effect_key,record_id,workspace_id,group_id,group_revision,bucket_start,due_at,partition_key,destination_scope,connector_id,connector_revision,table_name,payload,payload_digest,next_retry_at,created_at,updated_at) VALUES(?,'r','w','g','v',?,?,'part','target','c','cr','data','{}','digest',?,?,?)`, row.id, bucket, due, timestamp(t0), timestamp(t0), timestamp(t0))
		require.NoError(t, err)
	}
	// This is the exact production statement, including the partition-specific
	// variant; not a simplified stand-in and not a throughput benchmark.
	for _, one := range []bool{false, true} {
		args := []any{timestamp(at(200))}
		if one {
			args = append(args, "part")
		}
		rows, err := db.QueryContext(t.Context(), `EXPLAIN QUERY PLAN `+headStatement(one), args...)
		require.NoError(t, err)
		var plans []string
		for rows.Next() {
			var id, parent, unused int
			var detail string
			require.NoError(t, rows.Scan(&id, &parent, &unused, &detail))
			plans = append(plans, detail)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		require.Contains(t, strings.Join(plans, "\n"), "SEARCH e USING INDEX idx_wg_delivery_outbox_due_head")
		t.Logf("actual head partition=%v: %s", one, strings.Join(plans, "; "))
	}
	for _, want := range []string{"legacy-early", "mixed-a", "mixed-b", "legacy-last"} {
		heads, err := store.ReadyHeads(t.Context(), at(200))
		require.NoError(t, err)
		require.Len(t, heads, 1)
		require.Equal(t, want, heads[0].EffectKey)
		_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET state='sql_committed' WHERE effect_key=?`, want)
		require.NoError(t, err)
	}
}
