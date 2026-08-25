package modbusshare

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"testing"

	"go-gateway/internal/datalink"
	datalinksettings "go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func reservePort(t *testing.T) int {
	t.Helper()
	var listenConfig net.ListenConfig
	ln, err := listenConfig.Listen(context.Background(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

type failingShareSettingsRepository struct{}

func (failingShareSettingsRepository) Get(context.Context, string) (*datalinksettings.SettingItem, error) {
	return nil, errors.New("settings unavailable")
}
func (failingShareSettingsRepository) Set(context.Context, string, interface{}) error {
	return errors.New("settings unavailable")
}
func (failingShareSettingsRepository) List(context.Context) ([]*datalinksettings.SettingItem, error) {
	return nil, errors.New("settings unavailable")
}

type failingCASSettingsRepository struct{ failingShareSettingsRepository }

func (failingCASSettingsRepository) SetIfRevision(context.Context, string, string, interface{}) error {
	return errors.New("database unavailable")
}

func TestService_ApplySettingsCASStorageFailureIsNotRevisionConflict(t *testing.T) {
	svc := NewService(tag.NewService(tag.NewMemoryRepository()), 128)
	repo := failingCASSettingsRepository{}
	require.NoError(t, svc.SetSettingsRepository(repo))
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	err := svc.ApplySettingsCAS(context.Background(), Settings{Enabled: false}, "")
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeSettingsUpdateFailed, shareErr.Code)
}

func TestService_LoadPersistedSettingsFailureLeavesListenerStopped(t *testing.T) {
	svc := NewService(tag.NewService(tag.NewMemoryRepository()), 128)
	require.NoError(t, svc.SetSettingsRepository(failingShareSettingsRepository{}))
	err := svc.LoadPersistedSettings(context.Background())
	require.Error(t, err)
	assert.False(t, svc.Status().ConfiguredEnabled)
	assert.False(t, svc.Status().Enabled)
	assert.Equal(t, "disabled", svc.Status().BindState)
}

func TestService_ApplySettings_DisabledByDefaultDoesNotStartListener(t *testing.T) {
	ctx := context.Background()
	tagRepo := tag.NewMemoryRepository()
	tagSvc := tag.NewService(tagRepo)
	svc := NewService(tagSvc, 65536)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})

	settings := DefaultSettings() // Enabled: false
	err := svc.ApplySettings(ctx, settings)
	require.NoError(t, err)

	status := svc.Status()
	assert.False(t, status.Enabled)
	assert.Equal(t, "disabled", status.BindState)
	assert.Equal(t, 0, status.Port)
}

func TestService_ApplySettings_EnabledStartsListener(t *testing.T) {
	ctx := context.Background()
	tagRepo := tag.NewMemoryRepository()
	tagSvc := tag.NewService(tagRepo)
	svc := NewService(tagSvc, 65536)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	t.Cleanup(func() {
		_ = svc.Stop()
	})

	port := reservePort(t)
	settings := Settings{
		Enabled:           true,
		BindAddress:       "127.0.0.1",
		Port:              port,
		SlaveID:           1,
		CapacityRegisters: 32768,
		SettingsRevision:  "set-rev-1",
	}

	err := svc.ApplySettings(ctx, settings)
	require.NoError(t, err)

	status := svc.Status()
	assert.True(t, status.Enabled)
	assert.Equal(t, "pass", status.BindState)
	assert.Equal(t, port, status.Port)
	assert.Equal(t, uint8(1), status.SlaveID)
	assert.Equal(t, "set-rev-1", status.SettingsRevision)
}

