package runtime

import (
	"errors"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
)

// ErrCounterSampleClaimConflict is an incoming claim mismatch, not corruption
// of the immutable proof or a reason to quarantine an otherwise healthy group.
var ErrCounterSampleClaimConflict = errors.New("counter-sample-claim-conflict")

// CounterSampleClaimError identifies the affected member and claim field only;
// incoming values and sample data never appear in its safe error text.
type CounterSampleClaimError struct {
	Member counterproof.CounterBindingMemberIdentity
	Claim  string
}

func (e *CounterSampleClaimError) Error() string { return ErrCounterSampleClaimConflict.Error() }
func (e *CounterSampleClaimError) Unwrap() error { return ErrCounterSampleClaimConflict }

// counterClaimIdentity is a value-only projection of a validated applied proof,
// not a lookup key for mutable measurement definitions.
type counterClaimIdentity struct {
	measurementID, definitionRevision, sourceBindingRevision, seriesEpoch string
}

func validateCounterProofClaims(member boundaryMember, envelope measurement.SampleEnvelope) error {
	if member.counterProof.IsZero() {
		return nil // expansion compatibility; mandatory activation is a separate gate
	}
	for _, claim := range []struct{ field, incoming, applied string }{
		{"measurement_id", envelope.MeasurementID, member.counterClaims.measurementID},
		{"series_epoch", envelope.SeriesEpoch, member.counterClaims.seriesEpoch},
		{"definition_revision", envelope.DefinitionRevision, member.counterClaims.definitionRevision},
		{"source_binding_revision", envelope.SourceBindingRevision, member.counterClaims.sourceBindingRevision},
	} {
		if claim.applied == "" {
			cause := &counterproof.CounterBindingProofError{Reason: counterproof.ErrCounterBindingProofMissing}
			return &GroupSampleError{Outcome: snapshot.OfferRejected, Reason: cause.Error(), Cause: cause}
		}
		if claim.incoming != "" && claim.incoming != claim.applied {
			cause := &CounterSampleClaimError{Member: member.counterProof.Member, Claim: claim.field}
			return &GroupSampleError{Outcome: snapshot.OfferRejected, Reason: cause.Error(), Cause: cause}
		}
	}
	return nil
}
