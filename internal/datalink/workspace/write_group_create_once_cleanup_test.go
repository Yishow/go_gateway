package workspace

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/sqltx"

	"github.com/stretchr/testify/require"
)

func TestCreateOnceDoesNotRetryUnknownCommitOrLostConnection(t *testing.T) {
	for _, sentinel := range []error{sqltx.ErrCommitOutcomeUnknown, sqltx.ErrConnectionUnusable} {
		t.Run(sentinel.Error(), func(t *testing.T) {
			attempts := 0
			failed := errors.Join(sentinel, errors.New("database is locked"))
			_, err := retryWhileSQLiteBusy(t.Context(), func() (*WriteGroupSaveResult, error) {
				attempts++
				return nil, failed
			})
			require.ErrorIs(t, err, sentinel)
			require.Equal(t, 1, attempts)
		})
	}
}

type unknownCommitRepository struct{ *SQLRepository }
type unknownCommitTx struct{ DatabaseSetupTx }

func (r unknownCommitRepository) BeginDatabaseSetup(ctx context.Context) (DatabaseSetupTx, error) {
	tx, err := r.SQLRepository.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, err
	}
	return unknownCommitTx{tx}, nil
}
func (tx unknownCommitTx) Commit() error {
	if err := tx.DatabaseSetupTx.Commit(); err != nil {
		return err
	}
	return errors.Join(sqltx.ErrCommitOutcomeUnknown, errors.New("database is locked: acknowledgment lost"))
}

func TestEnsureBasicManagedDoesNotConvertUnknownCommitToReplaySuccess(t *testing.T) {
	db := openWorkspaceTestDB(t, filepath.Join(t.TempDir(), "unknown.db"))
	defer db.Close()
	fixture := seedWriteGroupFixture(t.Context(), t, db)
	repo := NewSQLRepository(db)
	workspaceSvc := NewService(unknownCommitRepository{repo})
	current, err := workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	result, err := service.EnsureBasicManaged(t.Context(), fixture.deviceID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicManagedDraft(fixture),
	})
	require.ErrorIs(t, err, sqltx.ErrCommitOutcomeUnknown)
	require.Nil(t, result)
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM write_group_basic_keys").Scan(&count))
	require.Equal(t, 1, count, "lost acknowledgment must not undo or repeat committed creation")
}
