package collector

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestFormatRecoveryStagedSlotBlocksScheduledManualAndCapturedRaw(t *testing.T) {
	config, requests := mcWordFixture(t, []uint16{0x1122, 0x3344})
	connector.Register(schema.ProtocolMC3E, adapters.NewMC3EConnector)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { manager.CloseAll() })
	s := NewScheduler(DefaultSchedulerConfig(), manager)
	d := &schema.Device{ID: "D", Protocol: schema.ProtocolMC3E, ConnectionConfig: config}
	s.AddDevice(d)
	group := "G"
	p := &schema.Point{ID: "P", DeviceID: d.ID, RevisionID: "p", Address: "D0", DataType: schema.DataTypeUint32, FormatMode: schema.PointFormatOverride, FormatOverride: "ABCD", FormatVersion: 1, Enabled: true, PollingGroupID: &group}
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, ConnectionRevision: connector.ReadConnectionRevision(config)})
	require.NoError(t, err)
	old := s.AddPointWithObservationLease(p, profile)
	require.NotZero(t, old)
	proof, ok := s.CaptureOrdinaryPoint(p.ID)
	require.True(t, ok)
	lease, ok := s.StagePointFormatRecovery(p, proof, profile)
	require.True(t, ok)
	s.pollGroup(group)
	require.Empty(t, s.PollNowContext(t.Context(), []string{p.ID}))
	require.Empty(t, requests, "staged scheduled/manual paths must not reach the actual adapter")
	require.False(t, s.AcceptsInstallationSample(p.ID, old, s.InstallationScope()))
	require.False(t, s.AcceptsInstallationSample(p.ID, lease, s.InstallationScope()))
	require.False(t, s.AcceptsDirectReadProof(proof, s.InstallationScope()))
	require.False(t, s.AcceptsLegacyUnscopedSample(p.ID))
	_, ok = s.CaptureReadablePoint(p.ID)
	require.False(t, ok)
	require.True(t, s.ActivatePointFormatRecovery(p.ID, lease, profile))
	s.pollGroup(group)
	got := s.PollNowContext(t.Context(), []string{p.ID})
	require.Len(t, got, 1)
	require.Equal(t, uint32(0x11223344), got[0].Value)
	require.Len(t, requests, 2)
}
