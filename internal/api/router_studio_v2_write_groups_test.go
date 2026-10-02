package api

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"sync"
	"testing"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

const writeGroupsPath = "/api/v1/datalink/studio-v2/workspace/write-groups"

type writeGroupRouterFixture struct {
	db        *sql.DB
	workspace *workspace.Service
	router    *gin.Engine
	record    *workspace.Record
}

func newWriteGroupRouterFixture(t *testing.T) writeGroupRouterFixture {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "configuration.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	ctx := t.Context()
	ws := workspace.NewService(workspace.NewSQLRepository(db))
	_, err = ws.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO devices (id,name,protocol,status,connection_config,readiness_status)
		VALUES ('device-A','PLC','modbus_tcp','draft','{}','{}')`)
	require.NoError(t, err)
	record, err := ws.AttachDevice(ctx, "device-A")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO points (id,device_id,name,address,function,data_type,mode,enabled)
		VALUES ('point-A','device-A','Temperature','40001','FC03','float32','read',1)`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels)
		VALUES ('tag-A','temperature','temperature','Temperature','float32','active','{}')`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled)
		VALUES ('mapping-A','point-A','tag-A','[]','active',1)`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO database_connectors (id,name,kind,connection_config,identity_revision,status,enabled)
		VALUES ('connector-A','Target','sqlite','{"database":"target.db"}','connector-1','ready',1)`)
	require.NoError(t, err)
	services := &DatalinkServices{Workspace: ws, WriteGroups: workspace.NewWriteGroupService(ws, workspace.NewSQLWriteGroupRepository(db))}
	return writeGroupRouterFixture{db: db, workspace: ws, router: NewRouter(services), record: record}
}

func (f writeGroupRouterFixture) createRequest() map[string]any {
	return map[string]any{
		"workspace_id": f.record.ID, "expected_workspace_revision": f.record.DatabaseSetupRevision,
		"expected_connector_revision": "connector-1",
		"group": map[string]any{
			"workspace_id": f.record.ID, "name": "Raw Tags",
			"id": "client-group", "revision": "client-revision", "applied_revision": "client-applied", "status": "running",
			"members":     []map[string]any{{"device_id": "device-A", "point_id": "point-A", "tag_id": "tag-A", "target_column": "temperature", "required": true}},
			"destination": map[string]any{"connector_id": "connector-A", "connector_revision": "connector-1", "table_name": "raw_values", "storage_strategy": "custom"},
			"row_policy":  map[string]any{"interval_seconds": 15},
		},
	}
}

func groupSaveData(t *testing.T, response map[string]any) (data, group map[string]any) {
	t.Helper()
	require.Equal(t, true, response["success"])
	data, ok := response["data"].(map[string]any)
	require.True(t, ok)
	group, ok = data["group"].(map[string]any)
	require.True(t, ok)
	return data, group
}

func groupEditRequest(f writeGroupRouterFixture, data, group map[string]any) map[string]any {
	return map[string]any{
		"workspace_id": f.record.ID, "expected_workspace_revision": data["workspace_revision"],
		"expected_group_revision": group["revision"], "expected_connector_revision": "connector-1", "group": group,
	}
}

func TestNewRouter_AtomicGroupSaveAndCAS(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, created))
	require.NotEqual(t, "client-group", group["id"])
	require.NotEqual(t, "client-revision", group["revision"])
	require.Equal(t, "", group["applied_revision"])
	require.Equal(t, "draft", group["status"])
	id := group["id"].(string)
	before, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Len(t, before.DatabaseRowGroups, 1)
	require.Equal(t, id, before.DatabaseRowGroups[0].ID)
	require.Equal(t, []string{"point-A"}, before.DatabaseRowGroups[0].MemberPointIDs)
	require.Equal(t, data["workspace_revision"], before.DatabaseSetupRevision)

	_, err = f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_group_projection BEFORE UPDATE ON system_settings
		WHEN NEW.key = 'studio_v2_workspace' BEGIN SELECT RAISE(ABORT, 'private-storage-diagnostic'); END`)
	require.NoError(t, err)
	group["name"] = "Must Roll Back"
	group["members"].([]any)[0].(map[string]any)["target_column"] = "changed_column"
	failed := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, groupEditRequest(f, data, group))
	require.Equal(t, http.StatusInternalServerError, failed.Code, failed.Body.String())
	require.NotContains(t, failed.Body.String(), "private-storage-diagnostic")
	assertWriteGroupError(t, decodeJSONBody(t, failed), "WRITE_GROUP_UNAVAILABLE")
	after, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, before, after)
	reloaded := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	require.Equal(t, http.StatusOK, reloaded.Code, reloaded.Body.String())
	oldData, oldGroup := groupSaveData(t, decodeJSONBody(t, reloaded))
	expectedData, expectedGroup := groupSaveData(t, decodeJSONBody(t, created))
	require.Equal(t, expectedData, oldData)
	require.Equal(t, expectedGroup, oldGroup)
	require.Equal(t, "Raw Tags", oldGroup["name"])
	require.Equal(t, data["workspace_revision"], oldData["workspace_revision"])
	require.Equal(t, group["revision"], oldGroup["revision"])
	_, err = f.db.ExecContext(t.Context(), `DROP TRIGGER fail_group_projection`)
	require.NoError(t, err)

	var wg sync.WaitGroup
	results := make(chan int, 2)
	for _, name := range []string{"Client A", "Client B"} {
		candidate := make(map[string]any, len(oldGroup))
		for key, value := range oldGroup {
			candidate[key] = value
		}
		candidate["name"] = name
		request := groupEditRequest(f, oldData, candidate)
		wg.Go(func() {
			results <- performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, request).Code
		})
	}
	wg.Wait()
	close(results)
	var codes []int
	for code := range results {
		codes = append(codes, code)
	}
	require.ElementsMatch(t, []int{http.StatusOK, http.StatusConflict}, codes)
	finalRead := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	finalData, finalGroup := groupSaveData(t, decodeJSONBody(t, finalRead))
	require.NotEqual(t, oldData["workspace_revision"], finalData["workspace_revision"])
	require.NotEqual(t, oldGroup["revision"], finalGroup["revision"])
	require.Contains(t, []string{"Client A", "Client B"}, finalGroup["name"])
}

func assertWriteGroupError(t *testing.T, body map[string]any, code string) {
	t.Helper()
	require.Equal(t, false, body["success"])
	envelope, ok := body["error"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, code, envelope["code"])
	require.NotEmpty(t, envelope["request_id"])
	require.NotEmpty(t, envelope["action"])
}
