package counterproof

import (
	"encoding/json"
	"errors"
	"strings"
	"unicode/utf8"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// CounterBindingMemberIdentity identifies one member of one immutable revision.
type CounterBindingMemberIdentity struct {
	WorkspaceID     string `json:"workspace_id"`
	GroupID         string `json:"group_id"`
	AppliedRevision string `json:"applied_revision"`
	DeviceID        string `json:"device_id"`
	PointID         string `json:"point_id"`
	TagID           string `json:"tag_id"`
}

// Key identifies a proof without relying on a tag label.
func (m CounterBindingMemberIdentity) Key() string {
	key, _ := json.Marshal([6]string{m.WorkspaceID, m.GroupID, m.AppliedRevision, m.DeviceID, m.PointID, m.TagID})
	return string(key)
}

// CounterScaleMode distinguishes an unscaled source from exact integer scaling.
type CounterScaleMode string

const (
	// CounterScaleUnscaled freezes the identity multiplier, explicitly "1".
	CounterScaleUnscaled CounterScaleMode = "unscaled"
	// CounterScaleExact freezes an integer-to-decimal conversion, even at "1".
	CounterScaleExact CounterScaleMode = "exact"
)

// CounterBindingFacts are accepted counter meaning, never editable references.
type CounterBindingFacts struct {
	Member                CounterBindingMemberIdentity `json:"member"`
	MeasurementID         string                       `json:"measurement_id"`
	DefinitionRevision    string                       `json:"definition_revision"`
	SourceBindingRevision string                       `json:"source_binding_revision"`
	SeriesEpoch           string                       `json:"series_epoch"`
	SemanticKind          measurement.SemanticKind     `json:"semantic_kind"`
	Unit                  string                       `json:"unit"`
	SourceRevision        string                       `json:"source_revision"`
	MappingRevision       string                       `json:"mapping_revision"`
	ScaleMode             CounterScaleMode             `json:"scale_mode"`
	ExactScale            string                       `json:"exact_scale"`
	SourceValueType       schema.DataType              `json:"source_value_type"`
	ResultValueType       measurement.ExactType        `json:"result_value_type"`
	CounterPolicy         measurement.CounterPolicy    `json:"counter_policy"`
}

// CounterBindingProof seals the v1 facts and their canonical digest.
// Callers must not edit sealed facts: serialization checks rather than renews
// the digest. Policy decimal text is preserved, including significant scale.
type CounterBindingProof struct {
	Version int `json:"version"`
	CounterBindingFacts
	Digest string `json:"digest,omitempty"`
}

// CounterBindingProofs is keyed by the complete member identity.
type CounterBindingProofs map[string]CounterBindingProof

// CounterBindingProofVersion is the only supported proof format.
const CounterBindingProofVersion = 1

var (
	// ErrCounterBindingProofConflict marks corrupt or contradictory frozen facts.
	ErrCounterBindingProofConflict = errors.New("counter-replay-proof-conflict")
	// ErrCounterBindingProofUnsupported marks a version this binary cannot read.
	ErrCounterBindingProofUnsupported = errors.New("counter-replay-proof-unsupported")
)

// CounterBindingProofError classifies an unusable proof without exposing facts.
type CounterBindingProofError struct {
	Reason  error
	Version int
}

func (e *CounterBindingProofError) Error() string { return e.Reason.Error() }
func (e *CounterBindingProofError) Unwrap() error { return e.Reason }

// NewCounterBindingProof seals independently supplied immutable facts.
func NewCounterBindingProof(facts CounterBindingFacts) (CounterBindingProof, error) {
	proof := CounterBindingProof{Version: CounterBindingProofVersion, CounterBindingFacts: facts}
	proof.CounterPolicy = *cloneExactCounterPolicy(&facts.CounterPolicy)
	if err := proof.validateFacts(); err != nil {
		return CounterBindingProof{}, err
	}
	digest, err := proof.canonicalDigest()
	if err != nil {
		return CounterBindingProof{}, err
	}
	proof.Digest = digest
	return proof, nil
}

func (p CounterBindingProof) conflict() error {
	return &CounterBindingProofError{Reason: ErrCounterBindingProofConflict, Version: p.Version}
}

func (p CounterBindingProof) validateFacts() error {
	for _, value := range []string{
		p.Member.WorkspaceID, p.Member.GroupID, p.Member.AppliedRevision,
		p.Member.DeviceID, p.Member.PointID, p.Member.TagID, p.MeasurementID,
		p.DefinitionRevision, p.SourceBindingRevision, p.SeriesEpoch, p.Unit,
		p.SourceRevision, p.MappingRevision,
	} {
		if strings.TrimSpace(value) == "" || !utf8.ValidString(value) {
			return p.conflict()
		}
	}
	exactType, ok := measurement.ExactTypeForTag(p.SourceValueType)
	if !ok || (exactType != measurement.ExactInt64 && exactType != measurement.ExactUint64 && exactType != measurement.ExactFloat64) {
		return p.conflict()
	}
	def := measurement.MeasurementDefinition{
		ID: p.MeasurementID, DefinitionRevision: p.DefinitionRevision,
		SourceBindingRevision: p.SourceBindingRevision, SeriesEpoch: p.SeriesEpoch,
		Unit: p.Unit, SemanticKind: p.SemanticKind, NumericEncoding: string(p.SourceValueType),
		CounterPolicy: &p.CounterPolicy,
	}
	switch p.ScaleMode {
	case CounterScaleUnscaled:
		if p.ExactScale != "1" || p.ResultValueType != exactType {
			return p.conflict()
		}
	case CounterScaleExact:
		if p.ExactScale == "" || p.ResultValueType != measurement.ExactDecimal {
			return p.conflict()
		}
		def.ExactRecordingScale = p.ExactScale
	default:
		return p.conflict()
	}
	if measurement.ConfirmCounter(def) != measurement.CounterConfirmed {
		return p.conflict()
	}
	return nil
}

func cloneExactCounterPolicy(policy *measurement.CounterPolicy) *measurement.CounterPolicy {
	cloned := *policy
	if policy.ResetThreshold != nil {
		value := *policy.ResetThreshold
		cloned.ResetThreshold = &value
	}
	if policy.MaxRollOver != nil {
		value := *policy.MaxRollOver
		cloned.MaxRollOver = &value
	}
	return &cloned
}
