package groupdelivery

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
)

// ErrHandoffProofConflict never exposes accepted values or descriptor contents.
var ErrHandoffProofConflict = &OwnershipError{Reason: "recording-handoff-proof-conflict"}

type handoffProofInput struct {
	Origin                      GroupKey
	Sequence                    int64
	Payload, Digest, AcceptedAt string
}

type handoffProofEvidence struct {
	OriginalReference, TargetReference       counterproof.CounterBindingReference
	OriginRuntimeDigest, TargetRuntimeDigest string
	TargetInterval                           time.Duration
}

type handoffProofMember struct {
	DeviceID        string          `json:"device_id"`
	PointID         string          `json:"point_id"`
	TagID           string          `json:"tag_id"`
	MeasurementID   *string         `json:"measurement_id"`
	SourceRevision  string          `json:"source_revision"`
	MappingRevision string          `json:"mapping_revision"`
	SourceValueType schema.DataType `json:"source_value_type"`
}
type handoffProofSnapshot struct {
	Group *struct {
		Members   []handoffProofMember `json:"members"`
		RowPolicy struct {
			IntervalSeconds int64 `json:"interval_seconds"`
		} `json:"row_policy"`
	} `json:"group"`
	Layout *struct {
		TagTypes map[string]schema.DataType `json:"tag_types"`
	} `json:"runtime_layout"`
}
type handoffResolvedProof struct {
	member   handoffProofMember
	evidence counterproof.SnapshotEvidence
	proof    counterproof.CounterBindingProof
	counter  bool
	input    schema.DataType
	interval time.Duration
	digest   string
}

// validateHandoffProof compares immutable input meaning, not the successor's
// output reduction. The caller proves journal/receipt lineage and raw bytes;
// this helper never relabels a reference or creates a legacy association.
func validateHandoffProof(ctx context.Context, q ownershipQuery, input handoffProofInput, target GroupKey) (handoffProofEvidence, error) {
	var result handoffProofEvidence
	if !input.Origin.Valid() || !target.Valid() || input.Origin.WorkspaceID != target.WorkspaceID || input.Sequence <= 0 {
		return result, ErrHandoffProofConflict
	}
	if accepted, err := time.Parse(time.RFC3339Nano, input.AcceptedAt); err != nil || accepted.IsZero() {
		return result, ErrHandoffProofConflict
	}
	sample, err := decodeVerifiedSample(input.Payload, input.Digest)
	if err != nil {
		return result, ErrHandoffProofConflict
	}
	origin, err := readHandoffProof(ctx, q, input.Origin, sample.MemberKey)
	if err != nil {
		return result, err
	}
	next, err := readHandoffProof(ctx, q, target, sample.MemberKey)
	if err != nil {
		return result, err
	}
	result.OriginRuntimeDigest, result.TargetRuntimeDigest, result.TargetInterval = origin.digest, next.digest, next.interval
	if origin.member.SourceRevision == "" || origin.member.MappingRevision == "" || next.member.SourceRevision == "" || next.member.MappingRevision == "" {
		return handoffProofEvidence{}, ErrOwnershipUnavailable
	}
	if origin.counter != next.counter || origin.member.SourceRevision != next.member.SourceRevision || origin.member.MappingRevision != next.member.MappingRevision ||
		sample.SourceRevision != origin.member.SourceRevision || sample.MappingRevision != origin.member.MappingRevision {
		return handoffProofEvidence{}, ErrHandoffProofConflict
	}
	if !origin.counter {
		kind, ok := measurement.ExactTypeForTag(origin.input)
		if !ok || origin.input != next.input || !sample.CounterProof.IsZero() ||
			!reflect.DeepEqual(origin.member.MeasurementID, next.member.MeasurementID) ||
			(sample.Value.Type() != "" && sample.Value.Type() != kind) {
			return handoffProofEvidence{}, ErrHandoffProofConflict
		}
		return result, nil
	}
	result.OriginalReference, err = origin.proof.Reference()
	if err != nil {
		return handoffProofEvidence{}, ErrHandoffProofConflict
	}
	result.TargetReference, err = next.proof.Reference()
	if err != nil {
		return handoffProofEvidence{}, ErrHandoffProofConflict
	}
	// Only the owner identity may differ. Do not reseal either immutable proof.
	a, b := origin.proof.CounterBindingFacts, next.proof.CounterBindingFacts
	b.Member.GroupID, b.Member.AppliedRevision = a.Member.GroupID, a.Member.AppliedRevision
	if !reflect.DeepEqual(a, b) || (sample.Value.Type() != "" && sample.Value.Type() != a.ResultValueType) {
		return handoffProofEvidence{}, ErrHandoffProofConflict
	}
	if !sample.CounterProof.IsZero() {
		if sample.CounterProof != result.OriginalReference {
			return handoffProofEvidence{}, ErrHandoffProofConflict
		}
	} else if err := verifyHandoffLegacyAssociation(ctx, q, input, sample, origin); err != nil {
		return handoffProofEvidence{}, err
	}
	return result, nil
}

