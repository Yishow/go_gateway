import { describe, expect, it } from 'vitest';
import { newRecordingRangeBuffer, parseRecordingRangeBuffer } from '@/features/datalink/workbench-v2/recording-editor/recordingRangeBuffer';
describe('range buffer explicit strict parsing', () => {
  it('keeps raw invalid numbers and uses data-item count, not register count', () => {
    const raw = { ...newRecordingRangeBuffer(), name: 'Meter', address: '40001', count: '10', dataType: 'float32' as const };
    const parsed = parseRecordingRangeBuffer(raw, 'D', 'modbus_tcp');
    expect(parsed?.rule.count).toBe(10); expect(parsed?.span).toBe(20); expect(parsed?.end).toBe('40020');
    for (const count of ['', '1x', '1.5', '0', '-1', 'Infinity']) {
      const invalid = { ...raw, count }; expect(parseRecordingRangeBuffer(invalid, 'D', 'modbus_tcp')).toBeUndefined(); expect(invalid.count).toBe(count);
    }
  });
  it('uses existing protocol geometry and rejects unsupported protocol / invalid address', () => {
    const raw = { ...newRecordingRangeBuffer(), name: 'Meter', address: 'WFFFFFE', count: '2', dataType: 'float32' as const };
    expect(parseRecordingRangeBuffer(raw, 'D', 'mc_3e')).toBeUndefined();
    expect(parseRecordingRangeBuffer({ ...raw, address: 'D0', count: '2' }, 'D', 'mc_3e')?.span).toBe(4);
    expect(parseRecordingRangeBuffer({ ...raw, address: '40001' }, 'D', 'mqtt')).toBeUndefined();
  });
  it('normalizes only the submitted backend prefix/address and preserves raw user input', () => {
    const raw = { ...newRecordingRangeBuffer(), name: 'Sensor room', address: ' d100 ', count: '1' };
    const parsed = parseRecordingRangeBuffer(raw, 'D', 'mc_3e');
    expect(parsed?.rule.naming_prefix).toBe('SENSOR-ROOM');
    expect(parsed?.rule.start_address).toBe('D100');
    expect(parsed?.rule.name).toBe('Sensor room');
    expect(raw.address).toBe(' d100 ');
  });

});
