package runtime

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"

	"github.com/stretchr/testify/require"
)

func TestService_ReconcileSourceRuleAppliesWorkspaceProjection(t *testing.T) {
	ctx := context.Background()
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	deviceRecord := &schema.Device{
		ID:               "dev-A",
		Name:             "Device A",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":1}`,
	}
	pointRecord := &schema.Point{
		ID:       "point-A",
		DeviceID: deviceRecord.ID,
		Name:     "Disabled point",
		Address:  "40001",
		DataType: schema.DataTypeInt16,
		Mode:     schema.PointModeReadOnly,
		Enabled:  false,
	}
	projection := &workspace.RuntimeProjection{
		WorkspaceID: "ws-1",
		Version:     "projection-v2",
		Alignment:   workspace.RuntimeProjectionAlignmentAligned,
		DeviceIDs:   []string{deviceRecord.ID},
		Devices:     []*schema.Device{deviceRecord},
		Points:      []*schema.Point{pointRecord},
	}
	svc := &Service{
		scheduler:      scheduler,
		workspace:      workspaceProjectionStub{projection: projection},
		mappingIndex:   map[string][]mappingBinding{},
		pointMetaIndex: map[string]pointMeta{},
	}
	svc.running.Store(true)

	outcome := svc.ReconcileSourceRule(ctx, sourcerule.RuntimeReconcileRequest{
		Operation: sourcerule.RuntimeReconcileOperationDisable,
		Scope: sourcerule.RuntimeReconcileScope{
			RuleID:   "rule-A",
			DeviceID: deviceRecord.ID,
		},
	})

	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	require.Equal(t, "rule-A", outcome.Scope.RuleID)
	_, deviceExists := scheduler.GetDeviceBreakerState(deviceRecord.ID)
	require.True(t, deviceExists)
	require.Empty(t, scheduler.PollNow([]string{pointRecord.ID}))
}

func TestService_ReconcileSourceRuleTargetsAffectedDeviceAndIsIdempotent(t *testing.T) {
	ctx := context.Background()
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	devA := &schema.Device{
		ID:               "dev-A",
		Name:             "Device A",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":1}`,
	}
	devB := &schema.Device{
		ID:               "dev-B",
		Name:             "Device B",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":1}`,
	}
	pointA := &schema.Point{ID: "point-A", DeviceID: devA.ID, Name: "A", Address: "40001", DataType: schema.DataTypeInt16, Enabled: true}
	pointB := &schema.Point{ID: "point-B", DeviceID: devB.ID, Name: "B", Address: "40002", DataType: schema.DataTypeInt16, Enabled: true}
	tagA := &schema.Tag{ID: "tag-A", Key: "a", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}
	projection := &workspace.RuntimeProjection{
		WorkspaceID: "ws-1",
		Version:     "projection-v3",
		Alignment:   workspace.RuntimeProjectionAlignmentAligned,
		DeviceIDs:   []string{devA.ID, devB.ID},
		Devices:     []*schema.Device{devA, devB},
		Points:      []*schema.Point{pointA, pointB},
		Mappings: []*schema.Mapping{
			{ID: "map-A", PointID: pointA.ID, TagID: tagA.ID, TransformPipeline: "[]", Enabled: true},
			{ID: "map-B", PointID: pointB.ID, TagID: "missing-tag", TransformPipeline: "[]", Enabled: true},
		},
		Tags: []*schema.Tag{tagA},
	}
	svc := &Service{
		scheduler:      scheduler,
		workspace:      workspaceProjectionStub{projection: projection},
		mappingIndex:   map[string][]mappingBinding{},
		pointMetaIndex: map[string]pointMeta{},
	}
	svc.running.Store(true)
	req := sourcerule.RuntimeReconcileRequest{
		Operation: sourcerule.RuntimeReconcileOperationUpdate,
		Scope: sourcerule.RuntimeReconcileScope{
			RuleID:   "rule-A",
			DeviceID: devA.ID,
		},
	}

	first := svc.ReconcileSourceRule(ctx, req)
	second := svc.ReconcileSourceRule(ctx, req)

	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, first.Status)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, second.Status)
	require.Len(t, svc.mappingBindingsForPoint(pointA.ID), 1)
	require.Empty(t, svc.mappingBindingsForPoint(pointB.ID))

	writer := &mockWriter{}
	svc.writer = writer
	svc.handleCollectedValue(ctx, collector.CollectedValue{
		PointID:   pointA.ID,
		DeviceID:  devA.ID,
		Value:     int16(0x1234),
		Timestamp: time.Date(2026, 8, 25, 14, 45, 0, 0, time.UTC),
		Quality:   schema.QualityGood,
	})
	require.Len(t, writer.records, 1)
	require.Equal(t, tagA.ID, writer.records[0].TagID)
	require.Equal(t, uint64(1), svc.writeSuccess.Load())
}

