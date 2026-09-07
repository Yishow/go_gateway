import * as React from 'react';
import { WorkbenchV2Shell } from '../../../features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { useActivateStudioV2WorkspaceMutation } from '../../../hooks/datalink/useStudioV2WorkspaceActivation';
import { useStudioV2WorkspaceQuery } from '../../../hooks/datalink/useStudioV2Workspace';
import { useStudioV2WorkspaceAuditHistoryQuery } from '../../../hooks/datalink/useStudioV2WorkspaceAuditHistory';
import { useModbusShareStatusQuery } from '../../../hooks/datalink/useModbusShareStatus';
import { useStudioV2AutosaveState } from './useStudioV2AutosaveState';
import {
  activateStudioV2WorkspaceWithShare,
  type StudioV2ShareActivationContext,
} from '../../../features/datalink/workbench-v2/state/studioV2ShareActivation';
import type { ModbusShareReconcileOutcome } from '../../../types/datalink';
import {
  isModbusShareConfiguredEnabled,
  modbusShareConfiguredValue,
  modbusShareStatusFromBootstrap,
} from '../../../types/modbusShare';
import { StudioV2ActivationBarrierError } from '../../../services/studioV2WorkspaceActivation';
import type { WorkbenchV2RuntimeReturnFocus } from '../../../features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { useStudioV2AutosaveSettlement } from './studioV2AutosaveBarrier';
import '../../../features/datalink/workbench-v2/styles/workbench-v2.css';

interface DatalinkWorkbenchV2PageProps {
  navigateTo?: (target: string) => void;
  runtimeReturnFocus?: WorkbenchV2RuntimeReturnFocus | null;
}

interface WorkbenchV2BootstrapStateProps {
  status: 'loading' | 'error';
}

function WorkbenchV2BootstrapState({ status }: WorkbenchV2BootstrapStateProps) {
  const isError = status === 'error';

  return (
    <div
      data-workbench-v2="true"
      data-testid={isError ? 'workbench-v2-bootstrap-error' : 'workbench-v2-bootstrap-loading'}
      role={isError ? 'alert' : 'status'}
      className="wbv2-bootstrap-state bg-canvas min-h-screen font-sans text-slate-100"
    >
      <section
        data-testid="workbench-v2-bootstrap-shell"
        className="wbv2-bootstrap-shell"
        aria-label={isError ? 'Studio V2 工作區載入失敗' : 'Studio V2 工作區載入中'}
      >
        <div className="flex items-start justify-between gap-6">
          <div>
            <p className="font-mono text-[11px] uppercase tracking-[0.16em] text-cyan-200/70">
              Studio V2
            </p>
            <h1 className="mt-2 text-2xl font-semibold leading-tight text-slate-50">
              {isError ? '無法載入工作區' : '正在整理工作區狀態'}
            </h1>
            <p className="mt-2 max-w-[34rem] text-sm leading-6 text-slate-400">
              {isError
                ? '無法載入 Studio V2 工作區。請檢查後端資料庫連線後重試。'
                : '正在讀取工作區、設備、接入規則與資料庫目標，完成後會直接進入設定流程。'}
            </p>
          </div>
          <div className={isError ? 'wbv2-bootstrap-mark is-error' : 'wbv2-bootstrap-mark'}>
            {isError ? 'ERR' : 'SYNC'}
          </div>
        </div>
        <div className="mt-8 grid gap-3 sm:grid-cols-3" aria-hidden="true">
          <span className="wbv2-skeleton-line" />
          <span className="wbv2-skeleton-line" />
          <span className="wbv2-skeleton-line" />
        </div>
      </section>
    </div>
  );
}

/**
 * Datalink Workbench V2 頁面入口
 * 
 * 路由策略：/studio/v2 是唯一設定入口；/studio 依既有 generic unknown-route policy 收斂。
 * 限制 scope 與狀態綁定。
 */
