package adapters

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
	"testing"
)

func TestModbusResolvedProfileOverridesLegacyProjection(t *testing.T) {
	transport := NewMockTransport()
	transport.HoldingRegs[100] = 0x1234
	transport.HoldingRegs[101] = 0x5678
	c := &ModbusTCPConnector{client: modbus.NewClient(transport, 1), connected: true, persistentMode: true}
	p := schema.Point{ID: "P", Address: "HR100", DataType: schema.DataTypeUint32, DataFormat: "ABCD", FormatMode: schema.PointFormatOverride, FormatOverride: "CDAB", FormatVersion: 1}
	profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolModbusTCP})
	require.NoError(t, err)
	result, err := c.Read(t.Context(), connector.ReadRequest{Address: p.Address, DataType: p.DataType, DataFormat: p.DataFormat, Profile: profile})
	require.NoError(t, err)
	require.Equal(t, uint32(0x56781234), result.Value)
}

func TestModbusResolvedProfilePreservesInherit16AndLegacy(t *testing.T) {
	for _, mode := range []schema.PointFormatMode{schema.PointFormatLegacy, schema.PointFormatInherit} {
		for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
			tr := &profileWordTransport{words: []uint16{0x1234}}
			c := profileProtocol(schema.ProtocolModbusTCP, tr)
			p := schema.Point{ID: "P", Address: "HR100", DataType: schema.DataTypeUint16, DataFormat: format, FormatMode: mode}
			profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolModbusTCP, RuleOwned: true, RuleFormat: format})
			require.NoError(t, err)
			req := connector.ReadRequest{Address: p.Address, DataType: p.DataType, DataFormat: format, Profile: profile}
			result, err := c.Read(t.Context(), req)
			require.NoError(t, err)
			require.Equal(t, convertModbusValue(tr.words, p.DataType, format), result.Value)
		}
	}
}
func TestModbusResolvedProfileRejectsIdentityBeforeWire(t *testing.T) {
	tr := &profileWordTransport{words: []uint16{1, 2}}
	c := profileProtocol(schema.ProtocolModbusTCP, tr)
	req := profileRequest(t, schema.ProtocolModbusTCP, schema.DataTypeUint32, "CDAB")
	req.Address = "HR101"
	result, err := c.Read(t.Context(), req)
	require.ErrorIs(t, err, connector.ErrReadProfile)
	require.NotEqual(t, schema.QualityGood, result.Quality)
	require.Zero(t, tr.sends)
}
