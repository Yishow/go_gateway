package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/api"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// These tests drive the production wiring over HTTP: real configuration
// database, real saved group and connector, real destination database.

const testWriteBase = "/api/v1/datalink/studio-v2/workspace"

func (e *outageEnv) router() http.Handler {
	return api.NewRouter(&api.DatalinkServices{
		Workspace: e.services.workspace, WriteGroups: e.services.writeGroups, RecordingPlan: e.services.recordingPlan,
		DBTarget: e.services.dbTarget, WriteGroupTestWrite: e.services.groupTestWrite,
	})
}

func httpJSON(t *testing.T, router http.Handler, method, path string, body any) (status int, decoded map[string]any) {
	t.Helper()
	var payload bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&payload).Encode(body))
	}
	request := httptest.NewRequestWithContext(t.Context(), method, path, &payload)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	decoded = map[string]any{}
	if recorder.Body.Len() > 0 {
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &decoded), recorder.Body.String())
	}
	return recorder.Code, decoded
}

// createTestWriteGroup saves (never applies) a group on the named destination.
func (e *outageEnv) createTestWriteGroup(t *testing.T, entityColumn, dedupe string) *workspace.WriteGroup {
	t.Helper()
	const name = "A"
	ctx := t.Context()
	record, err := e.services.workspace.GetOrCreate(ctx)
	require.NoError(t, err)
	created, err := e.services.writeGroups.Create(ctx, workspace.WriteGroupMutation{
		WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedConnectorRevision: "connector-1",
		Group: &workspace.WriteGroup{
			WorkspaceID: record.ID, Name: "Test write " + name,
			Members: []workspace.WriteGroupMember{
				{DeviceID: "device-1", PointID: "point-t-" + name, TagID: "tag-t-" + name, TargetColumn: "temperature", Required: true},
				{DeviceID: "device-1", PointID: "point-p-" + name, TagID: "tag-p-" + name, TargetColumn: "pressure", Required: true},
			},
			Destination: workspace.WriteGroupDestination{
				ConnectorID: "connector-" + name, ConnectorRevision: "connector-1", TableSchema: e.schemas[name], TableName: "readings",
				StorageStrategy: workspace.WriteGroupStorageStrategyCustom,
			},
			RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 10, EntityKeyColumn: entityColumn},
			WritePolicy: workspace.WriteGroupWritePolicy{DedupeCapability: dedupe},
		},
	})
	require.NoError(t, err)
	return created.Group
}

type testWriteTarget struct {
	count func(where string) int
}

// driveTestWriteChain is the operator's whole flow: preview, confirm, read the
// shared operation status, and confirm again.
func driveTestWriteChain(t *testing.T, env *outageEnv, groupID string, target testWriteTarget) {
	t.Helper()
	router := env.router()
	groupPath := testWriteBase + "/write-groups/" + groupID

	status, preview := httpJSON(t, router, http.MethodPost, groupPath+"/test-write-preview", nil)
	require.Equal(t, http.StatusOK, status, "%v", preview)
	data := preview["data"].(map[string]any)
	token, operationID := data["token"].(string), data["operation_id"].(string)
	require.Equal(t, "test_write", data["action"])
	require.Equal(t, 1, target.count(`TRUE`), "a preview inserts nothing")

	status, confirmed := httpJSON(t, router, http.MethodPost, groupPath+"/test-write", map[string]any{"token": token, "operation_id": operationID})
	require.Equal(t, http.StatusOK, status, "%v", confirmed)
	result := confirmed["data"].(map[string]any)
	require.Equal(t, "written_verified", result["write_outcome"], "%v", result)
	require.Equal(t, "cleaned", result["cleanup_status"], "%v", result)
	require.Equal(t, 1, target.count(`TRUE`), "only the production row remains")
	require.Equal(t, 1, target.count(`entity = 'line-a' AND note = 'production'`))

	status, stored := httpJSON(t, router, http.MethodGet, testWriteBase+"/database-operations/"+operationID, nil)
	require.Equal(t, http.StatusOK, status, "%v", stored)
	require.Equal(t, "written_verified", stored["data"].(map[string]any)["write_outcome"], "the shared status endpoint serves test-write operations")

	status, again := httpJSON(t, router, http.MethodPost, groupPath+"/test-write", map[string]any{"token": token, "operation_id": operationID})
	require.Equal(t, http.StatusOK, status, "%v", again)
	require.Equal(t, 1, target.count(`TRUE`), "a repeated confirmation never writes again")

	status, wrongKind := httpJSON(t, router, http.MethodPost, groupPath+"/test-write", map[string]any{"token": "tok-unknown", "operation_id": "op-x"})
	require.Equal(t, http.StatusNotFound, status, "%v", wrongKind)
}

