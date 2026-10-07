package api

import (
	"bufio"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/diagnostics"

	"github.com/gin-gonic/gin"
)

func TestRuntimeLogReplayLiveAndFilteredProgress(t *testing.T) {
	b := newRuntimeLogBroker(t)
	b.Emit(diagnostics.Input{Code: "runtime.started"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	snapshot, err := b.Snapshot(diagnostics.Query{})
	if err != nil {
		t.Fatal(err)
	}
	b.Emit(diagnostics.Input{Code: "startup.ready"})
	b.Emit(diagnostics.Input{Code: "runtime.degraded"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	h := &runtimeLogHandler{broker: b, heartbeat: 20 * time.Millisecond}
	r.GET("/logs", localLogGuard(), h.stream)
	server := httptest.NewServer(r)
	t.Cleanup(server.Close)
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/logs?source=runtime&after="+url.QueryEscape(snapshot.Latest), http.NoBody)
	client := server.Client()
	client.Timeout = 3 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	reader := bufio.NewReader(resp.Body)
	name, _, data := readRuntimeFrame(t, reader)
	var meta diagnostics.Metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		t.Fatal(err)
	}
	if name != "handshake" || meta.Progress == snapshot.InstanceID+":3" {
		t.Fatalf("handshake skipped replay: %s %+v", name, meta)
	}
	name, id, _ := readRuntimeFrame(t, reader)
	if name != "log" || id != snapshot.InstanceID+":3" {
		t.Fatalf("replay %s %s", name, id)
	}
	b.Emit(diagnostics.Input{Code: "startup.ready"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	name, _, data = readRuntimeFrame(t, reader)
	if err := json.Unmarshal(data, &meta); err != nil {
		t.Fatal(err)
	}
	if name != "heartbeat" || meta.Progress != snapshot.InstanceID+":4" {
		t.Fatalf("filtered progress %s %+v", name, meta)
	}
	b.Emit(diagnostics.Input{Code: "runtime.stopped"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	for name == "heartbeat" {
		name, id, _ = readRuntimeFrame(t, reader)
	}
	if name != "log" || id != snapshot.InstanceID+":5" {
		t.Fatalf("live %s %s", name, id)
	}
}

func TestRuntimeLogCapacityAndShutdown(t *testing.T) {
	b := newRuntimeLogBroker(t)
	r := gin.New()
	registerRuntimeLogRoutes(r.Group("/api/v1"), b)
	server := httptest.NewServer(r)
	t.Cleanup(server.Close)
	client := server.Client()
	client.Timeout = 3 * time.Second
	for range 8 {
		req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/api/v1/system/logs/stream", http.NoBody)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = resp.Body.Close() })
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("early capacity rejection %d", resp.StatusCode)
		}
	}
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/api/v1/system/logs/stream", http.NoBody)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("ninth stream %d", resp.StatusCode)
	}
	b.CloseSubscriptions()
	req, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/api/v1/system/logs/stream", http.NoBody)
	closed, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer closed.Body.Close()
	if closed.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("closed stream %d", closed.StatusCode)
	}
}

type logDeadlineRecorder struct {
	*httptest.ResponseRecorder
	deadlines []time.Time
	flushErr  error
}

func (w *logDeadlineRecorder) SetWriteDeadline(deadline time.Time) error {
	w.deadlines = append(w.deadlines, deadline)
	return nil
}
func (w *logDeadlineRecorder) FlushError() error { return w.flushErr }

func TestRuntimeLogFrameDeadlineAndFlushFailure(t *testing.T) {
	w := &logDeadlineRecorder{ResponseRecorder: httptest.NewRecorder()}
	start := time.Now()
	if err := writeRuntimeLogFrame(w, http.NewResponseController(w), "heartbeat", logTestInstance+":0", diagnostics.Metadata{}); err != nil {
		t.Fatal(err)
	}
	if len(w.deadlines) != 2 || w.deadlines[0].Before(start.Add(4*time.Second)) || w.deadlines[0].After(time.Now().Add(5*time.Second)) || !w.deadlines[1].IsZero() {
		t.Fatalf("deadlines=%v", w.deadlines)
	}
	failure := errors.New("transport failure")
	w.flushErr = failure
	if err := writeRuntimeLogFrame(w, http.NewResponseController(w), "heartbeat", "", diagnostics.Metadata{}); !errors.Is(err, failure) {
		t.Fatalf("flush error swallowed: %v", err)
	}
}

func TestRuntimeLogUnsupportedDeadlineFailsClosed(t *testing.T) {
	b := newRuntimeLogBroker(t)
	r := gin.New()
	registerRuntimeLogRoutes(r.Group("/api/v1"), b)
	req := localLogRequest(t)
	req.URL.Path = "/api/v1/system/logs/stream"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable || strings.Contains(w.Body.String(), "handshake") {
		t.Fatalf("unsupported writer %d %s", w.Code, w.Body.String())
	}
	subscriptions := make([]*diagnostics.Subscription, 0, 8)
	defer func() {
		for _, sub := range subscriptions {
			sub.Close()
		}
	}()
	for range 8 {
		sub, err := b.Subscribe(diagnostics.Query{}, "")
		if err != nil {
			t.Fatal(err)
		}
		subscriptions = append(subscriptions, sub)
	}
}

// blockingLogRecorder stalls the first log frame so the live queue overflows.
type blockingLogRecorder struct {
	logDeadlineRecorder
	entered, release chan struct{}
	blocked          bool
}

func (w *blockingLogRecorder) Write(p []byte) (int, error) {
	if !w.blocked && strings.HasPrefix(string(p), "event: log") {
		w.blocked = true
		close(w.entered)
		<-w.release
	}
	return w.logDeadlineRecorder.Write(p)
}

func TestRuntimeLogOverflowGapStaysWithItsOwnStream(t *testing.T) {
	b := newRuntimeLogBroker(t)
	r := gin.New()
	h := &runtimeLogHandler{broker: b, heartbeat: time.Minute}
	r.GET("/logs", h.stream)
	w := &blockingLogRecorder{logDeadlineRecorder: logDeadlineRecorder{ResponseRecorder: httptest.NewRecorder()}, entered: make(chan struct{}), release: make(chan struct{})}
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/logs", http.NoBody)
	done := make(chan struct{})
	go func() { defer close(done); r.ServeHTTP(w, req) }()
	b.Emit(diagnostics.Input{Code: "runtime.started"})
	select {
	case <-w.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("stream did not start writing")
	}
	for range diagnostics.LiveRecords + 1 {
		b.Emit(diagnostics.Input{Code: "startup.ready"})
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	close(w.release)
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("overflowed stream did not end")
	}
	if !strings.Contains(w.Body.String(), `event: gap`) || !strings.Contains(w.Body.String(), `"reason":"subscriber_overflow"`) {
		t.Fatalf("overflowed stream got no typed gap:\n%s", w.Body.String())
	}
	fresh, err := b.Subscribe(diagnostics.Query{}, "")
	if err != nil {
		t.Fatal(err)
	}
	defer fresh.Close()
	meta := fresh.Metadata()
	if meta.SubscriberOverflow != 1 {
		t.Fatalf("overflow counter hidden: %+v", meta)
	}
	for _, gap := range meta.Gaps {
		if gap.Reason == "subscriber_overflow" {
			t.Fatalf("another client's overflow reported as this stream's gap: %+v", meta.Gaps)
		}
	}
}
