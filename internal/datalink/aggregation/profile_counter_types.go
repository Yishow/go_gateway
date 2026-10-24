package aggregation

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/measurement"
)

// The profile counter path is the exact, typed counterpart of
// CalculateCounterUsage. It never converts through float64: integer and decimal
// readings subtract exactly, and a float source keeps its float type and basis.

// CounterKind is the cumulative semantics a confirmed measurement declares. It
// is never guessed from a tag name or a unit.
type CounterKind string

const (
	CounterKindMonotonic CounterKind = "counter"
	CounterKindSigned    CounterKind = "signed_counter"
)

// Valid reports whether the kind is one a confirmed measurement can declare.
func (k CounterKind) Valid() bool { return k == CounterKindMonotonic || k == CounterKindSigned }

// CounterStatus names the outcome of one closed window for one counter.
type CounterStatus string

const (
	// CounterBaseline is a first usable reading: raw is kept, delta is NULL.
	CounterBaseline CounterStatus = "baseline"
	// CounterOK is a trustworthy difference over the immediately preceding window.
	CounterOK CounterStatus = "ok"
	// CounterMissing is a window without an eligible reading; nothing advances.
	CounterMissing CounterStatus = "missing"
	// CounterUncertain is a window whose readings cannot prove a difference.
	CounterUncertain CounterStatus = "reset_or_uncertain"
	// CounterGap is a window that follows missing windows. A trustworthy
	// difference across them is kept as a span total over the real observation
	// interval, never as the usage of one window; an unproven one stays NULL.
	CounterGap CounterStatus = "gap"
	// CounterInvalid is a window whose readings are not numeric.
	CounterInvalid CounterStatus = "invalid"
)

// Reasons attached to a result that carries no delta.
const (
	ReasonFirstReading        = "first-reading"
	ReasonBindingChanged      = "binding-changed"
	ReasonNoEligibleReading   = "no-eligible-reading"
	ReasonDecreaseUnproven    = "decrease-unproven"
	ReasonContradictory       = "contradictory-same-time"
	ReasonNonMonotonicTime    = "non-monotonic-time"
	ReasonMixedBasis          = "mixed-numeric-basis"
	ReasonRolloverEvidence    = "rollover-evidence-insufficient"
	ReasonRolloverRange       = "rollover-range-invalid"
	ReasonRolloverMultiWrap   = "rollover-multiple-wraps-possible"
	ReasonRolloverRate        = "rollover-rate-bound-violated"
	ReasonNonNumericReading   = "non-numeric-reading"
	ReasonDeltaNotRepresented = "delta-not-representable"
)

// Errors that mean the caller broke the contract, as opposed to a window whose
// data is merely not comparable.
var (
	ErrCounterBindingIncomplete = errors.New("counter binding incomplete")
	ErrCounterPolicyInvalid     = errors.New("counter policy invalid")
	ErrCounterWindowInvalid     = errors.New("counter window invalid")
	ErrCounterWindowOrder       = errors.New("counter window already folded or out of order")
	ErrCounterReadingInvalid    = errors.New("counter reading invalid")
	ErrCounterCheckpointInvalid = errors.New("counter checkpoint invalid")
	// ErrCounterCheckpointVersion is the specific cause of a checkpoint written by
	// another version of the format. It always accompanies ErrCounterCheckpointInvalid.
	ErrCounterCheckpointVersion = errors.New("counter checkpoint version unsupported")
)

// CounterBinding is the identity under which readings are comparable: the
// workspace, group and profile revision, the source member, and the confirmed
// measurement definition, source binding and series epoch. Any change starts a
// new baseline; readings of different bindings are never subtracted.
type CounterBinding struct {
	WorkspaceID           string      `json:"workspace_id"`
	GroupID               string      `json:"group_id"`
	ProfileRevision       string      `json:"profile_revision"`
	MemberKey             string      `json:"member_key"`
	MeasurementID         string      `json:"measurement_id"`
	DefinitionRevision    string      `json:"definition_revision"`
	SourceBindingRevision string      `json:"source_binding_revision"`
	SeriesEpoch           string      `json:"series_epoch"`
	Kind                  CounterKind `json:"kind"`
	Unit                  string      `json:"unit"`
	// RolloverModulus and MaxRatePerSecond are the confirmed rollover evidence of
	// a monotonic counter. They are exact values, never floats, and part of the
	// identity: changing either starts a new baseline. A modulus alone is allowed
	// and proves nothing; a decrease then stays uncertain.
	RolloverModulus  *measurement.ExactValue `json:"rollover_modulus,omitempty"`
	MaxRatePerSecond *measurement.ExactValue `json:"max_rate_per_second,omitempty"`
}

