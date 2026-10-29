package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type runtimeSourcePageRig struct {
	f        writeGroupRouterFixture
	repo     *sourcerule.SQLRepository
	services *DatalinkServices
	router   *gin.Engine
	points   int
}

// Only the existing local SQLite configuration fixture is used. No recording
// provider is installed: source listing must not fan out into group details.
func newRuntimeSourcePageRig(t *testing.T) *runtimeSourcePageRig {
	t.Helper()
	f := newWriteGroupRouterFixture(t)
	_, err := f.db.Exec(`UPDATE devices SET description='' WHERE id='device-A'`)
	require.NoError(t, err)
	_, err = f.db.Exec(`INSERT INTO devices (id,name,description,protocol,status,connection_config,readiness_status) VALUES ('device-B','Device B','','modbus_tcp','draft','{}','{}')`)
	require.NoError(t, err)
	_, err = f.workspace.AttachDevice(t.Context(), "device-B")
	require.NoError(t, err)
	_, err = f.workspace.BindDatabaseConnector(t.Context(), "connector-A")
	require.NoError(t, err)
	devices := device.NewService(device.NewSQLRepository(f.db), nil)
	points := point.NewService(point.NewSQLRepository(f.db), pollinggroup.NewSQLRepository(f.db))
	repo := sourcerule.NewSQLRepository(f.db)
	services := &DatalinkServices{Workspace: f.workspace, Device: devices, Point: points, SourceRule: sourcerule.NewService(repo, devices, points, nil), Tag: tag.NewService(tag.NewSQLRepository(f.db)), Mapping: mapping.NewService(mapping.NewSQLRepository(f.db))}
	return &runtimeSourcePageRig{f: f, repo: repo, services: services, router: NewRouter(services)}
}

func (r *runtimeSourcePageRig) seed(t *testing.T, deviceID, ruleID string, mappings int) {
	t.Helper()
	now := time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)
	require.NoError(t, r.repo.Create(t.Context(), &schema.SourceRule{ID: ruleID, DeviceID: deviceID, StartAddress: fmt.Sprint(41000 + r.points), Count: max(1, mappings), DataType: schema.DataTypeFloat32, NamingPrefix: "Temperature " + ruleID, Enabled: true, Origin: "manual", RevisionID: "rev-" + ruleID, CreatedAt: now, UpdatedAt: now}))
	links := make([]*schema.SourceRuleLink, 0, mappings)
	for i := 0; i < mappings; i++ {
		address := fmt.Sprint(41000 + r.points)
		r.points++
		suffix := fmt.Sprintf("%s-%03d", ruleID, i)
		pointID, tagID, mappingID := "point-"+suffix, "tag-"+suffix, "mapping-"+suffix
		_, err := r.f.db.Exec(`INSERT INTO points(id,device_id,name,address,function,data_type,mode,enabled) VALUES(?,?,?,?, 'FC03','float32','read',1)`, pointID, deviceID, suffix, address)
		require.NoError(t, err)
		labels, err := json.Marshal(map[string]string{"source": "source-rule", "source_rule_id": ruleID, "source_rule_address": address})
		require.NoError(t, err)
		_, err = r.f.db.Exec(`INSERT INTO tags(id,key,key_lower,display_name,data_type,status,labels) VALUES(?,?,?,?,'float32','active',?)`, tagID, "temperature_"+suffix, "temperature_"+suffix, "Temperature "+suffix, string(labels))
		require.NoError(t, err)
		_, err = r.f.db.Exec(`INSERT INTO mappings(id,point_id,tag_id,transform_pipeline,status,enabled) VALUES(?,?,?,'[]','active',1)`, mappingID, pointID, tagID)
		require.NoError(t, err)
		links = append(links, &schema.SourceRuleLink{ID: "link-" + suffix, RuleID: ruleID, Address: address, PointID: pointID, TagID: &tagID, MappingID: &mappingID, CreatedAt: now, UpdatedAt: now})
	}
	require.NoError(t, r.repo.CreateLinks(t.Context(), links))
}
func (r *runtimeSourcePageRig) query(t *testing.T, kind string) url.Values {
	t.Helper()
	record, err := r.f.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	return url.Values{"workspace_id": {record.ID}, "workspace_revision": {record.DatabaseSetupRevision}, "device_id": {"device-A"}, "source_list": {kind}}
}
func (r *runtimeSourcePageRig) get(t *testing.T, q url.Values) *httptest.ResponseRecorder {
	t.Helper()
	return performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+q.Encode(), nil)
}

type runtimeSourcePageResponse struct {
	ContractVersion   int              `json:"contract_version"`
	WorkspaceID       string           `json:"workspace_id"`
	WorkspaceRevision string           `json:"workspace_revision"`
	DeviceID          string           `json:"device_id"`
	Kind              string           `json:"kind"`
	Query             string           `json:"query"`
	Limit             int              `json:"limit"`
	ListRevision      string           `json:"list_revision"`
	Rules             []map[string]any `json:"rules"`
	Mappings          []map[string]any `json:"mappings"`
	NextCursor        *string          `json:"next_cursor"`
}

func runtimeSourcePageBody(t *testing.T, response *httptest.ResponseRecorder) *runtimeSourcePageResponse {
	t.Helper()
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data struct {
			Page *runtimeSourcePageResponse `json:"source_page"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.NotNil(t, body.Data.Page, response.Body.String())
	require.Equal(t, 1, body.Data.Page.ContractVersion)
	require.NotEmpty(t, body.Data.Page.ListRevision)
	return body.Data.Page
}

// Inject only timing/failure, retaining actual SQL source reads and services.
type runtimeSourcePageFaultRepo struct {
	sourcerule.Repository
	calls    int
	failAt   int
	changeAt int
	change   func()
}

func (r *runtimeSourcePageFaultRepo) List(ctx context.Context, filter sourcerule.ListFilter) ([]*schema.SourceRule, error) {
	r.calls++
	if r.calls == r.failAt {
		return nil, errors.New("password=fixture-secret /private/list.sqlite")
	}
	result, err := r.Repository.List(ctx, filter)
	if err == nil && r.calls == r.changeAt && r.change != nil {
		r.change()
	}
	return result, err
}
func (r *runtimeSourcePageRig) fault(repo sourcerule.Repository) {
	r.services.SourceRule = sourcerule.NewService(repo, r.services.Device, r.services.Point, nil)
	r.router = NewRouter(r.services)
}
