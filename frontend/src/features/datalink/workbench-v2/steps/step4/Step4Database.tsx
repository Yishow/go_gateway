import * as React from 'react';
import { useMemo, useCallback } from 'react';
import { useStep4Activation } from './useStep4Activation';
import {
  buildSchemaPreviewSignature,
  clearTargetRowGroupAssignments,
  createPoolConnectorPatch,
  databaseConnectorNeedsPassword,
  isDatabaseConnectorIdentityChange,
  isRowGroupScopeChange,
  rowGroupsForConnector,
  syncTargetRowGroupMembership,
  useStep4Readonly,
} from './step4DatabaseHelpers';
import { TargetMappingTable } from './TargetMappingTable';
import { CommitSummary } from './CommitSummary';
import { CommitProgress } from './CommitProgress';
import { CommitSuccessCard } from './CommitSuccessCard';
import { DestinationOverviewCard } from './DestinationOverviewCard';
import { SchemaSetupSection } from './SchemaSetupSection';
import { Step4SupportPanels } from './Step4SupportPanels';
import { RowGroupPlanner } from './RowGroupPlanner';
import { ShareOutputSummary } from './ShareOutputSummary';
import { RecordingPlanSetupSection } from './RecordingPlanSetupSection';
import { ActivationNeutralSummary } from './ActivationNeutralSummary';
import { autoAssignTargets } from '../../state/autoAssignTargets';
import { getColumnsFor, getDefaultConnector, getDbKindPatch } from '../../state/dbSchemas';
import { hasRowGroupColumnConflict, hasUnsafeRowGroupUpsert } from '../../state/rowGroupValidation';
import type { WorkbenchV2State, DbConnector, DbRowGroup, DbTarget, SettingsConnector } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import type { StudioV2ActivationResponse } from '../../../../../types/studioV2Activation';
import type { StudioV2WorkspaceReadinessSummary } from '../../../../../types/studioV2WorkspaceReadiness';
import type { WorkspaceReadinessStepNumber } from '../../components/WorkspaceReadinessPanel';
import type { ModbusShareStatus } from '../../../../../types/modbusShare';

interface Step4DatabaseProps {
  state: WorkbenchV2State;
  dispatch: React.Dispatch<WorkbenchV2Action>;
  onCommit?: () => void;
  activateWorkspace?: () => Promise<StudioV2ActivationResponse>;
  workspaceReadiness?: StudioV2WorkspaceReadinessSummary | null;
  shareStatus?: ModbusShareStatus | null;
  onNavigateStep?: (step: WorkspaceReadinessStepNumber) => void;
}

export { useStep4Readonly };

/**
 * Step 4 Database 主頁面元件
 * 落地設計決策：「拆檔策略：8 個元件 + 3 個 state module」 與 「Commit log 序列：純函式 + reducer 串聯」
 */
