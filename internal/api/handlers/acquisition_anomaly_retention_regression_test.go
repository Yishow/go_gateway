package handlers

import (
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/anomaly"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionAnomalyHandlerRetentionOverTargetRemainsQueryable(t *testing.T) {
	f := newAnomalyHandlerFixture(t, anomaly.WithRetention(1))
	first := f.open(t, anomaly.Event{DeviceID: "dev-a", PointID: "p-d5", Kind: anomaly.KindStale, StartedAt: anomalyHandlerEpoch}, 0)
	second := f.open(t, anomaly.Event{DeviceID: "dev-a", PointID: "p-other", Kind: anomaly.KindStale, StartedAt: anomalyHandlerEpoch.Add(time.Second)}, 0)
	f.open(t, anomaly.Event{DeviceID: "dev-b", Kind: anomaly.KindDisconnected, StartedAt: anomalyHandlerEpoch}, 0)
	r := f.router(f.repo)

	items, total := anomalyItems(t, getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies?state=active&kind=stale&limit=1&offset=1", http.StatusOK))
	require.InDelta(t, 2, total, 0)
	require.Len(t, items, 1)
	require.InDelta(t, float64(first.ID), items[0]["id"], 0)
	summary := getAnomalyJSON(t, r, "/acquisition-anomalies/summary", http.StatusOK)["data"].(map[string]any)
	require.InDelta(t, 3, summary["total_active"], 0)
	require.Equal(t, []any{
		map[string]any{"device_id": "dev-a", "active_count": float64(2)},
		map[string]any{"device_id": "dev-b", "active_count": float64(1)},
		map[string]any{"device_id": "dev-idle", "active_count": float64(0)},
	}, summary["items"])

	require.NoError(t, f.repo.End(t.Context(), first.ID, anomalyHandlerEpoch.Add(time.Minute), nil, ""))
	items, total = anomalyItems(t, getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies?state=active&kind=stale&limit=1", http.StatusOK))
	require.InDelta(t, 1, total, 0)
	require.Len(t, items, 1)
	require.InDelta(t, float64(second.ID), items[0]["id"], 0)
	items, total = anomalyItems(t, getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies?state=ended", http.StatusOK))
	require.InDelta(t, 0, total, 0)
	require.Empty(t, items)
	summary = getAnomalyJSON(t, r, "/acquisition-anomalies/summary", http.StatusOK)["data"].(map[string]any)
	require.InDelta(t, 2, summary["total_active"], 0)
	require.Equal(t, []any{
		map[string]any{"device_id": "dev-a", "active_count": float64(1)},
		map[string]any{"device_id": "dev-b", "active_count": float64(1)},
		map[string]any{"device_id": "dev-idle", "active_count": float64(0)},
	}, summary["items"])
}
