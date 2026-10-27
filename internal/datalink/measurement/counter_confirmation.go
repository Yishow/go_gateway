package measurement

import "strings"

// CounterConfirmation says whether one saved definition confirms the counter
// semantics a cumulative profile counts from. The empty value is a confirmed
// counter; every other value names the single reason it is not one, so a caller
// can block safely and tell the operator what to repair.
type CounterConfirmation string

const (
	// CounterConfirmed marks a definition whose kind, policy, unit, revisions
	// and series epoch are all saved and consistent.
	CounterConfirmed CounterConfirmation = ""
	// CounterNotCounter marks a definition of another kind, such as a gauge. It
	// is never read as a counter.
	CounterNotCounter CounterConfirmation = "not_counter"
	// CounterUnconfirmed marks a counter or signed counter that still lacks
	// confirmed evidence: a unit, a revision, a series epoch, or a policy that
	// agrees with its kind and numeric encoding.
	CounterUnconfirmed CounterConfirmation = "unconfirmed"
)

// ConfirmCounter checks the saved definition itself. It never infers counter
// semantics from a name, a unit or a register: only the saved kind, policy and
// evidence count. Whether the definition belongs to a given source is a
// separate question the caller answers from the source identity.
func ConfirmCounter(def MeasurementDefinition) CounterConfirmation {
	allowNegative := def.CounterPolicy != nil && def.CounterPolicy.AllowNegative
	switch def.SemanticKind {
	case SemanticKindCounter:
		if allowNegative {
			return CounterUnconfirmed
		}
	case SemanticKindSignedCounter:
		if !allowNegative {
			return CounterUnconfirmed
		}
	default:
		return CounterNotCounter
	}
	if def.CounterPolicy.validateExactRollover(def.SemanticKind, def.NumericEncoding) != nil {
		return CounterUnconfirmed
	}
	if def.validateExactRecordingScale() != nil {
		return CounterUnconfirmed
	}
	for _, part := range []string{def.ID, def.DefinitionRevision, def.SourceBindingRevision, def.SeriesEpoch, def.Unit} {
		if strings.TrimSpace(part) == "" {
			return CounterUnconfirmed
		}
	}
	return CounterConfirmed
}
