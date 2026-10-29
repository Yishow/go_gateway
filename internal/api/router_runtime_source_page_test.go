package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntimeSourcePagesExposeSixRulesAndTenMappings(t *testing.T) {
	r := newRuntimeSourcePageRig(t)
	for i := 1; i <= 6; i++ {
		count := 1
		if i <= 4 {
			count = 2
		}
		r.seed(t, "device-A", fmt.Sprintf("A%d", i), count)
	}
	r.seed(t, "device-B", "B1", 1)
	for _, kind := range []string{"rules", "mappings"} {
		t.Run(kind, func(t *testing.T) {
			q := r.query(t, kind)
			q.Set("limit", "3")
			ids := []string{}
			pages := 0
			for {
				page := runtimeSourcePageBody(t, r.get(t, q))
				require.Equal(t, "device-A", page.DeviceID)
				require.Equal(t, kind, page.Kind)
				require.LessOrEqual(t, len(page.Rules)+len(page.Mappings), 3)
				rows := page.Rules
				if kind == "mappings" {
					require.Empty(t, page.Rules)
					rows = page.Mappings
				} else {
					require.Empty(t, page.Mappings)
				}
				for _, row := range rows {
					require.Equal(t, "device-A", row["device_id"])
					require.NotContains(t, fmt.Sprint(row), "B1")
					ids = append(ids, row["id"].(string))
					if kind == "mappings" {
						require.Equal(t, "rev-"+row["rule_id"].(string), row["rule_revision"])
					} else {
						require.Equal(t, "rev-"+row["id"].(string), row["revision_id"])
					}
				}
				pages++
				require.LessOrEqual(t, pages, 4)
				if page.NextCursor == nil {
					break
				}
				q.Set("cursor", *page.NextCursor)
			}
			expected := 6
			if kind == "mappings" {
				expected = 10
			}
			require.Len(t, ids, expected)
			unique := map[string]bool{}
			for _, id := range ids {
				require.False(t, unique[id], "duplicate row across pages")
				unique[id] = true
			}
		})
	}
}

func TestRuntimeSourcePageDefault50AndClamp100(t *testing.T) {
	r := newRuntimeSourcePageRig(t)
	for i := 1; i <= 120; i++ {
		r.seed(t, "device-A", fmt.Sprintf("A%03d", i), 0)
	}
	q := r.query(t, "rules")
	page := runtimeSourcePageBody(t, r.get(t, q))
	require.Equal(t, 50, page.Limit)
	require.Len(t, page.Rules, 50)
	require.NotNil(t, page.NextCursor)
	for _, tc := range []struct {
		raw   string
		limit int
	}{
		{"999999999999999999999999", 100}, {"000000000000000000000001", 1},
	} {
		q.Set("limit", tc.raw)
		bounded := runtimeSourcePageBody(t, r.get(t, q))
		require.Equal(t, tc.limit, bounded.Limit)
		require.Len(t, bounded.Rules, tc.limit)
		require.NotNil(t, bounded.NextCursor)
	}
	q.Set("limit", "500")
	page = runtimeSourcePageBody(t, r.get(t, q))
	require.Equal(t, 100, page.Limit)
	require.Len(t, page.Rules, 100)
	require.NotNil(t, page.NextCursor)
	first := map[string]bool{}
	for _, row := range page.Rules {
		first[row["id"].(string)] = true
	}
	q.Set("cursor", *page.NextCursor)
	remainder := runtimeSourcePageBody(t, r.get(t, q))
	require.Equal(t, page.ListRevision, remainder.ListRevision)
	require.Len(t, remainder.Rules, 20)
	require.Nil(t, remainder.NextCursor)
	for _, row := range remainder.Rules {
		require.False(t, first[row["id"].(string)], "remainder repeats earlier row")
	}
}

