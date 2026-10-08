package workspace

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type scopedReadinessGroupService struct {
	listed    *WriteGroupListResult
	readiness map[string]*WriteGroupReadiness
	calls     []string
}

func (s *scopedReadinessGroupService) List(context.Context) (*WriteGroupListResult, error) {
	return s.listed, nil
}

func (s *scopedReadinessGroupService) Readiness(_ context.Context, id string) (*WriteGroupReadiness, error) {
	s.calls = append(s.calls, id)
	return s.readiness[id], nil
}

type scopedReadinessRules struct {
	links []*schema.SourceRuleLink
}

func (s scopedReadinessRules) ListByDeviceIDs(context.Context, []string) ([]*schema.SourceRule, error) {
	return []*schema.SourceRule{{ID: "rule-A", DeviceID: "dev-A"}}, nil
}

func (s scopedReadinessRules) ListLinks(context.Context, string) ([]*schema.SourceRuleLink, error) {
	return s.links, nil
}

func TestService_ReadinessScopeIgnoresIncompleteUnselectedDevice(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A", "dev-B"))
	deviceSvc := &stubActivationDeviceService{
		devices: map[string]*schema.Device{
			"dev-A": {ID: "dev-A", Status: schema.DeviceStatusDraft},
			"dev-B": {ID: "dev-B", Status: schema.DeviceStatusDraft},
		},
		readiness: map[string]*schema.DeviceReadiness{
			"dev-A": {DeviceID: "dev-A", ActivationAllowed: true},
			"dev-B": {DeviceID: "dev-B", ActivationAllowed: false, BlockingReasons: []string{"incomplete"}},
		},
	}
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil)

	summary, err := workspaceSvc.ReadinessScope(ctx, []string{"dev-A"}, nil)
	require.NoError(t, err)
	require.True(t, summary.Ready)
	require.Zero(t, summary.BlockingCount)
	require.Empty(t, summary.Issues)
}

func TestService_ReadinessScopeEvaluatesOnlyOwnedSelectedGroup(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A", "dev-B"))
	record, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	deviceSvc := &stubActivationDeviceService{
		devices: map[string]*schema.Device{
			"dev-A": {ID: "dev-A", Status: schema.DeviceStatusDraft},
			"dev-B": {ID: "dev-B", Status: schema.DeviceStatusDraft},
		},
		readiness: map[string]*schema.DeviceReadiness{
			"dev-A": {DeviceID: "dev-A", ActivationAllowed: true},
			"dev-B": {DeviceID: "dev-B", ActivationAllowed: false, BlockingReasons: []string{"incomplete"}},
		},
	}
	groups := &scopedReadinessGroupService{
		listed: &WriteGroupListResult{
			WorkspaceID:       record.ID,
			WorkspaceRevision: record.DatabaseSetupRevision,
			Groups: []*WriteGroup{
				{ID: "group-A", WorkspaceID: record.ID, Revision: "group-rev-A", AppliedRevision: "applied-A", Members: []WriteGroupMember{{DeviceID: "dev-A"}}},
				{ID: "group-B", WorkspaceID: record.ID, Revision: "group-rev-B", Members: []WriteGroupMember{{DeviceID: "dev-B"}}},
			},
		},
		readiness: map[string]*WriteGroupReadiness{
			"group-A": {WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, GroupID: "group-A", GroupRevision: "group-rev-A", AppliedRevision: "applied-A", Ready: true},
			"group-B": {WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, GroupID: "group-B", GroupRevision: "group-rev-B", Ready: false, Issues: []ReadinessIssue{{Code: "incomplete", Severity: ReadinessSeverityBlocking}}},
		},
	}
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil).WithWriteGroupReadiness(groups)

	summary, err := workspaceSvc.ReadinessScope(ctx, []string{"dev-A"}, []string{"group-A"})
	require.NoError(t, err)
	require.True(t, summary.Ready)
	require.Equal(t, []string{"group-A"}, groups.calls)
}

func TestService_ReadinessScopeIgnoresIncompleteUnselectedGroupOnSelectedDevice(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A"))
	record, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	deviceSvc := &stubActivationDeviceService{
		devices:   map[string]*schema.Device{"dev-A": {ID: "dev-A", Status: schema.DeviceStatusDraft}},
		readiness: map[string]*schema.DeviceReadiness{"dev-A": {DeviceID: "dev-A", ActivationAllowed: true}},
	}
	groups := &scopedReadinessGroupService{
		listed: &WriteGroupListResult{
			WorkspaceID:       record.ID,
			WorkspaceRevision: record.DatabaseSetupRevision,
			Groups: []*WriteGroup{
				{ID: "group-A", WorkspaceID: record.ID, Revision: "group-rev-A", AppliedRevision: "applied-A", Members: []WriteGroupMember{{DeviceID: "dev-A", PointID: "point-A"}}},
				{ID: "group-B", WorkspaceID: record.ID, Revision: "group-rev-B", Members: []WriteGroupMember{{DeviceID: "dev-A", PointID: "point-B"}}},
			},
		},
		readiness: map[string]*WriteGroupReadiness{
			"group-A": {WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, GroupID: "group-A", GroupRevision: "group-rev-A", AppliedRevision: "applied-A", Ready: true},
			"group-B": {WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, GroupID: "group-B", GroupRevision: "group-rev-B", Ready: false, Issues: []ReadinessIssue{{Code: "incomplete", Severity: ReadinessSeverityBlocking}}},
		},
	}
	rules := scopedReadinessRulesByDevice{
		rules: map[string][]*schema.SourceRule{
			"dev-A": {{ID: "rule-A", DeviceID: "dev-A"}, {ID: "rule-B", DeviceID: "dev-A"}},
		},
		links: map[string][]*schema.SourceRuleLink{
			"rule-A": {{RuleID: "rule-A", PointID: "point-A", TagID: common.Ptr("tag-A"), MappingID: common.Ptr("mapping-A")}},
			"rule-B": {{RuleID: "rule-B", PointID: "point-B", TagID: common.Ptr("tag-B")}},
		},
	}
	workspaceSvc.WithReadinessServices(deviceSvc, rules, nil, nil).WithWriteGroupReadiness(groups)

	summary, err := workspaceSvc.ReadinessScope(ctx, []string{"dev-A"}, []string{"group-A"})
	require.NoError(t, err)
	require.True(t, summary.Ready)
	require.Zero(t, summary.BlockingCount)
}

