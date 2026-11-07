package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type activationSnapshotClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *activationSnapshotClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *activationSnapshotClock) Set(now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = now
}

func bugARequest(t *testing.T, f scheduledProductionFixture, method, path string, body any) (status int, payload map[string]any) {
	t.Helper()
	var data []byte
	var err error
	if body != nil {
		data, err = json.Marshal(body)
		require.NoError(t, err)
	}
	req := httptest.NewRequestWithContext(t.Context(), method, "/api/v1/datalink"+path, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	f.router.ServeHTTP(response, req)
	var parsed map[string]any
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &parsed), response.Body.String())
	return response.Code, parsed
}

func bugAActivationBody(t *testing.T, f scheduledProductionFixture) map[string]any {
	t.Helper()
	code, result := bugARequest(t, f, http.MethodGet, "/studio-v2/workspace", nil)
	require.Equal(t, http.StatusOK, code)
	data := result["data"].(map[string]any)
	share := data["modbus_share"].(map[string]any)
	return map[string]any{
		"workspace_revision": share["workspace_revision"],
		"settings_revision":  share["settings_revision"],
		"readiness_token":    share["readiness_token"],
	}
}

func TestProductionSnapshotWriter_ContinuesAfterFailedActivationAndSchemaPreparation(t *testing.T) {
	f, _, _ := newAcceptanceScheduledLoad(t, 1, 2, time.Millisecond, true)
	destinationPath := filepath.Join(t.TempDir(), "target.db")
	code, result := bugARequest(t, f, http.MethodPut, "/studio-v2/workspace/database-config", map[string]any{
		"kind": "sqlite", "name": "target", "database": destinationPath, "table": "readings", "write_mode": "append", "write_interval_seconds": 2,
	})
	require.Equal(t, http.StatusOK, code, result)
	bindings, err := f.services.mapping.List(t.Context(), mapping.ListFilter{})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(bindings), 2)
	for i, binding := range bindings[:2] {
		code, config := bugARequest(t, f, http.MethodGet, "/studio-v2/workspace/database-config", nil)
		require.Equal(t, http.StatusOK, code)
		revision := config["data"].(map[string]any)["setup_revision"]
		code, response := bugARequest(t, f, http.MethodPut, "/studio-v2/workspace/database-targets/"+binding.PointID, map[string]any{
			"column_name": []string{"first", "second"}[i], "enabled": true, "expected_setup_revision": revision,
		})
		require.Equal(t, http.StatusOK, code, response)
	}
	stored, err := f.services.dbMapping.List(t.Context(), dbtarget.TargetMappingListFilter{})
	require.NoError(t, err)
	require.Len(t, stored, 2)
	slices.SortFunc(stored, func(a, b *schema.DatabaseTargetMapping) int { return strings.Compare(a.ID, b.ID) })
	base := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	clock := &activationSnapshotClock{now: base}
	tick := make(chan time.Time, 1)
	// Use the production writer and ownership policy with a controlled tick;
	// the separate process-level reproducer exercises the scheduled ingestion.
	writer := dbtarget.NewWriterWithConfig(dbtarget.NewSQLConnectorRepository(f.db), dbtarget.NewSQLTargetMappingRepository(f.db), dbtarget.WriterConfig{
		FlushTick: tick, Now: clock.Now, TagReader: f.services.tag, SuppressMapping: f.services.groupPipe.Owns,
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	code, result = bugARequest(t, f, http.MethodPost, "/studio-v2/workspace/activate", bugAActivationBody(t, f))
	require.Equal(t, http.StatusUnprocessableEntity, code, result)
	require.Equal(t, "WORKSPACE_SCHEMA_PREPARATION_REQUIRED", result["error"].(map[string]any)["code"])
	code, result = bugARequest(t, f, http.MethodPost, "/studio-v2/workspace/database-schema/generate", map[string]any{"dry_run": true})
	require.Equal(t, http.StatusOK, code, result)
	destination, err := sql.Open("sqlite", destinationPath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = destination.Close() })
	for _, statement := range result["data"].(map[string]any)["statements"].([]any) {
		_, err = destination.ExecContext(t.Context(), statement.(string))
		require.NoError(t, err)
	}
	code, result = bugARequest(t, f, http.MethodPost, "/studio-v2/workspace/activate", bugAActivationBody(t, f))
	require.Equal(t, http.StatusOK, code, result)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), stored[0].TagID, 10, base.Add(time.Second)), dbtarget.ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), stored[1].TagID, 20, base.Add(time.Second)), dbtarget.ErrSnapshotAccepted)
	clock.Set(base.Add(2 * time.Second))
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), stored[0].TagID, 12, clock.Now()), dbtarget.ErrSnapshotAccepted)
	tick <- clock.Now()
	require.Eventually(t, func() bool {
		var count int
		return destination.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM readings WHERE "`+stored[1].ColumnName+`" IS NOT NULL`).Scan(&count) == nil && count == 1
	}, 3*time.Second, 10*time.Millisecond, "the worker must process the due sentinel")
	var count int
	require.NoError(t, destination.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM readings WHERE "`+stored[0].ColumnName+`" IS NOT NULL`).Scan(&count))
	require.Equal(t, 1, count, "new samples must not defer delivery forever after a failed activation")
	var value int
	require.NoError(t, destination.QueryRowContext(t.Context(), `SELECT "`+stored[0].ColumnName+`" FROM readings WHERE "`+stored[0].ColumnName+`" IS NOT NULL`).Scan(&value))
	require.Equal(t, 12, value, "the due snapshot must contain the latest value")
	require.Eventually(t, func() bool {
		history, historyErr := f.services.dbTarget.ListWriteHistory(t.Context(), stored[0].ConnectorID, 20)
		if historyErr != nil {
			return false
		}
		succeeded := 0
		for _, record := range history {
			if record.Status == "success" {
				succeeded++
			}
		}
		return succeeded >= 2
	}, 3*time.Second, 10*time.Millisecond, "successful writes must appear in connector history")
}
