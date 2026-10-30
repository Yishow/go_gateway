package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"go-gateway/internal/config"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/diagnostics"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRuntimeLogDomainFiltersStrictHTTPContract(t *testing.T) {
	b := newRuntimeLogBroker(t)
	router := gin.New()
	registerRuntimeLogRoutes(router.Group("/api/v1"), b)
	valid := []string{"workspace_id=W", "device_id=D", "rule_id=S", "group_id=G", "stage=blocked", "actionable=true", "actionable=false", "workspace_id=W&device_id=D&rule_id=S&group_id=G&stage=blocked&actionable=true", "device_id=opaque-1_a.2", "workspace_id=" + strings.Repeat("a", 128), "stage=unavailable", "stage=sql-committed"}
	for _, query := range valid {
		t.Run("valid/"+query[:min(60, len(query))], func(t *testing.T) {
			for _, stream := range []bool{false, true} {
				w := runtimeDomainFilterRequest(t, router, query, stream, nil)
				require.Equal(t, http.StatusOK, w.Code, w.Body.String())
			}
		})
	}
	invalid := make([]string, 0, 79)
	invalid = append(invalid, "actionable=1", "actionable=0", "actionable=True", "actionable=FALSE", "actionable=%20true", "stage=BLOCKED", "stage=not-a-stage", "workspace_id=%ZZ", "workspace_id=W;device_id=D", "workspace_id=W&workspace_id=%ZZ", "workspace%ZZ_id=W", "workspace_id=W&unknown=secret_fixture", "workspace_id=W&"+strings.Repeat("x", 4097))
	for _, key := range []string{"workspace_id", "device_id", "rule_id", "group_id", "stage", "actionable"} {
		invalid = append(invalid, key+"=", key, key+"=W&"+key+"=W")
	}
	for _, key := range []string{"workspace_id", "device_id", "rule_id", "group_id"} {
		for _, value := range []string{" leading", "trailing ", "a/b", "postgres://private-host/db", "private.example", "127.0.0.1", "[::1]", "host:502", "/home/private", "秘密", "a\x00b", strings.Repeat("a", 129)} {
			invalid = append(invalid, key+"="+url.QueryEscape(value))
		}
	}
	for i, query := range invalid {
		t.Run(fmt.Sprintf("invalid/%d", i), func(t *testing.T) {
			for _, stream := range []bool{false, true} {
				w := runtimeDomainFilterRequest(t, router, query, stream, nil)
				require.Equal(t, http.StatusBadRequest, w.Code)
				require.JSONEq(t, `{"code":"logs_invalid_query"}`, w.Body.String(), "never reflect rejected text")
			}
		})
	}
}
func runtimeDomainFilterRequest(t *testing.T, router http.Handler, query string, stream bool, edit func(*http.Request)) *httptest.ResponseRecorder {
	t.Helper()
	req := localLogRequest(t)
	req.URL.RawQuery = query
	if stream {
		req.URL.Path += "/stream"
	}
	if edit != nil {
		edit(req)
	}
	ctx, cancel := context.WithCancel(req.Context())
	cancel()
	req = req.WithContext(ctx)
	w := &logDeadlineRecorder{ResponseRecorder: httptest.NewRecorder()}
	router.ServeHTTP(w, req)
	return w.ResponseRecorder
}
func emitRuntimeFilterSource(t *testing.T, b *diagnostics.Broker, device, rule, revision string, success bool) {
	t.Helper()
	raw := runtimeDomainSource(t, rule, revision)
	var ref diagnostics.DomainEventRef
	require.NoError(t, json.Unmarshal(raw, &ref))
	ref.Locator.DeviceID = &device
	if success {
		ref.Code = "domain.source.applied"
		ref.Stage = "accepted"
		ref.Reason = "none"
		ref.Action = "none"
	}
	raw, err := json.Marshal(ref)
	require.NoError(t, err)
	scope := runtimeDomainScope()
	scope.DeviceID = device
	_, err = diagnostics.ParseDomainEventRef(raw, scope)
	require.NoError(t, err)
	require.True(t, b.EmitDomain(t.Context(), raw, scope))
	require.NoError(t, b.Flush(t.Context()))
}
func emitRuntimeFilterGroup(t *testing.T, b *diagnostics.Broker) {
	t.Helper()
	group, revision, connector, crev := "G", "G1", "C", "C1"
	ref := diagnostics.DomainEventRef{ContractVersion: 1, Code: "domain.delivery.blocked", Binding: runtimeview.Binding{WorkspaceID: "W", WorkspaceRevision: "W1", EntityRevision: revision}, Locator: runtimeview.EntityLocator{Step: 4, WorkspaceID: "W", GroupID: &group, GroupRevision: &revision}, ConnectorID: &connector, ConnectorRevision: &crev, Stage: "blocked", Reason: "unavailable", Action: "refresh", Evidence: diagnostics.DomainEvidenceRef{Kind: "durable-transition", ReferenceID: "effect-G1"}}
	raw, err := json.Marshal(ref)
	require.NoError(t, err)
	scope := runtimeview.Scope{WorkspaceID: "W", WorkspaceRevision: "W1"}
	_, err = diagnostics.ParseDomainEventRef(raw, scope)
	require.NoError(t, err)
	require.True(t, b.EmitDomain(t.Context(), raw, scope))
	require.NoError(t, b.Flush(t.Context()))
}
func runtimeFilterSnapshot(t *testing.T, r http.Handler, path, query string) diagnostics.Snapshot {
	t.Helper()
	req := localLogRequest(t)
	req.URL.Path = path
	req.URL.RawQuery = query
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	var snapshot diagnostics.Snapshot
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &snapshot))
	return snapshot
}
func TestRuntimeLogDomainFiltersActualSnapshotReplayAndLive(t *testing.T) {
	b := newRuntimeLogBroker(t)
	emitRuntimeFilterSource(t, b, "D", "S", "R1", false)
	emitRuntimeFilterSource(t, b, "E", "Other", "R1", false)
	emitRuntimeFilterSource(t, b, "D", "S", "R2", true)
	emitRuntimeFilterGroup(t, b)
	require.True(t, b.EmitDomainUnavailable(t.Context()))
	require.True(t, b.Emit(diagnostics.Input{Code: "runtime.degraded"}))
	require.NoError(t, b.Flush(t.Context()))
	router := NewRouter(&DatalinkServices{Diagnostics: b})
	path := strings.TrimRight(config.Get().API.BasePath, "/") + "/system/logs"
	cases := []struct {
		query string
		want  int
	}{
		{"", 6}, {"actionable=false", 6}, {"actionable=true", 4}, {"workspace_id=W", 4}, {"device_id=D", 2}, {"rule_id=S", 2}, {"group_id=G", 1},
		{"stage=blocked&actionable=true", 3}, {"stage=unavailable&actionable=true", 1}, {"stage=unavailable&workspace_id=W", 0}, {"device_id=D&group_id=G", 0},
		{"workspace_id=W&device_id=D&rule_id=S&stage=blocked&actionable=true&level=warn&source=runtime&q=REJECTED", 1}, {"workspace_id=other", 0}, {"rule_id=s", 0},
	}
	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) { require.Len(t, runtimeFilterSnapshot(t, router, path, tc.query).Records, tc.want) })
	}
	query := "workspace_id=W&device_id=D&rule_id=S&stage=blocked&actionable=true"
	snapshot := runtimeFilterSnapshot(t, router, path, query)
	require.Len(t, snapshot.Records, 1)
	first := snapshot.Records[0]
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	client := server.Client()
	client.Timeout = 3 * time.Second
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+path+"/stream?"+query+"&after="+url.QueryEscape(snapshot.InstanceID+":0"), http.NoBody)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	reader := bufio.NewReader(resp.Body)
	name, _, _ := readRuntimeFrame(t, reader)
	require.Equal(t, "handshake", name)
	name, id, data := readRuntimeFrame(t, reader)
	var replay diagnostics.Event
	require.NoError(t, json.Unmarshal(data, &replay))
	require.Equal(t, "log", name)
	require.Equal(t, first.Cursor(), id)
	require.Equal(t, first, replay)
	emitRuntimeFilterSource(t, b, "E", "Other", "R3", false)
	emitRuntimeFilterSource(t, b, "D", "S", "R3", false)
	name, id, data = readRuntimeFrame(t, reader)
	var live diagnostics.Event
	require.NoError(t, json.Unmarshal(data, &live))
	require.Equal(t, "log", name)
	require.Equal(t, live.Cursor(), id)
	require.Equal(t, "D", live.Fields["device_id"])
	require.Equal(t, "S", live.Fields["rule_id"])
	require.Equal(t, "R3", live.Fields["entity_revision"])
	require.Equal(t, "8", live.Sequence, "filtered rows advance native sequence, not a new filtered stream identity")
	latest := runtimeFilterSnapshot(t, router, path, query)
	require.Len(t, latest.Records, 2)
	require.Equal(t, live, latest.Records[1])
	require.Equal(t, live.Cursor(), latest.Progress)
	before := runtimeFilterSnapshot(t, router, path, query+"&before="+url.QueryEscape(live.Cursor()))
	require.Equal(t, []diagnostics.Event{first}, before.Records)
	all := runtimeFilterSnapshot(t, router, path, "")
	require.Len(t, all.Records, 8, "log requests must not amplify themselves")
}
func TestRuntimeLogDomainFiltersPreserveCursorPrecisionAndGap(t *testing.T) {
	cursor := logTestInstance + ":9007199254740993"
	req := localLogRequest(t)
	req.URL.RawQuery = "device_id=D&actionable=true&after=" + url.QueryEscape(cursor)
	req.Header.Set("Last-Event-ID", cursor)
	_, after, err := parseRuntimeLogQuery(req, true)
	require.NoError(t, err)
	require.Equal(t, cursor, after)
	req.Header.Set("Last-Event-ID", logTestInstance+":9007199254740992")
	_, _, err = parseRuntimeLogQuery(req, true)
	require.ErrorIs(t, err, diagnostics.ErrInvalidQuery, "adjacent >2^53 cursors remain distinct")
	b := newRuntimeLogBroker(t)
	emitRuntimeFilterSource(t, b, "D", "S", "R1", false)
	b.RecordDomainCaptureLoss(1)
	require.NoError(t, b.Flush(t.Context()))
	router := gin.New()
	registerRuntimeLogRoutes(router.Group("/api/v1"), b)
	snapshot := runtimeFilterSnapshot(t, router, "/api/v1/system/logs", "device_id=D&actionable=true")
	require.Len(t, snapshot.Records, 1)
	require.Equal(t, uint64(1), snapshot.CaptureDropped)
	require.NotEmpty(t, snapshot.Gaps)
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	client := server.Client()
	client.Timeout = 3 * time.Second
	req, err = http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/api/v1/system/logs/stream?device_id=D&actionable=true&after="+url.QueryEscape(cursor), http.NoBody)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	reader := bufio.NewReader(resp.Body)
	var reasons []string
	for {
		name, id, data := readRuntimeFrame(t, reader)
		if name == "handshake" {
			var m diagnostics.Metadata
			require.NoError(t, json.Unmarshal(data, &m))
			require.NotEqual(t, cursor, id)
			require.Equal(t, snapshot.InstanceID, m.InstanceID)
			require.Equal(t, uint64(1), m.CaptureDropped)
			break
		}
		require.Contains(t, []string{"gap", "reset"}, name)
		var gap diagnostics.Gap
		require.NoError(t, json.Unmarshal(data, &gap))
		reasons = append(reasons, gap.Reason)
	}
	require.Contains(t, reasons, "reset")
	require.Contains(t, reasons, "admission")
	name, _, data := readRuntimeFrame(t, reader)
	require.Equal(t, "log", name)
	var event diagnostics.Event
	require.NoError(t, json.Unmarshal(data, &event))
	require.Equal(t, "domain.source.rejected", event.Code, "retained fault stays historical despite incomplete history")
	future := runtimeDomainFilterRequest(t, router, "device_id=D&after="+snapshot.InstanceID+":9007199254740993", true, nil)
	require.Equal(t, http.StatusBadRequest, future.Code)
	require.JSONEq(t, `{"code":"logs_invalid_query"}`, future.Body.String())
}
func TestRuntimeLogDomainFiltersRetainLocalGuard(t *testing.T) {
	b := newRuntimeLogBroker(t)
	router := gin.New()
	registerRuntimeLogRoutes(router.Group("/api/v1"), b)
	for _, stream := range []bool{false, true} {
		for _, edit := range []func(*http.Request){func(r *http.Request) { r.RemoteAddr = "192.168.1.2:5000" }, func(r *http.Request) { r.Header.Set("Origin", "http://private.example") }, func(r *http.Request) { r.Header.Set("X-Forwarded-For", "127.0.0.1") }} {
			w := runtimeDomainFilterRequest(t, router, "workspace_id=W&device_id=D&stage=blocked&actionable=true", stream, edit)
			require.Equal(t, http.StatusForbidden, w.Code)
			require.JSONEq(t, `{"code":"logs_local_only"}`, w.Body.String())
			require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
			require.Equal(t, "", w.Header().Get("Access-Control-Allow-Origin"))
		}
	}
}
