package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

type runtimeSourceGroups struct{ record *workspace.Record }

func (g runtimeSourceGroups) List(context.Context) (*workspace.WriteGroupListResult, error) {
	return &workspace.WriteGroupListResult{WorkspaceID: g.record.ID, WorkspaceRevision: g.record.DatabaseSetupRevision, Groups: []*workspace.WriteGroup{}}, nil
}

type runtimeSourceApplied struct{}

func (runtimeSourceApplied) ResolveAppliedConfigurationAt(context.Context, string, time.Time) (*workspace.WriteGroupAppliedSnapshot, error) {
	return nil, nil
}

type runtimeSourceFixture struct {
	h       *RuntimeHandler
	record  *workspace.Record
	query   url.Values
	point   *schema.Point
	rule    *schema.SourceRule
	link    *schema.SourceRuleLink
	mapping *schema.Mapping
	tag     *schema.Tag
	polling *schema.PollingGroup
	points  *point.MemoryRepository
	groups  *pollinggroup.MemoryRepository
}

func newRuntimeSourceFixture(t *testing.T) *runtimeSourceFixture {
	t.Helper()
	h, record := runtimeEvidenceScopeFixture(t)
	pgid, mid, tid := "polling-P", "mapping-P", "tag-P"
	f := &runtimeSourceFixture{h: h, record: record, points: point.NewMemoryRepository(), groups: pollinggroup.NewMemoryRepository()}
	f.point = &schema.Point{ID: "point-P", DeviceID: "device-1", Address: "40001", Function: "03", DataType: schema.DataTypeInt16, DataFormat: "BADC", PollingGroupID: &pgid, Enabled: true}
	f.polling = &schema.PollingGroup{ID: pgid, IntervalMs: 1000, Enabled: true}
	require.NoError(t, f.points.Create(t.Context(), f.point))
	require.NoError(t, f.groups.Create(t.Context(), f.polling))
	h.pointSvc = point.NewService(f.points, f.groups)
	h.groupSvc = pollinggroup.NewService(f.groups)
	f.rule = &schema.SourceRule{ID: "rule-P", DeviceID: "device-1", RevisionID: "R2", Enabled: true, StartAddress: "40001", DataType: schema.DataTypeInt16}
	f.link = &schema.SourceRuleLink{ID: "link-P", RuleID: f.rule.ID, PointID: f.point.ID, Address: "40001", MappingID: &mid, TagID: &tid}
	f.mapping = &schema.Mapping{ID: mid, PointID: f.point.ID, TagID: tid, Enabled: true, Status: schema.MappingStatusActive, TransformPipeline: `[{"type":"scale","order":0,"params":{"scale":0.1,"offset":0}},{"type":"cast","order":1,"params":{"target_type":"float32"}}]`}
	f.tag = &schema.Tag{ID: tid, Key: "meter.voltage", DataType: schema.DataTypeFloat32, Status: schema.TagStatusActive}
	h.sourceRuleReader = runtimeSetupRuleReaderStub{rules: []*schema.SourceRule{f.rule}, links: map[string][]*schema.SourceRuleLink{f.rule.ID: {f.link}}}
	h.mappingReader = runtimeSetupMappingReaderStub{records: map[string]*schema.Mapping{mid: f.mapping}}
	h.tagReader = runtimeSetupTagReaderStub{records: map[string]*schema.Tag{tid: f.tag}}
	h.WithRuleEvidenceServices(runtimeSourceGroups{record}, runtimeSourceApplied{})
	f.query = url.Values{"workspace_id": {record.ID}, "workspace_revision": {record.DatabaseSetupRevision}, "device_id": {"device-1"}, "rule_id": {f.rule.ID}, "rule_revision": {"R2"}, "point_id": {f.point.ID}, "mapping_id": {mid}}
	return f
}
func (f *runtimeSourceFixture) source(t *testing.T) *runtimeview.SourceView {
	t.Helper()
	rec := runtimeEvidenceScopeRequest(t, f.h, "?"+f.query.Encode())
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var body struct {
		Data struct {
			Evidence runtimeview.RuleEvidenceView `json:"rule_evidence"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.NotNil(t, body.Data.Evidence.Source)
	return body.Data.Evidence.Source
}
func TestRuntimeSourceEvidenceSavedR2ScaleDoesNotBecomeEffectiveR1(t *testing.T) {
	f := newRuntimeSourceFixture(t)
	source := f.source(t)
	require.Equal(t, "saved-only", source.State)
	require.Equal(t, "unavailable", source.Proof)
	require.Nil(t, source.Effective)
	require.Nil(t, source.EffectiveRevision)
	require.Equal(t, "R2", *source.SavedRevision)
	require.Equal(t, "rule-P", source.RuleID)
	require.Equal(t, "point-P", source.PointID)
	require.Equal(t, "mapping-P", *source.MappingID)
	require.Equal(t, "tag-P", *source.TagID)
	require.Equal(t, int64(1000), source.Saved.Polling.IntervalMS)
	require.Equal(t, "polling-group", source.Saved.Polling.Source)
	require.Equal(t, "BADC", *source.Saved.Transforms[0].Order)
	require.Equal(t, 0.1, *source.Saved.Transforms[1].Multiplier)
	require.Equal(t, "float32", *source.Saved.Transforms[2].TargetType)
	require.Equal(t, "float32", *source.Saved.ResultType)
	require.NotEmpty(t, *source.Saved.MappingRevision)
}
func TestRuntimeSourceEvidenceStrictSelectionAndNoImplicitFirstMapping(t *testing.T) {
	for _, suffix := range []string{"rule_id=R", "mapping_id=M", "rule_id=R&rule_revision=R2&point_id=P", "workspace_id=W&rule_id=R&rule_revision=R2&point_id=P"} {
		t.Run(suffix, func(t *testing.T) {
			f := newRuntimeSourceFixture(t)
			rec := runtimeEvidenceScopeRequest(t, f.h, "?"+suffix)
			require.Equal(t, 400, rec.Code, rec.Body.String())
		})
	}
	for _, mutation := range []func(url.Values){func(q url.Values) { q.Del("rule_revision") }, func(q url.Values) { q.Add("point_id", "other") }, func(q url.Values) { q.Set("mapping_id", "") }, func(q url.Values) { q.Set("rule_id", "/private/source") }, func(q url.Values) { q.Del("rule_id"); q.Del("point_id"); q.Del("rule_revision") }} {
		f := newRuntimeSourceFixture(t)
		mutation(f.query)
		rec := runtimeEvidenceScopeRequest(t, f.h, "?"+f.query.Encode())
		require.Equal(t, 400, rec.Code, rec.Body.String())
	}
	f := newRuntimeSourceFixture(t)
	f.query.Del("mapping_id")
	s := f.source(t)
	require.Nil(t, s.MappingID)
	require.Nil(t, s.TagID)
	require.Nil(t, s.Saved.MappingRevision)
	require.Nil(t, s.Saved.ResultType)
	require.Len(t, s.Saved.Transforms, 1)
	require.Equal(t, "byte-order", s.Saved.Transforms[0].Kind)
}
func TestRuntimeSourceEvidenceRejectsForeignAndMismatchedBinding(t *testing.T) {
	cases := map[string]func(*runtimeSourceFixture){
		"foreign-rule": func(f *runtimeSourceFixture) { f.rule.DeviceID = "device-other" },
		"foreign-point": func(f *runtimeSourceFixture) {
			f.point.DeviceID = "device-other"
			require.NoError(t, f.points.Update(t.Context(), f.point))
		},
		"wrong-revision": func(f *runtimeSourceFixture) { f.query.Set("rule_revision", "R1") },
		"wrong-link":     func(f *runtimeSourceFixture) { f.link.PointID = "point-other" },
		"wrong-mapping":  func(f *runtimeSourceFixture) { f.mapping.PointID = "point-other" },
		"wrong-tag":      func(f *runtimeSourceFixture) { f.mapping.TagID = "tag-other" },
		"wrong-address":  func(f *runtimeSourceFixture) { f.link.Address = "40002" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			f := newRuntimeSourceFixture(t)
			mutate(f)
			rec := runtimeEvidenceScopeRequest(t, f.h, "?"+f.query.Encode())
			require.Contains(t, []int{403, 409, 503}, rec.Code, rec.Body.String())
			require.NotContains(t, rec.Body.String(), `"source"`)
		})
	}
}
func TestRuntimeSourceEvidenceOrdinaryContextCarriesObservedRevision(t *testing.T) {
	h, record := runtimeEvidenceScopeFixture(t)
	rec := runtimeEvidenceScopeRequest(t, h, "")
	require.Equal(t, 200, rec.Code)
	require.Contains(t, rec.Body.String(), `"workspace_revision":"`+record.DatabaseSetupRevision+`"`)
}

type runtimeSourceMutatingRules struct {
	runtimeWorkspaceSourceRuleReader
	calls  int
	mutate func()
}

func (r *runtimeSourceMutatingRules) ListByDeviceIDs(ctx context.Context, ids []string) ([]*schema.SourceRule, error) {
	r.calls++
	if r.calls == 2 {
		r.mutate()
	}
	return r.runtimeWorkspaceSourceRuleReader.ListByDeviceIDs(ctx, ids)
}
func TestRuntimeSourceEvidenceChangedPointPollingMappingOrTagIsStale(t *testing.T) {
	cases := map[string]func(*runtimeSourceFixture){
		"point-config": func(f *runtimeSourceFixture) {
			f.point.DataFormat = "DCBA"
			require.NoError(t, f.points.Update(t.Context(), f.point))
		},
		"polling-interval": func(f *runtimeSourceFixture) {
			f.polling.IntervalMs = 2000
			require.NoError(t, f.groups.Update(t.Context(), f.polling))
		},
		"polling-disabled": func(f *runtimeSourceFixture) {
			f.polling.Enabled = false
			require.NoError(t, f.groups.Update(t.Context(), f.polling))
		},
		"mapping-pipeline": func(f *runtimeSourceFixture) { f.mapping.TransformPipeline = "[]" },
		"tag-type":         func(f *runtimeSourceFixture) { f.tag.DataType = schema.DataTypeFloat64 },
		"rule-revision":    func(f *runtimeSourceFixture) { f.rule.RevisionID = "R3" },
		"link-address":     func(f *runtimeSourceFixture) { f.link.Address = "40002" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			f := newRuntimeSourceFixture(t)
			f.h.sourceRuleReader = &runtimeSourceMutatingRules{runtimeWorkspaceSourceRuleReader: f.h.sourceRuleReader, mutate: func() { mutate(f) }}
			rec := runtimeEvidenceScopeRequest(t, f.h, "?"+f.query.Encode())
			require.Equal(t, 409, rec.Code, rec.Body.String())
			require.Contains(t, rec.Body.String(), "stale-reference")
			require.NotContains(t, rec.Body.String(), `"source"`)
		})
	}
}
func TestRuntimeSourceEvidenceTelemetryChangesDoNotInvalidateConfiguration(t *testing.T) {
	f := newRuntimeSourceFixture(t)
	f.h.sourceRuleReader = &runtimeSourceMutatingRules{runtimeWorkspaceSourceRuleReader: f.h.sourceRuleReader, mutate: func() {
		require.NoError(t, f.points.UpdateReadResult(t.Context(), f.point.ID, 123, "driver /private/fixture-secret"))
	}}
	s := f.source(t)
	require.Equal(t, "saved-only", s.State)
	raw, _ := json.Marshal(s)
	require.NotContains(t, string(raw), "fixture-secret")
}

func TestRuntimeSourceEvidenceMissingReadersAndRecordsAreUnavailable(t *testing.T) {
	cases := map[string]func(*runtimeSourceFixture){
		"source-reader":   func(f *runtimeSourceFixture) { f.h.sourceRuleReader = nil },
		"point-reader":    func(f *runtimeSourceFixture) { f.h.pointSvc = nil },
		"polling-reader":  func(f *runtimeSourceFixture) { f.h.groupSvc = nil },
		"mapping-reader":  func(f *runtimeSourceFixture) { f.h.mappingReader = nil },
		"tag-reader":      func(f *runtimeSourceFixture) { f.h.tagReader = nil },
		"point-deleted":   func(f *runtimeSourceFixture) { require.NoError(t, f.points.Delete(t.Context(), f.point.ID)) },
		"polling-deleted": func(f *runtimeSourceFixture) { require.NoError(t, f.groups.Delete(t.Context(), f.polling.ID)) },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			f := newRuntimeSourceFixture(t)
			mutate(f)
			rec := runtimeEvidenceScopeRequest(t, f.h, "?"+f.query.Encode())
			require.Equal(t, 503, rec.Code, rec.Body.String())
			require.Contains(t, rec.Body.String(), `"code":"unavailable"`)
			require.NotContains(t, rec.Body.String(), `"source"`)
		})
	}
}
func TestRuntimeSourceEvidenceNeverBootstraps(t *testing.T) {
	repo := workspace.NewMemoryRepository()
	f := newRuntimeSourceFixture(t)
	f.h.workspaceSvc = workspace.NewService(repo)
	rec := runtimeEvidenceScopeRequest(t, f.h, "?"+f.query.Encode())
	_, err := repo.Get(t.Context())
	require.ErrorIs(t, err, workspace.ErrNotFound)
	require.Equal(t, 503, rec.Code)
}
func TestRuntimeSourceEvidenceUnknownWorkspaceRevisionStaysMissing(t *testing.T) {
	f := newRuntimeSourceFixture(t)
	record, err := f.h.workspaceSvc.ReadSnapshot(t.Context())
	require.NoError(t, err)
	record.DatabaseSetupRevision = ""
	repo := workspace.NewMemoryRepository()
	require.NoError(t, repo.Save(t.Context(), record))
	f.h.workspaceSvc = workspace.NewService(repo)
	rec := runtimeEvidenceScopeRequest(t, f.h, "")
	require.Equal(t, 200, rec.Code)
	require.NotContains(t, rec.Body.String(), `"workspace_revision"`)
	record.DatabaseSetupRevision = "/private/fixture-secret"
	require.NoError(t, repo.Save(t.Context(), record))
	rec = runtimeEvidenceScopeRequest(t, f.h, "")
	require.Equal(t, 200, rec.Code)
	require.NotContains(t, rec.Body.String(), `"workspace_revision"`)
	require.NotContains(t, rec.Body.String(), "fixture-secret")
}

type runtimeSourceChangingTagReader struct {
	runtimeWorkspaceTagReader
	mutate func()
}

func (r runtimeSourceChangingTagReader) GetByID(ctx context.Context, id string) (*schema.Tag, error) {
	r.mutate()
	return r.runtimeWorkspaceTagReader.GetByID(ctx, id)
}
func TestRuntimeSourceEvidenceDetachesMappingBeforeDependentTagRead(t *testing.T) {
	f := newRuntimeSourceFixture(t)
	f.tag.Labels = `{"source":"source-rule","source_rule_id":"rule-P","source_rule_address":"40001"}`
	f.h.tagReader = runtimeSourceChangingTagReader{runtimeWorkspaceTagReader: f.h.tagReader, mutate: func() { f.mapping.PointID = "point-foreign" }}
	records, err := f.h.readRuntimeSourceRecords(t.Context(), runtimeview.Scope{WorkspaceID: f.record.ID, WorkspaceRevision: f.record.DatabaseSetupRevision, DeviceID: "device-1"}, runtimeSourceSelection{RuleID: "rule-P", RuleRevision: "R2", PointID: "point-P", MappingID: "mapping-P"})
	require.NoError(t, err)
	require.Equal(t, "point-P", records.Mapping.PointID, "dependent reader must not alter already observed mapping identity")
}

type runtimeSourceFailingRuleReader struct {
	runtimeWorkspaceSourceRuleReader
	listCalls, linkCalls int
	failLinks            bool
}

func (r *runtimeSourceFailingRuleReader) ListByDeviceIDs(ctx context.Context, ids []string) ([]*schema.SourceRule, error) {
	r.listCalls++
	if !r.failLinks && r.listCalls == 2 {
		return nil, fmt.Errorf("storage password=fixture-secret")
	}
	return r.runtimeWorkspaceSourceRuleReader.ListByDeviceIDs(ctx, ids)
}
func (r *runtimeSourceFailingRuleReader) ListLinks(ctx context.Context, id string) ([]*schema.SourceRuleLink, error) {
	r.linkCalls++
	if r.failLinks && r.linkCalls == 2 {
		return nil, fmt.Errorf("storage password=fixture-secret")
	}
	return r.runtimeWorkspaceSourceRuleReader.ListLinks(ctx, id)
}

type runtimeSourceFailingPointReader struct {
	point.Repository
	calls int
}

func (r *runtimeSourceFailingPointReader) GetByID(ctx context.Context, id string) (*schema.Point, error) {
	r.calls++
	if r.calls == 2 {
		return nil, fmt.Errorf("storage password=fixture-secret")
	}
	return r.Repository.GetByID(ctx, id)
}

type runtimeSourceFailingPollingReader struct {
	pollinggroup.Repository
	calls int
}

func (r *runtimeSourceFailingPollingReader) GetByID(ctx context.Context, id string) (*schema.PollingGroup, error) {
	r.calls++
	if r.calls == 2 {
		return nil, fmt.Errorf("storage password=fixture-secret")
	}
	return r.Repository.GetByID(ctx, id)
}

type runtimeSourceFailingMappingReader struct {
	runtimeWorkspaceMappingReader
	calls int
}

func (r *runtimeSourceFailingMappingReader) GetByID(ctx context.Context, id string) (*schema.Mapping, error) {
	r.calls++
	if r.calls == 2 {
		return nil, fmt.Errorf("storage password=fixture-secret")
	}
	return r.runtimeWorkspaceMappingReader.GetByID(ctx, id)
}

type runtimeSourceFailingTagReader struct {
	runtimeWorkspaceTagReader
	calls int
}

func (r *runtimeSourceFailingTagReader) GetByID(ctx context.Context, id string) (*schema.Tag, error) {
	r.calls++
	if r.calls == 2 {
		return nil, fmt.Errorf("storage password=fixture-secret")
	}
	return r.runtimeWorkspaceTagReader.GetByID(ctx, id)
}

func TestRuntimeSourceEvidenceSecondObservationProviderFailureRemainsUnavailable(t *testing.T) {
	cases := map[string]func(*runtimeSourceFixture){
		"rules": func(f *runtimeSourceFixture) {
			f.h.sourceRuleReader = &runtimeSourceFailingRuleReader{runtimeWorkspaceSourceRuleReader: f.h.sourceRuleReader}
		},
		"links": func(f *runtimeSourceFixture) {
			f.h.sourceRuleReader = &runtimeSourceFailingRuleReader{runtimeWorkspaceSourceRuleReader: f.h.sourceRuleReader, failLinks: true}
		},
		"point": func(f *runtimeSourceFixture) {
			f.h.pointSvc = point.NewService(&runtimeSourceFailingPointReader{Repository: f.points}, f.groups)
		},
		"polling": func(f *runtimeSourceFixture) {
			f.h.groupSvc = pollinggroup.NewService(&runtimeSourceFailingPollingReader{Repository: f.groups})
		},
		"mapping": func(f *runtimeSourceFixture) {
			f.h.mappingReader = &runtimeSourceFailingMappingReader{runtimeWorkspaceMappingReader: f.h.mappingReader}
		},
		"tag": func(f *runtimeSourceFixture) {
			f.h.tagReader = &runtimeSourceFailingTagReader{runtimeWorkspaceTagReader: f.h.tagReader}
		},
	}
	for name, inject := range cases {
		t.Run(name, func(t *testing.T) {
			f := newRuntimeSourceFixture(t)
			inject(f)
			rec := runtimeEvidenceScopeRequest(t, f.h, "?"+f.query.Encode())
			require.Equal(t, http.StatusServiceUnavailable, rec.Code, rec.Body.String())
			require.Contains(t, rec.Body.String(), `"code":"unavailable"`)
			require.NotContains(t, rec.Body.String(), "fixture-secret")
			require.NotContains(t, rec.Body.String(), `"rule_evidence"`)
		})
	}
}
func TestRuntimeSourceEvidenceSecondObservationForeignPointRemainsForeign(t *testing.T) {
	f := newRuntimeSourceFixture(t)
	f.h.sourceRuleReader = &runtimeSourceMutatingRules{runtimeWorkspaceSourceRuleReader: f.h.sourceRuleReader, mutate: func() { f.point.DeviceID = "device-foreign"; require.NoError(t, f.points.Update(t.Context(), f.point)) }}
	rec := runtimeEvidenceScopeRequest(t, f.h, "?"+f.query.Encode())
	require.Equal(t, http.StatusForbidden, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), `"code":"foreign-scope"`)
	require.NotContains(t, rec.Body.String(), `"rule_evidence"`)
}
