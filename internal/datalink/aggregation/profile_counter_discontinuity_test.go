package aggregation

import (
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

func uintValue(v uint64) measurement.ExactValue { return measurement.NewUint64(v) }

func exactPtr(t *testing.T, value measurement.ExactValue) *measurement.ExactValue {
	t.Helper()
	return &value
}

// rolloverBinding is a monotonic counter with confirmed rollover evidence.
func rolloverBinding(t *testing.T, modulus, maxRate string) CounterBinding {
	t.Helper()
	binding := counterBindingFixture()
	if modulus != "" {
		binding.RolloverModulus = exactPtr(t, decimalOf(t, modulus))
	}
	if maxRate != "" {
		binding.MaxRatePerSecond = exactPtr(t, decimalOf(t, maxRate))
	}
	require.NoError(t, binding.Validate())
	return binding
}

// foldTwo closes a baseline window and then a window with the readings, and
// returns the second result.
func foldTwo(t *testing.T, binding CounterBinding, first measurement.ExactValue, second ...CounterReading) CounterWindowResult {
	t.Helper()
	results, _ := closeSeries(t, binding, []CounterReading{counterReading("base", 0, 9*time.Minute, first)}, second)
	return results[1]
}

// TestProfileCounterResetWithinWindow pins that the whole pre-decimation
// sequence is judged: a reset in the middle of a window is evidence, not
// something the newest value can paper over.
func TestProfileCounterResetWithinWindow(t *testing.T) {
	binding := counterBindingFixture()

	t.Run("100 then 10 then 110 keeps raw 110 and no delta", func(t *testing.T) {
		results, checkpoint := closeSeries(t, binding,
			[]CounterReading{counterReading("base", 0, 9*time.Minute, uintValue(90))},
			[]CounterReading{
				counterReading("r3", 1, 8*time.Minute, uintValue(110)),
				counterReading("r1", 1, 2*time.Minute, uintValue(100)),
				counterReading("r2", 1, 5*time.Minute, uintValue(10)),
			},
			[]CounterReading{counterReading("next", 2, 5*time.Minute, uintValue(125))},
		)
		window := results[1]
		require.Equal(t, CounterUncertain, window.Status)
		require.Equal(t, ReasonDecreaseUnproven, window.Reason)
		require.Equal(t, "110", digitsOf(t, window.Raw.Value))
		require.Equal(t, "r3", window.Raw.SampleID)
		require.Nil(t, window.Delta, "endpoints alone would have shown 20; the reset makes it unknown")

		require.Len(t, window.Events, 1)
		event := window.Events[0]
		require.Equal(t, EventDecrease, event.Kind)
		require.Equal(t, ReasonDecreaseUnproven, event.Reason)
		require.Equal(t, "r1", event.FromSampleID)
		require.Equal(t, "r2", event.ToSampleID)
		require.Equal(t, "100", digitsOf(t, event.FromValue))
		require.Equal(t, "10", digitsOf(t, event.ToValue))
		require.Equal(t, counterWindow(1).Start.Add(2*time.Minute), event.From)
		require.Equal(t, counterWindow(1).Start.Add(5*time.Minute), event.To)
		require.Zero(t, window.EventsOmitted)

		// The last reliable reading is the next baseline; the uncertain change is not counted.
		next := results[2]
		require.Equal(t, CounterOK, next.Status)
		require.Equal(t, "15", digitsOf(t, next.Delta.Value))
		require.Equal(t, "r3", next.Delta.BaselineSampleID)
		require.Equal(t, "next", checkpoint.Reading.SampleID)
	})

	t.Run("every discontinuity is an event and the list is bounded", func(t *testing.T) {
		readings := make([]CounterReading, 0, 40)
		for i := range 40 {
			value := uint64(100)
			if i%2 == 1 {
				value = 10
			}
			readings = append(readings, counterReading(fmt.Sprintf("s%02d", i), 1, time.Duration(i)*10*time.Second, uintValue(value)))
		}
		window := foldTwo(t, binding, uintValue(5), readings...)
		require.Equal(t, CounterUncertain, window.Status)
		require.Len(t, window.Events, maxCounterEvents)
		require.Equal(t, 20-maxCounterEvents, window.EventsOmitted, "20 drops happened; the rest are counted, not stored")
		require.Equal(t, "s01", window.Events[0].ToSampleID, "events keep time order")
	})

	t.Run("contradictory readings at one instant never make a delta", func(t *testing.T) {
		window := foldTwo(t, binding, uintValue(90),
			counterReading("a", 1, 2*time.Minute, uintValue(60)),
			counterReading("b", 1, 2*time.Minute, uintValue(61)),
			counterReading("c", 1, 8*time.Minute, uintValue(70)),
		)
		require.Equal(t, CounterUncertain, window.Status)
		require.Equal(t, ReasonDecreaseUnproven, window.Reason, "90 to 60 is also an unexplained drop and comes first")
		require.Nil(t, window.Delta)
		kinds := make([]CounterEventKind, 0, len(window.Events))
		for _, event := range window.Events {
			kinds = append(kinds, event.Kind)
		}
		require.Equal(t, []CounterEventKind{EventDecrease, EventContradictory}, kinds)
	})

	t.Run("a contradiction at the newest instant is not trusted as a baseline", func(t *testing.T) {
		results, checkpoint := closeSeries(t, binding,
			[]CounterReading{counterReading("base", 0, 9*time.Minute, uintValue(90))},
			[]CounterReading{
				counterReading("a", 1, 8*time.Minute, uintValue(95)),
				counterReading("b", 1, 8*time.Minute, uintValue(96)),
			},
			[]CounterReading{counterReading("next", 2, 5*time.Minute, uintValue(100))},
		)
		require.Equal(t, CounterUncertain, results[1].Status)
		require.Equal(t, ReasonContradictory, results[1].Reason)
		require.Nil(t, results[1].Delta)
		require.Equal(t, CounterGap, results[2].Status, "the unreliable instant is skipped, so the next total spans it")
		require.True(t, results[2].Delta.Span)
		require.Equal(t, "10", digitsOf(t, results[2].Delta.Value), "100 against the last reliable 90")
		require.Equal(t, "base", results[2].Delta.BaselineSampleID)
		require.Equal(t, "next", checkpoint.Reading.SampleID)
	})

	t.Run("a first window with a contradictory newest instant creates no baseline", func(t *testing.T) {
		result, checkpoint, err := CloseCounterWindow(nil, binding, counterWindow(0), []CounterReading{
			counterReading("a", 0, 3*time.Minute, uintValue(5)), counterReading("b", 0, 3*time.Minute, uintValue(6)),
		})
		require.NoError(t, err)
		require.Equal(t, CounterUncertain, result.Status)
		require.Equal(t, ReasonContradictory, result.Reason)
		require.Nil(t, result.Delta)
		require.Nil(t, checkpoint, "a baseline is never taken from an instant that disagrees with itself")
	})

	t.Run("the same value reported twice at one instant is not a contradiction", func(t *testing.T) {
		window := foldTwo(t, binding, uintValue(90),
			counterReading("a", 1, 4*time.Minute, uintValue(95)), counterReading("b", 1, 4*time.Minute, uintValue(95)))
		require.Equal(t, CounterOK, window.Status)
		require.Equal(t, "5", digitsOf(t, window.Delta.Value))
	})
}

// TestProfileSignedAndRollover pins that a negative step is accepted only for a
// confirmed signed counter and a wrap only with single-wrap evidence.
func TestProfileSignedAndRollover(t *testing.T) {
	t.Run("a confirmed signed counter allows 10 to 7 as -3", func(t *testing.T) {
		binding := counterBindingFixture()
		binding.Kind = CounterKindSigned
		window := foldTwo(t, binding, uintValue(10), counterReading("s", 1, 9*time.Minute, uintValue(7)))
		require.Equal(t, CounterOK, window.Status)
		require.Equal(t, "-3", digitsOf(t, window.Delta.Value))
		require.Equal(t, DeltaBasisExact, window.Delta.Basis)
		require.Empty(t, window.Events)

		decimal := foldTwo(t, binding, decimalOf(t, "12000.7"), counterReading("s", 1, 9*time.Minute, decimalOf(t, "12000.4")))
		require.Equal(t, "-0.3", digitsOf(t, decimal.Delta.Value))
	})

	t.Run("without the confirmation 10 to 7 is never guessed", func(t *testing.T) {
		window := foldTwo(t, counterBindingFixture(), uintValue(10), counterReading("s", 1, 9*time.Minute, uintValue(7)))
		require.Equal(t, CounterUncertain, window.Status)
		require.Nil(t, window.Delta, "not -3, not 3 and not zero")
	})

	t.Run("65530 to 4 with modulus 65536 is 10 only with the single-wrap bound", func(t *testing.T) {
		window := foldTwo(t, rolloverBinding(t, "65536", "100"), uintValue(65530), counterReading("s", 1, 9*time.Minute, uintValue(4)))
		require.Equal(t, CounterOK, window.Status)
		require.Equal(t, "10", digitsOf(t, window.Delta.Value))
		require.Equal(t, DeltaBasisExact, window.Delta.Basis)
		require.Equal(t, "base", window.Delta.BaselineSampleID)
		require.Len(t, window.Events, 1)
		require.Equal(t, EventRollover, window.Events[0].Kind)
		require.Empty(t, window.Events[0].Reason)
	})

	t.Run("evidence that does not exclude other wraps is uncertain", func(t *testing.T) {
		for name, tt := range map[string]struct {
			modulus, rate string
			from, to      uint64
			reason        string
		}{
			"modulus alone":              {"65536", "", 65530, 4, ReasonRolloverEvidence},
			"rate allows several wraps":  {"65536", "200", 65530, 4, ReasonRolloverMultiWrap},
			"wrap larger than the bound": {"65536", "50", 65530, 40000, ReasonRolloverRate},
			"reading outside the range":  {"65536", "100", 70000, 4, ReasonRolloverRange},
			"no rollover policy":         {"", "", 65530, 4, ReasonDecreaseUnproven},
		} {
			t.Run(name, func(t *testing.T) {
				window := foldTwo(t, rolloverBinding(t, tt.modulus, tt.rate), uintValue(tt.from), counterReading("s", 1, 9*time.Minute, uintValue(tt.to)))
				require.Equal(t, CounterUncertain, window.Status)
				require.Equal(t, tt.reason, window.Reason)
				require.Nil(t, window.Delta, "no abs() and no zero to hide a drop")
				require.Equal(t, EventDecrease, window.Events[0].Kind)
				require.Equal(t, tt.reason, window.Events[0].Reason)
			})
		}
	})

	t.Run("two proven wraps in one window are both compensated exactly", func(t *testing.T) {
		binding := rolloverBinding(t, "1000", "5")
		window := foldTwo(t, binding, uintValue(900),
			counterReading("r1", 1, 30*time.Second, uintValue(990)),
			counterReading("r2", 1, 90*time.Second, uintValue(40)),
			counterReading("r3", 1, 190*time.Second, uintValue(540)),
			counterReading("r4", 1, 280*time.Second, uintValue(990)),
			counterReading("r5", 1, 340*time.Second, uintValue(30)),
		)
		require.Equal(t, CounterOK, window.Status)
		require.Equal(t, "1130", digitsOf(t, window.Delta.Value), "90 + 50 + 500 + 450 + 40")
		require.Len(t, window.Events, 2)
	})

	t.Run("rollover arithmetic never goes through float64", func(t *testing.T) {
		binding := rolloverBinding(t, "18446744073709551616", "1000000000000000")
		window := foldTwo(t, binding, uintValue(18446744073709551610), counterReading("s", 1, 9*time.Minute, uintValue(4)))
		require.Equal(t, CounterOK, window.Status)
		require.Equal(t, "10", digitsOf(t, window.Delta.Value), "a float64 cannot tell these two values apart")

		decimal := rolloverBinding(t, "10000.0", "10")
		decimalWindow := foldTwo(t, decimal, decimalOf(t, "9999.8"), counterReading("s", 1, 9*time.Minute, decimalOf(t, "0.3")))
		require.Equal(t, "0.5", digitsOf(t, decimalWindow.Delta.Value))
	})

	t.Run("a float source cannot prove a wrap", func(t *testing.T) {
		low, _ := measurement.NewFloat64(65530)
		high, _ := measurement.NewFloat64(4)
		window := foldTwo(t, rolloverBinding(t, "65536", "100"), low, counterReading("s", 1, 9*time.Minute, high))
		require.Equal(t, CounterUncertain, window.Status)
		require.Equal(t, ReasonRolloverEvidence, window.Reason)
		require.Nil(t, window.Delta)
	})

	t.Run("the policy must be exact and belong to a monotonic counter", func(t *testing.T) {
		float, _ := measurement.NewFloat64(65536)
		for name, mutate := range map[string]func(*CounterBinding){
			"float modulus":    func(b *CounterBinding) { b.RolloverModulus = &float },
			"zero modulus":     func(b *CounterBinding) { b.RolloverModulus = exactPtr(t, uintValue(0)) },
			"negative rate":    func(b *CounterBinding) { b.MaxRatePerSecond = exactPtr(t, decimalOf(t, "-1")) },
			"rate alone":       func(b *CounterBinding) { b.RolloverModulus = nil },
			"signed + modulus": func(b *CounterBinding) { b.Kind = CounterKindSigned },
		} {
			t.Run(name, func(t *testing.T) {
				binding := rolloverBinding(t, "65536", "100")
				mutate(&binding)
				require.ErrorIs(t, binding.Validate(), ErrCounterPolicyInvalid)
				_, _, err := CloseCounterWindow(nil, binding, counterWindow(0),
					[]CounterReading{counterReading("s", 0, time.Minute, uintValue(1))})
				require.Error(t, err)
			})
		}
	})
}
