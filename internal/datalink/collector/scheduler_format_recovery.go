package collector

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// StagePointFormatRecovery consumes the original mutation proof but quarantines
// the new slot. No scheduler lock is retained across database acceptance/commit.
func (s *Scheduler) StagePointFormatRecovery(p *schema.Point, proof PointMutationLease, profile connector.ReadProfile) (uint64, bool) {
	if s == nil || p == nil || !p.Enabled || proof.PointID != p.ID || !profile.Present() {
		return 0, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pointObservationLease == ^uint64(0) || s.pointObservationLease != proof.Epoch || s.pointInfos[p.ID].observationLease != proof.Lease {
		return 0, false
	}
	lease := s.installPointLocked(p, profile)
	if lease == 0 {
		return 0, false
	}
	slot := s.pointInfos[p.ID]
	slot.recoveryPending = true
	s.pointInfos[p.ID] = slot
	return lease, true
}

// ActivatePointFormatRecovery must only be called after known durable acceptance
// and while the caller holds the matching hydration-generation publication fence.
func (s *Scheduler) ActivatePointFormatRecovery(id string, lease uint64, profile connector.ReadProfile) bool {
	if s == nil || lease == 0 || !profile.Present() {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	slot, ok := s.pointInfos[id]
	if !ok || !slot.recoveryPending || slot.observationLease != lease || slot.readProfile != profile || slot.configRevision != profile.Resolution().PointRevision {
		return false
	}
	device := s.deviceConfigs[slot.DeviceID]
	if !profile.MatchesConnection(device.Protocol, device.Config) {
		return false
	}
	slot.recoveryPending = false
	s.pointInfos[id] = slot
	return true
}

// CaptureReadablePoint prevents direct fallback from treating a quarantined slot
// as absence. Ordinary mutation capture still sees its real lease for safe CAS.
func (s *Scheduler) CaptureReadablePoint(id string) (PointMutationLease, bool) {
	if s == nil || id == "" {
		return PointMutationLease{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	slot := s.pointInfos[id]
	if slot.recoveryPending || s.pointObservationLease == ^uint64(0) {
		return PointMutationLease{}, false
	}
	return PointMutationLease{PointID: id, Revision: slot.configRevision, Lease: slot.observationLease, Epoch: s.pointObservationLease}, true
}

// MatchesFormatRecoveryLease is a read-only replay check; it cannot promote a
// staged slot or infer ownership from the point revision alone.
func (s *Scheduler) MatchesFormatRecoveryLease(id string, lease uint64, profile connector.ReadProfile) bool {
	if s == nil || lease == 0 || !profile.Present() {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	slot, ok := s.pointInfos[id]
	if !ok || slot.recoveryPending || slot.observationLease != lease || slot.readProfile != profile || slot.configRevision != profile.Resolution().PointRevision {
		return false
	}
	device := s.deviceConfigs[slot.DeviceID]
	return profile.MatchesConnection(device.Protocol, device.Config)
}

// MatchesFormatRecoveryAbsence preserves known inactive configuration without
// manufacturing an installation lease or removing a foreign active slot.
func (s *Scheduler) MatchesFormatRecoveryAbsence(proof PointMutationLease) bool {
	if s == nil || proof.PointID == "" || proof.Lease != 0 {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.pointInfos[proof.PointID]
	return !exists && s.pointObservationLease == proof.Epoch && s.pointObservationLease != ^uint64(0)
}
