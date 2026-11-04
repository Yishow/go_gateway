import { useStudioLocatorRecovery } from './useStudioLocatorRecovery';
import { RecordingHomeSourceHost, type HomeSourceControls, type HomeRuleObserver } from '../../../features/datalink/workbench-v2/recording-home/RecordingHomeSourceHost';
import { RecordingHomeSourceContent, RecordingHomeSourceLoading } from '../../../features/datalink/workbench-v2/recording-home/RecordingHomeSourceContent';
import { RecordingSourceNavigation } from '../../../features/datalink/workbench-v2/recording-editor/RecordingSourceNavigation';
import { useTranslation } from 'react-i18next';
import { RecordingDraftReadState, useRecordingNavigation } from '../../../features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { RecordingHomeRoute } from '../../../features/datalink/workbench-v2/recording-home/RecordingHomeRoute';
import { parseStudioPresentation, type StudioPresentation } from '../../../services/studioRecordingView';
import type { StudioEntityLocatorIntent } from '../../../services/studioEntityLocator';
import { useStudioLocatorEntry } from '../../../hooks/datalink/useStudioLocatorEntry';
import { StudioLocatorEntryNotice, type StudioLocatorNoticeReason } from '../../../features/datalink/workbench-v2/shell/StudioLocatorEntryNotice';
import * as React from 'react';
import { RuntimeLogsLink } from '../../../features/runtime-logs/RuntimeLogsLink';
import { useQueryClient } from '@tanstack/react-query';
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
import { studioV2RuntimeContextAPI } from '../../../services/studioV2RuntimeContext';
import { studioV2WorkspaceKeys } from '../../../hooks/datalink/keys';
import { boundedString, normalizeTypedEnvelope } from '../../../utils/safeJson';
import type { StudioV2ActivationRecovery } from '../../../types/studioV2Activation';
import '../../../features/datalink/workbench-v2/styles/workbench-v2.css';

