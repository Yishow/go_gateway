import { describe, expect, it } from 'vitest';
import {
  columnCompatibility,
  findColumnConflicts,
  suggestColumns,
} from '../../../src/features/datalink/workbench-v2/state/writeGroup/columns';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import type { DbColumn } from '../../../src/features/datalink/workbench-v2/state/dbSchemas';

function candidate(key: string, tagKey: string, targetType: GroupCandidate['target_type']): GroupCandidate {
  return {
    key, device_id: 'dev-1', point_id: `pt-${key}`, tag_id: `tag-${key}`, tag_key: tagKey, label: tagKey,
    device_name: 'PLC', address: '40001', target_type: targetType,
  };
}

const column = (name: string, type: string, extra: Partial<DbColumn> = {}): DbColumn =>
  ({ name, type, nullable: true, primary_key: false, ...extra });

describe('RealMetadataAndReviewableAssignment: compatibility', () => {
  it('judges a tag type against the real column type, and says unknown when it cannot tell', () => {
    expect(columnCompatibility('float32', column('t', 'double precision'))).toBe('compatible');
    expect(columnCompatibility('float64', column('t', 'REAL'))).toBe('compatible');
    expect(columnCompatibility('int16', column('t', 'integer'))).toBe('compatible');
    expect(columnCompatibility('uint64', column('t', 'bigint'))).toBe('compatible');
    expect(columnCompatibility('bool', column('t', 'boolean'))).toBe('compatible');
    expect(columnCompatibility('string', column('t', 'varchar(64)'))).toBe('compatible');
    expect(columnCompatibility('float32', column('t', 'text'))).toBe('incompatible');
    expect(columnCompatibility('string', column('t', 'double precision'))).toBe('incompatible');
    expect(columnCompatibility('int32', column('t', 'geometry'))).toBe('unknown');
  });
});

describe('RealMetadataAndReviewableAssignment: suggestions', () => {
  const columns = [
    column('ts', 'timestamptz', { primary_key: true, nullable: false }),
    column('temp_in_c', 'double precision'),
    column('pressure_kpa', 'double precision'),
    column('batch', 'text'),
  ];

  it('keeps a still-valid confirmed assignment and only suggests by name for the rest', () => {
    const items = [candidate('a', 'line.temp_in_c', 'float64'), candidate('b', 'line.pressure_kpa', 'float64'), candidate('c', 'line.batch', 'string')];
    const result = suggestColumns(items, columns, { a: 'temp_in_c' });
    expect(result.assignments.a).toEqual({ column: 'temp_in_c', status: 'confirmed' });
    expect(result.assignments.b).toEqual({ column: 'pressure_kpa', status: 'suggested' });
    expect(result.assignments.c).toEqual({ column: 'batch', status: 'suggested' });
    expect(result.unmatched).toEqual([]);
  });

  it('never wraps around, reuses a column or falls back to an index when there are fewer columns than tags', () => {
    const items = [
      candidate('a', 'x.alpha', 'float64'), candidate('b', 'x.beta', 'float64'),
      candidate('c', 'x.gamma', 'float64'), candidate('d', 'x.delta', 'float64'),
    ];
    const result = suggestColumns(items, [column('v1', 'double precision'), column('v2', 'double precision')], {});
    // Nothing matches by name, so nothing is assigned: guessing by position is not a suggestion.
    expect(Object.keys(result.assignments)).toEqual([]);
    expect(result.unmatched).toEqual(['a', 'b', 'c', 'd']);
  });

  it('never suggests a column twice, a primary key, or an incompatible column', () => {
    const items = [candidate('a', 'x.temp', 'float64'), candidate('b', 'x.temp_in', 'float64')];
    const result = suggestColumns(items, [column('temp', 'double precision'), column('id', 'bigint', { primary_key: true }), column('temp_text', 'text')], {});
    const used = Object.values(result.assignments).map((entry) => entry.column);
    expect(new Set(used).size).toBe(used.length);
    expect(used).not.toContain('id');
    expect(used).not.toContain('temp_text');
  });

  it('marks a confirmed assignment for repair when its column vanished or became incompatible', () => {
    const items = [candidate('a', 'x.temp', 'float64'), candidate('b', 'x.flow', 'float64')];
    const result = suggestColumns(items, [column('flow', 'text')], { a: 'gone', b: 'flow' });
    expect(result.assignments.a).toEqual({ column: 'gone', status: 'repair' });
    expect(result.assignments.b).toEqual({ column: 'flow', status: 'repair' });
  });
});

describe('RealMetadataAndReviewableAssignment: column conflicts', () => {
  const members = (...entries: [string, string, string?][]) => entries.map(([key, col, entity]) => ({ key, column: col, entity_key: entity ?? '' }));

  it('flags two members sharing a column in one row, and allows reuse across distinct entities', () => {
    expect(findColumnConflicts(members(['a', 'v'], ['b', 'v'])).map((c) => c.column)).toEqual(['v']);
    expect(findColumnConflicts(members(['a', 'v', 'line-1'], ['b', 'v', 'line-2']))).toEqual([]);
    expect(findColumnConflicts(members(['a', 'v', 'line-1'], ['b', 'v', 'line-1']))[0].member_keys).toEqual(['a', 'b']);
  });

  it('treats a missing entity key on either side as unproven identity', () => {
    expect(findColumnConflicts(members(['a', 'v', 'line-1'], ['b', 'v'])).length).toBe(1);
  });
});
