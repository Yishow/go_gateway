package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type runtimeSourceRig struct {
	f        managedGroupFixture
	services *DatalinkServices
	router   *gin.Engine
}

// Source rows and frozen group proof live in the existing disposable SQLite
// fixture. The recording group is created, schema-confirmed and applied through
// production routes; no effective descriptor is supplied by this test.
func newRuntimeSourceRig(t *testing.T) *runtimeSourceRig {
	t.Helper()
	f := newManagedDestinationFixture(t, "sqlite")
	f.groups.WithClock(func() time.Time { return time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC) })
	sources := []meanSource{{name: "source", tagType: "float32", role: "gauge"}}
	seedMeanSources(t, f, sources)
	ctx := t.Context()
	groupRepo := pollinggroup.NewSQLRepository(f.db)
	require.NoError(t, groupRepo.Create(ctx, &schema.PollingGroup{ID: "polling-P", Name: "Acquisition", IntervalMs: 1000, Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}))
	for _, statement := range []string{
		`UPDATE devices SET description='' WHERE id='device-A'`,
		`UPDATE points SET address='40005' WHERE id='point-A'`,
		`UPDATE points SET address='40001',function='FC03',data_format='CDAB',polling_group_id='polling-P' WHERE id='point-source'`,
		`UPDATE tags SET labels='{"source":"source-rule","source_rule_id":"rule-R","source_rule_address":"40001"}' WHERE id='tag-source'`,
		`UPDATE mappings SET transform_pipeline='[{"type":"scale","order":1,"params":{"scale":0.1,"offset":2}},{"type":"cast","order":2,"params":{"target_type":"float32"}}]' WHERE id='mapping-source'`,
	} {
		_, err := f.db.ExecContext(ctx, statement)
		require.NoError(t, err)
	}
	repo := sourcerule.NewSQLRepository(f.db)
	scale := 0.1
	require.NoError(t, repo.Create(ctx, &schema.SourceRule{ID: "rule-R", DeviceID: "device-A", StartAddress: "40001", Count: 1, DataType: schema.DataTypeFloat32, NamingPrefix: "Source", Enabled: true, Origin: "manual", RevisionID: "R2", ScaleMultiplier: &scale, DataFormat: "CDAB", CreatedAt: time.Now(), UpdatedAt: time.Now()}))
	tagID, mappingID := "tag-source", "mapping-source"
	require.NoError(t, repo.CreateLinks(ctx, []*schema.SourceRuleLink{{ID: "link-source", RuleID: "rule-R", Address: "40001", PointID: "point-source", TagID: &tagID, MappingID: &mappingID, CreatedAt: time.Now(), UpdatedAt: time.Now()}}))
	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, meanLongRequest(t, f, sources))
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	f.data, f.group = groupSaveData(t, decodeJSONBody(t, created))
	id := f.group["id"].(string)
	confirmLong(t, f, id, previewLong(t, f, id))
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	devices := device.NewService(device.NewSQLRepository(f.db), nil)
	points := point.NewService(point.NewSQLRepository(f.db), groupRepo)
	tags := tag.NewService(tag.NewSQLRepository(f.db))
	mappings := mapping.NewService(mapping.NewSQLRepository(f.db))
	rules := sourcerule.NewService(repo, devices, points, nil)
	// Install only the real runtime projection version marker. No scheduler is
	// started and no immutable R1 source facts exist, so R2 cannot become active.
	runtime, err := datalinkruntime.NewService(datalinkruntime.Config{Writer: storage.NewSQLiteWriter(f.db)})
	require.NoError(t, err)
	require.NoError(t, runtime.ApplyWorkspaceProjection(ctx, &workspace.RuntimeProjection{WorkspaceID: f.record.ID, Version: "R1", DeviceIDs: []string{"device-A"}}))
	services := &DatalinkServices{Workspace: f.workspace, WriteGroups: f.groups, WriteGroupDelivery: grouppipeline.New(grouppipeline.Dependencies{Groups: f.groups, Store: groupdelivery.NewStore(f.db)}, grouppipeline.Config{}), Device: devices, Point: points, PollingGroup: pollinggroup.NewService(groupRepo), SourceRule: rules, Tag: tags, Mapping: mappings, Runtime: runtime}
	return &runtimeSourceRig{f: f, services: services, router: NewRouter(services)}
}

func (r *runtimeSourceRig) query(t *testing.T) url.Values {
	t.Helper()
	record, err := r.f.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	return url.Values{"workspace_id": {record.ID}, "workspace_revision": {record.DatabaseSetupRevision}, "device_id": {"device-A"}, "rule_id": {"rule-R"}, "rule_revision": {"R2"}, "point_id": {"point-source"}, "mapping_id": {"mapping-source"}}
}
func (r *runtimeSourceRig) get(t *testing.T, query url.Values) *httptest.ResponseRecorder {
	t.Helper()
	return performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+query.Encode(), nil)
}
func runtimeSourceView(t *testing.T, response *httptest.ResponseRecorder) runtimeview.RuleEvidenceView {
	t.Helper()
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data struct {
			Evidence runtimeview.RuleEvidenceView `json:"rule_evidence"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.NotNil(t, body.Data.Evidence.Source, response.Body.String())
	return body.Data.Evidence
}

// This wrapper only times an actual SQL configuration change after the source
// read; successful group/proof reads still use their production services.
type runtimeSourceChangingGroups struct {
	groups *workspace.WriteGroupService
	change func()
	calls  int
	after  int
}

func (g *runtimeSourceChangingGroups) List(ctx context.Context) (*workspace.WriteGroupListResult, error) {
	result, err := g.groups.List(ctx)
	g.calls++
	if err == nil && g.calls == g.after {
		g.change()
	}
	return result, err
}
func (r *runtimeSourceRig) withChange(change func(), afterReads ...int) {
	after := 1
	if len(afterReads) > 0 {
		after = afterReads[0]
	}
	s := r.services
	applied := s.WriteGroupDelivery.(handlers.RuntimeAppliedConfigurationReader)
	h := handlers.NewRuntimeHandler(s.Device, s.Point, s.PollingGroup, nil, s.Runtime, s.Workspace).WithSetupContextServices(s.SourceRule, s.Tag, s.Mapping, nil, nil).WithRuleEvidenceServices(&runtimeSourceChangingGroups{groups: r.f.groups, change: change, after: after}, applied)
	r.router = gin.New()
	r.router.GET(runtimeRuleContextPath, h.WorkspaceContext)
}
