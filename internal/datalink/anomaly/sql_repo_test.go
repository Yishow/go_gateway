package anomaly

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSQLRepositoryOpenEndAndSingleActive(t *testing.T) {
	db := openAnomalyTestDB(t, "dev-a")
	repo := NewSQLRepository(db)
	ctx := t.Context()

	value := "42"
	opened, err := repo.Open(ctx, Event{DeviceID: "dev-a", PointID: "p-d5", Kind: KindStale, StartedAt: anomalyEpoch, LastValue: &value, Detail: "no-successful-read"})
	require.NoError(t, err)
	require.Positive(t, opened.ID)

	_, err = repo.Open(ctx, Event{DeviceID: "dev-a", PointID: "p-d5", Kind: KindStale, StartedAt: anomalyEpoch.Add(time.Second)})
	require.ErrorIs(t, err, ErrAlreadyActive, "one active event per device, point and kind")

	_, err = repo.Open(ctx, Event{DeviceID: "dev-a", PointID: "p-d5", Kind: KindFrozen, StartedAt: anomalyEpoch})
	require.NoError(t, err, "another kind on the same point is independent")

	active, err := repo.ListActive(ctx)
	require.NoError(t, err)
	require.Len(t, active, 2)

	ended := anomalyEpoch.Add(12 * time.Second)
	require.NoError(t, repo.End(ctx, opened.ID, ended, nil, DetailMonitoringStopped))
	require.ErrorIs(t, repo.End(ctx, opened.ID, ended, nil, ""), ErrNotActive, "ending twice must not report success")

	events, total, err := repo.List(ctx, ListFilter{DeviceID: "dev-a", State: StateEnded, Kind: KindStale, Limit: 10})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, events, 1)
	require.Equal(t, anomalyEpoch, events[0].StartedAt)
	require.NotNil(t, events[0].EndedAt)
	require.Equal(t, ended, *events[0].EndedAt)
	require.Equal(t, "42", *events[0].LastValue, "nil last value on end keeps the observed value")
	require.Equal(t, "no-successful-read; monitoring-stopped", events[0].Detail, "end detail is appended to the cause")

	_, err = repo.Open(ctx, Event{DeviceID: "dev-a", PointID: "p-d5", Kind: KindStale, StartedAt: anomalyEpoch.Add(time.Minute)})
	require.NoError(t, err, "a new active event may open after the previous one ended")

	counts, err := repo.ActiveCounts(ctx)
	require.NoError(t, err)
	require.Equal(t, map[string]int{"dev-a": 2}, counts)
}

func TestSQLRepositoryRejectsInvalidEvents(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
	ctx := t.Context()
	cases := map[string]Event{
		"unknown kind":          {DeviceID: "dev-a", PointID: "p", Kind: Kind("hot"), StartedAt: anomalyEpoch},
		"missing device":        {PointID: "p", Kind: KindStale, StartedAt: anomalyEpoch},
		"stale without point":   {DeviceID: "dev-a", Kind: KindStale, StartedAt: anomalyEpoch},
		"disconnect with point": {DeviceID: "dev-a", PointID: "p", Kind: KindDisconnected, StartedAt: anomalyEpoch},
		"zero start":            {DeviceID: "dev-a", PointID: "p", Kind: KindStale},
	}
	for name, event := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := repo.Open(ctx, event)
			require.ErrorIs(t, err, ErrInvalidEvent)
		})
	}
	_, _, err := repo.List(ctx, ListFilter{DeviceID: "dev-a", State: State("open")})
	require.ErrorIs(t, err, ErrInvalidFilter)
}

func TestSQLRepositoryRetentionKeepsActiveEvents(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a", "dev-b"), WithRetention(5))
	ctx := t.Context()
	activeIDs := make([]int64, 0, 3)
	for i := range 3 {
		ev, err := repo.Open(ctx, Event{DeviceID: "dev-a", PointID: "active-" + string(rune('a'+i)), Kind: KindStale, StartedAt: anomalyEpoch.Add(time.Duration(i) * time.Second)})
		require.NoError(t, err)
		activeIDs = append(activeIDs, ev.ID)
	}
	endedIDs := make([]int64, 0, 2)
	for i := range 2 {
		ev, err := repo.Open(ctx, Event{DeviceID: "dev-a", PointID: "ended", Kind: KindStale, StartedAt: anomalyEpoch.Add(time.Duration(10+i) * time.Second)})
		require.NoError(t, err)
		require.NoError(t, repo.End(ctx, ev.ID, ev.StartedAt.Add(time.Second), nil, ""))
		endedIDs = append(endedIDs, ev.ID)
	}
	_, err := repo.Open(ctx, Event{DeviceID: "dev-b", Kind: KindDisconnected, StartedAt: anomalyEpoch})
	require.NoError(t, err)

	_, err = repo.Open(ctx, Event{DeviceID: "dev-a", Kind: KindDisconnected, StartedAt: anomalyEpoch.Add(time.Hour)})
	require.NoError(t, err)

	events, total, err := repo.List(ctx, ListFilter{DeviceID: "dev-a", Limit: 100})
	require.NoError(t, err)
	require.Equal(t, 5, total)
	ids := map[int64]bool{}
	for _, ev := range events {
		ids[ev.ID] = true
	}
	for _, id := range activeIDs {
		require.True(t, ids[id], "active events are never deleted")
	}
	require.False(t, ids[endedIDs[0]], "oldest ended event is removed first")
	require.True(t, ids[endedIDs[1]])

	_, totalB, err := repo.List(ctx, ListFilter{DeviceID: "dev-b"})
	require.NoError(t, err)
	require.Equal(t, 1, totalB, "retention is per device")
}

func TestSQLRepositoryPrunesExcessWhenAnActiveEventEnds(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"), WithRetention(1))
	first, err := repo.Open(t.Context(), Event{DeviceID: "dev-a", PointID: "p1", Kind: KindStale, StartedAt: anomalyEpoch})
	require.NoError(t, err)
	second, err := repo.Open(t.Context(), Event{DeviceID: "dev-a", PointID: "p2", Kind: KindStale, StartedAt: anomalyEpoch.Add(time.Second)})
	require.NoError(t, err)
	require.NoError(t, repo.End(t.Context(), first.ID, anomalyEpoch.Add(time.Minute), nil, ""))
	events, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a"})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Equal(t, second.ID, events[0].ID, "the remaining active event must not be deleted")
}
