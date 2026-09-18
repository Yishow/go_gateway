import * as React from 'react';
import { useMemo, useCallback } from 'react';
import { useStep4Activation } from './useStep4Activation';
import {
  buildSchemaPreviewSignature,
  createPoolConnectorPatch,
  databaseConnectorNeedsPassword,
  isDatabaseConnectorIdentityChange,
  isRowGroupScopeChange,
  rowGroupsForConnector,
  syncTargetRowGroupMembership,
  useStep4Readonly,
} from './step4DatabaseHelpers';
import { TargetMappingTable } from './TargetMappingTable';
import { TargetMetadataStatus } from './TargetMetadataStatus';
import { useStep4TargetColumns } from './useStep4TargetColumns';
import { SafeQueryBoundary } from '@/utils/SafeQueryBoundary';
import { CommitSummary } from './CommitSummary';
import { CommitProgress } from './CommitProgress';
import { CommitSuccessCard } from './CommitSuccessCard';
import { DestinationOverviewCard } from './DestinationOverviewCard';
import { SchemaSetupSection } from './SchemaSetupSection';
import { Step4SupportPanels } from './Step4SupportPanels';
import { RowGroupPlanner } from './RowGroupPlanner';
import { ShareOutputSummary } from './ShareOutputSummary';
import { WorkspaceRecordingPlanSetupSection } from './WorkspaceRecordingPlanSetupSection';
import { ActivationNeutralSummary, ActivationRecoveryNotice } from './ActivationNeutralSummary';
import { autoAssignTargets } from '../../state/autoAssignTargets';
import { getDefaultConnector, getDbKindPatch } from '../../state/dbSchemas';
import { hasRowGroupColumnConflict, hasUnsafeRowGroupUpsert } from '../../state/rowGroupValidation';
import type { WorkbenchV2State, DbConnector, DbRowGroup, DbTarget, SettingsConnector } from '../../state/types';
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

  // 只使用已存目標資料表實際查得的欄位；未查到時不以示範欄位配對。
  const targetColumns = useStep4TargetColumns(connector);
  const { columns } = targetColumns;
  const columnNames = useMemo(() => columns.filter(c => !c.primary_key).map(c => c.name), [columns]);
  const columnNamesKey = columnNames.join('\u0000');

  // 1. 實際欄位可用時，替尚未配對的點位建議欄位；既有配對保持不變。
  React.useEffect(() => {
    if (columnNames.length === 0) return;
    const nextTargets = autoAssignTargets(enabledPoints, state.mappings, columnNames, targets);
    dispatch({ type: 'autoAssignDbTargets', targets: nextTargets });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [columnNamesKey]);

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

  // 3. 更新 Target 對應 Handler
  const handleUpdateTarget = useCallback((pointId: string, patch: Partial<DbTarget>) => {
    if (Object.prototype.hasOwnProperty.call(patch, 'row_group_id')) {
      dispatch({
        type: 'setDbRowGroups',
        rowGroups: syncTargetRowGroupMembership(scopedRowGroups, pointId, patch.row_group_id),
      });
    }
    dispatch({ type: 'updateDbTarget', pointId, patch });
  }, [dispatch, scopedRowGroups]);

  const handleSetRowGroups = useCallback((nextRowGroups: DbRowGroup[]) => {
    dispatch({ type: 'setDbRowGroups', rowGroups: nextRowGroups });
  }, [dispatch]);

  const hasConflict = useMemo(() => {
    return hasRowGroupColumnConflict(state.points, state.mappings, targets, scopedRowGroups) ||
      hasUnsafeRowGroupUpsert(state.points, state.mappings, targets, scopedRowGroups, connector.write_mode);
  }, [connector.write_mode, scopedRowGroups, state.points, state.mappings, targets]);

  // 6. 計算啟用中的 db targets 數量
  const enabledTargetCount = useMemo(() => {
    return Object.values(targets).filter(t => t.enabled).length;
  }, [targets]);

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
  const configurationSaved = saveStates.length > 0 && saveStates.every((saveState) => saveState === 'saved');

  return (
    <div className="space-y-6">
      <DestinationOverviewCard
        connector={connector}
        rules={state.rules}
        points={state.points}
        mappings={state.mappings}
        targets={targets}
        hasConflict={hasConflict}
      />

      <WorkspaceRecordingPlanSetupSection
        state={state}
        workspaceId={workspaceId}
        disabled={isReadonly}
      />

      <ShareOutputSummary state={state} shareStatus={shareStatus} />

      <RowGroupPlanner
        connector={connector}
        points={state.points}
        mappings={state.mappings}
        rowGroups={scopedRowGroups}
        onSetRowGroups={handleSetRowGroups}
        disabled={isReadonly}
      />

      <TargetMetadataStatus status={targetColumns.status} onRetry={targetColumns.refetch} />
      <TargetMappingTable
        points={state.points}
        mappings={state.mappings}
        targets={targets}
        rowGroups={scopedRowGroups}
        columns={columns}
        onUpdateTarget={handleUpdateTarget}
        onSetAllEnabled={(enabled) => dispatch({ type: 'setAllDbTargetsEnabled', enabled })}
        disabled={isReadonly}
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
        disabled={isReadonly}
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
            enabledTargetCount={enabledTargetCount}
            hasConflict={hasConflict}
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
  );
}

/** Step 4 以 React Query 讀取實際欄位；沒有 QueryClient 的嵌入情境使用備援 client。 */
export function Step4Database(props: Step4DatabaseProps) {
  return <SafeQueryBoundary><Step4DatabaseContent {...props} /></SafeQueryBoundary>;
}
