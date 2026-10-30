package sourcerule

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

type pointReadProfileOwner interface {
	CapturePointReadProfile(string, uint64) (connector.ReadProfile, bool)
	RestorePointProfileReceipt(*schema.Point, uint64, connector.ReadProfile) (uint64, bool)
	RestoreRemovedPointProfile(*schema.Point, uint64, uint64, connector.ReadProfile) (uint64, bool)
}

func (s *Service) captureBaselineReadProfile(j *pointWriteJournal, p *schema.Point, lease uint64) {
	runtime, ok := s.runtimeSync.(pointReadProfileOwner)
	if !ok {
		return
	}
	profile, captured := runtime.CapturePointReadProfile(p.ID, lease)
	if !captured {
		j.runtimeUncertain = true
		return
	}
	if j.baselineProfiles == nil {
		j.baselineProfiles = map[string]connector.ReadProfile{}
	}
	j.baselineProfiles[p.ID] = profile
}
func (s *Service) restorePointProfile(j *pointWriteJournal, p *schema.Point, lease uint64) (uint64, bool) {
	if runtime, ok := s.runtimeSync.(pointReadProfileOwner); ok {
		profile, captured := j.baselineProfiles[p.ID]
		if !captured {
			return 0, false
		}
		return runtime.RestorePointProfileReceipt(p, lease, profile)
	}
	return s.runtimeSync.(pointInstallationOwner).RestorePointLeaseReceipt(p, lease)
}
func (s *Service) restoreRemovedProfile(j *pointWriteJournal, p *schema.Point, lease, epoch uint64) (uint64, bool) {
	if runtime, ok := s.runtimeSync.(pointReadProfileOwner); ok {
		profile, captured := j.baselineProfiles[p.ID]
		if !captured {
			return 0, false
		}
		return runtime.RestoreRemovedPointProfile(p, lease, epoch, profile)
	}
	return s.runtimeSync.(pointInstallationOwner).RestoreRemovedPointOwned(p, lease, epoch)
}
