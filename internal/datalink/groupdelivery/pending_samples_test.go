package groupdelivery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPendingSampleCountUsesFullScopeAndExcludesConsumedInputs(t *testing.T) {
	store, db, _ := newTestStore(t)
	keys := []GroupKey{testKey, testKey, {WorkspaceID: "workspace-B", GroupID: testKey.GroupID, GroupRevision: testKey.GroupRevision}, otherKey, {WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID, GroupRevision: "rev-2"}}
	for i, key := range keys {
		require.NoError(t, store.AppendSample(t.Context(), key, t0, sampleAt(string(rune('a'+i)), "member", i+1, 1)))
	}
	_, err := db.ExecContext(t.Context(), `UPDATE wg_delivery_samples SET consumed=1 WHERE sample_id='a'`)
	require.NoError(t, err)
	count, err := store.PendingSampleCount(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, 1, count)
	empty := GroupKey{WorkspaceID: "absent", GroupID: testKey.GroupID, GroupRevision: testKey.GroupRevision}
	count, err = store.PendingSampleCount(t.Context(), empty)
	require.NoError(t, err)
	require.Zero(t, count)
	_, err = store.PendingSampleCount(t.Context(), GroupKey{GroupID: testKey.GroupID})
	require.ErrorIs(t, err, ErrInvalidGroupKey)
}
