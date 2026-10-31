package main

import (
	"context"
	"github.com/google/uuid"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/sourcerule"
	"log"
)

// hydrateStartupShare is the existing startup-only recovery boundary. It does
// not start recording or the Share listener; callers retain those lifecycle steps.
func hydrateStartupShare(ctx context.Context, services gatewayServices, revisionStore *modbusshare.SQLWorkspaceRevisionStore, reconciler *modbusshare.Reconciler) {
	shareSettingsLoaded, workspaceSvc := services.shareLoaded, services.workspace
	modbusShareSvc, sourceRuleSvc, tagSvc, mappingSvc := services.modbusShare, services.sourceRule, services.tag, services.mapping
	workspaceRecord, workspaceErr := workspaceSvc.GetOrCreate(ctx)
	//nolint:gocritic // Startup branches preserve the required hydration failure precedence.
	if !shareSettingsLoaded {
		modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateFailed, Readiness: false})
	} else if workspaceErr != nil {
		log.Printf("讀取 Studio V2 workspace 失敗，Share hydration 維持 failed: %v", workspaceErr)
		modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateFailed, Readiness: false})
	} else {
		reconciler.WithOwnershipValidator(sourcerule.NewShareDesiredMappingOwnershipChecker(workspaceSvc, sourceRuleSvc, tagSvc, modbusShareSvc.Settings, mappingSvc))

		settingsSnapshot := modbusShareSvc.Settings()
		workspaceRevision, _, revisionErr := revisionStore.GetRevision(ctx, workspaceRecord.ID)
		//nolint:gocritic // Revision, empty-state, and enabled-state checks intentionally remain ordered.
		if revisionErr != nil {
			log.Printf("讀取 Share workspace revision 失敗，hydration 維持 failed: %v", revisionErr)
			modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateFailed, WorkspaceID: workspaceRecord.ID, SettingsRevision: settingsSnapshot.SettingsRevision, Readiness: false})
		} else if workspaceRevision == "" {
			log.Printf("Share workspace revision 為空，hydration 維持 failed")
			modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateFailed, WorkspaceID: workspaceRecord.ID, SettingsRevision: settingsSnapshot.SettingsRevision, Readiness: false})
		} else if settingsSnapshot.Enabled {
			// Source-rule candidates are the sole restart authority. The durable
			// mapping snapshot is evidence for reconcile rollback/consistency,
			// never a competing desired-state source.
			if out, err := sourceRuleSvc.RestoreLocalModbusProjectionForDevices(ctx, workspaceRecord.ID, workspaceRevision, settingsSnapshot, workspaceRecord.OrderedDeviceIDs, reconciler); err != nil || (out.Outcome != shareOutcomeAligned && out.Outcome != shareOutcomeApplied) {
				if err != nil {
					log.Printf("本機 Modbus 分享服務 source-rule projection 還原失敗，hydration 維持 failed: %v", err)
				}
				modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateFailed, WorkspaceID: workspaceRecord.ID, WorkspaceRevision: workspaceRevision, SettingsRevision: settingsSnapshot.SettingsRevision, Readiness: false})
			} else {
				modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateReady, WorkspaceID: workspaceRecord.ID, WorkspaceRevision: workspaceRevision, SettingsRevision: settingsSnapshot.SettingsRevision, Readiness: true, ReadinessToken: uuid.NewString()})
			}
		} else if shareSettingsLoaded {
			modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateReady, WorkspaceID: workspaceRecord.ID, WorkspaceRevision: workspaceRevision, SettingsRevision: settingsSnapshot.SettingsRevision, Readiness: true, ReadinessToken: uuid.NewString()})
		}
	}
}
