package api

import (
	"net/http"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func newMigratedMappingReadFixture(t *testing.T) (writeGroupRouterFixture, *workspace.WriteGroup, *dbtarget.MappingService) {
	t.Helper()
	f, _ := newSingleMigrationRouterFixture(t)
	groups := workspace.NewSQLWriteGroupRepository(f.db)
	group := &workspace.WriteGroup{
		WorkspaceID: f.record.ID, Name: "Reviewed Temperature",
		Members:     []workspace.WriteGroupMember{{DeviceID: "device-A", PointID: "point-A", TagID: "tag-A", TargetColumn: "reviewed_temperature", Required: true}},
		Destination: workspace.WriteGroupDestination{ConnectorID: "connector-A", ConnectorRevision: "connector-1", TableSchema: "main", TableName: "reviewed_values", StorageStrategy: workspace.WriteGroupStorageStrategyCustom},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 30, IncompletePolicy: "skip_row"},
		WritePolicy: workspace.WriteGroupWritePolicy{Mode: "append"},
		Migration:   workspace.WriteGroupMigration{SourceKind: "legacy-single-mapping", SourceIDs: []string{"legacy-A"}, SourceRevision: "reviewed-source", AdapterVersion: "single-mapping-v1", ReviewResult: "confirmed_snapshot_conversion"},
	}
	require.NoError(t, groups.Create(t.Context(), group))
	connectors := dbtarget.NewSQLConnectorRepository(f.db)
	mappings := dbtarget.NewSQLTargetMappingRepository(f.db)
	mappingSvc := dbtarget.NewMappingService(mappings, connectors, nil)
	f.router = NewRouter(&DatalinkServices{
		Workspace: f.workspace, WriteGroups: workspace.NewWriteGroupService(f.workspace, groups),
		DBTarget: dbtarget.NewConnectorService(connectors, mappings), DBMapping: mappingSvc,
		SourceRule: sourcerule.NewService(sourcerule.NewSQLRepository(f.db), nil, nil, nil),
	})
	return f, group, mappingSvc
}

func TestNewRouter_LegacySingleMappingWorkspaceReadProjectsCanonicalDraft(t *testing.T) {
	f, group, _ := newMigratedMappingReadFixture(t)
	read := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-targets", nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	data := decodeJSONBody(t, read)["data"].([]any)
	require.Len(t, data, 1)
	row := data[0].(map[string]any)
	require.Equal(t, "legacy-A", row["id"])
	require.Equal(t, "point-A", row["point_id"])
	require.Equal(t, group.ID, row["row_group_id"])
	require.Equal(t, "reviewed_temperature", row["column_name"])
	require.Equal(t, group.ID, row["canonical_group"].(map[string]any)["id"])
}

func TestNewRouter_LegacySingleMappingReadProjectsCanonicalDraft(t *testing.T) {
	f, group, mappingSvc := newMigratedMappingReadFixture(t)
	read := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/db-targets/mappings/legacy-A", nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	data := decodeJSONBody(t, read)["data"].(map[string]any)
	require.Equal(t, "legacy-A", data["id"])
	require.Equal(t, "reviewed_temperature", data["column_name"])
	require.Equal(t, "reviewed_values", data["table_name"])
	require.Equal(t, float64(30), data["write_interval_seconds"])
	canonical := data["canonical_group"].(map[string]any)
	require.Equal(t, group.ID, canonical["id"])
	require.Equal(t, "", canonical["applied_revision"])
	intent := data["legacy_intent"].(map[string]any)
	require.Equal(t, "temperature", intent["column_name"])
	require.Equal(t, "observed_at", intent["timestamp_column"])
	legacy, err := mappingSvc.GetByID(t.Context(), "legacy-A")
	require.NoError(t, err)
	require.Equal(t, "temperature", legacy.ColumnName)
	require.Equal(t, "raw_values", legacy.TableName)
	require.True(t, legacy.Enabled)
}

