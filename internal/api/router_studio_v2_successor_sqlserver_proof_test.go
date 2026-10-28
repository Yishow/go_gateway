package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func TestSQLServerSuccessorActivationRejectsForeignOrMissingProof(t *testing.T) {
	sqlServerSuccessorFixtures(t) // Missing/unsafe real fixtures must mark the entrypoint NOT_RUN.
	for _, family := range []string{"sqlserver"} {
		t.Run(family, func(t *testing.T) {
			for _, mode := range []workspace.WriteGroupStorageMode{"plain", "native_monthly"} {
				t.Run(string(mode), func(t *testing.T) {
					for _, field := range []string{"missing_ledger", "failed_ledger", "foreign_scope", "foreign_token_schema", "foreign_connector", "stale_connector", "corrupted_lineage", "stale_source", "stale_target"} {
						t.Run(field, func(t *testing.T) {
							r := newSQLServerSuccessorRig(t, mode, false, true)
							f := r.activationAPI()
							confirmation := f.activationConfirmation(t)
							var query string
							switch field {
							case "missing_ledger":
								query = `DELETE FROM managed_schema_operations WHERE token IN (SELECT token FROM managed_schema_preview_tokens WHERE plan_id=?)`
							case "failed_ledger":
								query = `UPDATE managed_schema_operations SET status='failed' WHERE token IN (SELECT token FROM managed_schema_preview_tokens WHERE plan_id=?)`
							case "foreign_scope":
								query = `UPDATE managed_schema_operations SET scope_key='connector-A/foreign' WHERE token IN (SELECT token FROM managed_schema_preview_tokens WHERE plan_id=?)`
							case "foreign_token_schema":
								query = `UPDATE managed_schema_preview_tokens SET schema_name='foreign_schema' WHERE plan_id=?`
							case "foreign_connector":
								query = `UPDATE managed_schema_preview_tokens SET connector_id='connector-A' WHERE plan_id=?`
							case "stale_connector":
								query = `UPDATE database_connectors SET identity_revision='stale' WHERE id='connector-B' AND ?<>''`
							case "corrupted_lineage":
								query = `UPDATE write_group_lineages SET revision=revision+1 WHERE reserved_group_id=?`
							case "stale_source":
								query = `UPDATE write_groups SET revision='stale' WHERE id=(SELECT predecessor_group_id FROM write_group_successor_preparations WHERE group_id=?)`
							case "stale_target":
								query = `UPDATE write_groups SET revision='stale' WHERE id=?`
							}
							if field == "stale_source" {
								_, err := r.f.db.Exec(`UPDATE write_groups SET revision='stale' WHERE id=?`, r.groupID)
								require.NoError(t, err)
							} else {
								result, err := r.f.db.Exec(query, r.prepared.Group.ID)
								require.NoError(t, err)
								n, err := result.RowsAffected()
								require.NoError(t, err)
								require.Positive(t, n)
							}
							before := sqlServerSuccessorLocalState(t, r.f.db)
							external := r.targetState(t)
							response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", confirmation)
							require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
							if field == "corrupted_lineage" {
								require.Equal(t, http.StatusInternalServerError, response.Code)
							} else {
								require.Less(t, response.Code, 500, response.Body.String())
							}
							require.Equal(t, before, sqlServerSuccessorLocalState(t, r.f.db))
							require.Equal(t, external, r.targetState(t))
						})
					}
				})
			}
		})
	}
}

func TestSQLServerSuccessorSourceTokenCannotConfirmTarget(t *testing.T) {
	sqlServerSuccessorFixtures(t) // Missing/unsafe real fixtures must mark the entrypoint NOT_RUN.
	for _, family := range []string{"sqlserver"} {
		t.Run(family, func(t *testing.T) {
			for _, mode := range []workspace.WriteGroupStorageMode{"plain", "native_monthly"} {
				t.Run(string(mode), func(t *testing.T) {
					r := newSQLServerSuccessorRig(t, mode, false, false)
					sourceToken := previewLong(t, r.f, r.groupID)
					review := r.review(t)
					prepared := r.reserve(t, review)
					request := successorAPIStage(t, r.successorAPIFixture, review, prepared)
					request["action"], request["token"], request["digest"], request["schema_operation_id"] = "schema-confirm", sourceToken.Token, sourceToken.Digest, sourceToken.OperationID
					before, external := sqlServerSuccessorLocalState(t, r.f.db), r.targetState(t)
					response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
					require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
					require.Less(t, response.Code, 500, response.Body.String())
					require.Equal(t, before, sqlServerSuccessorLocalState(t, r.f.db))
					require.Equal(t, external, r.targetState(t))
				})
			}
		})
	}
}
