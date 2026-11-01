package fatek

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFrameMinimumLengthsPreserveLoopbackAndErrorPriority(t *testing.T) {
	for length := 0; length < 8; length++ {
		for _, command := range []string{"4E", "41"} {
			_, err := ParseResponse(make([]byte, length), command)
			require.ErrorIs(t, err, ErrResponseTooShort)
		}
	}
	loopback := BuildFrame(1, "4E", "")
	require.Len(t, loopback, 8)
	body, err := ParseResponse(loopback, "4E")
	require.NoError(t, err)
	require.Empty(t, body)
	_, err = ParseResponse(loopback, "41")
	require.ErrorIs(t, err, ErrResponseTooShort)
	ack := BuildFrame(1, "41", "0")
	require.Len(t, ack, 9)
	body, err = ParseResponse(ack, "41")
	require.NoError(t, err)
	require.Empty(t, body)
}
