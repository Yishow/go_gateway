package runtime

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

func (s *Service) CapturePointReadProfile(id string, lease uint64) (connector.ReadProfile, bool) {
	if s == nil || s.scheduler == nil {
		return connector.ReadProfile{}, false
	}
	return s.scheduler.CapturePointReadProfile(id, lease)
}
func (s *Service) RestorePointProfileReceipt(p *schema.Point, lease uint64, profile connector.ReadProfile) (uint64, bool) {
	if s == nil || s.scheduler == nil || connector.ValidatePointReadProfile(p, profile) != nil {
		return 0, false
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	installed, ok := s.scheduler.RestorePointLeaseReceipt(p, lease, profile)
	if ok {
		s.setPointMetaLocked(p)
	}
	return installed, ok
}
func (s *Service) RestoreRemovedPointProfile(p *schema.Point, removedLease, epoch uint64, profile connector.ReadProfile) (uint64, bool) {
	if s == nil || s.scheduler == nil || connector.ValidatePointReadProfile(p, profile) != nil {
		return 0, false
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	installed, ok := s.scheduler.RestoreRemovedPointOwned(p, removedLease, epoch, profile)
	if ok {
		s.setPointMetaLocked(p)
	}
	return installed, ok
}
