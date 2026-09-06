package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type hydrationBlockedWorkspaceActivator struct {
	called bool
}

func (a *hydrationBlockedWorkspaceActivator) ActivateEligible(context.Context) (*workspace.ActivationResponse, error) {
	a.called = true
	return &workspace.ActivationResponse{}, nil
}

func TestConfigureModbusShareActivation_BlocksFailedHydration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	share := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
	share.SetHydrationState(modbusshare.HydrationState{
		State:             modbusshare.HydrationStateFailed,
		WorkspaceID:       "ws-1",
		WorkspaceRevision: "rev-1",
		SettingsRevision:  "set-1",
		Readiness:         false,
	})
	settings := modbusshare.DefaultSettings()
	settings.Enabled = true
	settings.Port = 5020
	settings.SettingsRevision = "settings-1"
	require.NoError(t, share.ApplySettings(context.Background(), settings))
	activator := &hydrationBlockedWorkspaceActivator{}
	handler := handlers.NewStudioV2WorkspaceActivationHandler(activator)
	configureModbusShareActivation(handler, &DatalinkServices{ModbusShare: share})

	router := gin.New()
	router.POST("/activate", handler.Activate)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/activate", bytes.NewBufferString(`{"workspace_revision":"rev-1","settings_revision":"settings-1","readiness_token":"token-1"}`))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	require.Equal(t, http.StatusUnprocessableEntity, response.Code)
	require.Contains(t, response.Body.String(), modbusshare.ErrCodeHydrationRequired)
	require.False(t, activator.called)
}

func TestConfigureModbusShareActivation_DisabledShareSkipsRestoreAndActivatesWorkspace(t *testing.T) {
	gin.SetMode(gin.TestMode)
	share := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
	share.SetHydrationState(modbusshare.HydrationState{
		State:     modbusshare.HydrationStateFailed,
		Readiness: false,
	})
	activator := &hydrationBlockedWorkspaceActivator{}
	restoreCalled := false
	handler := handlers.NewStudioV2WorkspaceActivationHandler(activator)
	configureModbusShareActivation(handler, &DatalinkServices{
		ModbusShare: share,
		ShareRestore: func(context.Context, handlers.ActivateWorkspaceRequest) error {
			restoreCalled = true
			return &modbusshare.Error{Code: modbusshare.ErrCodeDisabled}
		},
	})

	router := gin.New()
	router.POST("/activate", handler.Activate)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/activate", bytes.NewBufferString(`{"workspace_revision":"rev-1","settings_revision":"default","readiness_token":"token-1"}`))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	require.Equal(t, http.StatusOK, response.Code)
	require.True(t, activator.called)
	require.False(t, restoreCalled)
}

// Share 停用時仍必須保留工作區層級的前置檢查，只略過 Share 專屬的 settings revision。
func TestConfigureModbusShareActivation_DisabledShareStillEnforcesReadinessToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	share := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
	share.SetHydrationState(modbusshare.HydrationState{
		State:             modbusshare.HydrationStateReady,
		WorkspaceID:       "ws-1",
		WorkspaceRevision: "rev-1",
		SettingsRevision:  "set-1",
		Readiness:         true,
		ReadinessToken:    "token-1",
	})
	activator := &hydrationBlockedWorkspaceActivator{}
	handler := handlers.NewStudioV2WorkspaceActivationHandler(activator)
	configureModbusShareActivation(handler, &DatalinkServices{ModbusShare: share})

	router := gin.New()
	router.POST("/activate", handler.Activate)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/activate", bytes.NewBufferString(`{"workspace_revision":"rev-1","settings_revision":"ignored","readiness_token":"stale-token"}`))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	require.Contains(t, response.Body.String(), modbusshare.ErrCodeSaveIncomplete)
	require.False(t, activator.called)
}

func TestConfigureModbusShareActivation_DisabledShareStillEnforcesWorkspaceRevision(t *testing.T) {
	gin.SetMode(gin.TestMode)
	share := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
	share.SetHydrationState(modbusshare.HydrationState{
		State:             modbusshare.HydrationStateReady,
		WorkspaceID:       "ws-1",
		WorkspaceRevision: "rev-1",
		SettingsRevision:  "set-1",
		Readiness:         true,
		ReadinessToken:    "token-1",
	})
	activator := &hydrationBlockedWorkspaceActivator{}
	handler := handlers.NewStudioV2WorkspaceActivationHandler(activator)
	configureModbusShareActivation(handler, &DatalinkServices{ModbusShare: share})

	router := gin.New()
	router.POST("/activate", handler.Activate)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/activate", bytes.NewBufferString(`{"workspace_revision":"stale-rev","settings_revision":"ignored","readiness_token":"token-1"}`))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	require.Contains(t, response.Body.String(), modbusshare.ErrCodeRevisionConflict)
	require.False(t, activator.called)
}

// Share 停用且 hydration 已就緒時，Share 專屬的 settings revision 檢查必須被略過。
func TestConfigureModbusShareActivation_DisabledShareSkipsSettingsRevisionGate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	share := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
	share.SetHydrationState(modbusshare.HydrationState{
		State:             modbusshare.HydrationStateReady,
		WorkspaceID:       "ws-1",
		WorkspaceRevision: "rev-1",
		SettingsRevision:  "set-1",
		Readiness:         true,
		ReadinessToken:    "token-1",
	})
	activator := &hydrationBlockedWorkspaceActivator{}
	handler := handlers.NewStudioV2WorkspaceActivationHandler(activator)
	configureModbusShareActivation(handler, &DatalinkServices{ModbusShare: share})

	router := gin.New()
	router.POST("/activate", handler.Activate)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/activate", bytes.NewBufferString(`{"workspace_revision":"rev-1","settings_revision":"whatever","readiness_token":"token-1"}`))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	require.Equal(t, http.StatusOK, response.Code)
	require.True(t, activator.called)
}