func TestService_ReconcileSourceRulePreservesHealthySameDeviceSharedTicker(t *testing.T) {
	ctx := t.Context()
	fixture := newSourceRuleReconcileRuntime(t, false, false)
	defer func() { require.NoError(t, fixture.service.Stop(ctx)) }()
	time.Sleep(700 * time.Millisecond)

	probeAt := time.Now().UTC()
	outcome := fixture.service.ReconcileSourceRule(ctx, sourcerule.RuntimeReconcileRequest{
		Operation: sourcerule.RuntimeReconcileOperationCreate,
		Scope: sourcerule.RuntimeReconcileScope{
			RuleID:   "rule-A",
			DeviceID: fixture.device.ID,
		},
	})

	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	require.Eventually(t, func() bool {
		return fixture.sink.pointCountAfter("point-B", probeAt) > 0
	}, 600*time.Millisecond, 10*time.Millisecond)
}

func TestService_ReconcileSourceRuleZeroPointCreatePreservesHealthySameDeviceSharedTicker(t *testing.T) {
	ctx := t.Context()
	fixture := newSourceRuleReconcileRuntime(t, false, false)
	defer func() { require.NoError(t, fixture.service.Stop(ctx)) }()
	time.Sleep(700 * time.Millisecond)

	probeAt := time.Now().UTC()
	outcome := fixture.service.ReconcileSourceRule(ctx, sourcerule.RuntimeReconcileRequest{
		Operation: sourcerule.RuntimeReconcileOperationCreate,
		Scope: sourcerule.RuntimeReconcileScope{
			RuleID:   "rule-zero-links",
			DeviceID: fixture.device.ID,
		},
	})

	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	require.Eventually(t, func() bool {
		return fixture.sink.pointCountAfter("point-B", probeAt) > 0
	}, 600*time.Millisecond, 10*time.Millisecond)
}

func TestService_ReconcileSourceRulePreservesHealthyOtherDeviceTicker(t *testing.T) {
	ctx := t.Context()
	fixture := newSourceRuleReconcileRuntime(t, true, false)
	defer func() { require.NoError(t, fixture.service.Stop(ctx)) }()
	time.Sleep(700 * time.Millisecond)

	probeAt := time.Now().UTC()
	outcome := fixture.service.ReconcileSourceRule(ctx, sourcerule.RuntimeReconcileRequest{
		Operation: sourcerule.RuntimeReconcileOperationCreate,
		Scope: sourcerule.RuntimeReconcileScope{
			RuleID:   "rule-A",
			DeviceID: fixture.device.ID,
		},
	})

	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	require.Eventually(t, func() bool {
		return fixture.sink.countAfter("dev-B", probeAt) > 0
	}, 600*time.Millisecond, 10*time.Millisecond)
}

func TestService_ReconcileSourceRuleDefersInactiveDeviceWithoutTouchingHealthyRuntime(t *testing.T) {
	ctx := t.Context()
	fixture := newSourceRuleReconcileRuntime(t, true, true)
	defer func() { require.NoError(t, fixture.service.Stop(ctx)) }()
	time.Sleep(700 * time.Millisecond)

	probeAt := time.Now().UTC()
	beforeProjection := fixture.reader.projection.Load()
	outcome := fixture.service.ReconcileSourceRule(ctx, sourcerule.RuntimeReconcileRequest{
		Operation: sourcerule.RuntimeReconcileOperationCreate,
		Scope: sourcerule.RuntimeReconcileScope{
			RuleID:   "rule-B",
			DeviceID: "dev-B",
		},
	})

	require.Equal(t, sourcerule.RuntimeReconcileStatusDeferred, outcome.Status)
	require.Same(t, beforeProjection, fixture.reader.projection.Load())
	require.Eventually(t, func() bool {
		return fixture.sink.pointCountAfter("point-A", probeAt) > 0
	}, 600*time.Millisecond, 10*time.Millisecond)
}

