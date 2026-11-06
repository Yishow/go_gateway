package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/anomaly"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var anomalyHandlerEpoch = time.Date(2026, 11, 5, 8, 0, 0, 0, time.UTC)

type anomalyHandlerFixture struct {
	db      *sql.DB
	repo    *anomaly.SQLRepository
	devices *device.Service
	points  *point.Service
}

func newAnomalyHandlerFixture(t *testing.T, opts ...anomaly.SQLOption) *anomalyHandlerFixture {
	t.Helper()
	path := filepath.ToSlash(filepath.Join(t.TempDir(), "anomaly-handler.db"))
	db, err := sql.Open("sqlite", strings.Replace(datalink.DefaultEmbeddedSQLiteDSN, "datalink.db", path, 1))
	require.NoError(t, err)
	datalink.ApplySQLitePoolDefaults(db)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	devRepo := device.NewSQLRepository(db)
	for _, id := range []string{"dev-a", "dev-b", "dev-idle"} {
		require.NoError(t, devRepo.Create(t.Context(), &schema.Device{ID: id, Name: id, Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusDraft, ConnectionConfig: "{}"}))
	}
	pointRepo := point.NewSQLRepository(db)
	now := time.Now().UTC()
	require.NoError(t, pointRepo.Create(t.Context(), &schema.Point{ID: "p-d5", DeviceID: "dev-a", Name: "Temp D5", Address: "D5", Function: "read_holding_registers", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, CreatedAt: now, UpdatedAt: now}))
	return &anomalyHandlerFixture{
		db: db, repo: anomaly.NewSQLRepository(db, opts...),
		devices: device.NewService(devRepo, nil), points: point.NewService(pointRepo, nil),
	}
}

func (f *anomalyHandlerFixture) open(t *testing.T, ev anomaly.Event, endAfter time.Duration) anomaly.Event {
	t.Helper()
	opened, err := f.repo.Open(t.Context(), ev)
	require.NoError(t, err)
	if endAfter > 0 {
		require.NoError(t, f.repo.End(t.Context(), opened.ID, opened.StartedAt.Add(endAfter), nil, ""))
	}
	return opened
}

func (f *anomalyHandlerFixture) router(reader AcquisitionAnomalyReader) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAcquisitionAnomalyHandler(reader, f.devices, f.points)
	r.GET("/devices/:id/acquisition-anomalies", h.ListDeviceAnomalies)
	r.GET("/acquisition-anomalies/summary", h.Summary)
	return r
}

func getAnomalyJSON(t *testing.T, r http.Handler, target string, wantStatus int) map[string]any {
	t.Helper()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, target, http.NoBody)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	require.Equal(t, wantStatus, rec.Code, rec.Body.String())
	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	return body
}

func anomalyItems(t *testing.T, body map[string]any) (items []map[string]any, total float64) {
	t.Helper()
	require.Equal(t, true, body["success"])
	data := body["data"].(map[string]any)
	raw := data["items"].([]any)
	items = make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		items = append(items, item.(map[string]any))
	}
	return items, data["total"].(float64)
}

