package aggregation

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

var counterBase = time.Date(2026, 10, 24, 10, 0, 0, 0, time.UTC)

const counterWindowLen = 10 * time.Minute

func counterBindingFixture() CounterBinding {
	return CounterBinding{
		WorkspaceID: "workspace-A", GroupID: "group-C", ProfileRevision: "rev-1", MemberKey: "member-kwh",
		MeasurementID: "meas-kwh", DefinitionRevision: "def-1", SourceBindingRevision: "bind-1", SeriesEpoch: "epoch-1",
		Kind: CounterKindMonotonic, Unit: "kWh",
	}
}

// counterWindow is the n-th ten minute window after counterBase.
func counterWindow(n int) CounterWindow {
	start := counterBase.Add(time.Duration(n) * counterWindowLen)
	return CounterWindow{Start: start, End: start.Add(counterWindowLen)}
}

func counterReading(id string, window int, offset time.Duration, value measurement.ExactValue) CounterReading {
	observed := counterWindow(window).Start.Add(offset)
	return CounterReading{SampleID: id, ObservedAt: observed, ReceivedAt: observed.Add(40 * time.Millisecond), Value: value}
}

func decimalOf(t *testing.T, digits string) measurement.ExactValue {
	t.Helper()
	value, err := measurement.NewDecimal(digits)
	require.NoError(t, err)
	return value
}

func digitsOf(t *testing.T, value measurement.ExactValue) string {
	t.Helper()
	digits, ok := value.DecimalDigits()
	require.True(t, ok, "value %v has decimal digits", value.Type())
	return digits
}

// closeSeries folds one window per entry, in order, and returns every result
// plus the final checkpoint.
func closeSeries(t *testing.T, binding CounterBinding, windows ...[]CounterReading) ([]CounterWindowResult, *CounterCheckpoint) {
	t.Helper()
	var checkpoint *CounterCheckpoint
	results := make([]CounterWindowResult, 0, len(windows))
	for i, readings := range windows {
		result, next, err := CloseCounterWindow(checkpoint, binding, counterWindow(i), readings)
		require.NoError(t, err, "window %d", i)
		results = append(results, result)
		checkpoint = next
	}
	return results, checkpoint
}

