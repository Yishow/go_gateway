package dbtarget

import "go-gateway/internal/datalink/aggregation"

// The layout owns the frozen binding projection, including its optional exact
// policy values. A caller editing a candidate cannot relabel encoded evidence.
func freezeMixedCounter(source *aggregation.CounterBinding) *aggregation.CounterBinding {
	if source == nil {
		return nil
	}
	copy := *source
	if source.RolloverModulus != nil {
		v := *source.RolloverModulus
		copy.RolloverModulus = &v
	}
	if source.MaxRatePerSecond != nil {
		v := *source.MaxRatePerSecond
		copy.MaxRatePerSecond = &v
	}
	return &copy
}
func freezeMixedWide(s GroupRowSpec) GroupRowSpec {
	if s.ResultShape != MixedResultShapeV1 {
		return s
	}
	s.Members = append([]GroupRowMember(nil), s.Members...)
	s.Columns = append([]ColumnInfo(nil), s.Columns...)
	for i := range s.Members {
		s.Members[i].Counter = freezeMixedCounter(s.Members[i].Counter)
	}
	return s
}
func freezeMixedLong(s LongRowSpec) LongRowSpec {
	if s.ResultShape != MixedResultShapeV1 {
		return s
	}
	s.Members = append([]LongRowMember(nil), s.Members...)
	s.Columns = append([]ColumnInfo(nil), s.Columns...)
	for i := range s.Members {
		s.Members[i].Counter = freezeMixedCounter(s.Members[i].Counter)
	}
	return s
}
