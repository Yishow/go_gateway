package api

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/config"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/diagnostics"

	"github.com/gin-gonic/gin"
)

type runtimeDomainEmitter interface {
	EmitDomain(context.Context, []byte, runtimeview.Scope) bool
	EmitDomainUnavailable(context.Context) bool
}

func requireRuntimeDomainEmitter(t *testing.T, b *diagnostics.Broker) runtimeDomainEmitter {
	t.Helper()
	emitter, ok := any(b).(runtimeDomainEmitter)
	if !ok {
		t.Fatal("production broker has no validated domain-event admission path")
	}
	return emitter
}

func runtimeDomainSource(t *testing.T, rule, revision string) []byte {
	t.Helper()
	device := "D"
	v := diagnostics.DomainEventRef{
		ContractVersion: 1, Code: "domain.source.rejected",
		Binding: runtimeview.Binding{WorkspaceID: "W", WorkspaceRevision: "W1", EntityRevision: revision},
		Locator: runtimeview.EntityLocator{Step: 2, WorkspaceID: "W", DeviceID: &device, RuleID: &rule},
		Stage:   "blocked", Reason: "unavailable", Action: "refresh",
		Evidence: diagnostics.DomainEvidenceRef{Kind: "source-runtime-transition", ReferenceID: "attempt-" + revision},
	}
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func runtimeDomainScope() runtimeview.Scope {
	return runtimeview.Scope{WorkspaceID: "W", WorkspaceRevision: "W1", DeviceID: "D"}
}

func newRuntimeDomainFileBroker(t *testing.T) (*diagnostics.Broker, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	db := filepath.Join(root, "protected.db")
	if err := os.WriteFile(db, []byte("protected test data"), 0o600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "logs", "runtime.jsonl")
	file, err := diagnostics.OpenFileSink(diagnostics.FileOptions{Path: path, DatabasePath: db, DatabaseID: strings.Repeat("a", 64)})
	if err != nil {
		t.Fatal(err)
	}
	b, err := diagnostics.New(diagnostics.Options{Level: "debug", File: file})
	if err != nil {
		_ = file.Close(t.Context())
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := b.Close(ctx); err != nil {
			t.Error(err)
		}
	})
	return b, path
}

func runtimeDomainSnapshot(t *testing.T, router http.Handler) diagnostics.Snapshot {
	t.Helper()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, localLogRequest(t))
	var snapshot diagnostics.Snapshot
	if err := json.Unmarshal(w.Body.Bytes(), &snapshot); err != nil || w.Code != http.StatusOK {
		t.Fatalf("snapshot status=%d error=%v body=%s", w.Code, err, w.Body.String())
	}
	return snapshot
}

func requireRuntimeDomainSafe(t *testing.T, raw []byte) {
	t.Helper()
	for _, forbidden := range []string{"password_fixture", "nested_secret", "driver_exception", "private-host", "postgres://", "SELECT secret", "/home/private", "sample_value", "raw_error", "credential", "<img", strings.Repeat("oversized", 20)} {
		if strings.Contains(string(raw), forbidden) {
			t.Fatalf("unsafe candidate text reached sink: %q", forbidden)
		}
	}
}

