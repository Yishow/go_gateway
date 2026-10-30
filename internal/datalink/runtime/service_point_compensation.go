package runtime

import "go-gateway/internal/datalink/schema"

type sourcePointInstallKey struct{}

func (s *Service) CapturePointLease(id, revision string) uint64 {
	if s == nil || s.scheduler == nil {
		return 0
	}
	return s.scheduler.CapturePointLease(id, revision)
}
func (s *Service) SuspendPointLease(id string, lease uint64) bool {
	if s == nil || s.scheduler == nil {
		return false
	}
	return s.scheduler.SuspendPointLease(id, lease)
}
func (s *Service) RestorePointLease(p *schema.Point, lease uint64) bool {
	if s == nil || s.scheduler == nil {
		return false
	}
	_, ok := s.RestorePointLeaseReceipt(p, lease)
	return ok
}

// InstallPointForMutation returns the exact installation lease created by this
// call. Disabled points have no installation proof and must not be resurrected
// by compensation based on absence alone.
func (s *Service) InstallPointForMutation(p *schema.Point) uint64 {
	if s == nil || p == nil {
		return 0
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	s.setPointMetaLocked(p)
	if s.scheduler == nil {
		return 0
	}
	if !p.Enabled {
		s.scheduler.RemovePoint(p.ID)
		return 0
	}
	return s.scheduler.AddPointWithObservationLease(p)
}

func (s *Service) InstallPointForMutationAtLease(p *schema.Point, expected uint64) (uint64, bool) {
	if s == nil || s.scheduler == nil || p == nil || !p.Enabled {
		return 0, false
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	lease, ok := s.scheduler.InstallPointAtLease(p, expected)
	if ok {
		s.setPointMetaLocked(p)
	}
	return lease, ok
}
func (s *Service) RemovePointOwned(id string, expected uint64) (uint64, bool) {
	if s == nil || s.scheduler == nil {
		return 0, false
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	epoch, ok := s.scheduler.RemovePointOwned(id, expected)
	if ok {
		delete(s.pointMetaIndex, id)
	}
	return epoch, ok
}
func (s *Service) RestoreRemovedPointOwned(p *schema.Point, removedLease, epoch uint64) (uint64, bool) {
	if s == nil || s.scheduler == nil {
		return 0, false
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	lease, ok := s.scheduler.RestoreRemovedPointOwned(p, removedLease, epoch)
	if ok {
		s.setPointMetaLocked(p)
	}
	return lease, ok
}

func (s *Service) RestorePointLeaseReceipt(p *schema.Point, lease uint64) (uint64, bool) {
	if s == nil || s.scheduler == nil {
		return 0, false
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	installed, ok := s.scheduler.RestorePointLeaseReceipt(p, lease)
	if ok {
		s.setPointMetaLocked(p)
	}
	return installed, ok
}

// Callers hold pointMetaMu before entering scheduler methods. Scheduler methods
// never invoke runtime callbacks while holding their lock.
func (s *Service) setPointMetaLocked(p *schema.Point) {
	if s.pointMetaIndex == nil {
		s.pointMetaIndex = map[string]pointMeta{}
	}
	s.pointMetaIndex[p.ID] = pointMeta{DeviceID: p.DeviceID, Address: p.Address}
}
