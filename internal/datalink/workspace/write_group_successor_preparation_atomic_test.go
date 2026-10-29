package workspace

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestSuccessorPreparationAtomicRollbackAndLostACK(t *testing.T) {
	for _, stage := range []string{"workspace_save", "before_commit", "lost_ack"} {
		t.Run(stage, func(t *testing.T) {
			h, id, req := successorActiveWideFixture(t)
			review, err := h.service.SuccessorReview(t.Context(), id, req)
			require.NoError(t, err)
			c := WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest}
			before := lineageReadSnapshot(t, h)
			fault := &recreateCommitRepository{SQLRepository: NewSQLRepository(h.db), stage: stage}
			svc := NewWriteGroupService(NewService(fault), h.service.repo).WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
			saved, err := svc.ReserveSuccessor(t.Context(), id, c)
			if stage == "lost_ack" {
				require.NoError(t, err)
				require.EqualValues(t, 2, saved.Generation)
				replay, err := h.service.ReserveSuccessor(t.Context(), id, c)
				require.NoError(t, err)
				require.Equal(t, saved, replay)
			} else {
				require.ErrorIs(t, err, errRecreateInjected)
				require.Equal(t, before, lineageReadSnapshot(t, h))
			}
		})
	}
}

func TestSuccessorPreparationConcurrentReservation(t *testing.T) {
	for _, different := range []bool{false, true} {
		t.Run(map[bool]string{false: "same", true: "different"}[different], func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "preparation.db")
			h, r := recreateServiceFixture(t, path)
			peer := newCircuitProfileHarnessOn(t, openWorkspaceTestDB(t, path), h)
			peer.db.SetMaxOpenConns(1)
			for _, db := range []*sql.DB{h.db, peer.db} {
				_, err := db.Exec("PRAGMA busy_timeout=1000")
				require.NoError(t, err)
			}
			req := WriteGroupSuccessorReviewRequest{WorkspaceID: r.WorkspaceID, ExpectedWorkspaceRevision: r.ExpectedWorkspaceRevision, LineageID: r.LineageID, ExpectedLineageRevision: r.ExpectedLineageRevision, ExpectedPredecessorRevision: r.ExpectedPredecessorRevision, ExpectedPredecessorConnectorRevision: r.ExpectedPredecessorConnectorRevision, ExpectedTargetConnectorRevision: r.ExpectedTargetConnectorRevision, OperationID: r.OperationID, Boundary: time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC), Group: r.Group}
			review, err := h.service.SuccessorReview(t.Context(), r.PredecessorGroupID, req)
			require.NoError(t, err)
			requests := []WriteGroupSuccessorReservationRequest{{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest}, {Envelope: review.Envelope, ReviewDigest: review.ReviewDigest}}
			if different {
				req.OperationID = "different-O"
				other, err := h.service.SuccessorReview(t.Context(), r.PredecessorGroupID, req)
				require.NoError(t, err)
				requests[1] = WriteGroupSuccessorReservationRequest{Envelope: other.Envelope, ReviewDigest: other.ReviewDigest}
			}
			results := make([]*WriteGroupSuccessorPreparation, 2)
			errs := make([]error, 2)
			services := []*WriteGroupService{h.service, peer.service}
			gate := make(chan struct{})
			var wg sync.WaitGroup
			for i := range services {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					<-gate
					results[i], errs[i] = services[i].ReserveSuccessor(t.Context(), r.PredecessorGroupID, requests[i])
				}(i)
			}
			close(gate)
			wg.Wait()
			wins := 0
			for _, err := range errs {
				if err == nil {
					wins++
				} else {
					require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
				}
			}
			if different {
				require.Equal(t, 1, wins)
			} else {
				require.Equal(t, 2, wins)
				require.Equal(t, results[0], results[1])
			}
			require.Equal(t, 2, h.countRows("SELECT COUNT(*) FROM write_groups"))
			require.Equal(t, 2, h.countRows("SELECT COUNT(*) FROM write_group_generations"))
		})
	}
}

func TestSuccessorPreparationCurrentCASAndForeignStageHaveNoEffects(t *testing.T) {
	for _, mode := range []string{"workspace", "lineage", "predecessor", "source-connector", "group", "group-revision", "target-connector", "envelope", "review"} {
		t.Run(mode, func(t *testing.T) {
			h, id, req := successorActiveWideFixture(t)
			review, err := h.service.SuccessorReview(t.Context(), id, req)
			require.NoError(t, err)
			saved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
			require.NoError(t, err)
			stage := WriteGroupSuccessorStageRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{WorkspaceRevision: saved.WorkspaceRevision, LineageRevision: saved.LineageRevision, PredecessorRevision: req.ExpectedPredecessorRevision, PredecessorConnectorRevision: req.ExpectedPredecessorConnectorRevision, GroupID: saved.Group.ID, GroupRevision: saved.Group.Revision, TargetConnectorRevision: saved.Group.Destination.ConnectorRevision}}
			switch mode {
			case "workspace":
				stage.Current.WorkspaceRevision = "stale"
			case "lineage":
				stage.Current.LineageRevision++
			case "predecessor":
				stage.Current.PredecessorRevision = "stale"
			case "source-connector":
				stage.Current.PredecessorConnectorRevision = "stale"
			case "group":
				stage.Current.GroupID = "foreign"
			case "group-revision":
				stage.Current.GroupRevision = "stale"
			case "target-connector":
				stage.Current.TargetConnectorRevision = "stale"
			case "envelope":
				stage.Envelope.Request.OperationID = "foreign"
			case "review":
				stage.ReviewDigest = "foreign"
			}
			before := lineageReadSnapshot(t, h)
			tx, err := h.db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			record, err := readWorkspaceScope(t.Context(), tx, h.service.repo.query, req.WorkspaceID)
			require.NoError(t, err)
			_, err = h.service.successorReservationInTx(t.Context(), tx, record, id, stage)
			require.Error(t, err)
			require.NoError(t, tx.Rollback())
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}
