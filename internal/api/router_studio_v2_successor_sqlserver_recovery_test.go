package api

import (
	"database/sql"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"
)

func TestSQLServerSuccessorActivationReplaySurvivesSQLiteReopen(t *testing.T) {
	sqlServerSuccessorFixtures(t) // Missing/unsafe real fixtures must mark the entrypoint NOT_RUN.
	for _, family := range []string{"sqlserver"} {
		t.Run(family, func(t *testing.T) {
			for _, mode := range []workspace.WriteGroupStorageMode{"plain", "native_monthly"} {
				t.Run(string(mode), func(t *testing.T) {
					r := newSQLServerSuccessorRig(t, mode, false, true)
					f := r.activationAPI()
					targetBefore := r.targetState(t)
					request := f.activationConfirmation(t)
					event := f.postActivation(t, request)
					before := sqlServerSuccessorLocalState(t, f.f.db)
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
					require.Equal(t, before, sqlServerSuccessorLocalState(t, reopened), "exact intent replay after process state loss must not apply again")
					require.Equal(t, targetBefore, r.targetState(t))
				})
			}
		})
	}
}

func TestSQLServerSuccessorActivationRollbackKeepsPreparedReservation(t *testing.T) {
	sqlServerSuccessorFixtures(t) // Missing/unsafe real fixtures must mark the entrypoint NOT_RUN.
	for _, family := range []string{"sqlserver"} {
		t.Run(family, func(t *testing.T) {
			for _, mode := range []workspace.WriteGroupStorageMode{"plain", "native_monthly"} {
				t.Run(string(mode), func(t *testing.T) {
					for _, fault := range []string{"head", "version", "fence"} {
						t.Run(fault, func(t *testing.T) {
							r := newSQLServerSuccessorRig(t, mode, false, true)
							f := r.activationAPI()
							targetBefore := r.targetState(t)
							request := f.activationConfirmation(t)
							statement := `CREATE TRIGGER reject_activation_head BEFORE UPDATE OF current_group_id ON write_group_lineages WHEN NEW.current_group_id<>OLD.current_group_id BEGIN SELECT RAISE(ABORT,'activation-head-failed'); END`
							if fault == "version" {
								statement = `CREATE TRIGGER reject_activation_head BEFORE INSERT ON write_group_versions BEGIN SELECT RAISE(ABORT,'activation-head-failed'); END`
							}
							if fault == "fence" {
								statement = `CREATE TRIGGER reject_activation_head BEFORE INSERT ON write_group_ownership_fences BEGIN SELECT RAISE(ABORT,'activation-head-failed'); END`
							}
							_, err := f.f.db.Exec(statement)
							require.NoError(t, err)
							before := sqlServerSuccessorLocalState(t, f.f.db)
							response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
							require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
							require.NotContains(t, response.Body.String(), "activation-head-failed")
							require.Equal(t, before, sqlServerSuccessorLocalState(t, f.f.db), "version, runtime, fences, event and reservation must roll back together")
							_, err = f.f.db.Exec(`DROP TRIGGER reject_activation_head`)
							require.NoError(t, err)
							event := f.postActivation(t, request)
							require.NotEmpty(t, event["id"])
							require.Equal(t, targetBefore, r.targetState(t))
						})
					}
				})
			}
		})
	}
}

func TestSQLServerSuccessorActivationLostCommitACK(t *testing.T) {
	sqlServerSuccessorFixtures(t) // Missing/unsafe real fixtures must mark the entrypoint NOT_RUN.
	for _, family := range []string{"sqlserver"} {
		t.Run(family, func(t *testing.T) {
			for _, mode := range []workspace.WriteGroupStorageMode{"plain", "native_monthly"} {
				t.Run(string(mode), func(t *testing.T) {
					r := newSQLServerSuccessorRig(t, mode, false, true)
					f := r.activationAPI()
					request := f.activationConfirmation(t)
					targetBefore := r.targetState(t)
					var seq int
					var name, path string
					require.NoError(t, r.f.db.QueryRow(`PRAGMA database_list`).Scan(&seq, &name, &path))
					lose := &atomic.Bool{}
					faults := &lostcommit.Faults{Lose: lose}
					var commits atomic.Int64
					afterCommit := func() { commits.Add(1) }
					faults.AfterCommit.Store(&afterCommit)
					faultDB := lostcommit.OpenFaults(path, faults)
					faultDB.SetMaxOpenConns(1)
					t.Cleanup(func() { require.NoError(t, faultDB.Close()) })
					ws := workspace.NewService(workspace.NewSQLRepository(faultDB))
					targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(faultDB))
					groups := workspace.NewWriteGroupService(ws, workspace.NewSQLWriteGroupRepository(faultDB)).WithClock(func() time.Time { return time.Unix(0, r.clock.Load()).UTC() }).WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{}).WithTableInspector(dbtarget.NewReadOnlyTableInspector(targets))
					groups.WithManagedTableInspector(groups.ManagedTableInspector(targets))
					router := NewRouter(&DatalinkServices{Workspace: ws, WriteGroups: groups, DBTarget: targets, RecordingPlan: recordingplan.NewService(recordingplan.NewSQLRepository(faultDB))})
					lose.Store(true)
					response := performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
					require.Equal(t, http.StatusOK, response.Code, response.Body.String())
					require.Equal(t, int64(1), commits.Load(), "real commit occurred once even though driver returned ACK loss")
					lose.Store(false)
					var id string
					require.NoError(t, r.f.db.QueryRow(`SELECT id FROM write_group_cutovers WHERE source_group_id=?`, r.groupID).Scan(&id))
					before := sqlServerSuccessorLocalState(t, r.f.db)
					require.Equal(t, id, f.postActivation(t, request)["id"])
					require.Equal(t, before, sqlServerSuccessorLocalState(t, r.f.db))
					require.Equal(t, targetBefore, r.targetState(t))
				})
			}
		})
	}
}
