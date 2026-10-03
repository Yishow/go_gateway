package groupdelivery

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

func TestGroupStatusReportsRecentScopedQualityCausesWithoutValues(t *testing.T) {
	store, _, _ := newTestStore(t)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, Closure{
		Destination: testDestination, NextClose: at(30),
		Buckets: []ClosedBucket{
			{Outcome: snapshot.Outcome{BucketStart: t0, Kind: snapshot.OutcomeSkipped, Reason: "incomplete-required-member", Members: []snapshot.MemberResult{
				{MemberKey: "missing", Status: snapshot.MemberMissing},
				{MemberKey: "bad", Status: snapshot.MemberBad, Reason: "private-dsn-must-not-escape"},
				{MemberKey: "stale", Status: snapshot.MemberStale},
			}}},
			silentBucket("", at(10)),
		},
	}))
	require.NoError(t, store.CommitClosure(t.Context(), otherKey, Closure{
		Destination: testDestination, NextClose: at(30), Buckets: []ClosedBucket{silentBucket("", at(20))},
	}))
	status, err := store.GroupStatus(t.Context(), testKey.GroupID)
	require.NoError(t, err)
	encoded, err := json.Marshal(status)
	require.NoError(t, err)
	var got struct {
		RecentBucketIssues []struct {
			GroupRevision string   `json:"group_revision"`
			Kind          string   `json:"kind"`
			Causes        []string `json:"causes"`
		} `json:"RecentBucketIssues"`
	}
	require.NoError(t, json.Unmarshal(encoded, &got))
	require.Len(t, got.RecentBucketIssues, 2)
	require.Equal(t, "no_data", got.RecentBucketIssues[0].Kind)
	require.Equal(t, []string{"no_data"}, got.RecentBucketIssues[0].Causes)
	require.Equal(t, testKey.GroupRevision, got.RecentBucketIssues[1].GroupRevision)
	require.Equal(t, []string{"bad", "missing", "stale"}, got.RecentBucketIssues[1].Causes)
	require.NotContains(t, string(encoded), "private-dsn-must-not-escape")
}

func TestGroupStatusBoundsRecentBucketIssuesAndKeepsTotalCounts(t *testing.T) {
	store, _, _ := newTestStore(t)
	for i := range 25 {
		require.NoError(t, store.CommitClosure(t.Context(), testKey, Closure{
			Destination: testDestination, NextClose: at((i + 1) * 10),
			Buckets: []ClosedBucket{silentBucket("", at(i*10))},
		}))
	}
	status, err := store.GroupStatus(t.Context(), testKey.GroupID)
	require.NoError(t, err)
	require.Equal(t, 25, status.NoDataBuckets)
	encoded, err := json.Marshal(status)
	require.NoError(t, err)
	var got struct{ RecentBucketIssues []json.RawMessage }
	require.NoError(t, json.Unmarshal(encoded, &got))
	require.Len(t, got.RecentBucketIssues, 20)
}
