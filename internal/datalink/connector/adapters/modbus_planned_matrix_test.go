package adapters

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestModbusPlannedCompleteScalarMatrix(t *testing.T) {
	// Literals describe the same complete values under all four point policies;
	// expected values do not come from the scalar adapter or decoder under test.
	for _, layout := range []struct {
		format             string
		w32, w64, f32, f64 []uint16
	}{
		{"ABCD", []uint16{0x8123, 0x4567}, []uint16{0x8123, 0x4567, 0x89ab, 0xcdef}, []uint16{0xc120, 0}, []uint16{0xc024, 0, 0, 0}},
		{"BADC", []uint16{0x2381, 0x6745}, []uint16{0x2381, 0x6745, 0xab89, 0xefcd}, []uint16{0x20c1, 0}, []uint16{0x24c0, 0, 0, 0}},
		{"CDAB", []uint16{0x4567, 0x8123}, []uint16{0xcdef, 0x89ab, 0x4567, 0x8123}, []uint16{0, 0xc120}, []uint16{0, 0, 0, 0xc024}},
		{"DCBA", []uint16{0x6745, 0x2381}, []uint16{0xefcd, 0xab89, 0x6745, 0x2381}, []uint16{0, 0x20c1}, []uint16{0, 0, 0, 0x24c0}},
	} {
		for _, tc := range []struct {
			dtype schema.DataType
			words []uint16
			want  any
		}{
			{schema.DataTypeInt32, layout.w32, int32(-2128394905)}, {schema.DataTypeUint32, layout.w32, uint32(0x81234567)},
			{schema.DataTypeFloat32, layout.f32, float32(-10)}, {schema.DataTypeInt64, layout.w64, int64(-9141386507638288913)},
			{schema.DataTypeUint64, layout.w64, uint64(0x8123456789abcdef)}, {schema.DataTypeFloat64, layout.f64, float64(-10)},
		} {
			for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
				for _, area := range []string{"HR", "IR"} {
					t.Run(fmt.Sprint(layout.format, "/", tc.dtype, "/", kind, "/", area), func(t *testing.T) {
						req := modbusPlannedRequest(t, kind, optimizer.ReadableCoverage{}, modbusPlannedPoint("scalar", area+"0", tc.dtype, layout.format))
						memory := map[modbusPlannedAddress]uint16{}
						function := byte(3)
						if area == "IR" {
							function = 4
						}
						for i, v := range tc.words {
							memory[modbusPlannedAddress{function, uint16(i)}] = v
						}
						peer := &modbusPlannedPeer{memory: memory}
						c := modbusPlannedProtocol(kind, peer)
						got, err := requireModbusPlannedReader(t, c).ReadPlanned(t.Context(), req)
						require.NoError(t, err)
						assertModbusPlannedResults(t, req, got, map[string]any{"scalar": tc.want})
						require.Len(t, peer.frames, 1)
						// Existing scalar API remains an additional compatibility check, never the oracle.
						b := req.Bindings[0]
						scalar, err := c.Read(t.Context(), connector.ReadRequest{Address: b.Point.Address, DataType: b.Point.DataType, Profile: b.Profile})
						require.NoError(t, err)
						require.Equal(t, tc.want, scalar.Value)
					})
				}
			}
		}
	}
}
func TestModbusPlannedLegacyAndNative16KeepScalarSemantics(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		for _, mode := range []schema.PointFormatMode{schema.PointFormatLegacy, schema.PointFormatInherit} {
			t.Run(fmt.Sprint(kind, "/", mode), func(t *testing.T) {
				points := []schema.Point{modbusPlannedPoint("signed", "HR0", schema.DataTypeInt16, ""), modbusPlannedPoint("unsigned", "HR0", schema.DataTypeUint16, ""), modbusPlannedPoint("bool", "HR0", schema.DataTypeBool, "")}
				for i := range points {
					points[i].FormatMode = mode
					points[i].DataFormat = "BADC"
				}
				req := modbusPlannedRequest(t, kind, optimizer.ReadableCoverage{}, points...)
				peer := &modbusPlannedPeer{memory: map[modbusPlannedAddress]uint16{{3, 0}: 0xfffe}}
				c := modbusPlannedProtocol(kind, peer)
				got, err := requireModbusPlannedReader(t, c).ReadPlanned(t.Context(), req)
				require.NoError(t, err)
				assertModbusPlannedResults(t, req, got, map[string]any{"signed": int16(-2), "unsigned": uint16(0xfffe), "bool": true})
			})
		}
	}
	for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		t.Run(string(kind)+"frozen-legacy", func(t *testing.T) {
			p := modbusPlannedPoint("legacy", "HR0", schema.DataTypeUint64, "")
			p.FormatMode = schema.PointFormatLegacy
			p.DataFormat = "CDAB"
			req := modbusPlannedRequest(t, kind, optimizer.ReadableCoverage{}, p)
			peer := &modbusPlannedPeer{memory: map[modbusPlannedAddress]uint16{{3, 0}: 0x0708, {3, 1}: 0x0506, {3, 2}: 0x0304, {3, 3}: 0x0102}}
			c := modbusPlannedProtocol(kind, peer)
			switch c := c.(type) {
			case *ModbusTCPConnector:
				c.config.DataFormat = "BADC"
			case *ModbusRTUConnector:
				c.config.DataFormat = "BADC"
			case *ModbusUDPConnector:
				c.config.DataFormat = "BADC"
			}
			got, err := requireModbusPlannedReader(t, c).ReadPlanned(t.Context(), req)
			require.NoError(t, err)
			assertModbusPlannedResults(t, req, got, map[string]any{"legacy": uint64(0x0102030405060708)})
		})
	}
}
