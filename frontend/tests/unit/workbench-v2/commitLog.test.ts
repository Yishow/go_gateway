import { describe, it, expect } from 'vitest';
import { buildCommitLogSequence } from '../../../src/features/datalink/workbench-v2/state/commitLog';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';

/**
 * @file commitLog.test.ts
 * @description 測試 10 步提交 Log 序列產生器，確保輸出順序與資訊正確。
 */

describe('commitLog', () => {
  const mockState: WorkbenchV2State = {
    view: 'flow',
    current: 4,
    completed: new Set([1, 2, 3]),
    sidebarCollapsed: false,
    showSummaryRail: true,
    devices: [
      { id: 'd-1', name: 'PLC-生產線-01', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null },
      { id: 'd-2', name: 'PLC-02', description: '', protocol: 'mc_3e', config: {}, status: 'draft', test: null }
    ],
    rules: [
      { id: 'r-1', device_id: 'd-1', name: 'Holding Registers', start_address: '40001', count: 4, data_type: 'int16', naming_prefix: 't_', enabled: true, scale_multiplier: 1, scale_offset: 0, data_format: '', skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null },
      { id: 'r-2', device_id: 'd-2', name: 'Coils', start_address: '00001', count: 4, data_type: 'bool', naming_prefix: 'c_', enabled: true, scale_multiplier: 1, scale_offset: 0, data_format: '', skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null }
    ],
    selectedRuleId: null,
    points: [
      { id: 'p-1', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 't_1', address: '40001', data_type: 'int16', function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
      { id: 'p-2', device_id: 'd-2', rule_id: 'r-2', rule_name: 'Coils', name: 'c_1', address: '00001', data_type: 'bool', function: 'coil', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 }
    ],
    mappings: {},
    db: {
      connector: {
        kind: 'postgres',
        name: 'PostgreSQL Connector',
        host: 'tsdb.internal',
        port: 5432,
        database: 'gateway_metrics',
        username: 'postgres',
        schema: 'public',
        table: 'sensor_readings',
        write_mode: 'insert',
        write_interval_seconds: 5,
        timestamp_column: 'ts',
        status: 'unknown'
      },
      targets: {
        'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: true },
        'p-2': { tag_id: 'tag.line1.c_1', column_name: 'temp_out_c', enabled: true }
      }
    },
    settings: {} as any,
    committed: false
  };

  it('應產生恰好 10 步 commit log 順序，且資訊與 state 連動', () => {
    const logs = buildCommitLogSequence(mockState);
    expect(logs).toHaveLength(10);

    // 第 1 步
    expect(logs[0].label).toBe('POST /devices × 2');
    expect(logs[0].detail).toContain('PLC-生產線-01 (modbus_tcp)');
    expect(logs[0].detail).toContain('PLC-02 (mc_3e)');

    // 第 2 步
    expect(logs[1].label).toBe('POST /devices/:id/activate × 2');

    // 第 3 步
    expect(logs[2].label).toBe('POST /source-rules × 2');
    expect(logs[2].detail).toContain('Holding Registers');
    expect(logs[2].detail).toContain('Coils');

    // 第 4 步
    expect(logs[3].label).toBe('POST /points × 2');

    // 第 8 步 (資料庫連線)
    expect(logs[7].label).toBe('POST /db-connectors/:id/test');
    expect(logs[7].detail).toBe('postgres tsdb.internal:5432');

    // 第 9 步 (目標寫入)
    expect(logs[8].label).toBe('POST /db-targets × 2');
    expect(logs[8].detail).toBe('→ public.sensor_readings');

    // 所有 log 預設狀態應為 pending
    logs.forEach(log => {
      expect(log.status).toBe('pending');
    });
  });

  it('當為 SQLite 時，資料庫詳細資訊應有不同呈現格式（無 host 且無 schema）', () => {
    const sqliteState: WorkbenchV2State = {
      ...mockState,
      db: {
        ...mockState.db,
        connector: {
          ...mockState.db.connector,
          kind: 'sqlite',
          host: '',
          port: 0,
          database: 'gateway.db',
          schema: '',
          table: 'sensor_readings'
        }
      }
    };

    const logs = buildCommitLogSequence(sqliteState);
    expect(logs[7].detail).toBe('sqlite gateway.db');
    expect(logs[8].detail).toBe('→ sensor_readings');
  });
});
