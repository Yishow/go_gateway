package mapping

import (
	"encoding/json"
	"go-gateway/internal/datalink/schema"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExactCastPreservesIntegers(t *testing.T) {
	for _, tc := range []struct {
		value  any
		target string
		want   any
	}{
		{uint64(9007199254740993), "uint64", uint64(9007199254740993)},
		{"18446744073709551615", "uint64", uint64(math.MaxUint64)},
		{json.Number("-9223372036854775808"), "int64", int64(math.MinInt64)},
		{uint64(32767), "int16", int16(32767)},
		{float64(42), "uint32", uint32(42)},
		{uint64(9007199254740992), "float64", float64(9007199254740992)},
	} {
		got, err := executeCast(tc.value, map[string]any{"target_type": tc.target})
		require.NoError(t, err)
		require.Equal(t, tc.want, got)
	}
}

func TestExactCastRejectsInvalidValues(t *testing.T) {
	for _, tc := range []struct {
		name   string
		value  any
		target string
	}{
		{"float64Underflow", "1e-400", "float64"}, {"float32Underflow", float64(1e-50), "float32"},
		{"invalid", "bad", "int64"}, {"negative", int64(-1), "uint64"},
		{"overflow16", 65536, "uint16"}, {"overflow32", uint64(4294967296), "uint32"},
		{"signedOverflow", uint64(math.MaxUint64), "int64"},
		{"floatOverflow", math.Exp2(64), "uint64"},
		{"fraction", 1.5, "int32"}, {"nan", math.NaN(), "float64"},
		{"inf", math.Inf(1), "float32"}, {"float32Overflow", math.MaxFloat64, "float32"},
		{"precision64", uint64(9007199254740993), "float64"},
		{"precision32", int64(16777217), "float32"}, {"invalidBool", "bad", "bool"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			steps, err := json.Marshal([]schema.TransformStep{{Type: schema.TransformCast, Params: map[string]any{"target_type": tc.target}}})
			require.NoError(t, err)
			ctx, err := ExecutePipeline(tc.value, string(steps))
			require.Error(t, err)
			require.Error(t, ctx.Error)
			require.Nil(t, ctx.StepResults[0].Output)
		})
	}
}
