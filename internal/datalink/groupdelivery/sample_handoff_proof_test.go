package groupdelivery

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/counterproof"
	"strings"
	"testing"
	"time"
)

func handoffProofDescriptor(t *testing.T, key GroupKey, proof *counterproof.CounterBindingProof, managed bool) []byte {
	t.Helper()
	member := map[string]any{"device_id": "device-1", "point_id": "point-1", "tag_id": "tag-1", "source_revision": "src-1", "mapping_revision": "map-1"}
	group := map[string]any{"id": key.GroupID, "workspace_id": key.WorkspaceID, "revision": key.GroupRevision, "applied_revision": key.GroupRevision, "members": []any{member}, "row_policy": map[string]any{"interval_seconds": 10}}
	bindings := counterproof.CounterBindingProofs{}
	if managed {
		group["recording_profile"] = map[string]any{"kind": "instantaneous"}
		member["value_role"], member["source_value_type"], member["result_value_type"] = "gauge", "int64", "decimal"
	}
	if proof != nil {
		group["recording_profile"] = map[string]any{"kind": "cumulative"}
		member["measurement_id"] = proof.MeasurementID
		bindings[proof.Member.Key()] = *proof
	}
	raw, err := json.Marshal(map[string]any{"workspace_id": key.WorkspaceID, "group_id": key.GroupID, "applied_revision": key.GroupRevision, "group": group, "runtime_layout": map[string]any{"tag_types": map[string]string{"tag-1": "int64"}}, "counter_bindings": bindings})
	require.NoError(t, err)
	return raw
}

func TestHandoffProofLegacyRawInputCompatibleWithMeanOutput(t *testing.T) {
	store, db, _ := newTestStore(t)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, handoffProofDescriptor(t, testKey, nil, false)))
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), otherKey, handoffProofDescriptor(t, otherKey, nil, true)))
	sample := sampleAt("raw", ownershipMember, 23, 7)
	payload, digest, err := encodeSample(sample)
	require.NoError(t, err)
	input := handoffProofInput{Origin: testKey, Sequence: 1, Payload: string(payload), Digest: digest, AcceptedAt: timestamp(at(21))}
	evidence, err := validateHandoffProof(t.Context(), db, input, otherKey)
	require.NoError(t, err)
	require.Equal(t, 10*time.Second, evidence.TargetInterval)
	require.True(t, evidence.TargetReference.IsZero())
	require.NotEmpty(t, evidence.OriginRuntimeDigest)
	require.NotEmpty(t, evidence.TargetRuntimeDigest)
}

func handoffCounterProofFixture(t *testing.T, legacy bool) (*Store, handoffProofInput, counterproof.CounterBindingProof, counterproof.CounterBindingProof) {
	t.Helper()
	store, db, _ := newTestStore(t)
	origin, sample, _ := journalProofFixture(t)
	if legacy {
		sample.CounterProof = counterproof.CounterBindingReference{}
		require.NoError(t, store.AppendSample(t.Context(), testKey, at(0), sample))
	}
	facts := origin.CounterBindingFacts
	facts.Member.GroupID, facts.Member.AppliedRevision = otherKey.GroupID, otherKey.GroupRevision
	target, err := counterproof.NewCounterBindingProof(facts)
	require.NoError(t, err)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, handoffProofDescriptor(t, testKey, &origin, false)))
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), otherKey, handoffProofDescriptor(t, otherKey, &target, false)))
	if !legacy {
		require.NoError(t, store.AppendSample(t.Context(), testKey, at(0), sample))
	} else {
		ref, err := origin.Reference()
		require.NoError(t, err)
		require.NoError(t, store.AssociateLegacyCounterProofs(t.Context(), testKey, []counterproof.CounterBindingReference{ref}))
	}
	input := handoffProofInput{Origin: testKey}
	require.NoError(t, db.QueryRow("SELECT seq,payload,payload_digest,created_at FROM wg_delivery_samples").Scan(&input.Sequence, &input.Payload, &input.Digest, &input.AcceptedAt))
	return store, input, origin, target
}