export default function DatalinkWorkbenchV2Page({
  navigateTo,
  runtimeReturnFocus,
}: DatalinkWorkbenchV2PageProps) {
  const workspaceQuery = useStudioV2WorkspaceQuery();
  const shareStatusQuery = useModbusShareStatusQuery(workspaceQuery.isSuccess);
  const bootstrapShareStatus = modbusShareStatusFromBootstrap(workspaceQuery.data?.modbus_share);
  const shareStatus = React.useMemo(() => {
    const liveStatus = shareStatusQuery.data;
    if (!bootstrapShareStatus || !liveStatus) {
      return liveStatus ?? bootstrapShareStatus;
    }

    const bootstrapConfigured = modbusShareConfiguredValue(bootstrapShareStatus);
    const liveConfigured = modbusShareConfiguredValue(liveStatus);
    const configuredEnabled = bootstrapConfigured === false
      ? false
      : liveConfigured ?? bootstrapConfigured;

    return {
      ...liveStatus,
      ...bootstrapShareStatus,
      workspace_revision: bootstrapShareStatus.workspace_revision ?? liveStatus.workspace_revision,
      settings_revision: bootstrapShareStatus.settings_revision ?? liveStatus.settings_revision,
      readiness_token: bootstrapShareStatus.readiness_token ?? liveStatus.readiness_token,
      configured_enabled: configuredEnabled ?? liveStatus.configured_enabled ?? bootstrapShareStatus.configured_enabled,
    };
  }, [bootstrapShareStatus, shareStatusQuery.data]);
  const auditHistoryQuery = useStudioV2WorkspaceAuditHistoryQuery(workspaceQuery.isSuccess);
  const autosave = useStudioV2AutosaveState(workspaceQuery.isSuccess);
  const waitForAutosaveSettlement = useStudioV2AutosaveSettlement(autosave.state);
  const activationMutation = useActivateStudioV2WorkspaceMutation();
  const activateWorkspace = async () => {
    const settlement = await waitForAutosaveSettlement();
    const barrier = settlement.barrier;
    if (!autosave.workspaceHydrated || barrier.pending_saves > 0 || barrier.save_error) {
      throw new StudioV2ActivationBarrierError('modbus_share_save_incomplete');
    }

    const [freshWorkspaceResult, freshShareResult] = await Promise.all([
      workspaceQuery.refetch(),
      shareStatusQuery.refetch(),
    ]);
    const freshWorkspace = freshWorkspaceResult.data ?? workspaceQuery.data;
    const freshBootstrapShare = modbusShareStatusFromBootstrap(freshWorkspace?.modbus_share);
    const freshLiveStatus = freshShareResult.data ?? shareStatusQuery.data;
    const currentShareStatus = (() => {
      if (!freshBootstrapShare || !freshLiveStatus) {
        return freshBootstrapShare ?? freshLiveStatus ?? shareStatus;
      }
      const bootstrapConfigured = modbusShareConfiguredValue(freshBootstrapShare);
      const liveConfigured = modbusShareConfiguredValue(freshLiveStatus);
      const configuredEnabled = bootstrapConfigured === false
        ? false
        : liveConfigured ?? bootstrapConfigured;
      return {
        ...freshLiveStatus,
        ...freshBootstrapShare,
        workspace_revision: freshBootstrapShare.workspace_revision ?? freshLiveStatus.workspace_revision,
        settings_revision: freshBootstrapShare.settings_revision ?? freshLiveStatus.settings_revision,
        readiness_token: freshBootstrapShare.readiness_token ?? freshLiveStatus.readiness_token,
        configured_enabled: configuredEnabled ?? freshLiveStatus.configured_enabled ?? freshBootstrapShare.configured_enabled,
      };
    })();

    if (!currentShareStatus) {
      throw new StudioV2ActivationBarrierError('modbus_share_hydration_required');
    }
    const shareConfigured = isModbusShareConfiguredEnabled(currentShareStatus);
    const shareIsStale = shareConfigured
      && (currentShareStatus.hydration_state !== 'ready' || currentShareStatus.readiness !== true);
    if (shareIsStale) {
      throw new StudioV2ActivationBarrierError(
        'modbus_share_revision_conflict',
        currentShareStatus.error?.retryable ?? true,
        currentShareStatus.error?.action,
        currentShareStatus.error?.request_id,
      );
    }

    const shareContext: StudioV2ShareActivationContext = {
      workspace_id: freshWorkspace?.id ?? workspaceQuery.data?.id ?? '',
      workspace_revision: currentShareStatus.workspace_revision ?? '',
      settings_revision: currentShareStatus.settings_revision ?? '',
      readiness_token: currentShareStatus.readiness_token ?? '',
      configured_enabled: shareConfigured,
      canonical_plan: currentShareStatus.canonical_plan,
    };

    return activateStudioV2WorkspaceWithShare(
      settlement.state,
      (projection?: ModbusShareReconcileOutcome) => activationMutation.mutateAsync({
        workspace_revision: projection?.new_workspace_revision || shareContext.workspace_revision,
        settings_revision: projection?.settings_revision || shareContext.settings_revision,
        readiness_token: projection?.new_readiness_token || shareContext.readiness_token,
        pending_saves: barrier.pending_saves,
        ...(barrier.save_error ? { save_error: barrier.save_error } : {}),
      }),
      shareContext,
    );
  };

  if (
    workspaceQuery.isError ||
    (workspaceQuery.isSuccess && !workspaceQuery.data) ||
    autosave.devicesQuery.isError ||
    autosave.rulesQuery.isError ||
    autosave.mappingsQuery.isError ||
    autosave.databaseConfigQuery.isError ||
    autosave.databaseTargetsQuery.isError
  ) {
    return <WorkbenchV2BootstrapState status="error" />;
  }

  if (
    workspaceQuery.isLoading ||
    (workspaceQuery.isSuccess && !autosave.workspaceHydrated) ||
    autosave.devicesQuery.isLoading ||
    autosave.rulesQuery.isLoading ||
    autosave.mappingsQuery.isLoading ||
    autosave.databaseConfigQuery.isLoading ||
    autosave.databaseTargetsQuery.isLoading
  ) {
    return <WorkbenchV2BootstrapState status="loading" />;
  }

  const workspace = workspaceQuery.data!;

  return (
    <div
      data-workbench-v2="true"
      data-testid="workbench-v2-root"
      data-workspace-id={workspace.id}
      data-workspace-status={workspace.status}
      className="bg-canvas min-h-screen font-sans text-slate-100 antialiased"
    >
      {autosave.draftLossWarning && (
        <div
          data-testid="workbench-v2-draft-loss-warning"
          role="alert"
          className="mx-4 mt-4 rounded-xl border border-amber-500/30 bg-amber-500/10 px-4 py-3 text-sm text-amber-200 sm:mx-5 lg:mx-7"
        >
          {autosave.draftLossWarning}
        </div>
      )}
      <WorkbenchV2Shell
        state={autosave.state}
        workspaceId={workspace.id}
        actions={autosave.actions}
        navigateTo={navigateTo}
        activateWorkspace={activateWorkspace}
        workspaceReadiness={workspace.readiness_summary}
        workspaceAuditHistory={auditHistoryQuery.data?.entries ?? []}
        workspaceAuditUnavailable={auditHistoryQuery.isError}
        shareStatus={shareStatus}
        runtimeReturnFocus={runtimeReturnFocus}
      />
    </div>
  );
}
