package runtime

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"

	"github.com/stretchr/testify/require"
)

type runningProjectionSampleSink struct {
	mu           sync.Mutex
	samples      map[string][]time.Time
	pointSamples map[string][]time.Time
}

func (s *runningProjectionSampleSink) AcceptSample(_ context.Context, sample measurement.SampleEnvelope) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.samples[sample.DeviceID] = append(s.samples[sample.DeviceID], sample.ReceivedAt)
	if s.pointSamples != nil {
		s.pointSamples[sample.PointID] = append(s.pointSamples[sample.PointID], sample.ReceivedAt)
	}
	return nil
}

func (s *runningProjectionSampleSink) countAfter(deviceID string, after time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := 0
	for _, receivedAt := range s.samples[deviceID] {
		if receivedAt.After(after) {
			count++
		}
	}
	return count
}

func (s *runningProjectionSampleSink) pointCountAfter(pointID string, after time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := 0
	for _, receivedAt := range s.pointSamples[pointID] {
		if receivedAt.After(after) {
			count++
		}
	}
	return count
}

func TestService_ApplyWorkspaceProjectionForDevicesPreservesUnselectedRuntimeState(t *testing.T) {
	ctx := t.Context()
	runtimeSvc := newRuntimeServiceForDeviceSyncTest(t)
	deviceB := &schema.Device{ID: "dev-B", Name: "Device B", Status: schema.DeviceStatusActive}
	pointB := &schema.Point{ID: "point-B", DeviceID: deviceB.ID, Name: "Point B", Address: "40002", DataType: schema.DataTypeInt16, Enabled: true}
	createdPoint, err := runtimeSvc.pointSvc.Create(ctx, point.CreatePointRequest{
		DeviceID: deviceB.ID, Name: pointB.Name, Address: pointB.Address,
		DataType: pointB.DataType, Mode: schema.PointModeReadOnly,
	})
	require.NoError(t, err)
	pointB.ID = createdPoint.ID

	// Seed B through the existing per-device sync path, then update only A from
	// one coherent workspace projection.
	require.NoError(t, runtimeSvc.UpsertDevice(ctx, deviceB))
	beforeState, beforeExists := runtimeSvc.scheduler.GetDeviceBreakerState(deviceB.ID)
	require.True(t, beforeExists)
	beforeMeta := runtimeSvc.lookupPointMeta(pointB.ID)
	runtimeSvc.markDeviceProjectionAligned(deviceB.ID, &workspace.RuntimeProjection{Version: "before"})

	projection := &workspace.RuntimeProjection{
		WorkspaceID: "workspace-1",
		Version:     "after",
		DeviceIDs:   []string{"dev-A", "dev-B"},
		Devices: []*schema.Device{
			{ID: "dev-A", Name: "Device A", Status: schema.DeviceStatusActive},
			deviceB,
		},
		Points: []*schema.Point{pointB},
	}
	require.NoError(t, runtimeSvc.ApplyWorkspaceProjectionForDevices(ctx, projection, []string{"dev-A"}))

	afterState, afterExists := runtimeSvc.scheduler.GetDeviceBreakerState(deviceB.ID)
	require.True(t, afterExists)
	require.Equal(t, beforeState, afterState)
	require.Equal(t, beforeMeta, runtimeSvc.lookupPointMeta(pointB.ID))
	require.Equal(t, "before", runtimeSvc.deviceProjectionState(deviceB.ID).RuntimeVersion)
	_, aExists := runtimeSvc.scheduler.GetDeviceBreakerState("dev-A")
	require.True(t, aExists)
}

func TestService_ApplyWorkspaceProjectionForDevicesRejectsUnscopedDevice(t *testing.T) {
	runtimeSvc := &Service{scheduler: collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)}
	err := runtimeSvc.ApplyWorkspaceProjectionForDevices(context.Background(), &workspace.RuntimeProjection{
		DeviceIDs: []string{"dev-A"},
		Devices:   []*schema.Device{{ID: "dev-A", Status: schema.DeviceStatusActive}},
	}, []string{"dev-B"})
	require.Error(t, err)
}

