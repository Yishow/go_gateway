package decoding

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"math"
	"testing"
)

func TestScalarAllNumericTypesAndOrders(t *testing.T) {
	for _, tc := range []struct {
		dtype schema.DataType
		bits  uint64
		want  any
		words int
	}{
		{schema.DataTypeInt32, 0x81234567, int32(-2128394905), 2},
		{schema.DataTypeUint32, 0x01020304, uint32(0x01020304), 2},
		{schema.DataTypeFloat32, 0xc1200000, float32(-10), 2},
		{schema.DataTypeInt64, 0x8123456789abcdef, int64(-9141386507638288913), 4},
		{schema.DataTypeUint64, 0x0102030405060708, uint64(0x0102030405060708), 4},
		{schema.DataTypeFloat64, 0xc024000000000000, float64(-10), 4},
	} {
		for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
			t.Run(string(tc.dtype)+format, func(t *testing.T) {
				// Encode independent byte fixture, then feed network-order words.
				raw := make([]byte, tc.words*2)
				for i := range raw {
					raw[len(raw)-1-i] = byte(tc.bits >> uint(i*8))
				}
				if format == "CDAB" || format == "DCBA" {
					for i, j := 0, len(raw)-2; i < j; i, j = i+2, j-2 {
						raw[i], raw[j] = raw[j], raw[i]
						raw[i+1], raw[j+1] = raw[j+1], raw[i+1]
					}
				}
				if format == "BADC" || format == "DCBA" {
					for i := 0; i < len(raw); i += 2 {
						raw[i], raw[i+1] = raw[i+1], raw[i]
					}
				}
				words := make([]uint16, tc.words)
				for i := range words {
					words[i] = uint16(raw[2*i])<<8 | uint16(raw[2*i+1])
				}
				before := append([]uint16(nil), words...)
				got, err := DecodeScalarWords(words, tc.dtype, format)
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
				require.Equal(t, before, words)
				for n := 0; n < tc.words; n++ {
					got, err = DecodeScalarWords(words[:n], tc.dtype, format)
					require.ErrorIs(t, err, ErrShortScalar)
					require.Nil(t, got)
				}
				for _, width := range []int{16, 32} {
					components := make([]int, 0, len(words))
					for i := 0; i < len(words); i += width / 16 {
						v := uint32(words[i])
						if width == 32 {
							v = v<<16 | uint32(words[i+1])
						}
						components = append(components, int(v))
					}
					normalized, err := NormalizeFatekComponents(components, width, tc.dtype)
					require.NoError(t, err)
					require.Equal(t, words, normalized)
					got, err = DecodeScalarWords(normalized, tc.dtype, format)
					require.NoError(t, err)
					require.Equal(t, tc.want, got)
					_, err = NormalizeFatekComponents(components[:len(components)-1], width, tc.dtype)
					require.ErrorIs(t, err, ErrShortScalar)
				}
			})
		}
	}
}

func TestScalarFailureAndComponentBitPatterns(t *testing.T) {
	for _, dtype := range []schema.DataType{schema.DataTypeBool, schema.DataTypeUint16, schema.DataTypeInt16, schema.DataTypeString, "secret://bad"} {
		got, err := DecodeScalarWords([]uint16{1, 2, 3, 4}, dtype, "ABCD")
		require.ErrorIs(t, err, ErrInapplicableFormat)
		require.Nil(t, got)
	}
	got, err := DecodeScalarWords([]uint16{1, 2}, schema.DataTypeUint32, "secret://bad")
	require.ErrorIs(t, err, ErrInvalidIntent)
	require.Nil(t, got)
	for _, width := range []int{0, 1, 8, 64} {
		got, err := NormalizeFatekComponents([]int{1, 2}, width, schema.DataTypeUint32)
		require.ErrorIs(t, err, ErrInapplicableFormat)
		require.Nil(t, got)
	}
	words, err := NormalizeFatekComponents([]int{-1}, 32, schema.DataTypeInt32)
	require.NoError(t, err)
	require.Equal(t, []uint16{math.MaxUint16, math.MaxUint16}, words)
}
