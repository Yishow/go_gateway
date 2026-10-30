package fatek

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type statusPayloadFixtureTransport struct {
	payload    string
	corruptLRC bool
	sent       [][]byte
}

func (t *statusPayloadFixtureTransport) Connect() error {
	return nil
}

func (t *statusPayloadFixtureTransport) Close() error {
	return nil
}

func (t *statusPayloadFixtureTransport) SendReceive(request []byte) ([]byte, error) {
	t.sent = append(t.sent, append([]byte(nil), request...))
	response := append([]byte{STX}, request[1:3]...)
	response = append(response, '4', '4', '0')
	response = append(response, t.payload...)
	response = append(response, []byte(CalculateLRC(response))...)
	response = append(response, ETX)
	if t.corruptLRC {
		response[len(response)-3] = '0'
	}
	return response, nil
}

func TestFatekStatusPayloadRejectsMalformed(t *testing.T) {
	for _, payload := range []string{"", "1", "101", "1X"} {
		t.Run(fmtStatusPayloadName(payload), func(t *testing.T) {
			transport := &statusPayloadFixtureTransport{payload: payload}
			client := NewClient(transport, 1)

			status, err := client.ReadStatus("X", 0, 2)
			require.Error(t, err)
			require.Nil(t, status)
			require.Len(t, transport.sent, 1)
		})
	}
}

func TestFatekStatusPayloadAcceptsValid(t *testing.T) {
	cases := []struct {
		name    string
		count   int
		payload string
		check   func(t *testing.T, status []bool)
	}{
		{
			name:    "scalar",
			count:   1,
			payload: "1",
			check: func(t *testing.T, status []bool) {
				require.Equal(t, []bool{true}, status)
			},
		},
		{
			name:    "batch",
			count:   2,
			payload: "01",
			check: func(t *testing.T, status []bool) {
				require.Equal(t, []bool{false, true}, status)
			},
		},
		{
			name:    "maximum",
			count:   255,
			payload: strings.Repeat("01", 127) + "1",
			check: func(t *testing.T, status []bool) {
				require.Len(t, status, 255)
				require.False(t, status[0])
				require.True(t, status[1])
				require.True(t, status[254])
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			transport := &statusPayloadFixtureTransport{payload: tc.payload}
			client := NewClient(transport, 1)

			status, err := client.ReadStatus("X", 0, tc.count)
			require.NoError(t, err)
			tc.check(t, status)
			require.Len(t, transport.sent, 1)
		})
	}
}

func TestFatekStatusPayloadRejectsInvalidCountBeforeWire(t *testing.T) {
	for _, count := range []int{0, 256} {
		t.Run(fmtStatusPayloadName(count), func(t *testing.T) {
			transport := &statusPayloadFixtureTransport{payload: "1"}
			client := NewClient(transport, 1)

			status, err := client.ReadStatus("X", 0, count)
			require.Error(t, err)
			require.Nil(t, status)
			require.Empty(t, transport.sent)
		})
	}
}

func TestFatekStatusPayloadRetainsFrameValidation(t *testing.T) {
	transport := &statusPayloadFixtureTransport{payload: "01", corruptLRC: true}
	client := NewClient(transport, 1)

	status, err := client.ReadStatus("X", 0, 2)
	require.Error(t, err)
	require.Nil(t, status)
	require.Len(t, transport.sent, 1)
}

func fmtStatusPayloadName(value any) string {
	name := strings.ReplaceAll(fmt.Sprint(value), " ", "_")
	if name == "" {
		return "empty"
	}
	return name
}
