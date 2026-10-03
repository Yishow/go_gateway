import { describe, expect, it } from 'vitest';
import {
  describeAddressConflicts,
  detectAddressConflicts,
  pointConflictKey,
} from '../../../src/features/datalink/workbench-v2/state/sourceRule';
import type { Point } from '../../../src/features/datalink/workbench-v2/state/types';

function point(over: Partial<Point> & Pick<Point, 'id' | 'device_id' | 'rule_id' | 'address'>): Point {
  return {
    rule_name: over.rule_id.toUpperCase(), name: over.id, data_type: 'int16', function: 'holding_register',
    width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0, ...over,
  };
}

describe('DeviceScopedAddressConflict', () => {
  it('two devices may both use 40001', () => {
    const conflicts = detectAddressConflicts([
      point({ id: 'a', device_id: 'dev-A', rule_id: 'r1', address: '40001' }),
      point({ id: 'b', device_id: 'dev-B', rule_id: 'r2', address: '40001' }),
    ]);
    expect(conflicts.size).toBe(0);
  });

  it('the same address on the same device across rules is a conflict, keyed by device', () => {
    const a = point({ id: 'a', device_id: 'dev-A', rule_id: 'r1', address: '40001' });
    const b = point({ id: 'b', device_id: 'dev-A', rule_id: 'r2', address: '40001' });
    const other = point({ id: 'c', device_id: 'dev-B', rule_id: 'r3', address: '40001' });
    const conflicts = detectAddressConflicts([a, b, other]);
    expect(conflicts.has(pointConflictKey(a))).toBe(true);
    expect(conflicts.has(pointConflictKey(other))).toBe(false);
  });

  it('a wide value overlapping a neighbouring address on the same device conflicts, in the same area only', () => {
    const wide = point({ id: 'w', device_id: 'dev-A', rule_id: 'r1', address: '40001', width: 2, data_type: 'int32' });
    const inside = point({ id: 'i', device_id: 'dev-A', rule_id: 'r2', address: '40002' });
    const after = point({ id: 'n', device_id: 'dev-A', rule_id: 'r3', address: '40003' });
    const otherArea = point({ id: 'o', device_id: 'dev-A', rule_id: 'r4', address: '30002' });
    const conflicts = detectAddressConflicts([wide, inside, after, otherArea]);
    expect(conflicts.has(pointConflictKey(wide))).toBe(true);
    expect(conflicts.has(pointConflictKey(inside))).toBe(true);
    expect(conflicts.has(pointConflictKey(after))).toBe(false);
    expect(conflicts.has(pointConflictKey(otherArea))).toBe(false);
  });

  it('works for letter-area protocols and ignores skipped or disabled points', () => {
    const d0 = point({ id: 'a', device_id: 'dev-A', rule_id: 'r1', address: 'D100', width: 2 });
    const d1 = point({ id: 'b', device_id: 'dev-A', rule_id: 'r2', address: 'D101' });
    const skipped = point({ id: 'c', device_id: 'dev-A', rule_id: 'r3', address: 'D100', skipped: true });
    const disabled = point({ id: 'd', device_id: 'dev-A', rule_id: 'r4', address: 'D100', enabled: false });
    const conflicts = detectAddressConflicts([d0, d1, skipped, disabled]);
    expect(conflicts.has(pointConflictKey(d0))).toBe(true);
    expect(conflicts.has(pointConflictKey(d1))).toBe(true);
    // Only the two enabled points count; the skipped and disabled twins of D100 add nothing.
    expect([...conflicts].sort()).toEqual(['dev-A|D100', 'dev-A|D101']);
    const onlyInactive = detectAddressConflicts([d0, skipped, disabled]);
    expect(onlyInactive.size).toBe(0);
  });

  it('falls back to exact text for addresses it cannot normalize', () => {
    const a = point({ id: 'a', device_id: 'dev-A', rule_id: 'r1', address: 'weird' });
    const b = point({ id: 'b', device_id: 'dev-A', rule_id: 'r2', address: 'weird' });
    const c = point({ id: 'c', device_id: 'dev-A', rule_id: 'r3', address: 'other' });
    const conflicts = detectAddressConflicts([a, b, c]);
    expect(conflicts.has(pointConflictKey(a))).toBe(true);
    expect(conflicts.has(pointConflictKey(c))).toBe(false);
  });

  it('describes each conflict by device, address and the rules involved so the operator can locate it', () => {
    const described = describeAddressConflicts([
      point({ id: 'a', device_id: 'dev-A', rule_id: 'r1', rule_name: 'Rule one', address: '40001' }),
      point({ id: 'b', device_id: 'dev-A', rule_id: 'r2', rule_name: 'Rule two', address: '40001' }),
      point({ id: 'c', device_id: 'dev-B', rule_id: 'r3', rule_name: 'Rule three', address: '40001' }),
    ]);
    expect(described).toEqual([{ device_id: 'dev-A', address: '40001', rule_names: ['Rule one', 'Rule two'] }]);
  });

  it('reads MC3E contacts as hex when the device protocol is known', () => {
    const protocols = { 'dev-A': 'mc_3e' as const };
    const x9 = point({ id: 'a', device_id: 'dev-A', rule_id: 'r1', address: 'X9', width: 2 });
    const x10 = point({ id: 'b', device_id: 'dev-A', rule_id: 'r2', address: 'X10' });
    const xa = point({ id: 'c', device_id: 'dev-A', rule_id: 'r3', address: 'XA' });
    // X9 spans hex 9 and A: it overlaps XA but not X10 (= 16).
    const conflicts = detectAddressConflicts([x9, x10, xa], protocols);
    expect(conflicts.has(pointConflictKey(x9))).toBe(true);
    expect(conflicts.has(pointConflictKey(xa))).toBe(true);
    expect(conflicts.has(pointConflictKey(x10))).toBe(false);
  });

  it('stays fast for thousands of points on one device', () => {
    const many = Array.from({ length: 5000 }, (_, i) => point({ id: `p${i}`, device_id: 'dev-A', rule_id: `r${i}`, address: String(40001 + i).padStart(5, '0') }));
    const started = performance.now();
    expect(detectAddressConflicts(many).size).toBe(0);
    expect(describeAddressConflicts(many)).toEqual([]);
    expect(performance.now() - started).toBeLessThan(1000);
  });
});

describe('new rule start address', () => {
  const rule = (start: string, count: number, data_type: 'int16' | 'uint64' | 'float32' = 'int16') => ({
    id: `r-${start}`, device_id: 'dev-A', name: 'r', start_address: start, count, data_type, naming_prefix: 'P_', enabled: true,
    scale_multiplier: 1, scale_offset: 0, data_format: '' as const, skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null,
  });
  it('starts after the last register another rule of the device covers', async () => {
    const { nextFreeStartAddress } = await import('../../../src/features/datalink/workbench-v2/state/sourceRule');
    expect(nextFreeStartAddress([], 'modbus_tcp', '40001')).toBe('40001');
    expect(nextFreeStartAddress([rule('40001', 8)], 'modbus_tcp', '40001')).toBe('40009');
    expect(nextFreeStartAddress([rule('40001', 1), rule('40003', 1, 'uint64')], 'modbus_tcp', '40001')).toBe('40007');
    expect(nextFreeStartAddress([rule('40001', 2, 'float32')], 'modbus_tcp', '40001')).toBe('40005');
    expect(nextFreeStartAddress([rule('garbage', 2)], 'modbus_tcp', '40001')).toBe('40001');
  });
});
