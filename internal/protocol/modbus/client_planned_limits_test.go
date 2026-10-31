package modbus

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/protocol/wirebudget"
)

func TestPlannedExactNativeLimitsAndLastAddress(t *testing.T) {
	for _, rtu := range []bool{false, true} {
		t.Run(fmt.Sprint(rtu), func(t *testing.T) {
			p := plannedBasicPeer(rtu)
			for i := 0; i < 2000; i++ {
				p.bits[plannedCell{1, uint16(63536 + i)}] = i == 0 || i == 7 || i == 1999
			}
			for i := 0; i < 125; i++ {
				p.words[plannedCell{4, uint16(65411 + i)}] = uint16(i + 1)
			}
			var responseLengths []int
			p.mutate = func(r []byte) []byte { responseLengths = append(responseLengths, len(r)); return r }
			ctx := wirebudget.With(t.Context(), 2, 0)
			c := plannedExactReader(t, NewClient(p, 247).WithContext(ctx))
			bits, err := c.ReadBitsExact(1, 63536, 2000)
			require.NoError(t, err)
			require.Len(t, bits, 2000)
			for i, v := range bits {
				require.Equal(t, i == 0 || i == 7 || i == 1999, v)
			}
			words, err := c.ReadWordsExact(4, 65411, 125)
			require.NoError(t, err)
			require.Len(t, words, 125)
			for i, v := range words {
				require.Equal(t, uint16(i+1), v)
			}
			reqBytes, respBytes := 12, 259
			if rtu {
				reqBytes, respBytes = 8, 255
			}
			require.Len(t, p.frames, 2)
			require.Len(t, p.frames[0], reqBytes)
			require.Len(t, p.frames[1], reqBytes)
			require.Equal(t, []int{respBytes, respBytes}, responseLengths)
			require.Equal(t, 2, wirebudget.From(ctx).Used())
		})
	}
}
