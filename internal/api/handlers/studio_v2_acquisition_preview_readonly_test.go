package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type previewMissingWorkspace struct{ saves int }

func (r *previewMissingWorkspace) Get(context.Context) (*workspace.Record, error) {
	return nil, workspace.ErrNotFound
}
func (r *previewMissingWorkspace) Save(context.Context, *workspace.Record) error {
	r.saves++
	return nil
}

func TestAcquisitionPreviewMissingWorkspaceNeverBootstraps(t *testing.T) {
	repo := &previewMissingWorkspace{}
	router := gin.New()
	router.GET("/preview", NewStudioV2AcquisitionPreviewHandler(workspace.NewService(repo), nil).Get)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, newHandlerTestRequest(http.MethodGet, "/preview?workspace_id=w&device_id=d&expected_workspace_revision=", http.NoBody))
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Zero(t, repo.saves)
	require.NotContains(t, response.Body.String(), "groups")
}
