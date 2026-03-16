package handlers

import (
	"context"
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
	ch           chan datalinkruntime.ValueEvent
	deviceID     string
	pointIDs     []string
	unsubscribed bool
}

func (s *stubRuntimeStreamSource) SubscribeValueEvents(deviceID string, pointIDs []string) (<-chan datalinkruntime.ValueEvent, func()) {
	s.deviceID = deviceID
	s.pointIDs = append([]string(nil), pointIDs...)
	return s.ch, func() {
		s.unsubscribed = true
	}
}

func TestRuntimeStreamHandler_StreamWritesValueEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	source := &stubRuntimeStreamSource{
		ch: make(chan datalinkruntime.ValueEvent, 1),
	}
	handler := NewRuntimeStreamHandler(source)

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

	source.ch <- datalinkruntime.ValueEvent{
		DeviceID:         "device-1",
		PointID:          "point-1",
		Address:          "40001",
		RawValue:         37.5,
		TransformedValue: 37.5,
		Quality:          schema.QualityGood,
		Stale:            false,
		Timestamp:        time.Date(2026, 3, 16, 6, 0, 0, 0, time.UTC),
	}

	time.Sleep(50 * time.Millisecond)
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
	if !strings.Contains(body, `"device_id":"device-1"`) {
		t.Fatalf("expected device payload in body, got %q", body)
	}
	if !strings.Contains(body, `"point_id":"point-1"`) {
		t.Fatalf("expected point payload in body, got %q", body)
	}
	if !strings.Contains(body, `"transformed_value":37.5`) {
		t.Fatalf("expected transformed value payload in body, got %q", body)
	}
}
