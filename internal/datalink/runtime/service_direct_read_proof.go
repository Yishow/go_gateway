package runtime

import "go-gateway/internal/datalink/collector"

func (s *Service) CaptureDirectReadProof(id string) (collector.PointMutationLease, string, bool) {
	if s == nil || s.scheduler == nil {
		return collector.PointMutationLease{}, "", false
	}
	proof, ok := s.scheduler.CaptureReadablePoint(id)
	return proof, s.scheduler.InstallationScope(), ok
}
func (s *Service) AcceptsDirectReadProof(proof collector.PointMutationLease, scope string) bool {
	return s != nil && s.scheduler != nil && s.scheduler.AcceptsDirectReadProof(proof, scope)
}
