package collector

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// PointMutationLease is a server-captured attempt token, not an operator revision.
// The existing global epoch also fences initially absent slots without tombstones.
type PointMutationLease struct {
	PointID      string
	Revision     string
	Lease, Epoch uint64
}

func (s *Scheduler) CaptureOrdinaryPoint(id string) (PointMutationLease, bool) {
	if s == nil || id == "" {
		return PointMutationLease{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.pointObservationLease == ^uint64(0) {
		return PointMutationLease{}, false
	}
	return PointMutationLease{PointID: id, Revision: s.pointInfos[id].configRevision, Lease: s.pointInfos[id].observationLease, Epoch: s.pointObservationLease}, true
}
func (s *Scheduler) ApplyOrdinaryPoint(p *schema.Point, proof PointMutationLease, profiles ...connector.ReadProfile) bool {
	_, ok := s.ApplyPointMutationReceipt(p, proof, profiles...)
	return ok
}

// ApplyPointMutationReceipt consumes the proof captured before provider I/O and
// returns the lease created by this atomic operation, never a later lookup.
func (s *Scheduler) ApplyPointMutationReceipt(p *schema.Point, proof PointMutationLease, profiles ...connector.ReadProfile) (uint64, bool) {
	if s == nil || p == nil || proof.PointID != p.ID {
		return 0, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pointObservationLease == ^uint64(0) || s.pointObservationLease != proof.Epoch || s.pointInfos[p.ID].observationLease != proof.Lease {
		return 0, false
	}
	if p.Enabled {
		lease := s.installPointLocked(p, profiles...)
		return lease, lease != 0
	}
	delete(s.pointInfos, p.ID)
	s.pointObservationLease++
	return 0, true
}

// A format invalidation fences pending ordinary callbacks even if no slot exists.
func (s *Scheduler) InvalidatePointFormat(id, expectedRevision string) bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pointObservationLease == ^uint64(0) {
		return false
	}
	s.pointObservationLease++
	current, exists := s.pointInfos[id]
	if !exists {
		return true
	}
	if current.configRevision != expectedRevision {
		return false
	}
	delete(s.pointInfos, id)
	return true
}

// InvalidateOrdinaryPoint closes a failed attempt, including captured absence.
// An unrelated epoch advance does not protect this point's matching old slot;
// only a different current lease is foreign and must survive.
func (s *Scheduler) InvalidateOrdinaryPoint(proof PointMutationLease) bool {
	if s == nil || proof.PointID == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, exists := s.pointInfos[proof.PointID]
	if exists && current.observationLease != proof.Lease {
		return false
	}
	if s.pointObservationLease < ^uint64(0) {
		s.pointObservationLease++
	}
	if exists {
		delete(s.pointInfos, proof.PointID)
	}
	return true
}
