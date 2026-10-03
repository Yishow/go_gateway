package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func (f managedGroupFixture) confirmPath() string {
	return writeGroupsPath + "/" + f.group["id"].(string) + "/schema-apply"
}

func assertNoManagedMutation(t *testing.T, f managedGroupFixture, kind string) {
	t.Helper()
	if kind == "sqlite" {
		require.NoFileExists(t, f.file)
		return
	}
	var count int
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT count(*) FROM information_schema.tables WHERE table_schema=$1`, f.namespace).Scan(&count))
	require.Zero(t, count)
}

func TestManagedGroupSchemaRejectsChangedScopeBeforeMutation(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			for _, scenario := range []string{"workspace", "group", "source", "connector", "destination", "foreign", "missing-group", "unknown-token", "missing-confirmation", "client-statements", "client-table", "client-dialect"} {
				t.Run(scenario, func(t *testing.T) {
					f := newManagedGroupFixture(t, kind)
					token := f.preview(t)
					request := f.schemaRequest()
					request["token"], request["operation_id"] = token.Token, token.OperationID
					path, want := f.confirmPath(), http.StatusConflict
					switch scenario {
					case "workspace":
						_, err := f.workspace.UpdateDatabaseSetup(t.Context(), request["expected_workspace_revision"].(string), func(context.Context, *sql.Tx, *workspace.Record) error { return nil })
						require.NoError(t, err)
					case "group", "destination":
						current, err := f.groups.Get(t.Context(), f.group["id"].(string))
						require.NoError(t, err)
						if scenario == "group" {
							current.Group.Name = "Another session changed this group"
							current.Group.RowPolicy.IntervalSeconds = 30
						} else {
							current.Group.Destination.TableName = "another_destination"
						}
						_, err = f.groups.Update(t.Context(), current.Group.ID, workspace.WriteGroupMutation{WorkspaceID: f.record.ID, ExpectedWorkspaceRevision: current.WorkspaceRevision, ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: "connector-1", Group: current.Group})
						require.NoError(t, err)
					case "source":
						_, err := f.db.ExecContext(t.Context(), `UPDATE points SET address='40100' WHERE id='point-A'`)
						require.NoError(t, err)
					case "connector":
						_, err := f.db.ExecContext(t.Context(), `UPDATE database_connectors SET identity_revision='connector-2' WHERE id='connector-A'`)
						require.NoError(t, err)
					case "foreign":
						request["workspace_id"], want = "foreign-workspace", http.StatusNotFound
					case "missing-group":
						path, want = writeGroupsPath+"/missing/schema-apply", http.StatusNotFound
					case "unknown-token":
						request["token"], want = "unknown-token", http.StatusNotFound
					case "missing-confirmation":
						delete(request, "operation_id")
						want = http.StatusBadRequest
					case "client-statements":
						request["statements"], want = []string{"DROP TABLE protected"}, http.StatusBadRequest
					case "client-table":
						request["table_name"], want = "foreign_table", http.StatusBadRequest
					case "client-dialect":
						request["dialect"], want = "mysql", http.StatusBadRequest
					}
					response := performJSONRequest(t, f.router, http.MethodPost, path, request)
					require.Equal(t, want, response.Code, response.Body.String())
					if f.file != "" {
						require.NotContains(t, response.Body.String(), f.file)
					}
					assertNoManagedMutation(t, f, kind)
				})
			}
		})
	}
}

func TestManagedGroupSchemaRejectsInternalDestinationBeforePreview(t *testing.T) {
	f := newManagedGroupFixture(t, "sqlite")
	var sequence int
	var name, internalFile string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&sequence, &name, &internalFile))
	config, err := json.Marshal(map[string]string{"dsn": internalFile})
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE database_connectors SET connection_config=? WHERE id='connector-A'`, string(config))
	require.NoError(t, err)
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+f.group["id"].(string)+"/schema-preview", f.schemaRequest())
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "WRITE_GROUP_INTERNAL_DATABASE")
	require.NotContains(t, response.Body.String(), internalFile)
	var count int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM sqlite_master WHERE name IN ('raw_values','gw_effect_receipts')`).Scan(&count))
	require.Zero(t, count)
}

func TestManagedGroupSQLiteRejectsNonMainSchemaBeforeMutation(t *testing.T) {
	f := newManagedGroupFixture(t, "sqlite")
	config, err := json.Marshal(map[string]string{"dsn": f.file, "schema": "other"})
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE database_connectors SET connection_config=? WHERE id='connector-A'`, string(config))
	require.NoError(t, err)
	f.group["destination"].(map[string]any)["table_schema"] = "other"
	saved := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+f.group["id"].(string), groupEditRequest(f.writeGroupRouterFixture, f.data, f.group))
	require.Equal(t, http.StatusOK, saved.Code, saved.Body.String())
	f.data, f.group = groupSaveData(t, decodeJSONBody(t, saved))
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+f.group["id"].(string)+"/schema-preview", f.schemaRequest())
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	require.NoFileExists(t, f.file)
	var operations int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM managed_schema_operations`).Scan(&operations))
	require.Zero(t, operations)
}
