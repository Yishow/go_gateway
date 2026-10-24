package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

var longTablePattern = regexp.MustCompile(`^gw_long_[0-9a-f]{16}$`)

// ensureLong saves one canonical managed_long_v1 group through the public route.
func ensureLong(t *testing.T, f managedGroupFixture, kind, unit string, interval int, source string) map[string]any {
	t.Helper()
	request := profileRequest(t, f.writeGroupRouterFixture, kind, interval, unit, "main", source)
	group := request["group"].(map[string]any)
	group["layout_version"] = "managed_long_v1"
	group["destination"].(map[string]any)["table_schema"] = f.namespace
	response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, response))
	return saved
}

// scopeBody is the revision-only body of schema preview and confirmation for
// the group's current saved state.
func scopeBody(t *testing.T, f managedGroupFixture, groupID string) map[string]any {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+groupID, nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, response))
	request := groupEditRequest(f.writeGroupRouterFixture, data, group)
	delete(request, "group")
	return request
}

func previewLong(t *testing.T, f managedGroupFixture, groupID string) recordingplan.SchemaPreviewToken {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-preview", scopeBody(t, f, groupID))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data recordingplan.SchemaPreviewToken `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	return body.Data
}

func confirmLong(t *testing.T, f managedGroupFixture, groupID string, token recordingplan.SchemaPreviewToken) recordingplan.SchemaOperation {
	t.Helper()
	request := scopeBody(t, f, groupID)
	request["token"], request["operation_id"] = token.Token, token.OperationID
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-apply", request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data recordingplan.SchemaOperation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	return body.Data
}

func destinationOf(group map[string]any) map[string]any { return group["destination"].(map[string]any) }

func columnsOf(t *testing.T, target *sql.DB, kind, namespace, table string) []string {
	t.Helper()
	query, args := `SELECT name FROM pragma_table_info(?) ORDER BY cid`, []any{table}
	if kind == "postgres" {
		query = `SELECT column_name FROM information_schema.columns WHERE table_schema=$1 AND table_name=$2 ORDER BY ordinal_position`
		args = []any{namespace, table}
	}
	rows, err := target.QueryContext(t.Context(), query, args...)
	require.NoError(t, err)
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		require.NoError(t, rows.Scan(&name))
		names = append(names, name)
	}
	require.NoError(t, rows.Err())
	return names
}