func TestHandoffProofCounterPreservesOriginalReference(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		t.Run(fmt.Sprint("legacy=", legacy), func(t *testing.T) {
			store, input, origin, target := handoffCounterProofFixture(t, legacy)
			evidence, err := validateHandoffProof(t.Context(), store.db, input, otherKey)
			require.NoError(t, err)
			wantOrigin, err := origin.Reference()
			require.NoError(t, err)
			wantTarget, err := target.Reference()
			require.NoError(t, err)
			require.Equal(t, wantOrigin, evidence.OriginalReference)
			require.Equal(t, wantTarget, evidence.TargetReference)
			require.NotEqual(t, wantOrigin, wantTarget)
			var payload, digest string
			require.NoError(t, store.db.QueryRow("SELECT payload,payload_digest FROM wg_delivery_samples").Scan(&payload, &digest))
			require.Equal(t, input.Payload, payload)
			require.Equal(t, input.Digest, digest)
			decoded, err := decodeSample([]byte(payload))
			require.NoError(t, err)
			require.Equal(t, legacy, decoded.CounterProof.IsZero())
		})
	}
}

func TestHandoffProofCounterRejectsChangedMeaning(t *testing.T) {
	for name, change := range map[string]func(*counterproof.CounterBindingFacts){
		"epoch":      func(f *counterproof.CounterBindingFacts) { f.SeriesEpoch = "epoch-2" },
		"unit":       func(f *counterproof.CounterBindingFacts) { f.Unit = "MWh" },
		"definition": func(f *counterproof.CounterBindingFacts) { f.DefinitionRevision = "def-2" },
		"binding":    func(f *counterproof.CounterBindingFacts) { f.SourceBindingRevision = "binding-2" },
		"source":     func(f *counterproof.CounterBindingFacts) { f.SourceRevision = "src-2" },
		"mapping":    func(f *counterproof.CounterBindingFacts) { f.MappingRevision = "map-2" },
		"policy": func(f *counterproof.CounterBindingFacts) {
			f.CounterPolicy.RolloverModulus, f.CounterPolicy.MaxRatePerSecond = "65536", "100"
		},
		"measurement": func(f *counterproof.CounterBindingFacts) { f.MeasurementID = "measurement-2" },
	} {
		t.Run(name, func(t *testing.T) {
			store, input, _, target := handoffCounterProofFixture(t, false)
			facts := target.CounterBindingFacts
			change(&facts)
			proof, err := counterproof.NewCounterBindingProof(facts)
			require.NoError(t, err)
			raw := handoffProofDescriptor(t, otherKey, &proof, false)
			_, err = store.db.Exec("UPDATE wg_runtime_versions SET payload=? WHERE group_id=?", string(raw), otherKey.GroupID)
			require.NoError(t, err)
			_, err = validateHandoffProof(t.Context(), store.db, input, otherKey)
			require.ErrorIs(t, err, ErrHandoffProofConflict)
			require.Equal(t, "recording-handoff-proof-conflict", err.Error())
		})
	}
}

func TestHandoffProofLegacyRequiresExactExistingAssociation(t *testing.T) {
	for _, change := range []string{
		"DELETE FROM wg_counter_proof_associations",
		"UPDATE wg_counter_proof_associations SET sample_digest='wrong'",
		"UPDATE wg_counter_proof_associations SET accepted_at='2026-01-01T00:00:00Z'",
		"UPDATE wg_counter_proof_associations SET runtime_digest='wrong'",
		"UPDATE wg_counter_proof_associations SET proof_digest='wrong'",
		"UPDATE wg_counter_proof_associations SET evidence_kind='invented'",
	} {
		t.Run(change, func(t *testing.T) {
			store, input, _, _ := handoffCounterProofFixture(t, true)
			_, err := store.db.Exec(change)
			require.NoError(t, err)
			_, err = validateHandoffProof(t.Context(), store.db, input, otherKey)
			if strings.HasPrefix(change, "DELETE") {
				require.ErrorIs(t, err, ErrOwnershipUnavailable)
			} else {
				require.ErrorIs(t, err, ErrHandoffProofConflict)
			}
			var payload string
			require.NoError(t, store.db.QueryRow("SELECT payload FROM wg_delivery_samples").Scan(&payload))
			require.Equal(t, input.Payload, payload)
		})
	}
}

