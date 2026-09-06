package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/history"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

func newStudioV2WorkspaceHistoryRouter(t *testing.T) (*gin.Engine, *workspace.Service, *history.Service, *history.MemoryHistoryRepository) {
	t.Helper()

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	historyRepo := history.NewMemoryHistoryRepository()
	historySvc := history.NewService(historyRepo)

	router := NewRouter(&DatalinkServices{
		Workspace: workspaceSvc,
		History:   historySvc,
	})

	return router, workspaceSvc, historySvc, historyRepo
}

func TestStudioV2WorkspaceHistory_QueryAndExport(t *testing.T) {
	router, workspaceSvc, _, historyRepo := newStudioV2WorkspaceHistoryRouter(t)

	wsRecord, err := workspaceSvc.GetOrCreate(context.Background())
	if err != nil {
		t.Fatalf("failed to get workspace: %v", err)
	}

	val := 125.5
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	_ = historyRepo.SavePoints(wsRecord.ID, "plan-1", []history.HistoryPoint{
		{
			MeasurementID: "meas-kw",
			ObservedAt:    now,
			ValueNumeric:  &val,
			Quality:       "good",
			CoverageRatio: 1.0,
		},
	})

	// 1. Query History
	queryBody := `{
		"plan_id": "plan-1",
		"measurement_ids": ["meas-kw"]
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/history/query", strings.NewReader(queryBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("query history failed with status %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "125.5") {
		t.Errorf("expected queried data to contain 125.5, got %s", w.Body.String())
	}

	// 2. Export CSV
	req = httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/history/export", strings.NewReader(queryBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("export history failed with status %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "observed_at,measurement_id") || !strings.Contains(w.Body.String(), "125.5") {
		t.Errorf("expected CSV export content, got %s", w.Body.String())
	}
}
