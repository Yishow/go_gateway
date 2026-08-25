package api

import (
	"context"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/gin-gonic/gin"
)

func registerModbusShareRoutes(datalinkGroup *gin.RouterGroup, services *DatalinkServices) {
	if services.ModbusShare == nil {
		return
	}
	if services.Workspace != nil {
		services.Workspace.SetMembershipScopeInvalidator(services.ModbusShare)
	}
	if services.ModbusShareReconciler != nil {
		services.ModbusShare.SetWorkspaceRevisionStore(services.ModbusShareReconciler.WorkspaceRevisionStore())
	}
	if services.SourceRule != nil && services.Workspace != nil {
		services.SourceRule.SetCandidateScope(services.Workspace, services.ModbusShare)
		services.ModbusShare.SetOwnershipChecker(sourcerule.NewShareOwnershipChecker(services.Workspace, services.SourceRule, services.Tag))
		services.ModbusShare.SetDesiredMappingOwnershipChecker(sourcerule.NewShareDesiredMappingOwnershipChecker(services.Workspace, services.SourceRule, services.Tag, services.ModbusShare.Settings, services.Mapping))
	}

	handler := handlers.NewModbusShareHandler(services.ModbusShare).
		WithGate(services.ModbusShare).WithReconciler(services.ModbusShareReconciler)
	if builder := canonicalShareDesiredMappingBuilder(services); builder != nil {
		handler.WithDesiredMappingBuilder(builder)
	}
	datalinkGroup.GET("/modbus-share/status", handler.Status)
	datalinkGroup.POST("/modbus-share/start", handler.Start)
	datalinkGroup.POST("/modbus-share/stop", handler.Stop)
	datalinkGroup.GET("/modbus-share/mappings", handler.ListMappings)
	datalinkGroup.PUT("/modbus-share/mappings/:tagId", handler.UpsertMapping)
	datalinkGroup.DELETE("/modbus-share/mappings/:tagId", handler.DeleteMapping)
	datalinkGroup.POST("/modbus-share/write-tag-value", handler.WriteTagValue)
	datalinkGroup.POST("/modbus-share/sync", handler.SyncFromMappings)
	datalinkGroup.POST("/modbus-share/reconcile", handler.Reconcile)
}

func canonicalShareDesiredMappingBuilder(services *DatalinkServices) handlers.DesiredMappingBuilder {
	if services == nil || services.SourceRule == nil || services.Workspace == nil {
		return nil
	}
	return func(ctx context.Context, workspaceID string, settings modbusshare.Settings) ([]modbusshare.DesiredMapping, error) {
		workspaceRecord, err := services.Workspace.GetOrCreate(ctx)
		if err != nil {
			return nil, err
		}
		desired, err := services.SourceRule.BuildDesiredShareMappingsForDevices(ctx, workspaceID, settings, workspaceRecord.OrderedDeviceIDs)
		if err != nil {
			return nil, err
		}
		for i := range desired {
			desired[i].OwnershipProof = &modbusshare.OwnershipProof{
				Verified:           true,
				WorkspaceID:        workspaceID,
				SourceRuleID:       desired[i].SourceRuleID,
				SourceRuleRevision: desired[i].SourceRuleRevision,
				Basis:              "persisted workspace/source-rule/tag relationship",
			}
		}
		return desired, nil
	}
}

func configureModbusShareActivation(handler *handlers.StudioV2WorkspaceActivationHandler, services *DatalinkServices) {
	if services.ModbusShare == nil {
		return
	}
	if services.ShareRestore != nil {
		handler.WithShareRestore(services.ShareRestore)
	}
	handler.WithActivationBarrierValidator(func(ctx context.Context, req handlers.ActivateWorkspaceRequest) error {
		hydration, err := services.ModbusShare.CheckHydration(ctx)
		if err != nil || hydration.State != modbusshare.HydrationStateReady || !hydration.Readiness {
			return &modbusshare.Error{Code: modbusshare.ErrCodeHydrationRequired, Message: "workspace hydration is required before activation", Retryable: true, Action: "complete workspace bootstrap and retry"}
		}
		settings, err := services.ModbusShare.GetSettings(ctx)
		if err != nil || !settings.Enabled {
			return &modbusshare.Error{Code: modbusshare.ErrCodeDisabled, Message: "modbus share is disabled in global settings", Retryable: false, SettingsRevision: settings.SettingsRevision}
		}
		if req.ReadinessToken == "" || req.ReadinessToken != hydration.ReadinessToken {
			return &modbusshare.Error{Code: modbusshare.ErrCodeSaveIncomplete, Message: "activation readiness token is missing or stale", Retryable: true, Action: "complete autosave and hydrate again"}
		}
		if req.SettingsRevision == "" || req.SettingsRevision != settings.SettingsRevision {
			return &modbusshare.Error{Code: modbusshare.ErrCodeRevisionConflict, Message: "settings revision conflict", Retryable: true, SettingsRevision: settings.SettingsRevision}
		}
		if req.WorkspaceRevision == "" || req.WorkspaceRevision != hydration.WorkspaceRevision {
			return &modbusshare.Error{Code: modbusshare.ErrCodeRevisionConflict, Message: "workspace revision conflict", Retryable: true, WorkspaceRevision: hydration.WorkspaceRevision, SettingsRevision: settings.SettingsRevision}
		}
		return nil
	})
}
