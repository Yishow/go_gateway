package device

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type acquisitionSettingsRepository interface {
	Repository
	UpdateAcquisitionSettings(ctx context.Context, id string, settings schema.DeviceAcquisitionSettings) error
}

func acquisitionTestDevice(id string) *schema.Device {
	return &schema.Device{ID: id, Name: id, Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusDraft, ConnectionConfig: "{}"}
}

func acquisitionRepos(t *testing.T) map[string]acquisitionSettingsRepository {
	db := setupTestDB(t)
	t.Cleanup(func() { _ = db.Close() })
	return map[string]acquisitionSettingsRepository{"sql": NewSQLRepository(db), "memory": NewMemoryRepository()}
}

func TestDeviceAcquisitionSettingsDefaultsAndRoundTrip(t *testing.T) {
	for name, repo := range acquisitionRepos(t) {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			require.NoError(t, repo.Create(ctx, acquisitionTestDevice("dev")))
			got, err := repo.GetByID(ctx, "dev")
			require.NoError(t, err)
			require.Nil(t, got.StaleThresholdMs, "nil means interval x5, at least 10s")
			require.Zero(t, got.FrozenThresholdMs, "frozen detection disabled by default")
			require.False(t, got.AcquisitionAutoInterval)

			stale := 15000
			require.NoError(t, repo.UpdateAcquisitionSettings(ctx, "dev", schema.DeviceAcquisitionSettings{StaleThresholdMs: &stale, FrozenThresholdMs: 60000, AcquisitionAutoInterval: true}))
			got, err = repo.GetByID(ctx, "dev")
			require.NoError(t, err)
			require.Equal(t, 15000, *got.StaleThresholdMs)
			require.Equal(t, 60000, got.FrozenThresholdMs)
			require.True(t, got.AcquisitionAutoInterval)
			listed, err := repo.List(ctx, ListFilter{})
			require.NoError(t, err)
			require.Len(t, listed, 1)
			require.Equal(t, 15000, *listed[0].StaleThresholdMs)
			require.Equal(t, 60000, listed[0].FrozenThresholdMs)
			require.True(t, listed[0].AcquisitionAutoInterval)

			// Whole-row updates (stats, readiness) read-modify-write devices; they
			// must not overwrite thresholds saved concurrently through the
			// dedicated setter.
			stale2 := 20000
			stalecopy := *got
			require.NoError(t, repo.UpdateAcquisitionSettings(ctx, "dev", schema.DeviceAcquisitionSettings{StaleThresholdMs: &stale2}))
			stalecopy.Name = "renamed"
			require.NoError(t, repo.Update(ctx, &stalecopy))
			got, err = repo.GetByID(ctx, "dev")
			require.NoError(t, err)
			require.Equal(t, "renamed", got.Name)
			require.Equal(t, 20000, *got.StaleThresholdMs)
			require.Zero(t, got.FrozenThresholdMs)
			require.False(t, got.AcquisitionAutoInterval)

			require.NoError(t, repo.UpdateAcquisitionSettings(ctx, "dev", schema.DeviceAcquisitionSettings{}))
			got, err = repo.GetByID(ctx, "dev")
			require.NoError(t, err)
			require.Nil(t, got.StaleThresholdMs, "clearing restores the derived default")
		})
	}
}

func TestDeviceAcquisitionSettingsRejectInvalidAndMissing(t *testing.T) {
	for name, repo := range acquisitionRepos(t) {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			require.NoError(t, repo.Create(ctx, acquisitionTestDevice("dev")))
			zero, negative := 0, -1
			for _, bad := range []schema.DeviceAcquisitionSettings{
				{StaleThresholdMs: &zero},
				{StaleThresholdMs: &negative},
				{FrozenThresholdMs: -1},
			} {
				require.ErrorIs(t, repo.UpdateAcquisitionSettings(ctx, "dev", bad), ErrValidation)
			}
			got, err := repo.GetByID(ctx, "dev")
			require.NoError(t, err)
			require.Nil(t, got.StaleThresholdMs)
			require.Error(t, repo.UpdateAcquisitionSettings(ctx, "missing", schema.DeviceAcquisitionSettings{}))
		})
	}
}
