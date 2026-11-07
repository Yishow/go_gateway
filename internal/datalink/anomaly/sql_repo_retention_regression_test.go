package anomaly

import (
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAnomalyRetentionDefaultPrunesOldestEndedAtFiveThousand(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a", "dev-b"))
	ctx := t.Context()

	activeIDs := make([]int64, 0, 4)
	for i, pointID := range []string{"active-a", "active-b", "active-c"} {
		event, err := repo.Open(ctx, Event{DeviceID: "dev-a", PointID: pointID, Kind: KindStale, StartedAt: anomalyEpoch.Add(time.Duration(i) * time.Second)})
		require.NoError(t, err)
		activeIDs = append(activeIDs, event.ID)
	}

	var oldestSurvivorID int64
	for i := range 4997 {
		startedAt := anomalyEpoch.Add(time.Hour + time.Duration(i)*time.Second)
		event, err := repo.Open(ctx, Event{DeviceID: "dev-a", PointID: "history", Kind: KindStale, StartedAt: startedAt})
		require.NoError(t, err)
		require.NoError(t, repo.End(ctx, event.ID, startedAt.Add(time.Second), nil, ""))
		if i == 1 {
			oldestSurvivorID = event.ID
		}
	}

	other, err := repo.Open(ctx, Event{DeviceID: "dev-b", Kind: KindDisconnected, StartedAt: anomalyEpoch})
	require.NoError(t, err)
	_, total, err := repo.List(ctx, ListFilter{DeviceID: "dev-a"})
	require.NoError(t, err)
	require.Equal(t, 5000, total)

	newEvent, err := repo.Open(ctx, Event{DeviceID: "dev-a", Kind: KindDisconnected, StartedAt: anomalyEpoch.Add(3 * time.Hour)})
	require.NoError(t, err)
	activeIDs = append(activeIDs, newEvent.ID)

	active, activeTotal, err := repo.List(ctx, ListFilter{DeviceID: "dev-a", State: StateActive})
	require.NoError(t, err)
	require.Equal(t, 4, activeTotal)
	gotActiveIDs := make([]int64, 0, len(active))
	for _, event := range active {
		gotActiveIDs = append(gotActiveIDs, event.ID)
	}
	require.ElementsMatch(t, activeIDs, gotActiveIDs)

	oldestEndedPage, endedTotal, err := repo.List(ctx, ListFilter{DeviceID: "dev-a", State: StateEnded, Kind: KindStale, Limit: 1, Offset: 4995})
	require.NoError(t, err)
	require.Equal(t, 4996, endedTotal)
	require.Len(t, oldestEndedPage, 1)
	require.Equal(t, oldestSurvivorID, oldestEndedPage[0].ID)

	otherEvents, otherTotal, err := repo.List(ctx, ListFilter{DeviceID: "dev-b"})
	require.NoError(t, err)
	require.Equal(t, 1, otherTotal)
	require.Len(t, otherEvents, 1)
	require.Equal(t, other.ID, otherEvents[0].ID)
	counts, err := repo.ActiveCounts(ctx)
	require.NoError(t, err)
	require.Equal(t, map[string]int{"dev-a": 4, "dev-b": 1}, counts)
}

func TestAnomalyRetentionOverFiveThousandActivePrunesOnEnd(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a", "dev-b"))
	ctx := t.Context()

	activeIDs := make([]int64, 0, 5001)
	for i := range 5001 {
		event, err := repo.Open(ctx, Event{DeviceID: "dev-a", PointID: "point-" + strconv.Itoa(i), Kind: KindStale, StartedAt: anomalyEpoch.Add(time.Duration(i) * time.Second)})
		require.NoError(t, err)
		activeIDs = append(activeIDs, event.ID)
	}
	other, err := repo.Open(ctx, Event{DeviceID: "dev-b", Kind: KindDisconnected, StartedAt: anomalyEpoch})
	require.NoError(t, err)

	active, activeTotal, err := repo.List(ctx, ListFilter{DeviceID: "dev-a", State: StateActive})
	require.NoError(t, err)
	require.Equal(t, 5001, activeTotal)
	gotActiveIDs := make([]int64, 0, len(active))
	for _, event := range active {
		gotActiveIDs = append(gotActiveIDs, event.ID)
	}
	expectedIDs := slices.Clone(activeIDs)
	slices.Sort(expectedIDs)
	slices.Sort(gotActiveIDs)
	require.Equal(t, expectedIDs, gotActiveIDs)
	page, pageTotal, err := repo.List(ctx, ListFilter{DeviceID: "dev-a", State: StateActive, Kind: KindStale, Limit: 2, Offset: 4999})
	require.NoError(t, err)
	require.Equal(t, 5001, pageTotal)
	require.Len(t, page, 2)
	require.Equal(t, activeIDs[1], page[0].ID)
	require.Equal(t, activeIDs[0], page[1].ID)
	counts, err := repo.ActiveCounts(ctx)
	require.NoError(t, err)
	require.Equal(t, map[string]int{"dev-a": 5001, "dev-b": 1}, counts)

	require.NoError(t, repo.End(ctx, activeIDs[0], anomalyEpoch.Add(2*time.Hour), nil, ""))
	remaining, total, err := repo.List(ctx, ListFilter{DeviceID: "dev-a"})
	require.NoError(t, err)
	require.Equal(t, 5000, total)
	gotActiveIDs = nil
	for _, event := range remaining {
		require.Nil(t, event.EndedAt)
		gotActiveIDs = append(gotActiveIDs, event.ID)
	}
	expectedIDs = slices.Clone(activeIDs[1:])
	slices.Sort(expectedIDs)
	slices.Sort(gotActiveIDs)
	require.Equal(t, expectedIDs, gotActiveIDs)
	ended, endedTotal, err := repo.List(ctx, ListFilter{DeviceID: "dev-a", State: StateEnded})
	require.NoError(t, err)
	require.Equal(t, 0, endedTotal)
	require.Empty(t, ended)
	otherEvents, otherTotal, err := repo.List(ctx, ListFilter{DeviceID: "dev-b"})
	require.NoError(t, err)
	require.Equal(t, 1, otherTotal)
	require.Len(t, otherEvents, 1)
	require.Equal(t, other.ID, otherEvents[0].ID)
	counts, err = repo.ActiveCounts(ctx)
	require.NoError(t, err)
	require.Equal(t, map[string]int{"dev-a": 5000, "dev-b": 1}, counts)
}