func TestHandoffProofUnknownAndConflictingNoncounterInputs(t *testing.T) {
	for name, mutate := range map[string]string{
		"missing layout":      "json_remove(payload,'$.runtime_layout')",
		"unknown input":       "json_set(payload,'$.runtime_layout.tag_types.tag-1','blob')",
		"contradicting input": "json_set(payload,'$.group.members[0].source_value_type','uint64')",
		"source revision":     "json_set(payload,'$.group.members[0].source_revision','src-2')",
		"mapping revision":    "json_set(payload,'$.group.members[0].mapping_revision','map-2')",
		"measurement":         "json_set(payload,'$.group.members[0].measurement_id','unexpected')",
		"missing interval":    "json_remove(payload,'$.group.row_policy')",
		"counter conversion":  "json_set(payload,'$.group.recording_profile.kind','cumulative')",
	} {
		t.Run(name, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, handoffProofDescriptor(t, testKey, nil, false)))
			require.NoError(t, store.SaveRuntimeVersion(t.Context(), otherKey, handoffProofDescriptor(t, otherKey, nil, true)))
			_, err := db.Exec("UPDATE wg_runtime_versions SET payload="+mutate+" WHERE group_id=?", otherKey.GroupID)
			require.NoError(t, err)
			payload, digest, err := encodeSample(sampleAt("raw", ownershipMember, 23, 7))
			require.NoError(t, err)
			_, err = validateHandoffProof(t.Context(), db, handoffProofInput{Origin: testKey, Sequence: 1, Payload: string(payload), Digest: digest, AcceptedAt: timestamp(at(21))}, otherKey)
			if name == "missing layout" || name == "missing interval" || name == "counter conversion" {
				require.ErrorIs(t, err, ErrOwnershipUnavailable)
			} else {
				require.ErrorIs(t, err, ErrHandoffProofConflict)
			}
		})
	}
}

func TestHandoffProofUnreadableAuthorityIsSafe(t *testing.T) {
	store, input, _, _ := handoffCounterProofFixture(t, false)
	_, err := store.db.Exec("DELETE FROM wg_runtime_versions WHERE group_id=?", otherKey.GroupID)
	require.NoError(t, err)
	_, err = validateHandoffProof(t.Context(), store.db, input, otherKey)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	require.Equal(t, "recording-ownership-unavailable", err.Error())
}

func compactHandoffProofFixture(t *testing.T, legacy bool) (*Store, handoffReceipt) {
	t.Helper()
	store, input, _, _ := handoffCounterProofFixture(t, legacy)
	evidence, err := validateHandoffProof(t.Context(), store.db, input, otherKey)
	require.NoError(t, err)
	sample, err := decodeSample([]byte(input.Payload))
	require.NoError(t, err)
	return store, handoffReceipt{Origin: input.Origin, Target: otherKey, OriginSeq: input.Sequence, SampleID: sample.SampleID,
		Member: sample.MemberKey, Digest: input.Digest, AcceptedAt: input.AcceptedAt, SampleProof: sample.CounterProof,
		OriginalProof: evidence.OriginalReference, TargetProof: evidence.TargetReference,
		OriginRuntimeDigest: evidence.OriginRuntimeDigest, TargetRuntimeDigest: evidence.TargetRuntimeDigest}
}

