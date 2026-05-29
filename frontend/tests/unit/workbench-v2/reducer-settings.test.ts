import { describe, it, expect } from 'vitest';
import { workbenchV2Reducer, INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { DEFAULT_SETTINGS } from '../../../src/features/datalink/workbench-v2/state/settingsDefaults';

/**
 * @file reducer-settings.test.ts
 * @description 測試 useWorkbenchV2State reducer 中的 8 個系統設定與連接器池 Action，確保資料轉換正確與關鍵變更時重置 status。
 */

describe('reducer-settings', () => {
  it('updateSettings: 應更新系統設定', () => {
    const newState = workbenchV2Reducer(INITIAL_STATE, {
      type: 'updateSettings',
      patch: {
        general: {
          ...INITIAL_STATE.settings.general,
          theme: 'light'
        }
      }
    });

    expect(newState.settings.general.theme).toBe('light');
    // 舊狀態不變
    expect(INITIAL_STATE.settings.general.theme).toBe('dark');
  });

  it('updateSettingsSection: 應只更新特定區塊設定', () => {
    const newState = workbenchV2Reducer(INITIAL_STATE, {
      type: 'updateSettingsSection',
      section: 'timeseries',
      patch: { write_precision: 'second' }
    });

    expect(newState.settings.timeseries.write_precision).toBe('second');
    expect(newState.settings.timeseries.batch_size).toBe(500); // 其它欄位應保留
  });

  it('addConnector: 應新增一筆預設 postgres 連接器，且 status 為 unknown', () => {
    const newState = workbenchV2Reducer(INITIAL_STATE, { type: 'addConnector' });
    expect(newState.settings.connectors).toHaveLength(2);
    expect(newState.settings.connectors[1].name).toBe('新連線 2');
    expect(newState.settings.connectors[1].status).toBe('unknown');
  });

  it('updateConnector: 應合併屬性變更，且當修改 kind/host/port 時，status 自動回歸 unknown', () => {
    // 1. 正常修改 name，status 不應重置
    let state = workbenchV2Reducer(INITIAL_STATE, {
      type: 'updateConnector',
      id: 'conn-prod',
      patch: { name: 'Prod DB' }
    });
    expect(state.settings.connectors[0].name).toBe('Prod DB');
    expect(state.settings.connectors[0].status).toBe('ready'); // 仍是 ready

    // 2. 修改 host 參數，status 應重置為 unknown
    state = workbenchV2Reducer(state, {
      type: 'updateConnector',
      id: 'conn-prod',
      patch: { host: 'new-db.internal' }
    });
    expect(state.settings.connectors[0].host).toBe('new-db.internal');
    expect(state.settings.connectors[0].status).toBe('unknown');
  });

  it('removeConnector: 應成功移除指定連接器', () => {
    const state = workbenchV2Reducer(INITIAL_STATE, {
      type: 'removeConnector',
      id: 'conn-prod'
    });
    expect(state.settings.connectors).toHaveLength(0);
  });

  it('Connector Test: startConnectorTest -> completeConnectorTest', () => {
    // 1. startConnectorTest
    let state = workbenchV2Reducer(INITIAL_STATE, {
      type: 'startConnectorTest',
      id: 'conn-prod'
    });
    expect(state.settings.connectors[0].status).toBe('testing');

    // 2. completeConnectorTest (成功)
    state = workbenchV2Reducer(state, {
      type: 'completeConnectorTest',
      id: 'conn-prod',
      result: {
        status: 'ready',
        last_check_at: '2026-05-29T12:00:00Z'
      }
    });
    expect(state.settings.connectors[0].status).toBe('ready');
    expect(state.settings.connectors[0].last_check_at).toBe('2026-05-29T12:00:00Z');
  });

  it('resetSettingsToDefaults: 應重置系統設定為預設', () => {
    // 先做修改
    let state = workbenchV2Reducer(INITIAL_STATE, {
      type: 'updateSettingsSection',
      section: 'scheduler',
      patch: { default_interval_ms: 5000 }
    });
    expect(state.settings.scheduler.default_interval_ms).toBe(5000);

    // 重置
    state = workbenchV2Reducer(state, { type: 'resetSettingsToDefaults' });
    expect(state.settings.scheduler.default_interval_ms).toBe(DEFAULT_SETTINGS.scheduler.default_interval_ms);
  });
});
