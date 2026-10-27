package groupdelivery

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

func legacyAssociationFixture(t *testing.T) (*Store, string, []byte, []counterproof.CounterBindingReference) {
	t.Helper()
	store, _, path := newTestStore(t)
	store.now = func() time.Time { return at(100) }
	for i, digits := range []string{"12000.0", "12000.4"} {
		value, err := measurement.NewDecimal(digits)
		require.NoError(t, err)
		sample := sampleAt([]string{"c-1", "c-2"}[i], `["device-1","point-1","tag-1"]`, i+2, 0)
		sample.Value = value
		require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
	}
	const descriptor = `{"workspace_id":"workspace-A","group_id":"group-G","applied_revision":"rev-1","group":{"id":"group-G","workspace_id":"workspace-A","revision":"rev-1","applied_revision":"rev-1","layout_version":"managed_long_v1","destination":{"storage_strategy":"managed"},"recording_profile":{"kind":"cumulative"},"members":[{"device_id":"device-1","point_id":"point-1","tag_id":"tag-1","measurement_id":"measurement-1","source_revision":"src-1","mapping_revision":"map-1","source_value_type":"int64","result_value_type":"decimal","exact_recording_source":{"measurement_id":"measurement-1","definition_revision":"def-1","source_binding_revision":"binding-1","series_epoch":"epoch-1","original_data_type":"int64","scale":"0.1","unit":"kWh","semantic_kind":"counter"}}]},"runtime_layout":{"tag_types":{"tag-1":"int64"}}}`
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, []byte(descriptor)))
	proofs, err := counterproof.ResolveSnapshot([]byte(descriptor), counterproof.SnapshotScope{WorkspaceID: "workspace-A", GroupID: "group-G", AppliedRevision: "rev-1"})
	require.NoError(t, err)
	refs := []counterproof.CounterBindingReference{}
	for _, proof := range proofs {
		ref, err := proof.Reference()
		require.NoError(t, err)
		refs = append(refs, ref)
	}
	return store, path, []byte(descriptor), refs
}

func associationJournal(t *testing.T, store *Store) []string {
	t.Helper()
	rows, err := store.DB().Query(`SELECT sample_id || '|' || payload || '|' || payload_digest || '|' || consumed || '|' || created_at FROM wg_delivery_samples ORDER BY seq`)
	require.NoError(t, err)
	defer rows.Close()
	var values []string
	for rows.Next() {
		var value string
		require.NoError(t, rows.Scan(&value))
		values = append(values, value)
	}
	require.NoError(t, rows.Err())
	return values
}

func TestLegacyCounterProofAssociationKeepsOriginalAcceptance(t *testing.T) {
	store, path, descriptor, refs := legacyAssociationFixture(t)
	before := associationJournal(t, store)
	const original = `{"sample_id":"c-1","member_key":"[\"device-1\",\"point-1\",\"tag-1\"]","observed_at":"2026-01-01T00:00:02Z","quality":"good","value":{"type":"decimal","encoding":"decimal-string","value":"12000.0"},"source_revision":"src-1","mapping_revision":"map-1"}`
	var payload, digest string
	require.NoError(t, store.DB().QueryRow(`SELECT payload,payload_digest FROM wg_delivery_samples WHERE sample_id='c-1'`).Scan(&payload, &digest))
	require.Equal(t, original, payload)
	sum := sha256.Sum256([]byte(original))
	require.Equal(t, hex.EncodeToString(sum[:]), digest)
	store.now = func() time.Time { return at(200) }
	require.NoError(t, store.AssociateLegacyCounterProofs(t.Context(), testKey, refs))
	require.Equal(t, 2, countRows(t, store.DB(), "wg_counter_proof_associations"), "both original acceptances need durable proof audit")
	require.Equal(t, before, associationJournal(t, store))
	stored, err := store.RuntimeVersion(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, descriptor, stored)
	var member, proofDigest, evidence, accepted string
	var version, associationVersion, seq int
	require.NoError(t, store.DB().QueryRow(`SELECT proof_member_key,proof_version,proof_digest,evidence_kind,accepted_at,association_version,sample_seq FROM wg_counter_proof_associations WHERE sample_id='c-1'`).Scan(&member, &version, &proofDigest, &evidence, &accepted, &associationVersion, &seq))
	require.Equal(t, `["workspace-A","group-G","rev-1","device-1","point-1","tag-1"]`, member)
	require.Equal(t, 1, version)
	require.Equal(t, 1, associationVersion)
	require.Equal(t, "legacy-exact-source-v1", evidence)
	require.Equal(t, "2026-01-01T00:01:40.000000000Z", accepted)
	require.Equal(t, 1, seq)
	require.Equal(t, "e8eba372c5147ae6486978cbd7771c47d4592216525ed0adce5d10b39a552551", proofDigest)
	var runtimeDigest, sampleDigest, associatedAt string
	require.NoError(t, store.DB().QueryRow(`SELECT runtime_digest,sample_digest,associated_at FROM wg_counter_proof_associations WHERE sample_id='c-1'`).Scan(&runtimeDigest, &sampleDigest, &associatedAt))
	evidenceSum := sha256.Sum256(descriptor)
	require.Equal(t, hex.EncodeToString(evidenceSum[:]), runtimeDigest)
	require.Equal(t, digest, sampleDigest)
	require.Equal(t, "2026-01-01T00:03:20.000000000Z", associatedAt)
	audit := associationAudit(t, store)
	require.NotContains(t, strings.Join(audit, ""), "12000.", "audit stores references, not raw readings")
	require.NoError(t, store.DB().Close())
	store = NewStore(openStoreDB(t, path))
	require.NoError(t, store.AssociateLegacyCounterProofs(t.Context(), testKey, refs))
	restored, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Len(t, restored.Samples, 2)
	resend := restored.Samples[0]
	resend.CounterProof = refs[0]
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, resend))
	require.Equal(t, before, associationJournal(t, store))
	require.Equal(t, 2, countRows(t, store.DB(), "wg_counter_proof_associations"))
	require.Equal(t, audit, associationAudit(t, store))
	// The old codec is still the stored truth, not a rewritten versioned sample.
	var wire map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(payload), &wire))
	require.NotContains(t, wire, "counter_proof")
}

func associationAudit(t *testing.T, store *Store) []string {
	t.Helper()
	rows, err := store.DB().Query(`SELECT workspace_id || '|' || group_id || '|' || group_revision || '|' || sample_id || '|' || sample_seq || '|' || sample_digest || '|' || accepted_at || '|' || member_key || '|' || proof_member_key || '|' || proof_version || '|' || proof_digest || '|' || runtime_digest || '|' || evidence_kind || '|' || association_version || '|' || associated_at FROM wg_counter_proof_associations ORDER BY sample_seq`)
	require.NoError(t, err)
	defer rows.Close()
	var values []string
	for rows.Next() {
		var value string
		require.NoError(t, rows.Scan(&value))
		values = append(values, value)
	}
	require.NoError(t, rows.Err())
	return values
}
