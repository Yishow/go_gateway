package modbusshare

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

type pointFormatRevisionStore struct {
	WorkspaceRevisionStore
	revision string
	dirty    bool
}

func (s *pointFormatRevisionStore) GetRevision(context.Context, string) (string, bool, error) {
	return s.revision, s.dirty, nil
}
func TestPointFormatReadinessInvalidationIsRevisionGuarded(t *testing.T) {
	svc := NewService(nil, 64)
	store := &pointFormatRevisionStore{revision: "new", dirty: true}
	svc.SetWorkspaceRevisionStore(store)
	original := HydrationState{State: HydrationStateReady, WorkspaceID: "ws", WorkspaceRevision: "old", Readiness: true, ReadinessToken: "proof"}
	svc.SetHydrationState(original)
	require.Error(t, svc.InvalidatePointFormatRevision(t.Context(), "foreign", "old", "new"))
	state, err := svc.CheckHydration(t.Context())
	require.NoError(t, err)
	require.Equal(t, original, state)
	require.NoError(t, svc.InvalidatePointFormatRevision(t.Context(), "ws", "old", "new"))
	state, err = svc.CheckHydration(t.Context())
	require.NoError(t, err)
	require.False(t, state.Readiness)
	require.Empty(t, state.ReadinessToken)
	require.Equal(t, "new", state.WorkspaceRevision)
	require.Equal(t, HydrationStatePending, state.State)
	latest := original
	latest.WorkspaceRevision = "latest"
	svc.SetHydrationState(latest)
	require.Error(t, svc.InvalidatePointFormatRevision(t.Context(), "ws", "old", "new"))
	state, err = svc.CheckHydration(t.Context())
	require.NoError(t, err)
	require.Equal(t, latest, state)
}

func TestPointFormatInvalidationUnknownDurabilityStillClosesMatchingReadiness(t *testing.T) {
	svc := NewService(nil, 64)
	store := &pointFormatRevisionStore{revision: "unconfirmed", dirty: true}
	svc.SetWorkspaceRevisionStore(store)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws", WorkspaceRevision: "old", Readiness: true, ReadinessToken: "old-proof"})
	require.Error(t, svc.InvalidatePointFormatRevision(t.Context(), "ws", "old", "new"))
	state, err := svc.CheckHydration(t.Context())
	require.NoError(t, err)
	require.False(t, state.Readiness)
	require.Empty(t, state.ReadinessToken)
}
