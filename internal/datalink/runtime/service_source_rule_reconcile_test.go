package runtime

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

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
}

func (s *Service) mappingBindingsForPoint(pointID string) []mappingBinding {
	s.mappingMu.RLock()
	defer s.mappingMu.RUnlock()
	return append([]mappingBinding{}, s.mappingIndex[pointID]...)
}
