package modbusshare

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFormatRecoveryPublicationReplaysOnlyOwnedCurrentAcceptance(t *testing.T) {
	s := NewService(nil, 64)
	settings := s.Settings()
	pending := HydrationState{State: HydrationStatePending, WorkspaceID: "ws", WorkspaceRevision: "rev", SettingsRevision: settings.SettingsRevision}
	s.SetHydrationState(pending)
	fence, err := s.CaptureFormatRecoveryFence("ws", "rev", settings.SettingsRevision)
	require.NoError(t, err)
	activations := 0
	state, err := s.PublishFormatRecovery(t.Context(), fence, "operation", 7, func() bool { activations++; return true })
	require.NoError(t, err)
	require.True(t, state.Readiness)
	require.NotEmpty(t, state.ReadinessToken)
	replay, ok := s.ReplayFormatRecovery("operation", func(lease uint64) bool { return lease == 7 })
	require.True(t, ok)
	require.Equal(t, state, replay)
	require.Equal(t, 1, activations)
	_, ok = s.ReplayFormatRecovery("other", func(uint64) bool { return true })
	require.False(t, ok)
	_, ok = s.ReplayFormatRecovery("operation", func(uint64) bool { return false })
	require.False(t, ok)
	s.FailClosed(state)
	_, ok = s.ReplayFormatRecovery("operation", func(uint64) bool { return true })
	require.False(t, ok)
}

func TestFormatRecoveryLatePublicationNeverActivates(t *testing.T) {
	for _, kind := range []string{"changed-generation", "installation-rejected"} {
		t.Run(kind, func(t *testing.T) {
			s := NewService(nil, 64)
			settings := s.Settings()
			pending := HydrationState{State: HydrationStatePending, WorkspaceID: "ws", WorkspaceRevision: "rev", SettingsRevision: settings.SettingsRevision}
			s.SetHydrationState(pending)
			fence, err := s.CaptureFormatRecoveryFence("ws", "rev", settings.SettingsRevision)
			require.NoError(t, err)
			if kind == "changed-generation" {
				s.SetHydrationState(pending)
			}
			calls := 0
			_, err = s.PublishFormatRecovery(t.Context(), fence, "op", 7, func() bool { calls++; return false })
			require.Error(t, err)
			if kind == "changed-generation" {
				require.Zero(t, calls)
			}
			current, e := s.CheckHydration(t.Context())
			require.NoError(t, e)
			require.False(t, current.Readiness)
			require.Empty(t, current.ReadinessToken)
		})
	}
}

func TestFormatRecoveryPublicationNeverWaitsForLifecycleLock(t *testing.T) {
	s := NewService(nil, 64)
	settings := s.Settings()
	s.SetHydrationState(HydrationState{State: HydrationStatePending, WorkspaceID: "ws", WorkspaceRevision: "rev", SettingsRevision: settings.SettingsRevision})
	fence, err := s.CaptureFormatRecoveryFence("ws", "rev", settings.SettingsRevision)
	require.NoError(t, err)
	s.mu.Lock()
	done := make(chan error, 1)
	go func() {
		_, e := s.PublishFormatRecovery(t.Context(), fence, "op", 7, func() bool { return true })
		done <- e
	}()
	select {
	case err = <-done:
		s.mu.Unlock()
		require.Error(t, err)
	case <-time.After(time.Second):
		s.mu.Unlock()
		<-done
		t.Fatal("publication waited on lifecycle lock while caller owns SQL writer reservation")
	}
}

func TestFormatRecoveryCleanupRevokesOnlyItsPublication(t *testing.T) {
	for _, foreign := range []bool{false, true} {
		s := NewService(nil, 64)
		settings := s.Settings()
		s.SetHydrationState(HydrationState{State: HydrationStatePending, WorkspaceID: "ws", WorkspaceRevision: "rev", SettingsRevision: settings.SettingsRevision})
		fence, err := s.CaptureFormatRecoveryFence("ws", "rev", settings.SettingsRevision)
		require.NoError(t, err)
		state, err := s.PublishFormatRecovery(t.Context(), fence, "op", 7, func() bool { return true })
		require.NoError(t, err)
		if foreign {
			replacement := state
			replacement.ReadinessToken = "foreign"
			s.SetHydrationState(replacement)
		}
		revoked := s.RevokeFormatRecoveryPublication("op", state.ReadinessToken, 7)
		require.Equal(t, !foreign, revoked)
		current, e := s.CheckHydration(t.Context())
		require.NoError(t, e)
		if foreign {
			require.Equal(t, "foreign", current.ReadinessToken)
			require.True(t, current.Readiness)
		} else {
			require.False(t, current.Readiness)
			require.Empty(t, current.ReadinessToken)
		}
	}
}
