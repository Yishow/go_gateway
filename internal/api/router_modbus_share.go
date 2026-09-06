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

// activationBootstrapAction 是 hydration 相關前置檢查一致的建議動作。
const activationBootstrapAction = "complete workspace bootstrap and retry"

func configureModbusShareActivation(handler *handlers.StudioV2WorkspaceActivationHandler, services *DatalinkServices) {
	if services.ModbusShare == nil {
		return
	}
	if services.ShareRestore != nil {
		handler.WithShareRestore(func(ctx context.Context, req handlers.ActivateWorkspaceRequest) error {
			settings, err := services.ModbusShare.GetSettings(ctx)
			if err != nil {
				return err
			}
			if !settings.Enabled {
				return nil
			}
			return services.ShareRestore(ctx, req)
		})
	}
	handler.WithActivationBarrierValidator(func(ctx context.Context, req handlers.ActivateWorkspaceRequest) error {
		return validateStudioV2ActivationBarrier(ctx, services.ModbusShare, req)
	})
	configureWorkspaceRevisionValidator(handler, services)
}

// configureWorkspaceRevisionValidator 安裝與 Share 設定 / hydration 無關的
// 工作區 revision 檢查。Share 的 barrier 在停用且未 hydrate 時會退場，若沒有
// 這一層，啟用就只剩前端自報的 pending_saves 可依賴。
func configureWorkspaceRevisionValidator(
	handler *handlers.StudioV2WorkspaceActivationHandler,
	services *DatalinkServices,
) {
	if services.Workspace == nil || services.ModbusShareReconciler == nil {
		return
	}
	revisionStore := services.ModbusShareReconciler.WorkspaceRevisionStore()
	if revisionStore == nil {
		return
	}
	handler.WithRevisionValidator(newWorkspaceRevisionValidator(
		func(ctx context.Context) (string, error) {
			workspaceRecord, err := services.Workspace.GetOrCreate(ctx)
			if err != nil {
				return "", err
			}
			return workspaceRecord.ID, nil
		},
		revisionStore,
	))
}

// newWorkspaceRevisionValidator 以持久化的工作區 revision 為基準檢查啟用請求，
// 不讀取任何 Share 設定或 hydration 狀態。
func newWorkspaceRevisionValidator(
	workspaceID func(context.Context) (string, error),
	revisionStore modbusshare.WorkspaceRevisionStore,
) func(ctx context.Context, expectedWorkspaceRevision, expectedSettingsRevision string) error {
	return func(ctx context.Context, expectedWorkspaceRevision, _ string) error {
		id, err := workspaceID(ctx)
		if err != nil {
			return &modbusshare.Error{Code: modbusshare.ErrCodeHydrationRequired, Message: "workspace could not be read before activation", Retryable: true, Action: activationBootstrapAction}
		}
		persistedRevision, _, err := revisionStore.GetRevision(ctx, id)
		if err != nil {
			return &modbusshare.Error{Code: modbusshare.ErrCodeHydrationRequired, Message: "workspace revision could not be read before activation", Retryable: true, Action: activationBootstrapAction}
		}
		if persistedRevision == "" {
			// 尚未產生持久化 revision（全新工作區）時沒有可比對的基準，
			// 交由 Share barrier 與前端 barrier 處理，不在此阻擋。
			return nil
		}
		if expectedWorkspaceRevision == "" || expectedWorkspaceRevision != persistedRevision {
			return &modbusshare.Error{Code: modbusshare.ErrCodeRevisionConflict, Message: "workspace revision conflict", Retryable: true, WorkspaceRevision: persistedRevision}
		}
		return nil
	}
}

// activationBarrierShareGate 是啟用前置檢查所需的 Share 讀取面，抽出以便測試
// 涵蓋設定讀取失敗這條分支。
type activationBarrierShareGate interface {
	GetSettings(context.Context) (modbusshare.Settings, error)
	CheckHydration(context.Context) (modbusshare.HydrationState, error)
}

