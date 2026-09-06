import type { DbConnector, DbRowGroup, DbTarget } from '../../state/types';

/** SQLite 走檔案路徑，不需要密碼；其餘類型的連線憑證都綁在身分上。 */
export function databaseConnectorNeedsPassword(kind: DbConnector['kind']): boolean {
  return kind !== 'sqlite';
}

/**
 * 判斷選取的連接器是否換掉了連線身分。身分一改，既有密碼就屬於前一組連線，
 * 不能沿用。
 */
export function isDatabaseConnectorIdentityChange(
  connector: DbConnector,
  next: Pick<DbConnector, 'kind' | 'host' | 'port' | 'database' | 'username'>,
): boolean {
  return (
    next.kind !== connector.kind ||
    (next.host ?? '') !== (connector.host ?? '') ||
    next.port !== connector.port ||
    (next.database ?? '') !== (connector.database ?? '') ||
    (next.username ?? '') !== (connector.username ?? '')
  );
}

export function isRowGroupScopeChange(connector: DbConnector, patch: Partial<DbConnector>): boolean {
  return (
    (patch.kind !== undefined && patch.kind !== connector.kind) ||
    (patch.schema !== undefined && patch.schema !== connector.schema) ||
    (patch.table !== undefined && patch.table !== connector.table)
  );
}

export function rowGroupsForConnector(connector: DbConnector, rowGroups: DbRowGroup[]): DbRowGroup[] {
  return rowGroups.filter((group) => (
    group.table_name === connector.table &&
    (group.table_schema ?? '') === (connector.schema ?? '')
  ));
}

export function clearTargetRowGroupAssignments(targets: Record<string, DbTarget>): Record<string, DbTarget> {
  return Object.fromEntries(
    Object.entries(targets).map(([pointId, target]) => [
      pointId,
      target.row_group_id
        ? { ...target, row_group_id: undefined }
        : target,
    ]),
  );
}

export function syncTargetRowGroupMembership(
  rowGroups: DbRowGroup[],
  pointId: string,
  rowGroupID: string | undefined,
): DbRowGroup[] {
  return rowGroups.map((group) => {
    const memberPointIDs = group.member_point_ids.filter((memberPointId) => memberPointId !== pointId);
    if (group.id !== rowGroupID) {
      return { ...group, member_point_ids: memberPointIDs };
    }
    return { ...group, member_point_ids: [...memberPointIDs, pointId] };
  });
}
