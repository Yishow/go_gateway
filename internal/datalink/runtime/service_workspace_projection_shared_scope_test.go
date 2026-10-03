package runtime

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"

	"github.com/stretchr/testify/require"
)

func TestService_ApplyWorkspaceProjectionForDevicesAllowsUnchangedSharedPollingGroup(t *testing.T) {
	ctx := t.Context()
	runtimeSvc, sink, deviceA, deviceB, group := newSharedPollingProjectionRuntime(t)
	defer func() { require.NoError(t, runtimeSvc.Stop(ctx)) }()

	projection := sharedPollingProjection(deviceA, deviceB, group, group.IntervalMs, "after")
	require.NoError(t, runtimeSvc.ApplyWorkspaceProjectionForDevices(ctx, projection, []string{deviceA.ID}))
	probeAt := time.Now().UTC()
	time.Sleep(100 * time.Millisecond)
	require.Greater(t, sink.countAfter(deviceB.ID, probeAt), 2)
}

func TestService_ApplyWorkspaceProjectionForDevicesRejectsChangedSharedPollingGroup(t *testing.T) {
	ctx := t.Context()
	runtimeSvc, sink, deviceA, deviceB, group := newSharedPollingProjectionRuntime(t)
	defer func() { require.NoError(t, runtimeSvc.Stop(ctx)) }()
	runtimeSvc.markDeviceProjectionAligned(deviceA.ID, &workspace.RuntimeProjection{Version: "before"})

	projection := sharedPollingProjection(deviceA, deviceB, group, 1000, "after")
	err := runtimeSvc.ApplyWorkspaceProjectionForDevices(ctx, projection, []string{deviceA.ID})
	require.Error(t, err)
	require.Equal(t, "before", runtimeSvc.deviceProjectionState(deviceA.ID).RuntimeVersion)
	probeAt := time.Now().UTC()
	time.Sleep(100 * time.Millisecond)
	require.Greater(t, sink.countAfter(deviceB.ID, probeAt), 2)
}

func TestService_ApplyWorkspaceProjectionForDevicesRejectsMissingSharedPollingGroup(t *testing.T) {
	ctx := t.Context()
	runtimeSvc, sink, deviceA, deviceB, group := newSharedPollingProjectionRuntime(t)
	defer func() { require.NoError(t, runtimeSvc.Stop(ctx)) }()
	runtimeSvc.markDeviceProjectionAligned(deviceA.ID, &workspace.RuntimeProjection{Version: "before"})

	projection := sharedPollingProjection(deviceA, deviceB, group, group.IntervalMs, "after")
	projection.PollingGroups = nil
	err := runtimeSvc.ApplyWorkspaceProjectionForDevices(ctx, projection, []string{deviceA.ID})
	require.Error(t, err)
	require.Equal(t, "before", runtimeSvc.deviceProjectionState(deviceA.ID).RuntimeVersion)
	probeAt := time.Now().UTC()
	time.Sleep(100 * time.Millisecond)
	require.Greater(t, sink.countAfter(deviceB.ID, probeAt), 2)
}

