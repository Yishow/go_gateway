package hsllogic

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestByteTransformIntegerBoundsPreserveFallback(t *testing.T) {
	buffer := []byte{0xff, 1, 2, 3, 4, 5, 6, 7, 8}
	for _, format := range []DataFormat{DataFormatABCD, DataFormatBADC, DataFormatCDAB, DataFormatDCBA} {
		transform := NewByteTransform(format)
		for _, size := range []int{2, 4, 8} {
			read := func(index int) uint64 {
				switch size {
				case 2:
					return uint64(transform.TransformUint16(buffer, index))
				case 4:
					return uint64(transform.TransformUint32(buffer, index))
				default:
					return transform.TransformUint64(buffer, index)
				}
			}
			for _, index := range []int{-1, math.MinInt, math.MaxInt, len(buffer) - size + 1, len(buffer)} {
				t.Run(fmt.Sprintf("%s/%d/%d", format, size, index), func(t *testing.T) {
					require.NotPanics(t, func() { require.Zero(t, read(index)) })
				})
			}
		}
	}
}

func TestByteTransformIntegerBoundsPreserveByteOrder(t *testing.T) {
	for _, tc := range []struct {
		format DataFormat
		word   uint16
		dword  uint32
		qword  uint64
	}{
		{DataFormatABCD, 0x0102, 0x01020304, 0x0102030405060708},
		{DataFormatBADC, 0x0201, 0x02010403, 0x0201040306050807},
		{DataFormatCDAB, 0x0102, 0x03040102, 0x0708050603040102},
		{DataFormatDCBA, 0x0201, 0x04030201, 0x0807060504030201},
	} {
		transform := NewByteTransform(tc.format)
		require.Equal(t, tc.word, transform.TransformUint16([]byte{0xff, 1, 2}, 1))
		require.Equal(t, tc.dword, transform.TransformUint32([]byte{0xff, 1, 2, 3, 4}, 1))
		require.Equal(t, tc.qword, transform.TransformUint64([]byte{0xff, 1, 2, 3, 4, 5, 6, 7, 8}, 1))
		require.Zero(t, transform.TransformUint16(nil, 0))
		require.Zero(t, transform.TransformUint32(nil, 0))
		require.Zero(t, transform.TransformUint64(nil, 0))
	}
}
