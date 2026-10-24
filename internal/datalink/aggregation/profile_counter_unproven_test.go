package aggregation

import (
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

// TestProfileCounterUnprovenStaysNull keeps the not-yet-modeled discontinuities
// safe: nothing is repaired with abs, zero or a guessed rollover.
func TestProfileCounterUnprovenStaysNull(t *testing.T) {
	binding := counterBindingFixture()
	uint64Of := measurement.NewUint64

	t.Run("a decrease is uncertain and rebases", func(t *testing.T) {
		results, checkpoint := closeSeries(t, binding,
			[]CounterReading{counterReading("a", 0, time.Minute, uint64Of(100))},
			[]CounterReading{counterReading("b", 1, time.Minute, uint64Of(10))},
			[]CounterReading{counterReading("c", 2, time.Minute, uint64Of(25))},
		)
		require.Equal(t, CounterUncertain, results[1].Status)
		require.Equal(t, ReasonDecreaseUnproven, results[1].Reason)
		require.Nil(t, results[1].Delta, "no abs() and no zero for an unexplained decrease")
		require.Equal(t, "10", digitsOf(t, results[1].Raw.Value), "the raw reading is kept")
		require.Equal(t, CounterOK, results[2].Status)
		require.Equal(t, "15", digitsOf(t, results[2].Delta.Value), "the post-decrease reading is the next baseline")
		require.Equal(t, "c", checkpoint.Reading.SampleID)
	})

	t.Run("a decrease in the middle of a window is not hidden by its endpoints", func(t *testing.T) {
		results, _ := closeSeries(t, binding,
			[]CounterReading{counterReading("a", 0, time.Minute, uint64Of(90))},
			[]CounterReading{
				counterReading("b1", 1, time.Minute, uint64Of(100)),
				counterReading("b2", 1, 2*time.Minute, uint64Of(10)),
				counterReading("b3", 1, 3*time.Minute, uint64Of(110)),
			},
		)
		require.Equal(t, CounterUncertain, results[1].Status)
		require.Nil(t, results[1].Delta)
		require.Equal(t, "110", digitsOf(t, results[1].Raw.Value))
	})

	t.Run("contradictory readings at one instant", func(t *testing.T) {
		results, _ := closeSeries(t, binding,
			[]CounterReading{counterReading("a", 0, time.Minute, uint64Of(1))},
			[]CounterReading{
				counterReading("b1", 1, time.Minute, uint64Of(5)),
				counterReading("b2", 1, time.Minute, uint64Of(6)),
			},
		)
		require.Equal(t, CounterUncertain, results[1].Status)
		require.Equal(t, ReasonContradictory, results[1].Reason)
		require.Nil(t, results[1].Delta)
	})

	t.Run("a window after missing windows is a span total, not one window of usage", func(t *testing.T) {
		var checkpoint *CounterCheckpoint
		_, checkpoint, err := CloseCounterWindow(checkpoint, binding, counterWindow(0),
			[]CounterReading{counterReading("a", 0, time.Minute, uint64Of(1000))})
		require.NoError(t, err)
		result, _, err := CloseCounterWindow(checkpoint, binding, counterWindow(2),
			[]CounterReading{counterReading("b", 2, time.Minute, uint64Of(1002))})
		require.NoError(t, err)
		require.Equal(t, CounterGap, result.Status)
		require.True(t, result.Delta.Span)
		require.Equal(t, "1002", digitsOf(t, result.Raw.Value))
	})

	t.Run("a changed binding starts a new baseline", func(t *testing.T) {
		_, checkpoint := closeSeries(t, binding,
			[]CounterReading{counterReading("a", 0, time.Minute, uint64Of(1000))})
		next := binding
		next.SeriesEpoch = "epoch-2"
		result, _, err := CloseCounterWindow(checkpoint, next, counterWindow(1),
			[]CounterReading{counterReading("b", 1, time.Minute, uint64Of(3))})
		require.NoError(t, err)
		require.Equal(t, CounterBaseline, result.Status)
		require.Equal(t, ReasonBindingChanged, result.Reason)
		require.Nil(t, result.Delta, "readings of different epochs are never subtracted")
	})

	t.Run("an incomplete binding is refused", func(t *testing.T) {
		incomplete := binding
		incomplete.Kind = ""
		_, _, err := CloseCounterWindow(nil, incomplete, counterWindow(0),
			[]CounterReading{counterReading("a", 0, time.Minute, uint64Of(1))})
		require.ErrorIs(t, err, ErrCounterBindingIncomplete)
	})

	t.Run("a non numeric reading carries no value", func(t *testing.T) {
		result, checkpoint, err := CloseCounterWindow(nil, binding, counterWindow(0),
			[]CounterReading{counterReading("a", 0, time.Minute, measurement.NewText("12"))})
		require.NoError(t, err)
		require.Equal(t, CounterInvalid, result.Status)
		require.Nil(t, result.Raw)
		require.Nil(t, result.Delta)
		require.Nil(t, checkpoint, "an invalid reading never becomes a baseline")
	})

	t.Run("checkpoint json refuses unknown shapes", func(t *testing.T) {
		_, checkpoint := closeSeries(t, binding,
			[]CounterReading{counterReading("a", 0, time.Minute, uint64Of(7))})
		encoded, err := EncodeCounterCheckpoint(*checkpoint)
		require.NoError(t, err)
		_, err = DecodeCounterCheckpoint([]byte(strings.Replace(string(encoded), `"version":2`, `"version":3`, 1)))
		require.ErrorIs(t, err, ErrCounterCheckpointInvalid)
		_, err = DecodeCounterCheckpoint([]byte(strings.Replace(string(encoded), `"version":2`, `"version":2,"extra":true`, 1)))
		require.ErrorIs(t, err, ErrCounterCheckpointInvalid)
		_, err = DecodeCounterCheckpoint([]byte(`{}`))
		require.ErrorIs(t, err, ErrCounterCheckpointInvalid)
	})
}