func TestRuntimeSourcePageSearchAndNestedRowsRemainBounded(t *testing.T) {
	r := newRuntimeSourcePageRig(t)
	r.seed(t, "device-A", "A1", 120)
	r.seed(t, "device-A", "A2", 2)
	r.seed(t, "device-B", "B1", 1)
	q := r.query(t, "rules")
	q.Set("q", "  temperature a2  ")
	page := runtimeSourcePageBody(t, r.get(t, q))
	require.Equal(t, "temperature a2", page.Query)
	require.Len(t, page.Rules, 1)
	require.Equal(t, "A2", page.Rules[0]["id"])
	require.NotContains(t, page.Rules[0], "mappings", "rule page must not smuggle an unbounded child array")
	q = r.query(t, "mappings")
	q.Set("rule_id", "A1")
	q.Set("rule_revision", "rev-A1")
	q.Set("limit", "100")
	page = runtimeSourcePageBody(t, r.get(t, q))
	require.Len(t, page.Mappings, 100)
	require.Empty(t, page.Rules)
	require.NotNil(t, page.NextCursor)
	for _, row := range page.Mappings {
		require.Equal(t, "A1", row["rule_id"])
		require.Equal(t, "rev-A1", row["rule_revision"])
		require.NotContains(t, row, "transforms")
	}
	q.Set("cursor", *page.NextCursor)
	remainder := runtimeSourcePageBody(t, r.get(t, q))
	require.Len(t, remainder.Mappings, 20)
	require.Nil(t, remainder.NextCursor)
	q.Del("cursor")
	q.Set("q", "temperature_a1-119")
	searched := runtimeSourcePageBody(t, r.get(t, q))
	require.Len(t, searched.Mappings, 1)
	require.Equal(t, "mapping-A1-119", searched.Mappings[0]["id"])
	q.Set("q", "absent")
	empty := runtimeSourcePageBody(t, r.get(t, q))
	require.Empty(t, empty.Mappings)
	require.Nil(t, empty.NextCursor)
}

func TestRuntimeSourcePageCursorRejectsChangedRevisionAndBindings(t *testing.T) {
	r := newRuntimeSourcePageRig(t)
	r.seed(t, "device-A", "A1", 3)
	r.seed(t, "device-A", "A2", 2)
	r.seed(t, "device-B", "B1", 2)
	q := r.query(t, "mappings")
	q.Set("limit", "1")
	q.Set("rule_id", "A1")
	q.Set("rule_revision", "rev-A1")
	first := runtimeSourcePageBody(t, r.get(t, q))
	require.NotNil(t, first.NextCursor)
	token := *first.NextCursor
	for _, tc := range []struct {
		name, key, value string
		status           int
	}{
		{"query", "q", "other", http.StatusConflict},
		{"limit", "limit", "2", http.StatusConflict},
		{"workspace-revision", "workspace_revision", "other-revision", http.StatusConflict},
		{"owned-other-device", "device_id", "device-B", http.StatusConflict},
		{"foreign-workspace", "workspace_id", "foreign", http.StatusForbidden},
		{"foreign-device", "device_id", "foreign", http.StatusForbidden},
	} {
		t.Run(tc.name, func(t *testing.T) {
			next := maps.Clone(q)
			next.Set("cursor", token)
			next.Set(tc.key, tc.value)
			response := r.get(t, next)
			require.Equal(t, tc.status, response.Code, response.Body.String())
			require.NotContains(t, response.Body.String(), `"source_page"`)
		})
	}
	next := maps.Clone(q)
	next.Set("cursor", token)
	next.Set("rule_id", "A2")
	next.Set("rule_revision", "rev-A2")
	response := r.get(t, next)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	next = maps.Clone(q)
	next.Set("cursor", token)
	next.Del("rule_id")
	next.Del("rule_revision")
	next.Set("source_list", "rules")
	response = r.get(t, next)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	raw, err := base64.RawURLEncoding.DecodeString(token)
	require.NoError(t, err)
	var cursor map[string]any
	require.NoError(t, json.Unmarshal(raw, &cursor))
	cursor["p"] = float64(999999)
	changed, err := json.Marshal(cursor)
	require.NoError(t, err)
	next = maps.Clone(q)
	next.Set("cursor", base64.RawURLEncoding.EncodeToString(changed))
	response = r.get(t, next)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	_, err = r.f.db.Exec(`UPDATE source_rules SET revision_id='changed-A1' WHERE id='A1'`)
	require.NoError(t, err)
	next = maps.Clone(q)
	next.Set("cursor", token)
	response = r.get(t, next)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
}

