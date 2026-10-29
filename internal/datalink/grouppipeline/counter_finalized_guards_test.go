package grouppipeline

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestCounterFinalizedLegacyBacklogNeedsNoUnfinalizedProof(t *testing.T) {
	f := newLongPipeline(t)
	require.NoError(t, f.reading(0, "legacy-finalized", at(1), 100))
	f.closeUntil(at(600))
	original := firstCounterEffect(t, f)
	damageCounterProof(t, f, "counter-replay-proof-missing")
	before := counterDurableState(t, f.store)
	f.groups.groups = nil
	reader := &unavailableExactCounterReader{}
	f.pipe = New(Dependencies{Groups: f.groups, Destinations: f.dest, Store: f.store, Counters: reader}, Config{Now: func() time.Time { return f.clock }})
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	result, err := counterFinalizedSender(f).Deliver(t.Context(), original.EffectKey)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StateCommitted, result.State)
	after, err := f.store.GetOutbox(t.Context(), original.EffectKey)
	require.NoError(t, err)
	require.Equal(t, original.Payload, after.Payload)
	require.Equal(t, original.PayloadDigest, after.PayloadDigest)
	require.Equal(t, original.Destination, after.Destination)
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Empty(t, entries, "no unfinalized accepted set exists to quarantine or reinterpret")
	now := counterDurableState(t, f.store)
	for _, table := range []string{"wg_delivery_samples", "wg_delivery_checkpoints", "wg_delivery_counter_baselines", "wg_runtime_versions", "wg_counter_proof_associations"} {
		require.Equal(t, before[table], now[table], table)
	}
	require.Equal(t, 1, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
	require.Zero(t, reader.calls)
}

func TestCounterFinalizedDamagedProofDoesNotAuthorizeInvalidOutputOrStorage(t *testing.T) {
	for _, fault := range []string{"payload digest", "destination", "owner marker", "scope alias"} {
		t.Run(fault, func(t *testing.T) {
			f := newLongPipeline(t)
			require.NoError(t, f.reading(0, "finalized", at(1), 100))
			f.closeUntil(at(600))
			original := firstCounterEffect(t, f)
			damageCounterProof(t, f, "counter-replay-proof-missing")
			code := "storage_owner_mismatch"
			switch fault {
			case "payload digest":
				tampered := strings.ReplaceAll(string(original.Payload), `"kWh"`, `"MWh"`)
				require.NotEqual(t, string(original.Payload), tampered)
				_, err := f.local.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET payload=? WHERE effect_key=?`, tampered, original.EffectKey)
				require.NoError(t, err)
				code = "payload-digest-mismatch"
			case "destination":
				_, err := f.dest.target.ExecContext(t.Context(), `CREATE TABLE forged_sink AS SELECT * FROM `+f.table+` WHERE 0`)
				require.NoError(t, err)
				_, err = f.local.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET table_name='forged_sink' WHERE effect_key=?`, original.EffectKey)
				require.NoError(t, err)
			case "owner marker":
				raw, err := f.store.RuntimeVersion(t.Context(), groupKey(f.snapshot))
				require.NoError(t, err)
				var snap workspace.WriteGroupAppliedSnapshot
				require.NoError(t, json.Unmarshal(raw, &snap))
				_, owner, ok := workspace.LongStorageIdentity(snap.Group)
				require.True(t, ok)
				for i, column := range snap.RuntimeLayout.Columns {
					if column.Name == owner {
						snap.RuntimeLayout.Columns = append(snap.RuntimeLayout.Columns[:i], snap.RuntimeLayout.Columns[i+1:]...)
						break
					}
				}
				raw, err = json.Marshal(snap)
				require.NoError(t, err)
				_, err = f.local.ExecContext(t.Context(), `UPDATE wg_runtime_versions SET payload=?`, string(raw))
				require.NoError(t, err)
			case "scope alias":
				raw, err := f.store.RuntimeVersion(t.Context(), groupKey(f.snapshot))
				require.NoError(t, err)
				raw = append(raw[:len(raw)-1], []byte(`,"Workspace_ID":"foreign"}`)...)
				_, err = f.local.ExecContext(t.Context(), `UPDATE wg_runtime_versions SET payload=?`, string(raw))
				require.NoError(t, err)
			}
			result, err := counterFinalizedSender(f).Deliver(t.Context(), original.EffectKey)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateBlocked, result.State)
			require.Equal(t, code, result.Code)
			require.Equal(t, 0, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
			require.Equal(t, 0, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_receipts`))
			if fault == "destination" {
				require.Equal(t, 0, f.count(f.plainTarget(), `SELECT COUNT(*) FROM forged_sink`))
			}
		})
	}
}

func TestCounterFinalizedStorageReadFailureDoesNotAuthorizeDelivery(t *testing.T) {
	f := newLongPipeline(t)
	require.NoError(t, f.reading(0, "finalized", at(1), 100))
	f.closeUntil(at(600))
	original := firstCounterEffect(t, f)
	damageCounterProof(t, f, "counter-replay-proof-missing")
	before := counterDurableState(t, f.store)
	injected := errors.New("private storage driver failure")
	calls := 0
	reject := func(query string) error {
		if strings.Contains(query, "SELECT payload FROM wg_runtime_versions") {
			calls++
			return injected
		}
		return nil
	}
	f.localFault.Reject.Store(&reject)
	result, err := counterFinalizedSender(f).Deliver(t.Context(), original.EffectKey)
	f.localFault.Reject.Store(nil)
	require.Positive(t, calls)
	require.ErrorIs(t, err, injected)
	require.Equal(t, groupdelivery.StatePending, result.State)
	require.Equal(t, before, counterDurableState(t, f.store))
	require.Equal(t, 0, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
	result, err = counterFinalizedSender(f).Deliver(t.Context(), original.EffectKey)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StateCommitted, result.State)
}
