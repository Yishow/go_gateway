package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func homeSourcesRequest(h *StudioV2RecordingHomeSourcesHandler, query string) *httptest.ResponseRecorder {
	out := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(out)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/recording-home/sources"+query, http.NoBody)
	h.Get(c)
	return out
}
func TestRecordingHomeSourcesUnconfiguredRealServices(t *testing.T) {
	old, _, _, _ := newWorkspaceMappingHandler(t)
	h := NewStudioV2RecordingHomeSourcesHandler(old.workspaceSvc, old.deviceSvc, old.ruleSvc, old.pointSvc, old.mappingSvc, old.tagSvc)
	before, err := old.workspaceSvc.ReadSnapshot(context.Background())
	require.NoError(t, err)
	linksBefore, err := old.ruleSvc.ListLinks(context.Background(), "rule-A")
	require.NoError(t, err)
	out := homeSourcesRequest(h, "?workspace_id="+before.ID)
	require.Equal(t, http.StatusOK, out.Code, out.Body.String())
	var body struct {
		Success bool `json:"success"`
		Data    struct {
			State      string `json:"state"`
			Candidates []any  `json:"candidates"`
			Devices    []any  `json:"devices"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(out.Body.Bytes(), &body))
	require.True(t, body.Success)
	require.Equal(t, "available", body.Data.State)
	require.Empty(t, body.Data.Candidates)
	require.Len(t, body.Data.Devices, 2)
	after, err := old.workspaceSvc.ReadSnapshot(context.Background())
	require.NoError(t, err)
	require.Equal(t, before, after)
	linksAfter, err := old.ruleSvc.ListLinks(context.Background(), "rule-A")
	require.NoError(t, err)
	require.Equal(t, linksBefore, linksAfter)
}
