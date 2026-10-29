package grouppipeline

import (
	"encoding/json"
	"strings"
	"testing"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestCounterBindingRecoveryRefusesUntrustedProofWithoutConsuming(t *testing.T) {
	for name, tc := range map[string]struct {
		mutate func(map[string]any)
		reason string
	}{
		"missing": {func(s map[string]any) { delete(s, "counter_bindings") }, "counter-replay-proof-missing"},
		"tampered unit": {func(s map[string]any) {
			for _, p := range s["counter_bindings"].(map[string]any) {
				p.(map[string]any)["unit"] = "MWh"
			}
		}, "counter-replay-proof-conflict"},
		"unsupported": {func(s map[string]any) {
			for _, p := range s["counter_bindings"].(map[string]any) {
				p.(map[string]any)["version"] = 77
			}
		}, "counter-replay-proof-unsupported"},
		"foreign group": {func(s map[string]any) { s["group"].(map[string]any)["workspace_id"] = "other" }, "counter-replay-proof-conflict"},
		"changed member mapping": {func(s map[string]any) {
			s["group"].(map[string]any)["members"].([]any)[0].(map[string]any)["mapping_revision"] = "mapping-2"
		}, "counter-replay-proof-conflict"},
		"nested case alias": {func(s map[string]any) {
			s["group"].(map[string]any)["workspace_id"] = "other"
			s["group"].(map[string]any)["Workspace_ID"] = "workspace-A"
		}, "counter-replay-proof-conflict"},
		"nested Unicode alias": {func(s map[string]any) {
			m := s["group"].(map[string]any)["members"].([]any)[0].(map[string]any)
			m["source_revision"] = "other"
			m["ſource_revision"] = "src-1"
		}, "counter-replay-proof-conflict"},
	} {
		t.Run(name, func(t *testing.T) {
			f := newLongPipeline(t)
			require.NoError(t, f.reading(0, "c-1", at(1), 100))
			f.closeUntil(at(600))
			require.NoError(t, f.reading(0, "c-2", at(601), 105))
			key := groupKey(f.snapshot)
			before, err := f.store.Restore(t.Context(), key)
			require.NoError(t, err)
			require.Len(t, before.Samples, 1)
			require.Len(t, before.Baselines, 1)
			var sample, digest, baseline, checkpoint string
			require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT payload,payload_digest FROM wg_delivery_samples WHERE consumed=0`).Scan(&sample, &digest))
			require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_counter_baselines`).Scan(&baseline))
			require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT next_close FROM wg_delivery_checkpoints`).Scan(&checkpoint))
			payload, err := f.store.RuntimeVersion(t.Context(), key)
			require.NoError(t, err)
			var fields map[string]any
			require.NoError(t, json.Unmarshal(payload, &fields))
			tc.mutate(fields)
			corrupted, err := json.Marshal(fields)
			require.NoError(t, err)
			_, err = f.local.ExecContext(t.Context(), `UPDATE wg_runtime_versions SET payload=?`, string(corrupted))
			require.NoError(t, err)
			status := f.reconcileStatus()
			require.Equal(t, StateBlocked, status.State)
			require.Equal(t, tc.reason, status.Reason)
			f.closeUntil(at(1200))
			var afterSample, afterDigest, afterBaseline, afterCheckpoint string
			require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT payload,payload_digest FROM wg_delivery_samples WHERE consumed=0`).Scan(&afterSample, &afterDigest))
			require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_counter_baselines`).Scan(&afterBaseline))
			require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT next_close FROM wg_delivery_checkpoints`).Scan(&afterCheckpoint))
			require.Equal(t, sample, afterSample)
			require.Equal(t, digest, afterDigest)
			require.Equal(t, baseline, afterBaseline)
			require.Equal(t, checkpoint, afterCheckpoint)
			after, err := f.store.RuntimeVersion(t.Context(), key)
			require.NoError(t, err)
			require.Equal(t, corrupted, after)
		})
	}
}

