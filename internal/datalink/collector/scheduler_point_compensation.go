package collector

import "go-gateway/internal/datalink/schema"

// CapturePointLease identifies an installation, not merely a configuration.
func (s *Scheduler) CapturePointLease(id, revision string) uint64 {
	if s == nil {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.pointInfos[id]
	if !ok || p.configRevision != revision {
		return 0
	}
	return p.observationLease
}
func (s *Scheduler) SuspendPointLease(id string, lease uint64) bool {
	if s == nil || lease == 0 {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.pointInfos[id]
	if !ok {
		return true
	}
	if p.observationLease != lease {
		return false
	}
	delete(s.pointInfos, id)
	return true
}

// RestorePointLease refuses absence: a missing slot is not ownership evidence.
func (s *Scheduler) RestorePointLease(point *schema.Point, lease uint64) bool {
	_, ok := s.RestorePointLeaseReceipt(point, lease)
	return ok
}
func (s *Scheduler) RestorePointLeaseReceipt(point *schema.Point, lease uint64) (uint64, bool) {
	if s == nil || point == nil || lease == 0 {
		return 0, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.pointInfos[point.ID]
	if !ok || p.observationLease != lease {
		return 0, false
	}
	if point.Enabled {
		if s.pointObservationLease == ^uint64(0) {
			return 0, false
		}
		return s.installPointLocked(point), true
	} else {
		delete(s.pointInfos, point.ID)
	}
	return 0, true
}

func (s *Scheduler) CaptureDevicePointLeases(deviceID string) map[string]uint64 {
	leases := map[string]uint64{}
	if s == nil {
		return leases
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for id, p := range s.pointInfos {
		if p.DeviceID == deviceID {
			leases[id] = p.observationLease
		}
	}
	return leases
}
func (s *Scheduler) InstallPointAtLease(p *schema.Point, expected uint64) (uint64, bool) {
	if s == nil || p == nil {
		return 0, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, exists := s.pointInfos[p.ID]
	if (expected == 0 && exists) || (expected != 0 && (!exists || current.observationLease != expected)) {
		return 0, false
	}
	if s.pointObservationLease == ^uint64(0) {
		return 0, false
	}
	return s.installPointLocked(p), true
}
func (s *Scheduler) RemovePointAtLease(id string, expected uint64) bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, exists := s.pointInfos[id]
	if expected == 0 {
		return !exists
	}
	if !exists || current.observationLease != expected {
		return false
	}
	delete(s.pointInfos, id)
	return true
}

// RemovePointOwned returns a bounded operation receipt. No tombstone map is kept.
// Global install epoch detects an absent-slot ABA, even when point revision repeats.
func (s *Scheduler) RemovePointOwned(id string, expected uint64) (uint64, bool) {
	if s == nil {
		return 0, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	p, exists := s.pointInfos[id]
	if expected == 0 {
		return s.pointObservationLease, !exists
	}
	if !exists || p.observationLease != expected {
		return 0, false
	}
	epoch := s.pointObservationLease
	if epoch == ^uint64(0) {
		return 0, false
	}
	delete(s.pointInfos, id)
	return epoch, true
}
func (s *Scheduler) RestoreRemovedPointOwned(p *schema.Point, removedLease, expectedEpoch uint64) (uint64, bool) {
	if s == nil || p == nil || removedLease == 0 || expectedEpoch == ^uint64(0) {
		return 0, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.pointInfos[p.ID]; exists || s.pointObservationLease != expectedEpoch {
		return 0, false
	}
	if !p.Enabled {
		return 0, false
	}
	return s.installPointLocked(p), true
}