export function Step4Database({
  state,
  dispatch,
  onCommit,
  activateWorkspace,
  workspaceReadiness,
  shareStatus,
  onNavigateStep,
}: Step4DatabaseProps) {
  const activation = useStep4Activation(activateWorkspace, dispatch);

  const isReadonly = useStep4Readonly(activation.phase);
  const { connector, targets } = state.db;
  const rowGroups = useMemo(() => state.db.row_groups ?? [], [state.db.row_groups]);
  const scopedRowGroups = useMemo(() => rowGroupsForConnector(connector, rowGroups), [connector, rowGroups]);
  const enabledPoints = useMemo(() => state.points.filter(p => p.enabled), [state.points]);

  // 取得目前資料庫種類的欄位定義
  const columns = useMemo(() => getColumnsFor(connector.kind), [connector.kind]);
  const columnNames = useMemo(() => columns.filter(c => !c.primary_key).map(c => c.name), [columns]);

  // 1. 初始化分配 (mount 時執行一次，確保 Step 1-3 變更同步)
  React.useEffect(() => {
    const initialTargets = autoAssignTargets(enabledPoints, state.mappings, columnNames, targets);
    dispatch({ type: 'autoAssignDbTargets', targets: initialTargets });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

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

    const scopeChanged = isRowGroupScopeChange(connector, patch);
    const baseTargets = scopeChanged ? clearTargetRowGroupAssignments(targets) : targets;

    if (scopeChanged) {
      clearRowGroupScope();
    }

    dispatch({ type: 'updateDbConnector', patch });

    // 同步重算欄位對應
    const nextColumns = getColumnsFor(kind);
    const nextColumnNames = nextColumns.filter(c => !c.primary_key).map(c => c.name);
    const nextTargets = autoAssignTargets(enabledPoints, state.mappings, nextColumnNames, baseTargets);
    dispatch({ type: 'autoAssignDbTargets', targets: nextTargets });
  }, [clearRowGroupScope, connector, enabledPoints, targets, state.mappings, dispatch]);

  const handleSelectConnectorPool = useCallback((poolConn: SettingsConnector) => {
    // 連接器清單的密碼一律被後端遮蔽為空字串，直接沿用會讓新連線帶著舊憑證。
    const patch = createPoolConnectorPatch(connector, poolConn);
    const scopeChanged = isRowGroupScopeChange(connector, patch);
    const baseTargets = scopeChanged ? clearTargetRowGroupAssignments(targets) : targets;
    if (scopeChanged) {
      clearRowGroupScope();
    }
    dispatch({ type: 'updateDbConnector', patch });
    const nextColumns = getColumnsFor(poolConn.kind);
    const nextColumnNames = nextColumns.filter(c => !c.primary_key).map(c => c.name);
    const nextTargets = autoAssignTargets(enabledPoints, state.mappings, nextColumnNames, baseTargets);
    dispatch({ type: 'autoAssignDbTargets', targets: nextTargets });
  }, [clearRowGroupScope, connector, enabledPoints, state.mappings, targets, dispatch]);


  const handleUpdateConnector = useCallback((patch: Partial<DbConnector>) => {
    if (isRowGroupScopeChange(connector, patch)) {
      clearRowGroupScope();
    }
    // 重新輸入密碼即解除「必須重新輸入」狀態；手動改掉身分欄位則重新要求。
    const nextPatch: Partial<DbConnector> = { ...patch };
    if (patch.password !== undefined && patch.password.trim() !== '') {
      nextPatch.password_required = false;
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

  const schemaPreviewSignature = useMemo(
    () => buildSchemaPreviewSignature(connector, targets),
    [connector, targets],
  );

  const hasActiveDevice = state.devices.some((d) => d.status === 'active' || d.running);
  const canContinueToRuntime = activation.canContinue || hasActiveDevice;
  const hasActivationSuccess = activation.canContinue || (
    activation.phase === 'done' &&
    !activation.logs.some((log) => log.status === 'failed') &&
    hasActiveDevice
  );

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

      <RecordingPlanSetupSection
        deviceId={state.devices[0]?.id}
        measurementIds={state.points.map((p) => p.id)}
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
            schemaActionsDisabled={schemaActionsDisabled}
            schemaPreviewSignature={schemaPreviewSignature}
          />
        }
        disabled={isReadonly}
      />


      {activation.phase === 'idle' ? (
        <CommitSummary
          deviceCount={state.devices.length}
          ruleCount={state.rules.filter(r => r.enabled).length}
          pointCount={enabledPoints.length}
          mappingCount={Object.values(state.mappings).filter(m => m.enabled).length}
          connector={connector}
          enabledTargetCount={enabledTargetCount}
          hasConflict={hasConflict}
          readinessSummary={workspaceReadiness}
          onActivate={activation.start}
          onNavigateStep={onNavigateStep}
        />
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
              onCommit={onCommit || (() => {})}
              onReset={activation.reset}
            />
          ) : (
            <ActivationNeutralSummary
              canContinue={canContinueToRuntime}
              onCommit={onCommit}
              onReset={activation.reset}
            />
          )}
        </>
      )}
    </div>
  );
}
