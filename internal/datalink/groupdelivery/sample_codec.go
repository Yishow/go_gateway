package groupdelivery

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
)

// samplePayload is the stable journal encoding of a snapshot sample. The exact
// value keeps its type and digits; an unset value (bad read) is omitted.
type samplePayload struct {
	SampleID        string                  `json:"sample_id"`
	MemberKey       string                  `json:"member_key"`
	ObservedAt      string                  `json:"observed_at"`
	Quality         string                  `json:"quality"`
	QualityReason   string                  `json:"quality_reason,omitempty"`
	Value           *measurement.ExactValue `json:"value,omitempty"`
	SourceRevision  string                  `json:"source_revision,omitempty"`
	MappingRevision string                  `json:"mapping_revision,omitempty"`
}

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
	payload, err = json.Marshal(wire)
	if err != nil {
		return nil, "", fmt.Errorf("encode journal sample: %w", err)
	}
	sum := sha256.Sum256(payload)
	return payload, hex.EncodeToString(sum[:]), nil
}

func decodeSample(payload []byte) (snapshot.Sample, error) {
	var wire samplePayload
	if err := json.Unmarshal(payload, &wire); err != nil {
		return snapshot.Sample{}, fmt.Errorf("decode journal sample: %w", err)
	}
	observed, err := time.Parse(time.RFC3339Nano, wire.ObservedAt)
	if err != nil {
		return snapshot.Sample{}, fmt.Errorf("decode journal sample time: %w", err)
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
	return sample, nil
}