func TestSharedStorageConfirmationScope(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newManagedDestinationFixture(t, kind)
			seedProfileSource(t, f.writeGroupRouterFixture, "energy")
			seedProfileSource(t, f.writeGroupRouterFixture, "voltage")

			power := ensureLong(t, f, "instantaneous", "second", 15, "")
			energy := ensureLong(t, f, "cumulative", "minute", 10, "energy")
			require.Equal(t, "managed_long_v1", power["layout_version"])
			require.Equal(t, "plain", power["storage_mode"])
			table := destinationOf(power)["table_name"].(string)
			require.Regexp(t, longTablePattern, table, "the table is named by storage scope, not by group")
			require.Equal(t, table, destinationOf(energy)["table_name"], "two profiles of one scope share one table")
			require.NotEqual(t, power["id"], energy["id"])

			// Group A prepares the shared table: table, composite key, index, receipts.
			powerToken := previewLong(t, f, power["id"].(string))
			require.Len(t, powerToken.Statements, 3, strings.Join(powerToken.Statements, "\n"))
			require.Contains(t, powerToken.Statements[0], `PRIMARY KEY ("bucket_start", "record_id")`)
			require.Contains(t, powerToken.Statements[1], "CREATE INDEX")
			for _, statement := range powerToken.Statements {
				require.NotRegexp(t, `(?i)retention|compress|delete|drop|timescale|create_hypertable`, statement)
			}
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, power["id"].(string), powerToken).Status)
			if kind == "sqlite" {
				var err error
				f.target, err = sql.Open("sqlite", f.file)
				require.NoError(t, err)
				t.Cleanup(func() { _ = f.target.Close() })
			}
			before := columnsOf(t, f.target, kind, f.namespace, table)
			require.Contains(t, before, "bucket_start")
			for _, column := range before {
				require.NotRegexp(t, `^v_[0-9a-f]{24}$`, column)
			}

			// Group B shares the table but is not confirmed by A's confirmation.
			powerReady, err := f.groups.Readiness(t.Context(), power["id"].(string))
			require.NoError(t, err)
			require.True(t, powerReady.ConfigReady)
			require.True(t, powerReady.SchemaReady, "%+v", powerReady.Issues)
			require.True(t, powerReady.Ready, "a confirmed plain instantaneous long group is a runtime-ready writer: %+v", powerReady.Issues)
			require.NotContains(t, issueCodes(powerReady.Issues), "long-runtime-unavailable")
			energyReady, err := f.groups.Readiness(t.Context(), energy["id"].(string))
			require.NoError(t, err)
			require.Contains(t, issueCodes(energyReady.Issues), "schema-confirmation-required")

			// An empty DDL preview must still be explicitly confirmed, per group.
			energyToken := previewLong(t, f, energy["id"].(string))
			require.Empty(t, energyToken.Statements)
			require.Equal(t, recordingplan.NoChangeSchemaCompatible, energyToken.NoChangeReason)
			energyOp := confirmLong(t, f, energy["id"].(string), energyToken)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, energyOp.Status)
			require.Zero(t, energyOp.ExecutedStatements)
			savedPower := getGroup(t, f, power["id"].(string))
			savedEnergy := getGroup(t, f, energy["id"].(string))
			require.NotEmpty(t, destinationOf(savedPower)["schema_revision"])
			require.NotEmpty(t, destinationOf(savedEnergy)["schema_revision"])
			require.NotEqual(t, destinationOf(savedPower)["schema_revision"], destinationOf(savedEnergy)["schema_revision"],
				"each group keeps its own revision-bound confirmation")

			// A new source in group A adds records, never columns, and only A needs a fresh preview.
			updated := cloneMap(savedPower)
			updated["members"] = append(updated["members"].([]any), map[string]any{
				"device_id": "device-A", "point_id": "point-voltage", "tag_id": "tag-voltage", "required": true,
				"circuit_id": "main", "circuit_name": "Circuit main",
			})
			data := getGroupData(t, f, power["id"].(string))
			response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+power["id"].(string), groupEditRequest(f.writeGroupRouterFixture, data, updated))
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			grownPower := getGroup(t, f, power["id"].(string))
			require.Empty(t, destinationOf(grownPower)["schema_revision"], "the changed group needs fresh confirmation")
			grownReady, err := f.groups.Readiness(t.Context(), power["id"].(string))
			require.NoError(t, err)
			require.Contains(t, issueCodes(grownReady.Issues), "schema-confirmation-required")
			require.Equal(t, destinationOf(savedEnergy)["schema_revision"], destinationOf(getGroup(t, f, energy["id"].(string)))["schema_revision"],
				"group B keeps its confirmation and applied state")
			growToken := previewLong(t, f, power["id"].(string))
			require.Empty(t, growToken.Statements, "a new point needs no DDL")
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, power["id"].(string), growToken).Status)
			require.Equal(t, before, columnsOf(t, f.target, kind, f.namespace, table), "no value column was added")
		})
	}
}

