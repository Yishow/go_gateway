package measurement

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func digitsOf(t *testing.T, value ExactValue) string {
	t.Helper()
	digits, ok := value.DecimalDigits()
	require.True(t, ok, "an exact mean must carry decimal digits")
	return digits
}

// TestMeanExactProjection pins the arithmetic mean contract: an exact sum and
// integer count, a decimal result independent of the source type, and at most
// 18 fractional digits rounded half-even.
func TestMeanExactProjection(t *testing.T) {
	t.Run("window of 10, 20 and 30", func(t *testing.T) {
		result, err := Mean([]ExactValue{NewInt64(10), NewInt64(20), NewInt64(30)})
		require.NoError(t, err)
		require.Equal(t, MeanAggregationVersion, result.Version)
		require.Equal(t, MeanAggregationMethod, result.Method)
		require.Equal(t, "20", digitsOf(t, result.Value))
		require.Equal(t, "60", digitsOf(t, result.Sum))
		require.Equal(t, int64(3), result.Count)
		require.Equal(t, MeanBasisExact, result.Basis)
		require.Equal(t, 0, result.Scale)
		require.Equal(t, MeanRoundingHalfEven, result.Rounding)
		require.False(t, result.Rounded)
	})

	t.Run("one half stays a decimal", func(t *testing.T) {
		result, err := Mean([]ExactValue{NewInt64(1), NewInt64(2)})
		require.NoError(t, err)
		require.Equal(t, "1.5", digitsOf(t, result.Value))
		require.Equal(t, int64(2), result.Count)
		require.Equal(t, 1, result.Scale)
		require.False(t, result.Rounded)
	})

	t.Run("repeating decimal keeps 18 digits", func(t *testing.T) {
		result, err := Mean([]ExactValue{NewInt64(0), NewInt64(0), NewInt64(1)})
		require.NoError(t, err)
		require.Equal(t, "0.333333333333333333", digitsOf(t, result.Value))
		require.Equal(t, "1", digitsOf(t, result.Sum))
		require.Equal(t, 18, result.Scale)
		require.True(t, result.Rounded)
	})

	t.Run("integers above 2^53 keep every digit", func(t *testing.T) {
		result, err := Mean([]ExactValue{NewInt64(9007199254740992), NewInt64(9007199254740993)})
		require.NoError(t, err)
		require.Equal(t, "9007199254740992.5", digitsOf(t, result.Value))
		require.Equal(t, "18014398509481985", digitsOf(t, result.Sum))
		require.Equal(t, MeanBasisExact, result.Basis)
		require.Equal(t, 1, result.Scale)
	})

	t.Run("uint64 and decimals stay exact", func(t *testing.T) {
		result, err := Mean([]ExactValue{NewUint64(math.MaxUint64), NewUint64(math.MaxUint64 - 1)})
		require.NoError(t, err)
		require.Equal(t, "18446744073709551614.5", digitsOf(t, result.Value))

		result, err = Mean([]ExactValue{mustDecimal(t, "0.1"), mustDecimal(t, "0.2")})
		require.NoError(t, err)
		require.Equal(t, "0.15", digitsOf(t, result.Value))
		require.Equal(t, "0.3", digitsOf(t, result.Sum))
		require.Equal(t, 2, result.Scale)
	})
}

// TestMeanHalfEvenTies pins both signs of an exact half-even tie at the 18th
// fractional digit: a tie goes to the even neighbor, never away from zero.
func TestMeanHalfEvenTies(t *testing.T) {
	tests := []struct {
		name   string
		values []ExactValue
		want   string
		scale  int
		round  bool
	}{
		{"positive tie rounds up to even", []ExactValue{mustDecimal(t, "0.000000000000000003"), NewInt64(0)}, "0.000000000000000002", 18, true},
		{"positive tie stays on even", []ExactValue{mustDecimal(t, "0.000000000000000005"), NewInt64(0)}, "0.000000000000000002", 18, true},
		{"negative tie rounds away from even", []ExactValue{mustDecimal(t, "-0.000000000000000003"), NewInt64(0)}, "-0.000000000000000002", 18, true},
		{"negative tie below one half is zero", []ExactValue{mustDecimal(t, "-0.000000000000000001"), NewInt64(0)}, "0", 0, true},
		{"positive tie below one half is zero", []ExactValue{mustDecimal(t, "0.000000000000000001"), NewInt64(0)}, "0", 0, true},
		{"exact fraction is not rounded", []ExactValue{mustDecimal(t, "0.25")}, "0.25", 2, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Mean(test.values)
			require.NoError(t, err)
			require.Equal(t, test.want, digitsOf(t, result.Value))
			require.Equal(t, test.scale, result.Scale)
			require.Equal(t, test.round, result.Rounded)
			require.Equal(t, MeanRoundingHalfEven, result.Rounding)
		})
	}
}

// TestMeanFloatBasis keeps a float window on its own basis: one float64 sum in
// the caller's order, explicitly declared, with no claimed extra precision.
func TestMeanFloatBasis(t *testing.T) {
	first, second := 0.1, 0.2
	sum := first + second // literal IEEE result, not a decimal expectation
	result, err := Mean([]ExactValue{mustFloat(t, first), mustFloat(t, second)})
	require.NoError(t, err)
	require.Equal(t, ExactFloat64, result.Value.Type(), "a float source keeps a float result")
	require.Equal(t, MeanBasisFloat64, result.Basis)
	require.Equal(t, int64(2), result.Count)
	require.Equal(t, sum, result.Sum.Value())
	require.Equal(t, sum/2, result.Value.Value())
	require.NotEqual(t, MeanRoundingHalfEven, result.Rounding, "a float mean is not a half-even decimal")

	// The answer depends on the given order; the caller supplies the stable
	// observed_at plus sample_id order.
	forward, err := Mean([]ExactValue{mustFloat(t, 0.1), mustFloat(t, 0.2), mustFloat(t, 0.3)})
	require.NoError(t, err)
	backward, err := Mean([]ExactValue{mustFloat(t, 0.3), mustFloat(t, 0.2), mustFloat(t, 0.1)})
	require.NoError(t, err)
	require.Equal(t, ExactFloat64, forward.Value.Type())
	require.NotEqual(t, forward.Sum.Value(), backward.Sum.Value(), "the sum follows the given order")
}

// TestMeanRejections keeps arithmetic and representation failures explicit:
// nothing is substituted by zero or silently reduced to the latest reading.
func TestMeanRejections(t *testing.T) {
	tests := []struct {
		name   string
		values []ExactValue
	}{
		{"no observations", nil},
		{"text is not numeric", []ExactValue{NewText("on")}},
		{"bool is not numeric", []ExactValue{NewBool(true)}},
		{"a float cannot mix with an exact value", []ExactValue{NewInt64(1), mustFloat(t, 2)}},
		{"a float sum that overflows", []ExactValue{mustFloat(t, math.MaxFloat64), mustFloat(t, math.MaxFloat64)}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Mean(test.values)
			require.ErrorIs(t, err, ErrExactValueInvalid)
		})
	}
}
