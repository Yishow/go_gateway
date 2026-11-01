package collector

// BeginFullAcquisitionGeometryMutation captures the original installed device
// set and quarantines it together with requested devices under one short lock.
// Later devices are not targets of this full projection's removal operation.
func (s *Scheduler) BeginFullAcquisitionGeometryMutation(requested []string) map[string]AcquisitionGeometryMutation {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make(map[string]AcquisitionGeometryMutation, len(s.deviceConfigs)+len(requested))
	for id := range s.deviceConfigs {
		result[id] = s.beginAcquisitionGeometryMutationLocked(id)
	}
	for _, id := range requested {
		if _, ok := result[id]; !ok {
			result[id] = s.beginAcquisitionGeometryMutationLocked(id)
		}
	}
	return result
}

// RemoveAcquisitionDevice removes only the originally captured complete device
// installation. The metadata/scope callback runs after all checks under the
// existing runtime -> scheduler lock order and must do assignments only.
func (s *Scheduler) RemoveAcquisitionDevice(mutation AcquisitionGeometryMutation, assign func()) bool {
	token := mutation.token
	if s == nil || token == nil || token.owner != s || assign == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.acquisitionGeometry[token.deviceID]
	cfg, present := s.deviceConfigs[token.deviceID]
	if state == nil || state.publication != token.marker || len(state.pending) != 1 || cfg != token.before.device || present != token.before.present || !s.acquisitionSlotsMatchLocked(token.deviceID, token.before.slots, false) {
		return false
	}
	if _, ok := state.pending[token]; !ok {
		return false
	}
	assign()
	s.removeDeviceLocked(token.deviceID)
	return true
}

// PublishEmptyAcquisitionProjection preserves the ordinary no-device install
// path without manufacturing a complete graph. Foreign new devices reject it.
func (s *Scheduler) PublishEmptyAcquisitionProjection(assign func()) bool {
	if s == nil || assign == nil || !s.mu.TryLock() {
		return false
	}
	defer s.mu.Unlock()
	if len(s.deviceConfigs) != 0 {
		return false
	}
	assign()
	return true
}
