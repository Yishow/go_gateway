package modbusshare

import "context"

// FormatRecoveryFence is an internal, process-local observation of the existing
// hydration/settings authority. Its fields cannot be injected through JSON.
type FormatRecoveryFence struct {
	epoch            uint64
	generation       uint64
	hydration        HydrationState
	settingsRevision string
}

func (s *Service) CaptureFormatRecoveryFence(workspaceID, workspaceRevision, settingsRevision string) (FormatRecoveryFence, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.hydrationGeneration == ^uint64(0) || s.hydrationErr != nil || workspaceID == "" || workspaceRevision == "" || settingsRevision == "" || s.hydration.WorkspaceID != workspaceID || s.hydration.WorkspaceRevision != workspaceRevision || s.settings.SettingsRevision != settingsRevision || (s.hydration.State != HydrationStatePending && s.hydration.State != HydrationStateReady) {
		return FormatRecoveryFence{}, NewError(ErrCodeHydrationRequired, "format recovery scope is unavailable", true)
	}
	return FormatRecoveryFence{epoch: s.lifecycleEpoch, generation: s.hydrationGeneration, hydration: s.hydration, settingsRevision: settingsRevision}, nil
}

func (s *Service) CheckFormatRecoveryFence(_ context.Context, fence FormatRecoveryFence) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.matchesFormatRecoveryFenceLocked(fence) {
		return NewError(ErrCodeRevisionConflict, "format recovery lifecycle changed", true)
	}
	return nil
}

func (s *Service) matchesFormatRecoveryFenceLocked(fence FormatRecoveryFence) bool {
	return fence.hydration.WorkspaceID != "" && s.hydrationErr == nil && s.hydrationGeneration == fence.generation && s.lifecycleEpoch == fence.epoch && s.settings.SettingsRevision == fence.settingsRevision && s.hydration == fence.hydration
}
func (s *Service) advanceHydrationGenerationLocked() {
	s.formatRecoveryCompletion = nil
	if s.hydrationGeneration < ^uint64(0) {
		s.hydrationGeneration++
	}
}
