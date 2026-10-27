package workspace

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// CounterBindingProofFromExactSource adapts immutable exact-source facts without
// consulting a measurement repository or converting a sample value.
func CounterBindingProofFromExactSource(identity CounterBindingMemberIdentity, member WriteGroupMember, tagType schema.DataType) (CounterBindingProof, error) {
	if identity.DeviceID != member.DeviceID || identity.PointID != member.PointID || identity.TagID != member.TagID ||
		ValidateExactRecordingSource(member, tagType) != nil {
		return CounterBindingProof{}, &CounterBindingProofError{Reason: ErrCounterBindingProofConflict}
	}
	def := member.ExactRecordingSource.Definition()
	facts := CounterBindingFacts{
		Member: identity, MeasurementID: def.ID,
		DefinitionRevision: def.DefinitionRevision, SourceBindingRevision: def.SourceBindingRevision,
		SeriesEpoch: def.SeriesEpoch, SemanticKind: def.SemanticKind, Unit: def.Unit,
		SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
		ScaleMode: CounterScaleExact, ExactScale: def.ExactRecordingScale,
		SourceValueType: tagType, ResultValueType: measurement.ExactDecimal,
	}
	if def.CounterPolicy != nil {
		facts.CounterPolicy = *def.CounterPolicy
	}
	return NewCounterBindingProof(facts)
}
