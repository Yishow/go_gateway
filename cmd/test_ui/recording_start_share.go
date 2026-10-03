package main

import (
	"context"
	"fmt"
	"strings"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
)

type recordingStartShareDesiredBuilder interface {
	BuildDesiredShareMappingsForDevices(context.Context, string, modbusshare.Settings, []string) ([]modbusshare.DesiredMapping, error)
	List(context.Context, sourcerule.ListFilter) ([]*schema.SourceRule, error)
}

type recordingStartShareProjectionChecker interface {
	Settings() modbusshare.Settings
	CheckDesiredProjectionForScopeWithSourceRules(context.Context, string, []string, []modbusshare.DesiredMapping, []*schema.SourceRule) error
}

func newRecordingStartShareBarrier(recordSource recordingStartShareDesiredBuilder, share recordingStartShareProjectionChecker) func(context.Context, workspace.RecordingStartRequest) error {
	return func(ctx context.Context, request workspace.RecordingStartRequest) error {
		if recordSource == nil || share == nil {
			return workspace.ErrRecordingStartShareNotReady
		}
		settings := share.Settings()
		if !settings.Enabled {
			return nil
		}
		if strings.TrimSpace(request.WorkspaceID) == "" || len(request.DeviceIDs) == 0 {
			return workspace.ErrRecordingStartShareNotReady
		}
		desired, err := recordSource.BuildDesiredShareMappingsForDevices(ctx, request.WorkspaceID, settings, request.DeviceIDs)
		if err != nil {
			return wrapRecordingStartShareNotReady(err)
		}
		if len(request.Groups) == 0 && len(desired) == 0 {
			return workspace.ErrRecordingStartShareNotReady
		}
		sourceRules, err := recordSource.List(ctx, sourcerule.ListFilter{})
		if err != nil {
			return wrapRecordingStartShareNotReady(err)
		}
		if err := share.CheckDesiredProjectionForScopeWithSourceRules(ctx, request.WorkspaceID, request.DeviceIDs, desired, sourceRules); err != nil {
			return wrapRecordingStartShareNotReady(err)
		}
		return nil
	}
}

func wrapRecordingStartShareNotReady(cause error) error {
	if cause == nil {
		return workspace.ErrRecordingStartShareNotReady
	}
	return fmt.Errorf("%w: %w", workspace.ErrRecordingStartShareNotReady, cause)
}

type workspaceShareRecordReader interface {
	GetOrCreate(context.Context) (*workspace.Record, error)
}

type workspaceShareRestoreSource interface {
	RestoreLocalModbusProjectionForDevices(context.Context, string, string, modbusshare.Settings, []string, sourcerule.ShareProjectionReconciler) (modbusshare.ReconcileOutcome, error)
}

type workspaceShareHydrationReader interface {
	CheckHydration(context.Context) (modbusshare.HydrationState, error)
	Settings() modbusshare.Settings
}

func newWorkspaceShareRestoreBarrier(workspaceService workspaceShareRecordReader, source workspaceShareRestoreSource, share workspaceShareHydrationReader, reconciler sourcerule.ShareProjectionReconciler) handlers.ShareRestoreBarrier {
	return func(ctx context.Context, request handlers.ActivateWorkspaceRequest) error {
		if workspaceService == nil || source == nil || share == nil {
			return modbusshare.NewError(modbusshare.ErrCodeProjectionRequired, "Share restore services are unavailable", true)
		}
		hydration, err := share.CheckHydration(ctx)
		if err != nil {
			return modbusshare.NewError(modbusshare.ErrCodeHydrationRequired, "workspace hydration is required before Share restore", true)
		}
		settingsSnapshot := share.Settings()
		if !settingsSnapshot.Enabled {
			return modbusshare.NewError(modbusshare.ErrCodeDisabled, "Modbus Share is disabled in global settings", false)
		}
		if request.WorkspaceRevision == "" || request.WorkspaceRevision != hydration.WorkspaceRevision {
			return modbusshare.NewError(modbusshare.ErrCodeRevisionConflict, "workspace revision conflict", true)
		}
		workspaceRecord, workspaceErr := workspaceService.GetOrCreate(ctx)
		if workspaceErr != nil || workspaceRecord == nil {
			return modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "workspace membership is unavailable", true)
		}
		// Restore always rebuilds desired state from persisted source-rule
		// candidates. Durable mapping rows are only consistency/rollback
		// evidence and must not become a competing authority.
		out, err := source.RestoreLocalModbusProjectionForDevices(
			ctx, hydration.WorkspaceID, request.WorkspaceRevision, settingsSnapshot,
			workspaceRecord.OrderedDeviceIDs, reconciler,
		)
		if err == nil && out.Outcome != shareOutcomeAligned && out.Outcome != shareOutcomeApplied {
			err = modbusshare.NewError(modbusshare.ErrCodeHydrationRequired, "Share restore did not reach an aligned projection", true)
		}
		return err
	}
}
