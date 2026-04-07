import { describe, expect, it } from 'vitest';
import type { Point } from '../../../../types/datalink';
import {
  buildAddressCanvasItems,
  buildPlannedPointAddresses,
  countEligibleSpans,
  formatSourceValue,
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

describe('formatSourceValue', () => {
  it('formats simple numbers in decimal mode', () => {
    expect(formatSourceValue(42, 'decimal')).toBe('42');
    expect(formatSourceValue(3.14159, 'float')).toBe('3.142');
    expect(formatSourceValue(255, 'hex')).toBe('0xFF');
    expect(formatSourceValue(5, 'binary')).toBe('0b101');
  });

  it('handles null, undefined, and empty string', () => {
    expect(formatSourceValue(null, 'decimal')).toBe('—');
    expect(formatSourceValue(undefined, 'decimal')).toBe('—');
    expect(formatSourceValue('', 'decimal')).toBe('—');
  });

  it('formats booleans', () => {
    expect(formatSourceValue(true, 'decimal')).toBe('true');
    expect(formatSourceValue(false, 'decimal')).toBe('false');
  });

  it('formats numeric strings using the selected number mode', () => {
    expect(formatSourceValue('255', 'hex')).toBe('0xFF');
    expect(formatSourceValue('5', 'binary')).toBe('0b101');
    expect(formatSourceValue('3.14159', 'float')).toBe('3.142');
  });

  it('extracts value from MC protocol JSON payload with raw_bytes', () => {
    const mcPayload = '{"value":123.45,"raw_bytes":"AEC3D3"}';
    expect(formatSourceValue(mcPayload, 'decimal')).toBe('123.45');
  });

  it('formats numeric strings nested inside JSON payloads', () => {
    const mcPayload = '{"value":"255","raw_bytes":"00FF"}';
    expect(formatSourceValue(mcPayload, 'hex')).toBe('0xFF');
  });

  it('handles nested JSON object with value key', () => {
    const nested = { value: { nested: true }, extra: 'ignored' };
    expect(formatSourceValue(nested, 'decimal')).toBe('{"nested":true}');
  });

  it('uses first element when value is an array (one logical value per canvas cell)', () => {
    expect(formatSourceValue([1, 2, 3], 'decimal')).toBe('1');
    expect(formatSourceValue([true, false], 'decimal')).toBe('true');
    expect(formatSourceValue([99.9, 123.12], 'decimal')).toBe('99.9');
  });

  it('handles JSON string arrays by taking the first element', () => {
    const arrayJson = '[10, 20, 30]';
    expect(formatSourceValue(arrayJson, 'decimal')).toBe('10');
  });

  it('returns raw string when JSON parsing fails', () => {
    expect(formatSourceValue('not-json', 'decimal')).toBe('not-json');
  });

  it('formats NaN as string', () => {
    expect(formatSourceValue(NaN, 'decimal')).toBe('NaN');
    expect(formatSourceValue(Infinity, 'decimal')).toBe('NaN');
  });
});

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

  it('preserves zero-based D registers for fatek planners', () => {
    expect(
      buildPlannedPointAddresses({
        startAddress: 'D0',
        count: 3,
        dataType: 'int16',
        protocol: 'fatek_fbs',
      }),
    ).toEqual(['D0', 'D1', 'D2']);
  });

  it('expands existing wide points across all occupied cells', () => {
    const items = buildAddressCanvasItems({
      points: [createPoint({ address: '40005', data_type: 'int32' })],
      plannedPointAddresses: ['40009'],
      plannedDataType: 'int16',
      protocol: 'modbus_tcp',
    });

    expect(items.find((item) => item.address === '40005')?.status).toBe('unmanaged');
    expect(items.find((item) => item.address === '40006')?.status).toBe('unmanaged');
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

  it('applies rule scale to live values shown in the canvas', () => {
    const items = buildAddressCanvasItems({
      points: [createPoint({ address: '40001', last_value: 10 })],
      rules: [{
        id: 'rule-1',
        startAddress: '40001',
        count: 1,
        dataType: 'int16',
        namingPrefix: 'SENSOR',
        enabled: true,
        locked: false,
        origin: 'manual',
        skippedAddresses: [],
        scaleMultiplier: 0.1,
        scaleOffset: 5,
      }],
      protocol: 'modbus_tcp',
    });

    expect(items.find((item) => item.address === '40001')?.liveValue).toBe(6);
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

  describe('32-bit mergeSpan in canvas items', () => {
    it('assigns mergeSpan=2 and mergeOffset 0,1 for int32 rule cells', () => {
      const items = buildAddressCanvasItems({
        points: [],
        rules: [
          {
            id: 'rule-1',
            startAddress: '40001',
            count: 1,
            dataType: 'int32',
            namingPrefix: 'T',
            enabled: true,
            locked: false,
            origin: 'manual',
            skippedAddresses: [],
          },
        ],
        protocol: 'modbus_tcp',
      });

      const root = items.find((item) => item.address === '40001');
      const cont = items.find((item) => item.address === '40002');

      expect(root?.mergeSpan).toBe(2);
      expect(root?.mergeOffset).toBe(0);
      expect(cont?.mergeSpan).toBe(2);
      expect(cont?.mergeOffset).toBe(1);
    });

    it('assigns mergeSpan=2 for float32 existing points', () => {
      const items = buildAddressCanvasItems({
        points: [createPoint({ address: '40010', data_type: 'float32' })],
        plannedPointAddresses: [],
        plannedDataType: 'int16',
        protocol: 'modbus_tcp',
      });

      const root = items.find((item) => item.address === '40010');
      const cont = items.find((item) => item.address === '40011');

      expect(root?.mergeSpan).toBe(2);
      expect(root?.mergeOffset).toBe(0);
      expect(root?.status).toBe('unmanaged');
      expect(cont?.mergeSpan).toBe(2);
      expect(cont?.mergeOffset).toBe(1);
      expect(cont?.status).toBe('unmanaged');
    });

    it('assigns mergeSpan=4 for int64 rules', () => {
      const items = buildAddressCanvasItems({
        points: [],
        rules: [
          {
            id: 'rule-1',
            startAddress: '40001',
            count: 1,
            dataType: 'int64',
            namingPrefix: 'T',
            enabled: true,
            locked: false,
            origin: 'manual',
            skippedAddresses: [],
          },
        ],
        protocol: 'modbus_tcp',
      });

      expect(items).toHaveLength(4);
      expect(items[0].mergeSpan).toBe(4);
      expect(items[0].mergeOffset).toBe(0);
      expect(items[3].mergeSpan).toBe(4);
      expect(items[3].mergeOffset).toBe(3);
    });
  });
});
