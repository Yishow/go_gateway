import { describe, it, expect } from 'vitest';
import { workbenchV2Reducer, INITIAL_STATE, cascadeRemoveDevice } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { Device, WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import type { ConnectionTestResult } from '../../../src/types/datalink';

describe('Workbench V2 Step 1 Reducer Actions', () => {
  const dummyDevice: Device = {
    id: 'dev-02',
    name: '設備-02',
    description: 'MQTT Broker',
    protocol: 'mqtt',
    config: { broker: 'mqtt://localhost:1883' },
    status: 'draft',
    test: null,
  };

  it('addDevice 應將新裝置加入 devices 陣列尾端', () => {
    const state = workbenchV2Reducer(INITIAL_STATE, { type: 'addDevice', device: dummyDevice });
    expect(state.devices).toHaveLength(2);
    expect(state.devices[1]).toEqual(dummyDevice);
  });

  it('removeDevice 應觸發 cascadeRemoveDevice 並移除指定裝置', () => {
    const stateWithTwo = workbenchV2Reducer(INITIAL_STATE, { type: 'addDevice', device: dummyDevice });
    const stateAfterRemove = workbenchV2Reducer(stateWithTwo, { type: 'removeDevice', deviceId: 'dev-02' });
    expect(stateAfterRemove.devices).toHaveLength(1);
    expect(stateAfterRemove.devices[0].id).toBe('dev-01');
  });

  it('updateDevice 應更新指定裝置之基礎欄位', () => {
    const state = workbenchV2Reducer(INITIAL_STATE, {
      type: 'updateDevice',
      deviceId: 'dev-01',
      patch: { description: 'Updated Desc', status: 'tested' },
    });
    expect(state.devices[0].description).toBe('Updated Desc');
    expect(state.devices[0].status).toBe('tested');
  });

  it('updateDeviceConfig 應與原 config 進行 merge 更新', () => {
    const state = workbenchV2Reducer(INITIAL_STATE, {
      type: 'updateDeviceConfig',
      deviceId: 'dev-01',
      patch: { host: '192.168.1.150' },
    });
    expect(state.devices[0].config).toHaveProperty('host', '192.168.1.150');
    expect(state.devices[0].config).toHaveProperty('port', 502); // 應保留舊的 port
  });

  it('renameDevice 應更新指定裝置名稱', () => {
    const state = workbenchV2Reducer(INITIAL_STATE, {
      type: 'renameDevice',
      deviceId: 'dev-01',
      name: 'New Name PLC',
    });
    expect(state.devices[0].name).toBe('New Name PLC');
  });

  it('changeDeviceProtocol 應更換協議、清空測試結果並替換為預設配置', () => {
    // 假設原狀態的設備已 tested
    const testedState: WorkbenchV2State = {
      ...INITIAL_STATE,
      devices: [
        {
          ...INITIAL_STATE.devices[0],
          status: 'tested',
          test: {
            status: 'success',
            stages: {},
          },
        },
      ],
    };

    const state = workbenchV2Reducer(testedState, {
      type: 'changeDeviceProtocol',
      deviceId: 'dev-01',
      protocol: 'mqtt',
    });

    expect(state.devices[0].protocol).toBe('mqtt');
    expect(state.devices[0].status).toBe('draft');
    expect(state.devices[0].test).toBeNull();
    expect(state.devices[0].config).toHaveProperty('broker', 'mqtt://127.0.0.1:1883');
  });


  it('changeDeviceProtocol 應連動更新使用預設起始位址之規則為新協議預設起始位址 (如 MC 3E -> D0)', () => {
    const state = workbenchV2Reducer(INITIAL_STATE, {
      type: 'changeDeviceProtocol',
      deviceId: 'dev-01',
      protocol: 'mc_3e',
    });

    expect(state.devices[0].protocol).toBe('mc_3e');
    expect(state.rules[0].start_address).toBe('D0');
  });

  it('startDeviceTest 應初始化 stages 並將狀態設為 running', () => {
    const state = workbenchV2Reducer(INITIAL_STATE, {
      type: 'startDeviceTest',
      deviceId: 'dev-01',
    });
    expect(state.devices[0].test).not.toBeNull();
    expect(state.devices[0].test?.status).toBe('running');
    expect(state.devices[0].test?.stages).toHaveProperty('connect');
    expect(state.devices[0].test?.stages.connect.status).toBe('running');
    expect(state.devices[0].test?.stages.probe.status).toBe('pending');
  });

  it('resolveDeviceTest 應以 backend diagnostics result 覆寫 connect/probe 階段與整體結果', () => {
    const state1 = workbenchV2Reducer(INITIAL_STATE, { type: 'startDeviceTest', deviceId: 'dev-01' });
    const result: ConnectionTestResult = {
      success: false,
      error: '讀取探測失敗: bad register',
      latency_ms: 18,
      connect: {
        status: 'success',
        message: 'connect ok',
        latency_ms: 12,
      },
      probe: {
        status: 'failed',
        error: 'bad register',
        latency_ms: 6,
      },
      can_activate: false,
      can_collect: false,
    };
    const state2 = workbenchV2Reducer(state1, {
      type: 'resolveDeviceTest',
      deviceId: 'dev-01',
      result,
    });
    expect(state2.devices[0].status).toBe('draft');
    expect(state2.devices[0].test?.status).toBe('failed');
    expect(state2.devices[0].test?.stages.connect.status).toBe('success');
    expect(state2.devices[0].test?.stages.connect.latency_ms).toBe(12);
    expect(state2.devices[0].test?.stages.probe.status).toBe('failed');
    expect(state2.devices[0].test?.stages.probe.message).toBe('bad register');
    expect(state2.devices[0].test?.latency_ms).toBe(18);
  });

  it('completeDeviceTest 應將整個測試設為 success 並記錄總延遲與時間', () => {
    const state1 = workbenchV2Reducer(INITIAL_STATE, { type: 'startDeviceTest', deviceId: 'dev-01' });
    const state2 = workbenchV2Reducer(state1, {
      type: 'completeDeviceTest',
      deviceId: 'dev-01',
      totalLatency: 48,
    });
    expect(state2.devices[0].status).toBe('tested');
    expect(state2.devices[0].test?.status).toBe('success');
    expect(state2.devices[0].test?.latency_ms).toBe(48);
    expect(state2.devices[0].test?.tested_at).toBeDefined();
  });

  it('failDeviceTest 應將特定步驟設為 failed 並把整體測試設為 failed', () => {
    const state1 = workbenchV2Reducer(INITIAL_STATE, { type: 'startDeviceTest', deviceId: 'dev-01' });
    const state2 = workbenchV2Reducer(state1, {
      type: 'failDeviceTest',
      deviceId: 'dev-01',
      stageId: 'connect',
      message: 'Connection timed out',
    });
    expect(state2.devices[0].status).toBe('draft');
    expect(state2.devices[0].test?.status).toBe('failed');
    expect(state2.devices[0].test?.stages.connect.status).toBe('failed');
    expect((state2.devices[0].test?.stages.connect as { message?: string }).message).toBe('Connection timed out');
  });

  it('cascadeRemoveDevice 應連動清除與該設備關聯的規則、點位、Mappings 與 DB Targets', () => {
    // 建立具備高度關聯的 mock 狀態
    const complexState: WorkbenchV2State = {
      ...INITIAL_STATE,
      devices: [
        { ...INITIAL_STATE.devices[0] },
        { ...dummyDevice },
      ],
      rules: [
        { ...INITIAL_STATE.rules[0] }, // dev-01 關聯 rule-01
        {
          id: 'rule-02',
          device_id: 'dev-02',
          name: 'Telemetry Rule',
          start_address: '0',
          count: 5,
          data_type: 'float32',
          naming_prefix: 'M_',
          enabled: true,
          scale_multiplier: 1,
          scale_offset: 0,
          data_format: '',
          skipped_addresses: [],
          share_enabled: false,
          share_start_register: null,
          share_stride: null,
        },
      ],
      points: [
        {
          id: 'p-01',
          device_id: 'dev-01',
          rule_id: 'rule-01',
          rule_name: 'Holding Registers',
          name: 'SENSOR_0',
          address: '40001',
          data_type: 'int16',
          function: 'holding_register',
          width: 1,
          enabled: true,
          skipped: false,
          _rule_scale: 0.1,
          _rule_offset: 0,
        },
        {
          id: 'p-02',
          device_id: 'dev-02',
          rule_id: 'rule-02',
          rule_name: 'Telemetry Rule',
          name: 'M_0',
          address: '0',
          data_type: 'float32',
          function: 'holding_register',
          width: 2,
          enabled: true,
          skipped: false,
          _rule_scale: 1,
          _rule_offset: 0,
        },
      ],
      mappings: {
        'p-01': {
          point_id: 'p-01',
          tag_key: 'line1.temp',
          display_name: 'Temp',
          unit: 'C',
          target_type: 'float64',
          scale: 1,
          offset: 0,
          enabled: true,
        },
        'p-02': {
          point_id: 'p-02',
          tag_key: 'line2.flow',
          display_name: 'Flow',
          unit: 'L/m',
          target_type: 'float64',
          scale: 1,
          offset: 0,
          enabled: true,
        },
      },
      db: {
        ...INITIAL_STATE.db,
        targets: {
          'p-01': { tag_id: 'line1.temp', column_name: 'temp', enabled: true },
          'p-02': { tag_id: 'line2.flow', column_name: 'flow', enabled: true },
        },
      },
    };

    // 刪除 dev-02
    const cleanState = cascadeRemoveDevice(complexState, 'dev-02');

    // 驗證
    expect(cleanState.devices).toHaveLength(1);
    expect(cleanState.devices[0].id).toBe('dev-01');

    // rules
    expect(cleanState.rules).toHaveLength(1);
    expect(cleanState.rules[0].id).toBe('rule-01');

    // points
    expect(cleanState.points).toHaveLength(1);
    expect(cleanState.points[0].id).toBe('p-01');

    // mappings
    expect(cleanState.mappings).toHaveProperty('p-01');
    expect(cleanState.mappings).not.toHaveProperty('p-02');

    // db targets
    expect(cleanState.db.targets).toHaveProperty('p-01');
    expect(cleanState.db.targets).not.toHaveProperty('p-02');
  });
});
