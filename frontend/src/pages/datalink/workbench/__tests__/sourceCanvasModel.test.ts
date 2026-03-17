import { describe, expect, it } from 'vitest';
import type { Point } from '../../../../types/datalink';
import {
  buildAddressCanvasItems,
  buildPlannedPointAddresses,
  countEligibleSpans,
} from '../sourceCanvasModel';
import type { SourceRule } from '../sourceCanvasModel';

function createPoint(overrides: Partial<Point>): Point {
  return {
    id: 'point-1',
    device_id: 'device-1',
    name: 'Existing Point',
    description: '',
    data_type: 'int16',
    address: '40001',
    enabled: true,
    polling_group_id: '',
    last_value: null,
    last_read_at: '',
    last_error: '',
    error_count: 0,
    created_at: '',
    updated_at: '',
    ...overrides,
  };
}

describe('sourceCanvasModel', () => {
  it('builds logical base addresses for wide data types', () => {
    expect(
      buildPlannedPointAddresses({
        startAddress: '40001',
        count: 2,
        dataType: 'float32',
        protocol: 'modbus_tcp',
      }),
    ).toEqual(['40001', '40003']);
  });

  it('expands existing wide points across all occupied cells', () => {
    const items = buildAddressCanvasItems({
      points: [createPoint({ address: '40005', data_type: 'int32' })],
      plannedPointAddresses: ['40009'],
      plannedDataType: 'int16',
      protocol: 'modbus_tcp',
    });

    expect(items.find((item) => item.address === '40005')?.status).toBe('used');
    expect(items.find((item) => item.address === '40006')?.status).toBe('used');
  });

  it('marks overlapping planned cells as conflict', () => {
    const items = buildAddressCanvasItems({
      points: [createPoint({ address: '40002', data_type: 'int16' })],
      plannedPointAddresses: ['40001'],
      plannedDataType: 'float32',
      protocol: 'modbus_tcp',
    });

    expect(items.find((item) => item.address === '40001')?.status).toBe('planned');
    expect(items.find((item) => item.address === '40002')?.status).toBe('conflict');
  });

  describe('countEligibleSpans', () => {
    function createRule(overrides: Partial<SourceRule>): SourceRule {
      return {
        id: 'rule-1',
        startAddress: '40001',
        count: 3,
        dataType: 'int16',
        namingPrefix: 'SENSOR',
        enabled: true,
        locked: false,
        origin: 'manual',
        skippedAddresses: [],
        ...overrides,
      };
    }

    it('returns 0 when no rules exist', () => {
      expect(
        countEligibleSpans({ rules: [], points: [], protocol: 'modbus_tcp' }),
      ).toBe(0);
    });

    it('counts all planned addresses when no points exist', () => {
      const rules = [createRule({ count: 3 })];
      expect(
        countEligibleSpans({ rules, points: [], protocol: 'modbus_tcp' }),
      ).toBe(3);
    });

    it('subtracts addresses already covered by existing points', () => {
      const rules = [createRule({ count: 3 })];
      const points = [createPoint({ address: '40001' })];
      expect(
        countEligibleSpans({ rules, points, protocol: 'modbus_tcp' }),
      ).toBe(2);
    });

    it('excludes skipped addresses from the count', () => {
      const rules = [createRule({ count: 3, skippedAddresses: ['40002'] })];
      expect(
        countEligibleSpans({ rules, points: [], protocol: 'modbus_tcp' }),
      ).toBe(2);
    });

    it('excludes disabled rules', () => {
      const rules = [createRule({ count: 5, enabled: false })];
      expect(
        countEligibleSpans({ rules, points: [], protocol: 'modbus_tcp' }),
      ).toBe(0);
    });

    it('handles wide data types correctly (counts logical spans, not cells)', () => {
      const rules = [createRule({ count: 2, dataType: 'float32' })];
      expect(
        countEligibleSpans({ rules, points: [], protocol: 'modbus_tcp' }),
      ).toBe(2);
    });

    it('excludes only the conflicting root spans, not partial-overlap roots', () => {
      // Rule plans a float32 at 40001 (occupies cells 40001-40002).
      // An existing point sits at 40002, creating a conflict on that cell.
      // The root cell 40001 is still 'planned' (matches Step 2 readyToCreateCount).
      const rules = [createRule({ count: 1, dataType: 'float32', startAddress: '40001' })];
      const points = [createPoint({ address: '40002', data_type: 'int16' })];
      expect(
        countEligibleSpans({ rules, points, protocol: 'modbus_tcp' }),
      ).toBe(1);
    });

    it('excludes spans from overlapping rules (rule-overlap conflict)', () => {
      // Two rules both plan an int16 at 40002 — that address is a conflict.
      // Non-overlapping addresses from each rule still count.
      const ruleA = createRule({ id: 'rule-a', startAddress: '40001', count: 2 });
      const ruleB = createRule({ id: 'rule-b', startAddress: '40002', count: 2 });
      // ruleA plans 40001, 40002; ruleB plans 40002, 40003
      // 40002 is a conflict; 40001 and 40003 are planned → 2 eligible
      expect(
        countEligibleSpans({ rules: [ruleA, ruleB], points: [], protocol: 'modbus_tcp' }),
      ).toBe(2);
    });

    it('excludes root-level point-overlap spans (existing point at planned address)', () => {
      // Rule plans int16 at 40001. Existing point also at 40001.
      // Canvas marks 40001 as 'used' (not 'planned'), so 0 eligible.
      const rules = [createRule({ count: 2 })];
      const points = [createPoint({ address: '40001' })];
      // 40001 → used (existing point covers it), 40002 → planned
      expect(
        countEligibleSpans({ rules, points, protocol: 'modbus_tcp' }),
      ).toBe(1);
    });

    it('excludes locked (protected) rule spans from the count', () => {
      const rules = [createRule({ count: 3, locked: true })];
      expect(
        countEligibleSpans({ rules, points: [], protocol: 'modbus_tcp' }),
      ).toBe(0);
    });

    it('counts unlocked spans while excluding locked spans from the same set', () => {
      const unlocked = createRule({ id: 'rule-u', startAddress: '40001', count: 2, locked: false });
      const locked = createRule({ id: 'rule-l', startAddress: '40010', count: 3, locked: true });
      expect(
        countEligibleSpans({ rules: [unlocked, locked], points: [], protocol: 'modbus_tcp' }),
      ).toBe(2);
    });

    it('correctly handles mixed scenario: points + conflicts + locked + skipped', () => {
      // Rule A: 40001-40003 (3 int16, unlocked)
      // Rule B: 40003 (1 int16, unlocked) — overlaps with A at 40003
      // Rule C: 40010-40011 (2 int16, locked/protected)
      // Existing point at 40001.
      // Skipped: 40002 in rule A.
      //
      // Expected:
      //   40001 → used (existing point) — excluded
      //   40002 → skipped in rule A — excluded
      //   40003 → conflict (rule A + rule B overlap) — excluded
      //   40010, 40011 → locked — excluded
      //   Result: 0 eligible from this set.
      //
      // But rule B also plans 40003 which is the conflict, so no additional spans.
      const ruleA = createRule({
        id: 'rule-a', startAddress: '40001', count: 3, skippedAddresses: ['40002'],
      });
      const ruleB = createRule({ id: 'rule-b', startAddress: '40003', count: 1 });
      const ruleC = createRule({ id: 'rule-c', startAddress: '40010', count: 2, locked: true });
      const points = [createPoint({ address: '40001' })];
      expect(
        countEligibleSpans({ rules: [ruleA, ruleB, ruleC], points, protocol: 'modbus_tcp' }),
      ).toBe(0);
    });
  });
});