func sqliteTestWriteDestinations(counts map[string]*sql.DB) destinationSetup {
	return func(t *testing.T, env *outageEnv, name string) (kind, configJSON, tableSchema string) {
		t.Helper()
		file := filepath.Join(t.TempDir(), "destination-"+name+".db")
		destination, err := sql.Open("sqlite", file)
		require.NoError(t, err)
		destination.SetMaxOpenConns(1)
		t.Cleanup(func() { _ = destination.Close() })
		_, err = destination.ExecContext(t.Context(), `CREATE TABLE readings (entity TEXT, temperature REAL NOT NULL, pressure INTEGER NOT NULL, note TEXT)`)
		require.NoError(t, err)
		_, err = destination.ExecContext(t.Context(), `INSERT INTO readings VALUES ('line-a', 20.5, 7, 'production')`)
		require.NoError(t, err)
		require.NoError(t, dbtarget.CreateEffectReceiptTable(t.Context(), destination, schema.DatabaseConnectorKindSQLite, ""))
		counts[name] = destination
		config, err := json.Marshal(map[string]string{"dsn": file})
		require.NoError(t, err)
		env.targets["connector-"+name] = file
		return "sqlite", string(config), "main"
	}
}

func TestProductionTestWriteOverHTTPOnSQLite(t *testing.T) {
	for _, dedupe := range []string{"", "receipt"} {
		t.Run("dedupe="+dedupe, func(t *testing.T) {
			destinations := map[string]*sql.DB{}
			env := newOutageEnvWith(t, sqliteTestWriteDestinations(destinations))
			group := env.createTestWriteGroup(t, "entity", dedupe)
			count := func(where string) int {
				var n int
				require.NoError(t, destinations["A"].QueryRowContext(t.Context(), `SELECT COUNT(*) FROM readings WHERE `+where).Scan(&n))
				return n
			}
			driveTestWriteChain(t, env, group.ID, testWriteTarget{count: count})
		})
	}
}

func TestProductionTestWriteOverHTTPRefusesUnsafeAndStaleRequests(t *testing.T) {
	destinations := map[string]*sql.DB{}
	env := newOutageEnvWith(t, sqliteTestWriteDestinations(destinations))
	router := env.router()
	count := func() int {
		var n int
		require.NoError(t, destinations["A"].QueryRowContext(t.Context(), `SELECT COUNT(*) FROM readings`).Scan(&n))
		return n
	}

	// A group with no entity key column cannot prove which row is the test's.
	unsafe := env.createTestWriteGroup(t, "", "")
	status, body := httpJSON(t, router, http.MethodPost, testWriteBase+"/write-groups/"+unsafe.ID+"/test-write-preview", nil)
	require.Equal(t, http.StatusUnprocessableEntity, status, "%v", body)
	require.Equal(t, "WRITE_GROUP_TEST_WRITE_UNSUPPORTED", body["error"].(map[string]any)["code"])
	require.Equal(t, 1, count(), "an unsupported preview mutates nothing")

	// Editing the group after the preview makes the confirmation stale.
	safe := env.createTestWriteGroup(t, "entity", "")
	status, preview := httpJSON(t, router, http.MethodPost, testWriteBase+"/write-groups/"+safe.ID+"/test-write-preview", nil)
	require.Equal(t, http.StatusOK, status, "%v", preview)
	data := preview["data"].(map[string]any)
	record, err := env.services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	edited := *safe
	edited.Name = "renamed after preview"
	_, err = env.services.writeGroups.Update(t.Context(), safe.ID, workspace.WriteGroupMutation{
		WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedConnectorRevision: "connector-1",
		ExpectedGroupRevision: safe.Revision, Group: &edited,
	})
	require.NoError(t, err)
	status, stale := httpJSON(t, router, http.MethodPost, testWriteBase+"/write-groups/"+safe.ID+"/test-write",
		map[string]any{"token": data["token"], "operation_id": data["operation_id"]})
	require.Equal(t, http.StatusConflict, status, "%v", stale)
	require.Equal(t, "WRITE_GROUP_TEST_WRITE_PREVIEW_STALE", stale["error"].(map[string]any)["code"])
	require.Equal(t, 1, count(), "a stale confirmation writes nothing")

	// A schema-apply token is a different action and cannot authorize a test write.
	schemaToken, err := env.services.recordingPlan.PrepareSchemaPreview(t.Context(), recordingplan.SchemaPreviewScope{
		WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, PlanID: "plan-1", PlanRevision: "rev-1",
		ConnectorID: "connector-A", ConnectorRevision: "connector-1", Dialect: "sqlite", Database: env.targets["connector-A"], Schema: "main", TablePrefix: "gw_record_",
	}, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
		return recordingplan.TargetTableInspection{Status: "missing"}, nil
	})
	require.NoError(t, err)
	status, wrong := httpJSON(t, router, http.MethodPost, testWriteBase+"/write-groups/"+safe.ID+"/test-write",
		map[string]any{"token": schemaToken.Token, "operation_id": schemaToken.OperationID})
	require.Equal(t, http.StatusUnprocessableEntity, status, "%v", wrong)
	require.Equal(t, "WRITE_GROUP_TEST_WRITE_TOKEN_KIND", wrong["error"].(map[string]any)["code"])
	require.Equal(t, 1, count(), "a wrong-kind token writes nothing")
}

