import * as React from 'react';
import { useMemo, useCallback } from 'react';
import { ConnectorSection } from './ConnectorSection';
import { TargetMappingTable } from './TargetMappingTable';
import { CommitSummary } from './CommitSummary';
import { CommitProgress } from './CommitProgress';
import { CommitSuccessCard } from './CommitSuccessCard';
import { autoAssignTargets } from '../../state/autoAssignTargets';
import { getColumnsFor, getDefaultConnector } from '../../state/dbSchemas';
import type { WorkbenchV2State, DbConnector, DbTarget, CommitLog } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import type { StudioV2ActivationResponse } from '../../../../../types/studioV2Activation';
import type { StudioV2WorkspaceReadinessSummary } from '../../../../../types/studioV2WorkspaceReadiness';
import type { WorkspaceReadinessStepNumber } from '../../components/WorkspaceReadinessPanel';

/**
 * Step4Database 元件屬性
 */
interface Step4DatabaseProps {
  state: WorkbenchV2State;
  dispatch: React.Dispatch<WorkbenchV2Action>;
  onCommit?: () => void;
  activateWorkspace?: () => Promise<StudioV2ActivationResponse>;
  workspaceReadiness?: StudioV2WorkspaceReadinessSummary | null;
  onNavigateStep?: (step: WorkspaceReadinessStepNumber) => void;
}

/**
 * 唯讀狀態自訂 Hook
 * 落地設計決策：「Commit 後 form 只讀」
 * @returns 是否唯讀
 */
export function useStep4Readonly(): boolean {
  return false;
}

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
  onNavigateStep,
}: Step4DatabaseProps) {
  const isReadonly = useStep4Readonly();
  const [activationState, setActivationState] = React.useState<{
    phase: 'idle' | 'activating' | 'done';
    response: StudioV2ActivationResponse | null;
  }>({
    phase: 'idle',
    response: null,
  });

  const { connector, targets } = state.db;
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

  // 2. 切換資料庫種類 Handler
  const handleKindChange = useCallback((kind: DbConnector['kind']) => {
    const defaultConn = getDefaultConnector(kind);
    const patch: Partial<DbConnector> = {
      kind,
      name: defaultConn.name,
      database: defaultConn.database,
      username: defaultConn.username,
      password: kind === 'sqlite' ? undefined : connector.password,
      schema: defaultConn.schema,
      table: defaultConn.table
    };

    // 非 SQLite 保留使用者輸入的 Host 與 Port
    if (kind !== 'sqlite') {
      patch.host = connector.host || defaultConn.host;
      patch.port = connector.port || defaultConn.port;
    }

    dispatch({ type: 'updateDbConnector', patch });

    // 同步重算欄位對應
    const nextColumns = getColumnsFor(kind);
    const nextColumnNames = nextColumns.filter(c => !c.primary_key).map(c => c.name);
    const nextTargets = autoAssignTargets(enabledPoints, state.mappings, nextColumnNames, targets);
    dispatch({ type: 'autoAssignDbTargets', targets: nextTargets });
  }, [connector, enabledPoints, targets, state.mappings, dispatch]);

  // 3. 更新 Target 對應 Handler
  const handleUpdateTarget = useCallback((pointId: string, patch: Partial<DbTarget>) => {
    dispatch({ type: 'updateDbTarget', pointId, patch });
  }, [dispatch]);

  const handleStartActivation = useCallback(async () => {
    if (!activateWorkspace) {
      return;
    }

    setActivationState({
      phase: 'activating',
      response: null,
    });

    try {
      const response = await activateWorkspace();
      response.results.forEach((result) => {
        dispatch({
          type: 'updateDevice',
          deviceId: result.device_id,
          patch: result.status === 'success'
            ? { status: 'active', running: true }
            : { running: false },
        });
      });
      setActivationState({
        phase: 'done',
        response,
      });
    } catch (error) {
      setActivationState({
        phase: 'done',
        response: {
          workspace_id: '',
          results: [],
          message: error instanceof Error ? error.message : 'activation failed',
        },
      });
    }
  }, [activateWorkspace, dispatch]);

  // 5. 計算是否有衝突 (有多個已啟用的對應指向同一個 column)
  const hasConflict = useMemo(() => {
    const counts: Record<string, number> = {};
    const visiblePoints = state.points.filter(p => p.enabled && state.mappings[p.id]);

    visiblePoints.forEach(p => {
      const target = targets[p.id];
      if (target && target.enabled) {
        counts[target.column_name] = (counts[target.column_name] || 0) + 1;
      }
    });

    return Object.values(counts).some(count => count > 1);
  }, [state.points, state.mappings, targets]);

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

  const schemaPreviewSignature = useMemo(() => JSON.stringify({
    kind: connector.kind,
    host: connector.host,
    port: connector.port,
    database: connector.database,
    username: connector.username,
    password: connector.password ?? '',
    schema: connector.schema,
    table: connector.table,
    write_mode: connector.write_mode,
    timestamp_column: connector.timestamp_column,
    targets: Object.entries(targets)
      .map(([pointId, target]) => `${pointId}:${target.tag_id}:${target.column_name}:${target.enabled}`)
      .sort(),
  }), [connector.kind, connector.host, connector.port, connector.database, connector.username, connector.password, connector.schema, connector.table, connector.write_mode, connector.timestamp_column, targets]);

  const activationLogs = useMemo<CommitLog[]>(() => {
    const response = activationState.response;
    if (!response) {
      return [];
    }

    return response.results.map((result) => ({
      label: `POST /studio-v2/workspace/activate → ${result.device_id}`,
      detail: result.message,
      status: result.status === 'success' ? 'success' : 'failed',
    }));
  }, [activationState.response]);

  const canContinueToRuntime = Boolean(
    activationState.response?.results.some((result) => result.status === 'success'),
  );

  return (
    <div className="space-y-6">
      {/* 連接器設定卡片 */}
      <ConnectorSection
        connector={connector}
        onUpdateConnector={(patch) => dispatch({ type: 'updateDbConnector', patch })}
        onKindChange={handleKindChange}
        disabled={isReadonly}
      />

      {/* 下方雙欄格局 (表格 7 欄，右側卡片 5 欄) */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div className="lg:col-span-7">
          <TargetMappingTable
            points={state.points}
            mappings={state.mappings}
            targets={targets}
            columns={columns}
            onUpdateTarget={handleUpdateTarget}
            disabled={isReadonly}
          />
        </div>

        {/* 右側三態面板 */}
        <div className="lg:col-span-5">
          {activationState.phase === 'idle' ? (
            <CommitSummary
              deviceCount={state.devices.length}
              ruleCount={state.rules.filter(r => r.enabled).length}
              pointCount={enabledPoints.length}
              mappingCount={Object.values(state.mappings).filter(m => m.enabled).length}
              connector={connector}
              enabledTargetCount={enabledTargetCount}
              hasConflict={hasConflict}
              schemaActionsDisabled={schemaActionsDisabled}
              schemaPreviewSignature={schemaPreviewSignature}
              readinessSummary={workspaceReadiness}
              onActivate={handleStartActivation}
              onNavigateStep={onNavigateStep}
            />
          ) : activationState.phase === 'activating' ? (
            <CommitProgress
              logs={activationLogs}
              status="committing"
            />
          ) : (
            <CommitSuccessCard
              response={activationState.response ?? { workspace_id: '', results: [] }}
              canContinue={canContinueToRuntime}
              onCommit={onCommit || (() => { })}
              onReset={() => {
                setActivationState({
                  phase: 'idle',
                  response: null,
                });
              }}
            />
          )}
        </div>
      </div>
    </div>
  );
}
