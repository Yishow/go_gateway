package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// roleDestination moves the fixture's saved connector to a throwaway database
// where TimescaleDB lives in the non-public schema tsext, and connects it as a
// non-superuser role that may CREATE in the target schema but has no USAGE on
// it. f.target becomes a superuser handle on that database, for catalog reads
// and for granting or revoking the role's rights.
func roleDestination(t *testing.T, f *managedGroupFixture) (role string) {
	t.Helper()
	var available int
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM pg_available_extensions WHERE name = 'timescaledb'`).Scan(&available))
	if available == 0 {
		t.Skip("this PostgreSQL does not offer the timescaledb extension; the role case is NOT RUN here")
	}
	fields := map[string]string{}
	for part := range strings.FieldsSeq(os.Getenv("POSTGRES_DSN")) {
		if k, v, ok := strings.Cut(part, "="); ok {
			fields[k] = v
		}
	}
	base := f.target
	stamp := time.Now().UnixNano()
	role = fmt.Sprintf("gwrole_user_%d", stamp)
	database := fmt.Sprintf("gwrole_db_%d", stamp)
	// The role is cluster-wide: register its removal first, so it runs after
	// the database that granted it rights has been dropped.
	_, err := base.ExecContext(t.Context(), `CREATE ROLE "`+role+`" LOGIN PASSWORD 'tsrole'`)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = base.ExecContext(context.Background(), `DROP ROLE IF EXISTS "`+role+`"`) })
	_, err = base.ExecContext(t.Context(), `CREATE DATABASE "`+database+`" TEMPLATE template0`)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = base.ExecContext(context.Background(), `DROP DATABASE IF EXISTS "`+database+`" WITH (FORCE)`)
	})

	admin, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", fields["host"], fields["port"], fields["user"], fields["password"], database))
	require.NoError(t, err)
	t.Cleanup(func() { _ = admin.Close() })
	for _, statement := range []string{
		`CREATE SCHEMA tsext`, `CREATE EXTENSION timescaledb SCHEMA tsext`,
		`CREATE SCHEMA "` + f.namespace + `"`,
		`GRANT USAGE ON SCHEMA tsext TO "` + role + `"`,
		`GRANT CREATE ON SCHEMA "` + f.namespace + `" TO "` + role + `"`,
	} {
		_, err := admin.ExecContext(t.Context(), statement)
		require.NoError(t, err, statement)
	}

	config, err := json.Marshal(map[string]string{
		"host": fields["host"], "port": fields["port"], "user": role, "password": "tsrole",
		"database": database, "sslmode": "disable", "schema": f.namespace,
	})
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE database_connectors SET connection_config=? WHERE id='connector-A'`, string(config))
	require.NoError(t, err)
	f.record, err = f.workspace.BindDatabaseConnector(t.Context(), "connector-A")
	require.NoError(t, err)
	f.target = admin
	return role
}

// TestTimescaleRoleNeedsTargetSchemaUsage runs the real flow as a non-superuser
// role against an extension installed outside public. A role that may CREATE in
// the target schema but cannot USE it is reported as timescale_unavailable at
// preview and again at confirmation, with no DDL; once the repair (USAGE) is
// made the same flow previews, confirms, applies and writes a record.
func TestTimescaleRoleNeedsTargetSchemaUsage(t *testing.T) {
	f := newManagedDestinationFixture(t, "postgres")
	role := roleDestination(t, &f)
	_, created := ensureTimescaleGroup(t, f)
	_, group := groupSaveData(t, created)
	groupID, table := group["id"].(string), destinationOf(group)["table_name"].(string)
	grant := func(statement string) {
		t.Helper()
		_, err := f.target.ExecContext(t.Context(), strings.NewReplacer("{ns}", `"`+f.namespace+`"`, "{role}", `"`+role+`"`).Replace(statement))
		require.NoError(t, err, statement)
	}
	requireUnavailable := func(response *httptest.ResponseRecorder) {
		t.Helper()
		require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
		failure := errorOf(t, decodeJSONBody(t, response))
		require.Equal(t, "timescale_unavailable", failure["code"])
		require.Equal(t, "repair_timescale_capability", failure["action"])
		require.Contains(t, failure["message"], "schema-usage-denied")
		require.NotContains(t, response.Body.String(), "tsrole", "no credential is exposed")
	}
	untouched := func(relations, extensions int) {
		t.Helper()
		afterRelations, afterExtensions := catalogObjects(t, f)
		require.Equal(t, [2]int{relations, extensions}, [2]int{afterRelations, afterExtensions}, "no table was created and no extension installed")
	}

	// 1. CREATE without USAGE: unavailable at preview, nothing created.
	relations, extensions := catalogObjects(t, f)
	requireUnavailable(performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-preview", scopeBody(t, f, groupID)))
	untouched(relations, extensions)
	require.Equal(t, "timescale", getGroup(t, f, groupID)["storage_mode"], "the selection is never rewritten to plain")

	// 2. A token prepared while the role could use the schema is refused at
	// confirmation once that right is gone, before any statement runs.
	grant(`GRANT USAGE ON SCHEMA {ns} TO {role}`)
	token := previewLong(t, f, groupID)
	require.Len(t, token.Statements, 4)
	require.Contains(t, token.Statements[1], `"tsext"."create_hypertable"`, "the extension outside public is addressed by its own schema")
	grant(`REVOKE USAGE ON SCHEMA {ns} FROM {role}`)
	request := scopeBody(t, f, groupID)
	request["token"], request["operation_id"] = token.Token, token.OperationID
	requireUnavailable(performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-apply", request))
	untouched(relations, extensions)
	require.Empty(t, destinationOf(getGroup(t, f, groupID))["schema_revision"], "an unconfirmed result is not saved as confirmed")

	// 3. Repaired: the whole flow works for the non-superuser role.
	grant(`GRANT USAGE ON SCHEMA {ns} TO {role}`)
	operation := confirmLong(t, f, groupID, previewLong(t, f, groupID))
	require.Equal(t, "succeeded", string(operation.Status))
	var owner, extensionSchema string
	require.NoError(t, f.target.QueryRowContext(t.Context(), `
		SELECT pg_get_userbyid(c.relowner), (SELECT n.nspname FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname = 'timescaledb')
		FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = $1 AND c.relname = $2`, f.namespace, table).Scan(&owner, &extensionSchema))
	require.Equal(t, role, owner, "the hypertable is owned by the connected non-superuser role")
	require.Equal(t, "tsext", extensionSchema)
	require.Equal(t, "604800000000.000000", chunkIntervalMicros(t, f, table))

	ready, err := f.groups.Readiness(t.Context(), groupID)
	require.NoError(t, err)
	require.True(t, ready.Ready, "%+v", ready.Issues)
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", scopeBody(t, f, groupID))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	writeOneSampleThroughPipeline(t, f, groupID, table)
}
