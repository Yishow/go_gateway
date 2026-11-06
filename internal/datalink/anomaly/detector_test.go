package anomaly

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func defaultInput(points map[string]string, thresholds map[string]Thresholds) ScanInput {
	return ScanInput{Monitored: points, Thresholds: thresholds}
}

func scanEachSecond(t *testing.T, d *Detector, clock *testClock, seconds int, input ScanInput) {
	t.Helper()
	for range seconds {
		clock.Advance(time.Second)
		require.NoError(t, d.Scan(t.Context(), input))
	}
}

func TestDefaultStaleThreshold(t *testing.T) {
	require.Equal(t, 10*time.Second, StaleThreshold(nil, 1000), "1000ms interval: 5s is below the 10s floor")
	require.Equal(t, 15*time.Second, StaleThreshold(nil, 3000))
	custom := 2500
	require.Equal(t, 2500*time.Millisecond, StaleThreshold(&custom, 1000), "a device setting wins over the default")
	require.Equal(t, time.Duration(0), FrozenThreshold(0), "frozen detection is disabled by default")
	require.Equal(t, 600*time.Second, FrozenThreshold(600000))
}

func TestStalePointOpensAndClosesOneEvent(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
	clock := &testClock{now: anomalyEpoch}
	d := NewDetector(repo, clock.Now)
	input := defaultInput(map[string]string{"D5": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: StaleThreshold(nil, 1000)}})

	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), OK: true, Value: int64(7)})
	require.NoError(t, d.Scan(t.Context(), input))
	scanEachSecond(t, d, clock, 10, input)
	active, err := repo.ListActive(t.Context())
	require.NoError(t, err)
	require.Empty(t, active, "exactly 10s without a read does not exceed the threshold")

	scanEachSecond(t, d, clock, 2, input)
	active, err = repo.ListActive(t.Context())
	require.NoError(t, err)
	require.Len(t, active, 1)
	require.Equal(t, KindStale, active[0].Kind)
	require.Equal(t, "D5", active[0].PointID)
	require.Equal(t, anomalyEpoch.Add(10*time.Second), active[0].StartedAt, "starts when the threshold is exceeded")
	require.Equal(t, "7", *active[0].LastValue)

	recovered := clock.Now()
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: recovered, OK: true, Value: int64(8)})
	scanEachSecond(t, d, clock, 3, input)

	events, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a", Limit: 10})
	require.NoError(t, err)
	require.Equal(t, 1, total, "one stale event for the whole outage")
	require.NotNil(t, events[0].EndedAt)
	require.Equal(t, recovered, *events[0].EndedAt, "ends at the successful read")
	require.Equal(t, DetailNoSuccessfulRead, events[0].Detail, "a natural recovery keeps the cause")
	require.Equal(t, "8", *events[0].LastValue)
}

func TestStaleCountsFromFirstFailureWhenPointNeverRead(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
	clock := &testClock{now: anomalyEpoch}
	d := NewDetector(repo, clock.Now)
	input := defaultInput(map[string]string{"D5": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: 10 * time.Second}})
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), Reason: "scan-deadline"})
	scanEachSecond(t, d, clock, 11, input)
	active, err := repo.ListActive(t.Context())
	require.NoError(t, err)
	require.Len(t, active, 1)
	require.Equal(t, anomalyEpoch.Add(10*time.Second), active[0].StartedAt)
	require.Nil(t, active[0].LastValue, "no value has ever been read")
}

func TestDisconnectIsRecordedOnce(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
	clock := &testClock{now: anomalyEpoch}
	d := NewDetector(repo, clock.Now)
	input := defaultInput(map[string]string{"D5": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: time.Hour}})

	firstFailure := clock.Now()
	for i := range 30 {
		reason := ReasonConnectionFailed
		if i%3 == 2 {
			reason = ReasonCircuitOpen
		}
		d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), Reason: reason})
		input.BreakerOpen = map[string]bool{"dev-a": i >= 5}
		clock.Advance(time.Second)
		require.NoError(t, d.Scan(t.Context(), input))
	}
	recovered := clock.Now()
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: recovered, OK: true, Value: 1})
	input.BreakerOpen = nil
	scanEachSecond(t, d, clock, 2, input)

	events, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a", Kind: KindDisconnected, Limit: 10})
	require.NoError(t, err)
	require.Equal(t, 1, total, "repeated failures for 30s produce exactly one disconnected event")
	require.Empty(t, events[0].PointID)
	require.Equal(t, firstFailure, events[0].StartedAt)
	require.NotNil(t, events[0].EndedAt)
	require.Equal(t, recovered, *events[0].EndedAt)
	require.Equal(t, ReasonConnectionFailed, events[0].Detail, "detail is a reason code, not a raw backend error")
}