func TestService_ApplySettings_BindConflictReportsTypedFailureWithoutFallbackTo5020(t *testing.T) {
	ctx := context.Background()
	tagRepo := tag.NewMemoryRepository()
	tagSvc := tag.NewService(tagRepo)
	svc := NewService(tagSvc, 65536)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	t.Cleanup(func() {
		_ = svc.Stop()
	})

	// Pre-occupy a port with a separate listener on all interfaces
	var listenConfig net.ListenConfig
	ln, err := listenConfig.Listen(context.Background(), "tcp", ":0")
	require.NoError(t, err)
	defer ln.Close()
	occupiedPort := ln.Addr().(*net.TCPAddr).Port

	settings := Settings{
		Enabled:           true,
		BindAddress:       "127.0.0.1",
		Port:              occupiedPort,
		SlaveID:           1,
		CapacityRegisters: 32768,
		SettingsRevision:  "set-conflict",
	}

	err = svc.ApplySettings(ctx, settings)
	require.Error(t, err)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeListenerBindFailed, shareErr.Code)

	// Status must report failed, not silently fallen back to 5020!
	status := svc.Status()
	assert.False(t, status.Enabled)
	assert.Equal(t, "fail", status.BindState)
	assert.NotEqual(t, 5020, status.Port)
	require.NotNil(t, status.LastFailure)
	assert.Equal(t, ErrCodeListenerBindFailed, status.LastFailure.Code)
	require.NoError(t, svc.ApplySettings(ctx, Settings{Enabled: false}))
	assert.Nil(t, svc.Status().LastFailure)
}

func TestService_ApplySettings_DisableStopsListener(t *testing.T) {
	ctx := context.Background()
	tagRepo := tag.NewMemoryRepository()
	tagSvc := tag.NewService(tagRepo)
	svc := NewService(tagSvc, 65536)
	t.Cleanup(func() {
		_ = svc.Stop()
	})

	port := reservePort(t)
	err := svc.ApplySettings(ctx, Settings{
		Enabled:           true,
		BindAddress:       "127.0.0.1",
		Port:              port,
		SlaveID:           1,
		CapacityRegisters: 32768,
	})
	require.NoError(t, err)
	require.True(t, svc.Status().Enabled)

	// Now disable
	err = svc.ApplySettings(ctx, Settings{
		Enabled: false,
	})
	require.NoError(t, err)

	status := svc.Status()
	assert.False(t, status.Enabled)
	assert.Equal(t, "disabled", status.BindState)
	assert.Equal(t, 0, status.Port)
}

func TestService_ApplySettings_CapacityIsPersistedContract(t *testing.T) {
	svc := NewService(tag.NewService(tag.NewMemoryRepository()), 16)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	settings := Settings{Enabled: true, BindAddress: "127.0.0.1", Port: reservePort(t), SlaveID: 1, CapacityRegisters: 8, SettingsRevision: "capacity-1"}
	require.NoError(t, svc.ApplySettings(context.Background(), settings))
	defer svc.Stop()
	assert.Equal(t, 8, svc.Status().CapacityRegisters)
}

func TestService_ApplySettings_RejectsCapacityBeyondBackingBank(t *testing.T) {
	svc := NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
	settings := Settings{Enabled: true, BindAddress: "127.0.0.1", Port: reservePort(t), SlaveID: 1, CapacityRegisters: 32769, SettingsRevision: "capacity-too-large"}
	err := svc.ApplySettings(context.Background(), settings)
	require.Error(t, err)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeCapacityExceeded, shareErr.Code)
	assert.False(t, svc.Status().Enabled)
}

