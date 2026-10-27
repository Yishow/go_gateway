package groupdelivery

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
)

// samplePayload is the stable journal encoding of a snapshot sample. The exact
// value keeps its type and digits; an unset value (bad read) is omitted.
type samplePayload struct {
	Version         int                                   `json:"version,omitzero"`
	SampleID        string                                `json:"sample_id"`
	MemberKey       string                                `json:"member_key"`
	ObservedAt      string                                `json:"observed_at"`
	Quality         string                                `json:"quality"`
	QualityReason   string                                `json:"quality_reason,omitempty"`
	Value           *measurement.ExactValue               `json:"value,omitempty"`
	SourceRevision  string                                `json:"source_revision,omitempty"`
	MappingRevision string                                `json:"mapping_revision,omitempty"`
	CounterProof    *counterproof.CounterBindingReference `json:"counter_proof,omitempty"`
	// ReceivedAt is optional so journals written before it was kept still
	// decode, as the unknown (zero) time. It is not part of the identity digest.
	ReceivedAt string `json:"received_at,omitempty"`
}

// encodeSample returns the journal payload and its identity digest. The digest
// covers everything but the receive time, so it is unchanged for samples that
// were journaled before ReceivedAt existed and a resend of a sample never
// becomes an identity conflict only because its receive time was recorded.
func encodeSample(sample snapshot.Sample) (payload []byte, digest string, err error) {
	wire := samplePayload{
		SampleID:        sample.SampleID,
		MemberKey:       sample.MemberKey,
		ObservedAt:      sample.ObservedAt.UTC().Format(time.RFC3339Nano),
		Quality:         string(sample.Quality),
		QualityReason:   sample.QualityReason,
		SourceRevision:  sample.SourceRevision,
		MappingRevision: sample.MappingRevision,
	}
	if sample.Value.Type() != "" {
		value := sample.Value
		wire.Value = &value
	}
	if !sample.CounterProof.IsZero() {
		if err := sample.CounterProof.Validate(); err != nil {
			return nil, "", err
		}
		wire.Version, wire.CounterProof = 2, &sample.CounterProof
	}
	identity, err := json.Marshal(wire)
	if err != nil {
		return nil, "", fmt.Errorf("encode journal sample: %w", err)
	}
	if !sample.ReceivedAt.IsZero() {
		wire.ReceivedAt = sample.ReceivedAt.UTC().Format(time.RFC3339Nano)
	}
	if payload, err = json.Marshal(wire); err != nil {
		return nil, "", fmt.Errorf("encode journal sample: %w", err)
	}
	sum := sha256.Sum256(identity)
	return payload, hex.EncodeToString(sum[:]), nil
}

func decodeSample(payload []byte) (snapshot.Sample, error) {
	conflict := &counterproof.CounterBindingProofError{Reason: counterproof.ErrCounterBindingProofConflict}
	if !counterproof.UniqueJSONObject(payload, false) {
		return snapshot.Sample{}, conflict
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(payload, &fields) != nil {
		return snapshot.Sample{}, conflict
	}
	version, versioned := fields["version"]
	_, referenced := fields["counter_proof"]
	if versioned {
		var v int
		if bytes.Equal(bytes.TrimSpace(version), []byte("null")) || json.Unmarshal(version, &v) != nil {
			return snapshot.Sample{}, conflict
		}
		if v != 2 {
			return snapshot.Sample{}, &counterproof.CounterBindingProofError{Reason: counterproof.ErrCounterBindingProofUnsupported, Version: v}
		}
		if !referenced || bytes.Equal(bytes.TrimSpace(fields["counter_proof"]), []byte("null")) {
			return snapshot.Sample{}, &counterproof.CounterBindingProofError{Reason: counterproof.ErrCounterBindingProofMissing}
		}
	} else if referenced {
		return snapshot.Sample{}, conflict
	}
	if !counterproof.UniqueJSONObject(payload, true) {
		return snapshot.Sample{}, conflict
	}
	var wire samplePayload
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&wire); err != nil {
		var proofErr *counterproof.CounterBindingProofError
		if errors.As(err, &proofErr) {
			return snapshot.Sample{}, err
		}
		return snapshot.Sample{}, conflict
	}
	// encoding/json case-folds field names; legacy has neither version nor reference.
	if !versioned && (wire.Version != 0 || wire.CounterProof != nil) {
		return snapshot.Sample{}, conflict
	}
	observed, err := time.Parse(time.RFC3339Nano, wire.ObservedAt)
	if err != nil {
		return snapshot.Sample{}, conflict
	}
	sample := snapshot.Sample{
		SampleID:        wire.SampleID,
		MemberKey:       wire.MemberKey,
		ObservedAt:      observed.UTC(),
		Quality:         schema.QualityFlag(wire.Quality),
		QualityReason:   wire.QualityReason,
		SourceRevision:  wire.SourceRevision,
		MappingRevision: wire.MappingRevision,
	}
	if wire.Value != nil {
		sample.Value = *wire.Value
	}
	if wire.CounterProof != nil {
		sample.CounterProof = *wire.CounterProof
	}
	if wire.ReceivedAt != "" {
		if sample.ReceivedAt, err = time.Parse(time.RFC3339Nano, wire.ReceivedAt); err != nil {
			return snapshot.Sample{}, conflict
		}
		sample.ReceivedAt = sample.ReceivedAt.UTC()
	}
	return sample, nil
}
