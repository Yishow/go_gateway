package runtime

import (
	"go-gateway/internal/datalink/collector"
	"slices"
)

// captureCollectedBindings linearizes raw acceptance with the installed target
// bindings. Retain the accepted immutable values across telemetry SQL I/O; do not
// reinterpret a previously accepted raw value with a later mapping definition.
func (s *Service) captureCollectedBindings(cv collector.CollectedValue) ([]mappingBinding, bool) {
	s.mappingMu.RLock()
	defer s.mappingMu.RUnlock()
	if cv.InstallationLease == 0 && s.scheduler != nil && !s.scheduler.AcceptsLegacyUnscopedSample(cv.PointID) {
		return nil, false
	}
	if cv.InstallationLease != 0 && (s.scheduler == nil || !s.scheduler.AcceptsInstallationSample(cv.PointID, cv.InstallationLease, cv.InstallationScope)) {
		return nil, false
	}
	return slices.Clone(s.mappingIndex[cv.PointID]), true
}
