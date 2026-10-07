package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"go-gateway/internal/diagnostics"

	"github.com/gin-gonic/gin"
)

func TestRuntimeLogBoundedReplayControlsAndLive(t *testing.T) {
	b := newRuntimeLogBroker(t)
	for range 600 {
		if !b.Emit(diagnostics.Input{Code: "runtime.started"}) {
			t.Fatal("unexpected admission drop")
		}
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	snapshot, err := b.Snapshot(diagnostics.Query{})
	if err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	registerRuntimeLogRoutes(r.Group("/api/v1"), b)
	server := httptest.NewServer(r)
	t.Cleanup(server.Close)
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/api/v1/system/logs/stream?after="+snapshot.InstanceID+":0", http.NoBody)
	client := server.Client()
	client.Timeout = 5 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	name, id, data := readRuntimeFrame(t, reader)
	var gap diagnostics.Gap
	if err := json.Unmarshal(data, &gap); err != nil {
		t.Fatal(err)
	}
	if name != "gap" || id != "" || gap.Reason != "replay_limit" {
		t.Fatalf("gap %s id=%q %+v", name, id, gap)
	}
	name, id, data = readRuntimeFrame(t, reader)
	var meta diagnostics.Metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		t.Fatal(err)
	}
	if name != "handshake" || id != snapshot.InstanceID+":100" || meta.Progress != id {
		t.Fatalf("handshake skipped replay %s %s %+v", name, id, meta)
	}
	for sequence := 101; sequence <= 600; sequence++ {
		name, id, _ = readRuntimeFrame(t, reader)
		if name != "log" || id != fmt.Sprintf("%s:%d", snapshot.InstanceID, sequence) {
			t.Fatalf("replay order %d: %s %s", sequence, name, id)
		}
	}
	b.Emit(diagnostics.Input{Code: "runtime.stopped"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	name, id, _ = readRuntimeFrame(t, reader)
	if name != "log" || id != snapshot.InstanceID+":601" {
		t.Fatalf("live after replay %s %s", name, id)
	}
}

func TestRuntimeLogResetAndFutureCursor(t *testing.T) {
	b := newRuntimeLogBroker(t)
	b.Emit(diagnostics.Input{Code: "runtime.started"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	snapshot, err := b.Snapshot(diagnostics.Query{})
	if err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	registerRuntimeLogRoutes(r.Group("/api/v1"), b)
	server := httptest.NewServer(r)
	t.Cleanup(server.Close)
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/api/v1/system/logs/stream?after="+url.QueryEscape(logTestInstance+":1"), http.NoBody)
	client := server.Client()
	client.Timeout = 3 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	name, id, _ := readRuntimeFrame(t, reader)
	if name != "reset" || id != "" {
		t.Fatalf("reset %s %s", name, id)
	}
	name, id, _ = readRuntimeFrame(t, reader)
	if name != "handshake" || id != snapshot.InstanceID+":0" {
		t.Fatalf("reset skipped replay %s %s", name, id)
	}
	name, id, _ = readRuntimeFrame(t, reader)
	if name != "log" || id != snapshot.InstanceID+":1" {
		t.Fatalf("new instance tail %s %s", name, id)
	}
	for _, path := range []string{"/api/v1/system/logs?before=", "/api/v1/system/logs/stream?after="} {
		req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+path+snapshot.InstanceID+":999999", http.NoBody)
		future, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if future.StatusCode != http.StatusBadRequest {
			t.Errorf("future cursor %d", future.StatusCode)
		}
		_ = future.Body.Close()
	}
}
