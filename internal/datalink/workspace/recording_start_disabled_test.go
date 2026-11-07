package workspace

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// Only device activation is simulated. Its original post-Apply scope guard is
// still evaluated against the real saved group and schema services.
type resumeRecordingActivation struct {
	workspace *Service
	calls     int
}

func (a *resumeRecordingActivation) ActivateScope(ctx context.Context, devices, groups []string) (*ActivationResponse, error) {
	a.calls++
	selection, ready, err := a.workspace.activationReadinessScope(ctx, devices, groups)
	if err != nil {
		return nil, err
	}
	if ready == nil || !ready.Ready {
		return nil, &ReadinessBlockedError{Operation: "activation", Summary: ready}
	}
	response := &ActivationResponse{WorkspaceID: selection.record.ID}
	for _, id := range devices {
		response.Results = append(response.Results, ActivationResult{DeviceID: id, Status: ActivationResultStatusSuccess})
	}
	return response, nil
}

func disabledMixedStartFixture(t *testing.T) (*circuitProfileHarness, *WriteGroup, *RecordingStartService, *resumeRecordingActivation) {
	t.Helper()
	ctx := t.Context()
	h := newCircuitProfileHarness(t, filepath.Join(t.TempDir(), "workspace.db"))
	target := filepath.Join(t.TempDir(), "recordings.db")
	external, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	_, err = external.ExecContext(ctx, `PRAGMA user_version=0`)
	require.NoError(t, err)
	require.NoError(t, external.Close())
	setSQLiteConnectorConfig(t, h.db, h.fixture.connectorID, target)
	now := time.Date(2026, 11, 1, 6, 40, 0, 0, time.UTC)
	h.service.WithClock(func() time.Time { return now })
	draft := mixedDraft(h, WriteGroupMemberCalculationMean, 5, "resume_readings")
	draft.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
	draft.Destination.Database = target
	second := h.seedSource(h.fixture.deviceID, "latest_source")
	draft.Members = append(draft.Members, WriteGroupMember{
		DeviceID: second.deviceID, PointID: second.pointID, TagID: second.tagID, Required: true,
		Calculation: WriteGroupMemberCalculationLatest,
		Cadence:     &WriteGroupMemberCadence{Interval: 10, Unit: WriteGroupIntervalUnitSecond, Alignment: WriteGroupAlignmentUTCEpoch},
	})
	created, err := h.service.Create(ctx, h.mutation(ctx, draft))
	require.NoError(t, err)
	confirmed := mixedHistoryConfirm(t, h, created.Group)
	applied, err := h.service.Apply(ctx, confirmed.ID, mixedHistoryRequest(t, h, confirmed))
	require.NoError(t, err)
	renamed := cloneWriteGroup(applied.Group)
	renamed.Name = "Reviewed name only"
	edit := mixedHistoryRequest(t, h, applied.Group)
	edit.Group = renamed
	saved, err := h.service.Update(ctx, renamed.ID, edit)
	require.NoError(t, err)
	now = now.Add(20 * time.Second)
	stopped, err := h.service.Disable(ctx, saved.Group.ID, mixedHistoryRequest(t, h, saved.Group))
	require.NoError(t, err)
	group := mixedHistoryConfirm(t, h, stopped.Group)
	require.Equal(t, WriteGroupStatusDisabled, group.Status)
	require.Equal(t, applied.Group.AppliedRevision, group.AppliedRevision)
	var statements int
	require.NoError(t, h.db.QueryRowContext(ctx, `SELECT executed_statements FROM managed_schema_operations WHERE operation_id=?`, group.Destination.SchemaRevision).Scan(&statements))
	require.Zero(t, statements, "explicit reconfirmation is a real no-DDL operation")
	devices := &stubActivationDeviceService{
		devices: map[string]*schema.Device{h.fixture.deviceID: {ID: h.fixture.deviceID, Status: schema.DeviceStatusActive}},
		readiness: map[string]*schema.DeviceReadiness{h.fixture.deviceID: {
			DeviceID: h.fixture.deviceID, ActivationAllowed: true, ApplyAllowed: true, PlanningAllowed: true,
		}},
	}
	h.workspace.WithReadinessServices(devices, nil, nil, nil).WithWriteGroupReadiness(h.service)
	activation := &resumeRecordingActivation{workspace: h.workspace}
	start := NewRecordingStartService(h.workspace, h.service, activation, recordingplan.NewService(recordingplan.NewSQLRepository(h.db)))
	return h, group, start, activation
}

