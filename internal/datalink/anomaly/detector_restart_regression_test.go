package anomaly

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFrozenEventEndsWhenValueChangedAcrossRestart(t *testing.T) {
	for _, observeBeforeLoad := range []bool{false, true} {
		t.Run(map[bool]string{false: "load first", true: "observation first"}[observeBeforeLoad], func(t *testing.T) {
			repo := NewSQLRepository(openAnomalyTestDB(t, "dev-a"))
			previous := "1"
			_, err := repo.Open(t.Context(), Event{DeviceID: "dev-a", PointID: "D5", Kind: KindFrozen, StartedAt: anomalyEpoch.Add(-time.Minute), LastValue: &previous})
			require.NoError(t, err)
			clock := &testClock{now: anomalyEpoch}
			d := NewDetector(repo, clock.Now)
			observe := func() { d.Observe(Observation{DeviceID: "dev-a", PointID: "D5", At: clock.Now(), OK: true, Value: 2}) }
			if observeBeforeLoad {
				observe()
			}
			require.NoError(t, d.Load(t.Context()))
			if !observeBeforeLoad {
				observe()
			}
			input := defaultInput(map[string]string{"D5": "dev-a"}, map[string]Thresholds{"dev-a": {Stale: 10 * time.Second, Frozen: time.Minute}})
			require.NoError(t, d.Scan(t.Context(), input))
			events, total, err := repo.List(t.Context(), ListFilter{DeviceID: "dev-a", Kind: KindFrozen, State: StateEnded})
			require.NoError(t, err)
			require.Equal(t, 1, total)
			require.Equal(t, anomalyEpoch, *events[0].EndedAt)
		})
	}
}