// TestProfileCounterExactDelta pins the cumulative contract: raw readings and
// exact differences, with a first window that only establishes a baseline.
func TestProfileCounterExactDelta(t *testing.T) {
	t.Run("decimal progression 12000.0 12000.4 12000.7", func(t *testing.T) {
		results, _ := closeSeries(t, counterBindingFixture(),
			[]CounterReading{counterReading("s1", 0, 9*time.Minute, decimalOf(t, "12000.0"))},
			[]CounterReading{counterReading("s2", 1, 9*time.Minute, decimalOf(t, "12000.4"))},
			[]CounterReading{counterReading("s3", 2, 9*time.Minute, decimalOf(t, "12000.7"))},
		)
		require.Len(t, results, 3)

		first := results[0]
		require.Equal(t, CounterBaseline, first.Status)
		require.Equal(t, ReasonFirstReading, first.Reason)
		require.NotNil(t, first.Raw)
		require.Equal(t, "12000.0", digitsOf(t, first.Raw.Value), "raw keeps the digits the source gave")
		require.Nil(t, first.Delta, "a baseline window has a NULL delta, not zero")

		second := results[1]
		require.Equal(t, CounterOK, second.Status)
		require.Equal(t, "12000.4", digitsOf(t, second.Raw.Value))
		require.NotNil(t, second.Delta)
		require.Equal(t, "0.4", digitsOf(t, second.Delta.Value))
		require.Equal(t, DeltaBasisExact, second.Delta.Basis)
		require.Equal(t, measurement.ExactDecimal, second.Delta.Value.Type())
		require.Equal(t, "s1", second.Delta.BaselineSampleID)
		require.Equal(t, "s2", second.Delta.EndSampleID)
		require.Equal(t, counterBase.Add(9*time.Minute), second.Delta.From, "the delta spans the real observation times")
		require.Equal(t, counterBase.Add(counterWindowLen+9*time.Minute), second.Delta.To)

		third := results[2]
		require.Equal(t, CounterOK, third.Status)
		require.Equal(t, "12000.7", digitsOf(t, third.Raw.Value))
		require.Equal(t, "0.3", digitsOf(t, third.Delta.Value))
		require.Equal(t, "s2", third.Delta.BaselineSampleID)
		require.Equal(t, "s3", third.Delta.EndSampleID)
	})

	t.Run("integer 9007199254740992 to 9007199254740993", func(t *testing.T) {
		for name, readings := range map[string][2]measurement.ExactValue{
			"int64":  {measurement.NewInt64(9007199254740992), measurement.NewInt64(9007199254740993)},
			"uint64": {measurement.NewUint64(9007199254740992), measurement.NewUint64(9007199254740993)},
		} {
			t.Run(name, func(t *testing.T) {
				results, checkpoint := closeSeries(t, counterBindingFixture(),
					[]CounterReading{counterReading("a", 0, time.Minute, readings[0])},
					[]CounterReading{counterReading("b", 1, time.Minute, readings[1])},
				)
				require.Equal(t, "9007199254740993", digitsOf(t, results[1].Raw.Value))
				require.Equal(t, "1", digitsOf(t, results[1].Delta.Value))
				require.Equal(t, DeltaBasisExact, results[1].Delta.Basis)

				encoded, err := EncodeCounterCheckpoint(*checkpoint)
				require.NoError(t, err)
				require.Contains(t, string(encoded), `"9007199254740993"`, "the digits travel as a string")
				require.NotContains(t, string(encoded), "9.007", "no float rendering anywhere")
				decoded, err := DecodeCounterCheckpoint(encoded)
				require.NoError(t, err)
				require.True(t, checkpoint.Reading.Value.Equal(decoded.Reading.Value))
				require.Equal(t, readings[1].Type(), decoded.Reading.Value.Type())
			})
		}
	})

	t.Run("large values never overflow", func(t *testing.T) {
		results, _ := closeSeries(t, counterBindingFixture(),
			[]CounterReading{counterReading("a", 0, time.Minute, measurement.NewUint64(math.MaxUint64-1))},
			[]CounterReading{counterReading("b", 1, time.Minute, measurement.NewUint64(math.MaxUint64))},
		)
		require.Equal(t, "1", digitsOf(t, results[1].Delta.Value))

		results, _ = closeSeries(t, counterBindingFixture(),
			[]CounterReading{counterReading("a", 0, time.Minute, measurement.NewInt64(math.MaxInt64-5))},
			[]CounterReading{counterReading("b", 1, time.Minute, measurement.NewInt64(math.MaxInt64))},
		)
		require.Equal(t, "5", digitsOf(t, results[1].Delta.Value))
	})

	t.Run("float source keeps its float basis", func(t *testing.T) {
		low, _ := measurement.NewFloat64(12000.4)
		high, _ := measurement.NewFloat64(12000.7)
		results, _ := closeSeries(t, counterBindingFixture(),
			[]CounterReading{counterReading("f1", 0, time.Minute, low)},
			[]CounterReading{counterReading("f2", 1, time.Minute, high)},
		)
		require.Equal(t, measurement.ExactFloat64, results[1].Raw.Value.Type())
		require.Equal(t, DeltaBasisFloat64, results[1].Delta.Basis, "a float delta is never labeled exact")
		require.Equal(t, measurement.ExactFloat64, results[1].Delta.Value.Type())
		require.InDelta(t, 0.3, results[1].Delta.Value.Value(), 1e-9)
		require.NotEqual(t, 0.3, results[1].Delta.Value.Value(),
			"the float result is the IEEE difference, not a reconstructed decimal")
	})

	t.Run("a window of readings is judged on the whole sequence", func(t *testing.T) {
		results, checkpoint := closeSeries(t, counterBindingFixture(),
			[]CounterReading{counterReading("w0", 0, time.Minute, measurement.NewUint64(90))},
			[]CounterReading{
				counterReading("w1-late", 1, 8*time.Minute, measurement.NewUint64(110)),
				counterReading("w1-early", 1, 2*time.Minute, measurement.NewUint64(100)),
				counterReading("w1-mid", 1, 5*time.Minute, measurement.NewUint64(105)),
			},
		)
		second := results[1]
		require.Equal(t, CounterOK, second.Status)
		require.Equal(t, "20", digitsOf(t, second.Delta.Value), "baseline 90 to the newest reading 110")
		require.Equal(t, "w0", second.Delta.BaselineSampleID)
		require.Equal(t, "w1-late", second.Delta.EndSampleID)
		require.Equal(t, "110", digitsOf(t, second.Raw.Value))
		require.Equal(t, "w1-late", checkpoint.Reading.SampleID, "the newest reading becomes the next baseline")
	})

	t.Run("no movement is a real zero but a missing window is not", func(t *testing.T) {
		results, _ := closeSeries(t, counterBindingFixture(),
			[]CounterReading{counterReading("a", 0, time.Minute, measurement.NewUint64(123))},
			[]CounterReading{counterReading("b", 1, time.Minute, measurement.NewUint64(123))},
			nil,
		)
		require.Equal(t, CounterOK, results[1].Status)
		require.Equal(t, "0", digitsOf(t, results[1].Delta.Value))
		require.Equal(t, CounterMissing, results[2].Status)
		require.Nil(t, results[2].Raw, "an absent reading has no raw value")
		require.Nil(t, results[2].Delta, "an absent reading is NULL, not zero")
	})

	t.Run("arrival order does not change the result", func(t *testing.T) {
		binding := counterBindingFixture()
		in := []CounterReading{
			counterReading("x", 1, time.Minute, measurement.NewUint64(10)),
			counterReading("y", 1, 2*time.Minute, measurement.NewUint64(20)),
			counterReading("z", 1, 3*time.Minute, measurement.NewUint64(30)),
		}
		base := []CounterReading{counterReading("base", 0, time.Minute, measurement.NewUint64(5))}
		forward, _ := closeSeries(t, binding, base, in)
		reversed, _ := closeSeries(t, binding, base, []CounterReading{in[2], in[1], in[0]})
		first, err := json.Marshal(forward[1].Delta)
		require.NoError(t, err)
		second, err := json.Marshal(reversed[1].Delta)
		require.NoError(t, err)
		require.JSONEq(t, string(first), string(second))
		require.Equal(t, "25", digitsOf(t, forward[1].Delta.Value))
	})

	t.Run("a folded window cannot be folded twice", func(t *testing.T) {
		binding := counterBindingFixture()
		_, checkpoint := closeSeries(t, binding,
			[]CounterReading{counterReading("a", 0, time.Minute, measurement.NewUint64(1))},
			[]CounterReading{counterReading("b", 1, time.Minute, measurement.NewUint64(2))},
		)
		_, _, err := CloseCounterWindow(checkpoint, binding, counterWindow(1),
			[]CounterReading{counterReading("b", 1, time.Minute, measurement.NewUint64(2))})
		require.ErrorIs(t, err, ErrCounterWindowOrder, "a replayed window must not consume its delta again")
	})
}
