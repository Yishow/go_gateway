package adapters

import (
	"encoding/binary"
	"fmt"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/lib/hsllogic"
	"sync"
	"testing"
)

type mcProfileTransport struct {
	mu       sync.Mutex
	words    []uint16
	extra    []byte
	requests []int
}

func (*mcProfileTransport) Connect() error { return nil }
func (*mcProfileTransport) Close() error   { return nil }
func (tr *mcProfileTransport) SendReceive(req []byte) ([]byte, error) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.requests = append(tr.requests, int(binary.LittleEndian.Uint16(req[len(req)-2:])))
	response := make([]byte, 2+len(tr.words)*2)
	for i, w := range tr.words {
		binary.LittleEndian.PutUint16(response[2+i*2:], w)
	}
	return append(response, tr.extra...), nil
}
func mcProfileConnector(tr *mcProfileTransport) *MC3EConnector {
	return &MC3EConnector{client: mcprotocol.NewClientWithTransport(tr), connected: true, persistentMode: true, dataConverter: hsllogic.NewDataConverter(hsllogic.DataFormatCDAB)}
}
func mcProfileRequest(t *testing.T, dtype schema.DataType, format string) connector.ReadRequest {
	t.Helper()
	p := schema.Point{ID: "P", DeviceID: "D", Address: "D100", DataType: dtype, DataFormat: "CDAB", FormatMode: schema.PointFormatOverride, FormatOverride: format, FormatVersion: 1}
	profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolMC3E, ConnectionFormat: "CDAB"})
	require.NoError(t, err)
	return connector.ReadRequest{Address: p.Address, DataType: dtype, DataFormat: p.DataFormat, Count: 1, Profile: profile}
}
func TestMC3EResolvedProfileOverridesConnectionConverter(t *testing.T) {
	tr := &mcProfileTransport{words: []uint16{0x1234, 0x5678}}
	c := mcProfileConnector(tr)
	result, err := c.Read(t.Context(), mcProfileRequest(t, schema.DataTypeUint32, "ABCD"))
	require.NoError(t, err)
	require.Equal(t, uint32(0x12345678), result.Value)
	require.Equal(t, []int{2}, tr.requests)
}

func TestMC3EResolvedProfileRejectsExtraPayload(t *testing.T) {
	for _, extra := range [][]byte{{0}, {0, 0}} {
		t.Run(fmt.Sprint(len(extra)), func(t *testing.T) {
			tr := &mcProfileTransport{words: []uint16{0x1234, 0x5678}, extra: extra}
			result, err := mcProfileConnector(tr).Read(t.Context(), mcProfileRequest(t, schema.DataTypeUint32, "CDAB"))
			require.Error(t, err)
			require.Equal(t, schema.QualityBad, result.Quality)
			require.Nil(t, result.Value)
			require.Len(t, tr.requests, 1, "malformed payload is not a connection retry")
		})
	}
}

func TestMC3EResolvedProfileRejectsUnknownModernType(t *testing.T) {
	p := schema.Point{ID: "P", DeviceID: "D", Address: "D100", DataType: schema.DataType("unknown"), FormatMode: schema.PointFormatInherit, FormatVersion: 1}
	profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolMC3E})
	require.NoError(t, err)
	tr := &mcProfileTransport{words: []uint16{1}}
	got, err := mcProfileConnector(tr).Read(t.Context(), connector.ReadRequest{Address: p.Address, DataType: p.DataType, Profile: profile})
	require.Error(t, err)
	require.Equal(t, schema.QualityBad, got.Quality)
	require.Empty(t, tr.requests)
}

func TestMC3EResolvedProfileRejectsForeignProtocolBeforeRead(t *testing.T) {
	p := schema.Point{ID: "P", DeviceID: "D", Address: "D100", DataType: schema.DataTypeUint32, FormatMode: schema.PointFormatOverride, FormatOverride: "ABCD", FormatVersion: 1}
	profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolFatekFBs})
	require.NoError(t, err)
	tr := &mcProfileTransport{words: []uint16{1, 2}}
	got, err := mcProfileConnector(tr).Read(t.Context(), connector.ReadRequest{Address: p.Address, DataType: p.DataType, Profile: profile})
	require.ErrorIs(t, err, connector.ErrReadProfile)
	require.Equal(t, schema.QualityBad, got.Quality)
	require.Empty(t, tr.requests)
}
