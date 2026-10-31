package adapters

import (
	"encoding/binary"
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
)

// profileWordTransport exercises the production client and adapter Read, without
// pretending to be a serial/TCP/UDP network integration test.
type profileWordTransport struct {
	words   []uint16
	sendsMu sync.Mutex // Protect only test instrumentation, not concurrent decoder calls.
	sends   int
}

func (*profileWordTransport) Connect() error               { return nil }
func (*profileWordTransport) Close() error                 { return nil }
func (*profileWordTransport) IsConnected() bool            { return true }
func (*profileWordTransport) GetNextTransactionID() uint16 { return 1 }
func (tr *profileWordTransport) SendReceive(request []byte) ([]byte, error) {
	tr.sendsMu.Lock()
	tr.sends++
	tr.sendsMu.Unlock()
	tx, unit, function, _, err := modbus.ParseTCPFrame(request)
	if err != nil {
		return nil, err
	}
	payload := make([]byte, 1+len(tr.words)*2)
	payload[0] = byte(len(tr.words) * 2)
	for i, w := range tr.words {
		binary.BigEndian.PutUint16(payload[1+2*i:], w)
	}
	return modbus.BuildTCPFrame(tx, unit, function, payload), nil
}
func profileProtocol(kind schema.ProtocolType, tr modbus.Transport) connector.Protocol {
	client := modbus.NewClient(tr, 1)
	switch kind {
	case schema.ProtocolModbusRTU:
		return &ModbusRTUConnector{client: client, connected: true, persistentMode: true}
	case schema.ProtocolModbusUDP:
		return &ModbusUDPConnector{client: client, connected: true, persistentMode: true}
	default:
		return &ModbusTCPConnector{client: client, connected: true, persistentMode: true}
	}
}
func profileRequest(t *testing.T, protocol schema.ProtocolType, dtype schema.DataType, format string) connector.ReadRequest {
	t.Helper()
	point := schema.Point{ID: "P", RevisionID: "revision", Address: "HR100", DataType: dtype, DataFormat: "ABCD", FormatMode: schema.PointFormatOverride, FormatOverride: format, FormatVersion: 1}
	profile, err := connector.ResolveReadProfile(point, decoding.FormatContext{Protocol: protocol})
	require.NoError(t, err)
	return connector.ReadRequest{Address: point.Address, DataType: dtype, DataFormat: point.DataFormat, Profile: profile}
}
func profileWireWords(bits uint64, count int, format string) []uint16 {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)
	bytes = bytes[8-count*2:]
	words := make([]uint16, count)
	for i := range words {
		words[i] = binary.BigEndian.Uint16(bytes[2*i:])
	}
	if format == "CDAB" || format == "DCBA" {
		for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
			words[i], words[j] = words[j], words[i]
		}
	}
	if format == "BADC" || format == "DCBA" {
		for i, w := range words {
			words[i] = (w << 8) | (w >> 8)
		}
	}
	return words
}
func TestModbusResolvedProfileNumericMatrix(t *testing.T) {
	fixtures := []struct {
		dtype schema.DataType
		bits  uint64
		want  any
		count int
	}{
		{schema.DataTypeInt32, 0xfedcba98, int32(-19088744), 2},
		{schema.DataTypeUint32, 0x12345678, uint32(0x12345678), 2},
		{schema.DataTypeFloat32, uint64(math.Float32bits(123.25)), float32(123.25), 2},
		{schema.DataTypeInt64, 0xfedcba9876543210, int64(-81985529216486896), 4},
		{schema.DataTypeUint64, 0x1122334455667788, uint64(0x1122334455667788), 4},
		{schema.DataTypeFloat64, math.Float64bits(-123456.25), float64(-123456.25), 4},
	}
	for _, protocol := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		for _, f := range fixtures {
			for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
				t.Run(string(protocol)+"/"+string(f.dtype)+"/"+format, func(t *testing.T) {
					tr := &profileWordTransport{words: profileWireWords(f.bits, f.count, format)}
					c := profileProtocol(protocol, tr)
					request := profileRequest(t, protocol, f.dtype, format)
					for _, address := range []string{"HR100", "IR100"} {
						p := schema.Point{ID: "P", Address: address, DataType: f.dtype, FormatMode: schema.PointFormatOverride, FormatOverride: format, FormatVersion: 1}
						var err error
						request.Profile, err = connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: protocol})
						require.NoError(t, err)
						request.Address = address
						result, err := c.Read(t.Context(), request)
						require.NoError(t, err)
						require.Equal(t, schema.QualityGood, result.Quality)
						require.Equal(t, f.want, result.Value)
						require.Equal(t, uint16SliceToBytes(tr.words), result.RawBytes)
					}
				})
			}
		}
	}
}
func TestModbusResolvedProfileInterleavedAndShort(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		t.Run(string(protocol), func(t *testing.T) {
			tr := &profileWordTransport{words: []uint16{0x1234, 0x5678}}
			c := profileProtocol(protocol, tr)
			a, b := profileRequest(t, protocol, schema.DataTypeUint32, "ABCD"), profileRequest(t, protocol, schema.DataTypeUint32, "CDAB")
			for _, test := range []struct {
				req  connector.ReadRequest
				want uint32
			}{{a, 0x12345678}, {b, 0x56781234}, {a, 0x12345678}, {b, 0x56781234}} {
				result, err := c.Read(t.Context(), test.req)
				require.NoError(t, err)
				require.Equal(t, test.want, result.Value)
			}
			for _, dtype := range []schema.DataType{schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32, schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64} {
				tr.words = []uint16{1}
				result, err := c.Read(t.Context(), profileRequest(t, protocol, dtype, "CDAB"))
				require.Error(t, err)
				require.NotEqual(t, schema.QualityGood, result.Quality)
				require.Nil(t, result.Value)
			}
		})
	}
}

func TestModbusResolvedProfilesConcurrentSameConnection(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		t.Run(string(protocol), func(t *testing.T) {
			tr := &profileWordTransport{words: []uint16{0x1234, 0x5678}}
			c := profileProtocol(protocol, tr)
			a, b := profileRequest(t, protocol, schema.DataTypeUint32, "ABCD"), profileRequest(t, protocol, schema.DataTypeUint32, "CDAB")
			type answer struct {
				value any
				err   error
				want  uint32
			}
			answers := make(chan answer, 40)
			for i := 0; i < 20; i++ {
				go func() { result, err := c.Read(t.Context(), a); answers <- answer{result.Value, err, 0x12345678} }()
				go func() { result, err := c.Read(t.Context(), b); answers <- answer{result.Value, err, 0x56781234} }()
			}
			for i := 0; i < 40; i++ {
				result := <-answers
				require.NoError(t, result.err)
				require.Equal(t, result.want, result.value)
			}
		})
	}
}
