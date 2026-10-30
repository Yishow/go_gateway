package runtime

import (
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
)

func (s *Service) CaptureOrdinaryPoint(id string) (collector.PointMutationLease, bool) {
	if s == nil || s.scheduler == nil {
		return collector.PointMutationLease{}, false
	}
	return s.scheduler.CaptureOrdinaryPoint(id)
}
func (s *Service) ApplyOrdinaryPoint(p *schema.Point, proof collector.PointMutationLease) bool {
	if s == nil || s.scheduler == nil || p == nil {
		return false
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	if !s.scheduler.ApplyOrdinaryPoint(p, proof) {
		return false
	}
	s.setPointMetaLocked(p)
	return true
}

func (s *Service) InvalidateOrdinaryPoint(proof collector.PointMutationLease) bool {
	return s != nil && s.scheduler != nil && s.scheduler.InvalidateOrdinaryPoint(proof)
}
