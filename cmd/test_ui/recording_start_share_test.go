package main

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
)

func TestRecordingStartShareBarrierChecksOnlySelectedDesiredProjection(t *testing.T) {
	desired := []modbusshare.DesiredMapping{{WorkspaceID: "ws-1", SourceRuleID: "rule-a", SourceRuleRevision: "rev-a", TagID: "tag-a"}}
	source := &recordingStartShareSourceStub{desired: desired, sourceRules: []*schema.SourceRule{{ID: "rule-a", DeviceID: "device-a"}}}
	share := &recordingStartShareServiceStub{settings: modbusshare.Settings{Enabled: true}}
	barrier := newRecordingStartShareBarrier(source, share)

	require.NoError(t, barrier(t.Context(), workspace.RecordingStartRequest{
		WorkspaceID: "ws-1",
		DeviceIDs:   []string{"device-a"},
		Groups:      []workspace.RecordingStartGroupIntent{{GroupID: "group-a"}},
	}))
	require.Equal(t, []string{"device-a"}, source.deviceIDs)
	require.Equal(t, desired, share.checked)
	require.Equal(t, []string{"device-a"}, share.checkedDeviceIDs)
	require.Equal(t, source.sourceRules, share.checkedSourceRules)
}

func TestRecordingStartShareBarrierRejectsSelectedProjectionMismatch(t *testing.T) {
	source := &recordingStartShareSourceStub{desired: []modbusshare.DesiredMapping{{WorkspaceID: "ws-1", TagID: "tag-a", SourceRuleID: "rule-a", SourceRuleRevision: "rev-a"}}}
	share := &recordingStartShareServiceStub{
		settings: modbusshare.Settings{Enabled: true},
		checkErr: modbusshare.NewError(modbusshare.ErrCodeProjectionRequired, "selected projection is stale", true),
	}
	barrier := newRecordingStartShareBarrier(source, share)

	err := barrier(t.Context(), workspace.RecordingStartRequest{WorkspaceID: "ws-1", DeviceIDs: []string{"device-a"}, Groups: []workspace.RecordingStartGroupIntent{{GroupID: "group-a"}}})
	require.Error(t, err)
	require.True(t, errors.Is(err, workspace.ErrRecordingStartShareNotReady))
	require.True(t, share.checkedCalled)
}

func TestRecordingStartShareBarrierRejectsEmptyShareOnlyProjection(t *testing.T) {
	source := &recordingStartShareSourceStub{desired: []modbusshare.DesiredMapping{}}
	share := &recordingStartShareServiceStub{settings: modbusshare.Settings{Enabled: true}}
	barrier := newRecordingStartShareBarrier(source, share)

	err := barrier(t.Context(), workspace.RecordingStartRequest{WorkspaceID: "ws-1", DeviceIDs: []string{"device-a"}})
	require.ErrorIs(t, err, workspace.ErrRecordingStartShareNotReady)
	require.False(t, share.checkedCalled, "empty Share-only scope must fail before a projection claim")
}

func TestRecordingStartShareBarrierAllowsEmptyDesiredForGroupScope(t *testing.T) {
	source := &recordingStartShareSourceStub{desired: []modbusshare.DesiredMapping{}}
	share := &recordingStartShareServiceStub{settings: modbusshare.Settings{Enabled: true}}
	barrier := newRecordingStartShareBarrier(source, share)

	require.NoError(t, barrier(t.Context(), workspace.RecordingStartRequest{
		WorkspaceID: "ws-1",
		DeviceIDs:   []string{"device-a"},
		Groups:      []workspace.RecordingStartGroupIntent{{GroupID: "group-a"}},
	}))
	require.True(t, share.checkedCalled)
}

