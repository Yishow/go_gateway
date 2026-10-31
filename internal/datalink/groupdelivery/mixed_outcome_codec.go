package groupdelivery

import (
	"bytes"
	"encoding/json"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/snapshot"
)

// Zero counter references mean absence, never a version-zero proof. This
// mixed-only wire adapter leaves legacy Sample/journal serialization unchanged.
func encodeMixedOutcome(out snapshot.Outcome) ([]byte, error) {
	type sampleWire struct {
		*snapshot.Sample
		CounterProof *counterproof.CounterBindingReference `json:"CounterProof,omitempty"`
	}
	type memberWire struct {
		snapshot.MemberResult
		Sample *sampleWire
	}
	type outcomeWire struct {
		snapshot.Outcome
		Members []memberWire
	}
	wire := outcomeWire{Outcome: out, Members: make([]memberWire, len(out.Members))}
	for i, m := range out.Members {
		wire.Members[i].MemberResult = m
		if m.Sample != nil {
			sample := &sampleWire{Sample: m.Sample}
			if !m.Sample.CounterProof.IsZero() {
				proof := m.Sample.CounterProof
				sample.CounterProof = &proof
			}
			wire.Members[i].Sample = sample
		}
	}
	return json.Marshal(wire)
}

// DecodeMixedOutcome verifies the retained immutable history before returning
// typed member results. It is not permission to replay or recompute them.
func DecodeMixedOutcome(raw []byte, digest string) (snapshot.Outcome, error) {
	var out snapshot.Outcome
	if handoffHash(string(raw)) != digest || !counterproof.UniqueJSONObject(raw, false) {
		return out, mixedInvalid("outcome digest or JSON")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&out); err != nil {
		return out, err
	}
	if out.DueAt.IsZero() || out.RecordID == "" || len(out.Members) == 0 {
		return out, mixedInvalid("outcome identity")
	}
	return out, nil
}
