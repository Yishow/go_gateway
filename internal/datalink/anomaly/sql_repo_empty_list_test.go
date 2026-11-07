package anomaly

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLRepositoryEmptyListsRemainNil(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
	events, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a"})
	require.NoError(t, err)
	require.Zero(t, total)
	require.Nil(t, events)
	active, err := repo.ListActive(t.Context())
	require.NoError(t, err)
	require.Nil(t, active)
}
