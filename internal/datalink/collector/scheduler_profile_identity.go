package collector

import "go-gateway/internal/datalink/schema"

// identityPoint reconstructs only the configuration captured by this install,
// never editable/latest state. The profile and configRevision belong to it.
func (pt pointInfo) identityPoint() schema.Point {
	p := schema.Point{ID: pt.ID, DeviceID: pt.DeviceID, RevisionID: pt.configRevision, Address: pt.Address, Function: pt.Function, DataType: pt.DataType, DataFormat: pt.DataFormat, Mode: pt.Mode, PollingGroupID: pollingGroupPointer(pt.PollingGroupID)}
	if pt.readProfile.Present() {
		r := pt.readProfile.Resolution()
		p.FormatMode = r.Mode
		p.FormatOverride = r.ConfiguredFormat
		p.FormatVersion = r.FormatVersion
	}
	return p
}

func (pt pointInfo) modernObservationLease() uint64 {
	if pt.readProfile.Present() && pt.readProfile.Resolution().Mode != schema.PointFormatLegacy {
		return pt.observationLease
	}
	return 0
}

// AcceptsInstallationSample checks an execution incarnation, not just a point
// revision. A removed/reinstalled same-revision point cannot accept old samples.
func (s *Scheduler) AcceptsInstallationSample(pointID string, lease uint64, scope string) bool {
	if s == nil || lease == 0 || (scope == "" || scope != s.installationScope) {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	pt, ok := s.pointInfos[pointID]
	return ok && !pt.recoveryPending && pt.observationLease == lease
}

func (s *Scheduler) InstallationScope() string {
	if s == nil {
		return ""
	}
	return s.installationScope
}

// Direct previews can start without a slot. The existing epoch proves absence
// has not changed during provider/wire I/O, without an unbounded tombstone map.
func (s *Scheduler) AcceptsDirectReadProof(proof PointMutationLease, scope string) bool {
	if s == nil || scope == "" || scope != s.installationScope || proof.PointID == "" {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	current := s.pointInfos[proof.PointID]
	if current.recoveryPending || current.observationLease != proof.Lease {
		return false
	}
	if proof.Lease != 0 {
		return current.configRevision == proof.Revision
	}
	return s.pointObservationLease == proof.Epoch
}

// A truly unscoped legacy API result cannot become proof for a modern install.
func (s *Scheduler) AcceptsLegacyUnscopedSample(pointID string) bool {
	if s == nil {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	pt, ok := s.pointInfos[pointID]
	return ok && !pt.recoveryPending && (!pt.readProfile.Present() || pt.readProfile.Resolution().Mode == schema.PointFormatLegacy)
}

func (s *Scheduler) observationNeedsCurrentInstall(captured pointInfo) bool {
	if captured.modernObservationLease() != 0 {
		return true
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	current, ok := s.pointInfos[captured.ID]
	return ok && current.modernObservationLease() != 0
}
