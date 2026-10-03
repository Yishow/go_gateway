package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type managedGroupFixture struct {
	writeGroupRouterFixture
	groups          *workspace.WriteGroupService
	target          *sql.DB
	file, namespace string
	data, group     map[string]any
}

func newManagedGroupFixture(t *testing.T, kind string, tagTypes ...string) managedGroupFixture {
	t.Helper()
	f := managedGroupFixture{writeGroupRouterFixture: newWriteGroupRouterFixture(t), namespace: "main"}
	config := map[string]string{}
	if kind == "sqlite" {
		f.file = filepath.Join(t.TempDir(), "managed.db")
		config["dsn"] = f.file
	} else {
		dsn := os.Getenv("POSTGRES_DSN")
		if dsn == "" {
			t.Skip("owned PostgreSQL fixture not configured")
		}
		fields := map[string]string{}
		for part := range strings.FieldsSeq(dsn) {
			k, v, ok := strings.Cut(part, "=")
			if ok {
				fields[k] = v
			}
		}
		require.Equal(t, "127.0.0.1", fields["host"])
		require.Equal(t, "55432", fields["port"])
		require.Equal(t, "gwtest", fields["dbname"])
		var err error
		f.target, err = sql.Open("pgx", dsn)
		require.NoError(t, err)
		f.namespace = fmt.Sprintf("gw_managed_%d", time.Now().UnixNano())
		_, err = f.target.ExecContext(t.Context(), `CREATE SCHEMA "`+f.namespace+`"`)
		require.NoError(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := f.target.ExecContext(ctx, `DROP SCHEMA "`+f.namespace+`" CASCADE`)
			require.NoError(t, err)
			require.NoError(t, f.target.Close())
		})
		config = map[string]string{"host": fields["host"], "port": fields["port"], "user": fields["user"], "password": fields["password"], "database": fields["dbname"], "sslmode": "disable"}
	}
	config["schema"] = f.namespace
	raw, err := json.Marshal(config)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE database_connectors SET kind=?,connection_config=? WHERE id='connector-A'`, kind, string(raw))
	require.NoError(t, err)
	f.record, err = f.workspace.BindDatabaseConnector(t.Context(), "connector-A")
	require.NoError(t, err)
	if len(tagTypes) > 0 {
		_, err = f.db.ExecContext(t.Context(), `UPDATE tags SET data_type=? WHERE id='tag-A'`, tagTypes[0])
		require.NoError(t, err)
	}
	targetSvc := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	f.groups = workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db)).WithTableInspector(dbtarget.NewReadOnlyTableInspector(targetSvc))
	f.groups.WithManagedTableInspector(f.groups.ManagedTableInspector(targetSvc))
	f.router = NewRouter(&DatalinkServices{
		Workspace: f.workspace, WriteGroups: f.groups, DBTarget: targetSvc,
		RecordingPlan: recordingplan.NewService(recordingplan.NewSQLRepository(f.db)),
		// Metadata registration shares these services with database setup; GET
		// metadata itself uses only the saved workspace, group and connector.
		DBMapping:  dbtarget.NewMappingService(dbtarget.NewSQLTargetMappingRepository(f.db), dbtarget.NewSQLConnectorRepository(f.db), nil),
		SourceRule: sourcerule.NewService(sourcerule.NewSQLRepository(f.db), nil, nil, nil),
	})
	request := f.createRequest()
	g := request["group"].(map[string]any)
	g["destination"].(map[string]any)["storage_strategy"] = "managed"
	g["destination"].(map[string]any)["table_schema"] = f.namespace
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	f.data, f.group = groupSaveData(t, decodeJSONBody(t, created))
	return f
}

func (f managedGroupFixture) preview(t *testing.T) recordingplan.SchemaPreviewToken {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+f.group["id"].(string)+"/schema-preview", f.schemaRequest())
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data recordingplan.SchemaPreviewToken `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	return body.Data
}

func (f managedGroupFixture) schemaRequest() map[string]any {
	request := groupEditRequest(f.writeGroupRouterFixture, f.data, f.group)
	delete(request, "group")
	return request
}

