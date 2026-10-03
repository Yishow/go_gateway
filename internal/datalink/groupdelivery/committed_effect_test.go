package groupdelivery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCommittedEffectRequiresMatchingReceiptAndAppliedRevision(t *testing.T) {
	store, _, _ := newTestStore(t)
	old := enqueueFor(t, store, testKey, testDestination, "entity", t0, 1)
	item, err := store.GetOutbox(t.Context(), old)
	require.NoError(t, err)
	claim, err := store.BeginDelivery(t.Context(), old, "writer", time.Minute)
	require.NoError(t, err)
	require.NoError(t, store.CompleteDelivery(t.Context(), claim, item.PayloadDigest))
	effect, err := store.LastCommittedEffect(t.Context(), testKey.GroupID, testKey.GroupRevision)
	require.NoError(t, err)
	require.NotNil(t, effect)
	require.Equal(t, item.RecordID, effect.RecordID)
	require.Equal(t, old, effect.EffectKey)
	require.Equal(t, item.PayloadDigest, effect.PayloadDigest)
	require.Equal(t, testDestination.ConnectorRevision, effect.ConnectorRevision)
	require.False(t, effect.CommittedAt.IsZero())
	for _, revision := range []string{"", "new-revision"} {
		got, err := store.LastCommittedEffect(t.Context(), testKey.GroupID, revision)
		require.NoError(t, err)
		require.Nil(t, got, "old SQL commits must not count as the new revision's first record")
	}
	got, err := store.LastCommittedEffect(t.Context(), otherKey.GroupID, testKey.GroupRevision)
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestRevisionStageEvidenceDoesNotAdoptHistoricalRows(t *testing.T) {
	store, _, _ := newTestStore(t)
	enqueueFor(t, store, testKey, testDestination, "old", t0, 1)
	key := testKey
	key.GroupRevision = "current-revision"
	enqueueFor(t, store, key, testDestination, "new", t0.Add(10*time.Second), 2)
	stages, err := store.RevisionStageCounts(t.Context(), key.GroupID, key.GroupRevision)
	require.NoError(t, err)
	require.Equal(t, 1, stages.Queued)
	stages, err = store.RevisionStageCounts(t.Context(), key.GroupID, "not-applied")
	require.NoError(t, err)
	require.Equal(t, StageCounts{}, stages)
}

func TestCommittedEffectRejectsQueuedUnknownAndMismatchedReceipt(t *testing.T) {
	store, _, _ := newTestStore(t)
	queued := enqueueFor(t, store, testKey, testDestination, "queued", t0, 1)
	unknown := enqueueFor(t, store, testKey, testDestination, "unknown", t0, 2)
	claim, err := store.BeginDelivery(t.Context(), unknown, "writer", time.Minute)
	require.NoError(t, err)
	require.NoError(t, store.MarkUnknown(t.Context(), claim, "commit-ambiguous"))
	got, err := store.LastCommittedEffect(t.Context(), testKey.GroupID, testKey.GroupRevision)
	require.NoError(t, err)
	require.Nil(t, got)
	claim, err = store.BeginDelivery(t.Context(), queued, "writer", time.Minute)
	require.NoError(t, err)
	require.NoError(t, store.CompleteDelivery(t.Context(), claim, "mismatched-digest"))
	got, err = store.LastCommittedEffect(t.Context(), testKey.GroupID, testKey.GroupRevision)
	require.NoError(t, err)
	require.Nil(t, got, "a mismatched receipt is not evidence for this payload")
}
