package modbusshare

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFormatRecoveryAlignedAndReplayStillRejectLateLifecycle(t *testing.T) {
	for _, replay := range []bool{false, true} {
		r, s, _, store := setupReconcilerFixture(t)
		r.revStore = recoveryUnknownStore{store}
		settings := s.Settings()
		pending := HydrationState{State: HydrationStatePending, WorkspaceID: "ws-1", WorkspaceRevision: "rev-1", SettingsRevision: settings.SettingsRevision}
		s.SetHydrationState(pending)
		fence, err := s.CaptureFormatRecoveryFence("ws-1", "rev-1", settings.SettingsRevision)
		require.NoError(t, err)
		req := ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: settings.SettingsRevision}
		if replay {
			out, e := r.ReconcileFormatRecovery(t.Context(), req, fence, "clean")
			require.NoError(t, e)
			require.Equal(t, "aligned", out.Outcome)
		}
		store.beforeGetDesired = func() { s.SetHydrationState(pending) }
		_, err = r.ReconcileFormatRecovery(t.Context(), req, fence, "clean")
		require.Error(t, err)
		require.Zero(t, store.commitCalls)
		state, e := s.CheckHydration(t.Context())
		require.NoError(t, e)
		require.False(t, state.Readiness)
		require.Empty(t, state.ReadinessToken)
	}
}
