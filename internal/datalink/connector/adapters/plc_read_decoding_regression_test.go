package adapters

import (
	"encoding/binary"
	"fmt"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/internal/protocol/modbus"
	"go-gateway/lib/hsllogic"

	"github.com/stretchr/testify/require"
)

const plcAdapterModbusStartAddress uint16 = 100

type plcAdapterWireVector struct {
	format   hsllogic.DataFormat
	float32  []uint16
	uint32   []uint16
	signed32 []uint16
	float64  []uint16
	uint64   []uint16
	signed64 []uint16
}

type plcAdapterValueCase struct {
	name     string
	dataType schema.DataType
	words    []uint16
	count    int
	expected any
}

func plcAdapterWireVectors() []plcAdapterWireVector {
	return []plcAdapterWireVector{
		{
			format:   hsllogic.DataFormatABCD,
			float32:  []uint16{0x3F80, 0x0000},
			uint32:   []uint16{0x0102, 0x0304},
			signed32: []uint16{0xFEDC, 0xBA98},
			float64:  []uint16{0x3FF0, 0x0000, 0x0000, 0x0000},
			uint64:   []uint16{0x0102, 0x0304, 0x0506, 0x0708},
			signed64: []uint16{0xFEDC, 0xBA98, 0x7654, 0x3210},
		},
		{
			format:   hsllogic.DataFormatBADC,
			float32:  []uint16{0x803F, 0x0000},
			uint32:   []uint16{0x0201, 0x0403},
			signed32: []uint16{0xDCFE, 0x98BA},
			float64:  []uint16{0xF03F, 0x0000, 0x0000, 0x0000},
			uint64:   []uint16{0x0201, 0x0403, 0x0605, 0x0807},
			signed64: []uint16{0xDCFE, 0x98BA, 0x5476, 0x1032},
		},
		{
			format:   hsllogic.DataFormatCDAB,
			float32:  []uint16{0x0000, 0x3F80},
			uint32:   []uint16{0x0304, 0x0102},
			signed32: []uint16{0xBA98, 0xFEDC},
			float64:  []uint16{0x0000, 0x0000, 0x0000, 0x3FF0},
			uint64:   []uint16{0x0708, 0x0506, 0x0304, 0x0102},
			signed64: []uint16{0x3210, 0x7654, 0xBA98, 0xFEDC},
		},
		{
			format:   hsllogic.DataFormatDCBA,
			float32:  []uint16{0x0000, 0x803F},
			uint32:   []uint16{0x0403, 0x0201},
			signed32: []uint16{0x98BA, 0xDCFE},
			float64:  []uint16{0x0000, 0x0000, 0x0000, 0xF03F},
			uint64:   []uint16{0x0807, 0x0605, 0x0403, 0x0201},
			signed64: []uint16{0x1032, 0x5476, 0x98BA, 0xDCFE},
		},
	}
}

func plcAdapterSigned32Value() int32 {
	bits := uint32(0xFEDCBA98)
	return int32(bits)
}

func plcAdapterSigned64Value() int64 {
	bits := uint64(0xFEDCBA9876543210)
	return int64(bits)
}

func plcAdapterValueCases(vector plcAdapterWireVector) []plcAdapterValueCase {
	return []plcAdapterValueCase{
		{name: "float32", dataType: schema.DataTypeFloat32, words: vector.float32, count: 2, expected: float32(1)},
		{name: "uint32", dataType: schema.DataTypeUint32, words: vector.uint32, count: 2, expected: uint32(0x01020304)},
		{name: "signed32", dataType: schema.DataTypeInt32, words: vector.signed32, count: 2, expected: plcAdapterSigned32Value()},
		{name: "float64", dataType: schema.DataTypeFloat64, words: vector.float64, count: 4, expected: float64(1)},
		{name: "uint64", dataType: schema.DataTypeUint64, words: vector.uint64, count: 4, expected: uint64(0x0102030405060708)},
		{name: "signed64", dataType: schema.DataTypeInt64, words: vector.signed64, count: 4, expected: plcAdapterSigned64Value()},
	}
}

type plcAdapterModbusRequest struct {
	function byte
	address  uint16
	count    uint16
}

type plcAdapterModbusTransport struct {
	*MockTransport
	requests []plcAdapterModbusRequest
}

func newPLCAdapterModbusTransport(words []uint16) *plcAdapterModbusTransport {
	transport := &plcAdapterModbusTransport{MockTransport: NewMockTransport()}
	for i, word := range words {
		transport.HoldingRegs[plcAdapterModbusStartAddress+uint16(i)] = word
	}
	return transport
}

func (t *plcAdapterModbusTransport) SendReceive(data []byte) ([]byte, error) {
	_, _, function, pdu, err := modbus.ParseTCPFrame(data)
	if err != nil {
		return nil, err
	}
	if len(pdu) < 4 {
		return nil, fmt.Errorf("modbus read request too short: %d", len(pdu))
	}
	t.requests = append(t.requests, plcAdapterModbusRequest{
		function: function,
		address:  binary.BigEndian.Uint16(pdu[0:2]),
		count:    binary.BigEndian.Uint16(pdu[2:4]),
	})
	return t.MockTransport.SendReceive(data)
}

func newPLCAdapterModbusConnector(name string, transport modbus.Transport) connector.Protocol {
	client := modbus.NewClient(transport, 1)
	switch name {
	case "tcp":
		return &ModbusTCPConnector{client: client, connected: true, persistentMode: true}
	case "rtu":
		return &ModbusRTUConnector{client: client, connected: true, persistentMode: true}
	case "udp":
		return &ModbusUDPConnector{client: client, connected: true, persistentMode: true}
	default:
		panic("unsupported Modbus connector: " + name)
	}
}

