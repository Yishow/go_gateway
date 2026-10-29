package groupdelivery

import (
	"encoding/json"
	"strings"
	"testing"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

func journalProofFixture(t *testing.T) (counterproof.CounterBindingProof, snapshot.Sample, []byte) {
	t.Helper()
	proof, err := counterproof.NewCounterBindingProof(counterproof.CounterBindingFacts{
		Member: proofSampleReference().Member, MeasurementID: "measurement-1", DefinitionRevision: "def-1", SourceBindingRevision: "binding-1",
		SeriesEpoch: "epoch-1", SemanticKind: measurement.SemanticKindCounter, Unit: "kWh",
		SourceRevision: "src-1", MappingRevision: "map-1", ScaleMode: counterproof.CounterScaleUnscaled,
		ExactScale: "1", SourceValueType: schema.DataTypeInt64, ResultValueType: measurement.ExactInt64,
	})
	require.NoError(t, err)
	sample := sampleAt("counter-1", `["device-1","point-1","tag-1"]`, 2, 7)
	sample.CounterProof = counterproof.CounterBindingReference{Version: proof.Version, Member: proof.Member, Digest: proof.Digest}
	payload, err := json.Marshal(map[string]any{
		"workspace_id": "workspace-A", "group_id": "group-G", "applied_revision": "rev-1",
		"counter_bindings": counterproof.CounterBindingProofs{proof.Member.Key(): proof},
		"group":            map[string]any{"id": "group-G", "workspace_id": "workspace-A", "revision": "rev-1", "applied_revision": "rev-1", "layout_version": "managed_long_v1", "recording_profile": map[string]any{"kind": "cumulative"}, "members": []any{map[string]any{"device_id": "device-1", "point_id": "point-1", "tag_id": "tag-1", "measurement_id": "measurement-1", "source_revision": "src-1", "mapping_revision": "map-1"}}},
		"runtime_layout":   map[string]any{"tag_types": map[string]any{"tag-1": "int64"}},
	})
	require.NoError(t, err)
	return proof, sample, payload
}

func TestCounterProofJournalRequiresActualPersistedProof(t *testing.T) {
	for name, mutate := range map[string]func(*snapshot.Sample, *[]byte){
		"missing snapshot":  func(_ *snapshot.Sample, p *[]byte) { *p = nil },
		"missing reference": func(s *snapshot.Sample, _ *[]byte) { s.CounterProof = counterproof.CounterBindingReference{} },
		"wrong digest":      func(s *snapshot.Sample, _ *[]byte) { s.CounterProof.Digest = strings.Repeat("b", 64) },
		"wrong workspace":   func(s *snapshot.Sample, _ *[]byte) { s.CounterProof.Member.WorkspaceID = "other" },
		"wrong group":       func(s *snapshot.Sample, _ *[]byte) { s.CounterProof.Member.GroupID = "other" },
		"wrong revision":    func(s *snapshot.Sample, _ *[]byte) { s.CounterProof.Member.AppliedRevision = "other" },
		"wrong member":      func(s *snapshot.Sample, _ *[]byte) { s.MemberKey = `["device-1","point-1","other"]` },
		"wrong mapping":     func(s *snapshot.Sample, _ *[]byte) { s.MappingRevision = "map-2" },
		"wrong source":      func(s *snapshot.Sample, _ *[]byte) { s.SourceRevision = "src-2" },
		"wrong result type": func(s *snapshot.Sample, _ *[]byte) { s.Value = measurement.NewUint64(7) },
		"tampered snapshot": func(_ *snapshot.Sample, p *[]byte) { *p = []byte(strings.Replace(string(*p), "kWh", "MWh", 1)) },
		"foreign snapshot": func(_ *snapshot.Sample, p *[]byte) {
			*p = []byte(strings.Replace(string(*p), `"workspace_id":"workspace-A"`, `"workspace_id":"other"`, 1))
		},
	} {
		t.Run(name, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			_, sample, payload := journalProofFixture(t)
			mutate(&sample, &payload)
			if payload != nil {
				require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, payload))
			}
			err := store.AppendSample(t.Context(), testKey, t0, sample)
			require.Error(t, err, "unverified counter evidence cannot be ACKed")
			var typed *counterproof.CounterBindingProofError
			require.ErrorAs(t, err, &typed)
			require.Zero(t, countRows(t, db, "wg_delivery_samples"))
		})
	}
}