func readHandoffProof(ctx context.Context, q ownershipQuery, key GroupKey, memberKey string) (handoffResolvedProof, error) {
	var result handoffResolvedProof
	var raw string
	if err := q.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions WHERE workspace_id=? AND group_id=? AND group_revision=?`,
		key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&raw); err != nil {
		return result, ErrOwnershipUnavailable
	}
	evidence, err := counterproof.ResolveSnapshotEvidence([]byte(raw), counterproof.SnapshotScope{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, AppliedRevision: key.GroupRevision})
	if err != nil {
		if errors.Is(err, counterproof.ErrCounterBindingProofMissing) {
			return result, ErrOwnershipUnavailable
		}
		return result, ErrHandoffProofConflict
	}
	var descriptor handoffProofSnapshot
	if json.Unmarshal([]byte(raw), &descriptor) != nil {
		return result, ErrHandoffProofConflict
	}
	if descriptor.Group == nil || descriptor.Layout == nil {
		return result, ErrOwnershipUnavailable
	}
	if !evidence.Members[memberKey] {
		return result, ErrHandoffProofConflict
	}
	seconds := descriptor.Group.RowPolicy.IntervalSeconds
	if seconds == 0 {
		return result, ErrOwnershipUnavailable
	}
	if seconds < 0 || seconds > int64((1<<63-1)/int64(time.Second)) {
		return result, ErrHandoffProofConflict
	}
	var found bool
	for _, member := range descriptor.Group.Members {
		identity := counterproof.CounterBindingMemberIdentity{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, AppliedRevision: key.GroupRevision,
			DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID}
		if identity.SampleMemberKey() != memberKey {
			continue
		}
		if found || member.DeviceID == "" || member.PointID == "" || member.TagID == "" {
			return result, ErrHandoffProofConflict
		}
		found = true
		result.member = member
		result.input = descriptor.Layout.TagTypes[member.TagID]
		if result.input == "" {
			return result, ErrOwnershipUnavailable
		}
		if _, ok := measurement.ExactTypeForTag(result.input); !ok ||
			(member.SourceValueType != "" && member.SourceValueType != result.input) {
			return result, ErrHandoffProofConflict
		}
		result.proof, result.counter = evidence.Bindings[identity.Key()]
	}
	if !found {
		return result, ErrHandoffProofConflict
	}
	sum := sha256.Sum256([]byte(raw))
	result.digest, result.interval, result.evidence = hex.EncodeToString(sum[:]), time.Duration(seconds)*time.Second, evidence
	return result, nil
}

func verifyHandoffLegacyAssociation(ctx context.Context, q ownershipQuery, input handoffProofInput, sample snapshot.Sample, origin handoffResolvedProof) error {
	var got legacyProofAssociation
	err := q.QueryRowContext(ctx, `SELECT workspace_id,group_id,group_revision,sample_id,sample_seq,sample_digest,accepted_at,member_key,proof_member_key,proof_version,proof_digest,runtime_digest,evidence_kind,association_version,associated_at FROM wg_counter_proof_associations WHERE sample_seq=?`, input.Sequence).
		Scan(&got.WorkspaceID, &got.GroupID, &got.GroupRevision, &got.SampleID, &got.SampleSeq, &got.SampleDigest, &got.AcceptedAt, &got.MemberKey,
			&got.ProofMemberKey, &got.ProofVersion, &got.ProofDigest, &got.RuntimeDigest, &got.EvidenceKind, &got.AssociationVersion, &got.AssociatedAt)
	if err != nil {
		return ErrOwnershipUnavailable
	}
	evidenceKind := runtimeCounterEvidenceKind
	if origin.evidence.LegacyExact {
		evidenceKind = legacyExactSourceEvidenceKind
	}
	want := legacyProofAssociation{WorkspaceID: input.Origin.WorkspaceID, GroupID: input.Origin.GroupID, GroupRevision: input.Origin.GroupRevision,
		SampleID: sample.SampleID, SampleSeq: input.Sequence, SampleDigest: input.Digest, AcceptedAt: input.AcceptedAt, MemberKey: sample.MemberKey,
		ProofMemberKey: origin.proof.Member.Key(), ProofVersion: origin.proof.Version, ProofDigest: origin.proof.Digest, RuntimeDigest: origin.digest,
		EvidenceKind: evidenceKind, AssociationVersion: legacyAssociationVersion, AssociatedAt: got.AssociatedAt}
	if associated, err := time.Parse(time.RFC3339Nano, got.AssociatedAt); err != nil || associated.IsZero() || got != want {
		return ErrHandoffProofConflict
	}
	return nil
}

// verifyCompactHandoffMeaning validates descriptor-backed proof authority after
// normal payload reclamation; it never reconstructs an accepted sample.
func verifyCompactHandoffMeaning(ctx context.Context, q ownershipQuery, r handoffReceipt) error {
	if !r.Origin.Valid() || !r.Target.Valid() || r.Origin.WorkspaceID != r.Target.WorkspaceID || r.OriginSeq <= 0 {
		return ErrHandoffProofConflict
	}
	if accepted, err := time.Parse(time.RFC3339Nano, r.AcceptedAt); err != nil || accepted.IsZero() {
		return ErrHandoffProofConflict
	}
	origin, err := readHandoffProof(ctx, q, r.Origin, r.Member)
	if err != nil {
		return err
	}
	target, err := readHandoffProof(ctx, q, r.Target, r.Member)
	if err != nil {
		return err
	}
	if origin.digest != r.OriginRuntimeDigest || target.digest != r.TargetRuntimeDigest {
		return ErrOwnershipUnavailable
	}
	if origin.member.SourceRevision == "" || origin.member.MappingRevision == "" || target.member.SourceRevision == "" || target.member.MappingRevision == "" {
		return ErrOwnershipUnavailable
	}
	if origin.counter != target.counter || origin.member.SourceRevision != target.member.SourceRevision || origin.member.MappingRevision != target.member.MappingRevision {
		return ErrHandoffProofConflict
	}
	if !origin.counter {
		if origin.input != target.input || !reflect.DeepEqual(origin.member.MeasurementID, target.member.MeasurementID) ||
			!r.SampleProof.IsZero() || !r.OriginalProof.IsZero() || !r.TargetProof.IsZero() {
			return ErrHandoffProofConflict
		}
		return nil
	}
	originalRef, err := origin.proof.Reference()
	if err != nil {
		return ErrHandoffProofConflict
	}
	targetRef, err := target.proof.Reference()
	if err != nil || originalRef != r.OriginalProof || targetRef != r.TargetProof {
		return ErrHandoffProofConflict
	}
	a, b := origin.proof.CounterBindingFacts, target.proof.CounterBindingFacts
	b.Member.GroupID, b.Member.AppliedRevision = a.Member.GroupID, a.Member.AppliedRevision
	if !reflect.DeepEqual(a, b) {
		return ErrHandoffProofConflict
	}
	if !r.SampleProof.IsZero() {
		if r.SampleProof != originalRef {
			return ErrHandoffProofConflict
		}
		return nil
	}
	// This metadata-only identity is solely an association query, not a sample
	// reconstruction: no value, observation or proof is offered to replay.
	identity := snapshot.Sample{SampleID: r.SampleID, MemberKey: r.Member}
	return verifyHandoffLegacyAssociation(ctx, q, handoffProofInput{Origin: r.Origin, Sequence: r.OriginSeq, Digest: r.Digest, AcceptedAt: r.AcceptedAt}, identity, origin)
}