func TestService_ReconcileSourceRuleDoesNotDeferInconsistentDeviceScope(t *testing.T) {
	tests := []struct {
		name                string
		includeOtherDevice  bool
		otherDeviceInactive bool
		mutateProjection    func(*workspace.RuntimeProjection)
	}{
		{
			name:               "device id without record",
			includeOtherDevice: true,
			mutateProjection: func(projection *workspace.RuntimeProjection) {
				projection.Devices = projection.Devices[:1]
			},
		},
		{
			name:                "draft record without device id",
			includeOtherDevice:  true,
			otherDeviceInactive: true,
			mutateProjection: func(projection *workspace.RuntimeProjection) {
				projection.DeviceIDs = projection.DeviceIDs[:1]
			},
		},
		{
			name:               "unknown device status",
			includeOtherDevice: true,
			mutateProjection: func(projection *workspace.RuntimeProjection) {
				projection.Devices[1].Status = schema.DeviceStatus("corrupt")
			},
		},
		{
			name:               "empty device status",
			includeOtherDevice: true,
			mutateProjection: func(projection *workspace.RuntimeProjection) {
				projection.Devices[1].Status = ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			fixture := newSourceRuleReconcileRuntime(t, tt.includeOtherDevice, tt.otherDeviceInactive)
			defer func() { require.NoError(t, fixture.service.Stop(ctx)) }()
			next := cloneRuntimeProjection(fixture.projection)
			tt.mutateProjection(next)
			fixture.reader.projection.Store(next)

			outcome := fixture.service.ReconcileSourceRule(ctx, sourcerule.RuntimeReconcileRequest{
				Operation: sourcerule.RuntimeReconcileOperationCreate,
				Scope: sourcerule.RuntimeReconcileScope{
					RuleID:   "rule-B",
					DeviceID: "dev-B",
				},
			})

			require.Equal(t, sourcerule.RuntimeReconcileStatusStale, outcome.Status)
		})
	}
}

func TestService_ReconcileSourceRuleUpdateRemovesRemovedPoints(t *testing.T) {
	ctx := t.Context()
	fixture := newSourceRuleReconcileRuntime(t, false, false)
	defer func() { require.NoError(t, fixture.service.Stop(ctx)) }()
	fixture.reader.projection.Store(projectionWithoutPointA(fixture.projection))

	outcome := fixture.service.ReconcileSourceRule(ctx, sourcerule.RuntimeReconcileRequest{
		Operation: sourcerule.RuntimeReconcileOperationUpdate,
		Scope: sourcerule.RuntimeReconcileScope{
			RuleID:   "rule-A",
			DeviceID: fixture.device.ID,
		},
	})

	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	require.Empty(t, fixture.service.lookupPointMeta("point-A"))
	require.Empty(t, fixture.service.scheduler.PollNow([]string{"point-A"}))
}

func TestService_ReconcileSourceRuleDeleteRemovesRemovedPoints(t *testing.T) {
	ctx := t.Context()
	fixture := newSourceRuleReconcileRuntime(t, false, false)
	defer func() { require.NoError(t, fixture.service.Stop(ctx)) }()
	fixture.reader.projection.Store(projectionWithoutPointA(fixture.projection))

	outcome := fixture.service.ReconcileSourceRule(ctx, sourcerule.RuntimeReconcileRequest{
		Operation: sourcerule.RuntimeReconcileOperationDelete,
		Scope: sourcerule.RuntimeReconcileScope{
			RuleID:   "rule-A",
			DeviceID: fixture.device.ID,
		},
	})

	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	require.Empty(t, fixture.service.lookupPointMeta("point-A"))
	require.Empty(t, fixture.service.scheduler.PollNow([]string{"point-A"}))
}

type sourceRuleReconcileRuntimeFixture struct {
	service    *Service
	sink       *runningProjectionSampleSink
	device     *schema.Device
	projection *workspace.RuntimeProjection
	reader     *sourceRuleProjectionReader
}

