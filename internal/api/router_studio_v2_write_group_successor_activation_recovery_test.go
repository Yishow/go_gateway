package api

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

func TestWriteGroupSuccessorActivationReplaySurvivesSQLiteReopen(t *testing.T) {
	f := newSuccessorActivationAPIFixture(t)
	request := f.activationConfirmation(t)
	event := f.postActivation(t, request)
	before := successorAPIAllState(t, f.f.db)
	var seq int
	var name, path string
	require.NoError(t, f.f.db.QueryRow(`PRAGMA database_list`).Scan(&seq, &name, &path))
	require.NotEmpty(t, path)
	require.NoError(t, f.f.db.Close())
	reopened, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reopened.Close()) })
	reopened.SetMaxOpenConns(1)
	ws := workspace.NewService(workspace.NewSQLRepository(reopened))
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(reopened))
	groups := workspace.NewWriteGroupService(ws, workspace.NewSQLWriteGroupRepository(reopened)).WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{}).WithTableInspector(dbtarget.NewReadOnlyTableInspector(targets))
	groups.WithManagedTableInspector(groups.ManagedTableInspector(targets))
	router := NewRouter(&DatalinkServices{Workspace: ws, WriteGroups: groups, DBTarget: targets, RecordingPlan: recordingplan.NewService(recordingplan.NewSQLRepository(reopened))})
	response := performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.Equal(t, event["id"], decodeJSONBody(t, response)["data"].(map[string]any)["id"])
	require.Equal(t, before, successorAPIAllState(t, reopened), "exact intent replay after process state loss must not apply again")
}

func TestWriteGroupSuccessorActivationRollbackKeepsPreparedReservation(t *testing.T) {
	f := newSuccessorActivationAPIFixture(t)
	request := f.activationConfirmation(t)
	_, err := f.f.db.Exec(`CREATE TRIGGER reject_activation_head BEFORE UPDATE OF current_group_id ON write_group_lineages WHEN NEW.current_group_id<>OLD.current_group_id BEGIN SELECT RAISE(ABORT,'activation-head-failed'); END`)
	require.NoError(t, err)
	before := successorAPIAllState(t, f.f.db)
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
	require.NotContains(t, response.Body.String(), "activation-head-failed")
	require.Equal(t, before, successorAPIAllState(t, f.f.db), "version, runtime, fences, event and reservation must roll back together")
	_, err = f.f.db.Exec(`DROP TRIGGER reject_activation_head`)
	require.NoError(t, err)
	event := f.postActivation(t, request)
	require.NotEmpty(t, event["id"])
}

func TestWriteGroupSuccessorActivationReadinessCannotUseMissingTarget(t *testing.T) {
	f := newSuccessorActivationAPIFixture(t)
	review := f.postActivation(t, f.activationRequest(t))
	saved, err := f.f.groups.Get(t.Context(), f.prepared.Group.ID)
	require.NoError(t, err)
	target := f.target(t)
	_, err = target.Exec(`ALTER TABLE "` + saved.Group.Destination.TableName + `" RENAME TO "missing_activation_target"`)
	require.NoError(t, err)
	before := successorAPIAllState(t, f.f.db)
	targetBefore := successorAPIAllState(t, target)
	ready := f.postActivation(t, map[string]any{"action": "activation-readiness", "envelope": review["envelope"], "review_digest": review["review_digest"]})
	require.Equal(t, false, ready["readiness"].(map[string]any)["ready"])
	require.Empty(t, ready["readiness_digest"], "unready target cannot yield activation grant")
	require.Equal(t, before, successorAPIAllState(t, f.f.db))
	require.Equal(t, targetBefore, successorAPIAllState(t, target), "readiness must not repair a missing target")
}

func TestWriteGroupSuccessorActivationConfirmationReprobesTarget(t *testing.T) {
	f := newSuccessorActivationAPIFixture(t)
	request := f.activationConfirmation(t)
	saved, err := f.f.groups.Get(t.Context(), f.prepared.Group.ID)
	require.NoError(t, err)
	target := f.target(t)
	_, err = target.Exec(`ALTER TABLE "` + saved.Group.Destination.TableName + `" RENAME TO "changed_after_readiness"`)
	require.NoError(t, err)
	before := successorAPIAllState(t, f.f.db)
	targetBefore := successorAPIAllState(t, target)
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
	require.Less(t, response.Code, 500, response.Body.String())
	require.Equal(t, before, successorAPIAllState(t, f.f.db), "an earlier readiness digest cannot grant stale activation")
	require.Equal(t, targetBefore, successorAPIAllState(t, target))
}
