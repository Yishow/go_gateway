package workspace

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
)

func TestSuccessorFirstStartAtomicEveryWriteRollsBack(t *testing.T) {
	for _, point := range []struct{ name, event, table, when string }{
		{"immutable_version", "INSERT", "write_group_versions", ""},
		{"runtime_version", "INSERT", "wg_runtime_versions", ""},
		{"applied_revision", "UPDATE", "write_groups", "WHEN NEW.applied_revision<>OLD.applied_revision"},
		{"head", "UPDATE", "write_group_lineages", "WHEN NEW.current_group_id<>OLD.current_group_id"},
		{"clear_reservation", "UPDATE", "write_group_lineages", "WHEN OLD.reserved_group_id IS NOT NULL AND NEW.reserved_group_id IS NULL"},
		{"workspace_projection", "UPDATE", "system_settings", "WHEN NEW.key='studio_v2_workspace'"},
		{"first_start_checkpoint", "INSERT", "managed_schema_operations", "WHEN NEW.action='successor_first_start'"},
	} {
		t.Run(point.name, func(t *testing.T) {
			h, id, req := firstStartAtomicFixture(t)
			_, err := h.db.Exec("CREATE TRIGGER reject_first_start AFTER " + point.event + " ON " + point.table + " " + point.when + " BEGIN SELECT RAISE(ABORT,'first-start-atomic-fault'); END")
			require.NoError(t, err)
			before := firstStartAtomicSnapshot(t, h)
			result, err := h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
			require.Error(t, err)
			require.Nil(t, result)
			require.ErrorContains(t, err, "first-start-atomic-fault", "the intended writepoint, not an earlier guard, must be reached")
			require.Equal(t, before, firstStartAtomicSnapshot(t, h))
			requireFirstStartNoPredecessorHistory(t, h, id)
			require.Zero(t, h.countRows("SELECT COUNT(*) FROM managed_schema_operations WHERE action='successor_first_start'"))
			_, err = h.db.Exec("DROP TRIGGER reject_first_start")
			require.NoError(t, err)
			result, err = h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
			require.NoError(t, err)
			require.NotNil(t, result)
		})
	}
}

func TestSuccessorFirstStartAtomicWorkspaceRollbackAndLostACK(t *testing.T) {
	for _, stage := range []string{"workspace_save", "before_commit", "lost_ack"} {
		t.Run(stage, func(t *testing.T) {
			h, id, req := firstStartAtomicFixture(t)
			before := firstStartAtomicSnapshot(t, h)
			fault := &recreateCommitRepository{SQLRepository: NewSQLRepository(h.db), stage: stage}
			service := activationTestService(h.service)
			service.workspaceSvc = NewService(fault)
			result, err := service.ConfirmSuccessorFirstStart(t.Context(), id, req)
			if stage == "lost_ack" {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, 1, fault.commits)
				replay, err := h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
				require.NoError(t, err)
				require.Equal(t, result, replay)
				require.Equal(t, 1, h.countRows("SELECT COUNT(*) FROM managed_schema_operations WHERE action='successor_first_start'"))
			} else {
				require.ErrorIs(t, err, errRecreateInjected)
				require.Nil(t, result)
				require.Equal(t, before, firstStartAtomicSnapshot(t, h))
			}
			requireFirstStartNoPredecessorHistory(t, h, id)
		})
	}
}