func TestService_ApplyWorkspaceProjectionForDevicesAndPointsPreservesSameDeviceUnselectedGroup(t *testing.T) {
	ctx := t.Context()
	bank := memory.NewMemoryBank(64)
	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	defer func() { require.NoError(t, server.Stop()) }()
	config, err := json.Marshal(schema.ConnectionConfigModbusTCP{Host: "127.0.0.1", Port: server.Port(), SlaveID: 1, Timeout: 1})
	require.NoError(t, err)

	device := &schema.Device{ID: "dev-A", Name: "Device A", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: string(config)}
	groupA := &schema.PollingGroup{ID: "poll-A", IntervalMs: 10, Enabled: true}
	groupB := &schema.PollingGroup{ID: "poll-B", IntervalMs: 10, Enabled: true}
	pointA := &schema.Point{ID: "point-A", DeviceID: device.ID, Address: "40001", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &groupA.ID}
	pointB := &schema.Point{ID: "point-B", DeviceID: device.ID, Address: "40002", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &groupB.ID}
	tagA := &schema.Tag{ID: "tag-A", Key: "tag-a", KeyLower: "tag-a", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}
	tagB := &schema.Tag{ID: "tag-B", Key: "tag-b", KeyLower: "tag-b", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}
	sink := &runningProjectionSampleSink{samples: make(map[string][]time.Time), pointSamples: make(map[string][]time.Time)}
	runtimeSvc, err := NewService(Config{
		WorkspaceID: "workspace-1", Writer: storage.NewMemoryStorage(64), SampleSink: sink,
		Snapshot: Snapshot{
			WorkspaceID: "workspace-1", Devices: []*schema.Device{device}, Points: []*schema.Point{pointA, pointB},
			PollingGroups: []*schema.PollingGroup{groupA, groupB},
			Mappings:      []*schema.Mapping{{ID: "mapping-A", PointID: pointA.ID, TagID: tagA.ID, TransformPipeline: "[]", Enabled: true}, {ID: "mapping-B", PointID: pointB.ID, TagID: tagB.ID, TransformPipeline: "[]", Enabled: true}},
			Tags:          []*schema.Tag{tagA, tagB},
		},
	})
	require.NoError(t, err)
	require.NoError(t, runtimeSvc.Start(ctx))
	defer func() { require.NoError(t, runtimeSvc.Stop(ctx)) }()
	require.Eventually(t, func() bool { return sink.countAfter(device.ID, time.Time{}) > 0 }, time.Second, 5*time.Millisecond)
	beforeMeta := runtimeSvc.lookupPointMeta(pointB.ID)

	projection := &workspace.RuntimeProjection{
		WorkspaceID: "workspace-1", Version: "after", DeviceIDs: []string{device.ID}, Devices: []*schema.Device{device},
		Points: []*schema.Point{pointA, pointB}, PollingGroups: []*schema.PollingGroup{groupA, groupB},
		Mappings: []*schema.Mapping{{ID: "mapping-A", PointID: pointA.ID, TagID: tagA.ID, TransformPipeline: "[]", Enabled: true}, {ID: "mapping-B", PointID: pointB.ID, TagID: tagB.ID, TransformPipeline: "[]", Enabled: true}},
		Tags:     []*schema.Tag{tagA, tagB},
	}
	require.NoError(t, runtimeSvc.ApplyWorkspaceProjectionForDevicesAndPoints(ctx, projection, []string{device.ID}, []string{pointA.ID}))
	probeAt := time.Now().UTC()
	time.Sleep(100 * time.Millisecond)
	require.Greater(t, sink.pointCountAfter(pointB.ID, probeAt), 2)
	require.Equal(t, beforeMeta, runtimeSvc.lookupPointMeta(pointB.ID))
}

func TestService_ApplyWorkspaceProjectionForDevicesAndPointsKeepsSameDeviceSharedTicker(t *testing.T) {
	ctx := t.Context()
	runtimeSvc, sink, device, group := newSameDeviceSharedPollingProjectionRuntime(t, 1000)
	defer func() { require.NoError(t, runtimeSvc.Stop(ctx)) }()
	require.Eventually(t, func() bool { return sink.pointCountAfter("point-B", time.Time{}) > 0 }, 1500*time.Millisecond, 10*time.Millisecond)
	time.Sleep(700 * time.Millisecond)

	projection := sameDeviceSharedPollingProjection(device, group, group.IntervalMs, "after")
	probeAt := time.Now().UTC()
	require.NoError(t, runtimeSvc.ApplyWorkspaceProjectionForDevicesAndPoints(ctx, projection, []string{device.ID}, []string{"point-A"}))
	require.Eventually(t, func() bool { return sink.pointCountAfter("point-B", probeAt) > 0 }, 600*time.Millisecond, 10*time.Millisecond)
}

