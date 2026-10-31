package snapshot

import (
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// mixedInputMatches preserves the raw frozen source type while recognizing the
// exact decimal input already produced by the confirmed counter scale. Only
// that member's immutable binding authorizes this transformation.
// MixedInputMatches is shared by selection and encoding so the frozen raw tag
// type and the scaled counter input cannot silently acquire different meanings.
func MixedInputMatches(sourceType schema.DataType, method Calculation, binding *aggregation.CounterBinding, value measurement.ExactValue) bool {
	if method == CalculationDelta && binding != nil && binding.ExactRecordingScale != "" {
		source, known := measurement.ExactTypeForTag(sourceType)
		return known && (source == measurement.ExactInt64 || source == measurement.ExactUint64) && value.Type() == measurement.ExactDecimal
	}
	return meanValueMatchesSource(sourceType, value)
}
func (a *Assembler) mixedInputMatches(member Member, value measurement.ExactValue) bool {
	binding, ok := a.config.Counters[member.Key]
	if !ok {
		return MixedInputMatches(member.SourceValueType, member.Calculation, nil, value)
	}
	return MixedInputMatches(member.SourceValueType, member.Calculation, &binding, value)
}