func TestCounterProofJournalCommitsReferenceAndRefusesChangedIdentity(t *testing.T) {
	store, db, _ := newTestStore(t)
	_, sample, payload := journalProofFixture(t)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, payload))
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
	restored, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
	require.Equal(t, sample.CounterProof, restored.Samples[0].CounterProof)
	sample.CounterProof.Digest = strings.Repeat("b", 64)
	require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, sample), ErrSampleConflict)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))
}

func TestCounterProofJournalLegacyResendKeepsOriginalBytes(t *testing.T) {
	store, db, _ := newTestStore(t)
	_, sample, runtimePayload := journalProofFixture(t)
	legacy := sample
	legacy.CounterProof = counterproof.CounterBindingReference{}
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, legacy))
	var before, digest string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT payload, payload_digest FROM wg_delivery_samples`).Scan(&before, &digest))
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, runtimePayload))
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample), "a duplicate does not migrate or relabel legacy history")
	var after, afterDigest string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT payload, payload_digest FROM wg_delivery_samples`).Scan(&after, &afterDigest))
	require.Equal(t, before, after)
	require.Equal(t, digest, afterDigest)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))
}

func TestCounterProofJournalRestoreRejectsChangedPayloadOrDigest(t *testing.T) {
	for _, change := range []string{`payload_digest='wrong'`, `payload=replace(payload, '"value":"7"', '"value":"8"')`} {
		t.Run(change, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			_, sample, payload := journalProofFixture(t)
			require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, payload))
			require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
			_, err := db.ExecContext(t.Context(), `UPDATE wg_delivery_samples SET `+change)
			require.NoError(t, err)
			_, err = store.Restore(t.Context(), testKey)
			require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict)
			require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"), "refusal retains the original damaged evidence")
		})
	}
}

func TestCounterProofJournalNeverTrustsIgnoredReplacementSnapshot(t *testing.T) {
	store, db, _ := newTestStore(t)
	proof, sample, original := journalProofFixture(t)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, original))
	facts := proof.CounterBindingFacts
	facts.Unit = "MWh"
	changed, err := counterproof.NewCounterBindingProof(facts)
	require.NoError(t, err)
	replacement := strings.Replace(string(original), proof.Digest, changed.Digest, 1)
	replacement = strings.Replace(replacement, "kWh", "MWh", 1)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, []byte(replacement)), "first descriptor wins")
	sample.CounterProof.Digest = changed.Digest
	require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, sample), counterproof.ErrCounterBindingProofConflict)
	require.Zero(t, countRows(t, db, "wg_delivery_samples"))
	stored, err := store.RuntimeVersion(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, original, stored)
}

func TestCounterProofJournalRestoreRequiresOriginalReferenceEvidence(t *testing.T) {
	store, db, _ := newTestStore(t)
	_, sample, payload := journalProofFixture(t)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, payload))
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
	_, err := db.ExecContext(t.Context(), `DELETE FROM wg_runtime_versions`)
	require.NoError(t, err)
	_, err = store.Restore(t.Context(), testKey)
	require.ErrorIs(t, err, counterproof.ErrCounterBindingProofMissing)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))
}

func TestCounterProofJournalRejectsAmbiguousSnapshotFields(t *testing.T) {
	for name, alias := range map[string]string{
		"scope case alias":    `,"Workspace_ID":"other"}`,
		"proof case alias":    `,"Counter_Bindings":null}`,
		"scope Unicode alias": `,"workſpace_id":"other"}`,
	} {
		t.Run(name, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			_, sample, payload := journalProofFixture(t)
			payload = append(payload[:len(payload)-1], []byte(alias)...)
			require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, payload))
			require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, sample), counterproof.ErrCounterBindingProofConflict)
			require.Zero(t, countRows(t, db, "wg_delivery_samples"))
		})
	}
}

func TestCounterProofJournalAllowsUnicodeDescriptorMapKeys(t *testing.T) {
	store, _, _ := newTestStore(t)
	_, sample, payload := journalProofFixture(t)
	payload = []byte(strings.Replace(string(payload), `"tag_types":{"tag-1":"int64"}`, `"tag_types":{"tag-1":"int64","量測ſ":"int64"}`, 1))
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, payload))
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
}