func TestService_ReadinessScopeAllowsConfirmedGroupBeforeApply(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A"))
	record, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	deviceSvc := &stubActivationDeviceService{
		devices:   map[string]*schema.Device{"dev-A": {ID: "dev-A", Status: schema.DeviceStatusDraft}},
		readiness: map[string]*schema.DeviceReadiness{"dev-A": {DeviceID: "dev-A", ActivationAllowed: true}},
	}
	groups := &scopedReadinessGroupService{
		listed: &WriteGroupListResult{
			WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision,
			Groups: []*WriteGroup{{ID: "group-A", WorkspaceID: record.ID, Revision: "group-rev-A", Members: []WriteGroupMember{{DeviceID: "dev-A"}}}},
		},
		readiness: map[string]*WriteGroupReadiness{
			"group-A": {WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, GroupID: "group-A", GroupRevision: "group-rev-A", Ready: true},
		},
	}
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil).WithWriteGroupReadiness(groups)

	summary, err := workspaceSvc.ReadinessScope(ctx, []string{"dev-A"}, []string{"group-A"})
	require.NoError(t, err)
	require.True(t, summary.Ready)
	require.NotContains(t, readinessIssueCodes(summary.Issues), "write-group-apply-required")
}

func TestService_ReadinessScopeRejectsForeignGroupAndDoesNotMutate(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A"))
	record, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	deviceSvc := &stubActivationDeviceService{
		devices:   map[string]*schema.Device{"dev-A": {ID: "dev-A", Status: schema.DeviceStatusDraft}},
		readiness: map[string]*schema.DeviceReadiness{"dev-A": {DeviceID: "dev-A", ActivationAllowed: true}},
	}
	groups := &scopedReadinessGroupService{listed: &WriteGroupListResult{
		WorkspaceID:       record.ID,
		WorkspaceRevision: record.DatabaseSetupRevision,
		Groups:            []*WriteGroup{{ID: "foreign", WorkspaceID: "other-workspace", Revision: "rev", Members: []WriteGroupMember{{DeviceID: "dev-A"}}}},
	}}
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil).WithWriteGroupReadiness(groups)

	_, err = workspaceSvc.ReadinessScope(ctx, []string{"dev-A"}, []string{"foreign"})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrReadinessScopeInvalid))
	require.Empty(t, groups.calls)
}

func TestService_ReadinessScopeShareOnlyDoesNotReadDatabaseTargets(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A"))
	deviceSvc := &stubActivationDeviceService{
		devices:   map[string]*schema.Device{"dev-A": {ID: "dev-A", Status: schema.DeviceStatusDraft}},
		readiness: map[string]*schema.DeviceReadiness{"dev-A": {DeviceID: "dev-A", ActivationAllowed: true}},
	}
	workspaceSvc.WithReadinessServices(deviceSvc, scopedReadinessRules{links: []*schema.SourceRuleLink{{RuleID: "rule-A", PointID: "point-A", TagID: common.Ptr("tag-A"), MappingID: common.Ptr("mapping-A")}}}, nil, panicReadinessTargetService{})

	summary, err := workspaceSvc.ReadinessScope(ctx, []string{"dev-A"}, nil)
	require.NoError(t, err)
	require.True(t, summary.Ready)
}

type panicReadinessTargetService struct{}

func (panicReadinessTargetService) List(context.Context, dbtarget.TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error) {
	panic("database target service must not be used by scoped Share readiness")
}

type scopedReadinessRulesByDevice struct {
	rules map[string][]*schema.SourceRule
	links map[string][]*schema.SourceRuleLink
}

func (s scopedReadinessRulesByDevice) ListByDeviceIDs(_ context.Context, deviceIDs []string) ([]*schema.SourceRule, error) {
	var result []*schema.SourceRule
	for _, deviceID := range deviceIDs {
		result = append(result, s.rules[deviceID]...)
	}
	return result, nil
}

func (s scopedReadinessRulesByDevice) ListLinks(_ context.Context, ruleID string) ([]*schema.SourceRuleLink, error) {
	return s.links[ruleID], nil
}

func attachScopeDevices(ctx context.Context, service *Service, deviceIDs ...string) error {
	for _, deviceID := range deviceIDs {
		if _, err := service.AttachDevice(ctx, deviceID); err != nil {
			return err
		}
	}
	return nil
}
