package workspace

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"sync"
	"testing"
)

func TestRecreateBasicRollbackAndLostACK(t *testing.T) {
	for _, stage := range []string{"workspace_save", "before_commit", "lost_ack"} {
		t.Run(stage, func(t *testing.T) {
			h, req := recreateBasicFixture(t, "sqlite")
			before := recreateRepositorySnapshot(t, h)
			fault := &recreateCommitRepository{SQLRepository: NewSQLRepository(h.db), stage: stage}
			svc := NewWriteGroupService(NewService(fault), h.service.repo)
			result, err := svc.reserveRecreatedDraft(t.Context(), req)
			if stage == "lost_ack" {
				require.NoError(t, err)
				require.EqualValues(t, 2, result.Generation)
				replay, err := h.service.reserveRecreatedDraft(t.Context(), req)
				require.NoError(t, err)
				require.Equal(t, result, replay)
				return
			}
			require.ErrorIs(t, err, errRecreateInjected)
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

func TestRecreateBasicConcurrentReservation(t *testing.T) {
	for _, different := range []bool{false, true} {
		t.Run(map[bool]string{false: "same", true: "different"}[different], func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "basic.db")
			h, req := recreateBasicFixtureAt(t, "sqlite", path)
			peer := newCircuitProfileHarnessOn(t, openWorkspaceTestDB(t, path), h)
			peer.db.SetMaxOpenConns(1)
			for _, db := range []*sql.DB{h.db, peer.db} {
				_, err := db.Exec(`PRAGMA busy_timeout=1000`)
				require.NoError(t, err)
			}
			requests := []writeGroupRecreateReservationRequest{req, req}
			if different {
				requests[1].OperationID = "O3"
			}
			services := []*WriteGroupService{h.service, peer.service}
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
			for _, err := range errs {
				if err == nil {
					successes++
				} else {
					require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
				}
			}
			if different {
				require.Equal(t, 1, successes)
			} else {
				require.Equal(t, 2, successes)
				require.Equal(t, results[0], results[1])
			}
			require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_groups`))
			require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_group_generations`))
			require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_basic_keys`))
		})
	}
}

func TestRecreateBasicPreservesAcceptedHistoryAndRequiresOwnership(t *testing.T) {
	for _, mode := range []string{"approved", "missing", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			h, req := recreateBasicFixture(t, "sqlite")
			source := seedRecreateEngineEvidence(t, h, req)
			seedLifecycleAcceptedData(t, h.db, source.Group)
			calls := 0
			if mode != "missing" {
				h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(ctx context.Context, tx *sql.Tx, g *WriteGroup) error {
					calls++
					require.Equal(t, req.PredecessorGroupID, g.ID)
					if mode == "rejected" {
						return errRecreateInjected
					}
					return nil
				}))
			}
			before := lifecyclePreservedData(t, h.db)
			state := recreateRepositorySnapshot(t, h)
			result, err := h.service.reserveRecreatedDraft(t.Context(), req)
			require.Equal(t, before, lifecyclePreservedData(t, h.db))
			if mode == "approved" {
				require.NoError(t, err)
				require.EqualValues(t, 2, result.Generation)
			} else {
				require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
				require.Equal(t, state, recreateRepositorySnapshot(t, h))
			}
			if mode == "missing" {
				require.Zero(t, calls)
			} else {
				require.Equal(t, 1, calls)
			}
		})
	}
}
