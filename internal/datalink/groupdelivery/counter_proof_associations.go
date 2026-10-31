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

const (
	legacyAssociationVersion      = 1
	runtimeCounterEvidenceKind    = "runtime-counter-bindings-v1"
	legacyExactSourceEvidenceKind = "legacy-exact-source-v1"
)

type legacyAcceptedSample struct {
	seq                                                                       int64
	workspaceID, sampleID, memberKey, observedAt, payload, digest, acceptedAt string
	sample                                                                    snapshot.Sample
}

// AssociateLegacyCounterProofs atomically audits proven legacy accepted samples.
// It never modifies original journal payloads, identity digests or descriptors.
// Expected references must match the actual persisted descriptor, not replace it.
func (s *Store) AssociateLegacyCounterProofs(ctx context.Context, key GroupKey, expected []counterproof.CounterBindingReference) error {
	if !key.Valid() {
		return ErrInvalidGroupKey
	}
	if len(expected) == 0 {
		return nil
	}
	references := map[string]counterproof.CounterBindingReference{}
	for _, ref := range expected {
		if err := ref.Validate(); err != nil {
			return err
		}
		if ref.Member.WorkspaceID != key.WorkspaceID || ref.Member.GroupID != key.GroupID || ref.Member.AppliedRevision != key.GroupRevision {
			return sampleProofError(counterproof.ErrCounterBindingProofConflict)
		}
		member := ref.Member.SampleMemberKey()
		if _, exists := references[member]; exists {
			return sampleProofError(counterproof.ErrCounterBindingProofConflict)
		}
		references[member] = ref
	}
	return s.inTx(ctx, "legacy counter proof", func(ctx context.Context, tx *sql.Tx) error {
		samples, err := readLegacyAcceptedSet(ctx, tx, key)
		if err != nil {
			return err
		}
		needsAssociation := false
		for _, accepted := range samples {
			if _, counter := references[accepted.memberKey]; counter && accepted.sample.CounterProof.IsZero() {
				needsAssociation = true
			}
		}
		if !needsAssociation {
			return nil
		}
		var raw string
		err = tx.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions WHERE workspace_id=? AND group_id=? AND group_revision=?`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&raw)
		if errors.Is(err, sql.ErrNoRows) {
			return sampleProofError(counterproof.ErrCounterBindingProofMissing)
		}
		if err != nil {
			return err
		}
		evidence, err := counterproof.ResolveSnapshotEvidence([]byte(raw), counterproof.SnapshotScope{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, AppliedRevision: key.GroupRevision})
		if err != nil {
			return err
		}
		if len(evidence.Bindings) != len(references) {
			return sampleProofError(counterproof.ErrCounterBindingProofConflict)
		}
		proofs := map[string]counterproof.CounterBindingProof{}
		for _, proof := range evidence.Bindings {
			ref, err := proof.Reference()
			if err != nil {
				return err
			}
			member := ref.Member.SampleMemberKey()
			if references[member] != ref {
				return sampleProofError(counterproof.ErrCounterBindingProofConflict)
			}
			proofs[member] = proof
		}
		sum := sha256.Sum256([]byte(raw))
		runtimeDigest := hex.EncodeToString(sum[:])
		evidenceKind := runtimeCounterEvidenceKind
		if evidence.LegacyExact {
			evidenceKind = legacyExactSourceEvidenceKind
		}
		var associations []legacyProofAssociation
		// Validate the entire applicable accepted set before publishing any audit.
		for _, accepted := range samples {
			sample := accepted.sample
			if !evidence.Members[sample.MemberKey] {
				return sampleProofError(counterproof.ErrCounterBindingProofConflict)
			}
			proof, counter := proofs[sample.MemberKey]
			if !counter {
				if !sample.CounterProof.IsZero() {
					return sampleProofError(counterproof.ErrCounterBindingProofConflict)
				}
				continue
			}
			if proof.SourceRevision != sample.SourceRevision || proof.MappingRevision != sample.MappingRevision || (sample.Value.Type() != "" && proof.ResultValueType != sample.Value.Type()) {
				return sampleProofError(counterproof.ErrCounterBindingProofConflict)
			}
			if !sample.CounterProof.IsZero() {
				if sample.CounterProof != references[sample.MemberKey] {
					return sampleProofError(counterproof.ErrCounterBindingProofConflict)
				}
				continue
			}
			associations = append(associations, legacyProofAssociation{
				WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, GroupRevision: key.GroupRevision, SampleID: sample.SampleID, SampleSeq: accepted.seq, SampleDigest: accepted.digest, AcceptedAt: accepted.acceptedAt,
				MemberKey: sample.MemberKey, ProofMemberKey: proof.Member.Key(), ProofVersion: proof.Version, ProofDigest: proof.Digest, RuntimeDigest: runtimeDigest, EvidenceKind: evidenceKind, AssociationVersion: legacyAssociationVersion, AssociatedAt: timestamp(s.now()),
			})
		}
		for _, association := range associations {
			if err := insertLegacyAssociation(ctx, tx, association); err != nil {
				return err
			}
		}
		return nil
	})
}

func readLegacyAcceptedSet(ctx context.Context, tx *sql.Tx, key GroupKey) ([]legacyAcceptedSample, error) {
	return readAcceptedSet(ctx, tx, key, time.Time{})
}

func readAcceptedSet(ctx context.Context, tx *sql.Tx, key GroupKey, before time.Time) ([]legacyAcceptedSample, error) {
	if err := validateHandoffScope(ctx, tx, key.GroupID, key.GroupRevision); err != nil {
		return nil, err
	}
	// Handoff receipts validate transferred samples and preserve their original
	// proof identity; they cannot be reassociated as native inputs.
	query := `SELECT seq,workspace_id,COALESCE(NULLIF(origin_sample_id,''),sample_id),member_key,observed_at,payload,payload_digest,created_at FROM wg_delivery_samples WHERE group_id=? AND group_revision=? AND consumed=0 AND NOT EXISTS(SELECT 1 FROM wg_delivery_sample_handoffs h WHERE h.source_seq=wg_delivery_samples.seq OR h.target_seq=wg_delivery_samples.seq)`
	args := []any{key.GroupID, key.GroupRevision}
	if !before.IsZero() {
		query += ` AND bucket_start<?`
		args = append(args, timestamp(before))
	}
	rows, err := tx.QueryContext(ctx, query+` ORDER BY seq`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var samples []legacyAcceptedSample
	for rows.Next() {
		var accepted legacyAcceptedSample
		if err := rows.Scan(&accepted.seq, &accepted.workspaceID, &accepted.sampleID, &accepted.memberKey, &accepted.observedAt, &accepted.payload, &accepted.digest, &accepted.acceptedAt); err != nil {
			return nil, err
		}
		accepted.sample, err = decodeVerifiedSample(accepted.payload, accepted.digest)
		if err != nil {
			return nil, err
		}
		observed, observedErr := time.Parse(time.RFC3339Nano, accepted.observedAt)
		_, acceptedErr := time.Parse(time.RFC3339Nano, accepted.acceptedAt)
		if accepted.workspaceID != key.WorkspaceID || accepted.sampleID != accepted.sample.SampleID || accepted.memberKey != accepted.sample.MemberKey || observedErr != nil || acceptedErr != nil || !observed.Equal(accepted.sample.ObservedAt) {
			return nil, sampleProofError(counterproof.ErrCounterBindingProofConflict)
		}
		samples = append(samples, accepted)
	}
	return samples, rows.Err()
}