func TestSharedStorageScopeIsolationAndForeignOwner(t *testing.T) {
	f := newManagedDestinationFixture(t, "sqlite")
	seedProfileSource(t, f.writeGroupRouterFixture, "energy")
	first := ensureLong(t, f, "instantaneous", "second", 15, "")
	table := destinationOf(first)["table_name"].(string)

	t.Run("another connector is another scope with another table", func(t *testing.T) {
		other := filepath.Join(t.TempDir(), "other.db")
		_, err := f.db.ExecContext(t.Context(), `INSERT INTO database_connectors (id,name,kind,connection_config,identity_revision,status,enabled)
			VALUES ('connector-B','Other','sqlite',?,'connector-1','ready',1)`, `{"dsn":"`+filepath.ToSlash(other)+`"}`)
		require.NoError(t, err)
		request := profileRequest(t, f.writeGroupRouterFixture, "cumulative", 10, "minute", "main", "energy")
		group := request["group"].(map[string]any)
		group["layout_version"] = "managed_long_v1"
		group["destination"].(map[string]any)["connector_id"] = "connector-B"
		response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		_, saved := groupSaveData(t, decodeJSONBody(t, response))
		require.Regexp(t, longTablePattern, destinationOf(saved)["table_name"])
		require.NotEqual(t, table, destinationOf(saved)["table_name"])
	})

	t.Run("a foreign table on the shared name fails closed before any DDL", func(t *testing.T) {
		foreign, err := sql.Open("sqlite", f.file)
		require.NoError(t, err)
		t.Cleanup(func() { _ = foreign.Close() })
		_, err = foreign.ExecContext(t.Context(), `CREATE TABLE "`+table+`" (bucket_start TEXT NOT NULL, record_id TEXT NOT NULL, PRIMARY KEY (bucket_start, record_id))`)
		require.NoError(t, err)
		response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+first["id"].(string)+"/schema-preview", scopeBody(t, f, first["id"].(string)))
		require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, response), "RECORDING_SCHEMA_INCOMPATIBLE")
		var tables int
		require.NoError(t, foreign.QueryRowContext(t.Context(), `SELECT count(*) FROM sqlite_master WHERE type='table'`).Scan(&tables))
		require.Equal(t, 1, tables, "the foreign table is untouched and nothing was created")
	})
}

func issueCodes(issues []workspace.ReadinessIssue) []string {
	codes := make([]string, len(issues))
	for i, issue := range issues {
		codes[i] = issue.Code
	}
	return codes
}

func getGroupData(t *testing.T, f managedGroupFixture, id string) map[string]any {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	data, _ := groupSaveData(t, decodeJSONBody(t, response))
	return data
}

func getGroup(t *testing.T, f managedGroupFixture, id string) map[string]any {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, response))
	return group
}

func cloneMap(source map[string]any) map[string]any {
	encoded, _ := json.Marshal(source)
	var cloned map[string]any
	_ = json.Unmarshal(encoded, &cloned)
	return cloned
}

// TestSharedStorageUnusableIndexBlocksConfirmation reproduces the catalog
// states of a failed, unfinished or dropping index on the disposable PostgreSQL
// fixture. An index with the confirmed name and columns but without usable
// state must fail preview and readiness, and a restored index passes again.
func TestSharedStorageUnusableIndexBlocksConfirmation(t *testing.T) {
	for _, flag := range []string{"indisvalid", "indisready", "indislive"} {
		t.Run(flag, func(t *testing.T) {
			f := newManagedDestinationFixture(t, "postgres")
			group := ensureLong(t, f, "instantaneous", "second", 15, "")
			groupID := group["id"].(string)
			table := destinationOf(group)["table_name"].(string)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, groupID, previewLong(t, f, groupID)).Status)
			ready, err := f.groups.Readiness(t.Context(), groupID)
			require.NoError(t, err)
			require.True(t, ready.SchemaReady, "%+v", ready.Issues)

			setIndexFlag := func(value bool) {
				_, err := f.target.ExecContext(t.Context(), `UPDATE pg_catalog.pg_index SET `+flag+` = $1
					WHERE indexrelid = (quote_ident($2) || '.' || quote_ident($3))::regclass`, value, f.namespace, "idx_"+table+"_src_time")
				require.NoError(t, err)
			}
			setIndexFlag(false)
			response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-preview", scopeBody(t, f, groupID))
			require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
			assertWriteGroupError(t, decodeJSONBody(t, response), "RECORDING_SCHEMA_INCOMPATIBLE")
			blocked, err := f.groups.Readiness(t.Context(), groupID)
			require.NoError(t, err)
			require.False(t, blocked.SchemaReady, "an unusable index must not satisfy the confirmed layout")
			require.Contains(t, issueCodes(blocked.Issues), "managed-schema-unverified")

			setIndexFlag(true)
			restored := previewLong(t, f, groupID)
			require.Empty(t, restored.Statements)
			recovered, err := f.groups.Readiness(t.Context(), groupID)
			require.NoError(t, err)
			require.True(t, recovered.SchemaReady, "%+v", recovered.Issues)
		})
	}
}
