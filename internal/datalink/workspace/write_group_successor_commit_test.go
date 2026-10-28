package workspace

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/testutil/lostcommit"
)

var errRecreateInjected = errors.New("private SQL password payload injected")

type recreateCommitRepository struct {
	*SQLRepository
	stage   string
	commits int
}

func (r *recreateCommitRepository) BeginDatabaseSetup(ctx context.Context) (DatabaseSetupTx, error) {
	tx, err := r.SQLRepository.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, err
	}
	return &recreateCommitTx{DatabaseSetupTx: tx, repo: r}, nil
}

type recreateCommitTx struct {
	DatabaseSetupTx
	repo *recreateCommitRepository
}

func (tx *recreateCommitTx) Save(ctx context.Context, record *Record) error {
	if err := tx.DatabaseSetupTx.Save(ctx, record); err != nil {
		return err
	}
	if tx.repo.stage == "workspace_save" {
		return errRecreateInjected
	}
	return nil
}
func (tx *recreateCommitTx) Commit() error {
	tx.repo.commits++
	if tx.repo.stage == "before_commit" {
		return errRecreateInjected
	}
	if err := tx.DatabaseSetupTx.Commit(); err != nil {
		return err
	}
	if tx.repo.stage == "lost_ack" {
		return lostcommit.ErrResponseLost
	}
	return nil
}

func TestRecreateCommitFailureRollbackAndLostACK(t *testing.T) {
	for _, stage := range []string{"workspace_save", "before_commit", "lost_ack"} {
		t.Run(stage, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			before := lineageReadSnapshot(t, h)
			fault := &recreateCommitRepository{SQLRepository: NewSQLRepository(h.db), stage: stage}
			service := NewWriteGroupService(NewService(fault), h.service.repo)
			result, err := service.reserveRecreatedDraft(t.Context(), req)
			if stage == "lost_ack" {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.EqualValues(t, 2, result.Generation)
				require.Equal(t, 1, fault.commits)
				durable, err := h.service.reserveRecreatedDraft(t.Context(), req)
				require.NoError(t, err)
				require.Equal(t, result, durable)
			} else {
				require.ErrorIs(t, err, errRecreateInjected)
				require.Equal(t, "write group successor reservation unavailable", err.Error())
				require.Nil(t, result)
				require.Equal(t, before, lineageReadSnapshot(t, h))
			}
		})
	}
}

func TestRecreateBacklogGuardPreservesAcceptedDataAndHistory(t *testing.T) {
	for _, mode := range []string{"approved", "missing", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			predecessor := cloneWriteGroup(req.Group)
			predecessor.AppliedRevision = predecessor.Revision
			predecessor.Status = WriteGroupStatusReady
			tx, err := h.db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			effective := time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)
			require.NoError(t, h.service.repo.insertVersionInTx(t.Context(), tx, predecessor, effective))
			require.NoError(t, h.service.saveRuntimeVersionInTx(t.Context(), tx, &WriteGroupAppliedSnapshot{WorkspaceID: predecessor.WorkspaceID, GroupID: predecessor.ID, GroupRevision: predecessor.Revision, AppliedRevision: predecessor.AppliedRevision, EffectiveAt: effective, Group: predecessor, RuntimeLayout: &WriteGroupRuntimeLayout{Dialect: dbtarget.SQLDialect("sqlite"), ReceiptTableReady: true}}))
			_, err = tx.Exec(`UPDATE write_groups SET applied_revision=? WHERE id=?`, predecessor.AppliedRevision, predecessor.ID)
			require.NoError(t, err)
			require.NoError(t, tx.Commit())
			seedLifecycleAcceptedData(t, h.db, predecessor)
			calls := 0
			if mode != "missing" {
				h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(ctx context.Context, tx *sql.Tx, group *WriteGroup) error {
					calls++
					require.Equal(t, predecessor.ID, group.ID)
					var samples int
					require.NoError(t, tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_samples`).Scan(&samples))
					require.Equal(t, 1, samples)
					if mode == "rejected" {
						return errRecreateInjected
					}
					return nil
				}))
			}
			before := lifecyclePreservedData(t, h.db)
			config := lineageReadSnapshot(t, h)
			result, err := h.service.reserveRecreatedDraft(t.Context(), req)
			require.Equal(t, before, lifecyclePreservedData(t, h.db))
			if mode == "approved" {
				require.NoError(t, err)
				require.Empty(t, result.Group.AppliedRevision)
				require.Equal(t, 0, h.countRows(`SELECT COUNT(*) FROM write_group_versions WHERE group_id=?`, result.Group.ID))
				require.Equal(t, 0, h.countRows(`SELECT COUNT(*) FROM wg_runtime_versions WHERE group_id=?`, result.Group.ID))
				require.Equal(t, 0, h.countRows(`SELECT COUNT(*) FROM write_group_ownership_fences WHERE group_id=?`, result.Group.ID))
				intake, err := h.service.IntakeEligibility(t.Context(), result.Group.ID)
				require.NoError(t, err)
				require.False(t, intake.Accepting)
				historical, err := h.service.ResolveAppliedAt(t.Context(), predecessor.ID, effective.Add(time.Hour))
				require.NoError(t, err)
				require.Equal(t, predecessor.AppliedRevision, historical.AppliedRevision)
			} else {
				require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
				require.Equal(t, config, lineageReadSnapshot(t, h))
			}
			if mode == "missing" {
				require.Zero(t, calls)
			} else {
				require.Equal(t, 1, calls)
			}
		})
	}
}
