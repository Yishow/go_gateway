package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/config"
	"go-gateway/internal/diagnostics"

	"github.com/gin-gonic/gin"
)

func TestRuntimeLogGuardPrecedesStoreAndWildcardCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(corsMiddleware(&config.Config{API: config.APIConfig{BasePath: "/custom"}, CORS: config.CORSConfig{AllowOrigins: []string{"*"}}}))
	registerRuntimeLogRoutes(r.Group("/custom"), nil)
	r.GET("/other", func(c *gin.Context) { c.Status(http.StatusOK) })
	for _, path := range []string{"/custom/system/logs", "/custom/system/logs/stream"} {
		req := localLogRequest(t)
		req.URL.Path = path
		req.Header.Set("Sec-Fetch-Site", "same-site")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden || strings.TrimSpace(w.Body.String()) != `{"code":"logs_local_only"}` {
			t.Fatalf("guard %d %s", w.Code, w.Body.String())
		}
		for _, key := range []string{"Access-Control-Allow-Origin", "Access-Control-Allow-Credentials"} {
			if w.Header().Get(key) != "" {
				t.Fatalf("CORS leak %s=%q", key, w.Header().Get(key))
			}
		}
		if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Cross-Origin-Resource-Policy") != "same-origin" || w.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Fatal("missing safe headers")
		}
		req = localLogRequest(t)
		req.URL.Path = path
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusServiceUnavailable {
			t.Fatalf("trusted nil source=%d", w.Code)
		}
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/other", http.NoBody))
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("unrelated CORS changed")
	}
	req := localLogRequest(t)
	req.Method = http.MethodOptions
	req.URL.Path = "/custom/system/logs"
	req.Header.Set("Origin", "https://evil.example")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("wildcard log preflight")
	}
}

func TestRuntimeLogSnapshotFiltersAndSecrets(t *testing.T) {
	b := newRuntimeLogBroker(t)
	b.Emit(diagnostics.Input{Code: "runtime.started"})
	b.Emit(diagnostics.Input{Code: "runtime.degraded"})
	b.Emit(diagnostics.Input{Code: "password=NEVER_EXPOSE;postgres://PRIVATE"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	registerRuntimeLogRoutes(r.Group("/api/v1"), b)
	req := localLogRequest(t)
	req.URL.RawQuery = "source=runtime&level=warn&q=DEGRADED"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var snapshot diagnostics.Snapshot
	if err := json.Unmarshal(w.Body.Bytes(), &snapshot); err != nil {
		t.Fatal(err)
	}
	if w.Code != http.StatusOK || len(snapshot.Records) != 1 || snapshot.Records[0].Code != "runtime.degraded" || snapshot.Latest == "" {
		t.Fatalf("snapshot %d %+v", w.Code, snapshot)
	}
	if strings.Contains(w.Body.String(), "NEVER_EXPOSE") {
		t.Fatal("secret reached snapshot")
	}
	req = localLogRequest(t)
	req.URL.RawQuery = "limit=501"
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest || strings.TrimSpace(w.Body.String()) != `{"code":"logs_invalid_query"}` {
		t.Fatalf("unsafe error %d %s", w.Code, w.Body.String())
	}
}

func TestProductionRouterWiresManagedLogRoutes(t *testing.T) {
	b := newRuntimeLogBroker(t)
	r := NewRouter(&DatalinkServices{Diagnostics: b})
	b.Emit(diagnostics.Input{Code: "startup.ready"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	req := localLogRequest(t)
	req.URL.Path = strings.TrimRight(config.Get().API.BasePath, "/") + "/system/logs"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "startup.ready") {
		t.Fatalf("production route %d %s", w.Code, w.Body.String())
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("wildcard CORS")
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	snapshot, err := b.Snapshot(diagnostics.Query{})
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Records) != 1 {
		t.Fatalf("log read amplified: %+v", snapshot.Records)
	}
}

func TestRuntimeLogUnicodeLiteralSearch(t *testing.T) {
	b := newRuntimeLogBroker(t)
	b.RegisterRoutes([]string{"/Straße/[literal]"})
	b.Emit(diagnostics.Input{Code: "http.access", Fields: map[string]any{"route": "/Straße/[literal]"}})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	registerRuntimeLogRoutes(r.Group("/api/v1"), b)
	req := localLogRequest(t)
	req.URL.RawQuery = "q=STRASSE%2F%5Bliteral%5D"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var snapshot diagnostics.Snapshot
	if err := json.Unmarshal(w.Body.Bytes(), &snapshot); err != nil {
		t.Fatal(err)
	}
	if w.Code != http.StatusOK || len(snapshot.Records) != 1 {
		t.Fatalf("Unicode literal query %d %+v", w.Code, snapshot)
	}
	req = localLogRequest(t)
	req.URL.RawQuery = "q=" + strings.Repeat("%C3%9F", 256)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("case fold expansion rejected %d %s", w.Code, w.Body.String())
	}
}
