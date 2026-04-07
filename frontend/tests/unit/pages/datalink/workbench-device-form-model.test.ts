import { describe, expect, it } from 'vitest';

import type { Device } from '../../../../src/types/datalink';
import {
  buildDeviceDraftFromDevice,
  createDefaultDeviceConnectionConfig,
  normalizeMc3eDataFormatValue,
  parseDeviceConnectionConfig,
  sanitizeDeviceConnectionConfig,
} from '../../../../src/pages/datalink/workbench/workbenchDeviceFormModel';

describe('workbenchDeviceFormModel', () => {
  it('fatek_fbs 預設連線含串列 framing（資料位元／停止位元／同位）', () => {
    const cfg = createDefaultDeviceConnectionConfig('fatek_fbs');
    expect(cfg.data_bits).toBe(7);
    expect(cfg.stop_bits).toBe(1);
    expect(cfg.parity).toBe('even');
  });

  it('mc_3e 預設 network_no 與 data_format 與後端／registry 一致', () => {
    const cfg = createDefaultDeviceConnectionConfig('mc_3e');
    expect(cfg.network_no).toBe(0);
    expect(cfg.data_format).toBe('CDAB');
    expect(cfg.port).toBe(5000);
  });

  it('modbus 預設 data_format 與後端／registry 一致', () => {
    expect(createDefaultDeviceConnectionConfig('modbus_tcp').data_format).toBe('ABCD');
    expect(createDefaultDeviceConnectionConfig('modbus_udp').data_format).toBe('ABCD');
    expect(createDefaultDeviceConnectionConfig('modbus_rtu').data_format).toBe('ABCD');
  });

  it('normalizeMc3eDataFormatValue 接受字節序並將舊 binary／ascii 轉為 CDAB', () => {
    expect(normalizeMc3eDataFormatValue('cdab')).toBe('CDAB');
    expect(normalizeMc3eDataFormatValue('binary')).toBe('CDAB');
    expect(normalizeMc3eDataFormatValue('ascii')).toBe('CDAB');
    expect(normalizeMc3eDataFormatValue(undefined)).toBe('CDAB');
  });

  it('parseDeviceConnectionConfig 無效 JSON 回傳空物件', () => {
    expect(parseDeviceConnectionConfig('not json')).toEqual({});
  });

  it('sanitizeDeviceConnectionConfig 會剔除空白字串並 trim 陣列元素', () => {
    const out = sanitizeDeviceConnectionConfig({
      host: '  ',
      topics: [' a ', '', 'b'],
      port: 502,
    });
    expect(out).toEqual({ port: 502, topics: ['a', 'b'] });
  });

  it('buildDeviceDraftFromDevice 載入 mc_3e 時正規化 data_format', () => {
    const device: Device = {
      id: '1',
      name: 'x',
      description: '',
      protocol: 'mc_3e',
      connection_config: JSON.stringify({ data_format: 'binary', host: 'h' }),
      status: 'active',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '',
      updated_at: '',
    };
    const draft = buildDeviceDraftFromDevice(device);
    expect(draft.connectionConfig.data_format).toBe('CDAB');
    expect(draft.connectionConfig.host).toBe('h');
  });

  it('buildDeviceDraftFromDevice 載入 modbus_tcp 缺省 data_format 時回退 ABCD', () => {
    const device: Device = {
      id: '2',
      name: 'modbus',
      description: '',
      protocol: 'modbus_tcp',
      connection_config: JSON.stringify({ host: '10.0.0.8', port: 502 }),
      status: 'active',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '',
      updated_at: '',
    };
    const draft = buildDeviceDraftFromDevice(device);
    expect(draft.connectionConfig.data_format).toBe('ABCD');
    expect(draft.connectionConfig.host).toBe('10.0.0.8');
  });
});