func newSharedPollingProjectionRuntime(t *testing.T) (runtimeSvc *Service, sink *runningProjectionSampleSink, deviceA, deviceB *schema.Device, group *schema.PollingGroup) {
	return newSharedPollingProjectionRuntimeWithInterval(t, 10)
}

func newSharedPollingProjectionRuntimeWithInterval(t *testing.T, interval int) (runtimeSvc *Service, sink *runningProjectionSampleSink, deviceA, deviceB *schema.Device, group *schema.PollingGroup) {
	t.Helper()
	ctx := t.Context()
	bank := memory.NewMemoryBank(64)
	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	t.Cleanup(func() { require.NoError(t, server.Stop()) })
	config, err := json.Marshal(schema.ConnectionConfigModbusTCP{
		Host: "127.0.0.1", Port: server.Port(), SlaveID: 1, Timeout: 1,
	})
	require.NoError(t, err)

	deviceA = &schema.Device{ID: "dev-A", Name: "Device A", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: string(config)}
	deviceB = &schema.Device{ID: "dev-B", Name: "Device B", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: string(config)}
	group = &schema.PollingGroup{ID: "group-shared", IntervalMs: interval, Enabled: true}
	pointA := &schema.Point{ID: "point-A", DeviceID: deviceA.ID, Name: "Point A", Address: "40001", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &group.ID}
	pointB := &schema.Point{ID: "point-B", DeviceID: deviceB.ID, Name: "Point B", Address: "40001", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &group.ID}
	tagA := &schema.Tag{ID: "tag-A", Key: "tag-a", KeyLower: "tag-a", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}
	tagB := &schema.Tag{ID: "tag-B", Key: "tag-b", KeyLower: "tag-b", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}
	mappingA := &schema.Mapping{ID: "mapping-A", PointID: pointA.ID, TagID: tagA.ID, TransformPipeline: "[]", Enabled: true}
	mappingB := &schema.Mapping{ID: "mapping-B", PointID: pointB.ID, TagID: tagB.ID, TransformPipeline: "[]", Enabled: true}
	sink = &runningProjectionSampleSink{samples: make(map[string][]time.Time)}
	runtimeSvc, err = NewService(Config{
		WorkspaceID: "workspace-1", Writer: storage.NewMemoryStorage(64), SampleSink: sink,
		Snapshot: Snapshot{
			WorkspaceID:   "workspace-1",
			Devices:       []*schema.Device{deviceA, deviceB},
			Points:        []*schema.Point{pointA, pointB},
			PollingGroups: []*schema.PollingGroup{group},
			Mappings:      []*schema.Mapping{mappingA, mappingB},
			Tags:          []*schema.Tag{tagA, tagB},
		},
	})
	require.NoError(t, err)
	require.NoError(t, runtimeSvc.Start(ctx))
	require.Eventually(t, func() bool { return sink.countAfter(deviceB.ID, time.Time{}) > 0 }, time.Duration(interval+500)*time.Millisecond, 5*time.Millisecond)
	return runtimeSvc, sink, deviceA, deviceB, group
}

