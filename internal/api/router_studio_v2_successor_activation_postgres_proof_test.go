package api

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
	"net/http"
	"testing"
)

func TestPostgresSuccessorActivationRejectsForeignOrMissingProof(t *testing.T) {
	for _, field := range []string{"missing_ledger", "failed_ledger", "foreign_scope", "foreign_token_schema", "foreign_connector", "stale_connector", "corrupted_lineage", "stale_source", "stale_target"} {
		t.Run(field, func(t *testing.T) {
			r := newPostgresSuccessorRuntimeRig(t, true, true, false, false)
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
				query = `UPDATE managed_schema_preview_tokens SET schema_name='public' WHERE plan_id=?`
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
			before := successorAPIAllState(t, r.f.db)
			external := r.targetState(t)
			response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", confirmation)
			require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
			if field == "corrupted_lineage" {
				require.Equal(t, http.StatusInternalServerError, response.Code)
			} else {
				require.Less(t, response.Code, 500, response.Body.String())
			}
			require.Equal(t, before, successorAPIAllState(t, r.f.db))
			require.Equal(t, external, r.targetState(t))
		})
	}
}

func TestPostgresSuccessorSameNameOtherSchemaCannotSatisfyReadiness(t *testing.T) {
	r := newPostgresSuccessorRuntimeRig(t, false, true, false, false)
	f := r.activationAPI()
	saved, err := r.f.groups.Get(t.Context(), r.prepared.Group.ID)
	require.NoError(t, err)
	table := saved.Group.Destination.TableName
	_, err = r.target.Exec(`CREATE TABLE "` + r.f.namespace + `"."` + table + `" (LIKE "` + r.targetNamespace + `"."` + table + `" INCLUDING ALL)`)
	require.NoError(t, err)
	request := f.activationRequest(t)
	review := f.postActivation(t, request)
	_, err = r.target.Exec(`DROP TABLE "` + r.targetNamespace + `"."` + table + `"`)
	require.NoError(t, err)
	before := successorAPIAllState(t, r.f.db)
	target := r.targetState(t)
	ready := f.postActivation(t, map[string]any{"action": "activation-readiness", "envelope": review["envelope"], "review_digest": review["review_digest"]})
	require.Equal(t, false, ready["readiness"].(map[string]any)["ready"])
	require.Empty(t, ready["readiness_digest"])
	require.Equal(t, before, successorAPIAllState(t, r.f.db))
	require.Equal(t, target, r.targetState(t), "same-named A table must not repair or satisfy B")
}

func TestPostgresSuccessorRejectsSourcePreviewTokenForTarget(t *testing.T) {
	r := newPostgresSuccessorRig(t, false, true, false, false, false)
	sourceToken := previewLong(t, r.f, r.groupID)
	review := r.review(t)
	prepared := r.reserve(t, review)
	request := successorAPIStage(t, r.successorAPIFixture, review, prepared)
	request["action"], request["token"], request["digest"], request["schema_operation_id"] = "schema-confirm", sourceToken.Token, sourceToken.Digest, sourceToken.OperationID
	before := successorAPIAllState(t, r.f.db)
	target := r.targetState(t)
	response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
	require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
	require.Less(t, response.Code, 500, response.Body.String())
	require.Equal(t, before, successorAPIAllState(t, r.f.db))
	require.Equal(t, target, r.targetState(t))
}

func TestPostgresSuccessorStaleCurrentLineageCAS(t *testing.T) {
	r := newPostgresSuccessorRuntimeRig(t, false, true, false, false)
	f := r.activationAPI()
	request := f.activationRequest(t)
	stage := request["stage"].(map[string]any)
	cas := stage["current"].(workspace.WriteGroupSuccessorCurrentCAS)
	cas.LineageRevision--
	stage["current"] = cas
	before := successorAPIAllState(t, r.f.db)
	target := r.targetState(t)
	response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Equal(t, before, successorAPIAllState(t, r.f.db))
	require.Equal(t, target, r.targetState(t))
}
