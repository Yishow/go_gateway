package groupdelivery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyReadSlicesRemainNil(t *testing.T) {
	store, db, _ := newTestStore(t)
	heads, err := store.ReadyHeads(t.Context(), at(0))
	require.NoError(t, err)
	require.Nil(t, heads)
	sources, err := store.HandoffSources(t.Context(), testKey.WorkspaceID, GroupKey{}, 8)
	require.NoError(t, err)
	require.Nil(t, sources)
	journals, err := store.OpenJournals(t.Context())
	require.NoError(t, err)
	require.Nil(t, journals)
	receipts, err := queryHandoffs(t.Context(), db, "1=1")
	require.NoError(t, err)
	require.Nil(t, receipts)
	rows, err := store.sourceRows(t.Context(), SourceScope{}, "missing-partition", "", 8)
	require.NoError(t, err)
	require.Nil(t, rows)
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, tx.Rollback()) })
	revisions, err := LifecyclePendingRevisions(t.Context(), tx, testKey.WorkspaceID, testKey.GroupID)
	require.NoError(t, err)
	require.Nil(t, revisions)
	accepted, err := readLegacyAcceptedSet(t.Context(), tx, testKey)
	require.NoError(t, err)
	require.Nil(t, accepted)
}
