package adapters

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"go-gateway/lib/hsllogic"

	"github.com/stretchr/testify/require"
)

// This checks the new policy description against real Read fixtures. It does
// not claim database migration, collector installation, or adapter rewiring.
func TestPLCLegacyResolver_ReconstructionAndUnrelatedEdit(t *testing.T) {
	for _, tc := range []struct {
		protocol                      schema.ProtocolType
		pointFormat, connectionFormat string
		want                          uint32
	}{
		{schema.ProtocolModbusTCP, "", "CDAB", 0x03040102},
		{schema.ProtocolModbusRTU, "ABCD", "CDAB", 0x01020304},
		{schema.ProtocolModbusUDP, "ABCD", "CDAB", 0x01020304},
		{schema.ProtocolMC3E, "ABCD", "CDAB", 0x03040102},
		{schema.ProtocolFatekFBs, "CDAB", "DCBA", 0x01020304},
	} {
		t.Run(string(tc.protocol), func(t *testing.T) {
			point := schema.Point{ID: "legacy-point", Name: "before", DataType: schema.DataTypeUint32, DataFormat: tc.pointFormat}
			var original decoding.Resolution
			for _, phase := range []string{"upgrade", "reconstruct", "unrelated-edit"} {
				if phase == "reconstruct" {
					raw, err := json.Marshal(point)
					require.NoError(t, err)
					var reloaded schema.Point
					require.NoError(t, json.Unmarshal(raw, &reloaded))
					point = reloaded
				}
				if phase == "unrelated-edit" {
					point.Name = "renamed"
					point.Description = "changed description"
				}
				resolved, err := decoding.ResolveLegacy(decoding.LegacyInput{
					Protocol: tc.protocol, PointFormat: point.DataFormat, ConnectionFormat: tc.connectionFormat,
				})
				require.NoError(t, err)
				if phase == "upgrade" {
					original = resolved
				} else {
					require.Equal(t, original, resolved)
				}
				var protocol connector.Protocol
				request := connector.ReadRequest{Address: "HR100", DataType: point.DataType, DataFormat: point.DataFormat, Count: 2}
				switch tc.protocol {
				case schema.ProtocolMC3E:
					protocol = newPLCAdapterMCConnector(hsllogic.DataFormat(tc.connectionFormat), &plcAdapterMCTransport{words: []uint16{0x0102, 0x0304}})
					request.Address, request.Count = "D0", 1
				case schema.ProtocolFatekFBs:
					protocol, _ = newPLCFatekReadFixture("01020304")
					request.Address = "D0"
					require.True(t, resolved.NativeCompatibility)
					require.Empty(t, resolved.EffectiveFormat)
				default:
					adapter := map[schema.ProtocolType]string{schema.ProtocolModbusTCP: "tcp", schema.ProtocolModbusRTU: "rtu", schema.ProtocolModbusUDP: "udp"}[tc.protocol]
					protocol = newPLCAdapterModbusConnector(adapter, newPLCAdapterModbusTransport([]uint16{0x0102, 0x0304}))
					loadConnectionConfigDataFormat(t, protocol, tc.connectionFormat)
				}
				result, err := protocol.Read(t.Context(), request)
				require.NoError(t, err)
				require.Equal(t, schema.QualityGood, result.Quality)
				require.Equal(t, tc.want, result.Value)
				if !resolved.NativeCompatibility {
					// Metadata describes, rather than changes, the existing adapter value.
					value := hsllogic.NewDataConverter(hsllogic.DataFormat(resolved.EffectiveFormat)).RegistersToValue([]uint16{0x0102, 0x0304}, hsllogic.DataTypeUint32)
					require.Equal(t, tc.want, value)
				}
				require.Equal(t, tc.pointFormat, point.DataFormat, "legacy metadata must not become an inferred override")
			}
		})
	}
}

// Previously accepted connection JSON can contain arbitrary optional formats.
// Exercise Read before the resolver, so its fallback is independent evidence.
func TestPLCLegacyResolver_ExistingNoncanonicalConnectionFormats(t *testing.T) {
	for _, protocolType := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP, schema.ProtocolMC3E} {
		for _, format := range []string{"unknown", "cdab", " CDAB ", " \t"} {
			t.Run(string(protocolType)+"/"+format, func(t *testing.T) {
				var protocol connector.Protocol
				request := connector.ReadRequest{Address: "HR100", DataType: schema.DataTypeUint32, Count: 2}
				want := uint32(0x01020304)
				effective := "ABCD"
				if protocolType == schema.ProtocolMC3E {
					protocol = newPLCAdapterMCConnector(hsllogic.DataFormat(format), &plcAdapterMCTransport{words: []uint16{0x0102, 0x0304}})
					request.Address, request.Count, request.DataFormat = "D0", 1, "ABCD"
				} else {
					adapter := map[schema.ProtocolType]string{schema.ProtocolModbusTCP: "tcp", schema.ProtocolModbusRTU: "rtu", schema.ProtocolModbusUDP: "udp"}[protocolType]
					protocol = newPLCAdapterModbusConnector(adapter, newPLCAdapterModbusTransport([]uint16{0x0102, 0x0304}))
					loadConnectionConfigDataFormat(t, protocol, format)
					if format == "cdab" || format == " CDAB " {
						want, effective = 0x03040102, "CDAB"
					}
				}
				result, err := protocol.Read(t.Context(), request)
				require.NoError(t, err)
				require.Equal(t, schema.QualityGood, result.Quality)
				require.Equal(t, want, result.Value)
				resolved, err := decoding.ResolveLegacy(decoding.LegacyInput{Protocol: protocolType, PointFormat: request.DataFormat, ConnectionFormat: format})
				require.NoError(t, err, "legacy metadata must describe accepted decoder fallback")
				require.Equal(t, effective, resolved.EffectiveFormat)
				require.Equal(t, want, hsllogic.NewDataConverter(hsllogic.DataFormat(resolved.EffectiveFormat)).RegistersToValue([]uint16{0x0102, 0x0304}, hsllogic.DataTypeUint32))
			})
		}
	}
}
