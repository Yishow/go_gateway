package adapters

import (
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/lib/hsllogic"

	"github.com/stretchr/testify/require"
)

// Reconstruct the real adapters to protect restart behavior without persisting
// a resolved connection value as a point override.
func TestPLCLegacyPolicy_ReconstructedMCConnectionIgnoresPoint(t *testing.T) {
	for _, stored := range []string{"ABCD", "DCBA", "ignored-invalid"} {
		for restart := 0; restart < 2; restart++ {
			transport := &plcAdapterMCTransport{words: []uint16{0x0102, 0x0304}}
			protocol := newPLCAdapterMCConnector(hsllogic.DataFormatCDAB, transport)
			result, err := protocol.Read(t.Context(), connector.ReadRequest{
				Address: "D0", DataType: schema.DataTypeUint32, DataFormat: stored,
			})
			require.NoError(t, err)
			require.Equal(t, schema.QualityGood, result.Quality)
			require.Equal(t, uint32(0x03040102), result.Value)
		}
	}
}

func TestPLCLegacyPolicy_FatekIgnoredPointFormat(t *testing.T) {
	for _, address := range []string{"D0", "DR0"} {
		for _, format := range []string{"ABCD", "CDAB", "DCBA", "ignored-invalid"} {
			protocol, transport := newPLCFatekReadFixture("01020304")
			count := 2
			if address == "DR0" {
				count = 1
			}
			result, err := protocol.Read(t.Context(), connector.ReadRequest{
				Address: address, DataType: schema.DataTypeUint32, DataFormat: format, Count: count,
			})
			require.NoError(t, err)
			require.Equal(t, schema.QualityGood, result.Quality)
			require.Equal(t, uint32(0x01020304), result.Value)
			require.Equal(t, count, plcFatekRequestCount(t, transport))
		}
	}
}

func TestPLCLegacyPolicy_SingleWordRegisterValues(t *testing.T) {
	for _, format := range []string{"BADC", "DCBA"} {
		for _, adapter := range []string{"tcp", "rtu", "udp", "mc3e"} {
			for _, dataType := range []schema.DataType{schema.DataTypeUint16, schema.DataTypeInt16} {
				var protocol connector.Protocol
				if adapter == "mc3e" {
					protocol = newPLCAdapterMCConnector(hsllogic.DataFormat(format), &plcAdapterMCTransport{words: []uint16{0x0102}})
				} else {
					protocol = newPLCAdapterModbusConnector(adapter, newPLCAdapterModbusTransport([]uint16{0x0102}))
					loadConnectionConfigDataFormat(t, protocol, format)
				}
				address := "HR100"
				if adapter == "mc3e" {
					address = "D0"
				}
				result, err := protocol.Read(t.Context(), connector.ReadRequest{Address: address, DataType: dataType})
				require.NoError(t, err)
				require.Equal(t, schema.QualityGood, result.Quality)
				if dataType == schema.DataTypeUint16 {
					require.Equal(t, uint16(0x0102), result.Value)
				} else {
					require.Equal(t, int16(0x0102), result.Value)
				}
			}
		}
	}
}

// Byte-oriented conversion is a separate legacy path: unlike register reads,
// it applies BADC/DCBA byte swaps. Neither behavior may be unified silently.
func TestPLCLegacyPolicy_ByteTransformSingleWordSwaps(t *testing.T) {
	for _, format := range []hsllogic.DataFormat{hsllogic.DataFormatBADC, hsllogic.DataFormatDCBA} {
		require.Equal(t, uint16(0x0201), hsllogic.NewByteTransform(format).TransformUint16([]byte{0x01, 0x02}, 0))
	}
}