type failingSettingsShareGate struct {
	hydration modbusshare.HydrationState
}

func (g failingSettingsShareGate) GetSettings(context.Context) (modbusshare.Settings, error) {
	return modbusshare.Settings{}, errors.New("settings storage unavailable")
}

func (g failingSettingsShareGate) CheckHydration(context.Context) (modbusshare.HydrationState, error) {
	return g.hydration, nil
}

// 讀取 Share 設定失敗屬於內部前置檢查失敗，不得誤報為「Share 已停用」。
func TestValidateStudioV2ActivationBarrier_SettingsReadFailureIsNotReportedAsDisabled(t *testing.T) {
	gate := failingSettingsShareGate{hydration: modbusshare.HydrationState{
		State:             modbusshare.HydrationStateReady,
		WorkspaceRevision: "rev-1",
		SettingsRevision:  "set-1",
		Readiness:         true,
		ReadinessToken:    "token-1",
	}}

	err := validateStudioV2ActivationBarrier(
		context.Background(),
		gate,
		handlers.ActivateWorkspaceRequest{WorkspaceRevision: "rev-1", SettingsRevision: "set-1", ReadinessToken: "token-1"},
	)

	require.Error(t, err)
	var shareErr *modbusshare.Error
	require.ErrorAs(t, err, &shareErr)
	require.NotEqual(t, modbusshare.ErrCodeDisabled, shareErr.Code)
	require.True(t, shareErr.Retryable)
	require.Empty(t, shareErr.SettingsRevision)
}

type failingHydrationShareGate struct {
	settings modbusshare.Settings
}

func (g failingHydrationShareGate) GetSettings(context.Context) (modbusshare.Settings, error) {
	return g.settings, nil
}

func (g failingHydrationShareGate) CheckHydration(context.Context) (modbusshare.HydrationState, error) {
	return modbusshare.HydrationState{}, errors.New("hydration storage unavailable")
}

// hydration 讀取失敗不得被當成「尚未 bootstrap」而放行，否則工作區層級的
// barrier 會在 Share 停用時被一併關掉。
func TestValidateStudioV2ActivationBarrier_HydrationReadFailureBlocksActivation(t *testing.T) {
	for _, shareEnabled := range []bool{false, true} {
		settings := modbusshare.DefaultSettings()
		settings.Enabled = shareEnabled

		err := validateStudioV2ActivationBarrier(
			context.Background(),
			failingHydrationShareGate{settings: settings},
			handlers.ActivateWorkspaceRequest{WorkspaceRevision: "rev-1", SettingsRevision: "set-1", ReadinessToken: "token-1"},
		)

		require.Error(t, err, "share enabled=%v", shareEnabled)
		var shareErr *modbusshare.Error
		require.ErrorAs(t, err, &shareErr)
		require.True(t, shareErr.Retryable)
	}
}

type stubWorkspaceRevisionStore struct {
	revision string
	err      error
}

func (s stubWorkspaceRevisionStore) GetRevision(context.Context, string) (revision string, isDirty bool, err error) {
	return s.revision, false, s.err
}

func (s stubWorkspaceRevisionStore) UpdateRevision(context.Context, string, string, string) error {
	return nil
}

func (s stubWorkspaceRevisionStore) MarkDirty(context.Context, string) error { return nil }

func (s stubWorkspaceRevisionStore) GetDesiredMappings(context.Context, string) ([]modbusshare.DesiredMapping, error) {
	return nil, nil
}

func (s stubWorkspaceRevisionStore) CommitDesiredMappings(context.Context, string, string, string, []modbusshare.DesiredMapping) error {
	return nil
}

// 工作區 revision 檢查必須獨立於 Share 設定與 hydration 之外成立。
func TestWorkspaceRevisionValidator_RejectsStaleRevisionRegardlessOfShareState(t *testing.T) {
	validator := newWorkspaceRevisionValidator(
		func(context.Context) (string, error) { return "ws-1", nil },
		stubWorkspaceRevisionStore{revision: "rev-2"},
	)

	require.NoError(t, validator(context.Background(), "rev-2", "ignored"))

	err := validator(context.Background(), "rev-1", "ignored")
	require.Error(t, err)
	var shareErr *modbusshare.Error
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, modbusshare.ErrCodeRevisionConflict, shareErr.Code)

	err = validator(context.Background(), "", "ignored")
	require.Error(t, err, "缺少 workspace_revision 同樣必須被擋下")
}

func TestWorkspaceRevisionValidator_AllowsWorkspaceWithoutPersistedRevision(t *testing.T) {
	validator := newWorkspaceRevisionValidator(
		func(context.Context) (string, error) { return "ws-1", nil },
		stubWorkspaceRevisionStore{revision: ""},
	)

	require.NoError(t, validator(context.Background(), "", "ignored"))
}

func TestWorkspaceRevisionValidator_BlocksOnReadFailure(t *testing.T) {
	validator := newWorkspaceRevisionValidator(
		func(context.Context) (string, error) { return "ws-1", nil },
		stubWorkspaceRevisionStore{err: errors.New("revision storage unavailable")},
	)

	err := validator(context.Background(), "rev-1", "ignored")
	require.Error(t, err)
	var shareErr *modbusshare.Error
	require.ErrorAs(t, err, &shareErr)
	require.True(t, shareErr.Retryable)
}
