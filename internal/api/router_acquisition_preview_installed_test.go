package api

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type changingAcquisitionReader struct {
	service *datalinkruntime.Service
	after   func()
}

func (r changingAcquisitionReader) ReadAcquisitionPreview(ctx context.Context, workspaceID, id, head string) (collector.AcquisitionPreview, func() bool, error) {
	view, verify, err := r.service.ReadAcquisitionPreview(ctx, workspaceID, id, head)
	if r.after != nil {
		r.after()
	}
	return view, verify, err
}
func TestAcquisitionPreviewInstalledHTTPHeadParityAndCaptureRace(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "preview.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	for _, statement := range []string{
		`INSERT INTO devices(id,name,description,protocol,status,connection_config) VALUES('d','PLC','','mc_3e','active','{"host":"127.0.0.1","port":5007}')`,
		`INSERT INTO polling_groups(id,name,interval_ms,enabled) VALUES('g','g',3600000,1)`,
		`INSERT INTO points(id,device_id,name,address,data_type,revision_id,polling_group_id) VALUES('p','d','p','D0','uint16','p1','g')`,
		`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id) VALUES('r','d','D0',1,'uint16','p','r1')`,
		`INSERT INTO source_rule_links(id,rule_id,address,point_id) VALUES('l','r','D0','p')`,
		`INSERT INTO system_settings(key,value) VALUES('studio_v2_workspace','{"id":"w","ordered_device_ids":["d"],"database_setup_revision":""}')`,
	} {
		_, err = db.ExecContext(t.Context(), statement)
		require.NoError(t, err)
	}
	devices := device.NewService(device.NewSQLRepository(db), nil)
	groups := pollinggroup.NewService(pollinggroup.NewSQLRepository(db))
	points := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	mappings := mapping.NewService(mapping.NewSQLRepository(db))
	tags := tag.NewService(tag.NewSQLRepository(db))
	rules := sourcerule.NewService(sourcerule.NewSQLRepository(db), devices, points, nil)
	ws := workspace.NewService(workspace.NewSQLRepository(db)).WithRuntimeProjectionServices(devices, rules, points, mappings, tags, nil, nil, groups)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), manager)
	runtime, err := datalinkruntime.NewService(datalinkruntime.DefaultConfig(), datalinkruntime.Dependencies{Scheduler: scheduler, Writer: storage.NewSQLiteWriter(db), WorkspaceProjection: ws, DeviceService: devices, PointService: points, MappingService: mappings, TagService: tags, PollingGroupService: groups})
	require.NoError(t, err)
	projection, err := ws.ReadRuntimeProjection(t.Context())
	require.NoError(t, err)
	graph, ok := projection.AcquisitionSourceGraph("d")
	require.True(t, ok, graph.Reason())
	require.NoError(t, runtime.ApplyWorkspaceProjection(t.Context(), projection))
	require.NoError(t, scheduler.Start(projection.PollingGroups))
	t.Cleanup(func() { require.NoError(t, scheduler.Stop()) })
	route := studioLocatorPath + "/acquisition-plan-preview"
	router := gin.New()
	router.GET(route, handlers.NewStudioV2AcquisitionPreviewHandler(ws, runtime).Get)
	query := url.Values{"workspace_id": {"w"}, "device_id": {"d"}, "expected_workspace_revision": {""}}
	before := manager.PeekAcquisitionConnection(graph.Device())
	require.Equal(t, "absent", before.State)
	response := performJSONRequest(t, router, http.MethodGet, route+"?"+query.Encode(), nil)
	data := locatorData(t, response)
	require.Equal(t, "available", data["status"], response.Body.String())
	head := data["preview_revision"].(string)
	require.NotEmpty(t, head)
	require.Equal(t, before, manager.PeekAcquisitionConnection(graph.Device()), "GET must not connect")
	query.Set("expected_preview_revision", head)
	require.Equal(t, http.StatusOK, performJSONRequest(t, router, http.MethodGet, route+"?"+query.Encode(), nil).Code)
	_, err = db.ExecContext(t.Context(), `UPDATE points SET revision_id='p2' WHERE id='p'`)
	require.NoError(t, err)
	response = performJSONRequest(t, router, http.MethodGet, route+"?"+query.Encode(), nil)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	_, err = db.ExecContext(t.Context(), `UPDATE points SET revision_id='p1' WHERE id='p'`)
	require.NoError(t, err)
	router = gin.New()
	router.GET(route, handlers.NewStudioV2AcquisitionPreviewHandler(ws, changingAcquisitionReader{runtime, func() { scheduler.RemovePoint("p") }}).Get)
	response = performJSONRequest(t, router, http.MethodGet, route+"?"+query.Encode(), nil)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
}
