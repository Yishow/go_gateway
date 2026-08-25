package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHydrationGate struct {
	hydrationState modbusshare.HydrationState
	settings       modbusshare.Settings
	ownershipProof func(ctx context.Context, workspaceID, tagID string) bool
}

func (m *mockHydrationGate) CheckHydration(ctx context.Context) (modbusshare.HydrationState, error) {
	return m.hydrationState, nil
}

func (m *mockHydrationGate) GetSettings(ctx context.Context) (modbusshare.Settings, error) {
	return m.settings, nil
}

func (m *mockHydrationGate) ValidateOwnership(ctx context.Context, workspaceID, tagID string) bool {
	if m.ownershipProof != nil {
		return m.ownershipProof(ctx, workspaceID, tagID)
	}
	return false
}

func setupGatedShareHandler(t *testing.T, gate *mockHydrationGate) (handler *ModbusShareHandler, shareService *modbusshare.Service, tagService *tag.Service) {
	t.Helper()
	tagRepo := tag.NewMemoryRepository()
	tagSvc := tag.NewService(tagRepo)
	svc := modbusshare.NewService(tagSvc, 65536)
	t.Cleanup(func() {
		_ = svc.Stop()
	})

	handler = NewModbusShareHandler(svc)
	if gate != nil {
		handler.WithGate(gate)
	}
	return handler, svc, tagSvc
}

func TestModbusShareHandler_BootstrapPendingOrFailedBlocksShareOperations(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gate := &mockHydrationGate{
		hydrationState: modbusshare.HydrationState{
			State:     "pending",
			Readiness: false,
		},
		settings: modbusshare.Settings{Enabled: true, Port: 5020},
	}

	handler, svc, _ := setupGatedShareHandler(t, gate)
	router := gin.New()
	router.GET("/status", handler.Status)
	router.GET("/mappings", handler.ListMappings)
	router.PUT("/mappings/:tagId", handler.UpsertMapping)
	router.DELETE("/mappings/:tagId", handler.DeleteMapping)
	router.POST("/write-tag-value", handler.WriteTagValue)
	router.POST("/sync", handler.SyncFromMappings)

	statusReq := newHandlerTestRequest(http.MethodGet, "/status", nil)
	statusW := httptest.NewRecorder()
	router.ServeHTTP(statusW, statusReq)
	require.Equal(t, http.StatusOK, statusW.Code)
	assert.Contains(t, statusW.Body.String(), `"hydration_state":"pending"`)

	listReq := newHandlerTestRequest(http.MethodGet, "/mappings", nil)
	listW := httptest.NewRecorder()
	router.ServeHTTP(listW, listReq)
	require.Equal(t, http.StatusUnprocessableEntity, listW.Code)
	assert.Contains(t, listW.Body.String(), modbusshare.ErrCodeHydrationRequired)

	reg := uint16(0)
	upsertBody, _ := json.Marshal(UpsertMirrorMappingRequest{Register: &reg})
	upsertReq := newHandlerTestRequest(http.MethodPut, "/mappings/tag-1", bytes.NewReader(upsertBody))
	upsertReq.Header.Set("Content-Type", "application/json")
	upsertW := httptest.NewRecorder()
	router.ServeHTTP(upsertW, upsertReq)
	require.Equal(t, http.StatusUnprocessableEntity, upsertW.Code)
	assert.Contains(t, upsertW.Body.String(), modbusshare.ErrCodeHydrationRequired)
	assert.False(t, svc.HasMapping("tag-1"))

	deleteReq := newHandlerTestRequest(http.MethodDelete, "/mappings/tag-1", nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)
	require.Equal(t, http.StatusUnprocessableEntity, deleteW.Code)
	assert.Contains(t, deleteW.Body.String(), modbusshare.ErrCodeHydrationRequired)

	gate.hydrationState.State = "failed"
	failReq := newHandlerTestRequest(http.MethodGet, "/mappings", nil)
	failW := httptest.NewRecorder()
	router.ServeHTTP(failW, failReq)
	require.Equal(t, http.StatusUnprocessableEntity, failW.Code)
	assert.Contains(t, failW.Body.String(), modbusshare.ErrCodeHydrationRequired)
}

