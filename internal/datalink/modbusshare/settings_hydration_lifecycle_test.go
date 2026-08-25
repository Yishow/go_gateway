package modbusshare

import (
	"context"
	"testing"

	datalinksettings "go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_ApplySettings_PersistsEnabledButWaitsForHydrationBeforeBinding(t *testing.T) {
	ctx := context.Background()
	repo := datalinksettings.NewMemoryRepository()
	svc := NewService(tag.NewService(tag.NewMemoryRepository()), 128)
	require.NoError(t, svc.SetSettingsRepository(repo))
	svc.SetHydrationState(HydrationState{State: HydrationStatePending, Readiness: false})

	settings := Settings{Enabled: true, BindAddress: "127.0.0.1", Port: 15020, SlaveID: 1, CapacityRegisters: 64, SettingsRevision: "hydration-pending"}
	require.NoError(t, svc.ApplySettingsCAS(ctx, settings, ""))

	persisted, err := repo.Get(ctx, settingsKey)
	require.NoError(t, err)
	assert.Equal(t, true, persisted.Value.(Settings).Enabled)
	assert.Equal(t, 0, svc.server.Port(), "listener must remain stopped while hydration is pending")
	assert.False(t, svc.Status().Enabled)
}
