package collector

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionInstallationRequiresCurrentDeviceConfiguration(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolMQTT} {
		t.Run(string(protocol), func(t *testing.T) {
			s := NewScheduler(DefaultSchedulerConfig(), nil)
			point := &schema.Point{ID: "p", DeviceID: "d", Address: "0", DataType: schema.DataTypeUint16}
			lease := s.AddPointWithObservationLease(point)
			require.False(t, s.AcceptsInstallationSample(point.ID, lease, s.InstallationScope()), "missing device cannot establish installation configuration")
			device := &schema.Device{ID: "d", Protocol: protocol, ConnectionConfig: `{}`}
			s.AddDevice(device)
			require.False(t, s.AcceptsInstallationSample(point.ID, lease, s.InstallationScope()), "adding a device cannot bless the earlier unknown installation")
			lease = s.AddPointWithObservationLease(point)
			require.True(t, s.AcceptsInstallationSample(point.ID, lease, s.InstallationScope()))
			s.AddDevice(device)
			require.True(t, s.AcceptsInstallationSample(point.ID, lease, s.InstallationScope()), "same configuration refresh does not manufacture a new lease")
			device.ConnectionConfig = `{"timeout":2}`
			s.AddDevice(device)
			require.False(t, s.AcceptsInstallationSample(point.ID, lease, s.InstallationScope()))
			device.ConnectionConfig = `{}`
			s.AddDevice(device)
			require.False(t, s.AcceptsInstallationSample(point.ID, lease, s.InstallationScope()), "configuration ABA must not revive an old installation")
			lease = s.AddPointWithObservationLease(point)
			require.True(t, s.AcceptsInstallationSample(point.ID, lease, s.InstallationScope()))
			device.Protocol = schema.ProtocolFatekFBs
			s.AddDevice(device)
			require.False(t, s.AcceptsInstallationSample(point.ID, lease, s.InstallationScope()))
			lease = s.AddPointWithObservationLease(point)
			require.True(t, s.AcceptsInstallationSample(point.ID, lease, s.InstallationScope()))
			s.RemoveDevice(device.ID)
			s.AddDevice(device)
			require.False(t, s.AcceptsInstallationSample(point.ID, lease, s.InstallationScope()))
			fresh := s.AddPointWithObservationLease(point)
			require.Greater(t, fresh, lease)
			require.True(t, s.AcceptsInstallationSample(point.ID, fresh, s.InstallationScope()))
			require.False(t, s.AcceptsInstallationSample(point.ID, lease, s.InstallationScope()))
		})
	}
}

func TestAcquisitionInstallationConfigurationABAAndPartialReinstall(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	d := &schema.Device{ID: "d", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{}`}
	other := &schema.Device{ID: "other", Protocol: d.Protocol, ConnectionConfig: d.ConnectionConfig}
	s.AddDevice(d)
	s.AddDevice(other)
	p := &schema.Point{ID: "p", DeviceID: d.ID, DataType: schema.DataTypeUint16}
	q := &schema.Point{ID: "q", DeviceID: d.ID, DataType: schema.DataTypeUint16}
	foreign := &schema.Point{ID: "foreign", DeviceID: other.ID, DataType: schema.DataTypeUint16}
	oldP, oldQ := s.AddPointWithObservationLease(p), s.AddPointWithObservationLease(q)
	otherLease := s.AddPointWithObservationLease(foreign)
	d.ConnectionConfig = `{"timeout":2}`
	s.AddDevice(d)
	d.ConnectionConfig = `{}`
	s.AddDevice(d)
	require.False(t, s.AcceptsInstallationSample(p.ID, oldP, s.InstallationScope()), "A->B->A before reinstall cannot bless A's old lease")
	require.False(t, s.AcceptsInstallationSample(q.ID, oldQ, s.InstallationScope()))
	require.True(t, s.AcceptsInstallationSample(foreign.ID, otherLease, s.InstallationScope()))
	// A failed provider or partially completed install leaves every untouched
	// slot invalid; successfully installed slots alone gain the current stamp.
	newP := s.AddPointWithObservationLease(p)
	require.True(t, s.AcceptsInstallationSample(p.ID, newP, s.InstallationScope()))
	require.False(t, s.AcceptsInstallationSample(q.ID, oldQ, s.InstallationScope()))
	require.ErrorIs(t, s.admitScheduled(scheduledSnapshot{device: s.deviceConfigs[d.ID]}, []pointInfo{s.pointInfos[q.ID]}), errAcquisitionSuperseded)
	require.True(t, s.AcceptsInstallationSample(foreign.ID, otherLease, s.InstallationScope()))
}
