package runtime

import "math"

// Metrics returns legacy signed counters, saturating instead of wrapping negative.
func (s *Service) Metrics() Metrics {
	st := s.Snapshot()
	return Metrics{
		TotalReads:        signedMetric(st.CollectedTotal),
		TotalWrites:       signedMetric(st.WriteSuccess),
		ErrorCount:        signedMetric(st.WriteError, st.MappingError, st.PointStateError),
		MappingErrorCount: signedMetric(st.MappingError),
		WriteErrorCount:   signedMetric(st.WriteError),
	}
}

func signedMetric(counters ...uint64) int64 {
	var total uint64
	for _, counter := range counters {
		if counter > math.MaxInt64-total {
			return math.MaxInt64
		}
		total += counter
	}
	return int64(total) // #nosec G115 -- Each addition above keeps total at or below MaxInt64.
}