func TestSuccessorFirstStartAtomicWriterBarrierAndCompetingCASRollback(t *testing.T) {
	for _, mode := range []string{"barrier", "lineage_cas"} {
		t.Run(mode, func(t *testing.T) {
			h, id, req := firstStartAtomicFixture(t)
			if mode == "barrier" {
				h.service.WithWriterOwnershipActivationBarrier(writerOwnershipActivationFunc(func(_ context.Context, tx *sql.Tx, request WriterOwnershipActivationRequest) error {
					require.Equal(t, req.Envelope.Request.Stage.Current.GroupID, request.GroupID)
					var count int
					require.NoError(t, tx.QueryRow(`SELECT COUNT(*) FROM wg_runtime_versions WHERE group_id=?`, request.GroupID).Scan(&count))
					require.Equal(t, 1, count)
					return errRecreateInjected
				}))
			} else {
				_, err := h.db.Exec(`CREATE TRIGGER compete_first_start AFTER INSERT ON wg_runtime_versions BEGIN UPDATE write_group_lineages SET revision=revision+1; END`)
				require.NoError(t, err)
			}
			before := firstStartAtomicSnapshot(t, h)
			result, err := h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
			require.Nil(t, result)
			if mode == "barrier" {
				require.ErrorIs(t, err, errRecreateInjected)
			} else {
				require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
			}
			require.Equal(t, before, firstStartAtomicSnapshot(t, h))
		})
	}
}

func TestSuccessorFirstStartAtomicExpiredUncommittedBoundaryNeedsFreshReview(t *testing.T) {
	h, id, req := firstStartAtomicFixture(t)
	_, err := h.db.Exec(`CREATE TRIGGER reject_first_start BEFORE INSERT ON managed_schema_operations WHEN NEW.action='successor_first_start' BEGIN SELECT RAISE(ABORT,'checkpoint-fault'); END`)
	require.NoError(t, err)
	before := firstStartAtomicSnapshot(t, h)
	_, err = h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
	require.ErrorContains(t, err, "checkpoint-fault")
	require.Equal(t, before, firstStartAtomicSnapshot(t, h))
	_, err = h.db.Exec("DROP TRIGGER reject_first_start")
	require.NoError(t, err)
	h.service.now = func() time.Time { return req.Envelope.Request.Boundary }
	_, err = h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
	require.Error(t, err)
	require.Equal(t, before, firstStartAtomicSnapshot(t, h))
	fresh := req.Envelope.Request
	fresh.Boundary = fresh.Boundary.Add(time.Hour)
	confirmation := firstStartAtomicConfirmation(t, h, id, fresh)
	require.NotEqual(t, req.ReviewDigest, confirmation.ReviewDigest)
	result, err := h.service.ConfirmSuccessorFirstStart(t.Context(), id, confirmation)
	require.NoError(t, err)
	require.Equal(t, fresh.Boundary, result.Boundary)
	after := firstStartAtomicSnapshot(t, h)
	_, err = h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	require.Equal(t, after, firstStartAtomicSnapshot(t, h))
}

func TestSuccessorFirstStartAtomicUnknownOriginalOutcomeCannotBeReplaced(t *testing.T) {
	h, id, req := firstStartAtomicFixture(t)
	scope := req.Envelope.Request.Stage.Envelope.Request
	root := recordingplan.SuccessorReserveOperationID(scope.WorkspaceID, scope.OperationID)
	operation := recordingplan.SuccessorFirstStartOperationID(scope.WorkspaceID, root)
	_, err := h.db.Exec(`INSERT INTO managed_schema_operations(operation_id,token,workspace_id,scope_key,owner,action,status,created_at,updated_at,detail,payload_digest) VALUES(?,?,?,'unknown-original-scope','owner','successor_first_start','unknown',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'{}',?)`, operation, "successor_first_start:"+operation, scope.WorkspaceID, req.ReviewDigest)
	require.NoError(t, err)
	before := firstStartAtomicSnapshot(t, h)
	result, err := h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
	require.Error(t, err)
	require.Nil(t, result)
	fresh := req.Envelope.Request
	fresh.Boundary = fresh.Boundary.Add(time.Hour)
	_, err = h.service.SuccessorActivationReview(t.Context(), id, fresh)
	require.Error(t, err, "new review cannot hide an unresolved original operation")
	require.Equal(t, before, firstStartAtomicSnapshot(t, h))
	require.Zero(t, h.countRows("SELECT COUNT(*) FROM write_group_versions"))
}
