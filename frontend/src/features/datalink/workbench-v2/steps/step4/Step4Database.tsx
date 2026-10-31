import type { DeviceSourcePort } from '../../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { markReturnRegion } from '@/utils/returnFocus';
import type { StudioLocatorEntryRead } from '../../../../../types/studioV2LocatorEntry';
import { useRuntimeRuleEvidence } from '../../../../../hooks/datalink/useRuntimeRuleEvidence';
import { RecordingScopeSummary } from './RecordingScopeSummary';
import type { RecordingViewedDevice } from './recordingScopeModel';
import { RecordingTestWriteRecoveryPanel } from './RecordingTestWriteRecoveryPanel';
import { useRecordingChangeController } from './writeGroup/lifecycle/useRecordingChangeController';
import { RecordingChangeDialog } from './writeGroup/lifecycle/RecordingChangeDialog';
import type { SavedRecordingActionHost } from './writeGroup/lifecycle/SavedRecordingActions';
import * as React from 'react';
import { useMemo, useCallback, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useStep4Activation } from './useStep4Activation';
import {
  createPoolConnectorPatch,
  databaseConnectorNeedsPassword,
  isDatabaseConnectorIdentityChange,
  isRowGroupScopeChange,
  useStep4Readonly,
} from './step4DatabaseHelpers';
import { WriteGroupSection } from './writeGroup/WriteGroupSection';
import { destinationIsSaved } from './basicRecordingPanelHelpers';
import { isModbusShareConfiguredEnabled } from '../../../../../types/modbusShare';
import { BasicRecordingPanel } from './BasicRecordingPanel';
import { useWriteGroupsQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import { SafeQueryBoundary } from '@/utils/SafeQueryBoundary';
import { CommitSummary } from './CommitSummary';
import { Step4SupportPanels } from './Step4SupportPanels';
import { ShareOutputSummary } from './ShareOutputSummary';
import { ActivationRecoveryNotice } from './ActivationNeutralSummary';
import { WorkspaceScopeNote } from './WorkspaceScopeNote';
import { getDefaultConnector, getDbKindPatch } from '../../state/dbSchemas';
import type { WorkbenchV2State, DbConnector, SettingsConnector } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import type {
  StudioV2ActivationRecovery,
  StudioV2ActivationResponse,
} from '../../../../../types/studioV2Activation';
import type { StudioV2WorkspaceReadinessSummary } from '../../../../../types/studioV2WorkspaceReadiness';
import type { WorkspaceReadinessStepNumber } from '../../components/WorkspaceReadinessPanel';
import type { ModbusShareStatus } from '../../../../../types/modbusShare';

interface Step4DatabaseProps {
  deviceSource?: DeviceSourcePort;
  locatorEntry?: StudioLocatorEntryRead;
  locatorOnly?: boolean;
  locatorSuspended?: boolean;
  onLocatorEditorMountedChange?: (mounted: boolean) => void;
  state: WorkbenchV2State;
  dispatch: React.Dispatch<WorkbenchV2Action>;
  workspaceId?: string;
  onCommit?: (confirmedDeviceIds?: string[]) => void;
  /**
   * The explicit workspace-wide activation of legacy callers. This page never starts it: recording starts only from
   * the reviewed device scope, and the overview below reports status.
   */
  activateWorkspace?: () => Promise<StudioV2ActivationResponse>;
  recoverActivationStatus?: () => Promise<StudioV2ActivationRecovery>;
  workspaceReadiness?: StudioV2WorkspaceReadinessSummary | null;
  shareStatus?: ModbusShareStatus | null;
  onNavigateStep?: (step: WorkspaceReadinessStepNumber) => void;
}

export { useStep4Readonly };

/**
 * Step 4 Database 主頁面元件
 * 落地設計決策：「拆檔策略：8 個元件 + 3 個 state module」 與 「Commit log 序列：純函式 + reducer 串聯」
 */
function Step4DatabaseContent({
  deviceSource, state,
  dispatch,
  workspaceId,
  onCommit,
  recoverActivationStatus,
  workspaceReadiness,
  shareStatus,
  onNavigateStep, locatorEntry, locatorSuspended, onLocatorEditorMountedChange,
}: Step4DatabaseProps) {
  const { t: guidedT } = useTranslation('workbench-v2-recording-guided');
  const { t } = useTranslation('workbench-v2');
  // Only the read-back of an earlier activation is shown; no workspace-wide start is wired to this page.
  const activation = useStep4Activation(undefined, recoverActivationStatus, workspaceId);
  const { connector, targets } = state.db;
  const rowGroups = useMemo(() => state.db.row_groups ?? [], [state.db.row_groups]);
  const enabledPoints = useMemo(() => state.points.filter(p => p.enabled), [state.points]);

  const clearRowGroupScope = useCallback(() => {
    if (rowGroups.length > 0) {
      dispatch({ type: 'setDbRowGroups', rowGroups: [] });
    }

    Object.entries(targets).forEach(([pointId, target]) => {
      if (!target.row_group_id) {
        return;
      }
      dispatch({ type: 'updateDbTarget', pointId, patch: { row_group_id: undefined } });
    });
  }, [dispatch, rowGroups, targets]);

  // 2. 切換資料庫種類 Handler
  const handleKindChange = useCallback((kind: DbConnector['kind']) => {
    const kindPatch = getDbKindPatch(kind, connector);
    const defaultConn = getDefaultConnector(kind);
    const patch: Partial<DbConnector> = {
      ...kindPatch,
      name: defaultConn.name,
      table: defaultConn.table,
    };

    if (isRowGroupScopeChange(connector, patch)) {
      clearRowGroupScope();
    }
    // 換了目標就等新目標儲存並查到實際欄位後再建議配對，不以示範欄位重配。
    dispatch({ type: 'updateDbConnector', patch });
  }, [clearRowGroupScope, connector, dispatch]);

  const handleSelectConnectorPool = useCallback((poolConn: SettingsConnector) => {
    // 已存連線以 ID 與版本交由後端解析，不複製遮蔽的憑證。
    const patch = createPoolConnectorPatch(poolConn);
    if (isRowGroupScopeChange(connector, patch)) {
      clearRowGroupScope();
    }
    dispatch({ type: 'updateDbConnector', patch });
  }, [clearRowGroupScope, connector, dispatch]);


  const handleUpdateConnector = useCallback((patch: Partial<DbConnector>) => {
    if (isRowGroupScopeChange(connector, patch)) {
      clearRowGroupScope();
    }
    // 重新輸入密碼即解除「必須重新輸入」狀態；手動改掉身分欄位則重新要求。
    const nextPatch: Partial<DbConnector> = { ...patch };
    if (patch.password !== undefined && patch.password !== '') {
      nextPatch.password_required = false;
      nextPatch.clear_password = false;
    } else if (
      connector.password_required !== true &&
      isDatabaseConnectorIdentityChange(connector, { ...connector, ...patch })
    ) {
      nextPatch.password_required = databaseConnectorNeedsPassword(nextPatch.kind ?? connector.kind);
      nextPatch.password = '';
    }
    dispatch({ type: 'updateDbConnector', patch: nextPatch });
  }, [clearRowGroupScope, connector, dispatch]);

  const saveStates = [
    ...state.devices.map((device) => device.save_state),
    ...state.rules.map((rule) => rule.save_state),
    ...Object.values(state.mappings).map((mapping) => mapping.save_state),
    connector.save_state,
    ...Object.values(targets).map((target) => target.save_state),
  ];
  const groupsQuery = useWriteGroupsQuery(Boolean(workspaceId));
  // Stable across Basic/legacy/circuit/disclosure remounts; one read for the viewed device and canonical revision.
  const [viewedDevice, setViewedDevice] = useState<RecordingViewedDevice>();
  const configurationState = groupsQuery.isError ? 'error' : groupsQuery.isLoading ? 'loading'
    : groupsQuery.data && groupsQuery.data.workspace_id === workspaceId && groupsQuery.data.workspace_revision ? 'ready' : 'unknown';
  const evidenceScope = workspaceId && groupsQuery.data && configurationState === 'ready' &&
    viewedDevice && viewedDevice.workspaceId === workspaceId &&
    state.devices.some(device => device.id === viewedDevice.deviceId && device.persisted && device.save_state === 'saved')
    ? { workspace_id: workspaceId, workspace_revision: groupsQuery.data.workspace_revision, device_id: viewedDevice.deviceId } : null;
  const recordingEvidence = useRuntimeRuleEvidence(evidenceScope);
  const liveGroups = (groupsQuery.data?.groups ?? []).filter((group) => group.status !== 'deleted');
  const groupSummary = {
    state: groupsQuery.isError ? 'error' as const : groupsQuery.isLoading ? 'loading' as const : 'ready' as const,
    total: liveGroups.length,
    applied: liveGroups.filter((group) => group.applied_revision !== '' && group.status !== 'disabled').length,
  };
  const configurationSaved = saveStates.length > 0 && saveStates.every((saveState) => saveState === 'saved');
  const [advancedOpen, setAdvancedOpen] = useState(false);
  React.useLayoutEffect(() => { if (locatorEntry?.selection.state === 'selected') setAdvancedOpen(true); }, [locatorEntry]);
  const [openGroupRequest, setOpenGroupRequest] = useState<{ id: string }>();
  const [reviewRequest, setReviewRequest] = useState<{ workspaceId: string; groupId: string }>();
  const [basicOperationActive, setBasicOperationActive] = useState(false);
  const [savedOperationActive, setSavedOperationActive] = useState(false);
  const savedOperationLock = useRef(false);
  const returnHeading = useRef<HTMLHeadingElement>(null);
  const advancedRef = useRef<HTMLDetailsElement>(null);
  const basicRef = useRef<HTMLDivElement>(null);
  const openAdvanced = useCallback(() => {
    setAdvancedOpen(true);
    advancedRef.current?.focus();
  }, []);
  // A managed profile group's effective action returns to the reviewed start of its own device; it starts nothing.
  const reviewScopedStart = useCallback((group: WriteGroup) => {
    setReviewRequest({ workspaceId: workspaceId ?? '', groupId: group.id });
    basicRef.current?.focus();
  }, [workspaceId]);

  const recordingChange = useRecordingChangeController({
    workspaceId: workspaceId ?? '', readonly: basicOperationActive || savedOperationActive,
    onChanged: (group) => { setOpenGroupRequest({ id: group.id }); void groupsQuery.refetch(); },
  });
  const lifecycleLocked = recordingChange.isOpen || recordingChange.busy || recordingChange.unresolved;
  const actionHost: SavedRecordingActionHost = {
    tryAcquire: () => {
      if (savedOperationLock.current || basicOperationActive || lifecycleLocked) return false;
      savedOperationLock.current = true; return true;
    },
    onBusyChange: (active) => { savedOperationLock.current = active; setSavedOperationActive(active); },
    onReviewScopedStart: reviewScopedStart,
    onAction: (intent) => {
      if (savedOperationLock.current || basicOperationActive || lifecycleLocked || intent.workspaceId !== workspaceId) return;
      if (intent.action === 'edit') { setOpenGroupRequest({ id: intent.groupId }); openAdvanced(); return; }
      if (intent.action === 'recreate' || intent.action === 'transition') {
        savedOperationLock.current = true;
        void recordingChange.requestAction(intent).finally(() => { savedOperationLock.current = false; });
      }
    },
  };
  const savedActionsLocked = savedOperationActive || lifecycleLocked;
  const destinations = state.settings.connectors.filter((item) => item.enabled && item.id && item.identity_revision)
    .map((item) => ({ id: item.id, name: item.name, revision: item.identity_revision! }));

  return (
    <div ref={node => markReturnRegion(node, recordingChange.controller)} className="space-y-6" data-focus-return-region={JSON.stringify([workspaceId, 'step4'])}>
      <h2 data-focus-return-heading ref={returnHeading} tabIndex={-1} className="sr-only focus:not-sr-only">{t('step4.overview_title')}</h2>
      {workspaceId && <RecordingTestWriteRecoveryPanel workspaceId={workspaceId} />}
      <RecordingChangeDialog controller={recordingChange} destinations={destinations} fallbackFocus={() => returnHeading.current} />
      <Step4SupportPanels
        destinationReady={destinationIsSaved(state)}
        connector={connector}
        connectors={state.settings.connectors}
        onUpdateConnector={handleUpdateConnector}
        onKindChange={handleKindChange}
        onSelectConnector={handleSelectConnectorPool}
        disabled={basicOperationActive || savedActionsLocked}
      />
      <div ref={basicRef} tabIndex={-1} data-testid="step4-basic-recording" hidden={!destinationIsSaved(state) && !isModbusShareConfiguredEnabled(shareStatus)}>
      <BasicRecordingPanel
        state={state}
        workspaceId={workspaceId}
        readonly={savedActionsLocked}
        lifecycleHost={actionHost}
        reviewRequest={reviewRequest}
        shareStatus={shareStatus}
        onCommit={onCommit}
        onOpenAdvanced={(groupId) => { if (groupId) setOpenGroupRequest({ id: groupId }); openAdvanced(); }}
        onOpenMemberSettings={(id) => { setOpenGroupRequest({ id }); openAdvanced(); }}
        onOperationActiveChange={setBasicOperationActive}
        onViewedDeviceChange={setViewedDevice}
        renderScopeSummary={(scope) => <RecordingScopeSummary scope={scope} state={state}
          groups={groupsQuery.data?.groups ?? []} evidence={recordingEvidence}
          configurationState={configurationState === 'ready' && (scope.workspaceId !== workspaceId ||
            scope.workspaceRevision !== groupsQuery.data?.workspace_revision) ? 'unknown' : configurationState} />}
      />

      </div>

      <ShareOutputSummary state={state} shareStatus={shareStatus} />

      <details
        ref={advancedRef}
        id="step4-advanced-recording"
        tabIndex={-1}
        open={advancedOpen}
        onToggle={(event) => setAdvancedOpen(event.currentTarget.open)}
        className="space-y-4 rounded-xl border border-slate-800 bg-slate-950/20 p-4"
        data-testid="step4-advanced-recording"
      >
        <summary className="cursor-pointer text-sm font-semibold text-slate-200">{guidedT('management.advancedTitle')}</summary>
        <div className="space-y-6 pt-2">
          <p data-testid="recording-workspace-status-note" className="text-xs text-slate-400">{guidedT('management.workspaceStatus')}</p>
          <WriteGroupSection deviceSource={deviceSource}
            locatorEntry={locatorEntry} locatorSuspended={locatorSuspended} onEditorMountedChange={onLocatorEditorMountedChange}
            state={state}
            workspaceId={workspaceId}
            readonly={basicOperationActive || savedActionsLocked}
            lifecycleHost={actionHost}
            onNavigateStep={onNavigateStep}
            openGroupRequest={openGroupRequest}
            onReviewScopedStart={reviewScopedStart}
          />

          <div id="step4-workspace-readiness" tabIndex={-1} data-testid="step4-workspace-readiness" className="space-y-4">
          <WorkspaceScopeNote groupSummary={groupSummary} />
          <ActivationRecoveryNotice
            recovery={activation.recovery}
            recoveryUnavailable={activation.recoveryUnavailable}
          />
          <CommitSummary
            deviceCount={state.devices.length}
            ruleCount={state.rules.filter(r => r.enabled).length}
            pointCount={enabledPoints.length}
            mappingCount={Object.values(state.mappings).filter(m => m.enabled).length}
            connector={connector}
            groupSummary={groupSummary}
            readinessSummary={workspaceReadiness}
            configurationSaved={configurationSaved}
            onNavigateStep={onNavigateStep}
          />
          </div>
        </div>
      </details>
    </div>
  );
}

/** Step 4 以 React Query 讀取實際欄位；沒有 QueryClient 的嵌入情境使用備援 client。 */
export function Step4Database(props: Step4DatabaseProps) {
  return <SafeQueryBoundary>{props.locatorOnly && props.locatorEntry
    ? <WriteGroupSection deviceSource={props.deviceSource} state={props.state} workspaceId={props.workspaceId} readonly={false}
        locatorEntry={props.locatorEntry} locatorSuspended={props.locatorSuspended} onEditorMountedChange={props.onLocatorEditorMountedChange} locatorOnly />
    : <Step4DatabaseContent {...props} />}</SafeQueryBoundary>;
}
