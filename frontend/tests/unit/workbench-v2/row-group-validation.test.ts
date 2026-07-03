import { describe, expect, it } from 'vitest';
import {
  getRowGroupColumnConflicts,
  hasRowGroupColumnConflict,
  hasUnsafeRowGroupUpsert,
} from '../../../src/features/datalink/workbench-v2/state/rowGroupValidation';
import type { DbRowGroup, DbTarget, Mapping, Point } from '../../../src/features/datalink/workbench-v2/state/types';

const points: Point[] = [
  { id: 'p-1', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Rule 1', name: 'T1', address: '40001', data_type: 'float64', function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
  { id: 'p-2', device_id: 'd-1', rule_id: 'r-2', rule_name: 'Rule 2', name: 'T2', address: '40002', data_type: 'float64', function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 },
];
const mappings: Record<string, Mapping> = {
  'p-1': { point_id: 'p-1', tag_key: 'line.t1', display_name: 'T1', unit: 'C', target_type: 'float64', scale: 1, offset: 0, enabled: true },
  'p-2': { point_id: 'p-2', tag_key: 'line.t2', display_name: 'T2', unit: 'C', target_type: 'float64', scale: 1, offset: 0, enabled: true },
};
const oneGroup: DbRowGroup[] = [{
  id: 'group-temp',
  table_schema: 'public',
  table_name: 'sensor_readings',
  member_point_ids: ['p-1', 'p-2'],
  group_key_columns: ['ts'],
}];

describe('row-group database column validation', () => {
  it('allows repeated columns inside one row group with row identity metadata', () => {
    const targets: Record<string, DbTarget> = {
      'p-1': { tag_id: 'tag-1', column_name: 'temperature_c', enabled: true, row_group_id: 'group-temp' },
      'p-2': { tag_id: 'tag-2', column_name: 'temperature_c', enabled: true, row_group_id: 'group-temp' },
    };

    expect(hasRowGroupColumnConflict(points, mappings, targets, oneGroup)).toBe(false);
    expect(getRowGroupColumnConflicts(points, mappings, targets, oneGroup)).toEqual(new Set());
  });

  it('allows repeated columns across different row groups when row identity metadata exists', () => {
    const targets: Record<string, DbTarget> = {
      'p-1': { tag_id: 'tag-1', column_name: 'temperature_c', enabled: true, row_group_id: 'group-a' },
      'p-2': { tag_id: 'tag-2', column_name: 'temperature_c', enabled: true, row_group_id: 'group-b' },
    };
    const rowGroups: DbRowGroup[] = [
      { ...oneGroup[0], id: 'group-a', member_point_ids: ['p-1'] },
      { ...oneGroup[0], id: 'group-b', member_point_ids: ['p-2'] },
    ];

    expect(hasRowGroupColumnConflict(points, mappings, targets, rowGroups)).toBe(false);
    expect(getRowGroupColumnConflicts(points, mappings, targets, rowGroups)).toEqual(new Set());
  });

  it('blocks shared-column row groups under upsert without uniqueness metadata', () => {
    const targets: Record<string, DbTarget> = {
      'p-1': { tag_id: 'tag-1', column_name: 'temperature_c', enabled: true, row_group_id: 'group-temp' },
      'p-2': { tag_id: 'tag-2', column_name: 'temperature_c', enabled: true, row_group_id: 'group-temp' },
    };

    expect(hasUnsafeRowGroupUpsert(points, mappings, targets, oneGroup, 'upsert')).toBe(true);
    expect(hasUnsafeRowGroupUpsert(points, mappings, targets, [{
      ...oneGroup[0],
      unique_key_columns: ['ts', 'line_id'],
    }], 'upsert')).toBe(false);
  });
});