func TestOpenBreakerAloneOpensDisconnect(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
	clock := &testClock{now: anomalyEpoch}
	d := NewDetector(repo, clock.Now)
	input := ScanInput{Monitored: map[string]string{"D5": "dev-a"}, Thresholds: map[string]Thresholds{"dev-a": {Stale: time.Hour}}, BreakerOpen: map[string]bool{"dev-a": true}}
	require.NoError(t, d.Scan(t.Context(), input))
	active, err := repo.ListActive(t.Context())
	require.NoError(t, err)
	require.Len(t, active, 1)
	require.Equal(t, KindDisconnected, active[0].Kind)
	require.Equal(t, ReasonBreakerOpen, active[0].Detail)
}

func TestFrozenDetectionIsOptIn(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
	clock := &testClock{now: anomalyEpoch}
	d := NewDetector(repo, clock.Now)
	input := defaultInput(map[string]string{"D0": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: 10 * time.Second}})
	tick := func(seconds int) {
		for range seconds {
			d.Observe(Observation{DeviceID: "dev-a", PointID: "D0", At: clock.Now(), OK: true, Value: 0})
			clock.Advance(time.Second)
			require.NoError(t, d.Scan(t.Context(), input))
		}
	}
	tick(3600)
	_, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a"})
	require.NoError(t, err)
	require.Zero(t, total, "no frozen event without a frozen threshold")

	input.Thresholds["dev-a"] = Thresholds{Stale: 10 * time.Second, Frozen: 600 * time.Second}
	tick(1)
	active, err := repo.ListActive(t.Context())
	require.NoError(t, err)
	require.Len(t, active, 1, "unchanged for longer than 600s opens one frozen event")
	require.Equal(t, KindFrozen, active[0].Kind)
	require.Equal(t, anomalyEpoch.Add(600*time.Second), active[0].StartedAt)
	tick(5)
	active, err = repo.ListActive(t.Context())
	require.NoError(t, err)
	require.Len(t, active, 1)

	changed := clock.Now()
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D0", At: changed, OK: true, Value: 1})
	clock.Advance(time.Second)
	require.NoError(t, d.Scan(t.Context(), input))
	events, _, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a", State: StateEnded, Kind: KindFrozen})
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, changed, *events[0].EndedAt, "ends when the value changes")
}

func TestUnmonitoredPointEndsItsEvents(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
	clock := &testClock{now: anomalyEpoch}
	d := NewDetector(repo, clock.Now)
	input := defaultInput(map[string]string{"D5": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: 10 * time.Second}})
	require.NoError(t, d.Scan(t.Context(), input))
	scanEachSecond(t, d, clock, 11, input)
	stopped := clock.Now().Add(time.Second)
	scanEachSecond(t, d, clock, 1, ScanInput{})
	events, _, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a", State: StateEnded})
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, stopped, *events[0].EndedAt)
	require.Equal(t, DetailNoSuccessfulRead+"; "+DetailMonitoringStopped, events[0].Detail)
}

type failingRepo struct {
	Repository
	failOpen error
	opens    int
}

func (r *failingRepo) Open(ctx context.Context, ev Event) (Event, error) {
	r.opens++
	if r.failOpen != nil {
		return Event{}, r.failOpen
	}
	return r.Repository.Open(ctx, ev)
}

func TestRepositoryFailureIsReportedAndRetried(t *testing.T) {
	repo := &failingRepo{Repository: NewSQLRepository(openAnomalyTestDB(t, "dev-a")), failOpen: errors.New("disk full")}
	clock := &testClock{now: anomalyEpoch}
	d := NewDetector(repo, clock.Now)
	input := ScanInput{Monitored: map[string]string{"D5": "dev-a"}, Thresholds: map[string]Thresholds{"dev-a": {Stale: time.Hour}}, BreakerOpen: map[string]bool{"dev-a": true}}
	require.ErrorContains(t, d.Scan(t.Context(), input), "disk full", "a failed write is never reported as recorded")
	repo.failOpen = nil
	clock.Advance(time.Second)
	require.NoError(t, d.Scan(t.Context(), input))
	active, err := repo.ListActive(t.Context())
	require.NoError(t, err)
	require.Len(t, active, 1)
	require.Equal(t, anomalyEpoch, active[0].StartedAt, "the retry keeps the original start")
	require.Equal(t, 2, repo.opens)
}

func TestLoadAdoptsActiveEventsAcrossRestart(t *testing.T) {
	repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
	clock := &testClock{now: anomalyEpoch}
	_, err := repo.Open(t.Context(), Event{DeviceID: "dev-a", PointID: "D5", Kind: KindStale, StartedAt: anomalyEpoch.Add(-time.Minute)})
	require.NoError(t, err)

	d := NewDetector(repo, clock.Now)
	require.NoError(t, d.Load(t.Context()))
	input := defaultInput(map[string]string{"D5": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: 10 * time.Second}})
	scanEachSecond(t, d, clock, 20, input)
	recovered := clock.Now()
	d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: recovered, OK: true, Value: 3})
	scanEachSecond(t, d, clock, 1, input)

	events, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a"})
	require.NoError(t, err)
	require.Equal(t, 1, total, "the adopted event is not duplicated")
	require.Equal(t, recovered, *events[0].EndedAt)
}
