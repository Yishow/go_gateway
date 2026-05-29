import * as React from 'react';
import { useMemo, useCallback } from 'react';
import { ConnectorSection } from './ConnectorSection';
import { TargetMappingTable } from './TargetMappingTable';
import { CommitSummary } from './CommitSummary';
import { CommitProgress } from './CommitProgress';
import { CommitSuccessCard } from './CommitSuccessCard';
import { autoAssignTargets } from '../../state/autoAssignTargets';
import { buildCommitLogSequence } from '../../state/commitLog';
import { getColumnsFor, getDefaultConnector } from '../../state/dbSchemas';
import type { WorkbenchV2State, DbConnector, DbTarget } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';

/**
 * Step4Database 元件屬性
 */
interface Step4DatabaseProps {
  state: WorkbenchV2State;
  dispatch: React.Dispatch<WorkbenchV2Action>;
  onCommit?: () => void;
}

/**
 * 唯讀狀態自訂 Hook
 * 落地設計決策：「Commit 後 form 只讀」
 * @param state 當前狀態
 * @returns 是否唯讀
 */
export function useStep4Readonly(state: WorkbenchV2State): boolean {
  return state.committed;
}

/**
 * Step 4 Database 主頁面元件
 * 落地設計決策：「拆檔策略：8 個元件 + 3 個 state module」 與 「Commit log 序列：純函式 + reducer 串聯」
 */
export function Step4Database({ state, dispatch, onCommit }: Step4DatabaseProps) {
  const isReadonly = useStep4Readonly(state);

  const { connector, targets } = state.db;
  const enabledPoints = useMemo(() => state.points.filter(p => p.enabled), [state.points]);

  // 取得目前資料庫種類的欄位定義
  const columns = useMemo(() => getColumnsFor(connector.kind), [connector.kind]);
  const columnNames = useMemo(() => columns.filter(c => !c.primary_key).map(c => c.name), [columns]);

  // 1. 初始化分配 (mount 時執行一次，確保 Step 1-3 變更同步)
  React.useEffect(() => {
    if (state.committed) return;
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

  // 4. 提交排程器
  const handleStartCommit = useCallback(() => {
    dispatch({ type: 'startCommit' });
  }, [dispatch]);

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

  // 7. Commit 動畫計時器 Effect
  const commitSeq = useMemo(() => buildCommitLogSequence(state), [state]);
  const commitState = state.commit;

  React.useEffect(() => {
    if (!commitState || commitState.status !== 'committing') return;

    const currentLogIndex = commitState.logs.length;

    if (currentLogIndex >= commitSeq.length) {
      dispatch({ type: 'completeCommit' });
      return;
    }

    const timer = setTimeout(() => {
      dispatch({
        type: 'appendCommitLog',
        log: commitSeq[currentLogIndex]
      });
    }, 280);

    return () => clearTimeout(timer);
  }, [commitState, commitSeq, dispatch]);

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
          {!commitState || commitState.status === 'idle' ? (
            <CommitSummary
              deviceCount={state.devices.length}
              ruleCount={state.rules.filter(r => r.enabled).length}
              pointCount={enabledPoints.length}
              mappingCount={Object.values(state.mappings).filter(m => m.enabled).length}
              connector={connector}
              enabledTargetCount={enabledTargetCount}
              hasConflict={hasConflict}
              onCommit={handleStartCommit}
            />
          ) : commitState.status === 'committing' ? (
            <CommitProgress
              logs={commitState.logs}
              status={commitState.status}
            />
          ) : (
            <CommitSuccessCard
              writeIntervalSeconds={connector.write_interval_seconds}
              onCommit={onCommit || (() => {})}
            />
          )}
        </div>
      </div>
    </div>
  );
}
