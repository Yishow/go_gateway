package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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
