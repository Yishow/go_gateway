package runtime

import (
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// InstallPreparedFormatRecovery performs no projection/database reads. The
// orchestrator resolves the profile before taking any publication locks, then
// consumes its original slot proof here. The staged lease cannot produce samples
// until explicit known-commit activation. Only this point's metadata is changed.
func (s *Service) InstallPreparedFormatRecovery(p *schema.Point, proof collector.PointMutationLease, profile connector.ReadProfile) (uint64, bool) {
	if s == nil || s.scheduler == nil || p == nil || !p.Enabled || !profile.Present() || connector.ValidatePointReadProfile(p, profile) != nil {
		return 0, false
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	lease, ok := s.scheduler.StagePointFormatRecovery(p, proof, profile)
	if !ok || lease == 0 {
		return 0, false
	}
	s.setPointMetaLocked(p)
	return lease, true
}

// ActivatePreparedFormatRecovery is the compatibility seam for an unmapped
// point only. Mapped targets require the captured SQL binding publication below.
func (s *Service) ActivatePreparedFormatRecovery(id string, lease uint64, profile connector.ReadProfile) bool {
	if s == nil || s.scheduler == nil {
		return false
	}
	s.mappingMu.Lock()
	defer s.mappingMu.Unlock()
	if len(s.mappingIndex[id]) != 0 {
		return false
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	return s.scheduler.ActivatePointFormatRecovery(id, lease, profile)
}

func (s *Service) MatchesFormatRecoveryLease(id string, lease uint64, profile connector.ReadProfile) bool {
	return s != nil && s.scheduler != nil && s.scheduler.MatchesFormatRecoveryLease(id, lease, profile)
}
func (s *Service) MatchesFormatRecoveryAbsence(proof collector.PointMutationLease) bool {
	return s != nil && s.scheduler != nil && s.scheduler.MatchesFormatRecoveryAbsence(proof)
}
