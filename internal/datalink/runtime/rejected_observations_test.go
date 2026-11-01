package runtime

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestRejectedObservationCounterDoesNotRefreshOrChangePublicStats(t *testing.T) {
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	scheduler.AddDevice(&schema.Device{ID: "device", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{}`})
	p := &schema.Point{ID: "point", DeviceID: "device", Address: "40001", DataType: schema.DataTypeUint16}
	oldLease := scheduler.AddPointWithObservationLease(p)
	scheduler.RemovePoint(p.ID)
	newLease := scheduler.AddPointWithObservationLease(p)
	require.Greater(t, newLease, oldLease)
	s := &Service{scheduler: scheduler}
	s.handleCollectedValue(t.Context(), collector.CollectedValue{PointID: p.ID, DeviceID: p.DeviceID, InstallationLease: oldLease, InstallationScope: scheduler.InstallationScope(), Quality: schema.QualityGood, Value: uint16(42)})
	require.EqualValues(t, 1, s.RejectedObservationCount())
	require.Zero(t, s.Snapshot().PointStateError)
	require.Zero(t, s.Snapshot().WriteError)
	encoded, err := json.Marshal(s.Snapshot())
	require.NoError(t, err)
	var fields map[string]any
	require.NoError(t, json.Unmarshal(encoded, &fields))
	require.Len(t, fields, 5, "native rejection evidence must not silently expand public metrics")
	require.NotContains(t, fields, "rejected_observations")
}

func TestRejectedObservationConfigPublishCannotAcceptOldInstalledLease(t *testing.T) {
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	scheduler.AddDevice(&schema.Device{ID: "device", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{}`})
	point := &schema.Point{ID: "point", DeviceID: "device", Address: "40001", DataType: schema.DataTypeUint16}
	lease := scheduler.AddPointWithObservationLease(point)
	service := &Service{scheduler: scheduler}
	value := collector.CollectedValue{PointID: point.ID, DeviceID: point.DeviceID, InstallationLease: lease, InstallationScope: scheduler.InstallationScope()}
	_, accepted := service.captureCollectedBindings(value)
	require.True(t, accepted)
	// UpsertDevice publishes device config before its fallible point-provider
	// reads/install loop. During that interval the old lease still exists.
	scheduler.AddDevice(&schema.Device{ID: "device", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{"timeout":2}`})
	_, accepted = service.captureCollectedBindings(value)
	require.False(t, accepted, "the old source config must not refresh new device configuration while Point reinstalls are pending")
	value.InstallationLease = 0
	_, accepted = service.captureCollectedBindings(value)
	require.False(t, accepted, "a legacy/manual result must not bypass the invalid installation stamp")
}
