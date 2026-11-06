package api

import (
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/anomaly"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"

	"github.com/stretchr/testify/require"
)

// The anomaly routes are registered by the same NewRouter entry used by
// cmd/test_ui and read the injected repository.
func TestNewRouter_AcquisitionAnomalyRoutes(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	_, err := f.db.ExecContext(t.Context(), `UPDATE devices SET description = '' WHERE id = 'device-A'`)
	require.NoError(t, err)
	repo := anomaly.NewSQLRepository(f.db)
	started := time.Date(2026, 11, 5, 8, 0, 0, 0, time.UTC)
	_, err = repo.Open(t.Context(), anomaly.Event{DeviceID: "device-A", PointID: "point-A", Kind: anomaly.KindStale, StartedAt: started})
	require.NoError(t, err)
	services := &DatalinkServices{
		Device:               device.NewService(device.NewSQLRepository(f.db), nil),
		Point:                point.NewService(point.NewSQLRepository(f.db), pollinggroup.NewSQLRepository(f.db)),
		Workspace:            f.workspace,
		AcquisitionAnomalies: repo,
	}
	router := NewRouter(services)

	resp := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/devices/device-A/acquisition-anomalies?state=active&kind=stale", nil)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	data := decodeJSONBody(t, resp)["data"].(map[string]any)
	require.InDelta(t, 1, data["total"], 0)
	item := data["items"].([]any)[0].(map[string]any)
	require.Equal(t, "40001", item["point_address"])

	resp = performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/acquisition-anomalies/summary", nil)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	summary := decodeJSONBody(t, resp)["data"].(map[string]any)
	require.InDelta(t, 1, summary["total_active"], 0)

	resp = performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/devices/unknown/acquisition-anomalies", nil)
	require.Equal(t, http.StatusNotFound, resp.Code, resp.Body.String())

	services.AcquisitionAnomalies = nil
	unwired := NewRouter(services)
	resp = performJSONRequest(t, unwired, http.MethodGet, "/api/v1/datalink/acquisition-anomalies/summary", nil)
	require.Equal(t, http.StatusServiceUnavailable, resp.Code, "an unwired repository is reported, not shown as zero anomalies")
}
