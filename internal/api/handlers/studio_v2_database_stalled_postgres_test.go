package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStudioDatabaseReadsReportStalledPostgresWithoutHanging(t *testing.T) {
	gin.SetMode(gin.TestMode)
	listener, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })
	serverContext := t.Context()
	go func() {
		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}
			go func() {
				stop := context.AfterFunc(serverContext, func() { _ = conn.Close() })
				defer stop()
				defer conn.Close()
				_, _ = io.Copy(io.Discard, conn)
			}()
		}
	}()

	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	connector, err := fixture.connectorSvc.GetByID(t.Context(), record.DatabaseConnectorID)
	require.NoError(t, err)
	config, err := json.Marshal(map[string]any{
		"host": "127.0.0.1", "port": listener.Addr().(*net.TCPAddr).Port, "user": "fixture",
		"password": "fixture-secret", "database": "fixture", "schema": "public", "table": "sensor_values", "sslmode": "disable",
	})
	require.NoError(t, err)
	connector.Kind = schema.DatabaseConnectorKindPostgres
	connector.ConnectionConfig = string(config)
	require.NoError(t, fixture.connectorRepo.Update(t.Context(), connector))

	common := url.Values{
		"expected_connector_revision": {connector.IdentityRevision},
		"locator_workspace_id":        {record.ID},
		"locator_workspace_revision":  {record.DatabaseSetupRevision},
	}
	catalogQuery := url.Values{}
	for key, values := range common {
		catalogQuery[key] = values
	}
	catalogQuery.Set("workspace_id", record.ID)
	catalogQuery.Set("expected_workspace_revision", record.DatabaseSetupRevision)
	catalogQuery.Set("connector_id", connector.ID)
	requests := []struct {
		name   string
		path   string
		handle func(*gin.Context)
	}{
		{"catalog", "/database-catalog?" + catalogQuery.Encode(), fixture.handler.GetCatalog},
		{"metadata", "/database-metadata?" + common.Encode(), fixture.handler.GetMetadata},
	}
	type response struct {
		name string
		code int
		body []byte
	}
	responses := make(chan response, len(requests))
	for _, request := range requests {
		go func() {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequestWithContext(t.Context(), http.MethodGet, request.path, http.NoBody)
			request.handle(ctx)
			responses <- response{name: request.name, code: recorder.Code, body: recorder.Body.Bytes()}
		}()
	}
	got := make(map[string]response, len(requests))
	deadline := time.NewTimer(7 * time.Second)
	defer deadline.Stop()
	for range requests {
		select {
		case result := <-responses:
			got[result.name] = result
		case <-deadline.C:
			t.Fatal("database read remained blocked after the connection timeout")
		}
	}
	require.Equal(t, http.StatusServiceUnavailable, got["catalog"].code)
	require.NotContains(t, string(got["catalog"].body), "fixture-secret")
	require.Equal(t, http.StatusOK, got["metadata"].code)
	var metadata workspaceDatabaseMetadataBody
	require.NoError(t, json.Unmarshal(got["metadata"].body, &metadata))
	require.NotNil(t, metadata.Data)
	require.Equal(t, "failed", metadata.Data.InspectionStatus)
	require.Equal(t, "connection_failed", metadata.Data.Reason)
	require.Empty(t, metadata.Data.Columns)
}
