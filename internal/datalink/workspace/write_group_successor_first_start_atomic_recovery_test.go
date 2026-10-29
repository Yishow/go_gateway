package workspace

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSuccessorFirstStartAtomicOfflineReplayAfterReopenAndHeadMoves(t *testing.T) {
	h, id, req := firstStartAtomicFixture(t)
	result, err := h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
	require.NoError(t, err)
	successor, err := h.service.repo.Get(t.Context(), req.Envelope.Request.Stage.Envelope.Request.WorkspaceID, result.GroupID)
	require.NoError(t, err)
	// A later local head is not the authority for the original immutable result.
	seedLineageGeneration(t, h, successor, result.LineageID, "later-G3", 3)
	_, err = h.db.Exec(`UPDATE write_group_lineages SET current_group_id='later-G3',revision=revision+1 WHERE lineage_id=?`, result.LineageID)
	require.NoError(t, err)
	before := firstStartAtomicSnapshot(t, h)
	var seq int
	var name, path string
	require.NoError(t, h.db.QueryRow("PRAGMA database_list").Scan(&seq, &name, &path))
	require.NotEmpty(t, path)
	require.NoError(t, h.db.Close())
	reopened := openWorkspaceTestDB(t, path)
	t.Cleanup(func() { reopened.Close() })
	h.db = reopened
	service := NewWriteGroupService(NewService(NewSQLRepository(reopened)), NewSQLWriteGroupRepository(reopened))
	// No target inspector, clock, or writer is installed. Replay is local only.
	replay, err := service.ConfirmSuccessorFirstStart(t.Context(), id, req)
	require.NoError(t, err)
	require.Equal(t, result, replay)
	require.Equal(t, before, firstStartAtomicSnapshot(t, h))
	requireFirstStartNoPredecessorHistory(t, h, id)
}

func TestSuccessorFirstStartAtomicConcurrentConfirmations(t *testing.T) {
	for _, same := range []bool{true, false} {
		name := "different_boundaries"
		if same {
			name = "exact_intent"
		}
		t.Run(name, func(t *testing.T) {
			h, id, req := firstStartAtomicFixture(t)
			requests := []WriteGroupSuccessorActivationConfirmation{req, req}
			if !same {
				other := req.Envelope.Request
				other.Boundary = other.Boundary.Add(time.Second)
				requests[1] = firstStartAtomicConfirmation(t, h, id, other)
			}
			var seq int
			var name, path string
			require.NoError(t, h.db.QueryRow("PRAGMA database_list").Scan(&seq, &name, &path))
			peerDB := openWorkspaceTestDB(t, path)
			t.Cleanup(func() { peerDB.Close() })
			peerDB.SetMaxOpenConns(1)
			for _, db := range []*sql.DB{h.db, peerDB} {
				_, err := db.Exec("PRAGMA busy_timeout=1000")
				require.NoError(t, err)
			}
			peer := activationTestService(h.service)
			peer.repo = NewSQLWriteGroupRepository(peerDB)
			peer.workspaceSvc = NewService(NewSQLRepository(peerDB))
			services := []*WriteGroupService{h.service, peer}
			results := make([]*WriteGroupSuccessorFirstStartResult, 2)
			errs := make([]error, 2)
			gate := make(chan struct{})
			var wg sync.WaitGroup
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			for i := range services {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					<-gate
					results[i], errs[i] = services[i].ConfirmSuccessorFirstStart(ctx, id, requests[i])
				}(i)
			}
			close(gate)
			wg.Wait()
			if same {
				require.NoError(t, errs[0])
				require.NoError(t, errs[1])
				require.Equal(t, results[0], results[1])
			} else {
				successes := 0
				for i := range results {
					if errs[i] == nil {
						successes++
						require.Equal(t, requests[i].Envelope.Request.Boundary, results[i].Boundary)
					} else {
						require.Nil(t, results[i])
					}
				}
				require.Equal(t, 1, successes, "only one complete boundary intent may win: %v", errs)
			}
			require.Equal(t, 1, h.countRows("SELECT COUNT(*) FROM managed_schema_operations WHERE action='successor_first_start'"))
			require.Equal(t, 1, h.countRows("SELECT COUNT(*) FROM write_group_versions"))
			require.Equal(t, 1, h.countRows("SELECT COUNT(*) FROM wg_runtime_versions"))
			require.Equal(t, 2, h.countRows("SELECT COUNT(*) FROM write_group_generations"))
			requireFirstStartNoPredecessorHistory(t, h, id)
		})
	}
}