func TestAcquisitionAnomalyHandlerFiltersByStateAndKind(t *testing.T) {
	f := newAnomalyHandlerFixture(t)
	value := "17"
	f.open(t, anomaly.Event{DeviceID: "dev-a", PointID: "p-d5", Kind: anomaly.KindStale, StartedAt: anomalyHandlerEpoch}, 12*time.Second)
	activeStale := f.open(t, anomaly.Event{DeviceID: "dev-a", PointID: "p-d5", Kind: anomaly.KindStale, StartedAt: anomalyHandlerEpoch.Add(time.Minute), LastValue: &value, Detail: anomaly.DetailNoSuccessfulRead}, 0)
	f.open(t, anomaly.Event{DeviceID: "dev-a", PointID: "p-d5", Kind: anomaly.KindFrozen, StartedAt: anomalyHandlerEpoch}, 0)
	f.open(t, anomaly.Event{DeviceID: "dev-a", Kind: anomaly.KindDisconnected, StartedAt: anomalyHandlerEpoch}, 30*time.Second)
	f.open(t, anomaly.Event{DeviceID: "dev-b", PointID: "p-other", Kind: anomaly.KindStale, StartedAt: anomalyHandlerEpoch}, 0)
	r := f.router(f.repo)

	items, total := anomalyItems(t, getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies?state=active&kind=stale", http.StatusOK))
	require.InDelta(t, 1, total, 0)
	require.Len(t, items, 1)
	got := items[0]
	require.InDelta(t, float64(activeStale.ID), got["id"], 0)
	require.Equal(t, "dev-a", got["device_id"])
	require.Equal(t, "p-d5", got["point_id"])
	require.Equal(t, "D5", got["point_address"], "the modal shows the point address")
	require.Equal(t, "Temp D5", got["point_name"])
	require.Equal(t, "stale", got["kind"])
	require.Equal(t, "active", got["state"])
	require.Equal(t, activeStale.StartedAt.Format(time.RFC3339Nano), got["started_at"])
	require.Nil(t, got["ended_at"])
	require.Nil(t, got["duration_ms"], "an active event has no final duration")
	require.Equal(t, "17", got["last_value"])
	require.Equal(t, anomaly.DetailNoSuccessfulRead, got["detail"])

	items, total = anomalyItems(t, getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies?state=ended", http.StatusOK))
	require.InDelta(t, 2, total, 0)
	require.Len(t, items, 2)
	for _, item := range items {
		require.Equal(t, "ended", item["state"])
		require.NotNil(t, item["ended_at"])
	}
	var disconnected map[string]any
	for _, item := range items {
		if item["kind"] == "disconnected" {
			disconnected = item
		}
	}
	require.NotNil(t, disconnected)
	require.Nil(t, disconnected["point_id"], "disconnected events are device-level")
	require.Nil(t, disconnected["point_address"])
	require.InDelta(t, 30000, disconnected["duration_ms"], 0)

	_, total = anomalyItems(t, getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies", http.StatusOK))
	require.InDelta(t, 4, total, 0, "other devices' events are excluded")

	items, total = anomalyItems(t, getAnomalyJSON(t, r, "/devices/dev-b/acquisition-anomalies", http.StatusOK))
	require.InDelta(t, 1, total, 0)
	require.Nil(t, items[0]["point_address"], "an unknown point keeps a null address instead of a guessed one")
}

func TestAcquisitionAnomalyHandlerPaginatesNewestFirstWithTotal(t *testing.T) {
	f := newAnomalyHandlerFixture(t)
	ids := make([]int64, 0, 5)
	for i := range 5 {
		ev := f.open(t, anomaly.Event{DeviceID: "dev-a", PointID: "p-d5", Kind: anomaly.KindStale, StartedAt: anomalyHandlerEpoch.Add(time.Duration(i) * time.Minute)}, time.Second)
		ids = append(ids, ev.ID)
	}
	r := f.router(f.repo)
	body := getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies?limit=2&offset=1", http.StatusOK)
	items, total := anomalyItems(t, body)
	require.InDelta(t, 5, total, 0)
	require.Len(t, items, 2)
	require.InDelta(t, float64(ids[3]), items[0]["id"], 0)
	require.InDelta(t, float64(ids[2]), items[1]["id"], 0)
	data := body["data"].(map[string]any)
	require.InDelta(t, 2, data["limit"], 0)
	require.InDelta(t, 1, data["offset"], 0)

	items, total = anomalyItems(t, getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies?offset=10", http.StatusOK))
	require.InDelta(t, 5, total, 0)
	require.Empty(t, items, "a page past the end is empty but still reports the total")

	body = getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies", http.StatusOK)
	require.InDelta(t, defaultAcquisitionAnomalyLimit, body["data"].(map[string]any)["limit"], 0)
}

func TestAcquisitionAnomalyHandlerRetentionKeepsActiveEvents(t *testing.T) {
	f := newAnomalyHandlerFixture(t, anomaly.WithRetention(5))
	active := map[float64]bool{}
	for i, kind := range []anomaly.Kind{anomaly.KindStale, anomaly.KindFrozen} {
		ev := f.open(t, anomaly.Event{DeviceID: "dev-a", PointID: "p-d5", Kind: kind, StartedAt: anomalyHandlerEpoch.Add(time.Duration(i) * time.Second)}, 0)
		active[float64(ev.ID)] = true
	}
	ev := f.open(t, anomaly.Event{DeviceID: "dev-a", Kind: anomaly.KindDisconnected, StartedAt: anomalyHandlerEpoch.Add(2 * time.Second)}, 0)
	active[float64(ev.ID)] = true
	oldestEnded := f.open(t, anomaly.Event{DeviceID: "dev-a", PointID: "p-old", Kind: anomaly.KindStale, StartedAt: anomalyHandlerEpoch.Add(time.Hour)}, time.Second)
	f.open(t, anomaly.Event{DeviceID: "dev-a", PointID: "p-mid", Kind: anomaly.KindStale, StartedAt: anomalyHandlerEpoch.Add(2 * time.Hour)}, time.Second)
	f.open(t, anomaly.Event{DeviceID: "dev-a", PointID: "p-new", Kind: anomaly.KindStale, StartedAt: anomalyHandlerEpoch.Add(3 * time.Hour)}, 0)

	r := f.router(f.repo)
	items, total := anomalyItems(t, getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies?limit=100", http.StatusOK))
	require.InDelta(t, 5, total, 0, "per-device cap is enforced")
	seen := map[float64]bool{}
	for _, item := range items {
		seen[item["id"].(float64)] = true
	}
	for id := range active {
		require.True(t, seen[id], "active events are never deleted by retention")
	}
	require.False(t, seen[float64(oldestEnded.ID)], "the oldest ended event is removed first")
	_, activeTotal := anomalyItems(t, getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies?state=active", http.StatusOK))
	require.InDelta(t, 4, activeTotal, 0)
}

func TestAcquisitionAnomalyHandlerRejectsInvalidQueries(t *testing.T) {
	f := newAnomalyHandlerFixture(t)
	r := f.router(f.repo)
	cases := map[string]string{
		"state=open":   "state",
		"kind=hot":     "kind",
		"limit=0":      "limit",
		"limit=501":    "limit",
		"limit=abc":    "limit",
		"limit=1.5":    "limit",
		"offset=-1":    "offset",
		"offset=x":     "offset",
		"state=ACTIVE": "state",
	}
	for query, field := range cases {
		t.Run(query, func(t *testing.T) {
			body := getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies?"+query, http.StatusBadRequest)
			require.Equal(t, false, body["success"])
			errBody := body["error"].(map[string]any)
			require.Equal(t, ErrCodeAcquisitionAnomalyInvalidFilter, errBody["code"])
			require.Equal(t, field, errBody["field"])
			require.Nil(t, body["data"], "a rejected filter never looks like an empty list")
		})
	}
}

func TestAcquisitionAnomalyHandlerUnknownDeviceIsNotFound(t *testing.T) {
	f := newAnomalyHandlerFixture(t)
	body := getAnomalyJSON(t, f.router(f.repo), "/devices/missing/acquisition-anomalies", http.StatusNotFound)
	require.Equal(t, false, body["success"])
	require.Equal(t, ErrCodeAcquisitionAnomalyDeviceNotFound, body["error"].(map[string]any)["code"])
}

type failingAnomalyReader struct{ err error }

func (r failingAnomalyReader) List(context.Context, anomaly.ListFilter) ([]anomaly.Event, int, error) {
	return nil, 0, r.err
}

func (r failingAnomalyReader) ActiveCounts(context.Context) (map[string]int, error) {
	return nil, r.err
}

func TestAcquisitionAnomalyHandlerStorageFailureIsNotAnEmptyList(t *testing.T) {
	f := newAnomalyHandlerFixture(t)
	r := f.router(failingAnomalyReader{err: errors.New("database is locked at /secret/path.db")})
	for _, target := range []string{"/devices/dev-a/acquisition-anomalies", "/acquisition-anomalies/summary"} {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, target, http.NoBody)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code, target)
		require.NotContains(t, rec.Body.String(), "/secret/path.db", "raw backend errors must not leak")
		var body map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		require.Equal(t, false, body["success"])
		require.Nil(t, body["data"])
		errBody := body["error"].(map[string]any)
		require.Equal(t, ErrCodeAcquisitionAnomalyUnavailable, errBody["code"])
		require.Equal(t, true, errBody["retryable"])
	}

	r = f.router(failingAnomalyReader{err: anomaly.ErrInvalidFilter})
	body := getAnomalyJSON(t, r, "/devices/dev-a/acquisition-anomalies", http.StatusBadRequest)
	require.Equal(t, ErrCodeAcquisitionAnomalyInvalidFilter, body["error"].(map[string]any)["code"])
}

func TestAcquisitionAnomalyHandlerWithoutRepositoryIsUnavailable(t *testing.T) {
	f := newAnomalyHandlerFixture(t)
	r := f.router(nil)
	for _, target := range []string{"/devices/dev-a/acquisition-anomalies", "/acquisition-anomalies/summary"} {
		body := getAnomalyJSON(t, r, target, http.StatusServiceUnavailable)
		require.Equal(t, false, body["success"])
		require.Equal(t, ErrCodeAcquisitionAnomalyCapabilityUnavailable, body["error"].(map[string]any)["code"])
	}
}

func TestAcquisitionAnomalyHandlerSummaryCountsActiveEventsPerDevice(t *testing.T) {
	f := newAnomalyHandlerFixture(t)
	f.open(t, anomaly.Event{DeviceID: "dev-a", PointID: "p-d5", Kind: anomaly.KindStale, StartedAt: anomalyHandlerEpoch}, 0)
	f.open(t, anomaly.Event{DeviceID: "dev-a", Kind: anomaly.KindDisconnected, StartedAt: anomalyHandlerEpoch}, 0)
	f.open(t, anomaly.Event{DeviceID: "dev-a", PointID: "p-d5", Kind: anomaly.KindFrozen, StartedAt: anomalyHandlerEpoch}, time.Second)
	f.open(t, anomaly.Event{DeviceID: "dev-b", PointID: "p-x", Kind: anomaly.KindStale, StartedAt: anomalyHandlerEpoch}, 0)

	body := getAnomalyJSON(t, f.router(f.repo), "/acquisition-anomalies/summary", http.StatusOK)
	require.Equal(t, true, body["success"])
	data := body["data"].(map[string]any)
	require.InDelta(t, 3, data["total_active"], 0)
	counts := map[string]float64{}
	for _, raw := range data["items"].([]any) {
		item := raw.(map[string]any)
		counts[item["device_id"].(string)] = item["active_count"].(float64)
	}
	require.Equal(t, map[string]float64{"dev-a": 2, "dev-b": 1, "dev-idle": 0}, counts, "devices without active events report an explicit zero")
}
