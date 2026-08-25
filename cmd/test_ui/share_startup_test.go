package main

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/modbusshare"
	datalinksettings "go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func TestStartConfiguredShareListenerPreservesReadyHydration(t *testing.T) {
	ctx := context.Background()
	repo := datalinksettings.NewMemoryRepository()
	port := freeTCPPort(t)
	require.NoError(t, repo.Set(ctx, datalinksettings.KeyModbusShare, modbusshare.Settings{
		Enabled: true, BindAddress: "127.0.0.1", Port: port, SlaveID: 1,
		CapacityRegisters: 64, SettingsRevision: "startup-ready",
	}))

	svc := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 128)
	require.NoError(t, svc.SetSettingsRepository(repo))
	require.NoError(t, svc.LoadPersistedSettings(ctx))
	svc.SetHydrationState(modbusshare.HydrationState{
		State: modbusshare.HydrationStateReady, WorkspaceID: "ws-startup",
		WorkspaceRevision: "rev-startup", Readiness: true,
	})
	t.Cleanup(func() { _ = svc.CloseRuntime() })

	startConfiguredShareListener(ctx, svc)

	status := svc.Status()
	require.True(t, status.Enabled)
	require.Equal(t, "pass", status.BindState)
	require.Equal(t, modbusshare.HydrationStateReady, status.HydrationState)
	require.True(t, status.Readiness)
	require.Equal(t, port, status.Port)
}

func TestConfiguredShareRuntimeReconcilerFailsClosed(t *testing.T) {
	ctx := context.Background()
	svc := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 64)
	svc.SetHydrationState(modbusshare.HydrationState{
		State: modbusshare.HydrationStateReady, WorkspaceID: "ws-fail-closed",
		WorkspaceRevision: "rev-fail-closed", Readiness: true,
	})
	reconciler := configuredShareRuntimeReconciler{share: svc}

	reconciler.FailClosed(ctx)

	state, err := svc.CheckHydration(ctx)
	require.NoError(t, err)
	require.Equal(t, modbusshare.HydrationStateFailed, state.State)
	require.False(t, state.Readiness)
}
