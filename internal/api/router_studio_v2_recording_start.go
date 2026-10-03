package api

import (
	"context"
	"fmt"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

func registerStudioV2RecordingStartRoutes(group *gin.RouterGroup, services *DatalinkServices) {
	var start *workspace.RecordingStartService
	if services.Workspace != nil && services.WriteGroups != nil && services.RecordingPlan != nil {
		start = workspace.NewRecordingStartService(services.Workspace, services.WriteGroups, nil, services.RecordingPlan)
		if services.Device != nil {
			start = workspace.NewRecordingStartService(services.Workspace, services.WriteGroups,
				workspace.NewActivationService(services.Workspace, services.Device, services.Runtime), services.RecordingPlan)
		}
		start.WithActivationBarrier(recordingStartActivationBarrier(services), recordingStartActivationResumeBarrier(services))
	}
	handler := handlers.NewStudioV2RecordingStartHandler(start, services.WriteGroups)
	group.POST("/studio-v2/workspace/recording-start", handler.Start)
	group.GET("/studio-v2/workspace/recording-start/operations/:id", handler.Get)
	group.POST("/studio-v2/workspace/write-groups/basic/:id", handler.EnsureBasic)
}

// A recorded successful barrier may obtain the new process's hydration proof,
// but only for unchanged persisted Share/workspace revisions. The original
// request and intent digest remain unchanged in the ledger.
func recordingStartActivationResumeBarrier(services *DatalinkServices) func(context.Context, workspace.RecordingStartRequest) error {
	initial := recordingStartActivationBarrier(services)
	return func(ctx context.Context, request workspace.RecordingStartRequest) error {
		if services.ModbusShare == nil {
			return initial(ctx, request)
		}
		if request.ReadinessToken == "" {
			return initial(ctx, request)
		}
		settings, err := services.ModbusShare.GetSettings(ctx)
		if err != nil {
			return fmt.Errorf("%w: %w", workspace.ErrRecordingStartShareNotReady, err)
		}
		hydration, err := services.ModbusShare.CheckHydration(ctx)
		if err != nil {
			return fmt.Errorf("%w: %w", workspace.ErrRecordingStartShareNotReady, err)
		}
		if hydration.State != "ready" || !hydration.Readiness || hydration.ReadinessToken == "" ||
			hydration.WorkspaceID != request.WorkspaceID || request.WorkspaceRevision == "" ||
			hydration.WorkspaceRevision != request.WorkspaceRevision ||
			settings.SettingsRevision != request.SettingsRevision {
			return workspace.ErrRecordingStartShareNotReady
		}
		request.ReadinessToken = hydration.ReadinessToken
		return initial(ctx, request)
	}
}

func recordingStartActivationBarrier(services *DatalinkServices) func(context.Context, workspace.RecordingStartRequest) error {
	return func(ctx context.Context, request workspace.RecordingStartRequest) error {
		if len(request.Groups) == 0 {
			if services.ModbusShare == nil {
				return workspace.ErrRecordingStartInvalid
			}
			settings, err := services.ModbusShare.GetSettings(ctx)
			if err != nil {
				return err
			}
			if !settings.Enabled {
				return workspace.ErrRecordingStartInvalid
			}
		}
		req := handlers.ActivateWorkspaceRequest{
			WorkspaceRevision: request.WorkspaceRevision, SettingsRevision: request.SettingsRevision,
			ReadinessToken: request.ReadinessToken,
		}
		if services.ModbusShareReconciler != nil {
			store := services.ModbusShareReconciler.WorkspaceRevisionStore()
			if store != nil {
				validator := newWorkspaceRevisionValidator(func(ctx context.Context) (string, error) {
					record, err := services.Workspace.GetOrCreate(ctx)
					if err != nil {
						return "", err
					}
					return record.ID, nil
				}, store)
				if err := validator(ctx, req.WorkspaceRevision, req.SettingsRevision); err != nil {
					return fmt.Errorf("%w: %w", workspace.ErrRecordingStartShareNotReady, err)
				}
			}
		}
		if services.ModbusShare == nil {
			return nil
		}
		if err := validateStudioV2ActivationBarrier(ctx, services.ModbusShare, req); err != nil {
			return fmt.Errorf("%w: %w", workspace.ErrRecordingStartShareNotReady, err)
		}
		settings, err := services.ModbusShare.GetSettings(ctx)
		if err != nil {
			return err
		}
		if settings.Enabled {
			if services.RecordingStartShareBarrier == nil {
				return workspace.ErrRecordingStartShareNotReady
			}
			return services.RecordingStartShareBarrier(ctx, request)
		}
		return nil
	}
}
