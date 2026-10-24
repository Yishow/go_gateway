package measurement

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func mustFloat(t *testing.T, v float64) ExactValue {
	t.Helper()
	value, err := NewFloat64(v)
	require.NoError(t, err)
	return value
}

// TestExactSubtractIsExact pins the arithmetic the cumulative path relies on:
// integers and decimals subtract without ever becoming a float64.
func TestExactSubtractIsExact(t *testing.T) {
	tests := []struct {
		name     string
		to, from ExactValue
		want     string
	}{
		{"int64 above 2^53", NewInt64(9007199254740993), NewInt64(9007199254740992), "1"},
		{"uint64 above 2^53", NewUint64(9007199254740993), NewUint64(9007199254740992), "1"},
		{"uint64 max from max-1", NewUint64(math.MaxUint64), NewUint64(math.MaxUint64 - 1), "1"},
		{"uint64 max from zero", NewUint64(math.MaxUint64), NewUint64(0), "18446744073709551615"},
		{"uint64 below zero", NewUint64(0), NewUint64(math.MaxUint64), "-18446744073709551615"},
		{"int64 full span does not overflow", NewInt64(math.MaxInt64), NewInt64(math.MinInt64), "18446744073709551615"},
		{"uint64 minus negative int64", NewUint64(math.MaxUint64), NewInt64(-1), "18446744073709551616"},
		{"decimal first step", mustDecimal(t, "12000.4"), mustDecimal(t, "12000.0"), "0.4"},
		{"decimal second step", mustDecimal(t, "12000.7"), mustDecimal(t, "12000.4"), "0.3"},
		{"decimal scales are aligned", mustDecimal(t, "12000.45"), mustDecimal(t, "12000.4"), "0.05"},
		{"decimal keeps the finer scale", mustDecimal(t, "0.30"), mustDecimal(t, "0.1"), "0.20"},
		{"decimal sign crossing", mustDecimal(t, "-0.5"), mustDecimal(t, "0.5"), "-1.0"},
		{"equal decimals are zero without a sign", mustDecimal(t, "1.0"), mustDecimal(t, "1.0"), "0.0"},
		{"integer against padded decimal", NewInt64(5), mustDecimal(t, "5.00"), "0.00"},
		{"decimal beyond float64 digits", mustDecimal(t, "123456789012345678901234567890.000000000000000002"), mustDecimal(t, "123456789012345678901234567890.000000000000000001"), "0.000000000000000001"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Subtract(tt.to, tt.from)
			require.NoError(t, err)
			require.Equal(t, ExactDecimal, got.Type(), "an exact difference stays a decimal")
			digits, ok := got.DecimalDigits()
			require.True(t, ok)
			require.Equal(t, tt.want, digits)
		})
	}
}

func TestExactSubtractFloatStaysFloat(t *testing.T) {
	to, from := 12000.7, 12000.4 // variables: a constant expression would be evaluated exactly
	got, err := Subtract(mustFloat(t, to), mustFloat(t, from))
	require.NoError(t, err)
	require.Equal(t, ExactFloat64, got.Type(), "a float source never yields an exact decimal")
	require.Equal(t, to-from, got.Value(), "one IEEE subtraction, no claimed extra precision")

	text, err := FloatDecimal(to - from)
	require.NoError(t, err)
	digits, _ := text.DecimalDigits()
	require.Equal(t, "0.3000000000010914", digits, "the decimal rendering shows the float's own result")
	_, err = FloatDecimal(math.Inf(1))
	require.ErrorIs(t, err, ErrExactValueInvalid)
}

func TestExactSubtractRefusesWhatItCannotProve(t *testing.T) {
	huge := mustFloat(t, math.MaxFloat64)
	negative := mustFloat(t, -math.MaxFloat64)
	for name, pair := range map[string][2]ExactValue{
		"float against integer": {mustFloat(t, 1), NewInt64(1)},
		"float against decimal": {mustFloat(t, 1), mustDecimal(t, "1.0")},
		"float overflow":        {huge, negative},
		"bool":                  {NewBool(true), NewBool(false)},
		"text":                  {NewText("2"), NewText("1")},
		"unset":                 {{}, NewInt64(1)},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := Subtract(pair[0], pair[1])
			require.ErrorIs(t, err, ErrExactValueInvalid)
		})
	}
}

func TestExactValueSign(t *testing.T) {
	for name, tt := range map[string]struct {
		value ExactValue
		want  int
	}{
		"positive int": {NewInt64(7), 1}, "negative int": {NewInt64(-7), -1}, "zero uint": {NewUint64(0), 0},
		"large uint": {NewUint64(math.MaxUint64), 1}, "padded zero": {mustDecimal(t, "-0.00"), 0},
		"small negative": {mustDecimal(t, "-0.000001"), -1}, "float": {mustFloat(t, -2.5), -1},
	} {
		t.Run(name, func(t *testing.T) {
			got, ok := tt.value.Sign()
			require.True(t, ok)
			require.Equal(t, tt.want, got)
		})
	}
	for _, value := range []ExactValue{NewBool(true), NewText("1"), {}} {
		_, ok := value.Sign()
		require.False(t, ok, "only numeric values have a sign")
	}
}

func TestExactAddIsExact(t *testing.T) {
	for _, tt := range []struct {
		name string
		a, b ExactValue
		want string
	}{
		{"above 2^53", NewUint64(9007199254740992), NewUint64(1), "9007199254740993"},
		{"uint64 max plus one does not wrap", NewUint64(math.MaxUint64), NewUint64(1), "18446744073709551616"},
		{"decimal scales aligned", mustDecimal(t, "0.30"), mustDecimal(t, "0.1"), "0.40"},
		{"negative and positive", NewInt64(-3), NewUint64(10), "7"},
		{"rollover compensation", NewInt64(-65526), NewUint64(65536), "10"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(tt.a, tt.b)
			require.NoError(t, err)
			digits, _ := got.DecimalDigits()
			require.Equal(t, tt.want, digits)
		})
	}
	_, err := Add(mustFloat(t, 1), NewInt64(1))
	require.ErrorIs(t, err, ErrExactValueInvalid)
}

func TestExactRatAndCanonicalDigits(t *testing.T) {
	rat, ok := mustDecimal(t, "0.250").Rat()
	require.True(t, ok)
	require.Equal(t, "1/4", rat.RatString())
	rat, ok = NewUint64(math.MaxUint64).Rat()
	require.True(t, ok)
	require.Equal(t, "18446744073709551615", rat.RatString())
	_, ok = mustFloat(t, 0.25).Rat()
	require.False(t, ok, "a float is not exact")
	_, ok = NewText("1").Rat()
	require.False(t, ok)

	canonical := func(v ExactValue) string {
		digits, ok := v.CanonicalDigits()
		require.True(t, ok)
		return digits
	}
	require.Equal(t, canonical(mustDecimal(t, "65536.000")), canonical(NewUint64(65536)), "equal numbers share one canonical form")
	require.Equal(t, "0.5", canonical(mustDecimal(t, "00.500")))
	require.Equal(t, "-1.25", canonical(mustDecimal(t, "-1.2500")))
	_, ok = mustFloat(t, 1).CanonicalDigits()
	require.False(t, ok)
}