func TestRuntimeSourcePageRejectsMalformedForeignAndMixedModeRequests(t *testing.T) {
	r := newRuntimeSourcePageRig(t)
	r.seed(t, "device-A", "A1", 2)
	for _, tc := range []struct{ key, value string }{
		{"source_list", "unknown"}, {"limit", "0"}, {"limit", "-1"}, {"limit", "1.5"}, {"limit", "no"}, {"limit", "1e3"}, {"limit", "000000"},
		{"q", strings.Repeat("測", 129)}, {"q", "password=secret"}, {"q", "/private/list.sqlite"}, {"q", "db.private.internal"}, {"q", "line\nbreak"},
		{"cursor", ""}, {"cursor", "not-json"}, {"cursor", strings.Repeat("a", 1025)},
		{"point_id", "point-A1-000"}, {"mapping_id", "mapping-A1-000"}, {"rule_id", "A1"}, {"rule_revision", "rev-A1"},
	} {
		t.Run(tc.key+tc.value[:min(10, len(tc.value))], func(t *testing.T) {
			q := r.query(t, "rules")
			q.Set(tc.key, tc.value)
			response := r.get(t, q)
			require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
			require.NotContains(t, response.Body.String(), `"source_page"`)
			require.NotContains(t, response.Body.String(), "private")
		})
	}
	for _, key := range []string{"workspace_id", "workspace_revision", "device_id", "source_list", "q", "limit", "cursor", "rule_id", "rule_revision"} {
		t.Run("duplicate-"+key, func(t *testing.T) {
			q := r.query(t, "mappings")
			q.Set("rule_id", "A1")
			q.Set("rule_revision", "rev-A1")
			if q.Get(key) == "" {
				q.Set(key, "x")
			}
			q.Add(key, q.Get(key))
			response := r.get(t, q)
			require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
		})
	}
	for _, key := range []string{"workspace_id", "workspace_revision", "device_id"} {
		q := r.query(t, "rules")
		q.Del(key)
		response := r.get(t, q)
		require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
	}
	for _, suffix := range []string{"&cursor=%ZZ", "&q=bad;query"} {
		response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t, "rules").Encode()+suffix, nil)
		require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
	}
	q := r.query(t, "rules")
	q.Set("limit", "1")
	first := runtimeSourcePageBody(t, r.get(t, q))
	// Two mappings supply a continuation even when the rule list has one row.
	q = r.query(t, "mappings")
	q.Set("limit", "1")
	first = runtimeSourcePageBody(t, r.get(t, q))
	require.NotNil(t, first.NextCursor)
	raw, err := base64.RawURLEncoding.DecodeString(*first.NextCursor)
	require.NoError(t, err)
	var cursor map[string]any
	require.NoError(t, json.Unmarshal(raw, &cursor))
	cursor["unknown"] = true
	raw, err = json.Marshal(cursor)
	require.NoError(t, err)
	q.Set("cursor", base64.RawURLEncoding.EncodeToString(raw))
	response := r.get(t, q)
	require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
	q.Set("cursor", *first.NextCursor+"=")
	response = r.get(t, q)
	require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
}