func TestManagedGroupPreparationEmptyNoPlan(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newManagedGroupFixture(t, kind)
			token := f.preview(t)
			if kind == "sqlite" {
				require.NoFileExists(t, f.file)
			}
			require.Equal(t, f.group["id"], token.PlanID)
			require.Equal(t, f.group["revision"], token.PlanRevision)
			require.Len(t, token.Tables, 2)
			require.NotContains(t, strings.Join(token.Statements, "\n"), "gw_record_samples")
			request := f.schemaRequest()
			request["token"], request["operation_id"] = token.Token, token.OperationID
			response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+f.group["id"].(string)+"/schema-apply", request)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			var body struct {
				Data recordingplan.SchemaOperation `json:"data"`
			}
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
			require.Equal(t, recordingplan.SchemaOperationSucceeded, body.Data.Status)
			require.Len(t, body.Data.VerifiedDigest, 64)
			if kind == "sqlite" {
				require.FileExists(t, f.file)
				var err error
				f.target, err = sql.Open("sqlite", f.file)
				require.NoError(t, err)
				t.Cleanup(func() { require.NoError(t, f.target.Close()) })
			}
			var plans int
			require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM recording_plans`).Scan(&plans))
			require.Zero(t, plans)
			var receipts int
			require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT count(*) FROM "`+f.namespace+`".gw_effect_receipts`).Scan(&receipts))
			require.Zero(t, receipts)
			ready, err := f.groups.Readiness(t.Context(), f.group["id"].(string))
			require.NoError(t, err)
			require.True(t, ready.Ready, "%+v", ready.Issues)
			current, err := f.groups.Get(t.Context(), f.group["id"].(string))
			require.NoError(t, err)
			f.data["workspace_revision"] = current.WorkspaceRevision
			noOp := f.preview(t)
			require.Empty(t, noOp.Statements)
			require.Equal(t, recordingplan.NoChangeSchemaCompatible, noOp.NoChangeReason)
			request = f.schemaRequest()
			request["token"], request["operation_id"] = noOp.Token, noOp.OperationID
			response = performJSONRequest(t, f.router, http.MethodPost, f.confirmPath(), request)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
			require.Equal(t, recordingplan.SchemaOperationSucceeded, body.Data.Status)
			require.Zero(t, body.Data.ExecutedStatements)
		})
	}
}

func TestManagedGroupSavesFrozenMetadataBindings(t *testing.T) {
	f := newManagedGroupFixture(t, "sqlite")
	policy := f.group["row_policy"].(map[string]any)
	for key, column := range map[string]string{"record_key_column": "record_id", "group_id_column": "group_id", "device_id_column": "device_id", "bucket_start_column": "bucket_start", "provenance_column": "provenance"} {
		require.Equal(t, column, policy[key], key)
	}
	member := f.group["members"].([]any)[0].(map[string]any)
	require.Regexp(t, `^v_[0-9a-f]{24}$`, member["target_column"])
	require.Equal(t, "receipt", f.group["write_policy"].(map[string]any)["dedupe_capability"])
	require.NoFileExists(t, f.file)
}

func TestManagedGroupSavesCanonicalNameWithEmptyTable(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newManagedGroupFixture(t, kind)
			request := f.createRequest()
			request["expected_workspace_revision"] = f.data["workspace_revision"]
			group := request["group"].(map[string]any)
			destination := group["destination"].(map[string]any)
			destination["storage_strategy"], destination["table_name"] = "managed", ""
			response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
			require.Equal(t, http.StatusCreated, response.Code, response.Body.String())
			_, saved := groupSaveData(t, decodeJSONBody(t, response))
			require.Regexp(t, `^gw_group_[0-9a-f]{32}$`, saved["destination"].(map[string]any)["table_name"])
			assertNoManagedMutation(t, f, kind)
		})
	}
}

func TestManagedGroupUpdateResolvesEmptyTableWithoutDDL(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newManagedGroupFixture(t, kind)
			f.group["destination"].(map[string]any)["table_name"] = ""
			response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+f.group["id"].(string), groupEditRequest(f.writeGroupRouterFixture, f.data, f.group))
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			_, saved := groupSaveData(t, decodeJSONBody(t, response))
			require.Equal(t, f.group["id"], saved["id"])
			require.Regexp(t, `^gw_group_[0-9a-f]{32}$`, saved["destination"].(map[string]any)["table_name"])
			assertNoManagedMutation(t, f, kind)
		})
	}
}
