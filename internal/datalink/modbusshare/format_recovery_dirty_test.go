package modbusshare

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

type recoveryDirtySpy struct {
	WorkspaceRevisionStore
	before func()
	err    error
}

func (s recoveryDirtySpy) MarkDirty(context.Context, string) error { s.before(); return s.err }

func TestFormatRecoveryDirtyEntryRevokesBeforeStoreEvenOnFailure(t *testing.T) {
	for _, fail := range []bool{false, true} {
		s := NewService(nil, 64)
		settings := s.Settings()
		ready := HydrationState{State: HydrationStateReady, Readiness: true, ReadinessToken: "old", WorkspaceID: "ws", WorkspaceRevision: "rev", SettingsRevision: settings.SettingsRevision}
		s.SetHydrationState(ready)
		fence, err := s.CaptureFormatRecoveryFence("ws", "rev", settings.SettingsRevision)
		require.NoError(t, err)
		spy := recoveryDirtySpy{before: func() {
			require.Error(t, s.CheckFormatRecoveryFence(t.Context(), fence), "DB must never observe an unrevoked lifecycle")
			current, e := s.CheckHydration(t.Context())
			require.NoError(t, e)
			require.False(t, current.Readiness)
		}}
		if fail {
			spy.err = errors.New("store failed")
		}
		r := NewReconciler(s, spy)
		err = r.markDirty(t.Context(), "ws", ready)
		require.Equal(t, fail, err != nil)
	}
}
