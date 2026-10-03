package workspace

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type scopedActivationRuntimeSyncer struct {
	stubActivationRuntimeSyncer
	upsertedScopes      [][]string
	upsertedPointScopes [][]string
	fullCalls           int
	applyErr            error
}

func (s *scopedActivationRuntimeSyncer) ApplyWorkspaceProjection(_ context.Context, _ *RuntimeProjection) error {
	s.fullCalls++
	return nil
}

func (s *scopedActivationRuntimeSyncer) ApplyWorkspaceProjectionForDevices(_ context.Context, _ *RuntimeProjection, deviceIDs []string) error {
	s.upsertedScopes = append(s.upsertedScopes, append([]string{}, deviceIDs...))
	return s.applyErr
}

func (s *scopedActivationRuntimeSyncer) ApplyWorkspaceProjectionForDevicesAndPoints(_ context.Context, _ *RuntimeProjection, deviceIDs, pointIDs []string) error {
	s.upsertedScopes = append(s.upsertedScopes, append([]string{}, deviceIDs...))
	s.upsertedPointScopes = append(s.upsertedPointScopes, append([]string{}, pointIDs...))
	return s.applyErr
}

func TestActivationService_ActivateScopeDoesNotBlockOnUnselectedDevice(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A", "dev-B"))
	deviceSvc := scopeActivationDevices()
	deviceSvc.readiness["dev-B"] = &schema.DeviceReadiness{DeviceID: "dev-B", ActivationAllowed: false, BlockingReasons: []string{"incomplete"}}
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil)

	service := NewActivationService(workspaceSvc, deviceSvc)
	response, err := service.ActivateScope(ctx, []string{"dev-A"}, nil)
	require.NoError(t, err)
	require.Len(t, response.Results, 1)
	require.Equal(t, "dev-A", response.Results[0].DeviceID)
	require.Equal(t, ActivationResultStatusSuccess, response.Results[0].Status)
	require.Equal(t, []string{"dev-A"}, deviceSvc.activateCalls)
}

func TestActivationService_ActivateScopePreservesHealthyUnrelatedDevice(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A", "dev-B"))
	deviceSvc := scopeActivationDevices()
	deviceSvc.devices["dev-B"].Status = schema.DeviceStatusActive
	deviceSvc.readiness["dev-B"] = &schema.DeviceReadiness{DeviceID: "dev-B", ActivationAllowed: true}
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil)
	workspaceSvc.WithRuntimeProjectionServices(
		deviceSvc,
		&runtimeProjectionRuleStub{},
		&runtimeProjectionPointStub{},
		&runtimeProjectionMappingStub{},
		&runtimeProjectionTagStub{},
		nil,
		nil,
		&runtimeProjectionGroupStub{},
	)
	runtimeSync := &scopedActivationRuntimeSyncer{}

	service := NewActivationService(workspaceSvc, deviceSvc, runtimeSync)
	response, err := service.ActivateScope(ctx, []string{"dev-A"}, nil)
	require.NoError(t, err)
	require.Len(t, response.Results, 1)
	require.Equal(t, ActivationResultStatusSuccess, response.Results[0].Status)
	require.Equal(t, []string{"dev-A"}, deviceSvc.activateCalls)
	require.Equal(t, [][]string{{"dev-A"}}, runtimeSync.upsertedScopes)
	require.Zero(t, runtimeSync.fullCalls)
	require.Empty(t, runtimeSync.upserted)
}

func TestActivationService_ActivateScopeAlreadyActiveIsIdempotent(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A", "dev-B"))
	deviceSvc := scopeActivationDevices()
	deviceSvc.devices["dev-A"].Status = schema.DeviceStatusActive
	deviceSvc.devices["dev-B"].Status = schema.DeviceStatusActive
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil)
	workspaceSvc.WithRuntimeProjectionServices(
		deviceSvc,
		&runtimeProjectionRuleStub{},
		&runtimeProjectionPointStub{},
		&runtimeProjectionMappingStub{},
		&runtimeProjectionTagStub{},
		nil,
		nil,
		&runtimeProjectionGroupStub{},
	)
	runtimeSync := &scopedActivationRuntimeSyncer{}

	service := NewActivationService(workspaceSvc, deviceSvc, runtimeSync)
	response, err := service.ActivateScope(ctx, []string{"dev-A"}, nil)
	require.NoError(t, err)
	require.Len(t, response.Results, 1)
	require.Equal(t, ActivationResultStatusSuccess, response.Results[0].Status)
	require.Empty(t, deviceSvc.activateCalls)
	require.Empty(t, runtimeSync.upserted)
	require.Equal(t, [][]string{{"dev-A"}}, runtimeSync.upsertedScopes)
	require.Zero(t, runtimeSync.fullCalls)
}

