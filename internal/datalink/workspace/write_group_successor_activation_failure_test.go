package workspace

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSuccessorActivationEveryLocalWriteRollsBack(t *testing.T) {
	for _, table := range []string{"write_group_versions", "wg_runtime_versions", "write_group_revision_ends", "write_group_cutovers", "write_group_ownership_fences", "write_group_lineages"} {
		t.Run(table, func(t *testing.T) {
			h, id, req := atomicActivationFixture(t)
			verb := "INSERT"
			if table == "write_group_lineages" {
				verb = "UPDATE"
			}
			_, err := h.db.Exec(fmt.Sprintf("CREATE TRIGGER reject_activation BEFORE %s ON %s BEGIN SELECT RAISE(ABORT,'injected activation failure'); END", verb, table))
			require.NoError(t, err)
			before := recreateRepositorySnapshot(t, h)
			_, err = h.service.ConfirmSuccessorActivation(t.Context(), id, req)
			require.Error(t, err)
			require.Equal(t, before, recreateRepositorySnapshot(t, h), "all version/runtime/end/event/fence/head writes are atomic")
			require.Equal(t, 0, h.countRows("SELECT COUNT(*) FROM write_group_cutovers"))
			require.Equal(t, 0, h.countRows("SELECT COUNT(*) FROM write_group_revision_ends"))
		})
	}
}

func TestSuccessorActivationWorkspaceRollbackAndLostACK(t *testing.T) {
	for _, stage := range []string{"workspace_save", "before_commit", "lost_ack"} {
		t.Run(stage, func(t *testing.T) {
			h, id, req := atomicActivationFixture(t)
			before := recreateRepositorySnapshot(t, h)
			fault := &recreateCommitRepository{SQLRepository: NewSQLRepository(h.db), stage: stage}
			service := activationTestService(h.service)
			service.workspaceSvc = NewService(fault)
			event, err := service.ConfirmSuccessorActivation(t.Context(), id, req)
			if stage == "lost_ack" {
				require.NoError(t, err)
				require.NotNil(t, event)
				require.Equal(t, 1, fault.commits)
				replay, e := h.service.ConfirmSuccessorActivation(t.Context(), id, req)
				require.NoError(t, e)
				require.Equal(t, event, replay)
			} else {
				require.ErrorIs(t, err, errRecreateInjected)
				require.Equal(t, before, recreateRepositorySnapshot(t, h))
				require.Equal(t, 0, h.countRows("SELECT COUNT(*) FROM write_group_cutovers"))
			}
		})
	}
}

func TestSuccessorActivationExactReplayAfterSQLiteReopen(t *testing.T) {
	h, id, req := atomicActivationFixture(t)
	event, err := h.service.ConfirmSuccessorActivation(t.Context(), id, req)
	require.NoError(t, err)
	var seq int
	var name, path string
	require.NoError(t, h.db.QueryRow("PRAGMA database_list").Scan(&seq, &name, &path))
	require.NotEmpty(t, path)
	require.NoError(t, h.db.Close())
	reopened := openWorkspaceTestDB(t, path)
	t.Cleanup(func() { reopened.Close() })
	service := NewWriteGroupService(NewService(NewSQLRepository(reopened)), NewSQLWriteGroupRepository(reopened))
	// No probe/clock/backlog dependencies are installed after restart: exact
	// durable replay must not re-authorize or execute the original transition.
	replay, err := service.ConfirmSuccessorActivation(t.Context(), id, req)
	require.NoError(t, err)
	require.Equal(t, event, replay)
	var count int
	require.NoError(t, reopened.QueryRow("SELECT COUNT(*) FROM write_group_cutovers").Scan(&count))
	require.Equal(t, 1, count)
}

func TestSuccessorActivationRejectsForeignTransactionAuthority(t *testing.T) {
	h, id, req := atomicActivationFixture(t)
	ready, err := h.service.SuccessorActivationReadiness(t.Context(), id, req)
	require.NoError(t, err)
	before := recreateRepositorySnapshot(t, h)
	tx, err := h.db.BeginTx(context.Background(), nil)
	require.NoError(t, err)
	record, err := readWorkspaceScope(t.Context(), tx, h.service.repo.query, req.Envelope.Request.Stage.Envelope.Request.WorkspaceID)
	require.NoError(t, err)
	_, err = h.service.activateSuccessorInTx(t.Context(), tx, record, req, "untrusted-intent", ready.snapshot, ready.Readiness)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	require.NoError(t, tx.Rollback())
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
}

func TestSuccessorActivationFinalizedOverlapRefusesWholeTransaction(t *testing.T) {
	h, id, req := atomicActivationFixture(t)
	boundary := req.Envelope.Request.Boundary
	_, err := h.db.Exec(`INSERT INTO wg_delivery_checkpoints(group_id,group_revision,next_close,updated_at) VALUES(?,?,?,?)`, id, req.Envelope.SourceAppliedRevision, boundary.Add(time.Second).Format(time.RFC3339Nano), time.Now().UTC())
	require.NoError(t, err)
	before := recreateRepositorySnapshot(t, h)
	_, err = h.service.ConfirmSuccessorActivation(t.Context(), id, req)
	require.ErrorIs(t, err, ErrCutoverFinalizedOverlap)
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
	require.Equal(t, 0, h.countRows("SELECT COUNT(*) FROM write_group_cutovers"))
}

func TestSuccessorActivationWriterBarrierFailureRollsBack(t *testing.T) {
	h, id, req := atomicActivationFixture(t)
	before := recreateRepositorySnapshot(t, h)
	h.service.WithWriterOwnershipActivationBarrier(writerOwnershipActivationFunc(func(_ context.Context, tx *sql.Tx, _ WriterOwnershipActivationRequest) error {
		var count int
		require.NoError(t, tx.QueryRow(`SELECT COUNT(*) FROM wg_runtime_versions WHERE group_id=?`, req.Envelope.Request.Stage.Current.GroupID).Scan(&count))
		require.Equal(t, 1, count, "runtime version is staged before ownership barrier")
		return errRecreateInjected
	}))
	_, err := h.service.ConfirmSuccessorActivation(t.Context(), id, req)
	require.ErrorIs(t, err, errRecreateInjected)
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
	require.Equal(t, 0, h.countRows("SELECT COUNT(*) FROM write_group_cutovers"))
}
