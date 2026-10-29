package groupdelivery

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/snapshot"
)

// verifyClosureCounterProofs runs after the closure transaction takes its write
// lock, so the descriptor cannot change between validation and consumption.
func verifyClosureCounterProofs(ctx context.Context, tx *sql.Tx, key GroupKey, expected []counterproof.CounterBindingReference, before time.Time) error {
	if len(expected) == 0 {
		return nil
	}
	var payload string
	err := tx.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions WHERE workspace_id=? AND group_id=? AND group_revision=?`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return sampleProofError(counterproof.ErrCounterBindingProofMissing)
	}
	if err != nil {
		return err
	}
	evidence, err := counterproof.ResolveSnapshotEvidence([]byte(payload), counterproof.SnapshotScope{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, AppliedRevision: key.GroupRevision})
	if err != nil {
		return err
	}
	if len(evidence.Bindings) != len(expected) {
		return sampleProofError(counterproof.ErrCounterBindingProofConflict)
	}
	seen := make(map[string]bool, len(expected))
	byMember := make(map[string]counterproof.CounterBindingReference, len(expected))
	for _, ref := range expected {
		if err := ref.Validate(); err != nil {
			return err
		}
		member := ref.Member.Key()
		if seen[member] {
			return sampleProofError(counterproof.ErrCounterBindingProofConflict)
		}
		seen[member] = true
		proof, ok := evidence.Bindings[member]
		if !ok || proof.Version != ref.Version || proof.Digest != ref.Digest {
			return sampleProofError(counterproof.ErrCounterBindingProofConflict)
		}
		byMember[ref.Member.SampleMemberKey()] = ref
	}
	return verifyClosureAcceptedSet(ctx, tx, key, before, payload, evidence, byMember)
}

func verifyClosureAcceptedSet(ctx context.Context, tx *sql.Tx, key GroupKey, before time.Time, payload string, evidence counterproof.SnapshotEvidence, refs map[string]counterproof.CounterBindingReference) error {
	samples, err := readAcceptedSet(ctx, tx, key, before)
	if err != nil {
		return err
	}
	for _, accepted := range samples {
		sample := accepted.sample
		if !evidence.Members[sample.MemberKey] {
			return sampleProofError(counterproof.ErrCounterBindingProofConflict)
		}
		ref, counter := refs[sample.MemberKey]
		if !counter {
			if !sample.CounterProof.IsZero() {
				return sampleProofError(counterproof.ErrCounterBindingProofConflict)
			}
			continue
		}
		proof := evidence.Bindings[ref.Member.Key()]
		if sample.SourceRevision != proof.SourceRevision || sample.MappingRevision != proof.MappingRevision || (sample.Value.Type() != "" && sample.Value.Type() != proof.ResultValueType) {
			return sampleProofError(counterproof.ErrCounterBindingProofConflict)
		}
		if !sample.CounterProof.IsZero() {
			if sample.CounterProof != ref {
				return sampleProofError(counterproof.ErrCounterBindingProofConflict)
			}
			continue
		}
		if err := requireLegacyAssociation(ctx, tx, key, accepted, payload, evidence, proof); err != nil {
			return err
		}
	}
	return nil
}

func requireLegacyAssociation(ctx context.Context, tx *sql.Tx, key GroupKey, accepted legacyAcceptedSample, payload string, evidence counterproof.SnapshotEvidence, proof counterproof.CounterBindingProof) error {
	sum := sha256.Sum256([]byte(payload))
	evidenceKind := "runtime-counter-bindings-v1"
	if evidence.LegacyExact {
		evidenceKind = "legacy-exact-source-v1"
	}
	association := legacyProofAssociation{
		WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, GroupRevision: key.GroupRevision, SampleID: accepted.sample.SampleID, SampleSeq: accepted.seq, SampleDigest: accepted.digest, AcceptedAt: accepted.acceptedAt,
		MemberKey: accepted.sample.MemberKey, ProofMemberKey: proof.Member.Key(), ProofVersion: proof.Version, ProofDigest: proof.Digest, RuntimeDigest: hex.EncodeToString(sum[:]), EvidenceKind: evidenceKind, AssociationVersion: legacyAssociationVersion,
	}
	found, err := validateLegacyAssociation(ctx, tx, association)
	if err != nil {
		return err
	}
	if !found {
		return sampleProofError(counterproof.ErrCounterBindingProofMissing)
	}
	return nil
}

func verifyLegacyDuplicateAssociation(ctx context.Context, tx *sql.Tx, key GroupKey, sample snapshot.Sample) error {
	accepted := legacyAcceptedSample{sample: sample}
	err := tx.QueryRowContext(ctx, `SELECT seq,payload_digest,created_at FROM wg_delivery_samples WHERE workspace_id=? AND group_id=? AND group_revision=? AND sample_id=? AND consumed=0`, key.WorkspaceID, key.GroupID, key.GroupRevision, sample.SampleID).Scan(&accepted.seq, &accepted.digest, &accepted.acceptedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return sampleProofError(counterproof.ErrCounterBindingProofMissing)
	}
	if err != nil {
		return err
	}
	var payload string
	err = tx.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions WHERE workspace_id=? AND group_id=? AND group_revision=?`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return sampleProofError(counterproof.ErrCounterBindingProofMissing)
	}
	if err != nil {
		return err
	}
	evidence, err := counterproof.ResolveSnapshotEvidence([]byte(payload), counterproof.SnapshotScope{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, AppliedRevision: key.GroupRevision})
	if err != nil {
		return err
	}
	proof, ok := evidence.Bindings[sample.CounterProof.Member.Key()]
	if !ok {
		return sampleProofError(counterproof.ErrCounterBindingProofMissing)
	}
	return requireLegacyAssociation(ctx, tx, key, accepted, payload, evidence, proof)
}