func TestService_ApplyWorkspaceProjectionForDevicesKeepsRunningUnselectedPollingGroup(t *testing.T) {
	ctx := t.Context()
	bank := memory.NewMemoryBank(64)
	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	defer server.Stop()
	config, err := json.Marshal(schema.ConnectionConfigModbusTCP{
		Host: "127.0.0.1", Port: server.Port(), SlaveID: 1, Timeout: 1,
	})
	require.NoError(t, err)

	deviceA := &schema.Device{ID: "dev-A", Name: "Device A", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: string(config)}
	deviceB := &schema.Device{ID: "dev-B", Name: "Device B", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: string(config)}
	groupA := &schema.PollingGroup{ID: "group-A", IntervalMs: 10, Enabled: true}
	groupB := &schema.PollingGroup{ID: "group-B", IntervalMs: 10, Enabled: true}
	pointA := &schema.Point{ID: "point-A", DeviceID: deviceA.ID, Name: "Point A", Address: "40001", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &groupA.ID}
	pointB := &schema.Point{ID: "point-B", DeviceID: deviceB.ID, Name: "Point B", Address: "40001", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &groupB.ID}
	tagA := &schema.Tag{ID: "tag-A", Key: "tag-a", KeyLower: "tag-a", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}
	tagB := &schema.Tag{ID: "tag-B", Key: "tag-b", KeyLower: "tag-b", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}
	mappingA := &schema.Mapping{ID: "mapping-A", PointID: pointA.ID, TagID: tagA.ID, TransformPipeline: "[]", Enabled: true}
	mappingB := &schema.Mapping{ID: "mapping-B", PointID: pointB.ID, TagID: tagB.ID, TransformPipeline: "[]", Enabled: true}
	sink := &runningProjectionSampleSink{samples: make(map[string][]time.Time), pointSamples: make(map[string][]time.Time)}
	runtimeSvc, err := NewService(Config{
		WorkspaceID: "workspace-1", Writer: storage.NewMemoryStorage(64), SampleSink: sink,
		Snapshot: Snapshot{
			WorkspaceID:   "workspace-1",
			Devices:       []*schema.Device{deviceA, deviceB},
			Points:        []*schema.Point{pointA, pointB},
			PollingGroups: []*schema.PollingGroup{groupA, groupB},
			Mappings:      []*schema.Mapping{mappingA, mappingB},
			Tags:          []*schema.Tag{tagA, tagB},
		},
	})
	require.NoError(t, err)
	require.NoError(t, runtimeSvc.Start(ctx))
	defer func() { require.NoError(t, runtimeSvc.Stop(ctx)) }()
	require.Eventually(t, func() bool { return sink.countAfter(deviceB.ID, time.Time{}) > 0 }, time.Second, 5*time.Millisecond)

	projection := &workspace.RuntimeProjection{
		WorkspaceID: "workspace-1",
		Version:     "after",
		DeviceIDs:   []string{deviceA.ID, deviceB.ID},
		Devices:     []*schema.Device{deviceA, deviceB},
		Points:      []*schema.Point{pointA, pointB},
		Mappings:    []*schema.Mapping{mappingA, mappingB},
		Tags:        []*schema.Tag{tagA, tagB},
		PollingGroups: []*schema.PollingGroup{
			{ID: groupA.ID, IntervalMs: groupA.IntervalMs, Enabled: true},
			{ID: groupB.ID, IntervalMs: 1000, Enabled: true},
		},
	}
	require.NoError(t, runtimeSvc.ApplyWorkspaceProjectionForDevices(ctx, projection, []string{deviceA.ID}))
	appliedAt := time.Now().UTC()
	// B remains on its original 10 ms ticker. A full-projection AddPollingGroup
	// would replace it with the tampered 1 s interval and yield no samples here.
	time.Sleep(150 * time.Millisecond)
	count := sink.countAfter(deviceB.ID, appliedAt)
	require.Greater(t, count, 2)
}
