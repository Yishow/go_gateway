import { describe, expect, it } from 'vitest';
import { resolveRuntimeDashboardDevice } from '../../../src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';

function createState(overrides?: Partial<WorkbenchV2State>): WorkbenchV2State {
  return {
    view: 'flow',
    current: 4,
    completed: new Set([1, 2, 3]),
    sidebarCollapsed: false,
    showSummaryRail: true,
    devices: [
      { id: 'device-A', name: 'Mixer PLC', description: '', protocol: 'modbus_tcp', config: {}, status: 'active', test: null },
      { id: 'device-B', name: 'Filler PLC', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null },
    ],
    rules: [
      { id: 'rule-A', device_id: 'device-A', name: 'Holding Registers', start_address: '40001', count: 1, data_type: 'int16', naming_prefix: 'a_', enabled: true, scale_multiplier: 1, scale_offset: 0, data_format: '', skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null },
      { id: 'rule-B', device_id: 'device-B', name: 'Input Registers', start_address: '30001', count: 1, data_type: 'int16', naming_prefix: 'b_', enabled: true, scale_multiplier: 1, scale_offset: 0, data_format: '', skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null },
    ],
    selectedRuleId: null,
    points: [],
    mappings: {},
    db: {
      connector: {
        kind: 'postgres',
        name: 'Postgres',
        host: 'localhost',
        port: 5432,
        database: 'gateway',
        username: 'postgres',
        schema: 'public',
        table: 'readings',
        write_mode: 'insert',
        write_interval_seconds: 5,
        timestamp_column: 'ts',
        status: 'unknown',
      },
      targets: {},
    },
    settings: {
      connectors: [],
      timeseries: { write_precision: 'millisecond', partition_interval: 'daily', batch_size: 100, retention_days: 30 },
      scheduler: { default_interval_ms: 1000, default_retry_count: 3, default_retry_delay_ms: 500, breaker_threshold: 5, auto_start: true },
      modbus_share: { enabled: false, bind_address: '127.0.0.1', port: 1502, slave_id: 1, base_register: 0 },
      general: { theme: 'dark', locale: 'zh-TW', addr_format: 'modbus', api_base: '/api', api_version: 'v1', timeout_seconds: 10, log_level: 'info', sse_heartbeat_seconds: 15, enable_debug_panel: false, enable_audit_log: false },
    },
    committed: true,
    ...overrides,
  };
}

describe('resolveRuntimeDashboardDevice', () => {
  it('prefers the selected rule device when selectedRuleId is present', () => {
    const state = createState({ selectedRuleId: 'rule-B' });

    expect(resolveRuntimeDashboardDevice(state)).toBe('device-B');
  });

  it('uses the single enabled-rule device when all enabled rules belong to one device', () => {
    const state = createState({
      rules: [
        { id: 'rule-A', device_id: 'device-A', name: 'Holding Registers', start_address: '40001', count: 1, data_type: 'int16', naming_prefix: 'a_', enabled: true, scale_multiplier: 1, scale_offset: 0, data_format: '', skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null },
        { id: 'rule-A2', device_id: 'device-A', name: 'Input Registers', start_address: '30001', count: 1, data_type: 'int16', naming_prefix: 'a2_', enabled: true, scale_multiplier: 1, scale_offset: 0, data_format: '', skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null },
      ],
    });

    expect(resolveRuntimeDashboardDevice(state)).toBe('device-A');
  });

  it('falls back to the single device when only one device exists', () => {
    const state = createState({
      devices: [
        { id: 'device-A', name: 'Mixer PLC', description: '', protocol: 'modbus_tcp', config: {}, status: 'active', test: null },
      ],
      rules: [],
    });

    expect(resolveRuntimeDashboardDevice(state)).toBe('device-A');
  });

  it('returns null when no deterministic device can be resolved', () => {
    const state = createState();

    expect(resolveRuntimeDashboardDevice(state)).toBeNull();
  });
});