func disabledStartRequest(t *testing.T, h *circuitProfileHarness, group *WriteGroup) RecordingStartRequest {
	t.Helper()
	mutation := mixedHistoryRequest(t, h, group)
	return RecordingStartRequest{RequestID: "explicit-resume", WorkspaceID: group.WorkspaceID,
		ExpectedWorkspaceRevision: mutation.ExpectedWorkspaceRevision, DeviceIDs: []string{h.fixture.deviceID},
		Groups: []RecordingStartGroupIntent{{GroupID: group.ID, ExpectedGroupRevision: group.Revision,
			ExpectedConnectorRevision: group.Destination.ConnectorRevision}}}
}

func TestRecordingStartResumesDisabledMixedGroup(t *testing.T) {
	h, group, start, activation := disabledMixedStartFixture(t)
	ctx := t.Context()
	request := disabledStartRequest(t, h, group)
	_, err := h.workspace.ReadinessScope(ctx, request.DeviceIDs, []string{group.ID})
	require.ErrorIs(t, err, ErrReadinessScopeInvalid, "ordinary readiness must not grant disabled scope")
	_, _, err = h.workspace.activationReadinessScope(ctx, request.DeviceIDs, []string{group.ID})
	require.ErrorIs(t, err, ErrReadinessScopeInvalid, "activation still requires the original Apply")
	before := h.countRows(`SELECT COUNT(*) FROM write_group_versions WHERE group_id=?`, group.ID)
	op, err := start.Start(ctx, request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, op.Status, "%+v", op)
	require.Equal(t, "complete", op.Stage)
	require.True(t, op.Groups[0].Ready)
	require.True(t, op.Groups[0].Applied)
	require.True(t, op.Devices[0].Activated)
	require.Equal(t, 1, activation.calls)
	current, err := h.service.Get(ctx, group.ID)
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusReady, current.Group.Status)
	require.Equal(t, group.Revision, current.Group.AppliedRevision)
	require.Equal(t, "Reviewed name only", current.Group.Name)
	require.Equal(t, group.Members, current.Group.Members)
	require.Equal(t, before+1, h.countRows(`SELECT COUNT(*) FROM write_group_versions WHERE group_id=?`, group.ID))
	old, err := h.service.repo.versionByRevisionInTx(ctx, h.db, group.WorkspaceID, group.ID, group.AppliedRevision)
	require.NoError(t, err)
	require.Equal(t, time.Date(2026, 11, 1, 6, 40, 0, 0, time.UTC), old.EffectiveAt)
	replayed, err := start.Start(ctx, request)
	require.NoError(t, err)
	require.Equal(t, op.OperationID, replayed.OperationID)
	require.Equal(t, 1, activation.calls, "completed replay cannot apply or activate again")
	require.Equal(t, before+1, h.countRows(`SELECT COUNT(*) FROM write_group_versions WHERE group_id=?`, group.ID))
}

func TestRecordingStartDisabledScopeRejectsInvalidIntent(t *testing.T) {
	for _, kind := range []string{"deleted", "foreign workspace", "foreign device", "stale group", "stale workspace"} {
		t.Run(kind, func(t *testing.T) {
			h, group, start, activation := disabledMixedStartFixture(t)
			if kind == "deleted" {
				group = seedLegacyRecordingTombstone(t, h.db, h.service, group.ID).Group
			}
			request := disabledStartRequest(t, h, group)
			switch kind {
			case "foreign workspace":
				request.WorkspaceID = "foreign-workspace"
			case "foreign device":
				request.DeviceIDs = []string{"foreign-device"}
			case "stale group":
				request.Groups[0].ExpectedGroupRevision = "stale-group-revision"
			case "stale workspace":
				request.ExpectedWorkspaceRevision = "stale-workspace-revision"
			}
			before := h.countRows(`SELECT COUNT(*) FROM write_group_versions WHERE group_id=?`, group.ID)
			op, err := start.Start(t.Context(), request)
			if kind == "stale group" || kind == "stale workspace" {
				require.NoError(t, err)
				require.Equal(t, recordingplan.SchemaOperationFailed, op.Status)
				require.Equal(t, "stale_intent", op.Reason)
			} else {
				require.ErrorIs(t, err, ErrRecordingStartInvalid)
			}
			require.Zero(t, activation.calls)
			require.Equal(t, before, h.countRows(`SELECT COUNT(*) FROM write_group_versions WHERE group_id=?`, group.ID))
			current, err := h.service.Get(t.Context(), group.ID)
			require.NoError(t, err)
			require.Equal(t, group, current.Group, "denial must preserve disabled/deleted intent and history")
		})
	}
}
