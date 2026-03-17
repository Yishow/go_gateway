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
  });
});
