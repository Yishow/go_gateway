package collector

// CollectionPointStatus is a detached observation of an installed scheduler
// slot, not a new lifecycle authority. Readable excludes quarantined recovery.
type CollectionPointStatus struct {
	DeviceID string
	Running  bool
	Readable bool
}

// CollectionStatus observes slots and their ticker/lifecycle under one lock.
// Stopped and quarantined slots remain visible, but cannot claim collection.
func (s *Scheduler) CollectionStatus() map[string]CollectionPointStatus {
	result := make(map[string]CollectionPointStatus)
	if s == nil {
		return result
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for id, point := range s.pointInfos {
		_, devicePresent := s.deviceConfigs[point.DeviceID]
		readable := devicePresent && !point.recoveryPending
		ticker := s.groupTickers[point.PollingGroupID]
		scheduled := ticker != nil && ticker.group != nil && ticker.group.Enabled
		result[id] = CollectionPointStatus{DeviceID: point.DeviceID, Readable: readable, Running: s.running && readable && scheduled}
	}
	return result
}