func newSameDeviceSharedPollingProjectionRuntime(t *testing.T, interval int) (runtimeSvc *Service, sink *runningProjectionSampleSink, device *schema.Device, group *schema.PollingGroup) {
	t.Helper()
	ctx := t.Context()
	bank := memory.NewMemoryBank(64)
	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	t.Cleanup(func() { require.NoError(t, server.Stop()) })
	config, err := json.Marshal(schema.ConnectionConfigModbusTCP{Host: "127.0.0.1", Port: server.Port(), SlaveID: 1, Timeout: 1})
	require.NoError(t, err)

	device = &schema.Device{ID: "dev-A", Name: "Device A", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: string(config)}
	group = &schema.PollingGroup{ID: "group-shared", IntervalMs: interval, Enabled: true}
	pointA := &schema.Point{ID: "point-A", DeviceID: device.ID, Address: "40001", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &group.ID}
	pointB := &schema.Point{ID: "point-B", DeviceID: device.ID, Address: "40002", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &group.ID}
	tagA := &schema.Tag{ID: "tag-A", Key: "tag-a", KeyLower: "tag-a", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}
	tagB := &schema.Tag{ID: "tag-B", Key: "tag-b", KeyLower: "tag-b", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}
	sink = &runningProjectionSampleSink{samples: make(map[string][]time.Time), pointSamples: make(map[string][]time.Time)}
	runtimeSvc, err = NewService(Config{
		WorkspaceID: "workspace-1", Writer: storage.NewMemoryStorage(64), SampleSink: sink,
		Snapshot: Snapshot{
			WorkspaceID: "workspace-1", Devices: []*schema.Device{device}, Points: []*schema.Point{pointA, pointB},
			PollingGroups: []*schema.PollingGroup{group},
			Mappings:      []*schema.Mapping{{ID: "mapping-A", PointID: pointA.ID, TagID: tagA.ID, TransformPipeline: "[]", Enabled: true}, {ID: "mapping-B", PointID: pointB.ID, TagID: tagB.ID, TransformPipeline: "[]", Enabled: true}},
			Tags:          []*schema.Tag{tagA, tagB},
		},
	})
	require.NoError(t, err)
	require.NoError(t, runtimeSvc.Start(ctx))
	require.Eventually(t, func() bool { return sink.pointCountAfter(pointB.ID, time.Time{}) > 0 }, time.Duration(interval+500)*time.Millisecond, 5*time.Millisecond)
	return runtimeSvc, sink, device, group
}

func sharedPollingProjection(deviceA, deviceB *schema.Device, group *schema.PollingGroup, interval int, version string) *workspace.RuntimeProjection {
	groupCopy := *group
	groupCopy.IntervalMs = interval
	return &workspace.RuntimeProjection{
		WorkspaceID:   "workspace-1",
		Version:       version,
		DeviceIDs:     []string{deviceA.ID, deviceB.ID},
		Devices:       []*schema.Device{deviceA, deviceB},
		Points:        []*schema.Point{{ID: "point-A", DeviceID: deviceA.ID, Address: "40001", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &groupCopy.ID}, {ID: "point-B", DeviceID: deviceB.ID, Address: "40001", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &groupCopy.ID}},
		Mappings:      []*schema.Mapping{{ID: "mapping-A", PointID: "point-A", TagID: "tag-A", TransformPipeline: "[]", Enabled: true}, {ID: "mapping-B", PointID: "point-B", TagID: "tag-B", TransformPipeline: "[]", Enabled: true}},
		Tags:          []*schema.Tag{{ID: "tag-A", Key: "tag-a", KeyLower: "tag-a", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}, {ID: "tag-B", Key: "tag-b", KeyLower: "tag-b", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}},
		PollingGroups: []*schema.PollingGroup{&groupCopy},
	}
}

func sameDeviceSharedPollingProjection(device *schema.Device, group *schema.PollingGroup, interval int, version string) *workspace.RuntimeProjection {
	groupCopy := *group
	groupCopy.IntervalMs = interval
	return &workspace.RuntimeProjection{
		WorkspaceID: "workspace-1", Version: version, DeviceIDs: []string{device.ID}, Devices: []*schema.Device{device},
		Points:        []*schema.Point{{ID: "point-A", DeviceID: device.ID, Address: "40001", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &groupCopy.ID}, {ID: "point-B", DeviceID: device.ID, Address: "40002", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &groupCopy.ID}},
		Mappings:      []*schema.Mapping{{ID: "mapping-A", PointID: "point-A", TagID: "tag-A", TransformPipeline: "[]", Enabled: true}, {ID: "mapping-B", PointID: "point-B", TagID: "tag-B", TransformPipeline: "[]", Enabled: true}},
		Tags:          []*schema.Tag{{ID: "tag-A", Key: "tag-a", KeyLower: "tag-a", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}, {ID: "tag-B", Key: "tag-b", KeyLower: "tag-b", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}},
		PollingGroups: []*schema.PollingGroup{&groupCopy},
	}
}
