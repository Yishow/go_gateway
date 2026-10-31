package adapters

import (
	"encoding/binary"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/lib/hsllogic"
)

type mcPlannedTransport struct {
	mu         sync.Mutex
	frames     [][]byte
	wordMemory map[mcPlannedWordAddress]uint16
	response   []byte
	err        error
	hook       func()
}

func (*mcPlannedTransport) Connect() error { return nil }
func (*mcPlannedTransport) Close() error   { return nil }
func (tr *mcPlannedTransport) SendReceive(frame []byte) ([]byte, error) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.frames = append(tr.frames, append([]byte(nil), frame...))
	if len(frame) >= 15 && binary.LittleEndian.Uint16(frame[11:]) == 0x0403 {
		addresses, err := mcPlannedRandomAddresses(frame)
		if err != nil {
			return nil, err
		}
		if tr.wordMemory != nil {
			words := make([]uint16, len(addresses))
			for i, address := range addresses {
				value, ok := tr.wordMemory[address]
				if !ok {
					return nil, fmt.Errorf("peer has no word at %02x:%d", address.code, address.address)
				}
				words[i] = value
			}
			return mcPlannedWords(words...), nil
		}
	}

	if tr.hook != nil {
		tr.hook()
	}
	return append([]byte(nil), tr.response...), tr.err
}
func mcPlannedConnector(tr *mcPlannedTransport) *MC3EConnector {
	return &MC3EConnector{client: mcprotocol.NewClientWithTransport(tr), connected: true, persistentMode: true, dataConverter: hsllogic.NewDataConverter(hsllogic.DataFormatDCBA)}
}
func mcPlannedWords(words ...uint16) []byte {
	raw := make([]byte, 2+2*len(words))
	for i, w := range words {
		binary.LittleEndian.PutUint16(raw[2+2*i:], w)
	}
	return raw
}
func mcPlannedPoint(id, address string, dtype schema.DataType, format string) schema.Point {
	p := schema.Point{ID: id, DeviceID: "device", RevisionID: "point-revision-" + id, Address: address, DataType: dtype, FormatMode: schema.PointFormatInherit, FormatVersion: 1}
	if format != "" {
		p.FormatMode = schema.PointFormatOverride
		p.FormatOverride = format
	}
	return p
}
func mcPlannedRequest(t *testing.T, points ...schema.Point) connector.PlannedReadRequest {
	t.Helper()
	return mcPlannedRequestCapability(t, "supported", points...)
}
func mcPlannedRequestCapability(t *testing.T, support string, points ...schema.Point) connector.PlannedReadRequest {
	t.Helper()
	source := optimizer.AcquisitionSource{Identity: optimizer.AcquisitionIdentity{DeviceID: "device", DeviceRevision: "device-rev", Station: "0", ConnectionRevision: connector.ReadConnectionRevision("{}"), Protocol: "mc", Frame: "3e-binary", ProtocolRevision: "mc-v1", PollingGroup: "fast", Cadence: time.Second}, Capability: optimizer.AcquisitionCapability{Revision: "cap-v1", RandomSupport: support, ProtocolLimit: 255}}
	req := connector.PlannedReadRequest{AcquisitionID: "scan-1", RequestID: "request-1"}
	for _, p := range points {
		profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolMC3E, RuleOwned: true, RuleID: "rule", RuleRevision: "source-rev", ConnectionFormat: "CDAB", ConnectionRevision: connector.ReadConnectionRevision("{}")})
		require.NoError(t, err)
		r := profile.Resolution()
		point := optimizer.EffectivePoint{ID: p.ID, SourceIdentity: "rule", SourceRevision: "source-rev", Address: p.Address, Function: p.Function, DataType: p.DataType, NativeWidth: schema.RegisterCountForDataType(p.DataType), Format: r.EffectiveFormat, FormatSource: string(r.Source), FormatRevision: p.RevisionID}
		source.Points = append(source.Points, point)
		req.Bindings = append(req.Bindings, connector.PlannedPointBinding{Point: point, Profile: profile, PointRevision: p.RevisionID, ConfigFingerprint: "fingerprint-" + p.ID, InstallationScope: "installation", InstallationLease: 19})
	}
	req.Plan = optimizer.BuildAcquisitionPlan(optimizer.AcquisitionInput{ConfigurationRevision: "config-v1", Sources: []optimizer.AcquisitionSource{source}})
	require.Empty(t, req.Plan.Errors())
	return req
}
func assertMCPlannedResults(t *testing.T, req connector.PlannedReadRequest, results []connector.PlannedPointResult, want map[string]any) {
	t.Helper()
	require.Len(t, results, len(want))
	seen := map[string]bool{}
	for _, r := range results {
		id := r.Binding.Point.ID
		require.False(t, seen[id])
		seen[id] = true
		require.Equal(t, want[id], r.Value, id)
		require.Equal(t, schema.QualityGood, r.Quality, id)
		require.Equal(t, req.AcquisitionID, r.AcquisitionID)
		require.Equal(t, req.RequestID, r.RequestID)
		require.Equal(t, req.Plan.Fingerprint(), r.PlanFingerprint)
		require.Equal(t, req.Plan.Requests()[req.RequestIndex].Group, r.Group)
		for _, b := range req.Bindings {
			if b.Point.ID == id {
				require.Equal(t, b, r.Binding)
			}
		}
		require.False(t, r.Timestamp.IsZero())
		require.Equal(t, results[0].Timestamp, r.Timestamp)
	}
}
func TestMCPlannedContiguousOverlapAndPerPointDecode(t *testing.T) {
	req := mcPlannedRequest(t, mcPlannedPoint("wide", "D0", schema.DataTypeInt64, "ABCD"), mcPlannedPoint("small", "D1", schema.DataTypeInt16, ""), mcPlannedPoint("abcd", "D0", schema.DataTypeUint32, "ABCD"), mcPlannedPoint("cdab", "D0", schema.DataTypeUint32, "CDAB"))
	require.Len(t, req.Plan.Requests(), 1)
	require.Equal(t, "contiguous", req.Plan.Requests()[0].Kind)
	tr := &mcPlannedTransport{response: mcPlannedWords(0x0102, 0x0304, 0x0506, 0x0708)}
	c := mcPlannedConnector(tr)
	results, err := c.ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	assertMCPlannedResults(t, req, results, map[string]any{"wide": int64(0x0102030405060708), "small": int16(0x0304), "abcd": uint32(0x01020304), "cdab": uint32(0x03040102)})
	require.Len(t, tr.frames, 1)
	require.Len(t, tr.frames[0], 21)
	require.Equal(t, []byte{1, 4, 0, 0, 0, 0, 0, 0xa8, 4, 0}, tr.frames[0][11:])
	require.Equal(t, uint32(0x04030201), c.dataConverter.RegistersToValue([]uint16{0x0102, 0x0304}, hsllogic.DataTypeUint32), "shared converter was mutated")
}
func TestMCPlannedRandomSparseDeduplicatesWholeScalars(t *testing.T) {
	req := mcPlannedRequest(t, mcPlannedPoint("a", "D0", schema.DataTypeUint32, "ABCD"), mcPlannedPoint("b", "D0", schema.DataTypeUint32, "CDAB"), mcPlannedPoint("c", "D10", schema.DataTypeInt16, ""))
	require.Len(t, req.Plan.Requests(), 1)
	require.Equal(t, "random", req.Plan.Requests()[0].Kind)
	tr := &mcPlannedTransport{response: mcPlannedWords(0x0102, 0x0304, 0x1234)}
	results, err := mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	assertMCPlannedResults(t, req, results, map[string]any{"a": uint32(0x01020304), "b": uint32(0x03040102), "c": int16(0x1234)})
	require.Len(t, tr.frames, 1)
	require.Len(t, tr.frames[0], 29)
	require.Equal(t, []byte{3, 4, 0, 0, 3, 0, 0, 0, 0, 0xa8, 1, 0, 0, 0xa8, 10, 0, 0, 0xa8}, tr.frames[0][11:])
}

func TestMCPlannedThreeSparseWordsUseOneRandomRequest(t *testing.T) {
	req := mcPlannedRequest(t, mcPlannedPoint("a", "D0", schema.DataTypeInt16, ""), mcPlannedPoint("b", "D10", schema.DataTypeInt16, ""), mcPlannedPoint("c", "D20", schema.DataTypeInt16, ""))
	require.Len(t, req.Plan.Requests(), 1)
	require.Equal(t, "random", req.Plan.Requests()[0].Kind)
	tr := &mcPlannedTransport{wordMemory: map[mcPlannedWordAddress]uint16{{0xa8, 0}: 11, {0xa8, 10}: 22, {0xa8, 20}: 33}}
	got, err := mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	assertMCPlannedResults(t, req, got, map[string]any{"a": int16(11), "b": int16(22), "c": int16(33)})
	require.Len(t, tr.frames, 1)
	require.Len(t, tr.frames[0], 29)
	require.Equal(t, []byte{3, 4, 0, 0, 3, 0, 0, 0, 0, 0xa8, 10, 0, 0, 0xa8, 20, 0, 0, 0xa8}, tr.frames[0][11:])
}
