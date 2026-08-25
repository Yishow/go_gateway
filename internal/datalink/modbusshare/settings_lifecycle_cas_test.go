package modbusshare

import (
	"context"
	"testing"

	datalinksettings "go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_StartCASChecksRevisionBeforeSamePortIdempotency(t *testing.T) {
	ctx := context.Background()
	svc := NewService(tag.NewService(tag.NewMemoryRepository()), 128)
	repo := datalinksettings.NewMemoryRepository()
	require.NoError(t, svc.SetSettingsRepository(repo))
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	port := reservePort(t)
	require.NoError(t, svc.ApplySettings(ctx, Settings{Enabled: true, BindAddress: "127.0.0.1", Port: port, SlaveID: 1, CapacityRegisters: 64, SettingsRevision: "settings-current"}))
	t.Cleanup(func() { _ = svc.Stop() })

	err := svc.StartCAS(ctx, port, "settings-stale")
	var conflict *Error
	require.ErrorAs(t, err, &conflict)
	assert.Equal(t, ErrCodeRevisionConflict, conflict.Code)
	assert.Equal(t, "settings-current", svc.Settings().SettingsRevision)

	require.NoError(t, svc.StartCAS(ctx, port, "settings-current"))
}
