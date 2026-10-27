package runtime

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
)

// sampleFrom converts an envelope into the typed sample a member contributes.
// A value that cannot be carried exactly becomes a bad sample, never a coerced
// or zeroed one. A cumulative member refuses an envelope that claims another
// measurement series than the one it was confirmed against.
func sampleFrom(member boundaryMember, envelope measurement.SampleEnvelope) (snapshot.Sample, error) {
	if member.exactScale != "" && (member.counter == nil || member.sourceRevision == "" || member.mappingRevision == "" ||
		envelope.SourceRevision != member.sourceRevision || envelope.MappingRevision != member.mappingRevision) {
		return snapshot.Sample{}, &GroupSampleError{Outcome: snapshot.OfferRejected, Reason: snapshot.ReasonRevisionMismatch}
	}
	if member.counter != nil && !counterIdentityMatches(member.counter, envelope) {
		return snapshot.Sample{}, &GroupSampleError{Outcome: snapshot.OfferRejected, Reason: snapshot.ReasonRevisionMismatch}
	}
	sample := snapshot.Sample{
		SampleID: envelope.SampleID, MemberKey: member.key, ObservedAt: envelope.ObservedAt,
		ReceivedAt: envelope.ReceivedAt.UTC(), Quality: envelope.Quality, QualityReason: envelope.QualityReason,
		SourceRevision: envelope.SourceRevision, MappingRevision: envelope.MappingRevision,
	}
	if envelope.Quality == schema.QualityGood && envelope.Value != nil {
		value, err := measurement.ExactFromGo(member.kind, envelope.Value)
		if err == nil && member.exactScale != "" {
			value, err = measurement.ScaleExactInteger(value, member.exactScale)
		}
		if err != nil {
			sample.Quality, sample.QualityReason = schema.QualityBad, reasonTypeMismatch
		} else {
			sample.Value = value
		}
	}
	return sample, nil
}
