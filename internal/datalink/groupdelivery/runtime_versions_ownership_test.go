package groupdelivery

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOpenJournalsVerifiesTransferredInputWithOneConnection(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, 1, batch.Accepted)
	healthy := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "healthy", GroupRevision: "rev-healthy"}
	require.NoError(t, store.AppendSample(t.Context(), healthy, at(10), sampleAt("healthy", "temp", 11, 6)))
	store.DB().SetMaxOpenConns(1)

	// Receipt verification reads cutover roots and their frozen versions. It
	// must share the pinned read connection, not wait for a second pool slot.
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	for range 2 {
		journals, err := store.OpenJournals(ctx)
		require.NoError(t, err, "valid transferred history must remain discoverable with one connection")
		require.ElementsMatch(t, []OpenJournal{
			{Key: otherKey, FirstBucket: at(20), LastBucket: at(20)},
			{Key: healthy, FirstBucket: at(10), LastBucket: at(10)},
		}, journals, "transferred source is excluded, successor and independent scope remain")
		require.Zero(t, store.DB().Stats().InUse, "discovery must release its read transaction")
	}
	require.Equal(t, 3, countRows(t, store.DB(), "wg_delivery_samples"))
	require.Zero(t, openCount(t, store.DB(), "wg_delivery_samples", "consumed=1"))
	require.Equal(t, 1, countRows(t, store.DB(), "wg_delivery_sample_handoffs"))
}
