package mcprotocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractBitFromWordPreservesLowWordBits(t *testing.T) {
	for _, tc := range []struct {
		name      string
		word, bit int
		want      bool
	}{
		{"zero", 0, 0, false}, {"low", 1, 0, true},
		{"high", 0x8000, 15, true}, {"outside word", 0x10000, 15, false},
		{"negative low", -1, 0, true}, {"negative high", -32768, 15, true},
		{"negative outside word", -65536, 15, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ExtractBitFromWord(tc.word, tc.bit)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
	for _, bit := range []int{-1, 16} {
		_, err := ExtractBitFromWord(1, bit)
		require.Error(t, err)
	}
}
