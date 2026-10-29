package workspace

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type readSnapshotRepository struct {
	record *Record
	err    error
	saves  int
}

func (r *readSnapshotRepository) Get(context.Context) (*Record, error) { return r.record, r.err }
func (r *readSnapshotRepository) Save(context.Context, *Record) error  { r.saves++; return nil }

func TestReadSnapshotAbsentAndFailedReadNeverSave(t *testing.T) {
	for _, cause := range []error{ErrNotFound, errors.New("read failure")} {
		repo := &readSnapshotRepository{err: cause}
		result, err := NewService(repo).ReadSnapshot(t.Context())
		require.ErrorIs(t, err, cause)
		require.Nil(t, result)
		require.Zero(t, repo.saves)
	}
}

func TestReadSnapshotDetachedAndConcurrentDeletionDoesNotRecreate(t *testing.T) {
	repo := &readSnapshotRepository{record: &Record{ID: "W", DatabaseSetupRevision: "R", OrderedDeviceIDs: []string{"D"}}}
	service := NewService(repo)
	result, err := service.ReadSnapshot(t.Context())
	require.NoError(t, err)
	result.ID, result.OrderedDeviceIDs[0] = "changed", "changed"
	require.Equal(t, "W", repo.record.ID)
	require.Equal(t, []string{"D"}, repo.record.OrderedDeviceIDs)
	// Deletion between two read observations must remain absent, not bootstrap.
	repo.record, repo.err = nil, ErrNotFound
	result, err = service.ReadSnapshot(t.Context())
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, result)
	require.Zero(t, repo.saves)
}

func TestReadSnapshotCanceledOrUnconfiguredDoesNotRead(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	repo := &readSnapshotRepository{}
	_, err := NewService(repo).ReadSnapshot(ctx)
	require.ErrorIs(t, err, context.Canceled)
	require.Zero(t, repo.saves)
	_, err = (*Service)(nil).ReadSnapshot(t.Context())
	require.ErrorIs(t, err, ErrDatabaseSetupUnavailable)
}
