package workspace

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRecreateConcurrentSeparateConnections(t *testing.T) {
	for _, mode := range []string{"same", "different_operation", "different_intent"} {
		t.Run(mode, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "race.db")
			h, req := recreateServiceFixture(t, path)
			peer := newCircuitProfileHarnessOn(t, openWorkspaceTestDB(t, path), h)
			peer.db.SetMaxOpenConns(1)
			for _, db := range []*sql.DB{h.db, peer.db} {
				_, err := db.Exec(`PRAGMA busy_timeout=1000`)
				require.NoError(t, err)
			}
			requests := []writeGroupRecreateReservationRequest{req, req}
			requests[1].Group = cloneWriteGroup(req.Group)
			if mode == "different_operation" {
				requests[1].OperationID = "O3"
			}
			if mode == "different_intent" {
				requests[1].Group.Name = "different"
			}
			services := []*WriteGroupService{h.service, peer.service}
			var revisionAdvances atomic.Int32
			for _, service := range services {
				service.workspaceSvc.newID = func() string { return fmt.Sprintf("workspace-race-%d", revisionAdvances.Add(1)) }
			}
			results := make([]*writeGroupReservedDraft, 2)
			errs := make([]error, 2)
			gate := make(chan struct{})
			var wg sync.WaitGroup
			for i := range services {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					<-gate
					results[i], errs[i] = services[i].reserveRecreatedDraft(t.Context(), requests[i])
				}(i)
			}
			close(gate)
			wg.Wait()
			successes := 0
			for i, err := range errs {
				if err == nil {
					successes++
					require.EqualValues(t, 2, results[i].Generation)
				} else {
					require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
				}
			}
			if mode == "same" {
				require.Equal(t, 2, successes)
				require.Equal(t, results[0], results[1])
			} else {
				require.Equal(t, 1, successes)
			}
			require.EqualValues(t, 1, revisionAdvances.Load(), "only the winner may advance workspace revision")
			require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_groups`))
			require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_group_generations`))
			var revision int64
			require.NoError(t, h.db.QueryRow(`SELECT revision FROM write_group_lineages`).Scan(&revision))
			require.EqualValues(t, 2, revision)
			winner := 0
			if errs[0] != nil {
				winner = 1
			}
			current, err := h.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)
			require.Equal(t, results[winner].WorkspaceRevision, current.DatabaseSetupRevision)
		})
	}
}
