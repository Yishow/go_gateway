package api

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const runtimeRuleContextPath = "/api/v1/datalink/studio-v2/workspace/runtime-context"

type runtimeRuleEvidenceRig struct {
	f      managedGroupFixture
	pipe   *grouppipeline.Pipeline
	router *gin.Engine
	clock  time.Time
}

// Both the configuration and destination are disposable local SQLite files.
// Apply, immutable versions and frozen layout proof are production services.
func newRuntimeRuleEvidenceRig(t *testing.T) *runtimeRuleEvidenceRig {
	t.Helper()
	f := newManagedDestinationFixture(t, "sqlite")
	create := f.createRequest()
	destination := create["group"].(map[string]any)["destination"].(map[string]any)
	destination["storage_strategy"], destination["table_schema"], destination["table_name"] = "managed", "main", "meter_old"
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, create)
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	f.data, f.group = groupSaveData(t, decodeJSONBody(t, created))
	// The older write-group fixture omits nullable description because it does
	// not read Device SQLRepository; runtime context does read that real service.
	_, err := f.db.Exec(`UPDATE devices SET description='' WHERE id='device-A'`)
	require.NoError(t, err)
	devices := device.NewService(device.NewSQLRepository(f.db), nil)
	_, err = devices.GetByID(t.Context(), "device-A")
	require.NoError(t, err)
	r := &runtimeRuleEvidenceRig{f: f, clock: time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC)}
	f.groups.WithClock(func() time.Time { return r.clock })
	token := f.preview(t)
	request := f.schemaRequest()
	request["token"], request["operation_id"] = token.Token, token.OperationID
	response := performJSONRequest(t, f.router, http.MethodPost, f.confirmPath(), request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	id := f.group["id"].(string)
	response = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	r.f.data = getGroupData(t, f, id)
	r.f.group = r.f.data["group"].(map[string]any)
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{Groups: f.groups, Store: groupdelivery.NewStore(f.db)}, grouppipeline.Config{})
	r.router = NewRouter(&DatalinkServices{
		Workspace: f.workspace, WriteGroups: f.groups, WriteGroupDelivery: r.pipe,
		Device: devices,
	})
	return r
}

func (r *runtimeRuleEvidenceRig) deterministicRouter(groups handlers.RuntimeRecordingGroupReader) *gin.Engine {
	h := handlers.NewRuntimeHandler(device.NewService(device.NewSQLRepository(r.f.db), nil), nil, nil, nil, nil, r.f.workspace).
		WithRuleEvidenceServices(groups, r.pipe).WithRuleEvidenceClock(func() time.Time { return r.clock })
	router := gin.New()
	router.GET(runtimeRuleContextPath, h.WorkspaceContext)
	return router
}

func (r *runtimeRuleEvidenceRig) update(t *testing.T, table string, interval int) *workspace.WriteGroupSaveResult {
	t.Helper()
	current, err := r.f.groups.Get(t.Context(), r.f.group["id"].(string))
	require.NoError(t, err)
	current.Group.Destination.TableName = table
	current.Group.RowPolicy.IntervalSeconds = interval
	updated, err := r.f.groups.Update(t.Context(), current.Group.ID, workspace.WriteGroupMutation{
		WorkspaceID: current.Group.WorkspaceID, ExpectedWorkspaceRevision: current.WorkspaceRevision,
		ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: current.Group.Destination.ConnectorRevision, Group: current.Group,
	})
	require.NoError(t, err)
	r.f.data = getGroupData(t, r.f, current.Group.ID)
	r.f.group = r.f.data["group"].(map[string]any)
	return updated
}

func (r *runtimeRuleEvidenceRig) apply(t *testing.T) *workspace.WriteGroupSaveResult {
	t.Helper()
	token := r.f.preview(t)
	request := r.f.schemaRequest()
	request["token"], request["operation_id"] = token.Token, token.OperationID
	response := performJSONRequest(t, r.f.router, http.MethodPost, r.f.confirmPath(), request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	current, err := r.f.groups.Get(t.Context(), r.f.group["id"].(string))
	require.NoError(t, err)
	applied, err := r.f.groups.Apply(t.Context(), current.Group.ID, workspace.WriteGroupMutation{
		WorkspaceID: current.Group.WorkspaceID, ExpectedWorkspaceRevision: current.WorkspaceRevision,
		ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: current.Group.Destination.ConnectorRevision,
	})
	require.NoError(t, err)
	r.f.data = getGroupData(t, r.f, current.Group.ID)
	r.f.group = r.f.data["group"].(map[string]any)
	return applied
}

// Only inject the timing of a real save; all successful reads still come from
// the actual workspace transaction and production immutable provider.
type runtimeRuleChangingGroups struct {
	groups     *workspace.WriteGroupService
	afterFirst func()
	calls      int
}

type runtimeRuleDeletingWorkspaceGroups struct {
	groups *workspace.WriteGroupService
	delete func()
	calls  int
}

func (g *runtimeRuleDeletingWorkspaceGroups) List(ctx context.Context) (*workspace.WriteGroupListResult, error) {
	result, err := g.groups.List(ctx)
	g.calls++
	if err == nil && g.calls == 2 {
		g.delete()
	}
	return result, err
}

func (g *runtimeRuleChangingGroups) List(ctx context.Context) (*workspace.WriteGroupListResult, error) {
	result, err := g.groups.List(ctx)
	g.calls++
	if err == nil && g.calls == 1 {
		g.afterFirst()
	}
	return result, err
}

func (r *runtimeRuleEvidenceRig) query(t *testing.T) url.Values {
	t.Helper()
	current, err := r.f.groups.Get(t.Context(), r.f.group["id"].(string))
	require.NoError(t, err)
	return url.Values{"workspace_id": {r.f.record.ID}, "workspace_revision": {current.WorkspaceRevision}, "device_id": {"device-A"}}
}
