package collector

import "go-gateway/internal/datalink/connector"

// CapturePointReadProfile captures only the policy on the exact installed lease.
func (s *Scheduler) CapturePointReadProfile(id string, lease uint64) (connector.ReadProfile, bool) {
	if s == nil || lease == 0 {
		return connector.ReadProfile{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.pointInfos[id]
	if !ok || p.observationLease != lease {
		return connector.ReadProfile{}, false
	}
	return p.readProfile, true
}