func TestModbusShareHandler_GatesRequireReadinessFlagWithReadyState(t *testing.T) {
	gate := &mockHydrationGate{
		hydrationState: modbusshare.HydrationState{State: modbusshare.HydrationStateReady, Readiness: false},
		settings:       modbusshare.Settings{Enabled: true},
	}
	handler, _, _ := setupGatedShareHandler(t, gate)
	router := gin.New()
	router.GET("/mappings", handler.ListMappings)
	req := newHandlerTestRequest(http.MethodGet, "/mappings", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	require.Equal(t, http.StatusUnprocessableEntity, response.Code)
	assert.Contains(t, response.Body.String(), modbusshare.ErrCodeHydrationRequired)
}

func TestModbusShareHandler_GlobalDisabledPreventsMutationAndHidesTruth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gate := &mockHydrationGate{
		hydrationState: modbusshare.HydrationState{
			State:             "ready",
			WorkspaceID:       "ws-1",
			WorkspaceRevision: "rev-1",
			SettingsRevision:  "set-1",
			Readiness:         true,
		},
		settings: modbusshare.Settings{Enabled: false, Port: 5020},
	}

	handler, svc, _ := setupGatedShareHandler(t, gate)
	router := gin.New()
	router.GET("/status", handler.Status)
	router.GET("/mappings", handler.ListMappings)
	router.PUT("/mappings/:tagId", handler.UpsertMapping)
	router.DELETE("/mappings/:tagId", handler.DeleteMapping)
	router.POST("/start", handler.Start)

	statusReq := newHandlerTestRequest(http.MethodGet, "/status", nil)
	statusW := httptest.NewRecorder()
	router.ServeHTTP(statusW, statusReq)
	require.Equal(t, http.StatusOK, statusW.Code)
	assert.Contains(t, statusW.Body.String(), `"enabled":false`)

	listReq := newHandlerTestRequest(http.MethodGet, "/mappings", nil)
	listW := httptest.NewRecorder()
	router.ServeHTTP(listW, listReq)
	require.Equal(t, http.StatusUnprocessableEntity, listW.Code)
	assert.Contains(t, listW.Body.String(), modbusshare.ErrCodeDisabled)

	reg := uint16(10)
	upsertBody, _ := json.Marshal(UpsertMirrorMappingRequest{Register: &reg})
	upsertReq := newHandlerTestRequest(http.MethodPut, "/mappings/tag-1", bytes.NewReader(upsertBody))
	upsertReq.Header.Set("Content-Type", "application/json")
	upsertW := httptest.NewRecorder()
	router.ServeHTTP(upsertW, upsertReq)
	require.Equal(t, http.StatusUnprocessableEntity, upsertW.Code)
	assert.Contains(t, upsertW.Body.String(), modbusshare.ErrCodeDisabled)
	assert.False(t, svc.HasMapping("tag-1"))

	deleteReq := newHandlerTestRequest(http.MethodDelete, "/mappings/tag-1", nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)
	require.Equal(t, http.StatusUnprocessableEntity, deleteW.Code)
	assert.Contains(t, deleteW.Body.String(), modbusshare.ErrCodeDisabled)

	// Start rejected with modbus_share_disabled
	startReq := newHandlerTestRequest(http.MethodPost, "/start", bytes.NewReader([]byte("{}")))
	startReq.Header.Set("Content-Type", "application/json")
	startW := httptest.NewRecorder()
	router.ServeHTTP(startW, startReq)
	require.Equal(t, http.StatusUnprocessableEntity, startW.Code)
	assert.Contains(t, startW.Body.String(), modbusshare.ErrCodeDisabled)
}

func TestModbusShareHandler_UnpersistedOrBrowserOnlyTagRejectsDeleteWithWorkspaceScope(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gate := &mockHydrationGate{
		hydrationState: modbusshare.HydrationState{
			State:             "ready",
			WorkspaceID:       "ws-1",
			WorkspaceRevision: "rev-1",
			SettingsRevision:  "set-1",
			Readiness:         true,
		},
		settings: modbusshare.Settings{Enabled: true, Port: 5020},
		ownershipProof: func(ctx context.Context, workspaceID, tagID string) bool {
			// Only tag-owned is proven
			return workspaceID == "ws-1" && tagID == "tag-owned"
		},
	}

	handler, svc, tagSvc := setupGatedShareHandler(t, gate)
	ctx := context.Background()
	tRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "t_unowned", DisplayName: "Unowned Tag", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	_, err = svc.UpsertMapping(ctx, tRecord.ID, 0)
	require.NoError(t, err)

	router := gin.New()
	router.DELETE("/mappings/:tagId", handler.DeleteMapping)

	// Deleting tRecord.ID must be rejected with modbus_share_workspace_scope
	req := newHandlerTestRequest(http.MethodDelete, "/mappings/"+tRecord.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), modbusshare.ErrCodeWorkspaceScope)
	// The mapping must be retained in memory!
	assert.True(t, svc.HasMapping(tRecord.ID))
}

func TestModbusShareHandler_CrossWorkspaceOperationRejection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gate := &mockHydrationGate{
		hydrationState: modbusshare.HydrationState{
			State:             "ready",
			WorkspaceID:       "ws-a",
			WorkspaceRevision: "rev-a",
			SettingsRevision:  "set-1",
			Readiness:         true,
		},
		settings: modbusshare.Settings{Enabled: true, Port: 5020},
		ownershipProof: func(ctx context.Context, workspaceID, tagID string) bool {
			// tag-ws-b belongs to ws-b only
			return workspaceID == "ws-b" && tagID == "tag-ws-b"
		},
	}

	handler, _, _ := setupGatedShareHandler(t, gate)
	router := gin.New()
	router.PUT("/mappings/:tagId", handler.UpsertMapping)

	reg := uint16(5)
	body, _ := json.Marshal(UpsertMirrorMappingRequest{Register: &reg})
	req := newHandlerTestRequest(http.MethodPut, "/mappings/tag-ws-b", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), modbusshare.ErrCodeWorkspaceScope)
}
