package modbusshare

import "context"

// InvalidatePointFormatRevision clears only matching live readiness after an
// atomic point-format commit advanced the existing durable workspace authority.
func (s *Service) InvalidatePointFormatRevision(ctx context.Context, workspaceID, expected, next string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if workspaceID == "" || expected == "" || next == "" {
		return NewError(ErrCodeHydrationRequired, "point format invalidation is unavailable", true)
	}
	if s.hydration.WorkspaceID != workspaceID || s.hydration.WorkspaceRevision != expected {
		return NewError(ErrCodeRevisionConflict, "point format readiness revision changed", true)
	}
	s.advanceHydrationGenerationLocked()
	s.hydration.State = HydrationStatePending
	s.hydration.Readiness = false
	s.hydration.ReadinessToken = ""
	if s.state == shareStateRunning {
		if s.server.Stop() != nil {
			s.state = HydrationStateFailed
			return NewError(ErrCodeHydrationRequired, "point format listener invalidation is uncertain", true)
		}
		s.state = shareStateStopped
	}
	if s.workspaceRevisionStore == nil {
		return NewError(ErrCodeHydrationRequired, "point format durability is unavailable", true)
	}
	revision, dirty, err := s.workspaceRevisionStore.GetRevision(ctx, workspaceID)
	if err != nil || revision != next || !dirty {
		return NewError(ErrCodeRevisionConflict, "point format durable revision changed", true)
	}

	s.hydration.WorkspaceRevision = next
	return nil
}
