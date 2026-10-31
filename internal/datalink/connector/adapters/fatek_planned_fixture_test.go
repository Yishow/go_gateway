package adapters

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

func fatekPlannedPoint(id, address string, dtype schema.DataType, format string) schema.Point {
	p := schema.Point{ID: id, DeviceID: "device", RevisionID: "point-revision-" + id, Address: address, DataType: dtype, FormatMode: schema.PointFormatInherit, FormatVersion: 1}
	if format != "" {
		p.FormatMode = schema.PointFormatOverride
		p.FormatOverride = format
	}
	return p
}
func fatekPlannedRequest(t *testing.T, points ...schema.Point) connector.PlannedReadRequest {
	t.Helper()
	source := optimizer.AcquisitionSource{Identity: optimizer.AcquisitionIdentity{DeviceID: "device", DeviceRevision: "device-rev", Station: "1", ConnectionRevision: connector.ReadConnectionRevision("{}"), Protocol: "fatek", Frame: "ascii", ProtocolRevision: "fatek-v1", PollingGroup: "fast", Cadence: time.Second}, Capability: optimizer.AcquisitionCapability{Revision: "cap-v1", RandomSupport: "supported", RandomCrossArea: true, ProtocolLimit: 64}}
	req := connector.PlannedReadRequest{AcquisitionID: "scan-1", RequestID: "request-1"}
	for _, p := range points {
		profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolFatekFBs, RuleOwned: true, RuleID: "opaque-rule", RuleRevision: "opaque-source-proof", ConnectionRevision: connector.ReadConnectionRevision("{}")})
		require.NoError(t, err)
		r := profile.Resolution()
		width := schema.RegisterCountForDataType(p.DataType)
		if strings.HasPrefix(p.Address, "DR") || strings.HasPrefix(p.Address, "DD") || strings.HasPrefix(p.Address, "DW") {
			width = (width + 1) / 2
		}
		point := optimizer.EffectivePoint{ID: p.ID, SourceIdentity: "opaque-rule", SourceRevision: "opaque-source-proof", Address: p.Address, Function: p.Function, DataType: p.DataType, NativeWidth: width, Format: r.EffectiveFormat, FormatSource: string(r.Source), FormatRevision: p.RevisionID}
		source.Points = append(source.Points, point)
		req.Bindings = append(req.Bindings, connector.PlannedPointBinding{Point: point, Profile: profile, PointRevision: p.RevisionID, ConfigFingerprint: "fingerprint-" + p.ID, InstallationScope: "installation", InstallationLease: 19})
	}
	req.Plan = optimizer.BuildAcquisitionPlan(optimizer.AcquisitionInput{ConfigurationRevision: "config-v1", Sources: []optimizer.AcquisitionSource{source}})
	require.Empty(t, req.Plan.Errors())
	return req
}
func readFatekPlanned(t *testing.T, ctx context.Context, c *FatekConnector, req connector.PlannedReadRequest) ([]connector.PlannedPointResult, error) {
	t.Helper()
	reader, ok := any(c).(connector.PlannedReader)
	require.True(t, ok, "formal FATEK adapter must implement PlannedReader")
	return reader.ReadPlanned(ctx, req)
}
func assertFatekPlannedResults(t *testing.T, req connector.PlannedReadRequest, results []connector.PlannedPointResult, want map[string]any) {
	t.Helper()
	require.Len(t, results, len(want))
	seen := map[string]bool{}
	for _, r := range results {
		id := r.Binding.Point.ID
		require.False(t, seen[id])
		seen[id] = true
		expected, ok := want[id]
		require.True(t, ok)
		require.Equal(t, expected, r.Value, id)
		require.Equal(t, schema.QualityGood, r.Quality, id)
		require.Empty(t, r.Error)
		require.Equal(t, req.AcquisitionID, r.AcquisitionID)
		require.Equal(t, req.RequestID, r.RequestID)
		require.Equal(t, req.Plan.Fingerprint(), r.PlanFingerprint)
		request, ok := req.Plan.Request(req.RequestIndex)
		require.True(t, ok)
		require.Equal(t, request.Group, r.Group)
		found := false
		for _, b := range req.Bindings {
			if b.Point.ID == id {
				require.Equal(t, b, r.Binding)
				found = true
			}
		}
		require.True(t, found)
		require.False(t, r.Timestamp.IsZero())
		require.Equal(t, results[0].Timestamp, r.Timestamp)
		require.Equal(t, r.Timestamp, r.ObservedAt)
		require.Equal(t, r.Timestamp, r.ReceivedAt)
		require.Equal(t, "gateway", r.TimeOrigin)
	}
}
func TestFatekPlannedThreeSparseWordsUseOneRandomRequest(t *testing.T) {
	req := fatekPlannedRequest(t, fatekPlannedPoint("a", "D0", schema.DataTypeInt16, ""), fatekPlannedPoint("b", "D10", schema.DataTypeInt16, ""), fatekPlannedPoint("c", "D20", schema.DataTypeInt16, ""))
	require.Len(t, req.Plan.Requests(), 1)
	require.Equal(t, "random", req.Plan.Requests()[0].Kind)
	memory := map[string]any{"D0": 11, "D10": 22, "D20": 33}
	baseline := &fatekPlannedPeer{memory: memory}
	for _, b := range req.Bindings {
		got, err := fatekPlannedConnector(baseline).Read(t.Context(), connector.ReadRequest{Address: b.Point.Address, DataType: b.Point.DataType, Count: 1, Profile: b.Profile})
		require.NoError(t, err)
		require.Equal(t, schema.QualityGood, got.Quality)
	}
	require.Len(t, baseline.frames, 3)
	peer := &fatekPlannedPeer{memory: memory}
	got, err := readFatekPlanned(t, t.Context(), fatekPlannedConnector(peer), req)
	require.NoError(t, err)
	assertFatekPlannedResults(t, req, got, map[string]any{"a": int16(11), "b": int16(22), "c": int16(33)})
	require.Len(t, peer.frames, 1)
	require.Equal(t, "\x02014803D00000D00010D00020D1\x03", string(peer.frames[0]))
	require.Equal(t, 28, len(peer.frames[0]))
	require.Equal(t, 21, len(peer.responses[0]))
	require.Equal(t, optimizer.ApplicationFrameBudget{RequestBytes: 28, ResponseBytes: 21}, req.Plan.Requests()[0].FrameBytes)
}
