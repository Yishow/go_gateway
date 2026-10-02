package api

import (
	"database/sql"
	"net/http"
	"os"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const legacyDatabaseConfigPath = "/api/v1/datalink/studio-v2/workspace/database-config"

func TestNewRouter_LegacyWriteAdapterConflictRowReplacement(t *testing.T) {
	for _, change := range []string{"empty", "members", "table", "group-keys", "unique-keys", "connector"} {
		t.Run(change, func(t *testing.T) {
			f, targetPath := newLegacyWriteRouterFixture(t, true)
			groups := append([]workspace.DatabaseRowGroup(nil), f.record.DatabaseRowGroups...)
			request := legacyDatabaseConfigRequest(f, targetPath)
			switch change {
			case "empty":
				groups = []workspace.DatabaseRowGroup{}
			case "members":
				groups[0].MemberPointIDs = []string{"point-A"}
			case "table":
				groups[0].TableName = "other_table"
				request["table"] = "other_table"
			case "group-keys":
				groups[0].GroupKeyColumns = []string{"other_entity"}
			case "unique-keys":
				groups[0].UniqueKeyColumns = []string{"entity_id"}
			case "connector":
				groups[0].ConnectorID = "connector-B"
				request["connector_id"] = "connector-B"
				request["expected_connector_revision"] = "connector-2"
				_, err := f.db.ExecContext(t.Context(), `INSERT INTO database_connectors
					(id,name,kind,connection_config,identity_revision,status,enabled)
					VALUES ('connector-B','Other target','sqlite','{}','connector-2','ready',1)`)
				require.NoError(t, err)
			}
			request["row_groups"] = groups
			before := legacyWriteLocalSnapshot(t, f.db)
			response := performJSONRequest(t, f.router, http.MethodPut, legacyDatabaseConfigPath, request)
			assertLegacyWriteConflict(t, response.Code, decodeJSONBody(t, response))
			require.Equal(t, before, legacyWriteLocalSnapshot(t, f.db))
			require.NotContains(t, response.Body.String(), targetPath)
			_, err := os.Stat(targetPath)
			require.ErrorIs(t, err, os.ErrNotExist)
		})
	}
}

func createLegacyWriteTarget(t *testing.T, targetPath string) {
	t.Helper()
	db, err := sql.Open("sqlite", targetPath)
	require.NoError(t, err)
	defer func() { require.NoError(t, db.Close()) }()
	_, err = db.ExecContext(t.Context(), `CREATE TABLE raw_values (
		temperature REAL, other_temperature REAL, value REAL, entity_id TEXT, observed_at TEXT
	)`)
	require.NoError(t, err)
}

func TestNewRouter_LegacyWriteAdapterConflictRowRetention(t *testing.T) {
	for _, replacement := range []string{"omitted", "equal-current"} {
		t.Run(replacement, func(t *testing.T) {
			f, targetPath := newLegacyWriteRouterFixture(t, true)
			createLegacyWriteTarget(t, targetPath)
			before := legacyWriteLocalSnapshot(t, f.db)
			request := legacyDatabaseConfigRequest(f, targetPath)
			if replacement == "equal-current" {
				request["row_groups"] = f.record.DatabaseRowGroups
			}
			response := performJSONRequest(t, f.router, http.MethodPut, legacyDatabaseConfigPath, request)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			after, err := f.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)
			require.Equal(t, f.record.DatabaseRowGroups, after.DatabaseRowGroups)
			require.Equal(t, f.record.DatabaseTargetRefs, after.DatabaseTargetRefs)
			require.NotEqual(t, f.record.DatabaseSetupRevision, after.DatabaseSetupRevision)
			afterSnapshot := legacyWriteLocalSnapshot(t, f.db)
			for _, unchanged := range []string{"targets", "groups", "members", "versions", "migration_maps"} {
				require.Equal(t, before[unchanged], afterSnapshot[unchanged], unchanged)
			}
		})
	}
}

func TestNewRouter_LegacyWriteAdapterConflictRowInheritedScope(t *testing.T) {
	f, targetPath := newLegacyWriteRouterFixture(t, true)
	createLegacyWriteTarget(t, targetPath)
	db, err := sql.Open("sqlite", targetPath)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `CREATE TABLE other_table (
		temperature REAL, value REAL, entity_id TEXT, observed_at TEXT
	)`)
	require.NoError(t, err)
	require.NoError(t, db.Close())
	request := legacyDatabaseConfigRequest(f, targetPath)
	request["table"] = "other_table"
	groups := append([]workspace.DatabaseRowGroup(nil), f.record.DatabaseRowGroups...)
	groups[0].TableName = ""
	groups[0].TableSchema = ""
	groups[0].ConnectorID = ""
	request["row_groups"] = groups
	before := legacyWriteLocalSnapshot(t, f.db)
	response := performJSONRequest(t, f.router, http.MethodPut, legacyDatabaseConfigPath, request)
	assertLegacyWriteConflict(t, response.Code, decodeJSONBody(t, response))
	require.Equal(t, before, legacyWriteLocalSnapshot(t, f.db))
}

func TestNewRouter_LegacyWriteAdapterConflictRowRetentionAfterCanonicalEdit(t *testing.T) {
	f, targetPath := newLegacyWriteRouterFixture(t, true)
	createLegacyWriteTarget(t, targetPath)
	list := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
	require.Equal(t, http.StatusOK, list.Code, list.Body.String())
	group := decodeJSONBody(t, list)["data"].(map[string]any)["groups"].([]any)[0].(map[string]any)
	read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+group["id"].(string), nil)
	data, group := groupSaveData(t, decodeJSONBody(t, read))
	group["members"] = group["members"].([]any)[:1]
	updated := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+group["id"].(string), groupEditRequest(f, data, group))
	require.Equal(t, http.StatusOK, updated.Code, updated.Body.String())
	var err error
	f.record, err = f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, []string{"point-A"}, f.record.DatabaseRowGroups[0].MemberPointIDs)
	request := legacyDatabaseConfigRequest(f, targetPath)
	request["row_groups"] = f.record.DatabaseRowGroups
	before := legacyWriteLocalSnapshot(t, f.db)
	response := performJSONRequest(t, f.router, http.MethodPut, legacyDatabaseConfigPath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	after := legacyWriteLocalSnapshot(t, f.db)
	for _, unchanged := range []string{"targets", "groups", "members", "versions", "migration_maps"} {
		require.Equal(t, before[unchanged], after[unchanged], unchanged)
	}
	current, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, f.record.DatabaseRowGroups, current.DatabaseRowGroups)
	require.Equal(t, f.record.DatabaseTargetRefs, current.DatabaseTargetRefs)
}