func TestCompactHandoffProofRejectsReferenceMutation(t *testing.T) {
	for name, mutate := range map[string]func(*handoffReceipt){
		"target reference":   func(r *handoffReceipt) { r.TargetProof.Digest = strings.Repeat("b", 64) },
		"original reference": func(r *handoffReceipt) { r.OriginalProof.Digest = strings.Repeat("b", 64) },
		"sample reference":   func(r *handoffReceipt) { r.SampleProof.Digest = strings.Repeat("b", 64) },
	} {
		t.Run(name, func(t *testing.T) {
			store, receipt := compactHandoffProofFixture(t, false)
			require.NoError(t, verifyCompactHandoffMeaning(t.Context(), store.db, receipt))
			mutate(&receipt)
			require.ErrorIs(t, verifyCompactHandoffMeaning(t.Context(), store.db, receipt), ErrHandoffProofConflict)
		})
	}
}

func TestCompactHandoffProofPreservesLegacyAssociationAfterReclaim(t *testing.T) {
	store, receipt := compactHandoffProofFixture(t, true)
	_, err := store.db.Exec("DELETE FROM wg_delivery_samples")
	require.NoError(t, err)
	require.NoError(t, verifyCompactHandoffMeaning(t.Context(), store.db, receipt), "valid proof does not depend on reclaimed value bytes")
	_, err = store.db.Exec("DELETE FROM wg_counter_proof_associations")
	require.NoError(t, err)
	require.ErrorIs(t, verifyCompactHandoffMeaning(t.Context(), store.db, receipt), ErrOwnershipUnavailable)
}

func TestCompactHandoffProofMissingReferenceIsNotLegacyEvidence(t *testing.T) {
	store, receipt := compactHandoffProofFixture(t, false)
	receipt.SampleProof = counterproof.CounterBindingReference{}
	require.ErrorIs(t, verifyCompactHandoffMeaning(t.Context(), store.db, receipt), ErrOwnershipUnavailable)
}

func TestCompactHandoffProofRequiresMeaningBeyondMatchingDescriptorHashes(t *testing.T) {
	store, receipt := compactHandoffProofFixture(t, false)
	target, err := readHandoffProof(t.Context(), store.db, otherKey, receipt.Member)
	require.NoError(t, err)
	facts := target.proof.CounterBindingFacts
	facts.Unit = "MWh"
	proof, err := counterproof.NewCounterBindingProof(facts)
	require.NoError(t, err)
	raw := handoffProofDescriptor(t, otherKey, &proof, false)
	_, err = store.db.Exec("UPDATE wg_runtime_versions SET payload=? WHERE group_id=?", string(raw), otherKey.GroupID)
	require.NoError(t, err)
	receipt.TargetRuntimeDigest = handoffHash(string(raw))
	receipt.TargetProof, err = proof.Reference()
	require.NoError(t, err)
	require.ErrorIs(t, verifyCompactHandoffMeaning(t.Context(), store.db, receipt), ErrHandoffProofConflict)
}

func TestCompactHandoffProofNoncounterInputProjection(t *testing.T) {
	store, db, _ := newTestStore(t)
	originRaw, targetRaw := handoffProofDescriptor(t, testKey, nil, false), handoffProofDescriptor(t, otherKey, nil, true)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, originRaw))
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), otherKey, targetRaw))
	receipt := handoffReceipt{Origin: testKey, Target: otherKey, OriginSeq: 1, Member: ownershipMember, AcceptedAt: timestamp(at(21)),
		OriginRuntimeDigest: handoffHash(string(originRaw)), TargetRuntimeDigest: handoffHash(string(targetRaw))}
	require.NoError(t, verifyCompactHandoffMeaning(t.Context(), db, receipt))
	_, err := db.Exec("UPDATE wg_runtime_versions SET payload=json_set(payload,'$.runtime_layout.tag_types.tag-1','uint64','$.group.members[0].source_value_type','uint64') WHERE group_id=?", otherKey.GroupID)
	require.NoError(t, err)
	var changed string
	require.NoError(t, db.QueryRow("SELECT payload FROM wg_runtime_versions WHERE group_id=?", otherKey.GroupID).Scan(&changed))
	receipt.TargetRuntimeDigest = handoffHash(changed)
	require.ErrorIs(t, verifyCompactHandoffMeaning(t.Context(), db, receipt), ErrHandoffProofConflict)
}