const counterBaselineKeyVersion = "profile-counter-baseline-v2"

// Validate reports the first missing identity part.
func (b CounterBinding) Validate() error {
	for _, field := range []struct{ name, value string }{
		{"workspace", b.WorkspaceID}, {"group", b.GroupID}, {"profile revision", b.ProfileRevision},
		{"member", b.MemberKey}, {"measurement", b.MeasurementID},
		{"definition revision", b.DefinitionRevision}, {"source binding revision", b.SourceBindingRevision},
		{"series epoch", b.SeriesEpoch}, {"unit", b.Unit},
	} {
		if strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("%w: %s is required", ErrCounterBindingIncomplete, field.name)
		}
	}
	if !b.Kind.Valid() {
		return fmt.Errorf("%w: counter semantics must be confirmed as counter or signed_counter", ErrCounterBindingIncomplete)
	}
	return b.validatePolicy()
}

// validatePolicy refuses rollover evidence that cannot be used exactly: it must
// be a positive integer or decimal, belong to a monotonic counter, and a rate
// bound has no meaning without the modulus it protects.
func (b CounterBinding) validatePolicy() error {
	invalid := func(reason string) error { return fmt.Errorf("%w: %s", ErrCounterPolicyInvalid, reason) }
	positive := func(value *measurement.ExactValue) bool {
		if value == nil {
			return true
		}
		rat, exact := value.Rat()
		return exact && rat.Sign() > 0
	}
	switch {
	case b.RolloverModulus == nil && b.MaxRatePerSecond == nil:
		return nil
	case b.Kind != CounterKindMonotonic:
		return invalid("rollover applies to a monotonic counter only")
	case b.RolloverModulus == nil:
		return invalid("a maximum rate needs the rollover modulus it protects")
	case !positive(b.RolloverModulus) || !positive(b.MaxRatePerSecond):
		return invalid("rollover modulus and maximum rate must be exact and positive")
	}
	return nil
}

// policyIdentity is the canonical text of the rollover evidence for the key.
func (b CounterBinding) policyIdentity() (modulus, maxRate string) {
	canonical := func(value *measurement.ExactValue) string {
		if value == nil {
			return "none"
		}
		digits, _ := value.CanonicalDigits()
		return digits
	}
	return canonical(b.RolloverModulus), canonical(b.MaxRatePerSecond)
}

// Key is the stable baseline key of the binding.
func (b CounterBinding) Key() (string, error) {
	if err := b.Validate(); err != nil {
		return "", err
	}
	modulus, maxRate := b.policyIdentity()
	return measurement.HashStringTuple(
		counterBaselineKeyVersion, b.WorkspaceID, b.GroupID, b.ProfileRevision, b.MemberKey, b.MeasurementID,
		b.DefinitionRevision, b.SourceBindingRevision, b.SeriesEpoch, string(b.Kind), b.Unit, modulus, maxRate,
	)
}

// CounterReading is one accepted, eligible source reading. Value keeps its
// exact type and digits.
type CounterReading struct {
	SampleID   string                 `json:"sample_id"`
	ObservedAt time.Time              `json:"observed_at"`
	ReceivedAt time.Time              `json:"received_at,omitzero"`
	Value      measurement.ExactValue `json:"value"`
}

// CounterWindow is the half-open UTC window [Start, End) being closed.
type CounterWindow struct {
	Start, End time.Time
}

// DeltaBasis says what the precision of a delta rests on.
type DeltaBasis string

const (
	// DeltaBasisExact is an exact difference of integer or decimal readings.
	DeltaBasisExact DeltaBasis = "exact"
	// DeltaBasisFloat64 is a difference of float readings. It carries no more
	// precision than the source offered and is not an exact measurement.
	DeltaBasisFloat64 DeltaBasis = "float64"
)

// CounterDelta is a trustworthy difference between two actual readings.
type CounterDelta struct {
	Value measurement.ExactValue
	Basis DeltaBasis
	// From and To are the observation times of the baseline and end readings,
	// not the window bounds.
	From, To time.Time
	// Span marks a span total: the difference across one or more missing windows
	// over its real interval, to be read as a total and never as one window.
	Span             bool
	BaselineSampleID string
	EndSampleID      string
}

// CounterWindowResult is the typed record of one closed window. Delta is nil
// when no trustworthy difference exists; that is NULL, never zero.
type CounterWindowResult struct {
	Status CounterStatus
	Reason string
	Raw    *CounterReading
	Delta  *CounterDelta
	// Events are the discontinuities and compensations found in the window's
	// pre-decimation sequence, bounded; EventsOmitted counts the ones dropped.
	Events        []CounterEvent
	EventsOmitted int
}