func TestProductionTestWriteOverHTTPOnPostgres(t *testing.T) {
	dsn, fields := postgresEnvFields(t)
	admin, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = admin.Close() })
	require.NoError(t, admin.PingContext(t.Context()))

	schemas := map[string]string{}
	setup := func(t *testing.T, env *outageEnv, name string) (kind, configJSON, tableSchema string) {
		t.Helper()
		pgSchema := fmt.Sprintf("gw_tw_%s_%d", strings.ToLower(name), time.Now().UnixNano())
		_, err := admin.ExecContext(t.Context(), `CREATE SCHEMA "`+pgSchema+`"`)
		require.NoError(t, err)
		t.Cleanup(func() { _, _ = admin.ExecContext(context.Background(), `DROP SCHEMA "`+pgSchema+`" CASCADE`) })
		_, err = admin.ExecContext(t.Context(), `CREATE TABLE "`+pgSchema+`".readings (entity TEXT, temperature DOUBLE PRECISION NOT NULL, pressure BIGINT NOT NULL, note TEXT)`)
		require.NoError(t, err)
		_, err = admin.ExecContext(t.Context(), `INSERT INTO "`+pgSchema+`".readings VALUES ('line-a', 20.5, 7, 'production')`)
		require.NoError(t, err)
		require.NoError(t, dbtarget.CreateEffectReceiptTable(t.Context(), admin, schema.DatabaseConnectorKindPostgres, pgSchema))
		port := fields["port"]
		if port == "" {
			port = "5432"
		}
		schemas[name] = pgSchema
		config, err := json.Marshal(map[string]string{
			"host": fields["host"], "port": port, "user": fields["user"], "password": fields["password"],
			"database": fields["dbname"], "sslmode": "disable",
		})
		require.NoError(t, err)
		return "postgres", string(config), pgSchema
	}
	env := newOutageEnvWith(t, setup)
	for _, dedupe := range []string{"", "receipt"} {
		group := env.createTestWriteGroup(t, "entity", dedupe)
		count := func(where string) int {
			var n int
			require.NoError(t, admin.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM "`+schemas["A"]+`".readings WHERE `+where).Scan(&n))
			return n
		}
		driveTestWriteChain(t, env, group.ID, testWriteTarget{count: count})
	}
}

func TestProductionGroupApplyOverHTTP(t *testing.T) {
	destinations := map[string]*sql.DB{}
	env := newOutageEnvWith(t, sqliteTestWriteDestinations(destinations))
	router := env.router()
	group := env.createTestWriteGroup(t, "entity", "")
	record, err := env.services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	request := map[string]any{
		"workspace_id": record.ID, "expected_workspace_revision": record.DatabaseSetupRevision,
		"expected_group_revision": group.Revision, "expected_connector_revision": "connector-1",
	}
	path := testWriteBase + "/write-groups/" + group.ID + "/apply"

	status, applied := httpJSON(t, router, http.MethodPost, path, request)
	require.Equal(t, http.StatusOK, status, "%v", applied)
	saved := applied["data"].(map[string]any)["group"].(map[string]any)
	require.Equal(t, saved["revision"], saved["applied_revision"], "the saved draft is now the applied revision")
	require.Equal(t, "ready", saved["status"])

	// The same request again is stale: the workspace revision moved with the apply.
	status, stale := httpJSON(t, router, http.MethodPost, path, request)
	require.Equal(t, http.StatusConflict, status, "%v", stale)
	require.Equal(t, "revision_mismatch", stale["error"].(map[string]any)["code"])
}
