package anomaly

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFrozenRecoveryBetweenScansIsRecorded(t *testing.T) {
	repo := &failingRepo{Repository: NewSQLRepository(openAnomalyTestDB(t, "dev-a"))}
	clock := &testClock{now: anomalyEpoch}
	detector := NewDetector(repo, clock.Now)
	input := defaultInput(map[string]string{"D5": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: time.Hour, Frozen: 200 * time.Millisecond}})
	detector.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), OK: true, Value: 0})
	require.NoError(t, detector.Scan(t.Context(), input))
	for i := 1; i <= 8; i++ {
		clock.now = anomalyEpoch.Add(time.Duration(i) * 100 * time.Millisecond)
		detector.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), OK: true, Value: 0})
	}
	clock.now = anomalyEpoch.Add(900 * time.Millisecond)
	detector.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), OK: true, Value: 1})
	clock.now = anomalyEpoch.Add(time.Second)
	repo.failOpen = errors.New("storage unavailable")
	require.ErrorContains(t, detector.Scan(t.Context(), input), "storage unavailable")
	repo.failOpen = nil
	require.NoError(t, detector.Scan(t.Context(), input))
	require.NoError(t, detector.Scan(t.Context(), input))
	events, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a", Kind: KindFrozen})
	require.NoError(t, err)
	require.Equal(t, 1, total, "the completed 200–900 ms frozen interval must be recorded")
	require.Len(t, events, 1)
	require.Equal(t, anomalyEpoch.Add(200*time.Millisecond), events[0].StartedAt)
	require.NotNil(t, events[0].EndedAt)
	require.Equal(t, anomalyEpoch.Add(900*time.Millisecond), *events[0].EndedAt)
}
