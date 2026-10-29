package grouppipeline

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"

	"github.com/stretchr/testify/require"
)

const proofRepairAction = "preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance"

func counterDurableState(t *testing.T, store *groupdelivery.Store) map[string][]string {
	t.Helper()
	state := map[string][]string{}
	for _, table := range []string{"wg_delivery_samples", "wg_delivery_checkpoints", "wg_delivery_counter_baselines", "wg_delivery_outbox", "wg_runtime_versions", "wg_counter_proof_associations"} {
		rows, err := store.DB().QueryContext(t.Context(), `SELECT * FROM `+table+` ORDER BY rowid`)
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		for rows.Next() {
			values := make([]any, len(columns))
			targets := make([]any, len(columns))
			for i := range values {
				targets[i] = &values[i]
			}
			require.NoError(t, rows.Scan(targets...))
			wire, err := json.Marshal(values)
			require.NoError(t, err)
			state[table] = append(state[table], string(wire))
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
	}
	return state
}

func damageCounterProof(t *testing.T, f *longPipeline, reason string) {
	t.Helper()
	raw, err := f.store.RuntimeVersion(t.Context(), groupKey(f.snapshot))
	require.NoError(t, err)
	var descriptor map[string]any
	require.NoError(t, json.Unmarshal(raw, &descriptor))
	if reason == "counter-replay-proof-missing" {
		delete(descriptor, "counter_bindings")
	} else {
		for _, entry := range descriptor["counter_bindings"].(map[string]any) {
			proof := entry.(map[string]any)
			if reason == "counter-replay-proof-conflict" {
				proof["unit"] = "private-reading-secret"
			} else {
				proof["version"] = 77
			}
		}
	}
	corrupted, err := json.Marshal(descriptor)
	require.NoError(t, err)
	_, err = f.local.ExecContext(t.Context(), `UPDATE wg_runtime_versions SET payload=? WHERE group_id=? AND group_revision=?`, string(corrupted), f.group.ID, f.group.Revision)
	require.NoError(t, err)
}

func TestCounterReplayQuarantineReportsExactPendingInputsAndRetainsEvidence(t *testing.T) {
	for _, reason := range []string{"counter-replay-proof-missing", "counter-replay-proof-conflict", "counter-replay-proof-unsupported"} {
		t.Run(reason, func(t *testing.T) {
			f := newLongPipeline(t)
			require.NoError(t, f.reading(0, "baseline", at(1), 100))
			f.closeUntil(at(600))
			for i := range 37 {
				require.NoError(t, f.reading(0, fmt.Sprintf("pending-%d", i), at(601+i), uint64(105+i)))
			}
			damageCounterProof(t, f, reason)
			before := counterDurableState(t, f.store)
			reader := &unavailableExactCounterReader{}
			restart := func() {
				f.pipe = New(Dependencies{Groups: f.groups, Destinations: f.dest, Store: f.store, Counters: reader}, Config{Now: func() time.Time { return f.clock }})
				require.NoError(t, f.pipe.Reconcile(t.Context()))
			}
			restart()
			for range 2 {
				view, err := f.pipe.Delivery(t.Context(), f.group.ID)
				require.NoError(t, err)
				require.Equal(t, []CounterReplayQuarantine{{WorkspaceID: "workspace-A", GroupID: "group-C", GroupRevision: "rev-1", PendingSamples: 37, Reason: reason, Action: proofRepairAction}}, view.CounterReplayQuarantines)
				require.Equal(t, 37, view.Stages.Collecting)
				require.Zero(t, view.Stages.Quarantined, "input counts must not be added to SQL outbox row counts")
				require.Equal(t, StateBlocked, view.Intake.State)
				require.Equal(t, reason, view.Intake.Reason)
				wire, err := json.Marshal(view.CounterReplayQuarantines)
				require.NoError(t, err)
				require.NotContains(t, string(wire), "private-reading-secret")
				require.NotContains(t, string(wire), "12000")
				f.closeUntil(at(1800))
				require.Equal(t, before, counterDurableState(t, f.store))
				require.Zero(t, reader.calls)
				restart()
			}
		})
	}
}
