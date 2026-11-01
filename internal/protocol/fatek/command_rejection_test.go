package fatek

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandRejectionRequiresValidatedRandomStatusNine(t *testing.T) {
	for _, tc := range []struct {
		name     string
		frame    []byte
		rejected bool
	}{
		{"unsupported", BuildFrame(1, "48", "9"), true},
		{"ambiguous-four", BuildFrame(1, "48", "4"), false},
		{"foreign-station", BuildFrame(2, "48", "9"), false},
		{"foreign-command", BuildFrame(1, "46", "9"), false},
		{"extra-payload", BuildFrame(1, "48", "90000"), false},
		{"truncated", []byte("\x0201489"), false},
		{"bad-lrc", []byte("\x020148900\x03"), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := &randomStrictTransport{response: tc.frame}
			_, err := NewClient(transport, 1).ReadRandom([]RandomReadItem{{Symbol: "D", Addr: 0}})
			require.Error(t, err)
			require.Equal(t, tc.rejected, IsCommandRejected(fmt.Errorf("read: %w", err)))
		})
	}
	require.False(t, IsCommandRejected(fmt.Errorf("unsupported 48 status 9")))
	require.False(t, IsCommandRejected(NewProtocolError("9", "46")))
}
