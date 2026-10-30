package collector

// RemovePointRevision stops only the matching configuration installation.
// Absence is already stopped. It does not provide the later acquisition-proof
// fingerprint contract or reinterpret any previously accepted sample.
func (s *Scheduler) RemovePointRevision(pointID, expectedRevision string) bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, exists := s.pointInfos[pointID]
	if !exists {
		return true
	}
	if current.configRevision != expectedRevision {
		return false
	}
	delete(s.pointInfos, pointID)
	return true
}
