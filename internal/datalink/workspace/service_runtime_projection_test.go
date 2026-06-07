package workspace

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
)

func TestService_RuntimeProjectionRebuildsSamePersistedWorkspaceForActivationAndRestart(t *testing.T) {
	ctx := context.Background()
	workspaceSvc := NewService(NewMemoryRepository())
	workspaceSvc.now = func() time.Time { return time.Unix(100, 0).UTC() }
	workspaceSvc.newID = func() string { return "ws-1" }

	devices := runtimeProjectionDeviceStub{records: map[string]*schema.Device{
		"dev-A": {ID: "dev-A", Name: "Device A", Status: schema.DeviceStatusActive},
		"dev-B": {ID: "dev-B", Name: "Device B", Status: schema.DeviceStatusActive},
		"dev-Z": {ID: "dev-Z", Name: "Outside workspace", Status: schema.DeviceStatusActive},
	}}
	rules := runtimeProjectionRuleStub{
		rules: []*schema.SourceRule{
			{ID: "rule-B", DeviceID: "dev-B", StartAddress: "40010", Count: 1, DataType: schema.DataTypeInt16, Enabled: true},
			{ID: "rule-A", DeviceID: "dev-A", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt16, Enabled: true},
			{ID: "rule-Z", DeviceID: "dev-Z", StartAddress: "40100", Count: 1, DataType: schema.DataTypeInt16, Enabled: true},
		},
		links: map[string][]*schema.SourceRuleLink{
			"rule-A": {{ID: "link-A", RuleID: "rule-A", Address: "40001", PointID: "point-A", TagID: stringPtr("tag-A"), MappingID: stringPtr("map-A")}},
			"rule-B": {{ID: "link-B", RuleID: "rule-B", Address: "40010", PointID: "point-B", TagID: stringPtr("tag-B"), MappingID: stringPtr("map-B")}},
		},
	}
	points := runtimeProjectionPointStub{records: []*schema.Point{
		{ID: "point-A", DeviceID: "dev-A", Name: "A", Address: "40001", DataType: schema.DataTypeInt16, Enabled: true, PollingGroupID: stringPtr("group-fast")},
		{ID: "point-B", DeviceID: "dev-B", Name: "B", Address: "40010", DataType: schema.DataTypeInt16, Enabled: true, PollingGroupID: stringPtr("group-fast")},
		{ID: "point-Z", DeviceID: "dev-Z", Name: "Z", Address: "40100", DataType: schema.DataTypeInt16, Enabled: true, PollingGroupID: stringPtr("group-fast")},
	}}
	mappings := runtimeProjectionMappingStub{records: []*schema.Mapping{
		{ID: "map-A", PointID: "point-A", TagID: "tag-A", Enabled: true},
		{ID: "map-B", PointID: "point-B", TagID: "tag-B", Enabled: true},
		{ID: "map-disabled", PointID: "point-B", TagID: "tag-disabled", Enabled: false},
		{ID: "map-Z", PointID: "point-Z", TagID: "tag-Z", Enabled: true},
	}}
	tags := runtimeProjectionTagStub{records: map[string]*schema.Tag{
		"tag-A": {ID: "tag-A", Key: "device.a", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive},
		"tag-B": {ID: "tag-B", Key: "device.b", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive},
	}}
	targets := runtimeProjectionTargetStub{records: []*schema.DatabaseTargetMapping{
		{ID: "target-A", ConnectorID: "db-main", TagID: "tag-A", TableName: "readings", ColumnName: "a", Enabled: true},
		{ID: "target-B", ConnectorID: "db-main", TagID: "tag-B", TableName: "readings", ColumnName: "b", Enabled: true},
		{ID: "target-Z", ConnectorID: "db-main", TagID: "tag-Z", TableName: "readings", ColumnName: "z", Enabled: true},
	}}
	groups := runtimeProjectionGroupStub{records: []*schema.PollingGroup{
		{ID: "group-fast", Name: "fast", IntervalMs: 1000, Enabled: true},
	}}

	workspaceSvc.WithRuntimeProjectionServices(&devices, &rules, &points, &mappings, &tags, nil, &targets, &groups)
	if _, err := workspaceSvc.AttachDevice(ctx, "dev-A"); err != nil {
		t.Fatalf("attach dev-A failed: %v", err)
	}
	if _, err := workspaceSvc.AttachDevice(ctx, "dev-B"); err != nil {
		t.Fatalf("attach dev-B failed: %v", err)
	}
	if _, err := workspaceSvc.BindDatabaseConnector(ctx, "db-main"); err != nil {
		t.Fatalf("bind database connector failed: %v", err)
	}

	activationProjection, err := workspaceSvc.RuntimeProjection(ctx)
	if err != nil {
		t.Fatalf("build activation projection failed: %v", err)
	}
	restartProjection, err := workspaceSvc.RuntimeProjection(ctx)
	if err != nil {
		t.Fatalf("build restart projection failed: %v", err)
	}

	if activationProjection.WorkspaceID != "ws-1" {
		t.Fatalf("expected workspace ws-1, got %q", activationProjection.WorkspaceID)
	}
	if activationProjection.Version == "" {
		t.Fatal("expected stable non-empty projection version")
	}
	if activationProjection.Alignment != RuntimeProjectionAlignmentAligned {
		t.Fatalf("expected aligned projection, got %q", activationProjection.Alignment)
	}
	if activationProjection.Version != restartProjection.Version {
		t.Fatalf("expected same projection version, got %q and %q", activationProjection.Version, restartProjection.Version)
	}

	got := summarizeRuntimeProjection(activationProjection)
	want := runtimeProjectionSummary{
		DeviceIDs:      []string{"dev-A", "dev-B"},
		RuleIDs:        []string{"rule-A", "rule-B"},
		RuleLinkIDs:    []string{"link-A", "link-B"},
		PointIDs:       []string{"point-A", "point-B"},
		MappingIDs:     []string{"map-A", "map-B"},
		TagIDs:         []string{"tag-A", "tag-B"},
		DatabaseTarget: []string{"target-A", "target-B"},
		PollingGroups:  []string{"group-fast"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected projection summary\nwant: %+v\n got: %+v", want, got)
	}
	if restartGot := summarizeRuntimeProjection(restartProjection); !reflect.DeepEqual(restartGot, got) {
		t.Fatalf("restart projection drifted\nactivation: %+v\n   restart: %+v", got, restartGot)
	}
}
