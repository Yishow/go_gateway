package runtime

import (
	"strings"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
)

// CounterSourceFromDefinition maps a saved measurement definition to the
// confirmed counter semantics a cumulative member needs, or reports false when
// the definition does not confirm them. Nothing is inferred from a name or a
// unit: only the saved semantic kind counts.
//
// A signed_counter is confirmed only when its counter policy explicitly allows
// negative differences, and a monotonic counter whose policy allows them is
// contradictory and not confirmed. Rollover evidence comes only from the exact
// decimal-text policy (rollover_modulus, max_rate_per_second), parsed without a
// float64; the legacy float MaxRollOver is never mapped, and saved evidence that
// does not parse, or contradicts the definition, fails closed as unconfirmed. A
// decrease without evidence stays unknown.
func CounterSourceFromDefinition(def measurement.MeasurementDefinition) (CounterSource, bool) {
	if measurement.ConfirmCounter(def) != measurement.CounterConfirmed {
		return CounterSource{}, false
	}
	kind := aggregation.CounterKindMonotonic
	if def.SemanticKind == measurement.SemanticKindSignedCounter {
		kind = aggregation.CounterKindSigned
	}
	modulus, maxRate, err := def.CounterPolicy.ExactRollover()
	if err != nil {
		return CounterSource{}, false
	}
	return CounterSource{
		RolloverModulus: modulus, MaxRatePerSecond: maxRate,
		MeasurementID: def.ID, DefinitionRevision: def.DefinitionRevision, SourceBindingRevision: def.SourceBindingRevision,
		SeriesEpoch: def.SeriesEpoch, Unit: strings.TrimSpace(def.Unit), Kind: kind,
	}, true
}
