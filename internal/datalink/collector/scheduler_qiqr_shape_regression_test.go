package collector

import (
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestScheduledQIQRUsesHexBitAddressBeforePlanning(t *testing.T) {
	config := `{"acquisition":{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1"}}`
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	device := &schema.Device{ID: "device", Protocol: schema.ProtocolMC3E, ConnectionConfig: config}
	s.AddDevice(device)
	point := &schema.Point{ID: "point", DeviceID: device.ID, Address: "X1A", DataType: schema.DataTypeBool, RevisionID: "p1", FormatMode: schema.PointFormatInherit, FormatVersion: 1}
	profile, err := connector.ResolveReadProfile(*point, decoding.FormatContext{Protocol: device.Protocol, ConnectionRevision: connector.ReadConnectionRevision(config)})
	require.NoError(t, err)
	require.NotZero(t, s.AddPointWithObservationLease(point, profile))
	snapshots := s.captureScheduledGroup("")
	require.Len(t, snapshots, 1)
	plan, err := buildScheduledPlan(snapshots[0])
	require.NoError(t, err)
	require.Empty(t, plan.failures)
	require.Empty(t, plan.plan.Errors())
	require.Len(t, plan.plan.Requests(), 1)
	require.Equal(t, 0x1a, plan.plan.Requests()[0].Start)
}
