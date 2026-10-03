package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func (f *managedGroupFixture) openTarget(t *testing.T) {
	t.Helper()
	if f.target != nil {
		return
	}
	var err error
	f.target, err = sql.Open("sqlite", f.file)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, f.target.Close()) })
}

func TestManagedGroupSchemaPreservesUnownedConflictingTable(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newManagedGroupFixture(t, kind)
			f.openTarget(t)
			table := `"` + f.namespace + `"."raw_values"`
			_, err := f.target.ExecContext(t.Context(), `CREATE TABLE `+table+` (sentinel TEXT NOT NULL)`)
			require.NoError(t, err)
			_, err = f.target.ExecContext(t.Context(), `INSERT INTO `+table+` VALUES ('preserve-existing-data')`)
			require.NoError(t, err)
			response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+f.group["id"].(string)+"/schema-preview", f.schemaRequest())
			require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
			require.Contains(t, response.Body.String(), "RECORDING_SCHEMA_INCOMPATIBLE")
			var sentinel string
			require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT sentinel FROM `+table).Scan(&sentinel))
			require.Equal(t, "preserve-existing-data", sentinel)
			var receipts int
			query, args := `SELECT count(*) FROM sqlite_master WHERE name='gw_effect_receipts'`, []any{}
			if kind == "postgres" {
				query, args = `SELECT count(*) FROM information_schema.tables WHERE table_schema=$1 AND table_name='gw_effect_receipts'`, []any{f.namespace}
			}
			require.NoError(t, f.target.QueryRowContext(t.Context(), query, args...).Scan(&receipts))
			require.Zero(t, receipts)
		})
	}
}

func TestManagedGroupSchemaPostgresCreatePermissionReportsVerifiedFailure(t *testing.T) {
	f := newManagedGroupFixture(t, "postgres")
	role := fmt.Sprintf("gw_managed_nocreate_%d", time.Now().UnixNano())
	_, err := f.target.ExecContext(t.Context(), `CREATE ROLE "`+role+`" LOGIN PASSWORD 'gwtest'`)
	require.NoError(t, err)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := f.target.ExecContext(ctx, `DROP ROLE "`+role+`"`)
		require.NoError(t, err)
	})
	_, err = f.target.ExecContext(t.Context(), `GRANT USAGE ON SCHEMA "`+f.namespace+`" TO "`+role+`"`)
	require.NoError(t, err)
	// Explicit schema privilege is the only grant; this role cannot CREATE.
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := f.target.ExecContext(ctx, `REVOKE USAGE ON SCHEMA "`+f.namespace+`" FROM "`+role+`"`)
		require.NoError(t, err)
	})
	var raw string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT connection_config FROM database_connectors WHERE id='connector-A'`).Scan(&raw))
	var config map[string]string
	require.NoError(t, json.Unmarshal([]byte(raw), &config))
	config["user"] = role
	encoded, err := json.Marshal(config)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE database_connectors SET connection_config=? WHERE id='connector-A'`, string(encoded))
	require.NoError(t, err)
	token := f.preview(t)
	request := f.schemaRequest()
	request["token"], request["operation_id"] = token.Token, token.OperationID
	response := performJSONRequest(t, f.router, http.MethodPost, f.confirmPath(), request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data recordingplan.SchemaOperation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Equal(t, recordingplan.SchemaOperationFailed, body.Data.Status)
	require.Equal(t, recordingplan.SchemaReasonPermissionDenied, body.Data.Reason)
	require.Empty(t, body.Data.VerifiedDigest)
	assertNoManagedMutation(t, f, "postgres")
	ready, err := f.groups.Readiness(t.Context(), f.group["id"].(string))
	require.NoError(t, err)
	require.False(t, ready.Ready)
}