func TestCounterBindingLegacyExactJournalUsesOriginalBytes(t *testing.T) {
	f := newLongPipelineOn(t, nil, true)
	key := groupKey(f.snapshot)
	raw, err := f.store.RuntimeVersion(t.Context(), key)
	require.NoError(t, err)
	var snap workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(raw, &snap))
	snap.CounterBindings = nil
	legacy, err := json.Marshal(snap)
	require.NoError(t, err)
	// Seed an actual pre-format journal while its legacy descriptor is absent.
	// The complete original exact descriptor is then restored, never current defs.
	_, err = f.local.ExecContext(t.Context(), `DELETE FROM wg_runtime_versions`)
	require.NoError(t, err)
	value, err := measurement.NewDecimal("100.0")
	require.NoError(t, err)
	sample := snapshot.Sample{SampleID: "legacy-c-1", MemberKey: `["device-1","point-1","tag-kwh-1"]`, ObservedAt: at(1), Quality: schema.QualityGood, Value: value, SourceRevision: "src-1", MappingRevision: "map-1"}
	require.NoError(t, f.store.AppendSample(t.Context(), key, at(0), sample))
	require.NoError(t, f.store.SaveRuntimeVersion(t.Context(), key, legacy))
	var original, digest string
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT payload,payload_digest FROM wg_delivery_samples`).Scan(&original, &digest))
	require.NotContains(t, original, "counter_proof")
	require.NotContains(t, original, `"version"`)
	reader := &unavailableExactCounterReader{}
	f.pipe.deps.Counters = reader
	f.pipe.boundaries = map[string]*managed{}
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	var associations int
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_counter_proof_associations`).Scan(&associations))
	require.Equal(t, 1, associations, "legacy activation must durably associate before duplicate ACK")
	require.Equal(t, StateActive, f.statusOf().State)
	require.NoError(t, f.reading(0, "legacy-c-1", at(1), 1000), "identical resend is not a new association or rewritten identity")
	f.closeUntil(at(600))
	require.Zero(t, reader.calls)
	var after, afterDigest string
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT payload,payload_digest FROM wg_delivery_samples`).Scan(&after, &afterDigest))
	require.Equal(t, original, after)
	require.Equal(t, digest, afterDigest)
	persisted, err := f.store.RuntimeVersion(t.Context(), key)
	require.NoError(t, err)
	require.Equal(t, legacy, persisted)
	var output string
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_outbox WHERE payload LIKE '%legacy-c-1%'`).Scan(&output))
	row, err := groupdelivery.DecodeRowPayload([]byte(output))
	require.NoError(t, err)
	cells := map[string]any{}
	for _, cell := range row.Cells {
		cells[cell.Column] = cell.Value
	}
	require.Equal(t, "100.0", cells["value_decimal"])
	require.Equal(t, "kWh", cells["unit"])
	require.Equal(t, "epoch-1", cells["series_epoch"])
	require.False(t, strings.Contains(string(persisted), "counter_bindings"), "association/audit migration remains separate")
}

func TestCounterBindingLegacyUnscaledJournalCannotBorrowCurrentDefinition(t *testing.T) {
	f := newLongPipeline(t)
	key := groupKey(f.snapshot)
	raw, err := f.store.RuntimeVersion(t.Context(), key)
	require.NoError(t, err)
	var snap workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(raw, &snap))
	snap.CounterBindings = nil
	legacy, err := json.Marshal(snap)
	require.NoError(t, err)
	_, err = f.local.ExecContext(t.Context(), `DELETE FROM wg_runtime_versions`)
	require.NoError(t, err)
	sample := snapshot.Sample{SampleID: "legacy-c-1", MemberKey: `["device-1","point-1","tag-kwh-1"]`, ObservedAt: at(1), Quality: schema.QualityGood, Value: measurement.NewUint64(100), SourceRevision: "src-1", MappingRevision: "map-1"}
	require.NoError(t, f.store.AppendSample(t.Context(), key, at(0), sample))
	require.NoError(t, f.store.SaveRuntimeVersion(t.Context(), key, legacy))
	var payload, digest string
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT payload,payload_digest FROM wg_delivery_samples`).Scan(&payload, &digest))
	// Current editable def-2 is not evidence for this original legacy input.
	current := f.counters.sources[longKwhFirst]
	current.DefinitionRevision = "def-2"
	f.counters.sources[longKwhFirst] = current
	reader := &changedCounterReader{sources: f.counters.sources}
	f.pipe.deps.Counters = reader
	f.pipe.boundaries = map[string]*managed{}
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	require.Equal(t, "counter-replay-proof-missing", f.statusOf().Reason)
	require.Empty(t, f.pipe.boundaries)
	require.Zero(t, reader.calls)
	require.False(t, f.pipe.WantsSample("device-1", "point-1", longKwhFirst), "there is no assembler whose duplicate path could ACK unproven history")
	var after, afterDigest string
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT payload,payload_digest FROM wg_delivery_samples WHERE consumed=0`).Scan(&after, &afterDigest))
	require.Equal(t, payload, after)
	require.Equal(t, digest, afterDigest)
	require.Zero(t, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_checkpoints`))
}
