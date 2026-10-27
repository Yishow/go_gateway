package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// TestSharedDestinationCutoverApplyCAS exercises the common saved-group
// Update/Apply schedule with local SQLite targets. SQL Server's frozen outbox
// proof is covered separately by the scripted pipeline regression.
func TestSharedDestinationCutoverApplyCAS(t *testing.T) {
	f := newManagedDestinationFixture(t, "sqlite")
	group := ensureLong(t, f, "instantaneous", "second", 15, "")
	groupID := group["id"].(string)
	token := previewLong(t, f, groupID)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, groupID, token).Status)
	initialApply := sqlServerCutoverApplyBody(t, f, groupID)
	resp := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", initialApply)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	first, err := f.groups.ResolveAppliedAt(t.Context(), groupID, time.Now().Add(24*time.Hour))
	require.NoError(t, err)
	require.Equal(t, "connector-A", first.Group.Destination.ConnectorID)

	const nextConnector, nextRevision = "connector-cutover", "connector-cutover-v1"
	config, err := json.Marshal(map[string]string{"dsn": f.file, "schema": f.namespace})
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `INSERT INTO database_connectors
		(id, name, kind, status, connection_config, identity_revision, enabled)
		VALUES (?, 'Cutover fixture', 'sqlite', 'ready', ?, ?, 1)`, nextConnector, string(config), nextRevision)
	require.NoError(t, err)

	edit := sqlServerCutoverGroupEdit(t, f, groupID, nextConnector, nextRevision, "")
	for _, field := range []string{"expected_workspace_revision", "expected_group_revision", "expected_connector_revision"} {
		fresh := edit[field]
		edit[field] = "stale-sqlserver-cutover-revision"
		stale := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+groupID, edit)
		require.Equal(t, http.StatusConflict, stale.Code, stale.Body.String())
		edit[field] = fresh
	}
	resp = performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+groupID, edit)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	_, moved := groupSaveData(t, decodeJSONBody(t, resp))
	require.Equal(t, nextConnector, destinationOf(moved)["connector_id"])
	require.NotEqual(t, destinationOf(group)["table_name"], destinationOf(moved)["table_name"])

	sqlServerCutoverConfirmSchema(t, f, groupID)
	f.groups.WithClock(func() time.Time { return first.EffectiveAt.Add(time.Hour) })
	t.Cleanup(func() { f.groups.WithClock(nil) })
	apply := sqlServerCutoverApplyBody(t, f, groupID)
	for _, field := range []string{"expected_workspace_revision", "expected_group_revision", "expected_connector_revision"} {
		fresh := apply[field]
		apply[field] = "stale-sqlserver-cutover-revision"
		stale := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", apply)
		require.Equal(t, http.StatusConflict, stale.Code, stale.Body.String())
		apply[field] = fresh
	}
	resp = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", apply)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	second, err := f.groups.ResolveAppliedAt(t.Context(), groupID, first.EffectiveAt.Add(2*time.Hour))
	require.NoError(t, err)
	require.Equal(t, nextConnector, second.Group.Destination.ConnectorID)
	beforeSecond, err := f.groups.ResolveAppliedAt(t.Context(), groupID, second.EffectiveAt.Add(-time.Nanosecond))
	require.NoError(t, err)
	require.Equal(t, "connector-A", beforeSecond.Group.Destination.ConnectorID)
	require.Equal(t, nextConnector, second.Group.Destination.ConnectorID, "the changed scope begins exactly at its Apply boundary")

	// Return to the original connector through another saved revision and Apply;
	// neither destination is edited in place and the second boundary is explicit.
	returnEdit := sqlServerCutoverGroupEdit(t, f, groupID, "connector-A", "connector-1", "")
	resp = performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+groupID, returnEdit)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	sqlServerCutoverConfirmSchema(t, f, groupID)
	f.groups.WithClock(func() time.Time { return second.EffectiveAt.Add(time.Hour) })
	returnApply := sqlServerCutoverApplyBody(t, f, groupID)
	resp = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", returnApply)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	third, err := f.groups.ResolveAppliedAt(t.Context(), groupID, second.EffectiveAt.Add(2*time.Hour))
	require.NoError(t, err)
	require.Equal(t, "connector-A", third.Group.Destination.ConnectorID)
	beforeReturn, err := f.groups.ResolveAppliedAt(t.Context(), groupID, third.EffectiveAt.Add(-time.Nanosecond))
	require.NoError(t, err)
	require.Equal(t, nextConnector, beforeReturn.Group.Destination.ConnectorID)
}

func sqlServerCutoverGroupEdit(
	t *testing.T, f managedGroupFixture, groupID, connectorID, connectorRevision, database string,
) map[string]any {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+groupID, nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, response))
	destination := destinationOf(group)
	destination["connector_id"], destination["connector_revision"] = connectorID, connectorRevision
	destination["database"], destination["table_schema"], destination["table_name"] = database, f.namespace, ""
	destination["schema_revision"], destination["schema_digest"] = "", ""
	request := groupEditRequest(f.writeGroupRouterFixture, data, group)
	request["expected_connector_revision"] = connectorRevision
	return request
}

func sqlServerCutoverApplyBody(t *testing.T, f managedGroupFixture, groupID string) map[string]any {
	t.Helper()
	body := scopeBody(t, f, groupID)
	response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+groupID, nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, response))
	body["expected_connector_revision"] = destinationOf(group)["connector_revision"]
	delete(body, "group")
	return body
}

