package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const recordingHomeSourcesPath = "/api/v1/datalink/studio-v2/recording-home/sources"

func homeRouterServices(t *testing.T) *DatalinkServices {
	t.Helper()
	wr := workspace.NewMemoryRepository()
	require.NoError(t, wr.Save(t.Context(), &workspace.Record{ID: "W1", OrderedDeviceIDs: []string{}}))
	d := device.NewService(device.NewMemoryRepository(), nil)
	p := point.NewService(point.NewMemoryRepository(), nil)
	return &DatalinkServices{Workspace: workspace.NewService(wr), Device: d, Point: p,
		SourceRule: sourcerule.NewService(sourcerule.NewMemoryRepository(), d, p, nil),
		Mapping:    mapping.NewService(mapping.NewMemoryRepository()), Tag: tag.NewService(tag.NewMemoryRepository())}
}

func TestRecordingHomeSourcesRouterMissingDependencies(t *testing.T) {
	for _, missing := range []string{"all", "workspace", "device", "rules", "points", "mappings", "tags"} {
		t.Run(missing, func(t *testing.T) {
			s := homeRouterServices(t)
			switch missing {
			case "all":
				s = &DatalinkServices{}
			case "workspace":
				s.Workspace = nil
			case "device":
				s.Device = nil
			case "rules":
				s.SourceRule = nil
			case "points":
				s.Point = nil
			case "mappings":
				s.Mapping = nil
			case "tags":
				s.Tag = nil
			}
			out := performJSONRequest(t, NewRouter(s), http.MethodGet, recordingHomeSourcesPath+"?workspace_id=W1", nil)
			require.Equal(t, http.StatusServiceUnavailable, out.Code, out.Body.String())
			require.JSONEq(t, `{"success":false,"code":"source-unavailable"}`, out.Body.String())
		})
	}
}
func TestRecordingHomeSourcesRouterScopeAndEmptyWorkspace(t *testing.T) {
	s := homeRouterServices(t)
	router := NewRouter(s)
	before, err := s.Workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	for _, query := range []string{"", "?workspace_id=", "?workspace_id=W1&workspace_id=W1", "?workspace_id=W1&extra=1"} {
		out := performJSONRequest(t, router, http.MethodGet, recordingHomeSourcesPath+query, nil)
		require.Equal(t, http.StatusBadRequest, out.Code, out.Body.String())
	}
	missing := performJSONRequest(t, router, http.MethodGet, recordingHomeSourcesPath+"?workspace_id=W2", nil)
	require.Equal(t, http.StatusNotFound, missing.Code, missing.Body.String())
	for range 2 {
		out := performJSONRequest(t, router, http.MethodGet, recordingHomeSourcesPath+"?workspace_id=W1", nil)
		require.Equal(t, http.StatusOK, out.Code, out.Body.String())
		var body handlers.RecordingHomeSourcesEnvelope
		require.NoError(t, json.Unmarshal(out.Body.Bytes(), &body))
		require.True(t, body.Success)
		require.Equal(t, "available", body.Data.State)
		require.Equal(t, "W1", body.Data.WorkspaceID)
		require.Len(t, body.Data.SourceDigest, 64)
		require.NotNil(t, body.Data.Devices)
		require.Empty(t, body.Data.Devices)
		require.NotNil(t, body.Data.Candidates)
		require.Empty(t, body.Data.Candidates)
		require.NotNil(t, body.Data.Excluded)
		require.Empty(t, body.Data.Excluded)
	}
	after, err := s.Workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	require.Equal(t, before, after)
}
