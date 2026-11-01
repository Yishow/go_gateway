package runtime

// RejectedObservationCount counts collected observations declined by the
// existing installation acceptance guard. It is not a source-quality or SQL
// failure count, and does not claim one particular rejection cause. Keeping
// this native counter separate leaves the existing public Stats/JSON unchanged.
func (s *Service) RejectedObservationCount() uint64 {
	return s.rejectedObservations.Load()
}