interface DatalinkWorkbenchV2PageProps {
  presentation?: StudioPresentation;
  presentationLocation?: { search: string; hash: string };
  navigateTo?: (target: string, options?: { replace?: boolean }) => void;
  locatorIntent?: StudioEntityLocatorIntent;
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
      className="wbv2-bootstrap-state bg-canvas min-h-screen font-sans text-[var(--gw-text)]"
    >
      <section
        data-testid="workbench-v2-bootstrap-shell"
        className="wbv2-bootstrap-shell"
        aria-label={isError ? 'Studio V2 工作區載入失敗' : 'Studio V2 工作區載入中'}
      >
        <div className="flex items-start justify-between gap-6">
          <div>
            <p className="font-mono text-xs text-[var(--gw-accent)]">
              Studio V2
            </p>
            <h1 className="mt-2 text-2xl font-semibold leading-tight text-[var(--gw-text)]">
              {isError ? '無法載入工作區' : '正在整理工作區狀態'}
            </h1>
            <p className="mt-2 max-w-[34rem] text-sm leading-6 text-[var(--gw-muted)]">
              {isError
                ? '無法載入 Studio V2 工作區。請檢查後端資料庫連線後重試。'
                : '正在讀取工作區、設備、接入規則與資料庫目標，完成後會直接進入設定流程。'}
            </p>
          </div>
          <div className={isError ? 'wbv2-bootstrap-mark is-error' : 'wbv2-bootstrap-mark'}>
            {isError ? 'ERR' : 'SYNC'}
          </div>
        </div>
        <div className="mt-5"><RuntimeLogsLink /></div>
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
function AdvancedWorkbenchV2Page({
  navigateTo,
  runtimeReturnFocus, locatorIntent, presentationLocation, recordingIntent, homeSource,
}: DatalinkWorkbenchV2PageProps & { recordingIntent?: 'new' | 'located' | 'copy'; homeSource?: HomeSourceControls }) {
  const { t: draftText } = useTranslation('recording-editor');
  const { blocked: hasUnsavedRecording } = useRecordingNavigation();
  const locatorMode = locatorIntent !== undefined && locatorIntent.kind !== 'absent';
  const locatorQuery = useStudioLocatorEntry(locatorIntent?.kind === 'locator' ? locatorIntent.locator : null);
  const freshLocatorEntry = locatorMode && !locatorQuery.isFetching && !locatorQuery.isError && locatorQuery.data?.ok ? locatorQuery.data.value : undefined;
  const recovery = useStudioLocatorRecovery(locatorIntent?.kind === 'locator' ? locatorIntent.locator : null, freshLocatorEntry, navigateTo);
  const locatorEntry = recovery.entry;
  const [locatorEditorMounted, setLocatorEditorMounted] = React.useState(false);
  const editorMountedRef = React.useRef(false), locatorPending = React.useRef(false);
  const { refetch: revalidateLocator } = locatorQuery;
  const retainLocatorEditor = React.useCallback((mounted: boolean) => {
    if (editorMountedRef.current === mounted) return;
    editorMountedRef.current = mounted; setLocatorEditorMounted(mounted);
    // Close resolves only the retention denial; a pending navigation still requires a fresh read.
    if (!mounted && locatorPending.current) void revalidateLocator();
  }, [revalidateLocator]);
  const workspaceQuery = useStudioV2WorkspaceQuery(!locatorMode);
  const shareStatusQuery = useModbusShareStatusQuery(homeSource?.target.kind === 'points' || !locatorMode && workspaceQuery.isSuccess);
  const bootstrapShare = workspaceQuery.data?.modbus_share;
  // 必須維持穩定引用：設定頁的 Share 區塊以 shareStatus 變更作為同步訊號，每次 render 產生新物件會把 Start/Stop 結果洗回舊快照。
  const bootstrapShareStatus = React.useMemo(() => modbusShareStatusFromBootstrap(bootstrapShare), [bootstrapShare]);
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
  const persistedMCWorkspaceScope = React.useMemo(() => {
    const workspaceId = locatorEntry?.workspace.id ?? workspaceQuery.data?.id;
    const workspaceRevision = shareStatus?.workspace_revision;
    if (!workspaceId || shareStatus?.workspace_id !== workspaceId || !workspaceRevision) return undefined;
    return { workspace_id: workspaceId, workspace_revision: workspaceRevision };
  }, [shareStatus?.workspace_id, shareStatus?.workspace_revision, workspaceQuery.data?.id, locatorEntry?.workspace.id]);
  const auditHistoryQuery = useStudioV2WorkspaceAuditHistoryQuery(!locatorMode && workspaceQuery.isSuccess);
  const homeRuleObserverRef = homeSource?.ruleObserverRef;
  const homeRuleObserver = React.useMemo<HomeRuleObserver | undefined>(() => homeRuleObserverRef ? {
    submitted: (rule, generation) => homeRuleObserverRef.current?.submitted?.(rule, generation),
    saved: (rule, saved) => homeRuleObserverRef.current?.saved(rule, saved),
    failed: (rule, error) => homeRuleObserverRef.current?.failed(rule, error),
  } : undefined, [homeRuleObserverRef]);
  const autosave = useStudioV2AutosaveState(!locatorMode && workspaceQuery.isSuccess, locatorEntry, locatorEditorMounted, persistedMCWorkspaceScope, locatorEntry?.workspace.id ?? workspaceQuery.data?.id, homeRuleObserver);
  const waitForAutosaveSettlement = useStudioV2AutosaveSettlement(autosave.state);
  const activationMutation = useActivateStudioV2WorkspaceMutation();
  const activationInFlight = React.useRef(false);
  const activationOperationId = React.useRef<string | undefined>(undefined);
  const queryClient = useQueryClient();
  const { refetch: refetchWorkspace } = workspaceQuery;
  const workspaceId = workspaceQuery.data?.id;
  const recoverActivationStatus = React.useCallback(async (): Promise<StudioV2ActivationRecovery> => {
    try {
      const [freshWorkspace, runtime] = await Promise.all([
        refetchWorkspace(),
        queryClient.fetchQuery({ queryKey: studioV2WorkspaceKeys.runtimeContext(),
          queryFn: studioV2RuntimeContextAPI.get, staleTime: 0, retry: false }),
      ]);
      if (freshWorkspace.isError) throw freshWorkspace.error;
      const confirmedWorkspaceId = boundedString(freshWorkspace.data?.id);
      if (!confirmedWorkspaceId || confirmedWorkspaceId !== workspaceId || !runtime ||
        runtime.workspace_id !== confirmedWorkspaceId || !Array.isArray(runtime.devices)) {
        throw new Error('activation recovery context invalid');
      }
      const devices: StudioV2ActivationRecovery['devices'] = [];
      for (const device of runtime.devices) {
        const deviceId = boundedString(device?.device_id);
        if (!deviceId || typeof device?.running !== 'boolean') throw new Error('activation recovery device invalid');
        devices.push({ device_id: deviceId, running: device.running });
      }
      const operationId = normalizeTypedEnvelope(runtime).operationId ?? activationOperationId.current;
      return { workspace_id: confirmedWorkspaceId, devices, ...(operationId ? { operation_id: operationId } : {}) };
    } catch (error) {
      throw new StudioV2ActivationBarrierError('activation_failed', false, undefined,
        normalizeTypedEnvelope(error).requestId, 'unconfirmed', activationOperationId.current);
    }
  }, [queryClient, refetchWorkspace, workspaceId]);

  const performActivation = async () => {
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

  const activateWorkspace = async () => {
    if (activationInFlight.current) {
      throw new StudioV2ActivationBarrierError('activation_failed', false, undefined, undefined, 'unconfirmed');
    }
    activationInFlight.current = true;
    activationOperationId.current = undefined;
    try {
      const response = await performActivation();
      activationOperationId.current = normalizeTypedEnvelope(response).operationId;
      return response;
    } catch (error) {
      activationOperationId.current = normalizeTypedEnvelope(error).operationId;
      throw error;
    } finally {
      activationInFlight.current = false;
    }
  };

  const { restoreRejected } = recovery;
  React.useLayoutEffect(() => {
    if (!homeSource && !recordingIntent) restoreRejected(autosave.locatorHydration.blocked, autosave.locatorHydration.lastEntry, autosave.locatorHydration.retainedRuleRevision);
  }, [homeSource, recordingIntent, restoreRejected, autosave.locatorHydration.blocked, autosave.locatorHydration.lastEntry, autosave.locatorHydration.retainedRuleRevision]);
  locatorPending.current = locatorMode && !autosave.locatorHydration.ready;
  const locatorReason: StudioLocatorNoticeReason = recovery.refused ? 'stale' : locatorIntent?.kind === 'invalid' ? 'invalid'
    : autosave.locatorHydration.blocked ? 'draftBlocked'
    : locatorEntry?.entry.write_groups?.groups.some(group => group.id === locatorEntry.locator.group_id && group.status === 'deleted') ? 'removed'
    : locatorEntry?.selection.state === 'unavailable' ? 'unavailable'
    : locatorQuery.data && !locatorQuery.data.ok ? locatorQuery.data.reason === 'foreign-scope' ? 'foreign'
      : locatorQuery.data.reason === 'stale-reference' ? 'stale' : 'unavailable'
    : locatorQuery.isError ? 'unavailable' : 'loading';
  const locatorNotice = <StudioLocatorEntryNotice reason={locatorReason} entry={locatorEntry} navigateTo={navigateTo} />;
  if (locatorMode && !autosave.locatorHydration.ready && !autosave.workspaceHydrated && !autosave.locatorHydration.lastEntry) {
    return homeSource && homeSource.target.kind !== 'create' && locatorReason === 'loading' ? <RecordingHomeSourceLoading onClose={homeSource.onClose} /> : locatorNotice;
  }

  const backgroundReadFailed = !locatorMode && (
    workspaceQuery.isError ||
    (workspaceQuery.isSuccess && !workspaceQuery.data) ||
    autosave.devicesQuery.isError ||
    autosave.rulesQuery.isError ||
    autosave.mappingsQuery.isError ||
    autosave.databaseConfigQuery.isError ||
    autosave.databaseTargetsQuery.isError
  );
  // Retain only a previously owned missing-mapping retry; the failed GET grants no write authority.
  const retainMappingRetry = autosave.mappingsQuery.isError && autosave.autosaveBarrier.pending_saves > 0 &&
    Object.values(autosave.state.mappings).some(mapping => Boolean(mapping.workspace_id) && mapping.workspace_id === workspaceQuery.data?.id &&
      mapping.save_state === 'save-error' && mapping.save_error_detail?.code === 'workspace_mapping_not_found');
  // A failed refresh must not destroy an already-hydrated editor's unsent intent.
  // This retains the existing form, not a second persisted draft or fresh readiness proof.
  if (backgroundReadFailed && !(workspaceQuery.data && autosave.workspaceHydrated && (hasUnsavedRecording() || retainMappingRetry))) {
    return <WorkbenchV2BootstrapState status="error" />;
  }

  if (!locatorMode && (
    workspaceQuery.isLoading ||
    (workspaceQuery.isSuccess && !autosave.workspaceHydrated) ||
    autosave.devicesQuery.isLoading ||
    autosave.rulesQuery.isLoading ||
    autosave.mappingsQuery.isLoading ||
    autosave.databaseConfigQuery.isLoading ||
    autosave.databaseTargetsQuery.isLoading
  )) {
    return <WorkbenchV2BootstrapState status="loading" />;
  }

  const shownEntry = autosave.locatorHydration.ready ? locatorEntry : autosave.locatorHydration.lastEntry;
  const workspace = locatorMode ? shownEntry?.workspace ?? workspaceQuery.data! : workspaceQuery.data!;

  if (homeSource) return <RecordingSourceNavigation workspaceId={workspace.id} source={autosave.deviceSource}><RecordingHomeSourceContent {...homeSource}
    shareEnabled={shareStatusQuery.isSuccess && !shareStatusQuery.isFetching && !shareStatusQuery.isError && shareStatusQuery.data?.workspace_id === homeSource.target.workspaceId ? modbusShareConfiguredValue(shareStatusQuery.data) : undefined}
    onRetryShare={() => { void shareStatusQuery.refetch(); }} onOpenShareSettings={() => navigateTo?.('/studio/v2?view=advanced&section=settings')}
    state={autosave.state} dispatch={autosave.actions.dispatch} source={autosave.deviceSource} mcEvidenceByRule={autosave.persistedMCPoints} suspended={homeSource.target.kind === 'create' ? !autosave.workspaceHydrated || workspace.id !== homeSource.target.workspaceId : !autosave.locatorHydration.ready} /></RecordingSourceNavigation>;

  return (
    <div
      data-workbench-v2="true"
      data-testid="workbench-v2-root"
      data-workspace-id={workspace.id}
      data-workspace-status={workspace.status}
      className="bg-canvas min-h-screen font-sans text-[var(--gw-text)] antialiased"
    >
      {backgroundReadFailed && <div role="alert" className="mx-4 my-3" data-testid="recording-draft-refresh-failed">
        <p>{draftText('navigation.refreshFailed')}</p>
        <button type="button" className="rounded-md border border-[var(--gw-control)] bg-[var(--gw-surface)] px-3 py-2 text-sm text-[var(--gw-text)] hover:bg-[var(--gw-active)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--gw-focus)]" onClick={() => { void Promise.all([workspaceQuery.refetch(), autosave.devicesQuery.refetch(), autosave.rulesQuery.refetch(), autosave.mappingsQuery.refetch(), autosave.databaseConfigQuery.refetch(), autosave.databaseTargetsQuery.refetch()]); }}>{draftText('navigation.refresh')}</button>
      </div>}
      {locatorMode && !autosave.locatorHydration.ready && !recovery.unavailable && locatorNotice}
      <div hidden={locatorMode && !autosave.locatorHydration.ready && !locatorEditorMounted && !autosave.locatorHydration.blocked && !recovery.unavailable}>
      {autosave.draftLossWarning && (
        <div
          data-testid="workbench-v2-draft-loss-warning"
          role="alert"
          className="mx-4 mt-4 rounded-xl border border-[var(--gw-warning-border)] bg-[var(--gw-warning-bg)] px-4 py-3 text-sm text-[var(--gw-warning-fg)] sm:mx-5 lg:mx-7"
        >
          {autosave.draftLossWarning}
        </div>
      )}
      <RecordingDraftReadState unavailable={backgroundReadFailed}><RecordingSourceNavigation workspaceId={workspace.id} source={autosave.deviceSource}><WorkbenchV2Shell
        state={autosave.state}
        recordingIntent={recordingIntent} presentationLocation={presentationLocation}
        workspaceId={workspace.id}
        actions={autosave.actions} deviceSource={autosave.deviceSource}
        navigateTo={navigateTo}
        activateWorkspace={activateWorkspace}
        recoverActivationStatus={recoverActivationStatus}
        workspaceReadiness={workspace.readiness_summary}
        workspaceAuditHistory={auditHistoryQuery.data?.entries ?? []}
        workspaceAuditUnavailable={auditHistoryQuery.isError}
        shareStatus={shareStatus}
        runtimeReturnFocus={locatorMode ? undefined : runtimeReturnFocus}
        locatorEntry={locatorMode ? shownEntry : undefined}
        locatorSuspended={locatorMode && !autosave.locatorHydration.ready}
        locatorDraftBlocked={locatorMode && autosave.locatorHydration.blocked}
        locatorOnly={locatorMode && !autosave.workspaceHydrated}
        locatorRecoveryUnavailable={recovery.unavailable} onRetryLocator={() => { void locatorQuery.refetch(); }}
        locatorCanLeave={autosave.locatorHydration.canLeave && autosave.autosaveBarrier.pending_saves === 0 && !autosave.autosaveBarrier.save_error}
        onLocatorEditorMountedChange={retainLocatorEditor}
      /></RecordingSourceNavigation></RecordingDraftReadState>
      </div>
    </div>
  );
}

/** Precise locator/runtime intent always outranks an optional presentation hint. */
export default function DatalinkWorkbenchV2Page(props: DatalinkWorkbenchV2PageProps) {
  const locator = props.locatorIntent?.kind;
  if (props.presentation === 'recordings' && (!locator || locator === 'absent') && !props.runtimeReturnFocus) {
    return <RecordingHomeSourceHost renderSource={controls => <AdvancedWorkbenchV2Page key={controls.target.epoch} homeSource={controls} navigateTo={props.navigateTo}
      locatorIntent={controls.target.kind === 'create' ? { kind: 'absent' } : { kind: 'locator', locator: { workspace_id: controls.target.workspaceId, step: controls.target.ruleId ? 2 : 1, device_id: controls.target.deviceId,
        rule_id: controls.target.ruleId ?? null, mapping_id: null, group_id: null, group_revision: null } }} />} >
      {(open, workspaceChanged, create) => <RecordingHomeRoute location={props.presentationLocation ?? { search: '?view=recordings', hash: '' }} onNavigate={props.navigateTo}
        onOpenDevice={open} onCreateDevice={create} onWorkspaceChange={workspaceChanged} />}
    </RecordingHomeSourceHost>;
  }
  const intent = parseStudioPresentation(props.presentationLocation?.search ?? '');
  const simple = props.presentation === 'recording' && !props.runtimeReturnFocus && intent.presentation === 'recording';
  return <AdvancedWorkbenchV2Page {...props} recordingIntent={simple ? intent.recordingIntent === 'copy' ? 'copy' : locator === 'locator' ? 'located' : 'new' : undefined} />;
}
