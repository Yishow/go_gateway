package collector

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestModernObservationProofRejectsReinstallAndRestart(t *testing.T) {
	observed := []AcquisitionObservation{}
	cfg := DefaultSchedulerConfig()
	cfg.AcquisitionObserver = func(o AcquisitionObservation) { observed = append(observed, o) }
	s := NewScheduler(cfg, nil)
	d := &schema.Device{ID: "D", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{"data_format":"CDAB"}`}
	p := &schema.Point{ID: "P", DeviceID: d.ID, RevisionID: "same", Address: "D0", DataType: schema.DataTypeUint32, FormatMode: schema.PointFormatOverride, FormatOverride: "ABCD", FormatVersion: 1, Enabled: true}
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, RuleID: "R", RuleOwned: true, RuleRevision: "r1", ConnectionRevision: connector.ReadConnectionRevision(d.ConnectionConfig)})
	require.NoError(t, err)
	s.AddDevice(d)
	lease := s.AddPointWithObservationLease(p, profile)
	require.NotZero(t, lease)
	fp, err := measurement.AcquisitionConfigFingerprint(*d, *p, profile)
	require.NoError(t, err)
	binding := PointObservationBinding{WorkspaceID: "W", WorkspaceRevision: "w1", DeviceID: d.ID, PointID: p.ID, RuleID: "R", RuleRevision: "r1", ConfigFingerprint: fp, InstallationScope: s.InstallationScope()}
	require.True(t, s.BindPointObservation(p.ID, lease, binding))
	original := s.pointInfos[p.ID]
	value := s.collectedValue(deviceConfig{ID: d.ID, Protocol: d.Protocol, Config: d.ConnectionConfig}, original, connector.ReadResult{Quality: schema.QualityGood}, nil, "old", s.now())
	require.Equal(t, fp, value.ConfigFingerprint)
	require.Equal(t, lease, value.InstallationLease)
	require.Equal(t, s.InstallationScope(), value.InstallationScope)
	observed = nil
	next := s.AddPointWithObservationLease(p, profile)
	require.NotEqual(t, lease, next)
	require.False(t, s.BindPointObservation(p.ID, lease, binding))
	s.observeAcquisition(original, value)
	require.Empty(t, observed)
	require.False(t, s.AcceptsInstallationSample(p.ID, lease, value.InstallationScope))
	restart := NewScheduler(cfg, nil)
	restart.AddDevice(d)
	newLease := restart.AddPointWithObservationLease(p, profile)
	require.Equal(t, lease, newLease, "restart may reuse numeric lease but never incarnation")
	observed = nil
	restart.observeAcquisition(original, value)
	require.Empty(t, observed, "captured old install cannot borrow the new scheduler scope")
	require.False(t, restart.BindPointObservation(p.ID, newLease, binding))
	require.False(t, restart.AcceptsInstallationSample(p.ID, newLease, value.InstallationScope))
	binding.InstallationScope = restart.InstallationScope()
	require.True(t, restart.BindPointObservation(p.ID, newLease, binding))
}

func TestLegacyObservationCannotRecoverModernReplacement(t *testing.T) {
	observed := []AcquisitionObservation{}
	cfg := DefaultSchedulerConfig()
	cfg.AcquisitionObserver = func(o AcquisitionObservation) { observed = append(observed, o) }
	s := NewScheduler(cfg, nil)
	d := &schema.Device{ID: "D", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{}`}
	s.AddDevice(d)
	p := &schema.Point{ID: "P", DeviceID: d.ID, Address: "D0", DataType: schema.DataTypeUint32, Enabled: true}
	require.NotZero(t, s.AddPointWithObservationLease(p))
	original := s.pointInfos[p.ID]
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "ABCD"
	p.FormatVersion = 1
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol})
	require.NoError(t, err)
	require.NotZero(t, s.AddPointWithObservationLease(p, profile))
	s.observeAcquisition(original, CollectedValue{PointID: p.ID, DeviceID: d.ID, Quality: schema.QualityGood})
	require.Empty(t, observed, "old legacy acquisition must not recover the modern point")
}