func TestService_PersistedSettingsRestoreAcrossServiceRestart(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	repo := datalinksettings.NewSQLRepository(db)
	ctx := context.Background()
	port := reservePort(t)
	persisted := Settings{
		Enabled:           true,
		BindAddress:       "127.0.0.1",
		Port:              port,
		SlaveID:           23,
		CapacityRegisters: 64,
		SettingsRevision:  "restart-revision",
	}

	svc1 := NewService(tag.NewService(tag.NewMemoryRepository()), 128)
	svc1.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	require.NoError(t, svc1.SetSettingsRepository(repo))
	require.NoError(t, svc1.ApplySettings(ctx, persisted))
	require.Equal(t, port, svc1.Status().Port)
	require.NoError(t, svc1.server.Stop())

	svc2 := NewService(tag.NewService(tag.NewMemoryRepository()), 128)
	require.NoError(t, svc2.SetSettingsRepository(repo))
	require.NoError(t, svc2.LoadPersistedSettings(ctx))
	assert.True(t, svc2.Settings().Enabled)
	assert.False(t, svc2.Status().Enabled)
	assert.Equal(t, "stopped", svc2.Status().BindState)
	svc2.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	require.NoError(t, svc2.StartConfiguredListener(ctx))
	t.Cleanup(func() { _ = svc2.Stop() })

	status := svc2.Status()
	assert.True(t, status.Enabled)
	assert.Equal(t, port, status.Port)
	assert.Equal(t, "127.0.0.1", status.BindAddress)
	assert.Equal(t, uint8(23), status.SlaveID)
	assert.Equal(t, 64, status.CapacityRegisters)
	assert.Equal(t, "restart-revision", status.SettingsRevision)
}

func TestService_PersistedDisabledSettingsRestoreAsDisabledAcrossServiceRestart(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	repo := datalinksettings.NewSQLRepository(db)
	ctx := context.Background()
	persisted := Settings{
		Enabled:           false,
		BindAddress:       "127.0.0.1",
		Port:              15020,
		SlaveID:           23,
		CapacityRegisters: 64,
		SettingsRevision:  "disabled-restart-revision",
	}
	require.NoError(t, repo.Set(ctx, datalinksettings.KeyModbusShare, persisted))

	svc := NewService(tag.NewService(tag.NewMemoryRepository()), 128)
	require.NoError(t, svc.SetSettingsRepository(repo))
	require.NoError(t, svc.LoadPersistedSettings(ctx))

	status := svc.Status()
	assert.False(t, status.ConfiguredEnabled)
	assert.False(t, status.Enabled)
	assert.Equal(t, "disabled", status.BindState)
	assert.Equal(t, "disabled", status.LifecycleState)
}

func TestService_StartConfiguredListenerRequiresReadyHydration(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	repo := datalinksettings.NewSQLRepository(db)
	ctx := context.Background()
	port := reservePort(t)
	svc := NewService(tag.NewService(tag.NewMemoryRepository()), 128)
	require.NoError(t, svc.SetSettingsRepository(repo))
	require.NoError(t, repo.Set(ctx, datalinksettings.KeyModbusShare, Settings{Enabled: true, BindAddress: "127.0.0.1", Port: port, SlaveID: 1, CapacityRegisters: 64, SettingsRevision: "pending-start"}))
	require.NoError(t, svc.LoadPersistedSettings(ctx))
	err = svc.StartConfiguredListener(ctx)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeHydrationRequired, shareErr.Code)
	status := svc.Status()
	assert.True(t, status.ConfiguredEnabled)
	assert.False(t, status.Enabled)
	assert.Equal(t, "fail", status.BindState)
	assert.Equal(t, 0, status.Port)
}

func TestService_CloseRuntimePreservesDurableEnabledSetting(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	repo := datalinksettings.NewSQLRepository(db)
	ctx := context.Background()
	port := reservePort(t)
	persisted := Settings{Enabled: true, BindAddress: "127.0.0.1", Port: port, SlaveID: 1, CapacityRegisters: 64, SettingsRevision: "close-runtime"}
	svc := NewService(tag.NewService(tag.NewMemoryRepository()), 128)
	require.NoError(t, svc.SetSettingsRepository(repo))
	require.NoError(t, svc.ApplySettings(ctx, persisted))
	require.NoError(t, svc.CloseRuntime())

	restored := NewService(tag.NewService(tag.NewMemoryRepository()), 128)
	require.NoError(t, restored.SetSettingsRepository(repo))
	require.NoError(t, restored.LoadPersistedSettings(ctx))
	restored.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	require.NoError(t, restored.StartConfiguredListener(ctx))
	assert.True(t, restored.Settings().Enabled)
	assert.Equal(t, port, restored.Status().Port)
	require.NoError(t, restored.CloseRuntime())
}
