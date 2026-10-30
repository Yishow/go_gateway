package snapshot

import "go-gateway/internal/datalink/aggregation"

// cloneCounterBinding takes ownership of the two optional policy values. Exact
// values hold immutable scalar/string payloads; copying their structs preserves
// exact digits and the existing binding key, without a JSON round trip.
func cloneCounterBinding(binding aggregation.CounterBinding) aggregation.CounterBinding {
	if binding.RolloverModulus != nil {
		value := *binding.RolloverModulus
		binding.RolloverModulus = &value
	}
	if binding.MaxRatePerSecond != nil {
		value := *binding.MaxRatePerSecond
		binding.MaxRatePerSecond = &value
	}
	return binding
}

func cloneCounterCheckpoint(checkpoint aggregation.CounterCheckpoint) aggregation.CounterCheckpoint {
	checkpoint.Binding = cloneCounterBinding(checkpoint.Binding)
	return checkpoint
}
