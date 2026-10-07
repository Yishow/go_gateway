package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/diagnostics"

	"github.com/gin-gonic/gin"
)

func TestManagedHTTPProjectionPreservesTemplatesAndSuppressesSecrets(t *testing.T) {
	b := newRuntimeLogBroker(t)
	r := gin.New()
	r.Use(managedHTTPMiddleware(b, "/api/v1"), managedRecoveryMiddleware(b))
	r.GET("/devices/:id", func(c *gin.Context) {
		_ = c.Error(fmt.Errorf("password=PRIVATE_ERROR"))
		c.Status(http.StatusBadGateway)
	})
	r.GET("/panic", func(*gin.Context) { panic("PRIVATE_PANIC") })
	r.GET("/api/v1/system/logs", func(c *gin.Context) { c.Status(http.StatusOK) })
	registerManagedHTTPRoutes(r, b)
	for _, path := range []string{"/devices/PRIVATE_PATH?password=PRIVATE_QUERY", "/panic", "/api/v1/system/logs?secret=PRIVATE_READ", "/PRIVATE_UNKNOWN"} {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, http.NoBody)
		req.Header.Set("Authorization", "Bearer PRIVATE_TOKEN")
		req.Header.Set("Cookie", "PRIVATE_COOKIE")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	snapshot, err := b.Snapshot(diagnostics.Query{})
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "PRIVATE_") {
		t.Fatalf("secret reached new sink %s", data)
	}
	access, recovered := 0, 0
	for _, event := range snapshot.Records {
		if event.Code == "http.access" {
			access++
			if route, ok := event.Fields["route"].(string); !ok || route == "" {
				t.Errorf("missing route template %+v", event)
			}
			if event.Fields["request_id"] == nil {
				t.Error("missing correlation")
			}
		}
		if event.Code == "http.recovered" {
			recovered++
		}
	}
	if access != 3 || recovered != 1 {
		t.Fatalf("access=%d recovered=%d %s", access, recovered, data)
	}
}

func TestManagedHTTPAccessSkipsRoutineReads(t *testing.T) {
	b := newRuntimeLogBroker(t)
	r := gin.New()
	r.Use(managedAccessMiddleware(b, "/api/v1", 20*time.Millisecond))
	r.GET("/status", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/slow", func(c *gin.Context) { time.Sleep(30 * time.Millisecond); c.Status(http.StatusOK) })
	r.POST("/status", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.HEAD("/status", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/missing", func(c *gin.Context) { c.Status(http.StatusNotFound) })
	registerManagedHTTPRoutes(r, b)
	requests := []struct{ method, path string }{
		{http.MethodGet, "/status"}, {http.MethodHead, "/status"}, {http.MethodGet, "/status"},
		{http.MethodGet, "/slow"}, {http.MethodPost, "/status"}, {http.MethodGet, "/missing"},
	}
	for _, request := range requests {
		r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequestWithContext(t.Context(), request.method, request.path, http.NoBody))
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	snapshot, err := b.Snapshot(diagnostics.Query{})
	if err != nil {
		t.Fatal(err)
	}
	got := make([]string, 0, len(snapshot.Records))
	for _, event := range snapshot.Records {
		got = append(got, fmt.Sprint(event.Fields["method"], " ", event.Fields["route"]))
	}
	want := []string{"GET /slow", "POST /status", "GET /missing"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("access events %v, want %v", got, want)
	}
}
