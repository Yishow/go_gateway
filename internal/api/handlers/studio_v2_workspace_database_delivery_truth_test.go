package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStudioV2WorkspaceDatabaseHandler_GetConfigSurfacesDeliveryTruth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)

	ctx := context.Background()
	record, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	connector, err := fixture.connectorSvc.GetByID(ctx, record.DatabaseConnectorID)
	require.NoError(t, err)

	schemaAt := time.Date(2026, 3, 16, 10, 3, 0, 0, time.UTC)
	writeAt := time.Date(2026, 3, 16, 10, 5, 0, 0, time.UTC)
	flushAt := time.Date(2026, 3, 16, 10, 6, 0, 0, time.UTC)
	connector.LastSchemaEnsureAt = &schemaAt
	connector.LastSchemaEnsureStatus = "success"
	connector.LastWriteAt = &writeAt
	connector.LastWriteStatus = "failed"
	connector.LastWriteError = "permission denied"
	connector.LastFlushAt = &flushAt
	connector.LastFlushStatus = "failed"
	connector.LastFlushError = "flush timeout"
	require.NoError(t, fixture.connectorRepo.Update(ctx, connector))

	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-config", http.NoBody)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	fixture.handler.GetConfig(c)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	body := decodeWorkspaceDatabaseBody(t, resp)
	data := body["data"].(map[string]any)
	require.Equal(t, "success", data["last_schema_ensure_status"])
	require.Equal(t, schemaAt.Format(time.RFC3339Nano), data["last_schema_ensure_at"])
	require.Equal(t, "failed", data["last_write_status"])
	require.Equal(t, writeAt.Format(time.RFC3339Nano), data["last_write_at"])
	require.Equal(t, "permission denied", data["last_write_error"])
	require.Equal(t, "failed", data["last_flush_status"])
	require.Equal(t, flushAt.Format(time.RFC3339Nano), data["last_flush_at"])
	require.Equal(t, "flush timeout", data["last_flush_error"])
}
