package handlers

import (
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

func TestModbusShareHandler_Diagnostics_HydrationRequiredEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gate := &mockHydrationGate{
		hydrationState: modbusshare.HydrationState{State: modbusshare.HydrationStatePending},
	}
	handler, _, _ := setupGatedShareHandler(t, gate)
	router := gin.New()
	router.GET("/mappings", handler.ListMappings)

	req := newHandlerTestRequest(http.MethodGet, "/mappings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp struct {
		Success bool `json:"success"`
		Error   struct {
			Code      string `json:"code"`
			Message   string `json:"message"`
			Retryable bool   `json:"retryable"`
			Action    string `json:"action"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.False(t, resp.Success)
	assert.Equal(t, modbusshare.ErrCodeHydrationRequired, resp.Error.Code)
	assert.True(t, resp.Error.Retryable)
	assert.NotEmpty(t, resp.Error.Action)
}

func TestModbusShareHandler_Diagnostics_RangeCollisionEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gate := &mockHydrationGate{
		hydrationState: modbusshare.HydrationState{State: modbusshare.HydrationStateReady, Readiness: true},
		settings:       modbusshare.Settings{Enabled: true, Port: 5020},
		ownershipProof: func(ctx context.Context, wsID, tagID string) bool { return true },
	}
	handler, svc, tagSvc := setupGatedShareHandler(t, gate)

	ctx := context.Background()
	t1, _ := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "t1", DisplayName: "Tag 1", DataType: schema.DataTypeInt32})
	t2, _ := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "t2", DisplayName: "Tag 2", DataType: schema.DataTypeInt32})

	_, err := svc.UpsertMapping(ctx, t1.ID, 0) // occupies [0, 2)
	require.NoError(t, err)

	router := gin.New()
	router.POST("/mappings/:tagId", handler.UpsertMapping)

	// Attempting to map t2 at register 1 (collides with t1)
	req := newHandlerTestRequest(http.MethodPost, "/mappings/"+t2.ID+"?register=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp struct {
		Success bool `json:"success"`
		Error   struct {
			Code      string `json:"code"`
			Message   string `json:"message"`
			Retryable bool   `json:"retryable"`
			Action    string `json:"action"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.False(t, resp.Success)
	assert.Equal(t, modbusshare.ErrCodeRangeCollision, resp.Error.Code)
	assert.NotEmpty(t, resp.Error.Action)
}

func TestModbusShareHandler_Diagnostics_CapacityExceededEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gate := &mockHydrationGate{
		hydrationState: modbusshare.HydrationState{State: modbusshare.HydrationStateReady, Readiness: true},
		settings:       modbusshare.Settings{Enabled: true, Port: 5020},
		ownershipProof: func(ctx context.Context, wsID, tagID string) bool { return true },
	}
	handler, _, tagSvc := setupGatedShareHandler(t, gate)

	ctx := context.Background()
	t1, _ := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "t1", DisplayName: "Tag 1", DataType: schema.DataTypeInt64})

	router := gin.New()
	router.POST("/mappings/:tagId", handler.UpsertMapping)

	// Attempting to map int64 (span=4) at 65535 (overflows 65536)
	req := newHandlerTestRequest(http.MethodPost, "/mappings/"+t1.ID+"?register=65535", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp struct {
		Success bool `json:"success"`
		Error   struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, modbusshare.ErrCodeCapacityExceeded, resp.Error.Code)
}