func newSourceRuleReconcileRuntime(t *testing.T, includeOtherDevice, otherDeviceInactive bool) sourceRuleReconcileRuntimeFixture {
	t.Helper()
	ctx := t.Context()
	interval := 1000
	bank := memory.NewMemoryBank(64)
	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	t.Cleanup(func() { require.NoError(t, server.Stop()) })
	config, err := json.Marshal(schema.ConnectionConfigModbusTCP{Host: "127.0.0.1", Port: server.Port(), SlaveID: 1, Timeout: 1})
	require.NoError(t, err)

	device := &schema.Device{ID: "dev-A", Name: "Device A", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: string(config)}
	otherDevice := device
	if includeOtherDevice {
		otherDevice = &schema.Device{ID: "dev-B", Name: "Device B", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: string(config)}
		if otherDeviceInactive {
			otherDevice.Status = schema.DeviceStatusDraft
		}
	}
	group := &schema.PollingGroup{ID: "group-shared", IntervalMs: interval, Enabled: true}
	pointA := &schema.Point{ID: "point-A", DeviceID: device.ID, Address: "40001", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &group.ID}
	pointB := &schema.Point{ID: "point-B", DeviceID: otherDevice.ID, Address: "40002", Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true, PollingGroupID: &group.ID}
	tagA := &schema.Tag{ID: "tag-A", Key: "tag-a", KeyLower: "tag-a", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}
	tagB := &schema.Tag{ID: "tag-B", Key: "tag-b", KeyLower: "tag-b", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}
	sink := &runningProjectionSampleSink{samples: make(map[string][]time.Time), pointSamples: make(map[string][]time.Time)}
	deviceIDs := []string{device.ID}
	devices := []*schema.Device{device}
	if includeOtherDevice {
		deviceIDs = append(deviceIDs, otherDevice.ID)
		devices = append(devices, otherDevice)
	}
	projection := &workspace.RuntimeProjection{
		WorkspaceID: "workspace-1", Version: "before", DeviceIDs: deviceIDs, Devices: devices,
		RuleLinks: []*schema.SourceRuleLink{{RuleID: "rule-A", PointID: "point-A"}},
		Points:    []*schema.Point{pointA, pointB}, PollingGroups: []*schema.PollingGroup{group},
		Mappings: []*schema.Mapping{{ID: "mapping-A", PointID: pointA.ID, TagID: tagA.ID, TransformPipeline: "[]", Enabled: true}, {ID: "mapping-B", PointID: pointB.ID, TagID: tagB.ID, TransformPipeline: "[]", Enabled: true}},
		Tags:     []*schema.Tag{tagA, tagB},
	}
	runtimeSvc, err := NewService(Config{
		WorkspaceID: "workspace-1", Writer: storage.NewMemoryStorage(64), SampleSink: sink,
		Snapshot: Snapshot{
			WorkspaceID: projection.WorkspaceID, Devices: projection.Devices, Points: projection.Points,
			PollingGroups: projection.PollingGroups, Mappings: projection.Mappings, Tags: projection.Tags,
		},
	})
	require.NoError(t, err)
	reader := &sourceRuleProjectionReader{}
	reader.projection.Store(projection)
	runtimeSvc.workspace = reader
	require.NoError(t, runtimeSvc.Start(ctx))
	waitPointID := pointB.ID
	if otherDevice.Status != schema.DeviceStatusActive {
		waitPointID = pointA.ID
	}
	require.Eventually(t, func() bool {
		return sink.pointCountAfter(waitPointID, time.Time{}) > 0
	}, time.Duration(interval+500)*time.Millisecond, 5*time.Millisecond)
	return sourceRuleReconcileRuntimeFixture{
		service: runtimeSvc, sink: sink, device: device, projection: projection, reader: reader,
	}
}

type sourceRuleProjectionReader struct {
	projection atomic.Pointer[workspace.RuntimeProjection]
}

func (r *sourceRuleProjectionReader) RuntimeProjection(context.Context) (*workspace.RuntimeProjection, error) {
	return r.projection.Load(), nil
}

func projectionWithoutPointA(projection *workspace.RuntimeProjection) *workspace.RuntimeProjection {
	next := *projection
	next.Version = "after-removal"
	next.RuleLinks = nil
	next.Points = []*schema.Point{projection.Points[1]}
	next.Mappings = []*schema.Mapping{projection.Mappings[1]}
	next.Tags = []*schema.Tag{projection.Tags[1]}
	return &next
}

func cloneRuntimeProjection(projection *workspace.RuntimeProjection) *workspace.RuntimeProjection {
	next := *projection
	next.DeviceIDs = append([]string(nil), projection.DeviceIDs...)
	next.Devices = append([]*schema.Device(nil), projection.Devices...)
	return &next
}

func (s *Service) mappingBindingsForPoint(pointID string) []mappingBinding {
	s.mappingMu.RLock()
	defer s.mappingMu.RUnlock()
	return append([]mappingBinding{}, s.mappingIndex[pointID]...)
}
