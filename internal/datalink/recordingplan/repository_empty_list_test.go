package recordingplan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepositoryPreallocationKeepsEmptyListsNil(t *testing.T) {
	repo := NewSQLRepository(recordingRepositoryDB(t))
	plans, err := repo.ListPlansByWorkspace(t.Context(), "empty-workspace")
	require.NoError(t, err)
	require.Nil(t, plans)
	operations, err := repo.listStaleStorageRepairs(t.Context(), time.Now().UTC())
	require.NoError(t, err)
	require.Nil(t, operations)
}
