package measurement

import (
	"math"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestExactMixedValueRoundTripTagTypeMapping(t *testing.T) {
	for dataType, want := range map[schema.DataType]ExactType{
		schema.DataTypeBool:    ExactBool,
		schema.DataTypeInt16:   ExactInt64,
		schema.DataTypeInt32:   ExactInt64,
		schema.DataTypeInt64:   ExactInt64,
		schema.DataTypeUint16:  ExactUint64,
		schema.DataTypeUint32:  ExactUint64,
		schema.DataTypeUint64:  ExactUint64,
		schema.DataTypeFloat32: ExactFloat64,
		schema.DataTypeFloat64: ExactFloat64,
		schema.DataTypeString:  ExactText,
	} {
		got, ok := ExactTypeForTag(dataType)
		require.True(t, ok, dataType)
		require.Equal(t, want, got, dataType)
	}
	_, ok := ExactTypeForTag(schema.DataType("decimal"))
	require.False(t, ok, "no tag type implies a decimal; it is never guessed")
}

func TestExactMixedValueRoundTripFromGoKeepsExactDigits(t *testing.T) {
	cases := []struct {
		name     string
		expected ExactType
		input    any
		want     ExactValue
	}{
		{"bool", ExactBool, true, NewBool(true)},
		{"text", ExactText, "batch-001", NewText("batch-001")},
		{"int", ExactInt64, int(-5), NewInt64(-5)},
		{"int16", ExactInt64, int16(7), NewInt64(7)},
		{"int64 min", ExactInt64, int64(math.MinInt64), NewInt64(math.MinInt64)},
		{"uint16 into int64", ExactInt64, uint16(9), NewInt64(9)},
		{"uint64 beyond float", ExactUint64, uint64(9007199254740993), NewUint64(9007199254740993)},
		{"uint64 max", ExactUint64, uint64(math.MaxUint64), NewUint64(math.MaxUint64)},
		{"positive int into uint64", ExactUint64, int64(12), NewUint64(12)},
	}
	for _, c := range cases {
		got, err := ExactFromGo(c.expected, c.input)
		require.NoError(t, err, c.name)
		require.True(t, c.want.Equal(got), c.name)
	}
	float, err := ExactFromGo(ExactFloat64, float32(1.5))
	require.NoError(t, err)
	require.Equal(t, ExactFloat64, float.Type())
	float, err = ExactFromGo(ExactFloat64, 2.25)
	require.NoError(t, err)
	require.Equal(t, 2.25, float.Value())
}

func TestExactMixedValueRoundTripFromGoRejectsCoercionAndLoss(t *testing.T) {
	for name, c := range map[string]struct {
		expected ExactType
		input    any
	}{
		"nil":                 {ExactInt64, nil},
		"bool as text":        {ExactText, true},
		"string as bool":      {ExactBool, "true"},
		"float into int":      {ExactInt64, 1.0},
		"int into float":      {ExactFloat64, int64(1)},
		"string into int":     {ExactInt64, "12"},
		"uint64 overflow":     {ExactInt64, uint64(math.MaxInt64) + 1},
		"negative into uint":  {ExactUint64, int64(-1)},
		"NaN":                 {ExactFloat64, math.NaN()},
		"+Inf":                {ExactFloat64, math.Inf(1)},
		"float32 NaN":         {ExactFloat64, float32(math.NaN())},
		"decimal unsupported": {ExactDecimal, 1.5},
		"unknown expected":    {ExactType("money"), 1},
		"struct":              {ExactText, struct{}{}},
	} {
		_, err := ExactFromGo(c.expected, c.input)
		require.ErrorIs(t, err, ErrExactValueInvalid, name)
	}
}
