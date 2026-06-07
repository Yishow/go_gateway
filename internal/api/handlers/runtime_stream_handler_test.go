package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

type stubRuntimeStreamSource struct {
	valueCh      chan datalinkruntime.ValueEvent
	statusCh     chan datalinkruntime.DeviceStatusEvent
	deviceID     string
	pointIDs     []string
	unsubscribed bool
}

func (s *stubRuntimeStreamSource) SubscribeValueEvents(deviceID string, pointIDs []string) (<-chan datalinkruntime.ValueEvent, func()) {
	s.deviceID = deviceID
	s.pointIDs = append([]string(nil), pointIDs...)
	return s.valueCh, func() {
		s.unsubscribed = true
	}
}

func (s *stubRuntimeStreamSource) SubscribeStatusEvents(deviceID string) (<-chan datalinkruntime.DeviceStatusEvent, func()) {
	s.deviceID = deviceID
	return s.statusCh, func() {
		s.unsubscribed = true
	}
}

func TestRuntimeStreamHandler_StreamRequiresDeviceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewRuntimeStreamHandler(&stubRuntimeStreamSource{
		valueCh:  make(chan datalinkruntime.ValueEvent, 1),
		statusCh: make(chan datalinkruntime.DeviceStatusEvent, 1),
	})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/runtime/stream", nil)
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	handler.Stream(c)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestRuntimeStreamHandler_StreamUnavailableIsExplicit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewRuntimeStreamHandler(&stubRuntimeStreamSource{})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/runtime/stream?device_id=device-1", nil)
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	handler.Stream(c)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d body=%s", recorder.Code, recorder.Body.String())
	}

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			DeviceID    string `json:"device_id"`
			StreamState struct {
				State       string `json:"state"`
				Unavailable bool   `json:"unavailable"`
			} `json:"stream_state"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response failed: %v body=%s", err, recorder.Body.String())
	}
	if body.Success {
		t.Fatalf("expected unsuccessful response, got %s", recorder.Body.String())
	}
	if body.Data.DeviceID != "device-1" {
		t.Fatalf("expected selected device id in stream state, got %s", recorder.Body.String())
	}
	if body.Data.StreamState.State != "unavailable" || !body.Data.StreamState.Unavailable {
		t.Fatalf("expected explicit unavailable stream state, got %s", recorder.Body.String())
	}
}

func TestRuntimeStreamHandler_StreamWritesValueStatusAndHeartbeatEvents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	source := &stubRuntimeStreamSource{
		valueCh:  make(chan datalinkruntime.ValueEvent, 1),
		statusCh: make(chan datalinkruntime.DeviceStatusEvent, 1),
	}
	handler := &RuntimeStreamHandler{
		source:            source,
		heartbeatInterval: 5 * time.Millisecond,
	}

	recorder := httptest.NewRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "/runtime/stream?device_id=device-1&point_ids=point-1,point-2", nil).WithContext(ctx)
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	done := make(chan struct{})
	go func() {
		handler.Stream(c)
		close(done)
	}()

	source.valueCh <- datalinkruntime.ValueEvent{
		DeviceID:         "device-1",
		PointID:          "point-1",
		Address:          "40001",
		RawValue:         37.5,
		TransformedValue: 37.5,
		Quality:          schema.QualityGood,
		Stale:            false,
		Timestamp:        time.Date(2026, 3, 16, 6, 0, 0, 0, time.UTC),
	}
	source.statusCh <- datalinkruntime.DeviceStatusEvent{
		DeviceID:      "device-1",
		Status:        "warning",
		PointsTotal:   2,
		PointsHealthy: 1,
		PointsStale:   1,
		PointsError:   0,
		BreakerState:  "closed",
	}

	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("expected runtime stream handler to exit after cancel")
	}

	if got := recorder.Header().Get("Content-Type"); !strings.Contains(got, "text/event-stream") {
		t.Fatalf("expected SSE content type, got %q", got)
	}
	if source.deviceID != "device-1" {
		t.Fatalf("expected device filter, got %q", source.deviceID)
	}
	if len(source.pointIDs) != 2 || source.pointIDs[0] != "point-1" || source.pointIDs[1] != "point-2" {
		t.Fatalf("unexpected point filters: %+v", source.pointIDs)
	}
	if !source.unsubscribed {
		t.Fatal("expected unsubscribe to be called")
	}

	body := recorder.Body.String()
	if !strings.Contains(body, "event: value") {
		t.Fatalf("expected value event body, got %q", body)
	}
	if !strings.Contains(body, "event: status") {
		t.Fatalf("expected status event body, got %q", body)
	}
	if !strings.Contains(body, "event: heartbeat") {
		t.Fatalf("expected heartbeat event body, got %q", body)
	}
	if !strings.Contains(body, `"device_id":"device-1"`) {
		t.Fatalf("expected device payload in body, got %q", body)
	}
	if !strings.Contains(body, `"point_id":"point-1"`) {
		t.Fatalf("expected point payload in body, got %q", body)
	}
	if !strings.Contains(body, `"transformed_value":37.5`) {
		t.Fatalf("expected transformed value payload in body, got %q", body)
	}
	if !strings.Contains(body, `"status":"warning"`) {
		t.Fatalf("expected status payload in body, got %q", body)
	}
}