func TestWorkspaceShareRestoreBarrierRetainsFullWorkspaceRestore(t *testing.T) {
	workspaceService := &workspaceShareWorkspaceStub{record: &workspace.Record{ID: "ws-1", OrderedDeviceIDs: []string{"device-a", "device-b"}}}
	source := &workspaceShareRestoreSourceStub{outcome: modbusshare.ReconcileOutcome{Outcome: shareOutcomeApplied}}
	share := &recordingStartShareServiceStub{
		settings:  modbusshare.Settings{Enabled: true, SettingsRevision: "settings-1"},
		hydration: modbusshare.HydrationState{State: modbusshare.HydrationStateReady, WorkspaceID: "ws-1", WorkspaceRevision: "workspace-1", Readiness: true},
	}
	reconciler := &workspaceShareReconcilerStub{}
	barrier := newWorkspaceShareRestoreBarrier(workspaceService, source, share, reconciler)

	require.NoError(t, barrier(t.Context(), handlers.ActivateWorkspaceRequest{WorkspaceRevision: "workspace-1"}))
	require.Equal(t, "ws-1", source.workspaceID)
	require.Equal(t, "workspace-1", source.workspaceRevision)
	require.Equal(t, []string{"device-a", "device-b"}, source.deviceIDs)
	require.Equal(t, share.settings, source.settings)
	require.Same(t, reconciler, source.reconciler)
}

type recordingStartShareSourceStub struct {
	desired     []modbusshare.DesiredMapping
	sourceRules []*schema.SourceRule
	err         error
	workspace   string
	deviceIDs   []string
}

func (s *recordingStartShareSourceStub) BuildDesiredShareMappingsForDevices(_ context.Context, workspaceID string, _ modbusshare.Settings, deviceIDs []string) ([]modbusshare.DesiredMapping, error) {
	s.workspace, s.deviceIDs = workspaceID, append([]string(nil), deviceIDs...)
	return s.desired, s.err
}

func (s *recordingStartShareSourceStub) List(context.Context, sourcerule.ListFilter) ([]*schema.SourceRule, error) {
	return s.sourceRules, s.err
}

type recordingStartShareServiceStub struct {
	settings           modbusshare.Settings
	hydration          modbusshare.HydrationState
	hydrationErr       error
	checkErr           error
	checked            []modbusshare.DesiredMapping
	checkedDeviceIDs   []string
	checkedSourceRules []*schema.SourceRule
	checkedCalled      bool
}

func (s *recordingStartShareServiceStub) Settings() modbusshare.Settings {
	return s.settings
}

func (s *recordingStartShareServiceStub) CheckHydration(context.Context) (modbusshare.HydrationState, error) {
	return s.hydration, s.hydrationErr
}

func (s *recordingStartShareServiceStub) CheckDesiredProjectionForScopeWithSourceRules(_ context.Context, _ string, deviceIDs []string, desired []modbusshare.DesiredMapping, sourceRules []*schema.SourceRule) error {
	s.checkedCalled = true
	s.checked = append([]modbusshare.DesiredMapping(nil), desired...)
	s.checkedDeviceIDs = append([]string(nil), deviceIDs...)
	s.checkedSourceRules = append([]*schema.SourceRule(nil), sourceRules...)
	return s.checkErr
}

type workspaceShareWorkspaceStub struct {
	record *workspace.Record
	err    error
}

func (s *workspaceShareWorkspaceStub) GetOrCreate(context.Context) (*workspace.Record, error) {
	return s.record, s.err
}

type workspaceShareRestoreSourceStub struct {
	outcome           modbusshare.ReconcileOutcome
	err               error
	workspaceID       string
	workspaceRevision string
	settings          modbusshare.Settings
	deviceIDs         []string
	reconciler        sourcerule.ShareProjectionReconciler
}

func (s *workspaceShareRestoreSourceStub) RestoreLocalModbusProjectionForDevices(_ context.Context, workspaceID, workspaceRevision string, settings modbusshare.Settings, deviceIDs []string, reconciler sourcerule.ShareProjectionReconciler) (modbusshare.ReconcileOutcome, error) {
	s.workspaceID, s.workspaceRevision, s.settings = workspaceID, workspaceRevision, settings
	s.deviceIDs, s.reconciler = append([]string(nil), deviceIDs...), reconciler
	return s.outcome, s.err
}

type workspaceShareReconcilerStub struct{}

func (workspaceShareReconcilerStub) Reconcile(context.Context, modbusshare.ReconcileRequest) (modbusshare.ReconcileOutcome, error) {
	return modbusshare.ReconcileOutcome{}, nil
}
