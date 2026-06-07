package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type failingWorkspaceRepository struct{}

func (failingWorkspaceRepository) Get(context.Context) (*workspace.Record, error) {
	return nil, workspace.ErrNotFound
}

func (failingWorkspaceRepository) Save(context.Context, *workspace.Record) error {
	return errors.New("write failed")
}

func TestStudioV2WorkspaceHandler_ReturnsActionableBootstrapError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewStudioV2WorkspaceHandler(workspace.NewService(failingWorkspaceRepository{}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Get(c)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d body=%s", resp.Code, resp.Body.String())
	}
	if body := resp.Body.String(); body == "" {
		t.Fatal("expected error response body")
	}
}

func TestStudioV2WorkspaceHandler_ReturnsReadinessSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)

	deviceSvc := device.NewService(device.NewMemoryRepository(), nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository()).WithReadinessServices(deviceSvc, nil, nil, nil)
	handler := NewStudioV2WorkspaceHandler(workspaceSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Get(c)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}
	if body := resp.Body.String(); body == "" || !strings.Contains(body, "readiness_summary") || !strings.Contains(body, "blocking_count") || !strings.Contains(body, "warning_count") {
		t.Fatalf("expected readiness summary payload, got %s", body)
	}
	if body := resp.Body.String(); !containsAll(body, `"ready":false`, `"blocking_count":1`, `"warning_count":0`, "workspace-device-missing") {
		t.Fatalf("expected empty workspace readiness blocker, got %s", body)
	}
}

func TestStudioV2WorkspaceHandler_ReturnsUnavailableReadinessSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewStudioV2WorkspaceHandler(workspace.NewService(workspace.NewMemoryRepository()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Get(c)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}
	if body := resp.Body.String(); !containsAll(body, `"ready":false`, `"blocking_count":1`, "readiness-unavailable", "workspace") {
		t.Fatalf("expected unavailable readiness summary, got %s", body)
	}
}