func sqlServerCutoverConfirmSchema(t *testing.T, f managedGroupFixture, groupID string) {
	t.Helper()
	body := sqlServerCutoverApplyBody(t, f, groupID)
	preview := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-preview", body)
	require.Equal(t, http.StatusOK, preview.Code, preview.Body.String())
	var tokenEnvelope struct {
		Data recordingplan.SchemaPreviewToken `json:"data"`
	}
	require.NoError(t, json.Unmarshal(preview.Body.Bytes(), &tokenEnvelope))
	body["token"], body["operation_id"] = tokenEnvelope.Data.Token, tokenEnvelope.Data.OperationID
	confirmed := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-apply", body)
	require.Equal(t, http.StatusOK, confirmed.Code, confirmed.Body.String())
}

// TestSQLServerAdapterCoexistence keeps existing custom and database-mode
// configurations intact while a separate SQL Server managed-long draft is
// saved. All connector data is local configuration; no external connector opens.
func TestSQLServerAdapterCoexistence(t *testing.T) {
	for _, tc := range []struct {
		name, kind, config, storage string
	}{
		{name: "sqlite custom", kind: "sqlite", config: `{"dsn":"local-only.db"}`},
		{name: "postgres managed", kind: "postgres", config: `{"host":"127.0.0.1","database":"legacy","schema":"public"}`, storage: "plain"},
		{name: "timescale", kind: "postgres", config: `{"host":"127.0.0.1","database":"legacy","schema":"public"}`, storage: "timescale"},
		{name: "mysql connector", kind: "mysql", config: `{"host":"127.0.0.1","database":"legacy","server_family":"mysql"}`},
		{name: "mariadb connector", kind: "mysql", config: `{"host":"127.0.0.1","database":"legacy","server_family":"mariadb"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newWriteGroupRouterFixture(t)
			createdCustom := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
			require.Equal(t, http.StatusCreated, createdCustom.Code, createdCustom.Body.String())
			_, customGroup := groupSaveData(t, decodeJSONBody(t, createdCustom))
			customBefore := pick(customGroup, savedFacts)

			var legacyModeGroup map[string]any
			var legacyModeBefore map[string]any
			if tc.name != "sqlite custom" {
				const legacyConnector, legacyRevision = "connector-legacy-mode", "legacy-mode-v1"
				_, err := f.db.ExecContext(t.Context(), `INSERT INTO database_connectors
					(id, name, kind, status, connection_config, identity_revision, enabled)
					VALUES (?, 'Legacy mode fixture', ?, 'ready', ?, ?, 1)`, legacyConnector, tc.kind, tc.config, legacyRevision)
				require.NoError(t, err)
				if tc.storage != "" {
					request := profileRequest(t, f, "instantaneous", 15, "second", "legacy-mode", "")
					group := request["group"].(map[string]any)
					group["layout_version"], group["storage_mode"] = "managed_long_v1", tc.storage
					destination := destinationOf(group)
					destination["connector_id"], destination["connector_revision"] = legacyConnector, legacyRevision
					if tc.kind == "postgres" {
						destination["database"], destination["table_schema"] = "legacy", "public"
					} else {
						destination["database"], destination["table_schema"] = "legacy", "legacy"
					}
					request["expected_connector_revision"] = legacyRevision
					created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
					require.Equal(t, http.StatusOK, created.Code, created.Body.String())
					_, legacyModeGroup = groupSaveData(t, decodeJSONBody(t, created))
					legacyModeBefore = pick(legacyModeGroup, savedFacts)
				}
			}

			_, err := f.db.ExecContext(t.Context(), `INSERT INTO database_connectors
				(id, name, kind, status, connection_config, identity_revision, enabled)
				VALUES ('connector-sqlserver', 'SQL Server fixture', 'sqlserver', 'ready',
				'{"host":"fixture.invalid","database":"recordings","schema":"dbo"}', 'sqlserver-v1', 1)`)
			require.NoError(t, err)
			request := profileRequest(t, f, "instantaneous", 15, "second", "sqlserver-new", "")
			group := request["group"].(map[string]any)
			group["layout_version"], group["storage_mode"] = "managed_long_v1", "plain"
			destination := destinationOf(group)
			destination["connector_id"], destination["connector_revision"] = "connector-sqlserver", "sqlserver-v1"
			destination["database"], destination["table_schema"] = "recordings", "dbo"
			request["expected_connector_revision"] = "sqlserver-v1"
			saved := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
			require.Equal(t, http.StatusOK, saved.Code, saved.Body.String())
			_, sqlServerGroup := groupSaveData(t, decodeJSONBody(t, saved))
			require.Equal(t, "managed_long_v1", sqlServerGroup["layout_version"])
			require.Equal(t, "plain", sqlServerGroup["storage_mode"])
			require.Equal(t, "dbo", destinationOf(sqlServerGroup)["table_schema"])
			require.Regexp(t, longTablePattern, destinationOf(sqlServerGroup)["table_name"])

			wrapper := managedGroupFixture{writeGroupRouterFixture: f}
			currentCustom := getGroup(t, wrapper, customGroup["id"].(string))
			require.Equal(t, customBefore, pick(currentCustom, savedFacts), "SQL Server setup does not rewrite prior custom settings")
			if legacyModeGroup != nil {
				currentMode := getGroup(t, wrapper, legacyModeGroup["id"].(string))
				require.Equal(t, legacyModeBefore, pick(currentMode, savedFacts), "SQL Server setup does not rewrite prior database mode")
			}
			if tc.kind != "sqlite" {
				var savedKind, savedConfig string
				require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT kind, connection_config FROM database_connectors WHERE id = 'connector-legacy-mode'`).Scan(&savedKind, &savedConfig))
				require.Equal(t, tc.kind, savedKind)
				require.JSONEq(t, tc.config, savedConfig, "the saved PostgreSQL/MySQL/MariaDB connector remains untouched")
			}
		})
	}
}
