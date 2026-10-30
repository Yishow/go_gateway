package adapters

import (
	"strings"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"

	"github.com/stretchr/testify/require"
)

type plcFatekStatusFixtureTransport struct {
	payload    string
	corruptLRC bool
	sent       [][]byte
}

func (t *plcFatekStatusFixtureTransport) Connect() error {
	return nil
}

func (t *plcFatekStatusFixtureTransport) Close() error {
	return nil
}

func (t *plcFatekStatusFixtureTransport) SendReceive(request []byte) ([]byte, error) {
	t.sent = append(t.sent, append([]byte(nil), request...))
	response := append([]byte{fatek.STX}, request[1:3]...)
	response = append(response, '4', '4', '0')
	response = append(response, t.payload...)
	response = append(response, []byte(fatek.CalculateLRC(response))...)
	response = append(response, fatek.ETX)
	if t.corruptLRC {
		response[len(response)-3] = '0'
	}
	return response, nil
}

func newPLCFatekStatusFixture(payload string) (*FatekConnector, *plcFatekStatusFixtureTransport) {
	transport := &plcFatekStatusFixtureTransport{payload: payload}
	protocol := &FatekConnector{
		client:         fatek.NewClient(transport, 1),
		connected:      true,
		persistentMode: true,
	}
	return protocol, transport
}

func TestPLCFatekStatusPayloadRejectsMalformed(t *testing.T) {
	for _, payload := range []string{"", "1", "101", "1X"} {
		t.Run(plcFatekStatusPayloadName(payload), func(t *testing.T) {
			protocol, transport := newPLCFatekStatusFixture(payload)
			result, err := protocol.Read(t.Context(), connector.ReadRequest{
				Address:  "X0",
				DataType: schema.DataTypeBool,
				Count:    2,
			})
			require.Error(t, err)
			require.Equal(t, schema.QualityBad, result.Quality)
			require.Nil(t, result.Value)
			require.Len(t, transport.sent, 1)
		})
	}
}

func TestPLCFatekStatusPayloadAcceptsValid(t *testing.T) {
	cases := []struct {
		name    string
		count   int
		payload string
		want    any
	}{
		{name: "scalar", count: 1, payload: "1", want: true},
		{name: "batch", count: 2, payload: "01", want: []bool{false, true}},
		{name: "maximum", count: 255, payload: strings.Repeat("01", 127) + "1", want: nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			protocol, transport := newPLCFatekStatusFixture(tc.payload)
			result, err := protocol.Read(t.Context(), connector.ReadRequest{
				Address:  "X0",
				DataType: schema.DataTypeBool,
				Count:    tc.count,
			})
			require.NoError(t, err)
			require.Equal(t, schema.QualityGood, result.Quality)
			if tc.want != nil {
				require.Equal(t, tc.want, result.Value)
			} else {
				values, ok := result.Value.([]bool)
				require.True(t, ok)
				require.Len(t, values, 255)
				require.False(t, values[0])
				require.True(t, values[1])
				require.True(t, values[254])
			}
			require.Len(t, transport.sent, 1)
		})
	}
}

func TestPLCFatekStatusPayloadRejectsInvalidFrame(t *testing.T) {
	protocol, transport := newPLCFatekStatusFixture("01")
	transport.corruptLRC = true

	result, err := protocol.Read(t.Context(), connector.ReadRequest{
		Address:  "X0",
		DataType: schema.DataTypeBool,
		Count:    2,
	})
	require.Error(t, err)
	require.Equal(t, schema.QualityBad, result.Quality)
	require.Nil(t, result.Value)
	require.Len(t, transport.sent, 1)
}

func plcFatekStatusPayloadName(payload string) string {
	if payload == "" {
		return "empty"
	}
	return strings.ReplaceAll(payload, " ", "_")
}
