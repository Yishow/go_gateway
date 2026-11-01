package collector

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestModernPointWithoutProfileCannotInstall(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	p := &schema.Point{ID: "P", DeviceID: "D", Address: "HR0", DataType: schema.DataTypeUint32, FormatMode: schema.PointFormatInherit, FormatVersion: 1, Enabled: true}
	require.Zero(t, s.AddPointWithObservationLease(p))
	require.NotContains(t, s.pointInfos, p.ID)
}

type readProfileProbe struct {
	acquisitionFactsProtocol
	requests []connector.ReadRequest
	planned  []connector.PlannedReadRequest
}

func (p *readProfileProbe) Read(_ context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	p.requests = append(p.requests, req)
	return connector.ReadResult{Value: uint32(7), Quality: schema.QualityGood}, nil
}
func (p *readProfileProbe) ReadPlanned(_ context.Context, req connector.PlannedReadRequest) ([]connector.PlannedPointResult, error) {
	p.planned = append(p.planned, req)
	results := executionResults(req, time.Now(), nil)
	for i := range results {
		results[i].Value = uint64(7)
	}
	return results, nil
}

func TestReadProfileScheduledAndManualShareInstalledSnapshot(t *testing.T) {
	old, exists := connector.GetFactory(schema.ProtocolModbusTCP)
	probe := &readProfileProbe{}
	connector.Register(schema.ProtocolModbusTCP, func() connector.Protocol { return probe })
	t.Cleanup(func() {
		if exists {
			connector.Register(schema.ProtocolModbusTCP, old)
		} else {
			connector.Unregister(schema.ProtocolModbusTCP)
		}
	})
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	s := NewScheduler(DefaultSchedulerConfig(), manager)
	d := &schema.Device{ID: "D", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{}`}
	s.AddDevice(d)
	p := &schema.Point{ID: "P", DeviceID: d.ID, RevisionID: "point-1", Address: "HR100", Function: "03", DataType: schema.DataTypeUint64, FormatMode: schema.PointFormatInherit, FormatVersion: 1, Enabled: true}
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, RuleOwned: true, RuleFormat: "DCBA", RuleRevision: "rule-1"})
	require.NoError(t, err)
	lease := s.AddPointWithObservationLease(p, profile)
	require.NotZero(t, lease)
	installed := s.pointInfos[p.ID]
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "ABCD"
	p.RevisionID = "point-2"
	values := s.PollNowContext(t.Context(), []string{p.ID})
	require.Len(t, values, 1)
	require.Equal(t, schema.QualityGood, values[0].Quality, values[0].Error)
	conn, err := manager.GetOrCreate(t.Context(), d.ID, d.Protocol, d.ConnectionConfig)
	require.NoError(t, err)
	require.NoError(t, s.pollPointWithAcquisition(t.Context(), conn, installed, "scheduled"))
	require.Len(t, probe.requests, 1)
	require.Len(t, probe.planned, 1)
	request, bindings, err := probe.planned[0].Snapshot()
	require.NoError(t, err)
	require.Equal(t, 4, request.Length)
	require.Len(t, bindings, 1)
	require.Equal(t, profile, bindings[0].Profile)
	require.Equal(t, "point-1", bindings[0].Profile.Resolution().PointRevision)
	for _, req := range probe.requests {
		require.Equal(t, profile, req.Profile)
		require.Equal(t, 4, req.Count)
		require.Equal(t, "point-1", req.Profile.Resolution().PointRevision)
	}
	require.Equal(t, lease, s.CapturePointLease(p.ID, "point-1"))
	require.False(t, s.BindPointObservation(p.ID, lease, PointObservationBinding{WorkspaceID: "W", WorkspaceRevision: "w", DeviceID: d.ID, PointID: p.ID, RuleID: "R", RuleRevision: "rule-1", ConfigFingerprint: "legacy-only-proof"}))
}

func TestReadProfileModernInstallRefusesExhaustedLease(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	s.AddDevice(&schema.Device{ID: "D", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{}`})
	p := &schema.Point{ID: "P", DeviceID: "D", Address: "HR0", DataType: schema.DataTypeUint32, FormatMode: schema.PointFormatOverride, FormatOverride: "CDAB", Enabled: true}
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: schema.ProtocolModbusTCP})
	require.NoError(t, err)
	s.pointObservationLease = ^uint64(0)
	require.Zero(t, s.AddPointWithObservationLease(p, profile))
	require.NotContains(t, s.pointInfos, p.ID)
}