func TestNewRouter_LegacySingleMappingReadFiltersCanonicalDisabledState(t *testing.T) {
	f, group, mappingSvc := newMigratedMappingReadFixture(t)
	disabled := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+group.ID+"/disable", map[string]any{
		"workspace_id": f.record.ID, "expected_workspace_revision": f.record.DatabaseSetupRevision,
		"expected_group_revision": group.Revision, "expected_connector_revision": "connector-1",
	})
	require.Equal(t, http.StatusOK, disabled.Code, disabled.Body.String())
	for _, enabled := range []string{"false", "true"} {
		read := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/db-targets/mappings?enabled="+enabled, nil)
		require.Equal(t, http.StatusOK, read.Code, read.Body.String())
		rows := decodeJSONBody(t, read)["data"].([]any)
		if enabled == "true" {
			require.Empty(t, rows)
		} else {
			require.Len(t, rows, 1)
			require.Equal(t, false, rows[0].(map[string]any)["enabled"])
		}
	}
	read := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-targets", nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	rows := decodeJSONBody(t, read)["data"].([]any)
	require.Len(t, rows, 1)
	require.Equal(t, false, rows[0].(map[string]any)["enabled"])
	legacy, err := mappingSvc.GetByID(t.Context(), "legacy-A")
	require.NoError(t, err)
	require.True(t, legacy.Enabled)
}

func TestNewRouter_LegacySingleMappingReadCannotFlattenEditedMembers(t *testing.T) {
	f, group, _ := newMigratedMappingReadFixture(t)
	_, err := f.db.ExecContext(t.Context(), `INSERT INTO points (id,device_id,name,address,function,data_type,mode,enabled)
		VALUES ('point-B','device-A','Pressure','40002','FC03','float32','read',1)`)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels)
		VALUES ('tag-B','pressure','pressure','Pressure','float32','active','{}')`)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled)
		VALUES ('mapping-B','point-B','tag-B','[]','active',1)`)
	require.NoError(t, err)
	group.Members = append(group.Members, workspace.WriteGroupMember{DeviceID: "device-A", PointID: "point-B", TagID: "tag-B", TargetColumn: "pressure", Required: true})
	edit := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+group.ID, map[string]any{
		"workspace_id": f.record.ID, "expected_workspace_revision": f.record.DatabaseSetupRevision,
		"expected_group_revision": group.Revision, "expected_connector_revision": "connector-1", "group": group,
	})
	require.Equal(t, http.StatusOK, edit.Code, edit.Body.String())
	for _, path := range []string{"/api/v1/datalink/db-targets/mappings/legacy-A", "/api/v1/datalink/db-targets/mappings", "/api/v1/datalink/studio-v2/workspace/database-targets"} {
		read := performJSONRequest(t, f.router, http.MethodGet, path, nil)
		require.Equal(t, http.StatusConflict, read.Code, read.Body.String())
		body := decodeJSONBody(t, read)
		assertWriteGroupError(t, body, "WRITE_GROUP_LEGACY_READ_CONFLICT")
		require.Equal(t, "open_write_groups", body["error"].(map[string]any)["action"])
	}
	read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+group.ID, nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	_, canonical := groupSaveData(t, decodeJSONBody(t, read))
	require.Len(t, canonical["members"], 2)
	read = performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/db-targets/mappings?connector_id=unrelated", nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	require.Empty(t, decodeJSONBody(t, read)["data"])
}

func TestNewRouter_LegacySingleMappingReadCannotTranslateConnectorChange(t *testing.T) {
	f, group, _ := newMigratedMappingReadFixture(t)
	_, err := f.db.ExecContext(t.Context(), `INSERT INTO database_connectors
		(id,name,kind,connection_config,identity_revision,status,enabled)
		VALUES ('connector-B','Other target','sqlite','{"database":"other-target.db"}','connector-2','ready',1)`)
	require.NoError(t, err)
	group.Destination.ConnectorID = "connector-B"
	group.Destination.ConnectorRevision = "connector-2"
	group.Destination.Database = ""
	edit := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+group.ID, map[string]any{
		"workspace_id": f.record.ID, "expected_workspace_revision": f.record.DatabaseSetupRevision,
		"expected_group_revision": group.Revision, "expected_connector_revision": "connector-2", "group": group,
	})
	require.Equal(t, http.StatusOK, edit.Code, edit.Body.String())
	for _, path := range []string{
		"/api/v1/datalink/db-targets/mappings/legacy-A",
		"/api/v1/datalink/db-targets/mappings?connector_id=connector-A",
		"/api/v1/datalink/db-targets/mappings?connector_id=connector-B",
		"/api/v1/datalink/studio-v2/workspace/database-targets",
	} {
		read := performJSONRequest(t, f.router, http.MethodGet, path, nil)
		require.Equal(t, http.StatusConflict, read.Code, read.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, read), "WRITE_GROUP_LEGACY_READ_CONFLICT")
	}
	read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+group.ID, nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, read))
	require.Equal(t, "connector-B", saved["destination"].(map[string]any)["connector_id"])
}
