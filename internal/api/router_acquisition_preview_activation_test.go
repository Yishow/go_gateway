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
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Only readiness for the already-applied recording group is supplied here.
// Activation, SQL source projection, runtime installation and preview GET are real.
type acquisitionActivationGroup struct{}

func (acquisitionActivationGroup) List(context.Context) (*workspace.WriteGroupListResult, error) {
	return &workspace.WriteGroupListResult{WorkspaceID: "w", Groups: []*workspace.WriteGroup{{ID: "record", WorkspaceID: "w", Revision: "r", AppliedRevision: "r", Members: []workspace.WriteGroupMember{{DeviceID: "d", PointID: "p"}}}}}, nil
}
func (acquisitionActivationGroup) Readiness(context.Context, string) (*workspace.WriteGroupReadiness, error) {
	return &workspace.WriteGroupReadiness{WorkspaceID: "w", GroupID: "record", GroupRevision: "r", AppliedRevision: "r", Ready: true}, nil
}
func TestAcquisitionPreviewActivationScopePublishesInstalledSet(t *testing.T) {
	for _, preserve := range []bool{false, true} {
		t.Run(map[bool]string{false: "fresh-subset", true: "preserve-other-installed"}[preserve], func(t *testing.T) {
			db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "preview.db"))
			require.NoError(t, err)
			db.SetMaxOpenConns(1)
			t.Cleanup(func() { require.NoError(t, db.Close()) })
			require.NoError(t, datalink.NewMigrator().Migrate(db))
			for _, statement := range []string{
				`INSERT INTO devices(id,name,description,protocol,status,connection_config) VALUES('d','PLC','','mc_3e','active','{"host":"127.0.0.1","port":5007}')`,
				`INSERT INTO polling_groups(id,name,interval_ms,enabled) VALUES('g','g',3600000,1)`,
				`INSERT INTO points(id,device_id,name,address,data_type,revision_id,polling_group_id) VALUES('p','d','p','D0','uint16','p1','g'),('q','d','q','D2','uint16','q1','g')`,
				`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id) VALUES('r','d','D0',1,'uint16','p','r1'),('rq','d','D2',1,'uint16','q','rq1')`,
				`INSERT INTO source_rule_links(id,rule_id,address,point_id) VALUES('l','r','D0','p'),('lq','rq','D2','q')`,
				`UPDATE devices SET last_test_success=1`,
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

			require.NoError(t, scheduler.Start(projection.PollingGroups))
			t.Cleanup(func() { require.NoError(t, scheduler.Stop()) })
			var oldLease uint64
			if preserve {
				require.NoError(t, runtime.ApplyWorkspaceProjectionForDevicesAndPoints(t.Context(), projection, []string{"d"}, []string{"q"}))
				oldLease = scheduler.CapturePointLease("q", "q1")
				require.NotZero(t, oldLease)
			}
			ws.WithReadinessServices(devices, nil, nil, nil).WithWriteGroupReadiness(acquisitionActivationGroup{})
			activated, err := workspace.NewActivationService(ws, devices, runtime).ActivateScope(t.Context(), []string{"d"}, []string{"record"})
			require.NoError(t, err)
			require.Len(t, activated.Results, 1)
			require.Equal(t, workspace.ActivationResultStatusSuccess, activated.Results[0].Status, activated.Results[0].Message)
			require.NotZero(t, scheduler.CapturePointLease("p", "p1"))
			require.Equal(t, oldLease, scheduler.CapturePointLease("q", "q1"), "record subset must neither install absent q nor restart installed q")
			router := gin.New()
			route := studioLocatorPath + "/acquisition-plan-preview"
			router.GET(route, handlers.NewStudioV2AcquisitionPreviewHandler(ws, runtime).Get)
			query := url.Values{"workspace_id": {"w"}, "device_id": {"d"}, "expected_workspace_revision": {""}}
			before := manager.PeekAcquisitionConnection(graph.Device())
			response := performJSONRequest(t, router, http.MethodGet, route+"?"+query.Encode(), nil)
			data := locatorData(t, response)
			require.Equal(t, "available", data["status"], response.Body.String())
			groupsData, ok := data["groups"].([]any)
			require.True(t, ok)
			require.Len(t, groupsData, 1)
			want := float64(1)
			if preserve {
				want = 2
			}
			groupData, ok := groupsData[0].(map[string]any)
			require.True(t, ok)
			require.Equal(t, want, groupData["acquisition_point_count"])
			require.Equal(t, before, manager.PeekAcquisitionConnection(graph.Device()), "GET never connects")
			require.Equal(t, schema.DeviceStatusActive, projection.Devices[0].Status)
		})
	}
}
