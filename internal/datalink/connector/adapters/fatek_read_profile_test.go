package adapters

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"strconv"
	"sync"
	"testing"
)

type fatekProfileTransport struct {
	mu          sync.Mutex
	payload     string
	badChecksum bool
	status      string
	counts      []int
	addresses   []string
}

func (*fatekProfileTransport) Connect() error { return nil }
func (*fatekProfileTransport) Close() error   { return nil }
func (tr *fatekProfileTransport) SendReceive(req []byte) ([]byte, error) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	count, err := strconv.ParseUint(string(req[5:7]), 16, 8)
	if err != nil {
		return nil, err
	}
	tr.counts = append(tr.counts, int(count))
	tr.addresses = append(tr.addresses, string(req[7:len(req)-3]))
	status := tr.status
	if status == "" {
		status = "0"
	}
	response := fatek.BuildFrame(1, string(req[3:5]), status+tr.payload)
	if tr.badChecksum {
		response[len(response)-3] = '?'
	}
	return response, nil
}
func fatekProfileConnector(tr *fatekProfileTransport) *FatekConnector {
	return &FatekConnector{client: fatek.NewClient(tr, 1), transport: tr, connected: true, persistentMode: true}
}
func fatekProfileRequest(t *testing.T, address string, dtype schema.DataType, format string) connector.ReadRequest {
	t.Helper()
	p := schema.Point{ID: "P", DeviceID: "F", Address: address, DataType: dtype, DataFormat: "ABCD", FormatMode: schema.PointFormatOverride, FormatOverride: format, FormatVersion: 1}
	profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolFatekFBs})
	require.NoError(t, err)
	count, err := connector.PointReadCountForAddress(schema.ProtocolFatekFBs, address, dtype)
	require.NoError(t, err)
	return connector.ReadRequest{Address: address, DataType: dtype, DataFormat: p.DataFormat, Count: count, Profile: profile}
}
func TestFatekResolvedProfileOverridesNativeAssembler(t *testing.T) {
	tr := &fatekProfileTransport{payload: "12345678"}
	got, err := fatekProfileConnector(tr).Read(t.Context(), fatekProfileRequest(t, "D100", schema.DataTypeUint32, "CDAB"))
	require.NoError(t, err)
	require.Equal(t, uint32(0x56781234), got.Value)
	require.Equal(t, []int{2}, tr.counts)
}
