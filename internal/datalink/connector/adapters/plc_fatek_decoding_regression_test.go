package adapters

import (
	"fmt"
	"strconv"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"

	"github.com/stretchr/testify/require"
)

type plcFatekReadFixtureTransport struct {
	payload  string
	requests [][]byte
}

func (t *plcFatekReadFixtureTransport) Connect() error {
	return nil
}

func (t *plcFatekReadFixtureTransport) Close() error {
	return nil
}

func (t *plcFatekReadFixtureTransport) SendReceive(request []byte) ([]byte, error) {
	if len(request) < 7 {
		return nil, fmt.Errorf("FATEK request too short: %d", len(request))
	}
	t.requests = append(t.requests, append([]byte(nil), request...))
	if string(request[3:5]) != "46" {
		return nil, fmt.Errorf("unexpected FATEK command %q", request[3:5])
	}

	response := append([]byte{fatek.STX}, request[1:3]...)
	response = append(response, '4', '6', '0')
	response = append(response, t.payload...)
	response = append(response, []byte(fatek.CalculateLRC(response))...)
	response = append(response, fatek.ETX)
	return response, nil
}

func newPLCFatekReadFixture(payload string) (*FatekConnector, *plcFatekReadFixtureTransport) {
	transport := &plcFatekReadFixtureTransport{payload: payload}
	protocol := &FatekConnector{
		client:         fatek.NewClient(transport, 1),
		connected:      true,
		persistentMode: true,
	}
	return protocol, transport
}

func plcFatekRequestCount(t *testing.T, transport *plcFatekReadFixtureTransport) int {
	t.Helper()
	require.Len(t, transport.requests, 1)
	count, err := strconv.ParseUint(string(transport.requests[0][5:7]), 16, 8)
	require.NoError(t, err)
	return int(count)
}

func plcFatekSigned32Value() int32 {
	bits := uint32(0xFEDCBA98)
	return int32(bits)
}

func plcFatekSigned64Value() int64 {
	bits := uint64(0xFEDCBA9876543210)
	return int64(bits)
}

func TestPLCFatekDecoding_RealReadTypedValues(t *testing.T) {
	cases := []struct {
		name     string
		address  string
		dataType schema.DataType
		count    int
		payload  string
		want     any
	}{
		{name: "D_float32", address: "D0", dataType: schema.DataTypeFloat32, count: 2, payload: "3F800000", want: float32(1)},
		{name: "DR_float32", address: "DR0", dataType: schema.DataTypeFloat32, count: 1, payload: "3F800000", want: float32(1)},
		{name: "D_int32", address: "D0", dataType: schema.DataTypeInt32, count: 2, payload: "FEDCBA98", want: plcFatekSigned32Value()},
		{name: "DR_int32", address: "DR0", dataType: schema.DataTypeInt32, count: 1, payload: "FEDCBA98", want: plcFatekSigned32Value()},
		{name: "D_uint32", address: "D0", dataType: schema.DataTypeUint32, count: 2, payload: "01020304", want: uint32(0x01020304)},
		{name: "DR_uint32", address: "DR0", dataType: schema.DataTypeUint32, count: 1, payload: "01020304", want: uint32(0x01020304)},
		{name: "D_int64", address: "D0", dataType: schema.DataTypeInt64, count: 4, payload: "FEDCBA9876543210", want: plcFatekSigned64Value()},
		{name: "DR_int64", address: "DR0", dataType: schema.DataTypeInt64, count: 2, payload: "FEDCBA9876543210", want: plcFatekSigned64Value()},
		{name: "D_uint64", address: "D0", dataType: schema.DataTypeUint64, count: 4, payload: "0102030405060708", want: uint64(0x0102030405060708)},
		{name: "DR_uint64", address: "DR0", dataType: schema.DataTypeUint64, count: 2, payload: "0102030405060708", want: uint64(0x0102030405060708)},
		{name: "D_float64", address: "D0", dataType: schema.DataTypeFloat64, count: 4, payload: "3FF0000000000000", want: float64(1)},
		{name: "DR_float64", address: "DR0", dataType: schema.DataTypeFloat64, count: 2, payload: "3FF0000000000000", want: float64(1)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			protocol, transport := newPLCFatekReadFixture(tc.payload)
			result, err := protocol.Read(t.Context(), connector.ReadRequest{
				Address:  tc.address,
				DataType: tc.dataType,
				Count:    tc.count,
			})
			require.NoError(t, err)
			require.Equal(t, schema.QualityGood, result.Quality)
			require.Equal(t, tc.want, result.Value)
			require.Equal(t, tc.count, plcFatekRequestCount(t, transport))
		})
	}
}

func TestPLCFatekDecoding_RejectsIncompleteOrUnsupportedValues(t *testing.T) {
	cases := []struct {
		name     string
		address  string
		dataType schema.DataType
		count    int
		payload  string
	}{
		{name: "D_float32_short", address: "D0", dataType: schema.DataTypeFloat32, count: 1, payload: "3F80"},
		{name: "D_int64_short", address: "D0", dataType: schema.DataTypeInt64, count: 3, payload: "010203040506"},
		{name: "D_uint64_short", address: "D0", dataType: schema.DataTypeUint64, count: 3, payload: "010203040506"},
		{name: "D_float64_short", address: "D0", dataType: schema.DataTypeFloat64, count: 3, payload: "3FF000000000"},
		{name: "DR_int64_short", address: "DR0", dataType: schema.DataTypeInt64, count: 1, payload: "01020304"},
		{name: "DR_uint64_short", address: "DR0", dataType: schema.DataTypeUint64, count: 1, payload: "01020304"},
		{name: "DR_float64_short", address: "DR0", dataType: schema.DataTypeFloat64, count: 1, payload: "3FF00000"},
		{name: "D_string_unsupported", address: "D0", dataType: schema.DataTypeString, count: 1, payload: "0001"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			protocol, transport := newPLCFatekReadFixture(tc.payload)
			result, err := protocol.Read(t.Context(), connector.ReadRequest{
				Address:  tc.address,
				DataType: tc.dataType,
				Count:    tc.count,
			})
			require.Error(t, err)
			require.Equal(t, schema.QualityBad, result.Quality)
			require.Nil(t, result.Value)
			require.Equal(t, tc.count, plcFatekRequestCount(t, transport))
		})
	}
}

func TestPLCFatekDecoding_SingleComponentCompatibility(t *testing.T) {
	cases := []struct {
		name     string
		dataType schema.DataType
		payload  string
		want     any
	}{
		{name: "bool", dataType: schema.DataTypeBool, payload: "0001", want: true},
		{name: "int16", dataType: schema.DataTypeInt16, payload: "FFFE", want: int16(-2)},
		{name: "uint16", dataType: schema.DataTypeUint16, payload: "FFFE", want: uint16(0xFFFE)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			protocol, transport := newPLCFatekReadFixture(tc.payload)
			result, err := protocol.Read(t.Context(), connector.ReadRequest{
				Address:  "D0",
				DataType: tc.dataType,
				Count:    1,
			})
			require.NoError(t, err)
			require.Equal(t, schema.QualityGood, result.Quality)
			require.Equal(t, tc.want, result.Value)
			require.Equal(t, 1, plcFatekRequestCount(t, transport))
		})
	}
}