func TestRuntimeSourcePageReadFailurePreservesUnavailable(t *testing.T) {
	for _, at := range []int{1, 2} {
		t.Run(fmt.Sprint(at), func(t *testing.T) {
			r := newRuntimeSourcePageRig(t)
			r.seed(t, "device-A", "A1", 1)
			r.fault(&runtimeSourcePageFaultRepo{Repository: r.repo, failAt: at})
			response := r.get(t, r.query(t, "rules"))
			require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
			require.Contains(t, response.Body.String(), `"code":"unavailable"`)
			require.NotContains(t, response.Body.String(), "fixture-secret")
			require.NotContains(t, response.Body.String(), "list.sqlite")
		})
	}
	r := newRuntimeSourcePageRig(t)
	r.seed(t, "device-A", "A1", 1)
	r.fault(&runtimeSourcePageFaultRepo{Repository: r.repo, changeAt: 1, change: func() {
		_, err := r.f.db.Exec(`UPDATE source_rules SET revision_id='R-new' WHERE id='A1'`)
		require.NoError(t, err)
	}})
	response := r.get(t, r.query(t, "rules"))
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), `"code":"stale-reference"`)
}

func TestRuntimeSourcePageNeverBootstrapsWritesOrLeaksRawFields(t *testing.T) {
	r := newRuntimeSourcePageRig(t)
	r.seed(t, "device-A", "A1", 2)
	for _, statement := range []string{
		`UPDATE source_rules SET naming_prefix='password=fixture-secret /private/source.sqlite' WHERE id='A1'`,
		`UPDATE tags SET display_name='db.private.internal',description='secret://raw' WHERE id='tag-A1-000'`,
		`UPDATE mappings SET transform_pipeline='[{"type":"formula","params":{"raw":"password=fixture-secret"}}]' WHERE id='mapping-A1-000'`,
	} {
		_, err := r.f.db.Exec(statement)
		require.NoError(t, err)
	}
	for _, kind := range []string{"rules", "mappings"} {
		q := r.query(t, kind)
		var before, after int
		require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&before))
		response := r.get(t, q)
		page := runtimeSourcePageBody(t, response)
		require.NotEmpty(t, page.ListRevision)
		require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&after))
		require.Equal(t, before, after)
		for _, secret := range []string{"fixture-secret", "source.sqlite", "db.private.internal", "secret://raw", "transform_pipeline", "connection_config"} {
			require.NotContains(t, response.Body.String(), secret)
		}
	}
	q := r.query(t, "rules")
	_, err := r.f.db.Exec(`DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
	require.NoError(t, err)
	response := r.get(t, q)
	var count int
	require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&count))
	require.Zero(t, count)
	require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
}

func TestRuntimeSourcePageCursorRejectsChangedListWithoutRuleRevisionChange(t *testing.T) {
	r := newRuntimeSourcePageRig(t)
	r.seed(t, "device-A", "A1", 3)
	q := r.query(t, "mappings")
	q.Set("limit", "1")
	first := runtimeSourcePageBody(t, r.get(t, q))
	require.NotNil(t, first.NextCursor)
	_, err := r.f.db.Exec(`UPDATE tags SET display_name='Changed visible name' WHERE id='tag-A1-002'`)
	require.NoError(t, err)
	q.Set("cursor", *first.NextCursor)
	response := r.get(t, q)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), `"code":"stale-reference"`)
}

func TestRuntimeSourcePageConcurrentWorkspaceDeletionDoesNotRecreate(t *testing.T) {
	r := newRuntimeSourcePageRig(t)
	r.seed(t, "device-A", "A1", 1)
	q := r.query(t, "rules")
	r.fault(&runtimeSourcePageFaultRepo{Repository: r.repo, changeAt: 2, change: func() {
		_, err := r.f.db.Exec(`DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
		require.NoError(t, err)
	}})
	response := r.get(t, q)
	var count int
	require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&count))
	require.Zero(t, count, "final list observation must not recreate a deleted workspace")
	require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
	require.NotContains(t, response.Body.String(), `"source_page"`)
}
