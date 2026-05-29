import { describe, it, expect } from 'vitest';
import { DEFAULT_SETTINGS, makeDefaultConnector } from '../../../src/features/datalink/workbench-v2/state/settingsDefaults';

/**
 * @file settingsDefaults.test.ts
 * @description 測試 Settings 預設值與預設連接器產生函數。
 */

describe('settingsDefaults', () => {
  it('DEFAULT_SETTINGS 應有正確的結構與預設值，且 Modbus Share 預設啟用', () => {
    expect(DEFAULT_SETTINGS.connectors).toHaveLength(1);
    expect(DEFAULT_SETTINGS.connectors[0].name).toBe('TimeSeries Prod');
    expect(DEFAULT_SETTINGS.connectors[0].status).toBe('ready');

    expect(DEFAULT_SETTINGS.timeseries.write_precision).toBe('millisecond');
    expect(DEFAULT_SETTINGS.scheduler.auto_start).toBe(true);
    expect(DEFAULT_SETTINGS.modbus_share.enabled).toBe(true);
    expect(DEFAULT_SETTINGS.modbus_share.base_register).toBe(40001);
    expect(DEFAULT_SETTINGS.general.theme).toBe('dark');
  });

  it('makeDefaultConnector 應能依 index 產生獨立連線，且 status 為 unknown', () => {
    const conn = makeDefaultConnector(2);
    expect(conn.name).toBe('新連線 2');
    expect(conn.kind).toBe('postgres');
    expect(conn.status).toBe('unknown');
    expect(conn.id).toContain('conn-');
  });
});
