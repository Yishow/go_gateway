package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

type firstStartLostResponse struct {
	header http.Header
	status int
	writes int
}

func (w *firstStartLostResponse) Header() http.Header       { return w.header }
func (w *firstStartLostResponse) WriteHeader(code int)      { w.status = code }
func (w *firstStartLostResponse) Write([]byte) (int, error) { w.writes++; return 0, io.ErrClosedPipe }

func TestSuccessorFirstStartHTTPLostResponseOfflineReopen(t *testing.T) {
	r := newFirstStartAPIRig(t, "sqlite", false, true, false, workspace.WriteGroupStorageModePlain)
	request := r.confirmation(t)
	raw, err := json.Marshal(request)
	require.NoError(t, err)
	lost := &firstStartLostResponse{header: make(http.Header)}
	req := httptest.NewRequest(http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	r.f.router.ServeHTTP(lost, req)
	require.Equal(t, http.StatusOK, lost.status)
	require.Positive(t, lost.writes)
	initial := successorStatus(t, r.successorAPIFixture)
	result := initial["first_start"].(map[string]any)["result"].(map[string]any)
	require.Equal(t, r.prepared.Group.ID, result["group_id"])
	require.NoError(t, r.pipe.Stop(time.Second))
	// A disconnected/revised target must not change historical status or replay.
	_, err = r.f.db.Exec(`UPDATE database_connectors SET enabled=0,identity_revision='offline-new-revision'`)
	require.NoError(t, err)
	before := successorAPIAllState(t, r.f.db)
	var seq int
	var name, path string
	require.NoError(t, r.f.db.QueryRow(`PRAGMA database_list`).Scan(&seq, &name, &path))
	require.NoError(t, r.f.db.Close())
	reopened, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reopened.Close()) })
	reopened.SetMaxOpenConns(1)
	ws := workspace.NewService(workspace.NewSQLRepository(reopened))
	groups := workspace.NewWriteGroupService(ws, workspace.NewSQLWriteGroupRepository(reopened)).WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{}).WithClock(func() time.Time { return r.boundary.Add(time.Hour) })
	forbidden := &successorAPIForbiddenInspector{}
	groups.WithTableInspector(forbidden)
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(reopened))
	groups.WithManagedTableInspector(groups.ManagedTableInspector(targets))
	router := NewRouter(&DatalinkServices{Workspace: ws, WriteGroups: groups, DBTarget: targets, RecordingPlan: recordingplan.NewService(recordingplan.NewSQLRepository(reopened))})
	for _, body := range []map[string]any{r.statusRequest(), request, {"action": "reserve", "envelope": r.preparationReview.Envelope, "review_digest": r.preparationReview.ReviewDigest}} {
		response := performJSONRequest(t, router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", body)
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		if body["action"] == "operation-status" {
			require.Equal(t, initial, decodeJSONBody(t, response)["data"])
		}
		if body["action"] == "first-start-confirm" {
			require.Equal(t, result, decodeJSONBody(t, response)["data"])
		}
		require.Equal(t, before, successorAPIAllState(t, reopened))
	}
	require.Zero(t, forbidden.calls, "readback/exact replay is local, never a new probe")
}

func TestSuccessorFirstStartHTTPRollbackExpiredBoundaryAndFreshReview(t *testing.T) {
	r := newFirstStartAPIRig(t, "sqlite", false, true, false, workspace.WriteGroupStorageModePlain)
	request := r.confirmation(t)
	_, err := r.f.db.Exec(`CREATE TRIGGER reject_first_start_head BEFORE UPDATE OF current_group_id ON write_group_lineages WHEN NEW.current_group_id<>OLD.current_group_id BEGIN SELECT RAISE(ABORT,'private-first-start-fault'); END`)
	require.NoError(t, err)
	before := successorAPIAllState(t, r.f.db)
	rejected := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
	require.GreaterOrEqual(t, rejected.Code, 400, rejected.Body.String())
	require.NotContains(t, rejected.Body.String(), "private-first-start-fault")
	require.Equal(t, before, successorAPIAllState(t, r.f.db))
	status := successorStatus(t, r.successorAPIFixture)
	require.Equal(t, "not_submitted", status["first_start"].(map[string]any)["status"])
	_, err = r.f.db.Exec(`DROP TRIGGER reject_first_start_head`)
	require.NoError(t, err)
	r.clock.Store(r.boundary.Add(time.Second).UnixNano())
	before = successorAPIAllState(t, r.f.db)
	rejected = performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
	require.GreaterOrEqual(t, rejected.Code, 400, rejected.Body.String())
	require.Less(t, rejected.Code, 500, rejected.Body.String())
	require.Equal(t, before, successorAPIAllState(t, r.f.db), "expired B must not silently become now")
	r.boundary = r.boundary.Add(time.Minute)
	result := r.start(t)
	require.Equal(t, r.boundary.Format(time.RFC3339Nano), result["boundary"])
}

func TestSuccessorFirstStartHTTPConcurrentIdenticalRequestOneResult(t *testing.T) {
	r := newFirstStartAPIRig(t, "sqlite", false, true, false, workspace.WriteGroupStorageModePlain)
	request := r.confirmation(t)
	raw, err := json.Marshal(request)
	require.NoError(t, err)
	responses := make(chan *httptest.ResponseRecorder, 2)
	start := make(chan struct{})
	var workers sync.WaitGroup
	for i := 0; i < 2; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			response := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", bytes.NewReader(raw))
			req.Header.Set("Content-Type", "application/json")
			r.f.router.ServeHTTP(response, req)
			responses <- response
		}()
	}
	close(start)
	workers.Wait()
	close(responses)
	var result any
	for response := range responses {
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		got := decodeJSONBody(t, response)["data"]
		if result == nil {
			result = got
		} else {
			require.Equal(t, result, got)
		}
	}
	var versions int
	require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM write_group_versions WHERE group_id=?`, r.prepared.Group.ID).Scan(&versions))
	require.Equal(t, 1, versions)
	for _, table := range []string{"write_group_cutovers", "write_group_revision_ends", "write_group_ownership_fences", "wg_delivery_sample_handoffs"} {
		var count int
		require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM `+table).Scan(&count))
		require.Zero(t, count, table)
	}
}
