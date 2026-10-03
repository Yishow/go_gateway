package api

import (
	"context"
	"testing"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func TestRecordingStartDatabaseOnlyRetainsHydratedWorkspaceBarrier(t *testing.T) {
	for _, validProof := range []bool{false, true} {
		name := "missing-proof"
		if validProof {
			name = "current-proof"
		}
		t.Run(name, func(t *testing.T) {
			f := newRecordingStartFixture(t)
			share := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
			share.SetHydrationState(modbusshare.HydrationState{
				State: modbusshare.HydrationStateReady, Readiness: true,
				WorkspaceID: f.record.ID, WorkspaceRevision: "persisted-share-revision",
				SettingsRevision: share.Settings().SettingsRevision, ReadinessToken: "current-proof",
			})
			restores := 0
			services := &DatalinkServices{
				Workspace: f.workspace, ModbusShare: share,
				ShareRestore: func(context.Context, handlers.ActivateWorkspaceRequest) error {
					restores++
					return nil
				},
			}
			f.start.WithActivationBarrier(recordingStartActivationBarrier(services), recordingStartActivationResumeBarrier(services))
			request := f.request("database-only-start")
			if validProof {
				request.WorkspaceRevision = "persisted-share-revision"
				request.SettingsRevision = share.Settings().SettingsRevision
				request.ReadinessToken = "current-proof"
			}
			operation, err := f.start.Start(t.Context(), request)
			require.NoError(t, err)
			if validProof {
				require.Equal(t, recordingplan.SchemaOperationSucceeded, operation.Status)
				require.Equal(t, 1, f.activation.Calls())
			} else {
				require.Equal(t, recordingplan.SchemaOperationFailed, operation.Status)
				require.Equal(t, "share_not_ready", operation.Reason)
				require.Equal(t, "refresh_share", operation.NextAction)
				require.False(t, operation.Groups[0].Applied)
				require.Zero(t, f.activation.Calls())
			}
			require.Zero(t, restores, "database-only start must not restore the full Share projection")
		})
	}
}

func TestRecordingStartResumeKeepsCurrentShareBarriers(t *testing.T) {
	for _, changed := range []string{"workspace-revision", "settings-revision", "failed-hydration"} {
		t.Run(changed, func(t *testing.T) {
			f := newRecordingStartFixture(t, recordingStartFailure(), recordingStartSuccess())
			share := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
			hydration := modbusshare.HydrationState{
				State: modbusshare.HydrationStateReady, Readiness: true, WorkspaceID: f.record.ID,
				WorkspaceRevision: "persisted-share-revision", SettingsRevision: share.Settings().SettingsRevision,
				ReadinessToken: "before-restart",
			}
			share.SetHydrationState(hydration)
			services := &DatalinkServices{Workspace: f.workspace, ModbusShare: share}
			f.start.WithActivationBarrier(recordingStartActivationBarrier(services), recordingStartActivationResumeBarrier(services))
			request := f.request("resume-with-current-barriers")
			request.WorkspaceRevision, request.SettingsRevision, request.ReadinessToken = hydration.WorkspaceRevision, hydration.SettingsRevision, hydration.ReadinessToken
			partial, err := f.start.Start(t.Context(), request)
			require.NoError(t, err)
			require.Equal(t, recordingplan.SchemaOperationPartial, partial.Status)
			hydration.ReadinessToken = "after-restart"
			switch changed {
			case "workspace-revision":
				hydration.WorkspaceRevision = "another-revision"
			case "settings-revision":
				settings := share.Settings()
				settings.SettingsRevision = "another-settings-revision"
				require.NoError(t, share.ApplySettings(t.Context(), settings))
			case "failed-hydration":
				hydration.State, hydration.Readiness = modbusshare.HydrationStateFailed, false
			}
			share.SetHydrationState(hydration)
			resumed, err := f.start.Start(t.Context(), request)
			require.NoError(t, err)
			require.Equal(t, "share_not_ready", resumed.Reason)
			require.Equal(t, "refresh_share", resumed.NextAction)
			require.Equal(t, partial.OperationID, resumed.OperationID)
			require.Equal(t, partial.Groups[0].AppliedRevision, resumed.Groups[0].AppliedRevision)
			require.Equal(t, 1, f.activation.Calls())
		})
	}
}

func TestRecordingStartPartialResumeReprovesUnchangedHydrationAfterRestart(t *testing.T) {
	f := newRecordingStartFixture(t, recordingStartFailure(), recordingStartSuccess())
	share := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
	hydration := modbusshare.HydrationState{
		State: modbusshare.HydrationStateReady, Readiness: true,
		WorkspaceID: f.record.ID, WorkspaceRevision: "persisted-share-revision",
		SettingsRevision: share.Settings().SettingsRevision, ReadinessToken: "before-restart",
	}
	share.SetHydrationState(hydration)
	services := &DatalinkServices{Workspace: f.workspace, ModbusShare: share}
	f.start.WithActivationBarrier(recordingStartActivationBarrier(services), recordingStartActivationResumeBarrier(services))
	request := f.request("partial-before-restart")
	request.WorkspaceRevision, request.SettingsRevision, request.ReadinessToken = hydration.WorkspaceRevision, hydration.SettingsRevision, hydration.ReadinessToken
	partial, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationPartial, partial.Status)
	require.True(t, partial.Groups[0].Applied)
	hydration.ReadinessToken = "after-restart"
	share.SetHydrationState(hydration)
	resumed, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, resumed.Status, "%+v", resumed)
	require.Equal(t, partial.OperationID, resumed.OperationID)
	require.Equal(t, partial.Groups[0].AppliedRevision, resumed.Groups[0].AppliedRevision)
	require.Equal(t, 2, f.activation.Calls())
}
