package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const studioLocatorPath = "/api/v1/datalink/studio-v2/workspace"

func locatorRequest(t *testing.T, r *runtimeSourceRig, step string, extra url.Values) *httptest.ResponseRecorder {
	t.Helper()
	q := url.Values{"entry": {"locator"}, "workspace_id": {r.f.record.ID}, "step": {step}}
	for k, v := range extra {
		q[k] = v
	}
	response := httptest.NewRecorder()
	r.router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, studioLocatorPath+"?"+q.Encode(), nil))
	return response
}

func locatorData(t *testing.T, response *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	body := decodeJSONBody(t, response)
	data, ok := body["data"].(map[string]any)
	require.True(t, ok, response.Body.String())
	return data
}

func locatorReaders(r *runtimeSourceRig) handlers.StudioV2LocatorReaders {
	s := r.services
	return handlers.StudioV2LocatorReaders{Connectors: s.DBTarget, Targets: s.DBMapping, Devices: s.Device, Rules: s.SourceRule, Points: s.Point, Mappings: s.Mapping, Tags: s.Tag, Groups: s.WriteGroups}
}
func useLocatorReaders(r *runtimeSourceRig, readers handlers.StudioV2LocatorReaders) {
	h := handlers.NewStudioV2WorkspaceHandler(r.f.workspace).WithLocatorEntryServices(readers)
	r.router = gin.New()
	r.router.GET(studioLocatorPath, h.Get)
}

type locatorDeviceRead func(context.Context, string) (*schema.Device, error)

func (f locatorDeviceRead) GetByID(ctx context.Context, id string) (*schema.Device, error) {
	return f(ctx, id)
}

type locatorRuleReads struct {
	*sourcerule.Service
	get func(context.Context, string) (*schema.SourceRule, error)
}

func (f locatorRuleReads) GetByID(ctx context.Context, id string) (*schema.SourceRule, error) {
	return f.get(ctx, id)
}
func locatorExec(t *testing.T, r *runtimeSourceRig, statement string, args ...any) {
	t.Helper()
	_, err := r.f.db.ExecContext(t.Context(), statement, args...)
	require.NoError(t, err)
}
func locatorTotal(t *testing.T, r *runtimeSourceRig) int64 {
	t.Helper()
	var n int64
	require.NoError(t, r.f.db.QueryRowContext(t.Context(), `SELECT total_changes()`).Scan(&n))
	return n
}
func locatorEntryData(t *testing.T, response *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	return locatorData(t, response)["locator_entry"].(map[string]any)
}
func locatorEditorQuery(t *testing.T, r *runtimeSourceRig) url.Values {
	t.Helper()
	record, err := r.f.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	group, err := r.services.WriteGroups.Get(t.Context(), r.f.group["id"].(string))
	require.NoError(t, err)
	return url.Values{"group_id": {group.Group.ID}, "editor": {"1"}, "expected_group_revision": {group.Group.Revision}, "expected_workspace_revision": {record.DatabaseSetupRevision}}
}
func wireLocatorDatabase(r *runtimeSourceRig) {
	r.services.DBTarget = dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))
	r.services.DBMapping = dbtarget.NewMappingService(dbtarget.NewSQLTargetMappingRepository(r.f.db), dbtarget.NewSQLConnectorRepository(r.f.db), nil)
	r.router = NewRouter(r.services)
}

type locatorGroupRead struct {
	*workspace.WriteGroupService
	get func(context.Context, string) (*workspace.WriteGroupSaveResult, error)
}

func (f locatorGroupRead) Get(ctx context.Context, id string) (*workspace.WriteGroupSaveResult, error) {
	return f.get(ctx, id)
}

type locatorConnectorRead func(context.Context, string) (*schema.DatabaseConnector, error)

func (f locatorConnectorRead) GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	return f(ctx, id)
}

func newLocatorMigratedRig(t *testing.T) *runtimeSourceRig {
	t.Helper()
	f, group, targets := newMigratedMappingReadFixture(t)
	_, descriptionErr := f.db.Exec(`UPDATE devices SET description='' WHERE id='device-A'`)
	require.NoError(t, descriptionErr)
	repo := sourcerule.NewSQLRepository(f.db)
	_, err := f.db.Exec(`UPDATE tags SET labels='{"source":"source-rule","source_rule_id":"rule-native","source_rule_address":"40001"}' WHERE id='tag-A'`)
	require.NoError(t, err)
	require.NoError(t, repo.Create(t.Context(), &schema.SourceRule{ID: "rule-native", DeviceID: "device-A", RevisionID: "R-native", StartAddress: "40001", Count: 1, DataType: schema.DataTypeFloat32, Enabled: true, Origin: "manual", CreatedAt: time.Now(), UpdatedAt: time.Now()}))
	tagID, mappingID := "tag-A", "mapping-A"
	require.NoError(t, repo.CreateLinks(t.Context(), []*schema.SourceRuleLink{{ID: "link-native", RuleID: "rule-native", Address: "40001", PointID: "point-A", TagID: &tagID, MappingID: &mappingID, CreatedAt: time.Now(), UpdatedAt: time.Now()}}))
	devices := device.NewService(device.NewSQLRepository(f.db), nil)
	points := point.NewService(point.NewSQLRepository(f.db), nil)
	groups := workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db))
	services := &DatalinkServices{Workspace: f.workspace, WriteGroups: groups, Device: devices, Point: points, Tag: tag.NewService(tag.NewSQLRepository(f.db)), Mapping: mapping.NewService(mapping.NewSQLRepository(f.db)), SourceRule: sourcerule.NewService(repo, devices, points, nil), DBTarget: dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db)), DBMapping: targets}
	_, err = services.Device.GetByID(t.Context(), "device-A")
	require.NoError(t, err)
	_, err = services.Point.GetByID(t.Context(), "point-A")
	require.NoError(t, err)
	_, err = services.Tag.GetByID(t.Context(), "tag-A")
	require.NoError(t, err)
	_, err = services.Mapping.GetByID(t.Context(), "mapping-A")
	require.NoError(t, err)
	_, err = services.SourceRule.GetByID(t.Context(), "rule-native")
	require.NoError(t, err)
	_, err = services.SourceRule.ListLinks(t.Context(), "rule-native")
	require.NoError(t, err)
	return &runtimeSourceRig{f: managedGroupFixture{writeGroupRouterFixture: f, groups: groups, group: map[string]any{"id": group.ID}}, services: services, router: NewRouter(services)}
}
