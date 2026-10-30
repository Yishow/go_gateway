package runtime

// SuspendPointRevision only removes a matching scheduler installation. Runtime
// metadata remains readable while configuration is pending; this does not install
// the new decoder or claim acquisition fingerprint migration has completed.
func (s *Service) SuspendPointRevision(pointID, expectedRevision string) bool {
	return s != nil && s.scheduler != nil && s.scheduler.RemovePointRevision(pointID, expectedRevision)
}
