import type { DbRowGroup, DbTarget, Mapping, Point } from './types';

function rowGroupByPoint(rowGroups: DbRowGroup[]): Map<string, DbRowGroup> {
  const out = new Map<string, DbRowGroup>();
  rowGroups.forEach((group) => {
    group.member_point_ids.forEach((pointId) => {
      out.set(pointId, group);
    });
  });
  return out;
}

function targetScope(pointId: string, target: DbTarget, groupsByPoint: Map<string, DbRowGroup>): string {
  if (!target.row_group_id) {
    return `single:${pointId}`;
  }
  const group = groupsByPoint.get(pointId);
  if (!group || group.id !== target.row_group_id || (group.group_key_columns?.length ?? 0) === 0) {
    return `invalid:${target.row_group_id}:${pointId}`;
  }
  return `group:${target.row_group_id}`;
}

function hasValidRowGroupScope(pointId: string, target: DbTarget, groupsByPoint: Map<string, DbRowGroup>): boolean {
  return targetScope(pointId, target, groupsByPoint).startsWith('group:');
}

export function getRowGroupColumnConflicts(
  points: Point[],
  mappings: Record<string, Mapping>,
  targets: Record<string, DbTarget>,
  rowGroups: DbRowGroup[],
): Set<string> {
  const groupsByPoint = rowGroupByPoint(rowGroups);
  const entriesByColumn = new Map<string, Array<{ pointId: string; target: DbTarget }>>();

  points.forEach((point) => {
    const target = targets[point.id];
    if (!point.enabled || !mappings[point.id] || !target?.enabled || !target.column_name) {
      return;
    }
    const entries = entriesByColumn.get(target.column_name) ?? [];
    entries.push({ pointId: point.id, target });
    entriesByColumn.set(target.column_name, entries);
  });

  const conflicts = new Set<string>();
  entriesByColumn.forEach((entries, columnName) => {
    if (entries.length <= 1) {
      return;
    }
    const allGrouped = entries.every(({ pointId, target }) => hasValidRowGroupScope(pointId, target, groupsByPoint));
    if (!allGrouped) {
      conflicts.add(columnName);
    }
  });
  return conflicts;
}

export function hasRowGroupColumnConflict(
  points: Point[],
  mappings: Record<string, Mapping>,
  targets: Record<string, DbTarget>,
  rowGroups: DbRowGroup[],
): boolean {
  return getRowGroupColumnConflicts(points, mappings, targets, rowGroups).size > 0;
}

export function hasUnsafeRowGroupUpsert(
  points: Point[],
  mappings: Record<string, Mapping>,
  targets: Record<string, DbTarget>,
  rowGroups: DbRowGroup[],
  writeMode: 'insert' | 'upsert',
): boolean {
  if (writeMode !== 'upsert') {
    return false;
  }
  const pointByID = new Set(points.filter((point) => point.enabled && mappings[point.id]).map((point) => point.id));
  return rowGroups.some((group) => {
    if ((group.unique_key_columns?.length ?? 0) > 0) {
      return false;
    }
    const counts = new Map<string, number>();
    group.member_point_ids.forEach((pointId) => {
      const target = targets[pointId];
      if (!pointByID.has(pointId) || !target?.enabled || target.row_group_id !== group.id || !target.column_name) {
        return;
      }
      counts.set(target.column_name, (counts.get(target.column_name) ?? 0) + 1);
    });
    return Array.from(counts.values()).some((count) => count > 1);
  });
}
