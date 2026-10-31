package modbusshare

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFormatRecoveryFenceRejectsHydrationABA(t *testing.T) {
	s := NewService(nil, 64)
	settings := s.Settings()
	state := HydrationState{State: HydrationStatePending, WorkspaceID: "ws", WorkspaceRevision: "rev", SettingsRevision: settings.SettingsRevision}
	s.SetHydrationState(state)
	proof, err := s.CaptureFormatRecoveryFence("ws", "rev", settings.SettingsRevision)
	require.NoError(t, err)
	require.NoError(t, s.CheckFormatRecoveryFence(t.Context(), proof))
	failed := state
	failed.State = HydrationStateFailed
	s.SetHydrationState(failed)
	s.SetHydrationState(state)
	require.Error(t, s.CheckFormatRecoveryFence(t.Context(), proof), "a later identical pending state is not the captured lifecycle")
}

func TestFormatRecoveryFenceRejectsEveryExistingInvalidationBoundary(t *testing.T) {
	for _, name := range []string{"membership", "format", "failed", "settings", "hydration_error"} {
		t.Run(name, func(t *testing.T) {
			s := NewService(nil, 64)
			settings := s.Settings()
			state := HydrationState{State: HydrationStateReady, Readiness: true, WorkspaceID: "ws", WorkspaceRevision: "rev", SettingsRevision: settings.SettingsRevision, ReadinessToken: "old"}
			s.SetHydrationState(state)
			proof, err := s.CaptureFormatRecoveryFence("ws", "rev", settings.SettingsRevision)
			require.NoError(t, err)
			switch name {
			case "membership":
				require.NoError(t, s.InvalidateWorkspaceMembership(t.Context(), "ws"))
			case "format":
				s.SetWorkspaceRevisionStore(&pointFormatRevisionStore{revision: "next", dirty: true})
				require.NoError(t, s.InvalidatePointFormatRevision(t.Context(), "ws", "rev", "next"))
			case "failed":
				s.FailClosed(state)
			case "settings":
				require.NoError(t, s.ApplySettings(t.Context(), settings))
			case "hydration_error":
				s.SetHydrationError(assertiveRecoveryError{})
				s.SetHydrationError(nil)
			}
			s.SetHydrationState(state)
			require.Error(t, s.CheckFormatRecoveryFence(t.Context(), proof))
		})
	}
}

type assertiveRecoveryError struct{}

func (assertiveRecoveryError) Error() string { return "injected hydration failure" }
