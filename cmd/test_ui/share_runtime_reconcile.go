package main

import (
	"context"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
)

func configureShareRuntimeReconciler(
	rules *sourcerule.Service,
	workspaceSvc *workspace.Service,
	share *modbusshare.Service,
	revisionStore *modbusshare.SQLWorkspaceRevisionStore,
	reconciler *modbusshare.Reconciler,
) {
	reconcile := sourcerule.ShareRuntimeReconcilerFunc(func(ctx context.Context, req sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		scope := req.Scope
		settingsSnapshot := share.Settings()
		if !settingsSnapshot.Enabled {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusAligned, Scope: scope}, nil
		}
		workspaceRecord, err := workspaceSvc.GetOrCreate(ctx)
		if err != nil {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusStale, Scope: scope, Message: "workspace hydration failed"}, err
		}
		workspaceRevision, _, err := revisionStore.GetRevision(ctx, workspaceRecord.ID)
		if err != nil {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusStale, Scope: scope, Message: "workspace revision read failed"}, err
		}
		hydration, err := share.CheckHydration(ctx)
		if err != nil {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusStale, Scope: scope, Message: "share hydration read failed"}, err
		}
		desired, err := rules.BuildDesiredShareMappingsForDevices(ctx, workspaceRecord.ID, settingsSnapshot, workspaceRecord.OrderedDeviceIDs)
		if err != nil {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusStale, Scope: scope, Message: "share candidate build failed"}, err
		}
		out, err := reconciler.Reconcile(ctx, modbusshare.ReconcileRequest{
			WorkspaceID:               workspaceRecord.ID,
			ExpectedWorkspaceRevision: workspaceRevision,
			ExpectedSettingsRevision:  settingsSnapshot.SettingsRevision,
			ReadinessToken:            hydration.ReadinessToken,
			DesiredMappings:           desired,
		})
		if err != nil {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusStale, Scope: scope, Message: "share projection reconcile failed"}, err
		}
		if out.Outcome != shareOutcomeAligned && out.Outcome != shareOutcomeApplied {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusStale, Scope: scope, Message: "share projection is not aligned"}, modbusshare.NewError(modbusshare.ErrCodeReconcileFailed, "share projection is not aligned", true)
		}
		return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusAligned, Scope: scope}, nil
	})
	rules.SetShareRuntimeReconciler(configuredShareRuntimeReconciler{reconcile: reconcile, share: share})
}

type configuredShareRuntimeReconciler struct {
	reconcile sourcerule.ShareRuntimeReconcilerFunc
	share     *modbusshare.Service
}

func (r configuredShareRuntimeReconciler) ReconcileSourceRuleShare(ctx context.Context, req sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
	return r.reconcile(ctx, req)
}

func (r configuredShareRuntimeReconciler) FailClosed(ctx context.Context) {
	if r.share == nil {
		return
	}
	state, err := r.share.CheckHydration(ctx)
	if err != nil {
		state = modbusshare.HydrationState{}
	}
	r.share.FailClosed(state)
}
