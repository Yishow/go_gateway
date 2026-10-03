package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func TestManagedGroupSchemaPostgresReportsRealPartialVerification(t *testing.T) {
	f := newManagedGroupFixture(t, "postgres")
	trigger := fmt.Sprintf("gw_managed_partial_%d", time.Now().UnixNano())
	function := `"` + f.namespace + `".remove_fixture_receipt`
	// Fault injection is confined to this disposable namespace: the DDL
	// transaction commits its data table, but its receipt table is removed.
	_, err := f.target.ExecContext(t.Context(), `CREATE FUNCTION `+function+`() RETURNS event_trigger LANGUAGE plpgsql AS $body$
	DECLARE command record;
	BEGIN
	  FOR command IN SELECT * FROM pg_event_trigger_ddl_commands() LOOP
	    IF command.schema_name = '`+f.namespace+`' AND command.object_type = 'table'
	       AND command.object_identity LIKE '%.gw_effect_receipts' THEN
	      EXECUTE 'DROP TABLE "`+f.namespace+`".gw_effect_receipts';
	    END IF;
	  END LOOP;
	END $body$`)
	require.NoError(t, err)
	_, err = f.target.ExecContext(t.Context(), `CREATE EVENT TRIGGER "`+trigger+`" ON ddl_command_end WHEN TAG IN ('CREATE TABLE') EXECUTE FUNCTION `+function+`() `)
	require.NoError(t, err)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := f.target.ExecContext(ctx, `DROP EVENT TRIGGER "`+trigger+`"`)
		require.NoError(t, err)
	})
	token := f.preview(t)
	request := f.schemaRequest()
	request["token"], request["operation_id"] = token.Token, token.OperationID
	confirmed := performJSONRequest(t, f.router, http.MethodPost, f.confirmPath(), request)
	require.Equal(t, http.StatusOK, confirmed.Code, confirmed.Body.String())
	var result struct {
		Data recordingplan.SchemaOperation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(confirmed.Body.Bytes(), &result))
	require.Equal(t, recordingplan.SchemaOperationPartial, result.Data.Status)
	require.Equal(t, recordingplan.SchemaReasonVerificationGap, result.Data.Reason)
	require.Empty(t, result.Data.VerifiedDigest)
	var tables int
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT count(*) FROM information_schema.tables WHERE table_schema=$1`, f.namespace).Scan(&tables))
	require.Equal(t, 1, tables)
	ready, err := f.groups.Readiness(t.Context(), f.group["id"].(string))
	require.NoError(t, err)
	require.False(t, ready.Ready)
	duplicate := performJSONRequest(t, f.router, http.MethodPost, f.confirmPath(), request)
	require.Equal(t, http.StatusOK, duplicate.Code, duplicate.Body.String())
	require.Contains(t, duplicate.Body.String(), `"status":"partial"`)
}
