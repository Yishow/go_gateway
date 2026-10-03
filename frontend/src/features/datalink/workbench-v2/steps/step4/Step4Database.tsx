import * as React from 'react';
import { useMemo, useCallback, useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useStep4Activation } from './useStep4Activation';
import {
  buildSchemaPreviewSignature,
  createPoolConnectorPatch,
  databaseConnectorNeedsPassword,
  isDatabaseConnectorIdentityChange,
  isRowGroupScopeChange,
  rowGroupsForConnector,
  useStep4Readonly,
} from './step4DatabaseHelpers';
import { WriteGroupSection } from './writeGroup/WriteGroupSection';
import { BasicRecordingPanel } from './BasicRecordingPanel';
import { useWriteGroupsQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import { SafeQueryBoundary } from '@/utils/SafeQueryBoundary';
import { CommitSummary } from './CommitSummary';
import { CommitProgress } from './CommitProgress';
import { CommitSuccessCard } from './CommitSuccessCard';
import { SchemaSetupSection } from './SchemaSetupSection';
import { Step4SupportPanels } from './Step4SupportPanels';
import { ShareOutputSummary } from './ShareOutputSummary';
import { ActivationNeutralSummary, ActivationRecoveryNotice } from './ActivationNeutralSummary';
import { getDefaultConnector, getDbKindPatch } from '../../state/dbSchemas';
import type { WorkbenchV2State, DbConnector, SettingsConnector } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import type {
  StudioV2ActivationRecovery,
  StudioV2ActivationResponse,
} from '../../../../../types/studioV2Activation';
import type { StudioV2WorkspaceReadinessSummary } from '../../../../../types/studioV2WorkspaceReadiness';
import type { WorkspaceReadinessStepNumber } from '../../components/WorkspaceReadinessPanel';
import type { ModbusShareStatus } from '../../../../../types/modbusShare';

interface Step4DatabaseProps {
  state: WorkbenchV2State;
  dispatch: React.Dispatch<WorkbenchV2Action>;
  workspaceId?: string;
  onCommit?: (confirmedDeviceIds?: string[]) => void;
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
  state,
  dispatch,
  workspaceId,
  onCommit,
  activateWorkspace,
  recoverActivationStatus,
  workspaceReadiness,
  shareStatus,
  onNavigateStep,
}: Step4DatabaseProps) {
  const { t } = useTranslation('workbench-v2');
  const activation = useStep4Activation(
    activateWorkspace,
    recoverActivationStatus,
    workspaceId,
  );

  const isReadonly = useStep4Readonly(activation.phase);
  const { connector, targets } = state.db;
  const rowGroups = useMemo(() => state.db.row_groups ?? [], [state.db.row_groups]);
  const scopedRowGroups = useMemo(() => rowGroupsForConnector(connector, rowGroups), [connector, rowGroups]);
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

  const schemaActionsDisabled = useMemo(() => {
    if (connector.save_state !== 'saved') {
      return true;
    }

    return Object.values(targets).some((target) => target.save_state !== 'saved');
  }, [connector.save_state, targets]);

  const schemaPreviewSignature = useMemo(() => buildSchemaPreviewSignature(connector, targets, {
    workspaceId, devices: state.devices, points: state.points, mappings: state.mappings, rowGroups: scopedRowGroups,
  }), [connector, targets, workspaceId, state.devices, state.points, state.mappings, scopedRowGroups]);

  const hasActiveDevice = state.devices.some((d) => d.status === 'active' || d.running) ||
    Boolean(activation.recovery?.devices.some((device) => device.running));
  const canContinueToRuntime = activation.canContinue || hasActiveDevice;
  const hasActivationSuccess = activation.canContinue;
  const saveStates = [
    ...state.devices.map((device) => device.save_state),
    ...state.rules.map((rule) => rule.save_state),
    ...Object.values(state.mappings).map((mapping) => mapping.save_state),
    connector.save_state,
    ...Object.values(targets).map((target) => target.save_state),
  ];
  const groupsQuery = useWriteGroupsQuery(Boolean(workspaceId));
  const liveGroups = (groupsQuery.data?.groups ?? []).filter((group) => group.status !== 'deleted');
  const groupSummary = {
    state: groupsQuery.isError ? 'error' as const : groupsQuery.isLoading ? 'loading' as const : 'ready' as const,
    total: liveGroups.length,
    applied: liveGroups.filter((group) => group.applied_revision !== '' && group.status !== 'disabled').length,
  };
  const configurationSaved = saveStates.length > 0 && saveStates.every((saveState) => saveState === 'saved');
  const [advancedOpen, setAdvancedOpen] = useState(false);
  const [basicOperationActive, setBasicOperationActive] = useState(false);
  const advancedRef = useRef<HTMLDetailsElement>(null);
  const openAdvanced = useCallback(() => {
    setAdvancedOpen(true);
    advancedRef.current?.focus();
  }, []);
  useEffect(() => {
    if (activation.phase !== 'idle') setAdvancedOpen(true);
  }, [activation.phase]);

  return (
    <div className="space-y-6">
      <BasicRecordingPanel
        state={state}
        workspaceId={workspaceId}
        readonly={isReadonly}
        shareStatus={shareStatus}
        onCommit={onCommit}
        onOpenAdvanced={openAdvanced}
        onOperationActiveChange={setBasicOperationActive}
      />

      <Step4SupportPanels
        connector={connector}
        connectors={state.settings.connectors}
        onUpdateConnector={handleUpdateConnector}
        onKindChange={handleKindChange}
        onSelectConnector={handleSelectConnectorPool}
        tableSetup={
          <SchemaSetupSection
            connector={connector}
            workspaceId={workspaceId}
            schemaActionsDisabled={schemaActionsDisabled || isReadonly}
            schemaPreviewSignature={schemaPreviewSignature}
          />
        }
        disabled={isReadonly || basicOperationActive}
      />

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
        <summary className="cursor-pointer text-sm font-semibold text-slate-200">{t('step4.basic.advanced_title')}</summary>
        <div className="space-y-6 pt-2">
          <WriteGroupSection
            state={state}
            workspaceId={workspaceId}
            readonly={isReadonly || basicOperationActive}
            onNavigateStep={onNavigateStep}
          />

          {activation.phase === 'idle' ? (
            <>
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
                onActivate={activation.start}
                onNavigateStep={onNavigateStep}
              />
            </>
          ) : activation.phase === 'activating' ? (
            <CommitProgress logs={activation.logs} status="committing" />
          ) : (
            <>
              {activation.logs.some((log) => log.status === 'failed') && (
                <CommitProgress logs={activation.logs} status="failed" onRetry={activation.start} />
              )}
              {hasActivationSuccess ? (
                <CommitSuccessCard
                  response={activation.response ?? { workspace_id: '', results: [] }}
                  canContinue={canContinueToRuntime}
                  onCommit={onCommit}
                  onReset={activation.reset}
                />
              ) : (
                <ActivationNeutralSummary
                  canContinue={canContinueToRuntime}
                  onCommit={onCommit}
                  response={activation.response}
                  recovery={activation.recovery}
                  recoveryUnavailable={activation.recoveryUnavailable}
                  onReset={activation.reset}
                />
              )}
            </>
          )}
        </div>
      </details>
    </div>
  );
}

/** Step 4 以 React Query 讀取實際欄位；沒有 QueryClient 的嵌入情境使用備援 client。 */
export function Step4Database(props: Step4DatabaseProps) {
  return <SafeQueryBoundary><Step4DatabaseContent {...props} /></SafeQueryBoundary>;
}
