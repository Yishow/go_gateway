package handlers

import (
	"context"
	"encoding/json"
	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type homeWorkspaceFunc func(context.Context) (*workspace.Record, error)

func (f homeWorkspaceFunc) ReadSnapshot(c context.Context) (*workspace.Record, error) { return f(c) }

type homeDeviceFunc func(context.Context, string) (*schema.Device, error)

func (f homeDeviceFunc) GetByID(c context.Context, id string) (*schema.Device, error) {
	return f(c, id)
}

type homePointFunc func(context.Context, point.ListFilter) ([]*schema.Point, error)

func (f homePointFunc) List(c context.Context, x point.ListFilter) ([]*schema.Point, error) {
	return f(c, x)
}

type homeTagFunc func(context.Context, string) (*schema.Tag, error)

func (f homeTagFunc) GetByID(c context.Context, id string) (*schema.Tag, error) { return f(c, id) }

type homeRuleStub struct {
	RecordingHomeRuleReader
	list  func(context.Context, []string) ([]*schema.SourceRule, error)
	links func(context.Context, string) ([]*schema.SourceRuleLink, error)
}

func (s *homeRuleStub) ListByDeviceIDs(c context.Context, ids []string) ([]*schema.SourceRule, error) {
	if s.list != nil {
		return s.list(c, ids)
	}
	return s.RecordingHomeRuleReader.ListByDeviceIDs(c, ids)
}
func (s *homeRuleStub) ListLinks(c context.Context, id string) ([]*schema.SourceRuleLink, error) {
	if s.links != nil {
		return s.links(c, id)
	}
	return s.RecordingHomeRuleReader.ListLinks(c, id)
}

type homeMappingStub struct {
	RecordingHomeMappingReader
	get  func(context.Context, string) (*schema.Mapping, error)
	list func(context.Context, mapping.ListFilter) ([]*schema.Mapping, error)
}

func (s *homeMappingStub) GetByID(c context.Context, id string) (*schema.Mapping, error) {
	if s.get != nil {
		return s.get(c, id)
	}
	return s.RecordingHomeMappingReader.GetByID(c, id)
}
func (s *homeMappingStub) List(c context.Context, f mapping.ListFilter) ([]*schema.Mapping, error) {
	if s.list != nil {
		return s.list(c, f)
	}
	return s.RecordingHomeMappingReader.List(c, f)
}

type homeRealFixture struct {
	old     *StudioV2WorkspaceMappingsHandler
	handler *StudioV2RecordingHomeSourcesHandler
	id      string
	mapping *schema.Mapping
	tag     *schema.Tag
}

func newHomeRealFixture(t *testing.T) homeRealFixture {
	t.Helper()
	ctx := t.Context()
	wr := workspace.NewMemoryRepository()
	require.NoError(t, wr.Save(ctx, &workspace.Record{ID: "W1", OrderedDeviceIDs: []string{"dev-A"}}))
	dr := device.NewMemoryRepository()
	seedWorkspaceMappingDevice(t, dr, "dev-A")
	pr := point.NewMemoryRepository()
	rr := sourcerule.NewMemoryRepository()
	tr := tag.NewMemoryRepository()
	mr := mapping.NewMemoryRepository()
	require.NoError(t, pr.Create(ctx, &schema.Point{ID: "P1", DeviceID: "dev-A", Address: "40001", Enabled: true, DataType: schema.DataTypeInt16, LastValue: common.Ptr("private-value")}))
	require.NoError(t, rr.Create(ctx, &schema.SourceRule{ID: "rule-A", DeviceID: "dev-A", RevisionID: "revision-1", Enabled: true, SkippedAddresses: "[]"}))
	savedTag := &schema.Tag{ID: "T1", Key: "line.a.saved", DisplayName: " Temperature ", DataType: schema.DataTypeFloat64, Labels: `{"source":"source-rule","source_rule_id":"rule-A","source_rule_address":"40001"}`}
	require.NoError(t, tr.Create(ctx, savedTag))
	pipeline, err := json.Marshal(buildWorkspaceMappingPipeline(schema.DataTypeInt16, schema.DataTypeFloat64, 2, 1))
	require.NoError(t, err)
	savedMapping := &schema.Mapping{ID: "M1", PointID: "P1", TagID: "T1", Enabled: true, TransformPipeline: string(pipeline)}
	require.NoError(t, mr.Create(ctx, savedMapping))
	require.NoError(t, rr.CreateLinks(ctx, []*schema.SourceRuleLink{{ID: "L1", RuleID: "rule-A", PointID: "P1", Address: "40001", MappingID: common.Ptr("M1"), TagID: common.Ptr("T1")}}))
	ps := point.NewService(pr, nil)
	ds := device.NewService(dr, nil)
	rs := sourcerule.NewService(rr, ds, ps, nil)
	ts := tag.NewService(tr)
	ms := mapping.NewService(mr)
	ws := workspace.NewService(wr)
	old := NewStudioV2WorkspaceMappingsHandler(ws, ds, rs, ps, ts, ms)
	return homeRealFixture{old, NewStudioV2RecordingHomeSourcesHandler(ws, ds, rs, ps, ms, ts), "W1", savedMapping, savedTag}
}
func (f homeRealFixture) snapshot(t *testing.T) string {
	t.Helper()
	ctx := t.Context()
	w, e := f.old.workspaceSvc.ReadSnapshot(ctx)
	require.NoError(t, e)
	rules, e := f.old.ruleSvc.ListByDeviceIDs(ctx, w.OrderedDeviceIDs)
	require.NoError(t, e)
	var links []*schema.SourceRuleLink
	devices := make([]*schema.Device, 0, len(w.OrderedDeviceIDs))
	for _, r := range rules {
		v, err := f.old.ruleSvc.ListLinks(ctx, r.ID)
		require.NoError(t, err)
		links = append(links, v...)
	}
	for _, id := range w.OrderedDeviceIDs {
		v, err := f.old.deviceSvc.GetByID(ctx, id)
		require.NoError(t, err)
		devices = append(devices, v)
	}
	p, e := f.old.pointSvc.List(ctx, point.ListFilter{})
	require.NoError(t, e)
	m, e := f.old.mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, e)
	tags, e := f.old.tagSvc.List(ctx, tag.ListFilter{})
	require.NoError(t, e)
	slices.SortFunc(p, func(a, b *schema.Point) int { return strings.Compare(a.ID, b.ID) })
	slices.SortFunc(m, func(a, b *schema.Mapping) int { return strings.Compare(a.ID, b.ID) })
	slices.SortFunc(tags, func(a, b *schema.Tag) int { return strings.Compare(a.ID, b.ID) })
	b, e := json.Marshal([]any{w, rules, links, devices, p, m, tags})
	require.NoError(t, e)
	return string(b)
}
func requireHomeUnknown(t *testing.T, f homeRealFixture) {
	t.Helper()
	out := homeSourcesRequest(f.handler, "?workspace_id="+f.id)
	require.Equal(t, 200, out.Code, out.Body.String())
	var body struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(out.Body.Bytes(), &body))
	require.Equal(t, "unknown", body.Data["state"])
	for _, key := range []string{"candidates", "excluded", "devices", "source_digest"} {
		require.NotContains(t, body.Data, key)
	}
}
