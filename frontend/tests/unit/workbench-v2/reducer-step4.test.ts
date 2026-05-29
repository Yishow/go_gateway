import { describe, it, expect } from 'vitest';
import { workbenchV2Reducer, INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { DbConnector, DbTarget, CommitLog } from '../../../src/features/datalink/workbench-v2/state/types';

/**
 * @file reducer-step4.test.ts
 * @description 測試 useWorkbenchV2State reducer 中的 8 個 db 與 commit 相關 action，確保狀態轉換正確且符合不可變性。
 */

describe('reducer-step4', () => {
  it('updateDbConnector: 應更新資料庫連接器設定', () => {
    const patch: Partial<DbConnector> = {
      host: 'mysql.local',
      port: 3306,
      kind: 'mysql'
    };
    const newState = workbenchV2Reducer(INITIAL_STATE, {
      type: 'updateDbConnector',
      patch
    });

    expect(newState.db.connector.host).toBe('mysql.local');
    expect(newState.db.connector.port).toBe(3306);
    expect(newState.db.connector.kind).toBe('mysql');
    // 舊狀態不應被影響
    expect(INITIAL_STATE.db.connector.host).toBe('tsdb.internal');
  });

  it('upsertDbTarget: 應新增點位的資料庫寫入目標', () => {
    const target: DbTarget = {
      tag_id: 'tag.line1.t_1',
      column_name: 'temp_in_c',
      enabled: true
    };
    const newState = workbenchV2Reducer(INITIAL_STATE, {
      type: 'upsertDbTarget',
      pointId: 'p-1',
      target
    });

    expect(newState.db.targets['p-1']).toEqual(target);
  });

  it('updateDbTarget: 應修改已存在的點位資料庫寫入目標，若不存在則忽略', () => {
    // 1. 不存在時
    const stateNoTarget = workbenchV2Reducer(INITIAL_STATE, {
      type: 'updateDbTarget',
      pointId: 'p-1',
      patch: { enabled: false }
    });
    expect(stateNoTarget.db.targets['p-1']).toBeUndefined();

    // 2. 存在時
    const stateWithTarget = workbenchV2Reducer(INITIAL_STATE, {
      type: 'upsertDbTarget',
      pointId: 'p-1',
      target: { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: true }
    });

    const newState = workbenchV2Reducer(stateWithTarget, {
      type: 'updateDbTarget',
      pointId: 'p-1',
      patch: { enabled: false }
    });
    expect(newState.db.targets['p-1'].enabled).toBe(false);
  });

  it('autoAssignDbTargets: 應整批更新/設定 targets', () => {
    const targets: Record<string, DbTarget> = {
      'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: true },
      'p-2': { tag_id: 'tag.line1.t_2', column_name: 'temp_out_c', enabled: true }
    };
    const newState = workbenchV2Reducer(INITIAL_STATE, {
      type: 'autoAssignDbTargets',
      targets
    });

    expect(newState.db.targets).toEqual(targets);
  });

  it('Commit Lifecycle: startCommit -> appendCommitLog -> completeCommit -> resetCommit', () => {
    // 1. startCommit
    let state = workbenchV2Reducer(INITIAL_STATE, { type: 'startCommit' });
    expect(state.commit).toBeDefined();
    expect(state.commit?.status).toBe('committing');
    expect(state.commit?.logs).toEqual([]);
    expect(state.committed).toBe(false);
    expect(state.commit?.started_at).toBeDefined();

    // 2. appendCommitLog (第一筆)
    const log1: CommitLog = { label: 'POST /devices × 1', detail: 'test-device', status: 'pending' };
    state = workbenchV2Reducer(state, { type: 'appendCommitLog', log: log1 });
    expect(state.commit?.logs).toHaveLength(1);
    expect(state.commit?.logs[0].label).toBe('POST /devices × 1');
    expect(state.commit?.logs[0].status).toBe('success'); // 應被轉換為 success 狀態

    // 3. appendCommitLog (防止重複新增同 label 的 log)
    state = workbenchV2Reducer(state, { type: 'appendCommitLog', log: log1 });
    expect(state.commit?.logs).toHaveLength(1);

    // 4. completeCommit
    state = workbenchV2Reducer(state, { type: 'completeCommit' });
    expect(state.committed).toBe(true);
    expect(state.commit?.status).toBe('success');
    expect(state.commit?.finished_at).toBeDefined();
    expect(state.completed.has(4)).toBe(true);

    // 5. resetCommit
    state = workbenchV2Reducer(state, { type: 'resetCommit' });
    expect(state.committed).toBe(false);
    expect(state.commit?.status).toBe('idle');
    expect(state.commit?.logs).toEqual([]);
  });
});
