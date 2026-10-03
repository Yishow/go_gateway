package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManagedGroupReadinessRejectsChangedPhysicalOwnership(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			for _, changed := range []string{"group-owner", "receipt-owner"} {
				t.Run(changed, func(t *testing.T) {
					f := newManagedGroupFixture(t, kind)
					token := f.preview(t)
					request := f.schemaRequest()
					request["token"], request["operation_id"] = token.Token, token.OperationID
					confirmed := performJSONRequest(t, f.router, http.MethodPost, f.confirmPath(), request)
					require.Equal(t, http.StatusOK, confirmed.Code, confirmed.Body.String())
					f.openTarget(t)
					ready, err := f.groups.Readiness(t.Context(), f.group["id"].(string))
					require.NoError(t, err)
					require.True(t, ready.Ready, "%+v", ready.Issues)
					table, column := token.GroupLayout.TableName, token.GroupLayout.OwnerColumn
					if changed == "receipt-owner" {
						table, column = "gw_effect_receipts", "_gw_effect_receipts_v1"
					}
					_, err = f.target.ExecContext(t.Context(), `ALTER TABLE "`+f.namespace+`"."`+table+`" DROP COLUMN "`+column+`"`)
					require.NoError(t, err)
					ready, err = f.groups.Readiness(t.Context(), f.group["id"].(string))
					require.NoError(t, err)
					require.False(t, ready.Ready, "stored proof cannot authorize a replaced ownership marker: %+v", ready)
					require.False(t, ready.SchemaReady)
				})
			}
		})
	}
}

func TestManagedGroupPostgresRejectsLossyNumericTable(t *testing.T) {
	f := newManagedGroupFixture(t, "postgres", "uint64")
	token := f.preview(t)
	request := f.schemaRequest()
	request["token"], request["operation_id"] = token.Token, token.OperationID
	confirmed := performJSONRequest(t, f.router, http.MethodPost, f.confirmPath(), request)
	require.Equal(t, http.StatusOK, confirmed.Code, confirmed.Body.String())
	member := f.group["members"].([]any)[0].(map[string]any)
	_, err := f.target.ExecContext(t.Context(), `ALTER TABLE "`+f.namespace+`"."`+token.GroupLayout.TableName+`" ALTER COLUMN "`+member["target_column"].(string)+`" TYPE NUMERIC(20,2)`)
	require.NoError(t, err)
	ready, err := f.groups.Readiness(t.Context(), f.group["id"].(string))
	require.NoError(t, err)
	require.False(t, ready.Ready, "NUMERIC(20,2) cannot preserve the full uint64 range")
	current, err := f.groups.Get(t.Context(), f.group["id"].(string))
	require.NoError(t, err)
	f.data["workspace_revision"] = current.WorkspaceRevision
	preview := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+f.group["id"].(string)+"/schema-preview", f.schemaRequest())
	require.Equal(t, http.StatusUnprocessableEntity, preview.Code, preview.Body.String())
	require.Contains(t, preview.Body.String(), "RECORDING_SCHEMA_INCOMPATIBLE")
}
