package workspace

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecordingLifecycleActionsRejectUnknownScopeAndState(t *testing.T) {
	for _, mode := range []string{"empty", "missing", "foreign", "unknown-status", "canceled", "no-transaction", "nil-service"} {
		t.Run(mode, func(t *testing.T) {
			db, service, group := lifecycleReadFixture(t, false)
			id := group.ID
			ctx := t.Context()
			var want error = ErrWriteGroupNotFound
			switch mode {
			case "empty":
				id = " "
			case "missing":
				id = "not-a-recording"
			case "foreign":
				_, err := db.Exec(`UPDATE write_groups SET workspace_id='another-workspace' WHERE id=?`, id)
				require.NoError(t, err)
			case "unknown-status":
				_, err := db.Exec(`UPDATE write_groups SET status='private-unknown-state' WHERE id=?`, id)
				require.NoError(t, err)
				want = ErrWriteGroupLifecycleBlocked
			case "canceled":
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
				want = context.Canceled
			case "no-transaction":
				service.workspaceSvc.repo = NewMemoryRepository()
				want = ErrWriteGroupServiceUnavailable
			case "nil-service":
				service = nil
				want = ErrWriteGroupServiceUnavailable
			}
			before := lifecycleReadState(t, db)
			view, err := service.Lifecycle(ctx, id)
			require.Nil(t, view)
			require.ErrorIs(t, err, want)
			require.Equal(t, before, lifecycleReadState(t, db))
		})
	}
}

type lifecycleReadFaultRepository struct {
	*SQLRepository
	beginErr, getErr, rollbackErr error
	commits, rollbacks            int
}

func (r *lifecycleReadFaultRepository) BeginDatabaseSetup(ctx context.Context) (DatabaseSetupTx, error) {
	if r.beginErr != nil {
		return nil, r.beginErr
	}
	tx, err := r.SQLRepository.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, err
	}
	return &lifecycleReadFaultTx{DatabaseSetupTx: tx, owner: r}, nil
}

type lifecycleReadFaultTx struct {
	DatabaseSetupTx
	owner *lifecycleReadFaultRepository
}

func (tx *lifecycleReadFaultTx) Get(ctx context.Context) (*Record, error) {
	if tx.owner.getErr != nil {
		return nil, tx.owner.getErr
	}
	return tx.DatabaseSetupTx.Get(ctx)
}

func (tx *lifecycleReadFaultTx) Commit() error {
	tx.owner.commits++
	return tx.DatabaseSetupTx.Commit()
}

func (tx *lifecycleReadFaultTx) Rollback() error {
	tx.owner.rollbacks++
	err := tx.DatabaseSetupTx.Rollback()
	return errors.Join(err, tx.owner.rollbackErr)
}

func TestRecordingLifecycleActionsReadFailureNeverReturnsAvailability(t *testing.T) {
	for _, mode := range []string{"begin", "workspace-missing", "workspace-read", "rollback"} {
		t.Run(mode, func(t *testing.T) {
			db, service, group := lifecycleReadFixture(t, false)
			fault := errors.New("injected read failure")
			repo := &lifecycleReadFaultRepository{SQLRepository: NewSQLRepository(db)}
			want := fault
			switch mode {
			case "begin":
				repo.beginErr = fault
			case "workspace-missing":
				repo.getErr, want = ErrNotFound, ErrWriteGroupNotFound
			case "workspace-read":
				repo.getErr = fault
			case "rollback":
				repo.rollbackErr = fault
			}
			service.workspaceSvc.repo = repo
			before := lifecycleReadState(t, db)
			view, err := service.Lifecycle(t.Context(), group.ID)
			require.Nil(t, view)
			require.ErrorIs(t, err, want)
			require.Zero(t, repo.commits)
			if mode != "begin" {
				require.Equal(t, 1, repo.rollbacks)
			}
			require.Equal(t, before, lifecycleReadState(t, db))
		})
	}
}

func TestRecordingLifecycleActionsAlwaysRollBackReadTransaction(t *testing.T) {
	db, service, group := lifecycleReadFixture(t, true)
	repo := &lifecycleReadFaultRepository{SQLRepository: NewSQLRepository(db)}
	service.workspaceSvc.repo = repo
	service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(ctx context.Context, tx *sql.Tx, actual *WriteGroup) error {
		// The guard contract is read-only. Even this deliberately misbehaving
		// test seam cannot persist an accidental write through the read endpoint.
		_, err := tx.ExecContext(ctx, `UPDATE write_groups SET name='must-roll-back' WHERE id=?`, actual.ID)
		return err
	}))
	before := lifecycleReadState(t, db)
	for range 2 {
		view, err := service.Lifecycle(t.Context(), group.ID)
		require.NoError(t, err)
		require.True(t, view.Actions[2].Allowed)
	}
	require.Zero(t, repo.commits)
	require.Equal(t, 2, repo.rollbacks)
	require.Equal(t, before, lifecycleReadState(t, db))
}
