package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

type runtimeSourcePageTestView struct {
	Kind     string           `json:"kind"`
	Query    string           `json:"query"`
	Limit    int              `json:"limit"`
	Revision string           `json:"list_revision"`
	Rules    []map[string]any `json:"rules"`
	Mappings []map[string]any `json:"mappings"`
	Next     *string          `json:"next_cursor"`
}

func runtimeSourcePageFixture(t *testing.T, count int) (*runtimeSourceFixture, url.Values) {
	t.Helper()
	f := newRuntimeSourceFixture(t)
	rules := make([]*schema.SourceRule, 0, count)
	for i := 1; i <= count; i++ {
		rules = append(rules, &schema.SourceRule{ID: fmt.Sprintf("A%03d", i), DeviceID: "device-1", RevisionID: fmt.Sprintf("R%03d", i), StartAddress: fmt.Sprint(40000 + i), Count: 1, DataType: schema.DataTypeInt16, NamingPrefix: fmt.Sprintf("Meter %03d", i), Enabled: true})
	}
	f.h.sourceRuleReader = runtimeSetupRuleReaderStub{rules: rules}
	f.h.ruleGroups = nil
	f.h.ruleApplied = nil
	q := url.Values{"workspace_id": {f.record.ID}, "workspace_revision": {f.record.DatabaseSetupRevision}, "device_id": {"device-1"}, "source_list": {"rules"}}
	return f, q
}
func runtimeSourcePageGet(t *testing.T, h *RuntimeHandler, q url.Values) runtimeSourcePageTestView {
	t.Helper()
	rec := runtimeEvidenceScopeRequest(t, h, "?"+q.Encode())
	require.Equal(t, 200, rec.Code, rec.Body.String())
	var envelope struct {
		Data struct {
			Page *runtimeSourcePageTestView `json:"source_page"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &envelope))
	require.NotNil(t, envelope.Data.Page, rec.Body.String())
	return *envelope.Data.Page
}
func TestRuntimeSourcePageDefault50AndClamp100(t *testing.T) {
	f, q := runtimeSourcePageFixture(t, 120)
	page := runtimeSourcePageGet(t, f.h, q)
	require.Equal(t, 50, page.Limit)
	require.Len(t, page.Rules, 50)
	require.Empty(t, page.Mappings)
	require.NotNil(t, page.Next)
	require.Equal(t, "A001", page.Rules[0]["id"])
	q.Set("limit", "500")
	page = runtimeSourcePageGet(t, f.h, q)
	require.Equal(t, 100, page.Limit)
	require.Len(t, page.Rules, 100)
	require.NotNil(t, page.Next)
	q.Set("cursor", *page.Next)
	last := runtimeSourcePageGet(t, f.h, q)
	require.Len(t, last.Rules, 20)
	require.Equal(t, "A101", last.Rules[0]["id"])
	require.Nil(t, last.Next)
	require.Equal(t, page.Revision, last.Revision)
}
func TestRuntimeSourcePageRejectsMalformedForeignAndMixedModeRequests(t *testing.T) {
	cases := map[string]func(url.Values){
		"no-scope":                 func(q url.Values) { q.Del("workspace_id"); q.Del("workspace_revision"); q.Del("device_id") },
		"bad-kind":                 func(q url.Values) { q.Set("source_list", "other") },
		"duplicate-kind":           func(q url.Values) { q.Add("source_list", "rules") },
		"detail":                   func(q url.Values) { q.Set("point_id", "point-P") },
		"mapping-detail":           func(q url.Values) { q.Set("mapping_id", "mapping-P") },
		"rule-filter-not-mappings": func(q url.Values) { q.Set("rule_id", "A001"); q.Set("rule_revision", "R001") },
		"unpaired-filter":          func(q url.Values) { q.Set("source_list", "mappings"); q.Set("rule_id", "A001") },
		"negative":                 func(q url.Values) { q.Set("limit", "-1") }, "zero": func(q url.Values) { q.Set("limit", "0") }, "float": func(q url.Values) { q.Set("limit", "1.5") }, "duplicate-limit": func(q url.Values) { q.Add("limit", "1"); q.Add("limit", "2") },
		"unsafe-query": func(q url.Values) { q.Set("q", "password=fixture-secret") }, "empty-cursor": func(q url.Values) { q.Set("cursor", "") },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			f, q := runtimeSourcePageFixture(t, 1)
			mutate(q)
			rec := runtimeEvidenceScopeRequest(t, f.h, "?"+q.Encode())
			require.Equal(t, 400, rec.Code, rec.Body.String())
			require.NotContains(t, rec.Body.String(), "fixture-secret")
		})
	}
	f, q := runtimeSourcePageFixture(t, 1)
	q.Set("device_id", "foreign-device")
	q.Set("cursor", "bad")
	rec := runtimeEvidenceScopeRequest(t, f.h, "?"+q.Encode())
	require.Equal(t, 403, rec.Code, rec.Body.String())
}
func TestRuntimeSourcePageSearchAndFilterAreRevisionBound(t *testing.T) {
	f, q := runtimeSourcePageFixture(t, 120)
	q.Set("q", "  Meter 11  ")
	q.Set("limit", "3")
	page := runtimeSourcePageGet(t, f.h, q)
	require.Equal(t, "Meter 11", page.Query)
	require.Len(t, page.Rules, 3)
	require.NotNil(t, page.Next)
	require.Equal(t, "A110", page.Rules[0]["id"])
	q.Set("cursor", *page.Next)
	q.Set("q", "Meter 10")
	rec := runtimeEvidenceScopeRequest(t, f.h, "?"+q.Encode())
	require.Equal(t, 409, rec.Code, rec.Body.String())
}
func TestRuntimeSourcePageMappingRowsCarryParentRevisionWithoutDetailCalls(t *testing.T) {
	f := newRuntimeSourceFixture(t)
	f.h.ruleGroups = nil
	f.h.ruleApplied = nil
	q := url.Values{"workspace_id": {f.record.ID}, "workspace_revision": {f.record.DatabaseSetupRevision}, "device_id": {"device-1"}, "source_list": {"mappings"}, "rule_id": {"rule-P"}, "rule_revision": {"R2"}}
	page := runtimeSourcePageGet(t, f.h, q)
	require.Len(t, page.Mappings, 1)
	require.Empty(t, page.Rules)
	require.Equal(t, "R2", page.Mappings[0]["rule_revision"])
	require.Equal(t, "mapping-P", page.Mappings[0]["id"])
	q.Set("rule_revision", "R1")
	rec := runtimeEvidenceScopeRequest(t, f.h, "?"+q.Encode())
	require.Equal(t, 409, rec.Code, rec.Body.String())
}

type runtimeSourcePageSharedPoints struct {
	point.Repository
	records []*schema.Point
}

func (r runtimeSourcePageSharedPoints) ListByDevice(context.Context, string) ([]*schema.Point, error) {
	return r.records, nil
}
func TestRuntimeSourcePageDoesNotMutateBorrowedPointSlice(t *testing.T) {
	f := newRuntimeSourceFixture(t)
	f.point.LastError = "private telemetry"
	shared := []*schema.Point{f.point}
	f.h.pointSvc = point.NewService(runtimeSourcePageSharedPoints{Repository: f.points, records: shared}, f.groups)
	q := url.Values{"workspace_id": {f.record.ID}, "workspace_revision": {f.record.DatabaseSetupRevision}, "device_id": {"device-1"}, "source_list": {"mappings"}}
	runtimeSourcePageGet(t, f.h, q)
	require.Same(t, f.point, shared[0])
	require.Equal(t, "private telemetry", shared[0].LastError)
}
func TestRuntimeSourcePageReadFailurePreservesUnavailable(t *testing.T) {
	for _, links := range []bool{false, true} {
		f := newRuntimeSourceFixture(t)
		f.h.sourceRuleReader = &runtimeSourceFailingRuleReader{runtimeWorkspaceSourceRuleReader: f.h.sourceRuleReader, failLinks: links}
		q := url.Values{"workspace_id": {f.record.ID}, "workspace_revision": {f.record.DatabaseSetupRevision}, "device_id": {"device-1"}, "source_list": {"mappings"}}
		rec := runtimeEvidenceScopeRequest(t, f.h, "?"+q.Encode())
		require.Equal(t, 503, rec.Code, rec.Body.String())
		require.NotContains(t, rec.Body.String(), "fixture-secret")
	}
}

func TestRuntimeSourcePageExplicitBindingsMatchDetailSelection(t *testing.T) {
	cases := map[string]func(*runtimeSourceFixture) string{
		"missing-explicit-mapping": func(f *runtimeSourceFixture) string {
			id := "mapping-missing"
			f.link.MappingID = &id
			return f.mapping.ID
		},
		"explicit-mapping-wrong-point": func(f *runtimeSourceFixture) string {
			alternative := *f.mapping
			alternative.ID = "mapping-alternative"
			f.mapping.PointID = "point-other"
			f.h.mappingReader = runtimeSetupMappingReaderStub{records: map[string]*schema.Mapping{f.mapping.ID: f.mapping, alternative.ID: &alternative}}
			return alternative.ID
		},
		"explicit-tag-mismatch": func(f *runtimeSourceFixture) string {
			id := "tag-other"
			f.link.TagID = &id
			return f.mapping.ID
		},
		"explicit-tag-foreign": func(f *runtimeSourceFixture) string {
			id := "tag-foreign"
			f.link.TagID = &id
			f.h.tagReader = runtimeSetupTagReaderStub{records: map[string]*schema.Tag{f.tag.ID: f.tag, id: {ID: id, Key: "foreign", Labels: `{"source":"source-rule","source_rule_id":"rule-foreign","source_rule_address":"40001"}`}}}
			return f.mapping.ID
		},
		"actual-tag-foreign": func(f *runtimeSourceFixture) string {
			f.tag.Labels = `{"source":"source-rule","source_rule_id":"rule-foreign","source_rule_address":"40001"}`
			return f.mapping.ID
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			f := newRuntimeSourceFixture(t)
			f.tag.Labels = `{"source":"source-rule","source_rule_id":"rule-P","source_rule_address":"40001"}`
			f.query.Set("mapping_id", mutate(f))
			// The existing detail reader already rejects the misleading row's identity.
			detail := runtimeEvidenceScopeRequest(t, f.h, "?"+f.query.Encode())
			require.Equal(t, 409, detail.Code, detail.Body.String())
			q := url.Values{"workspace_id": {f.record.ID}, "workspace_revision": {f.record.DatabaseSetupRevision}, "device_id": {"device-1"}, "source_list": {"mappings"}}
			page := runtimeEvidenceScopeRequest(t, f.h, "?"+q.Encode())
			require.Equal(t, 409, page.Code, page.Body.String())
			require.Contains(t, page.Body.String(), `"code":"stale-reference"`)
			require.NotContains(t, page.Body.String(), `"source_page"`)
		})
	}
	t.Run("valid-explicit-binding-opens-detail", func(t *testing.T) {
		f := newRuntimeSourceFixture(t)
		q := url.Values{"workspace_id": {f.record.ID}, "workspace_revision": {f.record.DatabaseSetupRevision}, "device_id": {"device-1"}, "source_list": {"mappings"}}
		page := runtimeSourcePageGet(t, f.h, q)
		require.Len(t, page.Mappings, 1)
		row := page.Mappings[0]
		f.query.Set("rule_id", row["rule_id"].(string))
		f.query.Set("rule_revision", row["rule_revision"].(string))
		f.query.Set("point_id", row["point_id"].(string))
		f.query.Set("mapping_id", row["id"].(string))
		source := f.source(t)
		require.Equal(t, row["id"], *source.MappingID)
		require.Equal(t, row["tag_id"], *source.TagID)
	})
}

func TestRuntimeSourcePageExplicitBindingProviderFailureRemainsUnavailable(t *testing.T) {
	for _, provider := range []string{"mapping", "tag"} {
		t.Run(provider, func(t *testing.T) {
			f := newRuntimeSourceFixture(t)
			if provider == "mapping" {
				f.h.mappingReader = &runtimeSourceFailingMappingReader{runtimeWorkspaceMappingReader: f.h.mappingReader}
			} else {
				f.h.tagReader = &runtimeSourceFailingTagReader{runtimeWorkspaceTagReader: f.h.tagReader}
			}
			q := url.Values{"workspace_id": {f.record.ID}, "workspace_revision": {f.record.DatabaseSetupRevision}, "device_id": {"device-1"}, "source_list": {"mappings"}}
			page := runtimeEvidenceScopeRequest(t, f.h, "?"+q.Encode())
			require.Equal(t, 503, page.Code, page.Body.String())
			require.Contains(t, page.Body.String(), `"code":"unavailable"`)
			require.NotContains(t, page.Body.String(), "fixture-secret")
			require.NotContains(t, page.Body.String(), `"source_page"`)
		})
	}
}
