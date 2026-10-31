package fatek

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNativeAddressCompleteComponentBoundsAndAlignment(t *testing.T) {
	for _, tc := range []struct {
		symbol  string
		address int
		valid   bool
	}{
		{"D", 65535, true}, {"D", 65536, false}, {"DR", 65534, true}, {"DR", 65535, false}, {"DR", 1, true},
		{"WM", 8, true}, {"WM", 24, true}, {"WM", 9, false}, {"WM", 9984, true}, {"WM", 9992, false},
		{"DWM", 8, true}, {"DWM", 40, true}, {"DWM", 9, false}, {"DWM", 9968, true}, {"DWM", 9976, false},
	} {
		comp, err := GetComponentType(tc.symbol)
		require.NoError(t, err)
		got, err := FormatAddress(comp, tc.address)
		if tc.valid {
			require.NoError(t, err, "%s%d", tc.symbol, tc.address)
			require.NotEmpty(t, got)
		} else {
			require.Error(t, err, "%s%d must not address a partial/misaligned component", tc.symbol, tc.address)
		}
	}
}

func TestNativeAddressClientRejectsIncompleteTailBeforeSend(t *testing.T) {
	tr := &randomStrictTransport{response: BuildFrame(1, "46", "00102030401020304")}
	got, err := NewClient(tr, 1).ReadRegisters("DR", 65534, 2)
	require.Error(t, err)
	require.Nil(t, got)
	require.Empty(t, tr.sent)
}
