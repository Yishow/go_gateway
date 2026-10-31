package adapters

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestModbusPlannedWordsOverlapAndFrozenFormats(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		t.Run(string(kind), func(t *testing.T) {
			req := modbusPlannedRequest(t, kind, optimizer.ReadableCoverage{}, modbusPlannedPoint("wide", "HR0", schema.DataTypeInt64, "ABCD"), modbusPlannedPoint("small", "HR1", schema.DataTypeInt16, ""), modbusPlannedPoint("abcd", "HR0", schema.DataTypeUint32, "ABCD"), modbusPlannedPoint("cdab", "HR0", schema.DataTypeUint32, "CDAB"))
			require.Len(t, req.Plan.Requests(), 1)
			peer := &modbusPlannedPeer{memory: map[modbusPlannedAddress]uint16{{3, 0}: 0x0102, {3, 1}: 0x0304, {3, 2}: 0x0506, {3, 3}: 0x0708}}
			got, err := requireModbusPlannedReader(t, modbusPlannedProtocol(kind, peer)).ReadPlanned(t.Context(), req)
			require.NoError(t, err)
			assertModbusPlannedResults(t, req, got, map[string]any{"wide": int64(0x0102030405060708), "small": int16(0x0304), "abcd": uint32(0x01020304), "cdab": uint32(0x03040102)})
			require.Len(t, peer.frames, 1)
			r := req.Plan.Requests()[0]
			require.Equal(t, r.FrameBytes.RequestBytes, len(peer.frames[0]))
			require.Equal(t, r.FrameBytes.ResponseBytes, len(peer.responses[0]))
			offset := 7
			if kind == schema.ProtocolModbusRTU {
				offset = 1
			}
			require.Equal(t, []byte{3, 0, 0, 0, 4}, peer.frames[0][offset:offset+5])
			for _, v := range got {
				if v.Binding.Point.ID == "cdab" {
					require.Equal(t, []byte{1, 2, 3, 4}, v.RawBytes)
				}
			}
		})
	}
}
func TestModbusPlannedAllBitsNotFirstBool(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		for _, area := range []string{"C", "DI"} {
			t.Run(string(kind)+area, func(t *testing.T) {
				bits := []bool{true, false, true, false, false, true, false, true, false, true}
				var points []schema.Point
				want := map[string]any{}
				memory := map[modbusPlannedAddress]uint16{}
				function := byte(1)
				if area == "DI" {
					function = 2
				}
				for i, b := range bits {
					id := string(rune('a' + i))
					points = append(points, modbusPlannedPoint(id, area+string(rune('0'+i)), schema.DataTypeBool, ""))
					want[id] = b
					var v uint16
					if b {
						v = 1
					}
					memory[modbusPlannedAddress{function, uint16(i)}] = v
				}
				req := modbusPlannedRequest(t, kind, optimizer.ReadableCoverage{}, points...)
				require.Len(t, req.Plan.Requests(), 1)
				peer := &modbusPlannedPeer{memory: memory}
				got, err := requireModbusPlannedReader(t, modbusPlannedProtocol(kind, peer)).ReadPlanned(t.Context(), req)
				require.NoError(t, err)
				assertModbusPlannedResults(t, req, got, want)
				require.Len(t, peer.frames, 1)
				offset := 8
				if kind == schema.ProtocolModbusRTU {
					offset = 2
				}
				require.Equal(t, []byte{2, 0xa5, 2}, peer.responses[0][offset:offset+3])
			})
		}
	}
}
