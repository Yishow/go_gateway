package modbusshare

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func TestService_HydrationFailureStopsRunningListener(t *testing.T) {
	svc := NewService(tag.NewService(tag.NewMemoryRepository()), 128)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	port := reservePort(t)
	require.NoError(t, svc.ApplySettings(context.Background(), Settings{Enabled: true, BindAddress: "127.0.0.1", Port: port, SlaveID: 1, CapacityRegisters: 64}))
	require.True(t, svc.Status().Enabled)

	svc.SetHydrationState(HydrationState{State: HydrationStateFailed, Readiness: false})
	status := svc.Status()
	require.False(t, status.Enabled)
	require.Equal(t, "stopped", status.LifecycleState)
	require.Equal(t, "stopped", status.BindState)
	require.Equal(t, 0, status.Port)
	t.Cleanup(func() { _ = svc.CloseRuntime() })
}
