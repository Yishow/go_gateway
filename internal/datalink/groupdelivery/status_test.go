package groupdelivery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func enqueueFor(t *testing.T, store *Store, key GroupKey, destination FrozenDestination, entity string, start time.Time, counter uint64) string {
	t.Helper()
	bucket := rowBucket(entity, start, counter)
	require.NoError(t, store.CommitClosure(t.Context(), key, Closure{Destination: destination, NextClose: start.Add(10 * time.Second), Buckets: []ClosedBucket{bucket}}))
	return bucket.Outcome.EffectKey
}

func TestRevisionBoundBacklogStatusCountsEachRowInExactlyOneStage(t *testing.T) {
	store, _, _ := newTestStore(t)
	seedSamples(t, store) // 3 accepted samples, bucket still open
	destination := testDestination

	committed := enqueueFor(t, store, testKey, destination, "e1", t0, 1)
	claim, err := store.BeginDelivery(t.Context(), committed, "w1", time.Minute)
	require.NoError(t, err)
	require.NoError(t, store.CompleteDelivery(t.Context(), claim, "d"))

	enqueueFor(t, store, testKey, destination, "e2", t0, 2) // pending

	retrying := enqueueFor(t, store, testKey, destination, "e3", t0, 3)
	claim, err = store.BeginDelivery(t.Context(), retrying, "w1", time.Minute)
	require.NoError(t, err)
	require.NoError(t, store.MarkRetry(t.Context(), claim, StateRetrying, "insert-failed", time.Now().Add(time.Hour)))

	blocked := enqueueFor(t, store, testKey, destination, "e4", t0, 4)
	require.NoError(t, store.MarkBlocked(t.Context(), blocked, "target-blocked"))

	quarantined := enqueueFor(t, store, testKey, destination, "e5", t0, 5)
	claim, err = store.BeginDelivery(t.Context(), quarantined, "w1", time.Minute)
	require.NoError(t, err)
	require.NoError(t, store.MarkQuarantined(t.Context(), claim, "destination-rejected-row"))

	unknown := enqueueFor(t, store, testKey, destination, "e6", t0, 6)
	claim, err = store.BeginDelivery(t.Context(), unknown, "w1", time.Minute)
	require.NoError(t, err)
	require.NoError(t, store.MarkUnknown(t.Context(), claim, "commit-ambiguous"))

	skipped := enqueueFor(t, store, testKey, destination, "e7", t0, 7)
	require.NoError(t, store.MarkBlocked(t.Context(), skipped, "x"))
	require.NoError(t, store.ResolveQuarantine(t.Context(), skipped, ResolutionSkip))

	sending := enqueueFor(t, store, testKey, destination, "e8", t0, 8)
	_, err = store.BeginDelivery(t.Context(), sending, "w1", time.Minute)
	require.NoError(t, err)

	enqueueFor(t, store, otherKey, destination, "other", t0, 9) // another group: never counted

	status, err := store.GroupStatus(t.Context(), testKey.GroupID)
	require.NoError(t, err)
	require.Equal(t, StageCounts{
		Collecting: 1, Queued: 2, Retrying: 1, Blocked: 1, Quarantined: 1, Unknown: 1, SQLCommitted: 1, Skipped: 1,
	}, status.Stages, "pending and sending are both 'queued'; each row is in exactly one stage")
	require.NotNil(t, status.LastSQLCommittedAt)
	require.Positive(t, status.OldestPendingSeconds)
}

func TestRevisionBoundBacklogBufferedIsNeverSQLCommitted(t *testing.T) {
	store, _, _ := newTestStore(t)
	enqueueFor(t, store, testKey, testDestination, "e1", t0, 1)
	enqueueFor(t, store, testKey, testDestination, "e2", t0.Add(10*time.Second), 2)

	status, err := store.GroupStatus(t.Context(), testKey.GroupID)
	require.NoError(t, err)
	require.Equal(t, 2, status.Stages.Queued)
	require.Zero(t, status.Stages.SQLCommitted)
	require.Nil(t, status.LastSQLCommittedAt, "accepted and queued data is not committed data: last committed stays unset")

	empty, err := store.GroupStatus(t.Context(), "never-seen")
	require.NoError(t, err)
	require.Equal(t, StageCounts{}, empty.Stages)
	require.Nil(t, empty.LastSQLCommittedAt)
	require.Empty(t, empty.Backlog)
}

func TestRevisionBoundBacklogKeepsOldDestinationAndRevisionPerBacklog(t *testing.T) {
	store, _, _ := newTestStore(t)
	oldDestination := testDestination
	newKey := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID, GroupRevision: "rev-2"}
	newDestination := testDestination
	newDestination.ConnectorID, newDestination.ConnectorRevision, newDestination.TableName = "connector-2", "crev-9", "readings_v2"

	old := enqueueFor(t, store, testKey, oldDestination, "e1", t0, 1)
	enqueueFor(t, store, testKey, oldDestination, "e1", t0.Add(10*time.Second), 2)
	require.NoError(t, store.MarkBlocked(t.Context(), old, "target-blocked"))
	enqueueFor(t, store, newKey, newDestination, "e1", t0.Add(20*time.Second), 3)

	status, err := store.GroupStatus(t.Context(), testKey.GroupID)
	require.NoError(t, err)
	require.Len(t, status.Backlog, 2)
	byRevision := map[string]RevisionBacklog{}
	for _, entry := range status.Backlog {
		byRevision[entry.GroupRevision] = entry
	}
	require.Equal(t, "connector-1", byRevision["rev-1"].ConnectorID, "old rows keep the destination they were accepted for")
	require.Equal(t, "crev-1", byRevision["rev-1"].ConnectorRevision)
	require.Equal(t, "readings", byRevision["rev-1"].TableName)
	require.Equal(t, 2, byRevision["rev-1"].Pending)
	require.Equal(t, []string{"target-blocked"}, byRevision["rev-1"].LastErrorCodes)
	require.Equal(t, "connector-2", byRevision["rev-2"].ConnectorID)
	require.Equal(t, "crev-9", byRevision["rev-2"].ConnectorRevision)
	require.Equal(t, 1, byRevision["rev-2"].Pending)
}

func TestRevisionBoundBacklogStatusCountsSilentAndSkippedBuckets(t *testing.T) {
	store, _, _ := newTestStore(t)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, Closure{
		Destination: testDestination, NextClose: at(30),
		Buckets: []ClosedBucket{silentBucket("", t0), silentBucket("", at(10)), {Outcome: skippedOutcome(at(20))}},
	}))
	status, err := store.GroupStatus(t.Context(), testKey.GroupID)
	require.NoError(t, err)
	require.Equal(t, 2, status.NoDataBuckets)
	require.Equal(t, 1, status.SkippedBuckets)
	require.Zero(t, status.Stages.Queued, "a bucket without a row is not backlog")
}
