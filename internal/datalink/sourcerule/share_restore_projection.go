package sourcerule

import (
	"context"

	"go-gateway/internal/datalink/modbusshare"
)

// ShareProjectionReconciler is the sole runtime projection seam for persisted
// Local Modbus candidates. Implementations must perform full desired-set CAS.
type ShareProjectionReconciler interface {
	Reconcile(context.Context, modbusshare.ReconcileRequest) (modbusshare.ReconcileOutcome, error)
}

// RestoreLocalModbusProjection rebuilds the desired set from persisted
// source-rule candidate snapshots and hands it to the canonical reconciler.
// Browser rows and direct runtime mapping calls are intentionally not inputs.
func (s *Service) RestoreLocalModbusProjection(
	ctx context.Context,
	workspaceID, workspaceRevision string,
	settings modbusshare.Settings,
	reconciler ShareProjectionReconciler,
) (modbusshare.ReconcileOutcome, error) {
	return s.restoreLocalModbusProjection(ctx, workspaceID, workspaceRevision, settings, nil, reconciler)
}

// RestoreLocalModbusProjectionForDevices restores only rules owned by the
// workspace's persisted device membership.
func (s *Service) RestoreLocalModbusProjectionForDevices(ctx context.Context, workspaceID, workspaceRevision string, settings modbusshare.Settings, deviceIDs []string, reconciler ShareProjectionReconciler) (modbusshare.ReconcileOutcome, error) {
	return s.restoreLocalModbusProjection(ctx, workspaceID, workspaceRevision, settings, deviceIDs, reconciler)
}

func (s *Service) restoreLocalModbusProjection(
	ctx context.Context, workspaceID, workspaceRevision string,
	settings modbusshare.Settings, deviceIDs []string,
	reconciler ShareProjectionReconciler,
) (modbusshare.ReconcileOutcome, error) {
	if workspaceID == "" {
		return modbusshare.ReconcileOutcome{Outcome: string(RuntimeReconcileStatusFailed)}, modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "workspace id is required", false)
	}
	if !settings.Enabled {
		return modbusshare.ReconcileOutcome{Outcome: "disabled", SettingsRevision: settings.SettingsRevision}, modbusshare.NewError(modbusshare.ErrCodeDisabled, "Modbus Share is disabled in global settings", false)
	}
	if reconciler == nil {
		return modbusshare.ReconcileOutcome{Outcome: string(RuntimeReconcileStatusFailed)}, modbusshare.NewError(modbusshare.ErrCodeProjectionRequired, "share reconciler is required", false)
	}
	var desired []modbusshare.DesiredMapping
	var err error
	if deviceIDs == nil {
		desired, err = s.BuildDesiredShareMappingsWithSettings(ctx, workspaceID, settings)
	} else {
		desired, err = s.BuildDesiredShareMappingsForDevices(ctx, workspaceID, settings, deviceIDs)
	}
	if err != nil {
		return modbusshare.ReconcileOutcome{Outcome: string(RuntimeReconcileStatusFailed)}, err
	}
	return reconciler.Reconcile(ctx, modbusshare.ReconcileRequest{
		WorkspaceID:               workspaceID,
		ExpectedWorkspaceRevision: workspaceRevision,
		ExpectedSettingsRevision:  settings.SettingsRevision,
		DesiredMappings:           desired,
		Restore:                   true,
	})
}
