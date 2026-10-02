package measurement

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func mustDecimal(t *testing.T, digits string) ExactValue {
	t.Helper()
	v, err := NewDecimal(digits)
	require.NoError(t, err)
	return v
}

func exactMixedFixture(t *testing.T) []ExactValue {
	t.Helper()
	return []ExactValue{
		NewBool(true),
		NewText("batch-001"),
		NewInt64(math.MinInt64),
		NewUint64(9007199254740993),
		NewUint64(math.MaxUint64),
		mustDecimal(t, "1234567890.123456789012345678"),
	}
}

func TestExactMixedValueRoundTripJSON(t *testing.T) {
	for _, original := range exactMixedFixture(t) {
		raw, err := json.Marshal(original)
		require.NoError(t, err)
		var decoded ExactValue
		require.NoError(t, json.Unmarshal(raw, &decoded), string(raw))
		require.Equal(t, original.Type(), decoded.Type(), string(raw))
		require.True(t, original.Equal(decoded), string(raw))
	}
}

func TestExactValueWireShapeKeepsDigitsAsStrings(t *testing.T) {
	raw, err := json.Marshal(NewUint64(9007199254740993))
	require.NoError(t, err)
	require.JSONEq(t, `{"type":"uint64","encoding":"decimal-string","value":"9007199254740993"}`, string(raw))

	raw, err = json.Marshal(NewBool(true))
	require.NoError(t, err)
	require.JSONEq(t, `{"type":"bool","encoding":"json","value":true}`, string(raw))
}

func TestExactValueRejectsNonFiniteAndMalformed(t *testing.T) {
	_, err := NewFloat64(math.NaN())
	require.ErrorIs(t, err, ErrExactValueInvalid)
	_, err = NewFloat64(math.Inf(-1))
	require.ErrorIs(t, err, ErrExactValueInvalid)
	for _, bad := range []string{"", "NaN", "Inf", "1e5", "1.2.3", "--1", " 1", "0x10", "1,5", "."} {
		_, err := NewDecimal(bad)
		require.ErrorIs(t, err, ErrExactValueInvalid, bad)
	}

	for name, raw := range map[string]string{
		"unknown type":        `{"type":"money","encoding":"json","value":1}`,
		"unknown encoding":    `{"type":"int64","encoding":"hex","value":"1"}`,
		"number for uint64":   `{"type":"uint64","encoding":"decimal-string","value":9007199254740993}`,
		"int64 overflow":      `{"type":"int64","encoding":"decimal-string","value":"9223372036854775808"}`,
		"uint64 negative":     `{"type":"uint64","encoding":"decimal-string","value":"-1"}`,
		"uint64 overflow":     `{"type":"uint64","encoding":"decimal-string","value":"18446744073709551616"}`,
		"bool as string":      `{"type":"bool","encoding":"json","value":"true"}`,
		"text as number":      `{"type":"text","encoding":"json","value":1}`,
		"decimal malformed":   `{"type":"decimal","encoding":"decimal-string","value":"1e3"}`,
		"missing encoding":    `{"type":"int64","value":"1"}`,
		"null value":          `{"type":"text","encoding":"json","value":null}`,
		"unknown field":       `{"type":"text","encoding":"json","value":"a","extra":1}`,
		"float64 nan as text": `{"type":"float64","encoding":"json","value":"NaN"}`,
		"not an object":       `[1]`,
		"trailing garbage":    `{"type":"text","encoding":"json","value":"a"} 1`,
	} {
		var v ExactValue
		require.Error(t, json.Unmarshal([]byte(raw), &v), name)
	}
}

func TestExactValueZeroValueDoesNotMarshalAsGood(t *testing.T) {
	_, err := json.Marshal(ExactValue{})
	require.ErrorIs(t, err, ErrExactValueInvalid)
}

func TestExactDecimalEqualityIgnoresTrailingZerosOnly(t *testing.T) {
	require.True(t, mustDecimal(t, "1.50").Equal(mustDecimal(t, "1.5")))
	require.True(t, mustDecimal(t, "+007.0").Equal(mustDecimal(t, "7")))
	require.True(t, mustDecimal(t, "-0").Equal(mustDecimal(t, "0.00")))
	require.False(t, mustDecimal(t, "1.5").Equal(mustDecimal(t, "1.51")))
	require.False(t, NewUint64(1).Equal(NewInt64(1)))
}
