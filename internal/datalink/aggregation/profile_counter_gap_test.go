package aggregation

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestProfileCounterGapAndEpoch pins missing windows, the gap span total and the
// revisions that start a new baseline.
func TestProfileCounterGapAndEpoch(t *testing.T) {
	binding := counterBindingFixture()

	t.Run("1000 at 10:00 and 1002 at 10:20 around a missing window", func(t *testing.T) {
		epoch2 := binding
		epoch2.SeriesEpoch = "epoch-2"
		var checkpoint *CounterCheckpoint
		fold := func(b CounterBinding, window int, readings ...CounterReading) CounterWindowResult {
			result, next, err := CloseCounterWindow(checkpoint, b, counterWindow(window), readings)
			require.NoError(t, err)
			checkpoint = next
			return result
		}
		first := fold(binding, 0, counterReading("m0", 0, 0, uintValue(1000)))
		require.Equal(t, CounterBaseline, first.Status)
		before := checkpoint

		missing := fold(binding, 1)
		require.Equal(t, CounterMissing, missing.Status)
		require.Nil(t, missing.Raw, "an absent reading has no raw value")
		require.Nil(t, missing.Delta, "and no delta, never zero")
		require.Same(t, before, checkpoint, "a missing window does not advance the reliable baseline")

		span := fold(binding, 2, counterReading("m2", 2, 0, uintValue(1002)))
		require.Equal(t, CounterGap, span.Status)
		require.Equal(t, "1002", digitsOf(t, span.Raw.Value))
		require.NotNil(t, span.Delta)
		require.True(t, span.Delta.Span, "labeled as a span total")
		require.Equal(t, "2", digitsOf(t, span.Delta.Value))
		require.Equal(t, time.Date(2026, 10, 24, 10, 0, 0, 0, time.UTC), span.Delta.From)
		require.Equal(t, time.Date(2026, 10, 24, 10, 20, 0, 0, time.UTC), span.Delta.To, "the real 10:00 to 10:20 interval")
		require.Equal(t, "m0", span.Delta.BaselineSampleID)
		require.Equal(t, "m2", span.Delta.EndSampleID)

		changed := fold(epoch2, 3, counterReading("m3", 3, time.Minute, uintValue(5)))
		require.Equal(t, CounterBaseline, changed.Status)
		require.Equal(t, ReasonBindingChanged, changed.Reason)
		require.Nil(t, changed.Delta, "readings of different epochs are never subtracted")
		require.Equal(t, "5", digitsOf(t, changed.Raw.Value))
	})

	t.Run("the window after a gap compares with the reading that closed it", func(t *testing.T) {
		results, _ := closeSeries(t, binding,
			[]CounterReading{counterReading("a", 0, 0, uintValue(1000))}, nil,
			[]CounterReading{counterReading("c", 2, 0, uintValue(1002))},
			[]CounterReading{counterReading("d", 3, 0, uintValue(1005))},
		)
		require.Equal(t, CounterOK, results[3].Status)
		require.False(t, results[3].Delta.Span)
		require.Equal(t, "3", digitsOf(t, results[3].Delta.Value))
		require.Equal(t, "c", results[3].Delta.BaselineSampleID)
	})

	t.Run("a span that cannot be trusted stays NULL", func(t *testing.T) {
		span := func(binding CounterBinding, first, last uint64) CounterWindowResult {
			results, _ := closeSeries(t, binding,
				[]CounterReading{counterReading("a", 0, 0, uintValue(first))},
				make([]CounterReading, 0),
				[]CounterReading{counterReading("z", 2, 0, uintValue(last))})
			return results[2]
		}
		decrease := span(binding, 1000, 900)
		require.Equal(t, CounterUncertain, decrease.Status)
		require.Nil(t, decrease.Delta)

		// 20 minutes at up to 100/s could wrap a 65536 modulus many times: not provable.
		wrapped := span(rolloverBinding(t, "65536", "100"), 65530, 4)
		require.Equal(t, CounterUncertain, wrapped.Status)
		require.Equal(t, ReasonRolloverMultiWrap, wrapped.Reason)
		require.Nil(t, wrapped.Delta)
	})

	t.Run("every revision that changes meaning starts a new baseline", func(t *testing.T) {
		for name, mutate := range map[string]func(*CounterBinding){
			"series epoch":        func(b *CounterBinding) { b.SeriesEpoch = "epoch-2" },
			"definition revision": func(b *CounterBinding) { b.DefinitionRevision = "def-2" },
			"source binding":      func(b *CounterBinding) { b.SourceBindingRevision = "bind-2" },
			"profile revision":    func(b *CounterBinding) { b.ProfileRevision = "rev-2" },
			"measurement":         func(b *CounterBinding) { b.MeasurementID = "meas-other" },
			"unit":                func(b *CounterBinding) { b.Unit = "MWh" },
			"counter kind":        func(b *CounterBinding) { b.Kind = CounterKindSigned },
			"rollover modulus": func(b *CounterBinding) {
				b.RolloverModulus = exactPtr(t, uintValue(65536))
				b.MaxRatePerSecond = exactPtr(t, uintValue(10))
			},
		} {
			t.Run(name, func(t *testing.T) {
				_, checkpoint := closeSeries(t, binding, []CounterReading{counterReading("a", 0, time.Minute, uintValue(1000))})
				next := binding
				mutate(&next)
				result, _, err := CloseCounterWindow(checkpoint, next, counterWindow(1),
					[]CounterReading{counterReading("b", 1, time.Minute, uintValue(1005))})
				require.NoError(t, err)
				require.Equal(t, CounterBaseline, result.Status)
				require.Equal(t, ReasonBindingChanged, result.Reason)
				require.Nil(t, result.Delta)
			})
		}

		first, second := rolloverBinding(t, "65536", "100"), rolloverBinding(t, "65536", "50")
		firstKey, err := first.Key()
		require.NoError(t, err)
		secondKey, err := second.Key()
		require.NoError(t, err)
		require.NotEqual(t, firstKey, secondKey, "a changed rate bound is a new baseline")
		same := rolloverBinding(t, "65536.00", "100.0")
		sameKey, err := same.Key()
		require.NoError(t, err)
		require.Equal(t, firstKey, sameKey, "an equal number written differently is the same evidence")
	})

	t.Run("a checkpoint from before the rollover policy existed is refused", func(t *testing.T) {
		_, checkpoint := closeSeries(t, binding, []CounterReading{counterReading("a", 0, time.Minute, uintValue(1))})
		encoded, err := EncodeCounterCheckpoint(*checkpoint)
		require.NoError(t, err)
		require.Contains(t, string(encoded), `"version":2`)
		_, err = DecodeCounterCheckpoint([]byte(strings.Replace(string(encoded), `"version":2`, `"version":1`, 1)))
		require.ErrorIs(t, err, ErrCounterCheckpointInvalid, "a baseline taken without the policy is never reinterpreted")
		require.ErrorIs(t, err, ErrCounterCheckpointVersion, "and the cause is named, so a caller can tell a format change from damage")
		_, err = DecodeCounterCheckpoint([]byte(`{"version":2,"binding":`))
		require.ErrorIs(t, err, ErrCounterCheckpointInvalid)
		require.NotErrorIs(t, err, ErrCounterCheckpointVersion, "damage is not a version change")

		policy := rolloverBinding(t, "65536", "100")
		_, withPolicy := closeSeries(t, policy, []CounterReading{counterReading("a", 0, time.Minute, uintValue(65530))})
		encoded, err = EncodeCounterCheckpoint(*withPolicy)
		require.NoError(t, err)
		decoded, err := DecodeCounterCheckpoint(encoded)
		require.NoError(t, err)
		require.True(t, decoded.Binding.RolloverModulus.Equal(*policy.RolloverModulus))
		require.True(t, decoded.Binding.MaxRatePerSecond.Equal(*policy.MaxRatePerSecond))
		decodedKey, _ := decoded.Binding.Key()
		policyKey, _ := policy.Key()
		require.Equal(t, policyKey, decodedKey, "the baseline identity survives serialization")
	})
}
