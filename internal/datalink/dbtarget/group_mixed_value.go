package dbtarget

import (
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
)

// mixedMemberValue is the single selected-result adapter used by both existing
// layouts. Raw samples and cross-window deltas remain evidence, never fallbacks.
func mixedMemberValue(method snapshot.Calculation, result snapshot.MemberResult) (measurement.ExactValue, error) {
	invalid := func() (measurement.ExactValue, error) {
		return measurement.ExactValue{}, &GroupRowError{Code: rowErrTypeMismatch}
	}
	if !result.Usable() || result.Window == nil || result.Window.Calculation != method {
		return invalid()
	}
	switch method {
	case snapshot.CalculationLatest:
		if result.Sample == nil || result.Reduction != nil || result.Counter != nil {
			return invalid()
		}
		return result.Sample.Value, nil
	case snapshot.CalculationMean:
		if result.Reduction == nil || result.Counter != nil || result.Reduction.Version != measurement.MeanAggregationVersion || result.Reduction.Method != measurement.MeanAggregationMethod || result.Reduction.Count <= 0 {
			return invalid()
		}
		return result.Reduction.Value, nil
	case snapshot.CalculationDelta:
		if result.Reduction != nil || result.Counter == nil || result.Counter.Status != aggregation.CounterOK || result.Counter.Delta == nil || result.Counter.Delta.Span {
			return invalid()
		}
		return result.Counter.Delta.Value, nil
	}
	return invalid()
}