func validateStudioV2ActivationBarrier(
	ctx context.Context,
	gate activationBarrierShareGate,
	req handlers.ActivateWorkspaceRequest,
) error {
	settings, err := gate.GetSettings(ctx)
	if err != nil {
		// 讀取失敗是內部前置檢查失敗，不是「Share 已停用」這個設定狀態；
		// 此時 settings 為零值，不得從中取出 revision 回報。
		return &modbusshare.Error{Code: modbusshare.ErrCodeHydrationRequired, Message: "modbus share settings could not be read before activation", Retryable: true, Action: activationBootstrapAction}
	}
	hydration, hydrationErr := gate.CheckHydration(ctx)
	if hydrationErr != nil {
		// 讀取失敗與「尚未完成 bootstrap」是兩回事：前者是內部錯誤，必須擋下
		// 啟用，不能因為讀不到狀態就把工作區層級的 barrier 一併關掉。
		return &modbusshare.Error{Code: modbusshare.ErrCodeHydrationRequired, Message: "modbus share hydration state could not be read before activation", Retryable: true, Action: activationBootstrapAction}
	}
	hydrationReady := hydration.State == modbusshare.HydrationStateReady && hydration.Readiness
	if !settings.Enabled {
		// Share 停用時只略過 Share 專屬檢查；工作區層級的 barrier 仍須成立。
		// Share 從未完成 bootstrap（hydration 未就緒）時不得綁架工作區啟用。
		if !hydrationReady {
			return nil
		}
		return validateWorkspaceActivationRevisions(req, hydration, settings)
	}
	if !hydrationReady {
		return &modbusshare.Error{Code: modbusshare.ErrCodeHydrationRequired, Message: "workspace hydration is required before activation", Retryable: true, Action: activationBootstrapAction}
	}
	if err := validateWorkspaceActivationReadinessToken(req, hydration); err != nil {
		return err
	}
	if req.SettingsRevision == "" || req.SettingsRevision != settings.SettingsRevision {
		return &modbusshare.Error{Code: modbusshare.ErrCodeRevisionConflict, Message: "settings revision conflict", Retryable: true, SettingsRevision: settings.SettingsRevision}
	}
	return validateWorkspaceActivationWorkspaceRevision(req, hydration, settings)
}

// validateWorkspaceActivationRevisions 檢查與 Share 設定無關的工作區層級 barrier：
// readiness token 與 workspace revision。
func validateWorkspaceActivationRevisions(
	req handlers.ActivateWorkspaceRequest,
	hydration modbusshare.HydrationState,
	settings modbusshare.Settings,
) error {
	if err := validateWorkspaceActivationReadinessToken(req, hydration); err != nil {
		return err
	}
	return validateWorkspaceActivationWorkspaceRevision(req, hydration, settings)
}

func validateWorkspaceActivationReadinessToken(
	req handlers.ActivateWorkspaceRequest,
	hydration modbusshare.HydrationState,
) error {
	if req.ReadinessToken == "" || req.ReadinessToken != hydration.ReadinessToken {
		return &modbusshare.Error{Code: modbusshare.ErrCodeSaveIncomplete, Message: "activation readiness token is missing or stale", Retryable: true, Action: "complete autosave and hydrate again"}
	}
	return nil
}

func validateWorkspaceActivationWorkspaceRevision(
	req handlers.ActivateWorkspaceRequest,
	hydration modbusshare.HydrationState,
	settings modbusshare.Settings,
) error {
	if req.WorkspaceRevision == "" || req.WorkspaceRevision != hydration.WorkspaceRevision {
		return &modbusshare.Error{Code: modbusshare.ErrCodeRevisionConflict, Message: "workspace revision conflict", Retryable: true, WorkspaceRevision: hydration.WorkspaceRevision, SettingsRevision: settings.SettingsRevision}
	}
	return nil
}