func TestRuntimeLogDomainMemoryFileSnapshotReplayAndLive(t *testing.T) {
	if runtimeLogHeartbeat != 15*time.Second {
		t.Fatal("domain summaries changed the existing 15-second stream heartbeat")
	}
	b, path := newRuntimeDomainFileBroker(t)
	emitter := requireRuntimeDomainEmitter(t, b)
	if !emitter.EmitDomain(t.Context(), runtimeDomainSource(t, "S", "R1"), runtimeDomainScope()) {
		t.Fatal("first source event was not captured")
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.Use(managedAccessMiddleware(b, "/api/v1", 0))
	registerRuntimeLogRoutes(router.Group("/api/v1"), b)
	initial := runtimeDomainSnapshot(t, router)
	if len(initial.Records) != 1 {
		t.Fatalf("first event missing: %+v", initial)
	}
	first := initial.Records[0]
	if first.Code != "domain.source.rejected" || first.Fields["workspace_revision"] != "W1" || first.Fields["entity_revision"] != "R1" || first.Fields["rule_id"] != "S" || first.Fields["reference_id"] != "attempt-R1" || first.Fields["occurrences"] != "1" {
		t.Fatalf("source proof lost or reprojected: %+v", first)
	}
	if first.InstanceID == "" || first.Sequence != "1" || first.Timestamp.IsZero() || first.Fields["stage"] != "blocked" || first.Fields["action"] != "refresh" {
		t.Fatalf("invalid envelope: %+v", first)
	}
	base := string(runtimeDomainSource(t, "S", "R2"))
	candidates := [][]byte{
		[]byte(strings.TrimSuffix(base, "}") + `,"raw_error":"password_fixture","credential":{"password":"nested_secret"}}`),
		[]byte(strings.Replace(base, "attempt-R2", "postgres://password_fixture@private-host/db", 1)),
		[]byte(strings.Replace(base, "attempt-R2", "/home/private/driver_exception", 1)),
		[]byte(strings.Replace(base, "attempt-R2", "SELECT secret", 1)),
		[]byte(strings.Replace(base, "attempt-R2", "<img src=x onerror=alert(1)>", 1)),
		[]byte(strings.TrimSuffix(base, "}") + `,"sample_value":123}`),
		[]byte(strings.Replace(base, "domain.source.rejected", "unknown.domain.password_fixture", 1)),
		[]byte(strings.Replace(base, "attempt-R2", strings.Repeat("oversized", 1100), 1)),
	}
	for _, raw := range candidates {
		if !emitter.EmitDomain(t.Context(), raw, runtimeDomainScope()) {
			t.Fatal("safe suppression was not admitted")
		}
		if err := b.Flush(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	if !emitter.EmitDomainUnavailable(t.Context()) {
		t.Fatal("fixed unattributed unavailable missing")
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	b.Emit(diagnostics.Input{Code: "domain.source.rejected", Fields: map[string]any{"raw_error": "password_fixture", "rule_id": "forged-rule"}})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	snapshot := runtimeDomainSnapshot(t, router)
	if len(snapshot.Records) != 11 || snapshot.CaptureDropped != 0 || snapshot.FileDropped != 0 {
		t.Fatalf("loss or log-read amplification: %+v", snapshot)
	}
	for _, e := range snapshot.Records[1:9] {
		if e.Code != "raw.suppressed" || len(e.Fields) != 0 {
			t.Fatalf("malformed candidate was attributed: %+v", e)
		}
	}
	if snapshot.Records[9].Code != "domain.unavailable" || len(snapshot.Records[9].Fields) != 0 || snapshot.Records[10].Code != "raw.suppressed" || len(snapshot.Records[10].Fields) != 0 {
		t.Fatal("unavailable or legacy unknown candidate retained attribution")
	}
	encoded, _ := json.Marshal(snapshot)
	requireRuntimeDomainSafe(t, encoded)
	fileRaw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	requireRuntimeDomainSafe(t, fileRaw)
	lines := strings.Split(strings.TrimSpace(string(fileRaw)), "\n")
	var ownership struct {
		Owner    string `json:"diagnostics_owner"`
		Database string `json:"database_id"`
		Set      string `json:"set_id"`
	}
	if err := json.Unmarshal([]byte(lines[0]), &ownership); err != nil || ownership.Owner != "gateway-runtime-v1" || ownership.Database != strings.Repeat("a", 64) || len(ownership.Set) != 64 {
		t.Fatalf("existing managed-file ownership header missing: %+v (%v)", ownership, err)
	}
	lines = lines[1:] // Native ownership metadata is not a captured event.
	if len(lines) != len(snapshot.Records) {
		t.Fatalf("file event count=%d snapshot=%d", len(lines), len(snapshot.Records))
	}
	for i, line := range lines {
		var stored diagnostics.Event
		if err := json.Unmarshal([]byte(line), &stored); err != nil || !reflect.DeepEqual(stored, snapshot.Records[i]) {
			t.Fatalf("file and memory differ at %d: %v", i, err)
		}
	}

	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	client := server.Client()
	client.Timeout = 3 * time.Second
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/api/v1/system/logs/stream?after="+url.QueryEscape(first.Cursor()), http.NoBody)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("stream status=%d", resp.StatusCode)
	}
	reader := bufio.NewReader(resp.Body)
	name, _, data := readRuntimeFrame(t, reader)
	if name != "handshake" {
		t.Fatalf("missing handshake: %s %s", name, data)
	}
	requireRuntimeDomainSafe(t, data)
	for _, expected := range snapshot.Records[1:] {
		name, id, raw := readRuntimeFrame(t, reader)
		requireRuntimeDomainSafe(t, raw)
		var actual diagnostics.Event
		if err := json.Unmarshal(raw, &actual); err != nil || name != "log" || id != expected.Cursor() || !reflect.DeepEqual(actual, expected) {
			t.Fatalf("replay identity/projection mismatch: %s %s %+v (%v)", name, id, actual, err)
		}
	}
	if !emitter.EmitDomain(t.Context(), runtimeDomainSource(t, "S", "R2"), runtimeDomainScope()) {
		t.Fatal("new revision live event missing")
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	name, id, data := readRuntimeFrame(t, reader)
	requireRuntimeDomainSafe(t, data)
	var live diagnostics.Event
	if err := json.Unmarshal(data, &live); err != nil || name != "log" || id != live.Cursor() || live.Fields["entity_revision"] != "R2" || live.Fields["occurrences"] != "1" || live.Code != "domain.source.rejected" {
		t.Fatalf("new-revision live proof lost: %s %s %+v (%v)", name, id, live, err)
	}
	if final := runtimeDomainSnapshot(t, router); len(final.Records) != 12 || final.CaptureDropped != 0 {
		t.Fatalf("log reads created events or loss: %+v", final)
	}
}

func TestRuntimeLogDomainGuardPreservesLocalOnly(t *testing.T) {
	b := newRuntimeLogBroker(t)
	if !requireRuntimeDomainEmitter(t, b).EmitDomain(t.Context(), runtimeDomainSource(t, "S", "R1"), runtimeDomainScope()) {
		t.Fatal("source event missing")
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.Use(corsMiddleware(&config.Config{API: config.APIConfig{BasePath: "/api/v1"}, CORS: config.CORSConfig{AllowOrigins: []string{"*"}}}))
	registerRuntimeLogRoutes(router.Group("/api/v1"), b)
	cases := []struct {
		name    string
		edit    func(*http.Request)
		allowed bool
	}{
		{"IPv4", func(r *http.Request) { r.Header.Set("Origin", "http://127.0.0.1:8080") }, true},
		{"IPv6", func(r *http.Request) {
			r.Host, r.RemoteAddr = "[::1]:8080", "[::1]:5000"
			r.Header.Set("Origin", "http://[::1]:8080")
			*r = *r.WithContext(context.WithValue(r.Context(), http.LocalAddrContextKey, &net.TCPAddr{IP: net.ParseIP("::1"), Port: 8080}))
		}, true},
		{"LAN peer", func(r *http.Request) { r.RemoteAddr = "192.168.1.2:5000" }, false},
		{"LAN destination", func(r *http.Request) {
			*r = *r.WithContext(context.WithValue(r.Context(), http.LocalAddrContextKey, &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 8080}))
		}, false},
		{"proxy", func(r *http.Request) { r.Header.Set("X-Forwarded-For", "127.0.0.1") }, false},
		{"forwarded", func(r *http.Request) { r.Header.Set("Forwarded", "for=127.0.0.1") }, false},
		{"foreign origin", func(r *http.Request) { r.Header.Set("Origin", "https://private-host") }, false},
		{"same-site", func(r *http.Request) { r.Header.Set("Sec-Fetch-Site", "same-site") }, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, path := range []string{"/api/v1/system/logs", "/api/v1/system/logs/stream"} {
				req := localLogRequest(t)
				req.URL.Path = path
				tc.edit(req)
				ctx, cancel := context.WithCancel(req.Context())
				cancel() // Allows a real stream handshake, then exits without waiting for a timer.
				req = req.WithContext(ctx)
				w := &logDeadlineRecorder{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(w, req)
				if !tc.allowed {
					if w.Code != http.StatusForbidden || strings.TrimSpace(w.Body.String()) != `{"code":"logs_local_only"}` {
						t.Fatalf("%s guard=%d %s", path, w.Code, w.Body.String())
					}
				} else if w.Code != http.StatusOK || (path == "/api/v1/system/logs" && !strings.Contains(w.Body.String(), "domain.source.rejected")) || (path != "/api/v1/system/logs" && !strings.Contains(w.Body.String(), "event: handshake")) {
					t.Fatalf("local projection missing: %s %d %s", path, w.Code, w.Body.String())
				}
				if w.Header().Get("Access-Control-Allow-Origin") != "" || w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Cross-Origin-Resource-Policy") != "same-origin" || w.Header().Get("X-Content-Type-Options") != "nosniff" {
					t.Fatalf("safe local headers changed: %v", w.Header())
				}
			}
		})
	}
}