func TestPLCReadDecoding_ModbusAdapters(t *testing.T) {
	for _, protocolName := range []string{"tcp", "rtu", "udp"} {
		for _, vector := range plcAdapterWireVectors() {
			for _, valueCase := range plcAdapterValueCases(vector) {
				t.Run(fmt.Sprintf("%s/%s/%s", protocolName, vector.format, valueCase.name), func(t *testing.T) {
					transport := newPLCAdapterModbusTransport(valueCase.words)
					protocol := newPLCAdapterModbusConnector(protocolName, transport)
					connectionFormat := string(vector.format)
					if vector.format == hsllogic.DataFormatABCD {
						connectionFormat = ""
					}
					loadConnectionConfigDataFormat(t, protocol, connectionFormat)

					result, err := protocol.Read(t.Context(), connector.ReadRequest{
						Address:  "HR100",
						DataType: valueCase.dataType,
					})
					require.NoError(t, err)
					require.Equal(t, schema.QualityGood, result.Quality)
					require.Equal(t, valueCase.expected, result.Value)
					require.Len(t, transport.requests, 1)
					require.Equal(t, byte(modbus.FuncReadHoldingRegisters), transport.requests[0].function)
					require.Equal(t, plcAdapterModbusStartAddress, transport.requests[0].address)
					require.Equal(t, uint16(valueCase.count), transport.requests[0].count)
				})
			}
		}
	}
}

func TestPLCReadDecoding_ModbusPointFormatOverridesConnection(t *testing.T) {
	vector := plcAdapterWireVectors()[1]
	transport := newPLCAdapterModbusTransport(vector.float32)
	protocol := newPLCAdapterModbusConnector("tcp", transport)
	loadConnectionConfigDataFormat(t, protocol, string(hsllogic.DataFormatCDAB))

	result, err := protocol.Read(t.Context(), connector.ReadRequest{
		Address:    "HR100",
		DataType:   schema.DataTypeFloat32,
		DataFormat: string(hsllogic.DataFormatBADC),
	})
	require.NoError(t, err)
	require.Equal(t, float32(1), result.Value)
	require.Len(t, transport.requests, 1)
	require.Equal(t, uint16(2), transport.requests[0].count)
}

type plcAdapterMCRequest struct {
	device  byte
	address uint32
	count   uint16
}

type plcAdapterMCTransport struct {
	words    []uint16
	requests []plcAdapterMCRequest
}

func (t *plcAdapterMCTransport) Connect() error { return nil }

func (t *plcAdapterMCTransport) Close() error { return nil }

func (t *plcAdapterMCTransport) SendReceive(request []byte) ([]byte, error) {
	if len(request) < 21 {
		return nil, fmt.Errorf("MC word request too short: %d", len(request))
	}
	address := uint32(request[15]) | uint32(request[16])<<8 | uint32(request[17])<<16
	count := binary.LittleEndian.Uint16(request[19:21])
	t.requests = append(t.requests, plcAdapterMCRequest{device: request[18], address: address, count: count})
	if int(count) > len(t.words) {
		return nil, fmt.Errorf("MC fixture has %d words, requested %d", len(t.words), count)
	}
	response := make([]byte, 2+int(count)*2)
	for i := range int(count) {
		binary.LittleEndian.PutUint16(response[2+i*2:], t.words[i])
	}
	return response, nil
}

func newPLCAdapterMCConnector(format hsllogic.DataFormat, transport *plcAdapterMCTransport) *MC3EConnector {
	return &MC3EConnector{
		client:         mcprotocol.NewClientWithTransport(transport),
		connected:      true,
		persistentMode: true,
		dataConverter:  hsllogic.NewDataConverter(format),
	}
}

func TestPLCReadDecoding_MC3EAdapter(t *testing.T) {
	for _, vector := range plcAdapterWireVectors() {
		for _, valueCase := range plcAdapterValueCases(vector) {
			t.Run(fmt.Sprintf("%s/%s", vector.format, valueCase.name), func(t *testing.T) {
				transport := &plcAdapterMCTransport{words: valueCase.words}
				protocol := newPLCAdapterMCConnector(vector.format, transport)

				result, err := protocol.Read(t.Context(), connector.ReadRequest{
					Address:  "D0",
					DataType: valueCase.dataType,
				})
				require.NoError(t, err)
				require.Equal(t, schema.QualityGood, result.Quality)
				require.Equal(t, valueCase.expected, result.Value)
				require.Len(t, transport.requests, 1)
				require.Equal(t, byte(0xA8), transport.requests[0].device)
				require.Zero(t, transport.requests[0].address)
				require.Equal(t, uint16(valueCase.count), transport.requests[0].count)
			})
		}
	}
}

func TestPLCReadDecoding_MCDirectTypedBatchCount(t *testing.T) {
	transport := &plcAdapterMCTransport{words: []uint16{34464, 1, 34465, 1}}
	protocol := newPLCAdapterMCConnector(hsllogic.DataFormatCDAB, transport)

	result, err := protocol.Read(t.Context(), connector.ReadRequest{
		Address:  "D3",
		DataType: schema.DataTypeUint32,
		Count:    2,
	})
	require.NoError(t, err)
	require.Equal(t, []any{uint32(100000), uint32(100001)}, result.Value)
	require.Len(t, transport.requests, 1)
	require.Equal(t, uint16(4), transport.requests[0].count)
}
