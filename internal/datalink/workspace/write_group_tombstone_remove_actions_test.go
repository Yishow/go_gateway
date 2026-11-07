package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDrainedTombstoneAllowsRemovalAndCanBePermanentlyDeleted(t *testing.T) {
	db, service, group := lifecycleReadFixture(t, false)
	group = seedLegacyRecordingTombstone(t, db, service, group.ID).Group

	view, err := service.Lifecycle(t.Context(), group.ID)
	require.NoError(t, err)
	require.Equal(t, WriteGroupLifecycleAction{Action: "remove", Allowed: true}, view.Actions[2])
	require.False(t, view.Actions[1].Allowed, "a tombstoned recording cannot be stopped")

	result, err := service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
	require.NoError(t, err)
	require.True(t, result.Deleted)
	var remaining int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_groups WHERE id=?`, group.ID).Scan(&remaining))
	require.Zero(t, remaining)
}
