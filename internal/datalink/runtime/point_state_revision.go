package runtime

import (
	"context"
	"go-gateway/internal/datalink/collector"
)

// Configuration CAS and runtime installation acceptance are separate checks.
// A raw value accepted before a later configuration change cannot overwrite
// that later point row. Already frozen recording envelopes are not revisited.
func (s *Service) persistCapturedPointState(ctx context.Context, cv collector.CollectedValue) bool {
	if !s.config.UpdatePointState || s.pointSvc == nil {
		return true
	}
	if cv.InstallationLease != 0 || s.scheduler != nil {
		if err := s.pointSvc.UpdateReadResultAtRevision(ctx, cv.PointID, cv.DeviceID, cv.PointRevision, cv.Value, cv.Error); err != nil {
			s.pointStateError.Add(1)
			return false
		}
		return true
	}
	// Historical injected/non-scheduler legacy source path; never a modern
	// scheduler result. Production scheduler producers always supply capture.
	if err := s.pointSvc.UpdateReadResult(ctx, cv.PointID, cv.Value, cv.Error); err != nil {
		s.pointStateError.Add(1)
	}
	return true
}
