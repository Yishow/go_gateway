package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/snapshot"
)

func sampleProofError(reason error) error {
	return &counterproof.CounterBindingProofError{Reason: reason}
}

func verifyExistingSample(ctx context.Context, tx *sql.Tx, key GroupKey, sample snapshot.Sample, digest string) (snapshot.Sample, bool, error) {
	var existing, payload, workspaceID string
	err := tx.QueryRowContext(ctx,
		`SELECT payload_digest, payload, workspace_id FROM wg_delivery_samples WHERE group_id = ? AND group_revision = ? AND acceptance_kind='native' AND COALESCE(NULLIF(origin_sample_id,''),sample_id) = ?`,
		key.GroupID, key.GroupRevision, sample.SampleID).Scan(&existing, &payload, &workspaceID)
	if errors.Is(err, sql.ErrNoRows) {
		return snapshot.Sample{}, false, nil
	}
	if err != nil {
		return snapshot.Sample{}, false, fmt.Errorf("check journal sample: %w", err)
	}
	if workspaceID != key.WorkspaceID {
		return snapshot.Sample{}, true, sampleProofError(counterproof.ErrCounterBindingProofConflict)
	}
	previous, err := decodeVerifiedSample(payload, existing)
	if err != nil {
		return snapshot.Sample{}, true, err
	}
	if existing != digest && !sameLegacySample(previous, sample, existing) {
		return snapshot.Sample{}, true, ErrSampleConflict
	}
	return previous, true, verifySampleProof(ctx, tx, key, sample, false)
}

// VerifyCounterDuplicate checks an already accepted sample without recreating
// a missing journal entry or changing its original bytes.
func (s *Store) VerifyCounterDuplicate(ctx context.Context, key GroupKey, sample snapshot.Sample) error {
	if !key.Valid() {
		return ErrInvalidGroupKey
	}
	if sample.CounterProof.IsZero() {
		return sampleProofError(counterproof.ErrCounterBindingProofMissing)
	}
	_, digest, err := encodeSample(sample)
	if err != nil {
		return err
	}
	return s.inTx(ctx, "counter duplicate", func(tx *sql.Tx) error {
		if err := refuseOrdinaryHandoffAppend(ctx, tx, key, sample); err != nil {
			return err
		}
		if err := validateCutoverAcceptance(ctx, tx, key, sample); err != nil {
			return err
		}
		previous, found, err := verifyExistingSample(ctx, tx, key, sample, digest)
		if err != nil {
			return err
		}
		if !found {
			return sampleProofError(counterproof.ErrCounterBindingProofMissing)
		}
		if previous.CounterProof.IsZero() {
			return verifyLegacyDuplicateAssociation(ctx, tx, key, sample)
		}
		return nil
	})
}

type proofQuery interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// verifySampleProof reads the immutable stored descriptor, not a caller's copy.
// Append passes its write transaction so proof verification and journal commit
// cannot be separated by a descriptor write. Legacy restore stays unassociated.
func verifySampleProof(ctx context.Context, runner proofQuery, key GroupKey, sample snapshot.Sample, requireReference bool) error {
	ref := sample.CounterProof
	if !ref.IsZero() {
		if err := ref.Validate(); err != nil {
			return err
		}
		if ref.Member.WorkspaceID != key.WorkspaceID || ref.Member.GroupID != key.GroupID ||
			ref.Member.AppliedRevision != key.GroupRevision || ref.Member.SampleMemberKey() != sample.MemberKey {
			return sampleProofError(counterproof.ErrCounterBindingProofConflict)
		}
	} else if !requireReference {
		return nil
	}
	var payload string
	err := runner.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions WHERE workspace_id=? AND group_id=? AND group_revision=?`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		if ref.IsZero() {
			return nil
		} // legacy/non-counter compatibility; production counter activation supplies a reference
		return sampleProofError(counterproof.ErrCounterBindingProofMissing)
	}
	if err != nil {
		return err
	}
	proofs, err := counterproof.ResolveSnapshot([]byte(payload), counterproof.SnapshotScope{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, AppliedRevision: key.GroupRevision})
	if err != nil {
		return err
	}
	if ref.IsZero() {
		for _, proof := range proofs {
			if proof.Member.SampleMemberKey() == sample.MemberKey {
				return sampleProofError(counterproof.ErrCounterBindingProofMissing)
			}
		}
		return nil
	}
	proof, ok := proofs[ref.Member.Key()]
	if !ok {
		return sampleProofError(counterproof.ErrCounterBindingProofMissing)
	}
	if proof.Version != ref.Version || proof.Digest != ref.Digest || proof.SourceRevision != sample.SourceRevision || proof.MappingRevision != sample.MappingRevision ||
		(sample.Value.Type() != "" && proof.ResultValueType != sample.Value.Type()) {
		return sampleProofError(counterproof.ErrCounterBindingProofConflict)
	}
	return nil
}

func decodeVerifiedSample(payload, digest string) (snapshot.Sample, error) {
	sample, err := decodeSample([]byte(payload))
	if err != nil {
		return snapshot.Sample{}, err
	}
	_, actual, err := encodeSample(sample)
	if err != nil {
		return snapshot.Sample{}, err
	}
	if actual != digest {
		return snapshot.Sample{}, sampleProofError(counterproof.ErrCounterBindingProofConflict)
	}
	return sample, nil
}

func sameLegacySample(previous, incoming snapshot.Sample, digest string) bool {
	if !previous.CounterProof.IsZero() || incoming.CounterProof.IsZero() {
		return false
	}
	incoming.CounterProof = counterproof.CounterBindingReference{}
	_, legacyDigest, err := encodeSample(incoming)
	return err == nil && digest == legacyDigest
}
