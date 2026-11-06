package anomaly

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDisconnectRecoveryBetweenScansIsRecorded(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
	clock := &testClock{now: anomalyEpoch}
	d := NewDetector(repo, clock.Now)
	input := defaultInput(map[string]string{"D5": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: 10 * time.Second}})
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), OK: true, Value: 7})
	require.NoError(t, d.Scan(t.Context(), input))
	clock.Advance(100 * time.Millisecond)
	started := clock.Now()
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: started, Reason: ReasonConnectionFailed})
	clock.Advance(100 * time.Millisecond)
	ended := clock.Now()
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: ended, OK: true, Value: 8})
	clock.Advance(800 * time.Millisecond)
	require.NoError(t, d.Scan(t.Context(), input))
	events, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a", Kind: KindDisconnected, State: StateEnded})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Equal(t, started, events[0].StartedAt)
	require.Equal(t, ended, *events[0].EndedAt)
}

func TestRecoveredDisconnectRetriesAfterStorageFailure(t *testing.T) {
	repo := &failingRepo{Repository: NewSQLRepository(openAnomalyTestDB(t, "dev-a"))}
	clock := &testClock{now: anomalyEpoch}
	d := NewDetector(repo, clock.Now)
	input := defaultInput(map[string]string{"D5": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: 10 * time.Second}})
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), OK: true})
	require.NoError(t, d.Scan(t.Context(), input))
	clock.Advance(time.Second)
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), Reason: ReasonConnectionFailed})
	clock.Advance(time.Second)
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), OK: true})
	repo.failOpen = errors.New("storage unavailable")
	require.ErrorContains(t, d.Scan(t.Context(), input), "storage unavailable")
	repo.failOpen = nil
	require.NoError(t, d.Scan(t.Context(), input))
	_, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a", Kind: KindDisconnected, State: StateEnded})
	require.NoError(t, err)
	require.Equal(t, 1, total)
}

func TestLoadAdoptsRecoveryObservedBeforeFirstScan(t *testing.T) {
	for _, kind := range []Kind{KindDisconnected, KindStale} {
		t.Run(string(kind), func(t *testing.T) {
			repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
			event := Event{DeviceID: "dev-a", Kind: kind, StartedAt: anomalyEpoch.Add(-time.Minute)}
			if kind == KindStale {
				event.PointID = "D5"
			}
			_, err := repo.Open(t.Context(), event)
			require.NoError(t, err)
			clock := &testClock{now: anomalyEpoch}
			d := NewDetector(repo, clock.Now)
			d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), OK: true, Value: 3})
			require.NoError(t, d.Load(t.Context()))
			input := defaultInput(map[string]string{"D5": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: 10 * time.Second}})
			require.NoError(t, d.Scan(t.Context(), input))
			_, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a", Kind: kind, State: StateEnded})
			require.NoError(t, err)
			require.Equal(t, 1, total)
		})
	}
}

func TestSeparateDisconnectsBetweenScansAreRecorded(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
	clock := &testClock{now: anomalyEpoch}
	d := NewDetector(repo, clock.Now)
	input := defaultInput(map[string]string{"D5": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: 10 * time.Second}})
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), OK: true, Value: 7})
	require.NoError(t, d.Scan(t.Context(), input))
	for range 2 {
		clock.Advance(100 * time.Millisecond)
		d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), Reason: ReasonConnectionFailed})
		clock.Advance(100 * time.Millisecond)
		d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), OK: true, Value: 7})
	}
	require.NoError(t, d.Scan(t.Context(), input))
	_, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a", Kind: KindDisconnected, State: StateEnded})
	require.NoError(t, err)
	require.Equal(t, 2, total)
}

func TestStaleRecoveryBetweenScansIsRecorded(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
	clock := &testClock{now: anomalyEpoch}
	d := NewDetector(repo, clock.Now)
	input := defaultInput(map[string]string{"D5": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: 10 * time.Second}})
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), OK: true, Value: 7})
	require.NoError(t, d.Scan(t.Context(), input))
	clock.Advance(10 * time.Second)
	require.NoError(t, d.Scan(t.Context(), input))
	clock.Advance(100 * time.Millisecond)
	ended := clock.Now()
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: ended, OK: true, Value: 8})
	clock.Advance(900 * time.Millisecond)
	require.NoError(t, d.Scan(t.Context(), input))
	events, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a", Kind: KindStale, State: StateEnded})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Equal(t, anomalyEpoch.Add(10*time.Second), events[0].StartedAt)
	require.Equal(t, ended, *events[0].EndedAt)
}