func TestActivationService_ActivateScopeProjectionFailureDoesNotDisableAlreadyActive(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A"))
	deviceSvc := scopeActivationDevices()
	deviceSvc.devices["dev-A"].Status = schema.DeviceStatusActive
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil)
	workspaceSvc.WithRuntimeProjectionServices(
		deviceSvc,
		&runtimeProjectionRuleStub{},
		&runtimeProjectionPointStub{},
		&runtimeProjectionMappingStub{},
		&runtimeProjectionTagStub{},
		nil,
		nil,
		&runtimeProjectionGroupStub{},
	)
	runtimeSync := &scopedActivationRuntimeSyncer{applyErr: errors.New("projection unavailable")}

	service := NewActivationService(workspaceSvc, deviceSvc, runtimeSync)
	response, err := service.ActivateScope(ctx, []string{"dev-A"}, nil)
	require.NoError(t, err)
	require.Len(t, response.Results, 1)
	require.Equal(t, ActivationResultStatusFailed, response.Results[0].Status)
	require.Empty(t, deviceSvc.activateCalls)
	require.Empty(t, deviceSvc.disableCalls)
	require.Equal(t, schema.DeviceStatusActive, deviceSvc.devices["dev-A"].Status)
}

func TestActivationService_ActivateScopeRequiresAppliedSelectedGroup(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A"))
	record, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	deviceSvc := scopeActivationDevices()
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

	service := NewActivationService(workspaceSvc, deviceSvc)
	_, err = service.ActivateScope(ctx, []string{"dev-A"}, []string{"group-A"})
	require.Error(t, err)
	var blockedErr *ReadinessBlockedError
	require.True(t, errors.As(err, &blockedErr))
	require.Empty(t, deviceSvc.activateCalls)
}

func TestActivationService_ActivateScopePassesSelectedGroupToRuntimeProjection(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A"))
	record, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	deviceSvc := scopeActivationDevices()
	groups := &scopedReadinessGroupService{
		listed: &WriteGroupListResult{
			WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision,
			Groups: []*WriteGroup{
				{ID: "write-group-A", WorkspaceID: record.ID, Revision: "group-rev-A", AppliedRevision: "applied-A", Members: []WriteGroupMember{{DeviceID: "dev-A", PointID: "point-A"}}},
				{ID: "write-group-B", WorkspaceID: record.ID, Revision: "group-rev-B", AppliedRevision: "applied-B", Members: []WriteGroupMember{{DeviceID: "dev-A", PointID: "point-B"}}},
			},
		},
		readiness: map[string]*WriteGroupReadiness{
			"write-group-A": {WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, GroupID: "write-group-A", GroupRevision: "group-rev-A", AppliedRevision: "applied-A", Ready: true},
		},
	}
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil).WithWriteGroupReadiness(groups)
	workspaceSvc.WithRuntimeProjectionServices(
		&runtimeProjectionDeviceStub{records: map[string]*schema.Device{"dev-A": deviceSvc.devices["dev-A"]}},
		&runtimeProjectionRuleStub{}, &runtimeProjectionPointStub{}, &runtimeProjectionMappingStub{}, &runtimeProjectionTagStub{}, nil, nil,
		&runtimeProjectionGroupStub{records: []*schema.PollingGroup{{ID: "poll-A"}, {ID: "poll-B"}}},
	)
	runtimeSync := &scopedActivationRuntimeSyncer{}

	service := NewActivationService(workspaceSvc, deviceSvc, runtimeSync)
	response, err := service.ActivateScope(ctx, []string{"dev-A"}, []string{"write-group-A"})
	require.NoError(t, err)
	require.Len(t, response.Results, 1)
	require.Equal(t, ActivationResultStatusSuccess, response.Results[0].Status)
	require.Equal(t, [][]string{{"dev-A"}}, runtimeSync.upsertedScopes)
	require.Equal(t, [][]string{{"point-A"}}, runtimeSync.upsertedPointScopes)
}

func TestActivationService_ActivateScopeRejectsUnknownDeviceBeforeMutation(t *testing.T) {
	ctx := t.Context()
	workspaceSvc := NewService(NewMemoryRepository())
	require.NoError(t, attachScopeDevices(ctx, workspaceSvc, "dev-A"))
	deviceSvc := scopeActivationDevices()
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil)

	service := NewActivationService(workspaceSvc, deviceSvc)
	_, err := service.ActivateScope(ctx, []string{"foreign"}, nil)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrReadinessScopeInvalid))
	require.Empty(t, deviceSvc.activateCalls)
}

func scopeActivationDevices() *stubActivationDeviceService {
	return &stubActivationDeviceService{
		devices: map[string]*schema.Device{
			"dev-A": {ID: "dev-A", Name: "Device A", Status: schema.DeviceStatusDraft},
			"dev-B": {ID: "dev-B", Name: "Device B", Status: schema.DeviceStatusDraft},
		},
		readiness: map[string]*schema.DeviceReadiness{
			"dev-A": {DeviceID: "dev-A", ActivationAllowed: true},
			"dev-B": {DeviceID: "dev-B", ActivationAllowed: true},
		},
		activateErrs: map[string]error{},
	}
}
