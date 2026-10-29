import { describe, expect, it } from 'vitest';
import { isSafeRuntimeSensorText, safeRuntimeAddress, safeRuntimeLabel } from '../../../src/utils/runtimeOperatorSafety';
import { parseRuntimeValueRecord } from '../../../src/utils/safeJson';
const event = { device_id: 'D', point_id: 'P', address: '40001', timestamp: '2026-10-29T00:00:00Z', raw_value: 12, transformed_value: 12, quality: 'good', stale: false };
describe('runtime-only numeric endpoint and paired DSN boundary', () => {
  it.each(['[::1]:502', '[fe80::1%eth0]:502', '[2001:db8::1]:5432', '[2001:DB8:0:0:0:0:0:1]', 'endpoint [fe80::1]:502',
    'Server=localhost;Database=metrics;Integrated Security=true', 'Data Source=localhost;Initial Catalog=metrics',
    'Host=localhost;Port=5432', 'server=plc;database=metrics', 'Driver={SQL Server};Server=localhost', 'host=localhost dbname=metrics', 'host=localhost port=5432', 'Server=localhost;User ID=operator'])('rejects %s as runtime text and streamed scalar', (value) => {
    expect(isSafeRuntimeSensorText(value)).toBe(false);
    expect(parseRuntimeValueRecord({ ...event, raw_value: value })).toBeNull();
    expect(parseRuntimeValueRecord({ ...event, transformed_value: value })).toBeNull();
  });
  it.each(['Server mode', 'Server=ready', 'Database=metrics', 'ratio 1:2', '[a:b:c]', '[12:30:00]', 'm/s', 'OK/NG', 'C: temperature', '2026-10-29T00:00:00Z', '9007199254740993', '0xFFFF', 'BFF', 'SB1A', 'W1A', 'reading=value', 'Pump A (°C)', '[normal]'])('preserves harmless scalar %s without a colon blacklist', (value) => {
    expect(isSafeRuntimeSensorText(value)).toBe(true); expect(parseRuntimeValueRecord({ ...event, raw_value: value })?.raw_value).toBe(value);
  });
  it.each(['40001', 'R100.2', 'BFF', 'SB1A', 'W7FFFFFFF'])('preserves source register %s', (address) => { expect(safeRuntimeAddress(address)).toBe(address); });
  it('preserves validated Unicode destination labels', () => { expect(safeRuntimeLabel('電力資料（A）_Café')).toBe('電力資料（A）_Café'); });
});
