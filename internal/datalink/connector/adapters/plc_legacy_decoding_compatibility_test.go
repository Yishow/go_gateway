package adapters

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/mcprotocol"

	"github.com/stretchr/testify/require"
)

// These characterize valid legacy reads, not future point-level format support
// or a unified Count contract: Modbus counts words, MC counts typed values,
// and FATEK counts components of the addressed device width.
func TestPLCLegacyCompatibility_ModbusFormatFallbackAndWordCount(t *testing.T) {
	cases := []struct {
		name, pointFormat, connectionFormat string
		words                               []uint16
	}{
		{name: "omitted_formats_default_ABCD", words: []uint16{0x0102, 0x0304}},
		{name: "omitted_point_inherits_CDAB", connectionFormat: "CDAB", words: []uint16{0x0304, 0x0102}},
		{name: "blank_point_inherits_CDAB", pointFormat: " \t", connectionFormat: "CDAB", words: []uint16{0x0304, 0x0102}},
		{name: "point_BADC_overrides_CDAB", pointFormat: "BADC", connectionFormat: "CDAB", words: []uint16{0x0201, 0x0403}},
	}
	for _, adapter := range []string{"tcp", "rtu", "udp"} {
		for _, tc := range cases {
			for _, count := range []int{0, 2} {
				t.Run(fmt.Sprintf("%s/%s/count_%d", adapter, tc.name, count), func(t *testing.T) {
					transport := newPLCAdapterModbusTransport(tc.words)
					protocol := newPLCAdapterModbusConnector(adapter, transport)
					loadConnectionConfigDataFormat(t, protocol, tc.connectionFormat)
					result, err := protocol.Read(t.Context(), connector.ReadRequest{
						Address: "HR100", DataType: schema.DataTypeUint32,
						DataFormat: tc.pointFormat, Count: count,
					})
					require.NoError(t, err)
					require.Equal(t, schema.QualityGood, result.Quality)
					require.Equal(t, uint32(0x01020304), result.Value)
					require.Len(t, transport.requests, 1)
					require.Equal(t, uint16(2), transport.requests[0].count)
				})
			}
		}
	}
}

func TestPLCLegacyCompatibility_MCConnectionDefaultAndTypedCount(t *testing.T) {
	for _, format := range []string{"", "ABCD"} {
		for _, count := range []int{0, 1, 2} {
			t.Run(fmt.Sprintf("format_%q/count_%d", format, count), func(t *testing.T) {
				// Only Connect initialization uses loopback. All protocol reads use
				// the existing in-memory wire fixture; no PLC endpoint is contacted.
				listener, err := net.Listen("tcp", "127.0.0.1:0")
				require.NoError(t, err)
				t.Cleanup(func() { _ = listener.Close() })
				config := map[string]any{
					"host": "127.0.0.1", "port": listener.Addr().(*net.TCPAddr).Port,
				}
				if format != "" {
					config["data_format"] = format
				}
				raw, err := json.Marshal(config)
				require.NoError(t, err)
				protocol := &MC3EConnector{}
				t.Cleanup(func() { _ = protocol.Close() })
				require.NoError(t, protocol.Connect(t.Context(), string(raw)))
				require.NoError(t, protocol.Close())

				words := []uint16{0x0304, 0x0102, 0x0708, 0x0506}
				if format == "ABCD" {
					words = []uint16{0x0102, 0x0304, 0x0506, 0x0708}
				}
				transport := &plcAdapterMCTransport{words: words}
				protocol.client = mcprotocol.NewClientWithTransport(transport)
				protocol.connected = true
				// Deliberately omit point DataFormat: this tests the connection
				// converter, not the currently unsupported point override contract.
				result, err := protocol.Read(t.Context(), connector.ReadRequest{
					Address: "D0", DataType: schema.DataTypeUint32, Count: count,
				})
				require.NoError(t, err)
				require.Equal(t, schema.QualityGood, result.Quality)
				wantWords := uint16(2)
				if count == 2 {
					require.Equal(t, []any{uint32(0x01020304), uint32(0x05060708)}, result.Value)
					wantWords = 4
				} else {
					require.Equal(t, uint32(0x01020304), result.Value)
				}
				require.Len(t, transport.requests, 1)
				require.Equal(t, wantWords, transport.requests[0].count)
			})
		}
	}
}

func TestPLCLegacyCompatibility_FatekComponentWidthAnd16BitDefaults(t *testing.T) {
	cases := []struct {
		name, address, payload string
		dataType               schema.DataType
		count, wireCount       int
		want                   any
	}{
		{"D_uint32_two_16bit_components", "D0", "01020304", schema.DataTypeUint32, 2, 2, uint32(0x01020304)},
		{"DR_uint32_one_32bit_component", "DR0", "01020304", schema.DataTypeUint32, 1, 1, uint32(0x01020304)},
		{"DR_uint32_default_one_component", "DR0", "01020304", schema.DataTypeUint32, 0, 1, uint32(0x01020304)},
		{"D_uint64_four_16bit_components", "D0", "0102030405060708", schema.DataTypeUint64, 4, 4, uint64(0x0102030405060708)},
		{"DR_uint64_two_32bit_components", "DR0", "0102030405060708", schema.DataTypeUint64, 2, 2, uint64(0x0102030405060708)},
		{"D_int16_default", "D0", "FFFE", schema.DataTypeInt16, 0, 1, int16(-2)},
		{"D_uint16_default", "D0", "FFFE", schema.DataTypeUint16, 0, 1, uint16(0xFFFE)},
		{"D_bool_default", "D0", "0001", schema.DataTypeBool, 0, 1, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			protocol, transport := newPLCFatekReadFixture(tc.payload)
			result, err := protocol.Read(t.Context(), connector.ReadRequest{
				Address: tc.address, DataType: tc.dataType, Count: tc.count,
			})
			require.NoError(t, err)
			require.Equal(t, schema.QualityGood, result.Quality)
			require.Equal(t, tc.want, result.Value)
			require.Equal(t, tc.wireCount, plcFatekRequestCount(t, transport))
		})
	}
}
