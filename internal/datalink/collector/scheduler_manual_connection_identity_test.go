package collector

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestSchedulerUsesManagedConnectionConfigurationForFingerprint(t *testing.T) {
	protocolType := schema.ProtocolType("acquisition-facts-managed-config-test")
	var connectedConfigs []string
	connector.Register(protocolType, func() connector.Protocol {
		return acquisitionFactsProtocol{connectedConfig: func(config string) { connectedConfigs = append(connectedConfigs, config) }}
	})
	defer connector.Unregister(protocolType)

	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	scheduler := NewScheduler(DefaultSchedulerConfig(), manager)
	scheduler.WithClock(func() time.Time { return time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC) })
	scheduler.WithAcquisitionIDFactory(func() string { return "acquisition-managed-config" })
	oldConfig := `{"host":"old-device-a"}`
	newConfig := `{"host":"new-device-a"}`
	scheduler.AddDevice(&schema.Device{ID: "device-managed", Protocol: protocolType, ConnectionConfig: oldConfig})
	point := &schema.Point{ID: "point-managed", DeviceID: "device-managed", Address: "40001", Function: "03", DataType: schema.DataTypeUint64, Mode: schema.PointModeReadOnly}
	scheduler.AddPoint(point)

	first := scheduler.PollNowContext(t.Context(), []string{point.ID})
	require.Len(t, first, 1)

	scheduler.AddDevice(&schema.Device{ID: "device-managed", Protocol: protocolType, ConnectionConfig: newConfig})
	stale := scheduler.PollNowContext(t.Context(), []string{point.ID})
	require.Len(t, stale, 1)
	require.Equal(t, schema.QualityBad, stale[0].Quality, "connection replacement invalidates the old point installation")
	require.Equal(t, []string{oldConfig}, connectedConfigs, "stale install must not connect as the new revision")
	scheduler.AddPoint(point)
	second := scheduler.PollNowContext(t.Context(), []string{point.ID})
	require.Len(t, second, 1)

	require.Equal(t, []string{oldConfig, newConfig}, connectedConfigs)
	oldFingerprint, err := measurement.AcquisitionConfigFingerprint(schema.Device{
		ID:               "device-managed",
		Protocol:         protocolType,
		ConnectionConfig: connectedConfigs[0],
	}, *point)
	require.NoError(t, err)
	newFingerprint, err := measurement.AcquisitionConfigFingerprint(schema.Device{
		ID:               "device-managed",
		Protocol:         protocolType,
		ConnectionConfig: connectedConfigs[1],
	}, *point)
	require.NoError(t, err)
	require.NotEqual(t, oldFingerprint, newFingerprint)
	require.Equal(t, oldFingerprint, first[0].ConfigFingerprint)
	require.Equal(t, newFingerprint, second[0].ConfigFingerprint)
}
