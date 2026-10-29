package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRuntimeSecurityLegacyAndScopedDeliveryRedactPersistedReasons(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	keys := r.historicalRows(t, 1)
	const hostile = "password=fixture-secret /private/raw.sql driver exception"
	require.NoError(t, r.store.MarkBlocked(t.Context(), keys[0], hostile))
	authority, err := r.pipe.DeliveryForRevision(t.Context(), r.f.group["id"].(string), r.f.group["applied_revision"].(string))
	require.NoError(t, err)
	require.Equal(t, hostile, authority.Attention[0].ErrorCode)
	require.Equal(t, []string{hostile}, authority.Backlog[0].ErrorCodes)
	before := lifecycleAPIState(t, r.f.db)
	var changes int
	require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&changes))
	for name, path := range map[string]string{"legacy": writeGroupsPath + "/" + r.f.group["id"].(string) + "/delivery", "scoped": runtimeRuleContextPath + "?" + r.query(t).Encode()} {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.RemoteAddr = "203.0.113.10:4567"
			response := httptest.NewRecorder()
			r.router.ServeHTTP(response, req)
			require.Equal(t, 200, response.Code, response.Body.String())
			require.NotContains(t, response.Body.String(), "fixture-secret")
			require.NotContains(t, response.Body.String(), "/private/")
			require.NotContains(t, response.Body.String(), "driver exception")
			require.Contains(t, response.Body.String(), "old_table")
			require.Contains(t, response.Body.String(), "C1")
			require.Contains(t, response.Body.String(), "rev7")
			require.Contains(t, response.Body.String(), "unavailable")
		})
	}
	require.Equal(t, before, lifecycleAPIState(t, r.f.db))
	var after int
	require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&after))
	require.Equal(t, changes, after)
}

func TestRuntimeSecurityRemoteSourceDetailsListsAndGroupKeepSafeFacts(t *testing.T) {
	r := newRuntimeSourceRig(t)
	for _, statement := range []string{
		`UPDATE devices SET connection_config='{"host":"private-host.invalid","password":"fixture-native-secret"}',last_test_error='driver fixture-driver-secret' WHERE id='device-A'`,
		`UPDATE points SET last_error='password=fixture-driver-secret',last_value='"[::1]:502"' WHERE id='point-source'`,
		`UPDATE source_rules SET naming_prefix='password=fixture-label-secret' WHERE id='rule-R'`,
		`UPDATE tags SET display_name='private-host.invalid',description='/private/fixture.sqlite' WHERE id='tag-source'`,
		`UPDATE mappings SET transform_pipeline='[{"type":"scale","params":{"scale":0.1,"offset":2,"password":"fixture-transform-secret"}}]' WHERE id='mapping-source'`,
	} {
		_, err := r.f.db.Exec(statement)
		require.NoError(t, err)
	}
	before := locatorTotal(t, r)
	for _, mode := range []string{"detail", "rules", "mappings", "group"} {
		t.Run(mode, func(t *testing.T) {
			q := r.query(t)
			if mode != "detail" {
				for _, key := range []string{"point_id", "mapping_id", "rule_id", "rule_revision"} {
					q.Del(key)
				}
			}
			if mode == "rules" || mode == "mappings" {
				q.Set("source_list", mode)
			}
			req := httptest.NewRequest(http.MethodGet, runtimeRuleContextPath+"?"+q.Encode(), nil)
			req.RemoteAddr = "203.0.113.11:4567"
			res := httptest.NewRecorder()
			r.router.ServeHTTP(res, req)
			require.Equal(t, 200, res.Code, res.Body.String())
			for _, unsafe := range []string{"fixture-native-secret", "fixture-driver-secret", "fixture-label-secret", "fixture-transform-secret", "private-host.invalid", "/private/", `"connection_config"`, `"password"`, `"last_error"`, `"last_value"`} {
				require.NotContains(t, res.Body.String(), unsafe)
			}
			if mode != "group" {
				require.Contains(t, res.Body.String(), "40001")
			}
			if mode == "detail" {
				require.Contains(t, res.Body.String(), `"kind":"unsupported"`)
				require.Contains(t, res.Body.String(), `"state":"saved-only"`)
			}
		})
	}
	require.Equal(t, before, locatorTotal(t, r))
	// The authorized native editing surface is deliberately not the runtime DTO.
	q := url.Values{"entry": {"locator"}, "workspace_id": {r.f.record.ID}, "step": {"1"}, "device_id": {"device-A"}}
	req := httptest.NewRequest(http.MethodGet, studioLocatorPath+"?"+q.Encode(), nil)
	req.RemoteAddr = "203.0.113.11:4567"
	res := httptest.NewRecorder()
	r.router.ServeHTTP(res, req)
	data := locatorData(t, res)
	require.Contains(t, res.Body.String(), "fixture-native-secret")
	require.NotContains(t, res.Body.String(), "fixture-driver-secret")
	entry := data["locator_entry"].(map[string]any)
	metadata, _ := json.Marshal(map[string]any{"scope": entry["hydration_scope"], "rule_heads": entry["rule_heads"], "group_heads": entry["group_heads"], "revalidated": entry["revalidated_mapping_ids"]})
	require.NotContains(t, string(metadata), "fixture-native-secret")
	require.Equal(t, before, locatorTotal(t, r))
}

func TestRuntimeSecurityRemoteMonitorCannotUseLocalProcessLogs(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	for _, path := range []string{"/api/v1/system/logs", "/api/v1/system/logs/stream"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.RemoteAddr = "203.0.113.11:4567"
		res := httptest.NewRecorder()
		r.router.ServeHTTP(res, req)
		require.Equal(t, 403, res.Code, res.Body.String())
		require.Contains(t, res.Body.String(), "logs_local_only")
		require.Empty(t, res.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestRuntimeSecuritySourcePageQueryRejectsIPv6AndPairedDSN(t *testing.T) {
	r := newRuntimeSourcePageRig(t)
	r.seed(t, "device-A", "A1", 1)
	for _, value := range []string{"[2001:db8::1]:5432", "[::1]:502", "Server=localhost;Database=metrics;Integrated Security=true"} {
		t.Run(value, func(t *testing.T) {
			q := r.query(t, "rules")
			q.Set("q", value)
			req := httptest.NewRequest(http.MethodGet, runtimeRuleContextPath+"?"+q.Encode(), nil)
			req.RemoteAddr = "203.0.113.11:4567"
			res := httptest.NewRecorder()
			r.router.ServeHTTP(res, req)
			require.Equal(t, 400, res.Code, res.Body.String())
			require.NotContains(t, res.Body.String(), value)
		})
	}
}

func TestRuntimeSecurityLegacyDeliveryErrorRequestIdentityIsSafe(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	router := NewRouter(&DatalinkServices{Workspace: r.f.workspace, WriteGroups: r.f.groups})
	for _, identity := range []string{"postgres://fixture-header-secret", "/private/fixture-header-secret", "[::1]:502", "req-safe-17"} {
		t.Run(identity, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, writeGroupsPath+"/"+r.f.group["id"].(string)+"/delivery", nil)
			req.RemoteAddr = "203.0.113.11:4567"
			req.Header.Set("X-Request-ID", identity)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			require.Equal(t, 503, response.Code, response.Body.String())
			if identity == "req-safe-17" {
				require.Contains(t, response.Body.String(), identity)
			} else {
				require.NotContains(t, response.Body.String(), identity)
			}
		})
	}
}
