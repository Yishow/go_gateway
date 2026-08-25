package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type statusRevisionStore struct {
	revision string
	dirty    bool
}

func (s *statusRevisionStore) GetRevision(context.Context, string) (revision string, dirty bool, err error) {
	return s.revision, s.dirty, nil
}

func (s *statusRevisionStore) UpdateRevision(context.Context, string, string, string) error {
	return nil
}

func (s *statusRevisionStore) MarkDirty(context.Context, string) error { return nil }

func TestModbusShareHandler_StatusScopesMappingsAndExposesLifecycleState(t *testing.T) {
	gate := &mockHydrationGate{
		hydrationState: modbusshare.HydrationState{
			State:             modbusshare.HydrationStateReady,
			WorkspaceID:       "ws-a",
			WorkspaceRevision: "rev-a",
			SettingsRevision:  "settings-a",
			Readiness:         true,
		},
		settings: modbusshare.Settings{Enabled: false, SettingsRevision: "settings-a"},
	}
	handler, svc, _ := setupGatedShareHandler(t, gate)
	svc.ReplaceMappings(map[string]modbusshare.TagMirrorMapping{
		"tag-a": {WorkspaceID: "ws-a", TagID: "tag-a"},
		"tag-b": {WorkspaceID: "ws-b", TagID: "tag-b"},
	})
	router := gin.New()
	router.GET("/status", handler.Status)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, newHandlerTestRequest(http.MethodGet, "/status", nil))
	require.Equal(t, http.StatusOK, response.Code)
	var envelope struct {
		Data struct {
			WorkspaceID    string `json:"workspace_id"`
			MappingCount   int    `json:"mapping_count"`
			ListenerState  string `json:"listener_state"`
			LifecycleState string `json:"lifecycle_state"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &envelope))
	require.Equal(t, "ws-a", envelope.Data.WorkspaceID)
	require.Equal(t, 1, envelope.Data.MappingCount)
	require.Equal(t, envelope.Data.LifecycleState, envelope.Data.ListenerState)
	require.Equal(t, "disabled", envelope.Data.ListenerState)
}

func TestModbusShareHandler_StatusShowsDurableDirtyRecoveryAfterRestart(t *testing.T) {
	gate := &mockHydrationGate{
		hydrationState: modbusshare.HydrationState{
			State:             modbusshare.HydrationStateReady,
			WorkspaceID:       "ws-a",
			WorkspaceRevision: "hydrated-revision",
			SettingsRevision:  "settings-a",
			Readiness:         true,
		},
		settings: modbusshare.Settings{Enabled: true, SettingsRevision: "settings-a"},
	}
	handler, _, _ := setupGatedShareHandler(t, gate)
	reconciler := modbusshare.NewReconciler(nil, &statusRevisionStore{revision: "durable-revision", dirty: true})
	handler.WithReconciler(reconciler)
	router := gin.New()
	router.GET("/status", handler.Status)

	request := newHandlerTestRequest(http.MethodGet, "/status", nil)
	request.Header.Set("X-Request-ID", "status-recovery-test")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	require.Equal(t, http.StatusOK, response.Code)
	var envelope struct {
		Data struct {
			WorkspaceRevision string `json:"workspace_revision"`
			Readiness         bool   `json:"readiness"`
			DirtyState        string `json:"dirty_state"`
			Recovery          struct {
				Code      string `json:"code"`
				Retryable bool   `json:"retryable"`
				Action    string `json:"action"`
				RequestID string `json:"request_id"`
			} `json:"recovery"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &envelope))
	require.Equal(t, "durable-revision", envelope.Data.WorkspaceRevision)
	require.False(t, envelope.Data.Readiness)
	require.Equal(t, "dirty_unknown", envelope.Data.DirtyState)
	require.Equal(t, modbusshare.ErrCodeDirtyUnknown, envelope.Data.Recovery.Code)
	require.True(t, envelope.Data.Recovery.Retryable)
	require.NotEmpty(t, envelope.Data.Recovery.Action)
	require.Equal(t, "status-recovery-test", envelope.Data.Recovery.RequestID)
}

func TestStudioV2WorkspaceBootstrapSharesDurableRevisionAndDirtyRecovery(t *testing.T) {
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository()).WithReadinessServices(
		device.NewService(device.NewMemoryRepository(), nil), nil, nil, nil,
	)
	record, err := workspaceSvc.GetOrCreate(context.Background())
	require.NoError(t, err)
	share := modbusshare.NewService(nil, 128)
	share.SetHydrationState(modbusshare.HydrationState{
		State:             modbusshare.HydrationStateReady,
		WorkspaceID:       record.ID,
		WorkspaceRevision: "hydrated-revision",
		SettingsRevision:  "settings-a",
		Readiness:         true,
	})
	handler := NewStudioV2WorkspaceHandler(workspaceSvc).
		WithModbusShare(share).
		WithModbusShareReconciler(modbusshare.NewReconciler(nil, &statusRevisionStore{revision: "durable-revision", dirty: true}))

	response := httptest.NewRecorder()
	contextValue, _ := gin.CreateTestContext(response)
	contextValue.Request = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/workspace", http.NoBody)
	handler.Get(contextValue)
	require.Equal(t, http.StatusOK, response.Code)
	var envelope struct {
		Data struct {
			ModbusShare struct {
				WorkspaceRevision string `json:"workspace_revision"`
				Readiness         bool   `json:"readiness"`
				Status            struct {
					WorkspaceRevision string `json:"workspace_revision"`
					DirtyState        string `json:"dirty_state"`
				} `json:"status"`
			} `json:"modbus_share"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &envelope))
	require.Equal(t, "durable-revision", envelope.Data.ModbusShare.WorkspaceRevision)
	require.Equal(t, envelope.Data.ModbusShare.WorkspaceRevision, envelope.Data.ModbusShare.Status.WorkspaceRevision)
	require.False(t, envelope.Data.ModbusShare.Readiness)
	require.Equal(t, "dirty_unknown", envelope.Data.ModbusShare.Status.DirtyState)
}
