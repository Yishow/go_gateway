package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func TestSuccessorFirstStartHTTPRechecksNoHistoryAfterReadiness(t *testing.T) {
	for _, owner := range []string{"predecessor", "reserved"} {
		t.Run(owner, func(t *testing.T) {
			r := newFirstStartAPIRig(t, "sqlite", false, true, false, workspace.WriteGroupStorageModePlain)
			request := r.confirmation(t)
			id := r.groupID
			if owner == "reserved" {
				id = r.prepared.Group.ID
			}
			// An undecodable runtime row is contradictory evidence, not proof of zero
			// history. Both sides of the same lineage must fail before target work.
			_, err := r.f.db.Exec(`INSERT INTO wg_runtime_versions(workspace_id,group_id,group_revision,payload) VALUES(?,?,?,?)`, r.f.record.ID, id, "late-history", "{invalid")
			require.NoError(t, err)
			before := successorAPIAllState(t, r.f.db)
			targetBefore := successorAPIAllState(t, r.target)
			response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
			require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
			require.Equal(t, before, successorAPIAllState(t, r.f.db))
			require.Equal(t, targetBefore, successorAPIAllState(t, r.target))
			require.NotContains(t, response.Body.String(), "{invalid")
		})
	}
}

func TestSuccessorFirstStartHTTPRechecksCASAndTarget(t *testing.T) {
	for _, field := range []string{"connector", "schema", "lineage", "source", "target", "missing_table"} {
		t.Run(field, func(t *testing.T) {
			r := newFirstStartAPIRig(t, "sqlite", false, true, false, workspace.WriteGroupStorageModePlain)
			request := r.confirmation(t)
			var err error
			switch field {
			case "connector":
				_, err = r.f.db.Exec(`UPDATE database_connectors SET identity_revision='changed' WHERE id='connector-A'`)
			case "schema":
				_, err = r.f.db.Exec(`UPDATE write_groups SET destination_schema_digest='changed' WHERE id=?`, r.prepared.Group.ID)
			case "lineage":
				_, err = r.f.db.Exec(`UPDATE write_group_lineages SET revision=revision+1 WHERE lineage_id=?`, r.request.LineageID)
			case "source":
				_, err = r.f.db.Exec(`UPDATE write_groups SET revision='changed' WHERE id=?`, r.groupID)
			case "target":
				_, err = r.f.db.Exec(`UPDATE write_groups SET revision='changed' WHERE id=?`, r.prepared.Group.ID)
			case "missing_table":
				_, err = r.target.Exec(`ALTER TABLE "` + r.prepared.Group.Destination.TableName + `" RENAME TO "missing_first_start_table"`)
			}
			require.NoError(t, err)
			before := successorAPIAllState(t, r.f.db)
			targetBefore := successorAPIAllState(t, r.target)
			response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
			require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
			require.Less(t, response.Code, 500, response.Body.String())
			require.Equal(t, before, successorAPIAllState(t, r.f.db))
			require.Equal(t, targetBefore, successorAPIAllState(t, r.target))
		})
	}
}
