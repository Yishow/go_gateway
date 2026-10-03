package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManagedGroupCompletedConfirmationRejectsChangedSource(t *testing.T) {
	for _, scenario := range []string{"source", "connector"} {
		t.Run(scenario, func(t *testing.T) {
			f := newManagedGroupFixture(t, "sqlite")
			token := f.preview(t)
			request := f.schemaRequest()
			request["token"], request["operation_id"] = token.Token, token.OperationID
			response := performJSONRequest(t, f.router, http.MethodPost, f.confirmPath(), request)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			response = performJSONRequest(t, f.router, http.MethodPost, f.confirmPath(), request)
			require.Equal(t, http.StatusOK, response.Code, "unchanged lost-response replay remains idempotent: %s", response.Body.String())
			statement := `UPDATE points SET address='40100' WHERE id='point-A'`
			if scenario == "connector" {
				statement = `UPDATE database_connectors SET identity_revision='connector-2' WHERE id='connector-A'`
			}
			_, err := f.db.ExecContext(t.Context(), statement)
			require.NoError(t, err)
			response = performJSONRequest(t, f.router, http.MethodPost, f.confirmPath(), request)
			require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
			lookup := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-operations/"+token.OperationID, nil)
			require.Equal(t, http.StatusOK, lookup.Code, lookup.Body.String())
			require.Contains(t, lookup.Body.String(), `"status":"succeeded"`, "the prior operation remains truthful historical evidence")
		})
	}
}
